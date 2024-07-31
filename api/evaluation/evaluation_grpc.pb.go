// Copyright 2023 Fraunhofer AISEC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//           $$\                           $$\ $$\   $$\
//           $$ |                          $$ |\__|  $$ |
//  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
// $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
// $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
// $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
// \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
//  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
//
// This file is part of Clouditor Community Edition.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
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
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	Evaluation_StartEvaluation_FullMethodName        = "/clouditor.evaluation.v1.Evaluation/StartEvaluation"
	Evaluation_StopEvaluation_FullMethodName         = "/clouditor.evaluation.v1.Evaluation/StopEvaluation"
	Evaluation_ListEvaluationResults_FullMethodName  = "/clouditor.evaluation.v1.Evaluation/ListEvaluationResults"
	Evaluation_CreateEvaluationResult_FullMethodName = "/clouditor.evaluation.v1.Evaluation/CreateEvaluationResult"
)

// EvaluationClient is the client API for Evaluation service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Manages the evaluation of Clouditor's assessment results
type EvaluationClient interface {
	// Evaluates periodically all assessment results of a cloud service id based
	// on the given catalog id. Part of the public API, also exposed as REST.
	StartEvaluation(ctx context.Context, in *StartEvaluationRequest, opts ...grpc.CallOption) (*StartEvaluationResponse, error)
	// StopEvaluation stops the evaluation for the given target of evaluation.
	// Part of the public API, also exposed as REST.
	StopEvaluation(ctx context.Context, in *StopEvaluationRequest, opts ...grpc.CallOption) (*StopEvaluationResponse, error)
	// List all evaluation results that the user can access. It can further be
	// restricted by various filtering options. Part of the public API, also
	// exposed as REST.
	ListEvaluationResults(ctx context.Context, in *ListEvaluationResultsRequest, opts ...grpc.CallOption) (*ListEvaluationResultsResponse, error)
	// Creates an evaluation result
	CreateEvaluationResult(ctx context.Context, in *CreateEvaluationResultRequest, opts ...grpc.CallOption) (*EvaluationResult, error)
}

type evaluationClient struct {
	cc grpc.ClientConnInterface
}

func NewEvaluationClient(cc grpc.ClientConnInterface) EvaluationClient {
	return &evaluationClient{cc}
}

func (c *evaluationClient) StartEvaluation(ctx context.Context, in *StartEvaluationRequest, opts ...grpc.CallOption) (*StartEvaluationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StartEvaluationResponse)
	err := c.cc.Invoke(ctx, Evaluation_StartEvaluation_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *evaluationClient) StopEvaluation(ctx context.Context, in *StopEvaluationRequest, opts ...grpc.CallOption) (*StopEvaluationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StopEvaluationResponse)
	err := c.cc.Invoke(ctx, Evaluation_StopEvaluation_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *evaluationClient) ListEvaluationResults(ctx context.Context, in *ListEvaluationResultsRequest, opts ...grpc.CallOption) (*ListEvaluationResultsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListEvaluationResultsResponse)
	err := c.cc.Invoke(ctx, Evaluation_ListEvaluationResults_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *evaluationClient) CreateEvaluationResult(ctx context.Context, in *CreateEvaluationResultRequest, opts ...grpc.CallOption) (*EvaluationResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EvaluationResult)
	err := c.cc.Invoke(ctx, Evaluation_CreateEvaluationResult_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EvaluationServer is the server API for Evaluation service.
// All implementations must embed UnimplementedEvaluationServer
// for forward compatibility
//
// Manages the evaluation of Clouditor's assessment results
type EvaluationServer interface {
	// Evaluates periodically all assessment results of a cloud service id based
	// on the given catalog id. Part of the public API, also exposed as REST.
	StartEvaluation(context.Context, *StartEvaluationRequest) (*StartEvaluationResponse, error)
	// StopEvaluation stops the evaluation for the given target of evaluation.
	// Part of the public API, also exposed as REST.
	StopEvaluation(context.Context, *StopEvaluationRequest) (*StopEvaluationResponse, error)
	// List all evaluation results that the user can access. It can further be
	// restricted by various filtering options. Part of the public API, also
	// exposed as REST.
	ListEvaluationResults(context.Context, *ListEvaluationResultsRequest) (*ListEvaluationResultsResponse, error)
	// Creates an evaluation result
	CreateEvaluationResult(context.Context, *CreateEvaluationResultRequest) (*EvaluationResult, error)
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
func (UnimplementedEvaluationServer) CreateEvaluationResult(context.Context, *CreateEvaluationResultRequest) (*EvaluationResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateEvaluationResult not implemented")
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
		FullMethod: Evaluation_StartEvaluation_FullMethodName,
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
		FullMethod: Evaluation_StopEvaluation_FullMethodName,
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
		FullMethod: Evaluation_ListEvaluationResults_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EvaluationServer).ListEvaluationResults(ctx, req.(*ListEvaluationResultsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Evaluation_CreateEvaluationResult_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateEvaluationResultRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EvaluationServer).CreateEvaluationResult(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Evaluation_CreateEvaluationResult_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EvaluationServer).CreateEvaluationResult(ctx, req.(*CreateEvaluationResultRequest))
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
		{
			MethodName: "CreateEvaluationResult",
			Handler:    _Evaluation_CreateEvaluationResult_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/evaluation/evaluation.proto",
}
