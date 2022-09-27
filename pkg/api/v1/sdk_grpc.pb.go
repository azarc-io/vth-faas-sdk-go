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

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceClient interface {
	GetStageStatus(ctx context.Context, in *GetStageRequest, opts ...grpc.CallOption) (*GetStageResponse, error)
	SetStageStatus(ctx context.Context, in *SetStageRequest, opts ...grpc.CallOption) (*Void, error)
	GetVariable(ctx context.Context, in *GetVariableRequest, opts ...grpc.CallOption) (*GetVariableResponse, error)
	SetVariable(ctx context.Context, in *SetVariableRequest, opts ...grpc.CallOption) (*Void, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) GetStageStatus(ctx context.Context, in *GetStageRequest, opts ...grpc.CallOption) (*GetStageResponse, error) {
	out := new(GetStageResponse)
	err := c.cc.Invoke(ctx, "/sdk_v1.Service/GetStageStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) SetStageStatus(ctx context.Context, in *SetStageRequest, opts ...grpc.CallOption) (*Void, error) {
	out := new(Void)
	err := c.cc.Invoke(ctx, "/sdk_v1.Service/SetStageStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) GetVariable(ctx context.Context, in *GetVariableRequest, opts ...grpc.CallOption) (*GetVariableResponse, error) {
	out := new(GetVariableResponse)
	err := c.cc.Invoke(ctx, "/sdk_v1.Service/GetVariable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) SetVariable(ctx context.Context, in *SetVariableRequest, opts ...grpc.CallOption) (*Void, error) {
	out := new(Void)
	err := c.cc.Invoke(ctx, "/sdk_v1.Service/SetVariable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServer is the server API for Service service.
// All implementations should embed UnimplementedServiceServer
// for forward compatibility
type ServiceServer interface {
	GetStageStatus(context.Context, *GetStageRequest) (*GetStageResponse, error)
	SetStageStatus(context.Context, *SetStageRequest) (*Void, error)
	GetVariable(context.Context, *GetVariableRequest) (*GetVariableResponse, error)
	SetVariable(context.Context, *SetVariableRequest) (*Void, error)
}

// UnimplementedServiceServer should be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (UnimplementedServiceServer) GetStageStatus(context.Context, *GetStageRequest) (*GetStageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStageStatus not implemented")
}
func (UnimplementedServiceServer) SetStageStatus(context.Context, *SetStageRequest) (*Void, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetStageStatus not implemented")
}
func (UnimplementedServiceServer) GetVariable(context.Context, *GetVariableRequest) (*GetVariableResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVariable not implemented")
}
func (UnimplementedServiceServer) SetVariable(context.Context, *SetVariableRequest) (*Void, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetVariable not implemented")
}

// UnsafeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceServer will
// result in compilation errors.
type UnsafeServiceServer interface {
	mustEmbedUnimplementedServiceServer()
}

func RegisterServiceServer(s grpc.ServiceRegistrar, srv ServiceServer) {
	s.RegisterService(&Service_ServiceDesc, srv)
}

func _Service_GetStageStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetStageStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sdk_v1.Service/GetStageStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetStageStatus(ctx, req.(*GetStageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_SetStageStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetStageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).SetStageStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sdk_v1.Service/SetStageStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).SetStageStatus(ctx, req.(*SetStageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_GetVariable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVariableRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetVariable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sdk_v1.Service/GetVariable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetVariable(ctx, req.(*GetVariableRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_SetVariable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetVariableRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).SetVariable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sdk_v1.Service/SetVariable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).SetVariable(ctx, req.(*SetVariableRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Service_ServiceDesc is the grpc.ServiceDesc for Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sdk_v1.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetStageStatus",
			Handler:    _Service_GetStageStatus_Handler,
		},
		{
			MethodName: "SetStageStatus",
			Handler:    _Service_SetStageStatus_Handler,
		},
		{
			MethodName: "GetVariable",
			Handler:    _Service_GetVariable_Handler,
		},
		{
			MethodName: "SetVariable",
			Handler:    _Service_SetVariable_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/api/v1/sdk.proto",
}
