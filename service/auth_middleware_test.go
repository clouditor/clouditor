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
	"fmt"
	"testing"

	"clouditor.io/clouditor/internal/testutil"

	oauth2 "github.com/oxisto/oauth2go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ValidClaimAssertion(tt assert.TestingT, i1 interface{}, _ ...interface{}) bool {
	ctx, ok := i1.(context.Context)
	if !ok {
		tt.Errorf("Return value is not a context")
		return false
	}

	claims, ok := ctx.Value(AuthContextKey).(*OpenIDConnectClaim)
	if !ok {
		tt.Errorf("Token value in context not a JWT claims object")
		return false
	}

	if claims.Subject != testutil.TestAuthClientID {
		tt.Errorf("Subject is not correct")
		return true
	}

	return true
}

func TestAuthConfig_AuthFunc(t *testing.T) {
	var (
		authSrv *oauth2.AuthorizationServer
		err     error
		port    uint16
	)

	// We need to start a REST server for JWKS (using our auth server)
	authSrv, port, err = testutil.StartAuthenticationServer()
	assert.NoError(t, err)
	defer authSrv.Close()

	// Some pre-work to retrieve a valid token
	token, err := authSrv.GenerateToken(testutil.TestAuthClientID, 0, 0)
	assert.NoError(t, err)
	assert.NotNil(t, token)

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
				opts: []AuthOption{WithJWKSURL(testutil.JWKSURL(port))},
			},
			args: args{
				ctx: metadata.NewIncomingContext(context.TODO(), metadata.MD{"Authorization": []string{fmt.Sprintf("bearer %s", token.AccessToken)}}),
			},
			wantCtx: ValidClaimAssertion,
		},
		{
			name: "Request with invalid bearer token using JWKS",
			configureArgs: configureArgs{
				opts: []AuthOption{WithJWKSURL(testutil.JWKSURL(port))},
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
				opts: []AuthOption{WithJWKSURL(testutil.JWKSURL(port))},
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
				opts: []AuthOption{WithPublicKey(authSrv.PublicKeys()[0])},
			},
			args: args{
				ctx: metadata.NewIncomingContext(context.TODO(), metadata.MD{"Authorization": []string{fmt.Sprintf("bearer %s", token.AccessToken)}}),
			},
			wantCtx: ValidClaimAssertion,
		},
		{
			name: "Request without bearer token using a public key",
			configureArgs: configureArgs{
				opts: []AuthOption{WithPublicKey(authSrv.PublicKeys()[0])},
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
