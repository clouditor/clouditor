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

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"

	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
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
	}

	return m.mockSender.Do(req)
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
	assert.Equal(t, 4, len(list))
	assert.NotEmpty(t, d.Name())
}

func TestVirtualMachine(t *testing.T) {
	d := NewAzureComputeDiscovery(
		WithSender(&mockComputeSender{}),
		WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()
	assert.NoError(t, err)

	virtualMachine, ok := list[0].(*voc.VirtualMachine)

	assert.True(t, ok)
	assert.Equal(t, "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1", string(virtualMachine.ID))
	assert.Equal(t, "vm1", virtualMachine.Name)
	assert.Equal(t, 2, len(virtualMachine.NetworkInterface))
	assert.Equal(t, 3, len(virtualMachine.BlockStorage))

	assert.Equal(t, "data_disk_1", string(virtualMachine.BlockStorage[1]))
	assert.Equal(t, "123", string(virtualMachine.NetworkInterface[0]))
	assert.Equal(t, "eastus", virtualMachine.GeoLocation.Region)
	assert.Equal(t, true, virtualMachine.BootLogging.Enabled)
	assert.Equal(t, voc.ResourceID("https://logstoragevm1.blob.core.windows.net/"), virtualMachine.BootLogging.LoggingService[0])
	assert.Equal(t, time.Duration(0), virtualMachine.BootLogging.RetentionPeriod)

	virtualMachine2, ok := list[1].(*voc.VirtualMachine)
	assert.True(t, ok)
	assert.Equal(t, []voc.ResourceID{}, virtualMachine2.BootLogging.LoggingService)

	virtualMachine3, ok := list[2].(*voc.VirtualMachine)
	assert.True(t, ok)
	assert.Equal(t, []voc.ResourceID{}, virtualMachine3.BlockStorage)
	assert.Equal(t, []voc.ResourceID{}, virtualMachine3.NetworkInterface)

}

func TestFunction(t *testing.T) {
	d := NewAzureComputeDiscovery(
		WithSender(&mockComputeSender{}),
		WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 4, len(list))

	function, ok := list[3].(*voc.Function)

	assert.True(t, ok)
	assert.Equal(t, "function1", function.Name)
}

func TestComputeDiscoverFunctionsWhenInputIsInvalid(t *testing.T) {
	d := azureComputeDiscovery{}

	discoverFunctionsResponse, err := d.discoverFunctions()

	assert.ErrorContains(t, err, ErrGettingNextPage.Error())
	assert.Nil(t, discoverFunctionsResponse)
}

func TestComputeDiscoverVirtualMachines(t *testing.T) {
	d := azureComputeDiscovery{}

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
				return assert.ErrorContains(t, err, "could not discover virtual machines:")
			},
		},
		{
			name: "Handle function error",
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
				return assert.ErrorContains(t, err, "could not discover virtual machines:")
			},
		},
		{
			name: "Without errors",
			fields: fields{
				azureDiscovery: azureDiscovery{
					cred: &mockAuthorizer{},
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: &mockComputeSender{},
						},
					},
				},
			},
			wantList: []voc.IsCloudResource{
				&voc.VirtualMachine{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1",
							Name:         "vm1",
							CreationTime: util.SafeTimestamp(&creationTime),
							Type:         []string{"VirtualMachine", "Compute", "Resource"},
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						NetworkInterface: []voc.ResourceID{"123", "234"},
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
					OSLogging: &voc.OSLogging{
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
							Name:         "vm2",
							CreationTime: util.SafeTimestamp(&time.Time{}),
							Type:         []string{"VirtualMachine", "Compute", "Resource"},
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						NetworkInterface: []voc.ResourceID{"987", "654"},
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
					OSLogging: &voc.OSLogging{
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
							Name:         "vm3",
							Type:         []string{"VirtualMachine", "Compute", "Resource"},
							CreationTime: util.SafeTimestamp(&time.Time{}),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						NetworkInterface: []voc.ResourceID{},
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
					OSLogging: &voc.OSLogging{
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
						NetworkInterface: []voc.ResourceID{},
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
							Transport: mockComputeSender{},
						},
					},
				},
			},

			want: []voc.IsCloudResource{
				&voc.Function{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function1",
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
						NetworkInterface: []voc.ResourceID{},
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
		azureDiscovery azureDiscovery
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
					NetworkInterface: []voc.ResourceID{},
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
		// TODO(anatheka): Wo kommt die panic her?
		//{
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
		//},
		{
			name: "No error",
			fields: fields{
				azureDiscovery: azureDiscovery{
					cred: &mockAuthorizer{},
					sub:  sub,
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockComputeSender{},
						},
					},
				},
			},
			want: []voc.IsCloudResource{
				&voc.VirtualMachine{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1",
							Name:         "vm1",
							CreationTime: util.SafeTimestamp(&creationTime),
							Type:         []string{"VirtualMachine", "Compute", "Resource"},
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						NetworkInterface: []voc.ResourceID{"123", "234"},
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
					OSLogging: &voc.OSLogging{
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
							Name:         "vm2",
							CreationTime: util.SafeTimestamp(&time.Time{}),
							Type:         []string{"VirtualMachine", "Compute", "Resource"},
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						NetworkInterface: []voc.ResourceID{"987", "654"},
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
					OSLogging: &voc.OSLogging{
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
							Name:         "vm3",
							Type:         []string{"VirtualMachine", "Compute", "Resource"},
							CreationTime: util.SafeTimestamp(&time.Time{}),
							Labels:       map[string]string{},
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
						},
						NetworkInterface: []voc.ResourceID{},
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
					OSLogging: &voc.OSLogging{
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
		azureDiscovery azureDiscovery
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
						Name:         "vm1",
						CreationTime: util.SafeTimestamp(&creationTime),
						Type:         []string{"VirtualMachine", "Compute", "Resource"},
						Labels:       map[string]string{},
						GeoLocation: voc.GeoLocation{
							Region: "eastus",
						},
					},
					NetworkInterface: []voc.ResourceID{"123", "234"},
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
				OSLogging: &voc.OSLogging{
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
