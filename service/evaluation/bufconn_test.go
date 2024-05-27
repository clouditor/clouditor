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

package evaluation

import (
	"context"
	"net"
	"syscall"

	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/persistence"
	service_orchestrator "clouditor.io/clouditor/v2/service/orchestrator"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const DefaultBufferSize = 1024 * 1024

func newBufConnDialer(storage persistence.Storage) func(context.Context, string) (net.Conn, error) {
	lis, _, _ := startBufConnServer(storage)

	return func(ctx context.Context, s string) (net.Conn, error) {
		return lis.Dial()
	}
}

func connectionRefusedDialer(context.Context, string) (net.Conn, error) {
	return nil, syscall.ECONNREFUSED
}

// startBufConnServer starts an gRPC listening on a bufconn listener. It exposes
// real functionality of the following services for testing purposes:
// * Auth
// * Orchestrator
func startBufConnServer(storage persistence.Storage) (*bufconn.Listener, *grpc.Server, *service_orchestrator.Service) {
	bufConnListener := bufconn.Listen(DefaultBufferSize)

	server := grpc.NewServer()

	orchestratorService := service_orchestrator.NewService(service_orchestrator.WithStorage(storage))
	orchestrator.RegisterOrchestratorServer(server, orchestratorService)

	go func() {
		if err := server.Serve(bufConnListener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	return bufConnListener, server, orchestratorService
}
