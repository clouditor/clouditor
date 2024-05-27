// Copyright 2024 Fraunhofer AISEC
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

	"clouditor.io/clouditor/v2/internal/util"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
)

// nsgFirewallEnabled checks if network security group (NSG) rules are configured. A NSG is a firewall that operates at OSI layers 3 and 4 to filter ingress and egress traffic. (https://learn.microsoft.com/en-us/azure/firewall/firewall-faq#what-is-the-difference-between-network-security-groups--nsgs--and-azure-firewall, Last access: 05/02/2023)
func (d *azureDiscovery) nsgFirewallEnabled(ni *armnetwork.Interface) bool {
	// initialize network interfaces client
	if err := d.initNetworkSecurityGroupClient(); err != nil {
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

// loadBalancerPorts returns the external endpoint ports
func loadBalancerPorts(lb *armnetwork.LoadBalancer) (loadBalancerPorts []uint32) {
	for _, item := range lb.Properties.LoadBalancingRules {
		loadBalancerPorts = append(loadBalancerPorts, uint32(util.Deref(item.Properties.FrontendPort)))
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
