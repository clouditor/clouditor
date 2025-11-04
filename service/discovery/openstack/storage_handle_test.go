// Copyright 2024 Fraunhofer AISEC
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

package openstack

import (
	"testing"
	"time"

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/util"
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_openstackDiscovery_handleBlockStorage(t *testing.T) {
	testTime := time.Date(2000, 01, 20, 9, 20, 12, 123, time.UTC)

	type fields struct {
		ctID     string
		clients  clients
		authOpts *gophercloud.AuthOptions
		region   string
		domain   *domain
		project  *project
		projects map[string]ontology.IsResource
	}
	type args struct {
		volume *volumes.Volume
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.Want[ontology.IsResource]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "error handling project",
			fields: fields{
				region:   "test region",
				domain:   &domain{},
				projects: map[string]ontology.IsResource{},
			},
			args: args{
				volume: &volumes.Volume{
					ID:        testdata.MockOpenstackVolumeID1,
					CreatedAt: testTime,
				},
			},
			want: assert.Nil[ontology.IsResource],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not handle project for block storage")
			},
		},
		{
			name: "Happy path: volume name missing",
			fields: fields{
				region: "test region",
				domain: &domain{
					domainID:   testdata.MockOpenstackDomainID1,
					domainName: testdata.MockOpenstackDomainName1,
				},
				projects: map[string]ontology.IsResource{},
			},
			args: args{
				volume: &volumes.Volume{
					ID:        testdata.MockOpenstackVolumeID1,
					TenantID:  testdata.MockOpenstackVolumeTenantID,
					CreatedAt: testTime,
				},
			},
			want: func(t *testing.T, got ontology.IsResource) bool {
				want := &ontology.BlockStorage{
					Id:           testdata.MockOpenstackVolumeID1,
					Name:         testdata.MockOpenstackVolumeID1,
					CreationTime: timestamppb.New(testTime),
					GeoLocation: &ontology.GeoLocation{
						Region: "test region",
					},
					ParentId: util.Ref(testdata.MockOpenstackVolumeTenantID),
				}

				gotNew, ok := got.(*ontology.BlockStorage)
				assert.True(t, ok)

				assert.NotEmpty(t, gotNew.GetRaw())
				gotNew.Raw = ""
				return assert.Equal(t, want, gotNew)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: volume name available",
			fields: fields{
				region: "test region",
				domain: &domain{
					domainID:   testdata.MockOpenstackDomainID1,
					domainName: testdata.MockOpenstackDomainName1,
				},
				projects: map[string]ontology.IsResource{},
			},
			args: args{
				volume: &volumes.Volume{
					ID:        testdata.MockOpenstackVolumeID1,
					Name:      testdata.MockOpenstackVolumeName1,
					TenantID:  testdata.MockOpenstackVolumeTenantID,
					CreatedAt: testTime,
				},
			},
			want: func(t *testing.T, got ontology.IsResource) bool {
				want := &ontology.BlockStorage{
					Id:           testdata.MockOpenstackVolumeID1,
					Name:         testdata.MockOpenstackVolumeName1,
					CreationTime: timestamppb.New(testTime),
					GeoLocation: &ontology.GeoLocation{
						Region: "test region",
					},
					ParentId: util.Ref(testdata.MockOpenstackVolumeTenantID),
				}

				gotNew, ok := got.(*ontology.BlockStorage)
				assert.True(t, ok)

				assert.NotEmpty(t, gotNew.GetRaw())
				gotNew.Raw = ""
				return assert.Equal(t, want, gotNew)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &openstackDiscovery{
				ctID:               tt.fields.ctID,
				clients:            tt.fields.clients,
				authOpts:           tt.fields.authOpts,
				region:             tt.fields.region,
				domain:             tt.fields.domain,
				configuredProject:  tt.fields.project,
				discoveredProjects: tt.fields.projects,
			}
			got, err := d.handleBlockStorage(tt.args.volume)

			tt.want(t, got)
			tt.wantErr(t, err)
		})
	}
}
