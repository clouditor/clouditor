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
	"time"

	"clouditor.io/clouditor/rest"
	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var log *logrus.Entry

type AuthConfig struct {
	jwksUrl string

	Jwks     *keyfunc.JWKS
	AuthFunc grpc_auth.AuthFunc
}

var DefaultJwksUrl = fmt.Sprintf("http://localhost:%d/.well-known/jwks.json", rest.DefaultAPIHTTPPort)

func init() {
	log = logrus.WithField("component", "service-auth")
}

type AuthOption func(*AuthConfig)

func WithJwksUrl(url string) AuthOption {
	return func(ac *AuthConfig) {
		ac.jwksUrl = url
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
		if config.Jwks == nil {
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
	parsedToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, authConfig.Jwks.Keyfunc)

	if err != nil {
		return nil, fmt.Errorf("could not validate JWT: %w", err)
	}

	return parsedToken.Claims, nil
}
