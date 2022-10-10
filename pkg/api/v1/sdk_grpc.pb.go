// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: pkg/api/v1/sdk.proto

package sdk_v1

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
	ExecuteJob(ctx context.Context, in *ExecuteJobRequest, opts ...grpc.CallOption) (*Void, error)
}

type agentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAgentServiceClient(cc grpc.ClientConnInterface) AgentServiceClient {
	return &agentServiceClient{cc}
}

func (c *agentServiceClient) ExecuteJob(ctx context.Context, in *ExecuteJobRequest, opts ...grpc.CallOption) (*Void, error) {
	out := new(Void)
	err := c.cc.Invoke(ctx, "/sdk_v1.AgentService/ExecuteJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AgentServiceServer is the server API for AgentService service.
// All implementations should embed UnimplementedAgentServiceServer
// for forward compatibility
type AgentServiceServer interface {
	ExecuteJob(context.Context, *ExecuteJobRequest) (*Void, error)
}

// UnimplementedAgentServiceServer should be embedded to have forward compatible implementations.
type UnimplementedAgentServiceServer struct {
}

func (UnimplementedAgentServiceServer) ExecuteJob(context.Context, *ExecuteJobRequest) (*Void, error) {
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
		FullMethod: "/sdk_v1.AgentService/ExecuteJob",
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
	ServiceName: "sdk_v1.AgentService",
	HandlerType: (*AgentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ExecuteJob",
			Handler:    _AgentService_ExecuteJob_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/api/v1/sdk.proto",
}

// ManagerServiceClient is the client API for ManagerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ManagerServiceClient interface {
	GetStageStatus(ctx context.Context, in *GetStageStatusRequest, opts ...grpc.CallOption) (*GetStageStatusResponse, error)
	SetStageStatus(ctx context.Context, in *SetStageStatusRequest, opts ...grpc.CallOption) (*Void, error)
	GetStageResult(ctx context.Context, in *GetStageResultRequest, opts ...grpc.CallOption) (*GetStageResultResponse, error)
	SetStageResult(ctx context.Context, in *SetStageResultRequest, opts ...grpc.CallOption) (*Void, error)
	GetVariable(ctx context.Context, in *GetVariableRequest, opts ...grpc.CallOption) (*GetVariableResponse, error)
	SetVariable(ctx context.Context, in *SetVariableRequest, opts ...grpc.CallOption) (*Void, error)
	SetJobStatus(ctx context.Context, in *SetJobStatusRequest, opts ...grpc.CallOption) (*Void, error)
	RegisterHeartbeats(ctx context.Context, in *RegisterHeartbeatsRequest, opts ...grpc.CallOption) (*Void, error)
}

type managerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewManagerServiceClient(cc grpc.ClientConnInterface) ManagerServiceClient {
	return &managerServiceClient{cc}
}

func (c *managerServiceClient) GetStageStatus(ctx context.Context, in *GetStageStatusRequest, opts ...grpc.CallOption) (*GetStageStatusResponse, error) {
	out := new(GetStageStatusResponse)
	err := c.cc.Invoke(ctx, "/sdk_v1.ManagerService/GetStageStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managerServiceClient) SetStageStatus(ctx context.Context, in *SetStageStatusRequest, opts ...grpc.CallOption) (*Void, error) {
	out := new(Void)
	err := c.cc.Invoke(ctx, "/sdk_v1.ManagerService/SetStageStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managerServiceClient) GetStageResult(ctx context.Context, in *GetStageResultRequest, opts ...grpc.CallOption) (*GetStageResultResponse, error) {
	out := new(GetStageResultResponse)
	err := c.cc.Invoke(ctx, "/sdk_v1.ManagerService/GetStageResult", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managerServiceClient) SetStageResult(ctx context.Context, in *SetStageResultRequest, opts ...grpc.CallOption) (*Void, error) {
	out := new(Void)
	err := c.cc.Invoke(ctx, "/sdk_v1.ManagerService/SetStageResult", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managerServiceClient) GetVariable(ctx context.Context, in *GetVariableRequest, opts ...grpc.CallOption) (*GetVariableResponse, error) {
	out := new(GetVariableResponse)
	err := c.cc.Invoke(ctx, "/sdk_v1.ManagerService/GetVariable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managerServiceClient) SetVariable(ctx context.Context, in *SetVariableRequest, opts ...grpc.CallOption) (*Void, error) {
	out := new(Void)
	err := c.cc.Invoke(ctx, "/sdk_v1.ManagerService/SetVariable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managerServiceClient) SetJobStatus(ctx context.Context, in *SetJobStatusRequest, opts ...grpc.CallOption) (*Void, error) {
	out := new(Void)
	err := c.cc.Invoke(ctx, "/sdk_v1.ManagerService/SetJobStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managerServiceClient) RegisterHeartbeats(ctx context.Context, in *RegisterHeartbeatsRequest, opts ...grpc.CallOption) (*Void, error) {
	out := new(Void)
	err := c.cc.Invoke(ctx, "/sdk_v1.ManagerService/RegisterHeartbeats", in, out, opts...)
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
	SetStageStatus(context.Context, *SetStageStatusRequest) (*Void, error)
	GetStageResult(context.Context, *GetStageResultRequest) (*GetStageResultResponse, error)
	SetStageResult(context.Context, *SetStageResultRequest) (*Void, error)
	GetVariable(context.Context, *GetVariableRequest) (*GetVariableResponse, error)
	SetVariable(context.Context, *SetVariableRequest) (*Void, error)
	SetJobStatus(context.Context, *SetJobStatusRequest) (*Void, error)
	RegisterHeartbeats(context.Context, *RegisterHeartbeatsRequest) (*Void, error)
}

// UnimplementedManagerServiceServer should be embedded to have forward compatible implementations.
type UnimplementedManagerServiceServer struct {
}

func (UnimplementedManagerServiceServer) GetStageStatus(context.Context, *GetStageStatusRequest) (*GetStageStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStageStatus not implemented")
}
func (UnimplementedManagerServiceServer) SetStageStatus(context.Context, *SetStageStatusRequest) (*Void, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetStageStatus not implemented")
}
func (UnimplementedManagerServiceServer) GetStageResult(context.Context, *GetStageResultRequest) (*GetStageResultResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStageResult not implemented")
}
func (UnimplementedManagerServiceServer) SetStageResult(context.Context, *SetStageResultRequest) (*Void, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetStageResult not implemented")
}
func (UnimplementedManagerServiceServer) GetVariable(context.Context, *GetVariableRequest) (*GetVariableResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVariable not implemented")
}
func (UnimplementedManagerServiceServer) SetVariable(context.Context, *SetVariableRequest) (*Void, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetVariable not implemented")
}
func (UnimplementedManagerServiceServer) SetJobStatus(context.Context, *SetJobStatusRequest) (*Void, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetJobStatus not implemented")
}
func (UnimplementedManagerServiceServer) RegisterHeartbeats(context.Context, *RegisterHeartbeatsRequest) (*Void, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterHeartbeats not implemented")
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
		FullMethod: "/sdk_v1.ManagerService/GetStageStatus",
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
		FullMethod: "/sdk_v1.ManagerService/SetStageStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagerServiceServer).SetStageStatus(ctx, req.(*SetStageStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManagerService_GetStageResult_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStageResultRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagerServiceServer).GetStageResult(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sdk_v1.ManagerService/GetStageResult",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagerServiceServer).GetStageResult(ctx, req.(*GetStageResultRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManagerService_SetStageResult_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetStageResultRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagerServiceServer).SetStageResult(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sdk_v1.ManagerService/SetStageResult",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagerServiceServer).SetStageResult(ctx, req.(*SetStageResultRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManagerService_GetVariable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVariableRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagerServiceServer).GetVariable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sdk_v1.ManagerService/GetVariable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagerServiceServer).GetVariable(ctx, req.(*GetVariableRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManagerService_SetVariable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetVariableRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagerServiceServer).SetVariable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sdk_v1.ManagerService/SetVariable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagerServiceServer).SetVariable(ctx, req.(*SetVariableRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManagerService_SetJobStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetJobStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagerServiceServer).SetJobStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sdk_v1.ManagerService/SetJobStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagerServiceServer).SetJobStatus(ctx, req.(*SetJobStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManagerService_RegisterHeartbeats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterHeartbeatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagerServiceServer).RegisterHeartbeats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sdk_v1.ManagerService/RegisterHeartbeats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagerServiceServer).RegisterHeartbeats(ctx, req.(*RegisterHeartbeatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ManagerService_ServiceDesc is the grpc.ServiceDesc for ManagerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ManagerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sdk_v1.ManagerService",
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
		{
			MethodName: "GetStageResult",
			Handler:    _ManagerService_GetStageResult_Handler,
		},
		{
			MethodName: "SetStageResult",
			Handler:    _ManagerService_SetStageResult_Handler,
		},
		{
			MethodName: "GetVariable",
			Handler:    _ManagerService_GetVariable_Handler,
		},
		{
			MethodName: "SetVariable",
			Handler:    _ManagerService_SetVariable_Handler,
		},
		{
			MethodName: "SetJobStatus",
			Handler:    _ManagerService_SetJobStatus_Handler,
		},
		{
			MethodName: "RegisterHeartbeats",
			Handler:    _ManagerService_RegisterHeartbeats_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/api/v1/sdk.proto",
}
