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
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

var (
	ErrEmptyStorageAccount        = errors.New("storage account is empty")
	ErrMissingDiskEncryptionSetID = errors.New("no disk encryption set ID was specified")
)

type azureStorageDiscovery struct {
	azureDiscovery
}

func NewAzureStorageDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &azureStorageDiscovery{
		azureDiscovery{
			discovererComponent: StorageComponent,
		},
	}

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
		return nil, fmt.Errorf("%s: %w", ErrCouldNotAuthenticate, err)
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
	client, err := armstorage.NewAccountsClient(util.Deref(d.sub.SubscriptionID), d.cred, &d.clientOptions)
	if err != nil {
		err = fmt.Errorf("could not get new storage accounts client: %w", err)
		return nil, err
	}

	// List all storage accounts across all resource groups
	listPager := client.NewListPager(&armstorage.AccountsClientListOptions{})
	for listPager.More() {
		pageResponse, err := listPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("%s: %v", ErrGettingNextPage, err)
			return nil, err
		}

		// Discover object and file storages
		for _, account := range pageResponse.Value {
			// Discover object storages
			objectStorages, err := d.discoverObjectStorages(account)
			if err != nil {
				return nil, fmt.Errorf("could not handle object storages: %w", err)
			}

			// Discover file storages
			fileStorages, err := d.discoverFileStorages(account)
			if err != nil {
				return nil, fmt.Errorf("could not handle file storages: %w", err)
			}

			storageResourcesList = append(storageResourcesList, objectStorages...)
			storageResourcesList = append(storageResourcesList, fileStorages...)

			// Create storage service for all storage account resources
			storageService, err := d.handleStorageAccount(account, storageResourcesList)
			if err != nil {
				return nil, fmt.Errorf("could not create storage service: %w", err)
			}

			storageResourcesList = append(storageResourcesList, storageService)
		}
	}

	return storageResourcesList, err
}

func (d *azureStorageDiscovery) discoverFileStorages(account *armstorage.Account) ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	// Create file shares client
	client, err := armstorage.NewFileSharesClient(util.Deref(d.sub.SubscriptionID), d.cred, &d.clientOptions)
	if err != nil {
		err = fmt.Errorf("could not get new virtual machines client: %w", err)
		return nil, err
	}

	// List all file shares in the specified resource group
	listPager := client.NewListPager(resourceGroupName(util.Deref(account.ID)), util.Deref(account.Name), &armstorage.FileSharesClientListOptions{})
	for listPager.More() {
		pageResponse, err := listPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("%s: %v", ErrGettingNextPage, err)
			return nil, err
		}

		for _, value := range pageResponse.Value {
			fileStorages, err := handleFileStorage(account, value)
			if err != nil {
				return nil, fmt.Errorf("could not handle file storage: %w", err)
			}

			log.Infof("Adding file storage '%s", fileStorages.Name)

			list = append(list, fileStorages)
		}
	}

	return list, nil
}

func (d *azureStorageDiscovery) discoverObjectStorages(account *armstorage.Account) ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	// Create blob containers client
	// TODO(oxisto): do not re-create client
	client, err := armstorage.NewBlobContainersClient(util.Deref(d.sub.SubscriptionID), d.cred, &d.clientOptions)
	if err != nil {
		err = fmt.Errorf("could not get new virtual machines client: %w", err)
		return nil, err
	}

	// List all blob containers in the specified resource group
	listPager := client.NewListPager(resourceGroupName(util.Deref(account.ID)), util.Deref(account.Name), &armstorage.BlobContainersClientListOptions{})
	for listPager.More() {
		pageResponse, err := listPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("%s: %v", ErrGettingNextPage, err)
			return nil, err
		}

		for _, value := range pageResponse.Value {
			objectStorages, err := handleObjectStorage(account, value)
			if err != nil {
				return nil, fmt.Errorf("could not handle object storage: %w", err)
			}
			log.Infof("Adding object storage '%s'", objectStorages.Name)

			list = append(list, objectStorages)
		}
	}

	return list, nil
}

func handleObjectStorage(account *armstorage.Account, container *armstorage.ListContainerItem) (*voc.ObjectStorage, error) {
	if account == nil {
		return nil, ErrEmptyStorageAccount
	}

	// It is possible that the container is not empty. In that case we have to check if a mandatory field is empty, so the whole disk is empty
	if container == nil || container.ID == nil {
		return nil, fmt.Errorf("container is nil")
	}

	enc, err := storageAtRestEncryption(account)
	if err != nil {
		return nil, fmt.Errorf("could not get object storage properties for the atRestEncryption: %w", err)
	}

	return &voc.ObjectStorage{
		Storage: &voc.Storage{
			Resource: &voc.Resource{
				ID:           voc.ResourceID(util.Deref(container.ID)),
				ServiceID:    discovery.DefaultCloudServiceID,
				Name:         util.Deref(container.Name),
				CreationTime: account.Properties.CreationTime.Unix(),
				Type:         []string{"ObjectStorage", "Storage", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: util.Deref(account.Location),
				},
				Labels: labels(account.Tags), // the storage account labels the object storage belongs to
			},
			AtRestEncryption: enc,
			Immutability: &voc.Immutability{
				Enabled: util.Deref(container.Properties.HasImmutabilityPolicy),
			},
		},
		PublicAccess: util.Deref(container.Properties.PublicAccess) != armstorage.PublicAccessNone,
	}, nil
}

func (*azureStorageDiscovery) handleStorageAccount(account *armstorage.Account, storagesList []voc.IsCloudResource) (*voc.StorageService, error) {
	var storageResourceIDs []voc.ResourceID

	if account == nil {
		return nil, ErrEmptyStorageAccount
	}

	// Get all object storage IDs
	for _, storage := range storagesList {
		if strings.Contains(string(storage.GetID()), accountName(util.Deref(account.ID))) {
			storageResourceIDs = append(storageResourceIDs, storage.GetID())
		}
	}

	te := &voc.TransportEncryption{
		Enforced:   util.Deref(account.Properties.EnableHTTPSTrafficOnly),
		Enabled:    true, // cannot be disabled
		TlsVersion: string(*account.Properties.MinimumTLSVersion),
		Algorithm:  "TLS",
	}

	storageService := &voc.StorageService{
		Storages: storageResourceIDs,
		NetworkService: &voc.NetworkService{
			Networking: &voc.Networking{
				Resource: &voc.Resource{
					ID:           voc.ResourceID(util.Deref(account.ID)),
					ServiceID:    discovery.DefaultCloudServiceID,
					Name:         util.Deref(account.Name),
					CreationTime: account.Properties.CreationTime.Unix(),
					Type:         []string{"StorageService", "NetworkService", "Networking", "Resource"},
					GeoLocation: voc.GeoLocation{
						Region: util.Deref(account.Location),
					},
					Labels: labels(account.Tags),
				},
			},
			TransportEncryption: te,
		},
		HttpEndpoint: &voc.HttpEndpoint{
			Url:                 generalizeURL(util.Deref(account.Properties.PrimaryEndpoints.Blob)),
			TransportEncryption: te,
		},
	}

	return storageService, nil
}

// generalizeURL generalizes the URL, because the URL depends on the storage type
func generalizeURL(url string) string {
	if url == "" {
		return ""
	}

	urlSplit := strings.Split(url, ".")
	urlSplit[1] = "[file,blob]"
	newURL := strings.Join(urlSplit, ".")

	return newURL
}

func handleFileStorage(account *armstorage.Account, fileshare *armstorage.FileShareItem) (*voc.FileStorage, error) {
	if account == nil {
		return nil, ErrEmptyStorageAccount
	}

	// It is possible that the fileshare is not empty. In that case we have to check if a mandatory field is empty, so the whole disk is empty
	if fileshare == nil || fileshare.ID == nil {
		return nil, fmt.Errorf("fileshare is nil")
	}

	enc, err := storageAtRestEncryption(account)
	if err != nil {
		return nil, fmt.Errorf("could not get file storage properties for the atRestEncryption: %w", err)
	}

	return &voc.FileStorage{
		Storage: &voc.Storage{
			Resource: &voc.Resource{
				ID:           voc.ResourceID(util.Deref(fileshare.ID)),
				ServiceID:    discovery.DefaultCloudServiceID,
				Name:         util.Deref(fileshare.Name),
				CreationTime: account.Properties.CreationTime.Unix(),
				Type:         []string{"FileStorage", "Storage", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: util.Deref(account.Location),
				},
				Labels: labels(account.Tags), // the storage account labels the file storage belongs to
			},
			AtRestEncryption: enc,
		},
	}, nil
}

// storageAtRestEncryption takes encryption properties of an armstorage.Account and converts it into our respective
// ontology object.
func storageAtRestEncryption(account *armstorage.Account) (enc voc.HasAtRestEncryption, err error) {
	if account == nil {
		return enc, ErrEmptyStorageAccount
	}

	if account.Properties == nil || account.Properties.Encryption.KeySource == nil {
		return enc, errors.New("keySource is empty")
	} else if *account.Properties.Encryption.KeySource == armstorage.KeySourceMicrosoftStorage {
		enc = voc.ManagedKeyEncryption{
			AtRestEncryption: &voc.AtRestEncryption{
				Algorithm: "AES256",
				Enabled:   true,
			},
		}
	} else if *account.Properties.Encryption.KeySource == armstorage.KeySourceMicrosoftKeyvault {
		enc = voc.CustomerKeyEncryption{
			AtRestEncryption: &voc.AtRestEncryption{
				Algorithm: "", // TODO(garuppel): TBD
				Enabled:   true,
			},
			KeyUrl: util.Deref(account.Properties.Encryption.KeyVaultProperties.KeyVaultURI),
		}
	}

	return enc, nil
}

// diskEncryptionSetName return the disk encryption set ID's name
func diskEncryptionSetName(diskEncryptionSetID string) string {
	if diskEncryptionSetID == "" {
		return ""
	}
	splitName := strings.Split(diskEncryptionSetID, "/")
	return splitName[8]
}

// accountName return the ID's account name
func accountName(id string) string {
	if id == "" {
		return ""
	}

	splitName := strings.Split(id, "/")
	return splitName[8]
}
