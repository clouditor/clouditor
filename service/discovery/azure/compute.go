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
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dataprotection/armdataprotection/v2"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/constants"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
)

var (
	ErrEmptyVirtualMachine = errors.New("virtual machine is empty")
)

type azureComputeDiscovery struct {
	*azureDiscovery
	defenderProperties map[string]*defenderProperties
}

func NewAzureComputeDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &azureComputeDiscovery{
		&azureDiscovery{
			discovererComponent: ComputeComponent,
			csID:                discovery.DefaultCloudServiceID,
			backupMap:           make(map[string]*backup),
		},
		make(map[string]*defenderProperties),
	}

	// Apply options
	for _, opt := range opts {
		opt(d.azureDiscovery)
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

	// initialize dataprotection factory client
	if err := d.initClientDataprotectionFactory(); err != nil {
		return nil, err
	}

	// Discover backup vaults
	err = d.azureDiscovery.discoverBackupVaults()
	if err != nil {
		log.Errorf("could not discover backup vaults: %v", err)
	}

	log.Info("Discover Azure block storage")
	// Discover block storage
	storage, err := d.discoverBlockStorages()
	if err != nil {
		return nil, fmt.Errorf("could not discover block storage: %w", err)
	}
	list = append(list, storage...)

	// Add backup block storages
	if d.backupMap[DataSourceTypeDisc] != nil && d.backupMap[DataSourceTypeDisc].backupStorages != nil {
		list = append(list, d.backupMap[DataSourceTypeDisc].backupStorages...)
	}

	log.Info("Discover Azure compute resources")
	// Discover virtual machines
	virtualMachines, err := d.discoverVirtualMachines()
	if err != nil {
		return nil, fmt.Errorf("could not discover virtual machines: %w", err)
	}
	list = append(list, virtualMachines...)

	// Discover functions and web apps
	resources, err := d.discoverFunctionsWebApps()
	if err != nil {
		return nil, fmt.Errorf("could not discover functions: %w", err)
	}
	// if resources != nil {
	list = append(list, resources...)
	// }

	return
}

// Discover functions and web apps
func (d *azureComputeDiscovery) discoverFunctionsWebApps() ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	// initialize appservice factory client
	if err := d.initClientAppserviceFactory(); err != nil {
		return nil, err
	}

	// List functions
	err := listPager(d.azureDiscovery,
		d.clients.appserviceFactory.NewWebAppsClient().NewListPager,
		d.clients.appserviceFactory.NewWebAppsClient().NewListByResourceGroupPager,
		func(res armappservice.WebAppsClientListResponse) []*armappservice.Site {
			return res.Value
		},
		func(res armappservice.WebAppsClientListByResourceGroupResponse) []*armappservice.Site {
			return res.Value
		},
		func(site *armappservice.Site) error {
			var r voc.IsCompute

			// Check kind of site (see https://github.com/Azure/app-service-linux-docs/blob/master/Things_You_Should_Know/kind_property.md)
			switch *site.Kind {
			case "app": // Windows Web App
				r = d.handleWebApp(site)
			case "app,linux": // Linux Web app
				r = d.handleWebApp(site)
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
				r = d.handleFunction(site)
			case "functionapp,linux": // Linux Consumption Function app
				r = d.handleFunction(site)
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
				log.Infof("Adding function %+v", r)
				list = append(list, r)
			}

			return nil
		})
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (d *azureComputeDiscovery) handleFunction(function *armappservice.Site) voc.IsCompute {
	var (
		runtimeLanguage     string
		runtimeVersion      string
		config              armappservice.WebAppsClientGetConfigurationResponse
		err                 error
		publicNetworkAccess = false
	)

	// If a mandatory field is empty, the whole function is empty
	if function == nil {
		return nil
	}

	if util.Deref(function.Properties.PublicNetworkAccess) == "Enabled" {
		publicNetworkAccess = true
	}

	if *function.Kind == "functionapp,linux" { // Linux function
		runtimeLanguage, runtimeVersion = runtimeInfo(*function.Properties.SiteConfig.LinuxFxVersion)
	} else if *function.Kind == "functionapp" { // Windows function, we need to get also the config information
		// Get site config
		config, err = d.clients.appserviceFactory.NewWebAppsClient().GetConfiguration(context.Background(), *function.Properties.ResourceGroup, *function.Name, &armappservice.WebAppsClientGetConfigurationOptions{})
		if err != nil {
			log.Errorf("error getting site config: %v", err)
		}

		// Check all runtime versions to get the used runtime language and runtime version
		if util.Deref(config.Properties.JavaVersion) != "" {
			runtimeLanguage = "Java"
			runtimeVersion = *config.Properties.JavaVersion
		} else if util.Deref(config.Properties.NodeVersion) != "" {
			runtimeLanguage = "Node.js"
			runtimeVersion = *config.Properties.NodeVersion
		} else if util.Deref(config.Properties.PowerShellVersion) != "" {
			runtimeLanguage = "PowerShell"
			runtimeVersion = *config.Properties.PowerShellVersion
		} else if util.Deref(config.Properties.PhpVersion) != "" {
			runtimeLanguage = "PHP"
			runtimeVersion = *config.Properties.PhpVersion
		} else if util.Deref(config.Properties.PythonVersion) != "" {
			runtimeLanguage = "Python"
			runtimeVersion = *config.Properties.PythonVersion
		} else if util.Deref(config.Properties.JavaContainer) != "" {
			runtimeLanguage = "JavaContainer"
			runtimeVersion = *config.Properties.JavaContainer
		} else if util.Deref(config.Properties.NetFrameworkVersion) != "" {
			runtimeLanguage = ".NET"
			runtimeVersion = *config.Properties.NetFrameworkVersion
		}
	}

	return &voc.Function{
		Compute: &voc.Compute{
			Resource: discovery.NewResource(d,
				voc.ResourceID(util.Deref(function.ID)),
				util.Deref(function.Name),
				// No creation time available
				nil,
				voc.GeoLocation{
					Region: util.Deref(function.Location),
				},
				labels(function.Tags),
				resourceGroupID(function.ID),
				voc.FunctionType,
				function,
				config,
			),
			NetworkInterfaces: []voc.ResourceID{},
		},
		HttpEndpoint: &voc.HttpEndpoint{
			TransportEncryption: getTransportEncryption(function.Properties),
		},
		RuntimeLanguage:     runtimeLanguage,
		RuntimeVersion:      runtimeVersion,
		PublicNetworkAccess: publicNetworkAccess,
	}
}

func (d *azureComputeDiscovery) handleWebApp(webApp *armappservice.Site) voc.IsCompute {
	var (
		ni                  []voc.ResourceID
		publicNetworkAccess = false
	)

	// If a mandatory field is empty, the whole function is empty
	if webApp == nil {
		return nil
	}

	// Get virtual network subnet ID
	if webApp.Properties.VirtualNetworkSubnetID != nil {
		ni = []voc.ResourceID{voc.ResourceID(*webApp.Properties.VirtualNetworkSubnetID)}
	}

	// Check if resource is public available
	if util.Deref(webApp.Properties.PublicNetworkAccess) == "Enabled" {
		publicNetworkAccess = true
	}

	return &voc.WebApp{
		Compute: &voc.Compute{
			Resource: discovery.NewResource(d,
				voc.ResourceID(util.Deref(webApp.ID)),
				util.Deref(webApp.Name),
				// No creation time available
				nil, // Only the last modified time is available
				voc.GeoLocation{
					Region: util.Deref(webApp.Location),
				},
				labels(webApp.Tags),
				resourceGroupID(webApp.ID),
				voc.WebAppType,
				webApp,
				// config,
			),
			NetworkInterfaces: ni, // Add the Virtual Network Subnet ID
		},
		HttpEndpoint: &voc.HttpEndpoint{
			TransportEncryption: getTransportEncryption(webApp.Properties),
		},
		PublicNetworkAccess: publicNetworkAccess,
	}
}

func getTransportEncryption(siteProps *armappservice.SiteProperties) (enc *voc.TransportEncryption) {
	var (
		tlsVersion string
	)

	switch util.Deref(siteProps.SiteConfig.MinTLSVersion) {
	case armappservice.SupportedTLSVersionsOne2:
		tlsVersion = constants.TLS1_2
	case armappservice.SupportedTLSVersionsOne1:
		tlsVersion = constants.TLS1_1
	case armappservice.SupportedTLSVersionsOne0:
		tlsVersion = constants.TLS1_0

	}
	// Check TLS version
	if tlsVersion != "" {
		enc = &voc.TransportEncryption{
			Enforced:   util.Deref(siteProps.HTTPSOnly),
			TlsVersion: tlsVersion,
			Algorithm:  string(util.Deref(siteProps.SiteConfig.MinTLSCipherSuite)),
			Enabled:    true,
		}
	} else {
		enc = &voc.TransportEncryption{
			Enforced:  util.Deref(siteProps.HTTPSOnly),
			Enabled:   false,
			Algorithm: string(util.Deref(siteProps.SiteConfig.MinTLSCipherSuite)),
		}
	}

	return
}

// runtimeInfo returns the runtime language and version
func runtimeInfo(runtime string) (runtimeLanguage string, runtimeVersion string) {
	if runtime == "" || !strings.Contains(runtime, "|") {
		return "", ""
	}
	split := strings.Split(runtime, "|")
	runtimeLanguage = split[0]
	runtimeVersion = split[1]

	return
}

// Discover virtual machines
func (d *azureComputeDiscovery) discoverVirtualMachines() ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	// initialize compute factory client
	if err := d.initClientComputeFactory(); err != nil {
		return nil, err
	}

	// List all VMs
	err := listPager(d.azureDiscovery,
		d.clients.computeFactory.NewVirtualMachinesClient().NewListAllPager,
		d.clients.computeFactory.NewVirtualMachinesClient().NewListPager,
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

func (d *azureComputeDiscovery) handleVirtualMachines(vm *armcompute.VirtualMachine) (voc.IsCompute, error) {
	var (
		bootLogging              = []voc.ResourceID{}
		osLoggingEnabled         bool
		autoUpdates              *voc.AutomaticUpdates
		monitoringLogDataEnabled bool
		securityAlertsEnabled    bool
	)

	// If a mandatory field is empty, the whole disk is empty
	if vm == nil || vm.ID == nil {
		return nil, ErrEmptyVirtualMachine
	}

	if bootLogOutput(vm) != "" {
		bootLogging = []voc.ResourceID{voc.ResourceID(bootLogOutput(vm))}
	}

	autoUpdates = automaticUpdates(vm)

	if d.defenderProperties[DefenderVirtualMachineType] != nil {
		monitoringLogDataEnabled = d.defenderProperties[DefenderVirtualMachineType].monitoringLogDataEnabled
		securityAlertsEnabled = d.defenderProperties[DefenderVirtualMachineType].securityAlertsEnabled
	}

	// Check extensions
	for _, extension := range vm.Resources {
		// Azure Monitor Agent (AMA) collects monitoring data from the guest operating system of Azure and hybrid virtual machines and delivers it to Azure Monitor for use (https://learn.microsoft.com/en-us/azure/azure-monitor/agents/agents-overview). The extension names are
		// * OMSAgentForLinux for Linux VMs and (legacy agent)
		// * MicrosoftMonitoringAgent for Windows VMs (legacy agent)
		// * AzureMonitoringWindowsAgent (new agent)
		// * AzureMonitoringLinuxAgent (new agent)
		if strings.Contains(*extension.ID, "OmsAgentForLinux") || strings.Contains(*extension.ID, "MicrosoftMonitoringAgent") || strings.Contains(*extension.ID, "AzureMonitoringWindowsAgent") || strings.Contains(*extension.ID, "AzureMonitoringLinuxAgent") {
			osLoggingEnabled = true
		}
	}

	r := &voc.VirtualMachine{
		Compute: &voc.Compute{
			Resource: discovery.NewResource(d,
				voc.ResourceID(util.Deref(vm.ID)),
				util.Deref(vm.Name),
				vm.Properties.TimeCreated,
				voc.GeoLocation{
					Region: util.Deref(vm.Location),
				},
				labels(vm.Tags),
				resourceGroupID(vm.ID),
				voc.VirtualMachineType,
				vm,
			),
			NetworkInterfaces: []voc.ResourceID{},
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
				MonitoringLogDataEnabled: monitoringLogDataEnabled,
				SecurityAlertsEnabled:    securityAlertsEnabled,
			},
		},
		OsLogging: &voc.OSLogging{
			Logging: &voc.Logging{
				Enabled:         osLoggingEnabled,
				RetentionPeriod: 0,
				LoggingService:  []voc.ResourceID{}, // TODO(all): TBD
				Auditing: &voc.Auditing{
					SecurityFeature: &voc.SecurityFeature{},
				},
				MonitoringLogDataEnabled: monitoringLogDataEnabled,
				SecurityAlertsEnabled:    monitoringLogDataEnabled,
			},
		},
		ActivityLogging: &voc.ActivityLogging{
			Logging: &voc.Logging{
				Enabled:         true, // is always enabled
				RetentionPeriod: RetentionPeriod90Days,
				LoggingService:  []voc.ResourceID{}, // TODO(all): TBD
			},
		},
		AutomaticUpdates: autoUpdates,
	}

	// Reference to networkInterfaces
	if vm.Properties.NetworkProfile != nil {
		for _, networkInterfaces := range vm.Properties.NetworkProfile.NetworkInterfaces {
			r.NetworkInterfaces = append(r.NetworkInterfaces, voc.ResourceID(util.Deref(networkInterfaces.ID)))
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

// automaticUpdates returns automaticUpdatesEnabled and automaticUpdatesInterval for a given VM.
func automaticUpdates(vm *armcompute.VirtualMachine) (automaticUpdates *voc.AutomaticUpdates) {
	automaticUpdates = &voc.AutomaticUpdates{}

	if vm == nil || vm.Properties == nil || vm.Properties.OSProfile == nil {
		return
	}

	// Check if Linux configuration is available
	if vm.Properties.OSProfile.LinuxConfiguration != nil &&
		vm.Properties.OSProfile.LinuxConfiguration.PatchSettings != nil {
		if util.Deref(vm.Properties.OSProfile.LinuxConfiguration.PatchSettings.PatchMode) == armcompute.LinuxVMGuestPatchModeAutomaticByPlatform {
			automaticUpdates.Enabled = true
			automaticUpdates.Interval = Duration30Days
			return
		}
	}

	// Check if Windows configuration is available
	if vm.Properties.OSProfile.WindowsConfiguration != nil &&
		vm.Properties.OSProfile.WindowsConfiguration.PatchSettings != nil {
		if util.Deref(vm.Properties.OSProfile.WindowsConfiguration.PatchSettings.PatchMode) == armcompute.WindowsVMGuestPatchModeAutomaticByOS && *vm.Properties.OSProfile.WindowsConfiguration.EnableAutomaticUpdates ||
			util.Deref(vm.Properties.OSProfile.WindowsConfiguration.PatchSettings.PatchMode) == armcompute.WindowsVMGuestPatchModeAutomaticByPlatform && *vm.Properties.OSProfile.WindowsConfiguration.EnableAutomaticUpdates {
			automaticUpdates.Enabled = true
			automaticUpdates.Interval = Duration30Days
			return

		} else {
			return

		}
	}

	return
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
		// TODO(oxisto): The issue here, is that this is an URL but not an ID of the object storage!
		// if vm.Properties.DiagnosticsProfile.BootDiagnostics.StorageURI != nil {
		// 	return util.Deref(vm.Properties.DiagnosticsProfile.BootDiagnostics.StorageURI)
		// }

		return ""
	}
	return ""
}

func (d *azureComputeDiscovery) discoverBlockStorages() ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	// initialize compute factory client
	if err := d.initClientComputeFactory(); err != nil {
		return nil, err
	}

	// List all disks
	err := listPager(d.azureDiscovery,
		d.clients.computeFactory.NewDisksClient().NewListPager,
		d.clients.computeFactory.NewDisksClient().NewListByResourceGroupPager,
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

func (d *azureComputeDiscovery) handleBlockStorage(disk *armcompute.Disk) (*voc.BlockStorage, error) {
	var (
		rawKeyUrl *armcompute.DiskEncryptionSet
		backups   []*voc.Backup
	)

	// If a mandatory field is empty, the whole disk is empty
	if disk == nil || disk.ID == nil {
		return nil, fmt.Errorf("disk is nil")
	}

	enc, rawKeyUrl, err := d.blockStorageAtRestEncryption(disk)
	if err != nil {
		return nil, fmt.Errorf("could not get block storage properties for the atRestEncryption: %w", err)
	}

	// Get voc.Backup
	if d.backupMap[DataSourceTypeDisc] != nil && d.backupMap[DataSourceTypeDisc].backup[util.Deref(disk.ID)] != nil {
		backups = d.backupMap[DataSourceTypeDisc].backup[util.Deref(disk.ID)]
	}
	backups = backupsEmptyCheck(backups)

	return &voc.BlockStorage{
		Storage: &voc.Storage{
			Resource: discovery.NewResource(d,
				voc.ResourceID(util.Deref(disk.ID)),
				util.Deref(disk.Name),
				disk.Properties.TimeCreated,
				voc.GeoLocation{
					Region: util.Deref(disk.Location),
				},
				labels(disk.Tags),
				resourceGroupID(disk.ID),
				voc.BlockStorageType,
				disk, rawKeyUrl,
			),
			AtRestEncryption: enc,
			Backups:          backups,
		},
	}, nil
}

// blockStorageAtRestEncryption takes encryption properties of an armcompute.Disk and converts it into our respective
// ontology object.
func (d *azureComputeDiscovery) blockStorageAtRestEncryption(disk *armcompute.Disk) (enc voc.IsAtRestEncryption, rawKeyUrl *armcompute.DiskEncryptionSet, err error) {
	var (
		diskEncryptionSetID string
		keyUrl              string
	)

	if disk == nil {
		return enc, nil, errors.New("disk is empty")
	}

	if disk.Properties.Encryption.Type == nil {
		return enc, nil, errors.New("error getting atRestEncryption properties of blockStorage")
	} else if util.Deref(disk.Properties.Encryption.Type) == armcompute.EncryptionTypeEncryptionAtRestWithPlatformKey {
		enc = &voc.ManagedKeyEncryption{AtRestEncryption: &voc.AtRestEncryption{
			Algorithm: "AES256",
			Enabled:   true,
		}}
	} else if util.Deref(disk.Properties.Encryption.Type) == armcompute.EncryptionTypeEncryptionAtRestWithCustomerKey {
		diskEncryptionSetID = util.Deref(disk.Properties.Encryption.DiskEncryptionSetID)

		keyUrl, rawKeyUrl, err = d.keyURL(diskEncryptionSetID)
		if err != nil {
			return nil, nil, fmt.Errorf("could not get keyVaultID: %w", err)
		}

		enc = &voc.CustomerKeyEncryption{
			AtRestEncryption: &voc.AtRestEncryption{
				Algorithm: "", // TODO(all): TBD
				Enabled:   true,
			},
			KeyUrl: keyUrl,
		}
	}

	return enc, rawKeyUrl, nil
}

func (d *azureComputeDiscovery) keyURL(diskEncryptionSetID string) (string, *armcompute.DiskEncryptionSet, error) {
	if diskEncryptionSetID == "" {
		return "", nil, ErrMissingDiskEncryptionSetID
	}

	if err := d.initClientComputeFactory(); err != nil {
		return "", nil, err
	}

	// Get disk encryption set
	kv, err := d.clients.computeFactory.NewDiskEncryptionSetsClient().Get(context.TODO(), resourceGroupName(diskEncryptionSetID), diskEncryptionSetName(diskEncryptionSetID), &armcompute.DiskEncryptionSetsClientGetOptions{})
	if err != nil {
		err = fmt.Errorf("could not get key vault: %w", err)
		return "", nil, err
	}

	keyURL := kv.DiskEncryptionSet.Properties.ActiveKey.KeyURL

	if keyURL == nil {
		return "", nil, fmt.Errorf("could not get keyURL")
	}

	return util.Deref(keyURL), &kv.DiskEncryptionSet, nil
}

// initClientNetworkFactory creates the client if not already exists
func (d *azureComputeDiscovery) initClientAppserviceFactory() (err error) {
	d.clients.appserviceFactory, err = initClient(d.clients.appserviceFactory, d.azureDiscovery, armappservice.NewClientFactory)

	return
}

// initClientComputeFactory creates the client if not already exists
func (d *azureComputeDiscovery) initClientComputeFactory() (err error) {
	d.clients.computeFactory, err = initClient(d.clients.computeFactory, d.azureDiscovery, armcompute.NewClientFactory)

	return
}

// initClientDataprotectionFactory creates the client if not already exists
func (d *azureComputeDiscovery) initClientDataprotectionFactory() (err error) {
	d.clients.dataprotectionFactory, err = initClient(d.clients.dataprotectionFactory, d.azureDiscovery, armdataprotection.NewClientFactory)

	return
}
