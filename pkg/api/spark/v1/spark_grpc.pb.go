// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package spark_v1

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
	err := c.cc.Invoke(ctx, "/sdk.spark.v1.AgentService/ExecuteJob", in, out, opts...)
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
		FullMethod: "/sdk.spark.v1.AgentService/ExecuteJob",
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
	ServiceName: "sdk.spark.v1.AgentService",
	HandlerType: (*AgentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ExecuteJob",
			Handler:    _AgentService_ExecuteJob_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/api/spark/v1/spark.proto",
}

// ManagerServiceClient is the client API for ManagerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ManagerServiceClient interface {
	GetStageStatus(ctx context.Context, in *GetStageStatusRequest, opts ...grpc.CallOption) (*GetStageStatusResponse, error)
	SetStageStatus(ctx context.Context, in *SetStageStatusRequest, opts ...grpc.CallOption) (*SetStageStatusResponse, error)
	GetStageResult(ctx context.Context, in *GetStageResultRequest, opts ...grpc.CallOption) (*GetStageResultResponse, error)
	SetStageResult(ctx context.Context, in *SetStageResultRequest, opts ...grpc.CallOption) (*SetStageResultResponse, error)
	GetVariables(ctx context.Context, in *GetVariablesRequest, opts ...grpc.CallOption) (*GetVariablesResponse, error)
	SetVariables(ctx context.Context, in *SetVariablesRequest, opts ...grpc.CallOption) (*SetVariablesResponse, error)
	SetJobStatus(ctx context.Context, in *SetJobStatusRequest, opts ...grpc.CallOption) (*SetJobStatusResponse, error)
	RegisterHeartbeat(ctx context.Context, in *RegisterHeartbeatRequest, opts ...grpc.CallOption) (*RegisterHeartbeatResponse, error)
}

type managerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewManagerServiceClient(cc grpc.ClientConnInterface) ManagerServiceClient {
	return &managerServiceClient{cc}
}

func (c *managerServiceClient) GetStageStatus(ctx context.Context, in *GetStageStatusRequest, opts ...grpc.CallOption) (*GetStageStatusResponse, error) {
	out := new(GetStageStatusResponse)
	err := c.cc.Invoke(ctx, "/sdk.spark.v1.ManagerService/GetStageStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managerServiceClient) SetStageStatus(ctx context.Context, in *SetStageStatusRequest, opts ...grpc.CallOption) (*SetStageStatusResponse, error) {
	out := new(SetStageStatusResponse)
	err := c.cc.Invoke(ctx, "/sdk.spark.v1.ManagerService/SetStageStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managerServiceClient) GetStageResult(ctx context.Context, in *GetStageResultRequest, opts ...grpc.CallOption) (*GetStageResultResponse, error) {
	out := new(GetStageResultResponse)
	err := c.cc.Invoke(ctx, "/sdk.spark.v1.ManagerService/GetStageResult", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managerServiceClient) SetStageResult(ctx context.Context, in *SetStageResultRequest, opts ...grpc.CallOption) (*SetStageResultResponse, error) {
	out := new(SetStageResultResponse)
	err := c.cc.Invoke(ctx, "/sdk.spark.v1.ManagerService/SetStageResult", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managerServiceClient) GetVariables(ctx context.Context, in *GetVariablesRequest, opts ...grpc.CallOption) (*GetVariablesResponse, error) {
	out := new(GetVariablesResponse)
	err := c.cc.Invoke(ctx, "/sdk.spark.v1.ManagerService/GetVariables", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managerServiceClient) SetVariables(ctx context.Context, in *SetVariablesRequest, opts ...grpc.CallOption) (*SetVariablesResponse, error) {
	out := new(SetVariablesResponse)
	err := c.cc.Invoke(ctx, "/sdk.spark.v1.ManagerService/SetVariables", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managerServiceClient) SetJobStatus(ctx context.Context, in *SetJobStatusRequest, opts ...grpc.CallOption) (*SetJobStatusResponse, error) {
	out := new(SetJobStatusResponse)
	err := c.cc.Invoke(ctx, "/sdk.spark.v1.ManagerService/SetJobStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managerServiceClient) RegisterHeartbeat(ctx context.Context, in *RegisterHeartbeatRequest, opts ...grpc.CallOption) (*RegisterHeartbeatResponse, error) {
	out := new(RegisterHeartbeatResponse)
	err := c.cc.Invoke(ctx, "/sdk.spark.v1.ManagerService/RegisterHeartbeat", in, out, opts...)
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
	GetStageResult(context.Context, *GetStageResultRequest) (*GetStageResultResponse, error)
	SetStageResult(context.Context, *SetStageResultRequest) (*SetStageResultResponse, error)
	GetVariables(context.Context, *GetVariablesRequest) (*GetVariablesResponse, error)
	SetVariables(context.Context, *SetVariablesRequest) (*SetVariablesResponse, error)
	SetJobStatus(context.Context, *SetJobStatusRequest) (*SetJobStatusResponse, error)
	RegisterHeartbeat(context.Context, *RegisterHeartbeatRequest) (*RegisterHeartbeatResponse, error)
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
func (UnimplementedManagerServiceServer) GetStageResult(context.Context, *GetStageResultRequest) (*GetStageResultResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStageResult not implemented")
}
func (UnimplementedManagerServiceServer) SetStageResult(context.Context, *SetStageResultRequest) (*SetStageResultResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetStageResult not implemented")
}
func (UnimplementedManagerServiceServer) GetVariables(context.Context, *GetVariablesRequest) (*GetVariablesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVariables not implemented")
}
func (UnimplementedManagerServiceServer) SetVariables(context.Context, *SetVariablesRequest) (*SetVariablesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetVariables not implemented")
}
func (UnimplementedManagerServiceServer) SetJobStatus(context.Context, *SetJobStatusRequest) (*SetJobStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetJobStatus not implemented")
}
func (UnimplementedManagerServiceServer) RegisterHeartbeat(context.Context, *RegisterHeartbeatRequest) (*RegisterHeartbeatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterHeartbeat not implemented")
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
		FullMethod: "/sdk.spark.v1.ManagerService/GetStageStatus",
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
		FullMethod: "/sdk.spark.v1.ManagerService/SetStageStatus",
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
		FullMethod: "/sdk.spark.v1.ManagerService/GetStageResult",
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
		FullMethod: "/sdk.spark.v1.ManagerService/SetStageResult",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagerServiceServer).SetStageResult(ctx, req.(*SetStageResultRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManagerService_GetVariables_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVariablesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagerServiceServer).GetVariables(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sdk.spark.v1.ManagerService/GetVariables",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagerServiceServer).GetVariables(ctx, req.(*GetVariablesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManagerService_SetVariables_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetVariablesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagerServiceServer).SetVariables(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sdk.spark.v1.ManagerService/SetVariables",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagerServiceServer).SetVariables(ctx, req.(*SetVariablesRequest))
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
		FullMethod: "/sdk.spark.v1.ManagerService/SetJobStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagerServiceServer).SetJobStatus(ctx, req.(*SetJobStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManagerService_RegisterHeartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterHeartbeatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagerServiceServer).RegisterHeartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sdk.spark.v1.ManagerService/RegisterHeartbeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagerServiceServer).RegisterHeartbeat(ctx, req.(*RegisterHeartbeatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ManagerService_ServiceDesc is the grpc.ServiceDesc for ManagerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ManagerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sdk.spark.v1.ManagerService",
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
			MethodName: "GetVariables",
			Handler:    _ManagerService_GetVariables_Handler,
		},
		{
			MethodName: "SetVariables",
			Handler:    _ManagerService_SetVariables_Handler,
		},
		{
			MethodName: "SetJobStatus",
			Handler:    _ManagerService_SetJobStatus_Handler,
		},
		{
			MethodName: "RegisterHeartbeat",
			Handler:    _ManagerService_RegisterHeartbeat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/api/spark/v1/spark.proto",
}
