// Copyright 2016-2023 Fraunhofer AISEC
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

package evidence

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
	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/prototest"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/evidencetest"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/orchestratortest"
	"clouditor.io/clouditor/v2/internal/util"
	"clouditor.io/clouditor/v2/launcher"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/persistence/gorm"
	"clouditor.io/clouditor/v2/service"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TestNewService is a simply test for NewService
func TestNewService(t *testing.T) {
	db, err := gorm.NewStorage(gorm.WithInMemory())
	assert.NoError(t, err)
	type args struct {
		opts []service.Option[*Service]
	}
	tests := []struct {
		name string
		args args
		want assert.Want[*Service]
	}{
		{
			name: "EvidenceStoreServer created without options",
			want: func(t *testing.T, got *Service) bool {
				// Storage should be default (in-memory storage). Hard to check since its type is not exported
				assert.NotNil(t, got.storage)
				return true
			},
		},
		{
			name: "EvidenceStoreServer created with storage option",
			args: args{opts: []service.Option[*Service]{WithStorage(db)}},
			want: func(t *testing.T, got *Service) bool {
				// Storage should be gorm (in-memory storage). Hard to check since its type is not exported
				assert.NotNil(t, got.storage)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewService()
			tt.want(t, got)
		})
	}
}

// TestStoreEvidence tests StoreEvidence
func TestService_StoreEvidence(t *testing.T) {
	type args struct {
		in0 context.Context
		req *evidence.StoreEvidenceRequest
	}
	tests := []struct {
		name    string
		args    args
		wantRes assert.Want[*evidence.StoreEvidenceResponse]
		want    assert.Want[*Service]
		wantErr assert.WantErr
	}{
		{
			name: "Store req to the map",
			args: args{
				in0: context.TODO(),
				req: &evidence.StoreEvidenceRequest{
					Evidence: &evidence.Evidence{
						Id:                   testdata.MockEvidenceID1,
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
						ToolId:               testdata.MockEvidenceToolID1,
						Timestamp:            timestamppb.Now(),
						Resource: prototest.NewAny(t, &ontology.VirtualMachine{
							Id: "mock-id",
						}),
					}},
			},
			wantRes: func(t *testing.T, got *evidence.StoreEvidenceResponse) bool {
				return assert.Empty[*evidence.StoreEvidenceResponse](t, got)
			},
			want: func(t *testing.T, s *Service) bool {
				e := &evidence.Evidence{}
				err := s.storage.Get(e)
				assert.NoError(t, err)
				return assert.Equal(t, testdata.MockEvidenceID1, e.Id)
			},
			wantErr: assert.Nil[error],
		},
		{
			name: "Store an evidence without toolId to the map",
			args: args{
				in0: context.TODO(),
				req: &evidence.StoreEvidenceRequest{
					Evidence: &evidence.Evidence{
						Id:                   testdata.MockEvidenceID1,
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
						Timestamp:            timestamppb.Now(),
						Resource: prototest.NewAny(t, &ontology.VirtualMachine{
							Id: "mock-id-1",
						}),
					},
				},
			},
			wantRes: assert.Nil[*evidence.StoreEvidenceResponse],
			wantErr: func(t *testing.T, err error) bool {
				return assert.ErrorContains(t, err, "evidence.tool_id: value length must be at least 1 characters")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			gotRes, err := s.StoreEvidence(tt.args.in0, tt.args.req)

			tt.wantErr(t, err)
			tt.wantRes(t, gotRes)
		})
	}
}

// TestStoreEvidences tests StoreEvidences
func TestService_StoreEvidences(t *testing.T) {
	type fields struct {
		count int
	}

	type args struct {
		streamToServer            *mockStreamer
		streamToClientWithSendErr *mockStreamerWithSendErr
		streamToServerWithRecvErr *mockStreamerWithRecvErr
	}

	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		wantErrMessage string
		wantResMessage *evidence.StoreEvidencesResponse
	}{
		{
			name:   "Store 1 evidence to the map",
			fields: fields{count: 1},
			args: args{
				streamToServer: createMockStream(createStoreEvidenceRequestMocks(t, 1))},
			wantErr: false,
			wantResMessage: &evidence.StoreEvidencesResponse{
				Status: true,
			},
		},
		{
			name:   "Store 2 evidences to the map",
			fields: fields{count: 2},
			args: args{
				streamToServer: createMockStream(createStoreEvidenceRequestMocks(t, 2))},
			wantErr: false,
			wantResMessage: &evidence.StoreEvidencesResponse{
				Status: true,
			},
		},
		{
			name:   "Store invalid evidence to the map",
			fields: fields{count: 1},
			args: args{
				streamToServer: createMockStream([]*evidence.StoreEvidenceRequest{
					{
						Evidence: &evidence.Evidence{
							Id:                   uuid.NewString(),
							TargetOfEvaluationId: "MockTargetOfEvaluationId",
							Timestamp:            timestamppb.Now(),
							Resource: prototest.NewAny(t, &ontology.VirtualMachine{
								Id: "mock-id-1",
							}),
						},
					},
				})},
			wantErr: false,
			wantResMessage: &evidence.StoreEvidencesResponse{
				Status:        false,
				StatusMessage: "evidence.target_of_evaluation_id: value must be a valid UUID",
			},
		},
		{
			name: "Error in stream to server - Recv()-err",
			args: args{
				streamToServerWithRecvErr: createMockStreamWithRecvErr(createStoreEvidenceRequestMocks(t, 1))},
			wantErr:        true,
			wantErrMessage: "rpc error: code = Unknown desc = cannot receive stream request",
		},
		{
			name: "Error in stream to client - Send()-err",
			args: args{
				streamToClientWithSendErr: createMockStreamWithSendErr(createStoreEvidenceRequestMocks(t, 1))},
			wantErr:        true,
			wantErrMessage: "rpc error: code = Unknown desc = cannot send response to the client:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				err                error
				responseFromServer *evidence.StoreEvidencesResponse
			)
			s := NewService()

			if tt.args.streamToServer != nil {
				err = s.StoreEvidences(tt.args.streamToServer)
				responseFromServer = <-tt.args.streamToServer.SentFromServer
			} else if tt.args.streamToClientWithSendErr != nil {
				err = s.StoreEvidences(tt.args.streamToClientWithSendErr)
			} else if tt.args.streamToServerWithRecvErr != nil {
				err = s.StoreEvidences(tt.args.streamToServerWithRecvErr)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Got AssessEvidence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.Nil(t, err)
				assert.Equal(t, tt.wantResMessage.Status, responseFromServer.Status)
				// We have to check both ways, as it fails if one StatusMessage is empty.
				assert.Contains(t, responseFromServer.StatusMessage, tt.wantResMessage.StatusMessage)
				assert.Contains(t, responseFromServer.StatusMessage, tt.wantResMessage.StatusMessage)
			} else {
				assert.ErrorContains(t, err, tt.wantErrMessage)
			}
		})
	}
}

// TestListEvidences tests List req
func TestService_ListEvidences(t *testing.T) {
	type fields struct {
		storage persistence.Storage
		authz   service.AuthorizationStrategy
	}
	type args struct {
		in0 context.Context
		req *evidence.ListEvidencesRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes assert.Want[*evidence.ListEvidencesResponse]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Successful List Of Evidences (with allowed target of evaluation)",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(false, evidencetest.MockEvidence1.TargetOfEvaluationId, evidencetest.MockEvidence2.TargetOfEvaluationId),
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(&evidencetest.MockEvidence1))
					assert.NoError(t, s.Create(&evidencetest.MockEvidence2))
				}),
			},
			args: args{
				in0: context.TODO(),
				req: &evidence.ListEvidencesRequest{},
			},
			wantErr: assert.NoError,
			wantRes: func(t *testing.T, got *evidence.ListEvidencesResponse) bool {
				return assert.Equal(t, len(got.Evidences), 2) &&
					assert.Equal(t, evidencetest.MockEvidence1.TargetOfEvaluationId, got.Evidences[0].TargetOfEvaluationId) &&
					assert.Equal(t, evidencetest.MockEvidence2.TargetOfEvaluationId, got.Evidences[1].TargetOfEvaluationId)
			},
		},
		{
			name: "Successful Filter Of Evidences (with allowed target of evaluation)",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(false, evidencetest.MockEvidence1.TargetOfEvaluationId),
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(&evidencetest.MockEvidence1))
					assert.NoError(t, s.Create(&evidencetest.MockEvidence2))
				}),
			},
			args: args{
				in0: context.TODO(),
				req: evidencetest.MockListEvidenceRequest1,
			},
			wantErr: assert.NoError,
			wantRes: func(t *testing.T, got *evidence.ListEvidencesResponse) bool {
				for _, e := range got.Evidences {
					assert.Equal(t, *evidencetest.MockListEvidenceRequest1.Filter.TargetOfEvaluationId, e.TargetOfEvaluationId)
					assert.Equal(t, *evidencetest.MockListEvidenceRequest1.Filter.ToolId, e.ToolId)
				}

				return true
			},
		},
		{
			name: "Only target_of_evaluation_Id filter applied, when Tool_Id filter off",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(true),
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(&evidencetest.MockEvidence1))
					assert.NoError(t, s.Create(&evidencetest.MockEvidence2))
				}),
			},
			args: args{
				in0: context.TODO(),
				req: &evidence.ListEvidencesRequest{
					PageSize:  evidencetest.MockListEvidenceRequest1.PageSize,
					PageToken: evidencetest.MockListEvidenceRequest1.PageToken,
					OrderBy:   evidencetest.MockListEvidenceRequest1.OrderBy,
					Asc:       evidencetest.MockListEvidenceRequest1.Asc,
					Filter: &evidence.Filter{
						TargetOfEvaluationId: evidencetest.MockListEvidenceRequest1.Filter.TargetOfEvaluationId,
					},
				},
			},
			wantErr: assert.NoError,
			wantRes: func(t *testing.T, got *evidence.ListEvidencesResponse) bool {
				for _, r := range got.Evidences {
					assert.Equal(t, *evidencetest.MockListEvidenceRequest1.Filter.TargetOfEvaluationId, r.TargetOfEvaluationId)
					assert.Equal(t, *evidencetest.MockListEvidenceRequest1.Filter.ToolId, r.ToolId)
				}

				return true
			},
		},
		{
			name: "Only Tool_Id filter applied, when target_of_evaluation_Id filter off",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(true),
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(&evidencetest.MockEvidence1))
					assert.NoError(t, s.Create(&evidencetest.MockEvidence2))
				}),
			},
			args: args{
				in0: context.TODO(),
				req: &evidence.ListEvidencesRequest{
					PageSize:  evidencetest.MockListEvidenceRequest2.PageSize,
					PageToken: evidencetest.MockListEvidenceRequest2.PageToken,
					OrderBy:   evidencetest.MockListEvidenceRequest2.OrderBy,
					Asc:       evidencetest.MockListEvidenceRequest2.Asc,
					Filter: &evidence.Filter{
						ToolId: evidencetest.MockListEvidenceRequest2.Filter.ToolId,
					},
				},
			},
			wantErr: assert.NoError,
			wantRes: func(t *testing.T, got *evidence.ListEvidencesResponse) bool {
				assert.Equal(t, 1, len(got.Evidences))

				// Loop through all received evidences and check whether tool and service ids are correct.
				for _, r := range got.Evidences {
					assert.Equal(t, *evidencetest.MockListEvidenceRequest2.Filter.TargetOfEvaluationId, r.TargetOfEvaluationId)
					assert.Equal(t, *evidencetest.MockListEvidenceRequest2.Filter.ToolId, r.ToolId)
				}

				return true
			},
		},
		{
			name: "Permission denied (target of evaluation id not allowed)",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID1), // allow only MockTargetOfEvaluationID
			},
			args: args{
				in0: context.TODO(),
				req: &evidence.ListEvidencesRequest{
					Filter: &evidence.Filter{
						TargetOfEvaluationId: util.Ref(testdata.MockTargetOfEvaluationID2),
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, status.Code(err), codes.PermissionDenied) // MockTargetOfEvaluationID2 is not allowed
				return assert.ErrorContains(t, err, service.ErrPermissionDenied.Error())
			},
			wantRes: assert.Nil[*evidence.ListEvidencesResponse],
		},
		{
			name: "Wrong Input handled correctly (req = nil)",
			args: args{
				in0: context.TODO(),
				req: nil,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.ErrorContains(t, err, api.ErrEmptyRequest.Error())
				return assert.Equal(t, status.Code(err), codes.InvalidArgument)
			},
			wantRes: assert.Nil[*evidence.ListEvidencesResponse],
		},
		{
			name: "Wrong Input handled correctly (tool_id not UUID)",
			args: args{
				in0: nil,
				req: &evidence.ListEvidencesRequest{
					PageSize:  evidencetest.MockListEvidenceRequest2.PageSize,
					PageToken: evidencetest.MockListEvidenceRequest2.PageToken,
					OrderBy:   evidencetest.MockListEvidenceRequest2.OrderBy,
					Asc:       evidencetest.MockListEvidenceRequest2.Asc,
					Filter: &evidence.Filter{
						ToolId: util.Ref("No UUID Format"),
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, status.Code(err), codes.InvalidArgument)
			},
			wantRes: assert.Nil[*evidence.ListEvidencesResponse],
		},
		{
			name: "Wrong Input handled correctly (target_of_evaluation_id not UUID)",
			args: args{
				in0: nil,
				req: &evidence.ListEvidencesRequest{
					PageSize:  evidencetest.MockListEvidenceRequest2.PageSize,
					PageToken: evidencetest.MockListEvidenceRequest2.PageToken,
					OrderBy:   evidencetest.MockListEvidenceRequest2.OrderBy,
					Asc:       evidencetest.MockListEvidenceRequest2.Asc,
					Filter: &evidence.Filter{
						TargetOfEvaluationId: util.Ref("No UUID Format"),
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, status.Code(err), codes.InvalidArgument)
			},
			wantRes: assert.Nil[*evidence.ListEvidencesResponse],
		},
		{
			name: "DB (pagination) error",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(true),
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(&orchestratortest.MockAssessmentResult1))
				}),
			},
			args: args{
				in0: context.TODO(),
				req: &evidence.ListEvidencesRequest{
					PageSize:  evidencetest.MockListEvidenceRequest2.PageSize,
					PageToken: evidencetest.MockListEvidenceRequest2.PageToken,
					OrderBy:   "Wrong Input",
					Asc:       evidencetest.MockListEvidenceRequest2.Asc,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.ErrorContains(t, err, "could not paginate results")
				return assert.Equal(t, status.Code(err), codes.Internal)
			},
			wantRes: assert.Nil[*evidence.ListEvidencesResponse],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				storage: tt.fields.storage,
				authz:   tt.fields.authz,
			}

			gotRes, err := svc.ListEvidences(tt.args.in0, tt.args.req)
			tt.wantErr(t, err)
			tt.wantRes(t, gotRes)
		})
	}
}

func TestService_EvidenceHook(t *testing.T) {
	var (
		hookCallCounter = 0
		wg              sync.WaitGroup
	)
	wg.Add(2)

	firstHookFunction := func(ctx context.Context, evidence *evidence.Evidence, err error) {
		hookCallCounter++
		log.Println("Hello from inside the firstHookFunction")

		wg.Done()
	}

	secondHookFunction := func(ctx context.Context, evidence *evidence.Evidence, err error) {
		hookCallCounter++
		log.Println("Hello from inside the secondHookFunction")

		wg.Done()
	}

	svc := NewService()
	svc.RegisterEvidenceHook(firstHookFunction)
	svc.RegisterEvidenceHook(secondHookFunction)

	// Check if first hook is registered
	funcName1 := runtime.FuncForPC(reflect.ValueOf(svc.evidenceHooks[0]).Pointer()).Name()
	funcName2 := runtime.FuncForPC(reflect.ValueOf(firstHookFunction).Pointer()).Name()
	assert.Equal(t, funcName1, funcName2)

	// Check if second hook is registered
	funcName1 = runtime.FuncForPC(reflect.ValueOf(svc.evidenceHooks[1]).Pointer()).Name()
	funcName2 = runtime.FuncForPC(reflect.ValueOf(secondHookFunction).Pointer()).Name()
	assert.Equal(t, funcName1, funcName2)

	// Check GRPC call
	type args struct {
		in0      context.Context
		evidence *evidence.StoreEvidenceRequest
	}
	tests := []struct {
		name    string
		args    args
		wantRes *evidence.StoreEvidenceResponse
		wantErr assert.WantErr
	}{
		{
			name: "Store an evidence to the map",
			args: args{
				in0: context.TODO(),
				evidence: &evidence.StoreEvidenceRequest{Evidence: &evidence.Evidence{
					Id:                   testdata.MockEvidenceID1,
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID2,
					Timestamp:            timestamppb.Now(),
					ToolId:               "mockToolId-1",
					Resource: prototest.NewAny(t, &ontology.VirtualMachine{
						Id: "mock-id-1",
					}),
				},
				},
			},
			wantErr: assert.Nil[error],
			wantRes: &evidence.StoreEvidenceResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := svc
			gotResp, err := s.StoreEvidence(tt.args.in0, tt.args.evidence)

			// wait for all hooks (2 hooks)
			wg.Wait()

			tt.wantErr(t, err)
			assert.Equal(t, tt.wantRes, gotResp)

			// Check if evidence is stored in storage DB
			var evidences []evidence.Evidence
			err = s.storage.List(&evidences, "", true, 0, -1)
			assert.NoError(t, err)
			assert.NotEmpty(t, evidences)
			assert.Equal(t, 2, hookCallCounter)
		})
	}

}

// createStoreEvidenceRequestMocks creates store evidence requests with random evidence IDs
func createStoreEvidenceRequestMocks(t *testing.T, count int) []*evidence.StoreEvidenceRequest {
	var mockRequests []*evidence.StoreEvidenceRequest

	for i := 0; i < count; i++ {
		evidenceRequest := &evidence.StoreEvidenceRequest{
			Evidence: &evidence.Evidence{
				Id:                   uuid.NewString(),
				ToolId:               fmt.Sprintf("MockToolId-%d", i),
				TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
				Timestamp:            timestamppb.Now(),
				Resource: prototest.NewAny(t, &ontology.VirtualMachine{
					Id: "mock-id-1",
				}),
			},
		}
		mockRequests = append(mockRequests, evidenceRequest)
	}

	return mockRequests
}

type mockStreamer struct {
	grpc.ServerStream
	RecvToServer   chan *evidence.StoreEvidenceRequest
	SentFromServer chan *evidence.StoreEvidencesResponse
}

func createMockStream(requests []*evidence.StoreEvidenceRequest) *mockStreamer {
	m := &mockStreamer{
		RecvToServer: make(chan *evidence.StoreEvidenceRequest, len(requests)),
	}
	for _, req := range requests {
		m.RecvToServer <- req
	}

	m.SentFromServer = make(chan *evidence.StoreEvidencesResponse, len(requests))
	return m
}

func (m *mockStreamer) Send(response *evidence.StoreEvidencesResponse) error {
	m.SentFromServer <- response
	return nil
}

func (*mockStreamer) SendAndClose(_ *emptypb.Empty) error {
	return nil
}

func (m *mockStreamer) Recv() (req *evidence.StoreEvidenceRequest, err error) {
	if len(m.RecvToServer) == 0 {
		return nil, io.EOF
	}
	req, more := <-m.RecvToServer
	if !more {
		return nil, errors.New("empty")
	}

	return req, nil
}

func (*mockStreamer) SetHeader(_ metadata.MD) error {
	panic("implement me")
}

func (*mockStreamer) SendHeader(_ metadata.MD) error {
	panic("implement me")
}

func (*mockStreamer) SetTrailer(_ metadata.MD) {
	panic("implement me")
}

func (*mockStreamer) Context() context.Context {
	return context.TODO()
}

func (*mockStreamer) SendMsg(_ interface{}) error {
	panic("implement me")
}

func (*mockStreamer) RecvMsg(_ interface{}) error {
	panic("implement me")
}

type mockStreamerWithRecvErr struct {
	grpc.ServerStream
	RecvToServer   chan *evidence.StoreEvidenceRequest
	SentFromServer chan *evidence.StoreEvidencesResponse
}

func (*mockStreamerWithRecvErr) Context() context.Context {
	return context.TODO()
}

func (*mockStreamerWithRecvErr) Send(*evidence.StoreEvidencesResponse) error {
	panic("implement me")
}

func (*mockStreamerWithRecvErr) Recv() (*evidence.StoreEvidenceRequest, error) {

	err := errors.New("Recv()-error")

	return nil, err
}

func createMockStreamWithRecvErr(requests []*evidence.StoreEvidenceRequest) *mockStreamerWithRecvErr {
	m := &mockStreamerWithRecvErr{
		RecvToServer: make(chan *evidence.StoreEvidenceRequest, len(requests)),
	}
	for _, req := range requests {
		m.RecvToServer <- req
	}

	m.SentFromServer = make(chan *evidence.StoreEvidencesResponse, len(requests))
	return m
}

type mockStreamerWithSendErr struct {
	grpc.ServerStream
	RecvToServer   chan *evidence.StoreEvidenceRequest
	SentFromServer chan *evidence.StoreEvidencesResponse
}

func (*mockStreamerWithSendErr) Context() context.Context {
	return context.TODO()
}

func (*mockStreamerWithSendErr) Send(*evidence.StoreEvidencesResponse) error {
	return errors.New("Send()-err")
}

func createMockStreamWithSendErr(requests []*evidence.StoreEvidenceRequest) *mockStreamerWithSendErr {
	m := &mockStreamerWithSendErr{
		RecvToServer: make(chan *evidence.StoreEvidenceRequest, len(requests)),
	}
	for _, req := range requests {
		m.RecvToServer <- req
	}

	m.SentFromServer = make(chan *evidence.StoreEvidencesResponse, len(requests))
	return m
}

func (m *mockStreamerWithSendErr) Recv() (req *evidence.StoreEvidenceRequest, err error) {
	if len(m.RecvToServer) == 0 {
		return nil, io.EOF
	}
	req, more := <-m.RecvToServer
	if !more {
		return nil, errors.New("empty")
	}

	return req, nil
}

func TestService_GetEvidence(t *testing.T) {
	type fields struct {
		storage persistence.Storage
		authz   service.AuthorizationStrategy
	}
	type args struct {
		ctx context.Context
		req *evidence.GetEvidenceRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.Want[*evidence.Evidence]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "permission denied (not found)",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(&evidence.Evidence{
						Id:                   testdata.MockEvidenceID1,
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID2,
						ToolId:               testdata.MockEvidenceToolID1,
						Resource:             nil,
						Timestamp:            timestamppb.Now(),
					}))
				}),
				authz: servicetest.NewAuthorizationStrategy(false),
			},
			args: args{
				ctx: context.TODO(),
				req: &evidence.GetEvidenceRequest{
					EvidenceId: testdata.MockEvidenceID1,
				},
			},
			want: assert.Nil[*evidence.Evidence],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.NotFound, status.Code(err))
				return assert.ErrorContains(t, err, "evidence not found")
			},
		},
		{
			name: "valid evidence",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(&evidence.Evidence{
						Id:                   testdata.MockEvidenceID1,
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
						ToolId:               testdata.MockEvidenceToolID1,
						Resource:             nil,
						Timestamp:            timestamppb.Now(),
					}))
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &evidence.GetEvidenceRequest{
					EvidenceId: testdata.MockEvidenceID1,
				},
			},
			wantErr: assert.NoError,
			want: func(t *testing.T, got *evidence.Evidence) bool {
				return assert.NoError(t, api.Validate(got))
			},
		},
		{
			name: "invalid UUID",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &evidence.GetEvidenceRequest{
					EvidenceId: "not valid",
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "evidence_id: value must be a valid UUID")
			},
			want: assert.Nil[*evidence.Evidence],
		},
		{
			name: "evidence not found",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &evidence.GetEvidenceRequest{
					EvidenceId: testdata.MockEvidenceID1,
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.NotFound, status.Code(err))
				return assert.ErrorContains(t, err, "evidence not found")
			},
			want: assert.Nil[*evidence.Evidence],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				storage: tt.fields.storage,
				authz:   tt.fields.authz,
			}
			gotRes, err := svc.GetEvidence(tt.args.ctx, tt.args.req)
			tt.wantErr(t, err)
			tt.want(t, gotRes)
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
