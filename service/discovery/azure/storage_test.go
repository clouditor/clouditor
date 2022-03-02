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
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"testing"

	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-07-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-02-01/storage"
	"k8s.io/apimachinery/pkg/util/json"

	"github.com/stretchr/testify/assert"
)

func init() {
	log = logrus.WithField("component", "azure-tests")
	log.Logger.Formatter = formatter.CapitalizeFormatter{Formatter: &logrus.TextFormatter{ForceColors: true, FullTimestamp: true}}
}

type mockStorageSender struct {
	mockSender
}

func newMockStorageSender() *mockStorageSender {
	m := &mockStorageSender{}
	return m
}

type responseStorageAccount struct {
	Value storage.Account `json:"value,omitempty"`
}

type responseDisk struct {
	Value []compute.Disk `json:"value,omitempty"`
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
							"keySource": storage.KeySourceMicrosoftStorage,
						},
						"minimumTlsVersion":        storage.MinimumTLSVersionTLS12,
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
							"keySource": storage.KeySourceMicrosoftKeyvault,
							"keyvaultproperties": map[string]interface{}{
								"keyvaulturi": "https://testvault.vault.azure.net/keys/testkey/123456",
							},
						},
						"minimumTlsVersion":        storage.MinimumTLSVersionTLS12,
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
						"keySource": storage.KeySourceMicrosoftStorage,
					},
					"minimumTlsVersion":        storage.MinimumTLSVersionTLS12,
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
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault1" { // ?api-version=2021-07-01?
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
	}

	return m.mockSender.Do(req)
}

func TestAzureStorageAuthorizer(t *testing.T) {

	d := NewAzureStorageDiscovery()
	list, err := d.List()

	assert.Error(t, err)
	assert.Nil(t, list)
	assert.Equal(t, "could not authorize Azure account: no authorized was available", err.Error())
}

func TestStorage(t *testing.T) {
	d := NewAzureStorageDiscovery(
		WithSender(&mockStorageSender{}),
		WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 8, len(list))
	assert.NotEmpty(t, d.Name())
}

func TestListObjectStorage(t *testing.T) {
	d := NewAzureStorageDiscovery(
		WithSender(&mockStorageSender{}),
		WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()
	assert.NoError(t, err)
	assert.NotNil(t, list)

	var counter int
	for _, item := range list {
		if item.GetType()[0] == "ObjectStorage" {
			counter++
		}
	}
	assert.Equal(t, 4, counter)
}

func TestObjectStorage(t *testing.T) {
	d := NewAzureStorageDiscovery(
		WithSender(&mockStorageSender{}),
		WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()
	assert.NoError(t, err)
	assert.NotNil(t, list)

	objectStorage, ok := list[0].(*voc.ObjectStorage)

	assert.True(t, ok)
	assert.Equal(t, "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container1", string(objectStorage.ID))
	assert.Equal(t, "container1", objectStorage.Name)
	assert.NotNil(t, objectStorage.CreationTime)
	assert.Equal(t, "ObjectStorage", objectStorage.Type[0])
	assert.Equal(t, "eastus", objectStorage.GeoLocation.Region)
	assert.Equal(t, true, objectStorage.AtRestEncryption.GetAtRestEncryption().Enabled)
	assert.Equal(t, "https://account1.blob.core.windows.net/container1", objectStorage.HttpEndpoint.Url)
	assert.Equal(t, true, objectStorage.HttpEndpoint.TransportEncryption.Enabled)
	assert.Equal(t, true, objectStorage.HttpEndpoint.TransportEncryption.Enforced)
	assert.Equal(t, "TLS1_2", objectStorage.HttpEndpoint.TransportEncryption.TlsVersion)
	assert.Equal(t, "TLS", objectStorage.HttpEndpoint.TransportEncryption.Algorithm)

	// Check ManagedKeyEncryption
	atRestEncryption := *objectStorage.GetAtRestEncryption()
	managedKeyEncryption, ok := atRestEncryption.(voc.ManagedKeyEncryption)
	assert.True(t, ok)
	assert.Equal(t, true, managedKeyEncryption.Enabled)
	assert.Equal(t, "AES256", managedKeyEncryption.Algorithm)

	// Check CustomerKeyEncryption
	objectStorage, ok = list[4].(*voc.ObjectStorage)
	assert.True(t, ok)
	atRestEncryption = *objectStorage.GetAtRestEncryption()
	customerKeyEncryption, ok := atRestEncryption.(voc.CustomerKeyEncryption)
	assert.True(t, ok)
	assert.Equal(t, true, customerKeyEncryption.Enabled)
	assert.Equal(t, "https://testvault.vault.azure.net/keys/testkey/123456", customerKeyEncryption.KeyUrl)
}

func TestListFileStorage(t *testing.T) {
	d := NewAzureStorageDiscovery(
		WithSender(&mockStorageSender{}),
		WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()
	assert.NoError(t, err)
	assert.NotNil(t, list)

	var counter int
	for _, item := range list {
		if item.GetType()[0] == "FileStorage" {
			counter++
		}
	}
	assert.Equal(t, 2, counter)
}

func TestFileStorage(t *testing.T) {
	d := NewAzureStorageDiscovery(
		WithSender(&mockStorageSender{}),
		WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()
	assert.NoError(t, err)
	assert.NotNil(t, list)

	fileStorage, ok := list[2].(*voc.FileStorage)

	assert.True(t, ok)
	assert.Equal(t, "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/fileshare1", string(fileStorage.ID))
	assert.Equal(t, "fileshare1", fileStorage.Name)
	assert.NotNil(t, fileStorage.CreationTime)
	assert.Equal(t, "FileStorage", fileStorage.Type[0])
	assert.Equal(t, "eastus", fileStorage.GeoLocation.Region)
	assert.Equal(t, true, fileStorage.AtRestEncryption.GetAtRestEncryption().Enabled)
	assert.Equal(t, "https://account1.file.core.windows.net/fileshare1", fileStorage.HttpEndpoint.Url)
	assert.Equal(t, true, fileStorage.HttpEndpoint.TransportEncryption.Enabled)
	assert.Equal(t, true, fileStorage.HttpEndpoint.TransportEncryption.Enforced)
	assert.Equal(t, "TLS1_2", fileStorage.HttpEndpoint.TransportEncryption.TlsVersion)
	assert.Equal(t, "TLS", fileStorage.HttpEndpoint.TransportEncryption.Algorithm)
}

func TestListBlockStorage(t *testing.T) {
	d := NewAzureStorageDiscovery(
		WithSender(&mockStorageSender{}),
		WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()
	assert.NoError(t, err)
	assert.NotNil(t, list)

	var counter int
	for _, item := range list {
		if item.GetType()[0] == "BlockStorage" {
			counter++
		}
	}
	assert.Equal(t, 2, counter)
}

func TestBlockStorage(t *testing.T) {
	d := NewAzureStorageDiscovery(
		WithSender(&mockStorageSender{}),
		WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()
	assert.NoError(t, err)
	assert.NotNil(t, list)

	blockStorage, ok := list[6].(*voc.BlockStorage)

	assert.True(t, ok)
	assert.Equal(t, "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk1", string(blockStorage.ID))
	assert.Equal(t, "disk1", blockStorage.Name)
	assert.NotNil(t, blockStorage.CreationTime)
	assert.Equal(t, "BlockStorage", blockStorage.Type[0])
	assert.Equal(t, "eastus", blockStorage.GeoLocation.Region)
	assert.Equal(t, true, blockStorage.AtRestEncryption.GetAtRestEncryption().Enabled)
}

func TestStorageHandleMethodsWhenInputIsInvalid(t *testing.T) {
	d := azureStorageDiscovery{}

	// Get mocked storage.Account
	reqURL := "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Storage/storageAccounts/account3"
	mockedStorageAccountObject, err := mockedStorageAccount(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}
	// Test method handleObjectStorage
	containerItem := storage.ListContainerItem{}

	// Clear KeySource
	mockedStorageAccountObject.Encryption.KeySource = ""

	handleObjectStorageRespone, err := handleObjectStorage(&mockedStorageAccountObject, containerItem)
	assert.Error(t, err)
	assert.Nil(t, handleObjectStorageRespone)

	// Test method handleFileStorage
	fileShare := storage.FileShareItem{}

	// Clear KeySource
	mockedStorageAccountObject.Encryption.KeySource = ""

	handleFileStorageRespone, err := handleFileStorage(&mockedStorageAccountObject, fileShare)
	assert.Error(t, err)
	assert.Nil(t, handleFileStorageRespone)

	// Test method handleBlockStorage
	reqURL = "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Compute/disks"
	disk, err := mockedDisk(reqURL)
	if err != nil {
		fmt.Println("error getting mocked disk object: %w", err)
	}
	// Clear KeySource
	disk.Encryption.Type = ""

	handleBlockStorageResponse, err := d.handleBlockStorage(disk)
	assert.Error(t, err)
	assert.Nil(t, handleBlockStorageResponse)
}

func TestStorageMethodsWhenInputIsInvalid(t *testing.T) {
	// Get mocked storage.Account
	reqURL := "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Storage/storageAccounts/account3"
	mockedStorageAccountObject, err := mockedStorageAccount(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	// Test method diskEncryptionSetName
	discEncryptionSetID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault1"
	assert.Equal(t, "encryptionkeyvault1", diskEncryptionSetName(discEncryptionSetID))

	// Test method storageAtRestEncryption
	atRestEncryption, err := storageAtRestEncryption(&mockedStorageAccountObject)
	assert.NoError(t, err)

	managedKeyEncryption := voc.ManagedKeyEncryption{AtRestEncryption: &voc.AtRestEncryption{Algorithm: "AES256", Enabled: true}}
	assert.Equal(t, managedKeyEncryption, atRestEncryption)

	// Test method blockStorageAtRestEncryption
	// Todo(garuppel): How to test? Problem: Azure call again
}

func TestStorageDiscoverMethodsWhenInputIsInvalid(t *testing.T) {
	d := azureStorageDiscovery{}

	// Get mocked storage.Account
	reqURL := "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Storage/storageAccounts/account3"
	mockedStorageAccountObject, err := mockedStorageAccount(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}
	// Test method discoverStorageAccounts
	discoverStorageAccountsResponse, err := d.discoverStorageAccounts()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not list storage accounts")
	assert.Nil(t, discoverStorageAccountsResponse)

	// Test method discoverObjectStorages
	discoverObjectStoragesResponse, err := d.discoverObjectStorages(&mockedStorageAccountObject)
	assert.Error(t, err)
	assert.Nil(t, discoverObjectStoragesResponse)

	// Test method discoverFileStorages
	discoverFileStoragesResponse, err := d.discoverFileStorages(&mockedStorageAccountObject)
	assert.Error(t, err)
	assert.Nil(t, discoverFileStoragesResponse)

	// Test method discoverBlockStorages
	discoverBlockStoragesResponse, err := d.discoverBlockStorages()
	assert.Error(t, err)
	assert.Nil(t, discoverBlockStoragesResponse)
}

// mockedDisk returns one mocked compute disk
func mockedDisk(reqUrl string) (disk compute.Disk, err error) {

	m := newMockStorageSender()
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return disk, errors.New("error creating new request")
	}
	resp, err := m.Do(req)
	if err != nil {
		return disk, fmt.Errorf("error getting mock http response: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("error io.ReadCloser: %w", err)
		}
	}(resp.Body)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return disk, fmt.Errorf("error read all: %w", err)
	}
	var disks responseDisk
	err = json.Unmarshal(responseBody, &disks)
	if err != nil {
		return disk, fmt.Errorf("error unmarshalling: %w", err)
	}

	return disks.Value[0], nil
}

// mockedStorageAccount returns one mocked storage account
func mockedStorageAccount(reqUrl string) (storageAccount storage.Account, err error) {
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

	storageAccount = storageAccountResponse.Value

	return storageAccount, nil
}
