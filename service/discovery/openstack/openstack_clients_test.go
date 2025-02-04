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

	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"github.com/gophercloud/gophercloud/v2"
)

func Test_openstackDiscovery_identityClient(t *testing.T) {
	type fields struct {
		ctID     string
		clients  clients
		authOpts *gophercloud.AuthOptions
	}
	tests := []struct {
		name       string
		fields     fields
		wantClient assert.Want[*gophercloud.ServiceClient]
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "identityClient not initialized",
			wantClient: assert.Nil[*gophercloud.ServiceClient],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "identity client not initialized")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				clients: clients{
					identityClient: &gophercloud.ServiceClient{},
				},
			},
			wantClient: func(t *testing.T, got *gophercloud.ServiceClient) bool {
				return assert.NotNil(t, got)
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
			gotClient, err := d.identityClient()

			tt.wantClient(t, gotClient)
			tt.wantErr(t, err)
		})
	}
}

func Test_openstackDiscovery_computeClient(t *testing.T) {
	type fields struct {
		ctID     string
		clients  clients
		authOpts *gophercloud.AuthOptions
	}
	tests := []struct {
		name       string
		fields     fields
		wantClient assert.Want[*gophercloud.ServiceClient]
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "computeClient not initialized",
			wantClient: assert.Nil[*gophercloud.ServiceClient],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "compute client not initialized")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				clients: clients{
					computeClient: &gophercloud.ServiceClient{},
				},
			},
			wantClient: func(t *testing.T, got *gophercloud.ServiceClient) bool {
				return assert.NotNil(t, got)
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
			gotClient, err := d.computeClient()

			tt.wantClient(t, gotClient)
			tt.wantErr(t, err)
		})
	}
}

func Test_openstackDiscovery_networkClient(t *testing.T) {
	type fields struct {
		ctID     string
		clients  clients
		authOpts *gophercloud.AuthOptions
	}
	tests := []struct {
		name       string
		fields     fields
		wantClient assert.Want[*gophercloud.ServiceClient]
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "networkClient not initialized",
			wantClient: assert.Nil[*gophercloud.ServiceClient],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "network client not initialized")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				clients: clients{
					networkClient: &gophercloud.ServiceClient{},
				},
			},
			wantClient: func(t *testing.T, got *gophercloud.ServiceClient) bool {
				return assert.NotNil(t, got)
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
			gotClient, err := d.networkClient()

			tt.wantClient(t, gotClient)
			tt.wantErr(t, err)
		})
	}
}

func Test_openstackDiscovery_storageClient(t *testing.T) {
	type fields struct {
		ctID     string
		clients  clients
		authOpts *gophercloud.AuthOptions
	}
	tests := []struct {
		name       string
		fields     fields
		wantClient assert.Want[*gophercloud.ServiceClient]
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "storageClient not initialized",
			wantClient: assert.Nil[*gophercloud.ServiceClient],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "storage client not initialized")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				clients: clients{
					storageClient: &gophercloud.ServiceClient{},
				},
			},
			wantClient: func(t *testing.T, got *gophercloud.ServiceClient) bool {
				return assert.NotNil(t, got)
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
			gotClient, err := d.storageClient()

			tt.wantClient(t, gotClient)
			tt.wantErr(t, err)
		})
	}
}
