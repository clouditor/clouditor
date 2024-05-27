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
	"errors"
	"fmt"

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v3"
)

var (
	ErrEmptyVirtualMachine = errors.New("virtual machine is empty")
)

// Discover virtual machines
func (d *azureDiscovery) discoverVirtualMachines() ([]ontology.IsResource, error) {
	var list []ontology.IsResource

	// initialize virtual machines client
	if err := d.initVirtualMachinesClient(); err != nil {
		return nil, err
	}

	// List all VMs
	err := listPager(d,
		d.clients.virtualMachinesClient.NewListAllPager,
		d.clients.virtualMachinesClient.NewListPager,
		func(res armcompute.VirtualMachinesClientListAllResponse) []*armcompute.VirtualMachine {
			return res.Value
		},
		func(res armcompute.VirtualMachinesClientListResponse) []*armcompute.VirtualMachine {
			return res.Value
		},
		func(vm *armcompute.VirtualMachine) error {
			r, err := d.handleVirtualMachines(vm)
			if err != nil {
				return fmt.Errorf("could not handle virtual machine: %w", err)
			}

			log.Infof("Adding virtual machine '%s'", r.GetName())

			list = append(list, r)

			return nil
		})
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (d *azureDiscovery) discoverBlockStorages() ([]ontology.IsResource, error) {
	var list []ontology.IsResource

	// initialize block storages client
	if err := d.initBlockStoragesClient(); err != nil {
		return nil, err
	}

	// List all disks
	err := listPager(d,
		d.clients.blockStorageClient.NewListPager,
		d.clients.blockStorageClient.NewListByResourceGroupPager,
		func(res armcompute.DisksClientListResponse) []*armcompute.Disk {
			return res.Value
		},
		func(res armcompute.DisksClientListByResourceGroupResponse) []*armcompute.Disk {
			return res.Value
		},
		func(disk *armcompute.Disk) error {
			blockStorage, err := d.handleBlockStorage(disk)
			if err != nil {
				return fmt.Errorf("could not handle block storage: %w", err)
			}

			log.Infof("Adding block storage '%s'", blockStorage.GetName())

			list = append(list, blockStorage)
			return nil
		})
	if err != nil {
		return nil, err
	}

	return list, nil
}

// Discover functions and web apps
func (d *azureDiscovery) discoverFunctionsWebApps() ([]ontology.IsResource, error) {
	var list []ontology.IsResource

	// initialize functions client
	if err := d.initWebAppsClient(); err != nil {
		return nil, err
	}

	// List functions
	err := listPager(d,
		d.clients.webAppsClient.NewListPager,
		d.clients.webAppsClient.NewListByResourceGroupPager,
		func(res armappservice.WebAppsClientListResponse) []*armappservice.Site {
			return res.Value
		},
		func(res armappservice.WebAppsClientListByResourceGroupResponse) []*armappservice.Site {
			return res.Value
		},
		func(site *armappservice.Site) error {
			var r ontology.IsResource

			// Get configuration for detailed properties
			config, err := d.clients.webAppsClient.GetConfiguration(context.Background(),
				util.Deref(site.Properties.ResourceGroup),
				util.Deref(site.Name),
				&armappservice.WebAppsClientGetConfigurationOptions{})
			if err != nil {
				log.Errorf("error getting site config: %v", err)
			}

			// Check kind of site (see https://github.com/Azure/app-service-linux-docs/blob/master/Things_You_Should_Know/kind_property.md)
			switch *site.Kind {
			case "app": // Windows Web App
				r = d.handleWebApp(site, config)
			case "app,linux": // Linux Web app
				r = d.handleWebApp(site, config)
			case "app,linux,container": // Linux Container Web App
				// TODO(all): TBD
				log.Debug("Linux Container Web App Web App currently not implemented.")
			case "hyperV": // Windows Container Web App
				// TODO(all): TBD
				log.Debug("Windows Container Web App currently not implemented.")
			case "app,container,windows": // Windows Container Web App
				// TODO(all): TBD
				log.Debug("Windows Web App currently not implemented.")
			case "app,linux,kubernetes": // Linux Web App on ARC
				// TODO(all): TBD
				log.Debug("Linux Web App on ARC currently not implemented.")
			case "app,linux,container,kubernetes": // Linux Container Web App on ARC
				// TODO(all): TBD
				log.Debug("Linux Container Web App on ARC currently not implemented.")
			case "functionapp": // Function Code App
				r = d.handleFunction(site, config)
			case "functionapp,linux": // Linux Consumption Function app
				r = d.handleFunction(site, config)
			case "functionapp,linux,container,kubernetes": // Function Container App on ARC
				// TODO(all): TBD
				log.Debug("Function Container App on ARC currently not implemented.")
			case "functionapp,linux,kubernetes": // Function Code App on ARC
				// TODO(all): TBD
				log.Debug("Function Code App on ARC currently not implemented.")
			default:
				log.Debugf("%s currently not supported.", *site.Kind)
			}

			if r != nil {
				log.Infof("Adding function %+v", r.GetName())
				list = append(list, r)
			}

			return nil
		})
	if err != nil {
		return nil, err
	}

	return list, nil
}
