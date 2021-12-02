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
	"fmt"
	"io"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	"testing"

	"clouditor.io/clouditor/service/discovery/azure"
	"clouditor.io/clouditor/voc"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
)

type mockIacTemplateSender struct {
	mockSender
}

func newMockArmTemplateSender() *mockIacTemplateSender {
	m := &mockIacTemplateSender{}
	return m
}

//type responseArmTemplate struct {
//	Value map[string]interface{} `json:"value,omitempty"`
//}

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
				"parameters": map[string]interface{}{
					"virtualMachines_vm_1_2_name": map[string]interface{}{
						"defaultValue": "vm1-2",
						"type":         "String",
					},
					"disks_disk1_name": map[string]interface{}{
						"defaultValue": "disk1",
						"type":         "String",
					},
					"virtualMachines_vm2_name": map[string]interface{}{
						"defaultValue": "vm-2",
						"type":         "String",
					},
					"storageAccounts_storage1_name": map[string]interface{}{
						"defaultValue": "storage1",
						"type":         "String",
					},
				},
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
				"parameters": map[string]interface{}{
					"virtualMachines_vm_3_name": map[string]interface{}{
						"defaultValue": "vm3",
						"type":         "String",
					},
					"storageAccounts_storage3_name": map[string]interface{}{
						"defaultValue": "storage3",
						"type":         "String",
					},
					"loadBalancers_kubernetes_name": map[string]interface{}{
						"defaultValue": "kubernetes",
						"type":         "String",
					},
				},
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
								"supportsHttpsTrafficOnly": false,
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

func TestAzureArmTemplateAuthorizer(t *testing.T) {

	d := azure.NewAzureArmTemplateDiscovery()
	list, err := d.List()

	assert.NotNil(t, err)
	assert.Nil(t, list)
	assert.Equal(t, "could not authorize Azure account: no authorized was available", err.Error())
}

func TestIaCTemplateDiscovery(t *testing.T) {
	d := azure.NewAzureArmTemplateDiscovery(
		azure.WithSender(&mockIacTemplateSender{}),
		azure.WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 7, len(list))
	assert.NotEmpty(t, d.Name())
}

func TestObjectStorageProperties(t *testing.T) {
	d := azure.NewAzureArmTemplateDiscovery(
		azure.WithSender(&mockIacTemplateSender{}),
		azure.WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.Nil(t, err)
	assert.NotNil(t, list)

	objectStorage, ok := list[2].(*voc.ObjectStorage)
	assert.True(t, ok)

	// That should be equal. The Problem is described in file 'service/discovery/azure/arm_template.go' TODO(all); do we need this comment any longer?
	assert.Equal(t, "container1", objectStorage.Name)
	assert.Equal(t, "TLS1_1", objectStorage.HttpEndpoint.TransportEncryption.TlsVersion)
	assert.Equal(t, "ObjectStorage", objectStorage.Type[0])
	assert.Equal(t, "eastus", objectStorage.GeoLocation.Region)
	assert.Equal(t, true, objectStorage.HttpEndpoint.TransportEncryption.Enabled)
	assert.Equal(t, true, objectStorage.HttpEndpoint.TransportEncryption.Enforced)
	assert.Equal(t, "TLS", objectStorage.HttpEndpoint.TransportEncryption.Algorithm)

	// Check ManagedKeyEncryption
	atRestEncryption := *objectStorage.GetAtRestEncryption()
	managedKeyEncryption, ok := atRestEncryption.(voc.ManagedKeyEncryption)
	assert.True(t, ok)
	assert.Equal(t, true, managedKeyEncryption.Enabled)
	assert.Equal(t, "AES256", managedKeyEncryption.Algorithm)

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
	d := azure.NewAzureArmTemplateDiscovery(
		azure.WithSender(&mockIacTemplateSender{}),
		azure.WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.Nil(t, err)
	assert.NotNil(t, list)

	fileStorage, ok := list[3].(*voc.FileStorage)
	assert.True(t, ok)

	// That should be equal. The Problem is described in file 'service/discovery/azure/arm_template.go' TODO(all); do we need this comment any longer?
	assert.Equal(t, "share1", fileStorage.Name)
	assert.Equal(t, "TLS1_1", fileStorage.HttpEndpoint.TransportEncryption.TlsVersion)
	assert.Equal(t, "FileStorage", fileStorage.Type[0])
	assert.Equal(t, "eastus", fileStorage.GeoLocation.Region)
	assert.Equal(t, true, fileStorage.HttpEndpoint.TransportEncryption.Enabled)
	assert.Equal(t, true, fileStorage.HttpEndpoint.TransportEncryption.Enforced)
	assert.Equal(t, "TLS", fileStorage.HttpEndpoint.TransportEncryption.Algorithm)

	// Check ManagedKeyEncryption
	atRestEncryption := *fileStorage.GetAtRestEncryption()
	atRest, ok := atRestEncryption.(voc.ManagedKeyEncryption)
	assert.True(t, ok)
	assert.Equal(t, true, atRest.Enabled)
}

func TestVmProperties(t *testing.T) {
	d := azure.NewAzureArmTemplateDiscovery(
		azure.WithSender(&mockIacTemplateSender{}),
		azure.WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.Nil(t, err)
	assert.NotNil(t, list)

	resourceVM, ok := list[0].(*voc.VirtualMachine)
	assert.True(t, ok)
	assert.Equal(t, "vm1-2", resourceVM.Name)
	assert.Equal(t, "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1-2", (string)(resourceVM.GetID()))
	assert.Equal(t, "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/blockStorage3", (string)(resourceVM.BlockStorage[0]))
	assert.Equal(t, "eastus", resourceVM.GeoLocation.Region)
	assert.True(t, resourceVM.BootLog.Enabled)
	// TODO(garuppel): Currently, we do not get the BootLog Output URI from the Azure call. Do we have to fix the mocking? Check the Azure API call.
	assert.Equal(t, voc.ResourceID("https://storage1.blob.core.windows.net/"), resourceVM.BootLog.Output[0])
	assert.False(t, resourceVM.OSLog.Enabled)
}

func TestLoadBalancerProperties(t *testing.T) {
	d := azure.NewAzureArmTemplateDiscovery(
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

// TestArmTemplateHandleObjectStorageMethodWhenInputIsInvalid tests the method handleObjectStorage w
func TestArmTemplateHandleObjectStorageMethodWhenInputIsInvalid(t *testing.T) {
	// Get mocked Azure Arm Template
	reqURL := "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/exportTemplate"
	mockedArmTemplate, err := getMockedArmTemplate(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	armTemplateResources := mockedArmTemplate["template"].(map[string]interface{})["resources"].([]interface{})

	// TODO(garuppel): d.sub.SubscriptionID = nil
	// check for dependsOn type assertion error
	armTemplateHandleObjectStorageResponse, err := getDependsOnTypeAssertionResponse("Object", armTemplateResources)

	assert.NotNil(t, err)
	assert.Equal(t, "dependsOn type assertion failed", err.Error())
	assert.Nil(t, armTemplateHandleObjectStorageResponse)

	// check for getStorageAccountResourceFromTemplate() response error
	armTemplateHandleObjectStorageResponse, err = getStorageAccountResourceFromTemplateResponse("Object", armTemplateResources)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "cannot get storage account resource from Azure ARM template:")
	assert.Nil(t, armTemplateHandleObjectStorageResponse)

	// check for getStorageAccountAtRestEncryptionFromArm() response error
	armTemplateHandleObjectStorageResponse, err = getStorageAccountAtRestEncryptionFromArmResponse("Object", armTemplateResources)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "cannot get atRestEncryption for storage account resource from Azure ARM template")
	assert.Nil(t, armTemplateHandleObjectStorageResponse)
}

func TestArmTemplateHandleFileStorageMethodWhenInputIsInvalid(t *testing.T) {
	// Get mocked Azure Arm Template
	reqURL := "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/exportTemplate"
	mockedArmTemplate, err := getMockedArmTemplate(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	armTemplateResources := mockedArmTemplate["template"].(map[string]interface{})["resources"].([]interface{})

	// Tests for method handleFileStorage
	// check for dependsOn type assertion error
	armTemplateHandleFileStorageResponse, err := getDependsOnTypeAssertionResponse("File", armTemplateResources)

	assert.NotNil(t, err)
	assert.Equal(t, "dependsOn type assertion failed", err.Error())
	assert.Nil(t, armTemplateHandleFileStorageResponse)

	// check for getStorageAccountResourceFromTemplate() response error
	armTemplateHandleFileStorageResponse, err = getStorageAccountResourceFromTemplateResponse("File", armTemplateResources)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "cannot get storage account resource from Azure ARM template:")
	assert.Nil(t, armTemplateHandleFileStorageResponse)

	// check for getStorageAccountAtRestEncryptionFromArm() response error
	armTemplateHandleFileStorageResponse, err = getStorageAccountAtRestEncryptionFromArmResponse("File", armTemplateResources)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "cannot get atRestEncryption for storage account resource from Azure ARM template")
	assert.Nil(t, armTemplateHandleFileStorageResponse)
}

func TestIsHttpsTrafficOnlyEnabled(t *testing.T) {
	// Get mocked Azure Arm Template
	reqURL := "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/exportTemplate"
	mockedArmTemplate, err := getMockedArmTemplate(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	armTemplateResources := mockedArmTemplate["template"].(map[string]interface{})["resources"].([]interface{})
	// Delete "supportsHttpsTrafficOnly" for test
	delete(armTemplateResources[3].(map[string]interface{})["properties"].(map[string]interface{})["encryption"].(map[string]interface{}), "supportsHttpsTrafficOnly")

	assert.False(t, azure.IsHttpsTrafficOnlyEnabled(armTemplateResources[3].(map[string]interface{})))
}

func TestIsServiceEncryptionEnabled(t *testing.T) {
	// Get mocked Azure Arm Template
	reqURL := "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/exportTemplate"
	mockedArmTemplate, err := getMockedArmTemplate(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	armTemplateResources := mockedArmTemplate["template"].(map[string]interface{})["resources"].([]interface{})
	// Delete "enabled" for test
	delete(armTemplateResources[3].(map[string]interface{})["properties"].(map[string]interface{})["encryption"].(map[string]interface{})["services"].(map[string]interface{})["blob"].(map[string]interface{}), "enabled")

	assert.False(t, azure.IsServiceEncryptionEnabled("blob", armTemplateResources[3].(map[string]interface{})))
}

func TestGetMinTlsVersionOfStorageAccount(t *testing.T) {

	// Get mocked Azure Arm Template
	reqURL := "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/exportTemplate"
	mockedArmTemplate, err := getMockedArmTemplate(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	armTemplateResources := mockedArmTemplate["template"].(map[string]interface{})["resources"].([]interface{})
	// Delete "minimumTlsVersion" for test
	delete(armTemplateResources[3].(map[string]interface{})["properties"].(map[string]interface{})["encryption"].(map[string]interface{}), "minimumTlsVersion")
	assert.Empty(t, azure.GetMinTlsVersionOfStorageAccount(armTemplateResources[3].(map[string]interface{})))
}

func TestMethodGetDefaultResourceNameFromParameter(t *testing.T) {
	var (
		modifiedArmTemplate                         map[string]interface{}
		err                                         error
		getDefaultResourceNameFromParameterResponse string
	)

	// URL mocked Azure Arm Template
	reqURL := "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/exportTemplate"

	// check getDefaultResourceNameFromParameter() type assertion fail
	// Get mocked Azure Arm Template
	modifiedArmTemplate, err = getMockedArmTemplate(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	// Copy Azure ARM template and change "parameters" type
	delete(modifiedArmTemplate["template"].(map[string]interface{}), "parameters")
	modifiedArmTemplate["template"].(map[string]interface{})["parameters"] = []map[string]interface{}{}
	getDefaultResourceNameFromParameterResponse, err = azure.GetDefaultResourceNameFromParameter(modifiedArmTemplate["template"].(map[string]interface{}), "")

	assert.NotNil(t, err.Error())
	assert.Contains(t, err.Error(), "templateValue type assertion failed")
	assert.Empty(t, getDefaultResourceNameFromParameterResponse)

	// check getDefaultResourceNameFromParameter() - error getting default resource name
	// Get mocked Azure Arm Template
	modifiedArmTemplate, err = getMockedArmTemplate(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}
	parameterResourceNameNotExisting := "[parameters('storageAccounts_storageFAIL_name')]"
	getDefaultResourceNameFromParameterResponse, err = azure.GetDefaultResourceNameFromParameter(modifiedArmTemplate["template"].(map[string]interface{}), parameterResourceNameNotExisting)
	assert.NotNil(t, err.Error())
	assert.Contains(t, err.Error(), "parameter resource type assertion failed")
	assert.Empty(t, getDefaultResourceNameFromParameterResponse)

	// check getDefaultResourceNameFromParameter() parameter resource type assertion fail
	// Get mocked Azure Arm Template
	modifiedArmTemplate, err = getMockedArmTemplate(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	// Update "parameters" to "parameter"
	modifiedArmTemplate["template"].(map[string]interface{})["parameter"] = modifiedArmTemplate["template"].(map[string]interface{})["parameters"]
	delete(modifiedArmTemplate["template"].(map[string]interface{}), "parameters")
	getDefaultResourceNameFromParameterResponse, err = azure.GetDefaultResourceNameFromParameter(modifiedArmTemplate["template"].(map[string]interface{}), "")

	assert.NotNil(t, err.Error())
	assert.Contains(t, err.Error(), "error getting default resource name")
	assert.Empty(t, getDefaultResourceNameFromParameterResponse)

	// check getDefaultResourceNameFromParameter() no "defaultValue" available
	// Get mocked Azure Arm Template
	modifiedArmTemplate, err = getMockedArmTemplate(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	// Delete "defaultValue" from parameters for resource "storageAccounts_storage1_name"
	delete(modifiedArmTemplate["template"].(map[string]interface{})["parameters"].(map[string]interface{})["storageAccounts_storage1_name"].(map[string]interface{}), "defaultValue")
	getDefaultResourceNameFromParameterResponse, err = azure.GetDefaultResourceNameFromParameter(modifiedArmTemplate["template"].(map[string]interface{}), "[parameters('storageAccounts_storage1_name')]")

	assert.Nil(t, err)
	assert.Equal(t, "storageAccounts_storage1_name", getDefaultResourceNameFromParameterResponse)
}

// TestMethodGetStorageUriFromArmTemplate tests the  method handleXStorage (X is object or file). Input is invalid and method getStorageAccountResourceFromTemplate() returns an error
func TestMethodGetStorageUriFromArmTemplate(t *testing.T) {
	// Get mocked Azure Arm Template
	reqURL := "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/exportTemplate"
	mockedArmTemplate, err := getMockedArmTemplate(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	bootDiagnostics := mockedArmTemplate["template"].(map[string]interface{})["resources"].([]interface{})[0].(map[string]interface{})["properties"].(map[string]interface{})["diagnosticsProfile"].(map[string]interface{})
	// Delete "storageUri"
	delete(bootDiagnostics["bootDiagnostics"].(map[string]interface{}), "storageUri")
	getStorageUriFromArmTemplateResponse := azure.GetStorageUriFromArmTemplate(bootDiagnostics)

	assert.Empty(t, getStorageUriFromArmTemplateResponse)
}

// getDependsOnTypeAssertionResponse returns the response from method handleXStorage() (X is object or file). Input is invalid and type assertion for "dependsOn" returns an error
func getDependsOnTypeAssertionResponse(storageType string, armTemplateResources []interface{}) (voc.IsCompute, error) {
	var (
		resource                       map[string]interface{}
		modifiedResource               map[string]interface{}
		dependsOnTypeAssertionResponse voc.IsCompute
		err                            error
	)

	// error check for dependsOn type assertion
	switch storageType {
	case "File":
		resource = armTemplateResources[5].(map[string]interface{}) // resource: [concat(parameters('storageAccounts_storage1_name'), 'default/share1')]
		// Copy ARM template resource and delete dependsOn from resource for test
		err = copier.Copy(&modifiedResource, &resource)
		if err != nil {
			fmt.Println("error deep copy")
		}
		delete(modifiedResource, "dependsOn")
		dependsOnTypeAssertionResponse, err = azure.ArmTemplateHandleObjectStorage(modifiedResource, armTemplateResources, "res1")
	case "Object":
		resource = armTemplateResources[4].(map[string]interface{}) // resource: [concat(parameters('storageAccounts_storage1_name'), 'default/container1')]
		// Copy ARM template resource and delete dependsOn from resource for test
		err = copier.Copy(&modifiedResource, &resource)
		if err != nil {
			fmt.Println("error deep copy")
		}
		delete(modifiedResource, "dependsOn")
		dependsOnTypeAssertionResponse, err = azure.ArmTemplateHandleFileStorage(modifiedResource, armTemplateResources, "res1")
	case "Block":
		//TBD
	}

	return dependsOnTypeAssertionResponse, err
}

// getStorageAccountResourceFromTemplateResponse returns the response from method handleXStorage (X is object or file). Input is invalid and method getStorageAccountResourceFromTemplate() returns an error
func getStorageAccountResourceFromTemplateResponse(storageType string, armTemplateResources []interface{}) (voc.IsCompute, error) {
	var (
		resource                                   map[string]interface{}
		modifiedArmTemplateResources               []interface{}
		storageAccountResourceFromTemplateResponse voc.IsCompute
		err                                        error
	)

	// Copy Azure ARM template resources and delete resource "storageAccounts_storage1_name"
	err = copier.Copy(&modifiedArmTemplateResources, &armTemplateResources)
	if err != nil {
		fmt.Println("error deep copy")
	}
	modifiedArmTemplateResources[3] = map[string]interface{}{}

	switch storageType {
	case "File":
		resource = armTemplateResources[5].(map[string]interface{}) // resource: [concat(parameters('storageAccounts_storage1_name'), 'default/share1')]
		storageAccountResourceFromTemplateResponse, err = azure.ArmTemplateHandleObjectStorage(resource, modifiedArmTemplateResources, "res1")
	case "Object":
		resource = armTemplateResources[4].(map[string]interface{}) // resource: [concat(parameters('storageAccounts_storage1_name'), 'default/container1')]
		storageAccountResourceFromTemplateResponse, err = azure.ArmTemplateHandleFileStorage(resource, modifiedArmTemplateResources, "res1")
	case "Block":
		//TBD
	}

	return storageAccountResourceFromTemplateResponse, err
}

// getStorageAccountAtRestEncryptionFromArmResponse returns the response from method handleXStorage (X is object or file). Input is invalid and method getStorageAccountAtRestEncryptionFromArm() returns an error.
func getStorageAccountAtRestEncryptionFromArmResponse(storageType string, armTemplateResources []interface{}) (voc.IsCompute, error) {
	var (
		resource                                      map[string]interface{}
		modifiedArmTemplateResources                  []interface{}
		storageAccountAtRestEncryptionFromArmResponse voc.IsCompute
		err                                           error
	)

	// Copy ARM template resources and delete 'keySource' from storage account resource for test
	err = copier.Copy(&modifiedArmTemplateResources, &armTemplateResources)
	if err != nil {
		fmt.Println("error deep copy")
	}
	delete(modifiedArmTemplateResources[3].(map[string]interface{})["properties"].(map[string]interface{})["encryption"].(map[string]interface{}), "keySource")

	switch storageType {
	case "File":
		resource = armTemplateResources[5].(map[string]interface{}) // resource: [concat(parameters('storageAccounts_storage1_name'), 'default/share1')]
		storageAccountAtRestEncryptionFromArmResponse, err = azure.ArmTemplateHandleObjectStorage(resource, armTemplateResources, "res1")
	case "Object":
		resource = armTemplateResources[4].(map[string]interface{}) // resource: [concat(parameters('storageAccounts_storage1_name'), 'default/container1')]
		storageAccountAtRestEncryptionFromArmResponse, err = azure.ArmTemplateHandleFileStorage(resource, armTemplateResources, "res1")
	case "Block":
		//TBD
	}

	return storageAccountAtRestEncryptionFromArmResponse, err
}

// getMockedArmTemplate returns the Azure ARM template defined in Do() with requestURL (reqUrl).
func getMockedArmTemplate(reqUrl string) (map[string]interface{}, error) {
	//var armTemplateResponse responseArmTemplate
	var armTemplateResponse map[string]interface{}

	m := newMockArmTemplateSender()
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %w", err)
	}
	resp, err := m.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting mock http response: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("error io.ReadCloser: %w", err)
		}
	}(resp.Body)
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error read all: %w", err)
	}
	err = json.Unmarshal(responseBody, &armTemplateResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling: %w", err)
	}

	return armTemplateResponse, nil
}
