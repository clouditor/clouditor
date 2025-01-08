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
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/discoverytest/openstacktest"
	"clouditor.io/clouditor/v2/internal/util"
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_openstackDiscovery_discoverServer(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	openstacktest.HandleServerListSuccessfully(t)
	openstacktest.HandleInterfaceListSuccessfully(t)

	type fields struct {
		ctID     string
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
				},
			},
			want: func(t *testing.T, got []ontology.IsResource) bool {
				assert.Equal(t, 3, len(got))

				t1, err := time.Parse(time.RFC3339, "2014-09-25T13:10:02Z")
				assert.NoError(t, err)

				want := &ontology.VirtualMachine{
					Id:           "ef079b0c-e610-4dfb-b1aa-b49f07ac48e5",
					Name:         "herp",
					CreationTime: timestamppb.New(t1),
					GeoLocation: &ontology.GeoLocation{
						Region: "unknown",
					},
					Labels:              map[string]string{},
					ParentId:            util.Ref("fcad67a6189847c4aecfa3c81a05783b"),
					BlockStorageIds:     []string{"2bdbc40f-a277-45d4-94ac-d9881c777d33"},
					NetworkInterfaceIds: []string{"8a5fe506-7e9f-4091-899b-96336909d93c"},
					MalwareProtection:   &ontology.MalwareProtection{},
					ActivityLogging:     &ontology.ActivityLogging{},
					AutomaticUpdates:    &ontology.AutomaticUpdates{},
				}

				got0 := got[0].(*ontology.VirtualMachine)

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
				ctID:     tt.fields.ctID,
				clients:  tt.fields.clients,
				authOpts: tt.fields.authOpts,
			}
			gotList, err := d.discoverServer()

			tt.want(t, gotList)
			tt.wantErr(t, err)
		})
	}
}
