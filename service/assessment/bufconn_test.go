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

package assessment

import (
	"context"
	"net"
	"syscall"

	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/api/orchestrator"
	service_evidence "clouditor.io/clouditor/service/evidence"
	service_orchestrator "clouditor.io/clouditor/service/orchestrator"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const DefaultBufferSize = 1024 * 1024

var (
	bufConnListener *bufconn.Listener
)

func bufConnDialer(context.Context, string) (net.Conn, error) {
	return bufConnListener.Dial()
}

func connectionRefusedDialer(context.Context, string) (net.Conn, error) {
	return nil, syscall.ECONNREFUSED
}

// startBufConnServer starts an gRPC listening on a bufconn listener. It exposes
// real functionality of the following services for testing purposes:
// * Auth
// * Orchestrator
// * Evidence Store
func startBufConnServer() (*grpc.Server, *service_orchestrator.Service, *service_evidence.Service) {
	bufConnListener = bufconn.Listen(DefaultBufferSize)

	server := grpc.NewServer()

	orchestratorService := service_orchestrator.NewService()
	orchestrator.RegisterOrchestratorServer(server, orchestratorService)

	evidenceService := service_evidence.NewService()
	evidence.RegisterEvidenceStoreServer(server, evidenceService)

	go func() {
		if err := server.Serve(bufConnListener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	return server, orchestratorService, evidenceService
}
