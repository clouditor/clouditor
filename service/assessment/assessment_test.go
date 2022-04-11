// Copyright 2021-2022 Fraunhofer AISEC
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

package assessment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"sync"
	"testing"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/internal/testutil/clitest"
	"clouditor.io/clouditor/voc"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	authPort int
)

func TestMain(m *testing.M) {
	clitest.AutoChdir()

	server, _, _ := startBufConnServer()

	code := m.Run()

	server.Stop()

	os.Exit(code)
}

// TestNewService is a simply test for NewService
func TestNewService(t *testing.T) {
	type args struct {
		opts []ServiceOption
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{
			name: "AssessmentServer created with empty results map",
			want: &Service{
				results:              make(map[string]*assessment.AssessmentResult),
				evidenceStoreAddress: "localhost:9090",
				evidenceStoreStreams: api.NewStreamsOf(api.WithLogger[evidence.EvidenceStore_StoreEvidencesClient, *evidence.StoreEvidenceRequest](log)),
				orchestratorAddress:  "localhost:9090",
				cachedConfigurations: make(map[string]cachedConfiguration),
				orchestratorChannel:  nil,
			},
		},
		{
			name: "AssessmentServer created with options",
			args: args{
				opts: []ServiceOption{
					WithEvidenceStoreAddress("localhost:9091"),
					WithOrchestratorAddress("localhost:9092"),
				},
			},
			want: &Service{
				results:              make(map[string]*assessment.AssessmentResult),
				evidenceStoreAddress: "localhost:9091",
				evidenceStoreStreams: api.NewStreamsOf(api.WithLogger[evidence.EvidenceStore_StoreEvidencesClient, *evidence.StoreEvidenceRequest](log)),
				orchestratorAddress:  "localhost:9092",
				cachedConfigurations: make(map[string]cachedConfiguration),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.args.opts...)

			// Check channels have been created
			assert.NotNil(t, s.orchestratorChannel)

			// Ignore pointers to channel in subsequent DeepEqual check
			s.orchestratorChannel = nil

			if got := s; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestAssessEvidence tests AssessEvidence
func TestAssessEvidence(t *testing.T) {
	type args struct {
		in0      context.Context
		evidence *evidence.Evidence
	}
	tests := []struct {
		name string
		args args
		// hasRPCConnection is true when connected to orchestrator and evidence store
		hasRPCConnection bool
		wantResp         *assessment.AssessEvidenceResponse
		wantErr          bool
	}{
		{
			name: "Assess resource without id",
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					ToolId:    "mock",
					Timestamp: timestamppb.Now(),
					Resource:  toStruct(voc.VirtualMachine{}, t),
				},
			},
			hasRPCConnection: true,
			wantResp: &assessment.AssessEvidenceResponse{
				Status:        assessment.AssessEvidenceResponse_FAILED,
				StatusMessage: "invalid evidence: " + evidence.ErrResourceIdFieldMissing.Error(),
			},
			wantErr: true,
		},
		{
			name: "Assess resource without tool id",
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					Timestamp: timestamppb.Now(),
					Resource:  toStruct(voc.VirtualMachine{}, t),
				},
			},
			hasRPCConnection: true,
			wantResp: &assessment.AssessEvidenceResponse{
				Status:        assessment.AssessEvidenceResponse_FAILED,
				StatusMessage: "invalid evidence: " + evidence.ErrResourceIdFieldMissing.Error(),
			},
			wantErr: true,
		},
		{
			name: "Assess resource without timestamp",
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					Id:       "11111111-1111-1111-1111-111111111111",
					ToolId:   "mock",
					Resource: toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: "my-resource-id", Type: []string{"VirtualMachine"}}}}, t),
				},
			},
			hasRPCConnection: true,
			wantResp: &assessment.AssessEvidenceResponse{
				Status:        assessment.AssessEvidenceResponse_FAILED,
				StatusMessage: "invalid evidence: " + evidence.ErrTimestampMissing.Error(),
			},
			wantErr: true,
		},
		{
			name: "Assess resource",
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					Id:        "11111111-1111-1111-1111-111111111111",
					ToolId:    "mock",
					Timestamp: timestamppb.Now(),
					Resource:  toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: "my-resource-id", Type: []string{"VirtualMachine"}}}}, t),
				},
			},
			hasRPCConnection: true,
			wantResp: &assessment.AssessEvidenceResponse{
				Status: assessment.AssessEvidenceResponse_ASSESSED,
			},
			wantErr: false,
		},
		{
			name: "No RPC connections",
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					Id:        "11111111-1111-1111-1111-111111111111",
					ToolId:    "mock",
					Timestamp: timestamppb.Now(),
					Resource:  toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: "my-resource-id", Type: []string{"VirtualMachine"}}}}, t),
				},
			},
			hasRPCConnection: false,
			wantResp: &assessment.AssessEvidenceResponse{
				Status:        assessment.AssessEvidenceResponse_FAILED,
				StatusMessage: "could not evaluate evidence: could not fetch metric configuration: could not retrieve metric configuration for",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			if tt.hasRPCConnection {
				s.grpcOpts = []grpc.DialOption{grpc.WithContextDialer(bufConnDialer)}
				assert.NoError(t, s.initOrchestratorStream(grpc.WithContextDialer(bufConnDialer)))
			} else {
				// clear the evidence URL, just to be sure
				s.evidenceStoreAddress = ""
				s.orchestratorAddress = ""
			}

			gotResp, err := s.AssessEvidence(tt.args.in0, &assessment.AssessEvidenceRequest{Evidence: tt.args.evidence})

			if (err != nil) != tt.wantErr {
				t.Errorf("AssessEvidence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check response
			assert.Equal(t, tt.wantResp.Status, gotResp.Status)
			assert.Contains(t, gotResp.StatusMessage, tt.wantResp.StatusMessage)
		})
	}
}

// TestAssessEvidences tests AssessEvidences
func TestAssessEvidences(t *testing.T) {
	type fields struct {
		ResultHooks                   []assessment.ResultHookFunc
		results                       map[string]*assessment.AssessmentResult
		UnimplementedAssessmentServer assessment.UnimplementedAssessmentServer
	}
	type args struct {
		streamToServer            *mockAssessmentServerStream
		streamToClientWithSendErr *mockAssessmentServerStreamWithSendErr
		streamToServerWithRecvErr *mockAssessmentServerStreamWithRecvErr
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantErr         bool
		wantErrMessage  string
		wantRespMessage *assessment.AssessEvidenceResponse
	}{
		{
			name: "Missing toolId",
			fields: fields{
				results: make(map[string]*assessment.AssessmentResult)},
			args: args{
				streamToServer: createMockAssessmentServerStream(&assessment.AssessEvidenceRequest{
					Evidence: &evidence.Evidence{
						Timestamp: timestamppb.Now(),
						Resource:  toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: "my-resource-id", Type: []string{"VirtualMachine"}}}}, t)}}),
			},
			wantErr: false,
			wantRespMessage: &assessment.AssessEvidenceResponse{
				Status:        assessment.AssessEvidenceResponse_FAILED,
				StatusMessage: "invalid evidence: " + evidence.ErrToolIdMissing.Error(),
			},
		},
		{
			name: "Missing evidenceID",
			fields: fields{
				results: make(map[string]*assessment.AssessmentResult)},
			args: args{
				streamToServer: createMockAssessmentServerStream(&assessment.AssessEvidenceRequest{
					Evidence: &evidence.Evidence{
						Timestamp: timestamppb.Now(),
						ToolId:    "2134",
						Resource:  toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: "my-resource-id", Type: []string{"VirtualMachine"}}}}, t)}}),
			},
			wantErr: false,
			wantRespMessage: &assessment.AssessEvidenceResponse{
				Status:        assessment.AssessEvidenceResponse_FAILED,
				StatusMessage: "invalid evidence: " + evidence.ErrEvidenceIdInvalidFormat.Error(),
			},
		},
		{
			name: "Assess evidences",
			fields: fields{
				results: make(map[string]*assessment.AssessmentResult)},
			args: args{
				streamToServer: createMockAssessmentServerStream(&assessment.AssessEvidenceRequest{
					Evidence: &evidence.Evidence{
						Timestamp: timestamppb.Now(),
						ToolId:    "2134",
						Resource:  toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: "my-resource-id", Type: []string{"VirtualMachine"}}}}, t)}}),
			},
			wantErr: false,
			wantRespMessage: &assessment.AssessEvidenceResponse{
				Status: assessment.AssessEvidenceResponse_ASSESSED,
			},
		},
		{
			name:   "Error in stream to client - Send()-err",
			fields: fields{},
			args: args{
				streamToClientWithSendErr: createMockAssessmentServerStreamWithSendErr(&assessment.AssessEvidenceRequest{
					Evidence: &evidence.Evidence{
						Timestamp: timestamppb.Now(),
						ToolId:    "2134",
						Resource:  toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: "my-resource-id", Type: []string{"VirtualMachine"}}}}, t)}}),
			},
			wantErr:        true,
			wantErrMessage: "rpc error: code = Unknown desc = cannot send response to the client",
		},
		{
			name:   "Error in stream to server - Recv()-err",
			fields: fields{},
			args: args{
				streamToServerWithRecvErr: createMockAssessmentServerStreamWithRecvErr(&assessment.AssessEvidenceRequest{
					Evidence: &evidence.Evidence{
						Timestamp: timestamppb.Now(),
						ToolId:    "2134",
						Resource:  toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: "my-resource-id", Type: []string{"VirtualMachine"}}}}, t)}}),
			},
			wantErr:        true,
			wantErrMessage: "rpc error: code = Unknown desc = cannot receive stream request",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				err                error
				responseFromServer *assessment.AssessEvidenceResponse
			)
			s := Service{
				resultHooks:                   tt.fields.ResultHooks,
				results:                       tt.fields.results,
				cachedConfigurations:          make(map[string]cachedConfiguration),
				UnimplementedAssessmentServer: tt.fields.UnimplementedAssessmentServer,
				orchestratorChannel:           make(chan *assessment.AssessmentResult, 1000),
				grpcOpts:                      []grpc.DialOption{grpc.WithContextDialer(bufConnDialer)},
			}

			assert.NoError(t, s.initOrchestratorStream(grpc.WithContextDialer(bufConnDialer)))

			if tt.args.streamToServer != nil {
				err = s.AssessEvidences(tt.args.streamToServer)
				responseFromServer = <-tt.args.streamToServer.SentFromServer
			} else if tt.args.streamToClientWithSendErr != nil {
				err = s.AssessEvidences(tt.args.streamToClientWithSendErr)
			} else if tt.args.streamToServerWithRecvErr != nil {
				err = s.AssessEvidences(tt.args.streamToServerWithRecvErr)
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

func TestAssessmentResultHooks(t *testing.T) {
	var (
		hookCallCounter = 0
		wg              sync.WaitGroup
		hookCounts      = 16
	)

	wg.Add(hookCounts)

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

	type args struct {
		in0         context.Context
		evidence    *assessment.AssessEvidenceRequest
		resultHooks []assessment.ResultHookFunc
	}
	tests := []struct {
		name     string
		args     args
		wantResp *assessment.AssessEvidenceResponse
		wantErr  bool
	}{
		{
			name: "Store evidence to the map",
			args: args{
				in0: context.TODO(),
				evidence: &assessment.AssessEvidenceRequest{
					Evidence: &evidence.Evidence{
						Id:        "11111111-1111-1111-1111-111111111111",
						ToolId:    "mock",
						Timestamp: timestamppb.Now(),
						Resource: toStruct(&voc.VirtualMachine{
							Compute: &voc.Compute{
								Resource: &voc.Resource{
									ID:   "my-resource-id",
									Type: []string{"VirtualMachine", "Compute", "Resource"}},
							},
							BootLogging: &voc.BootLogging{
								Logging: &voc.Logging{
									LoggingService:  []voc.ResourceID{"SomeResourceId2"},
									Enabled:         true,
									RetentionPeriod: 36,
								},
							},
							OSLogging: &voc.OSLogging{
								Logging: &voc.Logging{
									LoggingService:  []voc.ResourceID{"SomeResourceId2"},
									Enabled:         true,
									RetentionPeriod: 36,
								},
							},
							MalwareProtection: &voc.MalwareProtection{
								Enabled:              true,
								NumberOfThreatsFound: 5,
								DaysSinceActive:      20,
								ApplicationLogging: &voc.ApplicationLogging{
									Logging: &voc.Logging{
										Enabled:        true,
										LoggingService: []voc.ResourceID{"SomeAnalyticsService?"},
									},
								},
							},
						}, t),
					}},

				resultHooks: []assessment.ResultHookFunc{firstHookFunction, secondHookFunction},
			},
			wantErr:  false,
			wantResp: &assessment.AssessEvidenceResponse{Status: assessment.AssessEvidenceResponse_ASSESSED},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hookCallCounter = 0
			s := NewService(WithAdditionalGRPCOpts(grpc.WithContextDialer(bufConnDialer)))
			assert.NoError(t, s.initOrchestratorStream(grpc.WithContextDialer(bufConnDialer)))

			for i, hookFunction := range tt.args.resultHooks {
				s.RegisterAssessmentResultHook(hookFunction)

				// Check if hook is registered
				funcName1 := runtime.FuncForPC(reflect.ValueOf(s.resultHooks[i]).Pointer()).Name()
				funcName2 := runtime.FuncForPC(reflect.ValueOf(hookFunction).Pointer()).Name()
				assert.Equal(t, funcName1, funcName2)
			}

			// To test the hooks we have to call a function that calls the hook function
			gotResp, err := s.AssessEvidence(tt.args.in0, tt.args.evidence)

			// wait for all hooks (2 metrics * 2 hooks)
			wg.Wait()

			if (err != nil) != tt.wantErr {
				t.Errorf("AssessEvidence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("AssessEvidence() gotResp = %v, want %v", gotResp, tt.wantResp)
			}

			assert.Equal(t, tt.wantResp, gotResp)
			assert.NotEmpty(t, s.results)
			assert.Equal(t, hookCounts, hookCallCounter)
		})
	}
}

func TestListAssessmentResults(t *testing.T) {
	s := NewService(WithAdditionalGRPCOpts(grpc.WithContextDialer(bufConnDialer)))
	assert.NoError(t, s.initOrchestratorStream(grpc.WithContextDialer(bufConnDialer)))

	_, err := s.AssessEvidence(context.TODO(), &assessment.AssessEvidenceRequest{
		Evidence: &evidence.Evidence{
			Id:        "11111111-1111-1111-1111-111111111111",
			ToolId:    "mock",
			Timestamp: timestamppb.Now(),
			Resource:  toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: "my-resource-id", Type: []string{"VirtualMachine"}}}}, t),
		}})
	assert.NoError(t, err)
	var results *assessment.ListAssessmentResultsResponse
	results, err = s.ListAssessmentResults(context.TODO(), &assessment.ListAssessmentResultsRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, results)
}

// toStruct transforms r to a struct and asserts if it was successful
func toStruct(r voc.IsCloudResource, t *testing.T) (s *structpb.Value) {
	s, err := voc.ToStruct(r)
	if err != nil {
		assert.Error(t, err)
	}

	return
}

// mockAssessmentServerStream implements Assessment_AssessEvidencesServer which is used to mock incoming evidences as a stream
type mockAssessmentServerStream struct {
	grpc.ServerStream
	RecvToServer   chan *assessment.AssessEvidenceRequest
	SentFromServer chan *assessment.AssessEvidenceResponse
}

func (mockAssessmentServerStream) CloseSend() error {
	panic("implement me")
}

func createMockAssessmentServerStream(r *assessment.AssessEvidenceRequest) *mockAssessmentServerStream {
	m := &mockAssessmentServerStream{
		RecvToServer: make(chan *assessment.AssessEvidenceRequest, 1),
	}
	m.RecvToServer <- r

	m.SentFromServer = make(chan *assessment.AssessEvidenceResponse, 1)
	return m
}

func (m mockAssessmentServerStream) Send(response *assessment.AssessEvidenceResponse) error {
	m.SentFromServer <- response
	return nil
}

func (mockAssessmentServerStream) SendAndClose() error {
	return nil
}

// Stop, if no more evidences exist
// For now, just receive one evidence and directly stop the stream (EOF)
func (m *mockAssessmentServerStream) Recv() (req *assessment.AssessEvidenceRequest, err error) {
	if len(m.RecvToServer) == 0 {
		return nil, io.EOF
	}
	req, more := <-m.RecvToServer
	if !more {
		return nil, errors.New("empty")
	}

	return req, nil
}

func (mockAssessmentServerStream) SetHeader(metadata.MD) error {
	return nil
}

func (mockAssessmentServerStream) SendHeader(metadata.MD) error {
	return nil
}

func (mockAssessmentServerStream) SetTrailer(metadata.MD) {
}

func (mockAssessmentServerStream) Context() context.Context {
	return nil
}

func (mockAssessmentServerStream) SendMsg(interface{}) error {
	return nil
}

func (mockAssessmentServerStream) RecvMsg(interface{}) error {
	return nil
}

func createMockAssessmentServerStreamWithSendErr(r *assessment.AssessEvidenceRequest) *mockAssessmentServerStreamWithSendErr {
	m := &mockAssessmentServerStreamWithSendErr{
		RecvToServer: make(chan *assessment.AssessEvidenceRequest, 1),
	}
	m.RecvToServer <- r

	m.SentFromServer = make(chan *assessment.AssessEvidenceResponse, 1)
	return m
}

// mockAssessmentServerStreamWithSendErr implements Assessment_AssessEvidencesServer with error
type mockAssessmentServerStreamWithSendErr struct {
	grpc.ServerStream
	RecvToServer   chan *assessment.AssessEvidenceRequest
	SentFromServer chan *assessment.AssessEvidenceResponse
}

func (*mockAssessmentServerStreamWithSendErr) Send(*assessment.AssessEvidenceResponse) error {
	return errors.New("error sending response to client")
}

// Stop, if no more evidences exist
// For now, just receive one evidence and directly stop the stream (EOF)
func (m *mockAssessmentServerStreamWithSendErr) Recv() (req *assessment.AssessEvidenceRequest, err error) {
	if len(m.RecvToServer) == 0 {
		return nil, io.EOF
	}
	req, more := <-m.RecvToServer
	if !more {
		return nil, errors.New("empty")
	}

	return req, nil
}

type mockAssessmentServerStreamWithRecvErr struct {
	grpc.ServerStream
	RecvToServer   chan *assessment.AssessEvidenceRequest
	SentFromServer chan *assessment.AssessEvidenceResponse
}

func (mockAssessmentServerStreamWithRecvErr) Send(*assessment.AssessEvidenceResponse) error {
	panic("implement me")
}

func (mockAssessmentServerStreamWithRecvErr) Recv() (*assessment.AssessEvidenceRequest, error) {
	err := errors.New("Recv()-error")

	return nil, err
}

func createMockAssessmentServerStreamWithRecvErr(r *assessment.AssessEvidenceRequest) *mockAssessmentServerStreamWithRecvErr {
	m := &mockAssessmentServerStreamWithRecvErr{
		RecvToServer: make(chan *assessment.AssessEvidenceRequest, 1),
	}
	m.RecvToServer <- r

	m.SentFromServer = make(chan *assessment.AssessEvidenceResponse, 1)
	return m
}

func TestConvertTargetValue(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name                     string
		args                     args
		wantConvertedTargetValue *structpb.Value
		wantErr                  assert.ErrorAssertionFunc
	}{
		{
			name:                     "string",
			args:                     args{value: "TLS1.3"},
			wantConvertedTargetValue: &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: "TLS1.3"}},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name:                     "bool",
			args:                     args{value: false},
			wantConvertedTargetValue: &structpb.Value{Kind: &structpb.Value_BoolValue{BoolValue: false}},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name:                     "jsonNumber",
			args:                     args{value: json.Number("4")},
			wantConvertedTargetValue: &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: 4.}},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name:                     "int",
			args:                     args{value: 4},
			wantConvertedTargetValue: &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: 4.}},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name:                     "float64",
			args:                     args{value: 4.},
			wantConvertedTargetValue: &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: 4.}},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name:                     "float32",
			args:                     args{value: float32(4.)},
			wantConvertedTargetValue: &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: 4.}},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "list of strings",
			args: args{value: []string{"TLS1.2", "TLS1.3"}},
			wantConvertedTargetValue: &structpb.Value{Kind: &structpb.Value_ListValue{ListValue: &structpb.ListValue{Values: []*structpb.Value{
				{Kind: &structpb.Value_StringValue{StringValue: "TLS1.2"}},
				{Kind: &structpb.Value_StringValue{StringValue: "TLS1.3"}},
			}}}},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConvertedTargetValue, err := convertTargetValue(tt.args.value)
			if !tt.wantErr(t, err, fmt.Sprintf("convertTargetValue(%v)", tt.args.value)) {
				return
			}
			// Checking against 'String()' allows to compare the actual values instead of the respective pointers
			assert.Equalf(t, tt.wantConvertedTargetValue.String(), gotConvertedTargetValue.String(), "convertTargetValue(%v)", tt.args.value)
		})
	}
}

func TestHandleEvidence(t *testing.T) {
	type fields struct {
		hasEvidenceStoreStream bool
		hasOrchestratorStream  bool
	}
	type args struct {
		evidence   *evidence.Evidence
		resourceId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "correct evidence",
			fields: fields{
				hasOrchestratorStream:  true,
				hasEvidenceStoreStream: true,
			},
			args: args{
				evidence: &evidence.Evidence{
					Id:        "11111111-1111-1111-1111-111111111111",
					ToolId:    "mock",
					Timestamp: timestamppb.Now(),
					Resource:  toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: "my-resource-id", Type: []string{"VirtualMachine"}}}}, t),
				},
				resourceId: "my-resource-id",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NoError(t, err)
				return false
			},
		},
		{
			name: "missing type in evidence",
			fields: fields{
				hasOrchestratorStream:  true,
				hasEvidenceStoreStream: true,
			},
			args: args{
				evidence: &evidence.Evidence{
					Id:        "11111111-1111-1111-1111-111111111111",
					ToolId:    "mock",
					Timestamp: timestamppb.Now(),
					Resource:  toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: "my-resource-id", Type: []string{}}}}, t),
				},
				resourceId: "my-resource-id",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Error(t, err)
				// Check if error message contains "empty" (list of types)
				assert.Contains(t, err.Error(), "empty")
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			// Mock streams for target services if needed
			if tt.fields.hasEvidenceStoreStream {
				s.grpcOpts = []grpc.DialOption{grpc.WithContextDialer(bufConnDialer)}
			}
			if tt.fields.hasOrchestratorStream {
				assert.NoError(t, s.initOrchestratorStream(grpc.WithContextDialer(bufConnDialer)))
			}
			// Two tests: 1st) wantErr function. 2nd) if wantErr false then check if a result is added to map
			if !tt.wantErr(t, s.handleEvidence(tt.args.evidence, tt.args.resourceId), fmt.Sprintf("handleEvidence(%v, %v)", tt.args.evidence, tt.args.resourceId)) {
				assert.NotEmpty(t, s.results)
			}

		})
	}
}

func TestService_initOrchestratorStoreStream(t *testing.T) {
	type fields struct {
		opts []ServiceOption
	}
	type args struct {
		additionalOpts []grpc.DialOption
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Invalid RPC connection",
			fields: fields{
				opts: []ServiceOption{
					WithOrchestratorAddress("localhost:1"),
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				s, _ := status.FromError(errors.Unwrap(err))
				return assert.Equal(t, codes.Unavailable, s.Code())
			},
		},
		{
			name: "Authenticated RPC connection with valid user",
			fields: fields{
				opts: []ServiceOption{
					WithOrchestratorAddress("bufnet"),
					WithOAuth2Authorizer(testutil.AuthClientConfig(authPort)),
				},
			},
			args: args{
				[]grpc.DialOption{grpc.WithContextDialer(bufConnDialer)},
			},
		},
		{
			name: "Authenticated RPC connection with invalid user",
			fields: fields{
				opts: []ServiceOption{
					WithOrchestratorAddress("bufnet"),
					WithOAuth2Authorizer(testutil.AuthClientConfig(authPort)),
				},
			},
			args: args{
				[]grpc.DialOption{grpc.WithContextDialer(bufConnDialer)},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				s, _ := status.FromError(errors.Unwrap(err))
				return assert.Equal(t, codes.Unauthenticated, s.Code())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.fields.opts...)
			err := s.initOrchestratorStream(tt.args.additionalOpts...)

			if tt.wantErr != nil {
				tt.wantErr(t, err)
			}
		})
	}
}
