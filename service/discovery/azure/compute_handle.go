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

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/util"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v3"
)

func (d *azureDiscovery) handleVirtualMachines(vm *armcompute.VirtualMachine) (ontology.IsResource, error) {
	var (
		bootLogging              = []string{}
		osLoggingEnabled         bool
		autoUpdates              *ontology.AutomaticUpdates
		monitoringLogDataEnabled bool
		securityAlertsEnabled    bool
	)

	// If a mandatory field is empty, the whole disk is empty
	if vm == nil || vm.ID == nil {
		return nil, ErrEmptyVirtualMachine
	}

	if bootLogOutput(vm) != "" {
		bootLogging = []string{bootLogOutput(vm)}
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

	r := &ontology.VirtualMachine{
		Id:           resourceID(vm.ID),
		Name:         util.Deref(vm.Name),
		CreationTime: creationTime(vm.Properties.TimeCreated),
		GeoLocation: &ontology.GeoLocation{
			Region: util.Deref(vm.Location),
		},
		Labels:              labels(vm.Tags),
		ParentId:            resourceGroupID(vm.ID),
		Raw:                 discovery.Raw(vm),
		NetworkInterfaceIds: []string{}, // TODO(all): Discover network interface IDs
		BlockStorageIds:     []string{},
		MalwareProtection:   &ontology.MalwareProtection{},
		BootLogging: &ontology.BootLogging{
			Enabled:                  isBootDiagnosticEnabled(vm),
			LoggingServiceIds:        bootLogging,
			RetentionPeriod:          durationpb.New(0), // Currently, configuring the retention period for Managed Boot Diagnostics is not available. The logs will be overwritten after 1gb of space according to https://github.com/MicrosoftDocs/azure-docs/issues/69953
			MonitoringLogDataEnabled: monitoringLogDataEnabled,
			SecurityAlertsEnabled:    securityAlertsEnabled,
		},
		OsLogging: &ontology.OSLogging{
			Enabled:                  osLoggingEnabled,
			RetentionPeriod:          durationpb.New(0),
			LoggingServiceIds:        []string{}, // TODO(all): TBD
			MonitoringLogDataEnabled: monitoringLogDataEnabled,
			SecurityAlertsEnabled:    monitoringLogDataEnabled,
		},
		ActivityLogging: &ontology.ActivityLogging{
			Enabled:           true, // is always enabled
			RetentionPeriod:   durationpb.New(RetentionPeriod90Days),
			LoggingServiceIds: []string{}, // TODO(all): TBD
		},
		AutomaticUpdates: autoUpdates,
	}

	// Reference to networkInterfaces
	if vm.Properties.NetworkProfile != nil {
		for _, networkInterfaces := range vm.Properties.NetworkProfile.NetworkInterfaces {
			r.NetworkInterfaceIds = append(r.NetworkInterfaceIds, resourceID(networkInterfaces.ID))
		}
	}

	// Reference to blockstorage
	if vm.Properties.StorageProfile != nil && vm.Properties.StorageProfile.OSDisk != nil && vm.Properties.StorageProfile.OSDisk.ManagedDisk != nil {
		r.BlockStorageIds = append(r.BlockStorageIds, resourceID(vm.Properties.StorageProfile.OSDisk.ManagedDisk.ID))
	}

	if vm.Properties.StorageProfile != nil && vm.Properties.StorageProfile.DataDisks != nil {
		for _, blockstorage := range vm.Properties.StorageProfile.DataDisks {
			r.BlockStorageIds = append(r.BlockStorageIds, resourceID(blockstorage.ManagedDisk.ID))
		}
	}

	return r, nil
}

func (d *azureDiscovery) handleBlockStorage(disk *armcompute.Disk) (*ontology.BlockStorage, error) {
	var (
		rawKeyUrl *armcompute.DiskEncryptionSet
		backups   []*ontology.Backup
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

	return &ontology.BlockStorage{
		Id:               resourceID(disk.ID),
		Name:             util.Deref(disk.Name),
		CreationTime:     creationTime(disk.Properties.TimeCreated),
		GeoLocation:      location(disk.Location),
		Labels:           labels(disk.Tags),
		ParentId:         resourceGroupID(disk.ManagedBy),
		Raw:              discovery.Raw(disk, rawKeyUrl),
		AtRestEncryption: enc,
		Backups:          backups,
	}, nil
}

func (d *azureDiscovery) handleFunction(function *armappservice.Site, config armappservice.WebAppsClientGetConfigurationResponse) ontology.IsResource {
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

	return &ontology.Function{
		Id:           resourceID(function.ID),
		Name:         util.Deref(function.Name),
		CreationTime: nil, // No creation time available
		GeoLocation: &ontology.GeoLocation{
			Region: util.Deref(function.Location),
		},
		Labels:              labels(function.Tags),
		ParentId:            resourceGroupID(function.ID),
		Raw:                 discovery.Raw(function, config),
		NetworkInterfaceIds: getVirtualNetworkSubnetId(function), // Add the Virtual Network Subnet ID
		ResourceLogging:     d.getResourceLoggingWebApps(function),
		RuntimeLanguage:     runtimeLanguage,
		RuntimeVersion:      runtimeVersion,
		// TODO(oxisto): This is missing in the ontology
		/*HttpEndpoint: &ontology.HttpEndpoint{
			TransportEncryption: getTransportEncryption(function.Properties, config),
		},*/
		InternetAccessibleEndpoint: publicNetworkAccessStatus(function.Properties.PublicNetworkAccess),
		Redundancies:               getRedundancies(function),
	}
}

func (d *azureDiscovery) handleWebApp(webApp *armappservice.Site, config armappservice.WebAppsClientGetConfigurationResponse) ontology.IsResource {
	if webApp == nil || config == (armappservice.WebAppsClientGetConfigurationResponse{}) {
		log.Error("input parameter empty")
		return nil
	}

	return &ontology.Function{
		Id:           resourceID(webApp.ID),
		Name:         util.Deref(webApp.Name),
		CreationTime: nil, // Only the last modified time is available.
		GeoLocation: &ontology.GeoLocation{
			Region: util.Deref(webApp.Location),
		},
		Labels:              labels(webApp.Tags),
		ParentId:            resourceGroupID(webApp.ID),
		Raw:                 discovery.Raw(webApp, config),
		NetworkInterfaceIds: getVirtualNetworkSubnetId(webApp), // Add the Virtual Network Subnet ID
		ResourceLogging:     d.getResourceLoggingWebApps(webApp),
		// TODO(oxisto): This is missing in the ontology
		/*HttpEndpoint: &ontology.HttpEndpoint{
			TransportEncryption: getTransportEncryption(webApp.Properties, config),
		},*/
		InternetAccessibleEndpoint: publicNetworkAccessStatus(webApp.Properties.PublicNetworkAccess),
		Redundancies:               getRedundancies(webApp),
	}
}
