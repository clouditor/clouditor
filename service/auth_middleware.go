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
	"fmt"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var log *logrus.Entry

func init() {
	log = logrus.WithField("component", "auth-middleware")
}

type AuthConfig struct {
	jwksURL string
	useJWKS bool

	// Jwks contains a JSON Web Key Set, that is used if JWKS support is enabled. Otherwise a
	// stored public key will be used
	Jwks *keyfunc.JWKS

	// publicKey will be used to validate API tokens, if JWKS is not enabled
	publicKey *ecdsa.PublicKey

	AuthFunc grpc_auth.AuthFunc
}

// DefaultJWKSURL is the default JWKS url pointing to a local authentication server.
const DefaultJWKSURL = "http://localhost:8080/.well-known/jwks.json"

// AuthOption is a function-style option type to fine-tune authentication
type AuthOption func(*AuthConfig)

// WithJWKSURL is an option to provide a URL that contains a JSON Web Key Set (JWKS). The JWKS will be used
// to validate tokens coming from RPC clients against public keys contains in the JWKS.
func WithJWKSURL(url string) AuthOption {
	return func(ac *AuthConfig) {
		ac.jwksURL = url
		ac.useJWKS = true
	}
}

// WithPublicKey is an option to directly provide a ECDSA public key which is used to verify tokens coming from RPC clients.
func WithPublicKey(publicKey *ecdsa.PublicKey) AuthOption {
	return func(ac *AuthConfig) {
		ac.publicKey = publicKey
	}
}

// authContextKeyType is a key type that is used in context.WithValue to store the token info in the RPC context.
// It should exclusively be used with the value of AuthContextKey.
//
// Why is this needed? To avoid conflicts, the string type should not be used directly but they should be type-aliased.
type authContextKeyType string

// AuthContextKey is a key used in RPC context to retrieve the token info with using context.Value.
const AuthContextKey = authContextKeyType("token")

// ConfigureAuth creates a new AuthConfig, which can be used in gRPC middleware to provide an authentication layer.
func ConfigureAuth(opts ...AuthOption) *AuthConfig {
	var config = &AuthConfig{
		jwksURL: DefaultJWKSURL,
	}

	// Apply options
	for _, o := range opts {
		o(config)
	}

	config.AuthFunc = func(ctx context.Context) (newCtx context.Context, err error) {
		// Lazy loading of JWKS
		if config.Jwks == nil && config.useJWKS {
			config.Jwks, err = keyfunc.Get(config.jwksURL, keyfunc.Options{
				RefreshInterval: time.Hour,
			})
			if err != nil {
				log.Debugf("Could not retrieve JWKS. API authentication will fail: %v", err)
				return nil, status.Errorf(codes.FailedPrecondition, "could not retrieve JWKS: %v", err)
			}
		}

		token, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			// We do not want to disclose any error details which could be security related,
			// so we do not wrap the original error
			return nil, status.Error(codes.Unauthenticated, "invalid auth token")
		}

		tokenInfo, err := parseToken(token, config)
		if err != nil {
			// We do not want to disclose any error details which could be security related,
			// so we do not wrap the original error
			return nil, status.Errorf(codes.Unauthenticated, "invalid auth token")
		}

		newCtx = context.WithValue(ctx, AuthContextKey, tokenInfo)

		return newCtx, nil
	}

	return config
}

func parseToken(token string, authConfig *AuthConfig) (jwt.Claims, error) {
	var parsedToken *jwt.Token
	var err error

	// Use JWKS, if enabled
	if authConfig.useJWKS {
		parsedToken, err = jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, authConfig.Jwks.Keyfunc)
	} else {
		// Otherwise, we will use the supplied public key
		parsedToken, err = jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
			return authConfig.publicKey, nil
		})
	}

	if err != nil {
		return nil, fmt.Errorf("could not validate JWT: %w", err)
	}

	return parsedToken.Claims, nil
}
