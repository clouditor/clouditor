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
	"os"
	"reflect"
	"testing"
	"time"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evaluation"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/defaults"
	"clouditor.io/clouditor/internal/testutil/clitest"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/service"
	"github.com/go-co-op/gocron"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"
)

func TestMain(m *testing.M) {
	clitest.AutoChdir()

	server, _ := startBufConnServer()

	code := m.Run()

	server.Stop()

	os.Exit(code)
}

func TestNewService(t *testing.T) {

	type args struct {
		opts []service.Option[Service]
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{
			name: "WithOrchestratorAddress",
			args: args{
				opts: []service.Option[Service]{service.Option[Service](WithOrchestratorAddress(defaults.DefaultOrchestratorAddress))},
			},
			want: &Service{
				orchestratorAddress: grpcTarget{
					target: defaults.DefaultOrchestratorAddress,
				},
				results:   make(map[string]*evaluation.EvaluationResult),
				scheduler: gocron.NewScheduler(time.UTC),
				wg:        make(map[string]*WaitGroup),
			},
		},
		{
			name: "WithOAuth2Authorizer",
			args: args{
				opts: []service.Option[Service]{service.Option[Service](WithOAuth2Authorizer(&clientcredentials.Config{}))},
			},
			want: &Service{
				orchestratorAddress: grpcTarget{
					target: "localhost:9090",
				},
				results:    make(map[string]*evaluation.EvaluationResult),
				scheduler:  gocron.NewScheduler(time.UTC),
				authorizer: api.NewOAuthAuthorizerFromClientCredentials(&clientcredentials.Config{}),
				wg:         make(map[string]*WaitGroup),
			},
		},
		{
			name: "WithAuthorizer",
			args: args{
				opts: []service.Option[Service]{service.Option[Service](WithAuthorizer(api.NewOAuthAuthorizerFromClientCredentials(&clientcredentials.Config{})))},
			},
			want: &Service{
				orchestratorAddress: grpcTarget{
					target: "localhost:9090",
				},
				results:    make(map[string]*evaluation.EvaluationResult),
				scheduler:  gocron.NewScheduler(time.UTC),
				authorizer: api.NewOAuthAuthorizerFromClientCredentials(&clientcredentials.Config{}),
				wg:         make(map[string]*WaitGroup),
			},
		},
		{
			name: "Happy path",
			args: args{
				opts: []service.Option[Service]{},
			},
			want: &Service{
				orchestratorAddress: grpcTarget{
					target: "localhost:9090",
				},
				scheduler: gocron.NewScheduler(time.UTC),
				results:   make(map[string]*evaluation.EvaluationResult),
				wg:        make(map[string]*WaitGroup),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewService(tt.args.opts...)

			// Check if scheduler ist initialized and then remove
			if tt.want.scheduler != nil {
				assert.NotEmpty(t, tt.want.scheduler)
				tt.want.scheduler = nil
				got.scheduler = nil
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			}
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
				controlId: defaults.DefaultEUCSLowerLevelControlID137,
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
				controlId:      defaults.DefaultEUCSLowerLevelControlID137,
			},
			want: fmt.Sprintf("%s-%s", defaults.DefaultTargetCloudServiceID, defaults.DefaultEUCSLowerLevelControlID137),
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

func TestService_ListEvaluationResults(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestratorClient            orchestrator.OrchestratorClient
		orchestratorAddress           grpcTarget
		authorizer                    api.Authorizer
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
			name: "Filter cloud service id",
			fields: fields{results: map[string]*evaluation.EvaluationResult{
				"11111111-1111-1111-1111-111111111111": {
					Id: "11111111-1111-1111-1111-111111111111",
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CloudServiceId: "00000000-0000-0000-0000-000000000000",
					},
				},
				"22222222-2222-2222-2222-222222222222": {
					Id: "22222222-2222-2222-2222-222222222222",
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CloudServiceId: "99999999-9999-9999-9999-999999999999",
					},
				},
				"33333333-3333-3333-3333-333333333333": {
					Id: "33333333-3333-3333-3333-333333333333",
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CloudServiceId: "99999999-9999-9999-9999-999999999999",
					},
				},
				"44444444-4444-4444-4444-444444444444": {
					Id: "44444444-4444-4444-4444-444444444444",
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CloudServiceId: "00000000-0000-0000-0000-000000000000",
					},
				},
			}},
			args: args{
				in0: context.Background(),
				req: &evaluation.ListEvaluationResultsRequest{
					FilteredCloudServiceId: util.Ref("00000000-0000-0000-0000-000000000000"),
				},
			},
			wantRes: &evaluation.ListEvaluationResultsResponse{
				Results: []*evaluation.EvaluationResult{
					{
						Id: "11111111-1111-1111-1111-111111111111",
						TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
							CloudServiceId: "00000000-0000-0000-0000-000000000000",
						},
					},
					{
						Id: "44444444-4444-4444-4444-444444444444",
						TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
							CloudServiceId: "00000000-0000-0000-0000-000000000000",
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Filter control id",
			fields: fields{results: map[string]*evaluation.EvaluationResult{
				"11111111-1111-1111-1111-111111111111": {
					Id:        "11111111-1111-1111-1111-111111111111",
					ControlId: "control 1",
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CloudServiceId: "00000000-0000-0000-0000-000000000000",
					},
				},
				"22222222-2222-2222-2222-222222222222": {
					Id:        "22222222-2222-2222-2222-222222222222",
					ControlId: "control 1",
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CloudServiceId: "99999999-9999-9999-9999-999999999999",
					},
				},
				"33333333-3333-3333-3333-333333333333": {
					Id:        "33333333-3333-3333-3333-333333333333",
					ControlId: "control 2",
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CloudServiceId: "99999999-9999-9999-9999-999999999999",
					},
				},
				"44444444-4444-4444-4444-444444444444": {
					Id:        "44444444-4444-4444-4444-444444444444",
					ControlId: "control 2",
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CloudServiceId: "00000000-0000-0000-0000-000000000000",
					},
				},
			}},
			args: args{
				in0: context.Background(),
				req: &evaluation.ListEvaluationResultsRequest{
					FilteredControlId: util.Ref("control 1"),
				},
			},
			wantRes: &evaluation.ListEvaluationResultsResponse{
				Results: []*evaluation.EvaluationResult{
					{
						Id:        "11111111-1111-1111-1111-111111111111",
						ControlId: "control 1",
						TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
							CloudServiceId: "00000000-0000-0000-0000-000000000000",
						},
					},
					{
						Id:        "22222222-2222-2222-2222-222222222222",
						ControlId: "control 1",
						TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
							CloudServiceId: "99999999-9999-9999-9999-999999999999",
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Filter cloud service id and control id",
			fields: fields{results: map[string]*evaluation.EvaluationResult{
				"11111111-1111-1111-1111-111111111111": {
					Id:        "11111111-1111-1111-1111-111111111111",
					ControlId: "control 1",
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CloudServiceId: "00000000-0000-0000-0000-000000000000",
					},
				},
				"22222222-2222-2222-2222-222222222222": {
					Id:        "22222222-2222-2222-2222-222222222222",
					ControlId: "control 1",
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CloudServiceId: "99999999-9999-9999-9999-999999999999",
					},
				},
				"33333333-3333-3333-3333-333333333333": {
					Id:        "33333333-3333-3333-3333-333333333333",
					ControlId: "control 2",
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CloudServiceId: "99999999-9999-9999-9999-999999999999",
					},
				},
				"44444444-4444-4444-4444-444444444444": {
					Id:        "44444444-4444-4444-4444-444444444444",
					ControlId: "control 2",
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CloudServiceId: "00000000-0000-0000-0000-000000000000",
					},
				},
			}},
			args: args{
				in0: context.Background(),
				req: &evaluation.ListEvaluationResultsRequest{
					FilteredCloudServiceId: util.Ref("00000000-0000-0000-0000-000000000000"),
					FilteredControlId:      util.Ref("control 1"),
				},
			},
			wantRes: &evaluation.ListEvaluationResultsResponse{
				Results: []*evaluation.EvaluationResult{
					{
						Id:        "11111111-1111-1111-1111-111111111111",
						ControlId: "control 1",
						TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
							CloudServiceId: "00000000-0000-0000-0000-000000000000",
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
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
				results:                       tt.fields.results,
				storage:                       tt.fields.storage,
			}
			gotRes, err := s.ListEvaluationResults(tt.args.in0, tt.args.req)

			tt.wantErr(t, err)
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("ListEvaluationResults() gotResp = %v, want %v", gotRes, tt.wantRes)
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
			name:        "Control is missing",
			args:        args{},
			wantMetrics: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "control is missing")
			},
		},
		{
			name: "Sub-control level and metrics missing",
			args: args{
				control: &orchestrator.Control{
					Id:                "testId",
					CategoryName:      defaults.DefaultEUCSCategoryName,
					CategoryCatalogId: defaults.DefaultEUCSLowerLevelControlID137,
					Name:              "testId",
					Controls:          nil,
					Metrics:           nil,
				},
			},
			wantMetrics: nil,
			wantErr:     assert.NoError,
		},
		{
			name: "Error getting control from orchestrator",
			fields: fields{
				orchestratorAddress: grpcTarget{target: DefaultOrchestratorAddress},
			},
			args: args{
				control: &orchestrator.Control{
					Id:                "testId",
					CategoryName:      defaults.DefaultEUCSCategoryName,
					CategoryCatalogId: defaults.DefaultEUCSLowerLevelControlID137,
					Name:              "testId",
					Controls: []*orchestrator.Control{
						{
							Id:                "testId-subcontrol",
							CategoryName:      defaults.DefaultEUCSCategoryName,
							CategoryCatalogId: defaults.DefaultEUCSLowerLevelControlID137,
							Name:              "testId-subcontrol",
						},
					},
					Metrics: []*assessment.Metric{},
				},
			},
			wantMetrics: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "connection refused")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				UnimplementedEvaluationServer: tt.fields.UnimplementedEvaluationServer,
				orchestratorClient:            tt.fields.orchestratorClient,
				orchestratorAddress:           tt.fields.orchestratorAddress,
				authorizer:                    tt.fields.authorizer,
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

func TestService_StopEvaluation(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestratorClient            orchestrator.OrchestratorClient
		orchestratorAddress           grpcTarget
		authorizer                    api.Authorizer
		scheduler                     *gocron.Scheduler
		results                       map[string]*evaluation.EvaluationResult
		storage                       persistence.Storage
	}
	type args struct {
		in0              context.Context
		req              *evaluation.StopEvaluationRequest
		schedulerTag     string
		schedulerRunning bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *evaluation.StopEvaluationResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "ToE in request is missing",
			args: args{
				in0: context.Background(),
				req: &evaluation.StopEvaluationRequest{},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid StopEvaluationRequest.TargetOfEvaluation: value is required")
			},
		},
		// {
		// 	name: "Happy path",
		// 	args: args{
		// 		in0: context.Background(),
		// 		req: &evaluation.StopEvaluationRequest{
		// 			TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
		// 				CloudServiceId: defaults.DefaultTargetCloudServiceID,
		// 				CatalogId:      defaults.DefaultCatalogID,
		// 				AssuranceLevel: &defaults.AssuranceLevelHigh,
		// 			},
		// 			ControlId:    defaults.DefaultEUCSLowerLevelControlID137,
		// 			CategoryName: defaults.DefaultEUCSCategoryName,
		// 		},
		// 		schedulerTag:     createSchedulerTag(defaults.DefaultTargetCloudServiceID, defaults.DefaultEUCSLowerLevelControlID137),
		// 		schedulerRunning: true,
		// 	},
		// 	fields: fields{
		// 		scheduler: gocron.NewScheduler(time.UTC),
		// 	},
		// 	wantRes: &evaluation.StopEvaluationResponse{},
		// 	wantErr: assert.NoError,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				UnimplementedEvaluationServer: tt.fields.UnimplementedEvaluationServer,
				orchestratorClient:            tt.fields.orchestratorClient,
				orchestratorAddress:           tt.fields.orchestratorAddress,
				authorizer:                    tt.fields.authorizer,
				scheduler:                     tt.fields.scheduler,
				results:                       tt.fields.results,
				storage:                       tt.fields.storage,
			}

			// Start the scheduler
			if tt.args.schedulerRunning == true {
				_, err := s.scheduler.Every(1).Day().Tag(tt.args.schedulerTag).Do(func() { fmt.Println("Scheduler") })
				require.NoError(t, err)

			}

			gotRes, err := s.StopEvaluation(tt.args.in0, tt.args.req)

			tt.wantErr(t, err)
			assert.Equal(t, tt.wantRes, gotRes)
		})
	}
}

// func TestService_evaluateFirstLevelControl(t *testing.T) {
// 	type fields struct {
// 		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
// 		orchestratorClient            orchestrator.OrchestratorClient
// 		orchestratorAddress           grpcTarget
// 		authorizer                    api.Authorizer
// 		scheduler                     *gocron.Scheduler
// 		results                       map[string]*evaluation.EvaluationResult
// 		resultMutex                   sync.Mutex
// 		storage                       persistence.Storage
// 	}
// 	type args struct {
// 		toe          *orchestrator.TargetOfEvaluation
// 		categoryName string
// 		controlId    string
// 	}
// 	tests := []struct {
// 		name        string
// 		fields      fields
// 		args        args
// 		wantResults []*evaluation.EvaluationResult
// 	}{
// 		// TODO: TBD
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := &Service{
// 				UnimplementedEvaluationServer: tt.fields.UnimplementedEvaluationServer,
// 				orchestratorClient:            tt.fields.orchestratorClient,
// 				orchestratorAddress:           tt.fields.orchestratorAddress,
// 				authorizer:                    tt.fields.authorizer,
// 				scheduler:                     tt.fields.scheduler,
// 				results:                       tt.fields.results,
// 				resultMutex:                   tt.fields.resultMutex,
// 				storage:                       tt.fields.storage,
// 			}
// 			s.evaluateFirstLevelControl(tt.args.toe, tt.args.categoryName, tt.args.controlId)

// 			assert.Equal(t, tt.wantResults, s.results)
// 		})
// 	}
// }

func Test_getAllControlIdsFromControl(t *testing.T) {
	type args struct {
		control *orchestrator.Control
	}
	tests := []struct {
		name           string
		args           args
		wantControlIds []string
	}{
		{
			name:           "Control is missing",
			wantControlIds: []string{},
		},
		{
			name: "Control has no sub-controls",
			args: args{
				&orchestrator.Control{
					Id:                "testId-1",
					CategoryName:      defaults.DefaultEUCSCategoryName,
					CategoryCatalogId: defaults.DefaultCatalogID,
					Name:              "testId-1",
					Description:       "test test test",
					Controls:          []*orchestrator.Control{},
				},
			},
			wantControlIds: []string{"testId-1"},
		},
		{
			name: "Happy path",
			args: args{
				&orchestrator.Control{
					Id:                "testId-1",
					CategoryName:      defaults.DefaultEUCSCategoryName,
					CategoryCatalogId: defaults.DefaultCatalogID,
					Name:              "testId-1",
					Description:       "test test test",
					Controls: []*orchestrator.Control{
						{
							Id: "testId-1.1",
						},
						{
							Id: "testId-1.2",
						},
						{
							Id: "testId-1.3",
						},
					},
				},
			},
			wantControlIds: []string{"testId-1", "testId-1.1", "testId-1.2", "testId-1.3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotControlIds := getAllControlIdsFromControl(tt.args.control); !reflect.DeepEqual(gotControlIds, tt.wantControlIds) {
				t.Errorf("getAllControlIdsFromControl() = %v, want %v", gotControlIds, tt.wantControlIds)
			}
		})
	}
}

func TestService_StartEvaluation(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestratorClient            orchestrator.OrchestratorClient
		orchestratorAddress           grpcTarget
		authorizer                    api.Authorizer
		scheduler                     *gocron.Scheduler
		results                       map[string]*evaluation.EvaluationResult
		storage                       persistence.Storage
		wg                            map[string]*WaitGroup
	}
	type args struct {
		in0              context.Context
		req              *evaluation.StartEvaluationRequest
		schedulerTag     string
		schedulerRunning bool
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp *evaluation.StartEvaluationResponse
		wantErr  assert.ErrorAssertionFunc
	}{
		// {
		// 	name: "Start scheduler job for one control",
		// 	fields: fields{
		// 		scheduler:           gocron.NewScheduler(time.UTC),
		// 		wg:                  make(map[string]*WaitGroup),
		// 		results:             make(map[string]*evaluation.EvaluationResult),
		// 		orchestratorAddress: grpcTarget{target: DefaultOrchestratorAddress},
		// 	},
		// 	args: args{
		// 		in0: context.Background(),
		// 		req: &evaluation.StartEvaluationRequest{
		// 			TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
		// 				CloudServiceId: defaults.DefaultTargetCloudServiceID,
		// 				CatalogId:      defaults.DefaultCatalogID,
		// 				AssuranceLevel: &defaults.AssuranceLevelHigh,
		// 				ControlsInScope: []*orchestrator.Control{
		// 					{
		// 						Id:                defaults.DefaultEUCSUpperLevelControlID13,
		// 						CategoryName:      defaults.DefaultEUCSCategoryName,
		// 						CategoryCatalogId: defaults.DefaultCatalogID,
		// 						Name:              defaults.DefaultEUCSUpperLevelControlID13,
		// 					},
		// 				},
		// 			},
		// 		},
		// 		schedulerRunning: false,
		// 		schedulerTag:     createSchedulerTag(defaults.DefaultTargetCloudServiceID, defaults.DefaultEUCSUpperLevelControlID13),
		// 	},
		// 	wantResp: &evaluation.StartEvaluationResponse{Status: true},
		// 	wantErr:  assert.NoError,
		// },
		// {
		// 	name: "Scheduler job for one control already running",
		// 	fields: fields{
		// 		scheduler: gocron.NewScheduler(time.UTC),
		// 		wg: map[string]*WaitGroup{
		// 			fmt.Sprintf("%s-%s", defaults.DefaultTargetCloudServiceID, defaults.DefaultEUCSUpperLevelControlID13): {
		// 				wg:      &sync.WaitGroup{},
		// 				wgMutex: sync.Mutex{},
		// 			},
		// 		},
		// 		results: make(map[string]*evaluation.EvaluationResult),
		// 	},
		// 	args: args{
		// 		in0: context.Background(),
		// 		req: &evaluation.StartEvaluationRequest{
		// 			TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
		// 				CloudServiceId: defaults.DefaultTargetCloudServiceID,
		// 				CatalogId:      defaults.DefaultCatalogID,
		// 				AssuranceLevel: &defaults.AssuranceLevelHigh,
		// 				ControlsInScope: []*orchestrator.Control{
		// 					{
		// 						Id:                defaults.DefaultEUCSUpperLevelControlID13,
		// 						CategoryName:      defaults.DefaultEUCSCategoryName,
		// 						CategoryCatalogId: defaults.DefaultCatalogID,
		// 						Name:              defaults.DefaultEUCSUpperLevelControlID13,
		// 					},
		// 				},
		// 			},
		// 		},
		// 		schedulerRunning: true,
		// 		schedulerTag:     createSchedulerTag(defaults.DefaultTargetCloudServiceID, defaults.DefaultEUCSUpperLevelControlID13),
		// 	},
		// 	wantResp: &evaluation.StartEvaluationResponse{
		// 		Status:        false,
		// 		StatusMessage: fmt.Sprintf("evaluation for cloud service id '%s' and control id '%s' already started", defaults.DefaultTargetCloudServiceID, defaults.DefaultEUCSUpperLevelControlID13),
		// 	},
		// 	wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
		// 		return assert.ErrorContains(t, err, fmt.Sprintf("evaluation for cloud service id '%s' and control id '%s' already started", defaults.DefaultTargetCloudServiceID, defaults.DefaultEUCSUpperLevelControlID13))
		// 	},
		// },
		{
			name: "ToE missing in request",
			args: args{
				in0:              context.Background(),
				schedulerRunning: false,
				req:              &evaluation.StartEvaluationRequest{},
			},
			wantResp: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid StartEvaluationRequest.TargetOfEvaluation: value is required")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				UnimplementedEvaluationServer: tt.fields.UnimplementedEvaluationServer,
				orchestratorClient:            tt.fields.orchestratorClient,
				orchestratorAddress:           tt.fields.orchestratorAddress,
				authorizer:                    tt.fields.authorizer,
				scheduler:                     tt.fields.scheduler,
				results:                       tt.fields.results,
				storage:                       tt.fields.storage,
				wg:                            tt.fields.wg,
			}

			// Start the scheduler if needed
			if tt.args.schedulerRunning {
				_, err := s.scheduler.Every(1).Day().Tag(tt.args.schedulerTag).Do(func() { fmt.Println("Scheduler") })
				require.NoError(t, err)
			}

			gotResp, err := s.StartEvaluation(tt.args.in0, tt.args.req)
			tt.wantErr(t, err)
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("StartEvaluation() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func Test_getSchedulerTagsForControlIds(t *testing.T) {
	type args struct {
		controlIds     []string
		cloudServiceId string
	}
	tests := []struct {
		name              string
		args              args
		wantSchedulerTags []string
	}{
		{
			name: "Empty cloud service id",
			args: args{
				controlIds: []string{"OPS-13", "OPS-13.2", "OPS-13.3"},
			},
			wantSchedulerTags: []string{},
		},
		{
			name: "Empty control ids",
			args: args{
				cloudServiceId: defaults.DefaultTargetCloudServiceID,
			},
			wantSchedulerTags: []string{},
		},
		{
			name: "Happy path",
			args: args{
				controlIds:     []string{"OPS-13", "OPS-13.2", "OPS-13.3"},
				cloudServiceId: defaults.DefaultTargetCloudServiceID,
			},
			wantSchedulerTags: []string{
				fmt.Sprintf("%s-%s", defaults.DefaultTargetCloudServiceID, "OPS-13"),
				fmt.Sprintf("%s-%s", defaults.DefaultTargetCloudServiceID, "OPS-13.2"),
				fmt.Sprintf("%s-%s", defaults.DefaultTargetCloudServiceID, "OPS-13.3")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSchedulerTags := getSchedulerTagsForControlIds(tt.args.controlIds, tt.args.cloudServiceId); !reflect.DeepEqual(gotSchedulerTags, tt.wantSchedulerTags) {
				t.Errorf("getSchedulerTagsForControlIds() = %v, want %v", gotSchedulerTags, tt.wantSchedulerTags)
			}
		})
	}
}

func TestService_getMetrics(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestratorClient            orchestrator.OrchestratorClient
		orchestratorAddress           grpcTarget
		authorizer                    api.Authorizer
		scheduler                     *gocron.Scheduler
		wg                            map[string]*WaitGroup
		results                       map[string]*evaluation.EvaluationResult
		storage                       persistence.Storage
	}
	type args struct {
		catalogId    string
		categoryName string
		controlId    string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantMetrics []*assessment.Metric
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name:        "Input empty",
			wantMetrics: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not get control for control id")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				UnimplementedEvaluationServer: tt.fields.UnimplementedEvaluationServer,
				orchestratorClient:            tt.fields.orchestratorClient,
				orchestratorAddress:           tt.fields.orchestratorAddress,
				authorizer:                    tt.fields.authorizer,
				scheduler:                     tt.fields.scheduler,
				wg:                            tt.fields.wg,
				results:                       tt.fields.results,
				storage:                       tt.fields.storage,
			}
			gotMetrics, err := s.getMetrics(tt.args.catalogId, tt.args.categoryName, tt.args.controlId)
			tt.wantErr(t, err)
			assert.Equal(t, tt.wantMetrics, gotMetrics)
		})
	}
}

func Test_getMetricIds(t *testing.T) {
	type args struct {
		metrics []*assessment.Metric
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Empty input",
			args: args{},
			want: nil,
		},
		{
			name: "Happy path",
			args: args{
				metrics: []*assessment.Metric{
					{
						Id: defaults.DefaultEUCSLowerLevelControlID136,
					},
					{
						Id: defaults.DefaultEUCSLowerLevelControlID137,
					},
				},
			},
			want: []string{defaults.DefaultEUCSLowerLevelControlID136, defaults.DefaultEUCSLowerLevelControlID137},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMetricIds(tt.args.metrics); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getMetricIds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_controlContains(t *testing.T) {
	type args struct {
		controls  []*orchestrator.Control
		controlId string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty input",
			args: args{},
			want: false,
		},
		{
			name: "Contols list does not contain control ID",
			args: args{
				controls: []*orchestrator.Control{
					{
						Id: defaults.DefaultEUCSLowerLevelControlID137,
					},
				},
				controlId: defaults.DefaultEUCSLowerLevelControlID136,
			},
			want: false,
		},
		{
			name: "Contols list contains control ID",
			args: args{
				controls: []*orchestrator.Control{
					{
						Id: defaults.DefaultEUCSLowerLevelControlID136,
					},
				},
				controlId: defaults.DefaultEUCSLowerLevelControlID136,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := controlContains(tt.args.controls, tt.args.controlId); got != tt.want {
				t.Errorf("controlContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_stopSchedulerJobs(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestratorClient            orchestrator.OrchestratorClient
		orchestratorAddress           grpcTarget
		authorizer                    api.Authorizer
		scheduler                     *gocron.Scheduler
		wg                            map[string]*WaitGroup
		results                       map[string]*evaluation.EvaluationResult
		storage                       persistence.Storage
		schedulerTags                 []string
		schedulerRunning              bool
	}
	type args struct {
		schedulerTags []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "Empty input",
			args:    args{},
			wantErr: assert.NoError,
		},
		{
			name: "Stop not existing scheduler job",
			args: args{
				schedulerTags: []string{fmt.Sprintf("%s-%s", defaults.DefaultTargetCloudServiceID, defaults.DefaultEUCSLowerLevelControlID137)},
			},
			fields: fields{
				scheduler: gocron.NewScheduler(time.UTC),
				schedulerTags: []string{
					fmt.Sprintf("%s-%s", defaults.DefaultTargetCloudServiceID, defaults.DefaultEUCSLowerLevelControlID136)},
				schedulerRunning: true,
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "error when removing job for tag")
			},
		},
		{
			name: "Stopping two scheduler jobs",
			args: args{
				schedulerTags: []string{fmt.Sprintf("%s-%s", defaults.DefaultTargetCloudServiceID, defaults.DefaultEUCSLowerLevelControlID137)},
			},
			fields: fields{
				scheduler:        gocron.NewScheduler(time.UTC),
				schedulerRunning: true,
				schedulerTags: []string{
					fmt.Sprintf("%s-%s", defaults.DefaultTargetCloudServiceID, defaults.DefaultEUCSLowerLevelControlID136),
					fmt.Sprintf("%s-%s", defaults.DefaultTargetCloudServiceID, defaults.DefaultEUCSLowerLevelControlID137)},
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
				scheduler:                     tt.fields.scheduler,
				wg:                            tt.fields.wg,
				results:                       tt.fields.results,
				storage:                       tt.fields.storage,
			}

			// Start the scheduler
			if tt.fields.schedulerRunning == true {
				for _, tag := range tt.fields.schedulerTags {
					_, err := s.scheduler.Every(1).Day().Tag(tag).Do(func() { fmt.Println("Scheduler") })
					require.NoError(t, err)
				}
			}
			err := s.stopSchedulerJobs(tt.args.schedulerTags)
			tt.wantErr(t, err)

			if tt.fields.schedulerRunning {
				jobs := s.scheduler.Jobs()
				assert.NotContains(t, jobs, tt.args.schedulerTags)
			}
		})
	}
}

func TestService_stopSchedulerJob(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestratorClient            orchestrator.OrchestratorClient
		orchestratorAddress           grpcTarget
		authorizer                    api.Authorizer
		scheduler                     *gocron.Scheduler
		wg                            map[string]*WaitGroup
		results                       map[string]*evaluation.EvaluationResult
		storage                       persistence.Storage
	}
	type args struct {
		schedulerTag     string
		schedulerRunning bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Empty input",
			args: args{
				schedulerRunning: false,
			},
			fields: fields{
				scheduler: gocron.NewScheduler(time.UTC),
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "no jobs found with given tag")
			},
		},
		{
			name: "Happy path",
			args: args{
				schedulerRunning: true,
				schedulerTag:     fmt.Sprintf("%s-%s", defaults.DefaultTargetCloudServiceID, defaults.DefaultEUCSLowerLevelControlID137),
			},
			fields: fields{
				scheduler: gocron.NewScheduler(time.UTC),
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
				scheduler:                     tt.fields.scheduler,
				wg:                            tt.fields.wg,
				results:                       tt.fields.results,
				storage:                       tt.fields.storage,
			}

			// Start the scheduler
			if tt.args.schedulerRunning == true {
				_, err := s.scheduler.Every(1).Day().Tag(tt.args.schedulerTag).Do(func() { fmt.Println("Scheduler") })
				require.NoError(t, err)

			}

			err := s.stopSchedulerJob(tt.args.schedulerTag)
			tt.wantErr(t, err)
		})
	}
}
