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

package api

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"testing"
	"time"

	"clouditor.io/clouditor/api/auth"

	"github.com/golang-jwt/jwt/v4"
	oauth2 "github.com/oxisto/oauth2go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	tmpKey           *ecdsa.PrivateKey
	mockAccessToken  string
	mockRefreshToken string
	mockExpiry       time.Time
)

func init() {
	// Create a new temporary key
	tmpKey, _ = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	// Create a refresh token
	claims := jwt.NewWithClaims(jwt.SigningMethodES256, &jwt.RegisteredClaims{
		Subject: "clouditor",
	})

	mockRefreshToken, _ = claims.SignedString(tmpKey)

	mockExpiry = time.Now().Add(time.Hour * 1).Truncate(time.Second).UTC()
	// Create an access token
	claims = jwt.NewWithClaims(jwt.SigningMethodES256, &jwt.RegisteredClaims{
		Subject:   "clouditor",
		ExpiresAt: jwt.NewNumericDate(mockExpiry),
	})

	mockAccessToken, _ = claims.SignedString(tmpKey)
}

func TestInternalAuthorizer_Token(t *testing.T) {
	type fields struct {
		authURL        string
		grpcOptions    []grpc.DialOption
		username       string
		password       string
		client         auth.AuthenticationClient
		conn           grpc.ClientConnInterface
		protectedToken *protectedToken
	}
	tests := []struct {
		name    string
		fields  fields
		want    *oauth2.Token
		wantErr bool
	}{
		{
			name: "Fetch access token with refresh token",
			fields: fields{
				client: &mockAuthClient{},
				conn:   &mockConn{},
				protectedToken: &protectedToken{
					token: &oauth2.Token{
						RefreshToken: mockRefreshToken,
					},
				},
			},
			want: &oauth2.Token{
				AccessToken: mockAccessToken,
				Expiry:      mockExpiry,
				TokenType:   "Bearer",
			},
		},
		{
			name: "Fetch access token with username",
			fields: fields{
				client:         &mockAuthClient{},
				conn:           &mockConn{},
				username:       "mock",
				password:       "mock",
				protectedToken: &protectedToken{},
			},
			want: &oauth2.Token{
				AccessToken: mockAccessToken,
				Expiry:      mockExpiry,
				TokenType:   "Bearer",
			},
		},
		{
			name: "Token still valid",
			fields: fields{
				client: &mockAuthClient{},
				conn:   &mockConn{},
				protectedToken: &protectedToken{
					token: &oauth2.Token{
						AccessToken: mockAccessToken,
						TokenType:   "Bearer",
						Expiry:      mockExpiry,
					},
				},
			},
			want: &oauth2.Token{
				AccessToken: mockAccessToken,
				Expiry:      mockExpiry,
				TokenType:   "Bearer",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &internalAuthorizer{
				authURL:        tt.fields.authURL,
				grpcOptions:    tt.fields.grpcOptions,
				username:       tt.fields.username,
				password:       tt.fields.password,
				client:         tt.fields.client,
				conn:           tt.fields.conn,
				protectedToken: tt.fields.protectedToken,
			}

			i.fetchFunc = i.fetchToken

			got, err := i.Token()
			if (err != nil) != tt.wantErr {
				t.Errorf("InternalAuthorizer.Token() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InternalAuthorizer.Token() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_oauthAuthorizer_fetchToken(t *testing.T) {
	// start an embedded oauth server
	srv := oauth2.NewServer(":0", oauth2.WithClient("client", "secret", ""))

	ln, err := net.Listen("tcp", srv.Addr)
	assert.NoError(t, err)

	port := ln.Addr().(*net.TCPAddr).Port

	go srv.Serve(ln)
	defer func() {
		err = srv.Close()
		if err != http.ErrServerClosed {
			assert.NoError(t, err)
		}
	}()

	type fields struct {
		tokenURL       string
		clientID       string
		clientSecret   string
		protectedToken *protectedToken
	}
	type args struct {
		refreshToken string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp assert.ValueAssertionFunc
		wantErr  bool
	}{
		{
			name: "fetch token without refresh token",
			fields: fields{
				tokenURL:       fmt.Sprintf("http://localhost:%d/token", port),
				clientID:       "client",
				clientSecret:   "secret",
				protectedToken: &protectedToken{},
			},
			wantResp: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				token, ok := i1.(*auth.TokenResponse)
				assert.True(t, ok)
				assert.NotNil(tt, token)

				return assert.NotEmpty(tt, token.AccessToken)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &oauthAuthorizer{
				tokenURL:       tt.fields.tokenURL,
				clientID:       tt.fields.clientID,
				clientSecret:   tt.fields.clientSecret,
				protectedToken: tt.fields.protectedToken,
			}

			gotResp, err := o.fetchToken(tt.args.refreshToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("oauthAuthorizer.fetchToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantResp != nil {
				tt.wantResp(t, gotResp, tt.args.refreshToken)
			}
		})
	}
}

type mockAuthClient struct{}

func (mockAuthClient) Login(_ context.Context, _ *auth.LoginRequest, _ ...grpc.CallOption) (*auth.TokenResponse, error) {
	return &auth.TokenResponse{
		AccessToken: mockAccessToken,
		TokenType:   "Bearer",
		Expiry:      timestamppb.New(mockExpiry),
	}, nil
}

func (mockAuthClient) Token(_ context.Context, _ *auth.TokenRequest, _ ...grpc.CallOption) (*auth.TokenResponse, error) {
	return &auth.TokenResponse{
		AccessToken: mockAccessToken,
		TokenType:   "Bearer",
		Expiry:      timestamppb.New(mockExpiry),
	}, nil
}

func (mockAuthClient) ListPublicKeys(_ context.Context, _ *auth.ListPublicKeysRequest, _ ...grpc.CallOption) (*auth.ListPublicResponse, error) {
	return nil, nil
}

type mockConn struct{}

func (mockConn) Invoke(_ context.Context, _ string, _ interface{}, _ interface{}, _ ...grpc.CallOption) error {
	return nil
}

func (mockConn) NewStream(_ context.Context, _ *grpc.StreamDesc, _ string, _ ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, nil
}
