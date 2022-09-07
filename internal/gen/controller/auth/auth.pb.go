// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: controller/auth/v1/auth.proto

package auth

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// RequestInfo contains request parameters necessary for checking authn/authz
type RequestInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// path of the request
	Path string `protobuf:"bytes,10,opt,name=path,proto3" json:"path,omitempty"`
	// method (http verb) of the request
	Method string `protobuf:"bytes,20,opt,name=method,proto3" json:"method,omitempty"`
	// public_id for the request
	PublicId string `protobuf:"bytes,30,opt,name=public_id,json=publicId,proto3" json:"public_id,omitempty"`
	// encrypted_token from the request
	EncryptedToken string `protobuf:"bytes,40,opt,name=encrypted_token,json=encryptedToken,proto3" json:"encrypted_token,omitempty"`
	// token from the request
	Token string `protobuf:"bytes,50,opt,name=token,proto3" json:"token,omitempty"`
	// token_format of the request's token
	TokenFormat uint32 `protobuf:"varint,60,opt,name=token_format,json=tokenFormat,proto3" json:"token_format,omitempty"`
	// scope_id_override for the request (helpful for tests)
	ScopeIdOverride string `protobuf:"bytes,70,opt,name=scope_id_override,json=scopeIdOverride,proto3" json:"scope_id_override,omitempty"`
	// user_id_override for the request (helpful for tests)
	UserIdOverride string `protobuf:"bytes,80,opt,name=user_id_override,json=userIdOverride,proto3" json:"user_id_override,omitempty"`
	// disable_authz_failures for the request (helpful for tests)
	DisableAuthzFailures bool `protobuf:"varint,90,opt,name=disable_authz_failures,json=disableAuthzFailures,proto3" json:"disable_authz_failures,omitempty"`
	// disable_auth_entirely for the request (helpful for tests)
	DisableAuthEntirely bool `protobuf:"varint,100,opt,name=disable_auth_entirely,json=disableAuthEntirely,proto3" json:"disable_auth_entirely,omitempty"`
	// ticket is a unique id that allows the grpc-gateway to verify that the info
	// came from its companion http proxy
	Ticket string `protobuf:"bytes,110,opt,name=ticket,proto3" json:"ticket,omitempty"`
	// trace_id is the request's trace_id
	TraceId string `protobuf:"bytes,120,opt,name=trace_id,json=traceId,proto3" json:"trace_id,omitempty"`
	// event_id is the request's event id
	EventId string `protobuf:"bytes,130,opt,name=event_id,json=eventId,proto3" json:"event_id,omitempty"`
	// the client ip for the request
	ClientIp string `protobuf:"bytes,140,opt,name=client_ip,json=clientIp,proto3" json:"client_ip,omitempty"`
}

func (x *RequestInfo) Reset() {
	*x = RequestInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_controller_auth_v1_auth_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestInfo) ProtoMessage() {}

func (x *RequestInfo) ProtoReflect() protoreflect.Message {
	mi := &file_controller_auth_v1_auth_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestInfo.ProtoReflect.Descriptor instead.
func (*RequestInfo) Descriptor() ([]byte, []int) {
	return file_controller_auth_v1_auth_proto_rawDescGZIP(), []int{0}
}

func (x *RequestInfo) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *RequestInfo) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *RequestInfo) GetPublicId() string {
	if x != nil {
		return x.PublicId
	}
	return ""
}

func (x *RequestInfo) GetEncryptedToken() string {
	if x != nil {
		return x.EncryptedToken
	}
	return ""
}

func (x *RequestInfo) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *RequestInfo) GetTokenFormat() uint32 {
	if x != nil {
		return x.TokenFormat
	}
	return 0
}

func (x *RequestInfo) GetScopeIdOverride() string {
	if x != nil {
		return x.ScopeIdOverride
	}
	return ""
}

func (x *RequestInfo) GetUserIdOverride() string {
	if x != nil {
		return x.UserIdOverride
	}
	return ""
}

func (x *RequestInfo) GetDisableAuthzFailures() bool {
	if x != nil {
		return x.DisableAuthzFailures
	}
	return false
}

func (x *RequestInfo) GetDisableAuthEntirely() bool {
	if x != nil {
		return x.DisableAuthEntirely
	}
	return false
}

func (x *RequestInfo) GetTicket() string {
	if x != nil {
		return x.Ticket
	}
	return ""
}

func (x *RequestInfo) GetTraceId() string {
	if x != nil {
		return x.TraceId
	}
	return ""
}

func (x *RequestInfo) GetEventId() string {
	if x != nil {
		return x.EventId
	}
	return ""
}

func (x *RequestInfo) GetClientIp() string {
	if x != nil {
		return x.ClientIp
	}
	return ""
}

var File_controller_auth_v1_auth_proto protoreflect.FileDescriptor

var file_controller_auth_v1_auth_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2f, 0x61, 0x75, 0x74,
	0x68, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x12, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x61, 0x75, 0x74, 0x68,
	0x2e, 0x76, 0x31, 0x22, 0xe5, 0x03, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49,
	0x6e, 0x66, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x18, 0x14, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12,
	0x1b, 0x0a, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x1e, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x49, 0x64, 0x12, 0x27, 0x0a, 0x0f,
	0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18,
	0x28, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x32,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18, 0x3c, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x0b, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x2a,
	0x0a, 0x11, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x5f, 0x69, 0x64, 0x5f, 0x6f, 0x76, 0x65, 0x72, 0x72,
	0x69, 0x64, 0x65, 0x18, 0x46, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x73, 0x63, 0x6f, 0x70, 0x65,
	0x49, 0x64, 0x4f, 0x76, 0x65, 0x72, 0x72, 0x69, 0x64, 0x65, 0x12, 0x28, 0x0a, 0x10, 0x75, 0x73,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x5f, 0x6f, 0x76, 0x65, 0x72, 0x72, 0x69, 0x64, 0x65, 0x18, 0x50,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x4f, 0x76, 0x65, 0x72,
	0x72, 0x69, 0x64, 0x65, 0x12, 0x34, 0x0a, 0x16, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x5f,
	0x61, 0x75, 0x74, 0x68, 0x7a, 0x5f, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x73, 0x18, 0x5a,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x14, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x41, 0x75, 0x74,
	0x68, 0x7a, 0x46, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x73, 0x12, 0x32, 0x0a, 0x15, 0x64, 0x69,
	0x73, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x61, 0x75, 0x74, 0x68, 0x5f, 0x65, 0x6e, 0x74, 0x69, 0x72,
	0x65, 0x6c, 0x79, 0x18, 0x64, 0x20, 0x01, 0x28, 0x08, 0x52, 0x13, 0x64, 0x69, 0x73, 0x61, 0x62,
	0x6c, 0x65, 0x41, 0x75, 0x74, 0x68, 0x45, 0x6e, 0x74, 0x69, 0x72, 0x65, 0x6c, 0x79, 0x12, 0x16,
	0x0a, 0x06, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x6e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x74, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x74, 0x72, 0x61, 0x63, 0x65, 0x5f,
	0x69, 0x64, 0x18, 0x78, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x74, 0x72, 0x61, 0x63, 0x65, 0x49,
	0x64, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x82, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x1c, 0x0a,
	0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x70, 0x18, 0x8c, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x70, 0x42, 0x41, 0x5a, 0x3f, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x61, 0x73, 0x68, 0x69, 0x63,
	0x6f, 0x72, 0x70, 0x2f, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x61, 0x72, 0x79, 0x2f, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f,
	0x6c, 0x6c, 0x65, 0x72, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x3b, 0x61, 0x75, 0x74, 0x68, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_controller_auth_v1_auth_proto_rawDescOnce sync.Once
	file_controller_auth_v1_auth_proto_rawDescData = file_controller_auth_v1_auth_proto_rawDesc
)

func file_controller_auth_v1_auth_proto_rawDescGZIP() []byte {
	file_controller_auth_v1_auth_proto_rawDescOnce.Do(func() {
		file_controller_auth_v1_auth_proto_rawDescData = protoimpl.X.CompressGZIP(file_controller_auth_v1_auth_proto_rawDescData)
	})
	return file_controller_auth_v1_auth_proto_rawDescData
}

var file_controller_auth_v1_auth_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_controller_auth_v1_auth_proto_goTypes = []interface{}{
	(*RequestInfo)(nil), // 0: controller.auth.v1.RequestInfo
}
var file_controller_auth_v1_auth_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_controller_auth_v1_auth_proto_init() }
func file_controller_auth_v1_auth_proto_init() {
	if File_controller_auth_v1_auth_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_controller_auth_v1_auth_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_controller_auth_v1_auth_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_controller_auth_v1_auth_proto_goTypes,
		DependencyIndexes: file_controller_auth_v1_auth_proto_depIdxs,
		MessageInfos:      file_controller_auth_v1_auth_proto_msgTypes,
	}.Build()
	File_controller_auth_v1_auth_proto = out.File
	file_controller_auth_v1_auth_proto_rawDesc = nil
	file_controller_auth_v1_auth_proto_goTypes = nil
	file_controller_auth_v1_auth_proto_depIdxs = nil
}
