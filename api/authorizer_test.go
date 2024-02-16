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
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil"

	oauth2 "github.com/oxisto/oauth2go"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2/clientcredentials"
)

func Test_oauthAuthorizer_Token(t *testing.T) {
	var (
		srv  *oauth2.AuthorizationServer
		port uint16
		err  error
	)
	srv, port, err = testutil.StartAuthenticationServer()

	defer func() {
		err = srv.Close()
		if err != http.ErrServerClosed {
			assert.NoError(t, err)
		}
	}()

	type fields struct {
		TokenSource oauth2.TokenSource
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
				TokenSource: oauth2.ReuseTokenSource(nil,
					(&clientcredentials.Config{
						ClientID:     testdata.MockAuthClientID,
						ClientSecret: testdata.MockAuthClientSecret,
						TokenURL:     fmt.Sprintf("http://localhost:%d/v1/auth/token", port),
					}).TokenSource(context.Background()),
				),
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
				TokenSource: tt.fields.TokenSource,
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
		ClientID:     testdata.MockAuthClientID,
		ClientSecret: testdata.MockAuthClientSecret,
		TokenURL:     "/v1/auth/token",
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
				TokenSource: oauth2.ReuseTokenSource(nil, config.TokenSource(context.Background())),
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
