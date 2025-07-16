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
	"testing"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/persistence"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
)

const (
	TestCustomClaims   = "TargetOfEvaluationid"
	TestAllowAllClaims = "cladmin"
)

var (
	// TestContextOnlyService1 is an incoming context with a JWT that only allows access to target of evaluation ID
	// 11111111-1111-1111-1111-111111111111
	TestContextOnlyService1 context.Context

	// TestContextOnlyService1 is an incoming context with a JWT that allows access to all target of evaluations
	TestContextAllowAll context.Context

	// TestBrokenContext contains an invalid JWT
	TestBrokenContext = metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
		"authorization": "bearer what",
	}))

	// TestClaimsOnlyService1 contains claims that authorize the user for the target of evaluation
	// 11111111-1111-1111-1111-111111111111.
	TestClaimsOnlyService1 = jwt.MapClaims{
		"sub": "me",
		"TargetOfEvaluationid": []string{
			testdata.MockTargetOfEvaluationID1,
		},
		"other": []int{1, 2},
	}

	// TestClaimsOnlyService1 contains claims that authorize the user for all target of evaluations.
	TestClaimsAllowAll = jwt.MapClaims{
		"sub":     "me",
		"cladmin": true,
	}
)

func init() {
	var (
		err   error
		token *jwt.Token
		t     string
	)

	// Create a new token instead of hard-coding one
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, &TestClaimsOnlyService1)
	t, err = token.SignedString([]byte("mykey"))
	if err != nil {
		panic(err)
	}

	// Create a context containing our token
	TestContextOnlyService1 = metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
		"authorization": "bearer " + t,
	}))

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, &TestClaimsAllowAll)
	t, err = token.SignedString([]byte("mykey"))
	if err != nil {
		panic(err)
	}

	// Create a context containing our token
	TestContextAllowAll = metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
		"authorization": "bearer " + t,
	}))
}

func TestAuthorizationStrategyAllowAll_CheckAccess(t *testing.T) {
	type args struct {
		ctx context.Context
		typ RequestType
		req api.TargetOfEvaluationRequest
	}
	tests := []struct {
		name string
		a    *AuthorizationStrategyAllowAll
		args args
		want bool
	}{
		{
			name: "always true",
			a:    &AuthorizationStrategyAllowAll{},
			args: args{
				ctx: context.Background(),
				typ: AccessCreate,
				req: &orchestrator.GetTargetOfEvaluationRequest{TargetOfEvaluationId: config.DefaultTargetOfEvaluationID},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthorizationStrategyAllowAll{}
			if got := a.CheckAccess(tt.args.ctx, tt.args.typ, tt.args.req); got != tt.want {
				t.Errorf("AuthorizationStrategyAllowAll.CheckAccess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthorizationStrategyAllowAll_AllowedTargetOfEvaluations(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		a        *AuthorizationStrategyAllowAll
		args     args
		wantAll  bool
		wantList []string
	}{
		{
			name:     "all allowed",
			a:        &AuthorizationStrategyAllowAll{},
			args:     args{},
			wantAll:  true,
			wantList: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthorizationStrategyAllowAll{}
			gotAll, gotList := a.AllowedTargetOfEvaluations(tt.args.ctx)
			assert.Equal(t, tt.wantAll, gotAll)
			assert.Equal(t, tt.wantList, gotList)
		})
	}
}

func TestAuthorizationStrategyJWT_CheckAccess(t *testing.T) {
	type fields struct {
		TargetOfEvaluationsKey string
		AllowAllKey            string
	}
	type args struct {
		ctx context.Context
		typ RequestType
		req api.TargetOfEvaluationRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "valid context",
			fields: fields{
				TargetOfEvaluationsKey: TestCustomClaims,
			},
			args: args{
				ctx: TestContextOnlyService1,
				typ: AccessRead,
				req: &orchestrator.GetTargetOfEvaluationRequest{TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1},
			},
			want: true,
		},
		{
			name: "valid context, allow all",
			fields: fields{
				AllowAllKey: TestAllowAllClaims,
			},
			args: args{
				ctx: TestContextAllowAll,
				typ: AccessRead,
				req: &orchestrator.GetTargetOfEvaluationRequest{TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1},
			},
			want: true,
		},
		{
			name: "valid context, wrong claim",
			fields: fields{
				TargetOfEvaluationsKey: "sub",
			},
			args: args{
				ctx: TestContextOnlyService1,
				typ: AccessRead,
				req: &orchestrator.GetTargetOfEvaluationRequest{TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1},
			},
			want: false,
		},
		{
			name: "valid context, ignore non-string",
			fields: fields{
				TargetOfEvaluationsKey: "other",
			},
			args: args{
				ctx: TestContextOnlyService1,
				typ: AccessRead,
				req: &orchestrator.GetTargetOfEvaluationRequest{TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1},
			},
			want: false,
		},
		{
			name: "missing token",
			fields: fields{
				TargetOfEvaluationsKey: TestCustomClaims,
			},
			args: args{
				ctx: context.Background(),
				typ: AccessRead,
				req: &orchestrator.GetTargetOfEvaluationRequest{TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1},
			},
			want: false,
		},
		{
			name: "broken token",
			fields: fields{
				TargetOfEvaluationsKey: TestCustomClaims,
			},
			args: args{
				ctx: TestBrokenContext,
				typ: AccessRead,
				req: &orchestrator.GetTargetOfEvaluationRequest{TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthorizationStrategyJWT{
				TargetOfEvaluationsKey: tt.fields.TargetOfEvaluationsKey,
				AllowAllKey:            tt.fields.AllowAllKey,
			}
			if got := a.CheckAccess(tt.args.ctx, tt.args.typ, tt.args.req); got != tt.want {
				t.Errorf("AuthorizationStrategyJWT.CheckAccess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthorizationStrategyJWT_AllowedTargetOfEvaluations(t *testing.T) {
	type fields struct {
		TargetOfEvaluationsKey string
		AllowAllKey            string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantAll  bool
		wantList []string
	}{
		{
			name: "valid context",
			fields: fields{
				TargetOfEvaluationsKey: TestCustomClaims,
			},
			args: args{
				ctx: TestContextOnlyService1,
			},
			wantAll:  false,
			wantList: []string{testdata.MockTargetOfEvaluationID1},
		},
		{
			name: "valid context, allow all",
			fields: fields{
				AllowAllKey: TestAllowAllClaims,
			},
			args: args{
				ctx: TestContextAllowAll,
			},
			wantAll:  true,
			wantList: nil,
		},
		{
			name: "valid context, wrong claim",
			fields: fields{
				TargetOfEvaluationsKey: "sub",
			},
			args: args{
				ctx: TestContextOnlyService1,
			},
			wantAll:  false,
			wantList: nil,
		},
		{
			name: "valid context, ignore non-string",
			fields: fields{
				TargetOfEvaluationsKey: "other",
			},
			args: args{
				ctx: TestContextOnlyService1,
			},
			wantAll:  false,
			wantList: nil,
		},
		{
			name: "missing token",
			fields: fields{
				TargetOfEvaluationsKey: TestCustomClaims,
			},
			args: args{
				ctx: context.Background(),
			},
			wantAll:  false,
			wantList: nil,
		},
		{
			name: "broken token",
			fields: fields{
				TargetOfEvaluationsKey: TestCustomClaims,
			},
			args: args{
				ctx: TestBrokenContext,
			},
			wantAll:  false,
			wantList: nil,
		},
		{
			name:     "nil context",
			fields:   fields{},
			args:     args{ctx: nil},
			wantAll:  false,
			wantList: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthorizationStrategyJWT{
				TargetOfEvaluationsKey: tt.fields.TargetOfEvaluationsKey,
				AllowAllKey:            tt.fields.AllowAllKey,
			}
			gotAll, gotList := a.AllowedTargetOfEvaluations(tt.args.ctx)
			assert.Equal(t, tt.wantAll, gotAll)
			assert.Equal(t, tt.wantList, gotList)
		})
	}
}

func TestAuthorizationStrategyStorage_CheckAccess(t *testing.T) {
	type fields struct {
		DB persistence.Storage
	}
	type args struct {
		ctx context.Context
		in1 RequestType
		req api.TargetOfEvaluationRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "valid context, user not yet created",
			fields: fields{
				DB: testutil.NewInMemoryStorage(t),
			},
			args: args{
				ctx: TestContextOnlyService1,
				req: &orchestrator.GetTargetOfEvaluationRequest{TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthorizationStrategyStorage{
				DB: tt.fields.DB,
			}
			if got := a.CheckAccess(tt.args.ctx, tt.args.in1, tt.args.req); got != tt.want {
				t.Errorf("AuthorizationStrategyStorage.CheckAccess() = %v, want %v", got, tt.want)
			}
		})
	}
}
