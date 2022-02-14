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
	"reflect"
	"testing"
	"time"

	"clouditor.io/clouditor/api/auth"
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

	// Start at least an authentication server, so that we have something to forward
	sock, server, authService, err = StartDedicatedAuthServer(":0")
	if err != nil {
		panic(err)
	}

	grpcPort = sock.Addr().(*net.TCPAddr).Port

	exit := m.Run()

	sock.Close()
	server.Stop()

	os.Exit(exit)
}

func ValidClaimAssertion(tt assert.TestingT, i1 interface{}, _ ...interface{}) bool {
	ctx, ok := i1.(context.Context)
	if !ok {
		tt.Errorf("Return value is not a context")
		return false
	}

	claims, ok := ctx.Value(AuthContextKey).(*jwt.RegisteredClaims)
	if !ok {
		tt.Errorf("Token value in context not a JWT claims object")
		return false
	}

	if claims.Subject != "clouditor" {
		tt.Errorf("Subject is not correct")
		return true
	}

	return true
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
	assert.NoError(t, err)

	// Some pre-work to retrieve a valid token
	loginResponse, err := authService.Login(context.TODO(), &auth.LoginRequest{Username: "clouditor", Password: "clouditor"})
	assert.NoError(t, err)
	assert.NotNil(t, loginResponse)

	type configureArgs struct {
		opts []AuthOption
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name          string
		configureArgs configureArgs
		args          args
		wantJWKS      bool
		wantCtx       assert.ValueAssertionFunc
		wantErr       assert.ErrorAssertionFunc
	}{
		{
			name: "Request with valid bearer token using JWKS",
			configureArgs: configureArgs{
				opts: []AuthOption{WithJWKSURL(fmt.Sprintf("http://localhost:%d/.well-known/jwks.json", port))},
			},
			args: args{
				ctx: metadata.NewIncomingContext(context.TODO(), metadata.MD{"Authorization": []string{fmt.Sprintf("bearer %s", loginResponse.AccessToken)}}),
			},
			wantCtx: ValidClaimAssertion,
		},
		{
			name: "Request with invalid bearer token using JWKS",
			configureArgs: configureArgs{
				opts: []AuthOption{WithJWKSURL(fmt.Sprintf("http://localhost:%d/.well-known/jwks.json", port))},
			},
			args: args{
				ctx: metadata.NewIncomingContext(context.TODO(), metadata.MD{"Authorization": []string{"bearer not_really"}}),
			},
			wantErr: func(tt assert.TestingT, e error, i ...interface{}) bool {
				return assert.ErrorIs(tt, e, status.Error(codes.Unauthenticated, "invalid auth token"))
			},
		},
		{
			name: "Request without bearer token using JWKS",
			configureArgs: configureArgs{
				opts: []AuthOption{WithJWKSURL(fmt.Sprintf("http://localhost:%d/.well-known/jwks.json", port))},
			},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: func(tt assert.TestingT, e error, i ...interface{}) bool {
				return assert.ErrorIs(tt, e, status.Error(codes.Unauthenticated, "invalid auth token"))
			},
		},
		{
			name: "Request with valid bearer token using a public key",
			configureArgs: configureArgs{
				opts: []AuthOption{WithPublicKey(authService.GetPublicKey())},
			},
			args: args{
				ctx: metadata.NewIncomingContext(context.TODO(), metadata.MD{"Authorization": []string{fmt.Sprintf("bearer %s", loginResponse.AccessToken)}}),
			},
			wantCtx: ValidClaimAssertion,
		},
		{
			name: "Request without bearer token using a public key",
			configureArgs: configureArgs{
				opts: []AuthOption{WithPublicKey(authService.GetPublicKey())},
			},
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
			config := ConfigureAuth(tt.configureArgs.opts...)
			got, err := config.AuthFunc(tt.args.ctx)

			if tt.wantJWKS {
				assert.NotNil(t, config.Jwks)
			}

			if tt.wantErr != nil {
				tt.wantErr(t, err, tt.args.ctx)
			}

			if tt.wantCtx != nil {
				tt.wantCtx(t, got, tt.args.ctx)
			}
		})
	}
}

func TestStartDedicatedAuthServer(t *testing.T) {
	type args struct {
		address string
		opts    []service_auth.ServiceOption
	}
	tests := []struct {
		name            string
		args            args
		wantSock        net.Listener
		wantServer      *grpc.Server
		wantAuthService *service_auth.Service
		wantErr         bool
	}{
		{
			name: "Could not create default user",
			args: args{
				opts: []service_auth.ServiceOption{service_auth.WithStorage(mockStorage{
					createError: errors.New("could not create"),
				})},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSock, gotServer, gotAuthService, err := StartDedicatedAuthServer(tt.args.address, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("StartDedicatedAuthServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSock, tt.wantSock) {
				t.Errorf("StartDedicatedAuthServer() gotSock = %v, want %v", gotSock, tt.wantSock)
			}
			if !reflect.DeepEqual(gotServer, tt.wantServer) {
				t.Errorf("StartDedicatedAuthServer() gotServer = %v, want %v", gotServer, tt.wantServer)
			}
			if !reflect.DeepEqual(gotAuthService, tt.wantAuthService) {
				t.Errorf("StartDedicatedAuthServer() gotAuthService = %v, want %v", gotAuthService, tt.wantAuthService)
			}
		})
	}
}

// mockStorage is a mocked persistance.Storage implementation that returns errors at the specified
// operations.
//
// TODO(oxisto): Extract this struct into our new internal/testutils package
type mockStorage struct {
	createError error
	saveError   error
	updateError error
	getError    error
	listError   error
	countError  error
	deleteError error
}

func (m mockStorage) Create(interface{}) error { return m.createError }

func (m mockStorage) Save(interface{}, ...interface{}) error { return m.saveError }

func (m mockStorage) Update(interface{}, interface{}, ...interface{}) error {
	return m.updateError
}

func (m mockStorage) Get(interface{}, ...interface{}) error { return m.getError }

func (m mockStorage) List(interface{}, ...interface{}) error { return m.listError }

func (m mockStorage) Count(interface{}, ...interface{}) (int64, error) {
	return 0, m.countError
}

func (m mockStorage) Delete(interface{}, ...interface{}) error { return m.deleteError }
