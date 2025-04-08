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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/security/armsecurity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"google.golang.org/protobuf/types/known/durationpb"
)

type mockSender struct {
}

func newMockSender() *mockSender {
	m := &mockSender{}
	return m
}

func (mockSender) Do(req *http.Request) (res *http.Response, err error) {
	if req.URL.Path == "/subscriptions" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":             "/subscriptions/00000000-0000-0000-0000-000000000000",
					"subscriptionId": "00000000-0000-0000-0000-000000000000",
					"name":           "sub1",
					"displayName":    "displayName",
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1",
					"name": "res1",
					"tags": map[string]interface{}{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					"location": "westus",
				},
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res2",
					"name": "res2",
					"tags": map[string]interface{}{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					"location": "eastus",
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Security/pricings" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Security/pricings/VirtualMachines",
					"name": "VirtualMachines",
					"type": "Microsoft.Security/pricings",
					"properties": map[string]interface{}{
						"pricingTier": armsecurity.PricingTierStandard,
					},
				},
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Security/pricings/StorageAccounts",
					"name": "StorageAccounts",
					"type": "Microsoft.Security/pricings",
					"properties": map[string]interface{}{
						"pricingTier": armsecurity.PricingTierStandard,
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Storage/storageAccounts" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1",
					"name":     "account1",
					"location": "eastus",
					"properties": map[string]interface{}{
						"creationTime": "2017-05-24T13:28:53.004540398Z",
						"primaryEndpoints": map[string]interface{}{
							"blob": "https://account1.blob.core.windows.net/",
							"file": "https://account1.file.core.windows.net/",
						},
						"encryption": map[string]interface{}{
							"services": map[string]interface{}{
								"file": map[string]interface{}{
									"keyType":         "Account",
									"enabled":         true,
									"lastEnabledTime": "2019-12-11T20:49:31.7036140Z",
								},
								"blob": map[string]interface{}{
									"keyType":         "Account",
									"enabled":         true,
									"lastEnabledTime": "2019-12-11T20:49:31.7036140Z",
								},
							},
							"keySource": armstorage.KeySourceMicrosoftStorage,
						},
						"minimumTlsVersion":        armstorage.MinimumTLSVersionTLS12,
						"allowBlobPublicAccess":    false,
						"supportsHttpsTrafficOnly": true,
					},
				},
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2",
					"name":     "account2",
					"location": "eastus",
					"properties": map[string]interface{}{
						"creationTime": "2017-05-24T13:28:53.004540398Z",
						"primaryEndpoints": map[string]interface{}{
							"blob": "https://account1.blob.core.windows.net/",
							"file": "https://account1.file.core.windows.net/",
						},
						"encryption": map[string]interface{}{
							"services": map[string]interface{}{
								"file": map[string]interface{}{
									"keyType":         "Account",
									"enabled":         true,
									"lastEnabledTime": "2019-12-11T20:49:31.7036140Z",
								},
								"blob": map[string]interface{}{
									"keyType":         "Account",
									"enabled":         true,
									"lastEnabledTime": "2019-12-11T20:49:31.7036140Z",
								},
							},
							"keySource": armstorage.KeySourceMicrosoftKeyvault,
							"keyvaultproperties": map[string]interface{}{
								"keyvaulturi": "https://testvault.vault.azure.net/keys/testkey/123456",
							},
						},
						"minimumTlsVersion":        armstorage.MinimumTLSVersionTLS12,
						"allowBlobPublicAccess":    false,
						"supportsHttpsTrafficOnly": true,
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Storage/storageAccounts/account3" {
		return createResponse(req, map[string]interface{}{
			"value": &map[string]interface{}{
				"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account3",
				"name":     "account3",
				"location": "westus",
				"properties": map[string]interface{}{
					"creationTime": "2017-05-24T13:28:53.004540398Z",
					"primaryEndpoints": map[string]interface{}{
						"blob": "https://account3.blob.core.windows.net/",
						"file": "https://account3.file.core.windows.net/",
					},
					"encryption": map[string]interface{}{
						"services": map[string]interface{}{
							"file": map[string]interface{}{
								"keyType":         "Account",
								"enabled":         true,
								"lastEnabledTime": "2019-12-11T20:49:31.7036140Z",
							},
							"blob": map[string]interface{}{
								"keyType":         "Account",
								"enabled":         true,
								"lastEnabledTime": "2019-12-11T20:49:31.7036140Z",
							},
						},
						"keySource": armstorage.KeySourceMicrosoftStorage,
					},
					"minimumTlsVersion":        armstorage.MinimumTLSVersionTLS12,
					"allowBlobPublicAccess":    false,
					"supportsHttpsTrafficOnly": true,
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container1",
					"name": "container1",
					"type": "Microsoft.Storage/storageAccounts/blobServices/containers",
					"properties": map[string]interface{}{
						"hasImmutabilityPolicy": false,
						"publicAccess":          armstorage.PublicAccessContainer,
					},
				},
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container2",
					"name": "container2",
					"type": "Microsoft.Storage/storageAccounts/blobServices/containers",
					"properties": map[string]interface{}{
						"hasImmutabilityPolicy": false,
						"publicAccess":          armstorage.PublicAccessContainer,
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2/blobServices/default/containers" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2/blobServices/default/containers/container3",
					"name": "container3",
					"type": "Microsoft.Storage/storageAccounts/blobServices/containers",
					"properties": map[string]interface{}{
						"hasImmutabilityPolicy": false,
						"publicAccess":          armstorage.PublicAccessNone,
					},
				},
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2/blobServices/default/containers/container4",
					"name": "container4",
					"type": "Microsoft.Storage/storageAccounts/blobServices/containers",
					"properties": map[string]interface{}{
						"hasImmutabilityPolicy": false,
						"publicAccess":          armstorage.PublicAccessNone,
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/fileshare1",
					"name": "fileshare1",
					"type": "Microsoft.Storage/storageAccounts/fileServices/shares",
				},
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/fileshare2",
					"name": "fileshare2",
					"type": "Microsoft.Storage/storageAccounts/fileServices/shares",
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2/fileServices/default/shares" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{},
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
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupInstances/account1-account1-22222222-2222-2222-2222-222222222222",
					"name": "account1-account1-22222222-2222-2222-2222-222222222222",
					"properties": map[string]interface{}{
						"dataSourceInfo": map[string]interface{}{
							"resourceID":     "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1",
							"datasourceType": "Microsoft.Storage/storageAccounts/blobServices",
						},
						"policyInfo": map[string]interface{}{
							"policyId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupPolicies/backupPolicyContainer",
						},
					},
				},
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
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupPolicies/backupPolicyContainer" {
		return createResponse(req, map[string]interface{}{
			"properties": map[string]interface{}{
				"objectType": "BackupPolicy",
				"policyRules": []map[string]interface{}{
					{
						"objectType": "AzureRetentionRule",
						"lifecycles": []map[string]interface{}{
							{
								"deleteAfter": map[string]interface{}{
									"duration":   "P7D",
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
			// },
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Sql/servers" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer1",
					"name":     "SQLServer1",
					"location": "eastus",
					"properties": map[string]interface{}{
						"minimalTlsVersion": "1.2",
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer1/databases" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer1/databases/SqlDatabase1",
					"name":     "SqlDatabase1",
					"location": "eastus",
					"properties": map[string]interface{}{
						"isInfraEncryptionEnabled": true,
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer1/databases/SqlDatabase1/advancedThreatProtectionSettings" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer1/databases/SqlDatabase1/advancedThreatProtectionSettings/Default",
					"name":     "Default",
					"location": "eastus",
					"properties": map[string]interface{}{
						"state": "Enabled",
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer2/databases/SqlDatabase1/advancedThreatProtectionSettings" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer2/databases/SqlDatabase1/advancedThreatProtectionSettings/Default",
					"name":     "Default",
					"location": "eastus",
					"properties": map[string]interface{}{
						"state": "Disabled",
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.DocumentDB/databaseAccounts" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1",
					"name": "DBAccount1",
					"kind": "MongoDB",
					"type": "Microsoft.DocumentDB/databaseAccounts",
					"systemData": map[string]interface{}{
						"createdAt": "2017-05-24T13:28:53.004540398Z",
					},
					"location": "eastus",
					"tags": map[string]interface{}{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					"properties": map[string]interface{}{
						"keyVaultKeyUri": "https://testvault.vault.azure.net/keys/testkey/123456",
					},
				},
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount2",
					"name": "DBAccount2",
					"kind": "MongoDB",
					"type": "Microsoft.DocumentDB/databaseAccounts",
					"systemData": map[string]interface{}{
						"createdAt": "2017-05-24T13:28:53.004540398Z",
					},
					"location": "eastus",
					"tags": map[string]interface{}{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					"properties": map[string]interface{}{},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1",
					"name": "DBAccount1",
					"kind": "MongoDB",
					"type": "Microsoft.DocumentDB/databaseAccounts",
					"systemData": map[string]interface{}{
						"createdAt": "2017-05-24T13:28:53.004540398Z",
					},
					"location": "eastus",
					"tags": map[string]interface{}{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					"properties": map[string]interface{}{
						"keyVaultKeyUri": "https://testvault.vault.azure.net/keys/testkey/123456",
					},
				},
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount2",
					"name": "DBAccount2",
					"kind": "MongoDB",
					"type": "Microsoft.DocumentDB/databaseAccounts",
					"systemData": map[string]interface{}{
						"createdAt": "2017-05-24T13:28:53.004540398Z",
					},
					"location": "eastus",
					"tags": map[string]interface{}{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					"properties": map[string]interface{}{},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1/mongodbDatabases" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1/mongodbDatabases/mongoDB1",
					"name":     "mongoDB1",
					"location": "West Europe",
					"tags": map[string]interface{}{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					"properties": map[string]interface{}{},
				},
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1/mongodbDatabases/mongoDB2",
					"name":     "mongoDB2",
					"location": "eastus",
					"tags": map[string]interface{}{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					"properties": map[string]interface{}{},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Compute/virtualMachines" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1",
					"name":     "vm1",
					"location": "eastus",
					"properties": map[string]interface{}{
						"timeCreated": "2017-05-24T13:28:53.004540398Z",
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
						"timeCreated": "2017-05-24T13:28:53.004540398Z",
						"encryption": map[string]interface{}{
							"diskEncryptionSetId": "",
							"type":                "EncryptionAtRestWithPlatformKey",
						},
					},
					"managedBy": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1",
				},
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk2",
					"name":     "disk2",
					"type":     "Microsoft.Compute/disks",
					"location": "eastus",
					"properties": map[string]interface{}{
						"timeCreated": "2017-05-24T13:28:53.004540398Z",
						"encryption": map[string]interface{}{
							"diskEncryptionSetId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault1",
							"type":                "EncryptionAtRestWithCustomerKey",
						},
					},
					"managedBy": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1",
				},
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res2/providers/Microsoft.Compute/disks/disk3",
					"name":     "disk3",
					"type":     "Microsoft.Compute/disks",
					"location": "eastus",
					"properties": map[string]interface{}{
						"timeCreated": "2017-05-24T13:28:53.004540398Z",
						"encryption": map[string]interface{}{
							"diskEncryptionSetId": "",
							"type":                "EncryptionAtRestWithPlatformKey",
						},
					},
					"managedBy": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1",
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
						"timeCreated": "2017-05-24T13:28:53.004540398Z",
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
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Network/networkInterfaces" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkInterfaces/iface1",
					"name":     "iface1",
					"location": "eastus",
					"properties": map[string]interface{}{
						"networkSecurityGroup": map[string]interface{}{
							"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkSecurityGroups/nsg1",
							"location": "eastus",
						},
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkInterfaces" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkInterfaces/iface1",
					"name":     "iface1",
					"location": "eastus",
					"properties": map[string]interface{}{
						"networkSecurityGroup": map[string]interface{}{
							"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkSecurityGroups/nsg1",
							"location": "eastus",
						},
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkSecurityGroups/nsg1" {
		return createResponse(req, map[string]interface{}{
			"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkSecurityGroups/nsg1",
			"name":     "nsg1",
			"location": "eastus",
			"properties": map[string]interface{}{
				"securityRules": []map[string]interface{}{
					{
						"properties": map[string]interface{}{
							"access":          "Deny",
							"sourcePortRange": "*",
						},
					},
					{
						"properties": map[string]interface{}{
							"access":          "Deny",
							"sourcePortRange": "*",
						},
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.MachineLearningServices/workspaces" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.MachineLearningServices/workspaces/mlWorkspace",
					"name":     "mlWorkspace",
					"location": "westus",
					"properties": map[string]interface{}{
						"keyVault":            "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Keyvault/vaults/keyVault1",
						"applicationInsights": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.insights/components/appInsights1",
						"encryption": map[string]interface{}{
							"status": "Enabled",
							"keyVaultProperties": map[string]interface{}{
								"keyVaultArmId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Keyvault/vaults/keyVault1",
							},
						},
					},
					"systemData": map[string]interface{}{
						"createdAt": "2017-05-24T13:28:53.004540398Z",
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.MachineLearningServices/workspaces/mlWorkspace/computes" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.MachineLearningServices/workspaces/mlWorkspace/computes/compute1",
					"name":     "compute1",
					"location": "westus",
					"tags": map[string]interface{}{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					"properties": map[string]interface{}{
						"computeType": "VirtualMachine",
						"properties": map[string]interface{}{
							"computeLocation": "westus",
							"scaleSettings": map[string]interface{}{
								"minNodeCount":                0,
								"maxNodeCount":                10,
								"nodeIdleTimeBeforeScaleDown": "PT120S",
							},
							"vmPriority": "LowPriority",
							"vmSize":     "STANDARD_D2_V2",
						},
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Network/loadBalancers" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb1",
					"name":     "lb1",
					"location": "eastus",
					"properties": map[string]interface{}{
						"loadBalancingRules": []map[string]interface{}{
							{
								"properties": map[string]interface{}{
									"frontendPort": 1234,
								},
							},
							{
								"properties": map[string]interface{}{
									"frontendPort": 5678,
								},
							},
						},
						"frontendIPConfigurations": []map[string]interface{}{
							{
								"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb1/frontendIPConfigurations/b9cb3645-25d0-4288-910a-020563f63b1c",
								"name": "b9cb3645-25d0-4288-910a-020563f63b1c",
								"properties": map[string]interface{}{
									"publicIPAddress": map[string]interface{}{
										"id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/publicIPAddresses/test-b9cb3645-25d0-4288-910a-020563f63b1c",
										"properties": map[string]interface{}{
											"ipAddress": "111.222.333.444",
										},
									},
								},
							},
						},
					},
				},
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb2",
					"name":     "lb2",
					"location": "eastus",
					"properties": map[string]interface{}{
						"loadBalancingRules": []map[string]interface{}{
							{
								"properties": map[string]interface{}{
									"frontendPort": 1234,
								},
							},
							{
								"properties": map[string]interface{}{
									"frontendPort": 5678,
								},
							},
						},
						"frontendIPConfigurations": []map[string]interface{}{
							{
								"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb1/frontendIPConfigurations/b9cb3645-25d0-4288-910a-020563f63b1c",
								"name": "b9cb3645-25d0-4288-910a-020563f63b1c",
								"properties": map[string]interface{}{
									"publicIPAddress": nil,
								},
							},
						},
					},
				},
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb3",
					"name":     "lb3",
					"location": "eastus",
					"properties": map[string]interface{}{
						"loadBalancingRules": []map[string]interface{}{
							{
								"properties": map[string]interface{}{
									"frontendPort": 1234,
								},
							},
							{
								"properties": map[string]interface{}{
									"frontendPort": 5678,
								},
							},
						},
						"frontendIPConfigurations": []map[string]interface{}{
							{
								"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb1/frontendIPConfigurations/b9cb3645-25d0-4288-910a-020563f63b1c",
								"name": "b9cb3645-25d0-4288-910a-020563f63b1c",
								"properties": map[string]interface{}{
									"publicIPAddress": map[string]interface{}{
										"id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/publicIPAddresses/test-b9cb3645-25d0-4288-910a-020563f63b1d",
									},
								},
							},
						},
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/publicIPAddresses/test-b9cb3645-25d0-4288-910a-020563f63b1c" {
		return createResponse(req, map[string]interface{}{
			"properties": map[string]interface{}{
				"ipAddress": "111.222.333.444",
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/publicIPAddresses/test-b9cb3645-25d0-4288-910a-020563f63b1d" {
		return createResponse(req, map[string]interface{}{
			"properties": map[string]interface{}{
				"ipAddress": nil,
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Network/applicationGateways" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/applicationGateways/appgw1",
					"name":     "appgw1",
					"location": "eastus",
					"properties": map[string]interface{}{
						"webApplicationFirewallConfiguration": map[string]interface{}{
							"enabled": true,
						},
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/CosmosDB1/providers/Microsoft.Insights/diagnosticSettings" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/CosmosDB1",
					"name": "CosmosDB1",
					"type": "Microsoft.DocumentDB/databaseAccounts",
					"properties": map[string]interface{}{
						"logs": &[]map[string]interface{}{
							{
								"enabled": true,
							},
						},
						"workspaceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/insights-integration/providers/Microsoft.OperationalInsights/workspaces/workspace1",
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/providers/Microsoft.Insights/diagnosticSettings" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1",
					"name": "account1",
					"properties": map[string]interface{}{
						"logs": &[]map[string]interface{}{
							{
								"enabled": true,
							},
						},
						"workspaceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/insights-integration/providers/Microsoft.OperationalInsights/workspaces/workspace1",
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2/providers/Microsoft.Insights/diagnosticSettings" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2",
					"name": "account2",
					"properties": map[string]interface{}{
						"logs": &[]map[string]interface{}{
							{
								"enabled": true,
							},
						},
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account4/providers/Microsoft.Insights/diagnosticSettings" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account4",
					"name": "account4",
					"properties": map[string]interface{}{
						"logs":        &[]map[string]interface{}{},
						"workspaceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/insights-integration/providers/Microsoft.OperationalInsights/workspaces/workspace1",
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/providers/Microsoft.Insights/diagnosticSettings/blobServices/default" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1",
					"name": "account1",
					"properties": map[string]interface{}{
						"logs":        &[]map[string]interface{}{},
						"workspaceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/insights-integration/providers/Microsoft.OperationalInsights/workspaces/workspace1",
					},
				},
			},
		}, 200)
	} else {
		res, err = createResponse(req, map[string]interface{}{}, 404)
		log.Errorf("Not handling mock for %s yet", req.URL.Path)
	}

	return
}

type mockAuthorizer struct{}

func (*mockAuthorizer) GetToken(_ context.Context, _ policy.TokenRequestOptions) (azcore.AccessToken, error) {
	var token azcore.AccessToken

	return token, nil
}

func createResponse(req *http.Request, object map[string]interface{}, statusCode int) (res *http.Response, err error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)

	if err = enc.Encode(object); err != nil {
		return nil, fmt.Errorf("could not encode JSON object: %w", err)
	}

	body := io.NopCloser(buf)

	return &http.Response{
		StatusCode: statusCode,
		Body:       body,
		// We also need to fill out the request because the Azure SDK will
		// construct the error message out of this
		Request: req,
	}, nil
}

func Test_azureDiscovery_Name(t *testing.T) {
	d := NewAzureDiscovery()

	assert.Equal(t, "Azure", d.Name())
}

func TestNewAzureDiscovery(t *testing.T) {
	type args struct {
		opts []DiscoveryOption
	}
	tests := []struct {
		name string
		args args
		want discovery.Discoverer
	}{
		{
			name: "Happy path",
			args: args{},
			want: &azureDiscovery{
				ctID:               config.DefaultTargetOfEvaluationID,
				backupMap:          make(map[string]*backup),
				defenderProperties: make(map[string]*defenderProperties),
			},
		},
		{
			name: "Happy path: with target of evaluation id",
			args: args{
				opts: []DiscoveryOption{WithTargetOfEvaluationID(testdata.MockTargetOfEvaluationID1)},
			},
			want: &azureDiscovery{
				ctID:               testdata.MockTargetOfEvaluationID1,
				backupMap:          make(map[string]*backup),
				defenderProperties: make(map[string]*defenderProperties),
			},
		},
		{
			name: "Happy path: with resource group",
			args: args{
				opts: []DiscoveryOption{WithResourceGroup(testdata.MockResourceGroup)},
			},
			want: &azureDiscovery{
				rg:                 util.Ref(testdata.MockResourceGroup),
				ctID:               config.DefaultTargetOfEvaluationID,
				backupMap:          make(map[string]*backup),
				defenderProperties: make(map[string]*defenderProperties),
			},
		},
		{
			name: "Happy path: with sender",
			args: args{
				opts: []DiscoveryOption{WithSender(mockSender{})},
			},
			want: &azureDiscovery{
				clientOptions: arm.ClientOptions{
					ClientOptions: policy.ClientOptions{
						Transport: mockSender{},
					},
				},
				ctID:               config.DefaultTargetOfEvaluationID,
				backupMap:          make(map[string]*backup),
				defenderProperties: make(map[string]*defenderProperties),
			},
		},
		{
			name: "Happy path: with authorizer",
			args: args{
				opts: []DiscoveryOption{WithAuthorizer(&mockAuthorizer{})},
			},
			want: &azureDiscovery{
				cred:               &mockAuthorizer{},
				ctID:               config.DefaultTargetOfEvaluationID,
				backupMap:          make(map[string]*backup),
				defenderProperties: make(map[string]*defenderProperties),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAzureDiscovery(tt.args.opts...)
			assert.Equal(t, tt.want, got, assert.CompareAllUnexported())
		})
	}
}

func Test_azureDiscovery_List(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
	}
	tests := []struct {
		name    string
		fields  fields
		want    assert.Want[[]ontology.IsResource]
		wantErr assert.WantErr
	}{
		{
			name: "Authorize error: no credentials configured",
			fields: fields{
				&azureDiscovery{},
			},
			want: assert.Empty[[]ontology.IsResource],
			wantErr: func(t *testing.T, gotErr error) bool {
				assert.ErrorContains(t, gotErr, ErrNoCredentialsConfigured.Error())
				return assert.ErrorContains(t, gotErr, ErrCouldNotAuthenticate.Error())
			},
		},
		{
			name: "Happy path",
			fields: fields{
				NewMockAzureDiscovery(newMockSender()),
			},
			want: func(t *testing.T, got []ontology.IsResource) bool {
				return assert.True(t, len(got) > 31)
			},
			wantErr: assert.Nil[error],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery
			gotList, err := d.List()

			tt.wantErr(t, err)
			tt.want(t, gotList)
		})
	}
}

func Test_azureDiscovery_TargetOfEvaluationID(t *testing.T) {
	type fields struct {
		isAuthorized        bool
		sub                 *armsubscription.Subscription
		cred                azcore.TokenCredential
		rg                  *string
		clientOptions       arm.ClientOptions
		discovererComponent string
		clients             clients
		ctID                string
		backupMap           map[string]*backup
		defenderProperties  map[string]*defenderProperties
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy path",
			fields: fields{
				ctID: testdata.MockTargetOfEvaluationID1,
			},
			want: testdata.MockTargetOfEvaluationID1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &azureDiscovery{
				isAuthorized:        tt.fields.isAuthorized,
				sub:                 tt.fields.sub,
				cred:                tt.fields.cred,
				rg:                  tt.fields.rg,
				clientOptions:       tt.fields.clientOptions,
				discovererComponent: tt.fields.discovererComponent,
				clients:             tt.fields.clients,
				ctID:                tt.fields.ctID,
				backupMap:           tt.fields.backupMap,
				defenderProperties:  tt.fields.defenderProperties,
			}
			if got := a.TargetOfEvaluationID(); got != tt.want {
				t.Errorf("azureDiscovery.TargetOfEvaluationID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_azureDiscovery_authorize(t *testing.T) {
	type fields struct {
		isAuthorized  bool
		sub           *armsubscription.Subscription
		cred          azcore.TokenCredential
		clientOptions arm.ClientOptions
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Is authorized",
			fields: fields{
				isAuthorized: true,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, nil, err)
			},
		},
		{
			name: "No credentials configured",
			fields: fields{
				isAuthorized: false,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrNoCredentialsConfigured)
			},
		},
		{
			name: "Error getting subscriptions",
			fields: fields{
				isAuthorized: false,
				cred:         &mockAuthorizer{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrCouldNotGetSubscriptions.Error())
			},
		},
		{
			name: "Without errors",
			fields: fields{
				isAuthorized: false,
				cred:         &mockAuthorizer{},
				clientOptions: arm.ClientOptions{
					ClientOptions: policy.ClientOptions{
						Transport: mockSender{},
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &azureDiscovery{
				isAuthorized:  tt.fields.isAuthorized,
				sub:           tt.fields.sub,
				cred:          tt.fields.cred,
				clientOptions: tt.fields.clientOptions,
			}
			tt.wantErr(t, a.authorize())
		})
	}
}

func TestGetResourceGroupName(t *testing.T) {
	accountId := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account3"
	result := resourceGroupName(accountId)

	assert.Equal(t, "res1", result)
}

func Test_resourceGroupID(t *testing.T) {
	type args struct {
		ID *string
	}
	tests := []struct {
		name string
		args args
		want *string
	}{
		{
			name: "invalid",
			args: args{
				ID: util.Ref("this is not a resource ID but it should not crash the Clouditor"),
			},
			want: nil,
		},
		{
			name: "happy path",
			args: args{
				ID: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupInstances/account1-account1-22222222-2222-2222-2222-222222222222"),
			},
			want: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resourceGroupID(tt.args.ID)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_backupPolicyName(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "invalid",
			args: args{
				id: "this is not a resource ID but it should not crash the Clouditor",
			},
			want: "",
		},
		{
			name: "Happy path",
			args: args{
				id: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupPolicies/backupPolicyDisk",
			},
			want: "backupPolicyDisk",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := backupPolicyName(tt.args.id); got != tt.want {
				t.Errorf("backupPolicyName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_labels(t *testing.T) {

	testValue1 := "testValue1"
	testValue2 := "testValue2"
	testValue3 := "testValue3"

	type args struct {
		tags map[string]*string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "Empty map of tags",
			args: args{
				tags: map[string]*string{},
			},
			want: map[string]string{},
		},
		{
			name: "Tags are nil",
			args: args{},
			want: map[string]string{},
		},
		{
			name: "Valid tags",
			args: args{
				tags: map[string]*string{
					"testTag1": &testValue1,
					"testTag2": &testValue2,
					"testTag3": &testValue3,
				},
			},
			want: map[string]string{
				"testTag1": testValue1,
				"testTag2": testValue2,
				"testTag3": testValue3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, labels(tt.args.tags))
		})
	}
}

func Test_initClientWithSubID(t *testing.T) {
	var (
		subID      = "00000000-0000-0000-0000-000000000000"
		someError  = errors.New("some error")
		someClient = &armstorage.AccountsClient{}
	)

	type args struct {
		existingClient *armstorage.AccountsClient
		d              *azureDiscovery
		fun            ClientCreateFuncWithSubID[armstorage.AccountsClient]
	}
	tests := []struct {
		name       string
		args       args
		wantClient assert.Want[*armstorage.AccountsClient]
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "No error, client does not exist",
			args: args{
				existingClient: nil,
				d: &azureDiscovery{
					cred: &mockAuthorizer{},
					sub: &armsubscription.Subscription{
						SubscriptionID: &subID,
					},
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockSender{},
						},
					},
				},
				fun: armstorage.NewAccountsClient,
			},
			wantClient: assert.NotNil[*armstorage.AccountsClient],
			wantErr:    assert.NoError,
		},
		{
			name: "Some error, client does not exist",
			args: args{
				existingClient: nil,
				d: &azureDiscovery{
					cred: &mockAuthorizer{},
					sub: &armsubscription.Subscription{
						SubscriptionID: &subID,
					},
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockSender{},
						},
					},
				},
				fun: func(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*armstorage.AccountsClient, error) {
					return nil, someError
				},
			},
			wantClient: assert.Nil[*armstorage.AccountsClient],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, someError)
			},
		},
		{
			name: "No error, client already exists",
			args: args{
				existingClient: someClient,
				d: &azureDiscovery{
					cred: &mockAuthorizer{},
					sub: &armsubscription.Subscription{
						SubscriptionID: &subID,
					},
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockSender{},
						},
					},
				},
				fun: armstorage.NewAccountsClient,
			},
			wantClient: func(t *testing.T, got *armstorage.AccountsClient) bool {
				return assert.Same(t, someClient, got)
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClient, err := initClientWithSubID(tt.args.existingClient, tt.args.d, tt.args.fun)
			tt.wantErr(t, err)
			tt.wantClient(t, gotClient)
		})
	}
}

func Test_initClientWithoutSubID(t *testing.T) {
	var (
		someError  = errors.New("some error")
		someClient = &armsecurity.PricingsClient{}
	)

	type args struct {
		existingClient *armsecurity.PricingsClient
		d              *azureDiscovery
		fun            ClientCreateFuncWithoutSubID[armsecurity.PricingsClient]
	}
	tests := []struct {
		name       string
		args       args
		wantClient assert.Want[*armsecurity.PricingsClient]
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "No error, client does not exist",
			args: args{
				existingClient: nil,
				d: &azureDiscovery{
					cred: &mockAuthorizer{},
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockSender{},
						},
					},
				},
				fun: armsecurity.NewPricingsClient,
			},
			wantClient: assert.NotNil[*armsecurity.PricingsClient],
			wantErr:    assert.NoError,
		},
		{
			name: "Some error, client does not exist",
			args: args{
				existingClient: nil,
				d: &azureDiscovery{
					cred: &mockAuthorizer{},
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockSender{},
						},
					},
				},
				fun: func(credential azcore.TokenCredential, options *arm.ClientOptions) (*armsecurity.PricingsClient, error) {
					return nil, someError
				},
			},
			wantClient: assert.Nil[*armsecurity.PricingsClient],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, someError)
			},
		},
		{
			name: "No error, client already exists",
			args: args{
				existingClient: someClient,
				d: &azureDiscovery{
					cred: &mockAuthorizer{},
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockSender{},
						},
					},
				},
				fun: armsecurity.NewPricingsClient,
			},
			wantClient: func(t *testing.T, got *armsecurity.PricingsClient) bool {
				return assert.Same(t, someClient, got)
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClient, err := initClientWithoutSubID(tt.args.existingClient, tt.args.d, tt.args.fun)
			tt.wantErr(t, err)
			tt.wantClient(t, gotClient)
		})
	}
}

func Test_initClientWithoutSubscriptionID(t *testing.T) {
	var (
		someError  = errors.New("some error")
		someClient = &armmonitor.DiagnosticSettingsClient{}
	)
	type args struct {
		existingClient *armmonitor.DiagnosticSettingsClient
		d              *azureDiscovery
		fun            ClientCreateFuncWithoutSubID[armmonitor.DiagnosticSettingsClient]
	}
	tests := []struct {
		name       string
		args       args
		wantClient assert.Want[*armmonitor.DiagnosticSettingsClient]
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "No error, client does not exist",
			args: args{
				existingClient: nil,
				d: &azureDiscovery{
					cred: &mockAuthorizer{},
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockSender{},
						},
					},
				},
				fun: armmonitor.NewDiagnosticSettingsClient,
			},
			wantClient: assert.NotNil[*armmonitor.DiagnosticSettingsClient],
			wantErr:    assert.NoError,
		},
		{
			name: "Some error, client does not exist",
			args: args{
				existingClient: nil,
				d: &azureDiscovery{
					cred: &mockAuthorizer{},
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockSender{},
						},
					},
				},
				fun: func(credential azcore.TokenCredential, options *arm.ClientOptions) (*armmonitor.DiagnosticSettingsClient, error) {
					return nil, someError
				},
			},
			wantClient: assert.Nil[*armmonitor.DiagnosticSettingsClient],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, someError)
			},
		},
		{
			name: "No error, client already exists",
			args: args{
				existingClient: someClient,
				d: &azureDiscovery{
					cred: &mockAuthorizer{},

					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockSender{},
						},
					},
				},
				fun: armmonitor.NewDiagnosticSettingsClient,
			},
			wantClient: func(t *testing.T, got *armmonitor.DiagnosticSettingsClient) bool {
				return assert.Same(t, someClient, got)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClient, err := initClientWithoutSubID(tt.args.existingClient, tt.args.d, tt.args.fun)
			tt.wantErr(t, err)
			tt.wantClient(t, gotClient)
		})
	}
}

// WithDefenderProperties is a [DiscoveryOption] that adds the defender properties for our tests.
func WithDefenderProperties(dp map[string]*defenderProperties) DiscoveryOption {
	return func(d *azureDiscovery) {
		d.defenderProperties = dp
	}
}

// WithSubscription is a [DiscoveryOption] that adds the subscription to the discoverer for our tests.
func WithSubscription(sub *armsubscription.Subscription) DiscoveryOption {
	return func(d *azureDiscovery) {
		d.sub = sub
	}
}

func NewMockAzureDiscovery(transport policy.Transporter, opts ...DiscoveryOption) *azureDiscovery {
	sub := &armsubscription.Subscription{
		SubscriptionID: util.Ref(testdata.MockSubscriptionID),
		ID:             util.Ref(testdata.MockSubscriptionResourceID),
	}

	d := &azureDiscovery{
		cred: &mockAuthorizer{},
		sub:  sub,
		clientOptions: arm.ClientOptions{
			ClientOptions: policy.ClientOptions{
				Transport: transport,
			},
		},
		ctID:      testdata.MockTargetOfEvaluationID1,
		backupMap: make(map[string]*backup),
	}

	// Apply options
	for _, opt := range opts {
		opt(d)
	}

	return d
}

func Test_retentionDuration(t *testing.T) {
	type args struct {
		retention string
	}
	tests := []struct {
		name string
		args args
		want *durationpb.Duration
	}{
		{
			name: "Missing input",
			args: args{
				retention: "",
			},
			want: durationpb.New(time.Duration(0)),
		},
		{
			name: "Wrong input",
			args: args{
				retention: "TEST",
			},
			want: durationpb.New(time.Duration(0)),
		},
		{
			name: "Happy path",
			args: args{
				retention: "P30D",
			},
			want: durationpb.New(Duration30Days),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := retentionDuration(tt.args.retention)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_azureDiscovery_discoverDefender(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
	}
	tests := []struct {
		name    string
		fields  fields
		want    assert.Want[map[string]*defenderProperties]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			want: func(t *testing.T, got map[string]*defenderProperties) bool {
				want := &defenderProperties{
					monitoringLogDataEnabled: true,
					securityAlertsEnabled:    true,
				}

				return assert.Equal(t, want, got[DefenderStorageType], assert.CompareAllUnexported())
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			got, err := d.discoverDefender()

			tt.wantErr(t, err)

			if tt.want != nil {
				tt.want(t, got)
			}
		})
	}
}
