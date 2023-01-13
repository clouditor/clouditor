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
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/internal/testutil/clitest"
	"clouditor.io/clouditor/policies"
	"clouditor.io/clouditor/service"
	"clouditor.io/clouditor/voc"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	authPort uint16
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
		opts []service.Option[Service]
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{
			name: "AssessmentServer created with option rego package name",
			args: args{
				opts: []service.Option[Service]{
					WithRegoPackageName("testPkg"),
				},
			},
			want: &Service{
				results: make(map[string]*assessment.AssessmentResult),
				evidenceStoreAddress: grpcTarget{
					target: "localhost:9090",
				},
				orchestratorAddress: grpcTarget{
					target: "localhost:9090",
				},
				evidenceStoreStreams: nil,
				orchestratorStreams:  nil,
				cachedConfigurations: make(map[string]cachedConfiguration),
				evalPkg:              "testPkg",
			},
		},
		{
			name: "AssessmentServer created with option authorizer",
			args: args{
				opts: []service.Option[Service]{
					WithAuthorizer(api.NewOAuthAuthorizerFromClientCredentials(&clientcredentials.Config{})),
				},
			},
			want: &Service{
				results: make(map[string]*assessment.AssessmentResult),
				evidenceStoreAddress: grpcTarget{
					target: "localhost:9090",
				},
				orchestratorAddress: grpcTarget{
					target: "localhost:9090",
				},
				evidenceStoreStreams: nil,
				orchestratorStreams:  nil,
				cachedConfigurations: make(map[string]cachedConfiguration),
				evalPkg:              policies.DefaultRegoPackage,
				authorizer:           api.NewOAuthAuthorizerFromClientCredentials(&clientcredentials.Config{}),
			},
		},
		{
			name: "AssessmentServer created with option authorizer",
			args: args{
				opts: []service.Option[Service]{
					WithOAuth2Authorizer(&clientcredentials.Config{}),
				},
			},
			want: &Service{
				results: make(map[string]*assessment.AssessmentResult),
				evidenceStoreAddress: grpcTarget{
					target: "localhost:9090",
				},
				orchestratorAddress: grpcTarget{
					target: "localhost:9090",
				},
				evidenceStoreStreams: nil,
				orchestratorStreams:  nil,
				cachedConfigurations: make(map[string]cachedConfiguration),
				evalPkg:              policies.DefaultRegoPackage,
				authorizer:           api.NewOAuthAuthorizerFromClientCredentials(&clientcredentials.Config{}),
			},
		},
		{
			name: "AssessmentServer created with empty results map",
			want: &Service{
				results: make(map[string]*assessment.AssessmentResult),
				evidenceStoreAddress: grpcTarget{
					target: "localhost:9090",
				},
				orchestratorAddress: grpcTarget{
					target: "localhost:9090",
				},
				evidenceStoreStreams: nil,
				orchestratorStreams:  nil,
				cachedConfigurations: make(map[string]cachedConfiguration),
				evalPkg:              policies.DefaultRegoPackage,
			},
		},
		{
			name: "AssessmentServer created with options",
			args: args{
				opts: []service.Option[Service]{
					WithEvidenceStoreAddress("localhost:9091"),
					WithOrchestratorAddress("localhost:9092"),
				},
			},
			want: &Service{
				results: make(map[string]*assessment.AssessmentResult),
				evidenceStoreAddress: grpcTarget{
					target: "localhost:9091",
				},
				orchestratorAddress: grpcTarget{
					target: "localhost:9092",
				},
				evidenceStoreStreams: nil,
				orchestratorStreams:  nil,
				cachedConfigurations: make(map[string]cachedConfiguration),
				evalPkg:              policies.DefaultRegoPackage,
			},
		},
		{
			name: "AssessmentServer without EvidenceStore",
			args: args{
				opts: []service.Option[Service]{
					WithoutEvidenceStore(),
				},
			},
			want: &Service{
				results:                 make(map[string]*assessment.AssessmentResult),
				isEvidenceStoreDisabled: true,
				evidenceStoreAddress: grpcTarget{
					target: DefaultEvidenceStoreAddress,
				},
				orchestratorAddress: grpcTarget{
					target: DefaultOrchestratorAddress,
				},
				orchestratorStreams:  nil,
				cachedConfigurations: make(map[string]cachedConfiguration),
				evalPkg:              policies.DefaultRegoPackage,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.args.opts...)

			// Ignore pointers to storage and policy eval
			s.pe = nil

			// Check if stream are not nil and ignore them for the following deepEqual
			assert.NotNil(t, s.evidenceStoreStreams)
			assert.NotNil(t, s.orchestratorStreams)
			s.evidenceStoreStreams = nil
			s.orchestratorStreams = nil

			if got := s; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestAssessEvidence tests AssessEvidence
func TestService_AssessEvidence(t *testing.T) {
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
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "Missing evidence",
			args: args{
				in0: context.TODO(),
			},
			hasRPCConnection: false,
			wantResp:         nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid request: invalid AssessEvidenceRequest.Evidence: value is required")
			},
		},
		{
			name: "Empty evidence",
			args: args{
				in0:      context.TODO(),
				evidence: &evidence.Evidence{},
			},
			hasRPCConnection: false,
			wantResp:         nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid request: invalid AssessEvidenceRequest.Evidence: embedded message failed validation | caused by: invalid Evidence.Id: value must be a valid UUID | caused by: invalid uuid format")
			},
		},
		{
			name: "Assess resource without id",
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					ToolId:    testdata.MockEvidenceToolID,
					Timestamp: timestamppb.Now(),
					Resource:  toStruct(voc.VirtualMachine{}, t),
				},
			},
			hasRPCConnection: true,
			wantResp:         nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid request: invalid AssessEvidenceRequest.Evidence: embedded message failed validation | caused by: invalid Evidence.Id: value must be a valid UUID | caused by: invalid uuid format")
			},
		},
		{
			name: "Assess resource without tool id",
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					Id:             testdata.MockEvidenceID,
					Timestamp:      timestamppb.Now(),
					CloudServiceId: testdata.MockCloudServiceID,
					Resource:       toStruct(voc.VirtualMachine{}, t),
				},
			},
			hasRPCConnection: true,
			wantResp:         nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid request: invalid AssessEvidenceRequest.Evidence: embedded message failed validation | caused by: invalid Evidence.ToolId: value length must be at least 1 runes")
			},
		},
		{
			name: "Assess resource without timestamp",
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					Id:             testdata.MockEvidenceID,
					ToolId:         testdata.MockEvidenceToolID,
					CloudServiceId: testdata.MockCloudServiceID,
					Resource:       toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: testdata.MockResourceID, Type: []string{"VirtualMachine"}}}}, t),
				},
			},
			hasRPCConnection: true,
			wantResp:         nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid request: invalid AssessEvidenceRequest.Evidence: embedded message failed validation | caused by: invalid Evidence.Timestamp: value is required")
			},
		},
		{
			name: "Assess resource",
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					Id:             testdata.MockEvidenceID,
					ToolId:         testdata.MockEvidenceToolID,
					Timestamp:      timestamppb.Now(),
					Resource:       toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: testdata.MockResourceID, Type: []string{"VirtualMachine"}}}}, t),
					CloudServiceId: testdata.MockCloudServiceID},
			},
			hasRPCConnection: true,
			wantResp:         &assessment.AssessEvidenceResponse{},
			wantErr:          assert.NoError,
		},
		{
			name: "Assess resource without resource id",
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					Id:             testdata.MockEvidenceID,
					ToolId:         testdata.MockEvidenceToolID,
					Timestamp:      timestamppb.Now(),
					Resource:       toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{Type: []string{"VirtualMachine"}}}}, t),
					CloudServiceId: testdata.MockCloudServiceID,
				},
			},
			hasRPCConnection: true,
			wantResp:         nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid evidence: resource in evidence is missing the id field")
			},
		},
		{
			name: "No RPC connections",
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					Id:             testdata.MockEvidenceID,
					ToolId:         testdata.MockEvidenceToolID,
					Timestamp:      timestamppb.Now(),
					CloudServiceId: testdata.MockCloudServiceID,
					Resource:       toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: testdata.MockResourceID, Type: []string{"VirtualMachine"}}}}, t),
				},
			},
			hasRPCConnection: false,
			wantResp:         nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not evaluate evidence: could not retrieve metric definitions: could not retrieve metric list from orchestrator")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			if tt.hasRPCConnection {
				s.evidenceStoreAddress.opts = []grpc.DialOption{grpc.WithContextDialer(bufConnDialer)}
				s.orchestratorAddress.opts = []grpc.DialOption{grpc.WithContextDialer(bufConnDialer)}
			} else {
				// clear the evidence URL, just to be sure
				s.evidenceStoreAddress.target = ""
				s.orchestratorAddress.target = ""
			}

			gotResp, err := s.AssessEvidence(tt.args.in0, &assessment.AssessEvidenceRequest{Evidence: tt.args.evidence})

			tt.wantErr(t, err)

			// Check response
			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp.Status, gotResp.Status)
				assert.Contains(t, gotResp.StatusMessage, tt.wantResp.StatusMessage)
			}
		})
	}
}

// TestAssessEvidences tests AssessEvidences
func TestService_AssessEvidences(t *testing.T) {
	type fields struct {
		ResultHooks                   []assessment.ResultHookFunc
		results                       map[string]*assessment.AssessmentResult
		evidenceStoreStreams          *api.StreamsOf[evidence.EvidenceStore_StoreEvidencesClient, *evidence.StoreEvidenceRequest]
		orchestratorStreams           *api.StreamsOf[orchestrator.Orchestrator_StoreAssessmentResultsClient, *orchestrator.StoreAssessmentResultRequest]
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
						Id:             testdata.MockEvidenceID,
						Timestamp:      timestamppb.Now(),
						CloudServiceId: testdata.MockCloudServiceID,
						Resource:       toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: testdata.MockResourceID, Type: []string{"VirtualMachine"}}}}, t)}}),
			},
			wantErr: false,
			wantRespMessage: &assessment.AssessEvidenceResponse{
				Status:        assessment.AssessEvidenceResponse_FAILED,
				StatusMessage: "rpc error: code = InvalidArgument desc = rpc error: code = InvalidArgument desc = invalid request: invalid AssessEvidenceRequest.Evidence: embedded message failed validation | caused by: invalid Evidence.ToolId: value length must be at least 1 runes",
			},
		},
		{
			name: "Missing evidenceID",
			fields: fields{
				results: make(map[string]*assessment.AssessmentResult)},
			args: args{
				streamToServer: createMockAssessmentServerStream(&assessment.AssessEvidenceRequest{
					Evidence: &evidence.Evidence{
						Timestamp:      timestamppb.Now(),
						ToolId:         testdata.MockEvidenceToolID,
						CloudServiceId: testdata.MockCloudServiceID,
						Resource:       toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: testdata.MockResourceID, Type: []string{"VirtualMachine"}}}}, t)}}),
			},
			wantErr: false,
			wantRespMessage: &assessment.AssessEvidenceResponse{
				Status:        assessment.AssessEvidenceResponse_FAILED,
				StatusMessage: "rpc error: code = InvalidArgument desc = rpc error: code = InvalidArgument desc = invalid request: invalid AssessEvidenceRequest.Evidence: embedded message failed validation | caused by: invalid Evidence.Id: value must be a valid UUID | caused by: invalid uuid format",
			},
		},
		{
			name: "Assess evidences",
			fields: fields{
				results:              make(map[string]*assessment.AssessmentResult),
				evidenceStoreStreams: api.NewStreamsOf(api.WithLogger[evidence.EvidenceStore_StoreEvidencesClient, *evidence.StoreEvidenceRequest](log)),
				orchestratorStreams:  api.NewStreamsOf(api.WithLogger[orchestrator.Orchestrator_StoreAssessmentResultsClient, *orchestrator.StoreAssessmentResultRequest](log)),
			},
			args: args{
				streamToServer: createMockAssessmentServerStream(&assessment.AssessEvidenceRequest{
					Evidence: &evidence.Evidence{
						Id:             testdata.MockEvidenceID,
						Timestamp:      timestamppb.Now(),
						ToolId:         testdata.MockEvidenceToolID,
						CloudServiceId: testdata.MockCloudServiceID,
						Resource:       toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: testdata.MockResourceID, Type: []string{"VirtualMachine"}}}}, t)}}),
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
						Timestamp:      timestamppb.Now(),
						ToolId:         testdata.MockEvidenceToolID,
						CloudServiceId: testdata.MockCloudServiceID,
						Resource:       toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: testdata.MockResourceID, Type: []string{"VirtualMachine"}}}}, t)}}),
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
						Timestamp:      timestamppb.Now(),
						ToolId:         testdata.MockEvidenceToolID,
						CloudServiceId: testdata.MockCloudServiceID,
						Resource:       toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: testdata.MockResourceID, Type: []string{"VirtualMachine"}}}}, t)}}),
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
				evidenceStoreStreams:          tt.fields.evidenceStoreStreams,
				evidenceStoreAddress: grpcTarget{
					opts: []grpc.DialOption{grpc.WithContextDialer(bufConnDialer)},
				},
				orchestratorStreams: tt.fields.orchestratorStreams,
				orchestratorAddress: grpcTarget{
					opts: []grpc.DialOption{grpc.WithContextDialer(bufConnDialer)},
				},
				pe: policies.NewRegoEval(),
			}

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
				if !reflect.DeepEqual(responseFromServer, tt.wantRespMessage) {
					t.Errorf("AssessEvidences() gotResp = %v, want %v", responseFromServer, tt.wantRespMessage)
				}
			} else {
				assert.Contains(t, err.Error(), tt.wantErrMessage)
			}
		})
	}
}

func TestService_AssessmentResultHooks(t *testing.T) {
	var (
		hookCallCounter = 0
		wg              sync.WaitGroup
		hookCounts      = 18
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
						Id:             testdata.MockEvidenceID,
						ToolId:         testdata.MockEvidenceToolID,
						Timestamp:      timestamppb.Now(),
						CloudServiceId: testdata.MockCloudServiceID,
						Resource: toStruct(&voc.VirtualMachine{
							Compute: &voc.Compute{
								Resource: &voc.Resource{
									ID:   testdata.MockResourceID,
									Type: []string{"VirtualMachine", "Compute", "Resource"}},
							},
							BootLogging: &voc.BootLogging{
								Logging: &voc.Logging{
									LoggingService:  []voc.ResourceID{"SomeResourceId2"},
									Enabled:         true,
									RetentionPeriod: 36,
								},
							},
							OsLogging: &voc.OSLogging{
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
			wantResp: &assessment.AssessEvidenceResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hookCallCounter = 0
			s := NewService(WithEvidenceStoreAddress("", grpc.WithContextDialer(bufConnDialer)), WithOrchestratorAddress("", grpc.WithContextDialer(bufConnDialer)))

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

func TestService_ListAssessmentResults(t *testing.T) {
	type fields struct {
		evidenceStoreStreams  *api.StreamsOf[evidence.EvidenceStore_StoreEvidencesClient, *evidence.StoreEvidenceRequest]
		evidenceStoreAddress  string
		orchestratorStreams   *api.StreamsOf[orchestrator.Orchestrator_StoreAssessmentResultsClient, *orchestrator.StoreAssessmentResultRequest]
		orchestratorClient    orchestrator.OrchestratorClient
		orchestratorAddress   string
		metricEventStream     orchestrator.Orchestrator_SubscribeMetricChangeEventsClient
		resultHooks           []assessment.ResultHookFunc
		results               map[string]*assessment.AssessmentResult
		cachedConfigurations  map[string]cachedConfiguration
		authorizer            api.Authorizer
		grpcOptsEvidenceStore []grpc.DialOption
		grpcOptsOrchestrator  []grpc.DialOption
		pe                    policies.PolicyEval
	}
	type args struct {
		in0 context.Context
		req *assessment.ListAssessmentResultsRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *assessment.ListAssessmentResultsResponse
		wantErr bool
	}{
		{
			name: "Missing request",
			args: args{
				req: nil,
			},
			wantRes: nil,
			wantErr: true,
		},
		{
			name: "Empty request",
			fields: fields{
				results: map[string]*assessment.AssessmentResult{
					"1": {
						Id: "1",
					},
				},
			},
			args: args{
				req: &assessment.ListAssessmentResultsRequest{},
			},
			wantRes: &assessment.ListAssessmentResultsResponse{
				NextPageToken: "",
				Results: []*assessment.AssessmentResult{
					{
						Id: "1",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "single page result",
			fields: fields{
				results: map[string]*assessment.AssessmentResult{
					"1": {
						Id: "1",
					},
				},
			},
			args: args{
				req: &assessment.ListAssessmentResultsRequest{},
			},
			wantRes: &assessment.ListAssessmentResultsResponse{
				NextPageToken: "",
				Results: []*assessment.AssessmentResult{
					{
						Id: "1",
					},
				},
			},
		},
		{
			name: "multiple page result - first page",
			fields: fields{
				results: map[string]*assessment.AssessmentResult{
					"11111111-1111-1111-1111-111111111111": {Id: "11111111-1111-1111-1111-111111111111"},
					"22222222-2222-2222-2222-222222222222": {Id: "22222222-2222-2222-2222-222222222222"},
					"33333333-3333-3333-3333-333333333333": {Id: "33333333-3333-3333-3333-333333333333"},
					"44444444-4444-4444-4444-444444444444": {Id: "44444444-4444-4444-4444-444444444444"},
					"55555555-5555-5555-5555-555555555555": {Id: "55555555-5555-5555-5555-555555555555"},
					"66666666-6666-6666-6666-666666666666": {Id: "66666666-6666-6666-6666-666666666666"},
					"77777777-7777-7777-7777-777777777777": {Id: "77777777-7777-7777-7777-777777777777"},
					"88888888-8888-8888-8888-888888888888": {Id: "88888888-8888-8888-8888-888888888888"},
					"99999999-9999-9999-9999-999999999999": {Id: "99999999-9999-9999-9999-999999999999"},
				},
			},
			args: args{
				req: &assessment.ListAssessmentResultsRequest{
					PageSize: 2,
				},
			},
			wantRes: &assessment.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					{Id: "11111111-1111-1111-1111-111111111111"}, {Id: "22222222-2222-2222-2222-222222222222"},
				},
				NextPageToken: func() string {
					token, _ := (&api.PageToken{Start: 2, Size: 2}).Encode()
					return token
				}(),
			},
		},
		{
			name: "multiple page result - second page",
			fields: fields{
				results: map[string]*assessment.AssessmentResult{
					"11111111-1111-1111-1111-111111111111": {Id: "11111111-1111-1111-1111-111111111111"},
					"22222222-2222-2222-2222-222222222222": {Id: "22222222-2222-2222-2222-222222222222"},
					"33333333-3333-3333-3333-333333333333": {Id: "33333333-3333-3333-3333-333333333333"},
					"44444444-4444-4444-4444-444444444444": {Id: "44444444-4444-4444-4444-444444444444"},
					"55555555-5555-5555-5555-555555555555": {Id: "55555555-5555-5555-5555-555555555555"},
					"66666666-6666-6666-6666-666666666666": {Id: "66666666-6666-6666-6666-666666666666"},
					"77777777-7777-7777-7777-777777777777": {Id: "77777777-7777-7777-7777-777777777777"},
					"88888888-8888-8888-8888-888888888888": {Id: "88888888-8888-8888-8888-888888888888"},
					"99999999-9999-9999-9999-999999999999": {Id: "99999999-9999-9999-9999-999999999999"},
				},
			},
			args: args{
				req: &assessment.ListAssessmentResultsRequest{
					PageSize: 2,
					PageToken: func() string {
						token, _ := (&api.PageToken{Start: 2, Size: 2}).Encode()
						return token
					}(),
				},
			},
			wantRes: &assessment.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					{Id: "33333333-3333-3333-3333-333333333333"}, {Id: "44444444-4444-4444-4444-444444444444"},
				},
				NextPageToken: func() string {
					token, _ := (&api.PageToken{Start: 4, Size: 2}).Encode()
					return token
				}(),
			},
		},
		{
			name: "multiple page result - last page",
			fields: fields{
				results: map[string]*assessment.AssessmentResult{
					"11111111-1111-1111-1111-111111111111": {Id: "11111111-1111-1111-1111-111111111111"},
					"22222222-2222-2222-2222-222222222222": {Id: "22222222-2222-2222-2222-222222222222"},
					"33333333-3333-3333-3333-333333333333": {Id: "33333333-3333-3333-3333-333333333333"},
					"44444444-4444-4444-4444-444444444444": {Id: "44444444-4444-4444-4444-444444444444"},
					"55555555-5555-5555-5555-555555555555": {Id: "55555555-5555-5555-5555-555555555555"},
					"66666666-6666-6666-6666-666666666666": {Id: "66666666-6666-6666-6666-666666666666"},
					"77777777-7777-7777-7777-777777777777": {Id: "77777777-7777-7777-7777-777777777777"},
					"88888888-8888-8888-8888-888888888888": {Id: "88888888-8888-8888-8888-888888888888"},
					"99999999-9999-9999-9999-999999999999": {Id: "99999999-9999-9999-9999-999999999999"},
				},
			},
			args: args{
				req: &assessment.ListAssessmentResultsRequest{
					PageSize: 2,
					PageToken: func() string {
						token, _ := (&api.PageToken{Start: 8, Size: 2}).Encode()
						return token
					}(),
				},
			},
			wantRes: &assessment.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					{Id: "99999999-9999-9999-9999-999999999999"},
				},
				NextPageToken: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				evidenceStoreStreams: tt.fields.evidenceStoreStreams,
				evidenceStoreAddress: grpcTarget{
					target: tt.fields.evidenceStoreAddress,
					opts:   tt.fields.grpcOptsEvidenceStore,
				},
				orchestratorAddress: grpcTarget{
					target: tt.fields.orchestratorAddress,
					opts:   tt.fields.grpcOptsOrchestrator,
				},
				orchestratorStreams:  tt.fields.orchestratorStreams,
				orchestratorClient:   tt.fields.orchestratorClient,
				metricEventStream:    tt.fields.metricEventStream,
				resultHooks:          tt.fields.resultHooks,
				results:              tt.fields.results,
				cachedConfigurations: tt.fields.cachedConfigurations,
				authorizer:           tt.fields.authorizer,
				pe:                   tt.fields.pe,
			}
			gotRes, err := s.ListAssessmentResults(tt.args.in0, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ListAssessmentResults() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.ListAssessmentResults() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
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

func (*mockAssessmentServerStream) CloseSend() error {
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

func (m *mockAssessmentServerStream) Send(response *assessment.AssessEvidenceResponse) error {
	m.SentFromServer <- response
	return nil
}

func (*mockAssessmentServerStream) SendAndClose() error {
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

func (*mockAssessmentServerStream) SetHeader(metadata.MD) error {
	return nil
}

func (*mockAssessmentServerStream) SendHeader(metadata.MD) error {
	return nil
}

func (*mockAssessmentServerStream) SetTrailer(metadata.MD) {
}

func (*mockAssessmentServerStream) Context() context.Context {
	return nil
}

func (*mockAssessmentServerStream) SendMsg(interface{}) error {
	return nil
}

func (*mockAssessmentServerStream) RecvMsg(interface{}) error {
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

func (*mockAssessmentServerStreamWithRecvErr) Send(*assessment.AssessEvidenceResponse) error {
	panic("implement me")
}

func (*mockAssessmentServerStreamWithRecvErr) Recv() (*assessment.AssessEvidenceRequest, error) {
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

func TestService_HandleEvidence(t *testing.T) {
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
					Id:             testdata.MockEvidenceID,
					ToolId:         testdata.MockEvidenceToolID,
					Timestamp:      timestamppb.Now(),
					CloudServiceId: testdata.MockCloudServiceID,
					Resource: toStruct(voc.VirtualMachine{
						Compute: &voc.Compute{
							Resource: &voc.Resource{ID: testdata.MockResourceID, Type: []string{"VirtualMachine"}},
						},
						BootLogging: &voc.BootLogging{
							Logging: &voc.Logging{
								LoggingService:  nil,
								Enabled:         true,
								RetentionPeriod: 0,
							},
						},
					}, t),
				},
				resourceId: testdata.MockResourceID,
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
					Id:             testdata.MockEvidenceID,
					ToolId:         testdata.MockEvidenceToolID,
					Timestamp:      timestamppb.Now(),
					CloudServiceId: testdata.MockCloudServiceID,
					Resource:       toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: testdata.MockResourceID, Type: []string{}}}}, t),
				},
				resourceId: testdata.MockResourceID,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Error(t, err)
				// Check if error message contains "empty" (list of types)
				assert.Contains(t, err.Error(), "empty")
				return true
			},
		},
		{
			name: "evidence store stream error",
			fields: fields{
				hasOrchestratorStream:  true,
				hasEvidenceStoreStream: false,
			},
			args: args{
				evidence: &evidence.Evidence{
					Id:             testdata.MockEvidenceID,
					ToolId:         testdata.MockEvidenceToolID,
					Timestamp:      timestamppb.Now(),
					CloudServiceId: testdata.MockCloudServiceID,
					Resource: toStruct(voc.VirtualMachine{
						Compute: &voc.Compute{
							Resource: &voc.Resource{ID: testdata.MockResourceID, Type: []string{"VirtualMachine"}},
						},
						BootLogging: &voc.BootLogging{
							Logging: &voc.Logging{
								LoggingService:  nil,
								Enabled:         true,
								RetentionPeriod: 0,
							},
						}}, t),
				},
				resourceId: testdata.MockResourceID,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				if !assert.NotEmpty(t, err) {
					return false
				}

				return assert.Contains(t, err.Error(), "could not get stream to evidence store")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()

			// Mock streams for target services
			if tt.fields.hasEvidenceStoreStream {
				s.evidenceStoreAddress.opts = []grpc.DialOption{grpc.WithContextDialer(bufConnDialer)}
			} else {
				s.evidenceStoreAddress.opts = []grpc.DialOption{grpc.WithContextDialer(nil)}
			}
			if tt.fields.hasOrchestratorStream {
				s.orchestratorAddress.opts = []grpc.DialOption{grpc.WithContextDialer(bufConnDialer)}
			} else {
				s.orchestratorAddress.opts = []grpc.DialOption{grpc.WithContextDialer(nil)}
			}

			// Two tests: 1st) wantErr function. 2nd) if wantErr false then check if a result is added to map
			if !tt.wantErr(t, s.handleEvidence(tt.args.evidence, tt.args.resourceId), fmt.Sprintf("handleEvidence(%v, %v)", tt.args.evidence, tt.args.resourceId)) {
				assert.NotEmpty(t, s.results)
				// Check the result by validation
				for _, result := range s.results {
					err := result.Validate()
					assert.NoError(t, err)
				}
			}

		})
	}
}

func TestService_initOrchestratorStoreStream(t *testing.T) {
	type fields struct {
		opts []service.Option[Service]
	}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Invalid RPC connection",
			args: args{
				url: "localhost:1",
			},
			fields: fields{
				opts: []service.Option[Service]{
					WithOrchestratorAddress("localhost:1"),
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				s, _ := status.FromError(errors.Unwrap(err))
				return assert.Equal(t, codes.Unavailable, s.Code())
			},
		},
		// TODO: Fix test
		// {
		// 	name: "Authenticated RPC connection with valid user",
		// 	args: args{
		// 		url: "bufnet",
		// 	},
		// 	fields: fields{
		// 		opts: []ServiceOption{
		// 			WithOrchestratorAddress("bufnet"),
		// 			WithOAuth2Authorizer(testutil.AuthClientConfig(authPort)),
		// 			WithAdditionalGRPCOpts(grpc.WithContextDialer(bufConnDialer)),
		// 		},
		// 	},
		// },
		{
			name: "Authenticated RPC connection with invalid user",
			args: args{
				url: "bufnet",
			},
			fields: fields{
				opts: []service.Option[Service]{
					WithOrchestratorAddress("bufnet", grpc.WithContextDialer(bufConnDialer)),
					WithOAuth2Authorizer(testutil.AuthClientConfig(authPort)),
				},
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
			stream, err := s.initOrchestratorStream(tt.args.url, s.orchestratorAddress.opts...)

			if tt.wantErr != nil {
				tt.wantErr(t, err)
			} else {
				assert.NotEmpty(t, stream)
			}
		})
	}
}

func TestService_recvEventsLoop(t *testing.T) {
	type fields struct {
		evidenceStoreStreams  *api.StreamsOf[evidence.EvidenceStore_StoreEvidencesClient, *evidence.StoreEvidenceRequest]
		evidenceStoreAddress  string
		orchestratorStreams   *api.StreamsOf[orchestrator.Orchestrator_StoreAssessmentResultsClient, *orchestrator.StoreAssessmentResultRequest]
		orchestratorClient    orchestrator.OrchestratorClient
		orchestratorAddress   string
		metricEventStream     orchestrator.Orchestrator_SubscribeMetricChangeEventsClient
		resultHooks           []assessment.ResultHookFunc
		results               map[string]*assessment.AssessmentResult
		cachedConfigurations  map[string]cachedConfiguration
		authorizer            api.Authorizer
		grpcOptsEvidenceStore []grpc.DialOption
		grpcOptsOrchestrator  []grpc.DialOption
	}
	tests := []struct {
		name      string
		fields    fields
		wantEvent *orchestrator.MetricChangeEvent
	}{
		{
			name: "Receive event",
			fields: fields{
				metricEventStream: &testutil.ListRecvStreamerOf[*orchestrator.MetricChangeEvent]{Messages: []*orchestrator.MetricChangeEvent{
					{
						Type: orchestrator.MetricChangeEvent_TYPE_CONFIG_CHANGED,
					},
				}},
			},
			wantEvent: &orchestrator.MetricChangeEvent{
				Type: orchestrator.MetricChangeEvent_TYPE_CONFIG_CHANGED,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				evidenceStoreStreams: tt.fields.evidenceStoreStreams,
				evidenceStoreAddress: grpcTarget{
					target: tt.fields.evidenceStoreAddress,
					opts:   tt.fields.grpcOptsEvidenceStore,
				},
				orchestratorAddress: grpcTarget{
					target: tt.fields.orchestratorAddress,
					opts:   tt.fields.grpcOptsOrchestrator,
				},
				orchestratorStreams:  tt.fields.orchestratorStreams,
				orchestratorClient:   tt.fields.orchestratorClient,
				metricEventStream:    tt.fields.metricEventStream,
				resultHooks:          tt.fields.resultHooks,
				results:              tt.fields.results,
				cachedConfigurations: tt.fields.cachedConfigurations,
				authorizer:           tt.fields.authorizer,
			}
			rec := &eventRecorder{}
			svc.pe = rec
			svc.recvEventsLoop()

			if !proto.Equal(rec.event, tt.wantEvent) {
				t.Errorf("recvEventsLoop() = %v, want %v", rec.event, tt.wantEvent)
			}
		})
	}
}

type eventRecorder struct {
	event *orchestrator.MetricChangeEvent
	done  bool
}

func (*eventRecorder) Eval(_ *evidence.Evidence, _ policies.MetricsSource) (data []*policies.Result, err error) {
	return nil, nil
}

func (e *eventRecorder) HandleMetricEvent(event *orchestrator.MetricChangeEvent) (err error) {
	if e.done {
		return nil
	}

	e.event = event
	e.done = true

	return nil
}

func TestService_MetricImplementation(t *testing.T) {
	type fields struct {
		UnimplementedAssessmentServer assessment.UnimplementedAssessmentServer
		isEvidenceStoreDisabled       bool
		evidenceStoreStreams          *api.StreamsOf[evidence.EvidenceStore_StoreEvidencesClient, *evidence.StoreEvidenceRequest]
		evidenceStoreAddress          grpcTarget
		orchestratorStreams           *api.StreamsOf[orchestrator.Orchestrator_StoreAssessmentResultsClient, *orchestrator.StoreAssessmentResultRequest]
		orchestratorClient            orchestrator.OrchestratorClient
		orchestratorAddress           grpcTarget
		metricEventStream             orchestrator.Orchestrator_SubscribeMetricChangeEventsClient
		resultHooks                   []assessment.ResultHookFunc
		results                       map[string]*assessment.AssessmentResult
		cachedConfigurations          map[string]cachedConfiguration
		authorizer                    api.Authorizer
		pe                            policies.PolicyEval
		evalPkg                       string
	}
	type args struct {
		lang   assessment.MetricImplementation_Language
		metric string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantImpl *assessment.MetricImplementation
		wantErr  assert.ErrorAssertionFunc
	}{

		{
			name: "Unspecified language",
			args: args{
				lang: assessment.MetricImplementation_LANGUAGE_UNSPECIFIED,
			},
			wantImpl: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "unsupported language")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				UnimplementedAssessmentServer: tt.fields.UnimplementedAssessmentServer,
				isEvidenceStoreDisabled:       tt.fields.isEvidenceStoreDisabled,
				evidenceStoreStreams:          tt.fields.evidenceStoreStreams,
				evidenceStoreAddress:          tt.fields.evidenceStoreAddress,
				orchestratorStreams:           tt.fields.orchestratorStreams,
				orchestratorClient:            tt.fields.orchestratorClient,
				orchestratorAddress:           tt.fields.orchestratorAddress,
				metricEventStream:             tt.fields.metricEventStream,
				resultHooks:                   tt.fields.resultHooks,
				results:                       tt.fields.results,
				cachedConfigurations:          tt.fields.cachedConfigurations,
				authorizer:                    tt.fields.authorizer,
				pe:                            tt.fields.pe,
				evalPkg:                       tt.fields.evalPkg,
			}
			gotImpl, err := svc.MetricImplementation(tt.args.lang, tt.args.metric)
			tt.wantErr(t, err)
			assert.Equal(t, tt.wantImpl, gotImpl)
		})
	}
}
