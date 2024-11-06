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
	"testing"
	"time"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/evaluation"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/clitest"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/evaluationtest"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/orchestratortest"
	"clouditor.io/clouditor/v2/internal/util"
	"clouditor.io/clouditor/v2/launcher"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/service"
	"github.com/go-co-op/gocron"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
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
		opts []service.Option[*Service]
	}
	tests := []struct {
		name string
		args args
		want assert.Want[*Service]
	}{
		{
			name: "WithStorage",
			args: args{
				opts: []service.Option[*Service]{WithStorage(inmem)},
			},
			want: func(t *testing.T, got *Service) bool {
				return assert.Same(t, inmem, got.storage)
			},
		},
		{
			name: "WithOrchestratorAddress",
			args: args{
				opts: []service.Option[*Service]{WithOrchestratorAddress(testdata.MockGRPCTarget)},
			},
			want: func(t *testing.T, got *Service) bool {
				return assert.Equal(t, testdata.MockGRPCTarget, got.orchestrator.Target)
			},
		},
		{
			name: "WithOAuth2Authorizer",
			args: args{
				opts: []service.Option[*Service]{WithOAuth2Authorizer(&clientcredentials.Config{})},
			},
			want: func(t *testing.T, got *Service) bool {
				return assert.Equal(t, api.NewOAuthAuthorizerFromClientCredentials(&clientcredentials.Config{}), got.orchestrator.Authorizer(), assert.CompareAllUnexported())
			},
		},
		{
			name: "WithAuthorizer",
			args: args{
				opts: []service.Option[*Service]{WithAuthorizer(api.NewOAuthAuthorizerFromClientCredentials(&clientcredentials.Config{}))},
			},
			want: func(t *testing.T, got *Service) bool {
				return assert.Equal(t, api.NewOAuthAuthorizerFromClientCredentials(&clientcredentials.Config{}), got.orchestrator.Authorizer(), assert.CompareAllUnexported())
			},
		},
		{
			name: "Happy path: without additional services",
			args: args{
				opts: []service.Option[*Service]{},
			},
			want: func(t *testing.T, got *Service) bool {
				return assert.Equal(t, DefaultOrchestratorAddress, got.orchestrator.Target)
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

func TestService_ListEvaluationResults(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestrator                  *api.RPCConnection[orchestrator.OrchestratorClient]
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
			name: "Permission denied",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(evaluationtest.MockEvaluationResults))
				}),
				authz: servicetest.NewAuthorizationStrategy(false),
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.ListEvaluationResultsRequest{
					LatestByControlId: util.Ref(true),
					Filter: &evaluation.ListEvaluationResultsRequest_Filter{
						ControlId:             util.Ref(testdata.MockSubControlID11),
						SubControls:           util.Ref(testdata.MockControlID1),
						CertificationTargetId: util.Ref(testdata.MockCertificationTargetID1),
					},
				},
			},
			wantRes: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.PermissionDenied, status.Code(err))
				return assert.ErrorContains(t, err, service.ErrPermissionDenied.Error())
			},
		},
		{
			name: "Database error",
			fields: fields{
				storage: &testutil.StorageWithError{RawErr: gorm.ErrInvalidDB},
				authz:   &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.ListEvaluationResultsRequest{
					LatestByControlId: util.Ref(true),
					Filter: &evaluation.ListEvaluationResultsRequest_Filter{
						ControlId:             util.Ref(testdata.MockSubControlID11),
						SubControls:           util.Ref(testdata.MockControlID1),
						CertificationTargetId: util.Ref(testdata.MockCertificationTargetID1),
					},
				},
			},
			wantRes: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.Internal, status.Code(err))
				return assert.ErrorContains(t, err, persistence.ErrDatabase.Error())
			},
		},
		{
			name: "Paginate error",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(evaluationtest.MockEvaluationResults))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.ListEvaluationResultsRequest{
					LatestByControlId: util.Ref(false),
					OrderBy:           "Wrong input",
					Filter: &evaluation.ListEvaluationResultsRequest_Filter{
						ControlId:             util.Ref(testdata.MockSubControlID11),
						SubControls:           util.Ref(testdata.MockControlID1),
						CertificationTargetId: util.Ref(testdata.MockCertificationTargetID1),
					},
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.Internal, status.Code(err))
				return assert.ErrorContains(t, err, "could not paginate results")
			},
		},
		{
			name: "Filter latest_by_control_id, control_id, sub_controls, certification_target_id",
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
						ControlId:             util.Ref(testdata.MockSubControlID11),
						SubControls:           util.Ref(testdata.MockControlID1),
						CertificationTargetId: util.Ref(testdata.MockCertificationTargetID1),
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
			name: "Filter certification_target_id",
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
						CertificationTargetId: util.Ref(testdata.MockCertificationTargetID1),
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
				storage:                       tt.fields.storage,
				authz:                         tt.fields.authz,
			}
			gotRes, err := s.ListEvaluationResults(tt.args.in0, tt.args.req)

			tt.wantErr(t, err)
			assert.Equal(t, tt.wantRes, gotRes)
		})
	}
}

func TestService_getMetricsFromSubControls(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestrator                  *api.RPCConnection[orchestrator.OrchestratorClient]
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
					CategoryCatalogId: testdata.MockCatalogID1,
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
				storage:                       tt.fields.storage,
				catalogControls:               tt.fields.catalogControls,
			}
			gotMetrics, err := s.getMetricsFromSubcontrols(tt.args.control)

			tt.wantErr(t, err)

			assert.Equal(t, len(gotMetrics), len(tt.wantMetrics))
			for i := range gotMetrics {
				assert.Equal(t, tt.wantMetrics[i], gotMetrics[i])
			}
		})
	}
}

func TestService_StopEvaluation(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestrator                  *api.RPCConnection[orchestrator.OrchestratorClient]
		scheduler                     *gocron.Scheduler
		storage                       persistence.Storage
		authz                         service.AuthorizationStrategy
		auditScopeTag                 string
	}
	type args struct {
		in0              context.Context
		req              *evaluation.StopEvaluationRequest
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
			name: "error: request input missing",
			args: args{
				in0: context.Background(),
				req: &evaluation.StopEvaluationRequest{},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "audit_scope_id: value is empty, which is not a valid UUID")
			},
		},
		{
			name: "error: not authorized",
			fields: fields{
				authz:   servicetest.NewAuthorizationStrategy(false),
				storage: testutil.NewInMemoryStorage(t),
				orchestrator: api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Create(orchestratortest.MockAuditScopeCertTargetID1)
					assert.NoError(t, err)
				})))),
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.StopEvaluationRequest{
					AuditScopeId: testdata.MockAuditScopeID1,
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, service.ErrPermissionDenied)
			},
		},
		{
			name: "error: getting audit scope",
			fields: fields{
				authz:        servicetest.NewAuthorizationStrategy(true),
				storage:      testutil.NewInMemoryStorage(t),
				orchestrator: api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {})))),
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.StopEvaluationRequest{
					AuditScopeId: testdata.MockAuditScopeID1,
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.Internal, status.Code(err))
				return assert.ErrorContains(t, err, ErrAuditScopeNotFound.Error())
			},
		},
		{
			name: "Evaluation not running",
			args: args{
				in0: context.Background(),
				req: &evaluation.StopEvaluationRequest{
					AuditScopeId: testdata.MockAuditScopeID1,
				},
				schedulerRunning: false,
			},
			fields: fields{
				scheduler: gocron.NewScheduler(time.Local),
				orchestrator: api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewAuditScope(testdata.AssuranceLevelBasic, testdata.MockAuditScopeID1, testdata.MockCertificationTargetID1)))
				})))),
				authz:         &service.AuthorizationStrategyAllowAll{},
				auditScopeTag: fmt.Sprintf("%s-%s", testdata.MockCertificationTargetID1, testdata.MockCatalogID1),
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.MockAuditScopeCertTargetID1))
				}),
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, fmt.Sprintf("job for audit scope '%s' not running", testdata.MockAuditScopeID1))
			},
		},
		{
			name: "Happy path",
			args: args{
				in0: context.Background(),
				req: &evaluation.StopEvaluationRequest{
					AuditScopeId: testdata.MockAuditScopeID1,
				},
				schedulerRunning: true,
			},
			fields: fields{
				orchestrator: api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewAuditScope(testdata.AssuranceLevelBasic, testdata.MockAuditScopeID1, testdata.MockCertificationTargetID1)))
				})))),

				scheduler: func() *gocron.Scheduler {
					s := gocron.NewScheduler(time.UTC)
					_, err := s.Every(1).
						Day().
						Tag(testdata.MockAuditScopeID1).
						Do(func() { fmt.Println("Scheduler") })
					assert.NoError(t, err)

					return s
				}(),
				authz:         &service.AuthorizationStrategyAllowAll{},
				auditScopeTag: fmt.Sprintf("%s-%s", testdata.MockCertificationTargetID1, testdata.MockCatalogID1),
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
				scheduler:                     tt.fields.scheduler,
				storage:                       tt.fields.storage,
				authz:                         tt.fields.authz,
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
		scheduler       *gocron.Scheduler
		storage         persistence.Storage
		authz           service.AuthorizationStrategy
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
		want    assert.Want[*evaluation.StartEvaluationResponse]
		wantSvc assert.Want[*Service]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "error: request input missing",
			args: args{
				in0: context.Background(),
				req: &evaluation.StartEvaluationRequest{},
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "audit_scope_id: value is empty, which is not a valid UUID")
			},
		},
		{
			name: "error: not authorized",
			fields: fields{
				authz:   &service.AuthorizationStrategyJWT{},
				storage: testutil.NewInMemoryStorage(t),
				orchestrator: api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Create(orchestratortest.MockAuditScopeCertTargetID1)
					assert.NoError(t, err)
				})))),
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.StartEvaluationRequest{
					AuditScopeId: testdata.MockAuditScopeID1,
					Interval:     proto.Int32(5),
				},
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, service.ErrPermissionDenied)
			},
		},
		{
			name: "error: init orchestrator client",
			fields: fields{
				authz:        &service.AuthorizationStrategyAllowAll{},
				scheduler:    gocron.NewScheduler(time.Local),
				orchestrator: api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(connectionRefusedDialer)),
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.StartEvaluationRequest{
					AuditScopeId: testdata.MockAuditScopeID1,
					Interval:     proto.Int32(5),
				},
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "connection refused")
			},
		},
		{
			name: "error: get Audit Scope",
			fields: fields{
				orchestrator: api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
				})))),
				scheduler:       gocron.NewScheduler(time.Local),
				authz:           &service.AuthorizationStrategyAllowAll{},
				catalogControls: make(map[string]map[string]*orchestrator.Control),
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.StartEvaluationRequest{
					AuditScopeId: testdata.MockAuditScopeID1,
					Interval:     proto.Int32(5),
				},
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrAuditScopeNotFound.Error())
			},
		},
		{
			name: "error: getting cache controls",
			fields: fields{
				orchestrator: api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Create(orchestratortest.MockAuditScopeCertTargetID1)
					assert.NoError(t, err)
				})))),
				scheduler:       gocron.NewScheduler(time.Local),
				authz:           &service.AuthorizationStrategyAllowAll{},
				catalogControls: make(map[string]map[string]*orchestrator.Control),
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.StartEvaluationRequest{
					AuditScopeId: testdata.MockAuditScopeID1,
					Interval:     proto.Int32(5),
				},
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not cache controls:")
			},
		},
		{
			name: "error: scheduler for catalog started already",
			fields: fields{
				orchestrator: api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
					assert.NoError(t, s.Create(orchestratortest.NewAuditScope(testdata.AssuranceLevelBasic, testdata.MockAuditScopeID1, testdata.MockCertificationTargetID1)))
				})))),
				scheduler: func() *gocron.Scheduler {
					s := gocron.NewScheduler(time.Local)
					_, err := s.Every(1).
						Day().
						Tag(testdata.MockAuditScopeID1).
						Do(func() { fmt.Println("Scheduler") })
					assert.NoError(t, err)

					return s
				}(),
				authz:           &service.AuthorizationStrategyAllowAll{},
				catalogControls: make(map[string]map[string]*orchestrator.Control),
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.StartEvaluationRequest{
					AuditScopeId: testdata.MockAuditScopeID1,
					Interval:     proto.Int32(5),
				},
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "code = AlreadyExists desc = evaluation for Audit Scope")
			},
		},
		{
			name: "Happy path: scheduler added",
			fields: fields{
				orchestrator: api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
					assert.NoError(t, s.Create(orchestratortest.NewAuditScope(testdata.AssuranceLevelBasic, testdata.MockAuditScopeID1, testdata.MockCertificationTargetID1)))
				})))),
				scheduler:       gocron.NewScheduler(time.Local),
				authz:           &service.AuthorizationStrategyAllowAll{},
				storage:         testutil.NewInMemoryStorage(t),
				catalogControls: make(map[string]map[string]*orchestrator.Control),
			},
			args: args{
				in0: context.Background(),
				req: &evaluation.StartEvaluationRequest{
					AuditScopeId: testdata.MockAuditScopeID1,
				},
			},
			want: func(t *testing.T, got *evaluation.StartEvaluationResponse) bool {
				return assert.True(t, got.Successful)
			},
			wantSvc: func(t *testing.T, got *Service) bool {
				return assert.Equal(t, 1, len(got.scheduler.Jobs()))
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				orchestrator:    tt.fields.orchestrator,
				scheduler:       tt.fields.scheduler,
				storage:         tt.fields.storage,
				authz:           tt.fields.authz,
				catalogControls: tt.fields.catalogControls,
			}

			gotResp, err := svc.StartEvaluation(tt.args.in0, tt.args.req)
			tt.wantErr(t, err)
			assert.Optional(t, tt.want, gotResp)
			assert.Optional(t, tt.wantSvc, svc)
		})
	}
}

func TestService_getAllMetricsFromControl(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestrator                  *api.RPCConnection[orchestrator.OrchestratorClient]
		scheduler                     *gocron.Scheduler
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
				catalogId:    testdata.MockCatalogID1,
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
				catalogId:    testdata.MockCatalogID1,
				categoryName: testdata.MockCategoryName,
				controlId:    testdata.MockControlID1,
			},
			wantMetrics: []*assessment.Metric{
				{
					Id:          testdata.MockMetricID1,
					Category:    testdata.MockMetricCategory1,
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
				scheduler:                     tt.fields.scheduler,
				storage:                       tt.fields.storage,
				catalogControls:               tt.fields.catalogControls,
			}
			gotMetrics, err := s.getAllMetricsFromControl(tt.args.catalogId, tt.args.categoryName, tt.args.controlId)
			tt.wantErr(t, err)

			if assert.Equal(t, len(gotMetrics), len(tt.wantMetrics)) {
				for i := range gotMetrics {
					assert.Equal(t, tt.wantMetrics[i], gotMetrics[i])
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
			got := getMetricIds(tt.args.metrics)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_getControl(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestrator                  *api.RPCConnection[orchestrator.OrchestratorClient]
		scheduler                     *gocron.Scheduler
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
		want    assert.Want[*orchestrator.Control]
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
				catalogId: testdata.MockCatalogID1,
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
				catalogId:    testdata.MockCatalogID1,
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
				catalogId:    testdata.MockCatalogID1,
				categoryName: testdata.MockCategoryName,
				controlId:    testdata.MockControlID1,
			},
			want: func(t *testing.T, got *orchestrator.Control) bool {
				// We need to truncate the metric from the control because the control is only returned with its
				// sub-control but without the sub-control's metric.
				// TODO(oxisto): Use ignore fields instead
				wantControl := orchestratortest.MockControl1
				tmpMetrics := wantControl.Controls[0].Metrics
				wantControl.Controls[0].Metrics = nil

				if !assert.Equal(t, wantControl, got) {
					t.Errorf("Service.GetControl() = %v, want %v", got, wantControl)
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
				scheduler:                     tt.fields.scheduler,
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
		scheduler                     *gocron.Scheduler
		storage                       persistence.Storage
	}
	type args struct {
		ctx        context.Context
		auditScope *orchestrator.AuditScope
		catalog    *orchestrator.Catalog
		interval   int
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
				scheduler: gocron.NewScheduler(time.Local),
			},
			args: args{
				auditScope: &orchestrator.AuditScope{
					CertificationTargetId: testdata.MockCertificationTargetID1,
					CatalogId:             testdata.MockCatalogID1,
					AssuranceLevel:        &testdata.AssuranceLevelHigh,
				},
				catalog:  orchestratortest.NewCatalog(),
				interval: 0,
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "evaluation cannot be scheduled")
			},
		},
		{
			name: "Audit Scope input empty",
			args: args{
				auditScope: &orchestrator.AuditScope{
					CertificationTargetId: testdata.MockCertificationTargetID1,
					CatalogId:             testdata.MockCatalogID1,
					AssuranceLevel:        &testdata.AssuranceLevelHigh,
				},
				catalog: orchestratortest.NewCatalog(),
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "interval is invalid")
			},
		},
		{
			name: "Happy path: Add scheduler job for catalog",
			fields: fields{
				scheduler: gocron.NewScheduler(time.Local),
			},
			args: args{
				auditScope: &orchestrator.AuditScope{
					CertificationTargetId: testdata.MockCertificationTargetID1,
					CatalogId:             testdata.MockCatalogID1,
					AssuranceLevel:        &testdata.AssuranceLevelHigh,
				},
				catalog:  orchestratortest.NewCatalog(),
				interval: 2,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				UnimplementedEvaluationServer: tt.fields.UnimplementedEvaluationServer,
				orchestrator:                  tt.fields.orchestrator,
				scheduler:                     tt.fields.scheduler,
				storage:                       tt.fields.storage,
			}
			err := s.addJobToScheduler(tt.args.ctx, tt.args.auditScope, tt.args.catalog, tt.args.interval)
			tt.wantErr(t, err)
		})
	}
}

func TestService_evaluateControl(t *testing.T) {
	type fields struct {
		orchestrator    *api.RPCConnection[orchestrator.OrchestratorClient]
		scheduler       *gocron.Scheduler
		storage         persistence.Storage
		authz           service.AuthorizationStrategy
		catalogControls map[string]map[string]*orchestrator.Control
	}
	type args struct {
		ctx        context.Context
		auditScope *orchestrator.AuditScope
		catalog    *orchestrator.Catalog
		control    *orchestrator.Control
		manual     []*evaluation.EvaluationResult
	}
	tests := []struct {
		name                 string
		fields               fields
		args                 args
		newEvaluationResults *evaluation.EvaluationResult
		want                 func(t *testing.T, gotSvc *Service, gotResult *evaluation.EvaluationResult) bool
	}{
		{
			name: "AuthZ error in ListEvaluationResults",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(evaluationtest.MockEvaluationResultsWithoutResultsForParentControl))
				}),
				authz: &service.AuthorizationStrategyJWT{},
			},
			args: args{
				ctx: context.Background(),
				auditScope: &orchestrator.AuditScope{
					CertificationTargetId: testdata.MockCertificationTargetID1,
					CatalogId:             testdata.MockCatalogID1,
					AssuranceLevel:        &testdata.AssuranceLevelHigh,
				},
				catalog: orchestratortest.NewCatalog(),
				control: orchestratortest.MockControl1,
			},
			newEvaluationResults: &evaluation.EvaluationResult{},
			want: func(t *testing.T, gotSvc *Service, gotResult *evaluation.EvaluationResult) bool {
				evalResults, err := gotSvc.ListEvaluationResults(context.Background(), &evaluation.ListEvaluationResultsRequest{})
				assert.NoError(t, err)
				return assert.Equal(t, 0, len(evalResults.Results))
			},
		},
		{
			name: "No assessment results for evaluation available",
			fields: fields{
				storage:      testutil.NewInMemoryStorage(t),
				authz:        &service.AuthorizationStrategyAllowAll{},
				orchestrator: api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t)))),
				catalogControls: map[string]map[string]*orchestrator.Control{
					testdata.MockCatalogID1: {
						testdata.MockCategoryName + "-" + testdata.MockControlID1:     orchestratortest.MockControl1,
						testdata.MockCategoryName + "-" + testdata.MockSubControlID11: orchestratortest.MockControl11,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				auditScope: &orchestrator.AuditScope{
					CertificationTargetId: testdata.MockCertificationTargetID1,
					CatalogId:             testdata.MockCatalogID1,
					AssuranceLevel:        &testdata.AssuranceLevelHigh,
				},
				catalog: orchestratortest.NewCatalog(),
				control: orchestratortest.MockControl1,
			},
			newEvaluationResults: nil,
			want: func(t *testing.T, gotSvc *Service, gotResult *evaluation.EvaluationResult) bool {
				res, err := gotSvc.ListEvaluationResults(context.Background(), &evaluation.ListEvaluationResultsRequest{})
				return assert.NoError(t, err) &&
					assert.Equal(t, 2, len(res.Results)) &&
					assert.Equal(t, evaluation.EvaluationStatus_EVALUATION_STATUS_PENDING, res.Results[0].Status) &&
					assert.Equal(t, evaluation.EvaluationStatus_EVALUATION_STATUS_PENDING, res.Results[1].Status)
			},
		},
		{
			name: "Happy path",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   &service.AuthorizationStrategyAllowAll{},
				orchestrator: api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
				})))),
				catalogControls: map[string]map[string]*orchestrator.Control{
					testdata.MockCatalogID1: {
						testdata.MockCategoryName + "-" + testdata.MockControlID1:     orchestratortest.MockControl1,
						testdata.MockCategoryName + "-" + testdata.MockSubControlID11: orchestratortest.MockControl1.Controls[0],
					},
				},
			},
			args: args{
				ctx: context.Background(),
				auditScope: &orchestrator.AuditScope{
					CertificationTargetId: testdata.MockCertificationTargetID1,
					CatalogId:             testdata.MockCatalogID1,
					AssuranceLevel:        &testdata.AssuranceLevelHigh,
				},
				catalog: orchestratortest.NewCatalog(),
				control: orchestratortest.MockControl1,
				manual: []*evaluation.EvaluationResult{
					{
						Status: evaluation.EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT_MANUALLY,
					},
				},
			},
			newEvaluationResults: evaluationtest.MockEvaluationResult1,
			want: func(t *testing.T, gotSvc *Service, gotResult *evaluation.EvaluationResult) bool {
				evalResults, err := gotSvc.ListEvaluationResults(context.Background(), &evaluation.ListEvaluationResultsRequest{})
				assert.NoError(t, err)
				assert.Equal(t, 2, len(evalResults.Results))

				createdResult := evalResults.Results[len(evalResults.Results)-1]

				// Compare without ID and timestamp since they are random
				return assert.NotEmpty(t, gotResult.Id) &&
					assert.NotNil(t, gotResult.Timestamp) &&
					assert.Equal(t, createdResult, gotResult,
						protocmp.IgnoreFields(&evaluation.EvaluationResult{}, "id", "timestamp"),
					)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				orchestrator:    tt.fields.orchestrator,
				scheduler:       tt.fields.scheduler,
				storage:         tt.fields.storage,
				authz:           tt.fields.authz,
				catalogControls: tt.fields.catalogControls,
			}

			_ = s.evaluateControl(tt.args.ctx, tt.args.auditScope, tt.args.catalog, tt.args.control, tt.args.manual)

			tt.want(t, s, tt.newEvaluationResults)
		})
	}
}

func TestService_evaluateSubcontrol(t *testing.T) {
	type fields struct {
		orchestrator    *api.RPCConnection[orchestrator.OrchestratorClient]
		scheduler       *gocron.Scheduler
		storage         persistence.Storage
		authz           service.AuthorizationStrategy
		wgCounter       int
		catalogControls map[string]map[string]*orchestrator.Control
	}
	type args struct {
		ctx        context.Context
		auditScope *orchestrator.AuditScope
		control    *orchestrator.Control
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.Want[*evaluation.EvaluationResult]
		wantSvc assert.Want[*Service]
	}{
		{
			name: "Audit Scope input empty", // we do not check the other input parameters
			fields: fields{
				orchestrator: api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
				})))),
				storage:   testutil.NewInMemoryStorage(t),
				authz:     &service.AuthorizationStrategyAllowAll{},
				wgCounter: 2,
			},
			args: args{
				control: &orchestrator.Control{
					Id:           testdata.MockControlID1,
					CategoryName: testdata.MockCategoryName,
				},
			},
			wantSvc: func(t *testing.T, got *Service) bool {
				evalResults, err := got.ListEvaluationResults(context.Background(), &evaluation.ListEvaluationResultsRequest{})
				assert.NoError(t, err)
				return assert.Equal(t, 0, len(evalResults.Results))
			},
		},
		{
			name: "no assessment results available",
			fields: fields{
				wgCounter: 2,
				storage:   testutil.NewInMemoryStorage(t),
				authz:     &service.AuthorizationStrategyAllowAll{},
				orchestrator: api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
				})))),
			},
			args: args{
				auditScope: &orchestrator.AuditScope{
					CertificationTargetId: testdata.MockCertificationTargetID1,
					CatalogId:             testdata.MockCatalogID1,
					AssuranceLevel:        &testdata.AssuranceLevelHigh,
				},
				control: &orchestrator.Control{
					Id:           testdata.MockControlID1,
					CategoryName: testdata.MockCategoryName,
				},
			},
			wantSvc: func(t *testing.T, got *Service) bool {
				evalResults, err := got.ListEvaluationResults(context.Background(), &evaluation.ListEvaluationResultsRequest{})
				assert.NoError(t, err)
				return assert.Equal(t, 0, len(evalResults.Results))
			},
		},
		{
			name: "error getting metrics",
			fields: fields{
				wgCounter:    2,
				storage:      testutil.NewInMemoryStorage(t),
				authz:        &service.AuthorizationStrategyAllowAll{},
				orchestrator: api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t)))),
			},
			args: args{
				auditScope: &orchestrator.AuditScope{
					CertificationTargetId: testdata.MockCertificationTargetID1,
					CatalogId:             testdata.MockCatalogID1,
					AssuranceLevel:        &testdata.AssuranceLevelHigh,
				},
				control: &orchestrator.Control{
					Id:           testdata.MockSubControlID11,
					CategoryName: testdata.MockCategoryName,
				},
			},
			wantSvc: func(t *testing.T, got *Service) bool {
				evalResults, err := got.ListEvaluationResults(context.Background(), &evaluation.ListEvaluationResultsRequest{})
				assert.NoError(t, err)
				return assert.Equal(t, 0, len(evalResults.Results))
			},
		},
		{
			name: "error getting assessment results",
			fields: fields{
				wgCounter: 1,
				storage:   testutil.NewInMemoryStorage(t),
				authz:     &service.AuthorizationStrategyAllowAll{},
				orchestrator: api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
				})))),
			},
			args: args{
				auditScope: &orchestrator.AuditScope{
					CatalogId:      testdata.MockCatalogID1,
					AssuranceLevel: &testdata.AssuranceLevelHigh,
				},
				control: &orchestrator.Control{
					Id:           testdata.MockSubControlID11,
					CategoryName: testdata.MockCategoryName,
				},
			},
			wantSvc: func(t *testing.T, got *Service) bool {
				evalResults, err := got.ListEvaluationResults(context.Background(), &evaluation.ListEvaluationResultsRequest{})
				assert.NoError(t, err)
				return assert.Equal(t, 0, len(evalResults.Results))
			},
		},
		{
			name: "Happy path",
			fields: fields{
				wgCounter: 1,
				storage:   testutil.NewInMemoryStorage(t),
				authz:     &service.AuthorizationStrategyAllowAll{},
				orchestrator: api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
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
				auditScope: &orchestrator.AuditScope{
					CertificationTargetId: testdata.MockCertificationTargetID1,
					CatalogId:             testdata.MockCatalogID1,
					AssuranceLevel:        &testdata.AssuranceLevelHigh,
				},
				control: &orchestrator.Control{
					Id:           testdata.MockSubControlID11,
					CategoryName: testdata.MockCategoryName,
				},
			},
			wantSvc: func(t *testing.T, got *Service) bool {
				evalResults, err := got.ListEvaluationResults(context.Background(), &evaluation.ListEvaluationResultsRequest{})
				assert.NoError(t, err)
				return assert.Equal(t, 1, len(evalResults.Results))
			},
			want: func(t *testing.T, got *evaluation.EvaluationResult) bool {
				return assert.Equal(t, testdata.MockSubControlID11, got.ControlId)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				orchestrator:    tt.fields.orchestrator,
				scheduler:       tt.fields.scheduler,
				storage:         tt.fields.storage,
				authz:           tt.fields.authz,
				catalogControls: tt.fields.catalogControls,
			}

			got, _ := svc.evaluateSubcontrol(tt.args.ctx, tt.args.auditScope, tt.args.control)
			assert.Optional(t, tt.want, got)
			assert.Optional(t, tt.wantSvc, svc)
		})
	}
}

func TestService_cacheControls(t *testing.T) {
	type fields struct {
		UnimplementedEvaluationServer evaluation.UnimplementedEvaluationServer
		orchestrator                  *api.RPCConnection[orchestrator.OrchestratorClient]
		scheduler                     *gocron.Scheduler
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
		wantSvc assert.Want[*Service]
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
				orchestrator: api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(connectionRefusedDialer)),
			},
			args: args{
				catalogId: testdata.MockCatalogID1,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "connection refused")
			},
		},
		{
			name: "no controls available",
			fields: fields{
				orchestrator:    api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {})))),
				catalogControls: make(map[string]map[string]*orchestrator.Control),
			},
			args: args{
				catalogId: testdata.MockCatalogID1,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, fmt.Sprintf("no controls for catalog '%s' available", testdata.MockCatalogID1))
			},
		},
		{
			name: "Happy path",
			fields: fields{
				orchestrator: api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(newBufConnDialer(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
				})))),
				catalogControls: make(map[string]map[string]*orchestrator.Control),
			},
			args: args{
				catalogId: testdata.MockCatalogID1,
			},
			wantSvc: func(t *testing.T, got *Service) bool {
				assert.Equal(t, 1, len(got.catalogControls))
				return assert.Equal(t, 4, len(got.catalogControls[testdata.MockCatalogID1]))
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				UnimplementedEvaluationServer: tt.fields.UnimplementedEvaluationServer,
				orchestrator:                  tt.fields.orchestrator,
				scheduler:                     tt.fields.scheduler,
				authz:                         tt.fields.authz,
				storage:                       tt.fields.storage,
				catalogControls:               tt.fields.catalogControls,
			}
			err := svc.cacheControls(tt.args.catalogId)
			tt.wantErr(t, err)
			assert.Optional(t, tt.wantSvc, svc)
		})
	}
}

func TestService_CreateEvaluationResult(t *testing.T) {
	type fields struct {
		orchestrator    *api.RPCConnection[orchestrator.OrchestratorClient]
		scheduler       *gocron.Scheduler
		authz           service.AuthorizationStrategy
		storage         persistence.Storage
		catalogControls map[string]map[string]*orchestrator.Control
	}
	type args struct {
		ctx context.Context
		req *evaluation.CreateEvaluationResultRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes assert.Want[*evaluation.EvaluationResult]
		wantErr assert.WantErr
	}{
		{
			name: "Request validation error",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				req: &evaluation.CreateEvaluationResultRequest{
					Result: &evaluation.EvaluationResult{
						ControlCategoryName: orchestratortest.MockControl1.CategoryName,
						ControlCatalogId:    orchestratortest.MockControl1.CategoryCatalogId,
						Status:              evaluation.EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT_MANUALLY,
					},
				},
			},
			wantRes: assert.Nil[*evaluation.EvaluationResult],
			wantErr: func(t *testing.T, err error) bool {
				return assert.ErrorContains(t, err, " validation error:\n - result.control_id")
			},
		},
		{
			name: "Wrong status",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				req: &evaluation.CreateEvaluationResultRequest{
					Result: &evaluation.EvaluationResult{
						ControlId:           orchestratortest.MockControl1.Id,
						ControlCategoryName: orchestratortest.MockControl1.CategoryName,
						ControlCatalogId:    orchestratortest.MockControl1.CategoryCatalogId,
						Status:              evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT,
					},
				},
			},
			wantRes: assert.Nil[*evaluation.EvaluationResult],
			wantErr: func(t *testing.T, err error) bool {
				return assert.ErrorContains(t, err, "only manually set statuses are allowed")
			},
		},
		{
			name: "Missing validity",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				req: &evaluation.CreateEvaluationResultRequest{
					Result: &evaluation.EvaluationResult{
						ControlId:           orchestratortest.MockControl1.Id,
						ControlCategoryName: orchestratortest.MockControl1.CategoryName,
						ControlCatalogId:    orchestratortest.MockControl1.CategoryCatalogId,
						Status:              evaluation.EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT_MANUALLY,
					},
				},
			},
			wantRes: assert.Nil[*evaluation.EvaluationResult],
			wantErr: func(t *testing.T, err error) bool {
				return assert.ErrorContains(t, err, "validity must be set")
			},
		},
		{
			name: "Permission denied",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(false),
			},
			args: args{
				req: &evaluation.CreateEvaluationResultRequest{
					Result: &evaluation.EvaluationResult{
						ControlId:           orchestratortest.MockControl1.Id,
						ControlCategoryName: orchestratortest.MockControl1.CategoryName,
						ControlCatalogId:    orchestratortest.MockControl1.CategoryCatalogId,
						Status:              evaluation.EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT_MANUALLY,
						ValidUntil:          timestamppb.New(time.Now().Add(24 * time.Hour)),
					},
				},
			},
			wantRes: assert.Nil[*evaluation.EvaluationResult],
			wantErr: func(t *testing.T, err error) bool {
				assert.Equal(t, codes.PermissionDenied, status.Code(err))
				return assert.ErrorContains(t, err, service.ErrPermissionDenied.Error())
			},
		},
		{
			name: "Database error",
			fields: fields{
				storage: &testutil.StorageWithError{CreateErr: gorm.ErrInvalidDB},
				authz:   &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				req: &evaluation.CreateEvaluationResultRequest{
					Result: &evaluation.EvaluationResult{
						ControlId:           orchestratortest.MockControl1.Id,
						ControlCategoryName: orchestratortest.MockControl1.CategoryName,
						ControlCatalogId:    orchestratortest.MockControl1.CategoryCatalogId,
						Status:              evaluation.EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT_MANUALLY,
						ValidUntil:          timestamppb.New(time.Now().Add(24 * time.Hour)),
					},
				},
			},
			wantRes: assert.Nil[*evaluation.EvaluationResult],
			wantErr: func(t *testing.T, err error) bool {
				assert.Equal(t, codes.Internal, status.Code(err))
				return assert.ErrorContains(t, err, persistence.ErrDatabase.Error())
			},
		},
		{
			name: "Happy path",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				req: &evaluation.CreateEvaluationResultRequest{
					Result: &evaluation.EvaluationResult{
						ControlId:           orchestratortest.MockControl1.Id,
						ControlCategoryName: orchestratortest.MockControl1.CategoryName,
						ControlCatalogId:    orchestratortest.MockControl1.CategoryCatalogId,
						Status:              evaluation.EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT_MANUALLY,
						ValidUntil:          timestamppb.New(time.Now().Add(24 * time.Hour)),
					},
				},
			},
			wantRes: func(t *testing.T, got *evaluation.EvaluationResult) bool {
				return assert.Equal(t, orchestratortest.MockControl1.Id, got.ControlId)
			},
			wantErr: assert.Nil[error],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				orchestrator:    tt.fields.orchestrator,
				scheduler:       tt.fields.scheduler,
				authz:           tt.fields.authz,
				storage:         tt.fields.storage,
				catalogControls: tt.fields.catalogControls,
			}
			gotRes, err := svc.CreateEvaluationResult(tt.args.ctx, tt.args.req)

			tt.wantErr(t, err)
			tt.wantRes(t, gotRes)
		})
	}
}

func TestDefaultServiceSpec(t *testing.T) {
	tests := []struct {
		name string
		want assert.Want[launcher.ServiceSpec]
	}{
		{
			name: "Happy path",
			want: func(t *testing.T, got launcher.ServiceSpec) bool {
				return assert.NotNil(t, got)

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultServiceSpec()

			tt.want(t, got)
		})
	}
}

func Test_handlePending(t *testing.T) {
	type args struct {
		eval *evaluation.EvaluationResult
	}
	tests := []struct {
		name  string
		args  args
		want  assert.Want[evaluation.EvaluationStatus]
		want1 assert.Want[[]string]
	}{
		{
			name: "Status: Pending",
			args: args{
				eval: &evaluation.EvaluationResult{
					Status: evaluation.EvaluationStatus_EVALUATION_STATUS_PENDING,
				},
			},
			want: func(t *testing.T, got evaluation.EvaluationStatus) bool {
				return assert.Equal(t, evaluation.EvaluationStatus_EVALUATION_STATUS_PENDING, got)
			},
			want1: func(t *testing.T, got []string) bool {
				return assert.Nil(t, got)
			},
		},
		{
			name: "Status: Compliant",
			args: args{
				eval: &evaluation.EvaluationResult{
					Status: evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT,
				},
			},
			want: func(t *testing.T, got evaluation.EvaluationStatus) bool {
				return assert.Equal(t, evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT, got)
			},
			want1: func(t *testing.T, got []string) bool {
				return assert.Nil(t, got)
			},
		},
		{
			name: "Status: Compliant manually",
			args: args{
				eval: &evaluation.EvaluationResult{
					Status: evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT_MANUALLY,
				},
			},
			want: func(t *testing.T, got evaluation.EvaluationStatus) bool {
				return assert.Equal(t, evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT, got)
			},
			want1: func(t *testing.T, got []string) bool {
				return assert.Nil(t, got)
			},
		},
		{
			name: "Status: Not compliant manually without failing assessment results",
			args: args{
				eval: &evaluation.EvaluationResult{
					Status: evaluation.EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT_MANUALLY,
				},
			},
			want: func(t *testing.T, got evaluation.EvaluationStatus) bool {
				return assert.Equal(t, evaluation.EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT, got)
			},
			want1: func(t *testing.T, got []string) bool {
				return assert.Nil(t, got)
			},
		},
		{
			name: "Status: Not compliant with failing assessment results",
			args: args{
				eval: &evaluation.EvaluationResult{
					Status:                     evaluation.EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT,
					FailingAssessmentResultIds: []string{"fail1", "fail2"},
				},
			},
			want: func(t *testing.T, got evaluation.EvaluationStatus) bool {
				return assert.Equal(t, evaluation.EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT, got)
			},
			want1: func(t *testing.T, got []string) bool {
				return assert.Equal(t, []string{"fail1", "fail2"}, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := handlePending(tt.args.eval)
			tt.want(t, got)
			tt.want1(t, got1)
		})
	}
}

func Test_handleCompliant(t *testing.T) {
	type args struct {
		er *evaluation.EvaluationResult
	}
	tests := []struct {
		name  string
		args  args
		want  assert.Want[evaluation.EvaluationStatus]
		want1 assert.Want[[]string]
	}{
		{
			name: "Status: Pending",
			args: args{
				er: &evaluation.EvaluationResult{
					Status: evaluation.EvaluationStatus_EVALUATION_STATUS_PENDING,
				},
			},
			want: func(t *testing.T, got evaluation.EvaluationStatus) bool {
				return assert.Equal(t, evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT, got)
			},
			want1: func(t *testing.T, got []string) bool {
				return assert.Nil(t, got)
			},
		},
		{
			name: "Status: Compliant",
			args: args{
				er: &evaluation.EvaluationResult{
					Status: evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT,
				},
			},
			want: func(t *testing.T, got evaluation.EvaluationStatus) bool {
				return assert.Equal(t, evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT, got)
			},
			want1: func(t *testing.T, got []string) bool {
				return assert.Nil(t, got)
			},
		},
		{
			name: "Status: Compliant manually",
			args: args{
				er: &evaluation.EvaluationResult{
					Status: evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT_MANUALLY,
				},
			},
			want: func(t *testing.T, got evaluation.EvaluationStatus) bool {
				return assert.Equal(t, evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT, got)
			},
			want1: func(t *testing.T, got []string) bool {
				return assert.Nil(t, got)
			},
		},
		{
			name: "Status: Not compliant manually",
			args: args{
				er: &evaluation.EvaluationResult{
					Status:                     evaluation.EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT_MANUALLY,
					FailingAssessmentResultIds: []string{"fail1", "fail2"},
				},
			},
			want: func(t *testing.T, got evaluation.EvaluationStatus) bool {
				return assert.Equal(t, evaluation.EvaluationStatus_EVALUATION_STATUS_NOT_COMPLIANT, got)
			},
			want1: func(t *testing.T, got []string) bool {
				return assert.Equal(t, []string{"fail1", "fail2"}, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := handleCompliant(tt.args.er)
			tt.want(t, got)
			tt.want1(t, got1)
		})
	}
}
