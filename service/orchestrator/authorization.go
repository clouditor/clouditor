// Copyright 2021-2022 Fraunhofer AISEC
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

package orchestrator

import (
	"context"

	"clouditor.io/clouditor/api/orchestrator"
	"golang.org/x/exp/slices"

	"github.com/golang-jwt/jwt/v4"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

// AuthorizationStrategy is an interface that implements a function which whether the current cloud service request can
// be fulfilled using the current access strategy.
type AuthorizationStrategy interface {
	CheckAccess(ctx context.Context, req orchestrator.CloudServiceRequest) bool
}

// AuthorizationStrategyJWT is an AuthorizationStrategy that expects the cloud service ID to be in a specific JWT claim
// key.
type AuthorizationStrategyJWT struct {
	key string
}

// CheckAccess checks whether the current request can be fulfilled using the current access strategy.
func (a *AuthorizationStrategyJWT) CheckAccess(ctx context.Context, req orchestrator.CloudServiceRequest) bool {
	var err error
	var token string
	var claims jwt.MapClaims
	var list []string

	// Retrieve the raw token from the context
	token, err = grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return false
	}

	// We need to re-parse the already validiated claim to get a specific key from the claims map.
	parser := jwt.NewParser()
	_, _, err = parser.ParseUnverified(token, claims)
	if err != nil {
		return false
	}

	list = claims[a.key].([]string)

	return slices.Contains(list, req.GetCloudServiceId())
}

// WithAuthorizationStrategyJWT is an option that configures an JWT-based authorization strategy using a specific claim key.
func WithAuthorizationStrategyJWT(key string) ServiceOption {
	return func(s *Service) {
		s.authz = &AuthorizationStrategyJWT{key: key}
	}
}

// AuthorizationStrategyJWT is an AuthorizationStrategy that allows all requests
type AuthorizationStrategyAllowAll struct{}

// CheckAccess checks whether the current request can be fulfilled using the current access strategy.
func (a *AuthorizationStrategyAllowAll) CheckAccess(ctx context.Context, req orchestrator.CloudServiceRequest) bool {
	return true
}
