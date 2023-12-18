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
	"reflect"
	"testing"

	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
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

	return m.mockSender.Do(req)
}

// TODO(anatheka): Add more tests
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
			d := tt.fields.azureDiscovery

			got, err := d.discoverNetworkInterfaces()
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equalf(t, tt.want, got, "discoverNetworkInterfaces()")
		})
	}
}

// TODO(anatheka): Add more tests
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

// TODO(anatheka): Add more tests
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
			// d := &azureNetworkDiscovery{
			// 	azureDiscovery: tt.fields.azureDiscovery,
			// }
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

// TODO(anatheka): Add tests
func Test_azureDiscovery_handleLoadBalancer(t *testing.T) {
	type fields struct {
		isAuthorized        bool
		sub                 *armsubscription.Subscription
		cred                azcore.TokenCredential
		rg                  *string
		clientOptions       arm.ClientOptions
		discovererComponent string
		clients             clients
		csID                string
		backupMap           map[string]*backup
		defenderProperties  map[string]*defenderProperties
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureDiscovery{
				isAuthorized:        tt.fields.isAuthorized,
				sub:                 tt.fields.sub,
				cred:                tt.fields.cred,
				rg:                  tt.fields.rg,
				clientOptions:       tt.fields.clientOptions,
				discovererComponent: tt.fields.discovererComponent,
				clients:             tt.fields.clients,
				csID:                tt.fields.csID,
				backupMap:           tt.fields.backupMap,
				defenderProperties:  tt.fields.defenderProperties,
			}
			if got := d.handleLoadBalancer(tt.args.lb); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("azureDiscovery.handleLoadBalancer() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO(anatheka): Add tests
func Test_azureDiscovery_handleApplicationGateway(t *testing.T) {
	type fields struct {
		isAuthorized        bool
		sub                 *armsubscription.Subscription
		cred                azcore.TokenCredential
		rg                  *string
		clientOptions       arm.ClientOptions
		discovererComponent string
		clients             clients
		csID                string
		backupMap           map[string]*backup
		defenderProperties  map[string]*defenderProperties
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureDiscovery{
				isAuthorized:        tt.fields.isAuthorized,
				sub:                 tt.fields.sub,
				cred:                tt.fields.cred,
				rg:                  tt.fields.rg,
				clientOptions:       tt.fields.clientOptions,
				discovererComponent: tt.fields.discovererComponent,
				clients:             tt.fields.clients,
				csID:                tt.fields.csID,
				backupMap:           tt.fields.backupMap,
				defenderProperties:  tt.fields.defenderProperties,
			}
			if got := d.handleApplicationGateway(tt.args.ag); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("azureDiscovery.handleApplicationGateway() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO(anatheka): Add tests
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

// TODO(anatheka): Add tests
func Test_azureDiscovery_handleNetworkInterfaces(t *testing.T) {
	type fields struct {
		isAuthorized        bool
		sub                 *armsubscription.Subscription
		cred                azcore.TokenCredential
		rg                  *string
		clientOptions       arm.ClientOptions
		discovererComponent string
		clients             clients
		csID                string
		backupMap           map[string]*backup
		defenderProperties  map[string]*defenderProperties
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureDiscovery{
				isAuthorized:        tt.fields.isAuthorized,
				sub:                 tt.fields.sub,
				cred:                tt.fields.cred,
				rg:                  tt.fields.rg,
				clientOptions:       tt.fields.clientOptions,
				discovererComponent: tt.fields.discovererComponent,
				clients:             tt.fields.clients,
				csID:                tt.fields.csID,
				backupMap:           tt.fields.backupMap,
				defenderProperties:  tt.fields.defenderProperties,
			}
			if got := d.handleNetworkInterfaces(tt.args.ni); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("azureDiscovery.handleNetworkInterfaces() = %v, want %v", got, tt.want)
			}
		})
	}
}
