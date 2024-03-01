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
	"testing"

	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil"
	"clouditor.io/clouditor/v2/internal/testutil/assert"

	oauth2 "github.com/oxisto/oauth2go"
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
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantToken assert.Want[*oauth2.Token]
		wantErr   assert.WantErr
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
			wantToken: func(t *testing.T, token *oauth2.Token) bool {
				return assert.NotNil(t, token) && assert.NotEmpty(t, token.AccessToken)
			},
			wantErr: assert.Nil[error],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &oauthAuthorizer{
				TokenSource: tt.fields.TokenSource,
			}

			gotToken, err := o.Token()

			tt.wantErr(t, err)
			tt.wantToken(t, gotToken)
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
			got := NewOAuthAuthorizerFromClientCredentials(tt.args.config)
			assert.Equal(t, tt.want, got, assert.CompareAllUnexported())
		})
	}
}
