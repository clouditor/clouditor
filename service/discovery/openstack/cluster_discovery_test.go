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
	"github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_openstackDiscovery_discoverCluster(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	openstacktest.HandleListClusterSuccessfully(t)

	type fields struct {
		ctID     string
		clients  clients
		authOpts *gophercloud.AuthOptions
		region   string
		domain   *domain
		project  *project
	}
	tests := []struct {
		name     string
		fields   fields
		wantList assert.Want[[]ontology.IsResource]
		wantErr  assert.ErrorAssertionFunc
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
				domain: &domain{
					domainID: testdata.MockOpenstackDomainID1,
				},
				project: &project{
					projectID:   testdata.MockOpenstackProjectID1,
					projectName: testdata.MockOpenstackProjectName1,
				},
			},
			wantList: func(t *testing.T, got []ontology.IsResource) bool {
				assert.Equal(t, 2, len(got))

				t1, err := time.Parse(time.RFC3339, "2016-08-29T06:51:31Z")
				assert.NoError(t, err)

				t2, err := time.Parse(time.RFC3339, "0001-01-01T00:00:00Z")
				assert.NoError(t, err)

				want := &ontology.ContainerOrchestration{
					Id:           "746e779a-751a-456b-a3e9-c883d734946f",
					CreationTime: timestamppb.New(t1),
					Name:         "k8s",
					GeoLocation: &ontology.GeoLocation{
						Region: "test region",
					},
					ParentId: util.Ref(""),
				}

				want1 := &ontology.ContainerOrchestration{
					Id:           "846e779a-751a-456b-a3e9-c883d734946f",
					CreationTime: timestamppb.New(t2),
					Name:         "k8s",
					GeoLocation: &ontology.GeoLocation{
						Region: "test region",
					},
					ParentId: util.Ref(""),
				}

				// Check Raw field and skip it for comparison
				got0 := got[0].(*ontology.ContainerOrchestration)
				assert.NotEmpty(t, got0.GetRaw())
				got0.Raw = ""
				assert.Equal(t, want, got0)

				got1 := got[1].(*ontology.ContainerOrchestration)
				assert.NotEmpty(t, got1.GetRaw())
				got1.Raw = ""
				return assert.Equal(t, want1, got1)
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
			gotList, err := d.discoverCluster()

			tt.wantList(t, gotList)
			tt.wantErr(t, err)
		})
	}
}
