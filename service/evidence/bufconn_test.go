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

package evidence

import (
	"context"
	"net"

	"clouditor.io/clouditor/v2/api/assessment"
	service_assessment "clouditor.io/clouditor/v2/service/assessment"

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

// startBufConnServer starts an gRPC listening on a bufconn listener. It exposes
// real functionality of the following services for testing purposes:
// * Assessment Service
func startBufConnServer() (*grpc.Server, *service_assessment.Service) {
	bufConnListener = bufconn.Listen(DefaultBufferSize)

	server := grpc.NewServer()

	assessmentService := service_assessment.NewService()
	assessment.RegisterAssessmentServer(server, assessmentService)

	go func() {
		if err := server.Serve(bufConnListener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	return server, assessmentService
}
