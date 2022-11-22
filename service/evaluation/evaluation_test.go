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
	"time"

	"clouditor.io/clouditor/api/evaluation"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/defaults"
	"github.com/go-co-op/gocron"
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
			want: &Service{
				orchestratorAddress: grpcTarget{
					target: "localhost:9090",
				},
				evaluation: make(map[string]*EvaluationScheduler),
			},
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

func TestService_StartEvaluation(t *testing.T) {
	type fields struct {
		scheduler  *gocron.Scheduler
		evaluation map[string]*EvaluationScheduler
	}
	type args struct {
		ctx context.Context
		req *evaluation.StartEvaluationRequest
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp *evaluation.StartEvaluationResponse
		wantErr  assert.ErrorAssertionFunc
	}{
		// {
		// 	name: "Evaluation already started for cloud service",
		// 	fields: fields{
		// 		scheduler: gocron.NewScheduler(time.UTC),
		// 		evaluation: map[string]*EvaluationScheduler{defaults.DefaultTargetCloudServiceID: {
		// 			scheduler:           gocron.NewScheduler(time.UTC),
		// 			evaluatedControlIDs: []string{defaults.DefaultEUCSControlID},
		// 		}},
		// 	},
		// 	args: args{
		// 		ctx: context.Background(),
		// 		req: &evaluation.StartEvaluationRequest{
		// 			Toe: &orchestrator.TargetOfEvaluation{
		// 				CloudServiceId: defaults.DefaultTargetCloudServiceID,
		// 				CatalogId:      defaults.DefaultCatalogID,
		// 				AssuranceLevel: &defaults.AssuranceLevelHigh,
		// 			},
		// 			ControlId:    defaults.DefaultEUCSControlID,
		// 			CategoryName: defaults.DefaultEUCSCategoryName,
		// 		},
		// 	},
		// 	wantResp: &evaluation.StartEvaluationResponse{},
		// 	wantErr:  assert.NoError,
		// },
		{
			name: "Missing Control ID in request",
			args: args{
				ctx: context.Background(),
				req: &evaluation.StartEvaluationRequest{
					Toe: &orchestrator.TargetOfEvaluation{
						CloudServiceId: defaults.DefaultTargetCloudServiceID,
						CatalogId:      defaults.DefaultCatalogID,
						AssuranceLevel: &defaults.AssuranceLevelHigh,
					},
					CategoryName: defaults.DefaultEUCSCategoryName,
				},
			},
			wantResp: &evaluation.StartEvaluationResponse{
				Status:        false,
				StatusMessage: fmt.Sprintf(codes.InvalidArgument.String(), evaluation.ErrControlIDIsMissing.Error()),
			},
		},
		{
			name: "Happy path",
			fields: fields{
				scheduler: gocron.NewScheduler(time.UTC),
				// evaluation: map[string]*EvaluationScheduler{defaults.DefaultTargetCloudServiceID: {
				// 	scheduler:           gocron.NewScheduler(time.UTC),
				// 	evaluatedControlIDs: []string{defaults.DefaultEUCSControlID},
				// }},
			},
			args: args{
				ctx: context.Background(),
				req: &evaluation.StartEvaluationRequest{
					Toe: &orchestrator.TargetOfEvaluation{
						CloudServiceId: defaults.DefaultTargetCloudServiceID,
						CatalogId:      defaults.DefaultCatalogID,
						AssuranceLevel: &defaults.AssuranceLevelHigh,
					},
					ControlId:    defaults.DefaultEUCSControlID,
					CategoryName: defaults.DefaultEUCSCategoryName,
				},
			},
			wantResp: &evaluation.StartEvaluationResponse{},
			wantErr:  assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				orchestratorAddress: grpcTarget{
					target: DefaultOrchestratorAddress,
				},
				scheduler:  tt.fields.scheduler,
				evaluation: tt.fields.evaluation,
			}
			gotResp, err := s.StartEvaluation(tt.args.ctx, tt.args.req)
			if tt.wantErr != nil {
				tt.wantErr(t, err)
				assert.Equal(t, gotResp, tt.wantResp)
			}
		})
	}
}

func TestService_Evaluate(t *testing.T) {
	type args struct {
		req *evaluation.StartEvaluationRequest
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
			s.StartEvaluation(context.Background(), tt.args.req)
		})
	}
}

func TestService_Shutdown(t *testing.T) {
	service := NewService()
	service.Shutdown()

	assert.False(t, service.scheduler.IsRunning())
}
