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
	"os"
	"reflect"
	"runtime"
	"sync"
	"testing"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/clitest"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/evidencetest"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/orchestratortest"
	"clouditor.io/clouditor/v2/internal/util"
	"clouditor.io/clouditor/v2/launcher"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/persistence/gorm"
	"clouditor.io/clouditor/v2/service"
	"github.com/google/uuid"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	gormio "gorm.io/gorm"
)

func TestMain(m *testing.M) {
	clitest.AutoChdir()

	server, _ := startBufConnServer()

	code := m.Run()

	server.Stop()
	os.Exit(code)
}

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
			name: "EvidenceStoreServer created with option 'WithStorage'",
			args: args{opts: []service.Option[*Service]{WithStorage(db)}},
			want: func(t *testing.T, got *Service) bool {
				// Storage should be gorm (in-memory storage). Hard to check since its type is not exported
				assert.NotNil(t, got.storage)
				return true
			},
		},
		{
			name: "EvidenceStoreServer created with option 'WithAssessmentAddress'",
			args: args{opts: []service.Option[*Service]{WithAssessmentAddress("localhost:9091")}},
			want: func(t *testing.T, got *Service) bool {
				return assert.Equal(t, "localhost:9091", got.assessment.Target)
			},
		},
		{
			name: "EvidenceStoreServer created with option 'WithOAuth2Authorizer'",
			args: args{opts: []service.Option[*Service]{WithOAuth2Authorizer(&clientcredentials.Config{ClientID: "client"})}},
			want: func(t *testing.T, got *Service) bool {
				return assert.NotNil(t, got.assessment.Authorizer())
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

// TestStoreEvidence tests StoreEvidence
func TestService_StoreEvidence(t *testing.T) {
	// mock assessment stream
	mockStream := &mockAssessmentStream{connectionEstablished: true, expected: 2}
	mockStream.Prepare()

	type args struct {
		in0       context.Context
		req       *evidence.StoreEvidenceRequest
		addStream bool
	}
	type fields struct {
		storage persistence.Storage
		authz   service.AuthorizationStrategy
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		wantRes assert.Want[*evidence.StoreEvidenceResponse]
		want    assert.Want[*Service]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Error: invalid evidence",
			args: args{
				in0: context.TODO(),
				req: &evidence.StoreEvidenceRequest{
					Evidence: &evidence.Evidence{
						Id:                   testdata.MockEvidenceID1,
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
						Timestamp:            timestamppb.Now(),
						Resource: &ontology.Resource{
							Type: &ontology.Resource_VirtualMachine{
								VirtualMachine: &ontology.VirtualMachine{
									Id: "mock-id",
								},
							},
						},
					},
				},
			},
			wantRes: assert.Nil[*evidence.StoreEvidenceResponse],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "evidence.tool_id: value length must be at least 1 characters")
			},
			want: func(t *testing.T, s *Service) bool {
				return assert.NotEmpty(t, s)
			},
		},
		{
			name: "Error: authorization denied",
			args: args{
				in0: context.TODO(),
				req: &evidence.StoreEvidenceRequest{
					Evidence: evidencetest.MockEvidence1,
				},
			},
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID2),
			},
			wantRes: assert.Nil[*evidence.StoreEvidenceResponse],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.PermissionDenied, status.Code(err))
				return assert.ErrorContains(t, err, "access denied")
			},
		},
		{
			name: "Error: storage error (database error)",
			args: args{
				in0: context.TODO(),
				req: &evidence.StoreEvidenceRequest{
					Evidence: evidencetest.MockEvidence1,
				},
			},
			fields: fields{
				authz:   servicetest.NewAuthorizationStrategy(true, testdata.MockTargetOfEvaluationID1),
				storage: &testutil.StorageWithError{CreateErr: gormio.ErrInvalidDB},
			},
			wantRes: assert.Nil[*evidence.StoreEvidenceResponse],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.Internal, status.Code(err))
				return assert.ErrorContains(t, err, persistence.ErrDatabase.Error())
			},
			want: func(t *testing.T, s *Service) bool {
				return assert.NotEmpty(t, s)
			},
		},
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
						Resource: &ontology.Resource{
							Type: &ontology.Resource_VirtualMachine{
								VirtualMachine: &ontology.VirtualMachine{
									Id:   "mock-id",
									Name: "my-vm",
								},
							},
						},
					}},
				addStream: true,
			},
			fields: fields{
				authz:   servicetest.NewAuthorizationStrategy(true, testdata.MockTargetOfEvaluationID1),
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {}),
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
			wantErr: assert.NoError,
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
						Resource: &ontology.Resource{
							Type: &ontology.Resource_VirtualMachine{
								VirtualMachine: &ontology.VirtualMachine{
									Id:   "mock-id",
									Name: "mock-name",
								},
							},
						},
					},
				},
				addStream: true,
			},
			fields: fields{
				authz:   servicetest.NewAuthorizationStrategy(true, testdata.MockTargetOfEvaluationID1),
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {}),
			},
			wantRes: assert.Nil[*evidence.StoreEvidenceResponse],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "evidence.tool_id: value length must be at least 1 characters")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create service with assessment stream
			// The StoreEvidence method will use the assessment stream to send the evidence to the assessment service
			svc := NewService()
			svc.storage = tt.fields.storage
			svc.authz = tt.fields.authz
			svc.assessment = nil
			svc.assessmentStreams = nil

			// Add assessment stream if needed
			if tt.args.addStream {
				svc.assessment = &api.RPCConnection[assessment.AssessmentClient]{Target: "mock"}
				svc.assessmentStreams = api.NewStreamsOf[assessment.Assessment_AssessEvidencesClient, *assessment.AssessEvidenceRequest]()
				_, _ = svc.assessmentStreams.GetStream("mock", "Assessment", func(target string, additionalOpts ...grpc.DialOption) (stream assessment.Assessment_AssessEvidencesClient, err error) {
					return mockStream, nil
				})
			}

			gotRes, err := svc.StoreEvidence(tt.args.in0, tt.args.req)

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
		streamToServer            *mockEvidenceStoreStream
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
				Status: evidence.EvidenceStatus_EVIDENCE_STATUS_OK,
			},
		},
		{
			name:   "Store 2 evidences to the map",
			fields: fields{count: 2},
			args: args{
				streamToServer: createMockStream(createStoreEvidenceRequestMocks(t, 2))},
			wantErr: false,
			wantResMessage: &evidence.StoreEvidencesResponse{
				Status: evidence.EvidenceStatus_EVIDENCE_STATUS_OK,
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
							Resource: &ontology.Resource{
								Type: &ontology.Resource_VirtualMachine{
									VirtualMachine: &ontology.VirtualMachine{
										Id:   "mock-id-1",
										Name: "mock-name-1",
									},
								},
							},
						},
					},
				})},
			wantErr: false,
			wantResMessage: &evidence.StoreEvidencesResponse{
				Status:        evidence.EvidenceStatus_EVIDENCE_STATUS_ERROR,
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

			// mock assessment stream
			mockStream := &mockAssessmentStream{connectionEstablished: true, expected: 2}
			mockStream.Prepare()

			// create service with assessment stream
			// StoreEvidences sends the evidence to the StoreEvidence method which will use the assessment stream to send the evidence to the assessment service
			svc := NewService()
			svc.assessment = &api.RPCConnection[assessment.AssessmentClient]{Target: "mock"}
			svc.assessmentStreams = api.NewStreamsOf[assessment.Assessment_AssessEvidencesClient, *assessment.AssessEvidenceRequest]()
			_, _ = svc.assessmentStreams.GetStream("mock", "Assessment", func(target string, additionalOpts ...grpc.DialOption) (stream assessment.Assessment_AssessEvidencesClient, err error) {
				return mockStream, nil
			})

			if tt.args.streamToServer != nil {
				err = svc.StoreEvidences(tt.args.streamToServer)
				responseFromServer = <-tt.args.streamToServer.SentFromServer
			} else if tt.args.streamToClientWithSendErr != nil {
				err = svc.StoreEvidences(tt.args.streamToClientWithSendErr)
			} else if tt.args.streamToServerWithRecvErr != nil {
				err = svc.StoreEvidences(tt.args.streamToServerWithRecvErr)
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

	// add 2 hock functions
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

	// mock assessment stream
	mockStream := &mockAssessmentStream{connectionEstablished: true, expected: 2}
	mockStream.Prepare()

	// create service
	svc := NewService()
	svc.assessment = &api.RPCConnection[assessment.AssessmentClient]{Target: "mock"}
	svc.assessmentStreams = api.NewStreamsOf[assessment.Assessment_AssessEvidencesClient, *assessment.AssessEvidenceRequest]()
	_, _ = svc.assessmentStreams.GetStream("mock", "Assessment", func(target string, additionalOpts ...grpc.DialOption) (stream assessment.Assessment_AssessEvidencesClient, err error) {
		return mockStream, nil
	})
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
					Resource: &ontology.Resource{
						Type: &ontology.Resource_VirtualMachine{
							VirtualMachine: &ontology.VirtualMachine{
								Id:   "mock-id-1",
								Name: "mock-name-1",
							},
						},
					},
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
		{
			name: "DB error (database error)",
			fields: fields{
				storage: &testutil.StorageWithError{GetErr: gormio.ErrInvalidDB},
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &evidence.GetEvidenceRequest{
					EvidenceId: testdata.MockEvidenceID1,
				},
			},
			want: assert.Nil[*evidence.Evidence],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, persistence.ErrDatabase.Error())
			},
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

// createStoreEvidenceRequestMocks creates store evidence requests with random evidence IDs
func createStoreEvidenceRequestMocks(_ *testing.T, count int) []*evidence.StoreEvidenceRequest {
	var mockRequests []*evidence.StoreEvidenceRequest

	for i := 0; i < count; i++ {
		evidenceRequest := &evidence.StoreEvidenceRequest{
			Evidence: &evidence.Evidence{
				Id:                   uuid.NewString(),
				ToolId:               fmt.Sprintf("MockToolId-%d", i),
				TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
				Timestamp:            timestamppb.Now(),
				Resource: &ontology.Resource{
					Type: &ontology.Resource_VirtualMachine{
						VirtualMachine: &ontology.VirtualMachine{
							Id:   "mock-id-1",
							Name: "my-vm",
						},
					},
				},
			},
		}
		mockRequests = append(mockRequests, evidenceRequest)
	}

	return mockRequests
}

// mockEvidenceStoreStream implements Evidence_StoreEvidenceClient interface
// and is used to mock the server stream for testing purposes.
type mockEvidenceStoreStream struct {
	grpc.ServerStream
	RecvToServer   chan *evidence.StoreEvidenceRequest
	SentFromServer chan *evidence.StoreEvidencesResponse
}

func createMockStream(requests []*evidence.StoreEvidenceRequest) *mockEvidenceStoreStream {
	m := &mockEvidenceStoreStream{
		RecvToServer: make(chan *evidence.StoreEvidenceRequest, len(requests)),
	}
	for _, req := range requests {
		m.RecvToServer <- req
	}

	m.SentFromServer = make(chan *evidence.StoreEvidencesResponse, len(requests))
	return m
}

func (m *mockEvidenceStoreStream) Send(response *evidence.StoreEvidencesResponse) error {
	m.SentFromServer <- response
	return nil
}

func (*mockEvidenceStoreStream) SendAndClose(_ *emptypb.Empty) error {
	return nil
}

func (m *mockEvidenceStoreStream) Recv() (req *evidence.StoreEvidenceRequest, err error) {
	if len(m.RecvToServer) == 0 {
		return nil, io.EOF
	}
	req, more := <-m.RecvToServer
	if !more {
		return nil, errors.New("empty")
	}

	return req, nil
}

func (*mockEvidenceStoreStream) SetHeader(_ metadata.MD) error {
	panic("implement me")
}

func (*mockEvidenceStoreStream) SendHeader(_ metadata.MD) error {
	panic("implement me")
}

func (*mockEvidenceStoreStream) SetTrailer(_ metadata.MD) {
	panic("implement me")
}

func (*mockEvidenceStoreStream) Context() context.Context {
	return context.TODO()
}

func (*mockEvidenceStoreStream) SendMsg(_ interface{}) error {
	panic("implement me")
}

func (*mockEvidenceStoreStream) RecvMsg(_ interface{}) error {
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

// mockAssessmentStream implements Assessment_AssessEvidencesClient interface
type mockAssessmentStream struct {
	// We add sentEvidence field to test the evidence that would be sent over gRPC
	sentEvidences []*evidence.Evidence
	// We add connectionEstablished to differentiate between the case where evidences can be sent and not
	connectionEstablished bool
	counter               int
	expected              int
	wg                    sync.WaitGroup
}

func (m *mockAssessmentStream) Prepare() {
	m.wg.Add(m.expected)
}

func (m *mockAssessmentStream) Wait() {
	m.wg.Wait()
}

func (m *mockAssessmentStream) Recv() (*assessment.AssessEvidencesResponse, error) {
	if m.counter == 0 {
		m.counter++
		return &assessment.AssessEvidencesResponse{
			Status:        assessment.AssessmentStatus_ASSESSMENT_STATUS_FAILED,
			StatusMessage: "mockError1",
		}, nil
	} else if m.counter == 1 {
		m.counter++
		return &assessment.AssessEvidencesResponse{
			Status: assessment.AssessmentStatus_ASSESSMENT_STATUS_ASSESSED,
		}, nil
	} else {
		return nil, io.EOF
	}
}

func (m *mockAssessmentStream) Send(req *assessment.AssessEvidenceRequest) (err error) {
	return m.SendMsg(req)
}

func (*mockAssessmentStream) CloseAndRecv() (*emptypb.Empty, error) {
	return nil, nil
}

func (*mockAssessmentStream) Header() (metadata.MD, error) {
	return nil, nil
}

func (*mockAssessmentStream) Trailer() metadata.MD {
	return nil
}

func (*mockAssessmentStream) CloseSend() error {
	return nil
}

func (*mockAssessmentStream) Context() context.Context {
	return nil
}

func (m *mockAssessmentStream) SendMsg(req interface{}) (err error) {
	e := req.(*assessment.AssessEvidenceRequest).Evidence
	if m.connectionEstablished {
		m.sentEvidences = append(m.sentEvidences, e)
	} else {
		err = fmt.Errorf("mock send error")
	}
	m.wg.Done()
	return
}

func (*mockAssessmentStream) RecvMsg(_ interface{}) error {
	return nil
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

func TestService_initAssessmentStream(t *testing.T) {
	type fields struct {
		assessment *api.RPCConnection[assessment.AssessmentClient]
	}
	type args struct {
		target string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.WantErr
	}{
		{
			name: "Set up stream error",
			fields: fields{
				assessment: &api.RPCConnection[assessment.AssessmentClient]{
					Client: &mockAssessmentClient{failStream: true},
				},
			},
			args: args{
				target: "mock-target",
			},
			wantErr: func(t *testing.T, err error) bool {
				return assert.ErrorContains(t, err, "could not set up stream to assessment for assessing evidence")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				assessment: &api.RPCConnection[assessment.AssessmentClient]{
					Client: &mockAssessmentClient{},
				},
			},
			args: args{
				target: "mock-target",
			},
			wantErr: assert.Nil[error],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				assessment: tt.fields.assessment,
			}
			_, err := svc.initAssessmentStream(tt.args.target)
			tt.wantErr(t, err)
		})
	}
}

// mockAssessmentClient is a mock implementation of the AssessmentClient interface.
type mockAssessmentClient struct {
	failStream bool
}

func (m *mockAssessmentClient) AssessEvidences(ctx context.Context, opts ...grpc.CallOption) (assessment.Assessment_AssessEvidencesClient, error) {
	if m.failStream {
		return nil, fmt.Errorf("mock stream failure")
	}
	return &mockAssessmentStream{}, nil
}

// AssessEvidence is a stub implementation to satisfy the AssessmentClient interface.
func (m *mockAssessmentClient) AssessEvidence(ctx context.Context, req *assessment.AssessEvidenceRequest, opts ...grpc.CallOption) (*assessment.AssessEvidenceResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

// CalculateCompliance is a stub implementation to satisfy the AssessmentClient interface.
func (m *mockAssessmentClient) CalculateCompliance(ctx context.Context, req *assessment.CalculateComplianceRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, fmt.Errorf("not implemented")
}

func TestService_handleEvidence(t *testing.T) {
	// mock assessment stream
	mockStream := &mockAssessmentStream{connectionEstablished: true, expected: 2}
	mockStream.Prepare()

	type fields struct {
		storage    persistence.Storage
		assessment *api.RPCConnection[assessment.AssessmentClient]
		authz      service.AuthorizationStrategy
	}
	type args struct {
		evidence  *evidence.Evidence
		addStream bool
		target    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Error getting stream",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
				}),
				assessment: &api.RPCConnection[assessment.AssessmentClient]{
					Client: &mockAssessmentClient{failStream: true},
				},
			},
			args: args{
				addStream: true,
				evidence:  evidencetest.MockEvidence1,
				target:    "error",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not get stream to assessment service")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
				}),
				assessment: &api.RPCConnection[assessment.AssessmentClient]{
					Target: "mock",
				},
			},
			args: args{
				addStream: true,
				evidence:  evidencetest.MockEvidence1,
				target:    "mock",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService()
			svc.storage = tt.fields.storage
			svc.authz = tt.fields.authz
			svc.assessment = tt.fields.assessment
			svc.assessmentStreams = nil

			// Add assessment stream if needed
			if tt.args.addStream {
				svc.assessmentStreams = api.NewStreamsOf[assessment.Assessment_AssessEvidencesClient, *assessment.AssessEvidenceRequest]()
				_, _ = svc.assessmentStreams.GetStream("mock", "Assessment", func(target string, additionalOpts ...grpc.DialOption) (stream assessment.Assessment_AssessEvidencesClient, err error) {
					return mockStream, nil
				})
			}
			err := svc.handleEvidence(tt.args.evidence)
			tt.wantErr(t, err)
		})
	}
}

func TestService_ListSupportedResourceTypes(t *testing.T) {
	type fields struct {
		storage                          persistence.Storage
		assessmentStreams                *api.StreamsOf[assessment.Assessment_AssessEvidencesClient, *assessment.AssessEvidenceRequest]
		assessment                       *api.RPCConnection[assessment.AssessmentClient]
		channelEvidence                  chan *evidence.Evidence
		evidenceHooks                    []evidence.EvidenceHookFunc
		authz                            service.AuthorizationStrategy
		UnimplementedEvidenceStoreServer evidence.UnimplementedEvidenceStoreServer
	}
	type args struct {
		ctx context.Context
		req *evidence.ListSupportedResourceTypesRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes assert.Want[*evidence.ListSupportedResourceTypesResponse]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Request is nil",
			args: args{
				ctx: context.TODO(),
				req: nil,
			},
			wantRes: assert.Nil[*evidence.ListSupportedResourceTypesResponse],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.ErrorContains(t, err, api.ErrEmptyRequest.Error())
				return assert.Equal(t, status.Code(err), codes.InvalidArgument)
			},
		},
		{
			name: "Happy path",
			args: args{
				ctx: context.TODO(),
				req: &evidence.ListSupportedResourceTypesRequest{},
			},
			wantRes: func(t *testing.T, got *evidence.ListSupportedResourceTypesResponse) bool {
				assert.NotNil(t, got)
				return assert.NotEmpty(t, got.ResourceType)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				storage:           tt.fields.storage,
				assessmentStreams: tt.fields.assessmentStreams,
				assessment:        tt.fields.assessment,
				channelEvidence:   tt.fields.channelEvidence,
				evidenceHooks:     tt.fields.evidenceHooks,
				authz:             tt.fields.authz,
			}
			gotRes, err := svc.ListSupportedResourceTypes(tt.args.ctx, tt.args.req)

			tt.wantRes(t, gotRes)
			tt.wantErr(t, err)
		})
	}
}

func TestService_ListResources(t *testing.T) {
	type fields struct {
		storage           persistence.Storage
		assessmentStreams *api.StreamsOf[assessment.Assessment_AssessEvidencesClient, *assessment.AssessEvidenceRequest]
		assessment        *api.RPCConnection[assessment.AssessmentClient]
		channelEvidence   chan *evidence.Evidence
		evidenceHooks     []evidence.EvidenceHookFunc
		authz             service.AuthorizationStrategy
	}
	type args struct {
		ctx context.Context
		req *evidence.ListResourcesRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes assert.Want[*evidence.ListResourcesResponse]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "Request validation error",
			wantRes: assert.Nil[*evidence.ListResourcesResponse],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.ErrorContains(t, err, api.ErrEmptyRequest.Error())
				return assert.Equal(t, status.Code(err), codes.InvalidArgument)
			},
		},
		{
			name: "Filter: ToE not allowed",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID1), // allow only MockTargetOfEvaluationID1
			},
			args: args{
				ctx: context.TODO(),
				req: &evidence.ListResourcesRequest{
					Filter: &evidence.ListResourcesRequest_Filter{
						TargetOfEvaluationId: util.Ref(testdata.MockTargetOfEvaluationID2), // MockTargetOfEvaluationID2 is not allowed
					},
				},
			},
			wantRes: assert.Nil[*evidence.ListResourcesResponse],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, status.Code(err), codes.PermissionDenied)
				return assert.ErrorContains(t, err, service.ErrPermissionDenied.Error())
			},
		},
		{
			name: "Happy path: all filter options used",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(&evidencetest.MockVirtualMachineResource1))
					assert.NoError(t, s.Create(&evidencetest.MockVirtualMachineResource2))
					assert.NoError(t, s.Create(&evidencetest.MockBlockStorageResource1))
					assert.NoError(t, s.Create(&evidencetest.MockBlockStorageResource2))
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				ctx: context.Background(),
				req: &evidence.ListResourcesRequest{
					Filter: &evidence.ListResourcesRequest_Filter{
						TargetOfEvaluationId: util.Ref(testdata.MockTargetOfEvaluationID1),
						ToolId:               util.Ref(testdata.MockEvidenceToolID2),
						Type:                 util.Ref("VirtualMachine"),
					},
				},
			},
			wantRes: func(t *testing.T, got *evidence.ListResourcesResponse) bool {
				return assert.Equal(t, 1, len(got.Results))
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: list only resources for ToE2",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(&evidencetest.MockVirtualMachineResource1))
					assert.NoError(t, s.Create(&evidencetest.MockVirtualMachineResource2))
					assert.NoError(t, s.Create(&evidencetest.MockBlockStorageResource1))
					assert.NoError(t, s.Create(&evidencetest.MockBlockStorageResource2))
				}),
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID2), // allow only MockTargetOfEvaluationID2
			},
			args: args{
				ctx: context.TODO(),
				req: &evidence.ListResourcesRequest{},
			},
			wantRes: func(t *testing.T, got *evidence.ListResourcesResponse) bool {
				return assert.Equal(t, 2, len(got.Results))
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(&evidencetest.MockVirtualMachineResource1))
					assert.NoError(t, s.Create(&evidencetest.MockVirtualMachineResource2))
					assert.NoError(t, s.Create(&evidencetest.MockBlockStorageResource1))
					assert.NoError(t, s.Create(&evidencetest.MockBlockStorageResource2))
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				ctx: context.TODO(),
				req: &evidence.ListResourcesRequest{},
			},
			wantRes: func(t *testing.T, got *evidence.ListResourcesResponse) bool {
				return assert.Equal(t, 4, len(got.Results))
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				storage:           tt.fields.storage,
				assessmentStreams: tt.fields.assessmentStreams,
				assessment:        tt.fields.assessment,
				channelEvidence:   tt.fields.channelEvidence,
				evidenceHooks:     tt.fields.evidenceHooks,
				authz:             tt.fields.authz,
			}
			gotRes, err := svc.ListResources(tt.args.ctx, tt.args.req)

			tt.wantRes(t, gotRes)
			tt.wantErr(t, err)
		})
	}
}
