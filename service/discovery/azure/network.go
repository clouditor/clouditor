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
	"strings"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/go-autorest/autorest/to"
)

type azureNetworkDiscovery struct {
	azureDiscovery
}

func NewAzureNetworkDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &azureNetworkDiscovery{}

	for _, opt := range opts {
		if auth, ok := opt.(*authorizerOption); ok {
			d.authOption = auth
		} else {
			d.options = append(d.options, opt)
		}
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
		return nil, fmt.Errorf("could not authorize Azure account: %w", err)
	}

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

	clientNetworkInterfaces := network.NewInterfacesClient(to.String(d.sub.SubscriptionID))
	d.apply(&clientNetworkInterfaces.Client)

	resultNetworkInterfaces, err := clientNetworkInterfaces.ListAll(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not list network interfaces: %w", err)
	}

	interfaces := resultNetworkInterfaces.Values()
	for i := range interfaces {
		s := d.handleNetworkInterfaces(&interfaces[i])

		log.Infof("Adding network interfaces %+v", s)

		list = append(list, s)
	}

	return list, err
}

// Discover load balancer
func (d *azureNetworkDiscovery) discoverLoadBalancer() ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	clientLoadBalancer := network.NewLoadBalancersClient(to.String(d.sub.SubscriptionID))
	d.apply(&clientLoadBalancer.Client)

	resultLoadBalancer, err := clientLoadBalancer.ListAll(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not list load balancer: %w", err)
	}

	lbs := resultLoadBalancer.Values()
	for i := range lbs {
		s := d.handleLoadBalancer(&lbs[i])

		log.Infof("Adding load balancer %+v", s)

		list = append(list, s)
	}

	return list, err
}

func (d *azureNetworkDiscovery) handleLoadBalancer(lb *network.LoadBalancer) voc.IsNetwork {
	return &voc.LoadBalancer{
		NetworkService: &voc.NetworkService{
			Networking: &voc.Networking{
				CloudResource: &voc.CloudResource{
					ID:           voc.ResourceID(to.String(lb.ID)),
					Name:         to.String(lb.Name),
					CreationTime: 0, // No creation time available
					Type:         []string{"LoadBalancer", "NetworkService", "Resource"},
					GeoLocation: voc.GeoLocation{
						Region: *lb.Location,
					},
				},
			},
			Ips:   []string{d.publicIPAddressFromLoadBalancer(lb)},
			Ports: LoadBalancerPorts(lb),
		},
		// TODO(all): do we need the AccessRestriction for load balancers?
		AccessRestrictions: &[]voc.AccessRestriction{},
		// TODO(all): do we need the httpEndpoint for load balancers?
		HttpEndpoints: &[]voc.HttpEndpoint{},
	}
}

func (*azureNetworkDiscovery) handleNetworkInterfaces(ni *network.Interface) voc.IsNetwork {
	return &voc.NetworkInterface{
		Networking: &voc.Networking{
			CloudResource: &voc.CloudResource{
				ID:           voc.ResourceID(to.String(ni.ID)),
				Name:         to.String(ni.Name),
				CreationTime: 0, // No creation time available
				Type:         []string{"NetworkInterface", "Compute", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: *ni.Location,
				},
			},
		},
		//AccessRestriction: &voc.AccessRestriction{
		//	Inbound:         false, // TODO(garuppel): TBD
		//	RestrictedPorts: d.getRestrictedPorts(ni),
		//},
	}
}

func LoadBalancerPorts(lb *network.LoadBalancer) (loadBalancerPorts []int16) {

	for _, item := range *lb.LoadBalancingRules {
		loadBalancerPorts = append(loadBalancerPorts, int16(*item.FrontendPort))
	}

	return loadBalancerPorts
}

//// Returns all restricted ports for the network interface
//func (d *azureNetworkDiscovery) getRestrictedPorts(ni *network.Interface) string {
//
//	var restrictedPorts []string
//
//	if ni.InterfacePropertiesFormat == nil ||
//		ni.InterfacePropertiesFormat.NetworkSecurityGroup == nil ||
//		ni.InterfacePropertiesFormat.NetworkSecurityGroup.ID == nil {
//		return ""
//	}
//
//	nsgID := to.String(ni.NetworkSecurityGroup.ID)
//
//	client := network.NewSecurityGroupsClient(to.String(d.sub.SubscriptionID))
//	d.apply(&client.Client)
//
//	// Get the Security Group of the network interface ni
//	sg, err := client.Get(context.Background(), getResourceGroupName(nsgID), strings.Split(nsgID, "/")[8], "")
//
//	if err != nil {
//		log.Errorf("Could not get security group: %v", err)
//		return ""
//	}
//
//	if sg.SecurityGroupPropertiesFormat != nil && sg.SecurityGroupPropertiesFormat.SecurityRules != nil {
//		// Find all ports defined in the security rules with access property "Deny"
//		for _, securityRule := range *sg.SecurityRules {
//			if securityRule.Access == network.SecurityRuleAccessDeny {
//				restrictedPorts = append(restrictedPorts, *securityRule.SourcePortRange)
//			}
//		}
//	}
//
//	restrictedPortsClean := deleteDuplicatesFromSlice(restrictedPorts)
//
//	return strings.Join(restrictedPortsClean, ",")
//}
//
//func deleteDuplicatesFromSlice(intSlice []string) []string {
//	keys := make(map[string]bool)
//	var list []string
//	for _, entry := range intSlice {
//		if _, value := keys[entry]; !value {
//			keys[entry] = true
//			list = append(list, entry)
//		}
//	}
//	return list
//}

func (d *azureNetworkDiscovery) publicIPAddressFromLoadBalancer(lb *network.LoadBalancer) string {

	var publicIPAddresses []string

	// Get public IP resource
	client := network.NewPublicIPAddressesClient(to.String(d.sub.SubscriptionID))
	d.apply(&client.Client)

	if lb.LoadBalancerPropertiesFormat != nil && lb.LoadBalancerPropertiesFormat.FrontendIPConfigurations != nil {
		for _, publicIpProperties := range *lb.FrontendIPConfigurations {

			publicIPAddress, err := client.Get(context.Background(), getResourceGroupName(*publicIpProperties.ID), *publicIpProperties.Name, "")

			if err != nil {
				log.Infof("Error getting public IP address: %v", err)
				continue
			}

			publicIPAddresses = append(publicIPAddresses, *publicIPAddress.IPAddress)
		}
	}

	return strings.Join(publicIPAddresses, ",")
}
