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

//go:build exclude

package azure

import (
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dataprotection/armdataprotection"
)

func (d *azureDiscovery) handleInstances(vault *armdataprotection.BackupVaultResource, instance *armdataprotection.BackupInstanceResource) (resource voc.IsCloudResource, err error) {
	if vault == nil || instance == nil {
		return nil, ErrVaultInstanceIsEmpty
	}

	raw, err := voc.ToStringInterface([]interface{}{instance, vault})
	if err != nil {
		log.Error(err)
	}

	if *instance.Properties.DataSourceInfo.DatasourceType == "Microsoft.Storage/storageAccounts/blobServices" {
		resource = &voc.ObjectStorage{
			Storage: &voc.Storage{
				Resource: &voc.Resource{
					ID:           voc.ResourceID(*instance.ID),
					Name:         *instance.Name,
					CreationTime: 0,
					GeoLocation: voc.GeoLocation{
						Region: *vault.Location,
					},
					Labels:    nil,
					ServiceID: d.csID,
					Type:      voc.ObjectStorageType,
					Parent:    resourceGroupID(instance.ID),
					Raw:       raw,
				},
			},
		}
	} else if *instance.Properties.DataSourceInfo.DatasourceType == "Microsoft.Compute/disks" {
		resource = &voc.BlockStorage{
			Storage: &voc.Storage{
				Resource: &voc.Resource{
					ID:           voc.ResourceID(*instance.ID),
					Name:         *instance.Name,
					ServiceID:    d.csID,
					CreationTime: 0,
					Type:         voc.BlockStorageType,
					GeoLocation: voc.GeoLocation{
						Region: *vault.Location,
					},
					Labels: nil,
					Parent: resourceGroupID(instance.ID),
					Raw:    raw,
				},
			},
		}
	}

	return
}
