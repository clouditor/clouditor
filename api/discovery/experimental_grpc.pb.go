// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: api/discovery/experimental.proto

package discovery

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

// ExperimentalDiscoveryClient is the client API for ExperimentalDiscovery service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExperimentalDiscoveryClient interface {
	// UpdateResource updates a resource (or creates it, if it does not exist).
	// This is used to give third-party tools the possibility to add something to
	// the resource graph.
	//
	// Note: THIS API IS EXPERIMENTAL AND SUBJECT TO CHANGE
	UpdateResource(ctx context.Context, in *UpdateResourceRequest, opts ...grpc.CallOption) (*Resource, error)
	// ListGraphEdges returns the edges (relationship) between resources in our
	// resource graph.
	//
	// Note: THIS API IS EXPERIMENTAL AND SUBJECT TO CHANGE
	ListGraphEdges(ctx context.Context, in *ListGraphEdgesRequest, opts ...grpc.CallOption) (*ListGraphEdgesResponse, error)
}

type experimentalDiscoveryClient struct {
	cc grpc.ClientConnInterface
}

func NewExperimentalDiscoveryClient(cc grpc.ClientConnInterface) ExperimentalDiscoveryClient {
	return &experimentalDiscoveryClient{cc}
}

func (c *experimentalDiscoveryClient) UpdateResource(ctx context.Context, in *UpdateResourceRequest, opts ...grpc.CallOption) (*Resource, error) {
	out := new(Resource)
	err := c.cc.Invoke(ctx, "/clouditor.discovery.v1experimental.ExperimentalDiscovery/UpdateResource", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *experimentalDiscoveryClient) ListGraphEdges(ctx context.Context, in *ListGraphEdgesRequest, opts ...grpc.CallOption) (*ListGraphEdgesResponse, error) {
	out := new(ListGraphEdgesResponse)
	err := c.cc.Invoke(ctx, "/clouditor.discovery.v1experimental.ExperimentalDiscovery/ListGraphEdges", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExperimentalDiscoveryServer is the server API for ExperimentalDiscovery service.
// All implementations must embed UnimplementedExperimentalDiscoveryServer
// for forward compatibility
type ExperimentalDiscoveryServer interface {
	// UpdateResource updates a resource (or creates it, if it does not exist).
	// This is used to give third-party tools the possibility to add something to
	// the resource graph.
	//
	// Note: THIS API IS EXPERIMENTAL AND SUBJECT TO CHANGE
	UpdateResource(context.Context, *UpdateResourceRequest) (*Resource, error)
	// ListGraphEdges returns the edges (relationship) between resources in our
	// resource graph.
	//
	// Note: THIS API IS EXPERIMENTAL AND SUBJECT TO CHANGE
	ListGraphEdges(context.Context, *ListGraphEdgesRequest) (*ListGraphEdgesResponse, error)
	mustEmbedUnimplementedExperimentalDiscoveryServer()
}

// UnimplementedExperimentalDiscoveryServer must be embedded to have forward compatible implementations.
type UnimplementedExperimentalDiscoveryServer struct {
}

func (UnimplementedExperimentalDiscoveryServer) UpdateResource(context.Context, *UpdateResourceRequest) (*Resource, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateResource not implemented")
}
func (UnimplementedExperimentalDiscoveryServer) ListGraphEdges(context.Context, *ListGraphEdgesRequest) (*ListGraphEdgesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListGraphEdges not implemented")
}
func (UnimplementedExperimentalDiscoveryServer) mustEmbedUnimplementedExperimentalDiscoveryServer() {}

// UnsafeExperimentalDiscoveryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExperimentalDiscoveryServer will
// result in compilation errors.
type UnsafeExperimentalDiscoveryServer interface {
	mustEmbedUnimplementedExperimentalDiscoveryServer()
}

func RegisterExperimentalDiscoveryServer(s grpc.ServiceRegistrar, srv ExperimentalDiscoveryServer) {
	s.RegisterService(&ExperimentalDiscovery_ServiceDesc, srv)
}

func _ExperimentalDiscovery_UpdateResource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateResourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExperimentalDiscoveryServer).UpdateResource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.discovery.v1experimental.ExperimentalDiscovery/UpdateResource",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExperimentalDiscoveryServer).UpdateResource(ctx, req.(*UpdateResourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExperimentalDiscovery_ListGraphEdges_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListGraphEdgesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExperimentalDiscoveryServer).ListGraphEdges(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.discovery.v1experimental.ExperimentalDiscovery/ListGraphEdges",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExperimentalDiscoveryServer).ListGraphEdges(ctx, req.(*ListGraphEdgesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ExperimentalDiscovery_ServiceDesc is the grpc.ServiceDesc for ExperimentalDiscovery service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ExperimentalDiscovery_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "clouditor.discovery.v1experimental.ExperimentalDiscovery",
	HandlerType: (*ExperimentalDiscoveryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateResource",
			Handler:    _ExperimentalDiscovery_UpdateResource_Handler,
		},
		{
			MethodName: "ListGraphEdges",
			Handler:    _ExperimentalDiscovery_ListGraphEdges_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/discovery/experimental.proto",
}