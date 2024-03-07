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
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/constants"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/security/armsecurity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/json"
)

type mockStorageSender struct {
	mockSender
}

func newMockStorageSender() *mockStorageSender {
	m := &mockStorageSender{}
	return m
}

type responseStorageAccount struct {
	Value armstorage.Account `json:"value,omitempty"`
}

func (m mockStorageSender) Do(req *http.Request) (res *http.Response, err error) {
	if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Storage/storageAccounts" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1",
					"name":     "account1",
					"location": "eastus",
					"sku": map[string]interface{}{
						"name": "Standard_LRS",
					},
					"properties": map[string]interface{}{
						"creationTime": "2017-05-24T13:28:53.4540398Z",
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
					"sku": map[string]interface{}{
						"name": "Standard_LRS",
					},
					"properties": map[string]interface{}{
						"creationTime": "2017-05-24T13:28:53.4540398Z",
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
					"creationTime": "2017-05-24T13:28:53.4540398Z",
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
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/tableServices/default/tables" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2/tableServices/default/tables" {
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
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Sql/servers" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer1",
					"name":     "SQLServer1",
					"location": "eastus",
					"properties": map[string]interface{}{
						"publicNetworkAccess": "enabled",
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
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.DocumentDB/databaseAccounts" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/CosmosDB1",
					"name": "CosmosDB1",
					"kind": "MongoDB",
					"type": "Microsoft.DocumentDB/databaseAccounts",
					"systemData": map[string]interface{}{
						"createdAt": "2017-05-24T13:28:53.4540398Z",
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
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/CosmosDB2",
					"name": "CosmosDB2",
					"kind": "MongoDB",
					"type": "Microsoft.DocumentDB/databaseAccounts",
					"systemData": map[string]interface{}{
						"createdAt": "2017-05-24T13:28:53.4540398Z",
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
	}
	return m.mockSender.Do(req)
}

func TestNewAzureStorageDiscovery(t *testing.T) {
	type args struct {
		opts []DiscoveryOption
	}
	tests := []struct {
		name string
		args args
		want discovery.Discoverer
	}{
		{
			name: "Empty input",
			args: args{
				opts: nil,
			},
			want: &azureStorageDiscovery{
				azureDiscovery: &azureDiscovery{
					discovererComponent: StorageComponent,
					csID:                discovery.DefaultCloudServiceID,
					backupMap:           make(map[string]*backup),
				},
				defenderProperties: make(map[string]*defenderProperties),
			},
		},
		{
			name: "With sender",
			args: args{
				opts: []DiscoveryOption{WithSender(mockStorageSender{})},
			},
			want: &azureStorageDiscovery{
				azureDiscovery: &azureDiscovery{
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockStorageSender{},
						},
					},
					discovererComponent: StorageComponent,
					csID:                discovery.DefaultCloudServiceID,
					backupMap:           make(map[string]*backup),
				},
				defenderProperties: make(map[string]*defenderProperties),
			},
		},
		{
			name: "With authorizer",
			args: args{
				opts: []DiscoveryOption{WithAuthorizer(&mockAuthorizer{})},
			},
			want: &azureStorageDiscovery{
				azureDiscovery: &azureDiscovery{
					cred:                &mockAuthorizer{},
					discovererComponent: StorageComponent,
					csID:                discovery.DefaultCloudServiceID,
					backupMap:           make(map[string]*backup),
				},
				defenderProperties: make(map[string]*defenderProperties),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewAzureStorageDiscovery(tt.args.opts...)
			assert.Equal(t, tt.want, d)
			assert.Equal(t, "Azure Storage Account", d.Name())
		})
	}
}

func TestAzureStorageAuthorizer(t *testing.T) {

	d := NewAzureStorageDiscovery()
	list, err := d.List()

	assert.Error(t, err)
	assert.Nil(t, list)
	assert.ErrorIs(t, err, ErrNoCredentialsConfigured)
}

func TestStorage(t *testing.T) {
	d := NewAzureStorageDiscovery(
		WithSender(&mockStorageSender{}),
		WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 13, len(list))
	assert.NotEmpty(t, d.Name())
}

func Test_azureStorageDiscovery_List(t *testing.T) {
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)

	type fields struct {
		azureDiscovery *azureDiscovery
	}
	tests := []struct {
		name     string
		fields   fields
		wantList []voc.IsCloudResource
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "Authorize error",
			fields: fields{
				azureDiscovery: &azureDiscovery{
					cred: nil,
				},
			},
			wantList: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrCouldNotAuthenticate.Error())
			},
		},
		{
			name: "Without errors",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockStorageSender()),
			},
			wantList: []voc.IsCloudResource{
				&voc.ObjectStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1/blobservices/default/containers/container1",
							ServiceID:    testdata.MockCloudServiceID1,
							Name:         "container1",
							Type:         voc.ObjectStorageType,
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1"),
							Raw:    "{\"*armstorage.Account\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1\",\"location\":\"eastus\",\"name\":\"account1\",\"properties\":{\"allowBlobPublicAccess\":false,\"creationTime\":\"2017-05-24T13:28:53.4540398Z\",\"encryption\":{\"keySource\":\"Microsoft.Storage\",\"services\":{\"blob\":{\"enabled\":true,\"keyType\":\"Account\",\"lastEnabledTime\":\"2019-12-11T20:49:31.703614Z\"},\"file\":{\"enabled\":true,\"keyType\":\"Account\",\"lastEnabledTime\":\"2019-12-11T20:49:31.703614Z\"}}},\"minimumTlsVersion\":\"TLS1_2\",\"primaryEndpoints\":{\"blob\":\"https://account1.blob.core.windows.net/\",\"file\":\"https://account1.file.core.windows.net/\"},\"supportsHttpsTrafficOnly\":true},\"sku\":{\"name\":\"Standard_LRS\"}}],\"*armstorage.ListContainerItem\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container1\",\"name\":\"container1\",\"properties\":{\"hasImmutabilityPolicy\":false,\"publicAccess\":\"Container\"},\"type\":\"Microsoft.Storage/storageAccounts/blobServices/containers\"}]}",
						},
						AtRestEncryption: &voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
						Immutability: &voc.Immutability{Enabled: false},
						Backups: []*voc.Backup{{
							Enabled:         true,
							RetentionPeriod: Duration7Days,
							Storage:         voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupInstances/account1-account1-22222222-2222-2222-2222-222222222222"),
							TransportEncryption: &voc.TransportEncryption{
								Enforced:   true,
								Enabled:    true,
								TlsVersion: constants.TLS1_2,
								Algorithm:  constants.TLS,
							},
						},
						},
						ResourceLogging: &voc.ResourceLogging{
							Logging: &voc.Logging{
								MonitoringLogDataEnabled: true,
								SecurityAlertsEnabled:    true,
							},
						},
						Redundancy: &voc.Redundancy{
							Local: true,
						},
					},
					ContainerPublicAccess: true,
				},
				&voc.ObjectStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1/blobservices/default/containers/container2",
							ServiceID:    testdata.MockCloudServiceID1,
							Name:         "container2",
							Type:         voc.ObjectStorageType,
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1"),
							Raw:    "{\"*armstorage.Account\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1\",\"location\":\"eastus\",\"name\":\"account1\",\"properties\":{\"allowBlobPublicAccess\":false,\"creationTime\":\"2017-05-24T13:28:53.4540398Z\",\"encryption\":{\"keySource\":\"Microsoft.Storage\",\"services\":{\"blob\":{\"enabled\":true,\"keyType\":\"Account\",\"lastEnabledTime\":\"2019-12-11T20:49:31.703614Z\"},\"file\":{\"enabled\":true,\"keyType\":\"Account\",\"lastEnabledTime\":\"2019-12-11T20:49:31.703614Z\"}}},\"minimumTlsVersion\":\"TLS1_2\",\"primaryEndpoints\":{\"blob\":\"https://account1.blob.core.windows.net/\",\"file\":\"https://account1.file.core.windows.net/\"},\"supportsHttpsTrafficOnly\":true},\"sku\":{\"name\":\"Standard_LRS\"}}],\"*armstorage.ListContainerItem\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container2\",\"name\":\"container2\",\"properties\":{\"hasImmutabilityPolicy\":false,\"publicAccess\":\"Container\"},\"type\":\"Microsoft.Storage/storageAccounts/blobServices/containers\"}]}",
						},
						AtRestEncryption: &voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
						Immutability: &voc.Immutability{Enabled: false},
						Backups: []*voc.Backup{{
							Enabled:         true,
							RetentionPeriod: Duration7Days,
							Storage:         voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupInstances/account1-account1-22222222-2222-2222-2222-222222222222"),
							TransportEncryption: &voc.TransportEncryption{
								Enforced:   true,
								Enabled:    true,
								TlsVersion: constants.TLS1_2,
								Algorithm:  constants.TLS,
							},
						}},
						ResourceLogging: &voc.ResourceLogging{
							Logging: &voc.Logging{
								MonitoringLogDataEnabled: true,
								SecurityAlertsEnabled:    true,
							},
						},
						Redundancy: &voc.Redundancy{
							Local: true,
						},
					},
					ContainerPublicAccess: true,
				},
				&voc.FileStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1/fileservices/default/shares/fileshare1",
							ServiceID:    testdata.MockCloudServiceID1,
							Name:         "fileshare1",
							Type:         voc.FileStorageType,
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1"),
							Raw:    "{\"*armstorage.Account\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1\",\"location\":\"eastus\",\"name\":\"account1\",\"properties\":{\"allowBlobPublicAccess\":false,\"creationTime\":\"2017-05-24T13:28:53.4540398Z\",\"encryption\":{\"keySource\":\"Microsoft.Storage\",\"services\":{\"blob\":{\"enabled\":true,\"keyType\":\"Account\",\"lastEnabledTime\":\"2019-12-11T20:49:31.703614Z\"},\"file\":{\"enabled\":true,\"keyType\":\"Account\",\"lastEnabledTime\":\"2019-12-11T20:49:31.703614Z\"}}},\"minimumTlsVersion\":\"TLS1_2\",\"primaryEndpoints\":{\"blob\":\"https://account1.blob.core.windows.net/\",\"file\":\"https://account1.file.core.windows.net/\"},\"supportsHttpsTrafficOnly\":true},\"sku\":{\"name\":\"Standard_LRS\"}}],\"*armstorage.FileShareItem\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/fileshare1\",\"name\":\"fileshare1\",\"type\":\"Microsoft.Storage/storageAccounts/fileServices/shares\"}]}",
						},
						AtRestEncryption: &voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
						ResourceLogging: &voc.ResourceLogging{
							Logging: &voc.Logging{
								MonitoringLogDataEnabled: true,
								SecurityAlertsEnabled:    true,
							},
						},
						Redundancy: &voc.Redundancy{
							Local: true,
						},
					},
				},
				&voc.FileStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1/fileservices/default/shares/fileshare2",
							ServiceID:    testdata.MockCloudServiceID1,
							Name:         "fileshare2",
							Type:         voc.FileStorageType,
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1"),
							Raw:    "{\"*armstorage.Account\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1\",\"location\":\"eastus\",\"name\":\"account1\",\"properties\":{\"allowBlobPublicAccess\":false,\"creationTime\":\"2017-05-24T13:28:53.4540398Z\",\"encryption\":{\"keySource\":\"Microsoft.Storage\",\"services\":{\"blob\":{\"enabled\":true,\"keyType\":\"Account\",\"lastEnabledTime\":\"2019-12-11T20:49:31.703614Z\"},\"file\":{\"enabled\":true,\"keyType\":\"Account\",\"lastEnabledTime\":\"2019-12-11T20:49:31.703614Z\"}}},\"minimumTlsVersion\":\"TLS1_2\",\"primaryEndpoints\":{\"blob\":\"https://account1.blob.core.windows.net/\",\"file\":\"https://account1.file.core.windows.net/\"},\"supportsHttpsTrafficOnly\":true},\"sku\":{\"name\":\"Standard_LRS\"}}],\"*armstorage.FileShareItem\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/fileshare2\",\"name\":\"fileshare2\",\"type\":\"Microsoft.Storage/storageAccounts/fileServices/shares\"}]}",
						},
						AtRestEncryption: &voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
						ResourceLogging: &voc.ResourceLogging{
							Logging: &voc.Logging{
								MonitoringLogDataEnabled: true,
								SecurityAlertsEnabled:    true,
							},
						},
						Redundancy: &voc.Redundancy{
							Local: true,
						},
					},
				},
				&voc.ObjectStorageService{
					StorageService: &voc.StorageService{
						Storage: []voc.ResourceID{
							"/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1/blobservices/default/containers/container1",
							"/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1/blobservices/default/containers/container2",
							"/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1/fileservices/default/shares/fileshare1",
							"/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1/fileservices/default/shares/fileshare2",
						},
						NetworkService: &voc.NetworkService{
							Networking: &voc.Networking{
								Resource: &voc.Resource{
									ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1",
									ServiceID:    testdata.MockCloudServiceID1,
									Name:         "account1",
									Type:         voc.ObjectStorageServiceType,
									CreationTime: util.SafeTimestamp(&creationTime),
									Labels:       map[string]string{},
									GeoLocation: voc.GeoLocation{
										Region: "eastus",
									},
									Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
									Raw:    "{\"*armstorage.Account\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1\",\"location\":\"eastus\",\"name\":\"account1\",\"properties\":{\"allowBlobPublicAccess\":false,\"creationTime\":\"2017-05-24T13:28:53.4540398Z\",\"encryption\":{\"keySource\":\"Microsoft.Storage\",\"services\":{\"blob\":{\"enabled\":true,\"keyType\":\"Account\",\"lastEnabledTime\":\"2019-12-11T20:49:31.703614Z\"},\"file\":{\"enabled\":true,\"keyType\":\"Account\",\"lastEnabledTime\":\"2019-12-11T20:49:31.703614Z\"}}},\"minimumTlsVersion\":\"TLS1_2\",\"primaryEndpoints\":{\"blob\":\"https://account1.blob.core.windows.net/\",\"file\":\"https://account1.file.core.windows.net/\"},\"supportsHttpsTrafficOnly\":true},\"sku\":{\"name\":\"Standard_LRS\"}}]}",
								},
							},
							TransportEncryption: &voc.TransportEncryption{
								Enforced:   true,
								Enabled:    true,
								TlsVersion: constants.TLS1_2,
								Algorithm:  constants.TLS,
							},
						},
						Redundancy: &voc.Redundancy{
							Local: true,
						},
					},
					HttpEndpoint: &voc.HttpEndpoint{
						Url: "https://account1.[file,blob].core.windows.net/",
						TransportEncryption: &voc.TransportEncryption{
							Enforced:   true,
							Enabled:    true,
							TlsVersion: constants.TLS1_2,
							Algorithm:  constants.TLS,
						},
					},
				},
				&voc.ObjectStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account2/blobservices/default/containers/container3",
							ServiceID:    testdata.MockCloudServiceID1,
							Name:         "container3",
							Type:         voc.ObjectStorageType,
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account2"),
							Raw:    "{\"*armstorage.Account\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2\",\"location\":\"eastus\",\"name\":\"account2\",\"properties\":{\"allowBlobPublicAccess\":false,\"creationTime\":\"2017-05-24T13:28:53.4540398Z\",\"encryption\":{\"keySource\":\"Microsoft.Keyvault\",\"keyvaultproperties\":{\"keyvaulturi\":\"https://testvault.vault.azure.net/keys/testkey/123456\"},\"services\":{\"blob\":{\"enabled\":true,\"keyType\":\"Account\",\"lastEnabledTime\":\"2019-12-11T20:49:31.703614Z\"},\"file\":{\"enabled\":true,\"keyType\":\"Account\",\"lastEnabledTime\":\"2019-12-11T20:49:31.703614Z\"}}},\"minimumTlsVersion\":\"TLS1_2\",\"primaryEndpoints\":{\"blob\":\"https://account1.blob.core.windows.net/\",\"file\":\"https://account1.file.core.windows.net/\"},\"supportsHttpsTrafficOnly\":true},\"sku\":{\"name\":\"Standard_LRS\"}}],\"*armstorage.ListContainerItem\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2/blobServices/default/containers/container3\",\"name\":\"container3\",\"properties\":{\"hasImmutabilityPolicy\":false,\"publicAccess\":\"None\"},\"type\":\"Microsoft.Storage/storageAccounts/blobServices/containers\"}]}",
						},
						AtRestEncryption: &voc.CustomerKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "",
								Enabled:   true,
							},
							KeyUrl: "https://testvault.vault.azure.net/keys/testkey/123456",
						},
						Immutability: &voc.Immutability{Enabled: false},
						ResourceLogging: &voc.ResourceLogging{
							Logging: &voc.Logging{
								MonitoringLogDataEnabled: true,
								SecurityAlertsEnabled:    true,
							},
						},
						Backups: []*voc.Backup{
							{
								Enabled:         false,
								RetentionPeriod: -1,
								Interval:        -1,
							},
						},
						Redundancy: &voc.Redundancy{
							Local: true,
						},
					},
				},
				&voc.ObjectStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account2/blobservices/default/containers/container4",
							ServiceID:    testdata.MockCloudServiceID1,
							Name:         "container4",
							Type:         voc.ObjectStorageType,
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account2"),
							Raw:    "{\"*armstorage.Account\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2\",\"location\":\"eastus\",\"name\":\"account2\",\"properties\":{\"allowBlobPublicAccess\":false,\"creationTime\":\"2017-05-24T13:28:53.4540398Z\",\"encryption\":{\"keySource\":\"Microsoft.Keyvault\",\"keyvaultproperties\":{\"keyvaulturi\":\"https://testvault.vault.azure.net/keys/testkey/123456\"},\"services\":{\"blob\":{\"enabled\":true,\"keyType\":\"Account\",\"lastEnabledTime\":\"2019-12-11T20:49:31.703614Z\"},\"file\":{\"enabled\":true,\"keyType\":\"Account\",\"lastEnabledTime\":\"2019-12-11T20:49:31.703614Z\"}}},\"minimumTlsVersion\":\"TLS1_2\",\"primaryEndpoints\":{\"blob\":\"https://account1.blob.core.windows.net/\",\"file\":\"https://account1.file.core.windows.net/\"},\"supportsHttpsTrafficOnly\":true},\"sku\":{\"name\":\"Standard_LRS\"}}],\"*armstorage.ListContainerItem\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2/blobServices/default/containers/container4\",\"name\":\"container4\",\"properties\":{\"hasImmutabilityPolicy\":false,\"publicAccess\":\"None\"},\"type\":\"Microsoft.Storage/storageAccounts/blobServices/containers\"}]}",
						},
						AtRestEncryption: &voc.CustomerKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "",
								Enabled:   true,
							},
							KeyUrl: "https://testvault.vault.azure.net/keys/testkey/123456",
						},
						Immutability: &voc.Immutability{Enabled: false},
						ResourceLogging: &voc.ResourceLogging{
							Logging: &voc.Logging{
								MonitoringLogDataEnabled: true,
								SecurityAlertsEnabled:    true,
							},
						},
						Backups: []*voc.Backup{
							{
								Enabled:         false,
								RetentionPeriod: -1,
								Interval:        -1,
							},
						},
						Redundancy: &voc.Redundancy{
							Local: true,
						},
					},
				},
				&voc.ObjectStorageService{
					StorageService: &voc.StorageService{
						Storage: []voc.ResourceID{
							"/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account2/blobservices/default/containers/container3",
							"/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account2/blobservices/default/containers/container4",
						},
						NetworkService: &voc.NetworkService{
							Networking: &voc.Networking{
								Resource: &voc.Resource{
									ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account2",
									ServiceID:    testdata.MockCloudServiceID1,
									Name:         "account2",
									Type:         voc.ObjectStorageServiceType,
									CreationTime: util.SafeTimestamp(&creationTime),
									Labels:       map[string]string{},
									GeoLocation: voc.GeoLocation{
										Region: "eastus",
									},
									Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
									Raw:    "{\"*armstorage.Account\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2\",\"location\":\"eastus\",\"name\":\"account2\",\"properties\":{\"allowBlobPublicAccess\":false,\"creationTime\":\"2017-05-24T13:28:53.4540398Z\",\"encryption\":{\"keySource\":\"Microsoft.Keyvault\",\"keyvaultproperties\":{\"keyvaulturi\":\"https://testvault.vault.azure.net/keys/testkey/123456\"},\"services\":{\"blob\":{\"enabled\":true,\"keyType\":\"Account\",\"lastEnabledTime\":\"2019-12-11T20:49:31.703614Z\"},\"file\":{\"enabled\":true,\"keyType\":\"Account\",\"lastEnabledTime\":\"2019-12-11T20:49:31.703614Z\"}}},\"minimumTlsVersion\":\"TLS1_2\",\"primaryEndpoints\":{\"blob\":\"https://account1.blob.core.windows.net/\",\"file\":\"https://account1.file.core.windows.net/\"},\"supportsHttpsTrafficOnly\":true},\"sku\":{\"name\":\"Standard_LRS\"}}]}",
								},
							},
							TransportEncryption: &voc.TransportEncryption{
								Enforced:   true,
								Enabled:    true,
								TlsVersion: constants.TLS1_2,
								Algorithm:  constants.TLS,
							},
						},
						Redundancy: &voc.Redundancy{
							Local: true,
						},
					},

					HttpEndpoint: &voc.HttpEndpoint{
						Url: "https://account1.[file,blob].core.windows.net/",
						TransportEncryption: &voc.TransportEncryption{
							Enforced:   true,
							Enabled:    true,
							TlsVersion: constants.TLS1_2,
							Algorithm:  constants.TLS,
						},
					},
				},
				&voc.ObjectStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.dataprotection/backupvaults/backupaccount1/backupinstances/account1-account1-22222222-2222-2222-2222-222222222222",
							Name:         "account1-account1-22222222-2222-2222-2222-222222222222",
							ServiceID:    testdata.MockCloudServiceID1,
							CreationTime: 0,
							Type:         voc.ObjectStorageType,
							GeoLocation: voc.GeoLocation{
								Region: "westeurope",
							},
							//Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1"),
							Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
							Raw:    "{\"*armdataprotection.BackupInstanceResource\":[{\"properties\":{\"dataSourceInfo\":{\"resourceID\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1\",\"datasourceType\":\"Microsoft.Storage/storageAccounts/blobServices\"},\"policyInfo\":{\"policyId\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupPolicies/backupPolicyContainer\"}},\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupInstances/account1-account1-22222222-2222-2222-2222-222222222222\",\"name\":\"account1-account1-22222222-2222-2222-2222-222222222222\"}],\"*armdataprotection.BackupVaultResource\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1\",\"location\":\"westeurope\",\"name\":\"backupAccount1\"}]}",
						},
					},
				},
				&voc.DatabaseService{
					StorageService: &voc.StorageService{
						NetworkService: &voc.NetworkService{
							Networking: &voc.Networking{
								Resource: &voc.Resource{
									ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.sql/servers/sqlserver1",
									ServiceID:    testdata.MockCloudServiceID1,
									Name:         "SQLServer1",
									CreationTime: 0,
									Type:         voc.DatabaseServiceType,
									GeoLocation: voc.GeoLocation{
										Region: "eastus",
									},
									Labels: make(map[string]string),
									Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
									Raw:    "{\"*armsql.Server\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer1\",\"location\":\"eastus\",\"name\":\"SQLServer1\",\"properties\":{\"publicNetworkAccess\":\"enabled\"}}]}",
								},
							},
							TransportEncryption: &voc.TransportEncryption{
								Enabled:  true,
								Enforced: true,
							},
						},
					},
				},
				&voc.DatabaseStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.sql/servers/sqlserver1/databases/sqldatabase1",
							ServiceID:    testdata.MockCloudServiceID1,
							Name:         "SqlDatabase1",
							CreationTime: 0,
							Type:         voc.DatabaseStorageType,
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Labels: make(map[string]string),
							//Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer1"),
							Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.sql/servers/sqlserver1"),
							Raw:    "{\"*armsql.Database\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer1/databases/SqlDatabase1\",\"location\":\"eastus\",\"name\":\"SqlDatabase1\",\"properties\":{\"isInfraEncryptionEnabled\":true}}]}",
						},
						AtRestEncryption: &voc.AtRestEncryption{
							Enabled:   true,
							Algorithm: AES256,
						},
					},
				},

				&voc.DatabaseService{
					StorageService: &voc.StorageService{
						NetworkService: &voc.NetworkService{
							Networking: &voc.Networking{
								Resource: &voc.Resource{
									ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.documentdb/databaseaccounts/cosmosdb1",
									Name:         "CosmosDB1",
									ServiceID:    testdata.MockCloudServiceID1,
									CreationTime: util.SafeTimestamp(&creationTime),
									Type:         voc.DatabaseServiceType,
									GeoLocation: voc.GeoLocation{
										Region: "eastus",
									},
									Labels: map[string]string{
										"testKey1": "testTag1",
										"testKey2": "testTag2",
									},
									Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
									Raw:    "{\"*armcosmos.DatabaseAccountGetResults\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/CosmosDB1\",\"kind\":\"MongoDB\",\"location\":\"eastus\",\"name\":\"CosmosDB1\",\"properties\":{\"keyVaultKeyUri\":\"https://testvault.vault.azure.net/keys/testkey/123456\"},\"systemData\":{\"createdAt\":\"2017-05-24T13:28:53.4540398Z\"},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"},\"type\":\"Microsoft.DocumentDB/databaseAccounts\"}]}",
								},
							},
						},
						Redundancy: &voc.Redundancy{},
					},
				},
				&voc.DatabaseService{
					StorageService: &voc.StorageService{
						NetworkService: &voc.NetworkService{
							Networking: &voc.Networking{
								Resource: &voc.Resource{
									ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.documentdb/databaseaccounts/cosmosdb2",
									Name:         "CosmosDB2",
									ServiceID:    testdata.MockCloudServiceID1,
									CreationTime: util.SafeTimestamp(&creationTime),
									Type:         voc.DatabaseServiceType,
									GeoLocation: voc.GeoLocation{
										Region: "eastus",
									},
									Labels: map[string]string{
										"testKey1": "testTag1",
										"testKey2": "testTag2",
									},
									Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
									Raw:    "{\"*armcosmos.DatabaseAccountGetResults\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/CosmosDB2\",\"kind\":\"MongoDB\",\"location\":\"eastus\",\"name\":\"CosmosDB2\",\"properties\":{},\"systemData\":{\"createdAt\":\"2017-05-24T13:28:53.4540398Z\"},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"},\"type\":\"Microsoft.DocumentDB/databaseAccounts\"}]}",
								},
							},
						},
						Redundancy: &voc.Redundancy{},
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureStorageDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}
			gotList, err := d.List()
			if !tt.wantErr(t, err) {
				return
			}

			assert.Equal(t, tt.wantList, gotList)
		})
	}
}

func TestStorageHandleMethodsWhenInputIsInvalid(t *testing.T) {
	d := azureStorageDiscovery{&azureDiscovery{csID: testdata.MockCloudServiceID1}, make(map[string]*defenderProperties)}

	// Get mocked armstorage.Account
	reqURL := "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Storage/storageAccounts/account3"
	mockedStorageAccountObject, err := mockedStorageAccount(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	// Test method handleObjectStorage
	containerItem := armstorage.ListContainerItem{}
	handleObjectStorageRespone, err := d.handleObjectStorage(mockedStorageAccountObject, &containerItem, &voc.ActivityLogging{})
	assert.Error(t, err)
	assert.Nil(t, handleObjectStorageRespone)

	// Test method handleFileStorage
	fileShare := &armstorage.FileShareItem{}
	handleFileStorageRespone, err := d.handleFileStorage(mockedStorageAccountObject, fileShare, &voc.ActivityLogging{})
	assert.Error(t, err)
	assert.Nil(t, handleFileStorageRespone)
}

func TestStorageMethodsWhenInputIsInvalid(t *testing.T) {
	// Get mocked armstorage.Account
	reqURL := "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Storage/storageAccounts/account3"
	mockedStorageAccountObject, err := mockedStorageAccount(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	// Test method storageAtRestEncryption
	atRestEncryption, err := storageAtRestEncryption(mockedStorageAccountObject)
	assert.NoError(t, err)

	managedKeyEncryption := &voc.ManagedKeyEncryption{AtRestEncryption: &voc.AtRestEncryption{Algorithm: "AES256", Enabled: true}}
	assert.Equal(t, managedKeyEncryption, atRestEncryption)
}

func TestStorageDiscoverMethodsWhenInputIsInvalid(t *testing.T) {
	d := azureStorageDiscovery{&azureDiscovery{}, make(map[string]*defenderProperties)}

	// Get mocked armstorage.Account
	reqURL := "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Storage/storageAccounts/account3"
	mockedStorageAccountObject, err := mockedStorageAccount(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}
	// Test method discoverStorageAccounts
	discoverStorageAccountsResponse, err := d.discoverStorageAccounts()
	assert.Error(t, err)
	assert.ErrorContains(t, err, ErrGettingNextPage.Error())
	assert.Nil(t, discoverStorageAccountsResponse)

	// Test method discoverObjectStorages
	discoverObjectStoragesResponse, err := d.discoverObjectStorages(mockedStorageAccountObject, &voc.ActivityLogging{})
	assert.ErrorContains(t, err, ErrGettingNextPage.Error())
	assert.Nil(t, discoverObjectStoragesResponse)

	// Test method discoverFileStorages
	discoverFileStoragesResponse, err := d.discoverFileStorages(mockedStorageAccountObject, &voc.ActivityLogging{})
	assert.ErrorContains(t, err, ErrGettingNextPage.Error())
	assert.Nil(t, discoverFileStoragesResponse)
}

func Test_accountName(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Correct ID",
			args: args{
				id: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2/blobServices/default/containers/container4",
			},
			want: "account2",
		},
		{
			name: "Empty ID",
			args: args{
				id: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, accountName(tt.args.id))
		})
	}
}

// mockedStorageAccount returns one mocked storage account
func mockedStorageAccount(reqUrl string) (storageAccount *armstorage.Account, err error) {
	var storageAccountResponse responseStorageAccount

	m := newMockStorageSender()
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return storageAccount, fmt.Errorf("error creating new request: %w", err)
	}
	resp, err := m.Do(req)
	if err != nil {
		return storageAccount, fmt.Errorf("error getting mock http response: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("error io.ReadCloser: %w", err)
		}
	}(resp.Body)
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return storageAccount, fmt.Errorf("error read all: %w", err)
	}
	err = json.Unmarshal(responseBody, &storageAccountResponse)
	if err != nil {
		return storageAccount, fmt.Errorf("error unmarshalling: %w", err)
	}

	storageAccount = &storageAccountResponse.Value

	return storageAccount, nil
}

func Test_azureStorageDiscovery_discoverStorageAccounts(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
	}
	tests := []struct {
		name    string
		fields  fields
		want    []voc.IsCloudResource
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
				azureDiscovery: NewMockAzureDiscovery(newMockStorageSender()),
			},

			want:    nil,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureStorageDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}
			got, err := d.discoverStorageAccounts()
			if tt.wantErr != nil {
				if !tt.wantErr(t, err) {
					return
				}
			} else {
				assert.Nil(t, err)
				assert.Equal(t, 9, len(got))
			}
		})
	}
}

func Test_storageAtRestEncryption(t *testing.T) {
	keySource := armstorage.KeySourceMicrosoftStorage

	type args struct {
		account *armstorage.Account
	}
	tests := []struct {
		name    string
		args    args
		want    voc.IsAtRestEncryption
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Empty account",
			args: args{},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrEmptyStorageAccount)
			},
		},
		{
			name: "Empty KeySource",
			args: args{
				account: &armstorage.Account{
					Properties: &armstorage.AccountProperties{
						Encryption: &armstorage.Encryption{},
					},
				},
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "keySource is empty")
			},
		},
		{
			name: "No error",
			args: args{
				account: &armstorage.Account{
					Properties: &armstorage.AccountProperties{
						Encryption: &armstorage.Encryption{
							KeySource: &keySource,
						},
					},
				},
			},
			want: &voc.ManagedKeyEncryption{
				AtRestEncryption: &voc.AtRestEncryption{
					Algorithm: "AES256",
					Enabled:   true,
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := storageAtRestEncryption(tt.args.account)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_handleFileStorage(t *testing.T) {
	accountID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1"
	fileShareID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/fileshare1"
	fileShareName := "fileShare1"
	accountRegion := "eastus"
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)
	keySource := armstorage.KeySourceMicrosoftStorage

	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		account         *armstorage.Account
		fileshare       *armstorage.FileShareItem
		activityLogging *voc.ActivityLogging
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *voc.FileStorage
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Empty account",
			args: args{},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrEmptyStorageAccount)
			},
		},
		{
			name: "Empty fileShareItem",
			args: args{
				account: &armstorage.Account{
					ID: &accountID,
				},
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "fileshare is nil")
			},
		},
		{
			name: "Error getting atRestEncryption properties",
			args: args{
				account: &armstorage.Account{
					ID: &accountID,
				},
				fileshare: &armstorage.FileShareItem{
					ID: &fileShareID,
				},
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not get file storage properties for the atRestEncryption:")
			},
		},
		{
			name: "No error",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockStorageSender()),
			},
			args: args{
				account: &armstorage.Account{
					ID: util.Ref(accountID),
					Properties: &armstorage.AccountProperties{
						Encryption: &armstorage.Encryption{
							KeySource: &keySource,
						},
						CreationTime: &creationTime,
					},
					SKU: &armstorage.SKU{
						Name: util.Ref(armstorage.SKUNamePremiumLRS),
					},
					Location: &accountRegion,
				},
				fileshare: &armstorage.FileShareItem{
					ID:   &fileShareID,
					Name: &fileShareName,
				},
				activityLogging: &voc.ActivityLogging{
					Logging: &voc.Logging{
						Enabled: true,
					},
				},
			},
			want: &voc.FileStorage{
				Storage: &voc.Storage{
					Resource: &voc.Resource{
						ID:           voc.ResourceID(strings.ToLower(fileShareID)),
						ServiceID:    testdata.MockCloudServiceID1,
						Name:         fileShareName,
						CreationTime: util.SafeTimestamp(&creationTime),
						GeoLocation: voc.GeoLocation{
							Region: accountRegion,
						},
						Labels: map[string]string{},
						Type:   voc.FileStorageType,
						Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1"),
						Raw:    "{\"*armstorage.Account\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1\",\"location\":\"eastus\",\"properties\":{\"creationTime\":\"2017-05-24T13:28:53.004540398Z\",\"encryption\":{\"keySource\":\"Microsoft.Storage\"}},\"sku\":{\"name\":\"Premium_LRS\"}}],\"*armstorage.FileShareItem\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/fileshare1\",\"name\":\"fileShare1\"}]}",
					},
					AtRestEncryption: &voc.ManagedKeyEncryption{
						AtRestEncryption: &voc.AtRestEncryption{
							Algorithm: "AES256",
							Enabled:   true,
						},
					},
					ResourceLogging: &voc.ResourceLogging{
						Logging: &voc.Logging{
							Enabled:                  false,
							MonitoringLogDataEnabled: false,
							SecurityAlertsEnabled:    false,
						},
					},
					ActivityLogging: &voc.ActivityLogging{
						Logging: &voc.Logging{
							Enabled: true,
						},
					},
					Redundancy: &voc.Redundancy{
						Local: true,
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureStorageDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}

			got, err := d.handleFileStorage(tt.args.account, tt.args.fileshare, tt.args.activityLogging)
			if !tt.wantErr(t, err, fmt.Sprintf("handleFileStorage(%v, %v)", tt.args.account, tt.args.fileshare)) {
				return
			}
			assert.Equalf(t, tt.want, got, "handleFileStorage(%v, %v)", tt.args.account, tt.args.fileshare)
		})
	}
}

func Test_generalizeURL(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty input",
			want: "",
		},
		{
			name: "Correct input",
			args: args{
				url: "https://account1.file.core.windows.net/",
			},
			want: "https://account1.[file,blob].core.windows.net/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, generalizeURL(tt.args.url))
		})
	}
}

func Test_azureStorageDiscovery_handleStorageAccount(t *testing.T) {
	accountID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1"
	accountName := "account1"
	keySource := armstorage.KeySourceMicrosoftStorage
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)
	accountRegion := "eastus"
	minTLS := armstorage.MinimumTLSVersionTLS12
	endpointURL := "https://account1.blob.core.windows.net"
	httpsOnly := true

	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		account         *armstorage.Account
		storagesList    []voc.IsCloudResource
		activityLogging *voc.ActivityLogging
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *voc.ObjectStorageService
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Account is empty",
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrEmptyStorageAccount)
			},
		},
		{
			name: "No error",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockStorageSender()),
			},
			args: args{
				account: &armstorage.Account{
					ID:   &accountID,
					Name: &accountName,
					SKU: &armstorage.SKU{
						Name: util.Ref(armstorage.SKUNameStandardLRS),
					},
					Properties: &armstorage.AccountProperties{
						Encryption: &armstorage.Encryption{
							KeySource: &keySource,
						},
						MinimumTLSVersion: &minTLS,
						CreationTime:      &creationTime,
						PrimaryEndpoints: &armstorage.Endpoints{
							Blob: &endpointURL,
						},
						EnableHTTPSTrafficOnly: &httpsOnly,
					},
					Location: &accountRegion,
				},
				activityLogging: &voc.ActivityLogging{
					Logging: &voc.Logging{
						Enabled: true,
					},
				},
			},
			want: &voc.ObjectStorageService{
				StorageService: &voc.StorageService{
					NetworkService: &voc.NetworkService{
						Networking: &voc.Networking{
							Resource: &voc.Resource{
								ID:           voc.ResourceID(strings.ToLower(accountID)),
								ServiceID:    testdata.MockCloudServiceID1,
								Name:         accountName,
								CreationTime: util.SafeTimestamp(&creationTime),
								Type:         voc.ObjectStorageServiceType,
								GeoLocation: voc.GeoLocation{
									Region: accountRegion,
								},
								Labels: map[string]string{},
								Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
								Raw:    "{\"*armstorage.Account\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1\",\"location\":\"eastus\",\"name\":\"account1\",\"properties\":{\"creationTime\":\"2017-05-24T13:28:53.004540398Z\",\"encryption\":{\"keySource\":\"Microsoft.Storage\"},\"minimumTlsVersion\":\"TLS1_2\",\"primaryEndpoints\":{\"blob\":\"https://account1.blob.core.windows.net\"},\"supportsHttpsTrafficOnly\":true},\"sku\":{\"name\":\"Standard_LRS\"}}]}",
							},
						},
						TransportEncryption: &voc.TransportEncryption{
							Enforced:   true,
							Enabled:    true,
							TlsVersion: constants.TLS1_2,
							Algorithm:  constants.TLS,
						},
					},
					Redundancy: &voc.Redundancy{
						Local: true,
					},
					ActivityLogging: &voc.ActivityLogging{
						Logging: &voc.Logging{
							Enabled: true,
						},
					},
				},
				HttpEndpoint: &voc.HttpEndpoint{
					Url: "https://account1.[file,blob].core.windows.net",
					TransportEncryption: &voc.TransportEncryption{
						Enforced:   true,
						Enabled:    true,
						TlsVersion: constants.TLS1_2,
						Algorithm:  constants.TLS,
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			az := &azureStorageDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}
			got, err := az.handleStorageAccount(tt.args.account, tt.args.storagesList, tt.args.activityLogging)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_handleObjectStorage(t *testing.T) {
	accountID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1"
	accountRegion := "eastus"
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)
	keySource := armstorage.KeySourceMicrosoftStorage
	containerID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container1"
	containerName := "container1"
	immutability := false
	publicAccess := armstorage.PublicAccessNone

	type fields struct {
		azureDiscovery      *azureDiscovery
		blobContainerClient bool
	}
	type args struct {
		account         *armstorage.Account
		container       *armstorage.ListContainerItem
		activityLogging *voc.ActivityLogging
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *voc.ObjectStorage
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Account is empty",
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrEmptyStorageAccount)
			},
		},
		{
			name: "Container is empty",
			args: args{
				account: &armstorage.Account{
					ID: &accountID,
				},
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "container is nil")
			},
		},
		{
			name: "Error getting atRestEncryption properties",
			args: args{
				account: &armstorage.Account{
					ID: &accountID,
				},
				container: &armstorage.ListContainerItem{
					ID: &containerID,
				},
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not get object storage properties for the atRestEncryption:")
			},
		},
		{
			name: "No error",
			fields: fields{
				azureDiscovery:      NewMockAzureDiscovery(newMockStorageSender()),
				blobContainerClient: true,
			},
			args: args{
				account: &armstorage.Account{
					ID: &accountID,
					SKU: &armstorage.SKU{
						Name: util.Ref(armstorage.SKUNamePremiumLRS),
					},
					Properties: &armstorage.AccountProperties{
						Encryption: &armstorage.Encryption{
							KeySource: &keySource,
						},
						CreationTime: &creationTime,
					},
					Location: &accountRegion,
				},
				container: &armstorage.ListContainerItem{
					ID:   &containerID,
					Name: &containerName,
					Properties: &armstorage.ContainerProperties{
						HasImmutabilityPolicy: &immutability,
						PublicAccess:          &publicAccess,
					},
				},
				activityLogging: &voc.ActivityLogging{
					Logging: &voc.Logging{
						Enabled: true,
					},
				},
			},
			want: &voc.ObjectStorage{
				Storage: &voc.Storage{
					Resource: &voc.Resource{
						ID:           voc.ResourceID(strings.ToLower(containerID)),
						ServiceID:    testdata.MockCloudServiceID1,
						Name:         containerName,
						CreationTime: util.SafeTimestamp(&creationTime),
						GeoLocation: voc.GeoLocation{
							Region: accountRegion,
						},
						Labels: map[string]string{},
						Type:   voc.ObjectStorageType,
						Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1"),
						Raw:    "{\"*armstorage.Account\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1\",\"location\":\"eastus\",\"properties\":{\"creationTime\":\"2017-05-24T13:28:53.004540398Z\",\"encryption\":{\"keySource\":\"Microsoft.Storage\"}},\"sku\":{\"name\":\"Premium_LRS\"}}],\"*armstorage.ListContainerItem\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container1\",\"name\":\"container1\",\"properties\":{\"hasImmutabilityPolicy\":false,\"publicAccess\":\"None\"}}]}",
					},
					AtRestEncryption: &voc.ManagedKeyEncryption{
						AtRestEncryption: &voc.AtRestEncryption{
							Algorithm: "AES256",
							Enabled:   true,
						},
					},
					Immutability: &voc.Immutability{Enabled: false},
					ResourceLogging: &voc.ResourceLogging{
						Logging: &voc.Logging{
							MonitoringLogDataEnabled: false,
							SecurityAlertsEnabled:    false,
						},
					},
					ActivityLogging: &voc.ActivityLogging{
						Logging: &voc.Logging{
							Enabled: true,
						},
					},
					Backups: []*voc.Backup{
						{
							Enabled:         false,
							RetentionPeriod: -1,
							Interval:        -1,
						},
					},
					Redundancy: &voc.Redundancy{
						Local: true,
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureStorageDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}

			// initialize blob container client
			if tt.fields.blobContainerClient {
				_ = d.initBlobContainerClient()
			}

			got, err := d.handleObjectStorage(tt.args.account, tt.args.container, tt.args.activityLogging)
			if !tt.wantErr(t, err, fmt.Sprintf("handleObjectStorage(%v, %v)", tt.args.account, tt.args.container)) {
				return
			}
			assert.Equalf(t, tt.want, got, "handleObjectStorage(%v, %v)", tt.args.account, tt.args.container)
		})
	}
}

func Test_azureStorageDiscovery_discoverFileStorages(t *testing.T) {
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)
	accountID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1"
	accountName := "account1"
	accountRegion := "eastus"
	keySource := armstorage.KeySourceMicrosoftStorage

	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		account         *armstorage.Account
		activityLogging *voc.ActivityLogging
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []voc.IsCloudResource
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Error list pages",
			fields: fields{
				azureDiscovery: &azureDiscovery{
					cred: nil,
				},
			},
			args: args{
				account: &armstorage.Account{
					ID: &accountID,
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
				azureDiscovery: NewMockAzureDiscovery(newMockStorageSender()),
			},
			args: args{
				account: &armstorage.Account{
					ID:   &accountID,
					Name: &accountName,
					SKU: &armstorage.SKU{
						Name: util.Ref(armstorage.SKUNamePremiumLRS),
					},
					Properties: &armstorage.AccountProperties{
						Encryption: &armstorage.Encryption{
							KeySource: &keySource,
						},
						CreationTime: &creationTime,
					},
					Location: &accountRegion,
				},
				activityLogging: &voc.ActivityLogging{
					Logging: &voc.Logging{
						Enabled: true,
					},
				},
			},
			want: []voc.IsCloudResource{
				&voc.FileStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1/fileservices/default/shares/fileshare1",
							ServiceID:    testdata.MockCloudServiceID1,
							Name:         "fileshare1",
							Type:         voc.FileStorageType,
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1"),
							Raw:    "{\"*armstorage.Account\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1\",\"location\":\"eastus\",\"name\":\"account1\",\"properties\":{\"creationTime\":\"2017-05-24T13:28:53.004540398Z\",\"encryption\":{\"keySource\":\"Microsoft.Storage\"}},\"sku\":{\"name\":\"Premium_LRS\"}}],\"*armstorage.FileShareItem\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/fileshare1\",\"name\":\"fileshare1\",\"type\":\"Microsoft.Storage/storageAccounts/fileServices/shares\"}]}",
						},
						AtRestEncryption: &voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
						ResourceLogging: &voc.ResourceLogging{
							Logging: &voc.Logging{
								MonitoringLogDataEnabled: false,
								SecurityAlertsEnabled:    false,
							},
						},
						ActivityLogging: &voc.ActivityLogging{
							Logging: &voc.Logging{
								Enabled: true,
							},
						},
						Redundancy: &voc.Redundancy{
							Local: true,
						},
					},
				},
				&voc.FileStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1/fileservices/default/shares/fileshare2",
							ServiceID:    testdata.MockCloudServiceID1,
							Name:         "fileshare2",
							Type:         voc.FileStorageType,
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1"),
							Raw:    "{\"*armstorage.Account\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1\",\"location\":\"eastus\",\"name\":\"account1\",\"properties\":{\"creationTime\":\"2017-05-24T13:28:53.004540398Z\",\"encryption\":{\"keySource\":\"Microsoft.Storage\"}},\"sku\":{\"name\":\"Premium_LRS\"}}],\"*armstorage.FileShareItem\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/fileshare2\",\"name\":\"fileshare2\",\"type\":\"Microsoft.Storage/storageAccounts/fileServices/shares\"}]}",
						},
						AtRestEncryption: &voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
						ResourceLogging: &voc.ResourceLogging{
							Logging: &voc.Logging{
								MonitoringLogDataEnabled: false,
								SecurityAlertsEnabled:    false,
							},
						},
						ActivityLogging: &voc.ActivityLogging{
							Logging: &voc.Logging{
								Enabled: true,
							},
						},
						Redundancy: &voc.Redundancy{
							Local: true,
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureStorageDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}
			// initialize file share client
			_ = d.initFileStorageClient()

			got, err := d.discoverFileStorages(tt.args.account, tt.args.activityLogging)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_azureStorageDiscovery_discoverObjectStorages(t *testing.T) {
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)
	accountID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1"
	accountName := "account1"
	accountRegion := "eastus"
	keySource := armstorage.KeySourceMicrosoftStorage

	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		account         *armstorage.Account
		activityLogging *voc.ActivityLogging
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []voc.IsCloudResource
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Error list pages",
			fields: fields{
				azureDiscovery: &azureDiscovery{
					cred: nil,
				},
			},
			args: args{
				account: &armstorage.Account{
					ID: &accountID,
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
				azureDiscovery: NewMockAzureDiscovery(newMockStorageSender()),
			},
			args: args{
				account: &armstorage.Account{
					ID:   &accountID,
					Name: &accountName,
					SKU: &armstorage.SKU{
						Name: util.Ref(armstorage.SKUNamePremiumLRS),
					},
					Properties: &armstorage.AccountProperties{
						Encryption: &armstorage.Encryption{
							KeySource: &keySource,
						},
						CreationTime: &creationTime,
						PrimaryEndpoints: &armstorage.Endpoints{
							Blob: util.Ref("blob"),
						},
					},
					Location: &accountRegion,
				},
				activityLogging: &voc.ActivityLogging{
					Logging: &voc.Logging{
						Enabled: true,
					},
				},
			},
			want: []voc.IsCloudResource{
				&voc.ObjectStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1/blobservices/default/containers/container1",
							ServiceID:    testdata.MockCloudServiceID1,
							Name:         "container1",
							Type:         voc.ObjectStorageType,
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1"),
							Raw:    "{\"*armstorage.Account\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1\",\"location\":\"eastus\",\"name\":\"account1\",\"properties\":{\"creationTime\":\"2017-05-24T13:28:53.004540398Z\",\"encryption\":{\"keySource\":\"Microsoft.Storage\"},\"primaryEndpoints\":{\"blob\":\"blob\"}},\"sku\":{\"name\":\"Premium_LRS\"}}],\"*armstorage.ListContainerItem\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container1\",\"name\":\"container1\",\"properties\":{\"hasImmutabilityPolicy\":false,\"publicAccess\":\"Container\"},\"type\":\"Microsoft.Storage/storageAccounts/blobServices/containers\"}]}",
						},
						AtRestEncryption: &voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
						Immutability: &voc.Immutability{Enabled: false},
						ResourceLogging: &voc.ResourceLogging{
							Logging: &voc.Logging{
								MonitoringLogDataEnabled: false,
								SecurityAlertsEnabled:    false,
							},
						},
						ActivityLogging: &voc.ActivityLogging{
							Logging: &voc.Logging{
								Enabled: true,
							},
						},
						Backups: []*voc.Backup{
							{
								Enabled:         false,
								RetentionPeriod: -1,
								Interval:        -1,
							},
						},
						Redundancy: &voc.Redundancy{
							Local: true,
						},
					},
					ContainerPublicAccess: true,
				},
				&voc.ObjectStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1/blobservices/default/containers/container2",
							ServiceID:    testdata.MockCloudServiceID1,
							Name:         "container2",
							Type:         voc.ObjectStorageType,
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1"),
							Raw:    "{\"*armstorage.Account\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1\",\"location\":\"eastus\",\"name\":\"account1\",\"properties\":{\"creationTime\":\"2017-05-24T13:28:53.004540398Z\",\"encryption\":{\"keySource\":\"Microsoft.Storage\"},\"primaryEndpoints\":{\"blob\":\"blob\"}},\"sku\":{\"name\":\"Premium_LRS\"}}],\"*armstorage.ListContainerItem\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container2\",\"name\":\"container2\",\"properties\":{\"hasImmutabilityPolicy\":false,\"publicAccess\":\"Container\"},\"type\":\"Microsoft.Storage/storageAccounts/blobServices/containers\"}]}",
						},
						AtRestEncryption: &voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
						Immutability: &voc.Immutability{Enabled: false},
						ResourceLogging: &voc.ResourceLogging{
							Logging: &voc.Logging{
								MonitoringLogDataEnabled: false,
								SecurityAlertsEnabled:    false,
							},
						},
						ActivityLogging: &voc.ActivityLogging{
							Logging: &voc.Logging{
								Enabled: true,
							},
						},
						Backups: []*voc.Backup{
							{
								Enabled:         false,
								RetentionPeriod: -1,
								Interval:        -1,
							},
						},
						Redundancy: &voc.Redundancy{
							Local: true,
						},
					},
					ContainerPublicAccess: true,
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureStorageDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}

			// initialize blob container client
			_ = d.initBlobContainerClient()

			got, err := d.discoverObjectStorages(tt.args.account, tt.args.activityLogging)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_azureStorageDiscovery_handleSqlServer(t *testing.T) {
	type fields struct {
		azureDiscovery     *azureDiscovery
		defenderProperties map[string]*defenderProperties
	}
	type args struct {
		server *armsql.Server
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []voc.IsCloudResource
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockStorageSender()),
			},
			args: args{
				server: &armsql.Server{
					Location: util.Ref("eastus"),
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer1"),
					Name:     util.Ref("SQLServer1"),
					Properties: &armsql.ServerProperties{
						PublicNetworkAccess: util.Ref(armsql.ServerNetworkAccessFlagDisabled),
					},
				},
			},
			want: []voc.IsCloudResource{
				&voc.DatabaseService{
					StorageService: &voc.StorageService{
						NetworkService: &voc.NetworkService{
							Networking: &voc.Networking{
								Resource: &voc.Resource{
									ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.sql/servers/sqlserver1",
									ServiceID:    testdata.MockCloudServiceID1,
									Name:         "SQLServer1",
									CreationTime: 0,
									Type:         voc.DatabaseServiceType,
									GeoLocation: voc.GeoLocation{
										Region: "eastus",
									},
									Labels: make(map[string]string),
									Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
									Raw:    "{\"*armsql.Server\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer1\",\"location\":\"eastus\",\"name\":\"SQLServer1\",\"properties\":{\"publicNetworkAccess\":\"Disabled\"}}]}",
								},
							},
							TransportEncryption: &voc.TransportEncryption{
								Enabled:  true,
								Enforced: true,
							},
						},
					},
					PublicAccess: false,
				},
				&voc.DatabaseStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.sql/servers/sqlserver1/databases/sqldatabase1",
							Name:         "SqlDatabase1",
							ServiceID:    testdata.MockCloudServiceID1,
							CreationTime: 0,
							Type:         voc.DatabaseStorageType,
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Labels: make(map[string]string),
							Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.sql/servers/sqlserver1"),
							Raw:    "{\"*armsql.Database\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer1/databases/SqlDatabase1\",\"location\":\"eastus\",\"name\":\"SqlDatabase1\",\"properties\":{\"isInfraEncryptionEnabled\":true}}]}",
						},
						AtRestEncryption: &voc.AtRestEncryption{
							Enabled:   true,
							Algorithm: AES256,
						},
					},
					PublicAccess: false,
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureStorageDiscovery{
				azureDiscovery:     tt.fields.azureDiscovery,
				defenderProperties: tt.fields.defenderProperties,
			}

			got, err := d.handleSqlServer(tt.args.server)

			if !tt.wantErr(t, err, fmt.Sprintf("handleSqlServer(%v, %v)", tt.args.server, tt.args.server)) {
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_azureStorageDiscovery_anomalyDetectionEnabled(t *testing.T) {
	type fields struct {
		azureDiscovery     *azureDiscovery
		defenderProperties map[string]*defenderProperties
	}
	type args struct {
		server *armsql.Server
		db     *armsql.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "error list pager",
			fields: fields{
				azureDiscovery: &azureDiscovery{
					clients: clients{},
				},
			},
			args: args{
				server: &armsql.Server{
					Location: util.Ref("eastus"),
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer1"),
					Name:     util.Ref("SQLServer1"),
				},
				db: &armsql.Database{
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer1/databases/SqlDatabase1"),
					Name:     util.Ref("SqlDatabase1"),
					Location: util.Ref("eastus"),
					Properties: &armsql.DatabaseProperties{
						IsInfraEncryptionEnabled: util.Ref(true),
					},
				},
			},
			want: false,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "error getting next page: ")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockStorageSender()),
			},
			args: args{
				server: &armsql.Server{
					Location: util.Ref("eastus"),
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer1"),
					Name:     util.Ref("SQLServer1"),
				},
				db: &armsql.Database{
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer1/databases/SqlDatabase1"),
					Name:     util.Ref("SqlDatabase1"),
					Location: util.Ref("eastus"),
					Properties: &armsql.DatabaseProperties{
						IsInfraEncryptionEnabled: util.Ref(true),
					},
				},
			},
			want:    true,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureStorageDiscovery{
				azureDiscovery:     tt.fields.azureDiscovery,
				defenderProperties: tt.fields.defenderProperties,
			}
			got, err := d.anomalyDetectionEnabled(tt.args.server, tt.args.db)

			tt.wantErr(t, err)
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_azureStorageDiscovery_discoverCosmosDB(t *testing.T) {
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)

	type fields struct {
		azureDiscovery           *azureDiscovery
		defenderProperties       map[string]*defenderProperties
		diagnosticSettingsClient bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    []voc.IsCloudResource
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
			name: "Happy path",
			fields: fields{
				azureDiscovery:           NewMockAzureDiscovery(newMockStorageSender()),
				diagnosticSettingsClient: true,
			},
			want: []voc.IsCloudResource{
				&voc.DatabaseService{
					StorageService: &voc.StorageService{
						NetworkService: &voc.NetworkService{
							Networking: &voc.Networking{
								Resource: &voc.Resource{
									ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.documentdb/databaseaccounts/cosmosdb1",
									Name:         "CosmosDB1",
									ServiceID:    testdata.MockCloudServiceID1,
									CreationTime: util.SafeTimestamp(&creationTime),
									Type:         voc.DatabaseServiceType,
									GeoLocation: voc.GeoLocation{
										Region: "eastus",
									},
									Labels: map[string]string{
										"testKey1": "testTag1",
										"testKey2": "testTag2",
									},
									Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
									Raw:    "{\"*armcosmos.DatabaseAccountGetResults\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/CosmosDB1\",\"kind\":\"MongoDB\",\"location\":\"eastus\",\"name\":\"CosmosDB1\",\"properties\":{\"keyVaultKeyUri\":\"https://testvault.vault.azure.net/keys/testkey/123456\"},\"systemData\":{\"createdAt\":\"2017-05-24T13:28:53.4540398Z\"},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"},\"type\":\"Microsoft.DocumentDB/databaseAccounts\"}]}",
								},
							},
						},
						Redundancy: &voc.Redundancy{},
						ActivityLogging: &voc.ActivityLogging{
							Logging: &voc.Logging{
								Enabled:        true,
								LoggingService: []voc.ResourceID{"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/insights-integration/providers/Microsoft.OperationalInsights/workspaces/workspace1"},
							},
						},
					},
				},
				&voc.DatabaseService{
					StorageService: &voc.StorageService{
						NetworkService: &voc.NetworkService{
							Networking: &voc.Networking{
								Resource: &voc.Resource{
									ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.documentdb/databaseaccounts/cosmosdb2",
									Name:         "CosmosDB2",
									ServiceID:    testdata.MockCloudServiceID1,
									CreationTime: util.SafeTimestamp(&creationTime),
									Type:         voc.DatabaseServiceType,
									GeoLocation: voc.GeoLocation{
										Region: "eastus",
									},
									Labels: map[string]string{
										"testKey1": "testTag1",
										"testKey2": "testTag2",
									},
									Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
									Raw:    "{\"*armcosmos.DatabaseAccountGetResults\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/CosmosDB2\",\"kind\":\"MongoDB\",\"location\":\"eastus\",\"name\":\"CosmosDB2\",\"properties\":{},\"systemData\":{\"createdAt\":\"2017-05-24T13:28:53.4540398Z\"},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"},\"type\":\"Microsoft.DocumentDB/databaseAccounts\"}]}",
								},
							},
						},
						Redundancy: &voc.Redundancy{},
					},
				},
			},

			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureStorageDiscovery{
				azureDiscovery:     tt.fields.azureDiscovery,
				defenderProperties: tt.fields.defenderProperties,
			}

			// initialize diagnostic settings client
			if tt.fields.diagnosticSettingsClient {
				_ = d.initDiagnosticsSettingsClient()
			}

			got, err := d.discoverCosmosDB()
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_azureStorageDiscovery_handleCosmosDB(t *testing.T) {
	type fields struct {
		azureDiscovery     *azureDiscovery
		defenderProperties map[string]*defenderProperties
	}
	type args struct {
		account         *armcosmos.DatabaseAccountGetResults
		activityLogging *voc.ActivityLogging
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []voc.IsCloudResource
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path: ManagedKeyEncryption",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockStorageSender()),
			},
			args: args{
				account: &armcosmos.DatabaseAccountGetResults{
					Location: util.Ref("eastus"),
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/CosmosDB1"),
					Name:     util.Ref("CosmosDB1"),
					Tags: map[string]*string{
						"testKey1": util.Ref("testTag1"),
						"testKey2": util.Ref("testTag2"),
					},
					Properties: &armcosmos.DatabaseAccountGetProperties{},
					SystemData: &armcosmos.SystemData{
						CreatedAt: &time.Time{},
					},
				},
				activityLogging: &voc.ActivityLogging{
					Logging: &voc.Logging{
						Enabled: true,
					},
				},
			},
			want: []voc.IsCloudResource{
				&voc.DatabaseService{
					StorageService: &voc.StorageService{
						NetworkService: &voc.NetworkService{
							Networking: &voc.Networking{
								Resource: &voc.Resource{
									ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.documentdb/databaseaccounts/cosmosdb1",
									Name:         "CosmosDB1",
									ServiceID:    testdata.MockCloudServiceID1,
									CreationTime: util.SafeTimestamp(&time.Time{}),
									Type:         voc.DatabaseServiceType,
									GeoLocation: voc.GeoLocation{
										Region: "eastus",
									},
									Labels: map[string]string{
										"testKey1": "testTag1",
										"testKey2": "testTag2",
									},
									Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
									Raw:    "{\"*armcosmos.DatabaseAccountGetResults\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/CosmosDB1\",\"location\":\"eastus\",\"name\":\"CosmosDB1\",\"properties\":{},\"systemData\":{\"createdAt\":\"0001-01-01T00:00:00Z\"},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}]}",
								},
							},
						},
						Redundancy: &voc.Redundancy{},
						ActivityLogging: &voc.ActivityLogging{
							Logging: &voc.Logging{
								Enabled: true,
							},
						},
					},
				},
			},

			wantErr: assert.NoError,
		},
		{
			name: "Happy path: CustomerKeyEncryption",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockStorageSender()),
			},
			args: args{
				account: &armcosmos.DatabaseAccountGetResults{
					Location: util.Ref("eastus"),
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/CosmosDB2"),
					Name:     util.Ref("CosmosDB2"),
					Tags: map[string]*string{
						"testKey1": util.Ref("testTag1"),
						"testKey2": util.Ref("testTag2"),
					},
					Properties: &armcosmos.DatabaseAccountGetProperties{
						KeyVaultKeyURI: util.Ref("https://testvault.vault.azure.net/keys/testkey/123456"),
					},
					SystemData: &armcosmos.SystemData{
						CreatedAt: &time.Time{},
					},
				},
				activityLogging: &voc.ActivityLogging{
					Logging: &voc.Logging{
						Enabled: true,
					},
				},
			},
			want: []voc.IsCloudResource{
				&voc.DatabaseService{
					StorageService: &voc.StorageService{
						NetworkService: &voc.NetworkService{
							Networking: &voc.Networking{
								Resource: &voc.Resource{
									ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.documentdb/databaseaccounts/cosmosdb2",
									Name:         "CosmosDB2",
									ServiceID:    testdata.MockCloudServiceID1,
									CreationTime: util.SafeTimestamp(&time.Time{}),
									Type:         voc.DatabaseServiceType,
									GeoLocation: voc.GeoLocation{
										Region: "eastus",
									},
									Labels: map[string]string{
										"testKey1": "testTag1",
										"testKey2": "testTag2",
									},
									Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
									Raw:    "{\"*armcosmos.DatabaseAccountGetResults\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/CosmosDB2\",\"location\":\"eastus\",\"name\":\"CosmosDB2\",\"properties\":{\"keyVaultKeyUri\":\"https://testvault.vault.azure.net/keys/testkey/123456\"},\"systemData\":{\"createdAt\":\"0001-01-01T00:00:00Z\"},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}]}",
								},
							},
						},
						Redundancy: &voc.Redundancy{},
						ActivityLogging: &voc.ActivityLogging{
							Logging: &voc.Logging{
								Enabled: true,
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
			d := &azureStorageDiscovery{
				azureDiscovery:     tt.fields.azureDiscovery,
				defenderProperties: tt.fields.defenderProperties,
			}
			got, err := d.handleCosmosDB(tt.args.account, tt.args.activityLogging)
			if !tt.wantErr(t, err, fmt.Sprintf("handleCosmosDB(%v, %v)", tt.args.account, tt.args.account)) {
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_getCosmosDBRedundancy(t *testing.T) {
	type args struct {
		account *armcosmos.DatabaseAccountGetResults
	}
	tests := []struct {
		name  string
		args  args
		wantR *voc.Redundancy
	}{
		{
			name: "Happy path second location is zone redundant - return zone redundancy equals true ",
			args: args{account: &armcosmos.DatabaseAccountGetResults{
				Properties: &armcosmos.DatabaseAccountGetProperties{
					Locations: []*armcosmos.Location{
						{
							ID:              util.Ref("location1ID"),
							IsZoneRedundant: util.Ref(false),
						},
						{
							ID:              util.Ref("location2ID"),
							IsZoneRedundant: util.Ref(true),
						},
					},
				},
			}},
			wantR: &voc.Redundancy{Zone: true},
		},
		{
			name: "No location is zone redundant - return zone redundancy equals false",
			args: args{account: &armcosmos.DatabaseAccountGetResults{
				Properties: &armcosmos.DatabaseAccountGetProperties{
					Locations: []*armcosmos.Location{
						{
							ID:              util.Ref("location1ID"),
							IsZoneRedundant: util.Ref(false),
						},
						{
							ID:              util.Ref("location2ID"),
							IsZoneRedundant: util.Ref(false),
						},
					},
				},
			}},
			wantR: &voc.Redundancy{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantR, getCosmosDBRedundancy(tt.args.account), "getCosmosDBRedundancy(%v)", tt.args.account)
		})
	}
}
