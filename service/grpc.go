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

package service

import (
	"fmt"
	"net"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/api/orchestrator"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
)

type StartGRPCServerOption func(srv *grpc.Server)

func WithOrchestrator(svc orchestrator.OrchestratorServer) StartGRPCServerOption {
	return func(srv *grpc.Server) {
		orchestrator.RegisterOrchestratorServer(srv, svc)
	}
}

func WithEvidenceStore(svc evidence.EvidenceStoreServer) StartGRPCServerOption {
	return func(srv *grpc.Server) {
		evidence.RegisterEvidenceStoreServer(srv, svc)
	}
}

func WithDiscovery(svc discovery.DiscoveryServer) StartGRPCServerOption {
	return func(srv *grpc.Server) {
		discovery.RegisterDiscoveryServer(srv, svc)
	}
}

func StartGRPCServer(jwksURL string, opts ...StartGRPCServerOption) (sock net.Listener, srv *grpc.Server, err error) {
	var addr = ":0"

	// create a new socket for gRPC communication
	sock, err = net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, fmt.Errorf("could not listen: %w", err)
	}

	authConfig := ConfigureAuth(WithJWKSURL(jwksURL))

	// We also add our authentication middleware, because we usually add additional service later
	srv = grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_auth.UnaryServerInterceptor(authConfig.AuthFunc),
		),
	)

	for _, o := range opts {
		o(srv)
	}

	go func() {
		// serve the gRPC socket
		_ = srv.Serve(sock)
	}()

	return sock, srv, nil
}
