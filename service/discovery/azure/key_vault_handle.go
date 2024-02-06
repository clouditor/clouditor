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
	"time"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
)

// TODO(lebogg): Add handle functions used by the discover part here, e.g. `handleKeyVault`

// TODO(lebogg): Test handleKeyVault
func (d *azureDiscovery) handleKeyVault(kv *armkeyvault.Vault) (*voc.KeyVault, error) {
	var createdAt *time.Time

	if kv.SystemData != nil {
		createdAt = kv.SystemData.CreatedAt
	}

	return &voc.KeyVault{
		Resource: discovery.NewResource(d,
			voc.ResourceID(util.Deref(kv.ID)),
			util.Deref(kv.Name),
			createdAt,
			voc.GeoLocation{
				Region: util.Deref(kv.Location),
			},
			labels(kv.Tags),
			resourceGroupID(kv.ID),
			voc.KeyVaultType,
			kv),
		IsActive:     false,              // TODO(lebogg): Add `isActive`
		Keys:         []voc.ResourceID{}, // Will be added later when we retrieve the single keys
		PublicAccess: getPublicAccess(kv),
	}, nil
}
