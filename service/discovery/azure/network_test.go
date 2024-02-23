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

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/stretchr/testify/assert"
)

type mockNetworkSender struct {
	mockSender
}

func newMockNetworkSender() *mockNetworkSender {
	m := &mockNetworkSender{}
	return m
}

func (m mockNetworkSender) Do(req *http.Request) (res *http.Response, err error) {
	if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Network/networkInterfaces" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkInterfaces/iface1",
					"name":     "iface1",
					"location": "eastus",
					"properties": map[string]interface{}{
						"networkSecurityGroup": map[string]interface{}{
							"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkSecurityGroups/nsg1",
							"location": "eastus",
						},
					},
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkSecurityGroups/nsg1" {
		return createResponse(req, map[string]interface{}{
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
		return createResponse(req, map[string]interface{}{
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
										"properties": map[string]interface{}{
											"ipAddress": "111.222.333.444",
										},
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
		return createResponse(req, map[string]interface{}{
			"properties": map[string]interface{}{
				"ipAddress": "111.222.333.444",
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/publicIPAddresses/test-b9cb3645-25d0-4288-910a-020563f63b1d" {
		return createResponse(req, map[string]interface{}{
			"properties": map[string]interface{}{
				"ipAddress": nil,
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Network/applicationGateways" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/applicationGateways/appgw1",
					"name":     "appgw1",
					"location": "eastus",
					"properties": map[string]interface{}{
						"webApplicationFirewallConfiguration": map[string]interface{}{
							"enabled": true,
						},
					},
				},
			},
		}, 200)
	}
	// return createResponse(req, map[string]interface{}{
	// 	"value": &[]map[string]interface{}{

	return m.mockSender.Do(req)
}

func TestAzureNetworkAuthorizer(t *testing.T) {

	d := NewAzureNetworkDiscovery()
	list, err := d.List()

	assert.Error(t, err)
	assert.Nil(t, list)
	assert.ErrorIs(t, err, ErrNoCredentialsConfigured)
}

func Test_azureNetworkDiscovery_List(t *testing.T) {
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
			name: "Discovery error",
			fields: fields{
				// Intentionally use wrong sender
				azureDiscovery: NewMockAzureDiscovery(newMockStorageSender()),
			},
			wantList: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not discover network interfaces:")
			},
		},
		{
			name: "Without errors",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockNetworkSender()),
			},
			wantList: []voc.IsCloudResource{
				&voc.NetworkInterface{
					Networking: &voc.Networking{
						Resource: &voc.Resource{
							ID:        "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.network/networkinterfaces/iface1",
							ServiceID: testdata.MockCloudServiceID1,
							Name:      "iface1",
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Type:   voc.NetworkInterfaceType,
							Labels: map[string]string{},
							Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
							Raw:    "{\"*armnetwork.Interface\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkInterfaces/iface1\",\"location\":\"eastus\",\"name\":\"iface1\",\"properties\":{\"networkSecurityGroup\":{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkSecurityGroups/nsg1\",\"location\":\"eastus\"}}}]}",
						},
					},
					AccessRestriction: &voc.L3Firewall{
						Enabled: true,
					},
				},
				&voc.LoadBalancer{
					NetworkService: &voc.NetworkService{
						Networking: &voc.Networking{
							Resource: &voc.Resource{
								ID:        "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.network/loadbalancers/lb1",
								ServiceID: testdata.MockCloudServiceID1,
								Name:      "lb1",
								GeoLocation: voc.GeoLocation{
									Region: "eastus",
								},
								Type:   voc.LoadBalancerType,
								Labels: map[string]string{},
								Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
								Raw:    "{\"*armnetwork.LoadBalancer\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb1\",\"location\":\"eastus\",\"name\":\"lb1\",\"properties\":{\"frontendIPConfigurations\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb1/frontendIPConfigurations/b9cb3645-25d0-4288-910a-020563f63b1c\",\"name\":\"b9cb3645-25d0-4288-910a-020563f63b1c\",\"properties\":{\"publicIPAddress\":{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/publicIPAddresses/test-b9cb3645-25d0-4288-910a-020563f63b1c\",\"properties\":{\"ipAddress\":\"111.222.333.444\"}}}}],\"loadBalancingRules\":[{\"properties\":{\"frontendPort\":1234}},{\"properties\":{\"frontendPort\":5678}}]}}]}",
							},
						},
						Ips:   []string{"111.222.333.444"},
						Ports: []uint16{1234, 5678},
					},
					HttpEndpoints: []*voc.HttpEndpoint{},
				},
				&voc.LoadBalancer{
					NetworkService: &voc.NetworkService{
						Networking: &voc.Networking{
							Resource: &voc.Resource{
								ID:        "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.network/loadbalancers/lb2",
								ServiceID: testdata.MockCloudServiceID1,
								Name:      "lb2",
								GeoLocation: voc.GeoLocation{
									Region: "eastus",
								},
								Type:   voc.LoadBalancerType,
								Labels: map[string]string{},
								Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
								Raw:    "{\"*armnetwork.LoadBalancer\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb2\",\"location\":\"eastus\",\"name\":\"lb2\",\"properties\":{\"frontendIPConfigurations\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb1/frontendIPConfigurations/b9cb3645-25d0-4288-910a-020563f63b1c\",\"name\":\"b9cb3645-25d0-4288-910a-020563f63b1c\",\"properties\":{}}],\"loadBalancingRules\":[{\"properties\":{\"frontendPort\":1234}},{\"properties\":{\"frontendPort\":5678}}]}}]}",
							},
						},
						Ports: []uint16{1234, 5678},
						Ips:   []string{},
					},
					HttpEndpoints: []*voc.HttpEndpoint{},
				},
				&voc.LoadBalancer{
					NetworkService: &voc.NetworkService{
						Networking: &voc.Networking{
							Resource: &voc.Resource{
								ID:        "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.network/loadbalancers/lb3",
								ServiceID: testdata.MockCloudServiceID1,
								Name:      "lb3",
								GeoLocation: voc.GeoLocation{
									Region: "eastus",
								},
								Type:   voc.LoadBalancerType,
								Labels: map[string]string{},
								Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
								Raw:    "{\"*armnetwork.LoadBalancer\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb3\",\"location\":\"eastus\",\"name\":\"lb3\",\"properties\":{\"frontendIPConfigurations\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb1/frontendIPConfigurations/b9cb3645-25d0-4288-910a-020563f63b1c\",\"name\":\"b9cb3645-25d0-4288-910a-020563f63b1c\",\"properties\":{\"publicIPAddress\":{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/publicIPAddresses/test-b9cb3645-25d0-4288-910a-020563f63b1d\"}}}],\"loadBalancingRules\":[{\"properties\":{\"frontendPort\":1234}},{\"properties\":{\"frontendPort\":5678}}]}}]}",
							},
						},
						Ports: []uint16{1234, 5678},
						Ips:   []string{},
					},
					HttpEndpoints: []*voc.HttpEndpoint{},
				},
				&voc.LoadBalancer{
					NetworkService: &voc.NetworkService{
						Networking: &voc.Networking{
							Resource: &voc.Resource{
								ID:        "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.network/applicationgateways/appgw1",
								ServiceID: testdata.MockCloudServiceID1,
								Name:      "appgw1",
								GeoLocation: voc.GeoLocation{
									Region: "eastus",
								},
								Type:   voc.LoadBalancerType,
								Labels: map[string]string{},
								Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
								Raw:    "{\"*armnetwork.ApplicationGateway\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/applicationGateways/appgw1\",\"location\":\"eastus\",\"name\":\"appgw1\",\"properties\":{\"webApplicationFirewallConfiguration\":{\"enabled\":true}}}]}",
							},
						},
					},
					AccessRestriction: voc.WebApplicationFirewall{
						Enabled: true,
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureNetworkDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}
			gotList, err := d.List()

			assert.Equal(t, len(tt.wantList), len(gotList))
			if !tt.wantErr(t, err) {
				return
			}

			for i := 0; i < len(tt.wantList); i++ {
				assert.Equal(t, tt.wantList[i], gotList[i])
			}
		})
	}
}

func TestNewAzureNetworkDiscovery(t *testing.T) {
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
			want: &azureNetworkDiscovery{
				&azureDiscovery{
					discovererComponent: NetworkComponent,
					csID:                discovery.DefaultCloudServiceID,
					backupMap:           make(map[string]*backup),
				},
			},
		},
		{
			name: "With sender",
			args: args{
				opts: []DiscoveryOption{WithSender(mockNetworkSender{})},
			},
			want: &azureNetworkDiscovery{
				&azureDiscovery{
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockNetworkSender{},
						},
					},
					discovererComponent: NetworkComponent,
					csID:                discovery.DefaultCloudServiceID,
					backupMap:           make(map[string]*backup),
				},
			},
		},
		{
			name: "With authorizer",
			args: args{
				opts: []DiscoveryOption{WithAuthorizer(&mockAuthorizer{})},
			},
			want: &azureNetworkDiscovery{
				&azureDiscovery{
					cred:                &mockAuthorizer{},
					discovererComponent: NetworkComponent,
					csID:                discovery.DefaultCloudServiceID,
					backupMap:           make(map[string]*backup),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewAzureNetworkDiscovery(tt.args.opts...)
			assert.Equal(t, tt.want, d)
			assert.Equal(t, "Azure Network", d.Name())
		})

	}
}

func Test_azureNetworkDiscovery_discoverNetworkInterfaces(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureNetworkDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}
			got, err := d.discoverNetworkInterfaces()
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equalf(t, tt.want, got, "discoverNetworkInterfaces()")
		})
	}
}

func Test_azureNetworkDiscovery_discoverLoadBalancer(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureNetworkDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}
			got, err := d.discoverLoadBalancer()
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equalf(t, tt.want, got, "discoverLoadBalancer()")
		})
	}
}

func Test_publicIPAddressFromLoadBalancer(t *testing.T) {
	id := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb3"
	name := "lb3"
	location := "eastus"
	publicIPID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/publicIPAddresses/test-b9cb3645-25d0-4288-910a-020563f63b1c"
	publicIPName := "mockPublicName"
	publicIPAddress := "111.222.333.444"
	emptyString := ""

	lbComplete := &armnetwork.LoadBalancer{
		ID:       &id,
		Name:     &name,
		Location: &location,
		Properties: &armnetwork.LoadBalancerPropertiesFormat{
			FrontendIPConfigurations: []*armnetwork.FrontendIPConfiguration{
				{
					Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
						PublicIPAddress: &armnetwork.PublicIPAddress{
							ID:   &publicIPID,
							Name: &publicIPName,
							Properties: &armnetwork.PublicIPAddressPropertiesFormat{
								IPAddress: &publicIPAddress,
							},
						},
					},
				},
			},
		},
	}
	lbWithoutIPAddress := &armnetwork.LoadBalancer{
		ID:       &id,
		Name:     &name,
		Location: &location,
		Properties: &armnetwork.LoadBalancerPropertiesFormat{
			FrontendIPConfigurations: []*armnetwork.FrontendIPConfiguration{
				{
					Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
						PublicIPAddress: &armnetwork.PublicIPAddress{
							Properties: &armnetwork.PublicIPAddressPropertiesFormat{
								IPAddress: nil,
							},
						},
					},
				},
			},
		},
	}
	lbWithEmptyIPAddress := &armnetwork.LoadBalancer{
		ID:       &id,
		Name:     &name,
		Location: &location,
		Properties: &armnetwork.LoadBalancerPropertiesFormat{
			FrontendIPConfigurations: []*armnetwork.FrontendIPConfiguration{
				{
					Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
						PublicIPAddress: &armnetwork.PublicIPAddress{
							Properties: &armnetwork.PublicIPAddressPropertiesFormat{
								IPAddress: &emptyString,
							},
						},
					},
				},
			},
		},
	}
	lbWithoutPublicIPAddress := &armnetwork.LoadBalancer{
		ID:       &id,
		Name:     &name,
		Location: &location,
		Properties: &armnetwork.LoadBalancerPropertiesFormat{
			FrontendIPConfigurations: []*armnetwork.FrontendIPConfiguration{
				{
					Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
						PublicIPAddress: nil,
					},
				},
			},
		},
	}

	type args struct {
		lb *armnetwork.LoadBalancer
	}
	tests := []struct {
		name string
		args args
		want []string
	}{

		{
			name: "Empty input",
			args: args{
				lb: nil,
			},
			want: []string{},
		},
		{
			name: "Empty FrontendIPConfiguration",
			args: args{
				lb: &armnetwork.LoadBalancer{
					ID:       &id,
					Name:     &name,
					Location: &location,
					Properties: &armnetwork.LoadBalancerPropertiesFormat{
						FrontendIPConfigurations: nil,
					},
				},
			},
			want: []string{},
		},
		{
			name: "Empty PublicIPAddress",
			args: args{
				lb: lbWithoutPublicIPAddress,
			},
			want: []string{},
		},
		{
			name: "Empty IPAddress (== nil)",
			args: args{
				lb: lbWithoutIPAddress,
			},
			want: []string{},
		},
		{
			name: "Empty IPAddress string",
			args: args{
				lb: lbWithEmptyIPAddress,
			},
			want: []string{},
		},
		{
			name: "Correct IP",
			args: args{
				lb: lbComplete,
			},
			want: []string{publicIPAddress},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, publicIPAddressFromLoadBalancer(tt.args.lb))
		})
	}
}

func Test_azureNetworkDiscovery_discoverApplicationGateway(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureNetworkDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}
			got, err := d.discoverApplicationGateway()
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equalf(t, tt.want, got, "discoverApplicationGateway()")
		})
	}
}

func Test_nsgFirewallEnabled(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		ni *armnetwork.Interface
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Empty input",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockNetworkSender()),
			},
			args: args{},
			want: false,
		}, {
			name: "Error getting nsg",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockNetworkSender()),
			},
			args: args{
				ni: &armnetwork.Interface{
					Properties: &armnetwork.InterfacePropertiesFormat{
						NetworkSecurityGroup: &armnetwork.SecurityGroup{
							ID: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkInterfaces/false"),
						},
					},
				},
			},
			want: false,
		},
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockNetworkSender()),
			},
			args: args{
				ni: &armnetwork.Interface{
					Properties: &armnetwork.InterfacePropertiesFormat{
						NetworkSecurityGroup: &armnetwork.SecurityGroup{
							ID: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkInterfaces/nsg1"),
						},
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		d := &azureNetworkDiscovery{
			azureDiscovery: tt.fields.azureDiscovery,
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := d.nsgFirewallEnabled(tt.args.ni); got != tt.want {
				t.Errorf("nsgFirewallEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getIDName(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Missing input",
			args: args{},
			want: "",
		},
		{
			name: "Happy path",
			args: args{
				id: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkInterfaces/nsg1",
			},
			want: "nsg1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getName(tt.args.id); got != tt.want {
				t.Errorf("getIDName() = %v, want %v", got, tt.want)
			}
		})
	}
}
