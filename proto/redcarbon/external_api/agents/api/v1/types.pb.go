// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1-devel
// 	protoc        (unknown)
// source: redcarbon/external_api/agents/api/v1/types.proto

package agentsExternalApiV1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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

	AgentConfigurationId string                  `protobuf:"bytes,1,opt,name=agent_configuration_id,json=agentConfigurationId,proto3" json:"agent_configuration_id,omitempty"`
	CreatedAt            *timestamppb.Timestamp  `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt            *timestamppb.Timestamp  `protobuf:"bytes,3,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	Name                 string                  `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	Data                 *AgentConfigurationData `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
	Type                 string                  `protobuf:"bytes,6,opt,name=type,proto3" json:"type,omitempty"`
}

func (x *AgentConfiguration) Reset() {
	*x = AgentConfiguration{}
	if protoimpl.UnsafeEnabled {
		mi := &file_redcarbon_external_api_agents_api_v1_types_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AgentConfiguration) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AgentConfiguration) ProtoMessage() {}

func (x *AgentConfiguration) ProtoReflect() protoreflect.Message {
	mi := &file_redcarbon_external_api_agents_api_v1_types_proto_msgTypes[0]
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
	return file_redcarbon_external_api_agents_api_v1_types_proto_rawDescGZIP(), []int{0}
}

func (x *AgentConfiguration) GetAgentConfigurationId() string {
	if x != nil {
		return x.AgentConfigurationId
	}
	return ""
}

func (x *AgentConfiguration) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *AgentConfiguration) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *AgentConfiguration) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *AgentConfiguration) GetData() *AgentConfigurationData {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *AgentConfiguration) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

type AgentConfigurationData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Data:
	//
	//	*AgentConfigurationData_SentinelOne
	Data isAgentConfigurationData_Data `protobuf_oneof:"data"`
}

func (x *AgentConfigurationData) Reset() {
	*x = AgentConfigurationData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_redcarbon_external_api_agents_api_v1_types_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AgentConfigurationData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AgentConfigurationData) ProtoMessage() {}

func (x *AgentConfigurationData) ProtoReflect() protoreflect.Message {
	mi := &file_redcarbon_external_api_agents_api_v1_types_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AgentConfigurationData.ProtoReflect.Descriptor instead.
func (*AgentConfigurationData) Descriptor() ([]byte, []int) {
	return file_redcarbon_external_api_agents_api_v1_types_proto_rawDescGZIP(), []int{1}
}

func (m *AgentConfigurationData) GetData() isAgentConfigurationData_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *AgentConfigurationData) GetSentinelOne() *SentinelOneData {
	if x, ok := x.GetData().(*AgentConfigurationData_SentinelOne); ok {
		return x.SentinelOne
	}
	return nil
}

type isAgentConfigurationData_Data interface {
	isAgentConfigurationData_Data()
}

type AgentConfigurationData_SentinelOne struct {
	SentinelOne *SentinelOneData `protobuf:"bytes,1,opt,name=sentinel_one,json=sentinelOne,proto3,oneof"`
}

func (*AgentConfigurationData_SentinelOne) isAgentConfigurationData_Data() {}

type SentinelOneData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ApiToken string `protobuf:"bytes,1,opt,name=api_token,json=apiToken,proto3" json:"api_token,omitempty"`
	Url      string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *SentinelOneData) Reset() {
	*x = SentinelOneData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_redcarbon_external_api_agents_api_v1_types_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SentinelOneData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SentinelOneData) ProtoMessage() {}

func (x *SentinelOneData) ProtoReflect() protoreflect.Message {
	mi := &file_redcarbon_external_api_agents_api_v1_types_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SentinelOneData.ProtoReflect.Descriptor instead.
func (*SentinelOneData) Descriptor() ([]byte, []int) {
	return file_redcarbon_external_api_agents_api_v1_types_proto_rawDescGZIP(), []int{2}
}

func (x *SentinelOneData) GetApiToken() string {
	if x != nil {
		return x.ApiToken
	}
	return ""
}

func (x *SentinelOneData) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

var File_redcarbon_external_api_agents_api_v1_types_proto protoreflect.FileDescriptor

var file_redcarbon_external_api_agents_api_v1_types_proto_rawDesc = []byte{
	0x0a, 0x30, 0x72, 0x65, 0x64, 0x63, 0x61, 0x72, 0x62, 0x6f, 0x6e, 0x2f, 0x65, 0x78, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x73, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x24, 0x72, 0x65, 0x64, 0x63, 0x61, 0x72, 0x62, 0x6f, 0x6e, 0x2e, 0x65, 0x78,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74,
	0x73, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xba, 0x02, 0x0a, 0x12, 0x41, 0x67,
	0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x34, 0x0a, 0x16, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x14, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x50, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x3c,
	0x2e, 0x72, 0x65, 0x64, 0x63, 0x61, 0x72, 0x62, 0x6f, 0x6e, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x7c, 0x0a, 0x16, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x5a, 0x0a, 0x0c, 0x73, 0x65, 0x6e, 0x74, 0x69, 0x6e, 0x65, 0x6c, 0x5f, 0x6f, 0x6e, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x35, 0x2e, 0x72, 0x65, 0x64, 0x63, 0x61, 0x72, 0x62,
	0x6f, 0x6e, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x61, 0x70, 0x69, 0x2e,
	0x61, 0x67, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65,
	0x6e, 0x74, 0x69, 0x6e, 0x65, 0x6c, 0x4f, 0x6e, 0x65, 0x44, 0x61, 0x74, 0x61, 0x48, 0x00, 0x52,
	0x0b, 0x73, 0x65, 0x6e, 0x74, 0x69, 0x6e, 0x65, 0x6c, 0x4f, 0x6e, 0x65, 0x42, 0x06, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x22, 0x40, 0x0a, 0x0f, 0x53, 0x65, 0x6e, 0x74, 0x69, 0x6e, 0x65, 0x6c,
	0x4f, 0x6e, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1b, 0x0a, 0x09, 0x61, 0x70, 0x69, 0x5f, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x61, 0x70, 0x69, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x42, 0x51, 0x5a, 0x4f, 0x70, 0x6b, 0x67, 0x2e, 0x72, 0x65,
	0x64, 0x63, 0x61, 0x72, 0x62, 0x6f, 0x6e, 0x2e, 0x61, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x72, 0x65, 0x64, 0x63, 0x61, 0x72, 0x62, 0x6f, 0x6e, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x5f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x76, 0x31, 0x3b, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x73, 0x45, 0x78, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x41, 0x70, 0x69, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_redcarbon_external_api_agents_api_v1_types_proto_rawDescOnce sync.Once
	file_redcarbon_external_api_agents_api_v1_types_proto_rawDescData = file_redcarbon_external_api_agents_api_v1_types_proto_rawDesc
)

func file_redcarbon_external_api_agents_api_v1_types_proto_rawDescGZIP() []byte {
	file_redcarbon_external_api_agents_api_v1_types_proto_rawDescOnce.Do(func() {
		file_redcarbon_external_api_agents_api_v1_types_proto_rawDescData = protoimpl.X.CompressGZIP(file_redcarbon_external_api_agents_api_v1_types_proto_rawDescData)
	})
	return file_redcarbon_external_api_agents_api_v1_types_proto_rawDescData
}

var file_redcarbon_external_api_agents_api_v1_types_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_redcarbon_external_api_agents_api_v1_types_proto_goTypes = []interface{}{
	(*AgentConfiguration)(nil),     // 0: redcarbon.external_api.agents.api.v1.AgentConfiguration
	(*AgentConfigurationData)(nil), // 1: redcarbon.external_api.agents.api.v1.AgentConfigurationData
	(*SentinelOneData)(nil),        // 2: redcarbon.external_api.agents.api.v1.SentinelOneData
	(*timestamppb.Timestamp)(nil),  // 3: google.protobuf.Timestamp
}
var file_redcarbon_external_api_agents_api_v1_types_proto_depIdxs = []int32{
	3, // 0: redcarbon.external_api.agents.api.v1.AgentConfiguration.created_at:type_name -> google.protobuf.Timestamp
	3, // 1: redcarbon.external_api.agents.api.v1.AgentConfiguration.updated_at:type_name -> google.protobuf.Timestamp
	1, // 2: redcarbon.external_api.agents.api.v1.AgentConfiguration.data:type_name -> redcarbon.external_api.agents.api.v1.AgentConfigurationData
	2, // 3: redcarbon.external_api.agents.api.v1.AgentConfigurationData.sentinel_one:type_name -> redcarbon.external_api.agents.api.v1.SentinelOneData
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_redcarbon_external_api_agents_api_v1_types_proto_init() }
func file_redcarbon_external_api_agents_api_v1_types_proto_init() {
	if File_redcarbon_external_api_agents_api_v1_types_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_redcarbon_external_api_agents_api_v1_types_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_redcarbon_external_api_agents_api_v1_types_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AgentConfigurationData); i {
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
		file_redcarbon_external_api_agents_api_v1_types_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SentinelOneData); i {
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
	file_redcarbon_external_api_agents_api_v1_types_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*AgentConfigurationData_SentinelOne)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_redcarbon_external_api_agents_api_v1_types_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_redcarbon_external_api_agents_api_v1_types_proto_goTypes,
		DependencyIndexes: file_redcarbon_external_api_agents_api_v1_types_proto_depIdxs,
		MessageInfos:      file_redcarbon_external_api_agents_api_v1_types_proto_msgTypes,
	}.Build()
	File_redcarbon_external_api_agents_api_v1_types_proto = out.File
	file_redcarbon_external_api_agents_api_v1_types_proto_rawDesc = nil
	file_redcarbon_external_api_agents_api_v1_types_proto_goTypes = nil
	file_redcarbon_external_api_agents_api_v1_types_proto_depIdxs = nil
}