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

package azure

import (
	"clouditor.io/clouditor/voc"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault/fake"
	"testing"

	"clouditor.io/clouditor/internal/util"

	"github.com/stretchr/testify/assert"
)

// TODO(lebogg): Add tests for all KeyVault parts here. Try to use mocks provided by Azure

func Test_azureDiscovery_initKeyVaultsClient(t *testing.T) {
	// TODO(lebogg): Do it better (e.g. insert it before every test function (init) or pass as option)
	initKeyVaultTests()
	d := NewMockAzureDiscovery(
		fake.NewVaultsServerTransport(&FakeVaultsServer),
		WithResourceGroup("Non-Existent-Resource-Group"))

	err := d.initKeyVaultsClient()

	// Assert if no error was produced and if client is non-empty
	assert.NoError(t, err)
	assert.NotNil(t, d.clients.keyVaultClient)

	// Assert if the Client behaves as expected: Returning exemplary key vault
	get, err := d.clients.keyVaultClient.Get(context.TODO(), "fake-resource-group", "Fake-KeyVault-Name", nil)
	assert.NoError(t, err)
	assert.Equal(t, "Fake-KeyVault-Name", util.Deref(get.Name))
}

func Test_azureDiscovery_discoverKeyVaults(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
	}
	tests := []struct {
		name     string
		fields   fields
		wantList assert.ValueAssertionFunc
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "happy path - get two key vaults",
			fields: fields{azureDiscovery: NewMockAzureDiscovery(
				fake.NewVaultsServerTransport(&FakeVaultsServer))},
			wantList: func(t assert.TestingT, i1 interface{}, i ...interface{}) bool {
				gotKeyVaults, ok := i1.([]voc.IsCloudResource)
				assert.True(t, ok)
				assert.Len(t, gotKeyVaults, 2)
				assert.Equal(t, mockKeyVault1.Name, gotKeyVaults[0].GetName())
				assert.Equal(t, mockKeyVault2.Name, gotKeyVaults[1].GetName())
				return true
			},
			wantErr: assert.NoError,
		},
		{
			name: "happy path - get one key vault of one resource group",
			fields: fields{azureDiscovery: NewMockAzureDiscovery(
				fake.NewVaultsServerTransport(&FakeVaultsServer),
				WithResourceGroup(string(mockKeyVault1.Parent)))},
			wantList: func(t assert.TestingT, i1 interface{}, i ...interface{}) bool {
				gotKeyVaults, ok := i1.([]voc.IsCloudResource)
				assert.True(t, ok)
				assert.Len(t, gotKeyVaults, 1)
				assert.Equal(t, mockKeyVault1.Name, gotKeyVaults[0].GetName())
				return true
			},
			wantErr: assert.NoError,
		},
		{
			name: "error - wrong resource group",
			fields: fields{azureDiscovery: NewMockAzureDiscovery(
				fake.NewVaultsServerTransport(&FakeVaultsServer),
				WithResourceGroup("Non-Existent-Resource-Group"))},
			wantList: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.ErrorContains(t, err, "invalid resource group")
				return false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO(lebogg): Do it better (e.g. insert it before every test function (init) or pass as option)
			initKeyVaultTests()
			d := tt.fields.azureDiscovery
			gotList, err := d.discoverKeyVaults()
			if !tt.wantErr(t, err, fmt.Sprintf("discoverKeyVaults()")) {
				return
			}
			tt.wantList(t, gotList)
		})
	}
}
