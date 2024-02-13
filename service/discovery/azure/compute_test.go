//go:build exclude

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
	"fmt"
	"reflect"
	"testing"
	"time"

	"clouditor.io/clouditor/api/ontology"
	"clouditor.io/clouditor/internal/constants"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v3"
	"github.com/stretchr/testify/assert"
)

/*
TODO: check for changes
func (m mockComputeSender) Do(req *http.Request) (res *http.Response, err error) {
	if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Compute/virtualMachines" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1",
					"name":     "vm1",
					"location": "eastus",
					"properties": map[string]interface{}{
						"timeCreated": "2017-05-24T13:28:53.4540398Z",
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
						"osProfile": map[string]interface{}{
							"linuxConfiguration": map[string]interface{}{
								"patchSettings": map[string]interface{}{
									"patchMode": "AutomaticByPlatform",
								},
							},
						},
						"diagnosticsProfile": map[string]interface{}{
							"bootDiagnostics": map[string]interface{}{
								"enabled":    true,
								"storageUri": "https://logstoragevm1.blob.core.windows.net/",
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
					"resources": &[]map[string]interface{}{
						{
							"id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1/extensions/MicrosoftMonitoringAgent",
						},
					},
				},
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm2",
					"name":     "vm2",
					"location": "eastus",
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
						"osProfile": map[string]interface{}{
							"windowsConfiguration": map[string]interface{}{
								"patchSettings": map[string]interface{}{
									"patchMode": "AutomaticByOS",
								},
								"enableAutomaticUpdates": true,
							},
						},
						"diagnosticsProfile": map[string]interface{}{
							"bootDiagnostics": map[string]interface{}{
								"enabled":    true,
								"storageUri": nil,
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
					"resources": &[]map[string]interface{}{
						{
							"id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm2/extensions/OmsAgentForLinux",
						},
					},
				},
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm3",
					"name":     "vm3",
					"location": "eastus",
					"properties": map[string]interface{}{
						"diagnosticsProfile": map[string]interface{}{
							"bootDiagnostics": map[string]interface{}{},
						},
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Compute/disks" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk1",
					"name":     "disk1",
					"type":     "Microsoft.Compute/disks",
					"location": "eastus",
					"properties": map[string]interface{}{
						"timeCreated": "2017-05-24T13:28:53.4540398Z",
						"encryption": map[string]interface{}{
							"diskEncryptionSetId": "",
							"type":                "EncryptionAtRestWithPlatformKey",
						},
					},
				},
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk2",
					"name":     "disk2",
					"type":     "Microsoft.Compute/disks",
					"location": "eastus",
					"properties": map[string]interface{}{
						"timeCreated": "2017-05-24T13:28:53.4540398Z",
						"encryption": map[string]interface{}{
							"diskEncryptionSetId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault1",
							"type":                "EncryptionAtRestWithCustomerKey",
						},
					},
				},
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res2/providers/Microsoft.Compute/disks/disk3",
					"name":     "disk3",
					"type":     "Microsoft.Compute/disks",
					"location": "eastus",
					"properties": map[string]interface{}{
						"timeCreated": "2017-05-24T13:28:53.4540398Z",
						"encryption": map[string]interface{}{
							"diskEncryptionSetId": "",
							"type":                "EncryptionAtRestWithPlatformKey",
						},
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res2/providers/Microsoft.Compute/disks" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res2/providers/Microsoft.Compute/disks/disk3",
					"name":     "disk3",
					"type":     "Microsoft.Compute/disks",
					"location": "eastus",
					"properties": map[string]interface{}{
						"timeCreated": "2017-05-24T13:28:53.4540398Z",
						"encryption": map[string]interface{}{
							"diskEncryptionSetId": "",
							"type":                "EncryptionAtRestWithPlatformKey",
						},
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res2/providers/Microsoft.Compute/virtualMachines" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res2/providers/Microsoft.Web/sites" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Web/sites" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function1",
					"name":     "function1",
					"location": "West Europe",
					"kind":     "functionapp,linux",
					"tags": map[string]interface{}{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					"properties": map[string]interface{}{
						"siteConfig": map[string]interface{}{
							"linuxFxVersion": "PYTHON|3.8",
						},
						"resourceGroup":       "res1",
						"publicNetworkAccess": "Enabled",
					},
				},
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function2",
					"name":     "function2",
					"location": "West Europe",
					"kind":     "functionapp",
					"tags": map[string]interface{}{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					"properties": map[string]interface{}{
						"siteConfig":          map[string]interface{}{},
						"resourceGroup":       "res1",
						"publicNetworkAccess": "Disabled",
					},
				},
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp1",
					"name":     "WebApp1",
					"location": "West Europe",
					"kind":     "app",
					"tags": map[string]interface{}{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					"properties": map[string]interface{}{
						"siteConfig":             map[string]interface{}{},
						"httpsOnly":              true,
						"resourceGroup":          "res1",
						"virtualNetworkSubnetId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1",
						"publicNetworkAccess":    "Enabled",
					},
				},
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp2",
					"name":     "WebApp2",
					"location": "West Europe",
					"kind":     "app,linux",
					"tags": map[string]interface{}{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					"properties": map[string]interface{}{
						"siteConfig":             map[string]interface{}{},
						"httpsOnly":              false,
						"resourceGroup":          "res1",
						"virtualNetworkSubnetId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet2",
						"publicNetworkAccess":    "Disabled",
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function1/config/web" {
		return createResponse(req, map[string]interface{}{
			"type": "Microsoft.Web/sites/config",
			"name": "function1",
			"properties": map[string]interface{}{
				"minTlsVersion":     "1.1",
				"minTlsCipherSuite": "TLS_AES_128_GCM_SHA256",
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function2/config/web" {
		return createResponse(req, map[string]interface{}{
			"name": "function2",
			"type": "Microsoft.Web/sites/config",
			"properties": map[string]interface{}{
				"javaVersion": "1.8",
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp1/config/web" {
		return createResponse(req, map[string]interface{}{
			"name": "WebApp1",
			"type": "Microsoft.Web/sites/config",
			"properties": map[string]interface{}{
				"minTlsVersion":     "1.1",
				"minTlsCipherSuite": "TLS_AES_128_GCM_SHA256",
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp2/config/web" {
		return createResponse(req, map[string]interface{}{
			"name":       "WebApp2",
			"type":       "Microsoft.Web/sites/config",
			"properties": map[string]interface{}{},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp1/config/appsettings/list" {
		return createResponse(req, map[string]interface{}{
			"properties": map[string]interface{}{
				"APPLICATIONINSIGHTS_CONNECTION_STRING": "test_connection_string",
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp2/config/appsettings/list" {
		return createResponse(req, map[string]interface{}{
			"properties": map[string]interface{}{},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault1" {
		return createResponse(req, map[string]interface{}{
			"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryption-keyvault1",
			"type":     "Microsoft.Compute/diskEncryptionSets",
			"name":     "encryptionkeyvault1",
			"location": "germanywestcentral",
			"properties": map[string]interface{}{
				"activeKey": map[string]interface{}{
					"sourceVault": map[string]interface{}{
						"id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.KeyVault/vaults/keyvault1",
					},
					"keyUrl": "https://keyvault1.vault.azure.net/keys/customer-key/6273gdb374jz789hjm17819283748382",
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault2" {
		return createResponse(req, map[string]interface{}{
			"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryption-keyvault2",
			"type":     "Microsoft.Compute/diskEncryptionSets",
			"name":     "encryptionkeyvault2",
			"location": "germanywestcentral",
			"properties": map[string]interface{}{
				"activeKey": map[string]interface{}{
					"sourceVault": map[string]interface{}{
						"id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.KeyVault/vaults/keyvault2",
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.DataProtection/backupVaults" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1",
					"name":     "backupAccount1",
					"location": "westeurope",
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupInstances" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupInstances/disk1-disk1-22222222-2222-2222-2222-222222222222",
					"name": "disk1-disk1-22222222-2222-2222-2222-222222222222",
					"properties": map[string]interface{}{
						"dataSourceInfo": map[string]interface{}{
							"resourceID":     "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk1",
							"datasourceType": "Microsoft.Compute/disks",
						},
						"policyInfo": map[string]interface{}{
							"policyId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupPolicies/backupPolicyDisk",
						},
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupPolicies/backupPolicyDisk" {
		return createResponse(req, map[string]interface{}{
			"properties": map[string]interface{}{
				"objectType": "BackupPolicy",
				"policyRules": []map[string]interface{}{
					{
						"objectType": "AzureRetentionRule",
						"lifecycles": []map[string]interface{}{
							{
								"deleteAfter": map[string]interface{}{
									"duration":   "P30D",
									"objectType": "AbsoluteDeleteOption",
								},
								"sourceDataStore": map[string]interface{}{
									"objectType":    "OperationalStore",
									"DataStoreType": "DataStoreInfoBase",
								},
							},
						},
					},
				},
			},
		}, 200)
	}

	return m.mockSender.Do(req)
}*/

func Test_azureComputeDiscovery_discoverFunctionsWebApps(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
	}
	tests := []struct {
		name    string
		fields  fields
		want    []ontology.IsResource
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Error list pages",
			fields: fields{
				azureDiscovery: &azureDiscovery{
					cred: nil,
				},
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrGettingNextPage.Error())
			},
		},
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},

			want: []ontology.IsResource{
				&voc.Function{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function1",
							ServiceID:    testdata.MockCloudServiceID1,
							Name:         "function1",
							CreationTime: util.SafeTimestamp(&time.Time{}),
							Type:         voc.FunctionType,
							Labels: map[string]string{
								"testKey1": "testTag1",
								"testKey2": "testTag2",
							},
							GeoLocation: voc.GeoLocation{
								Region: "West Europe",
							},
							Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1"),
							Raw:    "{\"*armappservice.Site\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function1\",\"kind\":\"functionapp,linux\",\"location\":\"West Europe\",\"name\":\"function1\",\"properties\":{\"publicNetworkAccess\":\"Enabled\",\"resourceGroup\":\"res1\",\"siteConfig\":{\"linuxFxVersion\":\"PYTHON|3.8\"}},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}],\"armappservice.WebAppsClientGetConfigurationResponse\":[{\"name\":\"function1\",\"properties\":{\"minTlsCipherSuite\":\"TLS_AES_128_GCM_SHA256\",\"minTlsVersion\":\"1.1\"},\"type\":\"Microsoft.Web/sites/config\"}]}",
						},
						NetworkInterfaces: []voc.ResourceID{},
						ResourceLogging: &voc.ResourceLogging{
							Logging: &voc.Logging{
								Enabled: false,
							},
						},
					},
					HttpEndpoint: &voc.HttpEndpoint{
						TransportEncryption: &voc.TransportEncryption{
							Enabled:    true,
							Enforced:   false,
							TlsVersion: constants.TLS1_1,
							Algorithm:  string(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
						},
					},
					RuntimeVersion:  "3.8",
					RuntimeLanguage: "PYTHON",
					PublicAccess:    true,
					Redundancy:      &voc.Redundancy{},
				},
				&voc.Function{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function2",
							ServiceID:    testdata.MockCloudServiceID1,
							Name:         "function2",
							CreationTime: util.SafeTimestamp(&time.Time{}),
							Type:         voc.FunctionType,
							Labels: map[string]string{
								"testKey1": "testTag1",
								"testKey2": "testTag2",
							},
							GeoLocation: voc.GeoLocation{
								Region: "West Europe",
							},
							Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1"),
							Raw:    "{\"*armappservice.Site\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function2\",\"kind\":\"functionapp\",\"location\":\"West Europe\",\"name\":\"function2\",\"properties\":{\"publicNetworkAccess\":\"Disabled\",\"resourceGroup\":\"res1\",\"siteConfig\":{}},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}],\"armappservice.WebAppsClientGetConfigurationResponse\":[{\"name\":\"function2\",\"properties\":{\"javaVersion\":\"1.8\"},\"type\":\"Microsoft.Web/sites/config\"}]}",
						},
						NetworkInterfaces: []voc.ResourceID{},
						ResourceLogging: &voc.ResourceLogging{
							Logging: &voc.Logging{
								Enabled: false,
							},
						},
					},
					HttpEndpoint: &voc.HttpEndpoint{
						TransportEncryption: &voc.TransportEncryption{
							Enabled:    false,
							Enforced:   false,
							TlsVersion: "",
							Algorithm:  "",
						},
					},
					RuntimeVersion:  "1.8",
					RuntimeLanguage: "Java",
					PublicAccess:    false,
					Redundancy:      &voc.Redundancy{},
				},
				&voc.WebApp{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp1",
							ServiceID:    testdata.MockCloudServiceID1,
							Name:         "WebApp1",
							CreationTime: util.SafeTimestamp(&time.Time{}),
							Type:         []string{"WebApp", "Compute", "Resource"},
							Labels: map[string]string{
								"testKey1": "testTag1",
								"testKey2": "testTag2",
							},
							GeoLocation: voc.GeoLocation{
								Region: "West Europe",
							},
							Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1"),
							Raw:    "{\"*armappservice.Site\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp1\",\"kind\":\"app\",\"location\":\"West Europe\",\"name\":\"WebApp1\",\"properties\":{\"httpsOnly\":true,\"publicNetworkAccess\":\"Enabled\",\"resourceGroup\":\"res1\",\"siteConfig\":{},\"virtualNetworkSubnetId\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1\"},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}],\"armappservice.WebAppsClientGetConfigurationResponse\":[{\"name\":\"WebApp1\",\"properties\":{\"minTlsCipherSuite\":\"TLS_AES_128_GCM_SHA256\",\"minTlsVersion\":\"1.1\"},\"type\":\"Microsoft.Web/sites/config\"}]}",
						},
						NetworkInterfaces: []voc.ResourceID{"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1"},
						ResourceLogging: &voc.ResourceLogging{
							Logging: &voc.Logging{
								Enabled: true,
							},
						},
					},
					HttpEndpoint: &voc.HttpEndpoint{
						TransportEncryption: &voc.TransportEncryption{
							Enabled:    true,
							Enforced:   true,
							TlsVersion: constants.TLS1_1,
							Algorithm:  string(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
						},
					},
					PublicAccess: true,
					Redundancy:   &voc.Redundancy{},
				},
				&voc.WebApp{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp2",
							ServiceID:    testdata.MockCloudServiceID1,
							Name:         "WebApp2",
							CreationTime: util.SafeTimestamp(&time.Time{}),
							Type:         []string{"WebApp", "Compute", "Resource"},
							Labels: map[string]string{
								"testKey1": "testTag1",
								"testKey2": "testTag2",
							},
							GeoLocation: voc.GeoLocation{
								Region: "West Europe",
							},
							Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1"),
							Raw:    "{\"*armappservice.Site\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp2\",\"kind\":\"app,linux\",\"location\":\"West Europe\",\"name\":\"WebApp2\",\"properties\":{\"httpsOnly\":false,\"publicNetworkAccess\":\"Disabled\",\"resourceGroup\":\"res1\",\"siteConfig\":{},\"virtualNetworkSubnetId\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet2\"},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}],\"armappservice.WebAppsClientGetConfigurationResponse\":[{\"name\":\"WebApp2\",\"properties\":{},\"type\":\"Microsoft.Web/sites/config\"}]}",
						},
						NetworkInterfaces: []voc.ResourceID{"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet2"},
						ResourceLogging: &voc.ResourceLogging{
							Logging: &voc.Logging{
								Enabled: false,
							},
						},
					},
					HttpEndpoint: &voc.HttpEndpoint{
						TransportEncryption: &voc.TransportEncryption{
							Enabled:    false,
							Enforced:   false,
							TlsVersion: "",
							Algorithm:  "",
						},
					},
					PublicAccess: false,
					Redundancy:   &voc.Redundancy{},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			d := tt.fields.azureDiscovery

			got, err := d.discoverFunctionsWebApps()
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_azureComputeDiscovery_handleFunction(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
		clientWebApps  bool
	}
	type args struct {
		function *armappservice.Site
		config   armappservice.WebAppsClientGetConfigurationResponse
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   ontology.IsResource
	}{
		{
			name: "Empty input",
			args: args{
				function: nil,
			},
			want: nil,
		},
		{
			name: "Happy path: Linux function",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
				clientWebApps:  true,
			},
			args: args{
				function: &armappservice.Site{
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function1"),
					Name:     util.Ref("function1"),
					Location: util.Ref("West Europe"),
					Tags: map[string]*string{
						"testKey1": util.Ref("testTag1"),
						"testKey2": util.Ref("testTag2"),
					},
					Kind: util.Ref("functionapp,linux"),
					Properties: &armappservice.SiteProperties{
						SiteConfig: &armappservice.SiteConfig{
							LinuxFxVersion: util.Ref("PYTHON|3.8"),
						},
						HTTPSOnly:     util.Ref(true),
						ResourceGroup: util.Ref("res1"),
					},
				},
				config: armappservice.WebAppsClientGetConfigurationResponse{
					SiteConfigResource: armappservice.SiteConfigResource{
						Properties: &armappservice.SiteConfig{
							MinTLSVersion:     util.Ref(armappservice.SupportedTLSVersionsOne2),
							MinTLSCipherSuite: util.Ref(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
						},
					},
				},
			},
			want: &voc.Function{
				Compute: &voc.Compute{
					Resource: &voc.Resource{
						ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function1",
						ServiceID:    testdata.MockCloudServiceID1,
						Name:         "function1",
						CreationTime: util.SafeTimestamp(&time.Time{}),
						Type:         []string{"Function", "Compute", "Resource"},
						Labels: map[string]string{
							"testKey1": "testTag1",
							"testKey2": "testTag2",
						},
						GeoLocation: voc.GeoLocation{
							Region: "West Europe",
						},
						Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1"),
						Raw:    "{\"*armappservice.Site\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function1\",\"kind\":\"functionapp,linux\",\"location\":\"West Europe\",\"name\":\"function1\",\"properties\":{\"httpsOnly\":true,\"resourceGroup\":\"res1\",\"siteConfig\":{\"linuxFxVersion\":\"PYTHON|3.8\"}},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}],\"armappservice.WebAppsClientGetConfigurationResponse\":[{\"properties\":{\"minTlsCipherSuite\":\"TLS_AES_128_GCM_SHA256\",\"minTlsVersion\":\"1.2\"}}]}",
					},
					NetworkInterfaces: []voc.ResourceID{},
					ResourceLogging: &voc.ResourceLogging{
						Logging: &voc.Logging{
							Enabled: false,
						},
					},
				},
				HttpEndpoint: &voc.HttpEndpoint{
					TransportEncryption: &voc.TransportEncryption{
						Enforced:   true,
						Enabled:    true,
						TlsVersion: constants.TLS1_2,
						Algorithm:  string(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
					},
				},
				RuntimeVersion:  "3.8",
				RuntimeLanguage: "PYTHON",
				Redundancy:      &voc.Redundancy{},
				PublicAccess:    false,
			},
		},
		{
			name: "Happy path: Windows function",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
				clientWebApps:  true,
			},
			args: args{
				function: &armappservice.Site{
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function2"),
					Name:     util.Ref("function2"),
					Location: util.Ref("West Europe"),
					Tags: map[string]*string{
						"testKey1": util.Ref("testTag1"),
						"testKey2": util.Ref("testTag2"),
					},
					Kind: util.Ref("functionapp"),
					Properties: &armappservice.SiteProperties{
						SiteConfig:    &armappservice.SiteConfig{},
						ResourceGroup: util.Ref("res1"),
						HTTPSOnly:     util.Ref(true),
					},
				},
				config: armappservice.WebAppsClientGetConfigurationResponse{
					SiteConfigResource: armappservice.SiteConfigResource{
						Properties: &armappservice.SiteConfig{
							JavaVersion:       util.Ref("1.8"),
							MinTLSVersion:     util.Ref(armappservice.SupportedTLSVersionsOne2),
							MinTLSCipherSuite: util.Ref(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
						},
					},
				},
			},
			want: &voc.Function{
				Compute: &voc.Compute{
					Resource: &voc.Resource{
						ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function2",
						ServiceID:    testdata.MockCloudServiceID1,
						Name:         "function2",
						CreationTime: util.SafeTimestamp(&time.Time{}),
						Type:         []string{"Function", "Compute", "Resource"},
						Labels: map[string]string{
							"testKey1": "testTag1",
							"testKey2": "testTag2",
						},
						GeoLocation: voc.GeoLocation{
							Region: "West Europe",
						},
						Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1"),
						Raw:    "{\"*armappservice.Site\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function2\",\"kind\":\"functionapp\",\"location\":\"West Europe\",\"name\":\"function2\",\"properties\":{\"httpsOnly\":true,\"resourceGroup\":\"res1\",\"siteConfig\":{}},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}],\"armappservice.WebAppsClientGetConfigurationResponse\":[{\"properties\":{\"javaVersion\":\"1.8\",\"minTlsCipherSuite\":\"TLS_AES_128_GCM_SHA256\",\"minTlsVersion\":\"1.2\"}}]}",
					},
					NetworkInterfaces: []voc.ResourceID{},
					ResourceLogging: &voc.ResourceLogging{
						Logging: &voc.Logging{
							Enabled: false,
						},
					},
				},
				HttpEndpoint: &voc.HttpEndpoint{
					TransportEncryption: &voc.TransportEncryption{
						Enforced:   true,
						Enabled:    true,
						TlsVersion: constants.TLS1_2,
						Algorithm:  string(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
					},
				},
				RuntimeVersion:  "1.8",
				RuntimeLanguage: "Java",
				PublicAccess:    false,
				Redundancy:      &voc.Redundancy{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			// Set clients if needed
			if tt.fields.clientWebApps {
				// initialize webApps client
				_ = d.initWebAppsClient()
			}

			assert.Equalf(t, tt.want, d.handleFunction(tt.args.function, tt.args.config), "handleFunction(%v)", tt.args.function)
		})
	}
}

func Test_azureComputeDiscovery_discoverVirtualMachines(t *testing.T) {
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)

	type fields struct {
		azureDiscovery *azureDiscovery
	}
	tests := []struct {
		name    string
		fields  fields
		want    []ontology.IsResource
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Error list pages",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(nil),
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrGettingNextPage.Error())
			},
		},
		{
			name: "No error",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			want: []ontology.IsResource{
				&voc.VirtualMachine{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1",
							ServiceID:    testdata.MockCloudServiceID1,
							Name:         "vm1",
							CreationTime: util.SafeTimestamp(&creationTime),
							Type:         []string{"VirtualMachine", "Compute", "Resource"},
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1"),
							Raw:    "{\"*armcompute.VirtualMachine\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1\",\"location\":\"eastus\",\"name\":\"vm1\",\"properties\":{\"diagnosticsProfile\":{\"bootDiagnostics\":{\"enabled\":true,\"storageUri\":\"https://logstoragevm1.blob.core.windows.net/\"}},\"networkProfile\":{\"networkInterfaces\":[{\"id\":\"123\"},{\"id\":\"234\"}]},\"osProfile\":{\"linuxConfiguration\":{\"patchSettings\":{\"patchMode\":\"AutomaticByPlatform\"}}},\"storageProfile\":{\"dataDisks\":[{\"managedDisk\":{\"id\":\"data_disk_1\"}},{\"managedDisk\":{\"id\":\"data_disk_2\"}}],\"osDisk\":{\"managedDisk\":{\"id\":\"os_test_disk\"}}},\"timeCreated\":\"2017-05-24T13:28:53.4540398Z\"},\"resources\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1/extensions/MicrosoftMonitoringAgent\"}]}]}",
						},
						NetworkInterfaces: []voc.ResourceID{"123", "234"},
					},
					BlockStorage: []voc.ResourceID{"os_test_disk", "data_disk_1", "data_disk_2"},
					BootLogging: &voc.BootLogging{
						Logging: &voc.Logging{
							Enabled: true,
							//LoggingService: []voc.ResourceID{"https://logstoragevm1.blob.core.windows.net/"},
							LoggingService: []voc.ResourceID{},
							Auditing: &voc.Auditing{
								SecurityFeature: &voc.SecurityFeature{},
							},
							RetentionPeriod: 0,
						},
					},
					OsLogging: &voc.OSLogging{
						Logging: &voc.Logging{
							Enabled:         true,
							LoggingService:  []voc.ResourceID{},
							RetentionPeriod: 0,
							Auditing: &voc.Auditing{
								SecurityFeature: &voc.SecurityFeature{},
							},
						},
					},
					AutomaticUpdates: &voc.AutomaticUpdates{
						Enabled:  true,
						Interval: Duration30Days,
					},
					MalwareProtection: &voc.MalwareProtection{},
					ActivityLogging: &voc.ActivityLogging{
						Logging: &voc.Logging{
							Enabled:         true,
							RetentionPeriod: RetentionPeriod90Days,
							LoggingService:  []voc.ResourceID{},
						},
					},
				},
				&voc.VirtualMachine{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm2",
							ServiceID:    testdata.MockCloudServiceID1,
							Name:         "vm2",
							CreationTime: util.SafeTimestamp(&time.Time{}),
							Type:         []string{"VirtualMachine", "Compute", "Resource"},
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1"),
							Raw:    "{\"*armcompute.VirtualMachine\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm2\",\"location\":\"eastus\",\"name\":\"vm2\",\"properties\":{\"diagnosticsProfile\":{\"bootDiagnostics\":{\"enabled\":true}},\"networkProfile\":{\"networkInterfaces\":[{\"id\":\"987\"},{\"id\":\"654\"}]},\"osProfile\":{\"windowsConfiguration\":{\"enableAutomaticUpdates\":true,\"patchSettings\":{\"patchMode\":\"AutomaticByOS\"}}},\"storageProfile\":{\"dataDisks\":[{\"managedDisk\":{\"id\":\"data_disk_2\"}},{\"managedDisk\":{\"id\":\"data_disk_3\"}}],\"osDisk\":{\"managedDisk\":{\"id\":\"os_test_disk\"}}}},\"resources\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm2/extensions/OmsAgentForLinux\"}]}]}",
						},
						NetworkInterfaces: []voc.ResourceID{"987", "654"},
					},
					BlockStorage: []voc.ResourceID{"os_test_disk", "data_disk_2", "data_disk_3"},
					BootLogging: &voc.BootLogging{
						Logging: &voc.Logging{
							Enabled:        true,
							LoggingService: []voc.ResourceID{},
							Auditing: &voc.Auditing{
								SecurityFeature: &voc.SecurityFeature{},
							},
							RetentionPeriod: 0,
						},
					},
					OsLogging: &voc.OSLogging{
						Logging: &voc.Logging{
							Enabled:         true,
							LoggingService:  []voc.ResourceID{},
							RetentionPeriod: 0,
							Auditing: &voc.Auditing{
								SecurityFeature: &voc.SecurityFeature{},
							},
						},
					},
					AutomaticUpdates: &voc.AutomaticUpdates{
						Enabled:  true,
						Interval: Duration30Days,
					},
					MalwareProtection: &voc.MalwareProtection{},
					ActivityLogging: &voc.ActivityLogging{
						Logging: &voc.Logging{
							Enabled:         true,
							RetentionPeriod: RetentionPeriod90Days,
							LoggingService:  []voc.ResourceID{},
						},
					},
				},
				&voc.VirtualMachine{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm3",
							ServiceID:    testdata.MockCloudServiceID1,
							Name:         "vm3",
							Type:         []string{"VirtualMachine", "Compute", "Resource"},
							CreationTime: util.SafeTimestamp(&time.Time{}),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1"),
							Raw:    "{\"*armcompute.VirtualMachine\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm3\",\"location\":\"eastus\",\"name\":\"vm3\",\"properties\":{\"diagnosticsProfile\":{\"bootDiagnostics\":{}}}}]}",
						},
						NetworkInterfaces: []voc.ResourceID{},
					},
					BlockStorage: []voc.ResourceID{},
					BootLogging: &voc.BootLogging{
						Logging: &voc.Logging{
							Enabled:         false,
							LoggingService:  []voc.ResourceID{},
							RetentionPeriod: 0,
							Auditing: &voc.Auditing{
								SecurityFeature: &voc.SecurityFeature{},
							},
						},
					},
					OsLogging: &voc.OSLogging{
						Logging: &voc.Logging{
							Enabled:         false,
							LoggingService:  []voc.ResourceID{},
							RetentionPeriod: 0,
							Auditing: &voc.Auditing{
								SecurityFeature: &voc.SecurityFeature{},
							},
						},
					},
					AutomaticUpdates: &voc.AutomaticUpdates{
						Enabled:  false,
						Interval: time.Duration(0),
					},
					MalwareProtection: &voc.MalwareProtection{},
					ActivityLogging: &voc.ActivityLogging{
						Logging: &voc.Logging{
							Enabled:         true,
							RetentionPeriod: RetentionPeriod90Days,
							LoggingService:  []voc.ResourceID{},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			d := tt.fields.azureDiscovery

			got, err := d.discoverVirtualMachines()
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_azureComputeDiscovery_handleVirtualMachines(t *testing.T) {
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)

	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		vm *armcompute.VirtualMachine
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ontology.IsResource
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Virtual Machine is empty",
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrEmptyVirtualMachine)
			},
		},
		{
			name: "No error",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender(), WithDefenderProperties(map[string]*defenderProperties{
					DefenderVirtualMachineType: {
						monitoringLogDataEnabled: true,
						securityAlertsEnabled:    true,
					},
				})),
			},
			args: args{
				vm: &armcompute.VirtualMachine{
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1"),
					Name:     util.Ref("vm1"),
					Location: util.Ref("eastus"),
					Properties: &armcompute.VirtualMachineProperties{
						TimeCreated: &creationTime,
						NetworkProfile: &armcompute.NetworkProfile{
							NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
								{
									ID: util.Ref("123"),
								},
								{
									ID: util.Ref("234"),
								},
							},
						},
						StorageProfile: &armcompute.StorageProfile{
							OSDisk: &armcompute.OSDisk{
								ManagedDisk: &armcompute.ManagedDiskParameters{
									ID: util.Ref("os_test_disk"),
								},
							},
							DataDisks: []*armcompute.DataDisk{
								{
									ManagedDisk: &armcompute.ManagedDiskParameters{
										ID: util.Ref("data_disk_1"),
									},
								},
								{
									ManagedDisk: &armcompute.ManagedDiskParameters{
										ID: util.Ref("data_disk_2"),
									},
								},
							},
						},
						DiagnosticsProfile: &armcompute.DiagnosticsProfile{
							BootDiagnostics: &armcompute.BootDiagnostics{
								Enabled:    util.Ref(true),
								StorageURI: util.Ref("https://logstoragevm1.blob.core.windows.net/"),
							},
						},
					},
				},
			},
			want: &voc.VirtualMachine{
				Compute: &voc.Compute{
					Resource: &voc.Resource{
						ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1",
						ServiceID:    testdata.MockCloudServiceID1,
						Name:         "vm1",
						CreationTime: util.SafeTimestamp(&creationTime),
						Type:         []string{"VirtualMachine", "Compute", "Resource"},
						Labels:       map[string]string{},
						GeoLocation: voc.GeoLocation{
							Region: "eastus",
						},
						Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1"),
						Raw:    "{\"*armcompute.VirtualMachine\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1\",\"location\":\"eastus\",\"name\":\"vm1\",\"properties\":{\"diagnosticsProfile\":{\"bootDiagnostics\":{\"enabled\":true,\"storageUri\":\"https://logstoragevm1.blob.core.windows.net/\"}},\"networkProfile\":{\"networkInterfaces\":[{\"id\":\"123\"},{\"id\":\"234\"}]},\"storageProfile\":{\"dataDisks\":[{\"managedDisk\":{\"id\":\"data_disk_1\"}},{\"managedDisk\":{\"id\":\"data_disk_2\"}}],\"osDisk\":{\"managedDisk\":{\"id\":\"os_test_disk\"}}},\"timeCreated\":\"2017-05-24T13:28:53.004540398Z\"}}]}",
					},
					NetworkInterfaces: []voc.ResourceID{"123", "234"},
				},
				BlockStorage: []voc.ResourceID{"os_test_disk", "data_disk_1", "data_disk_2"},
				BootLogging: &voc.BootLogging{
					Logging: &voc.Logging{
						Enabled: true,
						//LoggingService: []voc.ResourceID{"https://logstoragevm1.blob.core.windows.net/"},
						LoggingService: []voc.ResourceID{},
						Auditing: &voc.Auditing{
							SecurityFeature: &voc.SecurityFeature{},
						},
						RetentionPeriod:          0,
						MonitoringLogDataEnabled: true,
						SecurityAlertsEnabled:    true,
					},
				},
				OsLogging: &voc.OSLogging{
					Logging: &voc.Logging{
						Enabled:         false,
						LoggingService:  []voc.ResourceID{},
						RetentionPeriod: 0,
						Auditing: &voc.Auditing{
							SecurityFeature: &voc.SecurityFeature{},
						},
						MonitoringLogDataEnabled: true,
						SecurityAlertsEnabled:    true,
					},
				},
				AutomaticUpdates: &voc.AutomaticUpdates{
					Enabled:  false,
					Interval: time.Duration(0),
				},
				MalwareProtection: &voc.MalwareProtection{},
				ActivityLogging: &voc.ActivityLogging{
					Logging: &voc.Logging{
						Enabled:         true,
						RetentionPeriod: RetentionPeriod90Days,
						LoggingService:  []voc.ResourceID{},
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			got, err := d.handleVirtualMachines(tt.args.vm)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_isBootDiagnosticEnabled(t *testing.T) {
	ID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1"
	name := "vm1"
	enabledTrue := true

	type args struct {
		vm *armcompute.VirtualMachine
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty input",
			args: args{
				vm: nil,
			},
			want: false,
		},
		{
			name: "Empty properties value",
			args: args{
				vm: &armcompute.VirtualMachine{
					ID:         &ID,
					Name:       &name,
					Properties: nil,
				},
			},
			want: false,
		},
		{
			name: "Empty DiagnosticsProfile value",
			args: args{
				vm: &armcompute.VirtualMachine{
					ID:   &ID,
					Name: &name,
					Properties: &armcompute.VirtualMachineProperties{
						DiagnosticsProfile: nil,
					},
				},
			},
			want: false,
		},
		{
			name: "Empty BootDiagnostics value",
			args: args{
				vm: &armcompute.VirtualMachine{
					ID:   &ID,
					Name: &name,
					Properties: &armcompute.VirtualMachineProperties{
						DiagnosticsProfile: &armcompute.DiagnosticsProfile{
							BootDiagnostics: nil,
						},
					},
				},
			},
			want: false,
		},
		{
			name: "Correct input",
			args: args{
				vm: &armcompute.VirtualMachine{
					ID:   &ID,
					Name: &name,
					Properties: &armcompute.VirtualMachineProperties{
						DiagnosticsProfile: &armcompute.DiagnosticsProfile{
							BootDiagnostics: &armcompute.BootDiagnostics{
								Enabled: &enabledTrue,
							},
						},
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isBootDiagnosticEnabled(tt.args.vm))
		})
	}
}

func Test_bootLogOutput(t *testing.T) {
	ID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1"
	name := "vm1"
	enabledTrue := true
	enabledFalse := false
	storageUri := "https://logstoragevm1.blob.core.windows.net/"
	emptyStorageUri := ""

	type args struct {
		vm *armcompute.VirtualMachine
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty input",
			args: args{
				vm: nil,
			},
			want: "",
		},
		{
			name: "StorageURI is nil",
			args: args{
				vm: &armcompute.VirtualMachine{
					ID:   &ID,
					Name: &name,
					Properties: &armcompute.VirtualMachineProperties{
						DiagnosticsProfile: &armcompute.DiagnosticsProfile{
							BootDiagnostics: &armcompute.BootDiagnostics{
								Enabled:    &enabledFalse,
								StorageURI: nil,
							},
						},
					},
				},
			},
			want: "",
		},
		{
			name: "BootDiagnostics disabled",
			args: args{
				vm: &armcompute.VirtualMachine{
					ID:   &ID,
					Name: &name,
					Properties: &armcompute.VirtualMachineProperties{
						DiagnosticsProfile: &armcompute.DiagnosticsProfile{
							BootDiagnostics: &armcompute.BootDiagnostics{
								Enabled:    &enabledFalse,
								StorageURI: &emptyStorageUri,
							},
						},
					},
				},
			},
			want: "",
		},
		{
			name: "BootDiagnostics enabled",
			args: args{
				vm: &armcompute.VirtualMachine{
					ID:   &ID,
					Name: &name,
					Properties: &armcompute.VirtualMachineProperties{
						DiagnosticsProfile: &armcompute.DiagnosticsProfile{
							BootDiagnostics: &armcompute.BootDiagnostics{
								Enabled:    &enabledTrue,
								StorageURI: &storageUri,
							},
						},
					},
				},
			},
			//want: storageUri,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, bootLogOutput(tt.args.vm))
		})
	}
}

func Test_azureComputeDiscovery_discoverBlockStorage(t *testing.T) {
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)

	type fields struct {
		azureDiscovery *azureDiscovery
	}
	tests := []struct {
		name    string
		fields  fields
		want    []ontology.IsResource
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Error list pages",
			fields: fields{
				azureDiscovery: &azureDiscovery{
					cred: nil,
				},
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrGettingNextPage.Error())
			},
		},
		{
			name: "No error",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			want: []ontology.IsResource{
				&voc.BlockStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk1",
							ServiceID:    testdata.MockCloudServiceID1,
							Name:         "disk1",
							CreationTime: util.SafeTimestamp(&creationTime),
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Type:   []string{"BlockStorage", "Storage", "Resource"},
							Labels: map[string]string{},
							Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1"),
							Raw:    "{\"*armcompute.Disk\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk1\",\"location\":\"eastus\",\"name\":\"disk1\",\"properties\":{\"encryption\":{\"diskEncryptionSetId\":\"\",\"type\":\"EncryptionAtRestWithPlatformKey\"},\"timeCreated\":\"2017-05-24T13:28:53.4540398Z\"},\"type\":\"Microsoft.Compute/disks\"}],\"*armcompute.DiskEncryptionSet\":[null]}",
						},
						AtRestEncryption: &voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
						Backups: []*voc.Backup{
							{
								Enabled:         false,
								RetentionPeriod: -1,
								Interval:        -1,
							},
						},
					},
				},
				&voc.BlockStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk2",
							ServiceID:    testdata.MockCloudServiceID1,
							Name:         "disk2",
							CreationTime: util.SafeTimestamp(&creationTime),
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Type:   []string{"BlockStorage", "Storage", "Resource"},
							Labels: map[string]string{},
							Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1"),
							Raw:    "{\"*armcompute.Disk\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk2\",\"location\":\"eastus\",\"name\":\"disk2\",\"properties\":{\"encryption\":{\"diskEncryptionSetId\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault1\",\"type\":\"EncryptionAtRestWithCustomerKey\"},\"timeCreated\":\"2017-05-24T13:28:53.4540398Z\"},\"type\":\"Microsoft.Compute/disks\"}],\"*armcompute.DiskEncryptionSet\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryption-keyvault1\",\"location\":\"germanywestcentral\",\"name\":\"encryptionkeyvault1\",\"properties\":{\"activeKey\":{\"keyUrl\":\"https://keyvault1.vault.azure.net/keys/customer-key/6273gdb374jz789hjm17819283748382\",\"sourceVault\":{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.KeyVault/vaults/keyvault1\"}}},\"type\":\"Microsoft.Compute/diskEncryptionSets\"}]}",
						},
						AtRestEncryption: &voc.CustomerKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "",
								Enabled:   true,
							},
							KeyUrl: "https://keyvault1.vault.azure.net/keys/customer-key/6273gdb374jz789hjm17819283748382",
						},
						Backups: []*voc.Backup{
							{
								Enabled:         false,
								RetentionPeriod: -1,
								Interval:        -1,
							},
						},
					},
				},
				&voc.BlockStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res2/providers/Microsoft.Compute/disks/disk3",
							ServiceID:    testdata.MockCloudServiceID1,
							Name:         "disk3",
							CreationTime: util.SafeTimestamp(&creationTime),
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Labels: map[string]string{},
							Type:   voc.BlockStorageType,
							Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res2"),
							Raw:    "{\"*armcompute.Disk\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res2/providers/Microsoft.Compute/disks/disk3\",\"location\":\"eastus\",\"name\":\"disk3\",\"properties\":{\"encryption\":{\"diskEncryptionSetId\":\"\",\"type\":\"EncryptionAtRestWithPlatformKey\"},\"timeCreated\":\"2017-05-24T13:28:53.4540398Z\"},\"type\":\"Microsoft.Compute/disks\"}],\"*armcompute.DiskEncryptionSet\":[null]}",
						},
						AtRestEncryption: &voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
						Backups: []*voc.Backup{
							{
								Enabled:         false,
								RetentionPeriod: -1,
								Interval:        -1,
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			d := tt.fields.azureDiscovery

			got, err := d.discoverBlockStorages()
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equalf(t, tt.want, got, "discoverBlockStorages()")
		})
	}
}

func Test_azureComputeDiscovery_handleBlockStorage(t *testing.T) {
	encType := armcompute.EncryptionTypeEncryptionAtRestWithCustomerKey
	diskID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk1"
	diskName := "disk1"
	diskRegion := "eastus"
	encSetID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault1"
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)

	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		disk *armcompute.Disk
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *voc.BlockStorage
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Empty input",
			args: args{
				disk: nil,
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "disk is nil")
			},
		},
		{
			name: "Empty diskID",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				disk: &armcompute.Disk{
					ID: &diskID,
					Properties: &armcompute.DiskProperties{
						Encryption: &armcompute.Encryption{
							Type: &encType,
						},
					},
				},
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not get block storage properties for the atRestEncryption:")
			},
		},
		{
			name: "Empty encryptionType",
			args: args{
				disk: &armcompute.Disk{
					ID: &diskID,
					Properties: &armcompute.DiskProperties{
						Encryption: &armcompute.Encryption{
							Type: nil,
						},
					},
				},
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "error getting atRestEncryption properties of blockStorage")
			},
		},
		{
			name: "No error",
			args: args{
				disk: &armcompute.Disk{
					ID:       &diskID,
					Name:     &diskName,
					Location: &diskRegion,
					Properties: &armcompute.DiskProperties{
						Encryption: &armcompute.Encryption{
							Type:                &encType,
							DiskEncryptionSetID: &encSetID,
						},
						TimeCreated: &creationTime,
					},
				},
			},
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			want: &voc.BlockStorage{
				Storage: &voc.Storage{
					Resource: &voc.Resource{
						ID:           voc.ResourceID(diskID),
						ServiceID:    testdata.MockCloudServiceID1,
						Name:         "disk1",
						CreationTime: util.SafeTimestamp(&creationTime),
						Type:         []string{"BlockStorage", "Storage", "Resource"},
						GeoLocation: voc.GeoLocation{
							Region: "eastus",
						},
						Labels: map[string]string{},
						Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1"),
						Raw:    "{\"*armcompute.Disk\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk1\",\"location\":\"eastus\",\"name\":\"disk1\",\"properties\":{\"encryption\":{\"diskEncryptionSetId\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault1\",\"type\":\"EncryptionAtRestWithCustomerKey\"},\"timeCreated\":\"2017-05-24T13:28:53.004540398Z\"}}],\"*armcompute.DiskEncryptionSet\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryption-keyvault1\",\"location\":\"germanywestcentral\",\"name\":\"encryptionkeyvault1\",\"properties\":{\"activeKey\":{\"keyUrl\":\"https://keyvault1.vault.azure.net/keys/customer-key/6273gdb374jz789hjm17819283748382\",\"sourceVault\":{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.KeyVault/vaults/keyvault1\"}}},\"type\":\"Microsoft.Compute/diskEncryptionSets\"}]}",
					},

					AtRestEncryption: &voc.CustomerKeyEncryption{
						AtRestEncryption: &voc.AtRestEncryption{
							Algorithm: "",
							Enabled:   true,
						},
						KeyUrl: "https://keyvault1.vault.azure.net/keys/customer-key/6273gdb374jz789hjm17819283748382",
					},
					Backups: []*voc.Backup{
						{
							Enabled:         false,
							RetentionPeriod: -1,
							Interval:        -1,
						},
					},
				},
			},

			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			d := tt.fields.azureDiscovery

			got, err := d.handleBlockStorage(tt.args.disk)
			if !tt.wantErr(t, err, fmt.Sprintf("handleBlockStorage(%v)", tt.args.disk)) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_azureComputeDiscovery_blockStorageAtRestEncryption(t *testing.T) {
	encType := armcompute.EncryptionTypeEncryptionAtRestWithCustomerKey
	diskID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk1"
	diskName := "disk1"
	diskRegion := "eastus"
	encSetID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault1"
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)

	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		disk *armcompute.Disk
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    voc.IsAtRestEncryption
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Empty disk",
			args: args{},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "disk is empty")
			},
		},
		{
			name: "Error getting atRestEncryptionProperties",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				disk: &armcompute.Disk{
					ID:       &diskID,
					Name:     &diskName,
					Location: &diskRegion,
					Properties: &armcompute.DiskProperties{
						Encryption:  &armcompute.Encryption{},
						TimeCreated: &creationTime,
					},
				},
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "error getting atRestEncryption properties of blockStorage")
			},
		},
		{
			name: "No error",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				disk: &armcompute.Disk{
					ID:       &diskID,
					Name:     &diskName,
					Location: &diskRegion,
					Properties: &armcompute.DiskProperties{
						Encryption: &armcompute.Encryption{
							Type:                &encType,
							DiskEncryptionSetID: &encSetID,
						},
						TimeCreated: &creationTime,
					},
				},
			},
			want: &voc.CustomerKeyEncryption{
				AtRestEncryption: &voc.AtRestEncryption{
					Algorithm: "",
					Enabled:   true,
				},
				KeyUrl: "https://keyvault1.vault.azure.net/keys/customer-key/6273gdb374jz789hjm17819283748382",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			d := tt.fields.azureDiscovery

			got, _, err := d.blockStorageAtRestEncryption(tt.args.disk)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_azureComputeDiscovery_keyURL(t *testing.T) {
	encSetID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault1"
	encSetID2 := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault2"

	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		diskEncryptionSetID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Empty input",
			want: "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrMissingDiskEncryptionSetID)
			},
		},
		{
			name:   "Error get disc encryption set",
			fields: fields{&azureDiscovery{}},
			args: args{
				diskEncryptionSetID: encSetID,
			},
			want: "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not get key vault:")
			},
		},
		{
			name: "Empty keyURL",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				diskEncryptionSetID: encSetID2,
			},
			want: "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not get keyURL")
			},
		},
		{
			name: "No error",
			args: args{
				diskEncryptionSetID: encSetID,
			},
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			want:    "https://keyvault1.vault.azure.net/keys/customer-key/6273gdb374jz789hjm17819283748382",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			d := tt.fields.azureDiscovery

			got, _, err := d.keyURL(tt.args.diskEncryptionSetID)
			if !tt.wantErr(t, err, fmt.Sprintf("keyURL(%v)", tt.args.diskEncryptionSetID)) {
				return
			}
			assert.Equalf(t, tt.want, got, "keyURL(%v)", tt.args.diskEncryptionSetID)
		})
	}
}

func Test_diskEncryptionSetName(t *testing.T) {
	type args struct {
		diskEncryptionSetID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Correct ID",
			args: args{
				diskEncryptionSetID: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault1",
			},
			want: "encryptionkeyvault1",
		},
		{
			name: "Empty ID",
			args: args{
				diskEncryptionSetID: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, diskEncryptionSetName(tt.args.diskEncryptionSetID))
		})
	}
}

func Test_runtimeInfo(t *testing.T) {
	type args struct {
		runtime string
	}
	tests := []struct {
		name                string
		args                args
		wantRuntimeLanguage string
		wantRuntimeVersion  string
	}{
		{
			name: "Empty input",
			args: args{
				runtime: "",
			},
			wantRuntimeLanguage: "",
			wantRuntimeVersion:  "",
		},
		{
			name: "Wrong input",
			args: args{
				runtime: "TEST",
			},
			wantRuntimeLanguage: "",
			wantRuntimeVersion:  "",
		},
		{
			name: "Happy path",
			args: args{
				runtime: "PYTHON|3.8",
			},
			wantRuntimeLanguage: "PYTHON",
			wantRuntimeVersion:  "3.8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRuntimeLanguage, gotRuntimeVersion := runtimeInfo(tt.args.runtime)
			if gotRuntimeLanguage != tt.wantRuntimeLanguage {
				t.Errorf("runtimeInfo() gotRuntimeLanguage = %v, want %v", gotRuntimeLanguage, tt.wantRuntimeLanguage)
			}
			if gotRuntimeVersion != tt.wantRuntimeVersion {
				t.Errorf("runtimeInfo() gotRuntimeVersion = %v, want %v", gotRuntimeVersion, tt.wantRuntimeVersion)
			}
		})
	}
}

func Test_automaticUpdates(t *testing.T) {
	type args struct {
		vm *armcompute.VirtualMachine
	}
	tests := []struct {
		name                 string
		args                 args
		wantAutomaticUpdates *voc.AutomaticUpdates
	}{
		{
			name:                 "Empty input",
			args:                 args{},
			wantAutomaticUpdates: &voc.AutomaticUpdates{},
		},
		{
			name: "No automatic update for the given VM",
			args: args{
				vm: &armcompute.VirtualMachine{
					Properties: &armcompute.VirtualMachineProperties{
						OSProfile: &armcompute.OSProfile{
							LinuxConfiguration: &armcompute.LinuxConfiguration{
								PatchSettings: &armcompute.LinuxPatchSettings{},
							},
						},
					},
				},
			},
			wantAutomaticUpdates: &voc.AutomaticUpdates{},
		},
		{
			name: "Happy path: Linux",
			args: args{
				vm: &armcompute.VirtualMachine{
					Properties: &armcompute.VirtualMachineProperties{
						OSProfile: &armcompute.OSProfile{
							LinuxConfiguration: &armcompute.LinuxConfiguration{
								PatchSettings: &armcompute.LinuxPatchSettings{
									PatchMode: util.Ref(armcompute.LinuxVMGuestPatchModeAutomaticByPlatform),
								},
							},
						},
					},
				},
			},
			wantAutomaticUpdates: &voc.AutomaticUpdates{
				Enabled:  true,
				Interval: Duration30Days,
			},
		},
		{
			name: "Happy path: Windows",
			args: args{
				vm: &armcompute.VirtualMachine{
					Properties: &armcompute.VirtualMachineProperties{
						OSProfile: &armcompute.OSProfile{
							WindowsConfiguration: &armcompute.WindowsConfiguration{
								PatchSettings: &armcompute.PatchSettings{
									PatchMode: util.Ref(armcompute.WindowsVMGuestPatchModeAutomaticByPlatform),
								},
								EnableAutomaticUpdates: util.Ref(true),
							},
						},
					},
				},
			},
			wantAutomaticUpdates: &voc.AutomaticUpdates{
				Enabled:  true,
				Interval: Duration30Days,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotAutomaticUpdates := automaticUpdates(tt.args.vm); !reflect.DeepEqual(gotAutomaticUpdates, tt.wantAutomaticUpdates) {
				t.Errorf("automaticUpdates() = %v, want %v", gotAutomaticUpdates, tt.wantAutomaticUpdates)
			}
		})
	}
}

func Test_azureComputeDiscovery_handleWebApp(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
		clientWebApps  bool
	}
	type args struct {
		webApp *armappservice.Site
		config armappservice.WebAppsClientGetConfigurationResponse
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   ontology.IsResource
	}{
		{
			name: "Empty input",
			args: args{
				webApp: nil,
			},
			want: nil,
		},
		{
			name: "Happy path: WebApp Windows",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
				clientWebApps:  true,
			},
			args: args{
				webApp: &armappservice.Site{
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp1"),
					Name:     util.Ref("WebApp1"),
					Location: util.Ref("West Europe"),
					Tags: map[string]*string{
						"testKey1": util.Ref("testTag1"),
						"testKey2": util.Ref("testTag2"),
					},
					Kind: util.Ref("app"),
					Properties: &armappservice.SiteProperties{
						HTTPSOnly:     util.Ref(true),
						ResourceGroup: util.Ref("res1"),
						SiteConfig: &armappservice.SiteConfig{
							MinTLSVersion: util.Ref(armappservice.SupportedTLSVersionsOne2),
						},
						VirtualNetworkSubnetID: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1"),
					},
				},
				config: armappservice.WebAppsClientGetConfigurationResponse{
					SiteConfigResource: armappservice.SiteConfigResource{
						Properties: &armappservice.SiteConfig{
							MinTLSVersion:     util.Ref(armappservice.SupportedTLSVersionsOne2),
							MinTLSCipherSuite: util.Ref(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
						},
					},
				},
			},
			want: &voc.WebApp{
				Compute: &voc.Compute{
					Resource: &voc.Resource{
						ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp1",
						ServiceID:    testdata.MockCloudServiceID1,
						Name:         "WebApp1",
						CreationTime: util.SafeTimestamp(&time.Time{}),
						Type:         []string{"WebApp", "Compute", "Resource"},
						Labels: map[string]string{
							"testKey1": "testTag1",
							"testKey2": "testTag2",
						},
						GeoLocation: voc.GeoLocation{
							Region: "West Europe",
						},
						Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1"),
						Raw:    "{\"*armappservice.Site\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp1\",\"kind\":\"app\",\"location\":\"West Europe\",\"name\":\"WebApp1\",\"properties\":{\"httpsOnly\":true,\"resourceGroup\":\"res1\",\"siteConfig\":{\"minTlsVersion\":\"1.2\"},\"virtualNetworkSubnetId\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1\"},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}],\"armappservice.WebAppsClientGetConfigurationResponse\":[{\"properties\":{\"minTlsCipherSuite\":\"TLS_AES_128_GCM_SHA256\",\"minTlsVersion\":\"1.2\"}}]}",
					},
					NetworkInterfaces: []voc.ResourceID{"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1"},
					ResourceLogging: &voc.ResourceLogging{
						Logging: &voc.Logging{
							Enabled: true,
						},
					},
				},
				HttpEndpoint: &voc.HttpEndpoint{
					TransportEncryption: &voc.TransportEncryption{
						Enabled:    true,
						Enforced:   true,
						TlsVersion: constants.TLS1_2,
						Algorithm:  string(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
					},
				},
				PublicAccess: false,
				Redundancy:   &voc.Redundancy{},
			},
		},
		{
			name: "Happy path: WebApp Linux",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
				clientWebApps:  true,
			},
			args: args{
				webApp: &armappservice.Site{
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp2"),
					Name:     util.Ref("WebApp2"),
					Location: util.Ref("West Europe"),
					Tags: map[string]*string{
						"testKey1": util.Ref("testTag1"),
						"testKey2": util.Ref("testTag2"),
					},
					Kind: util.Ref("app"),
					Properties: &armappservice.SiteProperties{
						HTTPSOnly: util.Ref(false),
						SiteConfig: &armappservice.SiteConfig{
							MinTLSVersion: nil,
						},
						ResourceGroup:          util.Ref("res1"),
						VirtualNetworkSubnetID: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet2"),
					},
				},
				config: armappservice.WebAppsClientGetConfigurationResponse{
					SiteConfigResource: armappservice.SiteConfigResource{
						Properties: &armappservice.SiteConfig{},
					},
				},
			},
			want: &voc.WebApp{
				Compute: &voc.Compute{
					Resource: &voc.Resource{
						ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp2",
						ServiceID:    testdata.MockCloudServiceID1,
						Name:         "WebApp2",
						CreationTime: util.SafeTimestamp(&time.Time{}),
						Type:         []string{"WebApp", "Compute", "Resource"},
						Labels: map[string]string{
							"testKey1": "testTag1",
							"testKey2": "testTag2",
						},
						GeoLocation: voc.GeoLocation{
							Region: "West Europe",
						},
						Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1"),
						Raw:    "{\"*armappservice.Site\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp2\",\"kind\":\"app\",\"location\":\"West Europe\",\"name\":\"WebApp2\",\"properties\":{\"httpsOnly\":false,\"resourceGroup\":\"res1\",\"siteConfig\":{},\"virtualNetworkSubnetId\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet2\"},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}],\"armappservice.WebAppsClientGetConfigurationResponse\":[{\"properties\":{}}]}",
					},
					NetworkInterfaces: []voc.ResourceID{"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet2"},
					ResourceLogging: &voc.ResourceLogging{
						Logging: &voc.Logging{
							Enabled: false,
						},
					},
				},
				HttpEndpoint: &voc.HttpEndpoint{
					TransportEncryption: &voc.TransportEncryption{
						Enabled:    false,
						Enforced:   false,
						TlsVersion: "",
						Algorithm:  "",
					},
				},
				PublicAccess: false,
				Redundancy:   &voc.Redundancy{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			// Set clients if needed
			if tt.fields.clientWebApps {
				// initialize webApps client
				_ = d.initWebAppsClient()
			}

			assert.Equalf(t, tt.want, d.handleWebApp(tt.args.webApp, tt.args.config), "handleWebApps(%v)", tt.args.webApp)
		})
	}
}

func Test_getTransportEncryption(t *testing.T) {
	type args struct {
		siteProps *armappservice.SiteProperties
		config    armappservice.WebAppsClientGetConfigurationResponse
	}
	tests := []struct {
		name    string
		args    args
		wantEnc *voc.TransportEncryption
	}{
		{
			name: "Happy path: TLSVersion/CipherSuite not available",
			args: args{
				siteProps: &armappservice.SiteProperties{
					SiteConfig: &armappservice.SiteConfig{},
					HTTPSOnly:  util.Ref(false),
				},
				config: armappservice.WebAppsClientGetConfigurationResponse{
					SiteConfigResource: armappservice.SiteConfigResource{
						Properties: &armappservice.SiteConfig{},
					},
				},
			},
			wantEnc: &voc.TransportEncryption{
				Enforced:   false,
				Enabled:    false,
				TlsVersion: "",
				Algorithm:  "",
			},
		},
		{
			name: "Happy path: TLSVersion/CipherSuite available, TLS version 1.0, TLS version 1.0",
			args: args{
				siteProps: &armappservice.SiteProperties{
					SiteConfig: &armappservice.SiteConfig{},
					HTTPSOnly:  util.Ref(true),
				},
				config: armappservice.WebAppsClientGetConfigurationResponse{
					SiteConfigResource: armappservice.SiteConfigResource{
						Properties: &armappservice.SiteConfig{
							MinTLSVersion:     util.Ref(armappservice.SupportedTLSVersionsOne0),
							MinTLSCipherSuite: util.Ref(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
						},
					},
				},
			},
			wantEnc: &voc.TransportEncryption{
				Enforced:   true,
				Enabled:    true,
				TlsVersion: constants.TLS1_0,
				Algorithm:  string(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
			},
		},
		{
			name: "Happy path: TLSVersion/CipherSuite available, TLS version 1.0, TLS version 1.1",
			args: args{
				siteProps: &armappservice.SiteProperties{
					SiteConfig: &armappservice.SiteConfig{},
					HTTPSOnly:  util.Ref(true),
				},
				config: armappservice.WebAppsClientGetConfigurationResponse{
					SiteConfigResource: armappservice.SiteConfigResource{
						Properties: &armappservice.SiteConfig{
							MinTLSVersion:     util.Ref(armappservice.SupportedTLSVersionsOne1),
							MinTLSCipherSuite: util.Ref(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
						},
					},
				},
			},
			wantEnc: &voc.TransportEncryption{
				Enforced:   true,
				Enabled:    true,
				TlsVersion: constants.TLS1_1,
				Algorithm:  string(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
			},
		},
		{
			name: "Happy path: TLSVersion/CipherSuite available, TLS version 1.2",
			args: args{
				siteProps: &armappservice.SiteProperties{
					SiteConfig: &armappservice.SiteConfig{},
					HTTPSOnly:  util.Ref(true),
				},
				config: armappservice.WebAppsClientGetConfigurationResponse{
					SiteConfigResource: armappservice.SiteConfigResource{
						Properties: &armappservice.SiteConfig{
							MinTLSVersion:     util.Ref(armappservice.SupportedTLSVersionsOne2),
							MinTLSCipherSuite: util.Ref(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
						},
					},
				},
			},
			wantEnc: &voc.TransportEncryption{
				Enforced:   true,
				Enabled:    true,
				TlsVersion: constants.TLS1_2,
				Algorithm:  string(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotEnc := getTransportEncryption(tt.args.siteProps, tt.args.config); !reflect.DeepEqual(gotEnc, tt.wantEnc) {
				t.Errorf("getTransportEncryption() = %v, want %v", gotEnc, tt.wantEnc)
			}
		})
	}
}

func Test_azureComputeDiscovery_getResourceLoggingWebApp(t *testing.T) {
	type fields struct {
		azureDiscovery     *azureDiscovery
		defenderProperties map[string]*defenderProperties
		clientWebApp       bool
	}
	type args struct {
		site *armappservice.Site
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantRl *voc.ResourceLogging
	}{
		{
			name: "Input empty",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				site: nil,
			},
			wantRl: &voc.ResourceLogging{
				Logging: &voc.Logging{},
			},
		},
		{
			name: "Happy path: logging disabled",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
				clientWebApp:   true,
			},
			args: args{
				site: &armappservice.Site{
					ID:   util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp2"),
					Name: util.Ref("WebApp2"),
					Kind: util.Ref("app"),
					Properties: &armappservice.SiteProperties{
						PublicNetworkAccess:    util.Ref("Enabled"),
						ResourceGroup:          util.Ref("res1"),
						VirtualNetworkSubnetID: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet2"),
					},
				},
			},
			wantRl: &voc.ResourceLogging{
				Logging: &voc.Logging{
					Enabled: false,
				},
			},
		},
		{
			name: "Happy path: logging enabled",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
				clientWebApp:   true,
			},
			args: args{
				site: &armappservice.Site{
					ID:   util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp1"),
					Name: util.Ref("WebApp1"),
					Kind: util.Ref("app"),
					Properties: &armappservice.SiteProperties{
						PublicNetworkAccess:    util.Ref("Enabled"),
						ResourceGroup:          util.Ref("res1"),
						VirtualNetworkSubnetID: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet2"),
					},
				},
			},
			wantRl: &voc.ResourceLogging{
				Logging: &voc.Logging{
					Enabled: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery
			d.defenderProperties = tt.fields.defenderProperties

			// Set clients if needed
			if tt.fields.clientWebApp {
				// initialize webApps client
				_ = d.initWebAppsClient()
			}

			if gotRl := d.getResourceLoggingWebApps(tt.args.site); !reflect.DeepEqual(gotRl, tt.wantRl) {
				t.Errorf("azureComputeDiscovery.getResourceLoggingWebApp() = %v, want %v", gotRl, tt.wantRl)
			}
		})
	}
}

func Test_getRedundancy(t *testing.T) {
	type args struct {
		app *armappservice.Site
	}
	tests := []struct {
		name string
		args args
		want *voc.Redundancy
	}{
		{
			name: "Happy path: no redundancy",
			args: args{
				app: &armappservice.Site{
					ID:   util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp1"),
					Name: util.Ref("WebApp1"),
					Properties: &armappservice.SiteProperties{
						RedundancyMode: util.Ref(armappservice.RedundancyModeNone),
					},
				},
			},
			want: &voc.Redundancy{
				Zone: false,
				Geo:  false,
			},
		},
		{
			name: "Happy path: zone redundancy",
			args: args{
				app: &armappservice.Site{
					ID:   util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp1"),
					Name: util.Ref("WebApp1"),
					Properties: &armappservice.SiteProperties{
						RedundancyMode: util.Ref(armappservice.RedundancyModeActiveActive),
					},
				},
			},
			want: &voc.Redundancy{
				Zone: true,
				Geo:  false,
			},
		},
		{
			name: "Happy path: zone and geo redundancy",
			args: args{
				app: &armappservice.Site{
					ID:   util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp1"),
					Name: util.Ref("WebApp1"),
					Properties: &armappservice.SiteProperties{
						RedundancyMode: util.Ref(armappservice.RedundancyModeGeoRedundant),
					},
				},
			},
			want: &voc.Redundancy{
				Zone: true,
				Geo:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRedundancy(tt.args.app); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRedundancy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPublicAccessStatus(t *testing.T) {
	type args struct {
		site *armappservice.Site
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty input",
			args: args{},
			want: false,
		},
		{
			name: "Happy path: Enabled",
			args: args{
				site: &armappservice.Site{
					Properties: &armappservice.SiteProperties{
						PublicNetworkAccess: util.Ref("Enabled"),
					},
				},
			},
			want: true,
		},
		{
			name: "Happy path: Empty String",
			args: args{
				site: &armappservice.Site{
					Properties: &armappservice.SiteProperties{
						PublicNetworkAccess: util.Ref(""),
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPublicAccessStatus(tt.args.site); got != tt.want {
				t.Errorf("getPublicAccessStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getVirtualNetworkSubnetId(t *testing.T) {
	type args struct {
		site *armappservice.Site
	}
	tests := []struct {
		name string
		args args
		want []voc.ResourceID
	}{
		{
			name: "Empty input",
			args: args{},
			want: []voc.ResourceID{},
		},
		{
			name: "Happy path: with virtual network subnet ID",
			args: args{
				site: &armappservice.Site{
					Properties: &armappservice.SiteProperties{
						VirtualNetworkSubnetID: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1"),
					},
				},
			},
			want: []voc.ResourceID{voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1")},
		},
		{
			name: "Happy path: without virtual network subnet ID",
			args: args{
				site: &armappservice.Site{
					Properties: &armappservice.SiteProperties{},
				},
			},
			want: []voc.ResourceID{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getVirtualNetworkSubnetId(tt.args.site); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getVirtualNetworkSubnetId() = %v, want %v", got, tt.want)
			}
		})
	}
}
