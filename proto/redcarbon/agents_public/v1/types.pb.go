// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: redcarbon/agents_public/v1/types.proto

package agents_publicv1

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

type AgentConfiguration struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QradarJobConfiguration      *QRadarJobConfiguration      `protobuf:"bytes,1,opt,name=qradar_job_configuration,json=qradarJobConfiguration,proto3,oneof" json:"qradar_job_configuration,omitempty"`
	SentineloneJobConfiguration *SentinelOneJobConfiguration `protobuf:"bytes,2,opt,name=sentinelone_job_configuration,json=sentineloneJobConfiguration,proto3,oneof" json:"sentinelone_job_configuration,omitempty"`
	FortisiemJobConfiguration   *FortiSIEMJobConfiguration   `protobuf:"bytes,3,opt,name=fortisiem_job_configuration,json=fortisiemJobConfiguration,proto3,oneof" json:"fortisiem_job_configuration,omitempty"`
}

func (x *AgentConfiguration) Reset() {
	*x = AgentConfiguration{}
	if protoimpl.UnsafeEnabled {
		mi := &file_redcarbon_agents_public_v1_types_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AgentConfiguration) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AgentConfiguration) ProtoMessage() {}

func (x *AgentConfiguration) ProtoReflect() protoreflect.Message {
	mi := &file_redcarbon_agents_public_v1_types_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AgentConfiguration.ProtoReflect.Descriptor instead.
func (*AgentConfiguration) Descriptor() ([]byte, []int) {
	return file_redcarbon_agents_public_v1_types_proto_rawDescGZIP(), []int{0}
}

func (x *AgentConfiguration) GetQradarJobConfiguration() *QRadarJobConfiguration {
	if x != nil {
		return x.QradarJobConfiguration
	}
	return nil
}

func (x *AgentConfiguration) GetSentineloneJobConfiguration() *SentinelOneJobConfiguration {
	if x != nil {
		return x.SentineloneJobConfiguration
	}
	return nil
}

func (x *AgentConfiguration) GetFortisiemJobConfiguration() *FortiSIEMJobConfiguration {
	if x != nil {
		return x.FortisiemJobConfiguration
	}
	return nil
}

type QRadarJobConfiguration struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host      string `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	Token     string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	VerifySsl bool   `protobuf:"varint,3,opt,name=verify_ssl,json=verifySsl,proto3" json:"verify_ssl,omitempty"`
}

func (x *QRadarJobConfiguration) Reset() {
	*x = QRadarJobConfiguration{}
	if protoimpl.UnsafeEnabled {
		mi := &file_redcarbon_agents_public_v1_types_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QRadarJobConfiguration) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QRadarJobConfiguration) ProtoMessage() {}

func (x *QRadarJobConfiguration) ProtoReflect() protoreflect.Message {
	mi := &file_redcarbon_agents_public_v1_types_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QRadarJobConfiguration.ProtoReflect.Descriptor instead.
func (*QRadarJobConfiguration) Descriptor() ([]byte, []int) {
	return file_redcarbon_agents_public_v1_types_proto_rawDescGZIP(), []int{1}
}

func (x *QRadarJobConfiguration) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *QRadarJobConfiguration) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *QRadarJobConfiguration) GetVerifySsl() bool {
	if x != nil {
		return x.VerifySsl
	}
	return false
}

type SentinelOneJobConfiguration struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host      string `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	Token     string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	VerifySsl bool   `protobuf:"varint,3,opt,name=verify_ssl,json=verifySsl,proto3" json:"verify_ssl,omitempty"`
}

func (x *SentinelOneJobConfiguration) Reset() {
	*x = SentinelOneJobConfiguration{}
	if protoimpl.UnsafeEnabled {
		mi := &file_redcarbon_agents_public_v1_types_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SentinelOneJobConfiguration) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SentinelOneJobConfiguration) ProtoMessage() {}

func (x *SentinelOneJobConfiguration) ProtoReflect() protoreflect.Message {
	mi := &file_redcarbon_agents_public_v1_types_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SentinelOneJobConfiguration.ProtoReflect.Descriptor instead.
func (*SentinelOneJobConfiguration) Descriptor() ([]byte, []int) {
	return file_redcarbon_agents_public_v1_types_proto_rawDescGZIP(), []int{2}
}

func (x *SentinelOneJobConfiguration) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *SentinelOneJobConfiguration) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *SentinelOneJobConfiguration) GetVerifySsl() bool {
	if x != nil {
		return x.VerifySsl
	}
	return false
}

type FortiSIEMJobConfiguration struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host      string `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	Username  string `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Password  string `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	VerifySsl bool   `protobuf:"varint,4,opt,name=verify_ssl,json=verifySsl,proto3" json:"verify_ssl,omitempty"`
}

func (x *FortiSIEMJobConfiguration) Reset() {
	*x = FortiSIEMJobConfiguration{}
	if protoimpl.UnsafeEnabled {
		mi := &file_redcarbon_agents_public_v1_types_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FortiSIEMJobConfiguration) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FortiSIEMJobConfiguration) ProtoMessage() {}

func (x *FortiSIEMJobConfiguration) ProtoReflect() protoreflect.Message {
	mi := &file_redcarbon_agents_public_v1_types_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FortiSIEMJobConfiguration.ProtoReflect.Descriptor instead.
func (*FortiSIEMJobConfiguration) Descriptor() ([]byte, []int) {
	return file_redcarbon_agents_public_v1_types_proto_rawDescGZIP(), []int{3}
}

func (x *FortiSIEMJobConfiguration) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *FortiSIEMJobConfiguration) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *FortiSIEMJobConfiguration) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *FortiSIEMJobConfiguration) GetVerifySsl() bool {
	if x != nil {
		return x.VerifySsl
	}
	return false
}

var File_redcarbon_agents_public_v1_types_proto protoreflect.FileDescriptor

var file_redcarbon_agents_public_v1_types_proto_rawDesc = []byte{
	0x0a, 0x26, 0x72, 0x65, 0x64, 0x63, 0x61, 0x72, 0x62, 0x6f, 0x6e, 0x2f, 0x61, 0x67, 0x65, 0x6e,
	0x74, 0x73, 0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x79, 0x70,
	0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1a, 0x72, 0x65, 0x64, 0x63, 0x61, 0x72,
	0x62, 0x6f, 0x6e, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x73, 0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69,
	0x63, 0x2e, 0x76, 0x31, 0x22, 0xe5, 0x03, 0x0a, 0x12, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x71, 0x0a, 0x18, 0x71,
	0x72, 0x61, 0x64, 0x61, 0x72, 0x5f, 0x6a, 0x6f, 0x62, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x32, 0x2e,
	0x72, 0x65, 0x64, 0x63, 0x61, 0x72, 0x62, 0x6f, 0x6e, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x73,
	0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x52, 0x61, 0x64, 0x61,
	0x72, 0x4a, 0x6f, 0x62, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x48, 0x00, 0x52, 0x16, 0x71, 0x72, 0x61, 0x64, 0x61, 0x72, 0x4a, 0x6f, 0x62, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x80,
	0x01, 0x0a, 0x1d, 0x73, 0x65, 0x6e, 0x74, 0x69, 0x6e, 0x65, 0x6c, 0x6f, 0x6e, 0x65, 0x5f, 0x6a,
	0x6f, 0x62, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x37, 0x2e, 0x72, 0x65, 0x64, 0x63, 0x61, 0x72, 0x62,
	0x6f, 0x6e, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x73, 0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63,
	0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x74, 0x69, 0x6e, 0x65, 0x6c, 0x4f, 0x6e, 0x65, 0x4a,
	0x6f, 0x62, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x48,
	0x01, 0x52, 0x1b, 0x73, 0x65, 0x6e, 0x74, 0x69, 0x6e, 0x65, 0x6c, 0x6f, 0x6e, 0x65, 0x4a, 0x6f,
	0x62, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01,
	0x01, 0x12, 0x7a, 0x0a, 0x1b, 0x66, 0x6f, 0x72, 0x74, 0x69, 0x73, 0x69, 0x65, 0x6d, 0x5f, 0x6a,
	0x6f, 0x62, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x35, 0x2e, 0x72, 0x65, 0x64, 0x63, 0x61, 0x72, 0x62,
	0x6f, 0x6e, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x73, 0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63,
	0x2e, 0x76, 0x31, 0x2e, 0x46, 0x6f, 0x72, 0x74, 0x69, 0x53, 0x49, 0x45, 0x4d, 0x4a, 0x6f, 0x62,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x02, 0x52,
	0x19, 0x66, 0x6f, 0x72, 0x74, 0x69, 0x73, 0x69, 0x65, 0x6d, 0x4a, 0x6f, 0x62, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x42, 0x1b, 0x0a,
	0x19, 0x5f, 0x71, 0x72, 0x61, 0x64, 0x61, 0x72, 0x5f, 0x6a, 0x6f, 0x62, 0x5f, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x20, 0x0a, 0x1e, 0x5f, 0x73,
	0x65, 0x6e, 0x74, 0x69, 0x6e, 0x65, 0x6c, 0x6f, 0x6e, 0x65, 0x5f, 0x6a, 0x6f, 0x62, 0x5f, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x1e, 0x0a, 0x1c,
	0x5f, 0x66, 0x6f, 0x72, 0x74, 0x69, 0x73, 0x69, 0x65, 0x6d, 0x5f, 0x6a, 0x6f, 0x62, 0x5f, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x61, 0x0a, 0x16,
	0x51, 0x52, 0x61, 0x64, 0x61, 0x72, 0x4a, 0x6f, 0x62, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x12, 0x1d, 0x0a, 0x0a, 0x76, 0x65, 0x72, 0x69, 0x66, 0x79, 0x5f, 0x73, 0x73, 0x6c, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x76, 0x65, 0x72, 0x69, 0x66, 0x79, 0x53, 0x73, 0x6c, 0x22,
	0x66, 0x0a, 0x1b, 0x53, 0x65, 0x6e, 0x74, 0x69, 0x6e, 0x65, 0x6c, 0x4f, 0x6e, 0x65, 0x4a, 0x6f,
	0x62, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12,
	0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f,
	0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x76, 0x65, 0x72, 0x69,
	0x66, 0x79, 0x5f, 0x73, 0x73, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x76, 0x65,
	0x72, 0x69, 0x66, 0x79, 0x53, 0x73, 0x6c, 0x22, 0x86, 0x01, 0x0a, 0x19, 0x46, 0x6f, 0x72, 0x74,
	0x69, 0x53, 0x49, 0x45, 0x4d, 0x4a, 0x6f, 0x62, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x76, 0x65, 0x72, 0x69, 0x66, 0x79, 0x5f, 0x73, 0x73, 0x6c, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x76, 0x65, 0x72, 0x69, 0x66, 0x79, 0x53, 0x73, 0x6c,
	0x42, 0xf5, 0x01, 0x0a, 0x1e, 0x63, 0x6f, 0x6d, 0x2e, 0x72, 0x65, 0x64, 0x63, 0x61, 0x72, 0x62,
	0x6f, 0x6e, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x73, 0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63,
	0x2e, 0x76, 0x31, 0x42, 0x0a, 0x54, 0x79, 0x70, 0x65, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50,
	0x01, 0x5a, 0x41, 0x70, 0x6b, 0x67, 0x2e, 0x72, 0x65, 0x64, 0x63, 0x61, 0x72, 0x62, 0x6f, 0x6e,
	0x2e, 0x61, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x64, 0x63, 0x61, 0x72,
	0x62, 0x6f, 0x6e, 0x2f, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x73, 0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69,
	0x63, 0x2f, 0x76, 0x31, 0x3b, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x73, 0x5f, 0x70, 0x75, 0x62, 0x6c,
	0x69, 0x63, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x52, 0x41, 0x58, 0xaa, 0x02, 0x19, 0x52, 0x65, 0x64,
	0x63, 0x61, 0x72, 0x62, 0x6f, 0x6e, 0x2e, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x73, 0x50, 0x75, 0x62,
	0x6c, 0x69, 0x63, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x19, 0x52, 0x65, 0x64, 0x63, 0x61, 0x72, 0x62,
	0x6f, 0x6e, 0x5c, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x73, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5c,
	0x56, 0x31, 0xe2, 0x02, 0x25, 0x52, 0x65, 0x64, 0x63, 0x61, 0x72, 0x62, 0x6f, 0x6e, 0x5c, 0x41,
	0x67, 0x65, 0x6e, 0x74, 0x73, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5c, 0x56, 0x31, 0x5c, 0x47,
	0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x1b, 0x52, 0x65, 0x64,
	0x63, 0x61, 0x72, 0x62, 0x6f, 0x6e, 0x3a, 0x3a, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x73, 0x50, 0x75,
	0x62, 0x6c, 0x69, 0x63, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_redcarbon_agents_public_v1_types_proto_rawDescOnce sync.Once
	file_redcarbon_agents_public_v1_types_proto_rawDescData = file_redcarbon_agents_public_v1_types_proto_rawDesc
)

func file_redcarbon_agents_public_v1_types_proto_rawDescGZIP() []byte {
	file_redcarbon_agents_public_v1_types_proto_rawDescOnce.Do(func() {
		file_redcarbon_agents_public_v1_types_proto_rawDescData = protoimpl.X.CompressGZIP(file_redcarbon_agents_public_v1_types_proto_rawDescData)
	})
	return file_redcarbon_agents_public_v1_types_proto_rawDescData
}

var file_redcarbon_agents_public_v1_types_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_redcarbon_agents_public_v1_types_proto_goTypes = []any{
	(*AgentConfiguration)(nil),          // 0: redcarbon.agents_public.v1.AgentConfiguration
	(*QRadarJobConfiguration)(nil),      // 1: redcarbon.agents_public.v1.QRadarJobConfiguration
	(*SentinelOneJobConfiguration)(nil), // 2: redcarbon.agents_public.v1.SentinelOneJobConfiguration
	(*FortiSIEMJobConfiguration)(nil),   // 3: redcarbon.agents_public.v1.FortiSIEMJobConfiguration
}
var file_redcarbon_agents_public_v1_types_proto_depIdxs = []int32{
	1, // 0: redcarbon.agents_public.v1.AgentConfiguration.qradar_job_configuration:type_name -> redcarbon.agents_public.v1.QRadarJobConfiguration
	2, // 1: redcarbon.agents_public.v1.AgentConfiguration.sentinelone_job_configuration:type_name -> redcarbon.agents_public.v1.SentinelOneJobConfiguration
	3, // 2: redcarbon.agents_public.v1.AgentConfiguration.fortisiem_job_configuration:type_name -> redcarbon.agents_public.v1.FortiSIEMJobConfiguration
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_redcarbon_agents_public_v1_types_proto_init() }
func file_redcarbon_agents_public_v1_types_proto_init() {
	if File_redcarbon_agents_public_v1_types_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_redcarbon_agents_public_v1_types_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*AgentConfiguration); i {
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
		file_redcarbon_agents_public_v1_types_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*QRadarJobConfiguration); i {
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
		file_redcarbon_agents_public_v1_types_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*SentinelOneJobConfiguration); i {
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
		file_redcarbon_agents_public_v1_types_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*FortiSIEMJobConfiguration); i {
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
	file_redcarbon_agents_public_v1_types_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_redcarbon_agents_public_v1_types_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_redcarbon_agents_public_v1_types_proto_goTypes,
		DependencyIndexes: file_redcarbon_agents_public_v1_types_proto_depIdxs,
		MessageInfos:      file_redcarbon_agents_public_v1_types_proto_msgTypes,
	}.Build()
	File_redcarbon_agents_public_v1_types_proto = out.File
	file_redcarbon_agents_public_v1_types_proto_rawDesc = nil
	file_redcarbon_agents_public_v1_types_proto_goTypes = nil
	file_redcarbon_agents_public_v1_types_proto_depIdxs = nil
}
