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
						"type":     "Microsoft.Compute/virtualMachines",
						"name":     "[parameters('virtualMachines_vm_1_2_name')]",
						"location": "eastus",
						"properties": map[string]interface{}{
							"storageProfile": map[string]interface{}{
								"dataDisks": []map[string]interface{}{
									{
										"name": "blockStorage3",
										"managedDisks": map[string]interface{}{
											"id": "[resourceId('Microsoft.Compute/disks', 'virtualMachines_blockStorage3_name')]",
										},
									},
									{
										"name": "blockStorage4",
										"managedDisks": map[string]interface{}{
											"id": "[resourceId('Microsoft.Compute/disks', 'virtualMachines_blockStorage4_name')]",
										},
									},
								},
							},
							"diagnosticsProfile": map[string]interface{}{
								"bootDiagnostics": map[string]interface{}{
									"enabled":    true,
									"storageUri": "[concat('https://', parameters('storageAccounts_storage1_name'), '.blob.core.windows.net/')]",
								},
							},
						},
					},
					{
						// TODO Do we need the OD disks?
						"type":     "Microsoft.Compute/disks",
						"name":     "[parameters('disks_disk1_name')]",
						"location": "eastus",
						"properties": map[string]interface{}{
							"encryption": map[string]interface{}{
								"type": "EncryptionAtRestWithPlatformKey",
							},
						},
					},
					{
						"type":     "Microsoft.Compute/virtualMachines",
						"name":     "[parameters('virtualMachines_vm2_name')]",
						"location": "eastus",
						"properties": map[string]interface{}{
							"storageProfile": map[string]interface{}{
								"dataDisks": []map[string]interface{}{
									{
										"name": "blockStorage1",
										"managedDisks": map[string]interface{}{
											"id": "[resourceId('Microsoft.Compute/disks', 'virtualMachines_blockStorage1_name')]",
										},
									},
									{
										"name": "blockStorage2",
										"managedDisks": map[string]interface{}{
											"id": "[resourceId('Microsoft.Compute/disks', 'virtualMachines_blockStorage2_name')]",
										},
									},
								},
							},
							"diagnosticsProfile": map[string]interface{}{
								"bootDiagnostics": map[string]interface{}{
									"enabled":    true,
									"storageUri": "[concat('https://', parameters('storageAccounts_storage_2_name'), '.blob.core.windows.net/')]",
								},
							},
						},
					},
					{
						"type":     "Microsoft.Storage/storageAccounts",
						"name":     "[parameters('storageAccounts_storage1_name')]",
						"location": "eastus",
						"properties": map[string]interface{}{
							"encryption": map[string]interface{}{
								"services": map[string]interface{}{
									"file": map[string]interface{}{
										"keyType": "Account",
										"enabled": true,
									},
									"blob": map[string]interface{}{
										"keyType": "Account",
										"enabled": true,
									},
								},
								"keySource":                "Microsoft.Storage",
								"minimumTlsVersion":        "TLS1_1",
								"supportsHttpsTrafficOnly": true,
							},
						},
					},
					{
						"type": "Microsoft.Storage/storageAccounts/blobServices/containers",
						"name": "[concat(parameters('storageAccounts_storage1_name'), 'default/container1')]",
						"dependsOn": []interface{}{
							"[resourceId('Microsoft.Storage/storageAccounts/blobServices', parameters('storageAccounts_storage1_name'), 'default')]",
							"[resourceId('Microsoft.Storage/storageAccounts', parameters('storageAccounts_storage1_name'))]",
						},
						"properties": map[string]interface{}{
							"defaultEncryptionScope":      "$account-encryption-key",
							"denyEncryptionScopeOverride": false,
							"publicAccess":                "None",
						},
					},
					{
						"type": "Microsoft.Storage/storageAccounts/fileServices/shares",
						"name": "[concat(parameters('storageAccounts_storage1_name'), 'default/share1')]",
						"dependsOn": []interface{}{
							"[resourceId('Microsoft.Storage/storageAccounts/fileServices', parameters('storageAccounts_storage1_name'), 'default')]",
							"[resourceId('Microsoft.Storage/storageAccounts', parameters('storageAccounts_storage1_name'))]",
						},
						"properties": map[string]interface{}{
							"defaultEncryptionScope":      "$account-encryption-key",
							"denyEncryptionScopeOverride": false,
							"publicAccess":                "None",
						},
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res2/exportTemplate" {
		return createResponse(map[string]interface{}{
			"template": &map[string]interface{}{
				"resources": []map[string]interface{}{
					{
						"type":     "Microsoft.Compute/virtualMachines",
						"name":     "[parameters('virtualMachines_vm_3_name')]",
						"location": "eastus",
						"properties": map[string]interface{}{
							"storageProfile": map[string]interface{}{
								"dataDisks": []map[string]interface{}{
									{
										"name": "blockStorage3",
										"managedDisk": map[string]interface{}{
											"id": "[resourceId('Microsoft.Compute/disks', 'virtualMachines_blockStorage3_name')]",
										},
									},
									{
										"name": "blockStorage4",
										"managedDisk": map[string]interface{}{
											"id": "[resourceId('Microsoft.Compute/disks', 'virtualMachines_blockStorage4_name')]",
										},
									},
								},
							},
							"diagnosticsProfile": map[string]interface{}{
								"bootDiagnostics": map[string]interface{}{
									"enabled":    true,
									"storageUri": "[concat('https://', parameters('storageAccounts_storage_3_name'), '.blob.core.windows.net/')]",
								},
							},
						},
					},
					{
						"type":     "Microsoft.Storage/storageAccounts",
						"name":     "[parameters('storageAccounts_storage3_name')]",
						"location": "eastus",
						"properties": map[string]interface{}{
							"encryption": map[string]interface{}{
								"services": map[string]interface{}{
									"file": map[string]interface{}{
										"keyType": "Account",
										"enabled": true,
									},
									"blob": map[string]interface{}{
										"keyType": "Account",
										"enabled": true,
									},
								},
								"keySource": "Microsoft.Keyvault",
								// TODO(all) Is the keyvaulturi correct? There is a difference between the keyUrl (URL), sourceVault (id) and keyvaulturi? Which do we need?
								"keyvaultproperties": map[string]interface{}{
									"keyvaulturi": "https://testvault.vault.azure.net/keys/testkey/123456",
								},
								"minimumTlsVersion":        "TLS1_1",
								"supportsHttpsTrafficOnly": true,
							},
						},
					},
					{
						"type": "Microsoft.Storage/storageAccounts/blobServices/containers",
						"name": "[concat(parameters('storageAccounts_storage3_name'), 'default/container3')]",
						"dependsOn": []interface{}{
							"[resourceId('Microsoft.Storage/storageAccounts/blobServices', parameters('storageAccounts_storage3_name'), 'default')]",
							"[resourceId('Microsoft.Storage/storageAccounts', parameters('storageAccounts_storage3_name'))]",
						},
						"properties": map[string]interface{}{
							"defaultEncryptionScope":      "$account-encryption-key",
							"denyEncryptionScopeOverride": false,
							"publicAccess":                "None",
						},
					},
					{
						"type":       "Microsoft.Network/loadBalancers",
						"name":       "[parameters('loadBalancers_kubernetes_name')]",
						"location":   "eastus",
						"properties": map[string]interface{}{},
					},
				},
			},
		}, 200)
	}

	return m.mockSender.Do(req)
}

func TestIaCTemplateDiscovery(t *testing.T) {
	d := azure.NewAzureIacTemplateDiscovery(
		azure.WithSender(&mockIacTemplateSender{}),
		azure.WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 7, len(list))

}

func TestObjectStorageProperties(t *testing.T) {
	d := azure.NewAzureIacTemplateDiscovery(
		azure.WithSender(&mockIacTemplateSender{}),
		azure.WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.Nil(t, err)
	assert.NotNil(t, list)

	objectStorage, ok := list[2].(*voc.ObjectStorage)
	assert.True(t, ok)

	// That should be equal. The Problem is described in file 'service/discovery/azure/iac_template.go' TODO(all); do we need this comment any longer?
	// TODO(garuppel): Tests for AtRestEncryption, ...
	assert.Equal(t, "container1", objectStorage.Name)
	assert.Equal(t, "TLS1_1", objectStorage.HttpEndpoint.TransportEncryption.TlsVersion)
	assert.Equal(t, "ObjectStorage", objectStorage.Type[0])
	assert.Equal(t, "eastus", objectStorage.GeoLocation.Region)
	assert.Equal(t, true, objectStorage.HttpEndpoint.TransportEncryption.Enabled)
	assert.Equal(t, true, objectStorage.HttpEndpoint.TransportEncryption.Enforced)

	// Check ManagedKeyEncryption
	atRestEncryption := *objectStorage.GetAtRestEncryption()
	managedKeyEncryption, ok := atRestEncryption.(voc.ManagedKeyEncryption)
	assert.True(t, ok)
	assert.Equal(t, true, managedKeyEncryption.Enabled)

	// Check CustomerKeyEncryption
	objectStorage, ok = list[5].(*voc.ObjectStorage)
	assert.True(t, ok)
	atRestEncryption = *objectStorage.GetAtRestEncryption()
	customerKeyEncryption, ok := atRestEncryption.(voc.CustomerKeyEncryption)
	assert.True(t, ok)
	assert.Equal(t, true, customerKeyEncryption.Enabled)
	assert.Equal(t, "https://testvault.vault.azure.net/keys/testkey/123456", customerKeyEncryption.KeyUrl)
}

func TestFileStorageProperties(t *testing.T) {
	d := azure.NewAzureIacTemplateDiscovery(
		azure.WithSender(&mockIacTemplateSender{}),
		azure.WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.Nil(t, err)
	assert.NotNil(t, list)

	fileStorage, ok := list[3].(*voc.FileStorage)
	assert.True(t, ok)

	// That should be equal. The Problem is described in file 'service/discovery/azure/iac_template.go' TODO(all); do we need this comment any longer?
	// TODO(garuppel): Tests for AtRestEncryption, ...
	assert.Equal(t, "share1", fileStorage.Name)
	assert.Equal(t, "TLS1_1", fileStorage.HttpEndpoint.TransportEncryption.TlsVersion)
	assert.Equal(t, "FileStorage", fileStorage.Type[0])
	assert.Equal(t, "eastus", fileStorage.GeoLocation.Region)
	assert.Equal(t, true, fileStorage.HttpEndpoint.TransportEncryption.Enabled)
	assert.Equal(t, true, fileStorage.HttpEndpoint.TransportEncryption.Enforced)

	// Check ManagedKeyEncryption
	atRestEncryption := *fileStorage.GetAtRestEncryption()
	atRest, ok := atRestEncryption.(voc.ManagedKeyEncryption)
	assert.True(t, ok)
	assert.Equal(t, true, atRest.Enabled)
}

func TestVmProperties(t *testing.T) {
	d := azure.NewAzureIacTemplateDiscovery(
		azure.WithSender(&mockIacTemplateSender{}),
		azure.WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.Nil(t, err)
	assert.NotNil(t, list)

	resourceVM, ok := list[0].(*voc.VirtualMachine)
	assert.True(t, ok)
	assert.Equal(t, "vm-1-2", resourceVM.Name)
	assert.Equal(t, "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm-1-2", (string)(resourceVM.GetID()))
	assert.Equal(t, "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/blockStorage3", (string)(resourceVM.BlockStorage[0]))
	assert.Equal(t, "eastus", resourceVM.GeoLocation.Region)
	assert.True(t, resourceVM.BootLog.Enabled)
	assert.False(t, resourceVM.OSLog.Enabled)
}

func TestLoadBalancerProperties(t *testing.T) {
	d := azure.NewAzureIacTemplateDiscovery(
		azure.WithSender(&mockIacTemplateSender{}),
		azure.WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.Nil(t, err)
	assert.NotNil(t, list)

	resourceLoadBalancer, ok := list[6].(*voc.LoadBalancer)
	assert.True(t, ok)
	assert.Equal(t, "kubernetes", resourceLoadBalancer.Name)
	assert.Equal(t, "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res2/providers/Microsoft.Network/loadBalancers/kubernetes", (string)(resourceLoadBalancer.GetID()))
	assert.Equal(t, "LoadBalancer", resourceLoadBalancer.Type[0])
	assert.Equal(t, "eastus", resourceLoadBalancer.GeoLocation.Region)
}
