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
	"clouditor.io/clouditor/v2/internal/util"
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/domains"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
)

func Test_openstackDiscovery_handleTenant(t *testing.T) {
	type fields struct {
		csID     string
		clients  clients
		authOpts *gophercloud.AuthOptions
	}
	type args struct {
		project *projects.Project
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
			args: args{
				project: &projects.Project{
					ID:          testdata.MockProjectID1,
					Name:        testdata.MockProjectName1,
					Description: testdata.MockProjectDescription1,
					Tags:        []string{},
					ParentID:    testdata.MockProjectParentID1,
				},
			},
			want: func(t *testing.T, got ontology.IsResource) bool {
				want := &ontology.ResourceGroup{
					Id:   testdata.MockProjectID1,
					Name: testdata.MockProjectName1,
					GeoLocation: &ontology.GeoLocation{
						Region: "unknown",
					},
					Description: testdata.MockProjectDescription1,
					Labels:      labels(util.Ref([]string{})),
					ParentId:    util.Ref(testdata.MockProjectParentID1),
				}

				gotNew := got.(*ontology.ResourceGroup)
				assert.NotEmpty(t, gotNew.GetRaw())
				gotNew.Raw = ""
				return assert.Equal(t, want, gotNew)
			},
			wantErr: assert.NoError,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &openstackDiscovery{
				csID:     tt.fields.csID,
				clients:  tt.fields.clients,
				authOpts: tt.fields.authOpts,
			}
			got, err := d.handleTenant(tt.args.project)

			tt.want(t, got)
			tt.wantErr(t, err)
		})
	}
}

func Test_openstackDiscovery_handleDomain(t *testing.T) {
	type fields struct {
		csID     string
		clients  clients
		authOpts *gophercloud.AuthOptions
	}
	type args struct {
		domain *domains.Domain
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
			args: args{
				domain: &domains.Domain{
					ID:          testdata.MockDomainID1,
					Name:        testdata.MockDomainName1,
					Description: testdata.MockDomainDescription1,
				},
			},
			want: func(t *testing.T, got ontology.IsResource) bool {
				want := &ontology.Account{
					Id:          testdata.MockDomainID1,
					Name:        testdata.MockDomainName1,
					Description: testdata.MockDomainDescription1,
				}

				gotNew := got.(*ontology.Account)
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
				csID:     tt.fields.csID,
				clients:  tt.fields.clients,
				authOpts: tt.fields.authOpts,
			}
			got, err := d.handleDomain(tt.args.domain)

			tt.want(t, got)
			tt.wantErr(t, err)
		})
	}
}
