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
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"k8s.io/apimachinery/pkg/util/json"

	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/stretchr/testify/assert"
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
		return createResponse(map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1",
					"name":     "account1",
					"location": "eastus",
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
		return createResponse(map[string]interface{}{
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
		return createResponse(map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container1",
					"name": "container1",
					"type": "Microsoft.Storage/storageAccounts/blobServices/containers",
				},
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container2",
					"name": "container2",
					"type": "Microsoft.Storage/storageAccounts/blobServices/containers",
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2/blobServices/default/containers" {
		return createResponse(map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2/blobServices/default/containers/container3",
					"name": "container3",
					"type": "Microsoft.Storage/storageAccounts/blobServices/containers",
				},
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2/blobServices/default/containers/container4",
					"name": "container4",
					"type": "Microsoft.Storage/storageAccounts/blobServices/containers",
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares" {
		return createResponse(map[string]interface{}{
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
		return createResponse(map[string]interface{}{
			"value": &[]map[string]interface{}{},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Compute/disks" {
		return createResponse(map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk1",
					"name":     "disk1",
					"type":     "Microsoft.Compute/disks",
					"location": "eastus",
					"properties": map[string]interface{}{
						"timeCreated": "2017-05-24T13:28:53.4540398Z",
						"encryption": map[string]interface{}{
							"diskEncryptionSetID": "",
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
							"diskEncryptionSetID": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault1",
							"type":                "EncryptionAtRestWithCustomerKey",
						},
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault1" {
		return createResponse(map[string]interface{}{
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
		return createResponse(map[string]interface{}{
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
	}

	return m.mockSender.Do(req)
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
	assert.Equal(t, 10, len(list))
	assert.NotEmpty(t, d.Name())
}

func Test_azureStorageDiscovery_List(t *testing.T) {
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)

	type fields struct {
		azureDiscovery azureDiscovery
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
				azureDiscovery: azureDiscovery{
					cred: nil,
				},
			},
			wantList: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrCouldNotAuthenticate.Error())
			},
		},
		{
			name: "Discovery error",
			fields: fields{
				azureDiscovery: azureDiscovery{
					cred: &mockAuthorizer{},
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: &mockNetworkSender{},
						},
					},
				},
			},
			wantList: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not discover storage accounts:")
			},
		},
		{
			name: "Without errors",
			fields: fields{
				azureDiscovery: azureDiscovery{
					cred: &mockAuthorizer{},
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: &mockStorageSender{},
						},
					},
				},
			},
			wantList: []voc.IsCloudResource{
				&voc.ObjectStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container1",
							Name:         "container1",
							Type:         []string{"ObjectStorage", "Storage", "Resource"},
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						AtRestEncryption: voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
					},
				},
				&voc.ObjectStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container2",
							Name:         "container2",
							Type:         []string{"ObjectStorage", "Storage", "Resource"},
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						AtRestEncryption: voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
					},
				},
				&voc.FileStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/fileshare1",
							Name:         "fileshare1",
							Type:         []string{"FileStorage", "Storage", "Resource"},
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						AtRestEncryption: voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
					},
				},
				&voc.FileStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/fileshare2",
							Name:         "fileshare2",
							Type:         []string{"FileStorage", "Storage", "Resource"},
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						AtRestEncryption: voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
					},
				},
				&voc.StorageService{
					Storages: []voc.ResourceID{
						"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container1",
						"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container2",
						"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/fileshare1",
						"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/fileshare2",
					},
					HttpEndpoint: &voc.HttpEndpoint{
						Url: "https://account1.[file,blob].core.windows.net/",
						TransportEncryption: &voc.TransportEncryption{
							Enforced:   true,
							Enabled:    true,
							TlsVersion: "TLS1_2",
							Algorithm:  "TLS",
						},
					},
					NetworkService: &voc.NetworkService{
						Networking: &voc.Networking{
							Resource: &voc.Resource{
								ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1",
								Name:         "account1",
								Type:         []string{"StorageService", "NetworkService", "Networking", "Resource"},
								CreationTime: util.SafeTimestamp(&creationTime),
								Labels:       map[string]string{},
								GeoLocation: voc.GeoLocation{
									Region: "eastus",
								},
							},
						},
						TransportEncryption: &voc.TransportEncryption{
							Enforced:   true,
							Enabled:    true,
							TlsVersion: "TLS1_2",
							Algorithm:  "TLS",
						},
					},
				},
				&voc.ObjectStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2/blobServices/default/containers/container3",
							Name:         "container3",
							Type:         []string{"ObjectStorage", "Storage", "Resource"},
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						AtRestEncryption: voc.CustomerKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "",
								Enabled:   true,
							},
							KeyUrl: "https://testvault.vault.azure.net/keys/testkey/123456",
						},
					},
				},
				&voc.ObjectStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2/blobServices/default/containers/container4",
							Name:         "container4",
							Type:         []string{"ObjectStorage", "Storage", "Resource"},
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						AtRestEncryption: voc.CustomerKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "",
								Enabled:   true,
							},
							KeyUrl: "https://testvault.vault.azure.net/keys/testkey/123456",
						},
					},
				},
				&voc.StorageService{
					Storages: []voc.ResourceID{
						"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2/blobServices/default/containers/container3",
						"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2/blobServices/default/containers/container4",
					},
					HttpEndpoint: &voc.HttpEndpoint{
						Url: "https://account1.[file,blob].core.windows.net/",
						TransportEncryption: &voc.TransportEncryption{
							Enforced:   true,
							Enabled:    true,
							TlsVersion: "TLS1_2",
							Algorithm:  "TLS",
						},
					},
					NetworkService: &voc.NetworkService{
						Networking: &voc.Networking{
							Resource: &voc.Resource{
								ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2",
								Name:         "account2",
								Type:         []string{"StorageService", "NetworkService", "Networking", "Resource"},
								CreationTime: util.SafeTimestamp(&creationTime),
								Labels:       map[string]string{},
								GeoLocation: voc.GeoLocation{
									Region: "eastus",
								},
							},
						},
						TransportEncryption: &voc.TransportEncryption{
							Enforced:   true,
							Enabled:    true,
							TlsVersion: "TLS1_2",
							Algorithm:  "TLS",
						},
					},
				},
				&voc.BlockStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk1",
							Name:         "disk1",
							CreationTime: util.SafeTimestamp(&creationTime),
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Labels: map[string]string{},
							Type:   []string{"BlockStorage", "Storage", "Resource"},
						},

						AtRestEncryption: voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
					},
				},
				&voc.BlockStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk2",
							Name:         "disk2",
							CreationTime: util.SafeTimestamp(&creationTime),
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Labels: map[string]string{},
							Type:   []string{"BlockStorage", "Storage", "Resource"},
						},
						AtRestEncryption: voc.CustomerKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "",
								Enabled:   true,
							},
							KeyUrl: "https://keyvault1.vault.azure.net/keys/customer-key/6273gdb374jz789hjm17819283748382",
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
			gotList, err := d.List()
			if !tt.wantErr(t, err) {
				return
			}

			assert.Equal(t, tt.wantList, gotList)
		})
	}
}

func TestStorageHandleMethodsWhenInputIsInvalid(t *testing.T) {
	d := azureStorageDiscovery{}

	// Get mocked armstorage.Account
	reqURL := "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Storage/storageAccounts/account3"
	mockedStorageAccountObject, err := mockedStorageAccount(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	// Test method handleObjectStorage
	containerItem := armstorage.ListContainerItem{}
	handleObjectStorageRespone, err := handleObjectStorage(mockedStorageAccountObject, &containerItem)
	assert.Error(t, err)
	assert.Nil(t, handleObjectStorageRespone)

	// Test method handleFileStorage
	fileShare := &armstorage.FileShareItem{}
	handleFileStorageRespone, err := handleFileStorage(mockedStorageAccountObject, fileShare)
	assert.Error(t, err)
	assert.Nil(t, handleFileStorageRespone)

	// Test method handleBlockStorage
	disk := &armcompute.Disk{}
	handleBlockStorageResponse, err := d.handleBlockStorage(disk)
	assert.Error(t, err)
	assert.Nil(t, handleBlockStorageResponse)
}

func TestStorageMethodsWhenInputIsInvalid(t *testing.T) {
	// Get mocked armstorage.Account
	reqURL := "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Storage/storageAccounts/account3"
	mockedStorageAccountObject, err := mockedStorageAccount(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	// Test method diskEncryptionSetName
	discEncryptionSetID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault1"
	assert.Equal(t, "encryptionkeyvault1", diskEncryptionSetName(discEncryptionSetID))

	// Test method storageAtRestEncryption
	atRestEncryption, err := storageAtRestEncryption(mockedStorageAccountObject)
	assert.NoError(t, err)

	managedKeyEncryption := voc.ManagedKeyEncryption{AtRestEncryption: &voc.AtRestEncryption{Algorithm: "AES256", Enabled: true}}
	assert.Equal(t, managedKeyEncryption, atRestEncryption)

}

func TestStorageDiscoverMethodsWhenInputIsInvalid(t *testing.T) {
	d := azureStorageDiscovery{}

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
	discoverObjectStoragesResponse, err := d.discoverObjectStorages(mockedStorageAccountObject)
	assert.ErrorContains(t, err, ErrGettingNextPage.Error())
	assert.Nil(t, discoverObjectStoragesResponse)

	// Test method discoverFileStorages
	discoverFileStoragesResponse, err := d.discoverFileStorages(mockedStorageAccountObject)
	assert.ErrorContains(t, err, ErrGettingNextPage.Error())
	assert.Nil(t, discoverFileStoragesResponse)

	// Test method discoverBlockStorages
	discoverBlockStoragesResponse, err := d.discoverBlockStorages()
	assert.ErrorContains(t, err, ErrGettingNextPage.Error())
	assert.Nil(t, discoverBlockStoragesResponse)
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
	var subID = "00000000-0000-0000-0000-000000000000"
	sub := armsubscription.Subscription{
		SubscriptionID: &subID,
	}

	type fields struct {
		azureDiscovery azureDiscovery
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
				azureDiscovery: azureDiscovery{
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
				azureDiscovery: azureDiscovery{
					cred: &mockAuthorizer{},
					sub:  sub,
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockStorageSender{},
						},
					},
				},
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
				assert.Equal(t, 10, len(got))
			}
		})
	}
}

func Test_azureStorageDiscovery_handleBlockStorage(t *testing.T) {
	var subID = "00000000-0000-0000-0000-000000000000"
	sub := &armsubscription.Subscription{
		SubscriptionID: &subID,
	}

	encType := armcompute.EncryptionTypeEncryptionAtRestWithCustomerKey
	diskID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk1"
	diskName := "disk1"
	diskRegion := "eastus"
	encSetID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault1"
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)

	type fields struct {
		azureDiscovery azureDiscovery
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
				azureDiscovery: azureDiscovery{
					cred: &mockAuthorizer{},
					sub:  *sub,
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockStorageSender{},
						},
					},
				},
			},
			want: &voc.BlockStorage{
				Storage: &voc.Storage{
					Resource: &voc.Resource{
						ID:           voc.ResourceID(diskID),
						Name:         "disk1",
						CreationTime: util.SafeTimestamp(&creationTime),
						Type:         []string{"BlockStorage", "Storage", "Resource"},
						GeoLocation: voc.GeoLocation{
							Region: "eastus",
						},
						Labels: map[string]string{},
					},

					AtRestEncryption: voc.CustomerKeyEncryption{
						AtRestEncryption: &voc.AtRestEncryption{
							Algorithm: "",
							Enabled:   true,
						},
						KeyUrl: "https://keyvault1.vault.azure.net/keys/customer-key/6273gdb374jz789hjm17819283748382",
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
			got, err := d.handleBlockStorage(tt.args.disk)
			if !tt.wantErr(t, err, fmt.Sprintf("handleBlockStorage(%v)", tt.args.disk)) {
				return
			}
			assert.Equal(t, tt.want, got)
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
		want    voc.HasAtRestEncryption
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
			want: voc.ManagedKeyEncryption{
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

func Test_azureStorageDiscovery_keyURL(t *testing.T) {
	encSetID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault1"
	encSetID2 := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault2"
	var subID = "00000000-0000-0000-0000-000000000000"
	sub := armsubscription.Subscription{
		SubscriptionID: &subID,
	}

	type fields struct {
		azureDiscovery azureDiscovery
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
				return assert.ErrorContains(t, err, "empty diskEncryptionSetID")
			},
		},
		{
			name: "Error get disc encryption set",
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
			args: args{
				diskEncryptionSetID: encSetID2,
			},
			fields: fields{
				azureDiscovery: azureDiscovery{
					cred: &mockAuthorizer{},
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockStorageSender{},
						},
					},
					sub: sub,
				},
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
				azureDiscovery: azureDiscovery{
					cred: &mockAuthorizer{},
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockStorageSender{},
						},
					},
					sub: sub,
				},
			},
			want:    "https://keyvault1.vault.azure.net/keys/customer-key/6273gdb374jz789hjm17819283748382",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureStorageDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}
			got, err := d.keyURL(tt.args.diskEncryptionSetID)
			if !tt.wantErr(t, err, fmt.Sprintf("keyURL(%v)", tt.args.diskEncryptionSetID)) {
				return
			}
			assert.Equalf(t, tt.want, got, "keyURL(%v)", tt.args.diskEncryptionSetID)
		})
	}
}

func Test_azureStorageDiscovery_blockStorageAtRestEncryption(t *testing.T) {
	var subID = "00000000-0000-0000-0000-000000000000"
	sub := armsubscription.Subscription{
		SubscriptionID: &subID,
	}

	encType := armcompute.EncryptionTypeEncryptionAtRestWithCustomerKey
	diskID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk1"
	diskName := "disk1"
	diskRegion := "eastus"
	encSetID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault1"
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)

	type fields struct {
		azureDiscovery azureDiscovery
	}
	type args struct {
		disk *armcompute.Disk
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    voc.HasAtRestEncryption
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
				azureDiscovery: azureDiscovery{
					cred: &mockAuthorizer{},
					sub:  sub,
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockStorageSender{},
						},
					},
				},
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
				azureDiscovery: azureDiscovery{
					cred: &mockAuthorizer{},
					sub:  sub,
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockStorageSender{},
						},
					},
				},
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
			want: voc.CustomerKeyEncryption{
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
			d := &azureStorageDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}
			got, err := d.blockStorageAtRestEncryption(tt.args.disk)
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

	type args struct {
		account   *armstorage.Account
		fileshare *armstorage.FileShareItem
	}
	tests := []struct {
		name    string
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
			args: args{
				account: &armstorage.Account{
					Properties: &armstorage.AccountProperties{
						Encryption: &armstorage.Encryption{
							KeySource: &keySource,
						},
						CreationTime: &creationTime,
					},
					Location: &accountRegion,
				},
				fileshare: &armstorage.FileShareItem{
					ID:   &fileShareID,
					Name: &fileShareName,
				},
			},
			want: &voc.FileStorage{
				Storage: &voc.Storage{
					Resource: &voc.Resource{
						ID:           voc.ResourceID(fileShareID),
						Name:         fileShareName,
						CreationTime: util.SafeTimestamp(&creationTime),
						GeoLocation: voc.GeoLocation{
							Region: accountRegion,
						},
						Labels: map[string]string{},
						Type:   []string{"FileStorage", "Storage", "Resource"},
					},
					AtRestEncryption: voc.ManagedKeyEncryption{
						AtRestEncryption: &voc.AtRestEncryption{
							Algorithm: "AES256",
							Enabled:   true,
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handleFileStorage(tt.args.account, tt.args.fileshare)
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
		azureDiscovery azureDiscovery
	}
	type args struct {
		account      *armstorage.Account
		storagesList []voc.IsCloudResource
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *voc.StorageService
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
			args: args{
				account: &armstorage.Account{
					ID:   &accountID,
					Name: &accountName,
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
			},
			want: &voc.StorageService{
				NetworkService: &voc.NetworkService{
					Networking: &voc.Networking{
						Resource: &voc.Resource{
							ID:           voc.ResourceID(accountID),
							Name:         accountName,
							CreationTime: util.SafeTimestamp(&creationTime),
							Type:         []string{"StorageService", "NetworkService", "Networking", "Resource"},
							GeoLocation: voc.GeoLocation{
								Region: accountRegion,
							},
							Labels: map[string]string{},
						},
					},
					TransportEncryption: &voc.TransportEncryption{
						Enforced:   true,
						Enabled:    true,
						TlsVersion: "TLS1_2",
						Algorithm:  "TLS",
					},
				},
				HttpEndpoint: &voc.HttpEndpoint{
					Url: "https://account1.[file,blob].core.windows.net",
					TransportEncryption: &voc.TransportEncryption{
						Enforced:   true,
						Enabled:    true,
						TlsVersion: "TLS1_2",
						Algorithm:  "TLS",
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
			got, err := az.handleStorageAccount(tt.args.account, tt.args.storagesList)
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

	type args struct {
		account   *armstorage.Account
		container *armstorage.ListContainerItem
	}
	tests := []struct {
		name    string
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
			args: args{
				account: &armstorage.Account{
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
				},
			},
			want: &voc.ObjectStorage{
				Storage: &voc.Storage{
					Resource: &voc.Resource{
						ID:           voc.ResourceID(containerID),
						Name:         containerName,
						CreationTime: util.SafeTimestamp(&creationTime),
						GeoLocation: voc.GeoLocation{
							Region: accountRegion,
						},
						Labels: map[string]string{},
						Type:   []string{"ObjectStorage", "Storage", "Resource"},
					},
					AtRestEncryption: voc.ManagedKeyEncryption{
						AtRestEncryption: &voc.AtRestEncryption{
							Algorithm: "AES256",
							Enabled:   true,
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handleObjectStorage(tt.args.account, tt.args.container)
			if !tt.wantErr(t, err, fmt.Sprintf("handleObjectStorage(%v, %v)", tt.args.account, tt.args.container)) {
				return
			}
			assert.Equalf(t, tt.want, got, "handleObjectStorage(%v, %v)", tt.args.account, tt.args.container)
		})
	}
}

func Test_azureStorageDiscovery_discoverBlockStorages(t *testing.T) {
	var subID = "00000000-0000-0000-0000-000000000000"
	sub := armsubscription.Subscription{
		SubscriptionID: &subID,
	}
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)

	type fields struct {
		azureDiscovery azureDiscovery
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
				azureDiscovery: azureDiscovery{
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
				azureDiscovery: azureDiscovery{
					cred: &mockAuthorizer{},
					sub:  sub,
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockStorageSender{},
						},
					},
				},
			},
			want: []voc.IsCloudResource{
				&voc.BlockStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk1",
							Name:         "disk1",
							CreationTime: util.SafeTimestamp(&creationTime),
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Type:   []string{"BlockStorage", "Storage", "Resource"},
							Labels: map[string]string{},
						},

						AtRestEncryption: voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
					},
				},
				&voc.BlockStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk2",
							Name:         "disk2",
							CreationTime: util.SafeTimestamp(&creationTime),
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Type:   []string{"BlockStorage", "Storage", "Resource"},
							Labels: map[string]string{},
						},
						AtRestEncryption: voc.CustomerKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "",
								Enabled:   true,
							},
							KeyUrl: "https://keyvault1.vault.azure.net/keys/customer-key/6273gdb374jz789hjm17819283748382",
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
			got, err := d.discoverBlockStorages()
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equalf(t, tt.want, got, "discoverBlockStorages()")
		})
	}
}

func Test_azureStorageDiscovery_discoverFileStorages(t *testing.T) {
	var subID = "00000000-0000-0000-0000-000000000000"
	sub := armsubscription.Subscription{
		SubscriptionID: &subID,
	}
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)
	accountID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1"
	accountName := "account1"
	accountRegion := "eastus"
	keySource := armstorage.KeySourceMicrosoftStorage

	type fields struct {
		azureDiscovery azureDiscovery
	}
	type args struct {
		account *armstorage.Account
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
				azureDiscovery: azureDiscovery{
					cred: nil,
					sub:  sub,
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
				azureDiscovery: azureDiscovery{
					cred: &mockAuthorizer{},
					sub:  sub,
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockStorageSender{},
						},
					},
				},
			},
			args: args{
				account: &armstorage.Account{
					ID:   &accountID,
					Name: &accountName,
					Properties: &armstorage.AccountProperties{
						Encryption: &armstorage.Encryption{
							KeySource: &keySource,
						},
						CreationTime: &creationTime,
					},
					Location: &accountRegion,
				},
			},
			want: []voc.IsCloudResource{
				&voc.FileStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/fileshare1",
							Name:         "fileshare1",
							Type:         []string{"FileStorage", "Storage", "Resource"},
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						AtRestEncryption: voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
					},
				},
				&voc.FileStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/fileshare2",
							Name:         "fileshare2",
							Type:         []string{"FileStorage", "Storage", "Resource"},
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						AtRestEncryption: voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
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
				azureDiscovery: tt.fields.azureDiscovery,
			}
			got, err := d.discoverFileStorages(tt.args.account)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_azureStorageDiscovery_discoverObjectStorages(t *testing.T) {
	var subID = "00000000-0000-0000-0000-000000000000"
	sub := armsubscription.Subscription{
		SubscriptionID: &subID,
	}
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)
	accountID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1"
	accountName := "account1"
	accountRegion := "eastus"
	keySource := armstorage.KeySourceMicrosoftStorage

	type fields struct {
		azureDiscovery azureDiscovery
	}
	type args struct {
		account *armstorage.Account
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
				azureDiscovery: azureDiscovery{
					cred: nil,
					sub:  sub,
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
				azureDiscovery: azureDiscovery{
					cred: &mockAuthorizer{},
					sub:  sub,
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockStorageSender{},
						},
					},
				},
			},
			args: args{
				account: &armstorage.Account{
					ID:   &accountID,
					Name: &accountName,
					Properties: &armstorage.AccountProperties{
						Encryption: &armstorage.Encryption{
							KeySource: &keySource,
						},
						CreationTime: &creationTime,
					},
					Location: &accountRegion,
				},
			},
			want: []voc.IsCloudResource{
				&voc.ObjectStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container1",
							Name:         "container1",
							Type:         []string{"ObjectStorage", "Storage", "Resource"},
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						AtRestEncryption: voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
					},
				},
				&voc.ObjectStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container2",
							Name:         "container2",
							Type:         []string{"ObjectStorage", "Storage", "Resource"},
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						AtRestEncryption: voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
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
				azureDiscovery: tt.fields.azureDiscovery,
			}
			got, err := d.discoverObjectStorages(tt.args.account)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
