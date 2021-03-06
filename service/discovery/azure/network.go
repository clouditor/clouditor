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

func (d *azureNetworkDiscovery) Name() string {
	return "Azure Network"
}

func (d *azureNetworkDiscovery) Description() string {
	return "Discovery Azure network resources."
}

// Discover network resources
func (d *azureNetworkDiscovery) List() (list []voc.IsResource, err error) {
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
func (d *azureNetworkDiscovery) discoverNetworkInterfaces() ([]voc.IsResource, error) {
	var list []voc.IsResource

	client_network_interfaces := network.NewInterfacesClient(to.String(d.sub.SubscriptionID))
	d.apply(&client_network_interfaces.Client)

	result_network_interfaces, err := client_network_interfaces.ListAll(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not list network interfaces: %w", err)
	}

	interfaces := result_network_interfaces.Values()
	for i := range interfaces {
		s := d.handleNetworkInterfaces(&interfaces[i])

		log.Infof("Adding network interfaces %+v", s)

		list = append(list, s)
	}

	return list, err
}

// Discover Load Balancer
func (d *azureNetworkDiscovery) discoverLoadBalancer() ([]voc.IsResource, error) {
	var list []voc.IsResource

	client_load_balancer := network.NewLoadBalancersClient(to.String(d.sub.SubscriptionID))
	d.apply(&client_load_balancer.Client)

	result_load_balancer, err := client_load_balancer.ListAll(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not list load balancer: %w", err)
	}

	lbs := result_load_balancer.Values()
	for i := range lbs {
		s := d.handleLoadBalancer(&lbs[i])

		log.Infof("Adding load balancer %+v", s)

		list = append(list, s)
	}

	return list, err
}

func (d *azureNetworkDiscovery) handleLoadBalancer(lb *network.LoadBalancer) voc.IsNetwork {
	return &voc.LoadBalancerResource{
		NetworkService: voc.NetworkService{
			NetworkResource: voc.NetworkResource{
				Resource: voc.Resource{
					ID:           voc.ResourceID(to.String(lb.ID)),
					Name:         to.String(lb.Name),
					CreationTime: 0, // No creation time available
					Type:         []string{"LoadBalancer", "NetworkService", "Resource"},
				},
			},
			IPs:   []string{d.GetPublicIPAddress(lb)},
			Ports: getLoadBalancerPorts(lb),
		},
		// TODO: do we need the AccessRestriction for load balancers?
		AccessRestriction: &voc.AccessRestriction{},
		// TODO: do we need the httpEndpoint for load balancers?
		HttpEndpoints: []*voc.HttpEndpoint{}}
}

func (d *azureNetworkDiscovery) handleNetworkInterfaces(ni *network.Interface) voc.IsNetwork {
	return &voc.NetworkInterface{
		NetworkResource: voc.NetworkResource{
			Resource: voc.Resource{
				ID:           voc.ResourceID(to.String(ni.ID)),
				Name:         to.String(ni.Name),
				CreationTime: 0, // No creation time available
				Type:         []string{"NetworkInterface", "Compute", "Resource"},
			},
		},
		AccessRestriction: &voc.AccessRestriction{
			Inbound:         false, //TBD
			RestrictedPorts: d.getRestrictedPorts(ni),
		},
	}
}

func getLoadBalancerPorts(lb *network.LoadBalancer) (loadBalancerPorts []int16) {

	for _, item := range *lb.LoadBalancingRules {
		loadBalancerPorts = append(loadBalancerPorts, int16(*item.FrontendPort))
	}

	return loadBalancerPorts
}

// Returns all restricted ports for the network interface
func (d *azureNetworkDiscovery) getRestrictedPorts(ni *network.Interface) string {

	var restrictedPorts []string

	if ni.InterfacePropertiesFormat == nil ||
		ni.InterfacePropertiesFormat.NetworkSecurityGroup == nil ||
		ni.InterfacePropertiesFormat.NetworkSecurityGroup.ID == nil {
		return ""
	}

	nsgID := to.String(ni.NetworkSecurityGroup.ID)

	client := network.NewSecurityGroupsClient(to.String(d.sub.SubscriptionID))
	d.apply(&client.Client)

	// Get the Security Group of the network interface ni
	sg, err := client.Get(context.Background(), GetResourceGroupName(nsgID), strings.Split(nsgID, "/")[8], "")

	if err != nil {
		log.Errorf("Could not get security group: %s", err)
		return ""
	}

	if sg.SecurityGroupPropertiesFormat != nil && sg.SecurityGroupPropertiesFormat.SecurityRules != nil {
		// Find all ports defined in the security rules with access property "Deny"
		for _, securityRule := range *sg.SecurityRules {
			if securityRule.Access == network.SecurityRuleAccessDeny {
				restrictedPorts = append(restrictedPorts, *securityRule.SourcePortRange)
			}
		}
	}

	restrictedPortsClean := deleteDuplicatesFromSlice(restrictedPorts)

	return strings.Join(restrictedPortsClean, ",")
}

func deleteDuplicatesFromSlice(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func GetResourceGroupName(nsgID string) string {
	log.Infof(strings.Split(nsgID, "/")[4])
	return strings.Split(nsgID, "/")[4]
}

func (d *azureNetworkDiscovery) GetPublicIPAddress(lb *network.LoadBalancer) string {

	var publicIPAddresses []string

	// Get public IP resource
	client := network.NewPublicIPAddressesClient(to.String(d.sub.SubscriptionID))
	d.apply(&client.Client)

	if lb.LoadBalancerPropertiesFormat != nil && lb.LoadBalancerPropertiesFormat.FrontendIPConfigurations != nil {
		for _, publicIpProperties := range *lb.FrontendIPConfigurations {

			publicIPAddress, err := client.Get(context.Background(), GetResourceGroupName(*publicIpProperties.ID), *publicIpProperties.Name, "")

			if err != nil {
				log.Errorf("Error getting public IP address: %v", err)
				continue
			}

			publicIPAddresses = append(publicIPAddresses, *publicIPAddress.IPAddress)
		}
	}

	return strings.Join(publicIPAddresses, ",")
}
