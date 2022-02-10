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
	"clouditor.io/clouditor/persistence"
	"context"
	"net"

	"clouditor.io/clouditor/api/auth"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/api/orchestrator"
	service_auth "clouditor.io/clouditor/service/auth"
	service_evidence "clouditor.io/clouditor/service/evidence"
	service_orchestrator "clouditor.io/clouditor/service/orchestrator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const DefaultBufferSize = 1024 * 1024

var (
	bufConnListener *bufconn.Listener
	db              persistence.IsDatabase
)

func bufConnDialer(context.Context, string) (net.Conn, error) {
	return bufConnListener.Dial()
}

// startBufConnServer starts an gRPC listening on a bufconn listener. It exposes
// real functionality of the following services for testing purposes:
// * Auth
// * Orchestrator
// * Evidence Store
func startBufConnServer() (*grpc.Server, *service_auth.Service, *service_orchestrator.Service, *service_evidence.Service) {
	bufConnListener = bufconn.Listen(DefaultBufferSize)

	server := grpc.NewServer()

	// We do not want a persistent key storage here
	db = new(persistence.GormX)
	err := db.Init(true, "", 0)
	if err != nil {
		panic(err)
	}
	authService := service_auth.NewService(db, service_auth.WithApiKeySaveOnCreate(false))
	auth.RegisterAuthenticationServer(server, authService)

	orchestratorService := service_orchestrator.NewService(nil)
	orchestrator.RegisterOrchestratorServer(server, orchestratorService)

	evidenceService := service_evidence.NewService()
	evidence.RegisterEvidenceStoreServer(server, evidenceService)

	go func() {
		if err := server.Serve(bufConnListener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	return server, authService, orchestratorService, evidenceService
}
