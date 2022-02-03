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

package rest

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"clouditor.io/clouditor/persistence"
	service_auth "clouditor.io/clouditor/service/auth"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

var (
	origins = []string{"clouditor.io", "localhost"}
	methods = []string{"GET", "POST"}
	headers = DefaultAllowedHeaders

	grpcPort int = 0
)

func TestMain(m *testing.M) {
	var (
		err    error
		server *grpc.Server
		sock   net.Listener
	)

	// A small embedded DB is needed for the server
	err = persistence.InitDB(true, "", 0)
	if err != nil {
		panic(err)
	}

	// Start at least an authentication server, so that we have something to forward
	sock, server, err = service_auth.StartDedicatedAuthServer(":0")
	if err != nil {
		panic(err)
	}

	grpcPort = sock.Addr().(*net.TCPAddr).Port

	exit := m.Run()

	sock.Close()
	server.Stop()

	os.Exit(exit)
}

func TestCORS(t *testing.T) {
	go func() {
		err := RunServer(
			context.Background(),
			grpcPort,
			0,
			WithAllowedOrigins(origins),
			WithAllowedMethods(methods),
			WithAllowedHeaders(headers),
		)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	defer StopServer(context.Background())

	// wait until server is ready to serve
	select {
	case <-ready:
		break
	case <-time.After(10 * time.Second):
		log.Println("Timeout while waiting for REST API")
	}

	assert.NotNil(t, sock)

	port, err := GetServerPort()

	assert.Nil(t, err)
	assert.NotEqual(t, 0, port)

	type args struct {
		origin    string
		method    string
		preflight bool
	}

	tests := []struct {
		name       string
		args       args
		statusCode int
		headers    map[string]string
	}{
		{
			name: "Preflight request from valid origin",
			args: args{
				origin:    "clouditor.io",
				method:    "POST",
				preflight: true,
			},
			statusCode: 200,
			headers: map[string]string{
				"Access-Control-Allow-Origin":  "clouditor.io",
				"Access-Control-Allow-Headers": strings.Join(headers, ","),
				"Access-Control-Allow-Methods": strings.Join(methods, ","),
			},
		},
		{
			name: "Actual request from valid origin",
			args: args{
				origin:    "clouditor.io",
				method:    "POST",
				preflight: false,
			},
			statusCode: 401, // because we are not supplying an actual login request
			headers: map[string]string{
				"Access-Control-Allow-Origin":  "clouditor.io",
				"Access-Control-Allow-Headers": "", // should only be part of preflight, not the actual request
				"Access-Control-Allow-Methods": "", // should only be part of preflight, not the actual request
			},
		},
		{
			name: "Preflight request from valid origin",
			args: args{
				origin:    "clouditor.com",
				method:    "POST",
				preflight: true,
			},
			statusCode: 501,
			headers: map[string]string{
				"Access-Control-Allow-Origin": "", // should not leak any origin
			},
		},
		{
			name: "Actual request from invalid origin",
			args: args{
				origin:    "clouditor.com",
				method:    "POST",
				preflight: false,
			},
			statusCode: 401, // because we are not supplying an actual login request
			headers: map[string]string{
				"Access-Control-Allow-Origin": "", // should not leak any origin
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &http.Client{}

			var method string
			if tt.args.preflight {
				method = "OPTIONS"
			} else {
				method = tt.args.method
			}

			req, err := http.NewRequest(method, fmt.Sprintf("http://localhost:%d/v1/auth/login", port), nil)
			assert.Nil(t, err)
			assert.NotNil(t, req)

			req.Header.Add("Origin", tt.args.origin)

			if tt.args.preflight {
				req.Header.Add("Access-Control-Request-Method", tt.args.method)
			}

			resp, err := client.Do(req)

			assert.Equal(t, tt.statusCode, resp.StatusCode)

			assert.Nil(t, err)
			assert.NotNil(t, resp)

			for key, value := range tt.headers {
				assert.Equal(t, value, resp.Header.Get(key))
			}
		})
	}

}
