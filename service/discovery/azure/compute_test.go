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
	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/compute/mgmt/compute"
	"io"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	"testing"
	"time"

	"clouditor.io/clouditor/voc"
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
	Value []compute.VirtualMachine `json:"value,omitempty"`
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
						"diagnosticsProfile": map[string]interface{}{
							"bootDiagnostics": map[string]interface{}{
								"enabled":    true,
								"storageUri": "https://logstoragevm1.blob.core.windows.net/",
							},
						},
					},
				},
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm2",
					"name":     "vm2",
					"location": "eastus",
					"properties": map[string]interface{}{
						"diagnosticsProfile": map[string]interface{}{
							"bootDiagnostics": map[string]interface{}{
								"enabled":    true,
								"storageUri": "",
							},
						},
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1" {
		return createResponse(map[string]interface{}{
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
						"enabled":    false,
						"storageUri": "test",
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
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm2" {
		return createResponse(map[string]interface{}{
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
						"storageUri": "test",
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
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Web/sites" {
		return createResponse(map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":         "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function1",
					"name":       "function1",
					"location":   "West Europe",
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

	assert.NotNil(t, err)
	assert.Nil(t, list)
	assert.Equal(t, "could not authorize Azure account: no authorized was available", err.Error())
}

func TestCompute(t *testing.T) {
	d := NewAzureComputeDiscovery(
		WithSender(&mockComputeSender{}),
		WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 3, len(list))
	assert.NotEmpty(t, d.Name())
}

func TestVirtualMachine(t *testing.T) {
	d := NewAzureComputeDiscovery(
		WithSender(&mockComputeSender{}),
		WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()
	assert.Nil(t, err)

	virtualMachine, ok := list[0].(*voc.VirtualMachine)

	assert.True(t, ok)
	assert.Equal(t, "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1", string(virtualMachine.ID))
	assert.Equal(t, "vm1", virtualMachine.Name)
	assert.Equal(t, 2, len(virtualMachine.NetworkInterface))
	assert.Equal(t, 3, len(virtualMachine.BlockStorage))

	assert.Equal(t, "data_disk_1", string(virtualMachine.BlockStorage[1]))
	assert.Equal(t, "123", string(virtualMachine.NetworkInterface[0]))
	assert.Equal(t, "eastus", virtualMachine.GeoLocation.Region)
	assert.Equal(t, true, virtualMachine.BootLog.Enabled)
	assert.Equal(t, voc.ResourceID("https://logstoragevm1.blob.core.windows.net/"), virtualMachine.BootLog.Output[0])
	assert.Equal(t, time.Duration(0), virtualMachine.BootLog.RetentionPeriod)

	virtualMachine2, ok := list[1].(*voc.VirtualMachine)
	assert.True(t, ok)
	assert.Equal(t, voc.ResourceID(""), virtualMachine2.BootLog.Output[0])
}

func TestFunction(t *testing.T) {
	d := NewAzureComputeDiscovery(
		WithSender(&mockComputeSender{}),
		WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 3, len(list))

	function, ok := list[2].(*voc.Function)

	assert.True(t, ok)
	assert.Equal(t, "function1", function.Name)
}

func TestComputeDiscoverFunctionWhenInputIsInvalid(t *testing.T) {
	d := azureComputeDiscovery{}

	discoverFunctionResponse, err := d.discoverFunction()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not list functions")
	assert.Nil(t, discoverFunctionResponse)
}

func TestComputeDiscoverVirtualMachines(t *testing.T) {
	d := azureComputeDiscovery{}

	discoverVirtualMachineResponse, err := d.discoverVirtualMachines()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not list virtual machines")
	assert.Nil(t, discoverVirtualMachineResponse)
}

func TestGetBootLogOutput(t *testing.T) {
	// Get mocked compute.VirtualMachine
	reqURL := "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Compute/virtualMachines"
	mockedVirtualMachinesResponse, err := getMockedVirtualMachines(reqURL)
	if err != nil {
		fmt.Println("error getting mocked storage account object: %w", err)
	}

	virtualMachine := mockedVirtualMachinesResponse[0]

	assert.NotEmpty(t, virtualMachine)
	// Delete the "diagnosticsProfile" property
	virtualMachine.DiagnosticsProfile = nil

	getBootLogOutputResponse := BootLogOutput(&virtualMachine)

	assert.Empty(t, getBootLogOutputResponse)
}

// getMockedVirtualMachines returns the mocked virtualMachines list
func getMockedVirtualMachines(reqUrl string) (virtualMachines []compute.VirtualMachine, err error) {
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
