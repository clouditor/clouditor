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

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-07-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-02-01/storage"
	"github.com/Azure/go-autorest/autorest/to"
)

type azureStorageDiscovery struct {
	azureDiscovery
}

func NewAzureStorageDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &azureStorageDiscovery{}

	for _, opt := range opts {
		if auth, ok := opt.(*authorizerOption); ok {
			d.authOption = auth
		} else {
			d.options = append(d.options, opt)
		}
	}

	return d
}

func (*azureStorageDiscovery) Name() string {
	return "Azure Storage Account"
}

func (*azureStorageDiscovery) Description() string {
	return "Discovery Azure storage accounts."
}

func (d *azureStorageDiscovery) List() (list []voc.IsCloudResource, err error) {
	// make sure, we are authorized
	if err = d.authorize(); err != nil {
		return nil, fmt.Errorf("could not authorize Azure account: %w", err)
	}

	// Discover storage accounts
	storageAccounts, err := d.discoverStorageAccounts()
	if err != nil {
		return nil, fmt.Errorf("could not discover storage accounts: %w", err)
	}
	list = append(list, storageAccounts...)

	return
}

func (d *azureStorageDiscovery) discoverStorageAccounts() ([]voc.IsCloudResource, error) {
	var storageResourcesList []voc.IsCloudResource

	client := storage.NewAccountsClient(to.String(d.sub.SubscriptionID))
	d.apply(&client.Client)

	result, err := client.List(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not list storage accounts: %w", err)
	}

	// Discover object and file storages
	accounts := result.Values()
	for _, account := range accounts {
		// Discover object storages
		objectStorages, err := d.discoverObjectStorages(&account)
		if err != nil {
			return nil, fmt.Errorf("could not handle object storages: %w", err)
		}
		log.Infof("Adding object storages %+v", objectStorages)

		// Discover file storages
		fileStorages, err := d.discoverFileStorages(&account)
		if err != nil {
			return nil, fmt.Errorf("could not handle file storages: %w", err)
		}
		log.Infof("Adding file storages %+v", fileStorages)

		storageResourcesList = append(storageResourcesList, objectStorages...)
		storageResourcesList = append(storageResourcesList, fileStorages...)

		// Create storage service for all storage account resources
		storageService, err := d.handleStorageAccount(&account, storageResourcesList)
		if err != nil {
			return nil, fmt.Errorf("could not create storage service: %w", err)
		}
		log.Infof("Adding storage account %+v", objectStorages)

		storageResourcesList = append(storageResourcesList, storageService)
	}

	// Discover block storages
	blockStorages, err := d.discoverBlockStorages()
	if err != nil {
		return nil, fmt.Errorf("could not handle block storages: %w", err)
	}
	storageResourcesList = append(storageResourcesList, blockStorages...)

	return storageResourcesList, err
}

func (d *azureStorageDiscovery) discoverBlockStorages() ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	client := compute.NewDisksClient(to.String(d.sub.SubscriptionID))
	d.apply(&client.Client)

	result, err := client.ListComplete(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not list block storages: %w", err)
	}

	for _, disk := range *result.Response().Value {
		blockStorages, err := d.handleBlockStorage(&disk)
		if err != nil {
			return nil, fmt.Errorf("could not handle block storage: %w", err)
		}
		log.Infof("Adding block storage %+v", blockStorages)

		list = append(list, blockStorages)
	}

	return list, nil
}

func (d *azureStorageDiscovery) discoverFileStorages(account *storage.Account) ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	client := storage.NewFileSharesClient(to.String(d.sub.SubscriptionID))
	d.apply(&client.Client)

	result, err := client.List(context.Background(), getResourceGroupName(*account.ID), *account.Name, "", "", "")
	if err != nil {
		return nil, fmt.Errorf("could not list file storages: %w", err)
	}

	for _, value := range result.Values() {
		fileStorages, err := handleFileStorage(account, value)
		if err != nil {
			return nil, fmt.Errorf("could not handle file storage: %w", err)
		}

		log.Infof("Adding file storage %+v", fileStorages)

		list = append(list, fileStorages)
	}

	return list, nil
}

func (d *azureStorageDiscovery) discoverObjectStorages(account *storage.Account) ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	client := storage.NewBlobContainersClient(to.String(d.sub.SubscriptionID))
	d.apply(&client.Client)

	result, err := client.List(context.Background(), getResourceGroupName(*account.ID), *account.Name, "", "", "")
	if err != nil {
		return nil, fmt.Errorf("could not list object storages: %w", err)
	}

	for _, value := range result.Values() {
		objectStorages, err := handleObjectStorage(account, value)
		if err != nil {
			return nil, fmt.Errorf("could not handle object storage: %w", err)
		}
		log.Infof("Adding object storage %+v", objectStorages)

		list = append(list, objectStorages)
	}

	return list, nil
}

func (d *azureStorageDiscovery) handleBlockStorage(disk *compute.Disk) (*voc.BlockStorage, error) {
	enc, err := d.blockStorageAtRestEncryption(disk)
	if err != nil {
		return nil, fmt.Errorf("could not get block storage properties for the atRestEncryption: %w", err)
	}

	return &voc.BlockStorage{
		Storage: &voc.Storage{
			Resource: &voc.Resource{
				ID:           voc.ResourceID(to.String(disk.ID)),
				Name:         to.String(disk.Name),
				CreationTime: disk.TimeCreated.Unix(),
				Type:         []string{"BlockStorage", "Storage", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: *disk.Location,
				},
			},
			AtRestEncryption: enc,
		},
	}, nil
}

func handleObjectStorage(account *storage.Account, container storage.ListContainerItem) (*voc.ObjectStorage, error) {
	enc, err := storageAtRestEncryption(account)
	if err != nil {
		return nil, fmt.Errorf("could not get object storage properties for the atRestEncryption: %w", err)
	}

	return &voc.ObjectStorage{
		Storage: &voc.Storage{
			Resource: &voc.Resource{
				ID:           voc.ResourceID(to.String(container.ID)),
				Name:         to.String(container.Name),
				CreationTime: account.CreationTime.Unix(),
				Type:         []string{"ObjectStorage", "Storage", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: *account.Location,
				},
			},
			AtRestEncryption: enc,
		},
	}, nil
}

func (*azureStorageDiscovery) handleStorageAccount(account *storage.Account, storagesList []voc.IsCloudResource) (*voc.StorageService, error) {
	var storageResourceIDs []voc.ResourceID

	// Get all object storage IDs
	for _, storage := range storagesList {
		storageResourceIDs = append(storageResourceIDs, storage.GetID())
	}

	storageService := &voc.StorageService{
		Storages: storageResourceIDs,
		NetworkService: &voc.NetworkService{
			Networking: &voc.Networking{
				Resource: &voc.Resource{
					ID:           voc.ResourceID(to.String(account.ID)),
					Name:         *account.Name,
					CreationTime: account.CreationTime.Unix(),
					Type:         []string{"StorageService", "NetworkService", "Networking", "Resource"},
					GeoLocation: voc.GeoLocation{
						Region: *account.Location,
					},
				},
			},
		},
		HttpEndpoint: &voc.HttpEndpoint{
			Url: generalizeURL(*account.PrimaryEndpoints.Blob),
			TransportEncryption: &voc.TransportEncryption{
				Enforced:   to.Bool(account.EnableHTTPSTrafficOnly),
				Enabled:    true, // cannot be disabled
				TlsVersion: string(account.MinimumTLSVersion),
				Algorithm:  "TLS",
			},
		},
	}

	return storageService, nil
}

//generalizeURL generalizes the URL, because the URL depends on the storage type
func generalizeURL(url string) string {
	urlSplit := strings.Split(url, ".")
	urlSplit[1] = "[file,blob]"
	newURL := strings.Join(urlSplit, ".")

	return newURL
}

func handleFileStorage(account *storage.Account, fileshare storage.FileShareItem) (*voc.FileStorage, error) {
	enc, err := storageAtRestEncryption(account)
	if err != nil {
		return nil, fmt.Errorf("could not get file storage properties for the atRestEncryption: %w", err)
	}

	return &voc.FileStorage{
		Storage: &voc.Storage{
			Resource: &voc.Resource{
				ID:           voc.ResourceID(to.String(fileshare.ID)),
				Name:         to.String(fileshare.Name),
				CreationTime: account.CreationTime.Unix(),
				Type:         []string{"FileStorage", "Storage", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: *account.Location,
				},
			},
			AtRestEncryption: enc,
		},
	}, nil
}

func (d *azureStorageDiscovery) blockStorageAtRestEncryption(disk *compute.Disk) (voc.HasAtRestEncryption, error) {

	var enc voc.HasAtRestEncryption

	if disk.Encryption.Type == compute.EncryptionTypeEncryptionAtRestWithPlatformKey {
		enc = voc.ManagedKeyEncryption{AtRestEncryption: &voc.AtRestEncryption{
			Algorithm: "AES256",
			Enabled:   true,
		}}
	} else if disk.Encryption.Type == compute.EncryptionTypeEncryptionAtRestWithCustomerKey {
		var keyUrl string
		discEncryptionSetID := disk.Encryption.DiskEncryptionSetID

		keyUrl, err := d.sourceVaultID(*discEncryptionSetID)
		if err != nil {
			return nil, fmt.Errorf("could not get keyVaultID: %w", err)
		}

		enc = voc.CustomerKeyEncryption{
			AtRestEncryption: &voc.AtRestEncryption{
				Algorithm: "", // TODO(garuppel): TBD
				Enabled:   true,
			},
			KeyUrl: keyUrl,
		}
	} else {
		return enc, errors.New("error getting atRestEncryption properties of blockStorage")
	}

	return enc, nil
}

func storageAtRestEncryption(account *storage.Account) (voc.HasAtRestEncryption, error) {

	var enc voc.HasAtRestEncryption

	if account.Encryption.KeySource == storage.KeySourceMicrosoftStorage {
		enc = voc.ManagedKeyEncryption{AtRestEncryption: &voc.AtRestEncryption{
			Algorithm: "AES256",
			Enabled:   true,
		}}
	} else if account.Encryption.KeySource == storage.KeySourceMicrosoftKeyvault {
		enc = voc.CustomerKeyEncryption{
			AtRestEncryption: &voc.AtRestEncryption{
				Algorithm: "", // TODO(garuppel): TBD
				Enabled:   true,
			},
			KeyUrl: to.String(account.Encryption.KeyVaultProperties.KeyVaultURI),
		}
	} else {
		return enc, errors.New("error getting atRestEncryption properties of storage account")
	}

	return enc, nil
}

func (d *azureStorageDiscovery) sourceVaultID(discEncryptionSetID string) (string, error) {
	client := compute.NewDiskEncryptionSetsClient(to.String(d.sub.SubscriptionID))
	d.apply(&client.Client)

	discEncryptionSet, err := client.Get(context.Background(), getResourceGroupName(discEncryptionSetID), diskEncryptionSetName(discEncryptionSetID))
	if err != nil {
		return "", fmt.Errorf("could not get discEncryptionSet: %w", err)
	}

	if discEncryptionSet.EncryptionSetProperties.ActiveKey.SourceVault.ID == nil {
		return "", fmt.Errorf("could not get sourceVaultID")
	}

	return *discEncryptionSet.ActiveKey.SourceVault.ID, nil
}

func diskEncryptionSetName(discEncryptionSetID string) string {
	splitName := strings.Split(discEncryptionSetID, "/")
	return splitName[8]
}
