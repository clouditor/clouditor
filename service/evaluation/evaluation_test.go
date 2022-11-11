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

package evaluation

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"clouditor.io/clouditor/api/evaluation"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/defaults"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestNewService(t *testing.T) {
	type args struct {
		opts []ServiceOption
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{
			name: "Happy path",
			args: args{
				opts: []ServiceOption{},
			},
			want: &Service{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewService(tt.args.opts...)

			// we cannot compare the scheduler, so we first check if it is not empty and then nil it
			assert.NotEmpty(t, got.scheduler)
			got.scheduler = nil
			tt.want.scheduler = nil

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Evaluate(t *testing.T) {
	type args struct {
		ctx context.Context
		req *evaluation.EvaluateRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *evaluation.EvaluateResponse
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "Missing Control ID in request",
			args: args{
				ctx: context.Background(),
				req: &evaluation.EvaluateRequest{
					Toe: &orchestrator.TargetOfEvaluation{
						CloudServiceId: defaults.DefaultTargetCloudServiceID,
						CatalogId:      defaults.DefaultCatalogID,
						AssuranceLevel: &defaults.AssuranceLevelHigh,
					},
				},
			},
			wantResp: &evaluation.EvaluateResponse{
				Status:        false,
				StatusMessage: fmt.Sprintf(codes.InvalidArgument.String(), evaluation.ErrControlIDIsMissing.Error()),
			},
		},
		{
			name: "Happy path",
			args: args{
				ctx: context.Background(),
				req: &evaluation.EvaluateRequest{
					Toe: &orchestrator.TargetOfEvaluation{
						CloudServiceId: defaults.DefaultTargetCloudServiceID,
						CatalogId:      defaults.DefaultCatalogID,
						AssuranceLevel: &defaults.AssuranceLevelHigh,
					},
					ControlId: defaults.DefaultControlId,
				},
			},
			wantResp: &evaluation.EvaluateResponse{},
			wantErr:  assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			gotResp, err := s.Evaluate(tt.args.ctx, tt.args.req)
			if tt.wantErr != nil {
				tt.wantErr(t, err)
				assert.Equal(t, gotResp, tt.wantResp)
			}
		})
	}
}

func TestService_StartEvaluation(t *testing.T) {
	type args struct {
		req *evaluation.EvaluateRequest
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			s.StartEvaluation(tt.args.req)

		})
	}
}

func TestService_Shutdown(t *testing.T) {
	service := NewService()
	service.Shutdown()

	assert.False(t, service.scheduler.IsRunning())
}
