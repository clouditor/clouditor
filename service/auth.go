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
	"net"
	"time"

	"clouditor.io/clouditor/api/auth"
	service_auth "clouditor.io/clouditor/service/auth"
	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var log *logrus.Entry

type AuthConfig struct {
	jwksUrl string
	useJwks bool

	// Jwks contains a JSON Web Key Set, that is used if JWKS support is enabled. Otherwise a
	// stored public key will be used
	Jwks *keyfunc.JWKS

	// publicKey will be used to validate API tokens, if JWKS is not enabled
	publicKey *ecdsa.PublicKey

	AuthFunc grpc_auth.AuthFunc
}

var DefaultJwksUrl = "http://localhost:8080/.well-known/jwks.json"

func init() {
	log = logrus.WithField("component", "service-auth")
}

type AuthOption func(*AuthConfig)

func WithJwksUrl(url string) AuthOption {
	return func(ac *AuthConfig) {
		ac.jwksUrl = url
		ac.useJwks = true
	}
}

func WithPublicKey(publicKey *ecdsa.PublicKey) AuthOption {
	return func(ac *AuthConfig) {
		ac.publicKey = publicKey
	}
}

type AuthContextKey string

func ConfigureAuth(opts ...AuthOption) *AuthConfig {
	var config = &AuthConfig{
		jwksUrl: DefaultJwksUrl,
	}

	// Apply options
	for _, o := range opts {
		o(config)
	}

	config.AuthFunc = func(ctx context.Context) (newCtx context.Context, err error) {
		// Lazy loading of JWKS
		if config.Jwks == nil && config.useJwks {
			config.Jwks, err = keyfunc.Get(config.jwksUrl, keyfunc.Options{
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

		newCtx = context.WithValue(ctx, AuthContextKey("token"), tokenInfo)

		return newCtx, nil
	}

	return config
}

func parseToken(token string, authConfig *AuthConfig) (jwt.Claims, error) {
	var parsedToken *jwt.Token
	var err error

	// Use JWKS, if enabled
	if authConfig.useJwks {
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

// StartDedicatedAuthServer starts a gRPC server containing just the auth service
func StartDedicatedAuthServer(address string) (sock net.Listener, server *grpc.Server, authService *service_auth.Service, err error) {
	// create a new socket for gRPC communication
	sock, err = net.Listen("tcp", address)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not listen: %w", err)
	}

	authService = service_auth.NewService()
	authService.CreateDefaultUser("clouditor", "clouditor")

	authConfig := ConfigureAuth(WithPublicKey(authService.GetPublicKey()))

	// We also add our authentication middleware, because we usually add additional service later
	server = grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_auth.UnaryServerInterceptor(authConfig.AuthFunc),
		),
	)
	auth.RegisterAuthenticationServer(server, authService)

	go func() {
		// serve the gRPC socket
		_ = server.Serve(sock)
	}()

	return sock, server, authService, nil
}
