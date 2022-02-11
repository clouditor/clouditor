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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"reflect"
	"testing"
	"time"

	"clouditor.io/clouditor/api/auth"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
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
		authURL     string
		grpcOptions []grpc.DialOption
		username    string
		password    string
		client      auth.AuthenticationClient
		conn        grpc.ClientConnInterface
		token       *oauth2.Token
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
				token: &oauth2.Token{
					RefreshToken: mockRefreshToken,
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
				authURL:     tt.fields.authURL,
				grpcOptions: tt.fields.grpcOptions,
				username:    tt.fields.username,
				password:    tt.fields.password,
				client:      tt.fields.client,
				conn:        tt.fields.conn,
				token:       tt.fields.token,
			}
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

type mockAuthClient struct{}

func (m mockAuthClient) Login(_ context.Context, req *auth.LoginRequest, opts ...grpc.CallOption) (*auth.TokenResponse, error) {
	return nil, nil
}

func (m mockAuthClient) Token(_ context.Context, req *auth.TokenRequest, opts ...grpc.CallOption) (*auth.TokenResponse, error) {
	return &auth.TokenResponse{
		AccessToken: mockAccessToken,
		TokenType:   "Bearer",
		Expiry:      timestamppb.New(mockExpiry),
	}, nil
}

func (m mockAuthClient) ListPublicKeys(_ context.Context, req *auth.ListPublicKeysRequest, opts ...grpc.CallOption) (*auth.ListPublicResponse, error) {
	return nil, nil
}

type mockConn struct{}

func (m mockConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	return nil
}

func (m mockConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, nil
}
