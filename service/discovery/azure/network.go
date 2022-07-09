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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"

	"github.com/Azure/go-autorest/autorest/to"
)

type azureNetworkDiscovery struct {
	azureDiscovery
}

func NewAzureNetworkDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &azureNetworkDiscovery{}

	// Apply options
	for _, opt := range opts {
		opt(&d.azureDiscovery)
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

	// Create network client
	client, err := armnetwork.NewInterfacesClient(to.String(d.sub.SubscriptionID), d.cred, &d.clientOptions)
	if err != nil {
		err = fmt.Errorf("could not get new virtual machines client: %w", err)
		return nil, err
	}

	// List all network interfaces accross all resource groups
	listPager := client.NewListAllPager(&armnetwork.InterfacesClientListAllOptions{})
	ni := make([]*armnetwork.Interface, 0)
	for listPager.More() {
		pageResponse, err := listPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("error getting next page: %v", err)
			return nil, err
		}
		ni = append(ni, pageResponse.Value...)
	}

	for i := range ni {
		s := d.handleNetworkInterfaces(ni[i])

		log.Infof("Adding network interfaces %+v", s)

		list = append(list, s)
	}

	return list, err
}

// Discover load balancer
func (d *azureNetworkDiscovery) discoverLoadBalancer() ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	// Create load balancer client
	client, err := armnetwork.NewLoadBalancersClient(to.String(d.sub.SubscriptionID), d.cred, &d.clientOptions)
	if err != nil {
		err = fmt.Errorf("could not get new load balancers client: %w", err)
		return nil, err
	}

	// List all load balancers accross all resource groups
	listPager := client.NewListAllPager(&armnetwork.LoadBalancersClientListAllOptions{})
	lbs := make([]*armnetwork.LoadBalancer, 0)
	for listPager.More() {
		pageResponse, err := listPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("error getting next page: %v", err)
			return nil, err
		}
		lbs = append(lbs, pageResponse.Value...)
	}

	for i := range lbs {
		s := d.handleLoadBalancer(lbs[i])

		log.Infof("Adding load balancer %+v", s)

		list = append(list, s)
	}

	return list, err
}

func (d *azureNetworkDiscovery) handleLoadBalancer(lb *armnetwork.LoadBalancer) voc.IsNetwork {
	return &voc.LoadBalancer{
		NetworkService: &voc.NetworkService{
			Networking: &voc.Networking{
				Resource: &voc.Resource{
					ID:           voc.ResourceID(to.String(lb.ID)),
					Name:         to.String(lb.Name),
					CreationTime: 0, // No creation time available
					Type:         []string{"LoadBalancer", "NetworkService", "Resource"},
					GeoLocation: voc.GeoLocation{
						Region: to.String(lb.Location),
					},
					Labels: labels(lb.Tags),
				},
			},
			Ips:   d.publicIPAddressFromLoadBalancer(lb),
			Ports: LoadBalancerPorts(lb),
		},
		// TODO(all): do we need the AccessRestriction for load balancers?
		AccessRestrictions: &[]voc.AccessRestriction{},
		// TODO(all): do we need the httpEndpoint for load balancers?
		HttpEndpoints: &[]voc.HttpEndpoint{},
	}
}

func (*azureNetworkDiscovery) handleNetworkInterfaces(ni *armnetwork.Interface) voc.IsNetwork {
	return &voc.NetworkInterface{
		Networking: &voc.Networking{
			Resource: &voc.Resource{
				ID:           voc.ResourceID(to.String(ni.ID)),
				Name:         to.String(ni.Name),
				CreationTime: 0, // No creation time available
				Type:         []string{"NetworkInterface", "Compute", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: to.String(ni.Location),
				},
				Labels: labels(ni.Tags),
			},
		},
		// AccessRestriction: &voc.AccessRestriction{
		// 	Inbound:         false, // TODO(garuppel): TBD
		// 	RestrictedPorts: d.getRestrictedPorts(ni),
		// },
	}
}

func LoadBalancerPorts(lb *armnetwork.LoadBalancer) (loadBalancerPorts []int16) {

	for _, item := range lb.Properties.LoadBalancingRules {
		loadBalancerPorts = append(loadBalancerPorts, int16(to.Int32(item.Properties.FrontendPort)))
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

func (d *azureNetworkDiscovery) publicIPAddressFromLoadBalancer(lb *armnetwork.LoadBalancer) []string {

	var (
		publicIPAddresses []string
	)

	// Create public IP address client
	client, err := armnetwork.NewPublicIPAddressesClient(to.String(d.sub.SubscriptionID), d.cred, &d.clientOptions)
	if err != nil {
		log.Debugf("could not get new public ip addresses client: %v", err)
		return []string{}
	}

	fIpConfig := lb.Properties.FrontendIPConfigurations
	for i := range fIpConfig {

		if fIpConfig[i].Properties.PublicIPAddress == nil {
			continue
		}

		publicIPAddressID := to.String(fIpConfig[i].Properties.PublicIPAddress.ID)
		if publicIPAddressID == "" {
			continue
		}

		publicIpAddressName := frontendPublicIPAddressName(publicIPAddressID)
		if publicIpAddressName == "" {
			continue
		}

		// Get public IP address
		publicIPAddress, err := client.Get(
			context.TODO(),
			resourceGroupName(publicIPAddressID),
			publicIpAddressName,
			&armnetwork.PublicIPAddressesClientGetOptions{})
		if err != nil {
			log.Debugf("could not get public ip address: %v", err)
			return []string{}
		}

		ipAddress := publicIPAddress.PublicIPAddress.Properties.IPAddress
		if ipAddress == nil {
			log.Infof("Error getting public IP address: %v", err)
			continue
		}

		publicIPAddresses = append(publicIPAddresses, to.String(ipAddress))
	}

	return publicIPAddresses
}

// frontendPublicIPAddressName returns the frontend public IP address name from the given public IP address ID
func frontendPublicIPAddressName(frontendPublicIPAddressID string) string {
	if frontendPublicIPAddressID == "" {
		log.Infof("Public IP address ID of frontend is empty.")
		return ""
	}

	split := strings.Split(frontendPublicIPAddressID, "/")
	if len(split) != 9 {
		log.Infof("Public IP address ID of frontend is not correct.")
		return ""
	}

	return split[8]
}
