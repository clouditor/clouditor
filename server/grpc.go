// Copyright 2022 Fraunhofer AISEC
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

package server

import (
	"context"
	"fmt"
	"net"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/api/evaluation"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/logging/formatter"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server is a typealias for [grpc.Server] so that users of this package do not need to import the grpc packages
// directly.
type Server = grpc.Server

// StartGRPCServerOption is a type for functional style options that can configure the [StartGRPCServer] function.
type StartGRPCServerOption func(c *config)

// config contains additional server configurations
type config struct {
	grpcOpts   []grpc.ServerOption
	services   map[*grpc.ServiceDesc]any
	ac         AuthConfig
	reflection bool
}

// WithOrchestrator is an option for [StartGRPCServer] to register a [orchestrator.OrchestratorServer] at start.
func WithOrchestrator(svc orchestrator.OrchestratorServer) StartGRPCServerOption {
	return func(c *config) {
		c.services[&orchestrator.Orchestrator_ServiceDesc] = svc
	}
}

// WithAssessment is an option for [StartGRPCServer] to register a [assessment.AssessmentServer] at start.
func WithAssessment(svc assessment.AssessmentServer) StartGRPCServerOption {
	return func(c *config) {
		c.services[&assessment.Assessment_ServiceDesc] = svc
	}
}

// WithEvidenceStore is an option for [StartGRPCServer] to register a [evidence.EvidenceStoreServer] at start.
func WithEvidenceStore(svc evidence.EvidenceStoreServer) StartGRPCServerOption {
	return func(c *config) {
		c.services[&evidence.EvidenceStore_ServiceDesc] = svc
	}
}

// WithDiscovery is an option for [StartGRPCServer] to register a [discovery.DiscoveryServer] at start.
func WithDiscovery(svc discovery.DiscoveryServer) StartGRPCServerOption {
	return func(c *config) {
		c.services[&discovery.Discovery_ServiceDesc] = svc
	}
}

// WithEvaluation is an option for [StartGRPCServer] to register a [evaluation.EvaluationServer] at start.
func WithEvaluation(svc evaluation.EvaluationServer) StartGRPCServerOption {
	return func(c *config) {
		c.services[&evaluation.Evaluation_ServiceDesc] = svc
	}
}

// WithReflection is an option for [StartGRPCServer] to enable gRPC reflection.
func WithReflection() StartGRPCServerOption {
	return func(c *config) {
		c.reflection = true
	}
}

// WithAdditionalGRPCOpts is an option to add an additional gRPC dial options in the REST server communication to the
// backend.
func WithAdditionalGRPCOpts(opts []grpc.ServerOption) StartGRPCServerOption {
	return func(c *config) {
		c.grpcOpts = append(c.grpcOpts, opts...)
	}
}

// StartGRPCServer starts a gRPC server listening on the given address. The server can be configured using the supplied
// opts, e.g., to register various Clouditor services. The server itself is started in a separate Go routine, therefore
// this function will NOT block.
func StartGRPCServer(addr string, opts ...StartGRPCServerOption) (sock net.Listener, srv *Server, err error) {
	// create a new socket for gRPC communication
	sock, err = net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, fmt.Errorf("could not listen: %w", err)
	}

	var c config

	grpcLogger := logrus.New()
	grpcLogger.Formatter = &formatter.GRPCFormatter{TextFormatter: logrus.TextFormatter{ForceColors: true}}
	grpcLoggerEntry := grpcLogger.WithField("component", "grpc")

	c.grpcOpts = []grpc.ServerOption{grpc.ChainUnaryInterceptor(
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_logrus.UnaryServerInterceptor(grpcLoggerEntry),
		UnaryServerInterceptorWithFilter(grpc_auth.UnaryServerInterceptor(c.ac.AuthFunc()), UnaryReflectionFilter),
	),
		grpc.ChainStreamInterceptor(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.StreamServerInterceptor(grpcLoggerEntry),
			StreamServerInterceptorWithFilter(grpc_auth.StreamServerInterceptor(c.ac.AuthFunc()), StreamReflectionFilter),
		),
	}
	c.services = map[*grpc.ServiceDesc]any{}

	for _, o := range opts {
		o(&c)
	}

	srv = grpc.NewServer(
		c.grpcOpts...,
	)

	// Register services
	for sd, svc := range c.services {
		srv.RegisterService(sd, svc)
	}

	// Enable reflection
	if c.reflection {
		reflection.Register(srv)
	}

	go func() {
		// serve the gRPC socket
		_ = srv.Serve(sock)
	}()

	return sock, srv, nil
}

// UnaryServerInterceptorWithFilter wraps a grpc.UnaryServerInterceptor and only invokes the interceptor, if the filter
// function does not return true.
func UnaryServerInterceptorWithFilter(in grpc.UnaryServerInterceptor, filter func(info *grpc.UnaryServerInfo) bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// If the filter evaluates to true, we directly return the handler and ignore the interceptor
		if filter(info) {
			return handler(ctx, req)
		}

		return in(ctx, req, info, handler)
	}
}

// StreamServerInterceptorWithFilter wraps a grpc.StreamServerInterceptor and only invokes the interceptor, if the
// filter function does not return true.
func StreamServerInterceptorWithFilter(in grpc.StreamServerInterceptor, filter func(info *grpc.StreamServerInfo) bool) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// If the filter evaluates to true, we directly return the handler and ignore the interceptor
		if filter(info) {
			return handler(srv, ss)
		}

		return in(srv, ss, info, handler)
	}
}

// UnaryReflectionFilter is a filter that ignores calls to the reflection endpoint
func UnaryReflectionFilter(info *grpc.UnaryServerInfo) bool {
	return info.FullMethod == "/grpc.reflection.v1.ServerReflection/ServerReflectionInfo"
}

// StreamReflectionFilter is a filter that ignores calls to the reflection endpoint
func StreamReflectionFilter(info *grpc.StreamServerInfo) bool {
	return info.FullMethod == "/grpc.reflection.v1.ServerReflection/ServerReflectionInfo"
}
