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
	"clouditor.io/clouditor/voc"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
)

// TODO(lebogg): Add main discover part here, e.g. `discoverKeyVaults` or `discoverKeys`

// TODO(lebogg): Add keys, secrets and certificates
// discoverKeyVaults discovers all key vaults including their associated keys, secrets and certificates
func (d *azureDiscovery) discoverKeyVaults() (list []voc.IsCloudResource, err error) {
	// TODO(lebogg): Implement
	// Initialize key vault client
	if err = d.initKeyVaultsClient(); err != nil {
		return nil, fmt.Errorf("could not initialize key vault client: %v", err)
	}

	// Initialize keys and other clients
	// TODO

	err = listPager(d,
		d.clients.keyVaultClient.NewListPager,
		d.clients.keyVaultClient.NewListByResourceGroupPager,
		func(res armkeyvault.VaultsClientListResponse) (vaults []*armkeyvault.Vault) {
			// Here, the `res.Value` list's type Resource instead of Vault -> Convert it
			vaults = make([]*armkeyvault.Vault, len(res.Value))
			for i, r := range res.Value {
				vaults[i] = &armkeyvault.Vault{
					Location: r.Location,
					Tags:     r.Tags,
					ID:       r.ID,
					Name:     r.Name,
					Type:     r.Type,
				}
			}
			return
		},
		func(res armkeyvault.VaultsClientListByResourceGroupResponse) []*armkeyvault.Vault {
			return res.Value
		},
		func(kv *armkeyvault.Vault) error {
			keyVault, err := d.handleKeyVault(kv)
			if err != nil {
				return fmt.Errorf("could not handle key vault: %w", err)
			}

			log.Infof("Adding key vault '%s'", keyVault.GetID())
			list = append(list, keyVault)

			return nil
		})
	if err != nil {
		list = nil
		return
	}
	return
}
