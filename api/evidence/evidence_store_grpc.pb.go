// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: api/evidence/evidence_store.proto

package evidence

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

// EvidenceStoreClient is the client API for EvidenceStore service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EvidenceStoreClient interface {
	// Stores an evidence to the evidence storage. Part of the public API, also
	// exposed as REST.
	StoreEvidence(ctx context.Context, in *StoreEvidenceRequest, opts ...grpc.CallOption) (*StoreEvidenceResponse, error)
	// Stores a stream of evidences to the evidence storage and returns a response
	// stream. Part of the public API, not exposed as REST.
	StoreEvidences(ctx context.Context, opts ...grpc.CallOption) (EvidenceStore_StoreEvidencesClient, error)
	// Returns all stored evidences. Part of the public API, also exposed as REST.
	ListEvidences(ctx context.Context, in *ListEvidencesRequest, opts ...grpc.CallOption) (*ListEvidencesResponse, error)
	// Returns a particular stored evidence. Part of the public API, also exposed
	// as REST.
	GetEvidence(ctx context.Context, in *GetEvidenceRequest, opts ...grpc.CallOption) (*Evidence, error)
}

type evidenceStoreClient struct {
	cc grpc.ClientConnInterface
}

func NewEvidenceStoreClient(cc grpc.ClientConnInterface) EvidenceStoreClient {
	return &evidenceStoreClient{cc}
}

func (c *evidenceStoreClient) StoreEvidence(ctx context.Context, in *StoreEvidenceRequest, opts ...grpc.CallOption) (*StoreEvidenceResponse, error) {
	out := new(StoreEvidenceResponse)
	err := c.cc.Invoke(ctx, "/clouditor.evidence.v1.EvidenceStore/StoreEvidence", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *evidenceStoreClient) StoreEvidences(ctx context.Context, opts ...grpc.CallOption) (EvidenceStore_StoreEvidencesClient, error) {
	stream, err := c.cc.NewStream(ctx, &EvidenceStore_ServiceDesc.Streams[0], "/clouditor.evidence.v1.EvidenceStore/StoreEvidences", opts...)
	if err != nil {
		return nil, err
	}
	x := &evidenceStoreStoreEvidencesClient{stream}
	return x, nil
}

type EvidenceStore_StoreEvidencesClient interface {
	Send(*StoreEvidenceRequest) error
	Recv() (*StoreEvidencesResponse, error)
	grpc.ClientStream
}

type evidenceStoreStoreEvidencesClient struct {
	grpc.ClientStream
}

func (x *evidenceStoreStoreEvidencesClient) Send(m *StoreEvidenceRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *evidenceStoreStoreEvidencesClient) Recv() (*StoreEvidencesResponse, error) {
	m := new(StoreEvidencesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *evidenceStoreClient) ListEvidences(ctx context.Context, in *ListEvidencesRequest, opts ...grpc.CallOption) (*ListEvidencesResponse, error) {
	out := new(ListEvidencesResponse)
	err := c.cc.Invoke(ctx, "/clouditor.evidence.v1.EvidenceStore/ListEvidences", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *evidenceStoreClient) GetEvidence(ctx context.Context, in *GetEvidenceRequest, opts ...grpc.CallOption) (*Evidence, error) {
	out := new(Evidence)
	err := c.cc.Invoke(ctx, "/clouditor.evidence.v1.EvidenceStore/GetEvidence", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EvidenceStoreServer is the server API for EvidenceStore service.
// All implementations must embed UnimplementedEvidenceStoreServer
// for forward compatibility
type EvidenceStoreServer interface {
	// Stores an evidence to the evidence storage. Part of the public API, also
	// exposed as REST.
	StoreEvidence(context.Context, *StoreEvidenceRequest) (*StoreEvidenceResponse, error)
	// Stores a stream of evidences to the evidence storage and returns a response
	// stream. Part of the public API, not exposed as REST.
	StoreEvidences(EvidenceStore_StoreEvidencesServer) error
	// Returns all stored evidences. Part of the public API, also exposed as REST.
	ListEvidences(context.Context, *ListEvidencesRequest) (*ListEvidencesResponse, error)
	// Returns a particular stored evidence. Part of the public API, also exposed
	// as REST.
	GetEvidence(context.Context, *GetEvidenceRequest) (*Evidence, error)
	mustEmbedUnimplementedEvidenceStoreServer()
}

// UnimplementedEvidenceStoreServer must be embedded to have forward compatible implementations.
type UnimplementedEvidenceStoreServer struct {
}

func (UnimplementedEvidenceStoreServer) StoreEvidence(context.Context, *StoreEvidenceRequest) (*StoreEvidenceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreEvidence not implemented")
}
func (UnimplementedEvidenceStoreServer) StoreEvidences(EvidenceStore_StoreEvidencesServer) error {
	return status.Errorf(codes.Unimplemented, "method StoreEvidences not implemented")
}
func (UnimplementedEvidenceStoreServer) ListEvidences(context.Context, *ListEvidencesRequest) (*ListEvidencesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEvidences not implemented")
}
func (UnimplementedEvidenceStoreServer) GetEvidence(context.Context, *GetEvidenceRequest) (*Evidence, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEvidence not implemented")
}
func (UnimplementedEvidenceStoreServer) mustEmbedUnimplementedEvidenceStoreServer() {}

// UnsafeEvidenceStoreServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EvidenceStoreServer will
// result in compilation errors.
type UnsafeEvidenceStoreServer interface {
	mustEmbedUnimplementedEvidenceStoreServer()
}

func RegisterEvidenceStoreServer(s grpc.ServiceRegistrar, srv EvidenceStoreServer) {
	s.RegisterService(&EvidenceStore_ServiceDesc, srv)
}

func _EvidenceStore_StoreEvidence_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoreEvidenceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EvidenceStoreServer).StoreEvidence(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.evidence.v1.EvidenceStore/StoreEvidence",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EvidenceStoreServer).StoreEvidence(ctx, req.(*StoreEvidenceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EvidenceStore_StoreEvidences_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(EvidenceStoreServer).StoreEvidences(&evidenceStoreStoreEvidencesServer{stream})
}

type EvidenceStore_StoreEvidencesServer interface {
	Send(*StoreEvidencesResponse) error
	Recv() (*StoreEvidenceRequest, error)
	grpc.ServerStream
}

type evidenceStoreStoreEvidencesServer struct {
	grpc.ServerStream
}

func (x *evidenceStoreStoreEvidencesServer) Send(m *StoreEvidencesResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *evidenceStoreStoreEvidencesServer) Recv() (*StoreEvidenceRequest, error) {
	m := new(StoreEvidenceRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _EvidenceStore_ListEvidences_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListEvidencesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EvidenceStoreServer).ListEvidences(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.evidence.v1.EvidenceStore/ListEvidences",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EvidenceStoreServer).ListEvidences(ctx, req.(*ListEvidencesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EvidenceStore_GetEvidence_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEvidenceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EvidenceStoreServer).GetEvidence(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.evidence.v1.EvidenceStore/GetEvidence",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EvidenceStoreServer).GetEvidence(ctx, req.(*GetEvidenceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EvidenceStore_ServiceDesc is the grpc.ServiceDesc for EvidenceStore service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EvidenceStore_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "clouditor.evidence.v1.EvidenceStore",
	HandlerType: (*EvidenceStoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StoreEvidence",
			Handler:    _EvidenceStore_StoreEvidence_Handler,
		},
		{
			MethodName: "ListEvidences",
			Handler:    _EvidenceStore_ListEvidences_Handler,
		},
		{
			MethodName: "GetEvidence",
			Handler:    _EvidenceStore_GetEvidence_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StoreEvidences",
			Handler:       _EvidenceStore_StoreEvidences_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "api/evidence/evidence_store.proto",
}
