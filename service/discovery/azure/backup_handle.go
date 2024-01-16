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
