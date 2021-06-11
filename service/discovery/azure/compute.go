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
	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/compute/mgmt/compute"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/Azure/go-autorest/autorest/to"
)

type azureComputeDiscovery struct {
	azureDiscovery
}

func NewAzureComputeDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &azureComputeDiscovery{}

	for _, opt := range opts {
		if auth, ok := opt.(*authorizerOption); ok {
			d.authOption = auth
		} else {
			d.options = append(d.options, opt)
		}
	}

	return d
}

func (d *azureComputeDiscovery) Name() string {
	return "Azure Compute"
}

func (d *azureComputeDiscovery) Description() string {
	return "Discovery Azure compute."
}

// Discover compute resources
func (d *azureComputeDiscovery) List() (list []voc.IsResource, err error) {
	if err = d.authorize(); err != nil {
		return nil, fmt.Errorf("could not authorize Azure account: %w", err)
	}

	// Discover virtual machines
	virtualMachines, _ := d.discoverVirtualMachines()
	list = append(list, virtualMachines...)

	// Discover network interfaces
	networkInterfaces, _ := d.discoverNetworkInterfaces()
	list = append(list, networkInterfaces...)

	// Discover Load Balancer
	loadBalancer, _ := d.discoverLoadBalancer()
	list = append(list, loadBalancer...)

	return
}

// Discover virtual machines
func (d *azureComputeDiscovery) discoverVirtualMachines() ([]voc.IsResource, error) {
	var list []voc.IsResource

	client := compute.NewVirtualMachinesClient(to.String(d.sub.SubscriptionID))
	d.apply(&client.Client)

	result, err := client.ListAllComplete(context.Background(), "true")
	if err != nil {
		return nil, fmt.Errorf("could not list virtual machines: %w", err)
	}

	vms := *result.Response().Value
	for i := range vms {
		s := d.handleVirtualMachines(&vms[i])

		log.Infof("Adding virtual machine %+v", s)

		list = append(list, s)
	}

	return list, err
}

// Discover network interfaces
func (d *azureComputeDiscovery) discoverNetworkInterfaces() ([]voc.IsResource, error) {
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
func (d *azureComputeDiscovery) discoverLoadBalancer() ([]voc.IsResource, error) {
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

//TBD
func (d *azureComputeDiscovery) handleLoadBalancer(lb *network.LoadBalancer) voc.IsCompute {
	return &voc.LoadBalancerResource{
		NetworkService: voc.NetworkService{
			Resource: voc.Resource{
				ID:           to.String(lb.ID),
				Name:         to.String(lb.Name),
				CreationTime: 0, // No creation time available
				Type:         []string{"LoadBalancer", "NetworkService", "Resource"},
			},
			IPs: []string{d.GetPublicIPAddress(lb)},
		},
		// TODO: fill out access restrictions
		AccessRestriction: &voc.AccessRestriction{
			Inbound:         false,
			RestrictedPorts: "", //TBD
		},
		// TODO: do we need the httpEndpoint?
		HttpEndpoints: []*voc.HttpEndpoint{{
			URL:                 d.GetPublicIPAddress(lb),                     // Get Public IP Address of the Load Balancer
			TransportEncryption: voc.NewTransportEncryption(false, false, ""), // No transport encryption defined for the Load Balancer
		}},
	}
}

func (d *azureComputeDiscovery) GetPublicIPAddress(lb *network.LoadBalancer) string {

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

func (d *azureComputeDiscovery) handleVirtualMachines(vm *compute.VirtualMachine) voc.IsCompute {
	return &voc.VirtualMachineResource{
		ComputeResource: voc.ComputeResource{
			Resource: voc.Resource{
				ID:           to.String(vm.ID),
				Name:         to.String(vm.Name),
				CreationTime: 0, // No creation time available
				Type:         []string{"VirtualMachine", "Compute", "Resource"},
			}},
		Log: &voc.Log{
			Enabled: IsBootDiagnosticEnabled(vm),
		},
	}
}

func IsBootDiagnosticEnabled(vm *compute.VirtualMachine) bool {
	if vm.DiagnosticsProfile == nil {
		return false
	} else {
		return to.Bool(vm.DiagnosticsProfile.BootDiagnostics.Enabled)
	}
}

func (d *azureComputeDiscovery) handleNetworkInterfaces(ni *network.Interface) voc.IsCompute {
	return &voc.NetworkInterfaceResource{
		ComputeResource: voc.ComputeResource{
			Resource: voc.Resource{
				ID:           to.String(ni.ID),
				Name:         to.String(ni.Name),
				CreationTime: 0, // No creation time available
				Type:         []string{"NetworkInterface", "Compute", "Resource"},
			}},
		VmID: GetVmID(*ni),
		AccessRestriction: &voc.AccessRestriction{
			Inbound:         false, //TBD
			RestrictedPorts: d.GetRestrictedPortsDefined(ni),
		},
	}
}

// Returns all restricted ports for the network interface
func (d *azureComputeDiscovery) GetRestrictedPortsDefined(ni *network.Interface) string {

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
				// TODO delete duplicates
				restrictedPorts = append(restrictedPorts, *securityRule.SourcePortRange)
			}
		}
	}

	return strings.Join(restrictedPorts, ",")
}

func GetResourceGroupName(nsgID string) string {
	log.Infof(strings.Split(nsgID, "/")[4])
	return strings.Split(nsgID, "/")[4]
}

func GetVmID(ni network.Interface) string {
	if ni.VirtualMachine == nil {
		return ""
	} else {
		return *ni.VirtualMachine.ID
	}

}
