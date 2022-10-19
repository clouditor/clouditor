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

package service

import (
	"context"

	"clouditor.io/clouditor/api/orchestrator"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang-jwt/jwt/v4"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

type RequestType int

const (
	AccessCreate RequestType = iota
	AccessRead
	AccessUpdate
	AccessDelete
)

// AuthorizationStrategy is an interface that implements a function which whether the current cloud service request can
// be fulfilled using the current access strategy.
type AuthorizationStrategy interface {
	CheckAccess(ctx context.Context, typ RequestType, req orchestrator.CloudServiceRequest) bool
	AllowedCloudServices(ctx context.Context) (all bool, IDs []string)
}

// AuthorizationStrategyJWT is an AuthorizationStrategy that expects the cloud service ID to be in a specific JWT claim
// key.
type AuthorizationStrategyJWT struct {
	Key string
}

// CheckAccess checks whether the current request can be fulfilled using the current access strategy.
func (a *AuthorizationStrategyJWT) CheckAccess(ctx context.Context, typ RequestType, req orchestrator.CloudServiceRequest) bool {
	var list []string

	return slices.Contains(list, req.GetCloudServiceId())
}

func (a *AuthorizationStrategyJWT) AllowedCloudServices(ctx context.Context) (all bool, list []string) {
	var err error
	var token string
	var claims jwt.MapClaims

	// Retrieve the raw token from the context
	token, err = grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return false, nil
	}

	// We need to re-parse the already validiated claim to get a specific key from the claims map.
	parser := jwt.NewParser()
	_, _, err = parser.ParseUnverified(token, &claims)
	if err != nil {
		return false, nil
	}

	for _, item := range claims[a.Key].([]interface{}) {
		list = append(list, item.(string))
	}

	return false, list
}

// AuthorizationStrategyJWT is an AuthorizationStrategy that allows all requests
type AuthorizationStrategyAllowAll struct{}

// CheckAccess checks whether the current request can be fulfilled using the current access strategy.
func (a *AuthorizationStrategyAllowAll) CheckAccess(ctx context.Context, typ RequestType, req orchestrator.CloudServiceRequest) bool {
	return true
}

func (a *AuthorizationStrategyAllowAll) AllowedCloudServices(ctx context.Context) (bool, []string) {
	return true, nil
}

var ErrPermissionDenied = status.Errorf(codes.PermissionDenied, "access denied")
