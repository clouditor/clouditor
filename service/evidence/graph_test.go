// Copyright 2023 Fraunhofer AISEC
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
	"testing"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest"
	"clouditor.io/clouditor/v2/internal/util"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/service"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestService_ListGraphEdges(t *testing.T) {
	type fields struct {
		assessmentStreams *api.StreamsOf[assessment.Assessment_AssessEvidencesClient, *assessment.AssessEvidenceRequest]
		assessment        *api.RPCConnection[assessment.AssessmentClient]
		storage           persistence.Storage
		authz             service.AuthorizationStrategy
		channelEvidence   chan *evidence.Evidence
		evidenceHooks     []evidence.EvidenceHookFunc
	}
	type args struct {
		ctx context.Context
		req *evidence.ListGraphEdgesRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *evidence.ListGraphEdgesResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "validation failed",
			fields: fields{},
			args:   args{},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "empty request")
			},
		},
		{
			name: "only allowed target of evaluation",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID1),
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(
						panicToDiscoveryResource(t, &ontology.ObjectStorage{
							Id:       "some-id",
							Name:     "some-name",
							ParentId: util.Ref("some-storage-account-id"),
						}, testdata.MockTargetOfEvaluationID2, testdata.MockEvidenceToolID1)))
					assert.NoError(t, s.Create(
						panicToDiscoveryResource(t, &ontology.ObjectStorageService{
							StorageIds: []string{"some-id"},
							Id:         "some-storage-account-id",
							Name:       "some-storage-account-name",
							HttpEndpoint: &ontology.HttpEndpoint{
								TransportEncryption: &ontology.TransportEncryption{
									Enforced:        false,
									Enabled:         true,
									ProtocolVersion: 1.2,
								},
							},
						}, testdata.MockTargetOfEvaluationID1, testdata.MockEvidenceToolID1)))
				}),
			},
			args: args{
				req: &evidence.ListGraphEdgesRequest{},
			},
			wantRes: &evidence.ListGraphEdgesResponse{
				Edges: []*evidence.GraphEdge{
					{
						Type:   "storage",
						Id:     "some-storage-account-id-some-id",
						Source: "some-storage-account-id",
						Target: "some-id",
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "happy path",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(true),
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(
						panicToDiscoveryResource(t, &ontology.ObjectStorage{
							Id:       "some-id",
							Name:     "some-name",
							ParentId: util.Ref("some-storage-account-id"),
						}, testdata.MockTargetOfEvaluationID2, testdata.MockEvidenceToolID1)))
					assert.NoError(t, s.Create(
						panicToDiscoveryResource(t, &ontology.ObjectStorageService{
							StorageIds: []string{"some-id"},
							Id:         "some-storage-account-id",
							Name:       "some-storage-account-name",
							HttpEndpoint: &ontology.HttpEndpoint{
								TransportEncryption: &ontology.TransportEncryption{
									Enforced:        false,
									Enabled:         true,
									ProtocolVersion: 1.2,
								},
							},
						}, testdata.MockTargetOfEvaluationID2, testdata.MockEvidenceToolID1)))
				}),
			},
			args: args{
				req: &evidence.ListGraphEdgesRequest{},
			},
			wantRes: &evidence.ListGraphEdgesResponse{
				Edges: []*evidence.GraphEdge{
					{
						Id:     "some-id-some-storage-account-id",
						Source: "some-id",
						Target: "some-storage-account-id",
						Type:   "parent",
					},
					{
						Type:   "storage",
						Id:     "some-storage-account-id-some-id",
						Source: "some-storage-account-id",
						Target: "some-id",
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				assessmentStreams: tt.fields.assessmentStreams,
				assessment:        tt.fields.assessment,
				storage:           tt.fields.storage,
				authz:             tt.fields.authz,
				channelEvidence:   tt.fields.channelEvidence,
				evidenceHooks:     tt.fields.evidenceHooks,
			}
			gotRes, err := svc.ListGraphEdges(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantRes, gotRes)
			tt.wantErr(t, err)
		})
	}
}

func panicToDiscoveryResource(t *testing.T, resource ontology.IsResource, ctID, collectorID string) *evidence.Resource {
	r, err := evidence.ToEvidenceResource(resource, ctID, collectorID)
	assert.NoError(t, err)

	return r
}

func TestService_UpdateResource(t *testing.T) {
	type fields struct {
		assessmentStreams *api.StreamsOf[assessment.Assessment_AssessEvidencesClient, *assessment.AssessEvidenceRequest]
		assessment        *api.RPCConnection[assessment.AssessmentClient]
		storage           persistence.Storage
		authz             service.AuthorizationStrategy
		channelEvidence   chan *evidence.Evidence
		evidenceHooks     []evidence.EvidenceHookFunc
	}
	type args struct {
		ctx context.Context
		req *evidence.UpdateResourceRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *evidence.Resource
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "validation failed",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID2),
			},
			args: args{
				req: &evidence.UpdateResourceRequest{
					Resource: panicToDiscoveryResource(t, &ontology.VirtualMachine{
						Name: "some-name",
					}, testdata.MockTargetOfEvaluationID1, testdata.MockEvidenceToolID1),
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "resource.id: value length must be at least 1 characters")
			},
		},
		{
			name: "auth failed",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID2),
			},
			args: args{
				req: &evidence.UpdateResourceRequest{
					Resource: panicToDiscoveryResource(t, &ontology.VirtualMachine{
						Id:   "my-id",
						Name: "some-name",
					}, testdata.MockTargetOfEvaluationID1, testdata.MockEvidenceToolID1),
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, service.ErrPermissionDenied)
			},
		},
		{
			name: "happy path",
			fields: fields{
				authz:   servicetest.NewAuthorizationStrategy(true),
				storage: testutil.NewInMemoryStorage(t),
			},
			args: args{
				req: &evidence.UpdateResourceRequest{
					Resource: panicToDiscoveryResource(t, &ontology.VirtualMachine{
						Id:   "my-id",
						Name: "some-name",
					}, testdata.MockTargetOfEvaluationID1, testdata.MockEvidenceToolID1),
				},
			},
			wantRes: panicToDiscoveryResource(t, &ontology.VirtualMachine{
				Id:   "my-id",
				Name: "some-name",
			}, testdata.MockTargetOfEvaluationID1, testdata.MockEvidenceToolID1),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				assessmentStreams: tt.fields.assessmentStreams,
				assessment:        tt.fields.assessment,
				storage:           tt.fields.storage,
				authz:             tt.fields.authz,
				channelEvidence:   tt.fields.channelEvidence,
				evidenceHooks:     tt.fields.evidenceHooks,
			}
			gotRes, err := svc.UpdateResource(tt.args.ctx, tt.args.req)
			assert.Empty(t, cmp.Diff(gotRes, tt.wantRes, protocmp.Transform()))
			tt.wantErr(t, err)
		})
	}
}
