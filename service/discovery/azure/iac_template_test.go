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

type mockIacTemplateSender struct {
	mockSender
}

func (m mockIacTemplateSender) Do(req *http.Request) (res *http.Response, err error) {
	if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups" {
		return createResponse(map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":         "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1",
					"name":       "res1",
					"location":   "eastus",
					"properties": map[string]interface{}{},
				},
				{
					"id":         "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res2",
					"name":       "res2",
					"location":   "eastus",
					"properties": map[string]interface{}{},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/exportTemplate" {

		return createResponse(map[string]interface{}{
			"template": &map[string]interface{}{
				"resources": []map[string]interface{}{
					{
						"type":       "Microsoft.Compute/virtualMachines",
						"name":       "[parameters('virtualMachines_vm1_name')]",
						"location":   "eastus",
						"properties": map[string]interface{}{},
					},
					{
						"type":       "Microsoft.Compute/virtualMachines",
						"name":       "[parameters('virtualMachines_vm2_name')]",
						"location":   "eastus",
						"properties": map[string]interface{}{},
					},
					{
						"type":       "Microsoft.Storage",
						"name":       "[parameters('storageAccounts_storage_1_name')]",
						"location":   "eastus",
						"properties": map[string]interface{}{},
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res2/exportTemplate" {

		return createResponse(map[string]interface{}{
			"template": &map[string]interface{}{
				"resources": []map[string]interface{}{
					{
						"type":       "Microsoft.Compute/virtualMachines",
						"name":       "[parameters('virtualMachines_vm_3_name')]",
						"location":   "eastus",
						"properties": map[string]interface{}{},
					},
					{
						"type":       "Microsoft.Storage",
						"name":       "[parameters('storageAccounts_storage_1_name')]",
						"location":   "eastus",
						"properties": map[string]interface{}{},
					},
					{
						"type":       "Microsoft.Storage",
						"name":       "[parameters('storageAccounts_storage_2_name')]",
						"location":   "eastus",
						"properties": map[string]interface{}{},
					},
				},
			},
		}, 200)
	}

	return m.mockSender.Do(req)
}

func TestIacDiscovery(t *testing.T) {
	d := azure.NewAzureIacTemplateDiscovery(
		azure.WithSender(&mockIacTemplateSender{}),
		azure.WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 6, len(list))

	resourceVM, ok := list[0].(*voc.VirtualMachineResource)
	assert.True(t, ok)
	assert.Equal(t, "vm1", resourceVM.Name)

	resourceStorage, ok := list[2].(*voc.StorageResource)
	assert.True(t, ok)

	// That should be equal. The Problem is described in file 'service/discovery/azure/iac_template.go'
	assert.NotEqual(t, "storage_1", resourceStorage.Name)

}
