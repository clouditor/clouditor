// Copyright 2016-2020 Fraunhofer AISEC
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

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/voc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

// TestNewService is a simply test for NewService
func TestNewService(t *testing.T) {
	tests := []struct {
		name string
		want evidence.EvidenceStoreServer
	}{
		{
			name: "EvidenceStoreServer created with empty req map",
			want: &Service{
				evidences:                        make(map[string]*evidence.Evidence),
				UnimplementedEvidenceStoreServer: evidence.UnimplementedEvidenceStoreServer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestStoreEvidence tests StoreEvidence
func TestStoreEvidence(t *testing.T) {
	type args struct {
		in0 context.Context
		req *evidence.StoreEvidenceRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *evidence.StoreEvidenceResponse
		wantErr  bool
	}{
		{
			name: "Store req to the map",
			args: args{
				in0: context.TODO(),
				req: &evidence.StoreEvidenceRequest{
					Evidence: &evidence.Evidence{
						Id:        "11111111-1111-1111-1111-111111111111",
						ServiceId: "MockServiceId",
						ToolId:    "MockTool",
						Timestamp: timestamppb.Now(),
						Raw:       "",
						Resource: toStruct(voc.VirtualMachine{
							Compute: &voc.Compute{Resource: &voc.Resource{
								ID: "mock-id",
							}},
						}, t),
					}},
			},
			wantErr:  false,
			wantResp: &evidence.StoreEvidenceResponse{Status: true},
		},
		{
			name: "Store an evidence without toolId to the map",
			args: args{
				in0: context.TODO(),
				req: &evidence.StoreEvidenceRequest{
					Evidence: &evidence.Evidence{
						Id:        "MockEvidenceId-1",
						ServiceId: "MockServiceId-1",
						Timestamp: timestamppb.Now(),
						Raw:       "",
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
			wantErr: true,
			wantResp: &evidence.StoreEvidenceResponse{
				Status:        false,
				StatusMessage: "invalid evidence: tool id in evidence is missing",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			gotResp, err := s.StoreEvidence(tt.args.in0, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreEvidence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("StoreEvidence() gotResp = %v, want %v", gotResp, tt.wantResp)
			}

			if gotResp.Status {
				assert.NotNil(t, s.evidences["11111111-1111-1111-1111-111111111111"])
			} else {
				assert.Empty(t, s.evidences)
			}
		})
	}
}

// TestStoreEvidences tests StoreEvidences
func TestStoreEvidences(t *testing.T) {
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
		wantRespMessage *evidence.StoreEvidenceResponse
	}{
		{
			name:   "Store 1 evidence to the map",
			fields: fields{count: 1},
			args: args{
				streamToServer: createMockStream(createStoreEvidenceRequestMocks(1))},
			wantErr: false,
			wantRespMessage: &evidence.StoreEvidenceResponse{
				Status: true,
			},
		},
		{
			name:   "Store 2 evidences to the map",
			fields: fields{count: 2},
			args: args{
				streamToServer: createMockStream(createStoreEvidenceRequestMocks(2))},
			wantErr: false,
			wantRespMessage: &evidence.StoreEvidenceResponse{
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
							Id:        uuid.NewString(),
							ServiceId: "MockServiceId",
							Timestamp: timestamppb.Now(),
							Raw:       "",
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
			wantRespMessage: &evidence.StoreEvidenceResponse{
				Status:        false,
				StatusMessage: "invalid evidence:",
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
				responseFromServer *evidence.StoreEvidenceResponse
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
				assert.Contains(t, err.Error(), tt.wantErrMessage)
			}
		})
	}
}

// TestListEvidences tests List req
func TestListEvidences(t *testing.T) {
	s := NewService()
	s.evidences["MockEvidenceId-1"] = &evidence.Evidence{
		Id:        "MockEvidenceId-1",
		ServiceId: "MockServiceId-1",
		Timestamp: timestamppb.Now(),
		Raw:       "",
		Resource:  nil,
	}
	s.evidences["MockEvidenceId-2"] = &evidence.Evidence{
		Id:        "MockEvidenceId-2",
		ServiceId: "MockServiceId-2",
		Timestamp: timestamppb.Now(),
		Raw:       "",
		Resource:  nil,
	}

	resp, err := s.ListEvidences(context.TODO(), &evidence.ListEvidencesRequest{})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(resp.Evidences))
}

func TestEvidenceHook(t *testing.T) {
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

	service := NewService()
	service.RegisterEvidenceHook(firstHookFunction)
	service.RegisterEvidenceHook(secondHookFunction)

	// Check if first hook is registered
	funcName1 := runtime.FuncForPC(reflect.ValueOf(service.evidenceHooks[0]).Pointer()).Name()
	funcName2 := runtime.FuncForPC(reflect.ValueOf(firstHookFunction).Pointer()).Name()
	assert.Equal(t, funcName1, funcName2)

	// Check if second hook is registered
	funcName1 = runtime.FuncForPC(reflect.ValueOf(service.evidenceHooks[1]).Pointer()).Name()
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
					Id:        "11111111-1111-1111-1111-111111111111",
					ServiceId: "MockServiceId-1",
					Timestamp: timestamppb.Now(),
					Raw:       "",
					ToolId:    "mockToolId-1",
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
			wantErr: false,
			wantResp: &evidence.StoreEvidenceResponse{
				Status: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hookCallCounter = 0
			s := service
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
			assert.NotEmpty(t, s.evidences)
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
				Id:        uuid.NewString(),
				ToolId:    fmt.Sprintf("MockToolId-%d", i),
				ServiceId: fmt.Sprintf("MockServiceId-%d", i),
				Timestamp: timestamppb.Now(),
				Raw:       "",
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
	SentFromServer chan *evidence.StoreEvidenceResponse
}

func createMockStream(requests []*evidence.StoreEvidenceRequest) *mockStreamer {
	m := &mockStreamer{
		RecvToServer: make(chan *evidence.StoreEvidenceRequest, len(requests)),
	}
	for _, req := range requests {
		m.RecvToServer <- req
	}

	m.SentFromServer = make(chan *evidence.StoreEvidenceResponse, len(requests))
	return m
}

func (m mockStreamer) Send(response *evidence.StoreEvidenceResponse) error {
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
	SentFromServer chan *evidence.StoreEvidenceResponse
}

func (mockStreamerWithRecvErr) Send(*evidence.StoreEvidenceResponse) error {
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

	m.SentFromServer = make(chan *evidence.StoreEvidenceResponse, len(requests))
	return m
}

type mockStreamerWithSendErr struct {
	grpc.ServerStream
	RecvToServer   chan *evidence.StoreEvidenceRequest
	SentFromServer chan *evidence.StoreEvidenceResponse
}

func (*mockStreamerWithSendErr) Send(*evidence.StoreEvidenceResponse) error {
	return errors.New("Send()-err")
}

func createMockStreamWithSendErr(requests []*evidence.StoreEvidenceRequest) *mockStreamerWithSendErr {
	m := &mockStreamerWithSendErr{
		RecvToServer: make(chan *evidence.StoreEvidenceRequest, len(requests)),
	}
	for _, req := range requests {
		m.RecvToServer <- req
	}

	m.SentFromServer = make(chan *evidence.StoreEvidenceResponse, len(requests))
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
		log.Errorf("eror getting struct of resource: %v", err)
	}

	return
}
