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

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/internal/testutil/servicetest"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/service"
	"clouditor.io/clouditor/voc"
	"github.com/go-co-op/gocron"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestService_ListGraphEdges(t *testing.T) {
	type fields struct {
		assessmentStreams *api.StreamsOf[assessment.Assessment_AssessEvidencesClient, *assessment.AssessEvidenceRequest]
		assessment        *api.RPCConnection[assessment.AssessmentClient]
		storage           persistence.Storage
		scheduler         *gocron.Scheduler
		authz             service.AuthorizationStrategy
		providers         []string
		Events            chan *DiscoveryEvent
		csID              string
	}
	type args struct {
		ctx context.Context
		req *discovery.ListGraphEdgesRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *discovery.ListGraphEdgesResponse
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
						panicToDiscoveryResource(t, &voc.ObjectStorage{
							Storage: &voc.Storage{
								Resource: discovery.NewResource(&mockDiscoverer{csID: testdata.MockCloudServiceID2},
									"some-id",
									"some-name",
									nil,
									voc.GeoLocation{},
									nil,
									"some-storage-account-id",
									[]string{"ObjectStorage", "Storage", "Resource"},
									map[string][]interface{}{"raw": {"raw"}},
								),
							},
						})))
					assert.NoError(t, s.Create(
						panicToDiscoveryResource(t, &voc.ObjectStorageService{
							StorageService: &voc.StorageService{
								Storage: []voc.ResourceID{"some-id"},
								NetworkService: &voc.NetworkService{
									Networking: &voc.Networking{
										Resource: discovery.NewResource(&mockDiscoverer{csID: testdata.MockCloudServiceID1},
											"some-storage-account-id",
											"some-storage-account-name",
											nil,
											voc.GeoLocation{},
											nil,
											"",
											[]string{"StorageService", "NetworkService", "Networking", "Resource"},
											map[string][]interface{}{"raw": {"raw"}},
										),
									},
								},
							},
							HttpEndpoint: &voc.HttpEndpoint{
								TransportEncryption: &voc.TransportEncryption{
									Enforced:   false,
									Enabled:    true,
									TlsVersion: "TLS1_2",
								},
							},
						})))
				}),
			},
			args: args{
				req: &discovery.ListGraphEdgesRequest{},
			},
			wantRes: &discovery.ListGraphEdgesResponse{
				Edges: []*discovery.GraphEdge{},
			},
			wantErr: assert.NoError,
		},
		{
			name: "happy path",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(true),
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(
						panicToDiscoveryResource(t, &voc.ObjectStorage{
							Storage: &voc.Storage{
								Resource: discovery.NewResource(&mockDiscoverer{},
									"some-id",
									"some-name",
									nil,
									voc.GeoLocation{},
									nil,
									"some-storage-account-id",
									[]string{"ObjectStorage", "Storage", "Resource"},
									map[string][]interface{}{"raw": {"raw"}},
								),
							},
						})))
					assert.NoError(t, s.Create(
						panicToDiscoveryResource(t, &voc.ObjectStorageService{
							StorageService: &voc.StorageService{
								Storage: []voc.ResourceID{"some-id"},
								NetworkService: &voc.NetworkService{
									Networking: &voc.Networking{
										Resource: discovery.NewResource(&mockDiscoverer{},
											"some-storage-account-id",
											"some-storage-account-name",
											nil,
											voc.GeoLocation{},
											nil,
											"",
											[]string{"StorageService", "NetworkService", "Networking", "Resource"},
											map[string][]interface{}{"raw": {"raw"}},
										),
									},
								},
							},
							HttpEndpoint: &voc.HttpEndpoint{
								TransportEncryption: &voc.TransportEncryption{
									Enforced:   false,
									Enabled:    true,
									TlsVersion: "TLS1_2",
								},
							},
						})))
				}),
			},
			args: args{
				req: &discovery.ListGraphEdgesRequest{},
			},
			wantRes: &discovery.ListGraphEdgesResponse{
				Edges: []*discovery.GraphEdge{
					{
						Id:     "some-id-some-storage-account-id",
						Source: "some-id",
						Target: "some-storage-account-id",
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
				scheduler:         tt.fields.scheduler,
				authz:             tt.fields.authz,
				providers:         tt.fields.providers,
				Events:            tt.fields.Events,
				csID:              tt.fields.csID,
			}
			gotRes, err := svc.ListGraphEdges(tt.args.ctx, tt.args.req)

			if !proto.Equal(gotRes, tt.wantRes) {
				t.Errorf("Service.ListGraphEdges() = %v, want %v", gotRes, tt.wantRes)
			}
			tt.wantErr(t, err)
		})
	}
}

func panicToDiscoveryResource(t *testing.T, resource voc.IsCloudResource) *discovery.Resource {
	r, _, err := toDiscoveryResource(resource)
	assert.NoError(t, err)

	return r
}

func TestService_UpdateResource(t *testing.T) {
	type fields struct {
		assessmentStreams *api.StreamsOf[assessment.Assessment_AssessEvidencesClient, *assessment.AssessEvidenceRequest]
		assessment        *api.RPCConnection[assessment.AssessmentClient]
		storage           persistence.Storage
		scheduler         *gocron.Scheduler
		authz             service.AuthorizationStrategy
		providers         []string
		Events            chan *DiscoveryEvent
		csID              string
	}
	type args struct {
		ctx context.Context
		req *discovery.UpdateResourceRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *discovery.Resource
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "validation failed",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockCloudServiceID2),
			},
			args: args{
				req: &discovery.UpdateResourceRequest{
					Resource: panicToDiscoveryResource(t, discovery.NewResource(&mockDiscoverer{
						csID: testdata.MockCloudServiceID1,
					},
						"",
						"some-name",
						nil,
						voc.GeoLocation{},
						nil,
						"",
						[]string{"Resource"},
					)),
				},
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
				req: &discovery.UpdateResourceRequest{
					Resource: panicToDiscoveryResource(t, discovery.NewResource(&mockDiscoverer{
						csID: testdata.MockCloudServiceID1,
					},
						"some-id",
						"some-name",
						nil,
						voc.GeoLocation{},
						nil,
						"",
						[]string{"Resource"},
					)),
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
				req: &discovery.UpdateResourceRequest{
					Resource: panicToDiscoveryResource(t, discovery.NewResource(&mockDiscoverer{
						csID: testdata.MockCloudServiceID1,
					},
						"some-id",
						"some-name",
						nil,
						voc.GeoLocation{},
						nil,
						"",
						[]string{"Resource"},
					)),
				},
			},
			wantRes: panicToDiscoveryResource(t, discovery.NewResource(&mockDiscoverer{
				csID: testdata.MockCloudServiceID1,
			},
				"some-id",
				"some-name",
				nil,
				voc.GeoLocation{},
				nil,
				"",
				[]string{"Resource"},
			)),
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

			if !proto.Equal(gotRes, tt.wantRes) {
				t.Errorf("Service.UpdateResource() = %v, want %v", gotRes, tt.wantRes)
			}
			tt.wantErr(t, err)
		})
	}
}
