// Copyright 2025 Fraunhofer AISEC
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
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/discoverytest/openstacktest"
	"clouditor.io/clouditor/v2/internal/util"
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/clusters"
	"github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_openstackDiscovery_handleCluster(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	openstacktest.HandleListClusterSuccessfully(t)

	t1, err := time.Parse(time.RFC3339, "2014-09-25T13:10:02Z")
	assert.NoError(t, err)

	type fields struct {
		ctID     string
		clients  clients
		authOpts *gophercloud.AuthOptions
		region   string
		domain   *domain
		project  *project
	}
	type args struct {
		cluster *clusters.Cluster
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.Want[ontology.IsResource]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path",
			fields: fields{
				authOpts: &gophercloud.AuthOptions{
					IdentityEndpoint: testdata.MockOpenstackIdentityEndpoint,
					Username:         testdata.MockOpenstackUsername,
					Password:         testdata.MockOpenstackPassword,
					TenantName:       testdata.MockOpenstackTenantName,
				},
				clients: clients{
					provider: &gophercloud.ProviderClient{
						TokenID: client.TokenID,
						EndpointLocator: func(eo gophercloud.EndpointOpts) (string, error) {
							return testhelper.Endpoint(), nil
						},
					},
					clusterClient: client.ServiceClient(),
				},
				region: "test region",
			},
			args: args{
				cluster: &clusters.Cluster{
					UUID:      "ef079b0c-e610-4dfb-b1aa-b49f07ac48e5",
					Name:      "test-cluster",
					CreatedAt: t1,
					ProjectID: "fcad67a6189847c4aecfa3c81a05783b",
					Labels: map[string]string{
						"label1": "value1",
						"label2": "value2",
					},
				},
			},
			want: func(t *testing.T, got ontology.IsResource) bool {
				assert.NotEmpty(t, got)

				want := &ontology.ContainerOrchestration{
					Id:           "ef079b0c-e610-4dfb-b1aa-b49f07ac48e5",
					Name:         "test-cluster",
					CreationTime: timestamppb.New(t1),
					GeoLocation: &ontology.GeoLocation{
						Region: "test region",
					},
					Labels: map[string]string{
						"label1": "value1",
						"label2": "value2",
					},
					ParentId: util.Ref("fcad67a6189847c4aecfa3c81a05783b"),
				}

				gotNew := got.(*ontology.ContainerOrchestration)

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
			}
			got, err := d.handleCluster(tt.args.cluster)

			tt.want(t, got)
			tt.wantErr(t, err)
		})
	}
}
