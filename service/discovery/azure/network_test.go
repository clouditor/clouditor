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
	"net/http"
	"testing"

	"clouditor.io/clouditor/voc"
	"github.com/stretchr/testify/assert"
)

type mockNetworkSender struct {
	mockSender
}

func (m mockNetworkSender) Do(req *http.Request) (res *http.Response, err error) {
	if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Network/networkInterfaces" {
		return createResponse(map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkInterfaces/iface1",
					"name":     "iface1",
					"location": "eastus",
					"properties": map[string]interface{}{
						"networkSecurityGroup": map[string]interface{}{
							"id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkSecurityGroups/nsg1",
						},
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkSecurityGroups/nsg1" {
		return createResponse(map[string]interface{}{
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
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Network/loadBalancers" {
		return createResponse(map[string]interface{}{
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
		return createResponse(map[string]interface{}{
			"properties": map[string]interface{}{
				"ipAddress": "111.222.333.444",
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/publicIPAddresses/test-b9cb3645-25d0-4288-910a-020563f63b1d" {
		return createResponse(map[string]interface{}{
			"properties": map[string]interface{}{
				"ipAddress": nil,
			},
		}, 200)
	}

	return m.mockSender.Do(req)
}

func TestAzureNetworkAuthorizer(t *testing.T) {

	d := NewAzureNetworkDiscovery()
	list, err := d.List()

	assert.Error(t, err)
	assert.Nil(t, list)
	assert.Equal(t, "could not authorize Azure account: no authorized was available", err.Error())
}

func TestNetwork(t *testing.T) {
	d := NewAzureNetworkDiscovery(
		WithSender(&mockNetworkSender{}),
		WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 4, len(list))
	assert.NotEmpty(t, d.Name())

	iface, ok := list[0].(*voc.NetworkInterface)

	assert.True(t, ok)
	assert.Equal(t, "iface1", iface.Name)
	//assert.Equal(t, "*", iface.AccessRestriction.RestrictedPorts)

	lb, ok := list[1].(*voc.LoadBalancer)

	assert.True(t, ok)
	assert.Equal(t, "lb1", lb.Name)
	assert.Equal(t, int16(1234), lb.Ports[0])
	assert.Equal(t, "111.222.333.444", lb.Ips[0])

	lb, ok = list[2].(*voc.LoadBalancer)
	assert.True(t, ok)
	assert.Equal(t, "lb2", lb.Name)
	assert.Nil(t, lb.NetworkService.Ips)
}

func TestComputeDiscoverMethodsWhenInputIsInvalid(t *testing.T) {
	d := azureNetworkDiscovery{}

	// Test method discoverNetworkInterfaces
	discoverNetworkInterfacesResponse, err := d.discoverNetworkInterfaces()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not list network interfaces")
	assert.Nil(t, discoverNetworkInterfacesResponse)

	// Test method discoverLoadBalancer
	discoverLoadBalancerResponse, err := d.discoverLoadBalancer()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not list load balancer")
	assert.Nil(t, discoverLoadBalancerResponse)

}
