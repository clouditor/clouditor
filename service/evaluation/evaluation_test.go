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
	"sync"
	"testing"
	"time"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evaluation"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/internal/testutil/clitest"
	"clouditor.io/clouditor/internal/testutil/evaluationtest"
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
				opts: []service.Option[Service]{service.Option[Service](WithOrchestratorAddress("localhost:1234"))},
			},
			want: &Service{
				orchestratorAddress: grpcTarget{
					target: "localhost:1234",
				},
				wg: make(map[string]*WaitGroup),
			},
		},
		{
			name: "WithOAuth2Authorizer",
			args: args{
				opts: []service.Option[Service]{service.Option[Service](WithOAuth2Authorizer(&clientcredentials.Config{}))},
			},
			want: &Service{
				orchestratorAddress: grpcTarget{
					target: testdata.MockOrchestratorAddress,
				},
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
					target: testdata.MockOrchestratorAddress,
				},
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
					target: testdata.MockOrchestratorAddress,
				},
				wg: make(map[string]*WaitGroup),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewService(tt.args.opts...)

			// Check if scheduler is initialized and then remove
			if got.scheduler != nil {
				assert.NotEmpty(t, got.scheduler)
				tt.want.scheduler = nil
				got.scheduler = nil
			}
			// Check if storage is initialized and then remove
			if got.storage != nil {
				assert.NotNil(t, got.storage)
				tt.want.storage = nil
				got.storage = nil
			}
			// Check if authz is initialized and then remove
			if got.authz != nil {
				assert.NotNil(t, got.authz)
				tt.want.authz = nil
				got.authz = nil
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
				controlId: testdata.MockSubControlID,
			},
			want: "",
		},
		{
			name: "Input cloudServiceId empty",
			args: args{
				cloudServiceId: testdata.MockCloudServiceID,
			},
			want: "",
		},
		{
			name: "Happy path",
			args: args{
				cloudServiceId: testdata.MockCloudServiceID,
				controlId:      testdata.MockSubControlID,
			},
			want: fmt.Sprintf("%s-%s", testdata.MockCloudServiceID, testdata.MockSubControlID),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSchedulerTag(tt.args.cloudServiceId, tt.args.controlId); got != tt.want {
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
		storage                       persistence.Storage
		authz                         service.AuthorizationStrategy
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
			name: "Missing request",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(evaluationtest.MockEvaluationResults))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				in0: context.Background(),
				req: nil,
			},
			wantRes: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.ErrorContains(t, err, api.ErrEmptyRequest.Error())
				return false
			},
		},
		{
			name: "Filter latest_by_resource_id, control_id, sub_controls, cloud_service_id",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(evaluationtest.MockEvaluationResults))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.ListEvaluationResultsRequest{
					LatestByResourceId:     util.Ref(true),
					FilteredControlId:      util.Ref(testdata.MockSubControlID11),
					FilteredSubControls:    util.Ref(testdata.MockControlID1),
					FilteredCloudServiceId: util.Ref(testdata.MockCloudServiceID),
				},
			},
			wantRes: &evaluation.ListEvaluationResultsResponse{
				Results: []*evaluation.EvaluationResult{
					evaluationtest.MockEvaluationResult22,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Filter latest_by_resource_id and control_id",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(evaluationtest.MockEvaluationResults))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.ListEvaluationResultsRequest{
					LatestByResourceId: util.Ref(true),
					FilteredControlId:  util.Ref(testdata.MockSubControlID11),
				},
			},
			wantRes: &evaluation.ListEvaluationResultsResponse{
				Results: []*evaluation.EvaluationResult{
					evaluationtest.MockEvaluationResult8,
					evaluationtest.MockEvaluationResult22,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Filter latest_by_resource_id",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(evaluationtest.MockEvaluationResults))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.ListEvaluationResultsRequest{LatestByResourceId: util.Ref(true)},
			},
			wantRes: &evaluation.ListEvaluationResultsResponse{
				Results: []*evaluation.EvaluationResult{
					evaluationtest.MockEvaluationResult7,
					evaluationtest.MockEvaluationResult1,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Filter control_id",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(evaluationtest.MockEvaluationResults))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.ListEvaluationResultsRequest{FilteredControlId: util.Ref(testdata.MockControlID1)},
			},
			wantRes: &evaluation.ListEvaluationResultsResponse{
				Results: []*evaluation.EvaluationResult{
					evaluationtest.MockEvaluationResult1,
					evaluationtest.MockEvaluationResult7,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Filter sub-controls - get all sub-control evaluation results for a given control",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(evaluationtest.MockEvaluationResults))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.ListEvaluationResultsRequest{FilteredSubControls: util.Ref(testdata.MockControlID1)},
			},
			wantRes: &evaluation.ListEvaluationResultsResponse{
				Results: []*evaluation.EvaluationResult{
					evaluationtest.MockEvaluationResult1,
					evaluationtest.MockEvaluationResult2,
					evaluationtest.MockEvaluationResult22,
					evaluationtest.MockEvaluationResult3,
					evaluationtest.MockEvaluationResult7,
					evaluationtest.MockEvaluationResult8,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Filter cloud_service_id",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(evaluationtest.MockEvaluationResults))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.ListEvaluationResultsRequest{FilteredCloudServiceId: util.Ref(testdata.MockCloudServiceID)},
			},
			wantRes: &evaluation.ListEvaluationResultsResponse{
				Results: []*evaluation.EvaluationResult{
					evaluationtest.MockEvaluationResult1,
					evaluationtest.MockEvaluationResult2,
					evaluationtest.MockEvaluationResult22,
					evaluationtest.MockEvaluationResult3,
					evaluationtest.MockEvaluationResult4,
					evaluationtest.MockEvaluationResult5,
					evaluationtest.MockEvaluationResult6,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "multiple page result - first page",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(evaluationtest.MockEvaluationResults))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.ListEvaluationResultsRequest{
					PageSize: 2,
				},
			},
			wantRes: &evaluation.ListEvaluationResultsResponse{
				Results: []*evaluation.EvaluationResult{
					evaluationtest.MockEvaluationResult1,
					evaluationtest.MockEvaluationResult2,
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
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(evaluationtest.MockEvaluationResults))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.ListEvaluationResultsRequest{
					PageSize: 6,
					PageToken: func() string {
						token, _ := (&api.PageToken{Start: 6, Size: 4}).Encode()
						return token
					}(),
				},
			},
			wantRes: &evaluation.ListEvaluationResultsResponse{
				Results: []*evaluation.EvaluationResult{
					evaluationtest.MockEvaluationResult6,
					evaluationtest.MockEvaluationResult7,
					evaluationtest.MockEvaluationResult8,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "List all results",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(evaluationtest.MockEvaluationResults))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.ListEvaluationResultsRequest{},
			},
			wantRes: &evaluation.ListEvaluationResultsResponse{
				Results: evaluationtest.MockEvaluationResults,
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
				storage:                       tt.fields.storage,
				authz:                         tt.fields.authz,
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
					target: testdata.MockOrchestratorAddress,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path",
			fields: fields{
				orchestratorAddress: grpcTarget{
					target: testdata.MockOrchestratorAddress,
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
					CategoryName:      testdata.MockCategoryName,
					CategoryCatalogId: testdata.MockSubControlID,
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
					CategoryName:      testdata.MockCategoryName,
					CategoryCatalogId: testdata.MockSubControlID,
					Name:              "testId",
					Controls: []*orchestrator.Control{
						{
							Id:                "testId-subcontrol",
							CategoryName:      testdata.MockCategoryName,
							CategoryCatalogId: testdata.MockSubControlID,
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
				storage:                       tt.fields.storage,
			}
			gotMetrics, err := s.getMetricsFromSubcontrols(tt.args.control)
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
				return assert.ErrorContains(t, err, "invalid StopEvaluationRequest.CloudServiceId: value must be a valid UUID")
			},
		},
		// {
		// 	name: "Happy path",
		// 	args: args{
		// 		in0: context.Background(),
		// 		req: &evaluation.StopEvaluationRequest{
		// 			TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
		// 				CloudServiceId:  testdata.MockCloudServiceID,
		// 				CatalogId:       testdata.MockCatalogID,
		// 				AssuranceLevel: &testdata.AssuranceLevelHigh,
		// 			},
		// 			ControlId:    defaults.DefaultEUCSLowerLevelControlID137,
		// 			CategoryName:  testdata.MockCategoryName,
		// 		},
		// 		schedulerTag:     createSchedulerTag( testdata.MockCloudServiceID, defaults.DefaultEUCSLowerLevelControlID137),
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
					CategoryName:      testdata.MockCategoryName,
					CategoryCatalogId: testdata.MockCatalogID,
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
					CategoryName:      testdata.MockCategoryName,
					CategoryCatalogId: testdata.MockCatalogID,
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
		// 			CloudServiceId: testdata.MockCloudServiceID,
		// 			CatalogId:      testdata.MockCatalogID,
		// 		},
		// 		schedulerRunning: false,
		// 		schedulerTag:     getSchedulerTag(testdata.MockCloudServiceID, testdata.MockControlID1),
		// 	},

		// 	wantResp: &evaluation.StartEvaluationResponse{Status: true},
		// 	wantErr:  assert.NoError,
		// },
		// {
		// 	name: "Scheduler job for one control already running",
		// 	fields: fields{
		// 		scheduler: gocron.NewScheduler(time.UTC),
		// 		wg: map[string]*WaitGroup{
		// 			fmt.Sprintf("%s-%s",  testdata.MockCloudServiceID, defaults.DefaultEUCSUpperLevelControlID13): {
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
		// 				CloudServiceId:  testdata.MockCloudServiceID,
		// 				CatalogId:       testdata.MockCatalogID,
		// 				AssuranceLevel: &testdata.AssuranceLevelHigh,
		// 				ControlsInScope: []*orchestrator.Control{
		// 					{
		// 						Id:                defaults.DefaultEUCSUpperLevelControlID13,
		// 						CategoryName:       testdata.MockCategoryName,
		// 						CategoryCatalogId:  testdata.MockCatalogID,
		// 						Name:              defaults.DefaultEUCSUpperLevelControlID13,
		// 					},
		// 				},
		// 			},
		// 		},
		// 		schedulerRunning: true,
		// 		schedulerTag:     createSchedulerTag( testdata.MockCloudServiceID, defaults.DefaultEUCSUpperLevelControlID13),
		// 	},
		// 	wantResp: &evaluation.StartEvaluationResponse{
		// 		Status:        false,
		// 		StatusMessage: fmt.Sprintf("evaluation for cloud service id '%s' and control id '%s' already started",  testdata.MockCloudServiceID, defaults.DefaultEUCSUpperLevelControlID13),
		// 	},
		// 	wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
		// 		return assert.ErrorContains(t, err, fmt.Sprintf("evaluation for cloud service id '%s' and control id '%s' already started",  testdata.MockCloudServiceID, defaults.DefaultEUCSUpperLevelControlID13))
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
				return assert.ErrorContains(t, err, "invalid StartEvaluationRequest.CloudServiceId: value must be a valid UUID")
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
				cloudServiceId: testdata.MockCloudServiceID,
			},
			wantSchedulerTags: []string{},
		},
		{
			name: "Happy path",
			args: args{
				controlIds:     []string{"OPS-13", "OPS-13.2", "OPS-13.3"},
				cloudServiceId: testdata.MockCloudServiceID,
			},
			wantSchedulerTags: []string{
				fmt.Sprintf("%s-%s", testdata.MockCloudServiceID, "OPS-13"),
				fmt.Sprintf("%s-%s", testdata.MockCloudServiceID, "OPS-13.2"),
				fmt.Sprintf("%s-%s", testdata.MockCloudServiceID, "OPS-13.3")},
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

func TestService_getAllMetricsFromControl(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestratorClient            orchestrator.OrchestratorClient
		orchestratorAddress           grpcTarget
		authorizer                    api.Authorizer
		scheduler                     *gocron.Scheduler
		wg                            map[string]*WaitGroup
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
		{
			name: "metric not exists",
			args: args{
				catalogId:    "test_catalog_id",
				categoryName: "test_category_id",
				controlId:    "test_control_id",
			},
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
				storage:                       tt.fields.storage,
			}
			gotMetrics, err := s.getAllMetricsFromControl(tt.args.catalogId, tt.args.categoryName, tt.args.controlId)
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
						Id: testdata.MockAnotherSubControlID,
					},
					{
						Id: testdata.MockSubControlID,
					},
				},
			},
			want: []string{testdata.MockAnotherSubControlID, testdata.MockSubControlID},
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

func TestService_stopSchedulerJobs(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestratorClient            orchestrator.OrchestratorClient
		orchestratorAddress           grpcTarget
		authorizer                    api.Authorizer
		scheduler                     *gocron.Scheduler
		wg                            map[string]*WaitGroup
		storage                       persistence.Storage
		schedulerTags                 []string
		schedulerRunning              bool
	}
	type args struct {
		schedulerTags  []string
		cloudServiceId string
		controlId      string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   assert.ValueAssertionFunc
	}{
		{
			name: "Empty input",
			args: args{},
			fields: fields{
				scheduler:        gocron.NewScheduler(time.UTC),
				schedulerRunning: true,
				schedulerTags: []string{
					fmt.Sprintf("%s-%s", testdata.MockCloudServiceID, testdata.MockAnotherSubControlID),
					fmt.Sprintf("%s-%s", testdata.MockCloudServiceID, testdata.MockSubControlID)},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				service, ok := i1.(*Service)
				if !assert.True(tt, ok) {
					return false
				}

				return assert.Equal(t, 2, len(service.scheduler.Jobs()))

			},
		},
		{
			name: "Stop not existing scheduler job",
			args: args{
				schedulerTags:  []string{fmt.Sprintf("%s-%s", testdata.MockCloudServiceID, testdata.MockControlID2)},
				cloudServiceId: testdata.MockCloudServiceID,
				controlId:      testdata.MockControlID2,
			},
			fields: fields{
				scheduler: gocron.NewScheduler(time.UTC),
				schedulerTags: []string{
					fmt.Sprintf("%s-%s", testdata.MockCloudServiceID, testdata.MockControlID1)},
				schedulerRunning: true,
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				service, ok := i1.(*Service)
				if !assert.True(tt, ok) {
					return false
				}

				return assert.Equal(t, 1, len(service.scheduler.Jobs()))

			},
		},
		{
			name: "Stopping two scheduler jobs",
			args: args{
				schedulerTags: []string{
					fmt.Sprintf("%s-%s", testdata.MockCloudServiceID, testdata.MockControlID1),
					fmt.Sprintf("%s-%s", testdata.MockCloudServiceID, testdata.MockSubControlID11)},
			},
			fields: fields{
				scheduler:        gocron.NewScheduler(time.UTC),
				schedulerRunning: true,
				schedulerTags: []string{
					fmt.Sprintf("%s-%s", testdata.MockCloudServiceID, testdata.MockControlID1),
					fmt.Sprintf("%s-%s", testdata.MockCloudServiceID, testdata.MockSubControlID11)},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				service, ok := i1.(*Service)
				if !assert.True(tt, ok) {
					return false
				}

				return assert.Empty(t, len(service.scheduler.Jobs()))
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
				storage:                       tt.fields.storage,
			}

			// Start the scheduler
			if tt.fields.schedulerRunning == true {
				for _, tag := range tt.fields.schedulerTags {
					_, err := s.scheduler.Every(1).Day().Tag(tag).Do(func() { fmt.Println("Scheduler") })
					require.NoError(t, err)
				}
			}

			s.stopSchedulerJobs(tt.args.cloudServiceId, tt.args.controlId, tt.args.schedulerTags)
			tt.want(t, s)
			if tt.fields.schedulerRunning {
				jobs := s.scheduler.Jobs()
				assert.NotContains(t, jobs, tt.args.schedulerTags)
			}
		})
	}
}

func TestService_getControl(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestratorClient            orchestrator.OrchestratorClient
		orchestratorAddress           grpcTarget
		authorizer                    api.Authorizer
		scheduler                     *gocron.Scheduler
		wg                            map[string]*WaitGroup
		storage                       persistence.Storage
		hasOrchestratorStream         bool
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
		wantControl *orchestrator.Control
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name:        "orchestrator address is missing",
			wantControl: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not connect to orchestrator service:")
			},
		},
		{
			name: "catalog_id is missing",
			fields: fields{
				hasOrchestratorStream: true,
			},
			wantControl: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid request: invalid GetControlRequest.CatalogId: value length must be at least 1 runes")
			},
		},
		{
			name: "category_name is missing",
			fields: fields{
				hasOrchestratorStream: true,
			},
			wantControl: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid request: invalid GetControlRequest.CatalogId: value length must be at least 1 runes")
			},
		},
		{
			name: "control_id is missing",
			fields: fields{
				hasOrchestratorStream: true,
			},
			wantControl: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid request: invalid GetControlRequest.CatalogId: value length must be at least 1 runes")
			},
		},
		{
			name: "control does not exist",
			fields: fields{
				hasOrchestratorStream: true,
			},
			args: args{
				catalogId:    "wrong_catalog_id",
				categoryName: "wrong_category_id",
				controlId:    "wrong_control_id",
			},
			wantControl: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "control not found")
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
				storage:                       tt.fields.storage,
			}

			// Mock streams for target services
			if tt.fields.hasOrchestratorStream {
				s.orchestratorAddress.opts = []grpc.DialOption{grpc.WithContextDialer(bufConnDialer)}
			} else {
				s.orchestratorAddress.opts = []grpc.DialOption{grpc.WithContextDialer(nil)}
			}

			gotControl, err := s.getControl(tt.args.catalogId, tt.args.categoryName, tt.args.controlId)
			tt.wantErr(t, err)

			if !reflect.DeepEqual(gotControl, tt.wantControl) {
				t.Errorf("Service.getControl() = %v, want %v", gotControl, tt.wantControl)
			}
		})
	}
}

func TestService_addJobToScheduler(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestratorClient            orchestrator.OrchestratorClient
		orchestratorAddress           grpcTarget
		authorizer                    api.Authorizer
		scheduler                     *gocron.Scheduler
		wg                            map[string]*WaitGroup
		storage                       persistence.Storage
		schedulerRunning              bool
		schedulerTag                  string
	}
	type args struct {
		c                  *orchestrator.Control
		toe                *orchestrator.TargetOfEvaluation
		parentSchedulerTag string
		interval           int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// Not necessary to check if control is empty, because method is only called if a control exists
		{
			name: "Add scheduler job for first level Control",
			fields: fields{
				scheduler: gocron.NewScheduler(time.UTC),
			},
			args: args{
				c: &orchestrator.Control{
					Id:                "control_id",
					CategoryName:      testdata.MockCategoryName,
					CategoryCatalogId: testdata.MockCatalogID,
					Name:              "control_id",
				},
				toe: &orchestrator.TargetOfEvaluation{
					CloudServiceId: testdata.MockCloudServiceID,
					CatalogId:      testdata.MockCatalogID,
					AssuranceLevel: &testdata.AssuranceLevelHigh,
				},
				parentSchedulerTag: testdata.MockCloudServiceID + "control_id",
				interval:           2,
			},
			wantErr: assert.NoError,
		},
		{
			name: "Add scheduler job for second level Control",
			fields: fields{
				scheduler: gocron.NewScheduler(time.UTC),
			},
			args: args{
				c: &orchestrator.Control{
					Id:                "sub_control_id",
					CategoryName:      testdata.MockCategoryName,
					CategoryCatalogId: testdata.MockCatalogID,
					Name:              "sub_control_id",
					ParentControlId:   util.Ref("control_id"),
				},
				toe: &orchestrator.TargetOfEvaluation{
					CloudServiceId: testdata.MockCloudServiceID,
					CatalogId:      testdata.MockCatalogID,
					AssuranceLevel: &testdata.AssuranceLevelHigh,
				},
				parentSchedulerTag: testdata.MockCloudServiceID + "control_id",
				interval:           2,
			},
			wantErr: assert.NoError,
		},
		{
			name: "Interval invalid",
			fields: fields{
				scheduler: gocron.NewScheduler(time.UTC),
			},
			args: args{
				c: &orchestrator.Control{
					Id:                "sub_control_id",
					CategoryName:      testdata.MockCategoryName,
					CategoryCatalogId: testdata.MockCatalogID,
					Name:              "sub_control_id",
					ParentControlId:   util.Ref("control_id"),
				},
				toe: &orchestrator.TargetOfEvaluation{
					CloudServiceId: testdata.MockCloudServiceID,
					CatalogId:      testdata.MockCatalogID,
					AssuranceLevel: &testdata.AssuranceLevelHigh,
				},
				parentSchedulerTag: testdata.MockCloudServiceID + "control_id",
				interval:           0,
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "evaluation cannot be scheduled")
			},
		},
		{
			name: "Empty input",
			args: args{},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "evaluation cannot be scheduled")
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
				storage:                       tt.fields.storage,
			}

			// Start the scheduler
			if tt.fields.schedulerRunning == true {
				_, err := s.scheduler.Every(1).Day().Tag(tt.fields.schedulerTag).Do(func() { fmt.Println("Scheduler") })
				require.NoError(t, err)
			}

			err := s.addJobToScheduler(tt.args.c, tt.args.toe, tt.args.parentSchedulerTag, tt.args.interval)
			tt.wantErr(t, err)

			if err == nil {
				tags, err := s.scheduler.FindJobsByTag()
				assert.NoError(t, err)
				assert.NotEmpty(t, tags)
			}
		})
	}
}

func TestService_evaluateControl(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestratorClient            orchestrator.OrchestratorClient
		orchestratorAddress           grpcTarget
		authorizer                    api.Authorizer
		scheduler                     *gocron.Scheduler
		wg                            map[string]*WaitGroup
		storage                       persistence.Storage
		authz                         service.AuthorizationStrategy
	}
	type args struct {
		toe          *orchestrator.TargetOfEvaluation
		categoryName string
		controlId    string
		schedulerTag string
		subControls  []*orchestrator.Control
	}
	tests := []struct {
		name                 string
		fields               fields
		args                 args
		newEvaluationResults *evaluation.EvaluationResult
		want                 assert.ValueAssertionFunc
	}{
		{
			name: "AuthZ error",
			fields: fields{
				wg: map[string]*WaitGroup{
					testdata.MockCloudServiceID + "-" + testdata.MockControlID1: {
						wg:      &sync.WaitGroup{},
						wgMutex: sync.Mutex{},
					},
				},
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(evaluationtest.MockEvaluationResultsWithoutResultsForParentControl))
				}),
				authz: &service.AuthorizationStrategyJWT{},
			},
			args: args{
				toe: &orchestrator.TargetOfEvaluation{
					CloudServiceId: testdata.MockCloudServiceID,
					CatalogId:      testdata.MockCatalogID,
					AssuranceLevel: &testdata.AssuranceLevelHigh,
				},
				categoryName: testdata.MockCategoryName,
				controlId:    testdata.MockControlID1,
				schedulerTag: testdata.MockCloudServiceID + "-" + testdata.MockControlID1,
				subControls:  make([]*orchestrator.Control, 2),
			},
			newEvaluationResults: &evaluation.EvaluationResult{},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				service, ok := i1.(*Service)
				if !assert.True(tt, ok) {
					return false
				}

				evalResults, err := service.ListEvaluationResults(context.Background(), &evaluation.ListEvaluationResultsRequest{})
				assert.NoError(t, err)
				return assert.Equal(t, 0, len(evalResults.Results))
			},
		},
		{
			name: "No evaluation results for evaluation available",
			fields: fields{
				wg: map[string]*WaitGroup{
					testdata.MockCloudServiceID + "-" + testdata.MockControlID1: {
						wg:      &sync.WaitGroup{},
						wgMutex: sync.Mutex{},
					},
				},
				storage: testutil.NewInMemoryStorage(t),
				authz:   &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				toe: &orchestrator.TargetOfEvaluation{
					CloudServiceId: testdata.MockCloudServiceID,
					CatalogId:      testdata.MockCatalogID,
					AssuranceLevel: &testdata.AssuranceLevelHigh,
				},
				categoryName: testdata.MockCategoryName,
				controlId:    testdata.MockControlID1,
				schedulerTag: testdata.MockCloudServiceID + "-" + testdata.MockControlID1,
				subControls:  make([]*orchestrator.Control, 2),
			},
			newEvaluationResults: nil,
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				service, ok := i1.(*Service)
				if !assert.True(tt, ok) {
					return false
				}

				evalResults, err := service.ListEvaluationResults(context.Background(), &evaluation.ListEvaluationResultsRequest{})
				assert.NoError(t, err)
				return assert.Equal(t, 0, len(evalResults.Results))
			},
		},
		{
			name: "Happy path",
			fields: fields{
				wg: map[string]*WaitGroup{
					testdata.MockCloudServiceID + "-" + testdata.MockControlID1: {
						wg:      &sync.WaitGroup{},
						wgMutex: sync.Mutex{},
					},
				},
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(evaluationtest.MockEvaluationResultsWithoutResultsForParentControl))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				toe: &orchestrator.TargetOfEvaluation{
					CloudServiceId: testdata.MockCloudServiceID,
					CatalogId:      testdata.MockCatalogID,
					AssuranceLevel: &testdata.AssuranceLevelHigh,
				},
				categoryName: testdata.MockCategoryName,
				controlId:    testdata.MockControlID1,
				schedulerTag: testdata.MockCloudServiceID + "-" + testdata.MockControlID1,
				subControls:  make([]*orchestrator.Control, 2),
			},
			newEvaluationResults: evaluationtest.MockEvaluationResult1,
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				service, ok := i1.(*Service)
				if !assert.True(tt, ok) {
					return false
				}

				newEvalResults, ok := i2[0].(*evaluation.EvaluationResult)
				if !assert.True(tt, ok) {
					return false
				}

				evalResults, err := service.ListEvaluationResults(context.Background(), &evaluation.ListEvaluationResultsRequest{})
				assert.NoError(t, err)
				assert.Equal(t, 7, len(evalResults.Results))

				createdResult := evalResults.Results[len(evalResults.Results)-1]

				// Delete ID and timestamp from the evaluation results
				assert.NotEmpty(t, createdResult.GetId())
				createdResult.Id = ""
				createdResult.Timestamp = nil
				newEvalResults.Id = ""
				newEvalResults.Timestamp = nil
				return assert.Equal(t, newEvalResults, createdResult)
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
				storage:                       tt.fields.storage,
				authz:                         tt.fields.authz,
			}

			s.evaluateControl(tt.args.toe, tt.args.categoryName, tt.args.controlId, tt.args.schedulerTag, tt.args.subControls)

			tt.want(t, s, tt.newEvaluationResults)
			assert.NotEmpty(t, s.wg[tt.args.schedulerTag])
		})
	}
}

func TestService_evaluateSubcontrol(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestratorClient            orchestrator.OrchestratorClient
		orchestratorAddress           grpcTarget
		authorizer                    api.Authorizer
		scheduler                     *gocron.Scheduler
		wg                            map[string]*WaitGroup
		storage                       persistence.Storage
		authz                         service.AuthorizationStrategy
		newEvaluationResults          []*evaluation.EvaluationResult
		schedulerTag                  string
		wgCounter                     int
	}
	type args struct {
		toe                *orchestrator.TargetOfEvaluation
		categoryName       string
		controlId          string
		parentSchedulerTag string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   assert.ValueAssertionFunc
	}{
		{
			name: "Error getting metrics",
			fields: fields{
				schedulerTag: testdata.MockCloudServiceID + "-" + testdata.MockControlID1,
				wgCounter:    2,
				wg: map[string]*WaitGroup{
					testdata.MockCloudServiceID + "-" + testdata.MockControlID1: {
						wg:      &sync.WaitGroup{},
						wgMutex: sync.Mutex{},
					},
				},
				storage:              testutil.NewInMemoryStorage(t),
				authz:                &service.AuthorizationStrategyAllowAll{},
				newEvaluationResults: nil,
			},
			args: args{
				toe: &orchestrator.TargetOfEvaluation{
					CloudServiceId: testdata.MockCloudServiceID,
					CatalogId:      testdata.MockCatalogID,
					AssuranceLevel: &testdata.AssuranceLevelHigh,
				},
				categoryName:       testdata.MockCategoryName,
				controlId:          testdata.MockControlID1,
				parentSchedulerTag: getSchedulerTag(testdata.MockCloudServiceID, testdata.MockControlID1),
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				service, ok := i1.(*Service)
				if !assert.True(tt, ok) {
					return false
				}

				evalResults, err := service.ListEvaluationResults(context.Background(), &evaluation.ListEvaluationResultsRequest{})
				assert.NoError(t, err)
				return assert.Equal(t, 0, len(evalResults.Results))
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
				storage:                       tt.fields.storage,
				authz:                         tt.fields.authz,
			}

			tt.fields.wg[tt.fields.schedulerTag].wg.Add(tt.fields.wgCounter)
			s.evaluateSubcontrol(tt.args.toe, tt.args.categoryName, tt.args.controlId, tt.args.parentSchedulerTag)

			tt.want(t, s)
		})
	}
}

func Test_getSchedulerTag(t *testing.T) {
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
			name: "cloud_service_id empty",
			args: args{
				controlId: testdata.MockControlID1,
			},
			want: "",
		},
		{
			name: "control_id empty",
			args: args{
				cloudServiceId: testdata.MockCloudServiceID,
			},
			want: "",
		},
		{
			name: "Happy path",
			args: args{
				cloudServiceId: testdata.MockCloudServiceID,
				controlId:      testdata.MockControlID1,
			},
			want: fmt.Sprintf("%s-%s", testdata.MockCloudServiceID, testdata.MockControlID1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSchedulerTag(tt.args.cloudServiceId, tt.args.controlId); got != tt.want {
				t.Errorf("getSchedulerTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
