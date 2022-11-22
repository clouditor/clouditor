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

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evaluation"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/defaults"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/persistence/inmemory"
	"github.com/go-co-op/gocron"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestNewService(t *testing.T) {
	var myStorage, err = inmemory.NewStorage()
	assert.NoError(t, err)

	type args struct {
		opts []ServiceOption
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{
			name: "WithOrchestratorAddress",
			args: args{
				opts: []ServiceOption{ServiceOption(WithStorage(myStorage))},
			},
			want: &Service{
				scheduler: gocron.NewScheduler(time.UTC),
				orchestratorAddress: grpcTarget{
					target: defaults.DefaultOrchestratorAddress,
				},
				evaluation: make(map[string]*EvaluationScheduler),
				results:    make(map[string]*evaluation.EvaluationResult),
				storage:    myStorage,
			},
		},
		{
			name: "WithOrchestratorAddress",
			args: args{
				opts: []ServiceOption{ServiceOption(WithOrchestratorAddress(defaults.DefaultOrchestratorAddress))},
			},
			want: &Service{
				scheduler: gocron.NewScheduler(time.UTC),
				orchestratorAddress: grpcTarget{
					target: defaults.DefaultOrchestratorAddress,
				},
				evaluation: make(map[string]*EvaluationScheduler),
				results:    make(map[string]*evaluation.EvaluationResult),
			},
		},
		{
			name: "WithOAuth2Authorizer",
			args: args{
				opts: []ServiceOption{ServiceOption(WithOAuth2Authorizer(&clientcredentials.Config{}))},
			},
			want: &Service{
				scheduler: gocron.NewScheduler(time.UTC),
				orchestratorAddress: grpcTarget{
					target: "localhost:9090",
				},
				evaluation: make(map[string]*EvaluationScheduler),
				results:    make(map[string]*evaluation.EvaluationResult),
				authorizer: api.NewOAuthAuthorizerFromClientCredentials(&clientcredentials.Config{}),
			},
		},
		{
			name: "WithAuthorizer",
			args: args{
				opts: []ServiceOption{ServiceOption(WithAuthorizer(api.NewOAuthAuthorizerFromClientCredentials(&clientcredentials.Config{})))},
			},
			want: &Service{
				scheduler: gocron.NewScheduler(time.UTC),
				orchestratorAddress: grpcTarget{
					target: "localhost:9090",
				},
				evaluation: make(map[string]*EvaluationScheduler),
				results:    make(map[string]*evaluation.EvaluationResult),
				authorizer: api.NewOAuthAuthorizerFromClientCredentials(&clientcredentials.Config{}),
			},
		},
		{
			name: "Happy path",
			args: args{
				opts: []ServiceOption{},
			},
			want: &Service{
				scheduler: gocron.NewScheduler(time.UTC),
				orchestratorAddress: grpcTarget{
					target: "localhost:9090",
				},
				evaluation: make(map[string]*EvaluationScheduler),
				results:    make(map[string]*evaluation.EvaluationResult),
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
		scheduler           *gocron.Scheduler
		evaluation          map[string]*EvaluationScheduler
		results             map[string]*evaluation.EvaluationResult
		orchestratorAddress grpcTarget
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
				StatusMessage: evaluation.ErrControlIDIsMissing.Error(),
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, evaluation.ErrControlIDIsMissing.Error())
			},
		},
		{
			name: "Happy path",
			fields: fields{
				scheduler:  gocron.NewScheduler(time.UTC),
				evaluation: make(map[string]*EvaluationScheduler),
				results:    make(map[string]*evaluation.EvaluationResult),
				orchestratorAddress: grpcTarget{
					target: defaults.DefaultOrchestratorAddress,
				},
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
	type fields struct {
		evaluation.UnimplementedEvaluationServer
		scheduler           *gocron.Scheduler
		orchestratorClient  orchestrator.OrchestratorClient
		orchestratorAddress grpcTarget
		authorizer          api.Authorizer
		evaluation          map[string]*EvaluationScheduler
		results             map[string]*evaluation.EvaluationResult
		storage             persistence.Storage
	}
	type args struct {
		req *evaluation.StartEvaluationRequest
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantFunc assert.ValueAssertionFunc
	}{
		{
			name: "Error getting metrics from control",
			args: args{
				req: &evaluation.StartEvaluationRequest{
					Toe: &orchestrator.TargetOfEvaluation{
						CloudServiceId: defaults.DefaultTargetCloudServiceID,
						CatalogId:      defaults.DefaultCatalogID,
						AssuranceLevel: &defaults.AssuranceLevelHigh,
					},
					CategoryName: defaults.DefaultEUCSCategoryName,
					ControlId:    defaults.DefaultEUCSControlID,
				},
			},
			wantFunc: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				service, ok := i1.(*Service)
				if !assert.True(tt, ok) {
					return false
				}

				return assert.Empty(t, service.evaluation[createSchedulerTag(defaults.DefaultTargetCloudServiceID, defaults.DefaultEUCSControlID)])
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				scheduler:           tt.fields.scheduler,
				orchestratorClient:  tt.fields.orchestratorClient,
				orchestratorAddress: tt.fields.orchestratorAddress,
				authorizer:          tt.fields.authorizer,
				evaluation:          tt.fields.evaluation,
				results:             tt.fields.results,
				storage:             tt.fields.storage,
			}
			s.Evaluate(tt.args.req)

			if tt.wantFunc != nil {
				tt.wantFunc(t, s)
			}
		})
	}
}

func TestService_Shutdown(t *testing.T) {
	service := NewService()
	service.Shutdown()

	assert.False(t, service.scheduler.IsRunning())
}

func TestService_SetAuthorizer(t *testing.T) {

	type args struct {
		auth api.Authorizer
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Happy path",
			args: args{
				auth: api.NewOAuthAuthorizerFromClientCredentials(&clientcredentials.Config{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			assert.Empty(t, s.authorizer)
			s.SetAuthorizer(tt.args.auth)
			assert.NotEmpty(t, s.authorizer)
		})
	}
}

func TestService_Authorizer(t *testing.T) {

	type fields struct {
		authorizer api.Authorizer
	}

	tests := []struct {
		name   string
		fields fields
		want   api.Authorizer
	}{
		{
			name: "No authorizer set",
			fields: fields{
				authorizer: nil,
			},
			want: nil,
		},
		{
			name: "Happy path",
			fields: fields{
				authorizer: api.NewOAuthAuthorizerFromClientCredentials(&clientcredentials.Config{}),
			},
			want: api.NewOAuthAuthorizerFromClientCredentials(&clientcredentials.Config{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				authorizer: tt.fields.authorizer,
			}
			if got := s.Authorizer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.Authorizer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createSchedulerTag(t *testing.T) {
	type args struct {
		cloudServiceId string
		controlId      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Input controlId empty",
			args: args{
				controlId: defaults.DefaultEUCSControlID,
			},
			want: "",
		},
		{
			name: "Input cloudServiceId empty",
			args: args{
				cloudServiceId: defaults.DefaultTargetCloudServiceID,
			},
			want: "",
		},
		{
			name: "Happy path",
			args: args{
				cloudServiceId: defaults.DefaultTargetCloudServiceID,
				controlId:      defaults.DefaultEUCSControlID,
			},
			want: fmt.Sprintf("%s-%s", defaults.DefaultTargetCloudServiceID, defaults.DefaultEUCSControlID),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createSchedulerTag(tt.args.cloudServiceId, tt.args.controlId); got != tt.want {
				t.Errorf("createSchedulerTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getMapping(t *testing.T) {
	type args struct {
		results []*assessment.AssessmentResult
		metrics []*assessment.Metric
	}
	tests := []struct {
		name            string
		args            args
		wantMappingList []*mappingResultMetric
	}{
		{
			name: "Empty metrics",
			args: args{
				results: []*assessment.AssessmentResult{
					{
						Id:             "11111111-1111-1111-1111-111111111111",
						MetricId:       "TransportEncryptionEnabled",
						CloudServiceId: defaults.DefaultTargetCloudServiceID,
						Compliant:      true,
					},
				},
				metrics: []*assessment.Metric{},
			},
			wantMappingList: nil,
		},
		{
			name: "Missing metrics",
			args: args{
				results: []*assessment.AssessmentResult{
					{
						Id:             "11111111-1111-1111-1111-111111111111",
						MetricId:       "TransportEncryptionEnabled",
						CloudServiceId: defaults.DefaultTargetCloudServiceID,
						Compliant:      true,
					},
				},
				metrics: nil,
			},
			wantMappingList: nil,
		},
		{
			name: "Empty assessment results",
			args: args{
				results: []*assessment.AssessmentResult{},
				metrics: []*assessment.Metric{
					{
						Id:          "TransportEncryptionEnabled",
						Name:        "Transport Encryption: Enabled",
						Description: "This metric describes, whether transport encryption is turned on or not",
						Scale:       assessment.Metric_ORDINAL,
						Range: &assessment.Range{
							Range: &assessment.Range_AllowedValues{AllowedValues: &assessment.AllowedValues{
								Values: []*structpb.Value{
									structpb.NewBoolValue(false),
									structpb.NewBoolValue(true),
								}}}},
					},
				},
			},
			wantMappingList: []*mappingResultMetric{
				{
					metricName: "TransportEncryptionEnabled",
					results:    []*assessment.AssessmentResult{},
				},
			},
		},
		{
			name: "Missing assessment results",
			args: args{
				results: nil,
				metrics: []*assessment.Metric{
					{
						Id:          "TransportEncryptionEnabled",
						Name:        "Transport Encryption: Enabled",
						Description: "This metric describes, whether transport encryption is turned on or not",
						Scale:       assessment.Metric_ORDINAL,
						Range: &assessment.Range{
							Range: &assessment.Range_AllowedValues{AllowedValues: &assessment.AllowedValues{
								Values: []*structpb.Value{
									structpb.NewBoolValue(false),
									structpb.NewBoolValue(true),
								}}}},
					},
				},
			},
			wantMappingList: []*mappingResultMetric{
				{
					metricName: "TransportEncryptionEnabled",
					results:    []*assessment.AssessmentResult{},
				},
			},
		},
		{
			name: "Happy path",
			args: args{
				results: []*assessment.AssessmentResult{
					{
						Id:             "11111111-1111-1111-1111-111111111111",
						MetricId:       "TransportEncryptionEnabled",
						CloudServiceId: defaults.DefaultTargetCloudServiceID,
						Compliant:      true,
					},
				},
				metrics: []*assessment.Metric{
					{
						Id:          "TransportEncryptionEnabled",
						Name:        "Transport Encryption: Enabled",
						Description: "This metric describes, whether transport encryption is turned on or not",
						Scale:       assessment.Metric_ORDINAL,
						Range: &assessment.Range{
							Range: &assessment.Range_AllowedValues{AllowedValues: &assessment.AllowedValues{
								Values: []*structpb.Value{
									structpb.NewBoolValue(false),
									structpb.NewBoolValue(true),
								}}}},
					},
				},
			},
			wantMappingList: []*mappingResultMetric{
				{
					metricName: "TransportEncryptionEnabled",
					results: []*assessment.AssessmentResult{
						{
							Id:             "11111111-1111-1111-1111-111111111111",
							MetricId:       "TransportEncryptionEnabled",
							CloudServiceId: defaults.DefaultTargetCloudServiceID,
							Compliant:      true,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMappingList := getMapping(tt.args.results, tt.args.metrics)
			assert.Equal(t, tt.wantMappingList, gotMappingList)
		})
	}
}
