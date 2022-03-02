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

package azure

import (
	"clouditor.io/clouditor/logging/formatter"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"testing"

	"k8s.io/apimachinery/pkg/util/json"

	"clouditor.io/clouditor/voc"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
)

type mockARMTemplateSender struct {
	mockSender
}

func init() {
	log = logrus.WithField("component", "azure-tests")
	log.Logger.Formatter = formatter.CapitalizeFormatter{Formatter: &logrus.TextFormatter{ForceColors: true, FullTimestamp: true}}
}

func newMockARMTemplateSender() *mockARMTemplateSender {
	m := &mockARMTemplateSender{}
	return m
}

func (m mockARMTemplateSender) Do(req *http.Request) (res *http.Response, err error) {
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

func TestAzureARMTemplateAuthorizer(t *testing.T) {

	d := NewAzureARMTemplateDiscovery()
	list, err := d.List()

	assert.Error(t, err)
	assert.Nil(t, list)
	assert.Equal(t, "could not authorize Azure account: no authorized was available", err.Error())
}

func TestARMTemplateDiscovery(t *testing.T) {
	d := NewAzureARMTemplateDiscovery(
		WithSender(&mockARMTemplateSender{}),
		WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 7, len(list))
	assert.NotEmpty(t, d.Name())
}

func TestObjectStorageProperties(t *testing.T) {
	d := NewAzureARMTemplateDiscovery(
		WithSender(&mockARMTemplateSender{}),
		WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.NoError(t, err)
	assert.NotNil(t, list)

	objectStorage, ok := list[2].(*voc.ObjectStorage)
	assert.True(t, ok)

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
	d := NewAzureARMTemplateDiscovery(
		WithSender(&mockARMTemplateSender{}),
		WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.NoError(t, err)
	assert.NotNil(t, list)

	fileStorage, ok := list[3].(*voc.FileStorage)
	assert.True(t, ok)

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
	d := NewAzureARMTemplateDiscovery(
		WithSender(&mockARMTemplateSender{}),
		WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.NoError(t, err)
	assert.NotNil(t, list)

	resourceVM, ok := list[0].(*voc.VirtualMachine)
	assert.True(t, ok)
	assert.Equal(t, "vm1-2", resourceVM.Name)
	assert.Equal(t, "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1-2", (string)(resourceVM.GetID()))
	assert.Equal(t, "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/blockStorage3", (string)(resourceVM.BlockStorage[0]))
	assert.Equal(t, "eastus", resourceVM.GeoLocation.Region)
	assert.True(t, resourceVM.BootLog.Enabled)
	assert.Equal(t, voc.ResourceID("https://storage1.blob.core.windows.net/"), resourceVM.BootLog.Output[0])
	assert.False(t, resourceVM.OSLog.Enabled)
}

func TestLoadBalancerProperties(t *testing.T) {
	d := NewAzureARMTemplateDiscovery(
		WithSender(&mockARMTemplateSender{}),
		WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.NoError(t, err)
	assert.NotNil(t, list)

	resourceLoadBalancer, ok := list[6].(*voc.LoadBalancer)
	assert.True(t, ok)
	assert.Equal(t, "kubernetes", resourceLoadBalancer.Name)
	assert.Equal(t, "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res2/providers/Microsoft.Network/loadBalancers/kubernetes", (string)(resourceLoadBalancer.GetID()))
	assert.Equal(t, "LoadBalancer", resourceLoadBalancer.Type[0])
	assert.Equal(t, "eastus", resourceLoadBalancer.GeoLocation.Region)
}

// TestARMTemplateHandleObjectStorageMethodWhenInputIsInvalid tests the method handleObjectStorage w
func TestARMTemplateHandleObjectStorageMethodWhenInputIsInvalid(t *testing.T) {
	// Get mocked Azure ARM Template
	reqURL := "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/exportTemplate"
	mockedARMTemplate, err := mockedARMTemplate(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	armTemplateResources := mockedARMTemplate["template"].(map[string]interface{})["resources"].([]interface{})

	// TODO(garuppel): d.sub.SubscriptionID = nil
	// check for dependsOn type assertion error
	armTemplateHandleObjectStorageResponse, err := dependsOnTypeAssertionResponse("Object", armTemplateResources)

	assert.Error(t, err)
	assert.Equal(t, "dependsOn type assertion failed", err.Error())
	assert.Nil(t, armTemplateHandleObjectStorageResponse)

	// check for storageAccountResourceFromARMTemplate() response error
	armTemplateHandleObjectStorageResponse, err = storageAccountResourceFromTemplateResponse("Object", armTemplateResources)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot get storage account resource from Azure ARM template:")
	assert.Nil(t, armTemplateHandleObjectStorageResponse)

	// check for storageAccountAtRestEncryptionFromARMtemplate() response error
	armTemplateHandleObjectStorageResponse, err = storageAccountAtRestEncryptionFromARMResponse("Object", armTemplateResources)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot get atRestEncryption for storage account resource from Azure ARM template")
	assert.Nil(t, armTemplateHandleObjectStorageResponse)
}

func TestARMTemplateHandleFileStorageMethodWhenInputIsInvalid(t *testing.T) {
	// Get mocked Azure ARM Template
	reqURL := "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/exportTemplate"
	mockedARMTemplate, err := mockedARMTemplate(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	armTemplateResources := mockedARMTemplate["template"].(map[string]interface{})["resources"].([]interface{})

	// Tests for method handleFileStorage
	// check for dependsOn type assertion error
	armTemplateHandleFileStorageResponse, err := dependsOnTypeAssertionResponse("File", armTemplateResources)

	assert.Error(t, err)
	assert.Equal(t, "dependsOn type assertion failed", err.Error())
	assert.Nil(t, armTemplateHandleFileStorageResponse)

	// check for storageAccountResourceFromARMTemplate() response error
	armTemplateHandleFileStorageResponse, err = storageAccountResourceFromTemplateResponse("File", armTemplateResources)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot get storage account resource from Azure ARM template:")
	assert.Nil(t, armTemplateHandleFileStorageResponse)

	// check for storageAccountAtRestEncryptionFromARMtemplate() response error
	armTemplateHandleFileStorageResponse, err = storageAccountAtRestEncryptionFromARMResponse("File", armTemplateResources)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot get atRestEncryption for storage account resource from Azure ARM template")
	assert.Nil(t, armTemplateHandleFileStorageResponse)
}

func TestIsHttpsTrafficOnlyEnabled(t *testing.T) {
	// Get mocked Azure ARM Template
	reqURL := "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/exportTemplate"
	mockedARMTemplate, err := mockedARMTemplate(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	armTemplateResources := mockedARMTemplate["template"].(map[string]interface{})["resources"].([]interface{})
	// Delete "supportsHttpsTrafficOnly" for test
	delete(armTemplateResources[3].(map[string]interface{})["properties"].(map[string]interface{})["encryption"].(map[string]interface{}), "supportsHttpsTrafficOnly")

	assert.False(t, isHttpsTrafficOnlyEnabled(armTemplateResources[3].(map[string]interface{})))
}

func TestIsServiceEncryptionEnabled(t *testing.T) {
	// Get mocked Azure ARM Template
	reqURL := "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/exportTemplate"
	mockedARMTemplate, err := mockedARMTemplate(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	armTemplateResources := mockedARMTemplate["template"].(map[string]interface{})["resources"].([]interface{})
	// Delete "enabled" for test
	delete(armTemplateResources[3].(map[string]interface{})["properties"].(map[string]interface{})["encryption"].(map[string]interface{})["services"].(map[string]interface{})["blob"].(map[string]interface{}), "enabled")

	assert.False(t, isServiceEncryptionEnabled("blob", armTemplateResources[3].(map[string]interface{})))
}

func TestMinTlsVersionOfStorageAccount(t *testing.T) {

	// Get mocked Azure ARM Template
	reqURL := "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/exportTemplate"
	mockedARMTemplate, err := mockedARMTemplate(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	armTemplateResources := mockedARMTemplate["template"].(map[string]interface{})["resources"].([]interface{})
	// Delete "minimumTlsVersion" for test
	delete(armTemplateResources[3].(map[string]interface{})["properties"].(map[string]interface{})["encryption"].(map[string]interface{}), "minimumTlsVersion")
	assert.Empty(t, minTlsVersionOfStorageAccount(armTemplateResources[3].(map[string]interface{})))
}

func TestMethodGetDefaultResourceNameFromParameter(t *testing.T) {
	var (
		modifiedARMTemplate                         map[string]interface{}
		err                                         error
		getDefaultResourceNameFromParameterResponse string
	)

	// URL mocked Azure ARM Template
	reqURL := "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/exportTemplate"

	// check defaultResourceNameFromParameter() type assertion fail
	// Get mocked Azure ARM Template
	modifiedARMTemplate, err = mockedARMTemplate(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	// Change "parameters" type
	delete(modifiedARMTemplate["template"].(map[string]interface{}), "parameters")
	modifiedARMTemplate["template"].(map[string]interface{})["parameters"] = []map[string]interface{}{}
	getDefaultResourceNameFromParameterResponse, err = defaultResourceNameFromParameter(modifiedARMTemplate["template"].(map[string]interface{}), "")

	assert.NotNil(t, err.Error())
	assert.Contains(t, err.Error(), "templateValue type assertion failed")
	assert.Empty(t, getDefaultResourceNameFromParameterResponse)

	// check defaultResourceNameFromParameter() - error getting default resource name
	// Get mocked Azure ARM Template
	modifiedARMTemplate, err = mockedARMTemplate(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	// Use parameter resource name that is not existing in 'parameters'
	parameterResourceNameNotExisting := "[parameters('storageAccounts_storageFAIL_name')]"
	getDefaultResourceNameFromParameterResponse, err = defaultResourceNameFromParameter(modifiedARMTemplate["template"].(map[string]interface{}), parameterResourceNameNotExisting)

	assert.NotNil(t, err.Error())
	assert.Contains(t, err.Error(), "parameter resource type assertion failed")
	assert.Empty(t, getDefaultResourceNameFromParameterResponse)

	// check defaultResourceNameFromParameter() parameter resource type assertion fail
	// Get mocked Azure ARM Template
	modifiedARMTemplate, err = mockedARMTemplate(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	// Change "parameters" to "parameter"
	modifiedARMTemplate["template"].(map[string]interface{})["parameter"] = modifiedARMTemplate["template"].(map[string]interface{})["parameters"]
	delete(modifiedARMTemplate["template"].(map[string]interface{}), "parameters")
	getDefaultResourceNameFromParameterResponse, err = defaultResourceNameFromParameter(modifiedARMTemplate["template"].(map[string]interface{}), "")

	assert.NotNil(t, err.Error())
	assert.Contains(t, err.Error(), "error getting default resource name")
	assert.Empty(t, getDefaultResourceNameFromParameterResponse)

	// check defaultResourceNameFromParameter() no "defaultValue" available
	// Get mocked Azure ARM Template
	modifiedARMTemplate, err = mockedARMTemplate(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	// Delete "defaultValue" from parameters for resource "storageAccounts_storage1_name"
	delete(modifiedARMTemplate["template"].(map[string]interface{})["parameters"].(map[string]interface{})["storageAccounts_storage1_name"].(map[string]interface{}), "defaultValue")
	getDefaultResourceNameFromParameterResponse, err = defaultResourceNameFromParameter(modifiedARMTemplate["template"].(map[string]interface{}), "[parameters('storageAccounts_storage1_name')]")

	assert.NoError(t, err)
	assert.Equal(t, "storageAccounts_storage1_name", getDefaultResourceNameFromParameterResponse)
}

// TestMethodGetStorageUriFromARMTemplate tests the  method handleXStorage (X is object or file). Input is invalid and method storageAccountResourceFromARMTemplate() returns an error
func TestMethodGetStorageUriFromARMTemplate(t *testing.T) {
	// Get mocked Azure ARM Template
	reqURL := "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/exportTemplate"
	mockedARMTemplate, err := mockedARMTemplate(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	bootDiagnostics := mockedARMTemplate["template"].(map[string]interface{})["resources"].([]interface{})[0].(map[string]interface{})["properties"].(map[string]interface{})["diagnosticsProfile"].(map[string]interface{})
	// Delete "storageUri"
	delete(bootDiagnostics["bootDiagnostics"].(map[string]interface{}), "storageUri")
	getStorageUriFromARMTemplateResponse := storageURIFromARMTemplate(bootDiagnostics)

	assert.Empty(t, getStorageUriFromARMTemplateResponse)
}

// dependsOnTypeAssertionResponse returns the response from method handleXStorage() (X is object or file). Input is invalid and type assertion for "dependsOn" returns an error
func dependsOnTypeAssertionResponse(storageType string, armTemplateResources []interface{}) (voc.IsCompute, error) {
	d := azureARMTemplateDiscovery{}

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
		dependsOnTypeAssertionResponse, err = d.handleObjectStorage(modifiedResource, armTemplateResources, "res1")
	case "Object":
		resource = armTemplateResources[4].(map[string]interface{}) // resource: [concat(parameters('storageAccounts_storage1_name'), 'default/container1')]
		// Copy ARM template resource and delete dependsOn from resource for test
		err = copier.Copy(&modifiedResource, &resource)
		if err != nil {
			fmt.Println("error deep copy")
		}
		delete(modifiedResource, "dependsOn")
		dependsOnTypeAssertionResponse, err = d.handleFileStorage(modifiedResource, armTemplateResources, "res1")
	case "Block":
		// TODO(garuppel): TBD
	}

	return dependsOnTypeAssertionResponse, err
}

// storageAccountResourceFromTemplateResponse returns the response from method handleXStorage (X is object or file). Input is invalid and method storageAccountResourceFromARMTemplate() returns an error
func storageAccountResourceFromTemplateResponse(storageType string, armTemplateResources []interface{}) (voc.IsCompute, error) {
	d := azureARMTemplateDiscovery{}

	var (
		resource                                   map[string]interface{}
		modifiedARMTemplateResources               []interface{}
		storageAccountResourceFromTemplateResponse voc.IsCompute
		err                                        error
	)

	// Copy Azure ARM template resources and delete resource "storageAccounts_storage1_name"
	err = copier.Copy(&modifiedARMTemplateResources, &armTemplateResources)
	if err != nil {
		fmt.Println("error deep copy")
	}
	modifiedARMTemplateResources[3] = map[string]interface{}{}

	switch storageType {
	case "File":
		resource = armTemplateResources[5].(map[string]interface{}) // resource: [concat(parameters('storageAccounts_storage1_name'), 'default/share1')]
		storageAccountResourceFromTemplateResponse, err = d.handleObjectStorage(resource, modifiedARMTemplateResources, "res1")
	case "Object":
		resource = armTemplateResources[4].(map[string]interface{}) // resource: [concat(parameters('storageAccounts_storage1_name'), 'default/container1')]
		storageAccountResourceFromTemplateResponse, err = d.handleFileStorage(resource, modifiedARMTemplateResources, "res1")
	case "Block":
		// TODO(garuppel): TBD
	}

	return storageAccountResourceFromTemplateResponse, err
}

// storageAccountAtRestEncryptionFromARMResponse returns the response from method handleXStorage (X is object or file). Input is invalid and method storageAccountAtRestEncryptionFromARMtemplate() returns an error.
func storageAccountAtRestEncryptionFromARMResponse(storageType string, armTemplateResources []interface{}) (voc.IsCompute, error) {
	d := azureARMTemplateDiscovery{}

	var (
		resource                                      map[string]interface{}
		modifiedARMTemplateResources                  []interface{}
		storageAccountAtRestEncryptionFromARMResponse voc.IsCompute
		err                                           error
	)

	// Copy ARM template resources and delete 'keySource' from storage account resource for test
	err = copier.Copy(&modifiedARMTemplateResources, &armTemplateResources)
	if err != nil {
		fmt.Println("error deep copy")
	}
	delete(modifiedARMTemplateResources[3].(map[string]interface{})["properties"].(map[string]interface{})["encryption"].(map[string]interface{}), "keySource")

	switch storageType {
	case "File":
		resource = armTemplateResources[5].(map[string]interface{}) // resource: [concat(parameters('storageAccounts_storage1_name'), 'default/share1')]
		storageAccountAtRestEncryptionFromARMResponse, err = d.handleObjectStorage(resource, armTemplateResources, "res1")
	case "Object":
		resource = armTemplateResources[4].(map[string]interface{}) // resource: [concat(parameters('storageAccounts_storage1_name'), 'default/container1')]
		storageAccountAtRestEncryptionFromARMResponse, err = d.handleFileStorage(resource, armTemplateResources, "res1")
	case "Block":
		// TODO(garuppel): TBD
	}

	return storageAccountAtRestEncryptionFromARMResponse, err
}

// mockedARMTemplate returns the Azure ARM template defined in Do() with requestURL (reqUrl).
func mockedARMTemplate(reqUrl string) (map[string]interface{}, error) {
	var armTemplateResponse map[string]interface{}

	m := newMockARMTemplateSender()
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
