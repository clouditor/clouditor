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

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/discoverytest/openstacktest"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func Test_openstackDiscovery_discoverProjects(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	openstacktest.HandleListProjectsSuccessfully(t)

	type fields struct {
		csID     string
		clients  clients
		authOpts *gophercloud.AuthOptions
	}
	tests := []struct {
		name    string
		fields  fields
		want    assert.Want[[]ontology.IsResource]
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
					identityClient: client.ServiceClient(),
				},
			},
			want: func(t *testing.T, got []ontology.IsResource) bool {
				assert.Equal(t, 2, len(got))

				want := &ontology.ResourceGroup{
					Id:          "1234",
					Name:        "Red Team",
					Description: "The team that is red",
					GeoLocation: &ontology.GeoLocation{
						Region: "unknown",
					},
					Labels: map[string]string{
						"Red":  "",
						"Team": "",
					},
					ParentId: util.Ref(""),
				}

				got0 := got[0].(*ontology.ResourceGroup)

				assert.NotEmpty(t, got0.GetRaw())
				got0.Raw = ""
				return assert.Equal(t, want, got0)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &openstackDiscovery{
				csID:     tt.fields.csID,
				clients:  tt.fields.clients,
				authOpts: tt.fields.authOpts,
			}

			gotList, err := d.discoverProjects()

			tt.want(t, gotList)
			tt.wantErr(t, err)
		})
	}
}

func Test_openstackDiscovery_discoverDomain(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	openstacktest.HandleListDomainsSuccessfully(t)

	type fields struct {
		csID     string
		clients  clients
		authOpts *gophercloud.AuthOptions
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
					identityClient: client.ServiceClient(),
				},
			},
			wantList: func(t *testing.T, got []ontology.IsResource) bool {
				assert.Equal(t, 2, len(got))

				want := &ontology.Account{
					Id:          "2844b2a08be147a08ef58317d6471f1f",
					Name:        "domain one",
					Description: "some description",
				}

				got0 := got[0].(*ontology.Account)

				assert.NotEmpty(t, got0.GetRaw())
				got0.Raw = ""
				return assert.Equal(t, want, got0)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &openstackDiscovery{
				csID:     tt.fields.csID,
				clients:  tt.fields.clients,
				authOpts: tt.fields.authOpts,
			}
			gotList, err := d.discoverDomains()

			tt.wantList(t, gotList)
			tt.wantErr(t, err)
		})
	}
}
