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
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"clouditor.io/clouditor/v2/internal/testutil"
	"clouditor.io/clouditor/v2/internal/testutil/clitest"
	"clouditor.io/clouditor/v2/server"
	service_orchestrator "clouditor.io/clouditor/v2/service/orchestrator"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

var (
	origins = []string{"clouditor.io", "localhost"}
	methods = []string{"GET", "POST"}
	headers = DefaultAllowedHeaders

	grpcPort uint16 = 0
)

func TestMain(m *testing.M) {
	var (
		err      error
		srv      *grpc.Server
		sock     net.Listener
		authPort uint16
	)

	clitest.AutoChdir()

	_, authPort, err = testutil.StartAuthenticationServer()
	if err != nil {
		panic(err)
	}

	// Start at least an orchestrator service, so that we have something to forward
	sock, srv, err = server.StartGRPCServer("127.0.0.1:0",
		server.WithJWKS(testutil.JWKSURL(authPort)),
		server.WithOrchestrator(service_orchestrator.NewService()),
	)
	if err != nil {
		panic(err)
	}

	grpcPort = sock.Addr().(*net.TCPAddr).AddrPort().Port()

	exit := m.Run()

	sock.Close()
	srv.Stop()

	os.Exit(exit)
}

func TestREST(t *testing.T) {
	go func() {
		err := RunServer(
			context.Background(),
			grpcPort,
			0,
			WithAllowedOrigins(origins),
			WithAllowedMethods(methods),
			WithAllowedHeaders(headers),
			WithAdditionalHandler("GET", "/test", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
				_, err := w.Write([]byte("just a test"))
				if err != nil {
					panic(err)
				}
			}),
		)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	defer StopServer(context.Background())

	// Wait until server is ready to serve
	select {
	case <-ready:
		break
	case <-time.After(10 * time.Second):
		log.Println("Timeout while waiting for REST API")
	}

	assert.NotNil(t, sock)

	port, err := GetServerPort()

	assert.NoError(t, err)
	assert.NotEqual(t, 0, port)

	type args struct {
		origin      string
		contentType string
		method      string
		url         string
		body        io.Reader
		preflight   bool
	}
	tests := []struct {
		name         string
		args         args
		statusCode   int
		headers      map[string]string
		wantResponse assert.ValueAssertionFunc
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
				url:       "v1/orchestrator/cloud_services",
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
				url:       "v1/orchestrator/cloud_services",
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
				url:       "v1/orchestrator/cloud_services",
				preflight: false,
			},
			statusCode: 401, // because we are not supplying an actual login request
			headers: map[string]string{
				"Access-Control-Allow-Origin": "", // should not leak any origin
			},
		},
		{
			name: "Actual request to additional handler",
			args: args{
				method:    "GET",
				url:       "test",
				preflight: false,
			},
			statusCode: 200,
			wantResponse: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				resp, ok := i1.(*http.Response)
				if !ok {
					return assert.True(tt, ok)
				}

				content, err := io.ReadAll(resp.Body)
				if !assert.ErrorIs(tt, err, nil) {
					return false
				}

				return assert.Equal(tt, []byte("just a test"), content)
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

			req, err := http.NewRequest(method, fmt.Sprintf("http://localhost:%d/%s", port, tt.args.url), tt.args.body)
			assert.NoError(t, err)
			assert.NotNil(t, req)

			req.Header.Add("Origin", tt.args.origin)
			req.Header.Add("Content-Type", tt.args.contentType)

			if tt.args.preflight {
				req.Header.Add("Access-Control-Request-Method", tt.args.method)
			}

			resp, err := client.Do(req)

			assert.Equal(t, tt.statusCode, resp.StatusCode)

			assert.NoError(t, err)
			assert.NotNil(t, resp)

			for key, value := range tt.headers {
				assert.Equal(t, value, resp.Header.Get(key))
			}

			if tt.wantResponse != nil {
				tt.wantResponse(t, resp)
			}
		})
	}
}

func Test_corsConfig_OriginAllowed(t *testing.T) {
	type fields struct {
		allowedOrigins []string
		allowedHeaders []string
		allowedMethods []string
	}
	type args struct {
		origin string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "Allow non-browser origin",
			fields: fields{},
			args: args{
				origin: "", // origin is only explicitly set by a browser
			},
			want: true,
		},
		{
			name: "Allowed origin",
			fields: fields{
				allowedOrigins: []string{"clouditor.io", "localhost"},
			},
			args: args{
				origin: "clouditor.io",
			},
			want: true,
		},
		{
			name: "Disallowed origin",
			fields: fields{
				allowedOrigins: []string{"clouditor.io", "localhost"},
			},
			args: args{
				origin: "clouditor.com",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cors := &corsConfig{
				allowedOrigins: tt.fields.allowedOrigins,
				allowedHeaders: tt.fields.allowedHeaders,
				allowedMethods: tt.fields.allowedMethods,
			}
			if got := cors.OriginAllowed(tt.args.origin); got != tt.want {
				t.Errorf("corsConfig.OriginAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}
