// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.20.3
// source: proto/service.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	RateModifier_RateLimiter_FullMethodName = "/RateModifier/RateLimiter"
)

// RateModifierClient is the client API for RateModifier service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RateModifierClient interface {
	RateLimiter(ctx context.Context, in *RateRequest, opts ...grpc.CallOption) (*RateResponse, error)
}

type rateModifierClient struct {
	cc grpc.ClientConnInterface
}

func NewRateModifierClient(cc grpc.ClientConnInterface) RateModifierClient {
	return &rateModifierClient{cc}
}

func (c *rateModifierClient) RateLimiter(ctx context.Context, in *RateRequest, opts ...grpc.CallOption) (*RateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RateResponse)
	err := c.cc.Invoke(ctx, RateModifier_RateLimiter_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RateModifierServer is the server API for RateModifier service.
// All implementations must embed UnimplementedRateModifierServer
// for forward compatibility.
type RateModifierServer interface {
	RateLimiter(context.Context, *RateRequest) (*RateResponse, error)
	mustEmbedUnimplementedRateModifierServer()
}

// UnimplementedRateModifierServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRateModifierServer struct{}

func (UnimplementedRateModifierServer) RateLimiter(context.Context, *RateRequest) (*RateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RateLimiter not implemented")
}
func (UnimplementedRateModifierServer) mustEmbedUnimplementedRateModifierServer() {}
func (UnimplementedRateModifierServer) testEmbeddedByValue()                      {}

// UnsafeRateModifierServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RateModifierServer will
// result in compilation errors.
type UnsafeRateModifierServer interface {
	mustEmbedUnimplementedRateModifierServer()
}

func RegisterRateModifierServer(s grpc.ServiceRegistrar, srv RateModifierServer) {
	// If the following call pancis, it indicates UnimplementedRateModifierServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&RateModifier_ServiceDesc, srv)
}

func _RateModifier_RateLimiter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RateModifierServer).RateLimiter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RateModifier_RateLimiter_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RateModifierServer).RateLimiter(ctx, req.(*RateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RateModifier_ServiceDesc is the grpc.ServiceDesc for RateModifier service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RateModifier_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "RateModifier",
	HandlerType: (*RateModifierServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RateLimiter",
			Handler:    _RateModifier_RateLimiter_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/service.proto",
}
