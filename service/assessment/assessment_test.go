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
	"io"
	"os"
	"reflect"
	"runtime"
	"sync"
	"testing"
	"time"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/clitest"
	"clouditor.io/clouditor/v2/internal/testutil/prototest"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/evidencetest"
	"clouditor.io/clouditor/v2/launcher"
	"clouditor.io/clouditor/v2/policies"
	"clouditor.io/clouditor/v2/service"

	"github.com/google/uuid"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
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
		opts []service.Option[*Service]
	}
	tests := []struct {
		name string
		args args
		want assert.Want[*Service]
	}{
		{
			name: "AssessmentServer created with option rego package name",
			args: args{
				opts: []service.Option[*Service]{
					WithRegoPackageName("testPkg"),
				},
			},
			want: func(t *testing.T, got *Service) bool {
				return assert.Equal(t, "testPkg", got.evalPkg)
			},
		},
		{
			name: "AssessmentServer created with option authorizer",
			args: args{
				opts: []service.Option[*Service]{
					WithAuthorizer(api.NewOAuthAuthorizerFromClientCredentials(&clientcredentials.Config{})),
				},
			},
			want: func(t *testing.T, got *Service) bool {
				return assert.Equal(t, api.NewOAuthAuthorizerFromClientCredentials(&clientcredentials.Config{}), got.orchestrator.Authorizer(), assert.CompareAllUnexported())
			},
		},
		{
			name: "AssessmentServer created with options",
			args: args{
				opts: []service.Option[*Service]{
					WithEvidenceStoreAddress("localhost:9091"),
					WithOrchestratorAddress("localhost:9092"),
				},
			},
			want: func(t *testing.T, got *Service) bool {
				return assert.Equal(t, "localhost:9091", got.evidenceStore.Target) &&
					assert.Equal(t, "localhost:9092", got.orchestrator.Target)
			},
		},
		{
			name: "AssessmentServer without EvidenceStore",
			args: args{
				opts: []service.Option[*Service]{
					WithoutEvidenceStore(),
				},
			},
			want: func(t *testing.T, got *Service) bool {
				return assert.True(t, got.isEvidenceStoreDisabled)
			},
		},
		{
			name: "AssessmentServer with oauth2 authorizer",
			args: args{
				opts: []service.Option[*Service]{
					WithOAuth2Authorizer(&clientcredentials.Config{ClientID: "client"}),
				},
			},
			want: func(t *testing.T, got *Service) bool {
				return assert.NotNil(t, got.orchestrator.Authorizer())
			},
		},
		{
			name: "AssessmentServer with authorization strategy",
			args: args{
				opts: []service.Option[*Service]{
					WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(true)),
				},
			},
			want: func(t *testing.T, got *Service) bool {
				return assert.NotNil(t, got.authz.(*servicetest.AuthorizationStrategyMock))
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
		authz               service.AuthorizationStrategy
		evidenceStore       *api.RPCConnection[evidence.EvidenceStoreClient]
		orchestrator        *api.RPCConnection[orchestrator.OrchestratorClient]
		evidenceResourceMap map[string]*evidence.Evidence
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
				return assert.ErrorContains(t, err, "evidence: value is required")
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
				return assert.ErrorContains(t, err, "evidence.id: value is empty, which is not a valid UUID")
			},
		},
		{
			name: "Assess evidence without id",
			fields: fields{
				evidenceStore: api.NewRPCConnection(testdata.MockGRPCTarget, evidence.NewEvidenceStoreClient, grpc.WithContextDialer(bufConnDialer)),
				orchestrator:  api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(bufConnDialer)),
			},
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					ToolId:    testdata.MockEvidenceToolID1,
					Timestamp: timestamppb.Now(),
					Resource:  prototest.NewAny(t, &ontology.VirtualMachine{}),
				},
			},
			wantResp: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "evidence.id: value is empty, which is not a valid UUID")
			},
		},
		{
			name: "Assess resource without tool id",
			fields: fields{
				evidenceStore: api.NewRPCConnection(testdata.MockGRPCTarget, evidence.NewEvidenceStoreClient, grpc.WithContextDialer(bufConnDialer)),
				orchestrator:  api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(bufConnDialer)),
			},
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					Id:                    testdata.MockEvidenceID1,
					Timestamp:             timestamppb.Now(),
					CertificationTargetId: testdata.MockCertificationTargetID1,
					Resource:              prototest.NewAny(t, &ontology.VirtualMachine{}),
				},
			},
			wantResp: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "evidence.tool_id: value length must be at least 1 characters")
			},
		},
		{
			name: "Assess resource without timestamp",
			fields: fields{
				evidenceStore: api.NewRPCConnection(testdata.MockGRPCTarget, evidence.NewEvidenceStoreClient, grpc.WithContextDialer(bufConnDialer)),
				orchestrator:  api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(bufConnDialer)),
			},
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					Id:                    testdata.MockEvidenceID1,
					ToolId:                testdata.MockEvidenceToolID1,
					CertificationTargetId: testdata.MockCertificationTargetID1,
					Resource:              prototest.NewAny(t, &ontology.VirtualMachine{Id: testdata.MockResourceID1}),
				},
			},
			wantResp: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "evidence.timestamp: value is required")
			},
		},
		{
			name: "Assess resource",
			fields: fields{
				evidenceStore:       api.NewRPCConnection(testdata.MockGRPCTarget, evidence.NewEvidenceStoreClient, grpc.WithContextDialer(bufConnDialer)),
				orchestrator:        api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(bufConnDialer)),
				authz:               servicetest.NewAuthorizationStrategy(true),
				evidenceResourceMap: make(map[string]*evidence.Evidence),
			},
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					Id:        testdata.MockEvidenceID1,
					ToolId:    testdata.MockEvidenceToolID1,
					Timestamp: timestamppb.Now(),
					Resource: prototest.NewAny(t, &ontology.VirtualMachine{
						Id:   testdata.MockResourceID1,
						Name: testdata.MockResourceName1,
					}),
					CertificationTargetId: testdata.MockCertificationTargetID1},
			},
			wantResp: &assessment.AssessEvidenceResponse{
				Status: assessment.AssessmentStatus_ASSESSMENT_STATUS_ASSESSED,
			},
			wantErr: assert.NoError,
		},
		{
			name: "Assess resource of wrong could service",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockCertificationTargetID2),
			},
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					Id:        testdata.MockEvidenceID1,
					ToolId:    testdata.MockEvidenceToolID1,
					Timestamp: timestamppb.Now(),
					Resource: prototest.NewAny(t, &ontology.VirtualMachine{
						Id:   testdata.MockResourceID1,
						Name: testdata.MockResourceName1,
					}),
					CertificationTargetId: testdata.MockCertificationTargetID1},
			},
			wantResp: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, service.ErrPermissionDenied.Error())
			},
		},
		{
			name: "Assess resource without resource id",
			fields: fields{
				evidenceStore:       api.NewRPCConnection(testdata.MockGRPCTarget, evidence.NewEvidenceStoreClient, grpc.WithContextDialer(bufConnDialer)),
				orchestrator:        api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(bufConnDialer)),
				authz:               servicetest.NewAuthorizationStrategy(true),
				evidenceResourceMap: make(map[string]*evidence.Evidence),
			},
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					Id:                    testdata.MockEvidenceID1,
					ToolId:                testdata.MockEvidenceToolID1,
					Timestamp:             timestamppb.Now(),
					Resource:              prototest.NewAny(t, &ontology.VirtualMachine{}),
					CertificationTargetId: testdata.MockCertificationTargetID1,
				},
			},
			wantResp: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "id: value is required")
			},
		},
		{
			name: "No RPC connections",
			fields: fields{
				evidenceStore:       api.NewRPCConnection(testdata.MockGRPCTarget, evidence.NewEvidenceStoreClient, grpc.WithContextDialer(connectionRefusedDialer)),
				orchestrator:        api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(connectionRefusedDialer)),
				authz:               servicetest.NewAuthorizationStrategy(true),
				evidenceResourceMap: make(map[string]*evidence.Evidence),
			},
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					Id:                    testdata.MockEvidenceID1,
					ToolId:                testdata.MockEvidenceToolID1,
					Timestamp:             timestamppb.Now(),
					CertificationTargetId: testdata.MockCertificationTargetID1,
					Resource: prototest.NewAny(t, &ontology.VirtualMachine{
						Id:   testdata.MockResourceID1,
						Name: testdata.MockResourceName1,
					}),
				},
			},
			wantResp: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "connection refused")
			},
		},
		{
			name: "Assess resource and wait existing related resources is already there",
			fields: fields{
				evidenceStore: api.NewRPCConnection(testdata.MockGRPCTarget, evidence.NewEvidenceStoreClient, grpc.WithContextDialer(bufConnDialer)),
				orchestrator:  api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(bufConnDialer)),
				authz:         servicetest.NewAuthorizationStrategy(true),
				evidenceResourceMap: map[string]*evidence.Evidence{
					"my-other-resource-id": {
						Id: testdata.MockEvidenceID2,
						Resource: prototest.NewAny(t, &ontology.VirtualMachine{
							Id: testdata.MockResourceID2,
						}),
					},
				},
			},
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					Id:                    testdata.MockEvidenceID1,
					ToolId:                testdata.MockEvidenceToolID1,
					Timestamp:             timestamppb.Now(),
					CertificationTargetId: testdata.MockCertificationTargetID1,
					Resource: prototest.NewAny(t, &ontology.VirtualMachine{
						Id:   testdata.MockResourceID1,
						Name: testdata.MockResourceName1,
					}),
					ExperimentalRelatedResourceIds: []string{"my-other-resource-id"},
				},
			},
			wantResp: &assessment.AssessEvidenceResponse{
				Status: assessment.AssessmentStatus_ASSESSMENT_STATUS_ASSESSED,
			},
			wantErr: assert.NoError,
		},
		{
			name: "Assess resource and wait existing related resources is not there",
			fields: fields{
				evidenceStore:       api.NewRPCConnection(testdata.MockGRPCTarget, evidence.NewEvidenceStoreClient, grpc.WithContextDialer(bufConnDialer)),
				orchestrator:        api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(bufConnDialer)),
				authz:               servicetest.NewAuthorizationStrategy(true),
				evidenceResourceMap: make(map[string]*evidence.Evidence),
			},
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					Id:                    testdata.MockEvidenceID1,
					ToolId:                testdata.MockEvidenceToolID1,
					Timestamp:             timestamppb.Now(),
					CertificationTargetId: testdata.MockCertificationTargetID1,
					Resource: prototest.NewAny(t, &ontology.VirtualMachine{
						Id:   testdata.MockResourceID1,
						Name: testdata.MockResourceName1,
					}),
					ExperimentalRelatedResourceIds: []string{"my-other-resource-id"},
				},
			},
			wantResp: &assessment.AssessEvidenceResponse{
				Status: assessment.AssessmentStatus_ASSESSMENT_STATUS_WAITING_FOR_RELATED,
			},
			wantErr: assert.NoError,
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
				evidenceResourceMap:  tt.fields.evidenceResourceMap,
				requests:             make(map[string]waitingRequest),
				pe:                   policies.NewRegoEval(policies.WithPackageName(policies.DefaultRegoPackage)),
				authz:                tt.fields.authz,
			}
			gotResp, err := s.AssessEvidence(tt.args.in0, &assessment.AssessEvidenceRequest{Evidence: tt.args.evidence})

			tt.wantErr(t, err)
			// Check response
			assert.Equal(t, tt.wantResp, gotResp)
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
			testdata.MockGRPCTarget, evidence.NewEvidenceStoreClient, grpc.WithContextDialer(bufConnDialer)),
		orchestrator: api.NewRPCConnection(
			testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(bufConnDialer)),
		authz: servicetest.NewAuthorizationStrategy(true),
		evidenceStoreStreams: api.NewStreamsOf(
			api.WithLogger[evidence.EvidenceStore_StoreEvidencesClient, *evidence.StoreEvidenceRequest](log)),
		orchestratorStreams: api.NewStreamsOf(
			api.WithLogger[orchestrator.Orchestrator_StoreAssessmentResultsClient,
				*orchestrator.StoreAssessmentResultRequest](log)),
		cachedConfigurations: make(map[string]cachedConfiguration),
		evidenceResourceMap:  make(map[string]*evidence.Evidence),
		pe:                   policies.NewRegoEval(policies.WithPackageName(policies.DefaultRegoPackage)),
	}
	// First assess evidence with a valid VM resource s.t. the cache is created for the combination of resource type and
	// tool id (="VirtualMachine-{testdata.MockEvidenceToolID}")
	e := evidencetest.MockEvidence1
	e.Resource = prototest.NewAny(t, &ontology.VirtualMachine{
		Id:   testdata.MockResourceID1,
		Name: testdata.MockResourceName1,
	})
	_, err := s.AssessEvidence(context.Background(), &assessment.AssessEvidenceRequest{Evidence: e})
	assert.NoError(t, err)

	// Now assess a new evidence which has not a valid format other than the resource type and tool id is set correctly
	// Prepare resource. Make sure both evidences have the same type (for caching key)
	a := prototest.NewAny(t, &ontology.VirtualMachine{
		Id:   uuid.NewString(),
		Name: "Some other name",
	})

	assert.NoError(t, err)
	_, err = s.AssessEvidence(context.Background(), &assessment.AssessEvidenceRequest{Evidence: &evidence.Evidence{
		Id:                    uuid.NewString(),
		Timestamp:             timestamppb.Now(),
		CertificationTargetId: testdata.MockCertificationTargetID1,
		// Make sure both evidences have the same tool id (for caching key)
		ToolId:   e.ToolId,
		Raw:      nil,
		Resource: a,
	}})
	assert.NoError(t, err)
}

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
		name    string
		fields  fields
		args    args
		wantErr assert.WantErr
		want    assert.Want[*assessment.AssessEvidencesResponse]
	}{
		{
			name: "Missing toolId",
			args: args{
				streamToServer: createMockAssessmentServerStream(&assessment.AssessEvidenceRequest{
					Evidence: &evidence.Evidence{
						Id:                    testdata.MockEvidenceID1,
						Timestamp:             timestamppb.Now(),
						CertificationTargetId: testdata.MockCertificationTargetID1,
						Resource:              prototest.NewAny(t, &ontology.VirtualMachine{Id: testdata.MockResourceID1}),
					},
				}),
			},
			wantErr: assert.Nil[error],
			want: func(t *testing.T, got *assessment.AssessEvidencesResponse) bool {
				assert.Equal(t, assessment.AssessmentStatus_ASSESSMENT_STATUS_FAILED, got.Status)
				return assert.Contains(t, got.StatusMessage, "evidence.tool_id: value length must be at least 1 characters")
			},
		},
		{
			name: "Missing evidenceID",
			args: args{
				streamToServer: createMockAssessmentServerStream(&assessment.AssessEvidenceRequest{
					Evidence: &evidence.Evidence{
						Timestamp:             timestamppb.Now(),
						ToolId:                testdata.MockEvidenceToolID1,
						CertificationTargetId: testdata.MockCertificationTargetID1,
						Resource:              prototest.NewAny(t, &ontology.VirtualMachine{Id: testdata.MockResourceID1}),
					},
				}),
			},
			wantErr: assert.Nil[error],
			want: func(t *testing.T, got *assessment.AssessEvidencesResponse) bool {
				assert.Equal(t, assessment.AssessmentStatus_ASSESSMENT_STATUS_FAILED, got.Status)
				return assert.Contains(t, got.StatusMessage, "evidence.id: value is empty, which is not a valid UUID")
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
						Id:                    testdata.MockEvidenceID1,
						Timestamp:             timestamppb.Now(),
						ToolId:                testdata.MockEvidenceToolID1,
						CertificationTargetId: testdata.MockCertificationTargetID1,
						Resource: prototest.NewAny(t, &ontology.VirtualMachine{
							Id:   testdata.MockResourceID1,
							Name: testdata.MockResourceName1,
						}),
					},
				}),
			},
			wantErr: assert.Nil[error],
			want: func(t *testing.T, got *assessment.AssessEvidencesResponse) bool {
				assert.Equal(t, assessment.AssessmentStatus_ASSESSMENT_STATUS_ASSESSED, got.Status)
				return assert.Empty(t, got.StatusMessage)
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
						Timestamp:             timestamppb.Now(),
						ToolId:                testdata.MockEvidenceToolID1,
						CertificationTargetId: testdata.MockCertificationTargetID1,
						Resource:              prototest.NewAny(t, &ontology.VirtualMachine{Id: testdata.MockResourceID1}),
					},
				}),
			},
			want: assert.Nil[*assessment.AssessEvidencesResponse],
			wantErr: func(t *testing.T, err error) bool {
				return assert.ErrorContains(t, err, "rpc error: code = Unknown desc = cannot send response to the client")
			},
		},
		{
			name: "Error in stream to server - Recv()-err",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				streamToServerWithRecvErr: createMockAssessmentServerStreamWithRecvErr(&assessment.AssessEvidenceRequest{
					Evidence: &evidence.Evidence{
						Timestamp:             timestamppb.Now(),
						ToolId:                testdata.MockEvidenceToolID1,
						CertificationTargetId: testdata.MockCertificationTargetID1,
						Resource:              prototest.NewAny(t, &ontology.VirtualMachine{Id: testdata.MockResourceID1}),
					},
				}),
			},
			want: assert.Nil[*assessment.AssessEvidencesResponse],
			wantErr: func(t *testing.T, err error) bool {
				return assert.ErrorContains(t, err, "rpc error: code = Unknown desc = cannot receive stream request")
			},
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
				evidenceStore:        api.NewRPCConnection(testdata.MockGRPCTarget, evidence.NewEvidenceStoreClient, grpc.WithContextDialer(bufConnDialer)),
				orchestrator:         api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(bufConnDialer)),
				orchestratorStreams:  tt.fields.orchestratorStreams,
				evidenceResourceMap:  make(map[string]*evidence.Evidence),
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

			tt.wantErr(t, err)
			tt.want(t, responseFromServer)
		})
	}
}

func TestService_AssessmentResultHooks(t *testing.T) {
	var (
		hookCallCounter = 0
		wg              sync.WaitGroup
		hookCounts      = 24
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
		wantErr  assert.WantErr
	}{
		{
			name: "Store evidence to the map",
			args: args{
				in0: context.TODO(),
				evidence: &assessment.AssessEvidenceRequest{
					Evidence: &evidence.Evidence{
						Id:                    testdata.MockEvidenceID1,
						ToolId:                testdata.MockEvidenceToolID1,
						Timestamp:             timestamppb.Now(),
						CertificationTargetId: testdata.MockCertificationTargetID1,
						Resource: prototest.NewAny(t, &ontology.VirtualMachine{
							Id:   testdata.MockResourceID1,
							Name: testdata.MockResourceName1,
							BootLogging: &ontology.BootLogging{
								LoggingServiceIds: []string{"SomeResourceId2"},
								Enabled:           true,
								RetentionPeriod:   durationpb.New(time.Hour * 24 * 36),
							},
							OsLogging: &ontology.OSLogging{
								LoggingServiceIds: []string{"SomeResourceId2"},
								Enabled:           true,
								RetentionPeriod:   durationpb.New(time.Hour * 24 * 36),
							},
							MalwareProtection: &ontology.MalwareProtection{
								Enabled:              true,
								NumberOfThreatsFound: 5,
								DurationSinceActive:  durationpb.New(time.Hour * 24 * 20),
								ApplicationLogging: &ontology.ApplicationLogging{
									Enabled:           true,
									LoggingServiceIds: []string{"SomeAnalyticsService?"},
								},
							},
						}),
					}},

				resultHooks: []assessment.ResultHookFunc{firstHookFunction, secondHookFunction},
			},
			wantErr: assert.Nil[error],
			wantResp: &assessment.AssessEvidenceResponse{
				Status: assessment.AssessmentStatus_ASSESSMENT_STATUS_ASSESSED,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hookCallCounter = 0
			s := NewService(
				WithEvidenceStoreAddress(testdata.MockGRPCTarget, grpc.WithContextDialer(bufConnDialer)),
				WithOrchestratorAddress(testdata.MockGRPCTarget, grpc.WithContextDialer(bufConnDialer)),
			)

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

			tt.wantErr(t, err)
			assert.Equal(t, tt.wantResp, gotResp)
			assert.Equal(t, hookCounts, hookCallCounter)
		})
	}
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

func TestService_handleEvidence(t *testing.T) {
	type fields struct {
		authz         service.AuthorizationStrategy
		evidenceStore *api.RPCConnection[evidence.EvidenceStoreClient]
		orchestrator  *api.RPCConnection[orchestrator.OrchestratorClient]
	}
	type args struct {
		evidence *evidence.Evidence
		related  map[string]ontology.IsResource
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.Want[[]*assessment.AssessmentResult]
		wantErr assert.WantErr
	}{
		{
			name: "correct evidence",
			fields: fields{
				evidenceStore: api.NewRPCConnection(testdata.MockGRPCTarget, evidence.NewEvidenceStoreClient, grpc.WithContextDialer(bufConnDialer)),
				orchestrator:  api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(bufConnDialer)),
			},
			args: args{
				evidence: &evidence.Evidence{
					Id:                    testdata.MockEvidenceID1,
					ToolId:                testdata.MockEvidenceToolID1,
					Timestamp:             timestamppb.Now(),
					CertificationTargetId: testdata.MockCertificationTargetID1,
					Resource: prototest.NewAny(t, &ontology.VirtualMachine{
						Id:   testdata.MockResourceID1,
						Name: testdata.MockResourceName1,
						BootLogging: &ontology.BootLogging{
							LoggingServiceIds: nil,
							Enabled:           true,
						},
					}),
				},
			},
			want: func(t *testing.T, got []*assessment.AssessmentResult) bool {
				for _, result := range got {
					err := api.Validate(result)
					assert.NoError(t, err)
				}
				return assert.Equal(t, 11, len(got))
			},
			wantErr: assert.Nil[error],
		},
		{
			name: "broken Any message",
			fields: fields{
				evidenceStore: api.NewRPCConnection(testdata.MockGRPCTarget, evidence.NewEvidenceStoreClient, grpc.WithContextDialer(bufConnDialer)),
				orchestrator:  api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(bufConnDialer)),
			},
			args: args{
				evidence: &evidence.Evidence{
					Id:                    testdata.MockEvidenceID1,
					ToolId:                testdata.MockEvidenceToolID1,
					Timestamp:             timestamppb.Now(),
					CertificationTargetId: testdata.MockCertificationTargetID1,
					Resource:              &anypb.Any{TypeUrl: "does-not-exist"},
				},
			},
			want: assert.Nil[[]*assessment.AssessmentResult],
			wantErr: func(t *testing.T, err error) bool {
				return assert.Contains(t, err.Error(), "could not unmarshal resource proto message")
			},
		},
		{
			name: "not an ontology resource",
			fields: fields{
				evidenceStore: api.NewRPCConnection(testdata.MockGRPCTarget, evidence.NewEvidenceStoreClient, grpc.WithContextDialer(bufConnDialer)),
				orchestrator:  api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(bufConnDialer)),
			},
			args: args{
				evidence: &evidence.Evidence{
					Id:                    testdata.MockEvidenceID1,
					ToolId:                testdata.MockEvidenceToolID1,
					Timestamp:             timestamppb.Now(),
					CertificationTargetId: testdata.MockCertificationTargetID1,
					Resource:              prototest.NewAny(t, &emptypb.Empty{}),
				},
			},
			want: assert.Nil[[]*assessment.AssessmentResult],
			wantErr: func(t *testing.T, err error) bool {
				return assert.Contains(t, err.Error(), discovery.ErrNotOntologyResource.Error())
			},
		},
		{
			name: "evidence store stream error",
			fields: fields{
				evidenceStore: api.NewRPCConnection(testdata.MockGRPCTarget, evidence.NewEvidenceStoreClient, grpc.WithContextDialer(connectionRefusedDialer)),
				orchestrator:  api.NewRPCConnection(testdata.MockGRPCTarget, orchestrator.NewOrchestratorClient, grpc.WithContextDialer(bufConnDialer)),
			},
			args: args{
				evidence: &evidence.Evidence{
					Id:                    testdata.MockEvidenceID1,
					ToolId:                testdata.MockEvidenceToolID1,
					Timestamp:             timestamppb.Now(),
					CertificationTargetId: testdata.MockCertificationTargetID1,
					Resource: prototest.NewAny(t, &ontology.VirtualMachine{
						Id:   testdata.MockResourceID1,
						Name: testdata.MockResourceName1,
						BootLogging: &ontology.BootLogging{
							LoggingServiceIds: nil,
							Enabled:           true,
						}}),
				},
			},
			want: assert.Nil[[]*assessment.AssessmentResult],
			wantErr: func(t *testing.T, err error) bool {
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

			results, err := s.handleEvidence(context.Background(), tt.args.evidence, tt.args.related)

			tt.wantErr(t, err)
			tt.want(t, results)
		})
	}
}

func TestService_initOrchestratorStoreStream(t *testing.T) {
	type fields struct {
		opts []service.Option[*Service]
	}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.Want[orchestrator.Orchestrator_StoreAssessmentResultsClient]
		wantErr assert.WantErr
	}{
		{
			name: "Invalid RPC connection",
			args: args{
				url: "localhost:1",
			},
			fields: fields{
				opts: []service.Option[*Service]{
					WithOrchestratorAddress("localhost:1"),
				},
			},
			want: assert.Empty[orchestrator.Orchestrator_StoreAssessmentResultsClient],
			wantErr: func(t *testing.T, err error) bool {
				s, _ := status.FromError(errors.Unwrap(err))
				return assert.Equal(t, codes.Unavailable, s.Code())
			},
		},
		{
			name: "Authenticated RPC connection with invalid user",
			args: args{
				url: testdata.MockGRPCTarget,
			},
			fields: fields{
				opts: []service.Option[*Service]{
					WithOrchestratorAddress(testdata.MockGRPCTarget, grpc.WithContextDialer(bufConnDialer)),
					WithOAuth2Authorizer(testutil.AuthClientConfig(authPort)),
				},
			},
			want: assert.Empty[orchestrator.Orchestrator_StoreAssessmentResultsClient],
			wantErr: func(t *testing.T, err error) bool {
				s, _ := status.FromError(errors.Unwrap(err))
				return assert.Equal(t, codes.Unauthenticated, s.Code())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.fields.opts...)
			stream, err := s.initOrchestratorStream(tt.args.url, s.orchestrator.Opts...)

			tt.wantErr(t, err)
			tt.want(t, stream)
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

			assert.Equal(t, tt.wantEvent, rec.event)
		})
	}
}

type eventRecorder struct {
	event *orchestrator.MetricChangeEvent
	done  bool
}

func (*eventRecorder) Eval(_ *evidence.Evidence, _ ontology.IsResource, _ map[string]ontology.IsResource, _ policies.MetricsSource) (data []*policies.Result, err error) {
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
		metric *assessment.Metric
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.Want[*assessment.MetricImplementation]
		wantErr assert.WantErr
	}{

		{
			name: "Unspecified language",
			args: args{
				lang: assessment.MetricImplementation_LANGUAGE_UNSPECIFIED,
			},
			want: assert.Nil[*assessment.MetricImplementation],
			wantErr: func(t *testing.T, err error) bool {
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
			tt.want(t, gotImpl)
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
