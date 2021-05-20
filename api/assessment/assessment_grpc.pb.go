// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package assessment

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AssessmentClient is the client API for Assessment service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AssessmentClient interface {
	// Triggers the assement. Part of the private API, not exposed as REST.
	TriggerAssessment(ctx context.Context, in *TriggerAssessmentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Stores the evidences coming from the discovery. Part of the public API,
	// also exposed as REST
	StoreEvidence(ctx context.Context, in *StoreEvidenceRequest, opts ...grpc.CallOption) (*Evidence, error)
}

type assessmentClient struct {
	cc grpc.ClientConnInterface
}

func NewAssessmentClient(cc grpc.ClientConnInterface) AssessmentClient {
	return &assessmentClient{cc}
}

func (c *assessmentClient) TriggerAssessment(ctx context.Context, in *TriggerAssessmentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/clouditor.Assessment/TriggerAssessment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *assessmentClient) StoreEvidence(ctx context.Context, in *StoreEvidenceRequest, opts ...grpc.CallOption) (*Evidence, error) {
	out := new(Evidence)
	err := c.cc.Invoke(ctx, "/clouditor.Assessment/StoreEvidence", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AssessmentServer is the server API for Assessment service.
// All implementations must embed UnimplementedAssessmentServer
// for forward compatibility
type AssessmentServer interface {
	// Triggers the assement. Part of the private API, not exposed as REST.
	TriggerAssessment(context.Context, *TriggerAssessmentRequest) (*emptypb.Empty, error)
	// Stores the evidences coming from the discovery. Part of the public API,
	// also exposed as REST
	StoreEvidence(context.Context, *StoreEvidenceRequest) (*Evidence, error)
	mustEmbedUnimplementedAssessmentServer()
}

// UnimplementedAssessmentServer must be embedded to have forward compatible implementations.
type UnimplementedAssessmentServer struct {
}

func (UnimplementedAssessmentServer) TriggerAssessment(context.Context, *TriggerAssessmentRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TriggerAssessment not implemented")
}
func (UnimplementedAssessmentServer) StoreEvidence(context.Context, *StoreEvidenceRequest) (*Evidence, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreEvidence not implemented")
}
func (UnimplementedAssessmentServer) mustEmbedUnimplementedAssessmentServer() {}

// UnsafeAssessmentServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AssessmentServer will
// result in compilation errors.
type UnsafeAssessmentServer interface {
	mustEmbedUnimplementedAssessmentServer()
}

func RegisterAssessmentServer(s grpc.ServiceRegistrar, srv AssessmentServer) {
	s.RegisterService(&Assessment_ServiceDesc, srv)
}

func _Assessment_TriggerAssessment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TriggerAssessmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AssessmentServer).TriggerAssessment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.Assessment/TriggerAssessment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AssessmentServer).TriggerAssessment(ctx, req.(*TriggerAssessmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Assessment_StoreEvidence_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoreEvidenceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AssessmentServer).StoreEvidence(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.Assessment/StoreEvidence",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AssessmentServer).StoreEvidence(ctx, req.(*StoreEvidenceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Assessment_ServiceDesc is the grpc.ServiceDesc for Assessment service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Assessment_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "clouditor.Assessment",
	HandlerType: (*AssessmentServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TriggerAssessment",
			Handler:    _Assessment_TriggerAssessment_Handler,
		},
		{
			MethodName: "StoreEvidence",
			Handler:    _Assessment_StoreEvidence_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "assessment.proto",
}
