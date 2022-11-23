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

// OutboundConnectorServiceClient is the client API for OutboundConnectorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OutboundConnectorServiceClient interface {
	HandleOutbound(ctx context.Context, in *HandleOutboundRequest, opts ...grpc.CallOption) (*HandleOutboundResponse, error)
}

type outboundConnectorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOutboundConnectorServiceClient(cc grpc.ClientConnInterface) OutboundConnectorServiceClient {
	return &outboundConnectorServiceClient{cc}
}

func (c *outboundConnectorServiceClient) HandleOutbound(ctx context.Context, in *HandleOutboundRequest, opts ...grpc.CallOption) (*HandleOutboundResponse, error) {
	out := new(HandleOutboundResponse)
	err := c.cc.Invoke(ctx, "/sdk.connector.v1.OutboundConnectorService/HandleOutbound", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OutboundConnectorServiceServer is the server API for OutboundConnectorService service.
// All implementations should embed UnimplementedOutboundConnectorServiceServer
// for forward compatibility
type OutboundConnectorServiceServer interface {
	HandleOutbound(context.Context, *HandleOutboundRequest) (*HandleOutboundResponse, error)
}

// UnimplementedOutboundConnectorServiceServer should be embedded to have forward compatible implementations.
type UnimplementedOutboundConnectorServiceServer struct {
}

func (UnimplementedOutboundConnectorServiceServer) HandleOutbound(context.Context, *HandleOutboundRequest) (*HandleOutboundResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleOutbound not implemented")
}

// UnsafeOutboundConnectorServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OutboundConnectorServiceServer will
// result in compilation errors.
type UnsafeOutboundConnectorServiceServer interface {
	mustEmbedUnimplementedOutboundConnectorServiceServer()
}

func RegisterOutboundConnectorServiceServer(s grpc.ServiceRegistrar, srv OutboundConnectorServiceServer) {
	s.RegisterService(&OutboundConnectorService_ServiceDesc, srv)
}

func _OutboundConnectorService_HandleOutbound_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HandleOutboundRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OutboundConnectorServiceServer).HandleOutbound(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sdk.connector.v1.OutboundConnectorService/HandleOutbound",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OutboundConnectorServiceServer).HandleOutbound(ctx, req.(*HandleOutboundRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// OutboundConnectorService_ServiceDesc is the grpc.ServiceDesc for OutboundConnectorService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OutboundConnectorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sdk.connector.v1.OutboundConnectorService",
	HandlerType: (*OutboundConnectorServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HandleOutbound",
			Handler:    _OutboundConnectorService_HandleOutbound_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "azarc/sdk/connector/v1/service.proto",
}

// InboundConnectorServiceClient is the client API for InboundConnectorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type InboundConnectorServiceClient interface {
	Forward(ctx context.Context, in *ForwardRequest, opts ...grpc.CallOption) (*ForwardResponse, error)
}

type inboundConnectorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewInboundConnectorServiceClient(cc grpc.ClientConnInterface) InboundConnectorServiceClient {
	return &inboundConnectorServiceClient{cc}
}

func (c *inboundConnectorServiceClient) Forward(ctx context.Context, in *ForwardRequest, opts ...grpc.CallOption) (*ForwardResponse, error) {
	out := new(ForwardResponse)
	err := c.cc.Invoke(ctx, "/sdk.connector.v1.InboundConnectorService/Forward", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InboundConnectorServiceServer is the server API for InboundConnectorService service.
// All implementations should embed UnimplementedInboundConnectorServiceServer
// for forward compatibility
type InboundConnectorServiceServer interface {
	Forward(context.Context, *ForwardRequest) (*ForwardResponse, error)
}

// UnimplementedInboundConnectorServiceServer should be embedded to have forward compatible implementations.
type UnimplementedInboundConnectorServiceServer struct {
}

func (UnimplementedInboundConnectorServiceServer) Forward(context.Context, *ForwardRequest) (*ForwardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Forward not implemented")
}

// UnsafeInboundConnectorServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to InboundConnectorServiceServer will
// result in compilation errors.
type UnsafeInboundConnectorServiceServer interface {
	mustEmbedUnimplementedInboundConnectorServiceServer()
}

func RegisterInboundConnectorServiceServer(s grpc.ServiceRegistrar, srv InboundConnectorServiceServer) {
	s.RegisterService(&InboundConnectorService_ServiceDesc, srv)
}

func _InboundConnectorService_Forward_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ForwardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InboundConnectorServiceServer).Forward(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sdk.connector.v1.InboundConnectorService/Forward",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InboundConnectorServiceServer).Forward(ctx, req.(*ForwardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// InboundConnectorService_ServiceDesc is the grpc.ServiceDesc for InboundConnectorService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var InboundConnectorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sdk.connector.v1.InboundConnectorService",
	HandlerType: (*InboundConnectorServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Forward",
			Handler:    _InboundConnectorService_Forward_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "azarc/sdk/connector/v1/service.proto",
}