// Copyright 2021 Fraunhofer AISEC
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
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidence"
	service_evidenceStore "clouditor.io/clouditor/service/evidence"
	"clouditor.io/clouditor/voc"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"net"
	"os"
	"reflect"
	"runtime"
	"testing"
	"time"
)

var lis *bufconn.Listener

func TestMain(m *testing.M) {
	// pre-configuration for mocking evidence store
	const bufSize = 1024 * 1024
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	evidence.RegisterEvidenceStoreServer(s, service_evidenceStore.NewService())
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	// make sure, that we are in the clouditor root folder to find the policies
	err := os.Chdir("../../")
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

// TestNewService is a simply test for NewService
func TestNewService(t *testing.T) {
	tests := []struct {
		name string
		want assessment.AssessmentServer
	}{
		{
			name: "AssessmentServer created with empty results map",
			want: &Service{
				results:                       make(map[string]*assessment.AssessmentResult),
				UnimplementedAssessmentServer: assessment.UnimplementedAssessmentServer{},
				Configuration:                 Configuration{evidenceStoreTargetAddress: "localhost:9090"},
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

func TestStart(t *testing.T) {
	type fields struct {
		UnimplementedAssessmentServer assessment.UnimplementedAssessmentServer
		evidenceStoreStream           evidence.EvidenceStore_StoreEvidencesClient
		ResultHook                    []assessment.ResultHookFunc
		results                       map[string]*assessment.AssessmentResult
		Configuration
	}
	type args struct {
		in0 context.Context
		in1 *assessment.StartAssessmentRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *assessment.StartAssessmentResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "dial-error",
			fields: fields{
				UnimplementedAssessmentServer: assessment.UnimplementedAssessmentServer{},
				// Setting invalid target causes grpc.dial to fail
				Configuration: Configuration{evidenceStoreTargetAddress: "\n"},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NotNil(t, err)
			},
		},
		{
			name: "evidenceStoreStream-error",
			fields: fields{
				UnimplementedAssessmentServer: assessment.UnimplementedAssessmentServer{},
				// Target is valid but causes `StoreEvidences` to fail with "missing address"
				Configuration: Configuration{evidenceStoreTargetAddress: ""},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NotNil(t, err)
			},
		},
		// TODO(lebogg): Test valid case. Have to mock evidenceStore component somehow
		//{
		//	name: "valid",
		//	fields: fields{
		//		UnimplementedAssessmentServer: assessment.UnimplementedAssessmentServer{},
		//		Configuration:                 Configuration{evidenceStoreTargetAddress: "localhost:9090"},
		//	},
		//	wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
		//		fmt.Print(err)
		//		return assert.Nil(t, err)
		//	},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				UnimplementedAssessmentServer: tt.fields.UnimplementedAssessmentServer,
				evidenceStoreStream:           tt.fields.evidenceStoreStream,
				resultHooks:                   tt.fields.ResultHook,
				results:                       tt.fields.results,
				Configuration:                 tt.fields.Configuration,
			}
			got, err := s.Start(tt.args.in0, tt.args.in1)
			if !tt.wantErr(t, err, fmt.Sprintf("Start(%v, %v)", tt.args.in0, tt.args.in1)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Start(%v, %v)", tt.args.in0, tt.args.in1)
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
		name     string
		args     args
		wantResp *assessment.AssessEvidenceResponse
		wantErr  bool
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
					ToolId:    "mock",
					Timestamp: timestamppb.Now(),
					Resource:  toStruct(voc.VirtualMachine{Compute: &voc.Compute{CloudResource: &voc.CloudResource{ID: "my-resource-id", Type: []string{"VirtualMachine"}}}}, t),
				},
			},
			wantResp: &assessment.AssessEvidenceResponse{
				Status: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			assert.NoError(t, s.mockEvidenceStream())

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

func TestService_AssessEvidences(t *testing.T) {
	type fields struct {
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
			name:   "Assessing evidences fails due to missing toolId",
			fields: fields{results: make(map[string]*assessment.AssessmentResult)},
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
			name:   "Assess evidences",
			fields: fields{results: make(map[string]*assessment.AssessmentResult)},
			args: args{stream: &mockAssessmentStream{
				evidence: &evidence.Evidence{
					ToolId:    "mock",
					Timestamp: timestamppb.Now(),
					Resource:  toStruct(voc.VirtualMachine{Compute: &voc.Compute{CloudResource: &voc.CloudResource{ID: "my-resource-id", Type: []string{"VirtualMachine"}}}}, t),
				},
			}},
			wantErr:        false,
			wantErrMessage: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				resultHooks:                   tt.fields.ResultHooks,
				results:                       tt.fields.results,
				UnimplementedAssessmentServer: tt.fields.UnimplementedAssessmentServer,
			}
			assert.NoError(t, s.mockEvidenceStream())

			err := s.AssessEvidences(tt.args.stream)
			if (err != nil) != tt.wantErr {
				t.Errorf("AssessEvidence() error = %v, wantErr %v", err, tt.wantErr)
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
	)

	firstHookFunction := func(assessmentResult *assessment.AssessmentResult, err error) {
		hookCallCounter++
		log.Println("Hello from inside the firstHookFunction")
	}

	secondHookFunction := func(assessmentResult *assessment.AssessmentResult, err error) {
		hookCallCounter++
		log.Println("Hello from inside the secondHookFunction")
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
						Id:        "MockEvidenceID",
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
			assert.NoError(t, s.mockEvidenceStream())

			for i, hookFunction := range tt.args.resultHooks {
				s.RegisterAssessmentResultHook(hookFunction)

				// Check if hook is registered
				funcName1 := runtime.FuncForPC(reflect.ValueOf(s.resultHooks[i]).Pointer()).Name()
				funcName2 := runtime.FuncForPC(reflect.ValueOf(hookFunction).Pointer()).Name()
				assert.Equal(t, funcName1, funcName2)
			}

			gotResp, err := s.AssessEvidence(tt.args.in0, tt.args.evidence)

			// That isn´t nice, but we have somehow to wait for the hook functions
			time.Sleep(3 * time.Second)

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
	assert.NoError(t, s.mockEvidenceStream())
	_, err := s.AssessEvidence(context.TODO(), &assessment.AssessEvidenceRequest{
		Evidence: &evidence.Evidence{
			ToolId:    "mock",
			Timestamp: timestamppb.Now(),
			Resource:  toStruct(voc.VirtualMachine{Compute: &voc.Compute{CloudResource: &voc.CloudResource{ID: "my-resource-id", Type: []string{"VirtualMachine"}}}}, t),
		}})
	assert.Nil(t, err)
	var results *assessment.ListAssessmentResultsResponse
	results, err = s.ListAssessmentResults(context.TODO(), &assessment.ListAssessmentResultsRequest{})
	assert.Nil(t, err)
	assert.NotNil(t, results)
}

// toStruct transforms r to a struct and asserts if it was successful
func toStruct(r voc.IsCloudResource, t *testing.T) (s *structpb.Value) {
	s, err := voc.ToStruct(r)
	if err != nil {
		assert.NotNil(t, err)
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

// Mocking evidence store service

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func (s *Service) mockEvidenceStream() error {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		return err
	}
	// defer conn.Close()
	client, err := evidence.NewEvidenceStoreClient(conn).StoreEvidences(ctx)
	if err != nil {
		return err
	}
	s.evidenceStoreStream = client
	return nil
}
