// Copyright 2021 Fraunhofer AISEC
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

package azure_test

import (
	"net/http"
	"testing"

	"clouditor.io/clouditor/service/discovery/azure"
	"clouditor.io/clouditor/voc"
	"github.com/stretchr/testify/assert"
)

type mockComputeSender struct {
	mockSender
}

func (m mockComputeSender) Do(req *http.Request) (res *http.Response, err error) {
	if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Compute/virtualMachines" {
		return createResponse(map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":         "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1",
					"name":       "vm1",
					"location":   "eastus",
					"properties": map[string]interface{}{
						"diagnosticsProfile": map[string]interface{}{
							"bootDiagnostics": map[string]interface{}{
								"enabled": true,
							},
						},
					},
				},
				{
					"id":         "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm2",
					"name":       "vm2",
					"location":   "eastus",
					"properties": map[string]interface{}{},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1" {
		return createResponse(map[string]interface{}{
			"properties": map[string]interface{}{
				"storageProfile": map[string]interface{}{
					"osDisk": map[string]interface{}{
						"managedDisk": map[string]interface{}{
							"id": "os_test_disk",
						},
					},
					"dataDisks": &[]map[string]interface{}{
						{
							"managedDisk": map[string]interface{}{
								"id": "data_disk_1",
							},
						},
						{
							"managedDisk": map[string]interface{}{
								"id": "data_disk_2",
							},
						},
					},
				},
				"networkProfile": map[string]interface{}{
					"networkInterfaces": &[]map[string]interface{}{
						{
							"id": "123",
						},
						{
							"id": "234",
						},
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm2" {
		return createResponse(map[string]interface{}{
			"properties": map[string]interface{}{
				"storageProfile": map[string]interface{}{
					"osDisk": map[string]interface{}{
						"managedDisk": map[string]interface{}{
							"id": "os_test_disk",
						},
					},
					"dataDisks": &[]map[string]interface{}{
						{
							"managedDisk": map[string]interface{}{
								"id": "data_disk_2",
							},
						},
						{
							"managedDisk": map[string]interface{}{
								"id": "data_disk_3",
							},
						},
					},
				},
				"networkProfile": map[string]interface{}{
					"networkInterfaces": &[]map[string]interface{}{
						{
							"id": "987",
						},
						{
							"id": "654",
						},
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Web/sites" {
		return createResponse(map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":         "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function1",
					"name":       "function1",
					"location":   "West Europe",
					"properties": map[string]interface{}{},
				},
			},
		}, 200)
	}

	return m.mockSender.Do(req)
}

func TestCompute(t *testing.T) {
	d := azure.NewAzureComputeDiscovery(
		azure.WithSender(&mockComputeSender{}),
		azure.WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 3, len(list))
}

func TestVirtualMachine(t *testing.T) {
	d := azure.NewAzureComputeDiscovery(
		azure.WithSender(&mockComputeSender{}),
		azure.WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()
	assert.Nil(t, err)
	
	virtualMachine, ok := list[0].(*voc.VirtualMachine)

	assert.True(t, ok)
	assert.Equal(t, "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1", string(virtualMachine.ID))
	assert.Equal(t, "vm1", virtualMachine.Name)
	assert.Equal(t, 2, len(virtualMachine.NetworkInterface))
	assert.Equal(t, 3, len(virtualMachine.BlockStorage))

	assert.Equal(t, "data_disk_1", string(virtualMachine.BlockStorage[1]))
	assert.Equal(t, "123", string(virtualMachine.NetworkInterface[0]))
	assert.Equal(t, "eastus", virtualMachine.GeoLocation.Region)
	assert.Equal(t, true, virtualMachine.BootLog.Enabled)
}

func TestFunction(t *testing.T) {
	d := azure.NewAzureComputeDiscovery(
		azure.WithSender(&mockComputeSender{}),
		azure.WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 3, len(list))

	function, ok := list[2].(*voc.Function)

	assert.True(t, ok)
	assert.Equal(t, "function1", function.Name)
}
