// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: redcarbon/public_apis/agents/api/v1/v1.proto

package agentsPublicApiV1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AgentsPublicApiV1SrvClient is the client API for AgentsPublicApiV1Srv service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AgentsPublicApiV1SrvClient interface {
	HZ(ctx context.Context, in *HZReq, opts ...grpc.CallOption) (*HZRes, error)
	PullConfigurations(ctx context.Context, in *PullConfigurationsReq, opts ...grpc.CallOption) (*PullConfigurationsRes, error)
	SendSentinelOneData(ctx context.Context, in *SendSentinelOneDataReq, opts ...grpc.CallOption) (*SendSentinelOneDataRes, error)
	SendGrayLogImpossibleTravelData(ctx context.Context, in *SendGrayLogImpossibleTravelDataReq, opts ...grpc.CallOption) (*SendGrayLogImpossibleTravelDataRes, error)
	SendGrayLogDatamineQueryResultsData(ctx context.Context, in *SendGrayLogDatamineQueryResultsDataReq, opts ...grpc.CallOption) (*SendGrayLogDatamineQueryResultsDataRes, error)
	SendGrayLogDatamineQueryErrorData(ctx context.Context, in *SendGrayLogDatamineQueryErrorDataReq, opts ...grpc.CallOption) (*SendGrayLogDatamineQueryErrorDataRes, error)
	RefreshToken(ctx context.Context, in *RefreshTokenReq, opts ...grpc.CallOption) (*RefreshTokenRes, error)
	GetGrayLogDataMinePendingQueries(ctx context.Context, in *GetGrayLogDataMinePendingQueriesReq, opts ...grpc.CallOption) (*GetGrayLogDataMinePendingQueriesRes, error)
}

type agentsPublicApiV1SrvClient struct {
	cc grpc.ClientConnInterface
}

func NewAgentsPublicApiV1SrvClient(cc grpc.ClientConnInterface) AgentsPublicApiV1SrvClient {
	return &agentsPublicApiV1SrvClient{cc}
}

func (c *agentsPublicApiV1SrvClient) HZ(ctx context.Context, in *HZReq, opts ...grpc.CallOption) (*HZRes, error) {
	out := new(HZRes)
	err := c.cc.Invoke(ctx, "/redcarbon.public_apis.agents.api.v1.AgentsPublicApiV1Srv/HZ", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentsPublicApiV1SrvClient) PullConfigurations(ctx context.Context, in *PullConfigurationsReq, opts ...grpc.CallOption) (*PullConfigurationsRes, error) {
	out := new(PullConfigurationsRes)
	err := c.cc.Invoke(ctx, "/redcarbon.public_apis.agents.api.v1.AgentsPublicApiV1Srv/PullConfigurations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentsPublicApiV1SrvClient) SendSentinelOneData(ctx context.Context, in *SendSentinelOneDataReq, opts ...grpc.CallOption) (*SendSentinelOneDataRes, error) {
	out := new(SendSentinelOneDataRes)
	err := c.cc.Invoke(ctx, "/redcarbon.public_apis.agents.api.v1.AgentsPublicApiV1Srv/SendSentinelOneData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentsPublicApiV1SrvClient) SendGrayLogImpossibleTravelData(ctx context.Context, in *SendGrayLogImpossibleTravelDataReq, opts ...grpc.CallOption) (*SendGrayLogImpossibleTravelDataRes, error) {
	out := new(SendGrayLogImpossibleTravelDataRes)
	err := c.cc.Invoke(ctx, "/redcarbon.public_apis.agents.api.v1.AgentsPublicApiV1Srv/SendGrayLogImpossibleTravelData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentsPublicApiV1SrvClient) SendGrayLogDatamineQueryResultsData(ctx context.Context, in *SendGrayLogDatamineQueryResultsDataReq, opts ...grpc.CallOption) (*SendGrayLogDatamineQueryResultsDataRes, error) {
	out := new(SendGrayLogDatamineQueryResultsDataRes)
	err := c.cc.Invoke(ctx, "/redcarbon.public_apis.agents.api.v1.AgentsPublicApiV1Srv/SendGrayLogDatamineQueryResultsData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentsPublicApiV1SrvClient) SendGrayLogDatamineQueryErrorData(ctx context.Context, in *SendGrayLogDatamineQueryErrorDataReq, opts ...grpc.CallOption) (*SendGrayLogDatamineQueryErrorDataRes, error) {
	out := new(SendGrayLogDatamineQueryErrorDataRes)
	err := c.cc.Invoke(ctx, "/redcarbon.public_apis.agents.api.v1.AgentsPublicApiV1Srv/SendGrayLogDatamineQueryErrorData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentsPublicApiV1SrvClient) RefreshToken(ctx context.Context, in *RefreshTokenReq, opts ...grpc.CallOption) (*RefreshTokenRes, error) {
	out := new(RefreshTokenRes)
	err := c.cc.Invoke(ctx, "/redcarbon.public_apis.agents.api.v1.AgentsPublicApiV1Srv/RefreshToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentsPublicApiV1SrvClient) GetGrayLogDataMinePendingQueries(ctx context.Context, in *GetGrayLogDataMinePendingQueriesReq, opts ...grpc.CallOption) (*GetGrayLogDataMinePendingQueriesRes, error) {
	out := new(GetGrayLogDataMinePendingQueriesRes)
	err := c.cc.Invoke(ctx, "/redcarbon.public_apis.agents.api.v1.AgentsPublicApiV1Srv/GetGrayLogDataMinePendingQueries", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AgentsPublicApiV1SrvServer is the server API for AgentsPublicApiV1Srv service.
// All implementations should embed UnimplementedAgentsPublicApiV1SrvServer
// for forward compatibility
type AgentsPublicApiV1SrvServer interface {
	HZ(context.Context, *HZReq) (*HZRes, error)
	PullConfigurations(context.Context, *PullConfigurationsReq) (*PullConfigurationsRes, error)
	SendSentinelOneData(context.Context, *SendSentinelOneDataReq) (*SendSentinelOneDataRes, error)
	SendGrayLogImpossibleTravelData(context.Context, *SendGrayLogImpossibleTravelDataReq) (*SendGrayLogImpossibleTravelDataRes, error)
	SendGrayLogDatamineQueryResultsData(context.Context, *SendGrayLogDatamineQueryResultsDataReq) (*SendGrayLogDatamineQueryResultsDataRes, error)
	SendGrayLogDatamineQueryErrorData(context.Context, *SendGrayLogDatamineQueryErrorDataReq) (*SendGrayLogDatamineQueryErrorDataRes, error)
	RefreshToken(context.Context, *RefreshTokenReq) (*RefreshTokenRes, error)
	GetGrayLogDataMinePendingQueries(context.Context, *GetGrayLogDataMinePendingQueriesReq) (*GetGrayLogDataMinePendingQueriesRes, error)
}

// UnimplementedAgentsPublicApiV1SrvServer should be embedded to have forward compatible implementations.
type UnimplementedAgentsPublicApiV1SrvServer struct {
}

func (UnimplementedAgentsPublicApiV1SrvServer) HZ(context.Context, *HZReq) (*HZRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HZ not implemented")
}
func (UnimplementedAgentsPublicApiV1SrvServer) PullConfigurations(context.Context, *PullConfigurationsReq) (*PullConfigurationsRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PullConfigurations not implemented")
}
func (UnimplementedAgentsPublicApiV1SrvServer) SendSentinelOneData(context.Context, *SendSentinelOneDataReq) (*SendSentinelOneDataRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendSentinelOneData not implemented")
}
func (UnimplementedAgentsPublicApiV1SrvServer) SendGrayLogImpossibleTravelData(context.Context, *SendGrayLogImpossibleTravelDataReq) (*SendGrayLogImpossibleTravelDataRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendGrayLogImpossibleTravelData not implemented")
}
func (UnimplementedAgentsPublicApiV1SrvServer) SendGrayLogDatamineQueryResultsData(context.Context, *SendGrayLogDatamineQueryResultsDataReq) (*SendGrayLogDatamineQueryResultsDataRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendGrayLogDatamineQueryResultsData not implemented")
}
func (UnimplementedAgentsPublicApiV1SrvServer) SendGrayLogDatamineQueryErrorData(context.Context, *SendGrayLogDatamineQueryErrorDataReq) (*SendGrayLogDatamineQueryErrorDataRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendGrayLogDatamineQueryErrorData not implemented")
}
func (UnimplementedAgentsPublicApiV1SrvServer) RefreshToken(context.Context, *RefreshTokenReq) (*RefreshTokenRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefreshToken not implemented")
}
func (UnimplementedAgentsPublicApiV1SrvServer) GetGrayLogDataMinePendingQueries(context.Context, *GetGrayLogDataMinePendingQueriesReq) (*GetGrayLogDataMinePendingQueriesRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGrayLogDataMinePendingQueries not implemented")
}

// UnsafeAgentsPublicApiV1SrvServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AgentsPublicApiV1SrvServer will
// result in compilation errors.
type UnsafeAgentsPublicApiV1SrvServer interface {
	mustEmbedUnimplementedAgentsPublicApiV1SrvServer()
}

func RegisterAgentsPublicApiV1SrvServer(s grpc.ServiceRegistrar, srv AgentsPublicApiV1SrvServer) {
	s.RegisterService(&AgentsPublicApiV1Srv_ServiceDesc, srv)
}

func _AgentsPublicApiV1Srv_HZ_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HZReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentsPublicApiV1SrvServer).HZ(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/redcarbon.public_apis.agents.api.v1.AgentsPublicApiV1Srv/HZ",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentsPublicApiV1SrvServer).HZ(ctx, req.(*HZReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _AgentsPublicApiV1Srv_PullConfigurations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PullConfigurationsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentsPublicApiV1SrvServer).PullConfigurations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/redcarbon.public_apis.agents.api.v1.AgentsPublicApiV1Srv/PullConfigurations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentsPublicApiV1SrvServer).PullConfigurations(ctx, req.(*PullConfigurationsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _AgentsPublicApiV1Srv_SendSentinelOneData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendSentinelOneDataReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentsPublicApiV1SrvServer).SendSentinelOneData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/redcarbon.public_apis.agents.api.v1.AgentsPublicApiV1Srv/SendSentinelOneData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentsPublicApiV1SrvServer).SendSentinelOneData(ctx, req.(*SendSentinelOneDataReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _AgentsPublicApiV1Srv_SendGrayLogImpossibleTravelData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendGrayLogImpossibleTravelDataReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentsPublicApiV1SrvServer).SendGrayLogImpossibleTravelData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/redcarbon.public_apis.agents.api.v1.AgentsPublicApiV1Srv/SendGrayLogImpossibleTravelData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentsPublicApiV1SrvServer).SendGrayLogImpossibleTravelData(ctx, req.(*SendGrayLogImpossibleTravelDataReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _AgentsPublicApiV1Srv_SendGrayLogDatamineQueryResultsData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendGrayLogDatamineQueryResultsDataReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentsPublicApiV1SrvServer).SendGrayLogDatamineQueryResultsData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/redcarbon.public_apis.agents.api.v1.AgentsPublicApiV1Srv/SendGrayLogDatamineQueryResultsData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentsPublicApiV1SrvServer).SendGrayLogDatamineQueryResultsData(ctx, req.(*SendGrayLogDatamineQueryResultsDataReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _AgentsPublicApiV1Srv_SendGrayLogDatamineQueryErrorData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendGrayLogDatamineQueryErrorDataReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentsPublicApiV1SrvServer).SendGrayLogDatamineQueryErrorData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/redcarbon.public_apis.agents.api.v1.AgentsPublicApiV1Srv/SendGrayLogDatamineQueryErrorData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentsPublicApiV1SrvServer).SendGrayLogDatamineQueryErrorData(ctx, req.(*SendGrayLogDatamineQueryErrorDataReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _AgentsPublicApiV1Srv_RefreshToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshTokenReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentsPublicApiV1SrvServer).RefreshToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/redcarbon.public_apis.agents.api.v1.AgentsPublicApiV1Srv/RefreshToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentsPublicApiV1SrvServer).RefreshToken(ctx, req.(*RefreshTokenReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _AgentsPublicApiV1Srv_GetGrayLogDataMinePendingQueries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGrayLogDataMinePendingQueriesReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentsPublicApiV1SrvServer).GetGrayLogDataMinePendingQueries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/redcarbon.public_apis.agents.api.v1.AgentsPublicApiV1Srv/GetGrayLogDataMinePendingQueries",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentsPublicApiV1SrvServer).GetGrayLogDataMinePendingQueries(ctx, req.(*GetGrayLogDataMinePendingQueriesReq))
	}
	return interceptor(ctx, in, info, handler)
}

// AgentsPublicApiV1Srv_ServiceDesc is the grpc.ServiceDesc for AgentsPublicApiV1Srv service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AgentsPublicApiV1Srv_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "redcarbon.public_apis.agents.api.v1.AgentsPublicApiV1Srv",
	HandlerType: (*AgentsPublicApiV1SrvServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HZ",
			Handler:    _AgentsPublicApiV1Srv_HZ_Handler,
		},
		{
			MethodName: "PullConfigurations",
			Handler:    _AgentsPublicApiV1Srv_PullConfigurations_Handler,
		},
		{
			MethodName: "SendSentinelOneData",
			Handler:    _AgentsPublicApiV1Srv_SendSentinelOneData_Handler,
		},
		{
			MethodName: "SendGrayLogImpossibleTravelData",
			Handler:    _AgentsPublicApiV1Srv_SendGrayLogImpossibleTravelData_Handler,
		},
		{
			MethodName: "SendGrayLogDatamineQueryResultsData",
			Handler:    _AgentsPublicApiV1Srv_SendGrayLogDatamineQueryResultsData_Handler,
		},
		{
			MethodName: "SendGrayLogDatamineQueryErrorData",
			Handler:    _AgentsPublicApiV1Srv_SendGrayLogDatamineQueryErrorData_Handler,
		},
		{
			MethodName: "RefreshToken",
			Handler:    _AgentsPublicApiV1Srv_RefreshToken_Handler,
		},
		{
			MethodName: "GetGrayLogDataMinePendingQueries",
			Handler:    _AgentsPublicApiV1Srv_GetGrayLogDataMinePendingQueries_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "redcarbon/public_apis/agents/api/v1/v1.proto",
}