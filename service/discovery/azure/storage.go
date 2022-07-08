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

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"

	"github.com/Azure/go-autorest/autorest/to"
)

type azureStorageDiscovery struct {
	azureDiscovery
}

func NewAzureStorageDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &azureStorageDiscovery{}

	// Apply options
	for _, opt := range opts {
		opt(&d.azureDiscovery)
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

	log.Info("Discover Azure storage resources")

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

	// Create storage accounts client
	client, err := armstorage.NewAccountsClient(to.String(d.sub.SubscriptionID), d.cred, &d.clientOptions)
	//d.apply(&policy.ClientOptions{})
	if err != nil {
		err = fmt.Errorf("could not get new storage accounts client: %w", err)
		return nil, err
	}

	// List all storage accounts accross all resource groups
	listPager := client.NewListPager(&armstorage.AccountsClientListOptions{})
	accounts := make([]*armstorage.Account, 0)
	for listPager.More() {
		pageResponse, err := listPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("error getting next page: %v", err)
			return nil, err
		}
		accounts = append(accounts, pageResponse.Value...)
	}

	// Discover object and file storages
	for _, account := range accounts {
		// Discover object storages
		objectStorages, err := d.discoverObjectStorages(account)
		if err != nil {
			return nil, fmt.Errorf("could not handle object storages: %w", err)
		}
		log.Infof("Adding object storages %+v", objectStorages)

		// Discover file storages
		fileStorages, err := d.discoverFileStorages(account)
		if err != nil {
			return nil, fmt.Errorf("could not handle file storages: %w", err)
		}
		log.Infof("Adding file storages %+v", fileStorages)

		storageResourcesList = append(storageResourcesList, objectStorages...)
		storageResourcesList = append(storageResourcesList, fileStorages...)

		// Create storage service for all storage account resources
		storageService, err := d.handleStorageAccount(account, storageResourcesList)
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

	// Create disks client
	client, err := armcompute.NewDisksClient(to.String(d.sub.SubscriptionID), d.cred, &d.clientOptions)
	if err != nil {
		err = fmt.Errorf("could not get new disks client: %w", err)
		return nil, err
	}

	// List all disks accross all resource groups
	listPager := client.NewListPager(&armcompute.DisksClientListOptions{})
	disks := make([]*armcompute.Disk, 0)
	for listPager.More() {
		pageResponse, err := listPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("error getting next page: %w", err)
			return nil, err
		}
		disks = append(disks, pageResponse.Value...)
	}

	for _, disk := range disks {
		blockStorages, err := d.handleBlockStorage(disk)
		if err != nil {
			return nil, fmt.Errorf("could not handle block storage: %w", err)
		}
		log.Infof("Adding block storage %+v", blockStorages)

		list = append(list, blockStorages)
	}

	return list, nil
}

func (d *azureStorageDiscovery) discoverFileStorages(account *armstorage.Account) ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	// Create file shares client
	client, err := armstorage.NewFileSharesClient(to.String(d.sub.SubscriptionID), d.cred, &d.clientOptions)
	if err != nil {
		err = fmt.Errorf("could not get new virtual machines client: %w", err)
		return nil, err
	}

	// List all file shares in the specified resource group
	listPager := client.NewListPager(resourceGroupName(to.String(account.ID)), to.String(account.Name), &armstorage.FileSharesClientListOptions{})
	fs := make([]*armstorage.FileShareItem, 0)
	for listPager.More() {
		pageResponse, err := listPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("error getting next page: %v", err)
			return nil, err
		}
		fs = append(fs, pageResponse.Value...)
	}

	for _, value := range fs {
		fileStorages, err := handleFileStorage(account, value)
		if err != nil {
			return nil, fmt.Errorf("could not handle file storage: %w", err)
		}

		log.Infof("Adding file storage %+v", fileStorages)

		list = append(list, fileStorages)
	}

	return list, nil
}

func (d *azureStorageDiscovery) discoverObjectStorages(account *armstorage.Account) ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	// Create blob containers client
	client, err := armstorage.NewBlobContainersClient(to.String(d.sub.SubscriptionID), d.cred, &d.clientOptions)
	if err != nil {
		err = fmt.Errorf("could not get new virtual machines client: %w", err)
		return nil, err
	}

	// List all file shares in the specified resource group
	listPager := client.NewListPager(resourceGroupName(to.String(account.ID)), to.String(account.Name), &armstorage.BlobContainersClientListOptions{})
	bc := make([]*armstorage.ListContainerItem, 0)
	for listPager.More() {
		pageResponse, err := listPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("error getting next page: %v", err)
			return nil, err
		}
		bc = append(bc, pageResponse.Value...)
	}

	for _, value := range bc {
		objectStorages, err := handleObjectStorage(account, value)
		if err != nil {
			return nil, fmt.Errorf("could not handle object storage: %w", err)
		}
		log.Infof("Adding object storage %+v", objectStorages)

		list = append(list, objectStorages)
	}

	return list, nil
}

func (d *azureStorageDiscovery) handleBlockStorage(disk *armcompute.Disk) (*voc.BlockStorage, error) {
	enc, err := d.blockStorageAtRestEncryption(disk)
	if err != nil {
		return nil, fmt.Errorf("could not get block storage properties for the atRestEncryption: %w", err)
	}

	return &voc.BlockStorage{
		Storage: &voc.Storage{
			Resource: &voc.Resource{
				ID:           voc.ResourceID(to.String(disk.ID)),
				Name:         to.String(disk.Name),
				CreationTime: disk.Properties.TimeCreated.Unix(),
				Type:         []string{"BlockStorage", "Storage", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: to.String(disk.Location),
				},
			},
			AtRestEncryption: enc,
		},
	}, nil
}

func handleObjectStorage(account *armstorage.Account, container *armstorage.ListContainerItem) (*voc.ObjectStorage, error) {
	enc, err := storageAtRestEncryption(account)
	if err != nil {
		return nil, fmt.Errorf("could not get object storage properties for the atRestEncryption: %w", err)
	}

	return &voc.ObjectStorage{
		Storage: &voc.Storage{
			Resource: &voc.Resource{
				ID:           voc.ResourceID(to.String(container.ID)),
				Name:         to.String(container.Name),
				CreationTime: account.Properties.CreationTime.Unix(),
				Type:         []string{"ObjectStorage", "Storage", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: to.String(account.Location),
				},
				Labels: labels(account.Tags), //the storage account labels the object storage belongs to
			},
			AtRestEncryption: enc,
		},
	}, nil
}

func (*azureStorageDiscovery) handleStorageAccount(account *armstorage.Account, storagesList []voc.IsCloudResource) (*voc.StorageService, error) {
	var storageResourceIDs []voc.ResourceID

	// Get all object storage IDs
	for _, storage := range storagesList {
		storageResourceIDs = append(storageResourceIDs, storage.GetID())
	}

	te := &voc.TransportEncryption{
		Enforced:   to.Bool(account.Properties.EnableHTTPSTrafficOnly),
		Enabled:    true, // cannot be disabled
		TlsVersion: string(*account.Properties.MinimumTLSVersion),
		Algorithm:  "TLS",
	}

	storageService := &voc.StorageService{
		Storages: storageResourceIDs,
		NetworkService: &voc.NetworkService{
			Networking: &voc.Networking{
				Resource: &voc.Resource{
					ID:           voc.ResourceID(to.String(account.ID)),
					Name:         to.String(account.Name),
					CreationTime: account.Properties.CreationTime.Unix(),
					Type:         []string{"StorageService", "NetworkService", "Networking", "Resource"},
					GeoLocation: voc.GeoLocation{
						Region: to.String(account.Location),
					},
					Labels: labels(account.Tags),
				},
			},
			TransportEncryption: te,
		},
		HttpEndpoint: &voc.HttpEndpoint{
			Url:                 generalizeURL(to.String(account.Properties.PrimaryEndpoints.Blob)),
			TransportEncryption: te,
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

func handleFileStorage(account *armstorage.Account, fileshare *armstorage.FileShareItem) (*voc.FileStorage, error) {
	enc, err := storageAtRestEncryption(account)
	if err != nil {
		return nil, fmt.Errorf("could not get file storage properties for the atRestEncryption: %w", err)
	}

	return &voc.FileStorage{
		Storage: &voc.Storage{
			Resource: &voc.Resource{
				ID:           voc.ResourceID(to.String(fileshare.ID)),
				Name:         to.String(fileshare.Name),
				CreationTime: account.Properties.CreationTime.Unix(),
				Type:         []string{"FileStorage", "Storage", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: to.String(account.Location),
				},
				Labels: labels(account.Tags), //the storage account labels the file storage belongs to
			},
			AtRestEncryption: enc,
		},
	}, nil
}

func (d *azureStorageDiscovery) blockStorageAtRestEncryption(disk *armcompute.Disk) (voc.HasAtRestEncryption, error) {

	var enc voc.HasAtRestEncryption

	if *disk.Properties.Encryption.Type == armcompute.EncryptionTypeEncryptionAtRestWithPlatformKey {
		enc = voc.ManagedKeyEncryption{AtRestEncryption: &voc.AtRestEncryption{
			Algorithm: "AES256",
			Enabled:   true,
		}}
	} else if *disk.Properties.Encryption.Type == armcompute.EncryptionTypeEncryptionAtRestWithCustomerKey {
		var keyUrl string
		discEncryptionSetID := disk.Properties.Encryption.DiskEncryptionSetID

		keyUrl, err := d.sourceVaultID(to.String(discEncryptionSetID))
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

func storageAtRestEncryption(account *armstorage.Account) (voc.HasAtRestEncryption, error) {

	var enc voc.HasAtRestEncryption

	if *account.Properties.Encryption.KeySource == armstorage.KeySourceMicrosoftStorage {
		enc = voc.ManagedKeyEncryption{AtRestEncryption: &voc.AtRestEncryption{
			Algorithm: "AES256",
			Enabled:   true,
		}}
	} else if *account.Properties.Encryption.KeySource == armstorage.KeySourceMicrosoftKeyvault {
		enc = voc.CustomerKeyEncryption{
			AtRestEncryption: &voc.AtRestEncryption{
				Algorithm: "", // TODO(garuppel): TBD
				Enabled:   true,
			},
			KeyUrl: to.String(account.Properties.Encryption.KeyVaultProperties.KeyVaultURI),
		}
	} else {
		return enc, errors.New("error getting atRestEncryption properties of storage account")
	}

	return enc, nil
}

func (d *azureStorageDiscovery) sourceVaultID(diskEncryptionSetID string) (string, error) {

	// Create Key Vault client
	client, err := armcompute.NewDiskEncryptionSetsClient(to.String(d.sub.SubscriptionID), d.cred, &d.clientOptions)
	if err != nil {
		err = fmt.Errorf("could not get new key vault client: %w", err)
		return "", err
	}

	// Get disk encryption set
	kv, err := client.Get(context.TODO(), resourceGroupName(diskEncryptionSetID), diskEncryptionSetName(diskEncryptionSetID), &armcompute.DiskEncryptionSetsClientGetOptions{})
	if err != nil {
		err = fmt.Errorf("could not get key vault: %w", err)
		return "", err
	}

	sourceVaultID := kv.DiskEncryptionSet.Properties.ActiveKey.SourceVault.ID

	if sourceVaultID == nil {
		return "", fmt.Errorf("could not get sourceVaultID")
	}

	return to.String(sourceVaultID), nil
}

func diskEncryptionSetName(diskEncryptionSetID string) string {
	splitName := strings.Split(diskEncryptionSetID, "/")
	return splitName[8]
}
