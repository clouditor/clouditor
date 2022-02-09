// Copyright 2016-2022 Fraunhofer AISEC
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

package auth

import (
	"clouditor.io/clouditor/api/auth"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/rest"
	"clouditor.io/clouditor/service"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
)

var (
	grpcPort    int
	authService *Service
	gormX       persistence.IsDatabase
)

func TestMain(m *testing.M) {
	var (
		err    error
		server *grpc.Server
		sock   net.Listener
	)

	// A small embedded DB is needed for the server
	gormX = new(persistence.GormX)
	err = gormX.Init(true, "", 0)
	if err != nil {
		panic(err)
	}
	authService = NewService(gormX)

	// Start at least an authentication server, so that we have something to forward
	sock, server, err = authService.StartDedicatedAuthServer(":0")
	if err != nil {
		panic(err)
	}

	grpcPort = sock.Addr().(*net.TCPAddr).Port

	exit := m.Run()

	sock.Close()
	server.Stop()

	os.Exit(exit)
}

func TestService_ListPublicKeys(t *testing.T) {
	type fields struct {
		apiKey *ecdsa.PrivateKey
	}
	type args struct {
		in0 context.Context
		in1 *auth.ListPublicKeysRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantResponse *auth.ListPublicResponse
		wantErr      bool
	}{
		{
			name: "List single public key",
			fields: fields{
				apiKey: &ecdsa.PrivateKey{
					PublicKey: ecdsa.PublicKey{
						Curve: elliptic.P256(),
						X:     big.NewInt(1),
						Y:     big.NewInt(2),
					},
				},
			},
			args: args{
				in0: context.TODO(),
				in1: &auth.ListPublicKeysRequest{},
			},
			wantResponse: &auth.ListPublicResponse{
				Keys: []*auth.JsonWebKey{
					{
						Kid: "1",
						Kty: "EC",
						Crv: "P-256",
						X:   "AQ",
						Y:   "Ag",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				apiKey: tt.fields.apiKey,
			}
			gotResponse, err := s.ListPublicKeys(tt.args.in0, tt.args.in1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ListPublicKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !proto.Equal(gotResponse, tt.wantResponse) {
				t.Errorf("Service.ListPublicKeys() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestService_recoverFromLoadApiKeyError(t *testing.T) {
	var tmpFile, _ = ioutil.TempFile("", "api.key")
	// Close it immediately , since we want to write to it
	tmpFile.Close()

	defer func() {
		os.Remove(tmpFile.Name())
	}()

	type fields struct {
		config struct {
			keySaveOnCreate bool
			keyPath         string
			keyPassword     string
		}
		apiKey *ecdsa.PrivateKey
	}
	type args struct {
		err         error
		defaultPath bool
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantService assert.ValueAssertionFunc
	}{
		{
			name: "Could not load key from custom path",
			fields: fields{
				config: struct {
					keySaveOnCreate bool
					keyPath         string
					keyPassword     string
				}{
					keySaveOnCreate: false,
					keyPath:         "doesnotexist",
					keyPassword:     "test",
				},
			},
			args: args{
				err:         os.ErrNotExist,
				defaultPath: false,
			},
			wantService: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				// A temporary key should be created
				return assert.NotNil(tt, i1.(*Service).apiKey)
			},
		},
		{
			name: "Could not load key from default path and save it",
			fields: fields{
				config: struct {
					keySaveOnCreate bool
					keyPath         string
					keyPassword     string
				}{
					keySaveOnCreate: true,
					keyPath:         tmpFile.Name(),
					keyPassword:     "test",
				},
			},
			args: args{
				err:         os.ErrNotExist,
				defaultPath: true,
			},
			wantService: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				// A temporary key should be created
				if !assert.NotNil(tt, i1.(*Service).apiKey) {
					return false
				}

				f, err := os.OpenFile(tmpFile.Name(), os.O_RDONLY, 0600)
				if !assert.ErrorIs(tt, err, nil) {
					return false
				}

				// Our tmp file should also contain something now
				data, err := ioutil.ReadAll(f)
				if !assert.ErrorIs(tt, err, nil) {
					return false
				}

				return assert.True(tt, len(data) > 0)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				config: tt.fields.config,
				apiKey: tt.fields.apiKey,
			}
			s.recoverFromLoadApiKeyError(tt.args.err, tt.args.defaultPath)

			if tt.wantService != nil {
				tt.wantService(t, s, tt.args.err, tt.args.defaultPath)
			}
		})
	}
}

func TestService_loadApiKey(t *testing.T) {
	// Prepare a tmp file that contains a new temporary private key
	var tmpFile, _ = ioutil.TempFile("", "api.key")
	tmpFile.Close()

	// Create a new temporary key
	tmpKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	defer func() {
		os.Remove(tmpFile.Name())
	}()

	// Save a key to it
	err := saveApiKey(tmpKey, tmpFile.Name(), "tmp")
	assert.ErrorIs(t, err, nil)

	type args struct {
		path     string
		password []byte
	}
	tests := []struct {
		name    string
		args    args
		wantKey *ecdsa.PrivateKey
		wantErr bool
	}{
		{
			name: "Load existing key",
			args: args{
				path:     tmpFile.Name(),
				password: []byte("tmp"),
			},
			wantKey: tmpKey,
			wantErr: false,
		},
		{
			name: "Load existing key with wrong password",
			args: args{
				path:     tmpFile.Name(),
				password: []byte("notpassword"),
			},
			wantKey: nil,
			wantErr: true,
		},
		{
			name: "Load not existing key",
			args: args{
				path:     "notexists",
				password: []byte("tmp"),
			},
			wantKey: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, err := loadApiKey(tt.args.path, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.loadApiKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotKey, tt.wantKey) {
				t.Errorf("Service.loadApiKey() = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
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

	type configureArgs struct {
		opts []service.AuthOption
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
				opts: []service.AuthOption{service.WithJWKSURL(fmt.Sprintf("http://localhost:%d/.well-known/jwks.json", port))},
			},
			args: args{
				ctx: metadata.NewIncomingContext(context.TODO(), metadata.MD{"Authorization": []string{fmt.Sprintf("bearer %s", loginResponse.AccessToken)}}),
			},
			wantCtx: ValidClaimAssertion,
		},
		{
			name: "Request with invalid bearer token using JWKS",
			configureArgs: configureArgs{
				opts: []service.AuthOption{service.WithJWKSURL(fmt.Sprintf("http://localhost:%d/.well-known/jwks.json", port))},
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
				opts: []service.AuthOption{service.WithJWKSURL(fmt.Sprintf("http://localhost:%d/.well-known/jwks.json", port))},
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
				opts: []service.AuthOption{service.WithPublicKey(authService.GetPublicKey())},
			},
			args: args{
				ctx: metadata.NewIncomingContext(context.TODO(), metadata.MD{"Authorization": []string{fmt.Sprintf("bearer %s", loginResponse.AccessToken)}}),
			},
			wantCtx: ValidClaimAssertion,
		},
		{
			name: "Request without bearer token using a public key",
			configureArgs: configureArgs{
				opts: []service.AuthOption{service.WithPublicKey(authService.GetPublicKey())},
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
			config := service.ConfigureAuth(tt.configureArgs.opts...)
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

func ValidClaimAssertion(tt assert.TestingT, i1 interface{}, _ ...interface{}) bool {
	ctx, ok := i1.(context.Context)
	if !ok {
		tt.Errorf("Return value is not a context")
		return false
	}

	claims, ok := ctx.Value(service.AuthContextKey).(*jwt.RegisteredClaims)
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
