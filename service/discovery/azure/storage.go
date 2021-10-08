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
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/compute/mgmt/compute"
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

func (d *azureStorageDiscovery) Name() string {
	return "Azure Storage Account"
}

func (d *azureStorageDiscovery) Description() string {
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
	var list []voc.IsCloudResource

	client := storage.NewAccountsClient(to.String(d.sub.SubscriptionID))
	d.apply(&client.Client)

	result, err := client.List(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not list storage accounts: %w", err)
	}

	accounts := result.Values()
	for i := range accounts {
		// Discover object storages
		objectStorages, err := d.discoverObjectStorages(&accounts[i])
		if err != nil {
			return nil, fmt.Errorf("could not handle object storages: %w", err)
		}
		
		// Discover file storages
		fileStorages, err := d.discoverFileStorages(&accounts[i])
		if err != nil {
			return nil, fmt.Errorf("could not handle file storages: %w", err)
		}

		//Discover block storages
		blockStorages, err := d.discoverBlockStorages()
		if err != nil {
			return nil, fmt.Errorf("could not handle block storages: %w", err)
		}


		log.Infof("Adding storage account %+v", objectStorages)

		list = append(list, objectStorages...)
		list = append(list, fileStorages...)
		list = append(list, blockStorages...)
	}

	return list, err
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
		blockStorages := handleBlockStorage(disk)
		log.Infof("Adding block storage %+v", blockStorages)

		list = append(list, blockStorages)
	}

	return list, nil
}

func (d *azureStorageDiscovery) discoverFileStorages(account *storage.Account) ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	client := storage.NewFileSharesClient(to.String(d.sub.SubscriptionID))
	d.apply(&client.Client)

	result, err := client.List(context.Background(), GetResourceGroupName(*account.ID), *account.Name, "", "", "")
	if err != nil {
	return nil, fmt.Errorf("could not list file storages: %w", err)
	}

	for _, value := range result.Values() {
		fileStorages := handleFileStorage(account, value)
		log.Infof("Adding file storage %+v", fileStorages)

		list = append(list, fileStorages)
	}

	return list, nil
}

func (d *azureStorageDiscovery) discoverObjectStorages(account *storage.Account) ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	client := storage.NewBlobContainersClient(to.String(d.sub.SubscriptionID))
	d.apply(&client.Client)

	result, err := client.List(context.Background(), GetResourceGroupName(*account.ID), *account.Name, "", "", "")
	if err != nil {
		return nil, fmt.Errorf("could not list object storages: %w", err)
	}

	for _, value := range result.Values() {
		objectStorages := handleObjectStorage(account, value)
		log.Infof("Adding object storage %+v", objectStorages)

		list = append(list, objectStorages)
	}

	return list, nil
}

func handleBlockStorage(disk compute.Disk) voc.IsStorage {

	return &voc.BlockStorage{
		Storage: &voc.Storage{
			CloudResource: &voc.CloudResource{
				ID:           voc.ResourceID(to.String(disk.ID)),
				Name:         to.String(disk.Name),
				CreationTime: disk.TimeCreated.Unix(),
				Type:         []string{"BlockStorage", "Storage", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: *disk.Location,
				},
			},
			AtRestEncryption: &voc.AtRestEncryption{
				//KeyManager: string(disk.Encryption.Type), // TODO(all): What do we do with the encryption type? Do we leave it like that?
				Enabled:    true, // is always enabled
				Algorithm: "", // not available
			},
		},
	}
}

func handleObjectStorage(account *storage.Account, container storage.ListContainerItem) voc.IsStorage {
	return &voc.ObjectStorage{
		Storage: &voc.Storage{
			CloudResource: &voc.CloudResource{
				ID:           voc.ResourceID(to.String(container.ID)),
				Name:         to.String(container.Name),
				CreationTime: account.CreationTime.Unix(),
				Type:         []string{"ObjectStorage", "Storage", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: *account.Location,
				},
			},
			AtRestEncryption: &voc.AtRestEncryption{
				//KeyManager: string(account.Encryption.KeySource),
				Algorithm:  "AES-265", // seems to be always AES-256
				Enabled:    to.Bool(account.Encryption.Services.Blob.Enabled),
			},
		},
		HttpEndpoint: &voc.HttpEndpoint{
			Url: to.String(account.PrimaryEndpoints.Blob) + to.String(container.Name),
			TransportEncryption: &voc.TransportEncryption{
				Enforced:   to.Bool(account.EnableHTTPSTrafficOnly),
				Enabled:    true, // cannot be disabled
				TlsVersion: string(account.MinimumTLSVersion),
				Algorithm:  "", // TBD
			},
		},
	}
}

func handleFileStorage(account *storage.Account, fileshare storage.FileShareItem) voc.IsStorage {
	return &voc.FileStorage{
		Storage: &voc.Storage{
			CloudResource: &voc.CloudResource{
				ID:           voc.ResourceID(to.String(fileshare.ID)),
				Name:         to.String(fileshare.Name),
				CreationTime: account.CreationTime.Unix(),
				Type:         []string{"FileStorage", "Storage", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: *account.Location,
				},
			},
			AtRestEncryption: &voc.AtRestEncryption{
				//KeyManager: string(account.Encryption.KeySource),
				Algorithm:  "AES-265", // seems to be always AES-256
				Enabled:    to.Bool(account.Encryption.Services.File.Enabled),
			},
		},
		// TODO(all) Uncomment as soon as the HttpEndpoint is added to voc/file_storage.go
		//HttpEndpoint: &voc.HttpEndpoint{
		//	Url: to.String(account.PrimaryEndpoints.File) + to.String(fileshare.Name),
		//	TransportEncryption: &voc.TransportEncryption{
		//		Enforced:   to.Bool(account.EnableHTTPSTrafficOnly),
		//		Enabled:    true, // cannot be disabled
		//		TlsVersion: string(account.MinimumTLSVersion),
		//		Algorithm:  "", // TBD
		//	},
		//},
	}
}
