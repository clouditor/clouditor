//go:build exclude

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
	"reflect"
	"testing"

	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/stretchr/testify/assert"
)

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
		{
			name: "Happy path: with resource group",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender(), WithResourceGroup("res1")),
			},
			want: []voc.IsCloudResource{
				&voc.NetworkInterface{
					Networking: &voc.Networking{
						Resource: &voc.Resource{
							ID:        "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkInterfaces/iface1",
							Name:      "iface1",
							ServiceID: testdata.MockCloudServiceID1,
							Type:      voc.NetworkInterfaceType,
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Labels: map[string]string{},
							Parent: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1",
							Raw:    "{\"*armnetwork.Interface\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkInterfaces/iface1\",\"location\":\"eastus\",\"name\":\"iface1\",\"properties\":{\"networkSecurityGroup\":{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkSecurityGroups/nsg1\",\"location\":\"eastus\"}}}]}",
						},
					},
					AccessRestriction: &voc.L3Firewall{
						Enabled: true,
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			want: []voc.IsCloudResource{
				&voc.NetworkInterface{
					Networking: &voc.Networking{
						Resource: &voc.Resource{
							ID:        "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkInterfaces/iface1",
							Name:      "iface1",
							ServiceID: testdata.MockCloudServiceID1,
							Type:      voc.NetworkInterfaceType,
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Labels: map[string]string{},
							Parent: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1",
							Raw:    "{\"*armnetwork.Interface\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkInterfaces/iface1\",\"location\":\"eastus\",\"name\":\"iface1\",\"properties\":{\"networkSecurityGroup\":{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkSecurityGroups/nsg1\",\"location\":\"eastus\"}}}]}",
						},
					},
					AccessRestriction: &voc.L3Firewall{
						Enabled: true,
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

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
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			want: []voc.IsCloudResource{
				&voc.LoadBalancer{
					NetworkService: &voc.NetworkService{
						Networking: &voc.Networking{
							Resource: &voc.Resource{
								ID:        "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb1",
								Name:      "lb1",
								ServiceID: testdata.MockCloudServiceID1,
								Type:      voc.LoadBalancerType,
								GeoLocation: voc.GeoLocation{
									Region: "eastus",
								},
								Labels: map[string]string{},
								Raw:    "{\"*armnetwork.LoadBalancer\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb1\",\"location\":\"eastus\",\"name\":\"lb1\",\"properties\":{\"frontendIPConfigurations\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb1/frontendIPConfigurations/b9cb3645-25d0-4288-910a-020563f63b1c\",\"name\":\"b9cb3645-25d0-4288-910a-020563f63b1c\",\"properties\":{\"publicIPAddress\":{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/publicIPAddresses/test-b9cb3645-25d0-4288-910a-020563f63b1c\",\"properties\":{\"ipAddress\":\"111.222.333.444\"}}}}],\"loadBalancingRules\":[{\"properties\":{\"frontendPort\":1234}},{\"properties\":{\"frontendPort\":5678}}]}}]}",
								Parent: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1",
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
								ID:        "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb2",
								Name:      "lb2",
								ServiceID: testdata.MockCloudServiceID1,
								Type:      voc.LoadBalancerType,
								GeoLocation: voc.GeoLocation{
									Region: "eastus",
								},
								Labels: map[string]string{},
								Raw:    "{\"*armnetwork.LoadBalancer\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb2\",\"location\":\"eastus\",\"name\":\"lb2\",\"properties\":{\"frontendIPConfigurations\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb1/frontendIPConfigurations/b9cb3645-25d0-4288-910a-020563f63b1c\",\"name\":\"b9cb3645-25d0-4288-910a-020563f63b1c\",\"properties\":{}}],\"loadBalancingRules\":[{\"properties\":{\"frontendPort\":1234}},{\"properties\":{\"frontendPort\":5678}}]}}]}",
								Parent: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1",
							},
						},
						Ips:   []string{},
						Ports: []uint16{1234, 5678},
					},
					HttpEndpoints: []*voc.HttpEndpoint{},
				},
				&voc.LoadBalancer{
					NetworkService: &voc.NetworkService{
						Networking: &voc.Networking{
							Resource: &voc.Resource{
								ID:        "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb3",
								Name:      "lb3",
								ServiceID: testdata.MockCloudServiceID1,
								Type:      voc.LoadBalancerType,
								GeoLocation: voc.GeoLocation{
									Region: "eastus",
								},
								Labels: map[string]string{},
								Raw:    "{\"*armnetwork.LoadBalancer\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb3\",\"location\":\"eastus\",\"name\":\"lb3\",\"properties\":{\"frontendIPConfigurations\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb1/frontendIPConfigurations/b9cb3645-25d0-4288-910a-020563f63b1c\",\"name\":\"b9cb3645-25d0-4288-910a-020563f63b1c\",\"properties\":{\"publicIPAddress\":{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/publicIPAddresses/test-b9cb3645-25d0-4288-910a-020563f63b1d\"}}}],\"loadBalancingRules\":[{\"properties\":{\"frontendPort\":1234}},{\"properties\":{\"frontendPort\":5678}}]}}]}",
								Parent: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1",
							},
						},
						Ips:   []string{},
						Ports: []uint16{1234, 5678},
					},
					HttpEndpoints: []*voc.HttpEndpoint{},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

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
					Location: util.Ref("eastus"),
					Properties: &armnetwork.LoadBalancerPropertiesFormat{
						FrontendIPConfigurations: nil,
					},
				},
			},
			want: []string{},
		},
		{
			name: "Missing PublicIPAddress",
			args: args{
				lb: &armnetwork.LoadBalancer{
					ID:       &id,
					Name:     &name,
					Location: util.Ref("eastus"),
					Properties: &armnetwork.LoadBalancerPropertiesFormat{
						FrontendIPConfigurations: []*armnetwork.FrontendIPConfiguration{
							{
								Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
									PublicIPAddress: nil,
								},
							},
						},
					},
				},
			},
			want: []string{},
		},
		{
			name: "Empty IPAddress (== nil)",
			args: args{
				lb: &armnetwork.LoadBalancer{
					ID:       &id,
					Name:     &name,
					Location: util.Ref("eastus"),
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
				},
			},
			want: []string{},
		},
		{
			name: "Empty IPAddress",
			args: args{
				lb: &armnetwork.LoadBalancer{
					ID:       &id,
					Name:     &name,
					Location: util.Ref("eastus"),
					Properties: &armnetwork.LoadBalancerPropertiesFormat{
						FrontendIPConfigurations: []*armnetwork.FrontendIPConfiguration{
							{
								Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
									PublicIPAddress: &armnetwork.PublicIPAddress{
										Properties: &armnetwork.PublicIPAddressPropertiesFormat{
											IPAddress: util.Ref(""),
										},
									},
								},
							},
						},
					},
				},
			},
			want: []string{},
		},
		{
			name: "Correct IP",
			args: args{
				lb: &armnetwork.LoadBalancer{
					ID:       &id,
					Name:     &name,
					Location: util.Ref("eastus"),
					Properties: &armnetwork.LoadBalancerPropertiesFormat{
						FrontendIPConfigurations: []*armnetwork.FrontendIPConfiguration{
							{
								Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
									PublicIPAddress: &armnetwork.PublicIPAddress{
										ID:   util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/publicIPAddresses/test-b9cb3645-25d0-4288-910a-020563f63b1c"),
										Name: util.Ref("publicName"),
										Properties: &armnetwork.PublicIPAddressPropertiesFormat{
											IPAddress: util.Ref("111.222.333.444"),
										},
									},
								},
							},
						},
					},
				},
			},
			want: []string{"111.222.333.444"},
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
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			want: []voc.IsCloudResource{
				&voc.LoadBalancer{
					NetworkService: &voc.NetworkService{
						Networking: &voc.Networking{
							Resource: &voc.Resource{
								ID:        "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/applicationGateways/appgw1",
								Name:      "appgw1",
								ServiceID: testdata.MockCloudServiceID1,
								Type:      voc.LoadBalancerType,
								GeoLocation: voc.GeoLocation{
									Region: "eastus",
								},
								Labels: map[string]string{},
								Raw:    "{\"*armnetwork.ApplicationGateway\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/applicationGateways/appgw1\",\"location\":\"eastus\",\"name\":\"appgw1\",\"properties\":{\"webApplicationFirewallConfiguration\":{\"enabled\":true}}}]}",
								Parent: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1",
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
			d := tt.fields.azureDiscovery

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
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{},
			want: false,
		}, {
			name: "Error getting nsg",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
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
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
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
		d := tt.fields.azureDiscovery

		t.Run(tt.name, func(t *testing.T) {
			if got := d.nsgFirewallEnabled(tt.args.ni); got != tt.want {
				t.Errorf("nsgFirewallEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getName(t *testing.T) {
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

func Test_azureDiscovery_handleLoadBalancer(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		lb *armnetwork.LoadBalancer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   voc.IsNetwork
	}{
		{
			name: "Happy path",
			fields: fields{
				NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				lb: &armnetwork.LoadBalancer{
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb1"),
					Name:     util.Ref("lb1"),
					Location: util.Ref("eastus"),
					Tags: map[string]*string{
						"tag1": util.Ref("value1"),
						"tag2": util.Ref("value2"),
					},
					Properties: &armnetwork.LoadBalancerPropertiesFormat{
						LoadBalancingRules: []*armnetwork.LoadBalancingRule{},
					},
				},
			},
			want: &voc.LoadBalancer{
				NetworkService: &voc.NetworkService{
					Networking: &voc.Networking{
						Resource: &voc.Resource{
							ID:        "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb1",
							Name:      "lb1",
							ServiceID: testdata.MockCloudServiceID1,
							Type:      voc.LoadBalancerType,
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Labels: map[string]string{
								"tag1": "value1",
								"tag2": "value2",
							},
							Raw:    "{\"*armnetwork.LoadBalancer\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/loadBalancers/lb1\",\"location\":\"eastus\",\"name\":\"lb1\",\"properties\":{\"loadBalancingRules\":[]},\"tags\":{\"tag1\":\"value1\",\"tag2\":\"value2\"}}]}",
							Parent: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1",
						},
					},
					Ips:   []string{},
					Ports: nil,
				},
				HttpEndpoints: []*voc.HttpEndpoint{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery
			got := d.handleLoadBalancer(tt.args.lb)

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_azureDiscovery_handleApplicationGateway(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		ag *armnetwork.ApplicationGateway
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   voc.IsNetwork
	}{
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				ag: &armnetwork.ApplicationGateway{
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/applicationGateways/appgw1"),
					Name:     util.Ref("appgw1"),
					Location: util.Ref("eastus"),
					Properties: &armnetwork.ApplicationGatewayPropertiesFormat{
						WebApplicationFirewallConfiguration: &armnetwork.ApplicationGatewayWebApplicationFirewallConfiguration{
							Enabled: util.Ref(true),
						},
					},
				},
			},
			want: &voc.LoadBalancer{
				NetworkService: &voc.NetworkService{
					Networking: &voc.Networking{
						Resource: &voc.Resource{
							ID:        "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/applicationGateways/appgw1",
							Name:      "appgw1",
							ServiceID: testdata.MockCloudServiceID1,
							GeoLocation: voc.GeoLocation{
								Region: "eastus",
							},
							Type:   voc.LoadBalancerType,
							Labels: map[string]string{},
							Raw:    "{\"*armnetwork.ApplicationGateway\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/applicationGateways/appgw1\",\"location\":\"eastus\",\"name\":\"appgw1\",\"properties\":{\"webApplicationFirewallConfiguration\":{\"enabled\":true}}}]}",
							Parent: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1",
						},
					},
				},
				AccessRestriction: voc.WebApplicationFirewall{
					Enabled: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery
			got := d.handleApplicationGateway(tt.args.ag)

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_loadBalancerPorts(t *testing.T) {
	type args struct {
		lb *armnetwork.LoadBalancer
	}
	tests := []struct {
		name                  string
		args                  args
		wantLoadBalancerPorts []uint16
	}{
		{
			name: "Happy path: empty input",
			args: args{
				lb: &armnetwork.LoadBalancer{
					Properties: &armnetwork.LoadBalancerPropertiesFormat{
						LoadBalancingRules: []*armnetwork.LoadBalancingRule{},
					},
				},
			},
			wantLoadBalancerPorts: nil,
		},
		{
			name: "Happy path",
			args: args{
				lb: &armnetwork.LoadBalancer{
					Properties: &armnetwork.LoadBalancerPropertiesFormat{
						LoadBalancingRules: []*armnetwork.LoadBalancingRule{
							{
								Properties: &armnetwork.LoadBalancingRulePropertiesFormat{
									FrontendPort: util.Ref(int32(99)),
								},
							},
						},
					},
				},
			},
			wantLoadBalancerPorts: []uint16{99},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLoadBalancerPorts := loadBalancerPorts(tt.args.lb); !reflect.DeepEqual(gotLoadBalancerPorts, tt.wantLoadBalancerPorts) {
				t.Errorf("LoadBalancerPorts() = %v, want %v", gotLoadBalancerPorts, tt.wantLoadBalancerPorts)
			}
		})
	}
}

func Test_azureDiscovery_handleNetworkInterfaces(t *testing.T) {
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
		want   voc.IsNetwork
	}{
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				ni: &armnetwork.Interface{
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkInterfaces/iface1"),
					Name:     util.Ref("iface1"),
					Location: util.Ref("eastus"),
					Properties: &armnetwork.InterfacePropertiesFormat{
						NetworkSecurityGroup: &armnetwork.SecurityGroup{
							ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkSecurityGroups/nsg1"),
							Location: util.Ref("eastus"),
						},
					},
				},
			},
			want: &voc.NetworkInterface{
				Networking: &voc.Networking{
					Resource: &voc.Resource{
						ID:        "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkInterfaces/iface1",
						Name:      "iface1",
						ServiceID: testdata.MockCloudServiceID1,
						Type:      voc.NetworkInterfaceType,
						GeoLocation: voc.GeoLocation{
							Region: "eastus",
						},
						Labels: map[string]string{},
						Raw:    "{\"*armnetwork.Interface\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkInterfaces/iface1\",\"location\":\"eastus\",\"name\":\"iface1\",\"properties\":{\"networkSecurityGroup\":{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkSecurityGroups/nsg1\",\"location\":\"eastus\"}}}]}",
						Parent: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1",
					},
				},
				AccessRestriction: &voc.L3Firewall{
					Enabled: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery
			got := d.handleNetworkInterfaces(tt.args.ni)

			assert.Equal(t, tt.want, got)
		})
	}
}
