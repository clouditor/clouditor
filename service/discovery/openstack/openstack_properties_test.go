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
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/domains"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func Test_labels(t *testing.T) {
	type args struct {
		tags *[]string
	}
	tests := []struct {
		name string
		args args
		want assert.Want[map[string]string]
	}{
		{
			name: "empty input",
			args: args{},
			want: func(t *testing.T, got map[string]string) bool {
				want := map[string]string{}

				return assert.Equal(t, want, got)
			},
		},
		{
			name: "Happy path",
			args: args{
				tags: &[]string{
					"tag1",
					"tag2",
					"tag3",
				},
			},
			want: func(t *testing.T, got map[string]string) bool {
				want := map[string]string{
					"tag1": "",
					"tag2": "",
					"tag3": "",
				}

				return assert.Equal(t, want, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := labels(tt.args.tags)

			tt.want(t, got)
		})
	}
}

func Test_openstackDiscovery_getAttachedNetworkInterfaces(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	type fields struct {
		ctID       string
		clients    clients
		authOpts   *gophercloud.AuthOptions
		testhelper bool
	}
	type args struct {
		serverID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.Want[[]string]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "error getting network interfaces",
			fields: fields{
				testhelper: false,
				clients: clients{
					provider: &gophercloud.ProviderClient{
						TokenID: client.TokenID,
						EndpointLocator: func(eo gophercloud.EndpointOpts) (string, error) {
							return testhelper.Endpoint(), nil
						},
					},
				},
			},
			args: args{
				serverID: "ef079b0c-e610-4dfb-b1aa-b49f07ac48e5",
			},
			want: assert.Nil[[]string],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not list network interfaces:")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				testhelper: true,
				clients: clients{
					provider: &gophercloud.ProviderClient{
						TokenID: client.TokenID,
						EndpointLocator: func(eo gophercloud.EndpointOpts) (string, error) {
							return testhelper.Endpoint(), nil
						},
					},
				},
			},
			args: args{
				serverID: "ef079b0c-e610-4dfb-b1aa-b49f07ac48e5",
			},
			want: func(t *testing.T, got []string) bool {
				return assert.Equal(t, "8a5fe506-7e9f-4091-899b-96336909d93c", got[0])
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

			if tt.fields.testhelper {
				openstacktest.HandleInterfaceListSuccessfully(t)

			}

			got, err := d.getAttachedNetworkInterfaces(tt.args.serverID)

			tt.want(t, got)
			tt.wantErr(t, err)
		})
	}
}

func Test_openstackDiscovery_setProjectInfo(t *testing.T) {
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
		x interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.Want[*openstackDiscovery]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "error: resource is nil",
			fields: fields{
				project: &project{},
			},
			args: args{
				x: nil,
			},
			want: func(t *testing.T, d *openstackDiscovery) bool {
				return assert.Empty(t, d.project)
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "unknown resource type:")
			},
		},
		{
			name: "error: unknown resource type",
			fields: fields{
				project: &project{},
			},
			args: args{
				x: struct{}{}, // Example of an unknown resource type
			},
			want: func(t *testing.T, d *openstackDiscovery) bool {
				return assert.Empty(t, d.project)
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "unknown resource type:")
			},
		},
		{
			name: "error networks: no tenant or project ID available",
			fields: fields{
				project: &project{},
			},
			args: args{
				x: []networks.Network{
					{},
				},
			},
			want: func(t *testing.T, d *openstackDiscovery) bool {
				return assert.Empty(t, d.project)
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "error getting project ID")
			},
		},
		{
			name: "Happy path: domains",
			fields: fields{
				project: &project{},
			},
			args: args{
				x: []domains.Domain{
					{
						ID:   testdata.MockOpenstackDomainID1,
						Name: testdata.MockOpenstackDomainName1,
					},
				},
			},
			want: func(t *testing.T, d *openstackDiscovery) bool {
				return assert.Empty(t, d.project) // Domain does not have a project ID or name, so we skip this
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: servers",
			fields: fields{
				project: &project{},
			},
			args: args{
				x: []servers.Server{
					{
						TenantID: testdata.MockOpenstackProjectID1,
					},
				},
			},
			want: func(t *testing.T, d *openstackDiscovery) bool {
				assert.Equal(t, testdata.MockOpenstackProjectID1, d.project.projectID)
				return assert.Equal(t, testdata.MockOpenstackProjectID1, d.project.projectName)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: networks TenantID",
			fields: fields{
				project: &project{},
			},
			args: args{
				x: []networks.Network{
					{
						TenantID: testdata.MockOpenstackProjectID1,
					},
				},
			},
			want: func(t *testing.T, d *openstackDiscovery) bool {
				assert.Equal(t, testdata.MockOpenstackProjectID1, d.project.projectID)
				return assert.Equal(t, testdata.MockOpenstackProjectID1, d.project.projectName)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: networks ProjectID",
			fields: fields{
				project: &project{},
			},
			args: args{
				x: []networks.Network{
					{
						ProjectID: testdata.MockOpenstackProjectID1,
					},
				},
			},
			want: func(t *testing.T, d *openstackDiscovery) bool {
				assert.Equal(t, testdata.MockOpenstackProjectID1, d.project.projectID)
				return assert.Equal(t, testdata.MockOpenstackProjectID1, d.project.projectName)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: volumes",
			fields: fields{
				project: &project{},
			},
			args: args{
				x: []volumes.Volume{
					{
						TenantID: testdata.MockOpenstackProjectID1,
					},
				},
			},
			want: func(t *testing.T, d *openstackDiscovery) bool {
				assert.Equal(t, testdata.MockOpenstackProjectID1, d.project.projectID)
				return assert.Equal(t, testdata.MockOpenstackProjectID1, d.project.projectName)
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
			err := d.setProjectInfo(tt.args.x)

			tt.want(t, d)
			tt.wantErr(t, err)
		})
	}
}

func Test_getProjectID(t *testing.T) {
	type args struct {
		call func() (string, error)
	}
	tests := []struct {
		name    string
		args    args
		want    assert.Want[string]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "error: no tenant or project ID available",
			args: args{
				call: func() (string, error) {
					return getProjectID(volumes.Volume{})
				},
			},
			want: assert.Empty[string],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "no tenant or project ID available")
			},
		},
		{
			name: "Happy path: volume",
			args: args{
				call: func() (string, error) {
					return getProjectID(volumes.Volume{
						TenantID: testdata.MockOpenstackProjectID1,
					})
				},
			},
			want: func(t *testing.T, got string) bool {
				return assert.Equal(t, testdata.MockOpenstackProjectID1, got)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: *volume",
			args: args{
				call: func() (string, error) {
					return getProjectID(&volumes.Volume{
						TenantID: testdata.MockOpenstackProjectID1,
					})
				},
			},
			want: func(t *testing.T, got string) bool {
				return assert.Equal(t, testdata.MockOpenstackProjectID1, got)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: networks project ID",
			args: args{
				call: func() (string, error) {
					return getProjectID(networks.Network{
						TenantID: testdata.MockOpenstackProjectID1,
					})
				},
			},
			want: func(t *testing.T, got string) bool {
				return assert.Equal(t, testdata.MockOpenstackProjectID1, got)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: networks tenant ID",
			args: args{
				call: func() (string, error) {
					return getProjectID(networks.Network{
						ProjectID: testdata.MockOpenstackProjectID1,
					})
				},
			},
			want: func(t *testing.T, got string) bool {
				return assert.Equal(t, testdata.MockOpenstackProjectID1, got)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: *networks tenant ID",
			args: args{
				call: func() (string, error) {
					return getProjectID(&networks.Network{
						ProjectID: testdata.MockOpenstackProjectID1,
					})
				},
			},
			want: func(t *testing.T, got string) bool {
				return assert.Equal(t, testdata.MockOpenstackProjectID1, got)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: servers",
			args: args{
				call: func() (string, error) {
					return getProjectID(servers.Server{
						TenantID: testdata.MockOpenstackProjectID1,
					})
				},
			},
			want: func(t *testing.T, got string) bool {
				return assert.Equal(t, testdata.MockOpenstackProjectID1, got)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: *servers",
			args: args{
				call: func() (string, error) {
					return getProjectID(&servers.Server{
						TenantID: testdata.MockOpenstackProjectID1,
					})
				},
			},
			want: func(t *testing.T, got string) bool {
				return assert.Equal(t, testdata.MockOpenstackProjectID1, got)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.call()

			tt.want(t, got)
			tt.wantErr(t, err)
		})
	}
}
