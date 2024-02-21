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

package discovery

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"

	"clouditor.io/clouditor/v2/api/assessment/assessmentconnect"
	service_assessment "clouditor.io/clouditor/v2/service/assessment"

	"go.akshayshah.org/memhttp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc/test/bufconn"
)

// TODO(oxisto): Adjust naming to memhttp

const DefaultBufferSize = 1024 * 1024

var (
	bufConnListener *bufconn.Listener
)

func client(srv *memhttp.Server) *http.Client {
	return &http.Client{
		Transport: &http2.Transport{
			DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
				return srv.Transport().DialContext(ctx, network, addr)
			},
			AllowHTTP: true,
		},
	}
}

// startBufConnServer starts an gRPC listening on a bufconn listener. It exposes
// real functionality of the following services for testing purposes:
// * Assessment Service
func startBufConnServer() (*memhttp.Server, *service_assessment.Service) {
	bufConnListener = bufconn.Listen(DefaultBufferSize)

	svc := service_assessment.NewService()

	mux := http.NewServeMux()

	mux.Handle(assessmentconnect.NewAssessmentHandler(svc))
	srv, err := memhttp.New(h2c.NewHandler(mux, &http2.Server{}), memhttp.WithoutTLS())
	if err != nil {
		log.Fatalf("Could not set up memhttp: %v", err)
	}

	return srv, svc
}
