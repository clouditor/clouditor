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
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-04-01/storage"
	"net/http"
	"testing"

	"clouditor.io/clouditor/service/discovery/azure"
	"github.com/stretchr/testify/assert"
)

type mockStorageSender struct {
	mockSender
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
						"minimumTlsVersion": storage.MinimumTLSVersionTLS12,
						"allowBlobPublicAccess": false,
						"supportsHttpsTrafficOnly": true,
					},
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
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares" {
		return createResponse(map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/shares/fileshare1",
					"name": "fileshare1",
					"type": "Microsoft.Storage/storageAccounts/fileServices/shares",
				},
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/fileshare2",
					"name": "fileshare2",
					"type": "Microsoft.Storage/storageAccounts/fileServices/shares",
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Compute/disks" {
		return createResponse(map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk1",
					"name": "disk1",
					"type": "Microsoft.Compute/disks",
					"location": "eastus",
					"properties": map[string]interface{}{
						"timeCreated": "2017-05-24T13:28:53.4540398Z",
						"encryption": map[string]interface{}{
							"diskEncryptionSetID": nil,
							"type":                "EncryptionAtRestWithPlatformKey",
						},
					},
				},
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk2",
					"name": "disk2",
					"type": "Microsoft.Compute/disks",
					"location": "eastus",
					"properties": map[string]interface{}{
						"timeCreated": "2017-05-24T13:28:53.4540398Z",
						"encryption": map[string]interface{}{
							"diskEncryptionSetID": nil,
							"type":                "EncryptionAtRestWithPlatformKey",
						},
					},
				},
			},
		}, 200)
	}

	return m.mockSender.Do(req)
}


// TODO tests for encryption stuff and extend basic tests
func TestListStorage(t *testing.T) {
	d := azure.NewAzureStorageDiscovery(
		azure.WithSender(&mockStorageSender{}),
		azure.WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 6, len(list))
}

func TestListObjectStorage(t *testing.T) {
	d := azure.NewAzureStorageDiscovery(
		azure.WithSender(&mockStorageSender{}),
		azure.WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()
	assert.Nil(t, err)
	assert.NotNil(t, list)

	var counter int
	for _, item := range list {
		if item.GetType()[0] == "ObjectStorage" {
			counter++
		}
	}
	assert.Equal(t, 2, counter)
}

func TestGetObjectStorage(t *testing.T) {
	d := azure.NewAzureStorageDiscovery(
		azure.WithSender(&mockStorageSender{}),
		azure.WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()
	assert.Nil(t, err)
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
	assert.Equal(t, "", objectStorage.HttpEndpoint.TransportEncryption.Algorithm)
}

func TestListFileStorage(t *testing.T) {
	d := azure.NewAzureStorageDiscovery(
		azure.WithSender(&mockStorageSender{}),
		azure.WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()
	assert.Nil(t, err)
	assert.NotNil(t, list)

	var counter int
	for _, item := range list {
		if item.GetType()[0] == "FileStorage" {
			counter++
		}
	}
	assert.Equal(t, 2, counter)
}

func TestGetFileStorage(t *testing.T) {
	d := azure.NewAzureStorageDiscovery(
		azure.WithSender(&mockStorageSender{}),
		azure.WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()
	assert.Nil(t, err)
	assert.NotNil(t, list)

	fileStorage, ok := list[2].(*voc.FileStorage)

	assert.True(t, ok)
	assert.Equal(t, "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/shares/fileshare1", string(fileStorage.ID))
	assert.Equal(t, "fileshare1", fileStorage.Name)
	assert.NotNil(t, fileStorage.CreationTime)
	assert.Equal(t, "FileStorage", fileStorage.Type[0])
	assert.Equal(t, "eastus", fileStorage.GeoLocation.Region)
	assert.Equal(t, true, fileStorage.AtRestEncryption.GetAtRestEncryption().Enabled)
	assert.Equal(t, "https://account1.file.core.windows.net/fileshare1", fileStorage.HttpEndpoint.Url)
	assert.Equal(t, true, fileStorage.HttpEndpoint.TransportEncryption.Enabled)
	assert.Equal(t, true, fileStorage.HttpEndpoint.TransportEncryption.Enforced)
	assert.Equal(t, "TLS1_2", fileStorage.HttpEndpoint.TransportEncryption.TlsVersion)
	assert.Equal(t, "", fileStorage.HttpEndpoint.TransportEncryption.Algorithm)
}

func TestListBlockStorage(t *testing.T) {
	d := azure.NewAzureStorageDiscovery(
		azure.WithSender(&mockStorageSender{}),
		azure.WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()
	assert.Nil(t, err)
	assert.NotNil(t, list)

	var counter int
	for _, item := range list {
		if item.GetType()[0] == "BlockStorage" {
			counter++
		}
	}
	assert.Equal(t, 2, counter)
}


func TestGetBlockStorage(t *testing.T) {
	d := azure.NewAzureStorageDiscovery(
		azure.WithSender(&mockStorageSender{}),
		azure.WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()
	assert.Nil(t, err)
	assert.NotNil(t, list)

	blockStorage, ok := list[4].(*voc.BlockStorage)

	assert.True(t, ok)
	assert.Equal(t, "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk1", string(blockStorage.ID))
	assert.Equal(t, "disk1", blockStorage.Name)
	assert.NotNil(t, blockStorage.CreationTime)
	assert.Equal(t, "BlockStorage", blockStorage.Type[0])
	assert.Equal(t, "eastus", blockStorage.GeoLocation.Region)
	assert.Equal(t, true, blockStorage.AtRestEncryption.GetAtRestEncryption().Enabled)
}

