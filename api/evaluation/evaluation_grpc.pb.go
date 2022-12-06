// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
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
	// StartEvaluation evaluates all assessment results of an service based on its Target of Evaluation (binding of a cloud service to a catalog). The assessment results are evaluated regarding the contol ID.
	StartEvaluation(ctx context.Context, in *StartEvaluationRequest, opts ...grpc.CallOption) (*StartEvaluationResponse, error)
	// StopEvaluation stop the evaluation for the given Cloud Service
	StopEvaluation(ctx context.Context, in *StopEvaluationRequest, opts ...grpc.CallOption) (*StopEvaluationResponse, error)
	// List all evaluation results. Part of the public API, also exposed as REST.
	ListEvaluationResults(ctx context.Context, in *ListEvaluationResultsRequest, opts ...grpc.CallOption) (*ListEvaluationResultsResponse, error)
}

type evaluationClient struct {
	cc grpc.ClientConnInterface
}

func NewEvaluationClient(cc grpc.ClientConnInterface) EvaluationClient {
	return &evaluationClient{cc}
}

func (c *evaluationClient) StartEvaluation(ctx context.Context, in *StartEvaluationRequest, opts ...grpc.CallOption) (*StartEvaluationResponse, error) {
	out := new(StartEvaluationResponse)
	err := c.cc.Invoke(ctx, "/clouditor.evaluation.v1.Evaluation/StartEvaluation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *evaluationClient) StopEvaluation(ctx context.Context, in *StopEvaluationRequest, opts ...grpc.CallOption) (*StopEvaluationResponse, error) {
	out := new(StopEvaluationResponse)
	err := c.cc.Invoke(ctx, "/clouditor.evaluation.v1.Evaluation/StopEvaluation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *evaluationClient) ListEvaluationResults(ctx context.Context, in *ListEvaluationResultsRequest, opts ...grpc.CallOption) (*ListEvaluationResultsResponse, error) {
	out := new(ListEvaluationResultsResponse)
	err := c.cc.Invoke(ctx, "/clouditor.evaluation.v1.Evaluation/ListEvaluationResults", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EvaluationServer is the server API for Evaluation service.
// All implementations must embed UnimplementedEvaluationServer
// for forward compatibility
type EvaluationServer interface {
	// StartEvaluation evaluates all assessment results of an service based on its Target of Evaluation (binding of a cloud service to a catalog). The assessment results are evaluated regarding the contol ID.
	StartEvaluation(context.Context, *StartEvaluationRequest) (*StartEvaluationResponse, error)
	// StopEvaluation stop the evaluation for the given Cloud Service
	StopEvaluation(context.Context, *StopEvaluationRequest) (*StopEvaluationResponse, error)
	// List all evaluation results. Part of the public API, also exposed as REST.
	ListEvaluationResults(context.Context, *ListEvaluationResultsRequest) (*ListEvaluationResultsResponse, error)
	mustEmbedUnimplementedEvaluationServer()
}

// UnimplementedEvaluationServer must be embedded to have forward compatible implementations.
type UnimplementedEvaluationServer struct {
}

func (UnimplementedEvaluationServer) StartEvaluation(context.Context, *StartEvaluationRequest) (*StartEvaluationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartEvaluation not implemented")
}
func (UnimplementedEvaluationServer) StopEvaluation(context.Context, *StopEvaluationRequest) (*StopEvaluationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopEvaluation not implemented")
}
func (UnimplementedEvaluationServer) ListEvaluationResults(context.Context, *ListEvaluationResultsRequest) (*ListEvaluationResultsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEvaluationResults not implemented")
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

func _Evaluation_StartEvaluation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartEvaluationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EvaluationServer).StartEvaluation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.evaluation.v1.Evaluation/StartEvaluation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EvaluationServer).StartEvaluation(ctx, req.(*StartEvaluationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Evaluation_StopEvaluation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopEvaluationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EvaluationServer).StopEvaluation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.evaluation.v1.Evaluation/StopEvaluation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EvaluationServer).StopEvaluation(ctx, req.(*StopEvaluationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Evaluation_ListEvaluationResults_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListEvaluationResultsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EvaluationServer).ListEvaluationResults(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.evaluation.v1.Evaluation/ListEvaluationResults",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EvaluationServer).ListEvaluationResults(ctx, req.(*ListEvaluationResultsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Evaluation_ServiceDesc is the grpc.ServiceDesc for Evaluation service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Evaluation_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "clouditor.evaluation.v1.Evaluation",
	HandlerType: (*EvaluationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StartEvaluation",
			Handler:    _Evaluation_StartEvaluation_Handler,
		},
		{
			MethodName: "StopEvaluation",
			Handler:    _Evaluation_StopEvaluation_Handler,
		},
		{
			MethodName: "ListEvaluationResults",
			Handler:    _Evaluation_ListEvaluationResults_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/evaluation/evaluation.proto",
}
