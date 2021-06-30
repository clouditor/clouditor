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

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/compute/mgmt/compute"
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
	virtualMachines, err := d.discoverVirtualMachines()
	if err != nil {
		return nil, fmt.Errorf("could not discover virtual machines: %w", err)
	}
	list = append(list, virtualMachines...)

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

func (d *azureComputeDiscovery) handleVirtualMachines(vm *compute.VirtualMachine) voc.IsCompute {
	return &voc.VirtualMachineResource{
		ComputeResource: voc.ComputeResource{
			Resource: voc.Resource{
				ID:           to.String(vm.ID),
				Name:         to.String(vm.Name),
				CreationTime: 0, // No creation time available
				Type:         []string{"VirtualMachine", "Compute", "Resource"},
			}},
		NetworkInterfaceResourceID: d.getNetworkInterfaceResourceID(vm),
		Log: &voc.Log{
			Enabled: IsBootDiagnosticEnabled(vm),
		},
	}
}

func (d *azureComputeDiscovery) getNetworkInterfaceResourceID(vm *compute.VirtualMachine) []string {

	var networkInterfaceIDs []string

	client := compute.NewVirtualMachinesClient(to.String(d.sub.SubscriptionID))
	d.apply(&client.Client)

	// TODO What do we do with this error? Ignore or handle?
	result, err := client.Get(context.Background(), GetResourceGroupName(*vm.ID), *vm.Name, "")
	if err != nil {
		return []string{}
	}

	networkInterfacesList := *result.NetworkProfile.NetworkInterfaces
	for i := range networkInterfacesList {
		networkInterfaceIDs = append(networkInterfaceIDs, *networkInterfacesList[i].ID)
	}

	return networkInterfaceIDs
}

func IsBootDiagnosticEnabled(vm *compute.VirtualMachine) bool {
	if vm.DiagnosticsProfile == nil {
		return false
	} else {
		return to.Bool(vm.DiagnosticsProfile.BootDiagnostics.Enabled)
	}
}
