// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package connectorv1

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

// AgentServiceClient is the client API for AgentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AgentServiceClient interface {
	ExecuteJob(ctx context.Context, in *ExecuteJobRequest, opts ...grpc.CallOption) (*ExecuteJobResponse, error)
}

type agentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAgentServiceClient(cc grpc.ClientConnInterface) AgentServiceClient {
	return &agentServiceClient{cc}
}

func (c *agentServiceClient) ExecuteJob(ctx context.Context, in *ExecuteJobRequest, opts ...grpc.CallOption) (*ExecuteJobResponse, error) {
	out := new(ExecuteJobResponse)
	err := c.cc.Invoke(ctx, "/sdk.connector.v1.AgentService/ExecuteJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AgentServiceServer is the server API for AgentService service.
// All implementations should embed UnimplementedAgentServiceServer
// for forward compatibility
type AgentServiceServer interface {
	ExecuteJob(context.Context, *ExecuteJobRequest) (*ExecuteJobResponse, error)
}

// UnimplementedAgentServiceServer should be embedded to have forward compatible implementations.
type UnimplementedAgentServiceServer struct {
}

func (UnimplementedAgentServiceServer) ExecuteJob(context.Context, *ExecuteJobRequest) (*ExecuteJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteJob not implemented")
}

// UnsafeAgentServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AgentServiceServer will
// result in compilation errors.
type UnsafeAgentServiceServer interface {
	mustEmbedUnimplementedAgentServiceServer()
}

func RegisterAgentServiceServer(s grpc.ServiceRegistrar, srv AgentServiceServer) {
	s.RegisterService(&AgentService_ServiceDesc, srv)
}

func _AgentService_ExecuteJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecuteJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServiceServer).ExecuteJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sdk.connector.v1.AgentService/ExecuteJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServiceServer).ExecuteJob(ctx, req.(*ExecuteJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AgentService_ServiceDesc is the grpc.ServiceDesc for AgentService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AgentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sdk.connector.v1.AgentService",
	HandlerType: (*AgentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ExecuteJob",
			Handler:    _AgentService_ExecuteJob_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "azarc/connector/v1/connector.proto",
}

// ManagerServiceClient is the client API for ManagerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ManagerServiceClient interface {
	GetStageStatus(ctx context.Context, in *GetStageStatusRequest, opts ...grpc.CallOption) (*GetStageStatusResponse, error)
	SetStageStatus(ctx context.Context, in *SetStageStatusRequest, opts ...grpc.CallOption) (*SetStageStatusResponse, error)
}

type managerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewManagerServiceClient(cc grpc.ClientConnInterface) ManagerServiceClient {
	return &managerServiceClient{cc}
}

func (c *managerServiceClient) GetStageStatus(ctx context.Context, in *GetStageStatusRequest, opts ...grpc.CallOption) (*GetStageStatusResponse, error) {
	out := new(GetStageStatusResponse)
	err := c.cc.Invoke(ctx, "/sdk.connector.v1.ManagerService/GetStageStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managerServiceClient) SetStageStatus(ctx context.Context, in *SetStageStatusRequest, opts ...grpc.CallOption) (*SetStageStatusResponse, error) {
	out := new(SetStageStatusResponse)
	err := c.cc.Invoke(ctx, "/sdk.connector.v1.ManagerService/SetStageStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ManagerServiceServer is the server API for ManagerService service.
// All implementations should embed UnimplementedManagerServiceServer
// for forward compatibility
type ManagerServiceServer interface {
	GetStageStatus(context.Context, *GetStageStatusRequest) (*GetStageStatusResponse, error)
	SetStageStatus(context.Context, *SetStageStatusRequest) (*SetStageStatusResponse, error)
}

// UnimplementedManagerServiceServer should be embedded to have forward compatible implementations.
type UnimplementedManagerServiceServer struct {
}

func (UnimplementedManagerServiceServer) GetStageStatus(context.Context, *GetStageStatusRequest) (*GetStageStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStageStatus not implemented")
}
func (UnimplementedManagerServiceServer) SetStageStatus(context.Context, *SetStageStatusRequest) (*SetStageStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetStageStatus not implemented")
}

// UnsafeManagerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ManagerServiceServer will
// result in compilation errors.
type UnsafeManagerServiceServer interface {
	mustEmbedUnimplementedManagerServiceServer()
}

func RegisterManagerServiceServer(s grpc.ServiceRegistrar, srv ManagerServiceServer) {
	s.RegisterService(&ManagerService_ServiceDesc, srv)
}

func _ManagerService_GetStageStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStageStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagerServiceServer).GetStageStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sdk.connector.v1.ManagerService/GetStageStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagerServiceServer).GetStageStatus(ctx, req.(*GetStageStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManagerService_SetStageStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetStageStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagerServiceServer).SetStageStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sdk.connector.v1.ManagerService/SetStageStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagerServiceServer).SetStageStatus(ctx, req.(*SetStageStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ManagerService_ServiceDesc is the grpc.ServiceDesc for ManagerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ManagerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sdk.connector.v1.ManagerService",
	HandlerType: (*ManagerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetStageStatus",
			Handler:    _ManagerService_GetStageStatus_Handler,
		},
		{
			MethodName: "SetStageStatus",
			Handler:    _ManagerService_SetStageStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "azarc/connector/v1/connector.proto",
}