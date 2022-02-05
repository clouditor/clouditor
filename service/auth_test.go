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
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"clouditor.io/clouditor/api/auth"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/rest"
	service_auth "clouditor.io/clouditor/service/auth"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	grpcPort    int
	authService *service_auth.Service
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
	sock, server, authService, err = service_auth.StartDedicatedAuthServer(":0")
	if err != nil {
		panic(err)
	}

	grpcPort = sock.Addr().(*net.TCPAddr).Port

	exit := m.Run()

	sock.Close()
	server.Stop()

	os.Exit(exit)
}

func TestAuthConfig_AuthFunc(t *testing.T) {
	// We need to start a REST server for JWKS (using our auth server)
	go func() {
		err := rest.RunServer(
			context.Background(),
			grpcPort,
			0,
		)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	defer rest.StopServer(context.Background())

	// Wait until server is ready to serve
	select {
	case <-rest.GetReadyChannel():
		break
	case <-time.After(10 * time.Second):
		log.Println("Timeout while waiting for REST API")
	}

	port, err := rest.GetServerPort()
	assert.ErrorIs(t, err, nil)

	// Some pre-work to retrieve a valid token
	loginResponse, err := authService.Login(context.TODO(), &auth.LoginRequest{Username: "clouditor", Password: "clouditor"})
	assert.ErrorIs(t, err, nil)
	assert.NotNil(t, loginResponse)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantCtx assert.ValueAssertionFunc
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Request with valid bearer token",
			args: args{
				ctx: metadata.NewIncomingContext(context.TODO(), metadata.MD{"Authorization": []string{fmt.Sprintf("bearer %s", loginResponse.GetToken())}}),
			},
			wantCtx: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				ctx, ok := i1.(context.Context)
				if !ok {
					tt.Errorf("Return value is not a context")
					return false
				}

				claims, ok := ctx.Value(AuthContextKey("token")).(*jwt.RegisteredClaims)
				if !ok {
					tt.Errorf("Token value in context not a JWT claims object")
					return false
				}

				if claims.Subject != "clouditor" {
					tt.Errorf("Subject is not correct")
					return true
				}

				return true
			},
		},
		{
			name: "Request with invalid bearer token",
			args: args{
				ctx: metadata.NewIncomingContext(context.TODO(), metadata.MD{"Authorization": []string{"bearer not_really"}}),
			},
			wantErr: func(tt assert.TestingT, e error, i ...interface{}) bool {
				return assert.ErrorIs(tt, e, status.Error(codes.Unauthenticated, "invalid auth token"))
			},
		},
		{
			name: "Request without bearer token",
			args: args{
				ctx: context.TODO(),
			},
			wantErr: func(tt assert.TestingT, e error, i ...interface{}) bool {
				return assert.ErrorIs(tt, e, status.Error(codes.Unauthenticated, "invalid auth token"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ConfigureAuth(WithJwksUrl(fmt.Sprintf("http://localhost:%d/.well-known/jwks.json", port)))
			got, err := config.AuthFunc(tt.args.ctx)

			if tt.wantErr != nil {
				tt.wantErr(t, err, tt.args.ctx)
			}

			if tt.wantCtx != nil {
				tt.wantCtx(t, got, tt.args.ctx)
			}
		})
	}
}
