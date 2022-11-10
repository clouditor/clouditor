// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: api/evaluation/evaluation.proto

package evaluation

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

// EvaluationClient is the client API for Evaluation service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EvaluationClient interface {
	Evaluate(ctx context.Context, in *EvaluateRequest, opts ...grpc.CallOption) (*EvaluateResponse, error)
}

type evaluationClient struct {
	cc grpc.ClientConnInterface
}

func NewEvaluationClient(cc grpc.ClientConnInterface) EvaluationClient {
	return &evaluationClient{cc}
}

func (c *evaluationClient) Evaluate(ctx context.Context, in *EvaluateRequest, opts ...grpc.CallOption) (*EvaluateResponse, error) {
	out := new(EvaluateResponse)
	err := c.cc.Invoke(ctx, "/clouditor.Evaluation/Evaluate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EvaluationServer is the server API for Evaluation service.
// All implementations must embed UnimplementedEvaluationServer
// for forward compatibility
type EvaluationServer interface {
	Evaluate(context.Context, *EvaluateRequest) (*EvaluateResponse, error)
	mustEmbedUnimplementedEvaluationServer()
}

// UnimplementedEvaluationServer must be embedded to have forward compatible implementations.
type UnimplementedEvaluationServer struct {
}

func (UnimplementedEvaluationServer) Evaluate(context.Context, *EvaluateRequest) (*EvaluateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Evaluate not implemented")
}
func (UnimplementedEvaluationServer) mustEmbedUnimplementedEvaluationServer() {}

// UnsafeEvaluationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EvaluationServer will
// result in compilation errors.
type UnsafeEvaluationServer interface {
	mustEmbedUnimplementedEvaluationServer()
}

func RegisterEvaluationServer(s grpc.ServiceRegistrar, srv EvaluationServer) {
	s.RegisterService(&Evaluation_ServiceDesc, srv)
}

func _Evaluation_Evaluate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EvaluateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EvaluationServer).Evaluate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.Evaluation/Evaluate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EvaluationServer).Evaluate(ctx, req.(*EvaluateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Evaluation_ServiceDesc is the grpc.ServiceDesc for Evaluation service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Evaluation_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "clouditor.Evaluation",
	HandlerType: (*EvaluationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Evaluate",
			Handler:    _Evaluation_Evaluate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/evaluation/evaluation.proto",
}
