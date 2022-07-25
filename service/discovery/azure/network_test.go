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

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"

	"github.com/stretchr/testify/assert"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
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
	assert.ErrorIs(t, err, ErrNoCredentialsConfigured)
}

func Test_azureNetworkDiscovery_List(t *testing.T) {
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
							Transport: &mockStorageSender{},
						},
					},
				},
			},
			wantList: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not discover network interfaces:")
			},
		},
		{
			name: "Without errors",
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
			wantList: []voc.IsCloudResource{
				&voc.NetworkInterface{
					Networking: &voc.Networking{
						Resource: &voc.Resource{
							ID:   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkInterfaces/iface1",
							Name: "iface1",
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Type:   []string{"NetworkInterface", "Compute", "Resource"},
							Labels: map[string]string{},
						},
					},
					AccessRestriction: nil,
				},
				&voc.LoadBalancer{
					NetworkService: &voc.NetworkService{
						Networking: &voc.Networking{
							Resource: &voc.Resource{
								ID:   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb1",
								Name: "lb1",
								GeoLocation: voc.GeoLocation{
									Region: "eastus",
								},
								Type:   []string{"LoadBalancer", "NetworkService", "Resource"},
								Labels: map[string]string{},
							},
						},
						Ips:   []string{"111.222.333.444"},
						Ports: []uint16{1234, 5678},
					},
					AccessRestrictions: &[]voc.AccessRestriction{},
					HttpEndpoints:      &[]voc.HttpEndpoint{},
				},
				&voc.LoadBalancer{
					NetworkService: &voc.NetworkService{
						Networking: &voc.Networking{
							Resource: &voc.Resource{
								ID:   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb2",
								Name: "lb2",
								GeoLocation: voc.GeoLocation{
									Region: "eastus",
								},
								Type:   []string{"LoadBalancer", "NetworkService", "Resource"},
								Labels: map[string]string{},
							},
						},
						Ports: []uint16{1234, 5678},
					},
					AccessRestrictions: &[]voc.AccessRestriction{},
					HttpEndpoints:      &[]voc.HttpEndpoint{},
				},
				&voc.LoadBalancer{
					NetworkService: &voc.NetworkService{
						Networking: &voc.Networking{
							Resource: &voc.Resource{
								ID:   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb3",
								Name: "lb3",
								GeoLocation: voc.GeoLocation{
									Region: "eastus",
								},
								Type:   []string{"LoadBalancer", "NetworkService", "Resource"},
								Labels: map[string]string{},
							},
						},
						Ports: []uint16{1234, 5678},
					},
					AccessRestrictions: &[]voc.AccessRestriction{},
					HttpEndpoints:      &[]voc.HttpEndpoint{},
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
			want: &azureNetworkDiscovery{},
		},
		{
			name: "With sender",
			args: args{
				opts: []DiscoveryOption{WithSender(mockNetworkSender{})},
			},
			want: &azureNetworkDiscovery{
				azureDiscovery{
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockNetworkSender{},
						},
					},
				},
			},
		},
		{
			name: "With authorizer",
			args: args{
				opts: []DiscoveryOption{WithAuthorizer(&mockAuthorizer{})},
			},
			want: &azureNetworkDiscovery{
				azureDiscovery{
					cred: &mockAuthorizer{},
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

func Test_frontendPublicIPAddressName(t *testing.T) {
	type args struct {
		frontendPublicIPAddressID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty ip",
			args: args{
				frontendPublicIPAddressID: "",
			},
			want: "",
		},
		{
			name: "IP incorrect",
			args: args{
				frontendPublicIPAddressID: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/test-b9cb3645-25d0-4288-910a-020563f63b1c",
			},
			want: "",
		},
		{
			name: "IP correct",
			args: args{
				frontendPublicIPAddressID: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/publicIPAddresses/test-b9cb3645-25d0-4288-910a-020563f63b1c",
			},
			want: "test-b9cb3645-25d0-4288-910a-020563f63b1c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, frontendPublicIPAddressName(tt.args.frontendPublicIPAddressID), "frontendPublicIPAddressName(%v)", tt.args.frontendPublicIPAddressID)
		})
	}
}

func Test_azureNetworkDiscovery_publicIPAddressFromLoadBalancer(t *testing.T) {
	id := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb3"
	name := "lb3"
	location := "eastus"
	publicIPID := "mockPublicID"

	lbWithoutPublicIPAddressID := &armnetwork.LoadBalancer{
		ID:       &id,
		Name:     &name,
		Location: &location,
		Properties: &armnetwork.LoadBalancerPropertiesFormat{
			FrontendIPConfigurations: []*armnetwork.FrontendIPConfiguration{
				{
					Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
						PublicIPAddress: &armnetwork.PublicIPAddress{},
					},
				},
			},
		},
	}
	lbWithoutPublicIPAddressName := &armnetwork.LoadBalancer{
		ID:       &id,
		Name:     &name,
		Location: &location,
		Properties: &armnetwork.LoadBalancerPropertiesFormat{
			FrontendIPConfigurations: []*armnetwork.FrontendIPConfiguration{
				{
					Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
						PublicIPAddress: &armnetwork.PublicIPAddress{
							ID: &publicIPID,
						},
					},
				},
			},
		},
	}

	type fields struct {
		azureDiscovery azureDiscovery
	}
	type args struct {
		lb *armnetwork.LoadBalancer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{

		{
			name: "Empty input",
			args: args{
				lb: nil,
			},
			want: []string{},
		},
		{
			name: "Empty FrontendIPConfiguation",
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
			fields: fields{
				azureDiscovery: azureDiscovery{
					cred: &mockAuthorizer{},
				},
			},
			want: []string(nil),
		},
		{
			name: "Empty PublicIPAdressID",
			args: args{
				lb: lbWithoutPublicIPAddressID,
			},
			fields: fields{
				azureDiscovery: azureDiscovery{
					cred: &mockAuthorizer{},
				},
			},
			want: []string(nil),
		},
		{
			name: "Empty PublicIPAdressName",
			args: args{
				lb: lbWithoutPublicIPAddressName,
			},
			fields: fields{
				azureDiscovery: azureDiscovery{
					cred: &mockAuthorizer{},
				},
			},
			want: []string(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureNetworkDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}
			assert.Equal(t, tt.want, d.publicIPAddressFromLoadBalancer(tt.args.lb))
		})
	}
}
