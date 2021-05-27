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

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/compute/mgmt/compute"
	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/resources/mgmt/subscriptions"
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

// Discover virtual machines
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

	client := compute.NewVirtualMachinesClient(*sub.SubscriptionID)
	client.Authorizer = authorizer

	ctx := context.Background()

	result, _ := client.ListAllComplete(ctx, "true")

	for _, v := range *result.Response().Value {
		s := handleVirtualMachines(v)

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
				CreationTime: 0, // VM has no creation time
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
