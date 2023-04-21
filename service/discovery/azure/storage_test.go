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
	"reflect"
	"testing"
	"time"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
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
		return createResponse(map[string]interface{}{
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
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.DataProtection/backupVaults" {
		return createResponse(map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1",
					"name":     "backupAccount1",
					"location": "westeurope",
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupInstances" {
		return createResponse(map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupInstances/account1-account1-22222222-2222-2222-2222-222222222222",
					"name": "account1-account1-22222222-2222-2222-2222-222222222222",
					"properties": map[string]interface{}{
						"dataSourceInfo": map[string]interface{}{
							"resourceID":     "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1",
							"datasourceType": "Microsoft.Storage/storageAccounts/blobServices",
						},
						"policyInfo": map[string]interface{}{
							"policyId": "policyId",
						},
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
					backupMap:           make(map[string]map[string]*voc.Backup),
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
					backupMap:           make(map[string]map[string]*voc.Backup),
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
					backupMap:           make(map[string]map[string]*voc.Backup),
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
	assert.Equal(t, 8, len(list))
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
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container1",
							ServiceID:    testdata.MockCloudServiceID,
							Name:         "container1",
							Type:         voc.ObjectStorageType,
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						AtRestEncryption: &voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
						Immutability: &voc.Immutability{Enabled: false},
						Backup: &voc.Backup{
							Enabled:         true,
							RetentionPeriod: 0,
							GeoLocation:     voc.GeoLocation{Region: "westeurope"},
							Policy:          "policyId",
							Storage:         voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupInstances/account1-account1-22222222-2222-2222-2222-222222222222"),
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
					},
					PublicAccess: true,
				},
				&voc.ObjectStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container2",
							ServiceID:    testdata.MockCloudServiceID,
							Name:         "container2",
							Type:         voc.ObjectStorageType,
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						AtRestEncryption: &voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
						Immutability: &voc.Immutability{Enabled: false},
						Backup: &voc.Backup{
							Enabled:         true,
							RetentionPeriod: 0,
							GeoLocation:     voc.GeoLocation{Region: "westeurope"},
							Policy:          "policyId",
							Storage:         voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupInstances/account1-account1-22222222-2222-2222-2222-222222222222"),
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
					},
					PublicAccess: true,
				},
				&voc.FileStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/fileshare1",
							ServiceID:    testdata.MockCloudServiceID,
							Name:         "fileshare1",
							Type:         voc.FileStorageType,
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						AtRestEncryption: &voc.ManagedKeyEncryption{
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
							ServiceID:    testdata.MockCloudServiceID,
							Name:         "fileshare2",
							Type:         voc.FileStorageType,
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						AtRestEncryption: &voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
					},
				},
				&voc.ObjectStorageService{
					StorageService: &voc.StorageService{
						Storage: []voc.ResourceID{
							"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container1",
							"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container2",
							"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/fileshare1",
							"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/fileshare2",
						},
						NetworkService: &voc.NetworkService{
							Networking: &voc.Networking{
								Resource: &voc.Resource{
									ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1",
									ServiceID:    testdata.MockCloudServiceID,
									Name:         "account1",
									Type:         voc.ObjectStorageServiceType,
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
					HttpEndpoint: &voc.HttpEndpoint{
						Url: "https://account1.[file,blob].core.windows.net/",
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
							ServiceID:    testdata.MockCloudServiceID,
							Name:         "container3",
							Type:         voc.ObjectStorageType,
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						AtRestEncryption: &voc.CustomerKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "",
								Enabled:   true,
							},
							KeyUrl: "https://testvault.vault.azure.net/keys/testkey/123456",
						},
						Immutability: &voc.Immutability{Enabled: false},
					},
					PublicAccess: false,
				},
				&voc.ObjectStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2/blobServices/default/containers/container4",
							ServiceID:    testdata.MockCloudServiceID,
							Name:         "container4",
							Type:         voc.ObjectStorageType,
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						AtRestEncryption: &voc.CustomerKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "",
								Enabled:   true,
							},
							KeyUrl: "https://testvault.vault.azure.net/keys/testkey/123456",
						},
						Immutability: &voc.Immutability{Enabled: false},
					},
					PublicAccess: false,
				},
				&voc.ObjectStorageService{
					StorageService: &voc.StorageService{
						Storage: []voc.ResourceID{
							"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2/blobServices/default/containers/container3",
							"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2/blobServices/default/containers/container4",
						},
						NetworkService: &voc.NetworkService{
							Networking: &voc.Networking{
								Resource: &voc.Resource{
									ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2",
									ServiceID:    testdata.MockCloudServiceID,
									Name:         "account2",
									Type:         voc.ObjectStorageServiceType,
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
					HttpEndpoint: &voc.HttpEndpoint{
						Url: "https://account1.[file,blob].core.windows.net/",
						TransportEncryption: &voc.TransportEncryption{
							Enforced:   true,
							Enabled:    true,
							TlsVersion: "TLS1_2",
							Algorithm:  "TLS",
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
	d := azureStorageDiscovery{&azureDiscovery{csID: testdata.MockCloudServiceID}, make(map[string]*defenderProperties)}

	// Get mocked armstorage.Account
	reqURL := "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Storage/storageAccounts/account3"
	mockedStorageAccountObject, err := mockedStorageAccount(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	// Test method handleObjectStorage
	containerItem := armstorage.ListContainerItem{}
	handleObjectStorageRespone, err := d.handleObjectStorage(mockedStorageAccountObject, &containerItem)
	assert.Error(t, err)
	assert.Nil(t, handleObjectStorageRespone)

	// Test method handleFileStorage
	fileShare := &armstorage.FileShareItem{}
	handleFileStorageRespone, err := d.handleFileStorage(mockedStorageAccountObject, fileShare)
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
	discoverObjectStoragesResponse, err := d.discoverObjectStorages(mockedStorageAccountObject)
	assert.ErrorContains(t, err, ErrGettingNextPage.Error())
	assert.Nil(t, discoverObjectStoragesResponse)

	// Test method discoverFileStorages
	discoverFileStoragesResponse, err := d.discoverFileStorages(mockedStorageAccountObject)
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
				assert.Equal(t, 8, len(got))
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
		account   *armstorage.Account
		fileshare *armstorage.FileShareItem
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
						ServiceID:    testdata.MockCloudServiceID,
						Name:         fileShareName,
						CreationTime: util.SafeTimestamp(&creationTime),
						GeoLocation: voc.GeoLocation{
							Region: accountRegion,
						},
						Labels: map[string]string{},
						Type:   voc.FileStorageType,
					},
					AtRestEncryption: &voc.ManagedKeyEncryption{
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
			d := &azureStorageDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}

			got, err := d.handleFileStorage(tt.args.account, tt.args.fileshare)
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
		account      *armstorage.Account
		storagesList []voc.IsCloudResource
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
			want: &voc.ObjectStorageService{
				StorageService: &voc.StorageService{
					NetworkService: &voc.NetworkService{
						Networking: &voc.Networking{
							Resource: &voc.Resource{
								ID:           voc.ResourceID(accountID),
								ServiceID:    testdata.MockCloudServiceID,
								Name:         accountName,
								CreationTime: util.SafeTimestamp(&creationTime),
								Type:         voc.ObjectStorageServiceType,
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
	// accountID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1"
	accountRegion := "eastus"
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)
	keySource := armstorage.KeySourceMicrosoftStorage
	containerID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container1"
	containerName := "container1"
	immutability := false
	publicAccess := armstorage.PublicAccessNone

	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		account   *armstorage.Account
		container *armstorage.ListContainerItem
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *voc.ObjectStorage
		wantErr assert.ErrorAssertionFunc
	}{
		// {
		// 	name: "Account is empty",
		// 	want: nil,
		// 	wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
		// 		return assert.ErrorIs(t, err, ErrEmptyStorageAccount)
		// 	},
		// },
		// {
		// 	name: "Container is empty",
		// 	args: args{
		// 		account: &armstorage.Account{
		// 			ID: &accountID,
		// 		},
		// 	},
		// 	want: nil,
		// 	wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
		// 		return assert.ErrorContains(t, err, "container is nil")
		// 	},
		// },
		// {
		// 	name: "Error getting atRestEncryption properties",
		// 	args: args{
		// 		account: &armstorage.Account{
		// 			ID: &accountID,
		// 		},
		// 		container: &armstorage.ListContainerItem{
		// 			ID: &containerID,
		// 		},
		// 	},
		// 	want: nil,
		// 	wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
		// 		return assert.ErrorContains(t, err, "could not get object storage properties for the atRestEncryption:")
		// 	},
		// },
		{
			name: "No error",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockStorageSender()),
			},
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
					Properties: &armstorage.ContainerProperties{
						HasImmutabilityPolicy: &immutability,
						PublicAccess:          &publicAccess,
					},
				},
			},
			want: &voc.ObjectStorage{
				Storage: &voc.Storage{
					Resource: &voc.Resource{
						ID:           voc.ResourceID(containerID),
						ServiceID:    testdata.MockCloudServiceID,
						Name:         containerName,
						CreationTime: util.SafeTimestamp(&creationTime),
						GeoLocation: voc.GeoLocation{
							Region: accountRegion,
						},
						Labels: map[string]string{},
						Type:   voc.ObjectStorageType,
					},
					AtRestEncryption: &voc.ManagedKeyEncryption{
						AtRestEncryption: &voc.AtRestEncryption{
							Algorithm: "AES256",
							Enabled:   true,
						},
					},
					Immutability: &voc.Immutability{Enabled: false},
				},
				PublicAccess: false,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureStorageDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}

			got, err := d.handleObjectStorage(tt.args.account, tt.args.container)
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
							ServiceID:    testdata.MockCloudServiceID,
							Name:         "fileshare1",
							Type:         voc.FileStorageType,
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						AtRestEncryption: &voc.ManagedKeyEncryption{
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
							ServiceID:    testdata.MockCloudServiceID,
							Name:         "fileshare2",
							Type:         voc.FileStorageType,
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						AtRestEncryption: &voc.ManagedKeyEncryption{
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
			// initialize file share client
			_ = d.initFileStorageClient()

			got, err := d.discoverFileStorages(tt.args.account)
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
							ServiceID:    testdata.MockCloudServiceID,
							Name:         "container1",
							Type:         voc.ObjectStorageType,
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						AtRestEncryption: &voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
						Immutability: &voc.Immutability{Enabled: false},
					},
					PublicAccess: true,
				},
				&voc.ObjectStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container2",
							ServiceID:    testdata.MockCloudServiceID,
							Name:         "container2",
							Type:         voc.ObjectStorageType,
							CreationTime: util.SafeTimestamp(&creationTime),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						AtRestEncryption: &voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
						Immutability: &voc.Immutability{Enabled: false},
					},
					PublicAccess: true,
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

			got, err := d.discoverObjectStorages(tt.args.account)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_azureStorageDiscovery_createResourceLogging(t *testing.T) {
	type fields struct {
		azureDiscovery     *azureDiscovery
		defenderProperties map[string]*defenderProperties
	}
	tests := []struct {
		name                string
		fields              fields
		wantResourceLogging *voc.ResourceLogging
	}{
		{
			name: "Missing defenderProperties",
			fields: fields{
				azureDiscovery:     NewMockAzureDiscovery(newMockStorageSender()),
				defenderProperties: make(map[string]*defenderProperties),
			},
			wantResourceLogging: nil,
		},
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockStorageSender()),
				defenderProperties: map[string]*defenderProperties{
					DefenderStorageType: {
						monitoringLogDataEnabled: true,
						securityAlertsEnabled:    true,
					},
				},
			},
			wantResourceLogging: &voc.ResourceLogging{
				MonitoringLogDataEnabled: true,
				SecurityAlertsEnabled:    true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureStorageDiscovery{
				azureDiscovery:     tt.fields.azureDiscovery,
				defenderProperties: tt.fields.defenderProperties,
			}
			if gotResourceLogging := d.createResourceLogging(); !reflect.DeepEqual(gotResourceLogging, tt.wantResourceLogging) {
				t.Errorf("azureStorageDiscovery.createResourceLogging() = %v, want %v", gotResourceLogging, tt.wantResourceLogging)
			}
		})
	}
}

func Test_idUpToStorageAccount(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty input",
			args: args{},
			want: "",
		},
		{
			name: "Wrong input",
			args: args{
				id: "teststring",
			},
			want: "",
		},
		{
			name: "Happy path",
			args: args{
				id: "/subscriptions/XXXXXXXXXXXX-XXXX-XXXX-XXXXXXXXXXXX/resourceGroups/resourceGroupName/providers/Microsoft.Storage/storageAccounts/containerName/blobServices/default/containers/testContainer",
			},
			want: "/subscriptions/XXXXXXXXXXXX-XXXX-XXXX-XXXXXXXXXXXX/resourceGroups/resourceGroupName/providers/Microsoft.Storage/storageAccounts/containerName",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := idUpToStorageAccount(tt.args.id); got != tt.want {
				t.Errorf("idUpToStorageAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}
