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
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
)

type azureNetworkDiscovery struct {
	*azureDiscovery
}

func NewAzureNetworkDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &azureNetworkDiscovery{
		&azureDiscovery{
			discovererComponent: NetworkComponent,
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(d.azureDiscovery)
	}

	return d
}

func (*azureNetworkDiscovery) Name() string {
	return "Azure Network"
}

func (*azureNetworkDiscovery) Description() string {
	return "Discovery Azure network resources."
}

// List network resources
func (d *azureNetworkDiscovery) List() (list []voc.IsCloudResource, err error) {
	if err = d.authorize(); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrCouldNotAuthenticate, err)
	}

	log.Info("Discover Azure network resources")

	// Discover network interfaces
	networkInterfaces, err := d.discoverNetworkInterfaces()
	if err != nil {
		return nil, fmt.Errorf("could not discover network interfaces: %w", err)
	}
	list = append(list, networkInterfaces...)

	// Discover Load Balancer
	loadBalancer, err := d.discoverLoadBalancer()
	if err != nil {
		return list, fmt.Errorf("could not discover load balancer: %w", err)
	}
	list = append(list, loadBalancer...)

	return
}

// Discover network interfaces
func (d *azureNetworkDiscovery) discoverNetworkInterfaces() ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	// initialize network interfaces client
	if err := d.initNetworkInterfacesClient(); err != nil {
		return nil, err
	}

	// List all network interfaces accross all resource groups
	listPager := d.clients.networkInterfacesClient.NewListAllPager(&armnetwork.InterfacesClientListAllOptions{})
	ni := make([]*armnetwork.Interface, 0)
	for listPager.More() {
		pageResponse, err := listPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("%s: %v", ErrGettingNextPage, err)
			return nil, err
		}
		ni = append(ni, pageResponse.Value...)
	}

	for i := range ni {
		s := d.handleNetworkInterfaces(ni[i])

		log.Infof("Adding network interface '%s'", s.GetName())

		list = append(list, s)
	}

	return list, nil
}

// Discover load balancer
func (d *azureNetworkDiscovery) discoverLoadBalancer() ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	// initialize load balancers client
	if err := d.initLoadBalancersClient(); err != nil {
		return nil, err
	}

	// List all load balancers accross all resource groups
	listPager := d.clients.loadBalancerClient.NewListAllPager(&armnetwork.LoadBalancersClientListAllOptions{})
	lbs := make([]*armnetwork.LoadBalancer, 0)
	for listPager.More() {
		pageResponse, err := listPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("%s: %v", ErrGettingNextPage, err)
			return nil, err
		}
		lbs = append(lbs, pageResponse.Value...)
	}

	for i := range lbs {
		s := d.handleLoadBalancer(lbs[i])

		log.Infof("Adding load balancer %+v", s)

		list = append(list, s)
	}

	return list, nil
}

func (d *azureNetworkDiscovery) handleLoadBalancer(lb *armnetwork.LoadBalancer) voc.IsNetwork {
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
					voc.LoadBalancerType,
				),
			},
			Ips:   publicIPAddressFromLoadBalancer(lb),
			Ports: LoadBalancerPorts(lb),
		},
		// TODO(all): do we need the AccessRestriction for load balancers?
		AccessRestrictions: &[]voc.AccessRestriction{},
		// TODO(all): do we need the httpEndpoint for load balancers?
		HttpEndpoints: &[]voc.HttpEndpoint{},
	}
}

func (d *azureNetworkDiscovery) handleNetworkInterfaces(ni *armnetwork.Interface) voc.IsNetwork {
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
				voc.NetworkInterfaceType,
			),
		},

		// AccessRestriction: &voc.AccessRestriction{
		// 	Inbound:         false, // TODO(garuppel): TBD
		// 	RestrictedPorts: d.getRestrictedPorts(ni),
		// },
	}
}

func LoadBalancerPorts(lb *armnetwork.LoadBalancer) (loadBalancerPorts []uint16) {

	for _, item := range lb.Properties.LoadBalancingRules {
		loadBalancerPorts = append(loadBalancerPorts, uint16(util.Deref(item.Properties.FrontendPort)))
	}

	return loadBalancerPorts
}

// // Returns all restricted ports for the network interface
// func (d *azureNetworkDiscovery) getRestrictedPorts(ni *network.Interface) string {
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

// initNetworkInterfacesClient creates the client if not already exists
func (d *azureNetworkDiscovery) initNetworkInterfacesClient() (err error) {
	d.clients.networkInterfacesClient, err = initClient(d.clients.networkInterfacesClient, d.azureDiscovery, armnetwork.NewInterfacesClient)
	return
}

// initLoadBalancersClient creates the client if not already exists
func (d *azureNetworkDiscovery) initLoadBalancersClient() (err error) {
	d.clients.loadBalancerClient, err = initClient(d.clients.loadBalancerClient, d.azureDiscovery, armnetwork.NewLoadBalancersClient)
	return
}
