// Copyright 2023 Fraunhofer AISEC
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
	"clouditor.io/clouditor/internal/testutil/servicetest/evaluationtest"
	"clouditor.io/clouditor/internal/testutil/servicetest/orchestratortest"
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
	var inmem = testutil.NewInMemoryStorage(t)

	type args struct {
		opts []service.Option[Service]
	}
	tests := []struct {
		name string
		args args
		want assert.ValueAssertionFunc
	}{
		{
			name: "WithStorage",
			args: args{
				opts: []service.Option[Service]{service.Option[Service](WithStorage(inmem))},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				s := i1.(*Service)
				return assert.Equal(t, inmem, s.storage)
			},
		},
		{
			name: "WithOrchestratorAddress",
			args: args{
				opts: []service.Option[Service]{service.Option[Service](WithOrchestratorAddress(testdata.MockOrchestratorAddress))},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				s := i1.(*Service)
				return assert.Equal(t, testdata.MockOrchestratorAddress, s.orchestrator.Target)
			},
		},
		{
			name: "WithOAuth2Authorizer",
			args: args{
				opts: []service.Option[Service]{service.Option[Service](WithOAuth2Authorizer(&clientcredentials.Config{}))},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				s := i1.(*Service)
				return assert.Equal(t, api.NewOAuthAuthorizerFromClientCredentials(&clientcredentials.Config{}), s.orchestrator.Authorizer())
			},
		},
		{
			name: "WithAuthorizer",
			args: args{
				opts: []service.Option[Service]{service.Option[Service](WithAuthorizer(api.NewOAuthAuthorizerFromClientCredentials(&clientcredentials.Config{})))},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				s := i1.(*Service)
				return assert.Equal(t, api.NewOAuthAuthorizerFromClientCredentials(&clientcredentials.Config{}), s.orchestrator.Authorizer())
			},
		},
		{
			name: "Happy path",
			args: args{
				opts: []service.Option[Service]{},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				s := i1.(*Service)
				return assert.Equal(t, DefaultOrchestratorAddress, s.orchestrator.Target)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewService(tt.args.opts...)
			tt.want(t, got)
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
				cloudServiceId: testdata.MockCloudServiceID1,
				catalogId:      testdata.MockCatalogID,
			},
			want: "",
		},
		{
			name: "Input catalogId empty",
			args: args{
				cloudServiceId: testdata.MockCloudServiceID1,
				controlId:      testdata.MockSubControlID,
			},
			want: "",
		},
		{
			name: "Happy path",
			args: args{
				cloudServiceId: testdata.MockCloudServiceID1,
				catalogId:      testdata.MockCatalogID,
				controlId:      testdata.MockSubControlID,
			},
			want: fmt.Sprintf("%s-%s-%s", testdata.MockCloudServiceID1, testdata.MockCatalogID, testdata.MockSubControlID),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createJobTag(tt.args.cloudServiceId, tt.args.catalogId, tt.args.controlId); got != tt.want {
				t.Errorf("createSchedulerTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_ListEvaluationResults(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestrator                  *api.RPCConnection[orchestrator.OrchestratorClient]
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
				return assert.ErrorContains(t, err, api.ErrEmptyRequest.Error())
			},
		},
		{
			name: "Filter latest_by_control_id, control_id, sub_controls, cloud_service_id",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(evaluationtest.MockEvaluationResults))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.ListEvaluationResultsRequest{
					LatestByControlId: util.Ref(true),
					Filter: &evaluation.ListEvaluationResultsRequest_Filter{
						ControlId:      util.Ref(testdata.MockSubControlID11),
						SubControls:    util.Ref(testdata.MockControlID1),
						CloudServiceId: util.Ref(testdata.MockCloudServiceID1),
					},
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
			name: "Filter latest_by_control_id and control_id",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(evaluationtest.MockEvaluationResults))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.ListEvaluationResultsRequest{
					LatestByControlId: util.Ref(true),
					Filter: &evaluation.ListEvaluationResultsRequest_Filter{
						ControlId: util.Ref(testdata.MockSubControlID11),
					},
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
			name: "Filter latest_by_control_id",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(evaluationtest.MockEvaluationResults))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.ListEvaluationResultsRequest{LatestByControlId: util.Ref(true)},
			},
			wantRes: &evaluation.ListEvaluationResultsResponse{
				Results: []*evaluation.EvaluationResult{
					evaluationtest.MockEvaluationResult1,
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
			name: "Filter latest_by_control_id, parents_only",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(evaluationtest.MockEvaluationResults))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.ListEvaluationResultsRequest{
					LatestByControlId: util.Ref(true),
					Filter: &evaluation.ListEvaluationResultsRequest_Filter{
						ParentsOnly: util.Ref(true),
					},
				},
			},
			wantRes: &evaluation.ListEvaluationResultsResponse{
				Results: []*evaluation.EvaluationResult{
					evaluationtest.MockEvaluationResult1,
					evaluationtest.MockEvaluationResult4,
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
				req: &evaluation.ListEvaluationResultsRequest{
					Filter: &evaluation.ListEvaluationResultsRequest_Filter{
						ControlId: util.Ref(testdata.MockControlID1),
					},
				},
			},
			wantRes: &evaluation.ListEvaluationResultsResponse{
				Results: []*evaluation.EvaluationResult{
					evaluationtest.MockEvaluationResult1,
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
				req: &evaluation.ListEvaluationResultsRequest{
					Filter: &evaluation.ListEvaluationResultsRequest_Filter{
						SubControls: util.Ref(testdata.MockControlID1),
					},
				},
			},
			wantRes: &evaluation.ListEvaluationResultsResponse{
				Results: []*evaluation.EvaluationResult{
					evaluationtest.MockEvaluationResult1,
					evaluationtest.MockEvaluationResult2,
					evaluationtest.MockEvaluationResult22,
					evaluationtest.MockEvaluationResult3,
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
				req: &evaluation.ListEvaluationResultsRequest{
					Filter: &evaluation.ListEvaluationResultsRequest_Filter{
						CloudServiceId: util.Ref(testdata.MockCloudServiceID1),
					},
				},
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
				orchestrator:                  tt.fields.orchestrator,
				authorizer:                    tt.fields.authorizer,
				storage:                       tt.fields.storage,
				authz:                         tt.fields.authz,
			}
			gotRes, err := s.ListEvaluationResults(tt.args.in0, tt.args.req)

			tt.wantErr(t, err)
			if !proto.Equal(gotRes, tt.wantRes) {
				t.Errorf("ListEvaluationResults() gotResp = %v, want %v", gotRes, tt.wantRes)
			}

		})
	}
}

func TestService_getMetricsFromSubControls(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestrator                  *api.RPCConnection[orchestrator.OrchestratorClient]
		authorizer                    api.Authorizer
		storage                       persistence.Storage
		catalogControls               map[string]map[string]*orchestrator.Control
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
			name:   "Error getting control",
			fields: fields{},
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
				return assert.ErrorIs(t, err, ErrControlNotAvailable)
			},
		},
		{
			name: "Happy path",
			fields: fields{
				catalogControls: map[string]map[string]*orchestrator.Control{
					orchestratortest.MockControl1.GetCategoryCatalogId(): {
						fmt.Sprintf("%s-%s", orchestratortest.MockControl1.GetCategoryName(), orchestratortest.MockControl1.GetId()):   orchestratortest.MockControl1,
						fmt.Sprintf("%s-%s", orchestratortest.MockControl11.GetCategoryName(), orchestratortest.MockControl11.GetId()): orchestratortest.MockControl11,
					},
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
				orchestrator:                  tt.fields.orchestrator,
				authorizer:                    tt.fields.authorizer,
				storage:                       tt.fields.storage,
				catalogControls:               tt.fields.catalogControls,
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
		orchestrator                  *api.RPCConnection[orchestrator.OrchestratorClient]
		authorizer                    api.Authorizer
		scheduler                     map[string]*gocron.Scheduler
		storage                       persistence.Storage
		authz                         service.AuthorizationStrategy
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
			name: "Request input missing",
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
			name: "Not authorized",
			fields: fields{
				authz: &service.AuthorizationStrategyJWT{},
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.StopEvaluationRequest{
					CloudServiceId: testdata.MockCloudServiceID1,
					CatalogId:      testdata.MockCatalogID,
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, service.ErrPermissionDenied)
			},
		},
		{
			name: "Evaluation not running",
			args: args{
				in0: context.Background(),
				req: &evaluation.StopEvaluationRequest{
					CloudServiceId: testdata.MockCloudServiceID1,
					CatalogId:      testdata.MockCatalogID,
				},
				schedulerRunning: false,
			},
			fields: fields{
				scheduler: make(map[string]*gocron.Scheduler),
				authz:     &service.AuthorizationStrategyAllowAll{},
				toeTag:    fmt.Sprintf("%s-%s", testdata.MockCloudServiceID1, testdata.MockCatalogID),
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, fmt.Sprintf("evaluation for Cloud Service '%s' and Catalog '%s' not running", testdata.MockCloudServiceID1, testdata.MockCatalogID))
			},
		},
		{
			name: "Happy path",
			args: args{
				in0: context.Background(),
				req: &evaluation.StopEvaluationRequest{
					CloudServiceId: testdata.MockCloudServiceID1,
					CatalogId:      testdata.MockCatalogID,
				},
				schedulerRunning: true,
			},
			fields: fields{
				scheduler: map[string]*gocron.Scheduler{
					fmt.Sprintf("%s-%s", testdata.MockCloudServiceID1, testdata.MockCatalogID): gocron.NewScheduler(time.UTC),
				},
				authz:  &service.AuthorizationStrategyAllowAll{},
				toeTag: fmt.Sprintf("%s-%s", testdata.MockCloudServiceID1, testdata.MockCatalogID),
			},
			wantRes: &evaluation.StopEvaluationResponse{},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				UnimplementedEvaluationServer: tt.fields.UnimplementedEvaluationServer,
				orchestrator:                  tt.fields.orchestrator,
				authorizer:                    tt.fields.authorizer,
				scheduler:                     tt.fields.scheduler,
				storage:                       tt.fields.storage,
				authz:                         tt.fields.authz,
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
		orchestrator    *api.RPCConnection[orchestrator.OrchestratorClient]
		authorizer      api.Authorizer
		scheduler       map[string]*gocron.Scheduler
		storage         persistence.Storage
		authz           service.AuthorizationStrategy
		wg              map[string]*sync.WaitGroup
		catalogControls map[string]map[string]*orchestrator.Control
	}
	type args struct {
		in0 context.Context
		req *evaluation.StartEvaluationRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.ValueAssertionFunc
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Request input missing",
			args: args{
				in0: context.Background(),
				req: &evaluation.StartEvaluationRequest{},
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid StartEvaluationRequest.CloudServiceId: value must be a valid UUID")
			},
		},
		{
			name: "Not authorized",
			fields: fields{
				authz: &service.AuthorizationStrategyJWT{},
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.StartEvaluationRequest{
					CloudServiceId: testdata.MockCloudServiceID2,
					CatalogId:      testdata.MockCatalogID,
					Interval:       proto.Int32(5),
				},
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, service.ErrPermissionDenied)
			},
		},
		{
			name: "error init orchestrator client",
			fields: fields{
				authz:        &service.AuthorizationStrategyAllowAll{},
				orchestrator: api.NewRPCConnection("bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(connectionRefusedDialer)),
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.StartEvaluationRequest{
					CloudServiceId: testdata.MockCloudServiceID2,
					CatalogId:      testdata.MockCatalogID,
					Interval:       proto.Int32(5),
				},
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "connection refused")
			},
		},
		{
			name: "error cache controls",
			fields: fields{
				orchestrator:    api.NewRPCConnection("bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t)))),
				authz:           &service.AuthorizationStrategyAllowAll{},
				catalogControls: make(map[string]map[string]*orchestrator.Control),
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.StartEvaluationRequest{
					CloudServiceId: testdata.MockCloudServiceID1,
					CatalogId:      testdata.MockCatalogID,
					Interval:       proto.Int32(5),
				},
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not cache controls:")
			},
		},
		{
			name: "error get ToE",
			fields: fields{
				orchestrator: api.NewRPCConnection("bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
				})))),
				authz:           &service.AuthorizationStrategyAllowAll{},
				catalogControls: make(map[string]map[string]*orchestrator.Control),
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.StartEvaluationRequest{
					CloudServiceId: testdata.MockCloudServiceID2,
					CatalogId:      testdata.MockCatalogID,
					Interval:       proto.Int32(5),
				},
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not get target of evaluation:")
			},
		},
		{
			name: "scheduler for control started already",
			fields: fields{
				orchestrator: api.NewRPCConnection("bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
					assert.NoError(t, s.Create(&orchestrator.CloudService{Id: testdata.MockCloudServiceID1}))
					assert.NoError(t, s.Create(orchestratortest.NewTargetOfEvaluation(testdata.AssuranceLevelBasic)))
				})))),
				scheduler: map[string]*gocron.Scheduler{
					fmt.Sprintf("%s-%s", testdata.MockCloudServiceID1, testdata.MockCatalogID): gocron.NewScheduler(time.UTC),
				},
				authz:           &service.AuthorizationStrategyAllowAll{},
				catalogControls: make(map[string]map[string]*orchestrator.Control),
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.StartEvaluationRequest{
					CloudServiceId: testdata.MockCloudServiceID1,
					CatalogId:      testdata.MockCatalogID,
					Interval:       proto.Int32(5),
				},
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "code = AlreadyExists desc = evaluation for Cloud Service ")
			},
		},
		{
			name: "No controls_in_scope and no scheduler added",
			fields: fields{
				orchestrator: api.NewRPCConnection("bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
					assert.NoError(t, s.Create(&orchestrator.CloudService{Id: testdata.MockCloudServiceID1}))
					assert.NoError(t, s.Create(orchestratortest.NewTargetOfEvaluationWithoutControlsInScope(testdata.AssuranceLevelBasic)))
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
				})))),
				scheduler:       make(map[string]*gocron.Scheduler),
				authz:           &service.AuthorizationStrategyAllowAll{},
				catalogControls: make(map[string]map[string]*orchestrator.Control),
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.StartEvaluationRequest{
					CloudServiceId: testdata.MockCloudServiceID1,
					CatalogId:      testdata.MockCatalogID,
				},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				_, ok := i1.(*Service)
				if !assert.True(tt, ok) {
					return false
				}

				gotResp, ok := i2[0].(*evaluation.StartEvaluationResponse)
				if !assert.True(tt, ok) {
					return false
				}
				return assert.True(t, gotResp.Successful)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: scheduler added",
			fields: fields{
				orchestrator: api.NewRPCConnection("bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
					assert.NoError(t, s.Create(&orchestrator.CloudService{Id: testdata.MockCloudServiceID1}))
					assert.NoError(t, s.Create(orchestratortest.NewTargetOfEvaluation(testdata.AssuranceLevelBasic)))
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
				})))),
				scheduler:       make(map[string]*gocron.Scheduler),
				wg:              make(map[string]*sync.WaitGroup),
				authz:           &service.AuthorizationStrategyAllowAll{},
				storage:         testutil.NewInMemoryStorage(t),
				catalogControls: make(map[string]map[string]*orchestrator.Control),
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.StartEvaluationRequest{
					CloudServiceId: testdata.MockCloudServiceID1,
					CatalogId:      testdata.MockCatalogID,
				},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				s, ok := i1.(*Service)
				if !assert.True(tt, ok) {
					return false
				}

				gotResp, ok := i2[0].(*evaluation.StartEvaluationResponse)
				if !assert.True(tt, ok) {
					return false
				}
				assert.True(t, gotResp.Successful)

				toeTag := createSchedulerTag(testdata.MockCloudServiceID1, testdata.MockCatalogID)
				return assert.Equal(t, 2, len(s.scheduler[toeTag].Jobs()))
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				orchestrator:    tt.fields.orchestrator,
				authorizer:      tt.fields.authorizer,
				scheduler:       tt.fields.scheduler,
				storage:         tt.fields.storage,
				wg:              tt.fields.wg,
				authz:           tt.fields.authz,
				catalogControls: tt.fields.catalogControls,
			}

			toeTag := createSchedulerTag(tt.args.req.GetCloudServiceId(), tt.args.req.GetCatalogId())
			// Start the scheduler
			if len(tt.fields.scheduler) != 0 {
				tt.fields.scheduler[toeTag] = gocron.NewScheduler(time.UTC)
				_, err := s.scheduler[toeTag].Every(1).Day().Tag(createJobTag(tt.args.req.GetCloudServiceId(), tt.args.req.GetCatalogId(), testdata.MockCatalogID)).Do(func() { fmt.Println("Scheduler") })
				assert.NoError(t, err)
				s.scheduler[toeTag].StartAsync()
			}

			gotResp, err := s.StartEvaluation(tt.args.in0, tt.args.req)
			tt.wantErr(t, err)

			if tt.want != nil {
				tt.want(t, s, gotResp)
			}
		})
	}
}

func TestService_getAllMetricsFromControl(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestrator                  *api.RPCConnection[orchestrator.OrchestratorClient]
		authorizer                    api.Authorizer
		scheduler                     map[string]*gocron.Scheduler
		wg                            map[string]*sync.WaitGroup
		storage                       persistence.Storage
		catalogControls               map[string]map[string]*orchestrator.Control
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
			fields:      fields{},
			wantMetrics: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not get control for control id")
			},
		},
		{
			name: "no sub-controls available",
			fields: fields{
				catalogControls: map[string]map[string]*orchestrator.Control{
					orchestratortest.MockControl1.GetCategoryCatalogId(): {
						fmt.Sprintf("%s-%s", orchestratortest.MockControl6.GetCategoryName(), orchestratortest.MockControl6.GetId()): orchestratortest.MockControl6,
					},
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
				catalogControls: map[string]map[string]*orchestrator.Control{
					orchestratortest.MockControl1.GetCategoryCatalogId(): {
						fmt.Sprintf("%s-%s", orchestratortest.MockControl1.GetCategoryName(), orchestratortest.MockControl1.GetId()):   orchestratortest.MockControl1,
						fmt.Sprintf("%s-%s", orchestratortest.MockControl11.GetCategoryName(), orchestratortest.MockControl11.GetId()): orchestratortest.MockControl11,
					},
				},
			},
			args: args{
				catalogId:    testdata.MockCatalogID,
				categoryName: testdata.MockCategoryName,
				controlId:    testdata.MockControlID1,
			},
			wantMetrics: []*assessment.Metric{
				{
					Id:          testdata.MockMetricID1,
					Name:        testdata.MockMetricName1,
					Description: testdata.MockMetricDescription1,
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
				orchestrator:                  tt.fields.orchestrator,
				authorizer:                    tt.fields.authorizer,
				scheduler:                     tt.fields.scheduler,
				wg:                            tt.fields.wg,
				storage:                       tt.fields.storage,
				catalogControls:               tt.fields.catalogControls,
			}
			gotMetrics, err := s.getAllMetricsFromControl(tt.args.catalogId, tt.args.categoryName, tt.args.controlId)
			tt.wantErr(t, err)

			if assert.Equal(t, len(gotMetrics), len(tt.wantMetrics)) {
				for i := range gotMetrics {
					reflect.DeepEqual(gotMetrics[i], tt.wantMetrics[i])
				}
			}
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
						Id: testdata.MockSubControlID11,
					},
					{
						Id: testdata.MockSubControlID,
					},
				},
			},
			want: []string{testdata.MockSubControlID11, testdata.MockSubControlID},
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
		orchestrator                  *api.RPCConnection[orchestrator.OrchestratorClient]
		authorizer                    api.Authorizer
		scheduler                     map[string]*gocron.Scheduler
		wg                            map[string]*sync.WaitGroup
		storage                       persistence.Storage
		catalogControls               map[string]map[string]*orchestrator.Control
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
			name:   "catalog_id is missing",
			fields: fields{},
			args: args{
				categoryName: testdata.MockCategoryName,
				controlId:    testdata.MockControlID1,
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrCatalogIdIsMissing)
			},
		},
		{
			name:   "category_name is missing",
			fields: fields{},
			args: args{
				catalogId: testdata.MockCatalogID,
				controlId: testdata.MockControlID1,
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrCategoryNameIsMissing)
			},
		},
		{
			name:   "control_id is missing",
			fields: fields{},
			args: args{
				catalogId:    testdata.MockCatalogID,
				categoryName: testdata.MockCategoryName,
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrControlIdIsMissing)
			},
		},
		{
			name:   "control does not exist",
			fields: fields{},
			args: args{
				catalogId:    "wrong_catalog_id",
				categoryName: "wrong_category_id",
				controlId:    "wrong_control_id",
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrControlNotAvailable)
			},
		},
		{
			name: "Happy path",
			fields: fields{
				catalogControls: map[string]map[string]*orchestrator.Control{
					orchestratortest.MockControl1.GetCategoryCatalogId(): {
						fmt.Sprintf("%s-%s", orchestratortest.MockControl1.GetCategoryName(), orchestratortest.MockControl1.GetId()): orchestratortest.MockControl1,
						fmt.Sprintf("%s-%s", orchestratortest.MockControl1.GetCategoryName(), orchestratortest.MockControl2.GetId()): orchestratortest.MockControl2,
					},
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

				// We need to truncate the metric from the control because the control is only returned with its
				// sub-control but without the sub-control's metric.
				wantControl := orchestratortest.MockControl1
				tmpMetrics := wantControl.Controls[0].Metrics
				wantControl.Controls[0].Metrics = nil

				if !proto.Equal(gotControl, wantControl) {
					t.Errorf("Service.GetControl() = %v, want %v", gotControl, wantControl)
					wantControl.Controls[0].Metrics = tmpMetrics
					return false
				}

				wantControl.Controls[0].Metrics = tmpMetrics
				return true
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				UnimplementedEvaluationServer: tt.fields.UnimplementedEvaluationServer,
				orchestrator:                  tt.fields.orchestrator,
				authorizer:                    tt.fields.authorizer,
				scheduler:                     tt.fields.scheduler,
				wg:                            tt.fields.wg,
				storage:                       tt.fields.storage,
				catalogControls:               tt.fields.catalogControls,
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
		orchestrator                  *api.RPCConnection[orchestrator.OrchestratorClient]
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
					CloudServiceId: testdata.MockCloudServiceID1,
					CatalogId:      testdata.MockCatalogID,
					AssuranceLevel: &testdata.AssuranceLevelHigh,
				},
				parentSchedulerTag: testdata.MockCloudServiceID1 + "control_id",
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
				parentSchedulerTag: testdata.MockCloudServiceID1 + "control_id",
				toe: &orchestrator.TargetOfEvaluation{
					CloudServiceId: testdata.MockCloudServiceID1,
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
					CloudServiceId: testdata.MockCloudServiceID1,
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
					CloudServiceId: testdata.MockCloudServiceID1,
					CatalogId:      testdata.MockCatalogID,
					AssuranceLevel: &testdata.AssuranceLevelHigh,
				},
				parentSchedulerTag: testdata.MockCloudServiceID1 + "control_id",
				interval:           2,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				UnimplementedEvaluationServer: tt.fields.UnimplementedEvaluationServer,
				orchestrator:                  tt.fields.orchestrator,
				authorizer:                    tt.fields.authorizer,
				scheduler:                     tt.fields.scheduler,
				wg:                            tt.fields.wg,
				storage:                       tt.fields.storage,
			}

			toeTag := createSchedulerTag(tt.args.toe.GetCloudServiceId(), tt.args.toe.GetCatalogId())
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
		orchestrator                  *api.RPCConnection[orchestrator.OrchestratorClient]
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
					testdata.MockCloudServiceID1 + "-" + testdata.MockControlID1: {},
				},
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(evaluationtest.MockEvaluationResultsWithoutResultsForParentControl))
				}),
				authz: &service.AuthorizationStrategyJWT{},
			},
			args: args{
				toe: &orchestrator.TargetOfEvaluation{
					CloudServiceId: testdata.MockCloudServiceID1,
					CatalogId:      testdata.MockCatalogID,
					AssuranceLevel: &testdata.AssuranceLevelHigh,
				},
				categoryName: testdata.MockCategoryName,
				controlId:    testdata.MockControlID1,
				schedulerTag: testdata.MockCloudServiceID1 + "-" + testdata.MockControlID1,
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
					testdata.MockCloudServiceID1 + "-" + testdata.MockControlID1: {},
				},
				storage: testutil.NewInMemoryStorage(t),
				authz:   &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				toe: &orchestrator.TargetOfEvaluation{
					CloudServiceId: testdata.MockCloudServiceID1,
					CatalogId:      testdata.MockCatalogID,
					AssuranceLevel: &testdata.AssuranceLevelHigh,
				},
				categoryName: testdata.MockCategoryName,
				controlId:    testdata.MockControlID1,
				schedulerTag: testdata.MockCloudServiceID1 + "-" + testdata.MockControlID1,
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
					testdata.MockCloudServiceID1 + "-" + testdata.MockControlID1: {},
				},
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(evaluationtest.MockEvaluationResultsWithoutResultsForParentControl))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				toe: &orchestrator.TargetOfEvaluation{
					CloudServiceId: testdata.MockCloudServiceID1,
					CatalogId:      testdata.MockCatalogID,
					AssuranceLevel: &testdata.AssuranceLevelHigh,
				},
				categoryName: testdata.MockCategoryName,
				controlId:    testdata.MockControlID1,
				schedulerTag: testdata.MockCloudServiceID1 + "-" + testdata.MockControlID1,
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
				assert.Equal(t, 6, len(evalResults.Results))

				createdResult := evalResults.Results[len(evalResults.Results)-1]

				// Delete ID and timestamp from the evaluation results
				assert.NotEmpty(t, createdResult.GetId())
				createdResult.Id = ""
				createdResult.Timestamp = nil
				newEvalResults.Id = ""
				newEvalResults.Timestamp = nil
				return proto.Equal(newEvalResults, createdResult)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				UnimplementedEvaluationServer: tt.fields.UnimplementedEvaluationServer,
				orchestrator:                  tt.fields.orchestrator,
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
		orchestrator    *api.RPCConnection[orchestrator.OrchestratorClient]
		authorizer      api.Authorizer
		scheduler       map[string]*gocron.Scheduler
		wg              map[string]*sync.WaitGroup
		storage         persistence.Storage
		authz           service.AuthorizationStrategy
		schedulerTag    string
		wgCounter       int
		catalogControls map[string]map[string]*orchestrator.Control
	}
	type args struct {
		toe                *orchestrator.TargetOfEvaluation
		control            *orchestrator.Control
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
				schedulerTag: createJobTag(testdata.MockCloudServiceID1, testdata.MockCatalogID, testdata.MockControlID1),
				orchestrator: api.NewRPCConnection("bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
				})))),
				storage:   testutil.NewInMemoryStorage(t),
				authz:     &service.AuthorizationStrategyAllowAll{},
				wgCounter: 2,
				wg: map[string]*sync.WaitGroup{
					createJobTag(testdata.MockCloudServiceID1, testdata.MockCatalogID, testdata.MockControlID1): {},
				},
			},
			args: args{
				control: &orchestrator.Control{
					Id:           testdata.MockControlID1,
					CategoryName: testdata.MockCategoryName,
				},
				parentSchedulerTag: createJobTag(testdata.MockCloudServiceID1, testdata.MockCatalogID, testdata.MockControlID1),
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
				schedulerTag: createJobTag(testdata.MockCloudServiceID1, testdata.MockCatalogID, testdata.MockControlID1),
				wgCounter:    2,
				wg: map[string]*sync.WaitGroup{
					createJobTag(testdata.MockCloudServiceID1, testdata.MockCatalogID, testdata.MockControlID1): {},
				},
				storage: testutil.NewInMemoryStorage(t),
				authz:   &service.AuthorizationStrategyAllowAll{},
				orchestrator: api.NewRPCConnection("bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
				})))),
			},
			args: args{
				toe: &orchestrator.TargetOfEvaluation{
					CloudServiceId: testdata.MockCloudServiceID1,
					CatalogId:      testdata.MockCatalogID,
					AssuranceLevel: &testdata.AssuranceLevelHigh,
				},
				control: &orchestrator.Control{
					Id:           testdata.MockControlID1,
					CategoryName: testdata.MockCategoryName,
				},
				parentSchedulerTag: createJobTag(testdata.MockCloudServiceID1, testdata.MockCatalogID, testdata.MockControlID1),
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
				schedulerTag: createJobTag(testdata.MockCloudServiceID1, testdata.MockCatalogID, testdata.MockControlID1),
				wgCounter:    2,
				wg: map[string]*sync.WaitGroup{
					createJobTag(testdata.MockCloudServiceID1, testdata.MockCatalogID, testdata.MockControlID1): {},
				},
				storage:      testutil.NewInMemoryStorage(t),
				authz:        &service.AuthorizationStrategyAllowAll{},
				orchestrator: api.NewRPCConnection("bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t)))),
			},
			args: args{
				toe: &orchestrator.TargetOfEvaluation{
					CloudServiceId: testdata.MockCloudServiceID1,
					CatalogId:      testdata.MockCatalogID,
					AssuranceLevel: &testdata.AssuranceLevelHigh,
				},
				control: &orchestrator.Control{
					Id:           testdata.MockSubControlID11,
					CategoryName: testdata.MockCategoryName,
				},
				parentSchedulerTag: createJobTag(testdata.MockCloudServiceID1, testdata.MockCatalogID, testdata.MockControlID1),
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
			name: "error getting assessment results",
			fields: fields{
				schedulerTag: createJobTag(testdata.MockCloudServiceID1, testdata.MockCatalogID, testdata.MockControlID1),
				wgCounter:    1,
				wg: map[string]*sync.WaitGroup{
					createJobTag(testdata.MockCloudServiceID1, testdata.MockCatalogID, testdata.MockControlID1): {},
				},
				storage: testutil.NewInMemoryStorage(t),
				authz:   &service.AuthorizationStrategyAllowAll{},
				orchestrator: api.NewRPCConnection("bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
				})))),
			},
			args: args{
				toe: &orchestrator.TargetOfEvaluation{
					CatalogId:      testdata.MockCatalogID,
					AssuranceLevel: &testdata.AssuranceLevelHigh,
				},
				control: &orchestrator.Control{
					Id:           testdata.MockSubControlID11,
					CategoryName: testdata.MockCategoryName,
				},
				parentSchedulerTag: createJobTag(testdata.MockCloudServiceID1, testdata.MockCatalogID, testdata.MockControlID1),
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
				schedulerTag: createJobTag(testdata.MockCloudServiceID1, testdata.MockCatalogID, testdata.MockControlID1),
				wgCounter:    1,
				wg: map[string]*sync.WaitGroup{
					createJobTag(testdata.MockCloudServiceID1, testdata.MockCatalogID, testdata.MockControlID1): {},
				},
				storage: testutil.NewInMemoryStorage(t),
				authz:   &service.AuthorizationStrategyAllowAll{},
				orchestrator: api.NewRPCConnection("bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
				})))),
				catalogControls: map[string]map[string]*orchestrator.Control{
					orchestratortest.MockControl1.GetCategoryCatalogId(): {
						fmt.Sprintf("%s-%s", orchestratortest.MockControl1.GetCategoryName(), orchestratortest.MockControl1.GetId()):   orchestratortest.MockControl1,
						fmt.Sprintf("%s-%s", orchestratortest.MockControl11.GetCategoryName(), orchestratortest.MockControl11.GetId()): orchestratortest.MockControl11,
					},
				},
			},
			args: args{
				toe: &orchestrator.TargetOfEvaluation{
					CloudServiceId: testdata.MockCloudServiceID1,
					CatalogId:      testdata.MockCatalogID,
					AssuranceLevel: &testdata.AssuranceLevelHigh,
				},
				control: &orchestrator.Control{
					Id:           testdata.MockSubControlID11,
					CategoryName: testdata.MockCategoryName,
				},
				parentSchedulerTag: createJobTag(testdata.MockCloudServiceID1, testdata.MockCatalogID, testdata.MockControlID1),
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
				orchestrator:    tt.fields.orchestrator,
				authorizer:      tt.fields.authorizer,
				scheduler:       tt.fields.scheduler,
				wg:              tt.fields.wg,
				storage:         tt.fields.storage,
				authz:           tt.fields.authz,
				catalogControls: tt.fields.catalogControls,
			}

			tt.fields.wg[tt.fields.schedulerTag].Add(tt.fields.wgCounter)
			s.evaluateSubcontrol(tt.args.toe, tt.args.control, tt.args.parentSchedulerTag)

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
				cloudServiceId: testdata.MockCloudServiceID1,
				catalogId:      testdata.MockCatalogID,
			},
			want: "",
		},
		{
			name: "catalog_id empty",
			args: args{
				cloudServiceId: testdata.MockCloudServiceID1,
				controlId:      testdata.MockControlID1,
			},
			want: "",
		},
		{
			name: "Happy path",
			args: args{
				cloudServiceId: testdata.MockCloudServiceID1,
				controlId:      testdata.MockControlID1,
				catalogId:      testdata.MockCatalogID,
			},
			want: fmt.Sprintf("%s-%s-%s", testdata.MockCloudServiceID1, testdata.MockCatalogID, testdata.MockControlID1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createJobTag(tt.args.cloudServiceId, tt.args.catalogId, tt.args.controlId); got != tt.want {
				t.Errorf("getSchedulerTag() = %v, want %v", got, tt.want)
			}
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
					testdata.MockResourceID1: {
						orchestratortest.MockAssessmentResult1,
						orchestratortest.MockAssessmentResult2,
						orchestratortest.MockAssessmentResult3,
					},
					testdata.MockResourceID2: {
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
			got := createAssessmentResultMap(tt.args.results)

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
			gotControlsHierarchy := createControlsInScopeHierarchy(tt.args.controls)

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
				cloudServiceId: testdata.MockCloudServiceID1,
			},
			want: "",
		},
		{
			name: "Happy path",
			args: args{
				cloudServiceId: testdata.MockCloudServiceID1,
				catalogId:      testdata.MockCatalogID,
			},
			want: fmt.Sprintf("%s-%s", testdata.MockCloudServiceID1, testdata.MockCatalogID),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createSchedulerTag(tt.args.cloudServiceId, tt.args.catalogId); got != tt.want {
				t.Errorf("getToeTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_cacheControls(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestrator                  *api.RPCConnection[orchestrator.OrchestratorClient]
		authorizer                    api.Authorizer
		scheduler                     map[string]*gocron.Scheduler
		wg                            map[string]*sync.WaitGroup
		authz                         service.AuthorizationStrategy
		storage                       persistence.Storage
		catalogControls               map[string]map[string]*orchestrator.Control
	}
	type args struct {
		catalogId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.ValueAssertionFunc
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "catalog_id missing",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrCatalogIdIsMissing)
			},
		},
		{
			name: "initOrchestratorClient fails",
			fields: fields{
				orchestrator: api.NewRPCConnection("bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(connectionRefusedDialer)),
			},
			args: args{
				catalogId: testdata.MockCatalogID,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "connection refused")
			},
		},
		{
			name: "no controls available",
			fields: fields{
				orchestrator:    api.NewRPCConnection("bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {})))),
				catalogControls: make(map[string]map[string]*orchestrator.Control),
			},
			args: args{
				catalogId: testdata.MockCatalogID,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, fmt.Sprintf("no controls for catalog '%s' available", testdata.MockCatalogID))
			},
		},
		{
			name: "Happy path",
			fields: fields{
				orchestrator: api.NewRPCConnection("bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
				})))),
				catalogControls: make(map[string]map[string]*orchestrator.Control),
			},
			args: args{
				catalogId: testdata.MockCatalogID,
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				service, ok := i1.(*Service)
				if !assert.True(tt, ok) {
					return false
				}

				assert.Equal(t, 1, len(service.catalogControls))
				return assert.Equal(t, 4, len(service.catalogControls[testdata.MockCatalogID]))
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				UnimplementedEvaluationServer: tt.fields.UnimplementedEvaluationServer,
				orchestrator:                  tt.fields.orchestrator,
				authorizer:                    tt.fields.authorizer,
				scheduler:                     tt.fields.scheduler,
				wg:                            tt.fields.wg,
				authz:                         tt.fields.authz,
				storage:                       tt.fields.storage,
				catalogControls:               tt.fields.catalogControls,
			}
			err := s.cacheControls(tt.args.catalogId)
			tt.wantErr(t, err)

			if tt.want != nil {
				tt.want(t, s)
			}
		})
	}
}
