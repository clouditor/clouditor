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

package server

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"testing"

	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil"
	"clouditor.io/clouditor/v2/internal/testutil/assert"

	oauth2 "github.com/oxisto/oauth2go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ValidClaimAssertion(t *testing.T, ctx context.Context) bool {
	claims, ok := ctx.Value(AuthContextKey).(*OpenIDConnectClaim)
	if !ok {
		t.Errorf("Token value in context not a JWT claims object")
		return false
	}

	if claims.Subject != testdata.MockAuthClientID {
		t.Errorf("Subject is not correct")
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
	token, err := authSrv.GenerateToken(testdata.MockAuthClientID, 0, 0)
	assert.NoError(t, err)
	assert.NotNil(t, token)

	type fields struct {
		jwksURL   string
		useJWKS   bool
		publicKey *ecdsa.PublicKey
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantJWKS bool
		wantCtx  assert.Want[context.Context]
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "Request with valid bearer token using JWKS",
			fields: fields{
				jwksURL: testutil.JWKSURL(port),
				useJWKS: true,
			},
			args: args{
				ctx: metadata.NewIncomingContext(context.TODO(), metadata.MD{"Authorization": []string{fmt.Sprintf("bearer %s", token.AccessToken)}}),
			},
			wantCtx: ValidClaimAssertion,
		},
		{
			name: "Request with invalid bearer token using JWKS",
			fields: fields{
				jwksURL: testutil.JWKSURL(port),
				useJWKS: true,
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
			fields: fields{
				jwksURL: testutil.JWKSURL(port),
				useJWKS: true,
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
			fields: fields{
				publicKey: authSrv.PublicKeys()[0],
			},
			args: args{
				ctx: metadata.NewIncomingContext(context.TODO(), metadata.MD{"Authorization": []string{fmt.Sprintf("bearer %s", token.AccessToken)}}),
			},
			wantCtx: ValidClaimAssertion,
		},
		{
			name: "Request without bearer token using a public key",
			fields: fields{
				publicKey: authSrv.PublicKeys()[0],
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
			config := &AuthConfig{
				jwksURL:   tt.fields.jwksURL,
				useJWKS:   tt.fields.useJWKS,
				publicKey: tt.fields.publicKey,
			}
			got, err := config.AuthFunc()(tt.args.ctx)

			if tt.wantJWKS {
				assert.NotNil(t, config.jwks)
			}

			if tt.wantErr != nil {
				tt.wantErr(t, err, tt.args.ctx)
			}

			assert.Optional(t, tt.wantCtx, got)
		})
	}
}
