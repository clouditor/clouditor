// Copyright 2024 Fraunhofer AISEC
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
	"fmt"
	"strings"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v3"
)

func (d *azureDiscovery) handleVirtualMachines(vm *armcompute.VirtualMachine) (voc.IsCompute, error) {
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

func (d *azureDiscovery) handleBlockStorage(disk *armcompute.Disk) (*voc.BlockStorage, error) {
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

func (d *azureDiscovery) handleFunction(function *armappservice.Site, config armappservice.WebAppsClientGetConfigurationResponse) voc.IsCompute {
	var (
		runtimeLanguage string
		runtimeVersion  string
	)

	// If a mandatory field is empty, the whole function is empty
	if function == nil || config == (armappservice.WebAppsClientGetConfigurationResponse{}) {
		log.Error("input parameter empty")
		return nil
	}

	if *function.Kind == "functionapp,linux" { // Linux function
		runtimeLanguage, runtimeVersion = runtimeInfo(util.Deref(function.Properties.SiteConfig.LinuxFxVersion))
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
			NetworkInterfaces: getVirtualNetworkSubnetId(function), // Add the Virtual Network Subnet ID
			ResourceLogging:   d.getResourceLoggingWebApps(function),
		},
		HttpEndpoint: &voc.HttpEndpoint{
			TransportEncryption: getTransportEncryption(function.Properties, config),
		},
		RuntimeLanguage: runtimeLanguage,
		RuntimeVersion:  runtimeVersion,
		PublicAccess:    getPublicAccessStatus(function),
		Redundancy:      getRedundancy(function),
	}
}

func (d *azureDiscovery) handleWebApp(webApp *armappservice.Site, config armappservice.WebAppsClientGetConfigurationResponse) voc.IsCompute {
	if webApp == nil || config == (armappservice.WebAppsClientGetConfigurationResponse{}) {
		log.Error("input parameter empty")
		return nil
	}

	return &voc.WebApp{
		Compute: &voc.Compute{
			Resource: discovery.NewResource(d,
				voc.ResourceID(util.Deref(webApp.ID)),
				util.Deref(webApp.Name),
				nil, // Only the last modified time is available.
				voc.GeoLocation{
					Region: util.Deref(webApp.Location),
				},
				labels(webApp.Tags),
				resourceGroupID(webApp.ID),
				voc.WebAppType,
				webApp,
				config,
			),
			NetworkInterfaces: getVirtualNetworkSubnetId(webApp), // Add the Virtual Network Subnet ID
			ResourceLogging:   d.getResourceLoggingWebApps(webApp),
		},
		HttpEndpoint: &voc.HttpEndpoint{
			TransportEncryption: getTransportEncryption(webApp.Properties, config),
		},
		PublicAccess: getPublicAccessStatus(webApp),
		Redundancy:   getRedundancy(webApp),
	}
}
