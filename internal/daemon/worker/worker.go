package worker

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/boundary/globals"
	"github.com/hashicorp/boundary/internal/cmd/base"
	"github.com/hashicorp/boundary/internal/cmd/config"
	"github.com/hashicorp/boundary/internal/daemon/worker/internal/metric"
	"github.com/hashicorp/boundary/internal/daemon/worker/session"
	"github.com/hashicorp/boundary/internal/errors"
	pb "github.com/hashicorp/boundary/internal/gen/controller/servers"
	pbs "github.com/hashicorp/boundary/internal/gen/controller/servers/services"
	"github.com/hashicorp/boundary/internal/observability/event"
	"github.com/hashicorp/boundary/internal/server"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/base62"
	"github.com/hashicorp/go-secure-stdlib/mlock"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/nodeenrollment"
	nodeenet "github.com/hashicorp/nodeenrollment/net"
	nodeefile "github.com/hashicorp/nodeenrollment/storage/file"
	"github.com/hashicorp/nodeenrollment/types"
	"github.com/mr-tron/base58"
	ua "go.uber.org/atomic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver/manual"
	"google.golang.org/protobuf/proto"
)

type randFn func(length int) (string, error)

// downstreamRouter defines a min interface which must be met by a
// Worker.downstreamRoutes field
type downstreamRouter interface {
	// StartRouteMgmtTicking starts a ticker which manages the router's
	// connections.
	StartRouteMgmtTicking(context.Context, func() string, int) error
}

// downstreamers provides at least a minimum interface that must be met by a
// Worker.downstreamWorkers field which is far better than allowing any (empty
// interface)
type downstreamers interface {
	// Root returns the root of the downstreamers' graph
	Root() string
}

// downstreamRouterFactory provides a simple factory which a Worker can use to
// create its downstreamRouter
var downstreamRouterFactory func() downstreamRouter

const (
	authenticationStatusNeverAuthenticated uint32 = iota
	authenticationStatusFirstAuthentication
	authenticationStatusFirstStatusRpcSuccessful
)

type Worker struct {
	conf   *Config
	logger hclog.Logger

	baseContext context.Context
	baseCancel  context.CancelFunc
	started     *ua.Bool

	tickerWg sync.WaitGroup

	// grpc.ClientConns are thread safe.
	// See https://github.com/grpc/grpc-go/blob/master/Documentation/concurrency.md#clients
	// This is exported for tests.
	GrpcClientConn *grpc.ClientConn

	// receives address updates and contains the grpc resolver.
	addressReceivers []addressReceiver

	sessionManager session.Manager

	controllerStatusConn *atomic.Value
	everAuthenticated    *ua.Uint32
	lastStatusSuccess    *atomic.Value
	workerStartTime      time.Time
	operationalState     *atomic.Value

	controllerMultihopConn *atomic.Value

	proxyListener *base.ServerListener

	// Used to generate a random nonce for Controller connections
	nonceFn randFn

	// We store the current set in an atomic value so that we can add
	// reload-on-sighup behavior later
	tags *atomic.Value
	// This stores whether or not to send updated tags on the next status
	// request. It can be set via startup in New below, or (eventually) via
	// SIGHUP.
	updateTags *ua.Bool

	// The storage for node enrollment
	WorkerAuthStorage             *nodeefile.Storage
	WorkerAuthCurrentKeyId        *ua.String
	WorkerAuthRegistrationRequest string
	workerAuthSplitListener       *nodeenet.SplitListener

	// downstream workers and routes to those workers
	downstreamWorkers downstreamers
	downstreamRoutes  downstreamRouter

	// Test-specific options (and possibly hidden dev-mode flags)
	TestOverrideX509VerifyDnsName  string
	TestOverrideX509VerifyCertPool *x509.CertPool
	TestOverrideAuthRotationPeriod time.Duration

	statusLock sync.Mutex
}

func New(conf *Config) (*Worker, error) {
	const op = "worker.New"
	metric.InitializeHttpCollectors(conf.PrometheusRegisterer)
	metric.InitializeWebsocketCollectors(conf.PrometheusRegisterer)
	metric.InitializeClusterClientCollectors(conf.PrometheusRegisterer)

	w := &Worker{
		conf:                   conf,
		logger:                 conf.Logger.Named("worker"),
		started:                ua.NewBool(false),
		controllerStatusConn:   new(atomic.Value),
		everAuthenticated:      ua.NewUint32(authenticationStatusNeverAuthenticated),
		lastStatusSuccess:      new(atomic.Value),
		controllerMultihopConn: new(atomic.Value),
		tags:                   new(atomic.Value),
		updateTags:             ua.NewBool(false),
		nonceFn:                base62.Random,
		WorkerAuthCurrentKeyId: new(ua.String),
		operationalState:       new(atomic.Value),
	}

	if downstreamRouterFactory != nil {
		w.downstreamRoutes = downstreamRouterFactory()
	}

	w.lastStatusSuccess.Store((*LastStatusInformation)(nil))
	scheme := strconv.FormatInt(time.Now().UnixNano(), 36)
	controllerResolver := manual.NewBuilderWithScheme(scheme)
	w.addressReceivers = []addressReceiver{&grpcResolverReceiver{controllerResolver}}

	if conf.RawConfig.Worker == nil {
		conf.RawConfig.Worker = new(config.Worker)
	}

	w.parseAndStoreTags(conf.RawConfig.Worker.Tags)

	if conf.SecureRandomReader == nil {
		conf.SecureRandomReader = rand.Reader
	}

	if !conf.RawConfig.DisableMlock {
		// Ensure our memory usage is locked into physical RAM
		if err := mlock.LockMemory(); err != nil {
			return nil, fmt.Errorf(
				"Failed to lock memory: %v\n\n"+
					"This usually means that the mlock syscall is not available.\n"+
					"Boundary uses mlock to prevent memory from being swapped to\n"+
					"disk. This requires root privileges as well as a machine\n"+
					"that supports mlock. Please enable mlock on your system or\n"+
					"disable Boundary from using it. To disable Boundary from using it,\n"+
					"set the `disable_mlock` configuration option in your configuration\n"+
					"file.",
				err)
		}
	}

	var listenerCount int
	for i := range conf.Listeners {
		l := conf.Listeners[i]
		if l == nil || l.Config == nil || l.Config.Purpose == nil {
			continue
		}
		if len(l.Config.Purpose) != 1 {
			return nil, fmt.Errorf("found listener with multiple purposes %q", strings.Join(l.Config.Purpose, ","))
		}
		switch l.Config.Purpose[0] {
		case "proxy":
			if w.proxyListener == nil {
				w.proxyListener = l
			}
			listenerCount++
		}
	}
	if listenerCount != 1 {
		return nil, fmt.Errorf("exactly one proxy listener is required")
	}

	return w, nil
}

// Reload will update a worker with a new Config. The worker will only use
// relevant parts of the new config, specifically:
// - Worker Tags
// - Initial Upstream addresses
func (w *Worker) Reload(ctx context.Context, newConf *config.Config) {
	const op = "worker.(Worker).Reload"

	w.parseAndStoreTags(newConf.Worker.Tags)

	if !strutil.EquivalentSlices(newConf.Worker.InitialUpstreams, w.conf.RawConfig.Worker.InitialUpstreams) {
		w.statusLock.Lock()
		defer w.statusLock.Unlock()

		upstreamsMessage := fmt.Sprintf(
			"Initial Upstreams has changed; old upstreams were: %s, new upstreams are: %s",
			w.conf.RawConfig.Worker.InitialUpstreams,
			newConf.Worker.InitialUpstreams,
		)
		event.WriteSysEvent(ctx, op, upstreamsMessage)
		w.conf.RawConfig.Worker.InitialUpstreams = newConf.Worker.InitialUpstreams

		for _, ar := range w.addressReceivers {
			ar.SetAddresses(w.conf.RawConfig.Worker.InitialUpstreams)
			// set InitialAddresses in case the worker has not successfully dialed yet
			ar.InitialAddresses(w.conf.RawConfig.Worker.InitialUpstreams)
		}
	}
}

func (w *Worker) Start() error {
	const op = "worker.(Worker).Start"

	w.baseContext, w.baseCancel = context.WithCancel(context.Background())

	if w.started.Load() {
		event.WriteSysEvent(w.baseContext, op, "already started, skipping")
		return nil
	}

	if w.conf.WorkerAuthKms == nil || w.conf.DevUsePkiForUpstream {
		// In this section, we look for existing worker credentials. The two
		// variables below store whether to create new credentials and whether
		// to create a fetch request so it can be displayed in the worker
		// startup info. These may be different because if initial creds have
		// been generated on the worker side but not yet authorized/fetched from
		// the controller, we don't want to invalidate that request on restart
		// by generating a new set of credentials. However it's safe to output a
		// new fetch request so we do in fact do that.
		//
		// Note that if a controller-generated activation token has been
		// supplied, we do not output a fetch request; we attempt to use that
		// directly later.
		var err error
		w.WorkerAuthStorage, err = nodeefile.New(w.baseContext,
			nodeefile.WithBaseDirectory(w.conf.RawConfig.Worker.AuthStoragePath))
		if err != nil {
			return errors.Wrap(w.baseContext, err, op, errors.WithMsg("error loading worker auth storage directory"))
		}

		var createNodeAuthCreds bool
		var createFetchRequest bool
		nodeCreds, err := types.LoadNodeCredentials(w.baseContext, w.WorkerAuthStorage, nodeenrollment.CurrentId, nodeenrollment.WithWrapper(w.conf.WorkerAuthStorageKms))
		switch {
		case err == nil:
			// It's unclear why this would ever happen -- it shouldn't -- so
			// this is simply safety against panics if something goes
			// catastrophically wrong
			if nodeCreds == nil {
				event.WriteSysEvent(w.baseContext, op, "no error loading worker auth creds but nil creds, creating new creds for registration")
				createNodeAuthCreds = true
				createFetchRequest = true
				break
			}

			// Check that we have valid creds, or that we have generated creds but
			// simply are still waiting on authentication (in which case we don't
			// want to invalidate what we've already sent)
			var validCreds bool
			switch len(nodeCreds.CertificateBundles) {
			case 0:
				// Still waiting on initial creds, so don't invalidate the request
				// by creating new credentials. However, we will generate and
				// display a new valid request in case the first was lost.
				createFetchRequest = true

			default:
				now := time.Now()
				for _, bundle := range nodeCreds.CertificateBundles {
					if bundle.CertificateNotBefore.AsTime().Before(now) && bundle.CertificateNotAfter.AsTime().After(now) {
						// If we have a certificate in its validity period,
						// everything is fine
						validCreds = true
						break
					}
				}

				// Certificates are both expired, so create new credentials and
				// output a request based on those
				createNodeAuthCreds = !validCreds
				createFetchRequest = !validCreds
			}

		case errors.Is(err, nodeenrollment.ErrNotFound):
			// Nothing was found on disk, so create
			createNodeAuthCreds = true
			createFetchRequest = true

		default:
			// Some other type of error happened, bail out
			return fmt.Errorf("error loading worker auth creds: %w", err)
		}

		// Don't output a fetch request if an activation token has been
		// provided. Technically we _could_ still output a fetch request, and it
		// would be valid to do so, but if a token was provided it may well be
		// confusing to a user if it seems like it was ignored because a fetch
		// request was still output.
		if actToken := w.conf.RawConfig.Worker.ControllerGeneratedActivationToken; actToken != "" {
			createFetchRequest = false
		}

		// NOTE: this block _must_ be before the `if createFetchRequest` block
		// or the fetch request may have no credentials to work with
		if createNodeAuthCreds {
			nodeCreds, err = types.NewNodeCredentials(
				w.baseContext,
				w.WorkerAuthStorage,
				nodeenrollment.WithRandomReader(w.conf.SecureRandomReader),
				nodeenrollment.WithWrapper(w.conf.WorkerAuthStorageKms),
			)
			if err != nil {
				return fmt.Errorf("error generating new worker auth creds: %w", err)
			}
		}

		if createFetchRequest {
			if nodeCreds == nil {
				return fmt.Errorf("need to create fetch request but worker auth creds are nil: %w", err)
			}
			req, err := nodeCreds.CreateFetchNodeCredentialsRequest(w.baseContext, nodeenrollment.WithRandomReader(w.conf.SecureRandomReader))
			if err != nil {
				return fmt.Errorf("error creating worker auth fetch credentials request: %w", err)
			}
			reqBytes, err := proto.Marshal(req)
			if err != nil {
				return fmt.Errorf("error marshaling worker auth fetch credentials request: %w", err)
			}
			w.WorkerAuthRegistrationRequest = base58.FastBase58Encoding(reqBytes)
			if err != nil {
				return fmt.Errorf("error encoding worker auth registration request: %w", err)
			}
			currentKeyId, err := nodeenrollment.KeyIdFromPkix(nodeCreds.CertificatePublicKeyPkix)
			if err != nil {
				return fmt.Errorf("error deriving worker auth key id: %w", err)
			}
			w.WorkerAuthCurrentKeyId.Store(currentKeyId)
		}
		// Regardless, we want to load the currentKeyId
		currentKeyId, err := nodeenrollment.KeyIdFromPkix(nodeCreds.CertificatePublicKeyPkix)
		if err != nil {
			return fmt.Errorf("error deriving worker auth key id: %w", err)
		}
		w.WorkerAuthCurrentKeyId.Store(currentKeyId)
	}

	if err := w.StartControllerConnections(); err != nil {
		return errors.Wrap(w.baseContext, err, op, errors.WithMsg("error making controller connections"))
	}

	var err error
	w.sessionManager, err = session.NewManager(pbs.NewSessionServiceClient(w.GrpcClientConn))
	if err != nil {
		return errors.Wrap(w.baseContext, err, op, errors.WithMsg("error creating session manager"))
	}

	if err := w.startListeners(w.sessionManager); err != nil {
		return errors.Wrap(w.baseContext, err, op, errors.WithMsg("error starting worker listeners"))
	}

	w.operationalState.Store(server.ActiveOperationalState)

	// Rather than deal with some of the potential error conditions for Add on
	// the waitgroup vs. Done (in case a function exits immediately), we will
	// always start rotation and simply exit early if we're using KMS
	w.tickerWg.Add(3)
	go func() {
		defer w.tickerWg.Done()
		w.startStatusTicking(w.baseContext, w.sessionManager, &w.addressReceivers)
	}()
	go func() {
		defer w.tickerWg.Done()
		w.startAuthRotationTicking(w.baseContext)
	}()
	go func() {
		defer w.tickerWg.Done()
		if w.downstreamRoutes != nil {
			err := w.downstreamRoutes.StartRouteMgmtTicking(
				w.baseContext,
				func() string {
					if s := w.LastStatusSuccess(); s != nil {
						return s.WorkerId
					}
					return "unknown worker id"
				},
				-1, // indicates the ticker should run until cancelled.
			)
			if err != nil {
				errors.Wrap(w.baseContext, err, op)
			}
		}
	}()

	w.workerStartTime = time.Now()
	w.started.Store(true)

	return nil
}

func (w *Worker) hasActiveConnection() bool {
	activeConnection := false
	w.sessionManager.ForEachLocalSession(
		func(s session.Session) bool {
			conns := s.GetLocalConnections()
			for _, v := range conns {
				if v.Status == pbs.CONNECTIONSTATUS_CONNECTIONSTATUS_CONNECTED {
					activeConnection = true
					return false
				}
			}
			return true
		})
	return activeConnection
}

// Graceful shutdown sets the worker state to "shutdown" and will wait to return until there
// are no longer any active connections.
func (w *Worker) GracefulShutdown() error {
	const op = "worker.(Worker).GracefulShutdown"
	event.WriteSysEvent(w.baseContext, op, "worker entering graceful shutdown")
	w.operationalState.Store(server.ShutdownOperationalState)

	// Wait for connections to drain
	for w.hasActiveConnection() {
		time.Sleep(time.Millisecond * 250)
	}
	event.WriteSysEvent(w.baseContext, op, "worker connections have drained")

	return nil
}

// Shutdown shuts down the workers. skipListeners can be used to not stop
// listeners, useful for tests if we want to stop and start a worker. In order
// to create new listeners we'd have to migrate listener setup logic here --
// doable, but work for later.
func (w *Worker) Shutdown() error {
	const op = "worker.(Worker).Shutdown"
	if !w.started.Load() {
		event.WriteSysEvent(w.baseContext, op, "already shut down, skipping")
		return nil
	}
	event.WriteSysEvent(w.baseContext, op, "worker shutting down")

	// Set state to shutdown
	w.operationalState.Store(server.ShutdownOperationalState)

	// Stop listeners first to prevent new connections to the
	// controller.
	defer w.started.Store(false)
	if err := w.stopServersAndListeners(); err != nil {
		return fmt.Errorf("error stopping worker servers and listeners: %w", err)
	}

	// Shut down all connections.
	w.cleanupConnections(w.baseContext, true, w.sessionManager)

	// Wait for next status request to succeed. Don't wait too long;
	// wrap the base context in a timeout equal to our status grace
	// period.
	waitStatusStart := time.Now()
	nextStatusCtx, nextStatusCancel := context.WithTimeout(w.baseContext, w.conf.StatusGracePeriodDuration)
	defer nextStatusCancel()
	for {
		if err := nextStatusCtx.Err(); err != nil {
			event.WriteError(w.baseContext, op, err, event.WithInfoMsg("error waiting for next status report to controller"))
			break
		}

		if w.lastSuccessfulStatusTime().Sub(waitStatusStart) > 0 {
			break
		}

		time.Sleep(time.Second)
	}

	// Proceed with remainder of shutdown.
	w.baseCancel()
	for _, ar := range w.addressReceivers {
		ar.SetAddresses(nil)
	}

	w.started.Store(false)
	w.tickerWg.Wait()
	if w.conf.Eventer != nil {
		if err := w.conf.Eventer.FlushNodes(context.Background()); err != nil {
			return fmt.Errorf("error flushing worker eventer nodes: %w", err)
		}
	}

	event.WriteSysEvent(w.baseContext, op, "worker finished shutting down")
	return nil
}

func (w *Worker) parseAndStoreTags(incoming map[string][]string) {
	if len(incoming) == 0 {
		w.tags.Store([]*pb.TagPair{})
		return
	}
	tags := []*pb.TagPair{}
	for k, vals := range incoming {
		for _, v := range vals {
			tags = append(tags, &pb.TagPair{
				Key:   k,
				Value: v,
			})
		}
	}
	w.tags.Store(tags)
	w.updateTags.Store(true)
}

func (w *Worker) getSessionTls(sessionManager session.Manager) func(hello *tls.ClientHelloInfo) (*tls.Config, error) {
	const op = "worker.(Worker).getSessionTls"
	return func(hello *tls.ClientHelloInfo) (*tls.Config, error) {
		ctx := w.baseContext
		var sessionId string
		switch {
		case strings.HasPrefix(hello.ServerName, globals.SessionPrefix):
			sessionId = hello.ServerName
		default:
			for _, proto := range hello.SupportedProtos {
				if strings.HasPrefix(proto, globals.SessionPrefix) {
					sessionId = proto
					break
				}
			}
		}

		if sessionId == "" {
			event.WriteSysEvent(ctx, op, "session_id not found in either SNI or ALPN protos", "server_name", hello.ServerName)
			return nil, fmt.Errorf("could not find session ID in SNI or ALPN protos")
		}

		lastSuccess := w.LastStatusSuccess()
		if lastSuccess == nil {
			event.WriteSysEvent(ctx, op, "no last status information found at session acceptance time")
			return nil, fmt.Errorf("no last status information found at session acceptance time")
		}

		timeoutContext, cancel := context.WithTimeout(w.baseContext, session.ValidateSessionTimeout)
		defer cancel()
		sess, err := sessionManager.LoadLocalSession(timeoutContext, sessionId, lastSuccess.GetWorkerId())
		if err != nil {
			return nil, fmt.Errorf("error refreshing session: %w", err)
		}

		certPool := x509.NewCertPool()
		certPool.AddCert(sess.GetCertificate())

		tlsConf := &tls.Config{
			Certificates: []tls.Certificate{
				{
					Certificate: [][]byte{sess.GetCertificate().Raw},
					PrivateKey:  ed25519.PrivateKey(sess.GetPrivateKey()),
					Leaf:        sess.GetCertificate(),
				},
			},
			NextProtos: []string{"http/1.1"},
			MinVersion: tls.VersionTLS13,

			// These two are set this way so we can make use of VerifyConnection,
			// which we set on this TLS config below. We are not skipping
			// verification!
			ClientAuth:         tls.RequireAnyClientCert,
			InsecureSkipVerify: true,
		}

		// We disable normal DNS SAN behavior as we don't rely on DNS or IP
		// addresses for security and want to avoid issues with including localhost
		// etc.
		verifyOpts := x509.VerifyOptions{
			DNSName: sessionId,
			Roots:   certPool,
			KeyUsages: []x509.ExtKeyUsage{
				x509.ExtKeyUsageClientAuth,
				x509.ExtKeyUsageServerAuth,
			},
		}
		if w.TestOverrideX509VerifyCertPool != nil {
			verifyOpts.Roots = w.TestOverrideX509VerifyCertPool
		}
		if w.TestOverrideX509VerifyDnsName != "" {
			verifyOpts.DNSName = w.TestOverrideX509VerifyDnsName
		}
		tlsConf.VerifyConnection = func(cs tls.ConnectionState) error {
			// Go will not run this without at least one peer certificate, but
			// doesn't hurt to check
			if len(cs.PeerCertificates) == 0 {
				return errors.New(ctx, errors.InvalidParameter, op, "no peer certificates provided")
			}
			_, err := cs.PeerCertificates[0].Verify(verifyOpts)
			return err
		}
		return tlsConf, nil
	}
}
