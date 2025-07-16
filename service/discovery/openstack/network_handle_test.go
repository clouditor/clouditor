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
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_openstackDiscovery_handleNetworkInterfaces(t *testing.T) {
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
		network *networks.Network
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.Want[ontology.IsResource]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path: projectID available",
			fields: fields{
				region: "test region",
				domain: &domain{
					domainID:   testdata.MockOpenStackDomainID,
					domainName: testdata.MockOpenStackDomainName,
				},
				projects: map[string]ontology.IsResource{},
			},
			args: args{
				network: &networks.Network{
					ID:        testdata.MockOpenstackNetworkID1,
					Name:      testdata.MockOpenstackNetworkName1,
					ProjectID: testdata.MockOpenstackServerTenantID,
					CreatedAt: testTime,
				},
			},
			want: func(t *testing.T, got ontology.IsResource) bool {
				want := &ontology.NetworkInterface{
					Id:           testdata.MockOpenstackNetworkID1,
					Name:         testdata.MockOpenstackNetworkName1,
					CreationTime: timestamppb.New(testTime),
					GeoLocation: &ontology.GeoLocation{
						Region: "test region",
					},
					ParentId: util.Ref(testdata.MockOpenstackServerTenantID),
				}

				gotNew, ok := got.(*ontology.NetworkInterface)
				assert.True(t, ok)

				assert.NotEmpty(t, gotNew.GetRaw())
				gotNew.Raw = ""
				return assert.Equal(t, want, gotNew)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: tenantID available",
			fields: fields{
				region: "test region",
				domain: &domain{
					domainID:   testdata.MockOpenStackDomainID,
					domainName: testdata.MockOpenStackDomainName,
				},
				projects: map[string]ontology.IsResource{},
			},
			args: args{
				network: &networks.Network{
					ID:        testdata.MockOpenstackNetworkID1,
					Name:      testdata.MockOpenstackNetworkName1,
					TenantID:  testdata.MockOpenstackServerTenantID,
					CreatedAt: testTime,
				},
			},
			want: func(t *testing.T, got ontology.IsResource) bool {
				want := &ontology.NetworkInterface{
					Id:           testdata.MockOpenstackNetworkID1,
					Name:         testdata.MockOpenstackNetworkName1,
					CreationTime: timestamppb.New(testTime),
					GeoLocation: &ontology.GeoLocation{
						Region: "test region",
					},
					ParentId: util.Ref(testdata.MockOpenstackServerTenantID),
				}

				gotNew, ok := got.(*ontology.NetworkInterface)
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
				ctID:     tt.fields.ctID,
				clients:  tt.fields.clients,
				authOpts: tt.fields.authOpts,
				region:   tt.fields.region,
				domain:   tt.fields.domain,
				project:  tt.fields.project,
				projects: tt.fields.projects,
			}
			got, err := d.handleNetworkInterfaces(tt.args.network)

			tt.want(t, got)
			tt.wantErr(t, err)
		})
	}
}
