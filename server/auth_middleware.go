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
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var log *logrus.Entry

func init() {
	log = logrus.WithField("component", "auth-middleware")
}

// ProfileClaim represents claims that are contained in the profile scope of OpenID Connect.
type ProfileClaim struct {
	PreferredUsername string `json:"preferred_username"`
	Name              string `json:"name"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
}

// OpenIDConnectClaim represents a claim that supports some aspects of a token issued by an OpenID Connect provider. It
// contains the regular registered JWT claims as well as some specific optional claims, which are empty if Open ID
// Connect is not used.
type OpenIDConnectClaim struct {
	*jwt.RegisteredClaims
	*ProfileClaim
}

// AuthConfig contains all necessary parameters that are needed to configure an authentication middleware.
type AuthConfig struct {
	jwksURL string
	useJWKS bool

	// jwks contains a JSON Web Key Set, that is used if JWKS support is enabled. Otherwise a
	// stored public key will be used
	jwks *keyfunc.JWKS

	// publicKey will be used to validate API tokens, if JWKS is not enabled
	publicKey *ecdsa.PublicKey
}

// DefaultJWKSURL is the default JWKS url pointing to a local authentication server.
const DefaultJWKSURL = "http://localhost:8080/v1/auth/certs"

// WithJWKS is an option to provide a URL that contains a JSON Web Key Set (JWKS). The JWKS will be used
// to validate tokens coming from RPC clients against public keys contains in the JWKS.
func WithJWKS(url string) StartGRPCServerOption {
	return func(c *config) {
		c.ac.jwksURL = url
		c.ac.useJWKS = true
	}
}

// WithPublicKey is an option to directly provide a ECDSA public key which is used to verify tokens coming from RPC clients.
func WithPublicKey(publicKey *ecdsa.PublicKey) StartGRPCServerOption {
	return func(c *config) {
		c.ac.publicKey = publicKey
	}
}

// AuthContextKeyType is a key type that is used in context.WithValue to store the token info in the RPC context.
// It should exclusively be used with the value of AuthContextKey.
//
// Why is this needed? To avoid conflicts, the string type should not be used directly but they should be type-aliased.
type AuthContextKeyType string

// AuthContextKey is a key used in RPC context to retrieve the token info with using context.Value.
const AuthContextKey = AuthContextKeyType("token")

// AuthFunc returns a [grpc_auth.AuthFunc] that authenticates incoming gRPC requests based on the configuration
func (config *AuthConfig) AuthFunc() grpc_auth.AuthFunc {
	return func(ctx context.Context) (newCtx context.Context, err error) {
		// Lazy loading of JWKS
		if config.jwks == nil && config.useJWKS {
			log.Debugf("Trying to retrieve JWKS from %s", config.jwksURL)
			config.jwks, err = keyfunc.Get(config.jwksURL, keyfunc.Options{
				RefreshInterval: time.Hour,
			})
			if err != nil {
				log.Debugf("Could not retrieve JWKS. API authentication will fail: %v", err)
				return nil, status.Errorf(codes.FailedPrecondition, "could not retrieve JWKS: %v", err)
			}
		}

		token, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			log.Debugf("Could not retrieve bearer token from header metadata: %v", err)

			// We do not want to disclose any error details which could be security related,
			// so we do not wrap the original error
			return nil, status.Error(codes.Unauthenticated, "invalid auth token")
		}

		tokenInfo, err := parseToken(token, config)
		if err != nil {
			log.Debugf("Could not parse token in request: %v", err)

			// We do not want to disclose any error details which could be security related,
			// so we do not wrap the original error
			return nil, status.Errorf(codes.Unauthenticated, "invalid auth token")
		}

		newCtx = context.WithValue(ctx, AuthContextKey, tokenInfo)

		return newCtx, nil
	}
}

func parseToken(token string, authConfig *AuthConfig) (jwt.Claims, error) {
	var parsedToken *jwt.Token
	var err error

	// Use JWKS, if enabled
	if authConfig.useJWKS {
		parsedToken, err = jwt.ParseWithClaims(token, &OpenIDConnectClaim{}, authConfig.jwks.Keyfunc)
	} else {
		// Otherwise, we will use the supplied public key
		parsedToken, err = jwt.ParseWithClaims(token, &OpenIDConnectClaim{}, func(t *jwt.Token) (interface{}, error) {
			return authConfig.publicKey, nil
		})
	}

	if err != nil {
		return nil, fmt.Errorf("could not validate JWT: %w", err)
	}

	return parsedToken.Claims, nil
}
