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

package discovery

import (
	"context"
	"testing"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/assessment/assessmentconnect"
	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest"
	"clouditor.io/clouditor/v2/internal/util"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/service"
	"connectrpc.com/connect"

	"github.com/go-co-op/gocron"
)

func TestService_ListGraphEdges(t *testing.T) {
	type fields struct {
		assessmentStreams *api.ConnectStreamsOf[
			assessmentconnect.AssessmentClient,
			assessment.AssessEvidenceRequest,
			assessment.AssessEvidencesResponse,
		]
		assessment assessmentconnect.AssessmentClient
		storage    persistence.Storage
		scheduler  *gocron.Scheduler
		authz      service.AuthorizationStrategy
		providers  []string
		Events     chan *DiscoveryEvent
		csID       string
	}
	type args struct {
		ctx context.Context
		req *connect.Request[discovery.ListGraphEdgesRequest]
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *connect.Response[discovery.ListGraphEdgesResponse]
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
			name: "only allowed cloud service",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockCloudServiceID1),
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(
						panicToDiscoveryResource(t, &ontology.ObjectStorage{
							Id:       "some-id",
							Name:     "some-name",
							ParentId: util.Ref("some-storage-account-id"),
						}, testdata.MockCloudServiceID2)))
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
						}, testdata.MockCloudServiceID1)))
				}),
			},
			args: args{
				req: connect.NewRequest(&discovery.ListGraphEdgesRequest{}),
			},
			wantRes: connect.NewResponse(&discovery.ListGraphEdgesResponse{
				Edges: []*discovery.GraphEdge{},
			}),
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
						}, testdata.MockCloudServiceID2)))
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
						}, testdata.MockCloudServiceID2)))
				}),
			},
			args: args{
				req: connect.NewRequest(&discovery.ListGraphEdgesRequest{}),
			},
			wantRes: connect.NewResponse(&discovery.ListGraphEdgesResponse{
				Edges: []*discovery.GraphEdge{
					{
						Id:     "some-id-some-storage-account-id",
						Source: "some-id",
						Target: "some-storage-account-id",
						Type:   "parent",
					},
				},
			}),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				assessmentStreams: tt.fields.assessmentStreams,
				assessment:        tt.fields.assessment,
				storage:           tt.fields.storage,
				scheduler:         tt.fields.scheduler,
				authz:             tt.fields.authz,
				providers:         tt.fields.providers,
				Events:            tt.fields.Events,
				csID:              tt.fields.csID,
			}
			gotRes, err := svc.ListGraphEdges(tt.args.ctx, tt.args.req)

			assert.Equal(t, tt.wantRes, gotRes, assert.CompareAllUnexported())
			tt.wantErr(t, err)
		})
	}
}

func panicToDiscoveryResource(t *testing.T, resource ontology.IsResource, csID string) *discovery.Resource {
	r, err := discovery.ToDiscoveryResource(resource, csID)
	assert.NoError(t, err)

	return r
}

func TestService_UpdateResource(t *testing.T) {
	type fields struct {
		assessmentStreams *api.ConnectStreamsOf[
			assessmentconnect.AssessmentClient,
			assessment.AssessEvidenceRequest,
			assessment.AssessEvidencesResponse,
		]
		assessment assessmentconnect.AssessmentClient
		storage    persistence.Storage
		scheduler  *gocron.Scheduler
		authz      service.AuthorizationStrategy
		providers  []string
		Events     chan *DiscoveryEvent
		csID       string
	}
	type args struct {
		ctx context.Context
		req *connect.Request[discovery.UpdateResourceRequest]
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *connect.Response[discovery.Resource]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "validation failed",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockCloudServiceID2),
			},
			args: args{
				req: connect.NewRequest(&discovery.UpdateResourceRequest{
					Resource: panicToDiscoveryResource(t, &ontology.VirtualMachine{
						Name: "some-name",
					}, testdata.MockCloudServiceID1),
				}),
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "resource.id: value length must be at least 1 characters")
			},
		},
		{
			name: "auth failed",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockCloudServiceID2),
			},
			args: args{
				req: connect.NewRequest(&discovery.UpdateResourceRequest{
					Resource: panicToDiscoveryResource(t, &ontology.VirtualMachine{
						Id:   "my-id",
						Name: "some-name",
					}, testdata.MockCloudServiceID1),
				}),
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
				req: connect.NewRequest(&discovery.UpdateResourceRequest{
					Resource: panicToDiscoveryResource(t, &ontology.VirtualMachine{
						Id:   "my-id",
						Name: "some-name",
					}, testdata.MockCloudServiceID1),
				}),
			},
			wantRes: connect.NewResponse(panicToDiscoveryResource(t, &ontology.VirtualMachine{
				Id:   "my-id",
				Name: "some-name",
			}, testdata.MockCloudServiceID1)),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				assessmentStreams: tt.fields.assessmentStreams,
				assessment:        tt.fields.assessment,
				storage:           tt.fields.storage,
				scheduler:         tt.fields.scheduler,
				authz:             tt.fields.authz,
				providers:         tt.fields.providers,
				Events:            tt.fields.Events,
				csID:              tt.fields.csID,
			}
			gotRes, err := svc.UpdateResource(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantRes, gotRes)
			tt.wantErr(t, err)
		})
	}
}
