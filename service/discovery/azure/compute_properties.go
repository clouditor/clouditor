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

	"clouditor.io/clouditor/internal/constants"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v3"
)

// blockStorageAtRestEncryption takes encryption properties of an armcompute.Disk and converts it into our respective
// ontology object.
func (d *azureDiscovery) blockStorageAtRestEncryption(disk *armcompute.Disk) (enc voc.IsAtRestEncryption, rawKeyUrl *armcompute.DiskEncryptionSet, err error) {
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

func (d *azureDiscovery) keyURL(diskEncryptionSetID string) (string, *armcompute.DiskEncryptionSet, error) {
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

// diskEncryptionSetName return the disk encryption set ID's name
func diskEncryptionSetName(diskEncryptionSetID string) string {
	if diskEncryptionSetID == "" {
		return ""
	}
	splitName := strings.Split(diskEncryptionSetID, "/")
	return splitName[8]
}

// getVirtualNetworkSubnetId returns the virtual network subnet ID for webApp und function
func getVirtualNetworkSubnetId(site *armappservice.Site) []voc.ResourceID {
	var ni = []voc.ResourceID{}

	// Check if a mandatory field is empty
	if site == nil {
		return ni
	}

	// Get virtual network subnet ID
	if site.Properties.VirtualNetworkSubnetID != nil {
		ni = []voc.ResourceID{voc.ResourceID(util.Deref(site.Properties.VirtualNetworkSubnetID))}
	}

	return ni
}

// getPublicAccessStatus returns the public access status for webApp and function
func getPublicAccessStatus(site *armappservice.Site) bool {
	// Check if a mandatory field is empty
	if site == nil {
		return false
	}

	// Check if resource is public available
	if util.Deref(site.Properties.PublicNetworkAccess) == "Enabled" {
		return true
	}

	return false
}

// getResourceLoggingWebApps determines if logging is activated for given web app or function by checking the respective app setting
func (d *azureDiscovery) getResourceLoggingWebApps(site *armappservice.Site) (rl *voc.ResourceLogging) {
	rl = &voc.ResourceLogging{Logging: &voc.Logging{}}

	if site == nil {
		log.Error("given parameter is empty")
		return
	}

	appSettings, err := d.clients.webAppsClient.ListApplicationSettings(context.Background(),
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

// getRedundancy returns the redundancy status
func getRedundancy(app *armappservice.Site) *voc.Redundancy {
	r := &voc.Redundancy{}
	switch util.Deref(app.Properties.RedundancyMode) {
	case armappservice.RedundancyModeNone:
		break
	case armappservice.RedundancyModeActiveActive:
		r.Zone = true
	case armappservice.RedundancyModeFailover, armappservice.RedundancyModeGeoRedundant:
		r.Zone = true
		r.Geo = true
	}
	return r
}

// We really need both parameters since config is indeed more precise but it does not include the `httpsOnly` property
func getTransportEncryption(siteProperties *armappservice.SiteProperties, config armappservice.WebAppsClientGetConfigurationResponse) (enc *voc.TransportEncryption) {
	var (
		tlsVersion string
	)

	// Check TLS version
	switch util.Deref(config.Properties.MinTLSVersion) {
	case armappservice.SupportedTLSVersionsOne2:
		tlsVersion = constants.TLS1_2
	case armappservice.SupportedTLSVersionsOne1:
		tlsVersion = constants.TLS1_1
	case armappservice.SupportedTLSVersionsOne0:
		tlsVersion = constants.TLS1_0
	}

	// Create transportEncryption voc object
	if tlsVersion != "" {
		enc = &voc.TransportEncryption{
			Enforced:   util.Deref(siteProperties.HTTPSOnly),
			TlsVersion: tlsVersion,
			Algorithm:  string(util.Deref(config.Properties.MinTLSCipherSuite)), // MinTLSCipherSuite is a new property and currently not filled from Azure side
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

func (d *azureDiscovery) bootLogOutput(vm *armcompute.VirtualMachine) string {
	if isBootDiagnosticEnabled(vm) {
		// If storageUri is not specified while enabling boot diagnostics, managed storage will be used.
		if vm.Properties.DiagnosticsProfile.BootDiagnostics.StorageURI != nil {
			return d.getResourceId(util.Deref(vm.Properties.DiagnosticsProfile.BootDiagnostics.StorageURI), resourceGroupName(util.Deref(vm.ID)))
		}

		return ""
	}
	return ""
}

// getResourceId returns the resource ID of a given URI
func (d *azureDiscovery) getResourceId(uri string, rg string) string {
	// Check if needed values are available
	if uri == "" || rg == "" || d.sub == nil || util.Deref(d.sub.SubscriptionID) == "" {
		return ""
	}

	// Get storage account name from URI
	// Example of the given URI: "https://YYYY.blob.core.windows.net/"
	tmp := strings.Split(uri, ".")
	accountName := strings.Split(tmp[0], "/")[2]

	// return the resource ID
	// Example resource ID:/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/XXXX/providers/Microsoft.Storage/storageAccounts/YYYY
	return "/subscriptions/" + util.Deref(d.sub.SubscriptionID) + "/resourceGroups/" + rg + "/providers/Microsoft.Storage/storageAccounts/" + accountName
}
