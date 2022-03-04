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

	"github.com/golang-jwt/jwt/v4"
	oauth2 "github.com/oxisto/oauth2go"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2/clientcredentials"
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

func Test_oauthAuthorizer_Token(t *testing.T) {
	// start an embedded oauth server
	srv := oauth2.NewServer(":0", oauth2.WithClient("client", "secret", ""))

	ln, err := net.Listen("tcp", srv.Addr)
	assert.NoError(t, err)

	port := ln.Addr().(*net.TCPAddr).Port

	go func() {
		err = srv.Serve(ln)
		if err != http.ErrServerClosed {
			assert.NoError(t, err)
		}
	}()

	defer func() {
		err = srv.Close()
		if err != http.ErrServerClosed {
			assert.NoError(t, err)
		}
	}()

	type fields struct {
		Config         *clientcredentials.Config
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
				protectedToken: &protectedToken{
					oauth2.ReuseTokenSource(nil,
						(&clientcredentials.Config{
							ClientID:     "client",
							ClientSecret: "secret",
							TokenURL:     fmt.Sprintf("http://localhost:%d/token", port),
						}).TokenSource(context.Background()),
					),
				},
			},
			wantResp: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				token, ok := i1.(*oauth2.Token)
				assert.True(t, ok)
				assert.NotNil(tt, token)

				return assert.NotEmpty(tt, token.AccessToken)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &oauthAuthorizer{
				protectedToken: tt.fields.protectedToken,
			}

			gotResp, err := o.Token()
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

func TestNewOAuthAuthorizerFromClientCredentials(t *testing.T) {
	var config = clientcredentials.Config{
		ClientID:     "client",
		ClientSecret: "secret",
		TokenURL:     "/token",
	}

	type args struct {
		config *clientcredentials.Config
	}
	tests := []struct {
		name string
		args args
		want Authorizer
	}{
		{
			name: "new",
			args: args{
				&config,
			},
			want: &oauthAuthorizer{
				authURL: "/token",
				protectedToken: &protectedToken{
					TokenSource: oauth2.ReuseTokenSource(nil, config.TokenSource(context.Background())),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOAuthAuthorizerFromClientCredentials(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOAuthAuthorizerFromClientCredentials() = %v, want %v", got, tt.want)
			}
		})
	}
}
