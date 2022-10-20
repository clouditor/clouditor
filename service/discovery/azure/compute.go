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

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v3"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
)

var (
	ErrEmptyVirtualMachine = errors.New("virtual machine is empty")
)

type azureComputeDiscovery struct {
	azureDiscovery
}

func NewAzureComputeDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &azureComputeDiscovery{
		azureDiscovery{
			discovererComponent: ComputeComponent,
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(&d.azureDiscovery)
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
		return nil, fmt.Errorf("%s: %w", ErrCouldNotAuthenticate, err)
	}

	log.Info("Discover Azure block storage")
	// Discover block storage
	storage, err := d.discoverBlockStorages()
	if err != nil {
		return nil, fmt.Errorf("could not discover block storage: %w", err)
	}
	list = append(list, storage...)

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

	client, err := armappservice.NewWebAppsClient(util.Deref(d.sub.SubscriptionID), d.cred, &d.clientOptions)
	if err != nil {
		err = fmt.Errorf("could not get new web apps client: %w", err)
	}

	// List functions
	listPager := client.NewListPager(&armappservice.WebAppsClientListOptions{})
	functionApps := make([]*armappservice.Site, 0)
	for listPager.More() {
		pageResponse, err := listPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("%s: %v", ErrGettingNextPage, err)
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

	// If a mandatory field is empty, the whole function is empty
	if function == nil || function.ID == nil {
		return nil
	}

	return &voc.Function{
		Compute: &voc.Compute{
			Resource: &voc.Resource{
				ID:           voc.ResourceID(util.Deref(function.ID)),
				ServiceID:    discovery.DefaultCloudServiceID,
				Name:         util.Deref(function.Name),
				CreationTime: 0, // No creation time available
				Type:         []string{"Function", "Compute", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: util.Deref(function.Location),
				},
				Labels: labels(function.Tags),
			},
			NetworkInterface: []voc.ResourceID{},
		},
		RuntimeLanguage: "",
		RuntimeVersion:  "",
	}
}

// Discover virtual machines
func (d *azureComputeDiscovery) discoverVirtualMachines() ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	// Create VM client
	client, err := armcompute.NewVirtualMachinesClient(util.Deref(d.sub.SubscriptionID), d.cred, &d.clientOptions)
	if err != nil {
		err = fmt.Errorf("could not get new virtual machines client: %w", err)
		return nil, err
	}

	// List all VMs across all resource groups
	listPager := client.NewListAllPager(&armcompute.VirtualMachinesClientListAllOptions{})

	for listPager.More() {
		pageResponse, err := listPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("%s: %v", ErrGettingNextPage, err)
			return nil, err
		}

		for _, vm := range pageResponse.Value {
			r, err := d.handleVirtualMachines(vm)
			if err != nil {
				return nil, fmt.Errorf("could not handle virtual machine: %w", err)
			}

			log.Infof("Adding virtual machine '%s'", r.GetName())

			list = append(list, r)
		}
	}

	return list, err
}

func (d *azureComputeDiscovery) handleVirtualMachines(vm *armcompute.VirtualMachine) (voc.IsCompute, error) {
	var bootLogging = []voc.ResourceID{}
	var osLogging = []voc.ResourceID{}

	// If a mandatory field is empty, the whole disk is empty
	if vm == nil || vm.ID == nil {
		return nil, ErrEmptyVirtualMachine
	}

	if bootLogOutput(vm) != "" {
		bootLogging = []voc.ResourceID{voc.ResourceID(bootLogOutput(vm))}
	}

	r := &voc.VirtualMachine{
		Compute: &voc.Compute{
			Resource: &voc.Resource{
				ID:           voc.ResourceID(util.Deref(vm.ID)),
				ServiceID:    discovery.DefaultCloudServiceID,
				Name:         util.Deref(vm.Name),
				CreationTime: util.SafeTimestamp(vm.Properties.TimeCreated),
				Type:         []string{"VirtualMachine", "Compute", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: util.Deref(vm.Location),
				},
				Labels: labels(vm.Tags),
			},
			NetworkInterface: []voc.ResourceID{},
		},
		BlockStorage:      []voc.ResourceID{},
		MalwareProtection: &voc.MalwareProtection{},
		BootLogging: &voc.BootLogging{
			Logging: &voc.Logging{
				Enabled:         isBootDiagnosticEnabled(vm),
				LoggingService:  bootLogging,
				RetentionPeriod: 0, // Currently, configuring the retention period for Managed Boot Diagnostics is not available. The logs will be overwritten after 1gb of space according to https://github.com/MicrosoftDocs/azure-docs/issues/69953
				Auditing: &voc.Auditing{
					SecurityFeature: &voc.SecurityFeature{},
				},
			},
		},
		OSLogging: &voc.OSLogging{
			Logging: &voc.Logging{
				Enabled:         false,
				RetentionPeriod: 0,
				LoggingService:  osLogging,
				Auditing: &voc.Auditing{
					SecurityFeature: &voc.SecurityFeature{},
				},
			},
		},
	}

	// Reference to networkInterfaces
	if vm.Properties.NetworkProfile != nil {
		for _, networkInterfaces := range vm.Properties.NetworkProfile.NetworkInterfaces {
			r.NetworkInterface = append(r.NetworkInterface, voc.ResourceID(util.Deref(networkInterfaces.ID)))
		}
	}

	// Reference to blockstorage
	if vm.Properties.StorageProfile != nil && vm.Properties.StorageProfile.OSDisk != nil && vm.Properties.StorageProfile.OSDisk.ManagedDisk != nil {
		r.BlockStorage = append(r.BlockStorage, voc.ResourceID(util.Deref(vm.Properties.StorageProfile.OSDisk.ManagedDisk.ID)))
	}

	if vm.Properties.StorageProfile != nil && vm.Properties.StorageProfile.DataDisks != nil {
		for _, blockstorage := range vm.Properties.StorageProfile.DataDisks {
			r.BlockStorage = append(r.BlockStorage, voc.ResourceID(util.Deref(blockstorage.ManagedDisk.ID)))
		}
	}

	return r, nil
}

func isBootDiagnosticEnabled(vm *armcompute.VirtualMachine) bool {
	if vm == nil || vm.Properties == nil || vm.Properties.DiagnosticsProfile == nil || vm.Properties.DiagnosticsProfile.BootDiagnostics == nil {
		return false
	} else {
		return util.Deref(vm.Properties.DiagnosticsProfile.BootDiagnostics.Enabled)
	}
}

func bootLogOutput(vm *armcompute.VirtualMachine) string {
	if isBootDiagnosticEnabled(vm) {
		// If storageUri is not specified while enabling boot diagnostics, managed storage will be used.
		if vm.Properties.DiagnosticsProfile.BootDiagnostics.StorageURI != nil {
			return util.Deref(vm.Properties.DiagnosticsProfile.BootDiagnostics.StorageURI)
		}

		return ""
	}
	return ""
}

func (d *azureComputeDiscovery) discoverBlockStorages() ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	// Create disks client
	client, err := armcompute.NewDisksClient(util.Deref(d.sub.SubscriptionID), d.cred, &d.clientOptions)
	if err != nil {
		err = fmt.Errorf("could not get new disks client: %w", err)
		return nil, err
	}

	// List all disks across all resource groups
	listPager := client.NewListPager(&armcompute.DisksClientListOptions{})
	for listPager.More() {
		pageResponse, err := listPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("%s: %w", ErrGettingNextPage, err)
			return nil, err
		}

		for _, disk := range pageResponse.Value {
			blockStorages, err := d.handleBlockStorage(disk)
			if err != nil {
				return nil, fmt.Errorf("could not handle block storage: %w", err)
			}
			log.Infof("Adding block storage '%s'", blockStorages.Name)

			list = append(list, blockStorages)
		}
	}

	return list, nil
}

func (d *azureComputeDiscovery) handleBlockStorage(disk *armcompute.Disk) (*voc.BlockStorage, error) {
	// If a mandatory field is empty, the whole disk is empty
	if disk == nil || disk.ID == nil {
		return nil, fmt.Errorf("disk is nil")
	}

	enc, err := d.blockStorageAtRestEncryption(disk)
	if err != nil {
		return nil, fmt.Errorf("could not get block storage properties for the atRestEncryption: %w", err)
	}

	return &voc.BlockStorage{
		Storage: &voc.Storage{
			Resource: &voc.Resource{
				ID:           voc.ResourceID(util.Deref(disk.ID)),
				ServiceID:    discovery.DefaultCloudServiceID,
				Name:         util.Deref(disk.Name),
				CreationTime: disk.Properties.TimeCreated.Unix(),
				Type:         []string{"BlockStorage", "Storage", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: util.Deref(disk.Location),
				},
				Labels: labels(disk.Tags),
			},
			AtRestEncryption: enc,
		},
	}, nil
}

// blockStorageAtRestEncryption takes encryption properties of an armcompute.Disk and converts it into our respective
// ontology object.
func (d *azureComputeDiscovery) blockStorageAtRestEncryption(disk *armcompute.Disk) (enc voc.IsAtRestEncryption, err error) {
	var (
		diskEncryptionSetID string
		keyUrl              string
	)

	if disk == nil {
		return enc, errors.New("disk is empty")
	}

	if disk.Properties.Encryption.Type == nil {
		return enc, errors.New("error getting atRestEncryption properties of blockStorage")
	} else if *disk.Properties.Encryption.Type == armcompute.EncryptionTypeEncryptionAtRestWithPlatformKey {
		enc = &voc.ManagedKeyEncryption{AtRestEncryption: &voc.AtRestEncryption{
			Algorithm: "AES256",
			Enabled:   true,
		}}
	} else if *disk.Properties.Encryption.Type == armcompute.EncryptionTypeEncryptionAtRestWithCustomerKey {
		diskEncryptionSetID = util.Deref(disk.Properties.Encryption.DiskEncryptionSetID)

		keyUrl, err = d.keyURL(diskEncryptionSetID)
		if err != nil {
			return nil, fmt.Errorf("could not get keyVaultID: %w", err)
		}

		enc = &voc.CustomerKeyEncryption{
			AtRestEncryption: &voc.AtRestEncryption{
				Algorithm: "", // TODO(garuppel): TBD
				Enabled:   true,
			},
			KeyUrl: keyUrl,
		}
	}

	return enc, nil
}

func (d *azureComputeDiscovery) keyURL(diskEncryptionSetID string) (string, error) {
	if diskEncryptionSetID == "" {
		return "", ErrMissingDiskEncryptionSetID
	}

	// Create Key Vault client
	client, err := armcompute.NewDiskEncryptionSetsClient(util.Deref(d.sub.SubscriptionID), d.cred, &d.clientOptions)
	if err != nil {
		err = fmt.Errorf("could not get new key vault client: %w", err)
		return "", err
	}

	// Get disk encryption set
	kv, err := client.Get(context.TODO(), resourceGroupName(diskEncryptionSetID), diskEncryptionSetName(diskEncryptionSetID), &armcompute.DiskEncryptionSetsClientGetOptions{})
	if err != nil {
		err = fmt.Errorf("could not get key vault: %w", err)
		return "", err
	}

	keyURL := kv.DiskEncryptionSet.Properties.ActiveKey.KeyURL

	if keyURL == nil {
		return "", fmt.Errorf("could not get keyURL")
	}

	return util.Deref(keyURL), nil
}
