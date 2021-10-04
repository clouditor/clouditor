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
							"keySource": "Microsoft.Storage",
						},
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

	storage, ok := list[0].(*voc.ObjectStorage)

	assert.True(t, ok)
	assert.Equal(t, "container1", storage.Name)
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

	storage, ok := list[2].(*voc.FileStorage)

	assert.True(t, ok)
	assert.Equal(t, "fileshare1", storage.Name)
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

	storage, ok := list[4].(*voc.BlockStorage)

	assert.True(t, ok)
	assert.Equal(t, "disk1", storage.Name)
}

