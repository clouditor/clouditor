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

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/voc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestMain(m *testing.M) {
	// make sure, that we are in the clouditor root folder to find the policies
	err := os.Chdir("../../")
	if err != nil {
		panic(err)
	}

	server, authService, _, _ := startBufConnServer()
	err = authService.CreateDefaultUser("clouditor", "clouditor")
	if err != nil {
		panic(err)
	}

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
				orchestratorAddress:  "localhost:9090",
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
				orchestratorAddress:  "localhost:9092",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewService(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
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
				Status: false,
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
				Status: false,
			},
			wantErr: true,
		},
		{
			name: "Assess resource without timestamp",
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					ToolId:   "mock",
					Resource: toStruct(voc.VirtualMachine{}, t),
				},
			},
			hasRPCConnection: true,
			wantResp: &assessment.AssessEvidenceResponse{
				Status: false,
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
					Resource:  toStruct(voc.VirtualMachine{Compute: &voc.Compute{CloudResource: &voc.CloudResource{ID: "my-resource-id", Type: []string{"VirtualMachine"}}}}, t),
				},
			},
			hasRPCConnection: true,
			wantResp: &assessment.AssessEvidenceResponse{
				Status: true,
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
					Resource:  toStruct(voc.VirtualMachine{Compute: &voc.Compute{CloudResource: &voc.CloudResource{ID: "my-resource-id", Type: []string{"VirtualMachine"}}}}, t),
				},
			},
			hasRPCConnection: false,
			wantResp:         &assessment.AssessEvidenceResponse{Status: false},
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			if tt.hasRPCConnection {
				assert.NoError(t, s.initEvidenceStoreStream(grpc.WithContextDialer(bufConnDialer)))
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

			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("AssessEvidence() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func TestAssessEvidences(t *testing.T) {
	type fields struct {
		hasRPCConnection              bool
		ResultHooks                   []assessment.ResultHookFunc
		results                       map[string]*assessment.AssessmentResult
		UnimplementedAssessmentServer assessment.UnimplementedAssessmentServer
	}
	type args struct {
		stream assessment.Assessment_AssessEvidencesServer
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		wantErrMessage string
	}{
		{
			name: "Missing toolId",
			fields: fields{
				hasRPCConnection: true,
				results:          make(map[string]*assessment.AssessmentResult)},
			args: args{stream: &mockAssessmentStream{
				evidence: &evidence.Evidence{
					Timestamp: timestamppb.Now(),
					Resource:  toStruct(voc.VirtualMachine{Compute: &voc.Compute{CloudResource: &voc.CloudResource{ID: "my-resource-id", Type: []string{"VirtualMachine"}}}}, t),
				},
			}},
			wantErr:        true,
			wantErrMessage: "invalid evidence",
		},
		{
			name: "Assess evidences",
			fields: fields{
				hasRPCConnection: true,
				results:          make(map[string]*assessment.AssessmentResult)},
			args: args{stream: &mockAssessmentStream{
				evidence: &evidence.Evidence{
					Id:        "11111111-1111-1111-1111-111111111111",
					ToolId:    "mock",
					Timestamp: timestamppb.Now(),
					Resource:  toStruct(voc.VirtualMachine{Compute: &voc.Compute{CloudResource: &voc.CloudResource{ID: "my-resource-id", Type: []string{"VirtualMachine"}}}}, t),
				},
			}},
			wantErr:        false,
			wantErrMessage: "",
		},
		{
			name: "No RPC connections",
			fields: fields{
				hasRPCConnection: false,
			},
			args:           args{stream: &mockAssessmentStreamWithRecvErr{}},
			wantErr:        true,
			wantErrMessage: codes.Internal.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				resultHooks:                   tt.fields.ResultHooks,
				results:                       tt.fields.results,
				UnimplementedAssessmentServer: tt.fields.UnimplementedAssessmentServer,
			}
			if tt.fields.hasRPCConnection {
				assert.NoError(t, s.initEvidenceStoreStream(grpc.WithContextDialer(bufConnDialer)))
				assert.NoError(t, s.initOrchestratorStream(grpc.WithContextDialer(bufConnDialer)))
			}

			err := s.AssessEvidences(tt.args.stream)
			fmt.Println(err)
			if (err != nil) != tt.wantErr {
				t.Errorf("Got AssessEvidence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				assert.Contains(t, err.Error(), tt.wantErrMessage)
			}
		})
	}
}

func TestAssessmentResultHooks(t *testing.T) {
	var (
		hookCallCounter = 0
		wg              sync.WaitGroup
	)

	wg.Add(12)

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

	// Check GRPC call
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
						Resource: toStruct(voc.VirtualMachine{
							Compute: &voc.Compute{
								CloudResource: &voc.CloudResource{
									ID:   "my-resource-id",
									Type: []string{"VirtualMachine"}},
							},
						}, t),
					}},

				resultHooks: []assessment.ResultHookFunc{firstHookFunction, secondHookFunction},
			},
			wantErr:  false,
			wantResp: &assessment.AssessEvidenceResponse{Status: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hookCallCounter = 0
			s := NewService()
			assert.NoError(t, s.initEvidenceStoreStream(grpc.WithContextDialer(bufConnDialer)))
			assert.NoError(t, s.initOrchestratorStream(grpc.WithContextDialer(bufConnDialer)))

			for i, hookFunction := range tt.args.resultHooks {
				s.RegisterAssessmentResultHook(hookFunction)

				// Check if hook is registered
				funcName1 := runtime.FuncForPC(reflect.ValueOf(s.resultHooks[i]).Pointer()).Name()
				funcName2 := runtime.FuncForPC(reflect.ValueOf(hookFunction).Pointer()).Name()
				assert.Equal(t, funcName1, funcName2)
			}

			gotResp, err := s.AssessEvidence(tt.args.in0, tt.args.evidence)

			// wait for all hooks (6 metrics * 2 hooks)
			wg.Wait()

			if (err != nil) != tt.wantErr {
				t.Errorf("AssessEvidence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("AssessEvidence() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
			assert.NotEmpty(t, s.results)
			assert.Equal(t, 12, hookCallCounter)
		})
	}
}

func TestListAssessmentResults(t *testing.T) {
	s := NewService()
	assert.NoError(t, s.initEvidenceStoreStream(grpc.WithContextDialer(bufConnDialer)))
	assert.NoError(t, s.initOrchestratorStream(grpc.WithContextDialer(bufConnDialer)))
	_, err := s.AssessEvidence(context.TODO(), &assessment.AssessEvidenceRequest{
		Evidence: &evidence.Evidence{
			Id:        "11111111-1111-1111-1111-111111111111",
			ToolId:    "mock",
			Timestamp: timestamppb.Now(),
			Resource:  toStruct(voc.VirtualMachine{Compute: &voc.Compute{CloudResource: &voc.CloudResource{ID: "my-resource-id", Type: []string{"VirtualMachine"}}}}, t),
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

// mockAssessmentStream implements Assessment_AssessEvidencesServer which is used to mock incoming evidences as a stream
type mockAssessmentStream struct {
	evidence         *evidence.Evidence
	receivedEvidence bool
}

func (mockAssessmentStream) SendAndClose(*emptypb.Empty) error {
	return nil
}

// For now, just receive one evidence and directly stop the stream (EOF)
func (m *mockAssessmentStream) Recv() (req *assessment.AssessEvidenceRequest, err error) {
	if !m.receivedEvidence {
		req = new(assessment.AssessEvidenceRequest)
		req.Evidence = m.evidence
		m.receivedEvidence = true
	} else {
		err = io.EOF
	}
	return
}

func (mockAssessmentStream) SetHeader(metadata.MD) error {
	return nil
}

func (mockAssessmentStream) SendHeader(metadata.MD) error {
	return nil
}

func (mockAssessmentStream) SetTrailer(metadata.MD) {
}

func (mockAssessmentStream) Context() context.Context {
	return nil
}

func (mockAssessmentStream) SendMsg(interface{}) error {
	return nil
}

func (mockAssessmentStream) RecvMsg(interface{}) error {
	return nil
}

// mockAssessmentStream implements Assessment_AssessEvidencesServer which directly throws error on Recv
type mockAssessmentStreamWithRecvErr struct {
}

func (mockAssessmentStreamWithRecvErr) SendAndClose(*emptypb.Empty) error {
	return nil
}

func (mockAssessmentStreamWithRecvErr) Recv() (*assessment.AssessEvidenceRequest, error) {
	return nil, status.Errorf(codes.Internal, "receiving internal error")
}

func (mockAssessmentStreamWithRecvErr) SetHeader(metadata.MD) error {
	return nil
}

func (mockAssessmentStreamWithRecvErr) SendHeader(metadata.MD) error {
	return nil
}

func (mockAssessmentStreamWithRecvErr) SetTrailer(metadata.MD) {
}

func (mockAssessmentStreamWithRecvErr) Context() context.Context {
	return nil
}

func (mockAssessmentStreamWithRecvErr) SendMsg(interface{}) error {
	return nil
}

func (mockAssessmentStreamWithRecvErr) RecvMsg(interface{}) error {
	return nil
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
		//results                       map[string]*assessment.AssessmentResult
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
					Resource:  toStruct(voc.VirtualMachine{Compute: &voc.Compute{CloudResource: &voc.CloudResource{ID: "my-resource-id", Type: []string{"VirtualMachine"}}}}, t),
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
					Resource:  toStruct(voc.VirtualMachine{Compute: &voc.Compute{CloudResource: &voc.CloudResource{ID: "my-resource-id", Type: []string{}}}}, t),
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
		{
			name: "no connection to Evidence Store",
			fields: fields{
				hasOrchestratorStream:  true,
				hasEvidenceStoreStream: false,
			},
			args: args{
				evidence: &evidence.Evidence{
					Id:        "11111111-1111-1111-1111-111111111111",
					ToolId:    "mock",
					Timestamp: timestamppb.Now(),
					Resource:  toStruct(voc.VirtualMachine{Compute: &voc.Compute{CloudResource: &voc.CloudResource{ID: "my-resource-id", Type: []string{"VirtualMachine"}}}}, t),
				},
				resourceId: "my-resource-id",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Error(t, err)
				// Check if error message contains "empty" (list of types)
				assert.Contains(t, err.Error(), "Evidence Store")
				assert.Contains(t, err.Error(), "Unavailable desc")
				return true
			},
		},
		{
			name: "no connection to Orchestrator",
			fields: fields{
				hasOrchestratorStream:  false,
				hasEvidenceStoreStream: true,
			},
			args: args{
				evidence: &evidence.Evidence{
					Id:        "11111111-1111-1111-1111-111111111111",
					ToolId:    "mock",
					Timestamp: timestamppb.Now(),
					Resource:  toStruct(voc.VirtualMachine{Compute: &voc.Compute{CloudResource: &voc.CloudResource{ID: "my-resource-id", Type: []string{"VirtualMachine"}}}}, t),
				},
				resourceId: "my-resource-id",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Error(t, err)
				// Check if error message contains "empty" (list of types)
				assert.Contains(t, err.Error(), "Orchestrator")
				assert.Contains(t, err.Error(), "Unavailable desc")
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			// Mock streams for target services if needed
			if tt.fields.hasEvidenceStoreStream {
				assert.NoError(t, s.initEvidenceStoreStream(grpc.WithContextDialer(bufConnDialer)))
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

func TestService_initEvidenceStoreStream(t *testing.T) {
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
					WithEvidenceStoreAddress("localhost:1"),
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				// We are looking for a connection refused message
				innerErr := errors.Unwrap(err)
				s, _ := status.FromError(innerErr)

				if s.Code() != codes.Unavailable {
					tt.Errorf("Status should be codes.Unavailable: %v", s.Code())
					return false
				}

				return true
			},
		},
		{
			name: "Authenticated RPC connection with valid user",
			fields: fields{
				opts: []ServiceOption{
					WithEvidenceStoreAddress("bufnet"),
					WithInternalAuthorizer(
						"bufnet",
						"clouditor",
						"clouditor",
						grpc.WithContextDialer(bufConnDialer),
					),
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
					WithEvidenceStoreAddress("bufnet"),
					WithInternalAuthorizer(
						"bufnet",
						"not_clouditor",
						"clouditor",
						grpc.WithContextDialer(bufConnDialer),
					),
				},
			},
			args: args{
				[]grpc.DialOption{grpc.WithContextDialer(bufConnDialer)},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(tt, err.Error(), "could not set up stream for storing evidences: rpc error: code = Unauthenticated desc = transport: per-RPC creds failed due to error: error while logging in: rpc error: code = Unauthenticated desc = login failed")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.fields.opts...)
			err := s.initEvidenceStoreStream(tt.args.additionalOpts...)

			if tt.wantErr != nil {
				tt.wantErr(t, err)
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
				// We are looking for a connection refused message
				innerErr := errors.Unwrap(err)
				s, _ := status.FromError(innerErr)

				if s.Code() != codes.Unavailable {
					tt.Errorf("Status should be codes.Unavailable: %v", s.Code())
					return false
				}

				return true
			},
		},
		{
			name: "Authenticated RPC connection with valid user",
			fields: fields{
				opts: []ServiceOption{
					WithOrchestratorAddress("bufnet"),
					WithInternalAuthorizer(
						"bufnet",
						"clouditor",
						"clouditor",
						grpc.WithContextDialer(bufConnDialer),
					),
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
					WithInternalAuthorizer(
						"bufnet",
						"not_clouditor",
						"clouditor",
						grpc.WithContextDialer(bufConnDialer),
					),
				},
			},
			args: args{
				[]grpc.DialOption{grpc.WithContextDialer(bufConnDialer)},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(tt, err.Error(), "could not set up stream for storing assessment results: rpc error: code = Unauthenticated desc = transport: per-RPC creds failed due to error: error while logging in: rpc error: code = Unauthenticated desc = login failed")
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
