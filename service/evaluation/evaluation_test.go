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
	"sort"
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
	"clouditor.io/clouditor/internal/testutil/orchestratortest"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/service"
	"github.com/go-co-op/gocron"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestMain(m *testing.M) {
	clitest.AutoChdir()

	/*	server, orch := startBufConnServer()
		orch.CreateCatalog(context.TODO(), &orchestrator.CreateCatalogRequest{Catalog: orchestratortest.NewCatalog()})
		orch.StoreAssessmentResult(context.TODO(), &orchestrator.StoreAssessmentResultRequest{
			Result: orchestratortest.MockAssessmentResult1,
		})
		orch.StoreAssessmentResult(context.TODO(), &orchestrator.StoreAssessmentResultRequest{
			Result: orchestratortest.MockAssessmentResult2,
		})
		orch.StoreAssessmentResult(context.TODO(), &orchestrator.StoreAssessmentResultRequest{
			Result: orchestratortest.MockAssessmentResult3,
		})
		orch.StoreAssessmentResult(context.TODO(), &orchestrator.StoreAssessmentResultRequest{
			Result: orchestratortest.MockAssessmentResult4,
		})*/

	code := m.Run()

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
			name: "WithStorage",
			args: args{
				opts: []service.Option[Service]{service.Option[Service](WithStorage(testutil.NewInMemoryStorage(t)))},
			},
			want: &Service{
				storage: testutil.NewInMemoryStorage(t),
				orchestratorAddress: grpcTarget{
					target: testdata.MockOrchestratorAddress,
				},
				wg: make(map[string]*sync.WaitGroup),
			},
		},
		{
			name: "WithOrchestratorAddress",
			args: args{
				opts: []service.Option[Service]{service.Option[Service](WithOrchestratorAddress("localhost:1234"))},
			},
			want: &Service{
				orchestratorAddress: grpcTarget{
					target: "localhost:1234",
				},
				wg: make(map[string]*sync.WaitGroup),
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
				wg:         make(map[string]*sync.WaitGroup),
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
				wg:         make(map[string]*sync.WaitGroup),
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
				wg: make(map[string]*sync.WaitGroup),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewService(tt.args.opts...)

			// Check if scheduler is initialized and then remove
			if assert.NotNil(t, got.scheduler) {
				tt.want.scheduler = nil
				got.scheduler = nil
			}
			// Check if storage is initialized and then remove
			if assert.NotNil(t, got.storage) {
				tt.want.storage = nil
				got.storage = nil
			}
			// Check if authz is initialized and then remove
			if assert.NotNil(t, got.authz) {
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
		catalogId      string
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
				catalogId: testdata.MockCatalogID,
			},
			want: "",
		},
		{
			name: "Input cloudServiceId empty",
			args: args{
				cloudServiceId: testdata.MockCloudServiceID,
				catalogId:      testdata.MockCatalogID,
			},
			want: "",
		},
		{
			name: "Input catalogId empty",
			args: args{
				cloudServiceId: testdata.MockCloudServiceID,
				controlId:      testdata.MockSubControlID,
			},
			want: "",
		},
		{
			name: "Happy path",
			args: args{
				cloudServiceId: testdata.MockCloudServiceID,
				catalogId:      testdata.MockCatalogID,
				controlId:      testdata.MockSubControlID,
			},
			want: fmt.Sprintf("%s-%s-%s", testdata.MockCloudServiceID, testdata.MockCatalogID, testdata.MockSubControlID),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSchedulerTag(tt.args.cloudServiceId, tt.args.catalogId, tt.args.controlId); got != tt.want {
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
					Id:                testdata.MockControlID1,
					CategoryName:      testdata.MockCategoryName,
					CategoryCatalogId: testdata.MockCatalogID,
					Name:              testdata.MockControlName,
					Description:       testdata.MockControlDescription,
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
		{
			name: "Happy path",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				orchestratorAddress: grpcTarget{
					opts: []grpc.DialOption{grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
					})))},
				},
			},
			args: args{
				control: orchestratortest.MockControl1,
			},
			wantMetrics: orchestratortest.MockControl1.Controls[0].GetMetrics(),
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
				storage:                       tt.fields.storage,
			}
			gotMetrics, err := s.getMetricsFromSubcontrols(tt.args.control)

			tt.wantErr(t, err)

			assert.Equal(t, len(gotMetrics), len(tt.wantMetrics))
			for i := range gotMetrics {
				if !proto.Equal(gotMetrics[i], tt.wantMetrics[i]) {
					t.Errorf("Service.GetControl() = %v, want %v", gotMetrics[i], tt.wantMetrics[i])
				}
			}
		})
	}
}

func TestService_StopEvaluation(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestratorClient            orchestrator.OrchestratorClient
		orchestratorAddress           grpcTarget
		authorizer                    api.Authorizer
		scheduler                     map[string]*gocron.Scheduler
		storage                       persistence.Storage
		toeTag                        string
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
		{
			name: "Happy path",
			args: args{
				in0: context.Background(),
				req: &evaluation.StopEvaluationRequest{
					CloudServiceId: testdata.MockCloudServiceID,
					CatalogId:      testdata.MockCatalogID,
				},
				schedulerRunning: true,
			},
			fields: fields{
				scheduler: map[string]*gocron.Scheduler{
					fmt.Sprintf("%s-%s", testdata.MockCloudServiceID, testdata.MockCatalogID): gocron.NewScheduler(time.UTC),
				},
				toeTag: fmt.Sprintf("%s-%s", testdata.MockCloudServiceID, testdata.MockCatalogID),
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
				scheduler:                     tt.fields.scheduler,
				storage:                       tt.fields.storage,
			}

			// Start the scheduler
			if tt.args.schedulerRunning == true {
				_, err := s.scheduler[tt.fields.toeTag].Every(1).Day().Tag(tt.args.schedulerTag).Do(func() { fmt.Println("Scheduler") })
				require.NoError(t, err)

			}

			gotRes, err := s.StopEvaluation(tt.args.in0, tt.args.req)

			tt.wantErr(t, err)
			assert.Equal(t, tt.wantRes, gotRes)
		})
	}
}

func TestService_StartEvaluation(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestratorClient            orchestrator.OrchestratorClient
		orchestratorAddress           grpcTarget
		authorizer                    api.Authorizer
		scheduler                     map[string]*gocron.Scheduler
		storage                       persistence.Storage
		wg                            map[string]*sync.WaitGroup
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
		// 		wg:                  make(map[string]*sync.WaitGroup),
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
				_, err := s.scheduler[getToeTag(tt.args.req.GetCloudServiceId(), tt.args.req.GetCatalogId())].Every(1).Day().Tag(tt.args.schedulerTag).Do(func() { fmt.Println("Scheduler") })
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

func TestService_getAllMetricsFromControl(t *testing.T) {
	catalogWithoutSubControls := &orchestrator.Catalog{
		Name:            testdata.MockCatalogName,
		Id:              testdata.MockCatalogID,
		Description:     testdata.MockCatalogDescription,
		AllInScope:      true,
		AssuranceLevels: []string{testdata.AssuranceLevelBasic, testdata.AssuranceLevelSubstantial, testdata.AssuranceLevelHigh},
		Categories: []*orchestrator.Category{{
			Name:        testdata.MockCategoryName,
			Description: testdata.MockCategoryDescription,
			CatalogId:   testdata.MockCatalogID,
			Controls: []*orchestrator.Control{
				{
					Id:                testdata.MockControlID1,
					Name:              testdata.MockControlName,
					CategoryName:      testdata.MockCategoryName,
					CategoryCatalogId: testdata.MockCatalogID,
					Description:       testdata.MockControlDescription,
					AssuranceLevel:    &testdata.AssuranceLevelHigh,
				},
			},
		}}}

	catalogWithControlAnd2Metrics := orchestratortest.NewCatalog()
	catalogWithControlAnd2Metrics.Categories[0].Controls[0].Controls[0].Metrics = append(catalogWithControlAnd2Metrics.Categories[0].Controls[0].Controls[0].Metrics, &assessment.Metric{
		Id:          testdata.MockAnotherMetricID,
		Name:        testdata.MockAnotherMetricID,
		Description: testdata.MockMetricDescription,
		Scale:       assessment.Metric_ORDINAL,
		Range: &assessment.Range{
			Range: &assessment.Range_AllowedValues{
				AllowedValues: &assessment.AllowedValues{
					Values: []*structpb.Value{
						structpb.NewBoolValue(false),
						structpb.NewBoolValue(true),
					},
				},
			},
		},
	})

	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestratorClient            orchestrator.OrchestratorClient
		orchestratorAddress           grpcTarget
		authorizer                    api.Authorizer
		scheduler                     map[string]*gocron.Scheduler
		wg                            map[string]*sync.WaitGroup
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
			name: "Input empty",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				orchestratorAddress: grpcTarget{
					opts: []grpc.DialOption{grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
					})))},
				},
			},
			wantMetrics: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "code = InvalidArgument desc = invalid request: invalid GetControlRequest.CatalogId: value length must be at least 1 runes")
			},
		},
		{
			name: "metric not exists",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				orchestratorAddress: grpcTarget{
					opts: []grpc.DialOption{grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
					})))},
				},
			},
			args: args{
				catalogId:    "test_catalog_id",
				categoryName: "test_category_id",
				controlId:    "test_control_id",
			},
			wantMetrics: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "code = NotFound desc = control not found")
			},
		},
		{
			name: "no sub-controls available",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				orchestratorAddress: grpcTarget{
					opts: []grpc.DialOption{grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						assert.NoError(t, s.Create(catalogWithoutSubControls))
					})))},
				},
			},
			args: args{
				catalogId:    testdata.MockCatalogID,
				categoryName: testdata.MockCategoryName,
				controlId:    testdata.MockControlID1,
			},
			wantMetrics: nil,
			wantErr:     assert.NoError,
		},
		{
			name: "Happy path",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				orchestratorAddress: grpcTarget{
					opts: []grpc.DialOption{grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						assert.NoError(t, s.Create(catalogWithControlAnd2Metrics))
					})))},
				},
			},
			args: args{
				catalogId:    testdata.MockCatalogID,
				categoryName: testdata.MockCategoryName,
				controlId:    testdata.MockControlID1,
			},
			wantMetrics: []*assessment.Metric{
				{
					Id:          testdata.MockAnotherMetricID,
					Name:        testdata.MockAnotherMetricID,
					Description: testdata.MockMetricDescription,
					Scale:       assessment.Metric_ORDINAL,
					Range: &assessment.Range{
						Range: &assessment.Range_AllowedValues{
							AllowedValues: &assessment.AllowedValues{
								Values: []*structpb.Value{
									structpb.NewBoolValue(false),
									structpb.NewBoolValue(true),
								},
							},
						},
					},
				},
				{
					Id:          testdata.MockMetricID,
					Name:        testdata.MockMetricName,
					Description: testdata.MockMetricDescription,
					Scale:       assessment.Metric_ORDINAL,
					Range: &assessment.Range{
						Range: &assessment.Range_AllowedValues{
							AllowedValues: &assessment.AllowedValues{
								Values: []*structpb.Value{
									structpb.NewBoolValue(false),
									structpb.NewBoolValue(true),
								},
							},
						},
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
				scheduler:                     tt.fields.scheduler,
				wg:                            tt.fields.wg,
				storage:                       tt.fields.storage,
			}
			gotMetrics, err := s.getAllMetricsFromControl(tt.args.catalogId, tt.args.categoryName, tt.args.controlId)
			tt.wantErr(t, err)

			assert.Equal(t, len(gotMetrics), len(tt.wantMetrics))
			for i := range gotMetrics {
				reflect.DeepEqual(gotMetrics[i], tt.wantMetrics[i])
			}
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

func TestService_getControl(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestratorClient            orchestrator.OrchestratorClient
		orchestratorAddress           grpcTarget
		authorizer                    api.Authorizer
		scheduler                     map[string]*gocron.Scheduler
		wg                            map[string]*sync.WaitGroup
		storage                       persistence.Storage
	}
	type args struct {
		catalogId    string
		categoryName string
		controlId    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.ValueAssertionFunc
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "orchestrator address is missing",
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not connect to orchestrator service:")
			},
		},
		{
			name: "catalog_id is missing",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				orchestratorAddress: grpcTarget{
					opts: []grpc.DialOption{grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
					})))},
				},
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid request: invalid GetControlRequest.CatalogId: value length must be at least 1 runes")
			},
		},
		{
			name: "category_name is missing",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				orchestratorAddress: grpcTarget{
					opts: []grpc.DialOption{grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
					})))},
				},
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid request: invalid GetControlRequest.CatalogId: value length must be at least 1 runes")
			},
		},
		{
			name: "control_id is missing",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				orchestratorAddress: grpcTarget{
					opts: []grpc.DialOption{grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
					})))},
				},
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid request: invalid GetControlRequest.CatalogId: value length must be at least 1 runes")
			},
		},
		{
			name: "control does not exist",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				orchestratorAddress: grpcTarget{
					opts: []grpc.DialOption{grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
					})))},
				},
			},
			args: args{
				catalogId:    "wrong_catalog_id",
				categoryName: "wrong_category_id",
				controlId:    "wrong_control_id",
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "control not found")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				orchestratorAddress: grpcTarget{
					opts: []grpc.DialOption{grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
					})))},
				},
			},
			args: args{
				catalogId:    testdata.MockCatalogID,
				categoryName: testdata.MockCategoryName,
				controlId:    testdata.MockControlID1,
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				gotControl, ok := i1.(*orchestrator.Control)
				if !assert.True(tt, ok) {
					return false
				}

				// We need to truncate the metric from the control because the control is only returned with its sub-control but without the sub-control's metric.
				wantControl := orchestratortest.MockControl1
				wantControl.Controls[0].Metrics = nil

				if !proto.Equal(gotControl, wantControl) {
					t.Errorf("Service.GetControl() = %v, want %v", gotControl, wantControl)
					return false
				}

				return true
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
				storage:                       tt.fields.storage,
			}

			gotControl, err := s.getControl(tt.args.catalogId, tt.args.categoryName, tt.args.controlId)
			tt.wantErr(t, err)

			if gotControl != nil {
				tt.want(t, gotControl)
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
		scheduler                     map[string]*gocron.Scheduler
		wg                            map[string]*sync.WaitGroup
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
			name: "Empty input",
			args: args{},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "evaluation cannot be scheduled")
			},
		},
		{
			name: "Interval invalid",
			fields: fields{
				scheduler: make(map[string]*gocron.Scheduler),
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
			name: "ToE input empty",
			args: args{
				c: &orchestrator.Control{
					Id:                "sub_control_id",
					CategoryName:      testdata.MockCategoryName,
					CategoryCatalogId: testdata.MockCatalogID,
					Name:              "sub_control_id",
					ParentControlId:   util.Ref("control_id"),
				},
				parentSchedulerTag: testdata.MockCloudServiceID + "control_id",
				toe: &orchestrator.TargetOfEvaluation{
					CloudServiceId: testdata.MockCloudServiceID,
					CatalogId:      testdata.MockCatalogID,
					AssuranceLevel: &testdata.AssuranceLevelHigh,
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "interval is invalid")
			},
		},
		{
			name: "Happy path: Add scheduler job for control",
			fields: fields{
				scheduler: make(map[string]*gocron.Scheduler),
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
				parentSchedulerTag: "",
				interval:           2,
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: Add scheduler job for sub-control",
			fields: fields{
				scheduler: make(map[string]*gocron.Scheduler),
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

			toeTag := getToeTag(tt.args.toe.GetCloudServiceId(), tt.args.toe.GetCatalogId())
			// Start the scheduler
			if tt.fields.schedulerRunning == true {
				_, err := s.scheduler[toeTag].Every(1).Day().Tag(tt.fields.schedulerTag).Do(func() { fmt.Println("Scheduler") })
				require.NoError(t, err)
			}
			if tt.fields.scheduler != nil {
				tt.fields.scheduler[toeTag] = gocron.NewScheduler(time.UTC)
			}

			err := s.addJobToScheduler(tt.args.c, tt.args.toe, tt.args.parentSchedulerTag, tt.args.interval)
			tt.wantErr(t, err)

			if err == nil {
				tags, err := s.scheduler[toeTag].FindJobsByTag()
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
		scheduler                     map[string]*gocron.Scheduler
		wg                            map[string]*sync.WaitGroup
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
			name: "AuthZ error in ListEvaluationResults",
			fields: fields{
				wg: map[string]*sync.WaitGroup{
					testdata.MockCloudServiceID + "-" + testdata.MockControlID1: {},
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
				wg: map[string]*sync.WaitGroup{
					testdata.MockCloudServiceID + "-" + testdata.MockControlID1: {},
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
				wg: map[string]*sync.WaitGroup{
					testdata.MockCloudServiceID + "-" + testdata.MockControlID1: {},
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
		orchestratorClient  orchestrator.OrchestratorClient
		orchestratorAddress grpcTarget
		authorizer          api.Authorizer
		scheduler           map[string]*gocron.Scheduler
		wg                  map[string]*sync.WaitGroup
		storage             persistence.Storage
		authz               service.AuthorizationStrategy
		schedulerTag        string
		wgCounter           int
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
			name: "ToE input empty", // we do not check the other input parameters
			fields: fields{
				schedulerTag: getSchedulerTag(testdata.MockCloudServiceID, testdata.MockCatalogID, testdata.MockControlID1),
				orchestratorAddress: grpcTarget{
					opts: []grpc.DialOption{grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
						assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
					})))},
				},
				storage:   testutil.NewInMemoryStorage(t),
				authz:     &service.AuthorizationStrategyAllowAll{},
				wgCounter: 2,
				wg: map[string]*sync.WaitGroup{
					getSchedulerTag(testdata.MockCloudServiceID, testdata.MockCatalogID, testdata.MockControlID1): {},
				},
			},
			args: args{
				categoryName:       testdata.MockCategoryName,
				controlId:          testdata.MockControlID1,
				parentSchedulerTag: getSchedulerTag(testdata.MockCloudServiceID, testdata.MockCatalogID, testdata.MockControlID1),
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
		{
			name: "no assessment results available",
			fields: fields{
				schedulerTag: getSchedulerTag(testdata.MockCloudServiceID, testdata.MockCatalogID, testdata.MockControlID1),
				wgCounter:    2,
				wg: map[string]*sync.WaitGroup{
					getSchedulerTag(testdata.MockCloudServiceID, testdata.MockCatalogID, testdata.MockControlID1): {},
				},
				storage: testutil.NewInMemoryStorage(t),
				authz:   &service.AuthorizationStrategyAllowAll{},
				orchestratorAddress: grpcTarget{
					opts: []grpc.DialOption{grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t)))},
				},
			},
			args: args{
				toe: &orchestrator.TargetOfEvaluation{
					CloudServiceId: testdata.MockCloudServiceID,
					CatalogId:      testdata.MockCatalogID,
					AssuranceLevel: &testdata.AssuranceLevelHigh,
				},
				categoryName:       testdata.MockCategoryName,
				controlId:          testdata.MockSubControlID11,
				parentSchedulerTag: getSchedulerTag(testdata.MockCloudServiceID, testdata.MockCatalogID, testdata.MockControlID1),
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
		{
			name: "error getting metrics",
			fields: fields{
				schedulerTag: getSchedulerTag(testdata.MockCloudServiceID, testdata.MockCatalogID, testdata.MockControlID1),
				wgCounter:    2,
				wg: map[string]*sync.WaitGroup{
					getSchedulerTag(testdata.MockCloudServiceID, testdata.MockCatalogID, testdata.MockControlID1): {},
				},
				storage: testutil.NewInMemoryStorage(t),
				authz:   &service.AuthorizationStrategyAllowAll{},
				orchestratorAddress: grpcTarget{
					opts: []grpc.DialOption{grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t)))},
				},
			},
			args: args{
				toe: &orchestrator.TargetOfEvaluation{
					CloudServiceId: testdata.MockCloudServiceID,
					CatalogId:      testdata.MockCatalogID,
					AssuranceLevel: &testdata.AssuranceLevelHigh,
				},
				categoryName:       testdata.MockCategoryName,
				controlId:          testdata.MockSubControlID11,
				parentSchedulerTag: getSchedulerTag(testdata.MockCloudServiceID, testdata.MockCatalogID, testdata.MockControlID1),
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
		{
			name: "Happy path",
			fields: fields{
				schedulerTag: getSchedulerTag(testdata.MockCloudServiceID, testdata.MockCatalogID, testdata.MockControlID1),
				wgCounter:    2,
				wg: map[string]*sync.WaitGroup{
					getSchedulerTag(testdata.MockCloudServiceID, testdata.MockCatalogID, testdata.MockControlID1): {},
				},
				storage: testutil.NewInMemoryStorage(t),
				authz:   &service.AuthorizationStrategyAllowAll{},
				orchestratorAddress: grpcTarget{
					opts: []grpc.DialOption{grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
						assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
					})))},
				},
			},
			args: args{
				toe: &orchestrator.TargetOfEvaluation{
					CloudServiceId: testdata.MockCloudServiceID,
					CatalogId:      testdata.MockCatalogID,
					AssuranceLevel: &testdata.AssuranceLevelHigh,
				},
				categoryName:       testdata.MockCategoryName,
				controlId:          testdata.MockSubControlID11,
				parentSchedulerTag: getSchedulerTag(testdata.MockCloudServiceID, testdata.MockCatalogID, testdata.MockControlID1),
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				service, ok := i1.(*Service)
				if !assert.True(tt, ok) {
					return false
				}

				evalResults, err := service.ListEvaluationResults(context.Background(), &evaluation.ListEvaluationResultsRequest{})
				assert.NoError(t, err)
				return assert.Equal(t, 1, len(evalResults.Results))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				orchestratorClient:  tt.fields.orchestratorClient,
				orchestratorAddress: tt.fields.orchestratorAddress,
				authorizer:          tt.fields.authorizer,
				scheduler:           tt.fields.scheduler,
				wg:                  tt.fields.wg,
				storage:             tt.fields.storage,
				authz:               tt.fields.authz,
			}

			err := s.initOrchestratorClient()
			assert.NoError(t, err)

			tt.fields.wg[tt.fields.schedulerTag].Add(tt.fields.wgCounter)
			s.evaluateSubcontrol(tt.args.toe, tt.args.categoryName, tt.args.controlId, tt.args.parentSchedulerTag)

			if tt.want != nil {
				tt.want(t, s)
			}
		})
	}
}

func Test_getSchedulerTag(t *testing.T) {
	type args struct {
		cloudServiceId string
		catalogId      string
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
				catalogId: testdata.MockCatalogID,
			},
			want: "",
		},
		{
			name: "control_id empty",
			args: args{
				cloudServiceId: testdata.MockCloudServiceID,
				catalogId:      testdata.MockCatalogID,
			},
			want: "",
		},
		{
			name: "catalog_id empty",
			args: args{
				cloudServiceId: testdata.MockCloudServiceID,
				controlId:      testdata.MockControlID1,
			},
			want: "",
		},
		{
			name: "Happy path",
			args: args{
				cloudServiceId: testdata.MockCloudServiceID,
				controlId:      testdata.MockControlID1,
				catalogId:      testdata.MockCatalogID,
			},
			want: fmt.Sprintf("%s-%s-%s", testdata.MockCloudServiceID, testdata.MockCatalogID, testdata.MockControlID1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSchedulerTag(tt.args.cloudServiceId, tt.args.catalogId, tt.args.controlId); got != tt.want {
				t.Errorf("getSchedulerTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getEvaluationResultMap(t *testing.T) {
	type args struct {
		results []*evaluation.EvaluationResult
	}
	tests := []struct {
		name string
		args args
		want assert.ValueAssertionFunc
	}{
		{
			name: "Empty input",
			args: args{results: []*evaluation.EvaluationResult{}},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				gotResults, ok := i1.(map[string][]*evaluation.EvaluationResult)
				if !assert.True(tt, ok) {
					return false
				}
				return assert.Empty(t, gotResults)
			},
		},
		{
			name: "Happy path",
			args: args{
				results: evaluationtest.MockEvaluationResults,
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				gotResults, ok := i1.(map[string][]*evaluation.EvaluationResult)
				if !assert.True(tt, ok) {
					return false
				}

				wantResults := map[string][]*evaluation.EvaluationResult{
					testdata.MockResourceID: {
						evaluationtest.MockEvaluationResult6,
						evaluationtest.MockEvaluationResult5,
						evaluationtest.MockEvaluationResult4,
						evaluationtest.MockEvaluationResult3,
						evaluationtest.MockEvaluationResult2,
						evaluationtest.MockEvaluationResult22,
						evaluationtest.MockEvaluationResult1,
					},
					testdata.MockAnotherResourceID: {
						evaluationtest.MockEvaluationResult8,
						evaluationtest.MockEvaluationResult7,
					},
				}

				for _, r := range gotResults {
					sort.SliceStable(r, func(i, j int) bool {
						return r[i].GetId() < r[j].GetId()
					})
				}
				for _, r := range wantResults {
					sort.SliceStable(r, func(i, j int) bool {
						return r[i].GetId() < r[j].GetId()
					})
				}

				assert.Equal(t, len(gotResults), len(wantResults))

				for key, results := range gotResults {
					if !assert.Equal(t, len(results), len(wantResults[key])) {
						return false
					}
					for i := range results {
						reflect.DeepEqual(results[i], wantResults[key][i])
					}
				}

				return true
			},
			//
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getEvaluationResultMap(tt.args.results)

			tt.want(t, got)
		})
	}
}

func Test_getAssessmentResultMap(t *testing.T) {
	type args struct {
		results []*assessment.AssessmentResult
	}
	tests := []struct {
		name string
		args args
		want assert.ValueAssertionFunc
	}{
		{
			name: "Empty input",
			args: args{results: []*assessment.AssessmentResult{}},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				gotResults, ok := i1.(map[string][]*assessment.AssessmentResult)
				if !assert.True(tt, ok) {
					return false
				}
				return assert.Empty(t, gotResults)
			},
		},
		{
			name: "Happy path",
			args: args{
				results: orchestratortest.MockAssessmentResults,
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				gotResults, ok := i1.(map[string][]*assessment.AssessmentResult)
				if !assert.True(tt, ok) {
					return false
				}

				wantResults := map[string][]*assessment.AssessmentResult{
					testdata.MockResourceID: {
						orchestratortest.MockAssessmentResult1,
						orchestratortest.MockAssessmentResult2,
						orchestratortest.MockAssessmentResult3,
					},
					testdata.MockAnotherResourceID: {
						orchestratortest.MockAssessmentResult4,
					},
				}

				for _, r := range gotResults {
					sort.SliceStable(r, func(i, j int) bool {
						return r[i].GetId() < r[j].GetId()
					})
				}
				for _, r := range wantResults {
					sort.SliceStable(r, func(i, j int) bool {
						return r[i].GetId() < r[j].GetId()
					})
				}

				assert.Equal(t, len(gotResults), len(wantResults))

				for key, results := range gotResults {
					if !assert.Equal(t, len(results), len(wantResults[key])) {
						return false
					}
					for i := range results {
						reflect.DeepEqual(results[i], wantResults[key][i])
					}
				}

				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getAssessmentResultMap(tt.args.results)

			tt.want(t, got)
		})
	}
}

func Test_getControlsInScopeHierarchy(t *testing.T) {
	type args struct {
		controls []*orchestrator.Control
	}
	tests := []struct {
		name string
		args args
		want assert.ValueAssertionFunc
	}{
		{
			name: "Empty input",
			args: args{controls: []*orchestrator.Control{}},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				gotResults, ok := i1.([]*orchestrator.Control)
				if !assert.True(tt, ok) {
					return false
				}
				return assert.Empty(t, gotResults)
			},
		},
		{
			name: "Happy path",
			args: args{
				controls: orchestratortest.MockControlsInScope,
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				gotResults, ok := i1.([]*orchestrator.Control)
				if !assert.True(tt, ok) {
					return false
				}

				wantResults := []*orchestrator.Control{
					orchestratortest.MockControlsInScope5,
					orchestratortest.MockControlsInScope4,
					orchestratortest.MockControlsInScope3,
					orchestratortest.MockControlsInScope2,
					orchestratortest.MockControlsInScope1,
				}
				wantResults[3].Controls = append(wantResults[3].Controls, orchestratortest.MockControlsInScopeSubControl11)
				wantResults[2].Controls = append(wantResults[2].Controls, orchestratortest.MockControlsInScopeSubControl21)
				wantResults[1].Controls = append(wantResults[1].Controls, orchestratortest.MockControlsInScopeSubControl31)
				wantResults[0].Controls = append(wantResults[0].Controls, orchestratortest.MockControlsInScopeSubControl32)

				// Sort lists
				sort.SliceStable(gotResults, func(i, j int) bool {
					return gotResults[i].GetId() < gotResults[j].GetId()
				})
				sort.SliceStable(wantResults, func(i, j int) bool {
					return wantResults[i].GetId() < wantResults[j].GetId()
				})
				assert.Equal(t, len(wantResults), len(gotResults))

				for i := range gotResults {
					reflect.DeepEqual(gotResults[i], wantResults[i])
				}

				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotControlsHierarchy := getControlsInScopeHierarchy(tt.args.controls)

			tt.want(t, gotControlsHierarchy)
		})
	}
}

func Test_getToeTag(t *testing.T) {
	type args struct {
		cloudServiceId string
		catalogId      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty cloud_service_id",
			args: args{
				catalogId: testdata.MockCatalogID,
			},
			want: "",
		},
		{
			name: "Empty catalog_id",
			args: args{
				cloudServiceId: testdata.MockCloudServiceID,
			},
			want: "",
		},
		{
			name: "Happy path",
			args: args{
				cloudServiceId: testdata.MockCloudServiceID,
				catalogId:      testdata.MockCatalogID,
			},
			want: fmt.Sprintf("%s-%s", testdata.MockCloudServiceID, testdata.MockCatalogID),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getToeTag(tt.args.cloudServiceId, tt.args.catalogId); got != tt.want {
				t.Errorf("getToeTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
