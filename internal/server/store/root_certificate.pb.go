// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: controller/storage/servers/store/v1/root_certificate.proto

// Package store provides protobufs for storing types in the pki package.

package store

import (
	timestamp "github.com/hashicorp/boundary/internal/db/timestamp"
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

// CertificateAuthority is a versioned entity used to lock the database when rotation RootCertificates
type CertificateAuthority struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// @inject_tag: `gorm:"primary_key"`
	PrivateId string `protobuf:"bytes,10,opt,name=private_id,json=privateId,proto3" json:"private_id,omitempty" gorm:"primary_key"`
	// version allows optimistic locking of the resource.
	// @inject_tag: `gorm:"default:null"`
	Version uint32 `protobuf:"varint,20,opt,name=version,proto3" json:"version,omitempty" gorm:"default:null"`
}

func (x *CertificateAuthority) Reset() {
	*x = CertificateAuthority{}
	if protoimpl.UnsafeEnabled {
		mi := &file_controller_storage_servers_store_v1_root_certificate_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CertificateAuthority) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CertificateAuthority) ProtoMessage() {}

func (x *CertificateAuthority) ProtoReflect() protoreflect.Message {
	mi := &file_controller_storage_servers_store_v1_root_certificate_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CertificateAuthority.ProtoReflect.Descriptor instead.
func (*CertificateAuthority) Descriptor() ([]byte, []int) {
	return file_controller_storage_servers_store_v1_root_certificate_proto_rawDescGZIP(), []int{0}
}

func (x *CertificateAuthority) GetPrivateId() string {
	if x != nil {
		return x.PrivateId
	}
	return ""
}

func (x *CertificateAuthority) GetVersion() uint32 {
	if x != nil {
		return x.Version
	}
	return 0
}

// RootCertificate contains all fields related to a RootCertificate resource
type RootCertificate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The serial number of the root certificate
	// @inject_tag: `gorm:"not_null"`
	SerialNumber uint64 `protobuf:"varint,10,opt,name=serial_number,json=serialNumber,proto3" json:"serial_number,omitempty" gorm:"not_null"`
	// Certificate is the PEM encoded certificate.
	// @inject_tag: `gorm:"not_null"`
	Certificate []byte `protobuf:"bytes,20,opt,name=certificate,proto3" json:"certificate,omitempty" gorm:"not_null"`
	// Not valid before is the timestamp at which this certificate's validity period starts
	NotValidBefore *timestamp.Timestamp `protobuf:"bytes,30,opt,name=not_valid_before,json=notValidBefore,proto3" json:"not_valid_before,omitempty"`
	// Not valid after is the timestamp at which this certificate's validity period ends
	NotValidAfter *timestamp.Timestamp `protobuf:"bytes,40,opt,name=not_valid_after,json=notValidAfter,proto3" json:"not_valid_after,omitempty"`
	// The public key associated with this certificate
	// @inject_tag: `gorm:"primary_key;not_null"`
	PublicKey []byte `protobuf:"bytes,50,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty" gorm:"primary_key;not_null"`
	// The private key associated with this certificate
	// This is a ciphertext field
	// @inject_tag: `gorm:"not_null"`
	PrivateKey []byte `protobuf:"bytes,60,opt,name=private_key,json=privateKey,proto3" json:"private_key,omitempty" gorm:"not_null"`
	// The id of the kms database key used for encrypting this entry.
	// @inject_tag: `gorm:"not_null"`
	KeyId string `protobuf:"bytes,70,opt,name=key_id,json=keyId,proto3" json:"key_id,omitempty" gorm:"not_null"`
	// State is an enum value indicating if this is the next or current root cert
	// @inject_tag: `gorm:"not_null"`
	State string `protobuf:"bytes,80,opt,name=state,proto3" json:"state,omitempty" gorm:"not_null"`
	// A reference to the CertificateAuthority
	// @inject_tag: `gorm:"not_null"`
	IssuingCa string `protobuf:"bytes,90,opt,name=issuing_ca,json=issuingCa,proto3" json:"issuing_ca,omitempty" gorm:"not_null"`
}

func (x *RootCertificate) Reset() {
	*x = RootCertificate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_controller_storage_servers_store_v1_root_certificate_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RootCertificate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RootCertificate) ProtoMessage() {}

func (x *RootCertificate) ProtoReflect() protoreflect.Message {
	mi := &file_controller_storage_servers_store_v1_root_certificate_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RootCertificate.ProtoReflect.Descriptor instead.
func (*RootCertificate) Descriptor() ([]byte, []int) {
	return file_controller_storage_servers_store_v1_root_certificate_proto_rawDescGZIP(), []int{1}
}

func (x *RootCertificate) GetSerialNumber() uint64 {
	if x != nil {
		return x.SerialNumber
	}
	return 0
}

func (x *RootCertificate) GetCertificate() []byte {
	if x != nil {
		return x.Certificate
	}
	return nil
}

func (x *RootCertificate) GetNotValidBefore() *timestamp.Timestamp {
	if x != nil {
		return x.NotValidBefore
	}
	return nil
}

func (x *RootCertificate) GetNotValidAfter() *timestamp.Timestamp {
	if x != nil {
		return x.NotValidAfter
	}
	return nil
}

func (x *RootCertificate) GetPublicKey() []byte {
	if x != nil {
		return x.PublicKey
	}
	return nil
}

func (x *RootCertificate) GetPrivateKey() []byte {
	if x != nil {
		return x.PrivateKey
	}
	return nil
}

func (x *RootCertificate) GetKeyId() string {
	if x != nil {
		return x.KeyId
	}
	return ""
}

func (x *RootCertificate) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *RootCertificate) GetIssuingCa() string {
	if x != nil {
		return x.IssuingCa
	}
	return ""
}

var File_controller_storage_servers_store_v1_root_certificate_proto protoreflect.FileDescriptor

var file_controller_storage_servers_store_v1_root_certificate_proto_rawDesc = []byte{
	0x0a, 0x3a, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2f, 0x73, 0x74, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x73, 0x2f, 0x73, 0x74, 0x6f,
	0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x6f, 0x6f, 0x74, 0x5f, 0x63, 0x65, 0x72, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x23, 0x63, 0x6f,
	0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x73, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76,
	0x31, 0x1a, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2f, 0x73, 0x74,
	0x6f, 0x72, 0x61, 0x67, 0x65, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2f,
	0x76, 0x31, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x4f, 0x0a, 0x14, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x65, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72,
	0x69, 0x76, 0x61, 0x74, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x18, 0x14, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x22, 0x8e, 0x03, 0x0a, 0x0f, 0x52, 0x6f, 0x6f, 0x74, 0x43, 0x65, 0x72, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x65, 0x72, 0x69, 0x61,
	0x6c, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c,
	0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x20, 0x0a, 0x0b,
	0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x18, 0x14, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x0b, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x12, 0x54,
	0x0a, 0x10, 0x6e, 0x6f, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x5f, 0x62, 0x65, 0x66, 0x6f,
	0x72, 0x65, 0x18, 0x1e, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72,
	0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x0e, 0x6e, 0x6f, 0x74, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x42, 0x65,
	0x66, 0x6f, 0x72, 0x65, 0x12, 0x52, 0x0a, 0x0f, 0x6e, 0x6f, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x5f, 0x61, 0x66, 0x74, 0x65, 0x72, 0x18, 0x28, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x76, 0x31, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0d, 0x6e, 0x6f, 0x74, 0x56, 0x61,
	0x6c, 0x69, 0x64, 0x41, 0x66, 0x74, 0x65, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x75, 0x62, 0x6c,
	0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x32, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x70, 0x75,
	0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x72, 0x69, 0x76, 0x61,
	0x74, 0x65, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x3c, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x70, 0x72,
	0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x12, 0x15, 0x0a, 0x06, 0x6b, 0x65, 0x79, 0x5f,
	0x69, 0x64, 0x18, 0x46, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6b, 0x65, 0x79, 0x49, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x50, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x73, 0x73, 0x75, 0x69, 0x6e, 0x67,
	0x5f, 0x63, 0x61, 0x18, 0x5a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x69, 0x73, 0x73, 0x75, 0x69,
	0x6e, 0x67, 0x43, 0x61, 0x42, 0x3b, 0x5a, 0x39, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x68, 0x61, 0x73, 0x68, 0x69, 0x63, 0x6f, 0x72, 0x70, 0x2f, 0x62, 0x6f, 0x75,
	0x6e, 0x64, 0x61, 0x72, 0x79, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x3b, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_controller_storage_servers_store_v1_root_certificate_proto_rawDescOnce sync.Once
	file_controller_storage_servers_store_v1_root_certificate_proto_rawDescData = file_controller_storage_servers_store_v1_root_certificate_proto_rawDesc
)

func file_controller_storage_servers_store_v1_root_certificate_proto_rawDescGZIP() []byte {
	file_controller_storage_servers_store_v1_root_certificate_proto_rawDescOnce.Do(func() {
		file_controller_storage_servers_store_v1_root_certificate_proto_rawDescData = protoimpl.X.CompressGZIP(file_controller_storage_servers_store_v1_root_certificate_proto_rawDescData)
	})
	return file_controller_storage_servers_store_v1_root_certificate_proto_rawDescData
}

var file_controller_storage_servers_store_v1_root_certificate_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_controller_storage_servers_store_v1_root_certificate_proto_goTypes = []interface{}{
	(*CertificateAuthority)(nil), // 0: controller.storage.servers.store.v1.CertificateAuthority
	(*RootCertificate)(nil),      // 1: controller.storage.servers.store.v1.RootCertificate
	(*timestamp.Timestamp)(nil),  // 2: controller.storage.timestamp.v1.Timestamp
}
var file_controller_storage_servers_store_v1_root_certificate_proto_depIdxs = []int32{
	2, // 0: controller.storage.servers.store.v1.RootCertificate.not_valid_before:type_name -> controller.storage.timestamp.v1.Timestamp
	2, // 1: controller.storage.servers.store.v1.RootCertificate.not_valid_after:type_name -> controller.storage.timestamp.v1.Timestamp
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_controller_storage_servers_store_v1_root_certificate_proto_init() }
func file_controller_storage_servers_store_v1_root_certificate_proto_init() {
	if File_controller_storage_servers_store_v1_root_certificate_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_controller_storage_servers_store_v1_root_certificate_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CertificateAuthority); i {
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
		file_controller_storage_servers_store_v1_root_certificate_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RootCertificate); i {
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
			RawDescriptor: file_controller_storage_servers_store_v1_root_certificate_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_controller_storage_servers_store_v1_root_certificate_proto_goTypes,
		DependencyIndexes: file_controller_storage_servers_store_v1_root_certificate_proto_depIdxs,
		MessageInfos:      file_controller_storage_servers_store_v1_root_certificate_proto_msgTypes,
	}.Build()
	File_controller_storage_servers_store_v1_root_certificate_proto = out.File
	file_controller_storage_servers_store_v1_root_certificate_proto_rawDesc = nil
	file_controller_storage_servers_store_v1_root_certificate_proto_goTypes = nil
	file_controller_storage_servers_store_v1_root_certificate_proto_depIdxs = nil
}
