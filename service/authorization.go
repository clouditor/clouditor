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
	"slices"

	"clouditor.io/clouditor/v2/api"
	"github.com/golang-jwt/jwt/v5"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RequestType specifies the type of request, usually CRUD.
type RequestType int

const (
	AccessCreate RequestType = iota
	AccessRead
	AccessUpdate
	AccessDelete
)

// ErrPermissionDenied represents an error, where permission to fulfill the request is denied.
var ErrPermissionDenied = status.Errorf(codes.PermissionDenied, "access denied")

// AuthorizationStrategy is an interface that implements a function which
// checkers whether the current target of evaluation request can be fulfilled using the
// supplied context (e.g., based on the authenticated user).
type AuthorizationStrategy interface {
	CheckAccess(ctx context.Context, typ RequestType, req api.TargetOfEvaluationRequest) bool
	AllowedTargetOfEvaluations(ctx context.Context) (all bool, IDs []string)
}

// AuthorizationStrategyJWT is an AuthorizationStrategy that expects a list of target of evaluation IDs to be in a specific JWT
// claim key.
type AuthorizationStrategyJWT struct {
	TargetOfEvaluationsKey string
	AllowAllKey            string
}

// CheckAccess checks whether the current request can be fulfilled using the current access strategy.
func (a *AuthorizationStrategyJWT) CheckAccess(ctx context.Context, _ RequestType, req api.TargetOfEvaluationRequest) bool {
	var (
		list []string
		all  bool
	)

	// Retrieve the list of allowed target of evaluations.
	all, list = a.AllowedTargetOfEvaluations(ctx)

	if all {
		return true
	}

	return slices.Contains(list, req.GetTargetOfEvaluationId())
}

// AllowedTargetOfEvaluations retrieves a list of allowed target of evaluation IDs according to the current access strategy.
func (a *AuthorizationStrategyJWT) AllowedTargetOfEvaluations(ctx context.Context) (all bool, list []string) {
	var (
		err    error
		ok     bool
		token  string
		claims jwt.MapClaims
		l      []interface{}
		s      string
	)

	// Check, if the context is nil
	if ctx == nil {
		log.Debugf("Retrieving allowed target of evaluations failed because of an empty context")
		return false, nil
	}

	// Retrieve the raw token from the context
	token, err = grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		log.Debugf("Retrieving allowed target of evaluations from token failed: %v", err)
		return false, nil
	}

	// We need to re-parse the already validated claim to get a specific key from the claims map.
	parser := jwt.NewParser()
	_, _, err = parser.ParseUnverified(token, &claims)
	if err != nil {
		log.Debugf("Retrieving allowed target of evaluations from token failed: %v", err)
		return false, nil
	}

	// Let's look for an allow all key
	if b, ok := claims[a.AllowAllKey].(bool); ok && b {
		return true, nil
	}

	// We are looking for an array claim
	if l, ok = claims[a.TargetOfEvaluationsKey].([]interface{}); !ok {
		log.Debug("Retrieving allowed target of evaluations from token failed: specified claims key is not an array", err)
		return false, nil
	}

	// Loop through the array claim and add all string values to our list
	for _, item := range l {
		if s, ok = item.(string); !ok {
			continue
		}

		list = append(list, s)
	}

	return false, list
}

// AuthorizationStrategyAllowAll is an AuthorizationStrategy that allows all requests.
type AuthorizationStrategyAllowAll struct{}

// CheckAccess checks whether the current request can be fulfilled using the current access strategy. Returns true since
// strategy is `AuthorizationStrategyAllowAll`
func (*AuthorizationStrategyAllowAll) CheckAccess(_ context.Context, _ RequestType, _ api.TargetOfEvaluationRequest) bool {
	return true
}

// AllowedTargetOfEvaluations retrieves a list of allowed target of evaluation IDs according to the current access strategy. Returns
// `all = true` since strategy is `AuthorizationStrategyAllowAll`
func (*AuthorizationStrategyAllowAll) AllowedTargetOfEvaluations(_ context.Context) (all bool, list []string) {
	return true, nil
}
