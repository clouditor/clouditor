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

package evidences

import (
	"context"
	"errors"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"sync"
	"testing"

	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/persistence/gorm"
	"clouditor.io/clouditor/service"
	"clouditor.io/clouditor/voc"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TestNewService is a simply test for NewService
func TestNewService(t *testing.T) {
	db, err := gorm.NewStorage(gorm.WithInMemory())
	assert.NoError(t, err)
	type args struct {
		opts []service.Option[Service]
	}
	tests := []struct {
		name string
		args args
		want assert.ValueAssertionFunc
	}{
		{
			name: "EvidenceStoreServer created without options",
			want: func(t assert.TestingT, i interface{}, i3 ...interface{}) bool {
				svc, ok := i.(*Service)
				assert.True(t, ok)
				// Storage should be default (in-memory storage). Hard to check since its type is not exported
				assert.NotNil(t, svc.storage)
				return true
			},
		},
		{
			name: "EvidenceStoreServer created with storage option",
			args: args{opts: []service.Option[Service]{WithStorage(db)}},
			want: func(t assert.TestingT, i interface{}, i3 ...interface{}) bool {
				svc, ok := i.(*Service)
				assert.True(t, ok)
				// Storage should be gorm (in-memory storage). Hard to check since its type is not exported
				assert.NotNil(t, svc.storage)
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
		name     string
		args     args
		wantResp *evidence.StoreEvidenceResponse
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "Store req to the map",
			args: args{
				in0: context.TODO(),
				req: &evidence.StoreEvidenceRequest{
					Evidence: &evidence.Evidence{
						Id:             testdata.MockEvidenceID,
						CloudServiceId: testdata.MockCloudServiceID,
						ToolId:         testdata.MockEvidenceToolID,
						Timestamp:      timestamppb.Now(),
						Raw:            nil,
						Resource: toStruct(voc.VirtualMachine{
							Compute: &voc.Compute{Resource: &voc.Resource{
								ID: "mock-id",
							}},
						}, t),
					}},
			},
			wantErr:  assert.NoError,
			wantResp: &evidence.StoreEvidenceResponse{},
		},
		{
			name: "Store an evidence without toolId to the map",
			args: args{
				in0: context.TODO(),
				req: &evidence.StoreEvidenceRequest{
					Evidence: &evidence.Evidence{
						Id:             testdata.MockEvidenceID,
						CloudServiceId: testdata.MockCloudServiceID,
						Timestamp:      timestamppb.Now(),
						Raw:            nil,
						Resource: toStruct(voc.VirtualMachine{
							Compute: &voc.Compute{
								Resource: &voc.Resource{
									ID: "mock-id-1",
								},
							},
						}, t),
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "Evidence.ToolId: value length must be at least 1 runes")
			},
			wantResp: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			gotResp, err := s.StoreEvidence(tt.args.in0, tt.args.req)

			tt.wantErr(t, err)

			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("StoreEvidence() gotResp = %v, want %v", gotResp, tt.wantResp)
			}

			if gotResp != nil {
				e := &evidence.Evidence{}
				err := s.storage.Get(e)
				assert.NoError(t, err)
				assert.Equal(t, tt.args.req.Evidence.Id, e.Id)
			}
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
		name            string
		fields          fields
		args            args
		wantErr         bool
		wantErrMessage  string
		wantRespMessage *evidence.StoreEvidencesResponse
	}{
		{
			name:   "Store 1 evidence to the map",
			fields: fields{count: 1},
			args: args{
				streamToServer: createMockStream(createStoreEvidenceRequestMocks(1))},
			wantErr: false,
			wantRespMessage: &evidence.StoreEvidencesResponse{
				Status: true,
			},
		},
		{
			name:   "Store 2 evidences to the map",
			fields: fields{count: 2},
			args: args{
				streamToServer: createMockStream(createStoreEvidenceRequestMocks(2))},
			wantErr: false,
			wantRespMessage: &evidence.StoreEvidencesResponse{
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
							Id:             uuid.NewString(),
							CloudServiceId: "MockCloudServiceId",
							Timestamp:      timestamppb.Now(),
							Raw:            nil,
							Resource: toStructWithoutTest(voc.VirtualMachine{
								Compute: &voc.Compute{
									Resource: &voc.Resource{
										ID: "mock-id-1",
									},
								},
							}),
						},
					},
				})},
			wantErr: false,
			wantRespMessage: &evidence.StoreEvidencesResponse{
				Status:        false,
				StatusMessage: "rpc error: code = InvalidArgument desc = invalid request: invalid StoreEvidenceRequest.Evidence: embedded message failed validation | caused by: invalid Evidence.CloudServiceId: value must be a valid UUID | caused by: invalid uuid format",
			},
		},
		{
			name: "Error in stream to server - Recv()-err",
			args: args{
				streamToServerWithRecvErr: createMockStreamWithRecvErr(createStoreEvidenceRequestMocks(1))},
			wantErr:        true,
			wantErrMessage: "rpc error: code = Unknown desc = cannot receive stream request",
		},
		{
			name: "Error in stream to client - Send()-err",
			args: args{
				streamToClientWithSendErr: createMockStreamWithSendErr(createStoreEvidenceRequestMocks(1))},
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
				assert.Contains(t, responseFromServer.StatusMessage, tt.wantRespMessage.StatusMessage)
			} else {
				assert.ErrorContains(t, err, tt.wantErrMessage)
			}
		})
	}
}

// TestListEvidences tests List req
func TestService_ListEvidences(t *testing.T) {
	// TODO(oxisto): Convert this test to a table test
	s := NewService()
	err := s.storage.Create(&evidence.Evidence{
		Id:             testdata.MockEvidenceID,
		CloudServiceId: testdata.MockCloudServiceID,
		Timestamp:      timestamppb.Now(),
		Raw:            util.Ref(""),
		Resource:       structpb.NewNullValue(),
	})
	assert.NoError(t, err)
	err = s.storage.Create(&evidence.Evidence{
		Id:             testdata.MockAnotherEvidenceID,
		CloudServiceId: testdata.MockAnotherCloudServiceID,
		Timestamp:      timestamppb.Now(),
		Raw:            util.Ref(""),
		Resource:       structpb.NewNullValue(),
	})
	assert.NoError(t, err)

	resp, err := s.ListEvidences(context.TODO(), &evidence.ListEvidencesRequest{})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(resp.Evidences))

	s.authz = &service.AuthorizationStrategyJWT{CloudServicesKey: testutil.TestCustomClaims, AllowAllKey: testutil.TestAllowAllClaims}

	resp, err = s.ListEvidences(testutil.TestContextOnlyService1, &evidence.ListEvidencesRequest{})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Evidences))
}

func TestService_EvidenceHook(t *testing.T) {
	var (
		hookCallCounter = 0
		wg              sync.WaitGroup
	)
	wg.Add(2)

	firstHookFunction := func(evidence *evidence.Evidence, err error) {
		hookCallCounter++
		log.Println("Hello from inside the firstHookFunction")

		wg.Done()
	}

	secondHookFunction := func(evidence *evidence.Evidence, err error) {
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
		name     string
		args     args
		wantResp *evidence.StoreEvidenceResponse
		wantErr  bool
	}{
		{
			name: "Store an evidence to the map",
			args: args{
				in0: context.TODO(),
				evidence: &evidence.StoreEvidenceRequest{Evidence: &evidence.Evidence{
					Id:             testdata.MockEvidenceID,
					CloudServiceId: testdata.MockAnotherCloudServiceID,
					Timestamp:      timestamppb.Now(),
					Raw:            nil,
					ToolId:         "mockToolId-1",
					Resource: toStruct(voc.VirtualMachine{
						Compute: &voc.Compute{
							Resource: &voc.Resource{
								ID: "mock-id-1",
							},
						},
					}, t),
				},
				},
			},
			wantErr:  false,
			wantResp: &evidence.StoreEvidenceResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hookCallCounter = 0
			s := svc
			gotResp, err := s.StoreEvidence(tt.args.in0, tt.args.evidence)

			// wait for all hooks (2 hooks)
			wg.Wait()

			if (err != nil) != tt.wantErr {
				t.Errorf("StoreEvidence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("StoreEvidence() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
			var evidences []evidence.Evidence
			err = s.storage.List(&evidences, "", true, 0, -1)
			assert.NoError(t, err)
			assert.NotEmpty(t, evidences)
			assert.Equal(t, 2, hookCallCounter)
		})
	}

}

// createStoreEvidenceRequestMocks creates store evidence requests with random evidence IDs
func createStoreEvidenceRequestMocks(count int) []*evidence.StoreEvidenceRequest {
	var mockRequests []*evidence.StoreEvidenceRequest

	for i := 0; i < count; i++ {
		evidenceRequest := &evidence.StoreEvidenceRequest{
			Evidence: &evidence.Evidence{
				Id:             uuid.NewString(),
				ToolId:         fmt.Sprintf("MockToolId-%d", i),
				CloudServiceId: fmt.Sprintf("MockCloudServiceId-%d", i),
				Timestamp:      timestamppb.Now(),
				Raw:            nil,
				Resource: toStructWithoutTest(voc.VirtualMachine{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID: "mock-id-1",
						},
					},
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

func (m mockStreamer) Send(response *evidence.StoreEvidencesResponse) error {
	m.SentFromServer <- response
	return nil
}

func (mockStreamer) SendAndClose(_ *emptypb.Empty) error {
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

type mockStreamerWithRecvErr struct {
	grpc.ServerStream
	RecvToServer   chan *evidence.StoreEvidenceRequest
	SentFromServer chan *evidence.StoreEvidencesResponse
}

func (mockStreamerWithRecvErr) Send(*evidence.StoreEvidencesResponse) error {
	panic("implement me")
}

func (mockStreamerWithRecvErr) Recv() (*evidence.StoreEvidenceRequest, error) {

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

func toStruct(r voc.IsCloudResource, t *testing.T) (s *structpb.Value) {
	s, err := voc.ToStruct(r)
	if err != nil {
		assert.Error(t, err)
	}

	return
}

func toStructWithoutTest(r voc.IsCloudResource) (s *structpb.Value) {
	s, err := voc.ToStruct(r)
	if err != nil {
		log.Errorf("error getting struct of resource: %v", err)
	}

	return
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
		want    assert.ValueAssertionFunc
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "permission denied (not found)",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(&evidence.Evidence{
						Id:             testdata.MockEvidenceID,
						CloudServiceId: testdata.MockAnotherCloudServiceID,
						ToolId:         testdata.MockEvidenceToolID,
						Resource:       structpb.NewNullValue(),
						Timestamp:      timestamppb.Now(),
					}))
				}),
				authz: &service.AuthorizationStrategyJWT{CloudServicesKey: testutil.TestCustomClaims, AllowAllKey: testutil.TestAllowAllClaims},
			},
			args: args{
				ctx: testutil.TestContextOnlyService1,
				req: &evidence.GetEvidenceRequest{
					EvidenceId: testdata.MockEvidenceID,
				},
			},
			want: assert.Nil,
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
						Id:             testdata.MockEvidenceID,
						CloudServiceId: testdata.MockCloudServiceID,
						ToolId:         testdata.MockEvidenceToolID,
						Resource:       structpb.NewNullValue(),
						Timestamp:      timestamppb.Now(),
					}))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				req: &evidence.GetEvidenceRequest{
					EvidenceId: testdata.MockEvidenceID,
				},
			},
			wantErr: assert.NoError,
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				res := i1.(*evidence.Evidence)
				return assert.NoError(t, res.Validate())
			},
		},
		{
			name: "invalid UUID",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				req: &evidence.GetEvidenceRequest{
					EvidenceId: "not valid",
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "EvidenceId: value must be a valid UUID")
			},
			want: assert.Nil,
		},
		{
			name: "evidence not found",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				req: &evidence.GetEvidenceRequest{
					EvidenceId: testdata.MockEvidenceID,
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.NotFound, status.Code(err))
				return assert.ErrorContains(t, err, "evidence not found")
			},
			want: assert.Nil,
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
