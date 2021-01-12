package schema

import (
	"context"
	"database/sql"
	"testing"

	"github.com/hashicorp/boundary/internal/db/schema/postgres"
	"github.com/hashicorp/boundary/internal/docker"
	"github.com/hashicorp/boundary/internal/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	c, u, _, err := docker.StartDbInDocker("postgres")
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, c())
	})
	d, err := sql.Open("postgres", u)
	require.NoError(t, err)

	ctx := context.Background()
	_, err = NewManager(ctx, "postgres", d)
	require.NoError(t, err)
	_, err = NewManager(ctx, "unknown", d)
	assert.True(t, errors.Match(errors.T(errors.InvalidParameter), err))

	d.Close()
	_, err = NewManager(ctx, "postgres", d)
	assert.True(t, errors.Match(errors.T(errors.Op("schema.NewManager")), err))
}

func TestSetup(t *testing.T) {
	_, _, _, err := docker.StartDbInDocker("postgres")
	require.NoError(t, err)
}

func TestRollForward(t *testing.T) {
	dialect := "postgres"
	c, u, _, err := docker.StartDbInDocker(dialect)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, c())
	})
	d, err := sql.Open(dialect, u)
	require.NoError(t, err)

	ctx := context.Background()
	m, err := NewManager(ctx, dialect, d)
	require.NoError(t, err)
	assert.NoError(t, m.RollForward(ctx))

	// Now set to dirty at an early version
	testDriver, err := postgres.NewPostgres(ctx, d)
	require.NoError(t, err)
	testDriver.SetVersion(ctx, 0, true)
	assert.Error(t, m.RollForward(ctx))
}

func TestRollForward_NotFromFresh(t *testing.T) {
	dialect := "postgres"
	oState := migrationStates[dialect]
	nState := migrationState{
		devMigration:   oState.devMigration,
		upMigrations:   make(map[int][]byte),
		downMigrations: make(map[int][]byte),
	}
	for k := range oState.upMigrations {
		if k > 8 {
			// Don't store any versions past our test version.
			continue
		}
		nState.upMigrations[k] = oState.upMigrations[k]
		nState.downMigrations[k] = oState.downMigrations[k]
		if nState.binarySchemaVersion < k {
			nState.binarySchemaVersion = k
		}
	}
	migrationStates[dialect] = nState

	c, u, _, err := docker.StartDbInDocker(dialect)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, c())
	})
	d, err := sql.Open(dialect, u)
	require.NoError(t, err)

	// Initialize the DB with only a portion of the current sql scripts.
	ctx := context.Background()
	m, err := NewManager(ctx, dialect, d)
	require.NoError(t, err)
	assert.NoError(t, m.RollForward(ctx))

	ver, dirty, err := m.driver.Version(ctx)
	assert.NoError(t, err)
	assert.Equal(t, nState.binarySchemaVersion, ver)
	assert.False(t, dirty)

	// Restore the full set of sql scripts and roll the rest of the way forward.
	migrationStates[dialect] = oState

	newM, err := NewManager(ctx, dialect, d)
	require.NoError(t, err)
	assert.NoError(t, newM.RollForward(ctx))
	ver, dirty, err = newM.driver.Version(ctx)
	assert.NoError(t, err)
	assert.Equal(t, oState.binarySchemaVersion, ver)
	assert.False(t, dirty)
}

func TestManager_ExclusiveLock(t *testing.T) {
	dialect := "postgres"
	c, u, _, err := docker.StartDbInDocker(dialect)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, c())
	})
	ctx := context.TODO()
	d1, err := sql.Open(dialect, u)
	require.NoError(t, err)
	m1, err := NewManager(ctx, dialect, d1)
	require.NoError(t, err)

	d2, err := sql.Open(dialect, u)
	require.NoError(t, err)
	m2, err := NewManager(ctx, dialect, d2)
	require.NoError(t, err)

	assert.NoError(t, m1.ExclusiveLock(ctx))
	assert.NoError(t, m1.ExclusiveLock(ctx))
	assert.Error(t, m2.ExclusiveLock(ctx))
	assert.Error(t, m2.SharedLock(ctx))
}

func TestManager_SharedLock(t *testing.T) {
	dialect := "postgres"
	c, u, _, err := docker.StartDbInDocker(dialect)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, c())
	})
	ctx := context.TODO()
	d1, err := sql.Open(dialect, u)
	require.NoError(t, err)
	m1, err := NewManager(ctx, dialect, d1)
	require.NoError(t, err)

	d2, err := sql.Open(dialect, u)
	require.NoError(t, err)
	m2, err := NewManager(ctx, dialect, d2)
	require.NoError(t, err)

	assert.NoError(t, m1.SharedLock(ctx))
	assert.NoError(t, m2.SharedLock(ctx))
	assert.NoError(t, m1.SharedLock(ctx))
	assert.NoError(t, m2.SharedLock(ctx))

	assert.Error(t, m1.ExclusiveLock(ctx))
	assert.Error(t, m2.ExclusiveLock(ctx))
}
