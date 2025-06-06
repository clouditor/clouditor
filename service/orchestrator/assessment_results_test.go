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

package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"sync"
	"testing"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/orchestratortest"
	"clouditor.io/clouditor/v2/internal/util"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/service"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestService_GetAssessmentResult(t *testing.T) {
	type fields struct {
		storage persistence.Storage
		authz   service.AuthorizationStrategy
	}
	type args struct {
		ctx context.Context
		req *orchestrator.GetAssessmentResultRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		res     *assessment.AssessmentResult
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "request is missing",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
			},
			args: args{},
			res:  nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Contains(t, err.Error(), "empty request")
			},
		},
		{
			name: "permission denied because of non authorized target_of_evaluation_id",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult1))
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult2))
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult3))
				}),
				authz: servicetest.NewAuthorizationStrategy(false, orchestratortest.MockAssessmentResult1.TargetOfEvaluationId),
			},
			args: args{
				req: orchestratortest.MockAssessmentResultRequest2,
				ctx: context.TODO(),
			},
			res: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, service.ErrPermissionDenied.Error())
			},
		},
		{
			name: "record not found in database",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult2))
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult3))
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				ctx: context.TODO(),
				req: orchestratortest.MockAssessmentResultRequest1,
			},
			res: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.NotFound, status.Code(err))
				return assert.ErrorContains(t, err, "assessment result not found")
			},
		},
		{
			name: "database error handling",
			fields: fields{
				storage: &testutil.StorageWithError{GetErr: ErrSomeError},
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{req: orchestratortest.MockAssessmentResultRequest1},
			res:  nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.Internal, status.Code(err))
				return assert.Contains(t, err.Error(), ErrSomeError.Error())
			},
		},
		{
			name: "Happy path with 'allow all key'",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult1))
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult2))
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult3))
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: orchestratortest.MockAssessmentResultRequest1,
			},
			res:     orchestratortest.MockAssessmentResult1,
			wantErr: assert.NoError,
		},
		{
			name: "Happy path with allowed target_of_evaluation_id",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult1))
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult2))
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult3))
				}),
				authz: servicetest.NewAuthorizationStrategy(false, orchestratortest.MockAssessmentResult1.TargetOfEvaluationId),
			},
			args: args{
				req: orchestratortest.MockAssessmentResultRequest1,
			},
			res:     orchestratortest.MockAssessmentResult1,
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				storage: tt.fields.storage,
				authz:   tt.fields.authz,
			}
			gotRes, err := svc.GetAssessmentResult(context.Background(), tt.args.req)
			tt.wantErr(t, err)

			if tt.res == nil {
				assert.Nil(t, gotRes)
			} else {
				assert.NoError(t, api.Validate(gotRes))
				assert.Equal(t, tt.res, gotRes)
			}
		})
	}
}

func TestService_ListAssessmentResults(t *testing.T) {
	type fields struct {
		storage persistence.Storage
		authz   service.AuthorizationStrategy
	}
	type args struct {
		ctx context.Context
		req *orchestrator.ListAssessmentResultsRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *orchestrator.ListAssessmentResultsResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "request is missing",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args:    args{},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "empty request")
			},
		},
		{
			name: "request is empty",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.ListAssessmentResultsRequest{},
			},
			wantRes: &orchestrator.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{},
			},
			wantErr: assert.NoError,
		},
		{
			name: "list all with allow all",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult1))
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult2))
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{req: &orchestrator.ListAssessmentResultsRequest{}},
			wantRes: &orchestrator.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					orchestratortest.MockAssessmentResult1,
					orchestratortest.MockAssessmentResult2,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "list all only with allowed target of evaluation",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult1))
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult2))
				}),
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID1),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.ListAssessmentResultsRequest{},
			},
			wantRes: &orchestrator.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					orchestratortest.MockAssessmentResult1,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "specify filtered target of evaluation ID which is not allowed",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult1))
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult2))
				}),
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID1),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.ListAssessmentResultsRequest{
					Filter: &orchestrator.ListAssessmentResultsRequest_Filter{
						TargetOfEvaluationId: util.Ref(testdata.MockTargetOfEvaluationID2),
					},
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, service.ErrPermissionDenied)
			},
		},
		{
			name: "return filtered target of evaluation ID",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult1))
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult2))
				}),
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID1),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.ListAssessmentResultsRequest{
					Filter: &orchestrator.ListAssessmentResultsRequest_Filter{
						TargetOfEvaluationId: util.Ref(testdata.MockTargetOfEvaluationID1),
					},
				},
			},
			wantRes: &orchestrator.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					orchestratortest.MockAssessmentResult1,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "return filtered target of evaluation ID and filtered compliant assessment results",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
				}),
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID1),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.ListAssessmentResultsRequest{
					Filter: &orchestrator.ListAssessmentResultsRequest_Filter{
						TargetOfEvaluationId: util.Ref(testdata.MockTargetOfEvaluationID1),
						Compliant:            util.Ref(true),
					},
				},
			},
			wantRes: &orchestrator.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					orchestratortest.MockAssessmentResult1,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "return filtered target of evaluation ID and filtered non-compliant assessment results",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
				}),
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID1),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.ListAssessmentResultsRequest{
					Filter: &orchestrator.ListAssessmentResultsRequest_Filter{
						TargetOfEvaluationId: util.Ref(testdata.MockTargetOfEvaluationID1),
						Compliant:            util.Ref(false),
					},
				},
			},
			wantRes: &orchestrator.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					orchestratortest.MockAssessmentResult3,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "return filtered compliant assessment results",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.ListAssessmentResultsRequest{
					Filter: &orchestrator.ListAssessmentResultsRequest_Filter{
						Compliant: util.Ref(true),
					},
				},
			},
			wantRes: &orchestrator.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					orchestratortest.MockAssessmentResult1,
					orchestratortest.MockAssessmentResult2,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "return filtered assessment results with specific tool_id",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.ListAssessmentResultsRequest{
					Filter: &orchestrator.ListAssessmentResultsRequest_Filter{
						ToolId: util.Ref(testdata.MockAssessmentResultToolID),
					},
				},
			},
			wantRes: &orchestrator.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					orchestratortest.MockAssessmentResult4,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "return filtered non-compliant assessment results",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))

				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.ListAssessmentResultsRequest{
					Filter: &orchestrator.ListAssessmentResultsRequest_Filter{
						Compliant: util.Ref(false),
					},
				},
			},
			wantRes: &orchestrator.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					orchestratortest.MockAssessmentResult3,
					orchestratortest.MockAssessmentResult4,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "return filtered target of evaluation ID and one filtered metric ID",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
				}),
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID1),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.ListAssessmentResultsRequest{
					Filter: &orchestrator.ListAssessmentResultsRequest_Filter{
						TargetOfEvaluationId: util.Ref(testdata.MockTargetOfEvaluationID1),
						MetricIds:            []string{testdata.MockMetricID1, testdata.MockMetricID2},
					},
				},
			},
			wantRes: &orchestrator.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					orchestratortest.MockAssessmentResult1,
					orchestratortest.MockAssessmentResult3,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "return filtered target of evaluation ID and two filtered metric IDs",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
				}),
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID1),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.ListAssessmentResultsRequest{
					Filter: &orchestrator.ListAssessmentResultsRequest_Filter{
						TargetOfEvaluationId: util.Ref(testdata.MockTargetOfEvaluationID1),
						MetricIds:            []string{testdata.MockMetricID1, testdata.MockMetricID2},
					},
				},
			},
			wantRes: &orchestrator.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					orchestratortest.MockAssessmentResult1,
					orchestratortest.MockAssessmentResult3,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "return one filtered metric ID",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.ListAssessmentResultsRequest{
					Filter: &orchestrator.ListAssessmentResultsRequest_Filter{
						MetricIds: []string{testdata.MockMetricID1},
					},
				},
			},
			wantRes: &orchestrator.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					orchestratortest.MockAssessmentResult1,
					orchestratortest.MockAssessmentResult2,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Invalid target of evaluation id request",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.ListAssessmentResultsRequest{
					Filter: &orchestrator.ListAssessmentResultsRequest_Filter{
						TargetOfEvaluationId: util.Ref("No Valid UUID"),
					},
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "value must be a valid UUID")
			},
		},
		{
			name: "grouped by resource ID",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.ListAssessmentResultsRequest{
					LatestByResourceId: util.Ref(true),
				},
			},
			wantRes: &orchestrator.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					orchestratortest.MockAssessmentResult4,
					orchestratortest.MockAssessmentResult1,
					orchestratortest.MockAssessmentResult3,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "grouped by resource ID with filter",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.ListAssessmentResultsRequest{
					LatestByResourceId: util.Ref(true),
					Filter: &orchestrator.ListAssessmentResultsRequest_Filter{
						TargetOfEvaluationId: util.Ref(testdata.MockTargetOfEvaluationID1),
					},
				},
			},
			wantRes: &orchestrator.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					orchestratortest.MockAssessmentResult1,
					orchestratortest.MockAssessmentResult3,
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				storage: tt.fields.storage,
				authz:   tt.fields.authz,
			}
			gotRes, err := svc.ListAssessmentResults(tt.args.ctx, tt.args.req)
			tt.wantErr(t, err)

			if tt.wantRes == nil {
				assert.Nil(t, gotRes)
			} else {
				assert.NoError(t, api.Validate(gotRes))
				assert.Equal(t, tt.wantRes, gotRes)
			}
		})
	}
}

func TestAssessmentResultHook(t *testing.T) {
	var (
		hookCallCounter = 0
		wg              sync.WaitGroup
	)
	wg.Add(2)

	firstHookFunction := func(ctx context.Context, assessmentResult *assessment.AssessmentResult, err error) {
		hookCallCounter++
		log.Println("Hello from inside the firstHookFunction")

		wg.Done()
	}

	secondHookFunction := func(ctx context.Context, assessmentResult *assessment.AssessmentResult, err error) {
		hookCallCounter++
		log.Println("Hello from inside the secondHookFunction")

		wg.Done()
	}

	service := NewService()
	service.RegisterAssessmentResultHook(firstHookFunction)
	service.RegisterAssessmentResultHook(secondHookFunction)

	// Check if first hook is registered
	funcName1 := runtime.FuncForPC(reflect.ValueOf(service.AssessmentResultHooks[0]).Pointer()).Name()
	funcName2 := runtime.FuncForPC(reflect.ValueOf(firstHookFunction).Pointer()).Name()
	assert.Equal(t, funcName1, funcName2)

	// Check if second hook is registered
	funcName1 = runtime.FuncForPC(reflect.ValueOf(service.AssessmentResultHooks[1]).Pointer()).Name()
	funcName2 = runtime.FuncForPC(reflect.ValueOf(secondHookFunction).Pointer()).Name()
	assert.Equal(t, funcName1, funcName2)

	// Check GRPC call
	type args struct {
		in0        context.Context
		assessment *orchestrator.StoreAssessmentResultRequest
	}
	tests := []struct {
		name    string
		args    args
		wantRes *orchestrator.StoreAssessmentResultResponse
		wantErr assert.WantErr
	}{
		{
			name: "Store first assessment result to the map",
			args: args{
				in0: context.TODO(),
				assessment: &orchestrator.StoreAssessmentResultRequest{
					Result: orchestratortest.MockAssessmentResult1,
				},
			},
			wantErr: assert.Nil[error],
			wantRes: &orchestrator.StoreAssessmentResultResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hookCallCounter = 0
			s := service
			gotRes, err := s.StoreAssessmentResult(tt.args.in0, tt.args.assessment)

			// wait for all hooks (2 hooks)
			wg.Wait()

			tt.wantErr(t, err)
			assert.Equal(t, tt.wantRes, gotRes)

			var results []*assessment.AssessmentResult
			assert.NoError(t, s.storage.List(&results, "", true, 0, -1))

			assert.NotEmpty(t, results)
			assert.Equal(t, 2, hookCallCounter)
		})
	}
}

func TestStoreAssessmentResult(t *testing.T) {
	type args struct {
		in0        context.Context
		assessment *orchestrator.StoreAssessmentResultRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *orchestrator.StoreAssessmentResultResponse
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "Request validation error",
			args: args{
				assessment: &orchestrator.StoreAssessmentResultRequest{
					Result: &assessment.AssessmentResult{},
				},
			},
			wantResp: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, " validation error:\n - result.id: value is empty, which is not a valid UUID")
			},
		},
		{
			name: "Store assessment without metricId to the map",
			args: args{
				in0: context.TODO(),
				assessment: &orchestrator.StoreAssessmentResultRequest{
					Result: &assessment.AssessmentResult{
						Id:                   uuid.NewString(),
						EvidenceId:           testdata.MockEvidenceID1,
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
						Timestamp:            timestamppb.Now(),
						MetricConfiguration: &assessment.MetricConfiguration{
							TargetValue:          toStruct(1.0),
							Operator:             "<=",
							IsDefault:            true,
							TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
							MetricId:             testdata.MockMetricID1,
						},
						ComplianceComment: testdata.MockAssessmentResultNonComplianceComment,
						Compliant:         true,
						ResourceId:        testdata.MockVirtualMachineID1,
						ResourceTypes:     testdata.MockVirtualMachineTypes,
						ToolId:            util.Ref(assessment.AssessmentToolId),
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "result.metric_id: value length must be at least 1 characters")
			},
			wantResp: nil,
		},
		{
			name: "Happy path",
			args: args{
				in0: context.TODO(),
				assessment: &orchestrator.StoreAssessmentResultRequest{
					Result: orchestratortest.MockAssessmentResult1,
				},
			},
			wantErr:  assert.NoError,
			wantResp: &orchestrator.StoreAssessmentResultResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			gotResp, err := s.StoreAssessmentResult(tt.args.in0, tt.args.assessment)
			tt.wantErr(t, err)
			assert.Equal(t, tt.wantResp, gotResp)

			var results []*assessment.AssessmentResult
			assert.NoError(t, s.storage.List(&results, "", true, 0, -1))

			if err == nil {
				assert.NotNil(t, results)
			} else {
				assert.Empty(t, results)
			}
		})
	}
}

func TestStoreAssessmentResults(t *testing.T) {
	const (
		count1 = 1
		count2 = 2
	)

	type fields struct {
		countElementsInMock    int
		countElementsInResults int
	}

	type args struct {
		streamToServer            *mockStreamer
		streamToClientWithSendErr *mockStreamerWithSendErr
		streamToServerWithRecvErr *mockStreamerWithRecvErr
	}

	tests := []struct {
		name            string
		fields          fields
		args            args
		wantErr         bool
		wantRespMessage []orchestrator.StoreAssessmentResultsResponse
		wantErrMessage  string
	}{
		{
			name: "Store 2 assessment results to the map",
			fields: fields{
				countElementsInMock:    count2,
				countElementsInResults: count2,
			},
			args:    args{streamToServer: createMockStream(createStoreAssessmentResultRequestsMock(count2))},
			wantErr: false,
			wantRespMessage: []orchestrator.StoreAssessmentResultsResponse{
				{
					Status: true,
				},
				{
					Status: true,
				},
			},
		},
		{
			name: "Missing MetricID",
			fields: fields{
				countElementsInMock:    count1,
				countElementsInResults: 0,
			},
			args:    args{streamToServer: createMockStream(createStoreAssessmentResultRequestMockWithMissingMetricID(count1))},
			wantErr: false,
			wantRespMessage: []orchestrator.StoreAssessmentResultsResponse{
				{
					Status:        false,
					StatusMessage: "result.metric_id: value length must be at least 1 characters",
				},
			},
		},
		{
			name: "Error in stream to server - Recv()-err",
			fields: fields{
				countElementsInMock:    count1,
				countElementsInResults: 0,
			},
			args:           args{streamToServerWithRecvErr: createMockStreamWithRecvErr(createStoreAssessmentResultRequestsMock(count1))},
			wantErr:        true,
			wantErrMessage: "rpc error: code = Unknown desc = cannot receive stream request",
		},
		{
			name: "Error in stream to client - Send()-err",
			fields: fields{
				countElementsInMock:    count1,
				countElementsInResults: 0,
			},
			args:           args{streamToClientWithSendErr: createMockStreamWithSendErr(createStoreAssessmentResultRequestsMock(count1))},
			wantErr:        true,
			wantErrMessage: "rpc error: code = Unknown desc = cannot stream response to the client",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()

			var err error

			if tt.args.streamToServer != nil {
				err = s.StoreAssessmentResults(tt.args.streamToServer)
			} else if tt.args.streamToClientWithSendErr != nil {
				err = s.StoreAssessmentResults(tt.args.streamToClientWithSendErr)
			} else if tt.args.streamToServerWithRecvErr != nil {
				err = s.StoreAssessmentResults(tt.args.streamToServerWithRecvErr)
			}

			var results []*assessment.AssessmentResult
			assert.NoError(t, s.storage.List(&results, "", true, 0, -1))

			if (err != nil) != tt.wantErr {
				t.Errorf("Got StoreAssessmentResults() error = %v, wantErrMessage %v", err, tt.wantErr)
				assert.Equal(t, tt.fields.countElementsInResults, len(results))
				return
			} else if tt.wantErr {
				assert.Contains(t, err.Error(), tt.wantErrMessage)
			} else {
				// Close stream for testing
				close(tt.args.streamToServer.SentFromServer)
				assert.Nil(t, err)
				assert.Equal(t, tt.fields.countElementsInResults, len(results))

				// Check all stream responses from server to client
				i := 0
				for elem := range tt.args.streamToServer.SentFromServer {
					assert.Contains(t, elem.StatusMessage, tt.wantRespMessage[i].StatusMessage)
					assert.Equal(t, tt.wantRespMessage[i].Status, elem.Status)
					i++
				}
			}

		})
	}
}

// createStoreAssessmentResultRequestMockWithMissingMetricID create one StoreAssessmentResultRequest without ToolID
func createStoreAssessmentResultRequestMockWithMissingMetricID(count int) []*orchestrator.StoreAssessmentResultRequest {
	req := createStoreAssessmentResultRequestsMock(count)

	req[0].Result.MetricId = ""

	return req
}

// createStoreAssessmentResultRequestMocks creates store assessment result requests with random assessment result IDs
func createStoreAssessmentResultRequestsMock(count int) []*orchestrator.StoreAssessmentResultRequest {
	var mockRequests []*orchestrator.StoreAssessmentResultRequest

	for i := 0; i < count; i++ {
		storeAssessmentResultRequest := &orchestrator.StoreAssessmentResultRequest{
			Result: &assessment.AssessmentResult{
				Id:                   uuid.NewString(),
				MetricId:             fmt.Sprintf("assessmentResultMetricID-%d", i),
				EvidenceId:           testdata.MockEvidenceID1,
				TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
				Timestamp:            timestamppb.Now(),
				MetricConfiguration: &assessment.MetricConfiguration{
					TargetValue:          toStruct(1.0),
					Operator:             "<=",
					IsDefault:            true,
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					MetricId:             testdata.MockMetricID1,
				},
				ComplianceComment: testdata.MockAssessmentResultNonComplianceComment,
				Compliant:         true,
				ResourceId:        testdata.MockVirtualMachineID1,
				ResourceTypes:     testdata.MockVirtualMachineTypes,
				ToolId:            util.Ref(assessment.AssessmentToolId),
			},
		}

		mockRequests = append(mockRequests, storeAssessmentResultRequest)
	}

	return mockRequests
}

type mockStreamer struct {
	grpc.ServerStream
	RecvToServer   chan *orchestrator.StoreAssessmentResultRequest
	SentFromServer chan *orchestrator.StoreAssessmentResultsResponse
}

func createMockStream(requests []*orchestrator.StoreAssessmentResultRequest) *mockStreamer {
	m := &mockStreamer{
		RecvToServer: make(chan *orchestrator.StoreAssessmentResultRequest, len(requests)),
	}
	for _, req := range requests {
		m.RecvToServer <- req
	}

	m.SentFromServer = make(chan *orchestrator.StoreAssessmentResultsResponse, len(requests))
	return m
}

func (m mockStreamer) Send(response *orchestrator.StoreAssessmentResultsResponse) error {
	m.SentFromServer <- response
	return nil
}

func (m mockStreamer) Recv() (*orchestrator.StoreAssessmentResultRequest, error) {
	if len(m.RecvToServer) == 0 {
		return nil, io.EOF
	}
	req, more := <-m.RecvToServer
	if !more {
		return nil, errors.New("empty")
	}

	return req, nil
}

func (mockStreamer) SendAndClose(_ *emptypb.Empty) error {
	return nil
}

func (mockStreamer) SetHeader(_ metadata.MD) error {
	panic("implement me")
}

func (mockStreamer) SendHeader(_ metadata.MD) error {
	panic("implement me")
}

func (mockStreamer) SetTrailer(_ metadata.MD) {
	panic("implement me")
}

func (mockStreamer) Context() context.Context {
	return context.TODO()
}

func (mockStreamer) SendMsg(_ interface{}) error {
	panic("implement me")
}

func (mockStreamer) RecvMsg(_ interface{}) error {
	panic("implement me")
}

type mockStreamerWithSendErr struct {
	grpc.ServerStream
	RecvToServer   chan *orchestrator.StoreAssessmentResultRequest
	SentFromServer chan *orchestrator.StoreAssessmentResultsResponse
}

func (mockStreamerWithSendErr) Context() context.Context {
	return context.TODO()
}

func (mockStreamerWithSendErr) Send(*orchestrator.StoreAssessmentResultsResponse) error {
	return errors.New("Send()-err")
}

func (m mockStreamerWithSendErr) Recv() (*orchestrator.StoreAssessmentResultRequest, error) {
	if len(m.RecvToServer) == 0 {
		return nil, io.EOF
	}
	req, more := <-m.RecvToServer
	if !more {
		return nil, errors.New("empty")
	}

	return req, nil
}

func createMockStreamWithSendErr(requests []*orchestrator.StoreAssessmentResultRequest) *mockStreamerWithSendErr {
	m := &mockStreamerWithSendErr{
		RecvToServer: make(chan *orchestrator.StoreAssessmentResultRequest, len(requests)),
	}
	for _, req := range requests {
		m.RecvToServer <- req
	}

	m.SentFromServer = make(chan *orchestrator.StoreAssessmentResultsResponse, len(requests))
	return m
}

type mockStreamerWithRecvErr struct {
	grpc.ServerStream
	RecvToServer   chan *orchestrator.StoreAssessmentResultRequest
	SentFromServer chan *orchestrator.StoreAssessmentResultsResponse
}

func (mockStreamerWithRecvErr) Context() context.Context {
	return context.TODO()
}

func (mockStreamerWithRecvErr) Send(*orchestrator.StoreAssessmentResultsResponse) error {
	panic("implement me")
}

func (mockStreamerWithRecvErr) Recv() (*orchestrator.StoreAssessmentResultRequest, error) {
	err := errors.New("Recv()-error")

	return nil, err
}

func createMockStreamWithRecvErr(requests []*orchestrator.StoreAssessmentResultRequest) *mockStreamerWithRecvErr {
	m := &mockStreamerWithRecvErr{
		RecvToServer: make(chan *orchestrator.StoreAssessmentResultRequest, len(requests)),
	}
	for _, req := range requests {
		m.RecvToServer <- req
	}

	m.SentFromServer = make(chan *orchestrator.StoreAssessmentResultsResponse, len(requests))
	return m
}
