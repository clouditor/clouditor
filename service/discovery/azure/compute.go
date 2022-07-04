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
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/go-autorest/autorest/to"
)

type azureComputeDiscovery struct {
	azureDiscovery
}

func NewAzureComputeDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &azureComputeDiscovery{}

	for _, opt := range opts {
		if auth, ok := opt.(*authOption); ok {
			d.authCredentials = auth
		} else {
			d.options = append(d.options, opt)
		}
	}

	return d
}

func (*azureComputeDiscovery) Name() string {
	return "Azure Compute"
}

func (*azureComputeDiscovery) Description() string {
	return "Discovery Azure compute."
}

// List compute resources
func (d *azureComputeDiscovery) List() (list []voc.IsCloudResource, err error) {
	if err = d.authorize(); err != nil {
		return nil, fmt.Errorf("could not authorize Azure account: %w", err)
	}

	log.Info("Discover Azure compute resources")

	// Discover virtual machines
	virtualMachines, err := d.discoverVirtualMachines()
	if err != nil {
		return nil, fmt.Errorf("could not discover virtual machines: %w", err)
	}
	list = append(list, virtualMachines...)

	// Discover functions
	function, err := d.discoverFunctions()
	if err != nil {
		return nil, fmt.Errorf("could not discover functions: %w", err)
	}
	list = append(list, function...)

	return
}

// Discover function
func (d *azureComputeDiscovery) discoverFunctions() ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	client, err := armappservice.NewWebAppsClient(to.String(d.sub.SubscriptionID), d.authCredentials.credential, &arm.ClientOptions{})
	if err != nil {
		err = fmt.Errorf("could not get new web apps client: %w", err)
	}

	// List functions
	listPager := client.NewListPager(&armappservice.WebAppsClientListOptions{})
	functionApps := make([]*armappservice.Site, 0)
	for listPager.More() {
		pageResponse, err := listPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("error getting next page: %v", err)
			return nil, err
		}
		functionApps = append(functionApps, pageResponse.Value...)
	}

	// functionApp := *result.Response().Value
	for i := range functionApps {
		r := d.handleFunction(functionApps[i])

		log.Infof("Adding function %+v", r)

		list = append(list, r)
	}

	return list, err
}

func (*azureComputeDiscovery) handleFunction(function *armappservice.Site) voc.IsCompute {

	return &voc.Function{
		Compute: &voc.Compute{
			Resource: &voc.Resource{
				ID:           voc.ResourceID(to.String(function.ID)),
				Name:         to.String(function.Name),
				CreationTime: 0, // No creation time available
				Type:         []string{"Function", "Compute", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: to.String(function.Location),
				},
				Labels: labels(function.Tags),
			},
		},
	}
}

// Discover virtual machines
func (d *azureComputeDiscovery) discoverVirtualMachines() ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	// Create VM client
	client, err := armcompute.NewVirtualMachinesClient(to.String(d.sub.SubscriptionID), d.authCredentials.credential, &arm.ClientOptions{})
	if err != nil {
		err = fmt.Errorf("could not get new virtual machines client: %w", err)
		return nil, err
	}

	// List all VMs accross all resource groups
	listPager := client.NewListAllPager(&armcompute.VirtualMachinesClientListAllOptions{})
	vms := make([]*armcompute.VirtualMachine, 0)
	for listPager.More() {
		pageResponse, err := listPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("error getting next page: %v", err)
			return nil, err
		}
		vms = append(vms, pageResponse.Value...)
	}

	for i := range vms {
		r, err := d.handleVirtualMachines(vms[i])
		if err != nil {
			return nil, fmt.Errorf("could not handle virtual machine: %w", err)
		}

		log.Infof("Adding virtual machine %+v", r)

		list = append(list, r)
	}

	return list, err
}

func (d *azureComputeDiscovery) handleVirtualMachines(vm *armcompute.VirtualMachine) (voc.IsCompute, error) {

	r := &voc.VirtualMachine{
		Compute: &voc.Compute{
			Resource: &voc.Resource{
				ID:           voc.ResourceID(to.String(vm.ID)),
				Name:         to.String(vm.Name),
				CreationTime: vm.Properties.TimeCreated.Unix(),
				Type:         []string{"VirtualMachine", "Compute", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: to.String(vm.Location),
				},
				Labels: labels(vm.Tags),
			}},
		BootLogging: &voc.BootLogging{Logging: &voc.Logging{
			Enabled:         isBootDiagnosticEnabled(vm),
			LoggingService:  []voc.ResourceID{voc.ResourceID(bootLogOutput(vm))},
			RetentionPeriod: 0, // Currently, configuring the retention period for Managed Boot Diagnostics is not available. The logs will be overwritten after 1gb of space according to https://github.com/MicrosoftDocs/azure-docs/issues/69953
		}},
		OSLogging: &voc.OSLogging{
			Logging: &voc.Logging{
				Enabled:         false,
				RetentionPeriod: 0,
				LoggingService:  []voc.ResourceID{},
			},
		},
		BlockStorage: []voc.ResourceID{},
	}

	// Reference to networkInterfaces
	for _, networkInterfaces := range vm.Properties.NetworkProfile.NetworkInterfaces {
		r.NetworkInterface = append(r.NetworkInterface, voc.ResourceID(to.String(networkInterfaces.ID)))
	}

	// Reference to blockstorage
	r.BlockStorage = append(r.BlockStorage, voc.ResourceID(to.String(vm.Properties.StorageProfile.OSDisk.ManagedDisk.ID)))
	for _, blockstorage := range vm.Properties.StorageProfile.DataDisks {
		r.BlockStorage = append(r.BlockStorage, voc.ResourceID(to.String(blockstorage.ManagedDisk.ID)))
	}

	return r, nil
}

func isBootDiagnosticEnabled(vm *armcompute.VirtualMachine) bool {
	if vm.Properties.DiagnosticsProfile == nil {
		return false
	} else {
		return to.Bool(vm.Properties.DiagnosticsProfile.BootDiagnostics.Enabled)
	}
}

func bootLogOutput(vm *armcompute.VirtualMachine) string {
	if isBootDiagnosticEnabled(vm) {
		// If storageUri is not specified while enabling boot diagnostics, managed storage will be used.
		if vm.Properties.DiagnosticsProfile.BootDiagnostics.StorageURI != nil {
			return to.String(vm.Properties.DiagnosticsProfile.BootDiagnostics.StorageURI)
		}

		return ""
	}
	return ""
}
