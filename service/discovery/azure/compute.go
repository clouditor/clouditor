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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dataprotection/armdataprotection"

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

	// initialize backup policies client
	if err := d.initBackupPoliciesClient(); err != nil {
		return nil, err
	}

	// initialize backup vaults client
	if err := d.initBackupVaultsClient(); err != nil {
		return nil, err
	}

	// initialize backup instances client
	if err := d.initBackupInstancesClient(); err != nil {
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

	// initialize functions client
	if err := d.initWebAppsClient(); err != nil {
		return nil, err
	}

	// initialize farms client
	if err := d.initAppServiceFarmsClient(); err != nil {
		return nil, err
	}

	// List functions
	err := listPager(d.azureDiscovery,
		d.clients.sitesClient.NewListPager,
		d.clients.sitesClient.NewListByResourceGroupPager,
		func(res armappservice.WebAppsClientListResponse) []*armappservice.Site {
			return res.Value
		},
		func(res armappservice.WebAppsClientListByResourceGroupResponse) []*armappservice.Site {
			return res.Value
		},
		func(site *armappservice.Site) error {
			var r voc.IsCompute

			// Get configuration
			config, err := d.clients.sitesClient.GetConfiguration(context.Background(), *site.Properties.ResourceGroup, *site.Name, &armappservice.WebAppsClientGetConfigurationOptions{})
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

func (d *azureComputeDiscovery) handleFunction(function *armappservice.Site, config armappservice.WebAppsClientGetConfigurationResponse) voc.IsCompute {
	var (
		runtimeLanguage     string
		runtimeVersion      string
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
	resourceLogging := d.getResourceLoggingWebApp(function)

	return &voc.Function{
		Compute: &voc.Compute{
			Resource: discovery.NewResource(d,
				voc.ResourceID(resourceID(function.ID)),
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
			ResourceLogging:   resourceLogging,
		},
		HttpEndpoint: &voc.HttpEndpoint{
			TransportEncryption: getTransportEncryption(function.Properties, config),
		},
		RuntimeLanguage: runtimeLanguage,
		RuntimeVersion:  runtimeVersion,
		PublicAccess:    publicNetworkAccess,
		Redundancy:      d.getRedundancy(function),
	}
}

func (d *azureComputeDiscovery) handleWebApp(webApp *armappservice.Site, config armappservice.WebAppsClientGetConfigurationResponse) voc.IsCompute {
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
		ni = []voc.ResourceID{voc.ResourceID(resourceID(webApp.Properties.VirtualNetworkSubnetID))}
	}

	// Check if resource is public available
	if util.Deref(webApp.Properties.PublicNetworkAccess) == "Enabled" {
		publicNetworkAccess = true
	}

	resourceLogging := d.getResourceLoggingWebApp(webApp)

	// Check if secrets are used and if so, add them to the 'secretUsage' dictionary
	// Using NewGetAppSettingsKeyVaultReferencesPager would be optimal but is bugged, see https://github.com/Azure/azure-sdk-for-go/issues/14509
	settings, err := d.clients.sitesClient.ListApplicationSettings(context.TODO(), util.Deref(d.rg),
		util.Deref(webApp.Name), &armappservice.WebAppsClientListApplicationSettingsOptions{})
	if err != nil {
		// Maybe returning error better here
		log.Warnf("Could not get application settings: %v", err)
	}
	addSecretUsages(webApp.ID, settings)

	return &voc.WebApp{
		Compute: &voc.Compute{
			Resource: discovery.NewResource(d,
				voc.ResourceID(resourceID(webApp.ID)),
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
				config,
			),
			NetworkInterfaces: ni, // Add the Virtual Network Subnet ID
			ResourceLogging:   resourceLogging,
		},
		HttpEndpoint: &voc.HttpEndpoint{
			TransportEncryption: getTransportEncryption(webApp.Properties, config),
		},
		PublicAccess: publicNetworkAccess,
		Redundancy:   d.getRedundancy(webApp),
	}
}

// getRedundancy returns for a given web app/function the redundancy in the voc format
func (d *azureComputeDiscovery) getRedundancy(app *armappservice.Site) (r *voc.Redundancy) {
	r = &voc.Redundancy{}

	if app.Properties == nil || app.Properties.ServerFarmID == nil {
		log.Errorf("Could not look at properties or the Server Farm ID because one of them is empty")
		return
	}
	planName := getAppServicePlanName(app.Properties.ServerFarmID)
	farm, err := d.clients.plansClient.Get(context.TODO(), util.Deref(d.rg), planName, &armappservice.PlansClientGetOptions{})
	if err != nil {
		log.Errorf("Could not get App Service Farm '%s', zone redundancy of web app '%s' is assumed to be false: %v",
			util.Deref(app.Properties.ServerFarmID), util.Deref(app.Name), err)
		return
	}
	r.Zone = util.Deref(farm.Properties.ZoneRedundant)
	return
}

// getAppServicePlanName returns the name for the given ID of an app service plan (formerly farm). If it is wrongly
// formatted, the empty string will be returned
// A plan/farm id has the following form:
// "/subscriptions/{subscriptionID}/resourceGroups/{groupName}/providers/Microsoft.Web/serverfarms/{appServicePlanName}"
func getAppServicePlanName(id *string) (farmName string) {
	var ok bool
	_, farmName, ok = strings.Cut(util.Deref(id),
		"/Microsoft.Web/serverfarms/")
	if !ok {
		log.Errorf("Could not cut the ID '%s' correctly. Probably it is not formatted correctly", util.Deref(id))
		farmName = ""
		return
	}
	return
}

// We really need both parameters since config is indeed more precise but it does not include the `httpsOnly` property
func getTransportEncryption(siteProperties *armappservice.SiteProperties, config armappservice.WebAppsClientGetConfigurationResponse) (enc *voc.TransportEncryption) {
	var (
		tlsVersion string
	)

	switch util.Deref(config.Properties.MinTLSVersion) {
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
			Enforced:   util.Deref(siteProperties.HTTPSOnly),
			TlsVersion: tlsVersion,
			Algorithm:  string(util.Deref(config.Properties.MinTLSCipherSuite)),
			Enabled:    true,
		}
	} else {
		enc = &voc.TransportEncryption{
			Enforced:  util.Deref(siteProperties.HTTPSOnly),
			Enabled:   false,
			Algorithm: string(util.Deref(config.Properties.MinTLSCipherSuite)),
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

	// initialize virtual machines client
	if err := d.initVirtualMachinesClient(); err != nil {
		return nil, err
	}

	// List all VMs
	err := listPager(d.azureDiscovery,
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
				voc.ResourceID(resourceID(vm.ID)),
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
			r.NetworkInterfaces = append(r.NetworkInterfaces, voc.ResourceID(resourceID(networkInterfaces.ID)))
		}
	}

	// Reference to blockstorage
	if vm.Properties.StorageProfile != nil && vm.Properties.StorageProfile.OSDisk != nil && vm.Properties.StorageProfile.OSDisk.ManagedDisk != nil {
		r.BlockStorage = append(r.BlockStorage, voc.ResourceID(resourceID(vm.Properties.StorageProfile.OSDisk.ManagedDisk.ID)))
	}

	if vm.Properties.StorageProfile != nil && vm.Properties.StorageProfile.DataDisks != nil {
		for _, blockstorage := range vm.Properties.StorageProfile.DataDisks {
			r.BlockStorage = append(r.BlockStorage, voc.ResourceID(resourceID(blockstorage.ManagedDisk.ID)))
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

	// initialize block storages client
	if err := d.initBlockStoragesClient(); err != nil {
		return nil, err
	}

	// List all disks
	err := listPager(d.azureDiscovery,
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

func (d *azureComputeDiscovery) handleBlockStorage(disk *armcompute.Disk) (*voc.BlockStorage, error) {
	var (
		rawKeyUrl           *armcompute.DiskEncryptionSet
		backups             []*voc.Backup
		publicNetworkAccess = false
	)

	// If a mandatory field is empty, the whole disk is empty
	if disk == nil || disk.ID == nil {
		return nil, fmt.Errorf("disk is nil")
	}

	enc, rawKeyUrl, err := d.blockStorageAtRestEncryption(disk)
	if err != nil {
		return nil, fmt.Errorf("could not get block storage properties for the atRestEncryption: %w", err)
	}

	// Check if resource is public available
	if util.Deref(disk.Properties.PublicNetworkAccess) == "Enabled" {
		publicNetworkAccess = true
	}

	// Get voc.Backup
	if d.backupMap[DataSourceTypeDisc] != nil && d.backupMap[DataSourceTypeDisc].backup[util.Deref(disk.ID)] != nil {
		backups = d.backupMap[DataSourceTypeDisc].backup[util.Deref(disk.ID)]
	}
	backups = backupsEmptyCheck(backups)

	return &voc.BlockStorage{
		Storage: &voc.Storage{
			Resource: discovery.NewResource(d,
				voc.ResourceID(resourceID(disk.ID)),
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
			// Todo(lebogg): Add tests
			Redundancy: getDiskRedundancy(disk),
		},
		PublicAccess: publicNetworkAccess,
	}, nil
}

// Todo(lebogg): Add tests
// getDiskRedundancy maps Azure SKUs to the redundancy model provided by the ontology. Geo redundancy is currently
// not supported in Azure. Therefore, an Azure disk can be either locally redundant, zone redundant or not redundant at
// all.
func getDiskRedundancy(disk *armcompute.Disk) (r *voc.Redundancy) {
	r = &voc.Redundancy{}
	// If SKU is nil, no redundancy is set. Therefore, we return false for all redundant options in the ontology
	if disk.SKU == nil {
		return
	}
	// The SKU (formerly called account types) name tells us which redundancy model is used. We compare all account type
	// constants with the given name.
	name := util.Deref(disk.SKU.Name)
	switch name {
	// Check constants indicating local redundancy
	case armcompute.DiskStorageAccountTypesStandardLRS, armcompute.DiskStorageAccountTypesStandardSSDLRS,
		armcompute.DiskStorageAccountTypesPremiumLRS, armcompute.DiskStorageAccountTypesPremiumV2LRS,
		armcompute.DiskStorageAccountTypesUltraSSDLRS:
		r.Local = true
	// Check constants indicating zone redundancy
	case armcompute.DiskStorageAccountTypesStandardSSDZRS, armcompute.DiskStorageAccountTypesPremiumZRS:
		r.Zone = true
	// When there are new names in the future we will probably miss it. Print out a warning if there is a name we don't
	// consider so far.
	default:
		log.Warnf("Unknown redundancy model (via SKU) for disk '%s': '%s'. Probably, we should add it.",
			util.Deref(disk.Name), name)
		// consideredAccountTypes shows how many account types (SKUs) we consider so far. It has to be a "magic" number.
		consideredAccountTypes := 7
		log.Warnf("Currently there are %d different SKU types/name. We consider %d so far",
			len(armcompute.PossibleDiskStorageAccountTypesValues()), consideredAccountTypes)
	}
	return
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

	if err := d.initDiskEncryptonSetClient(); err != nil {
		return "", nil, err
	}

	// Get disk encryption set
	kv, err := d.clients.diskEncSetClient.Get(context.TODO(), resourceGroupName(diskEncryptionSetID), diskEncryptionSetName(diskEncryptionSetID), &armcompute.DiskEncryptionSetsClientGetOptions{})
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

// initWebAppsClient creates the client if not already exists
func (d *azureComputeDiscovery) initWebAppsClient() (err error) {
	d.clients.sitesClient, err = initClient(d.clients.sitesClient, d.azureDiscovery, armappservice.NewWebAppsClient)
	return
}

// initWebAppsClient creates the client if not already exists
func (d *azureComputeDiscovery) initAppServiceFarmsClient() (err error) {
	d.clients.plansClient, err = initClient(d.clients.plansClient, d.azureDiscovery, armappservice.NewPlansClient)
	return
}

// initVirtualMachinesClient creates the client if not already exists
func (d *azureComputeDiscovery) initVirtualMachinesClient() (err error) {
	d.clients.virtualMachinesClient, err = initClient(d.clients.virtualMachinesClient, d.azureDiscovery, armcompute.NewVirtualMachinesClient)
	return
}

// initBlockStoragesClient creates the client if not already exists
func (d *azureComputeDiscovery) initBlockStoragesClient() (err error) {
	d.clients.blockStorageClient, err = initClient(d.clients.blockStorageClient, d.azureDiscovery, armcompute.NewDisksClient)
	return
}

// initBlockStoragesClient creates the client if not already exists
func (d *azureComputeDiscovery) initDiskEncryptonSetClient() (err error) {
	d.clients.diskEncSetClient, err = initClient(d.clients.diskEncSetClient, d.azureDiscovery, armcompute.NewDiskEncryptionSetsClient)
	return
}

// initBackupPoliciesClient creates the client if not already exists
func (d *azureComputeDiscovery) initBackupPoliciesClient() (err error) {
	d.clients.backupPoliciesClient, err = initClient(d.clients.backupPoliciesClient, d.azureDiscovery, armdataprotection.NewBackupPoliciesClient)

	return
}

// initBackupVaultsClient creates the client if not already exists
func (d *azureComputeDiscovery) initBackupVaultsClient() (err error) {
	d.clients.backupVaultClient, err = initClient(d.clients.backupVaultClient, d.azureDiscovery, armdataprotection.NewBackupVaultsClient)

	return
}

// initBackupInstancesClient creates the client if not already exists
func (d *azureComputeDiscovery) initBackupInstancesClient() (err error) {
	d.clients.backupInstancesClient, err = initClient(d.clients.backupInstancesClient, d.azureDiscovery, armdataprotection.NewBackupInstancesClient)

	return
}

// getResourceLoggingWebApp determines if logging is activated for given web app by checking the respective app setting
func (d *azureComputeDiscovery) getResourceLoggingWebApp(site *armappservice.Site) (rl *voc.ResourceLogging) {
	rl = &voc.ResourceLogging{Logging: &voc.Logging{}}

	appSettings, err := d.clients.sitesClient.ListApplicationSettings(context.Background(),
		*site.Properties.ResourceGroup, *site.Name, &armappservice.WebAppsClientListApplicationSettingsOptions{})
	if err != nil {
		log.Errorf("could not get application settings for '%s': %v", util.Deref(site.Name), err)
		return
	}

	if appSettings.Properties["APPLICATIONINSIGHTS_CONNECTION_STRING"] != nil {
		rl.Enabled = true
		// TODO: Get id of logging service and add it (currently not possible via app settings): rl.LoggingService

	}

	return
}

// addSecretUsages checks if secrets are used in the given web app and, if so, adds them to secretUsage
func addSecretUsages(webAppID *string, settings armappservice.WebAppsClientListApplicationSettingsResponse) {
	for _, v := range settings.Properties {
		if s := util.Deref(v); strings.Contains(s, "@Microsoft.KeyVault") {
			sURI := getSecretURI(s)
			secretUsage[sURI] = append(secretUsage[sURI], util.Deref(webAppID))
		}
	}
}

// getSecretURI gets the URI of the given secret s. There can be two options how a secret attribute is stored:
// Option 1: @Microsoft.KeyVault(VaultName=myvault;SecretName=mysecret)
// Option 2: @Microsoft.KeyVault(SecretUri=https://myvault.vault.azure.net/secrets/mysecret/)
func getSecretURI(s string) (secretURI string) {
	var (
		vaultName  string
		secretName string
		ok         bool
	)

	// If s contains "VaultName" it is Option 1
	if strings.Contains(s, "VaultName") {
		s, ok = strings.CutPrefix(s, "@Microsoft.KeyVault(VaultName=")
		if !ok {
			log.Error("Could not find prefix '@Microsoft.KeyVault(VaultName=' in secret:", s)
			return ""
		}
		splits := strings.Split(s, ";")
		if len(splits) < 2 {
			log.Error("Splitting ';' should give at least two strings but didn't:", s)
			return ""
		}
		vaultName = splits[0]
		s = splits[1]
		s, ok = strings.CutPrefix(s, "SecretName=")
		if !ok {
			log.Error("Could not find prefix 'SecretName=' in:", s)
			return ""
		}
		s, ok = strings.CutSuffix(s, ")")
		if !ok {
			log.Error("Could not find suffix ')' in:", s)
			return ""
		}
		secretName = s
		secretURI = "https://" + vaultName + ".vault.azure.net/secrets/" + secretName
	} else { // Option 2
		s, ok := strings.CutPrefix(s, "@Microsoft.KeyVault(SecretUri=")
		if !ok {
			log.Error("Could not find prefix '@Microsoft.KeyVault(SecretUri=' in secret:", s)
			return ""
		}
		s, ok = strings.CutSuffix(s, "/)")
		if !ok {
			log.Error("Could not find suffix ')' in:", s)
			return ""
		}
		secretURI = s
	}

	return
}
