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
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
)

// type azureDiscovery struct {
// 	*azureDiscovery
// }

// func NewazureDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
// 	d := &azureDiscovery{
// 		&azureDiscovery{
// 			discovererComponent: NetworkComponent,
// 			csID:                discovery.DefaultCloudServiceID,
// 			backupMap:           make(map[string]*backup),
// 		},
// 	}

// 	// Apply options
// 	for _, opt := range opts {
// 		opt(d)
// 	}

// 	return d
// }

// func (*azureDiscovery) Name() string {
// 	return "Azure Network"
// }

// func (*azureDiscovery) Description() string {
// 	return "Discovery Azure network resources."
// }

// // List network resources
// func (d *azureDiscovery) List() (list []voc.IsCloudResource, err error) {
// 	// if err = d.authorize(); err != nil {
// 	// 	return nil, fmt.Errorf("%s: %w", ErrCouldNotAuthenticate, err)
// 	// }

// 	log.Info("Discover Azure network resources")

// 	// Discover network interfaces
// 	networkInterfaces, err := d.discoverNetworkInterfaces()
// 	if err != nil {
// 		return nil, fmt.Errorf("could not discover network interfaces: %w", err)
// 	}
// 	list = append(list, networkInterfaces...)

// 	// Discover Load Balancer
// 	loadBalancer, err := d.discoverLoadBalancer()
// 	if err != nil {
// 		return list, fmt.Errorf("could not discover load balancer: %w", err)
// 	}
// 	list = append(list, loadBalancer...)

// 	// Discover Application Gateway
// 	ag, err := d.discoverApplicationGateway()
// 	if err != nil {
// 		return list, fmt.Errorf("could not discover application gateways: %w", err)
// 	}
// 	list = append(list, ag...)

// 	return
// }

// discoverNetworkInterfaces discovers network interfaces
func (d *azureDiscovery) discoverNetworkInterfaces() ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	// initialize network interfaces client
	if err := d.initNetworkInterfacesClient(); err != nil {
		return nil, err
	}

	// List all network interfaces
	err := listPager(d,
		d.clients.networkInterfacesClient.NewListAllPager,
		d.clients.networkInterfacesClient.NewListPager,
		func(res armnetwork.InterfacesClientListAllResponse) []*armnetwork.Interface {
			return res.Value
		},
		func(res armnetwork.InterfacesClientListResponse) []*armnetwork.Interface {
			return res.Value
		},
		func(ni *armnetwork.Interface) error {
			s := d.handleNetworkInterfaces(ni)

			log.Infof("Adding network interface '%s'", s.GetName())

			list = append(list, s)

			return nil
		})
	if err != nil {
		return nil, err
	}

	return list, nil
}

// discoverApplicationGateway discovers application gateways
func (d *azureDiscovery) discoverApplicationGateway() ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	// initialize application gateway client
	if err := d.initApplicationGatewayClient(); err != nil {
		return nil, err
	}

	// List all application gateways
	err := listPager(d,
		d.clients.applicationGatewayClient.NewListAllPager,
		d.clients.applicationGatewayClient.NewListPager,
		func(res armnetwork.ApplicationGatewaysClientListAllResponse) []*armnetwork.ApplicationGateway {
			return res.Value
		},
		func(res armnetwork.ApplicationGatewaysClientListResponse) []*armnetwork.ApplicationGateway {
			return res.Value
		},
		func(ags *armnetwork.ApplicationGateway) error {
			s := d.handleApplicationGateway(ags)

			log.Infof("Adding application gateway %+v", s)

			list = append(list, s)

			return nil
		})
	if err != nil {
		return nil, err
	}

	return list, nil
}

// discoverLoadBalancer discovers load balancer
func (d *azureDiscovery) discoverLoadBalancer() ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	// initialize load balancers client
	if err := d.initLoadBalancersClient(); err != nil {
		return nil, err
	}

	// List all load balancers
	err := listPager(d,
		d.clients.loadBalancerClient.NewListAllPager,
		d.clients.loadBalancerClient.NewListPager,
		func(res armnetwork.LoadBalancersClientListAllResponse) []*armnetwork.LoadBalancer {
			return res.Value
		},
		func(res armnetwork.LoadBalancersClientListResponse) []*armnetwork.LoadBalancer {
			return res.Value
		},
		func(lbs *armnetwork.LoadBalancer) error {
			s := d.handleLoadBalancer(lbs)

			log.Infof("Adding load balancer %+v", s)

			list = append(list, s)

			return nil
		})
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (d *azureDiscovery) handleLoadBalancer(lb *armnetwork.LoadBalancer) voc.IsNetwork {
	return &voc.LoadBalancer{
		NetworkService: &voc.NetworkService{
			Networking: &voc.Networking{
				Resource: discovery.NewResource(d,
					voc.ResourceID(util.Deref(lb.ID)),
					util.Deref(lb.Name),
					// No creation time available
					nil,
					voc.GeoLocation{
						Region: util.Deref(lb.Location),
					},
					labels(lb.Tags),
					resourceGroupID(lb.ID),
					voc.LoadBalancerType,
					lb,
				),
			},
			Ips:   publicIPAddressFromLoadBalancer(lb),
			Ports: LoadBalancerPorts(lb),
		},
		// TODO(all): do we need the httpEndpoint for load balancers?
		HttpEndpoints: []*voc.HttpEndpoint{},
	}
}

// handleApplicationGateway returns the application gateway with its properties
// NOTE: handleApplicationGateway uses the LoadBalancer for now until there is a own resource
func (d *azureDiscovery) handleApplicationGateway(ag *armnetwork.ApplicationGateway) voc.IsNetwork {
	return &voc.LoadBalancer{
		NetworkService: &voc.NetworkService{
			Networking: &voc.Networking{
				Resource: discovery.NewResource(
					d,
					voc.ResourceID(util.Deref(ag.ID)),
					util.Deref(ag.Name),
					nil,
					voc.GeoLocation{Region: util.Deref(ag.Location)},
					labels(ag.Tags),
					resourceGroupID(ag.ID),
					voc.LoadBalancerType,
					ag,
				),
			},
		},
		AccessRestriction: voc.WebApplicationFirewall{
			Enabled: util.Deref(ag.Properties.WebApplicationFirewallConfiguration.Enabled),
		},
	}
}

func (d *azureDiscovery) handleNetworkInterfaces(ni *armnetwork.Interface) voc.IsNetwork {
	return &voc.NetworkInterface{
		Networking: &voc.Networking{
			Resource: discovery.NewResource(d,
				voc.ResourceID(util.Deref(ni.ID)),
				util.Deref(ni.Name),
				// No creation time available
				nil,
				voc.GeoLocation{
					Region: util.Deref(ni.Location),
				},
				labels(ni.Tags),
				resourceGroupID(ni.ID),
				voc.NetworkInterfaceType,
				ni,
			),
		},
		AccessRestriction: &voc.L3Firewall{
			Enabled: d.nsgFirewallEnabled(ni),
			// Inbound: ,
			// RestrictedPorts: ,
		},
	}
}

// nsgFirewallEnabled checks if network security group (NSG) rules are configured. A NSG is a firewall that operates at OSI layers 3 and 4 to filter ingress and egress traffic. (https://learn.microsoft.com/en-us/azure/firewall/firewall-faq#what-is-the-difference-between-network-security-groups--nsgs--and-azure-firewall, Last access: 05/02/2023)
func (d *azureDiscovery) nsgFirewallEnabled(ni *armnetwork.Interface) bool {
	// initialize network interfaces client
	if err := d.initNetworkSecurityGroupClient(); err != nil {
		log.Error(err)
		return false
	}

	if ni != nil && ni.Properties != nil && ni.Properties.NetworkSecurityGroup != nil {
		vmNsg := ni.Properties.NetworkSecurityGroup
		nsg, err := d.clients.networkSecurityGroupsClient.Get(context.Background(), resourceGroupName(*vmNsg.ID), getName(*vmNsg.ID), &armnetwork.SecurityGroupsClientGetOptions{})
		if err != nil {
			log.Errorf("error getting network security group: %v", err)
			return false
		}

		// TODO(all): We have to check more than len(securityRules) > 0. But what is a good check?
		if len(nsg.SecurityGroup.Properties.SecurityRules) > 0 {
			return true
		}

	}

	return false
}

func LoadBalancerPorts(lb *armnetwork.LoadBalancer) (loadBalancerPorts []uint16) {

	for _, item := range lb.Properties.LoadBalancingRules {
		loadBalancerPorts = append(loadBalancerPorts, uint16(util.Deref(item.Properties.FrontendPort)))
	}

	return loadBalancerPorts
}

// // Returns all restricted ports for the network interface
// func (d *azureDiscovery) getRestrictedPorts(ni *network.Interface) string {
//
//     var restrictedPorts []string
//
//     if ni.InterfacePropertiesFormat == nil ||
//             ni.InterfacePropertiesFormat.NetworkSecurityGroup == nil ||
//             ni.InterfacePropertiesFormat.NetworkSecurityGroup.ID == nil {
//             return ""
//     }
//
//	nsgID := util.Deref(ni.NetworkSecurityGroup.ID)
//
//	client := network.NewSecurityGroupsClient(util.Deref(d.sub.SubscriptionID))
//
//     // Get the Security Group of the network interface ni
//     sg, err := client.Get(context.Background(), getResourceGroupName(nsgID), strings.Split(nsgID, "/")[8], "")
//
//     if err != nil {
//             log.Errorf("Could not get security group: %v", err)
//             return ""
//     }
//
//     if sg.SecurityGroupPropertiesFormat != nil && sg.SecurityGroupPropertiesFormat.SecurityRules != nil {
//             // Find all ports defined in the security rules with access property "Deny"
//             for _, securityRule := range *sg.SecurityRules {
//                     if securityRule.Access == network.SecurityRuleAccessDeny {
//                             restrictedPorts = append(restrictedPorts, *securityRule.SourcePortRange)
//                     }
//             }
//     }
//
//     restrictedPortsClean := deleteDuplicatesFromSlice(restrictedPorts)
//
//     return strings.Join(restrictedPortsClean, ",")
// }
//
// func deleteDuplicatesFromSlice(intSlice []string) []string {
//     keys := make(map[string]bool)
//     var list []string
//     for _, entry := range intSlice {
//             if _, value := keys[entry]; !value {
//                     keys[entry] = true
//                     list = append(list, entry)
//             }
//     }
//     return list
// }

func publicIPAddressFromLoadBalancer(lb *armnetwork.LoadBalancer) []string {
	var publicIPAddresses = []string{}

	if lb == nil || lb.Properties == nil || lb.Properties.FrontendIPConfigurations == nil {
		return []string{}
	}

	fIpConfig := lb.Properties.FrontendIPConfigurations
	for i := range fIpConfig {

		if fIpConfig[i].Properties.PublicIPAddress == nil || fIpConfig[i].Properties.PublicIPAddress.Properties == nil || fIpConfig[i].Properties.PublicIPAddress.Properties.IPAddress == nil {
			continue
		}

		// Get public IP address
		ipAddress := util.Deref(fIpConfig[i].Properties.PublicIPAddress.Properties.IPAddress)
		if ipAddress == "" {
			log.Infof("No public IP adress available.")
			continue
		}

		publicIPAddresses = append(publicIPAddresses, ipAddress)
	}

	return publicIPAddresses
}

// getName returns the name of a given Azure ID
func getName(id string) string {
	if id == "" {
		return ""
	}
	return strings.Split(id, "/")[8]
}
