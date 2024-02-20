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

	"clouditor.io/clouditor/internal/testutil/servicetest/evidencetest"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/internal/testutil/clitest"
	"clouditor.io/clouditor/internal/testutil/servicetest"
	"clouditor.io/clouditor/policies"
	"clouditor.io/clouditor/service"
	"clouditor.io/clouditor/voc"

	"github.com/google/uuid"
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
		want assert.ValueAssertionFunc
	}{
		{
			name: "AssessmentServer created with option rego package name",
			args: args{
				opts: []service.Option[Service]{
					WithRegoPackageName("testPkg"),
				},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				s := i1.(*Service)
				return assert.Equal(t, "testPkg", s.evalPkg)
			},
		},
		{
			name: "AssessmentServer created with option authorizer",
			args: args{
				opts: []service.Option[Service]{
					WithAuthorizer(api.NewOAuthAuthorizerFromClientCredentials(&clientcredentials.Config{})),
				},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				s := i1.(*Service)
				return assert.Equal(t, api.NewOAuthAuthorizerFromClientCredentials(&clientcredentials.Config{}), s.orchestrator.Authorizer())
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
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				s := i1.(*Service)
				return assert.Equal(t, "localhost:9091", s.evidenceStore.Target) &&
					assert.Equal(t, "localhost:9092", s.orchestrator.Target)
			},
		},
		{
			name: "AssessmentServer without EvidenceStore",
			args: args{
				opts: []service.Option[Service]{
					WithoutEvidenceStore(),
				},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				s := i1.(*Service)
				return assert.True(t, s.isEvidenceStoreDisabled)
			},
		},
		{
			name: "AssessmentServer with oauth2 authorizer",
			args: args{
				opts: []service.Option[Service]{
					WithOAuth2Authorizer(&clientcredentials.Config{ClientID: "client"}),
				},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				s := i1.(*Service)
				return assert.NotNil(t, s.orchestrator.Authorizer())
			},
		},
		{
			name: "AssessmentServer with authorization strategy",
			args: args{
				opts: []service.Option[Service]{
					WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(true)),
				},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				s := i1.(*Service)
				a := s.authz.(*servicetest.AuthorizationStrategyMock)
				return assert.NotNil(t, a)
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

// TestAssessEvidence tests AssessEvidence
func TestService_AssessEvidence(t *testing.T) {
	type fields struct {
		authz         service.AuthorizationStrategy
		evidenceStore *api.RPCConnection[evidence.EvidenceStoreClient]
		orchestrator  *api.RPCConnection[orchestrator.OrchestratorClient]
	}
	type args struct {
		in0      context.Context
		evidence *evidence.Evidence
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp *assessment.AssessEvidenceResponse
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "Missing evidence",
			args: args{
				in0: context.TODO(),
			},
			wantResp: nil,
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
			wantResp: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid request: invalid AssessEvidenceRequest.Evidence: embedded message failed validation | caused by: invalid Evidence.Id: value must be a valid UUID | caused by: invalid uuid format")
			},
		},
		{
			name: "Assess resource without id",
			fields: fields{
				evidenceStore: api.NewRPCConnection("bufnet", evidence.NewEvidenceStoreClient, grpc.WithContextDialer(bufConnDialer)),
				orchestrator:  api.NewRPCConnection("bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(bufConnDialer)),
			},
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					ToolId:    testdata.MockEvidenceToolID1,
					Timestamp: timestamppb.Now(),
					Resource:  toStruct(voc.VirtualMachine{}, t),
				},
			},
			wantResp: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid request: invalid AssessEvidenceRequest.Evidence: embedded message failed validation | caused by: invalid Evidence.Id: value must be a valid UUID | caused by: invalid uuid format")
			},
		},
		{
			name: "Assess resource without tool id",
			fields: fields{
				evidenceStore: api.NewRPCConnection("bufnet", evidence.NewEvidenceStoreClient, grpc.WithContextDialer(bufConnDialer)),
				orchestrator:  api.NewRPCConnection("bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(bufConnDialer)),
			},
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					Id:             testdata.MockEvidenceID1,
					Timestamp:      timestamppb.Now(),
					CloudServiceId: testdata.MockCloudServiceID1,
					Resource:       toStruct(voc.VirtualMachine{}, t),
				},
			},
			wantResp: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid request: invalid AssessEvidenceRequest.Evidence: embedded message failed validation | caused by: invalid Evidence.ToolId: value length must be at least 1 runes")
			},
		},
		{
			name: "Assess resource without timestamp",
			fields: fields{
				evidenceStore: api.NewRPCConnection("bufnet", evidence.NewEvidenceStoreClient, grpc.WithContextDialer(bufConnDialer)),
				orchestrator:  api.NewRPCConnection("bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(bufConnDialer)),
			},
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					Id:             testdata.MockEvidenceID1,
					ToolId:         testdata.MockEvidenceToolID1,
					CloudServiceId: testdata.MockCloudServiceID1,
					Resource:       toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: testdata.MockResourceID1, Type: []string{"VirtualMachine"}}}}, t),
				},
			},
			wantResp: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid request: invalid AssessEvidenceRequest.Evidence: embedded message failed validation | caused by: invalid Evidence.Timestamp: value is required")
			},
		},
		{
			name: "Assess resource",
			fields: fields{
				evidenceStore: api.NewRPCConnection("bufnet", evidence.NewEvidenceStoreClient, grpc.WithContextDialer(bufConnDialer)),
				orchestrator:  api.NewRPCConnection("bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(bufConnDialer)),
				authz:         servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					Id:             testdata.MockEvidenceID1,
					ToolId:         testdata.MockEvidenceToolID1,
					Timestamp:      timestamppb.Now(),
					Resource:       toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: testdata.MockResourceID1, Type: []string{"VirtualMachine"}}}}, t),
					CloudServiceId: testdata.MockCloudServiceID1},
			},
			wantResp: &assessment.AssessEvidenceResponse{},
			wantErr:  assert.NoError,
		},
		{
			name: "Assess resource of wrong could service",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockCloudServiceID2),
			},
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					Id:             testdata.MockEvidenceID1,
					ToolId:         testdata.MockEvidenceToolID1,
					Timestamp:      timestamppb.Now(),
					Resource:       toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: testdata.MockResourceID1, Type: []string{"VirtualMachine"}}}}, t),
					CloudServiceId: testdata.MockCloudServiceID1},
			},
			wantResp: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, service.ErrPermissionDenied.Error())
			},
		},
		{
			name: "Assess resource without resource id",
			fields: fields{
				evidenceStore: api.NewRPCConnection("bufnet", evidence.NewEvidenceStoreClient, grpc.WithContextDialer(bufConnDialer)),
				orchestrator:  api.NewRPCConnection("bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(bufConnDialer)),
			},
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					Id:             testdata.MockEvidenceID1,
					ToolId:         testdata.MockEvidenceToolID1,
					Timestamp:      timestamppb.Now(),
					Resource:       toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{Type: []string{"VirtualMachine"}}}}, t),
					CloudServiceId: testdata.MockCloudServiceID1,
				},
			},
			wantResp: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, evidence.ErrResourceIdIsEmpty.Error())
			},
		},
		{
			name: "No RPC connections",
			fields: fields{
				evidenceStore: api.NewRPCConnection("bufnet", evidence.NewEvidenceStoreClient, grpc.WithContextDialer(connectionRefusedDialer)),
				orchestrator:  api.NewRPCConnection("bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(connectionRefusedDialer)),
				authz:         servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					Id:             testdata.MockEvidenceID1,
					ToolId:         testdata.MockEvidenceToolID1,
					Timestamp:      timestamppb.Now(),
					CloudServiceId: testdata.MockCloudServiceID1,
					Resource:       toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: testdata.MockResourceID1, Type: []string{"VirtualMachine"}}}}, t),
				},
			},
			wantResp: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "connection refused")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				evidenceStore:        tt.fields.evidenceStore,
				orchestrator:         tt.fields.orchestrator,
				evidenceStoreStreams: api.NewStreamsOf(api.WithLogger[evidence.EvidenceStore_StoreEvidencesClient, *evidence.StoreEvidenceRequest](log)),
				orchestratorStreams:  api.NewStreamsOf(api.WithLogger[orchestrator.Orchestrator_StoreAssessmentResultsClient, *orchestrator.StoreAssessmentResultRequest](log)),
				cachedConfigurations: make(map[string]cachedConfiguration),
				pe:                   policies.NewRegoEval(policies.WithPackageName(policies.DefaultRegoPackage)),
				authz:                tt.fields.authz,
			}
			gotResp, err := s.AssessEvidence(tt.args.in0, &assessment.AssessEvidenceRequest{Evidence: tt.args.evidence})

			tt.wantErr(t, err)

			// Check response
			assert.Empty(t, gotResp)
		})
	}
}

// TestService_AssessEvidence_DetectMisconfiguredEvidenceEvenWhenAlreadyCached tests the following workflow: First an
// evidence with a VM resource is assessed. The resource contains all required fields s.t. the metric cache is filled
// with all applicable metrics. In a second step we assess another evidence. It is also of type "VirtualMachine" but all
// other fields are not set (e.g. MalwareProtection). Thus, metric will be applied and therefore no error occurs in
// AssessEvidence-handleEvidence (assessment.go) which loops over all evaluations
// Todo: Add it to table test above (would probably need some function injection in test cases like we do with storage)
func TestService_AssessEvidence_DetectMisconfiguredEvidenceEvenWhenAlreadyCached(t *testing.T) {
	s := &Service{
		evidenceStore: api.NewRPCConnection(
			"bufnet", evidence.NewEvidenceStoreClient, grpc.WithContextDialer(bufConnDialer)),
		orchestrator: api.NewRPCConnection(
			"bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(bufConnDialer)),
		authz: servicetest.NewAuthorizationStrategy(true),
		evidenceStoreStreams: api.NewStreamsOf(
			api.WithLogger[evidence.EvidenceStore_StoreEvidencesClient, *evidence.StoreEvidenceRequest](log)),
		orchestratorStreams: api.NewStreamsOf(
			api.WithLogger[orchestrator.Orchestrator_StoreAssessmentResultsClient,
				*orchestrator.StoreAssessmentResultRequest](log)),
		cachedConfigurations: make(map[string]cachedConfiguration),
		pe:                   policies.NewRegoEval(policies.WithPackageName(policies.DefaultRegoPackage)),
	}
	// First assess evidence with a valid VM resource s.t. the cache is created for the combination of resource type and
	// tool id (="VirtualMachine-{testdata.MockEvidenceToolID}")
	e := evidencetest.MockEvidence1
	e.Resource = toStruct(voc.VirtualMachine{
		Compute: &voc.Compute{
			Resource: &voc.Resource{ID: testdata.MockResourceID1, Type: []string{"VirtualMachine"}}}}, t)
	_, err := s.AssessEvidence(context.Background(), &assessment.AssessEvidenceRequest{Evidence: e})
	assert.NoError(t, err)

	// Now assess a new evidence which has not a valid format other than the resource type and tool id is set correctly
	// Prepare resource
	r := map[string]any{
		// Make sure both evidences have the same type (for caching key)
		"type": e.Resource.GetStructValue().AsMap()["type"],
		"id":   uuid.NewString(),
	}
	v, err := structpb.NewValue(r)
	assert.NoError(t, err)
	_, err = s.AssessEvidence(context.Background(), &assessment.AssessEvidenceRequest{Evidence: &evidence.Evidence{
		Id:             uuid.NewString(),
		Timestamp:      timestamppb.Now(),
		CloudServiceId: testdata.MockCloudServiceID1,
		// Make sure both evidences have the same tool id (for caching key)
		ToolId:   e.ToolId,
		Raw:      nil,
		Resource: v,
	}})
	assert.NoError(t, err)
}

// TestAssessEvidences tests AssessEvidences
func TestService_AssessEvidences(t *testing.T) {
	type fields struct {
		ResultHooks          []assessment.ResultHookFunc
		evidenceStoreStreams *api.StreamsOf[evidence.EvidenceStore_StoreEvidencesClient, *evidence.StoreEvidenceRequest]
		orchestratorStreams  *api.StreamsOf[orchestrator.Orchestrator_StoreAssessmentResultsClient, *orchestrator.StoreAssessmentResultRequest]
		authz                service.AuthorizationStrategy
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
		wantRespMessage *assessment.AssessEvidencesResponse
	}{
		{
			name: "Missing toolId",
			args: args{
				streamToServer: createMockAssessmentServerStream(&assessment.AssessEvidenceRequest{
					Evidence: &evidence.Evidence{
						Id:             testdata.MockEvidenceID1,
						Timestamp:      timestamppb.Now(),
						CloudServiceId: testdata.MockCloudServiceID1,
						Resource:       toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: testdata.MockResourceID1, Type: []string{"VirtualMachine"}}}}, t)}}),
			},
			wantErr: false,
			wantRespMessage: &assessment.AssessEvidencesResponse{
				Status:        assessment.AssessEvidencesResponse_FAILED,
				StatusMessage: "rpc error: code = InvalidArgument desc = invalid request: invalid AssessEvidenceRequest.Evidence: embedded message failed validation | caused by: invalid Evidence.ToolId: value length must be at least 1 runes",
			},
		},
		{
			name: "Missing evidenceID",
			args: args{
				streamToServer: createMockAssessmentServerStream(&assessment.AssessEvidenceRequest{
					Evidence: &evidence.Evidence{
						Timestamp:      timestamppb.Now(),
						ToolId:         testdata.MockEvidenceToolID1,
						CloudServiceId: testdata.MockCloudServiceID1,
						Resource:       toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: testdata.MockResourceID1, Type: []string{"VirtualMachine"}}}}, t)}}),
			},
			wantErr: false,
			wantRespMessage: &assessment.AssessEvidencesResponse{
				Status:        assessment.AssessEvidencesResponse_FAILED,
				StatusMessage: "rpc error: code = InvalidArgument desc = invalid request: invalid AssessEvidenceRequest.Evidence: embedded message failed validation | caused by: invalid Evidence.Id: value must be a valid UUID | caused by: invalid uuid format",
			},
		},
		{
			name: "Assess evidences",
			fields: fields{
				evidenceStoreStreams: api.NewStreamsOf(api.WithLogger[evidence.EvidenceStore_StoreEvidencesClient, *evidence.StoreEvidenceRequest](log)),
				orchestratorStreams:  api.NewStreamsOf(api.WithLogger[orchestrator.Orchestrator_StoreAssessmentResultsClient, *orchestrator.StoreAssessmentResultRequest](log)),
				authz:                servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				streamToServer: createMockAssessmentServerStream(&assessment.AssessEvidenceRequest{
					Evidence: &evidence.Evidence{
						Id:             testdata.MockEvidenceID1,
						Timestamp:      timestamppb.Now(),
						ToolId:         testdata.MockEvidenceToolID1,
						CloudServiceId: testdata.MockCloudServiceID1,
						Resource:       toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: testdata.MockResourceID1, Type: []string{"VirtualMachine"}}}}, t)}}),
			},
			wantErr: false,
			wantRespMessage: &assessment.AssessEvidencesResponse{
				Status: assessment.AssessEvidencesResponse_ASSESSED,
			},
		},
		{
			name: "Error in stream to client - Send()-err",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				streamToClientWithSendErr: createMockAssessmentServerStreamWithSendErr(&assessment.AssessEvidenceRequest{
					Evidence: &evidence.Evidence{
						Timestamp:      timestamppb.Now(),
						ToolId:         testdata.MockEvidenceToolID1,
						CloudServiceId: testdata.MockCloudServiceID1,
						Resource:       toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: testdata.MockResourceID1, Type: []string{"VirtualMachine"}}}}, t)}}),
			},
			wantErr:        true,
			wantErrMessage: "rpc error: code = Unknown desc = cannot send response to the client",
		},
		{
			name: "Error in stream to server - Recv()-err",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				streamToServerWithRecvErr: createMockAssessmentServerStreamWithRecvErr(&assessment.AssessEvidenceRequest{
					Evidence: &evidence.Evidence{
						Timestamp:      timestamppb.Now(),
						ToolId:         testdata.MockEvidenceToolID1,
						CloudServiceId: testdata.MockCloudServiceID1,
						Resource:       toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: testdata.MockResourceID1, Type: []string{"VirtualMachine"}}}}, t)}}),
			},
			wantErr:        true,
			wantErrMessage: "rpc error: code = Unknown desc = cannot receive stream request",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				err                error
				responseFromServer *assessment.AssessEvidencesResponse
			)
			s := Service{
				resultHooks:          tt.fields.ResultHooks,
				cachedConfigurations: make(map[string]cachedConfiguration),
				evidenceStoreStreams: tt.fields.evidenceStoreStreams,
				evidenceStore:        api.NewRPCConnection("bufnet", evidence.NewEvidenceStoreClient, grpc.WithContextDialer(bufConnDialer)),
				orchestrator:         api.NewRPCConnection("bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(bufConnDialer)),
				orchestratorStreams:  tt.fields.orchestratorStreams,
				pe:                   policies.NewRegoEval(),
				authz:                tt.fields.authz,
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
				assert.Contains(t, responseFromServer.StatusMessage, tt.wantRespMessage.StatusMessage)
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
						Id:             testdata.MockEvidenceID1,
						ToolId:         testdata.MockEvidenceToolID1,
						Timestamp:      timestamppb.Now(),
						CloudServiceId: testdata.MockCloudServiceID1,
						Resource: toStruct(&voc.VirtualMachine{
							Compute: &voc.Compute{
								Resource: &voc.Resource{
									ID:   testdata.MockResourceID1,
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
			assert.Equal(t, hookCounts, hookCallCounter)
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
	SentFromServer chan *assessment.AssessEvidencesResponse
}

func (*mockAssessmentServerStream) CloseSend() error {
	panic("implement me")
}

func createMockAssessmentServerStream(r *assessment.AssessEvidenceRequest) *mockAssessmentServerStream {
	m := &mockAssessmentServerStream{
		RecvToServer: make(chan *assessment.AssessEvidenceRequest, 1),
	}
	m.RecvToServer <- r

	m.SentFromServer = make(chan *assessment.AssessEvidencesResponse, 1)
	return m
}

func (m *mockAssessmentServerStream) Send(response *assessment.AssessEvidencesResponse) error {
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
	return context.TODO()
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

	m.SentFromServer = make(chan *assessment.AssessEvidencesResponse, 1)
	return m
}

// mockAssessmentServerStreamWithSendErr implements Assessment_AssessEvidencesServer with error
type mockAssessmentServerStreamWithSendErr struct {
	grpc.ServerStream
	RecvToServer   chan *assessment.AssessEvidenceRequest
	SentFromServer chan *assessment.AssessEvidencesResponse
}

func (*mockAssessmentServerStreamWithSendErr) Send(*assessment.AssessEvidencesResponse) error {
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

func (*mockAssessmentServerStreamWithSendErr) Context() context.Context {
	return context.TODO()
}

type mockAssessmentServerStreamWithRecvErr struct {
	grpc.ServerStream
	RecvToServer   chan *assessment.AssessEvidenceRequest
	SentFromServer chan *assessment.AssessEvidencesResponse
}

func (*mockAssessmentServerStreamWithRecvErr) Send(*assessment.AssessEvidencesResponse) error {
	panic("implement me")
}

func (*mockAssessmentServerStreamWithRecvErr) Recv() (*assessment.AssessEvidenceRequest, error) {
	err := errors.New("Recv()-error")

	return nil, err
}

func (*mockAssessmentServerStreamWithRecvErr) Context() context.Context {
	return context.TODO()
}

func createMockAssessmentServerStreamWithRecvErr(r *assessment.AssessEvidenceRequest) *mockAssessmentServerStreamWithRecvErr {
	m := &mockAssessmentServerStreamWithRecvErr{
		RecvToServer: make(chan *assessment.AssessEvidenceRequest, 1),
	}
	m.RecvToServer <- r

	m.SentFromServer = make(chan *assessment.AssessEvidencesResponse, 1)
	return m
}

func TestService_HandleEvidence(t *testing.T) {
	type fields struct {
		authz         service.AuthorizationStrategy
		evidenceStore *api.RPCConnection[evidence.EvidenceStoreClient]
		orchestrator  *api.RPCConnection[orchestrator.OrchestratorClient]
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
				evidenceStore: api.NewRPCConnection("bufnet", evidence.NewEvidenceStoreClient, grpc.WithContextDialer(bufConnDialer)),
				orchestrator:  api.NewRPCConnection("bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(bufConnDialer)),
			},
			args: args{
				evidence: &evidence.Evidence{
					Id:             testdata.MockEvidenceID1,
					ToolId:         testdata.MockEvidenceToolID1,
					Timestamp:      timestamppb.Now(),
					CloudServiceId: testdata.MockCloudServiceID1,
					Resource: toStruct(voc.VirtualMachine{
						Compute: &voc.Compute{
							Resource: &voc.Resource{ID: testdata.MockResourceID1, Type: []string{"VirtualMachine"}},
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
				resourceId: testdata.MockResourceID1,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NoError(t, err)
				return false
			},
		},
		{
			name: "missing type in evidence",
			fields: fields{
				evidenceStore: api.NewRPCConnection("bufnet", evidence.NewEvidenceStoreClient, grpc.WithContextDialer(bufConnDialer)),
				orchestrator:  api.NewRPCConnection("bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(bufConnDialer)),
			},
			args: args{
				evidence: &evidence.Evidence{
					Id:             testdata.MockEvidenceID1,
					ToolId:         testdata.MockEvidenceToolID1,
					Timestamp:      timestamppb.Now(),
					CloudServiceId: testdata.MockCloudServiceID1,
					Resource:       toStruct(voc.VirtualMachine{Compute: &voc.Compute{Resource: &voc.Resource{ID: testdata.MockResourceID1, Type: []string{}}}}, t),
				},
				resourceId: testdata.MockResourceID1,
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
				evidenceStore: api.NewRPCConnection("bufnet", evidence.NewEvidenceStoreClient, grpc.WithContextDialer(connectionRefusedDialer)),
				orchestrator:  api.NewRPCConnection("bufnet", orchestrator.NewOrchestratorClient, grpc.WithContextDialer(bufConnDialer)),
			},
			args: args{
				evidence: &evidence.Evidence{
					Id:             testdata.MockEvidenceID1,
					ToolId:         testdata.MockEvidenceToolID1,
					Timestamp:      timestamppb.Now(),
					CloudServiceId: testdata.MockCloudServiceID1,
					Resource: toStruct(voc.VirtualMachine{
						Compute: &voc.Compute{
							Resource: &voc.Resource{ID: testdata.MockResourceID1, Type: []string{"VirtualMachine"}},
						},
						BootLogging: &voc.BootLogging{
							Logging: &voc.Logging{
								LoggingService:  nil,
								Enabled:         true,
								RetentionPeriod: 0,
							},
						}}, t),
				},
				resourceId: testdata.MockResourceID1,
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
			s := &Service{
				evidenceStore:        tt.fields.evidenceStore,
				orchestrator:         tt.fields.orchestrator,
				evidenceStoreStreams: api.NewStreamsOf(api.WithLogger[evidence.EvidenceStore_StoreEvidencesClient, *evidence.StoreEvidenceRequest](log)),
				orchestratorStreams:  api.NewStreamsOf(api.WithLogger[orchestrator.Orchestrator_StoreAssessmentResultsClient, *orchestrator.StoreAssessmentResultRequest](log)),
				cachedConfigurations: make(map[string]cachedConfiguration),
				pe:                   policies.NewRegoEval(policies.WithPackageName(policies.DefaultRegoPackage)),
				authz:                tt.fields.authz,
			}

			// Two tests: 1st) wantErr function. 2nd) if wantErr false then check if the results are valid
			results, err := s.handleEvidence(context.Background(), tt.args.evidence, tt.args.resourceId)
			if !tt.wantErr(t, err, fmt.Sprintf("handleEvidence(%v, %v)", tt.args.evidence, tt.args.resourceId)) {
				assert.NotEmpty(t, results)
				// Check the result by validation
				for _, result := range results {
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
			stream, err := s.initOrchestratorStream(tt.args.url, s.orchestrator.Opts...)

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
		evidenceStoreStreams *api.StreamsOf[evidence.EvidenceStore_StoreEvidencesClient, *evidence.StoreEvidenceRequest]
		orchestratorStreams  *api.StreamsOf[orchestrator.Orchestrator_StoreAssessmentResultsClient, *orchestrator.StoreAssessmentResultRequest]
		orchestrator         *api.RPCConnection[orchestrator.OrchestratorClient]
		evidenceStore        *api.RPCConnection[evidence.EvidenceStoreClient]
		metricEventStream    orchestrator.Orchestrator_SubscribeMetricChangeEventsClient
		resultHooks          []assessment.ResultHookFunc
		cachedConfigurations map[string]cachedConfiguration
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
				evidenceStore: api.NewRPCConnection(DefaultEvidenceStoreAddress, evidence.NewEvidenceStoreClient),
				orchestrator:  api.NewRPCConnection(DefaultOrchestratorAddress, orchestrator.NewOrchestratorClient),
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
				evidenceStore:        tt.fields.evidenceStore,
				orchestrator:         tt.fields.orchestrator,
				orchestratorStreams:  tt.fields.orchestratorStreams,
				metricEventStream:    tt.fields.metricEventStream,
				resultHooks:          tt.fields.resultHooks,
				cachedConfigurations: tt.fields.cachedConfigurations,
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
		isEvidenceStoreDisabled bool
		evidenceStoreStreams    *api.StreamsOf[evidence.EvidenceStore_StoreEvidencesClient, *evidence.StoreEvidenceRequest]
		evidenceStore           *api.RPCConnection[evidence.EvidenceStoreClient]
		orchestrator            *api.RPCConnection[orchestrator.OrchestratorClient]
		metricEventStream       orchestrator.Orchestrator_SubscribeMetricChangeEventsClient
		resultHooks             []assessment.ResultHookFunc
		cachedConfigurations    map[string]cachedConfiguration
		pe                      policies.PolicyEval
		evalPkg                 string
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
				isEvidenceStoreDisabled: tt.fields.isEvidenceStoreDisabled,
				evidenceStoreStreams:    tt.fields.evidenceStoreStreams,
				evidenceStore:           tt.fields.evidenceStore,
				orchestrator:            tt.fields.orchestrator,
				metricEventStream:       tt.fields.metricEventStream,
				resultHooks:             tt.fields.resultHooks,
				cachedConfigurations:    tt.fields.cachedConfigurations,
				pe:                      tt.fields.pe,
				evalPkg:                 tt.fields.evalPkg,
			}
			gotImpl, err := svc.MetricImplementation(tt.args.lang, tt.args.metric)
			tt.wantErr(t, err)
			assert.Equal(t, tt.wantImpl, gotImpl)
		})
	}
}
