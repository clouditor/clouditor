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
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"
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

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_StartEvaluation(t *testing.T) {
	schedulerTag := createSchedulerTag(defaults.DefaultTargetCloudServiceID, defaults.DefaultEUCSControlID)

	type fields struct {
		evaluation          map[string]*EvaluationScheduler
		results             map[string]*evaluation.EvaluationResult
		orchestratorAddress grpcTarget
	}
	type args struct {
		ctx                context.Context
		req                *evaluation.StartEvaluationRequest
		schedulerTag       string
		evaluatedControlID string
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
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
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
			name: "Scheduler already running",
			fields: fields{
				evaluation: map[string]*EvaluationScheduler{
					schedulerTag: {
						scheduler:          gocron.NewScheduler(time.UTC),
						evaluatedControlID: defaults.DefaultEUCSControlID,
					},
				},
				results: make(map[string]*evaluation.EvaluationResult),
				orchestratorAddress: grpcTarget{
					target: defaults.DefaultOrchestratorAddress,
				},
			},
			args: args{
				ctx:          context.Background(),
				schedulerTag: schedulerTag,
				req: &evaluation.StartEvaluationRequest{
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CloudServiceId: defaults.DefaultTargetCloudServiceID,
						CatalogId:      defaults.DefaultCatalogID,
						AssuranceLevel: &defaults.AssuranceLevelHigh,
					},
					ControlId:    defaults.DefaultEUCSControlID,
					CategoryName: defaults.DefaultEUCSCategoryName,
				},
			},
			wantResp: &evaluation.StartEvaluationResponse{},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, fmt.Sprintf("Cloud Service '%s' is being evaluated with Control %s already.", defaults.DefaultTargetCloudServiceID, defaults.DefaultEUCSControlID))
			},
		},
		{
			name: "Happy path",
			fields: fields{
				evaluation: make(map[string]*EvaluationScheduler),
				results:    make(map[string]*evaluation.EvaluationResult),
				orchestratorAddress: grpcTarget{
					target: defaults.DefaultOrchestratorAddress,
				},
			},
			args: args{
				ctx: context.Background(),
				req: &evaluation.StartEvaluationRequest{
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CloudServiceId: defaults.DefaultTargetCloudServiceID,
						CatalogId:      defaults.DefaultCatalogID,
						AssuranceLevel: &defaults.AssuranceLevelHigh,
					},
					ControlId:    defaults.DefaultEUCSControlID,
					CategoryName: defaults.DefaultEUCSCategoryName,
				},
				schedulerTag:       schedulerTag,
				evaluatedControlID: defaults.DefaultEUCSControlID,
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
				evaluation: tt.fields.evaluation,
			}

			// Start the scheduler
			if s.evaluation[tt.args.schedulerTag] != nil {
				s.evaluation[tt.args.schedulerTag].scheduler.StartAsync()
			}

			gotResp, err := s.StartEvaluation(tt.args.ctx, tt.args.req)
			tt.wantErr(t, err)
			assert.Equal(t, gotResp, tt.wantResp)

			if err == nil {
				// We must check if the contol ID is stored
				assert.Equal(t, s.evaluation[tt.args.schedulerTag].evaluatedControlID, tt.args.evaluatedControlID)

				// We must check if the scheduler is running
				assert.True(t, s.evaluation[tt.args.schedulerTag].scheduler.IsRunning())
			}
		})
	}
}

func TestService_Evaluate(t *testing.T) {
	type fields struct {
		evaluation.UnimplementedEvaluationServer
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
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
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
				orchestratorClient:  tt.fields.orchestratorClient,
				orchestratorAddress: tt.fields.orchestratorAddress,
				authorizer:          tt.fields.authorizer,
				evaluation:          tt.fields.evaluation,
				results:             tt.fields.results,
				storage:             tt.fields.storage,
			}
			s.Evaluate(tt.args.req)

			tt.wantFunc(t, s)

		})
	}
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

func TestService_ListEvaluationResults(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestratorClient            orchestrator.OrchestratorClient
		orchestratorAddress           grpcTarget
		authorizer                    api.Authorizer
		evaluation                    map[string]*EvaluationScheduler
		results                       map[string]*evaluation.EvaluationResult
		storage                       persistence.Storage
	}
	type args struct {
		in0 context.Context
		req *evaluation.ListEvaluationResultsRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *evaluation.ListEvaluationResultsResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "multiple page result - first page",
			fields: fields{results: map[string]*evaluation.EvaluationResult{
				"11111111-1111-1111-1111-111111111111": {Id: "11111111-1111-1111-1111-111111111111"},
				"22222222-2222-2222-2222-222222222222": {Id: "22222222-2222-2222-2222-222222222222"},
				"33333333-3333-3333-3333-333333333333": {Id: "33333333-3333-3333-3333-333333333333"},
				"44444444-4444-4444-4444-444444444444": {Id: "44444444-4444-4444-4444-444444444444"},
			}},
			args: args{
				in0: context.Background(),
				req: &evaluation.ListEvaluationResultsRequest{
					PageSize: 2,
				},
			},
			wantRes: &evaluation.ListEvaluationResultsResponse{
				Results: []*evaluation.EvaluationResult{
					{
						Id: "11111111-1111-1111-1111-111111111111",
					},
					{
						Id: "22222222-2222-2222-2222-222222222222",
					},
				},
				NextPageToken: func() string {
					token, _ := (&api.PageToken{Start: 2, Size: 2}).Encode()
					return token
				}(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "multiple page result - second page",
			fields: fields{results: map[string]*evaluation.EvaluationResult{
				"11111111-1111-1111-1111-111111111111": {Id: "11111111-1111-1111-1111-111111111111"},
				"22222222-2222-2222-2222-222222222222": {Id: "22222222-2222-2222-2222-222222222222"},
				"33333333-3333-3333-3333-333333333333": {Id: "33333333-3333-3333-3333-333333333333"},
				"44444444-4444-4444-4444-444444444444": {Id: "44444444-4444-4444-4444-444444444444"},
			}},
			args: args{
				in0: context.Background(),
				req: &evaluation.ListEvaluationResultsRequest{
					PageSize: 2,
					PageToken: func() string {
						token, _ := (&api.PageToken{Start: 2, Size: 2}).Encode()
						return token
					}(),
				},
			},
			wantRes: &evaluation.ListEvaluationResultsResponse{
				Results: []*evaluation.EvaluationResult{
					{
						Id: "33333333-3333-3333-3333-333333333333",
					},
					{
						Id: "44444444-4444-4444-4444-444444444444",
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "List all results",
			fields: fields{
				results: map[string]*evaluation.EvaluationResult{
					"11111111-1111-1111-1111-111111111111": {Id: "11111111-1111-1111-1111-111111111111"},
					"22222222-2222-2222-2222-222222222222": {Id: "22222222-2222-2222-2222-222222222222"},
					"33333333-3333-3333-3333-333333333333": {Id: "33333333-3333-3333-3333-333333333333"},
					"44444444-4444-4444-4444-444444444444": {Id: "44444444-4444-4444-4444-444444444444"},
				},
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.ListEvaluationResultsRequest{},
			},
			wantRes: &evaluation.ListEvaluationResultsResponse{
				Results: []*evaluation.EvaluationResult{
					{
						Id: "11111111-1111-1111-1111-111111111111",
					},
					{
						Id: "22222222-2222-2222-2222-222222222222",
					},
					{
						Id: "33333333-3333-3333-3333-333333333333",
					},
					{
						Id: "44444444-4444-4444-4444-444444444444",
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				UnimplementedEvaluationServer: tt.fields.UnimplementedEvaluationServer,
				orchestratorClient:            tt.fields.orchestratorClient,
				orchestratorAddress:           tt.fields.orchestratorAddress,
				authorizer:                    tt.fields.authorizer,
				evaluation:                    tt.fields.evaluation,
				results:                       tt.fields.results,
				storage:                       tt.fields.storage,
			}
			gotRes, err := s.ListEvaluationResults(tt.args.in0, tt.args.req)

			tt.wantErr(t, err)
			assert.Equal(t, gotRes, tt.wantRes)

		})
	}
}

func TestService_StopEvaluation(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestratorClient            orchestrator.OrchestratorClient
		orchestratorAddress           grpcTarget
		authorizer                    api.Authorizer
		evaluation                    map[string]*EvaluationScheduler
		results                       map[string]*evaluation.EvaluationResult
		storage                       persistence.Storage
	}
	type args struct {
		in0          context.Context
		req          *evaluation.StopEvaluationRequest
		schedulerTag string
		// schedulerRunning is a bool for starting the scheduler in the test
		schedulerRunning bool
		// tag is a bool for deciding if the scheduler job needs a tag
		tag bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *evaluation.StopEvaluationResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Control Id in request is missing",
			args: args{
				in0: context.Background(),
				req: &evaluation.StopEvaluationRequest{},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, evaluation.ErrControlIDIsMissing.Error())
			},
		},
		{
			name: "No entry for cloud service available in evaluation map",
			args: args{
				in0: context.Background(),
				req: &evaluation.StopEvaluationRequest{
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CloudServiceId: defaults.DefaultTargetCloudServiceID,
						CatalogId:      defaults.DefaultCatalogID,
						AssuranceLevel: &defaults.AssuranceLevelHigh,
					},
					ControlId:    defaults.DefaultEUCSControlID,
					CategoryName: defaults.DefaultEUCSCategoryName,
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, fmt.Sprintf("evaluation of cloud service %s not running", defaults.DefaultTargetCloudServiceID))
			},
		},
		{
			name: "Evaluation for cloud service id is not running",
			args: args{
				in0: context.Background(),
				req: &evaluation.StopEvaluationRequest{
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CloudServiceId: defaults.DefaultTargetCloudServiceID,
						CatalogId:      defaults.DefaultCatalogID,
						AssuranceLevel: &defaults.AssuranceLevelHigh,
					},
					ControlId:    defaults.DefaultEUCSControlID,
					CategoryName: defaults.DefaultEUCSCategoryName,
				},
				schedulerTag: createSchedulerTag(defaults.DefaultTargetCloudServiceID, defaults.DefaultEUCSControlID),
			},
			fields: fields{
				evaluation: map[string]*EvaluationScheduler{
					createSchedulerTag(defaults.DefaultTargetCloudServiceID, defaults.DefaultEUCSControlID): {
						scheduler:          gocron.NewScheduler(time.UTC),
						evaluatedControlID: defaults.DefaultEUCSControlID,
					},
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, fmt.Sprintf("evaluation of cloud service %s not running", defaults.DefaultTargetCloudServiceID))
			},
		},
		{
			name: "Missing tag",
			args: args{
				in0: context.Background(),
				req: &evaluation.StopEvaluationRequest{
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CloudServiceId: defaults.DefaultTargetCloudServiceID,
						CatalogId:      defaults.DefaultCatalogID,
						AssuranceLevel: &defaults.AssuranceLevelHigh,
					},
					ControlId:    defaults.DefaultEUCSControlID,
					CategoryName: defaults.DefaultEUCSCategoryName,
				},
				schedulerTag:     createSchedulerTag(defaults.DefaultTargetCloudServiceID, defaults.DefaultEUCSControlID),
				schedulerRunning: true,
			},
			fields: fields{
				evaluation: map[string]*EvaluationScheduler{
					createSchedulerTag(defaults.DefaultTargetCloudServiceID, defaults.DefaultEUCSControlID): {
						scheduler:          gocron.NewScheduler(time.UTC),
						evaluatedControlID: defaults.DefaultEUCSControlID,
					},
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "error at stopping scheduler")
			},
		},
		{
			name: "Happy path",
			args: args{
				in0: context.Background(),
				req: &evaluation.StopEvaluationRequest{
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CloudServiceId: defaults.DefaultTargetCloudServiceID,
						CatalogId:      defaults.DefaultCatalogID,
						AssuranceLevel: &defaults.AssuranceLevelHigh,
					},
					ControlId:    defaults.DefaultEUCSControlID,
					CategoryName: defaults.DefaultEUCSCategoryName,
				},
				schedulerTag:     createSchedulerTag(defaults.DefaultTargetCloudServiceID, defaults.DefaultEUCSControlID),
				schedulerRunning: true,
				tag:              true,
			},
			fields: fields{
				evaluation: map[string]*EvaluationScheduler{
					createSchedulerTag(defaults.DefaultTargetCloudServiceID, defaults.DefaultEUCSControlID): {
						scheduler:          gocron.NewScheduler(time.UTC),
						evaluatedControlID: defaults.DefaultEUCSControlID,
					},
				},
			},
			wantRes: &evaluation.StopEvaluationResponse{},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				UnimplementedEvaluationServer: tt.fields.UnimplementedEvaluationServer,
				orchestratorClient:            tt.fields.orchestratorClient,
				orchestratorAddress:           tt.fields.orchestratorAddress,
				authorizer:                    tt.fields.authorizer,
				evaluation:                    tt.fields.evaluation,
				results:                       tt.fields.results,
				storage:                       tt.fields.storage,
			}

			// Start the scheduler
			if tt.args.schedulerRunning == true {
				if tt.args.tag {
					j, err := s.evaluation[tt.args.schedulerTag].scheduler.Every(1).Day().Tag(tt.args.schedulerTag).Do(func() { log.Debug("Scheduler") })
					require.NoError(t, err)

					s.evaluation[tt.args.schedulerTag].scheduler.Job(j)
				}
				s.evaluation[tt.args.schedulerTag].scheduler.StartAsync()
			}

			gotRes, err := s.StopEvaluation(tt.args.in0, tt.args.req)
			tt.wantErr(t, err)
			assert.Equal(t, tt.wantRes, gotRes)

			if err == nil {
				assert.Empty(t, s.evaluation[tt.args.schedulerTag])
			}
		})
	}
}

func TestService_initOrchestratorClient(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestratorClient            orchestrator.OrchestratorClient
		orchestratorAddress           grpcTarget
		authorizer                    api.Authorizer
		evaluation                    map[string]*EvaluationScheduler
		results                       map[string]*evaluation.EvaluationResult
		storage                       persistence.Storage
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "OrchestratorClient already exists",
			fields: fields{
				orchestratorClient: orchestrator.NewOrchestratorClient(&grpc.ClientConn{}),
				orchestratorAddress: grpcTarget{
					target: defaults.DefaultOrchestratorAddress,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path",
			fields: fields{
				orchestratorAddress: grpcTarget{
					target: defaults.DefaultOrchestratorAddress,
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				UnimplementedEvaluationServer: tt.fields.UnimplementedEvaluationServer,
				orchestratorClient:            tt.fields.orchestratorClient,
				orchestratorAddress:           tt.fields.orchestratorAddress,
				authorizer:                    tt.fields.authorizer,
				evaluation:                    tt.fields.evaluation,
				results:                       tt.fields.results,
				storage:                       tt.fields.storage,
			}

			err := s.initOrchestratorClient()
			if tt.wantErr != nil {
				tt.wantErr(t, err)
			}
		})
	}
}

func TestService_getMetricsFromSubControls(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestratorClient            orchestrator.OrchestratorClient
		orchestratorAddress           grpcTarget
		authorizer                    api.Authorizer
		evaluation                    map[string]*EvaluationScheduler
		results                       map[string]*evaluation.EvaluationResult
		storage                       persistence.Storage
	}
	type args struct {
		control *orchestrator.Control
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantMetrics []*assessment.Metric
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "Sub-control level and metrics missing",
			args: args{
				control: &orchestrator.Control{
					Id:                "testId",
					CategoryName:      defaults.DefaultEUCSCategoryName,
					CategoryCatalogId: defaults.DefaultEUCSControlID,
					Name:              "testId",
					Controls:          nil,
					Metrics:           nil,
				},
			},
			wantMetrics: nil,
			wantErr:     assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				UnimplementedEvaluationServer: tt.fields.UnimplementedEvaluationServer,
				orchestratorClient:            tt.fields.orchestratorClient,
				orchestratorAddress:           tt.fields.orchestratorAddress,
				authorizer:                    tt.fields.authorizer,
				evaluation:                    tt.fields.evaluation,
				results:                       tt.fields.results,
				storage:                       tt.fields.storage,
			}
			gotMetrics, err := s.getMetricsFromSubControls(tt.args.control)
			if tt.wantErr != nil {
				tt.wantErr(t, err)
			}
			assert.Equal(t, tt.wantMetrics, gotMetrics)
		})
	}
}
