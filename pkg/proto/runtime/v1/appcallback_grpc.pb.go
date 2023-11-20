// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.3
// source: dapr/proto/runtime/v1/appcallback.proto

package runtime

import (
	context "context"
	v1 "github.com/liuxd6825/dapr/pkg/proto/common/v1"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// UnsafeAppCallbackServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AppCallbackServer will
// result in compilation errors.
type UnsafeAppCallbackServer interface {
	mustEmbedUnimplementedAppCallbackServer()
}

func RegisterAppCallbackServer(s grpc.ServiceRegistrar, srv AppCallbackServer) {
	s.RegisterService(&AppCallback_ServiceDesc, srv)
}

func _AppCallback_OnInvoke_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.InvokeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppCallbackServer).OnInvoke(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dapr.proto.runtime.v1.AppCallback/OnInvoke",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppCallbackServer).OnInvoke(ctx, req.(*v1.InvokeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AppCallback_ListTopicSubscriptions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppCallbackServer).ListTopicSubscriptions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dapr.proto.runtime.v1.AppCallback/ListTopicSubscriptions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppCallbackServer).ListTopicSubscriptions(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AppCallback_OnTopicEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TopicEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppCallbackServer).OnTopicEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dapr.proto.runtime.v1.AppCallback/OnTopicEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppCallbackServer).OnTopicEvent(ctx, req.(*TopicEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AppCallback_ListInputBindings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppCallbackServer).ListInputBindings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dapr.proto.runtime.v1.AppCallback/ListInputBindings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppCallbackServer).ListInputBindings(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AppCallback_OnBindingEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BindingEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppCallbackServer).OnBindingEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dapr.proto.runtime.v1.AppCallback/OnBindingEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppCallbackServer).OnBindingEvent(ctx, req.(*BindingEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AppCallback_ServiceDesc is the grpc.ServiceDesc for AppCallback service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AppCallback_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dapr.proto.runtime.v1.AppCallback",
	HandlerType: (*AppCallbackServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "OnInvoke",
			Handler:    _AppCallback_OnInvoke_Handler,
		},
		{
			MethodName: "ListTopicSubscriptions",
			Handler:    _AppCallback_ListTopicSubscriptions_Handler,
		},
		{
			MethodName: "OnTopicEvent",
			Handler:    _AppCallback_OnTopicEvent_Handler,
		},
		{
			MethodName: "ListInputBindings",
			Handler:    _AppCallback_ListInputBindings_Handler,
		},
		{
			MethodName: "OnBindingEvent",
			Handler:    _AppCallback_OnBindingEvent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dapr/proto/runtime/v1/appcallback.proto",
}

// AppCallbackHealthCheckClient is the client API for AppCallbackHealthCheck service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AppCallbackHealthCheckClient interface {
	// Health check.
	HealthCheck(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*HealthCheckResponse, error)
}

type appCallbackHealthCheckClient struct {
	cc grpc.ClientConnInterface
}

func NewAppCallbackHealthCheckClient(cc grpc.ClientConnInterface) AppCallbackHealthCheckClient {
	return &appCallbackHealthCheckClient{cc}
}

// AppCallbackHealthCheckServer is the server API for AppCallbackHealthCheck service.
// All implementations must embed UnimplementedAppCallbackHealthCheckServer
// for forward compatibility
type AppCallbackHealthCheckServer interface {
	// Health check.
	HealthCheck(context.Context, *emptypb.Empty) (*HealthCheckResponse, error)
	mustEmbedUnimplementedAppCallbackHealthCheckServer()
}

// UnimplementedAppCallbackHealthCheckServer must be embedded to have forward compatible implementations.
type UnimplementedAppCallbackHealthCheckServer struct {
}

func (UnimplementedAppCallbackHealthCheckServer) HealthCheck(context.Context, *emptypb.Empty) (*HealthCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HealthCheck not implemented")
}
func (UnimplementedAppCallbackHealthCheckServer) mustEmbedUnimplementedAppCallbackHealthCheckServer() {
}

// UnsafeAppCallbackHealthCheckServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AppCallbackHealthCheckServer will
// result in compilation errors.
type UnsafeAppCallbackHealthCheckServer interface {
	mustEmbedUnimplementedAppCallbackHealthCheckServer()
}

func RegisterAppCallbackHealthCheckServer(s grpc.ServiceRegistrar, srv AppCallbackHealthCheckServer) {
	s.RegisterService(&AppCallbackHealthCheck_ServiceDesc, srv)
}

func _AppCallbackHealthCheck_HealthCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppCallbackHealthCheckServer).HealthCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dapr.proto.runtime.v1.AppCallbackHealthCheck/HealthCheck",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppCallbackHealthCheckServer).HealthCheck(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// AppCallbackHealthCheck_ServiceDesc is the grpc.ServiceDesc for AppCallbackHealthCheck service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AppCallbackHealthCheck_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dapr.proto.runtime.v1.AppCallbackHealthCheck",
	HandlerType: (*AppCallbackHealthCheckServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HealthCheck",
			Handler:    _AppCallbackHealthCheck_HealthCheck_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dapr/proto/runtime/v1/appcallback.proto",
}

// AppCallbackAlphaClient is the client API for AppCallbackAlpha service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AppCallbackAlphaClient interface {
	// Subscribes bulk events from Pubsub
	OnBulkTopicEventAlpha1(ctx context.Context, in *TopicEventBulkRequest, opts ...grpc.CallOption) (*TopicEventBulkResponse, error)
}

type appCallbackAlphaClient struct {
	cc grpc.ClientConnInterface
}

func NewAppCallbackAlphaClient(cc grpc.ClientConnInterface) AppCallbackAlphaClient {
	return &appCallbackAlphaClient{cc}
}

// AppCallbackAlphaServer is the server API for AppCallbackAlpha service.
// All implementations must embed UnimplementedAppCallbackAlphaServer
// for forward compatibility
type AppCallbackAlphaServer interface {
	// Subscribes bulk events from Pubsub
	OnBulkTopicEventAlpha1(context.Context, *TopicEventBulkRequest) (*TopicEventBulkResponse, error)
	mustEmbedUnimplementedAppCallbackAlphaServer()
}

// UnimplementedAppCallbackAlphaServer must be embedded to have forward compatible implementations.
type UnimplementedAppCallbackAlphaServer struct {
}

func (UnimplementedAppCallbackAlphaServer) mustEmbedUnimplementedAppCallbackAlphaServer() {}

// UnsafeAppCallbackAlphaServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AppCallbackAlphaServer will
// result in compilation errors.
type UnsafeAppCallbackAlphaServer interface {
	mustEmbedUnimplementedAppCallbackAlphaServer()
}

func RegisterAppCallbackAlphaServer(s grpc.ServiceRegistrar, srv AppCallbackAlphaServer) {
	s.RegisterService(&AppCallbackAlpha_ServiceDesc, srv)
}

func _AppCallbackAlpha_OnBulkTopicEventAlpha1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TopicEventBulkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppCallbackAlphaServer).OnBulkTopicEventAlpha1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dapr.proto.runtime.v1.AppCallbackAlpha/OnBulkTopicEventAlpha1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppCallbackAlphaServer).OnBulkTopicEventAlpha1(ctx, req.(*TopicEventBulkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AppCallbackAlpha_ServiceDesc is the grpc.ServiceDesc for AppCallbackAlpha service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AppCallbackAlpha_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dapr.proto.runtime.v1.AppCallbackAlpha",
	HandlerType: (*AppCallbackAlphaServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "OnBulkTopicEventAlpha1",
			Handler:    _AppCallbackAlpha_OnBulkTopicEventAlpha1_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dapr/proto/runtime/v1/appcallback.proto",
}
