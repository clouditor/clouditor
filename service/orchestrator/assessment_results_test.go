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

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/internal/testutil/orchestratortest"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
				authz:   &service.AuthorizationStrategyAllowAll{},
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
				authz:   &service.AuthorizationStrategyAllowAll{},
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
				authz: &service.AuthorizationStrategyAllowAll{},
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
			name: "list all only with allowed cloud service",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult1))
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult2))
				}),
				authz: &service.AuthorizationStrategyJWT{Key: testutil.TestCustomClaims},
			},
			args: args{
				ctx: testutil.TestContextOnlyService1,
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
			name: "specify filtered cloud service ID which is not allowed",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult1))
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult2))
				}),
				authz: &service.AuthorizationStrategyJWT{Key: testutil.TestCustomClaims},
			},
			args: args{
				ctx: testutil.TestContextOnlyService1,
				req: &orchestrator.ListAssessmentResultsRequest{
					FilteredCloudServiceId: util.Ref(testdata.MockAnotherCloudServiceID),
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, err, service.ErrPermissionDenied)
			},
		},
		{
			name: "return filtered cloud service ID",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult1))
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult2))
				}),
				authz: &service.AuthorizationStrategyJWT{Key: testutil.TestCustomClaims},
			},
			args: args{
				ctx: testutil.TestContextOnlyService1,
				req: &orchestrator.ListAssessmentResultsRequest{
					FilteredCloudServiceId: util.Ref(testdata.MockCloudServiceID),
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
			name: "return filtered cloud service ID and filtered compliant assessment results",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
				}),
				authz: &service.AuthorizationStrategyJWT{Key: testutil.TestCustomClaims},
			},
			args: args{
				ctx: testutil.TestContextOnlyService1,
				req: &orchestrator.ListAssessmentResultsRequest{
					FilteredCloudServiceId: util.Ref(testdata.MockCloudServiceID),
					FilteredCompliant:      util.Ref(true),
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
			name: "return filtered cloud service ID and filtered non-compliant assessment results",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
				}),
				authz: &service.AuthorizationStrategyJWT{Key: testutil.TestCustomClaims},
			},
			args: args{
				ctx: testutil.TestContextOnlyService1,
				req: &orchestrator.ListAssessmentResultsRequest{
					FilteredCloudServiceId: util.Ref(testdata.MockCloudServiceID),
					FilteredCompliant:      util.Ref(false),
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
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				ctx: testutil.TestContextOnlyService1,
				req: &orchestrator.ListAssessmentResultsRequest{
					FilteredCompliant: util.Ref(true),
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
			name: "return filtered non-compliant assessment results",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))

				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				ctx: testutil.TestContextOnlyService1,
				req: &orchestrator.ListAssessmentResultsRequest{
					FilteredCompliant: util.Ref(false),
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
			name: "return filtered cloud service ID and one filtered metric ID",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
				}),
				authz: &service.AuthorizationStrategyJWT{Key: testutil.TestCustomClaims},
			},
			args: args{
				ctx: testutil.TestContextOnlyService1,
				req: &orchestrator.ListAssessmentResultsRequest{
					FilteredCloudServiceId: util.Ref(testdata.MockCloudServiceID),
					FilteredMetricId:       []string{testdata.MockMetricID},
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
			name: "return filtered cloud service ID and two filtered metric IDs",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
				}),
				authz: &service.AuthorizationStrategyJWT{Key: testutil.TestCustomClaims},
			},
			args: args{
				ctx: testutil.TestContextOnlyService1,
				req: &orchestrator.ListAssessmentResultsRequest{
					FilteredCloudServiceId: util.Ref(testdata.MockCloudServiceID),
					FilteredMetricId:       []string{testdata.MockMetricID, testdata.MockAnotherMetricID},
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
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				ctx: testutil.TestContextOnlyService1,
				req: &orchestrator.ListAssessmentResultsRequest{
					FilteredMetricId: []string{testdata.MockMetricID},
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
			name: "Invalid cloud service id request",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				ctx: testutil.TestContextOnlyService1,
				req: &orchestrator.ListAssessmentResultsRequest{
					FilteredCloudServiceId: util.Ref("testCloudServiceID"),
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "FilteredCloudServiceId: value must be a valid UUID")
			},
		},
		{
			name: "grouped by resource ID",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.MockAssessmentResults))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				req: &orchestrator.ListAssessmentResultsRequest{
					LatestByResourceId: util.Ref(true),
				},
			},
			wantRes: &orchestrator.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					orchestratortest.MockAssessmentResult4,
					orchestratortest.MockAssessmentResult3,
					orchestratortest.MockAssessmentResult1,
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
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				req: &orchestrator.ListAssessmentResultsRequest{
					LatestByResourceId:     util.Ref(true),
					FilteredCloudServiceId: util.Ref(testdata.MockCloudServiceID),
				},
			},
			wantRes: &orchestrator.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					orchestratortest.MockAssessmentResult3,
					orchestratortest.MockAssessmentResult1,
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
				assert.NoError(t, gotRes.Validate())
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

	firstHookFunction := func(assessmentResult *assessment.AssessmentResult, err error) {
		hookCallCounter++
		log.Println("Hello from inside the firstHookFunction")

		wg.Done()
	}

	secondHookFunction := func(assessmentResult *assessment.AssessmentResult, err error) {
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
		name     string
		args     args
		wantResp *orchestrator.StoreAssessmentResultResponse
		wantErr  bool
	}{
		{
			name: "Store first assessment result to the map",
			args: args{
				in0: context.TODO(),
				assessment: &orchestrator.StoreAssessmentResultRequest{
					Result: &assessment.AssessmentResult{
						Id:             testdata.MockAssessmentResultID,
						MetricId:       testdata.MockMetricID,
						EvidenceId:     testdata.MockEvidenceID,
						CloudServiceId: testdata.MockCloudServiceID,
						Timestamp:      timestamppb.Now(),
						MetricConfiguration: &assessment.MetricConfiguration{
							TargetValue:    toStruct(1.0),
							Operator:       ">=",
							IsDefault:      true,
							CloudServiceId: testdata.MockCloudServiceID,
							MetricId:       testdata.MockMetricID,
						},
						NonComplianceComments: testdata.MockAssessmentResultNonComplianceComment,
						Compliant:             true,
						ResourceId:            testdata.MockResourceID,
						ResourceTypes:         []string{"ResourceType"},
					},
				},
			},
			wantErr:  false,
			wantResp: &orchestrator.StoreAssessmentResultResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hookCallCounter = 0
			s := service
			gotResp, err := s.StoreAssessmentResult(tt.args.in0, tt.args.assessment)

			// wait for all hooks (2 hooks)
			wg.Wait()

			if (err != nil) != tt.wantErr {
				t.Errorf("StoreAssessmentResult() error = %v, wantErrMessage %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("StoreAssessmentResult() gotResp = %v, want %v", gotResp, tt.wantResp)
			}

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
			name: "Store assessment to the map",
			args: args{
				in0: context.TODO(),
				assessment: &orchestrator.StoreAssessmentResultRequest{
					Result: &assessment.AssessmentResult{
						Id:             testdata.MockAssessmentResultID,
						MetricId:       "assessmentResultMetricID",
						EvidenceId:     testdata.MockEvidenceID,
						CloudServiceId: testdata.MockCloudServiceID,
						Timestamp:      timestamppb.Now(),
						MetricConfiguration: &assessment.MetricConfiguration{
							TargetValue:    toStruct(1.0),
							Operator:       "<=",
							IsDefault:      true,
							CloudServiceId: testdata.MockCloudServiceID,
							MetricId:       testdata.MockMetricID,
						},
						NonComplianceComments: testdata.MockAssessmentResultNonComplianceComment,
						Compliant:             true,
						ResourceId:            testdata.MockResourceID,
						ResourceTypes:         []string{"ResourceType"},
					},
				},
			},
			wantErr:  assert.NoError,
			wantResp: &orchestrator.StoreAssessmentResultResponse{},
		},
		{
			name: "Store assessment without metricId to the map",
			args: args{
				in0: context.TODO(),
				assessment: &orchestrator.StoreAssessmentResultRequest{
					Result: &assessment.AssessmentResult{
						Id:             testdata.MockAssessmentResultID,
						EvidenceId:     testdata.MockEvidenceID,
						CloudServiceId: testdata.MockCloudServiceID,
						Timestamp:      timestamppb.Now(),
						MetricConfiguration: &assessment.MetricConfiguration{
							TargetValue:    toStruct(1.0),
							Operator:       "<=",
							IsDefault:      true,
							CloudServiceId: testdata.MockCloudServiceID,
							MetricId:       testdata.MockMetricID,
						},
						NonComplianceComments: testdata.MockAssessmentResultNonComplianceComment,
						Compliant:             true,
						ResourceId:            testdata.MockResourceID,
						ResourceTypes:         []string{"ResourceType"},
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "caused by: invalid AssessmentResult.MetricId: value length must be at least 1 runes")
			},
			wantResp: nil,
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
					StatusMessage: "MetricId: value length must be at least 1 runes",
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
				Id:             uuid.NewString(),
				MetricId:       fmt.Sprintf("assessmentResultMetricID-%d", i),
				EvidenceId:     testdata.MockEvidenceID,
				CloudServiceId: testdata.MockCloudServiceID,
				Timestamp:      timestamppb.Now(),
				MetricConfiguration: &assessment.MetricConfiguration{
					TargetValue:    toStruct(1.0),
					Operator:       "<=",
					IsDefault:      true,
					CloudServiceId: testdata.MockCloudServiceID,
					MetricId:       testdata.MockMetricID,
				},
				NonComplianceComments: testdata.MockAssessmentResultNonComplianceComment,
				Compliant:             true,
				ResourceId:            testdata.MockResourceID,
				ResourceTypes:         []string{"ResourceType"},
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
	panic("implement me")
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
