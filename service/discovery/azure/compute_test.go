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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v3"
	"github.com/stretchr/testify/assert"
)

type mockComputeSender struct {
	mockSender
}

func newMockComputeSender() *mockComputeSender {
	m := &mockComputeSender{}
	return m
}

type mockedVirtualMachinesResponse struct {
	Value []armcompute.VirtualMachine `json:"value,omitempty"`
}

func (m mockComputeSender) Do(req *http.Request) (res *http.Response, err error) {
	if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Compute/virtualMachines" {
		return createResponse(map[string]interface{}{
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
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res2/providers/Microsoft.Compute/disks/anotherdisk",
					"name":     "anotherdisk",
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
		return createResponse(map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res2/providers/Microsoft.Compute/disks/anotherdisk",
					"name":     "anotherdisk",
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
		return createResponse(map[string]interface{}{
			"value": &[]map[string]interface{}{},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res2/providers/Microsoft.Web/sites" {
		return createResponse(map[string]interface{}{
			"value": &[]map[string]interface{}{},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Web/sites" {
		return createResponse(map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function1",
					"name":     "function1",
					"location": "West Europe",
					"tags": map[string]interface{}{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					"properties": map[string]interface{}{},
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

func TestNewAzureComputeDiscovery(t *testing.T) {
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
			want: &azureComputeDiscovery{
				&azureDiscovery{
					discovererComponent: ComputeComponent,
					csID:                discovery.DefaultCloudServiceID,
				},
				make(map[string]*defenderProperties),
			},
		},
		{
			name: "With sender",
			args: args{
				opts: []DiscoveryOption{WithSender(mockComputeSender{})},
			},
			want: &azureComputeDiscovery{
				&azureDiscovery{
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockComputeSender{},
						},
					},
					discovererComponent: ComputeComponent,
					csID:                discovery.DefaultCloudServiceID,
				},
				make(map[string]*defenderProperties),
			},
		},
		{
			name: "With authorizer",
			args: args{
				opts: []DiscoveryOption{WithAuthorizer(&mockAuthorizer{})},
			},
			want: &azureComputeDiscovery{
				&azureDiscovery{
					cred:                &mockAuthorizer{},
					discovererComponent: ComputeComponent,
					csID:                discovery.DefaultCloudServiceID,
				},
				make(map[string]*defenderProperties),
			},
		},
		{
			name: "With cloud service ID",
			args: args{
				opts: []DiscoveryOption{WithCloudServiceID(testdata.MockCloudServiceID)},
			},
			want: &azureComputeDiscovery{
				&azureDiscovery{
					discovererComponent: ComputeComponent,
					csID:                testdata.MockCloudServiceID,
				},
				make(map[string]*defenderProperties),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewAzureComputeDiscovery(tt.args.opts...)
			assert.Equal(t, tt.want, d)
			assert.Equal(t, "Azure Compute", d.Name())
		})
	}
}

func TestAzureComputeAuthorizer(t *testing.T) {

	d := NewAzureComputeDiscovery()
	list, err := d.List()

	assert.Error(t, err)
	assert.Nil(t, list)
	assert.ErrorIs(t, err, ErrNoCredentialsConfigured)
}

func TestCompute(t *testing.T) {
	d := NewAzureComputeDiscovery(
		WithSender(&mockComputeSender{}),
		WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 7, len(list))
	assert.NotEmpty(t, d.Name())
}

func TestDiscoverer_List(t *testing.T) {
	d := NewAzureComputeDiscovery(
		WithSender(&mockComputeSender{}),
		WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()
	assert.NoError(t, err)

	virtualMachine, ok := list[3].(*voc.VirtualMachine)

	assert.True(t, ok)
	assert.Equal(t, "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1", string(virtualMachine.ID))
	assert.Equal(t, "vm1", virtualMachine.Name)
	assert.Equal(t, 2, len(virtualMachine.NetworkInterfaces))
	assert.Equal(t, 3, len(virtualMachine.BlockStorage))

	assert.Equal(t, "data_disk_1", string(virtualMachine.BlockStorage[1]))
	assert.Equal(t, "123", string(virtualMachine.NetworkInterfaces[0]))
	assert.Equal(t, "eastus", virtualMachine.GeoLocation.Region)
	assert.Equal(t, true, virtualMachine.BootLogging.Enabled)
	assert.Equal(t, voc.ResourceID("https://logstoragevm1.blob.core.windows.net/"), virtualMachine.BootLogging.LoggingService[0])
	assert.Equal(t, time.Duration(0), virtualMachine.BootLogging.RetentionPeriod)

	virtualMachine2, ok := list[4].(*voc.VirtualMachine)
	assert.True(t, ok)
	assert.Equal(t, []voc.ResourceID{}, virtualMachine2.BootLogging.LoggingService)

	virtualMachine3, ok := list[5].(*voc.VirtualMachine)
	assert.True(t, ok)
	assert.Equal(t, []voc.ResourceID{}, virtualMachine3.BlockStorage)
	assert.Equal(t, []voc.ResourceID{}, virtualMachine3.NetworkInterfaces)

}

func TestFunction(t *testing.T) {
	d := NewAzureComputeDiscovery(
		WithSender(&mockComputeSender{}),
		WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 7, len(list))

	function, ok := list[6].(*voc.Function)

	assert.True(t, ok)
	assert.Equal(t, "function1", function.Name)
}

func TestComputeDiscoverFunctionsWhenInputIsInvalid(t *testing.T) {
	d := azureComputeDiscovery{azureDiscovery: &azureDiscovery{}}

	discoverFunctionsResponse, err := d.discoverFunctions()

	assert.ErrorContains(t, err, ErrGettingNextPage.Error())
	assert.Nil(t, discoverFunctionsResponse)
}

func TestComputeDiscoverVirtualMachines(t *testing.T) {
	d := azureComputeDiscovery{azureDiscovery: &azureDiscovery{}}

	discoverVirtualMachineResponse, err := d.discoverVirtualMachines()

	assert.ErrorContains(t, err, ErrGettingNextPage.Error())
	assert.Nil(t, discoverVirtualMachineResponse)
}

func TestBootLogOutput(t *testing.T) {
	// Get mocked compute.VirtualMachine
	reqURL := "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Compute/virtualMachines"
	mockedVirtualMachinesResponse, err := mockedVirtualMachines(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	virtualMachine := mockedVirtualMachinesResponse[0]

	assert.NotEmpty(t, virtualMachine)
	// Delete the "diagnosticsProfile" property
	virtualMachine.Properties.DiagnosticsProfile = nil

	getBootLogOutputResponse := bootLogOutput(&virtualMachine)

	assert.Empty(t, getBootLogOutputResponse)
}

// mockedVirtualMachines returns the mocked virtualMachines list
func mockedVirtualMachines(reqUrl string) (virtualMachines []armcompute.VirtualMachine, err error) {
	var mockedVirtualMachinesResponse mockedVirtualMachinesResponse

	m := newMockComputeSender()
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return virtualMachines, fmt.Errorf("error creating new request: %w", err)
	}
	resp, err := m.Do(req)
	if err != nil || resp.StatusCode == 404 {
		return virtualMachines, fmt.Errorf("error getting mock http response: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("error io.ReadCloser: %w", err)
		}
	}(resp.Body)
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return virtualMachines, fmt.Errorf("error read all: %w", err)
	}
	err = json.Unmarshal(responseBody, &mockedVirtualMachinesResponse)
	if err != nil {
		return virtualMachines, fmt.Errorf("error unmarshalling: %w", err)
	}

	virtualMachines = mockedVirtualMachinesResponse.Value

	return virtualMachines, nil
}

func Test_azureComputeDiscovery_List(t *testing.T) {
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
		// {
		// 	name: "Authorize error",
		// 	fields: fields{
		// 		azureDiscovery: azureDiscovery{
		// 			cred: nil,
		// 		},
		// 	},
		// 	wantList: nil,
		// 	wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
		// 		return assert.ErrorContains(t, err, ErrCouldNotAuthenticate.Error())
		// 	},
		// },
		{
			name: "Discovery error",
			fields: fields{
				// Intentionally use wrong sender
				azureDiscovery: NewMockAzureDiscovery(newMockNetworkSender()),
			},
			wantList: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not discover block storage:")
			},
		},
		{
			name: "Without errors",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockComputeSender()),
			},
			wantList: []voc.IsCloudResource{
				&voc.BlockStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk1",
							ServiceID:    testdata.MockCloudServiceID,
							Name:         "disk1",
							CreationTime: util.SafeTimestamp(&creationTime),
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Labels: map[string]string{},
							Type:   voc.BlockStorageType,
						},

						AtRestEncryption: &voc.ManagedKeyEncryption{
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
							ServiceID:    testdata.MockCloudServiceID,
							Name:         "disk2",
							CreationTime: util.SafeTimestamp(&creationTime),
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Labels: map[string]string{},
							Type:   voc.BlockStorageType,
						},
						AtRestEncryption: &voc.CustomerKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "",
								Enabled:   true,
							},
							KeyUrl: "https://keyvault1.vault.azure.net/keys/customer-key/6273gdb374jz789hjm17819283748382",
						},
					},
				},
				&voc.BlockStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res2/providers/Microsoft.Compute/disks/anotherdisk",
							ServiceID:    testdata.MockCloudServiceID,
							Name:         "anotherdisk",
							CreationTime: util.SafeTimestamp(&creationTime),
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Labels: map[string]string{},
							Type:   voc.BlockStorageType,
						},

						AtRestEncryption: &voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
					},
				},
				&voc.VirtualMachine{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1",
							ServiceID:    testdata.MockCloudServiceID,
							Name:         "vm1",
							CreationTime: util.SafeTimestamp(&creationTime),
							Type:         voc.VirtualMachineType,
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						NetworkInterfaces: []voc.ResourceID{"123", "234"},
					},
					BlockStorage: []voc.ResourceID{"os_test_disk", "data_disk_1", "data_disk_2"},
					BootLogging: &voc.BootLogging{
						Logging: &voc.Logging{
							Enabled:        true,
							LoggingService: []voc.ResourceID{"https://logstoragevm1.blob.core.windows.net/"},
							Auditing: &voc.Auditing{
								SecurityFeature: &voc.SecurityFeature{},
							},
							RetentionPeriod: 0,
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
					MalwareProtection: &voc.MalwareProtection{},
				},
				&voc.VirtualMachine{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm2",
							ServiceID:    testdata.MockCloudServiceID,
							Name:         "vm2",
							CreationTime: util.SafeTimestamp(&time.Time{}),
							Type:         voc.VirtualMachineType,
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
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
							Enabled:         false,
							LoggingService:  []voc.ResourceID{},
							RetentionPeriod: 0,
							Auditing: &voc.Auditing{
								SecurityFeature: &voc.SecurityFeature{},
							},
						},
					},
					MalwareProtection: &voc.MalwareProtection{},
				},
				&voc.VirtualMachine{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm3",
							ServiceID:    testdata.MockCloudServiceID,
							Name:         "vm3",
							Type:         voc.VirtualMachineType,
							CreationTime: util.SafeTimestamp(&time.Time{}),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
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
					MalwareProtection: &voc.MalwareProtection{},
				},
				&voc.Function{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function1",
							ServiceID:    testdata.MockCloudServiceID,
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
						},
						NetworkInterfaces: []voc.ResourceID{},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "With resource group",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockComputeSender(), WithResourceGroup("res2")),
			},
			wantList: []voc.IsCloudResource{
				&voc.BlockStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res2/providers/Microsoft.Compute/disks/anotherdisk",
							ServiceID:    testdata.MockCloudServiceID,
							Name:         "anotherdisk",
							CreationTime: util.SafeTimestamp(&creationTime),
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Labels: map[string]string{},
							Type:   voc.BlockStorageType,
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
			d := &azureComputeDiscovery{
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

func Test_azureComputeDiscovery_discoverFunctions(t *testing.T) {
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
				azureDiscovery: NewMockAzureDiscovery(newMockComputeSender()),
			},

			want: []voc.IsCloudResource{
				&voc.Function{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function1",
							ServiceID:    testdata.MockCloudServiceID,
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
						},
						NetworkInterfaces: []voc.ResourceID{},
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureComputeDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}
			got, err := d.discoverFunctions()
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_azureComputeDiscovery_handleFunction(t *testing.T) {
	functionID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function1"
	functionName := "function1"
	diskRegion := "West Europe"
	testTag1 := "testTag1"
	testTag2 := "testTag2"

	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		function *armappservice.Site
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   voc.IsCompute
	}{
		{
			name: "Empty input",
			args: args{
				function: nil,
			},
			want: nil,
		},
		{
			name: "Empty functionID",
			args: args{
				function: &armappservice.Site{
					ID: nil,
				},
			},
			want: nil,
		},
		{
			name: "No error",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockComputeSender()),
			},
			args: args{
				function: &armappservice.Site{
					ID:       &functionID,
					Name:     &functionName,
					Location: &diskRegion,
					Tags: map[string]*string{
						"testKey1": &testTag1,
						"testKey2": &testTag2,
					},
				},
			},
			want: &voc.Function{
				Compute: &voc.Compute{
					Resource: &voc.Resource{
						ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function1",
						ServiceID:    testdata.MockCloudServiceID,
						Name:         "function1",
						CreationTime: util.SafeTimestamp(&time.Time{}),
						Type:         []string{"Function", "Compute", "Resource"},
						Labels: map[string]string{
							"testKey1": testTag1,
							"testKey2": testTag2,
						},
						GeoLocation: voc.GeoLocation{
							Region: "West Europe",
						},
					},
					NetworkInterfaces: []voc.ResourceID{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			az := &azureComputeDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}
			assert.Equalf(t, tt.want, az.handleFunction(tt.args.function), "handleFunction(%v)", tt.args.function)
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
		want    []voc.IsCloudResource
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO(anatheka): Wo kommt die panic her?
		// {
		//	name: "Error list pages",
		//	fields: fields{
		//		azureDiscovery: azureDiscovery{
		//			cred: nil,
		//			sub:  sub,
		//		},
		//	},
		//	want: nil,
		//	wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
		//		return assert.ErrorContains(t, err, ErrGettingNextPage.Error())
		//	},
		// },
		{
			name: "No error",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockComputeSender()),
			},
			want: []voc.IsCloudResource{
				&voc.VirtualMachine{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1",
							ServiceID:    testdata.MockCloudServiceID,
							Name:         "vm1",
							CreationTime: util.SafeTimestamp(&creationTime),
							Type:         []string{"VirtualMachine", "Compute", "Resource"},
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						NetworkInterfaces: []voc.ResourceID{"123", "234"},
					},
					BlockStorage: []voc.ResourceID{"os_test_disk", "data_disk_1", "data_disk_2"},
					BootLogging: &voc.BootLogging{
						Logging: &voc.Logging{
							Enabled:        true,
							LoggingService: []voc.ResourceID{"https://logstoragevm1.blob.core.windows.net/"},
							Auditing: &voc.Auditing{
								SecurityFeature: &voc.SecurityFeature{},
							},
							RetentionPeriod: 0,
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
					MalwareProtection: &voc.MalwareProtection{},
				},
				&voc.VirtualMachine{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm2",
							ServiceID:    testdata.MockCloudServiceID,
							Name:         "vm2",
							CreationTime: util.SafeTimestamp(&time.Time{}),
							Type:         []string{"VirtualMachine", "Compute", "Resource"},
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
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
							Enabled:         false,
							LoggingService:  []voc.ResourceID{},
							RetentionPeriod: 0,
							Auditing: &voc.Auditing{
								SecurityFeature: &voc.SecurityFeature{},
							},
						},
					},
					MalwareProtection: &voc.MalwareProtection{},
				},
				&voc.VirtualMachine{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm3",
							ServiceID:    testdata.MockCloudServiceID,
							Name:         "vm3",
							Type:         []string{"VirtualMachine", "Compute", "Resource"},
							CreationTime: util.SafeTimestamp(&time.Time{}),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
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
					MalwareProtection: &voc.MalwareProtection{},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureComputeDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}
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
	ID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1"
	name := "vm1"
	region := "eastus"
	netInterface1 := "123"
	netInterface2 := "234"
	netInterfaces := armcompute.NetworkProfile{
		NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
			{
				ID: &netInterface1,
			},
			{
				ID: &netInterface2,
			},
		},
	}
	dataDisk1 := "data_disk_1"
	dataDisk2 := "data_disk_2"
	dataDisks := []*armcompute.DataDisk{
		{
			ManagedDisk: &armcompute.ManagedDiskParameters{
				ID: &dataDisk1,
			},
		},
		{
			ManagedDisk: &armcompute.ManagedDiskParameters{
				ID: &dataDisk2,
			},
		},
	}
	osDisk := "os_test_disk"
	storageUri := "https://logstoragevm1.blob.core.windows.net/"
	enabledTrue := true

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
		want    voc.IsCompute
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
				azureDiscovery: NewMockAzureDiscovery(newMockComputeSender()),
			},
			args: args{
				vm: &armcompute.VirtualMachine{
					ID:       &ID,
					Name:     &name,
					Location: &region,
					Properties: &armcompute.VirtualMachineProperties{
						TimeCreated:    &creationTime,
						NetworkProfile: &netInterfaces,
						StorageProfile: &armcompute.StorageProfile{
							OSDisk: &armcompute.OSDisk{
								ManagedDisk: &armcompute.ManagedDiskParameters{
									ID: &osDisk,
								},
							},
							DataDisks: dataDisks,
						},
						DiagnosticsProfile: &armcompute.DiagnosticsProfile{
							BootDiagnostics: &armcompute.BootDiagnostics{
								Enabled:    &enabledTrue,
								StorageURI: &storageUri,
							},
						},
					},
				},
			},
			want: &voc.VirtualMachine{
				Compute: &voc.Compute{
					Resource: &voc.Resource{
						ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1",
						ServiceID:    testdata.MockCloudServiceID,
						Name:         "vm1",
						CreationTime: util.SafeTimestamp(&creationTime),
						Type:         []string{"VirtualMachine", "Compute", "Resource"},
						Labels:       map[string]string{},
						GeoLocation: voc.GeoLocation{
							Region: "eastus",
						},
					},
					NetworkInterfaces: []voc.ResourceID{"123", "234"},
				},
				BlockStorage: []voc.ResourceID{"os_test_disk", "data_disk_1", "data_disk_2"},
				BootLogging: &voc.BootLogging{
					Logging: &voc.Logging{
						Enabled:        true,
						LoggingService: []voc.ResourceID{"https://logstoragevm1.blob.core.windows.net/"},
						Auditing: &voc.Auditing{
							SecurityFeature: &voc.SecurityFeature{},
						},
						RetentionPeriod: 0,
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
				MalwareProtection: &voc.MalwareProtection{},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureComputeDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}
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
			want: storageUri,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, bootLogOutput(tt.args.vm))
		})
	}
}

func TestBlockStoragesHandleMethodsWhenInputIsInvalid(t *testing.T) {
	d := azureComputeDiscovery{}

	// Test method handleBlockStorage
	disk := &armcompute.Disk{}
	handleBlockStorageResponse, err := d.handleBlockStorage(disk)
	assert.Error(t, err)
	assert.Nil(t, handleBlockStorageResponse)
}

func Test_azureComputeDiscovery_discoverBlockStorage(t *testing.T) {
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)

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
				azureDiscovery: NewMockAzureDiscovery(newMockComputeSender()),
			},
			want: []voc.IsCloudResource{
				&voc.BlockStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk1",
							ServiceID:    testdata.MockCloudServiceID,
							Name:         "disk1",
							CreationTime: util.SafeTimestamp(&creationTime),
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Type:   []string{"BlockStorage", "Storage", "Resource"},
							Labels: map[string]string{},
						},

						AtRestEncryption: &voc.ManagedKeyEncryption{
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
							ServiceID:    testdata.MockCloudServiceID,
							Name:         "disk2",
							CreationTime: util.SafeTimestamp(&creationTime),
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Type:   []string{"BlockStorage", "Storage", "Resource"},
							Labels: map[string]string{},
						},
						AtRestEncryption: &voc.CustomerKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "",
								Enabled:   true,
							},
							KeyUrl: "https://keyvault1.vault.azure.net/keys/customer-key/6273gdb374jz789hjm17819283748382",
						},
					},
				},
				&voc.BlockStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res2/providers/Microsoft.Compute/disks/anotherdisk",
							ServiceID:    testdata.MockCloudServiceID,
							Name:         "anotherdisk",
							CreationTime: util.SafeTimestamp(&creationTime),
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Labels: map[string]string{},
							Type:   voc.BlockStorageType,
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
			d := &azureComputeDiscovery{
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
				azureDiscovery: NewMockAzureDiscovery(newMockComputeSender()),
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
				azureDiscovery: NewMockAzureDiscovery(newMockComputeSender()),
			},
			want: &voc.BlockStorage{
				Storage: &voc.Storage{
					Resource: &voc.Resource{
						ID:           voc.ResourceID(diskID),
						ServiceID:    testdata.MockCloudServiceID,
						Name:         "disk1",
						CreationTime: util.SafeTimestamp(&creationTime),
						Type:         []string{"BlockStorage", "Storage", "Resource"},
						GeoLocation: voc.GeoLocation{
							Region: "eastus",
						},
						Labels: map[string]string{},
					},

					AtRestEncryption: &voc.CustomerKeyEncryption{
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
			d := &azureComputeDiscovery{
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
				azureDiscovery: NewMockAzureDiscovery(newMockComputeSender()),
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
				azureDiscovery: NewMockAzureDiscovery(newMockComputeSender()),
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
			d := &azureComputeDiscovery{
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
				azureDiscovery: NewMockAzureDiscovery(newMockComputeSender()),
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
				azureDiscovery: NewMockAzureDiscovery(newMockComputeSender()),
			},
			want:    "https://keyvault1.vault.azure.net/keys/customer-key/6273gdb374jz789hjm17819283748382",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureComputeDiscovery{
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
