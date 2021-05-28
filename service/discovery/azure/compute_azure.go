/*
 * Copyright 2021 Fraunhofer AISEC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *           $$\                           $$\ $$\   $$\
 *           $$ |                          $$ |\__|  $$ |
 *  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
 * $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
 * $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
 * $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
 * \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
 *  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
 *
 * This file is part of Clouditor Community Edition.
 */

package azure

import (
	"context"
	"strings"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/compute/mgmt/compute"
	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/resources/mgmt/subscriptions"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

type azureComputeDiscovery struct{}

func NewAzureComputeDiscovery() discovery.Discoverer {
	return &azureComputeDiscovery{}
}

func (d *azureComputeDiscovery) Name() string {
	return "Azure Compute"
}

func (d *azureComputeDiscovery) Description() string {
	return "Discovery Azure compute."
}

// Discover compute resources
func (d *azureComputeDiscovery) List() (list []voc.IsResource, err error) {
	// create an authorizer from env vars or Azure Managed Service Identity
	authorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		log.Errorf("Could not authenticate to Azure: %s", err)
		return
	}

	subClient := subscriptions.NewClient()
	subClient.Authorizer = authorizer

	// get first subcription
	page, _ := subClient.List(context.Background())
	sub := page.Values()[0]

	log.Infof("Using %s as subscription", *sub.SubscriptionID)

	ctx := context.Background()

	// Discover virtual machines
	client := compute.NewVirtualMachinesClient(*sub.SubscriptionID)
	client.Authorizer = authorizer

	result, _ := client.ListAllComplete(ctx, "true")

	for _, v := range *result.Response().Value {
		s := handleVirtualMachines(v)

		log.Infof("Adding virtual machine %+v", s)

		list = append(list, s)
	}

	// Discover network interfaces
	client_network_interfaces := network.NewInterfacesClient(*sub.SubscriptionID)
	client_network_interfaces.Authorizer = authorizer

	result_network_interfaces, _ := client_network_interfaces.ListAll(ctx)

	for _, ni := range result_network_interfaces.Values() {
		s := handleNetworkInterfaces(ni)

		log.Infof("Adding virtual machine %+v", s)

		list = append(list, s)
	}

	return
}

func handleVirtualMachines(vm compute.VirtualMachine) voc.IsCompute {
	return &voc.VirtualMachineResource{
		ComputeResource: voc.ComputeResource{
			Resource: voc.Resource{
				ID:           *vm.ID,
				Name:         *vm.Name,
				CreationTime: 0, // No creation time available
			}},
		Log: &voc.Log{
			Enabled: IsBootDiagnosticEnabled(vm),
		},
	}
}

func IsBootDiagnosticEnabled(vm compute.VirtualMachine) bool {
	if vm.DiagnosticsProfile == nil {
		return false
	} else {
		return *vm.DiagnosticsProfile.BootDiagnostics.Enabled
	}
}

func handleNetworkInterfaces(ni network.Interface) voc.IsCompute {
	return &voc.NetworkInterfaceResource{
		ComputeResource: voc.ComputeResource{
			Resource: voc.Resource{
				ID:           *ni.ID,
				Name:         *ni.Name,
				CreationTime: 0, // No creation time available
			}},
		VmID: GetVmID(ni),
		AccessRestriction: &voc.AccessRestriction{
			Inbound:         false, //TBD
			RestrictedPorts: AreRestrictedPortsDefined(ni),
		},
	}
}

// TODO What is the definition of restricted ports?
// For now it is checked if an user-configured inbound security rule is enabled
func AreRestrictedPortsDefined(ni network.Interface) bool {

	if ni.NetworkSecurityGroup.ID == nil {
		return false
	}

	nsgID := *ni.NetworkSecurityGroup.ID

	// TODO refactor authorizer
	// create an authorizer from env vars or Azure Managed Service Identity
	authorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		log.Errorf("Could not authenticate to Azure: %s", err)
		return false
	}

	subClient := subscriptions.NewClient()
	subClient.Authorizer = authorizer

	// get first subcription
	page, _ := subClient.List(context.Background())
	sub := page.Values()[0]

	ctx := context.Background()

	client2 := network.NewSecurityGroupsClient(*sub.SubscriptionID)
	client2.Authorizer = authorizer
	sg, err := client2.Get(ctx, GetResourceGroupName(nsgID), strings.Split(nsgID, "/")[8], "")

	if err != nil {
		log.Errorf("Could not get security group: %s", err)
		return false
	}

	for _, securityRule := range *sg.SecurityRules {
		if securityRule.Access == network.SecurityRuleAccessAllow {
			return true
		}
	}

	return false
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
