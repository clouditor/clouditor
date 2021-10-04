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
	"github.com/Azure/azure-sdk-for-go/profiles/latest/web/mgmt/web"
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

// List compute resources
func (d *azureComputeDiscovery) List() (list []voc.IsCloudResource, err error) {
	if err = d.authorize(); err != nil {
		return nil, fmt.Errorf("could not authorize Azure account: %w", err)
	}

	// Discover virtual machines
	virtualMachines, err := d.discoverVirtualMachines()
	if err != nil {
		return nil, fmt.Errorf("could not discover virtual machines: %w", err)
	}
	list = append(list, virtualMachines...)

	// Discover functions
	function, err := d.discoverFunction()
	if err != nil {
		return nil, fmt.Errorf("could not discover functions: %w", err)
	}
	list = append(list, function...)

	return
}

// Discover function
func (d *azureComputeDiscovery) discoverFunction() ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	client := web.NewAppsClient(to.String(d.sub.SubscriptionID))
	d.apply(&client.Client)

	result, err := client.ListComplete(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not list functions: %w", err)
	}

	functionApp := *result.Response().Value
	for i := range functionApp {
		functionResource := d.handleFunction(&functionApp[i])
		list = append(list, functionResource)
	}

	return list, err
}

func (d *azureComputeDiscovery) handleFunction(function *web.Site) voc.IsCompute {
	return &voc.Function{
		Compute: &voc.Compute{
			CloudResource: &voc.CloudResource{
				ID:           voc.ResourceID(to.String(function.ID)),
				Name:         to.String(function.Name),
				CreationTime: 0, // No creation time available
				Type:         []string{"Function", "Compute", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: *function.Location,
				},
			},
		},
	}

}

// Discover virtual machines
func (d *azureComputeDiscovery) discoverVirtualMachines() ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	client := compute.NewVirtualMachinesClient(to.String(d.sub.SubscriptionID))
	d.apply(&client.Client)

	result, err := client.ListAllComplete(context.Background(), "true")
	if err != nil {
		return nil, fmt.Errorf("could not list virtual machines: %w", err)
	}

	vms := *result.Response().Value
	for i := range vms {
		s, err := d.handleVirtualMachines(&vms[i])
		if err != nil {
			return nil, fmt.Errorf("could not handle virtual machine: %w", err)
		}

		log.Infof("Adding virtual machine %+v", s)

		list = append(list, s)
	}

	return list, err
}

func (d *azureComputeDiscovery) handleVirtualMachines(vm *compute.VirtualMachine) (voc.IsCompute, error) {

	r := &voc.VirtualMachine{
		Compute: &voc.Compute{
			CloudResource: &voc.CloudResource{
				ID:           voc.ResourceID(to.String(vm.ID)),
				Name:         to.String(vm.Name),
				CreationTime: 0, // No creation time available
				Type:         []string{"VirtualMachine", "Compute", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: *vm.Location,
				},
			}},
		Log: &voc.Log{
			Activated: IsBootDiagnosticEnabled(vm),
			Auditing: &voc.Auditing{
				SecurityFeature: &voc.SecurityFeature{},
			},
		},
	}

	vmExtended, err := d.getExtendedVirtualMachine(vm)
	if err != nil {
		return nil, fmt.Errorf("could not get virtual machine with extended information: %w", err)
	}

	// Reference to networkInterfaces
	for _, networkInterfaces := range *vmExtended.VirtualMachineProperties.NetworkProfile.NetworkInterfaces {
		r.NetworkInterface = append(r.NetworkInterface, voc.ResourceID(to.String(networkInterfaces.ID)))
	}

	// Reference to blockstorage
	r.BlockStorage = append(r.BlockStorage, voc.ResourceID(*vmExtended.StorageProfile.OsDisk.ManagedDisk.ID))
	for _, blockstorage := range *vmExtended.StorageProfile.DataDisks {
		r.BlockStorage = append(r.BlockStorage, voc.ResourceID(*blockstorage.ManagedDisk.ID))
	}

	return r, nil
}

// Get virtual machine with extended information, e.g., managed disk ID, network interface ID
func (d *azureComputeDiscovery) getExtendedVirtualMachine(vm *compute.VirtualMachine) (*compute.VirtualMachine, error) {
	client := compute.NewVirtualMachinesClient(to.String(d.sub.SubscriptionID))
	d.apply(&client.Client)

	vmExtended, err := client.Get(context.Background(), GetResourceGroupName(*vm.ID), *vm.Name, "")
	if err != nil {
		return nil, fmt.Errorf("could not get virtual machine: %w", err)
	}
	return &vmExtended, nil
}

func IsBootDiagnosticEnabled(vm *compute.VirtualMachine) bool {
	if vm.DiagnosticsProfile == nil {
		return false
	} else {
		return to.Bool(vm.DiagnosticsProfile.BootDiagnostics.Enabled)
	}
}
