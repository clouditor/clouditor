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

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dataprotection/armdataprotection"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

var (
	ErrEmptyStorageAccount        = errors.New("storage account is empty")
	ErrMissingDiskEncryptionSetID = errors.New("no disk encryption set ID was specified")
)

type azureStorageDiscovery struct {
	*azureDiscovery
	defenderProperties map[string]*defenderProperties
	backupMap          map[string]*voc.Backup
}

func NewAzureStorageDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &azureStorageDiscovery{
		&azureDiscovery{
			discovererComponent: StorageComponent,
			csID:                discovery.DefaultCloudServiceID,
		},
		make(map[string]*defenderProperties),
		make(map[string]*voc.Backup),
	}

	// Apply options
	for _, opt := range opts {
		opt(d.azureDiscovery)
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

	// Discover Defender for X properties to add it to the required resource properties
	d.defenderProperties, err = d.discoverDefender()
	if err != nil {
		log.Errorf("Could not discover Defender for X")
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

	// // initialize data protection client
	// if err := d.initDataProtectionClient(); err != nil {
	// 	return nil, err
	// }

	// initialize backup instances client
	if err := d.initBackupInstancesClient(); err != nil {
		return nil, err
	}

	// initialize backup vaults client
	if err := d.initBackupVaultsClient(); err != nil {
		return nil, err
	}

	// initialize storage accounts client
	if err := d.initAccountsClient(); err != nil {
		return nil, err
	}

	// initialize blob container client
	if err := d.initBlobContainerClient(); err != nil {
		return nil, err
	}

	// initialize file share client
	if err := d.initFileStorageClient(); err != nil {
		return nil, err
	}

	// Discover object and file storages
	err := listPager(d.azureDiscovery,
		d.clients.accountsClient.NewListPager,
		d.clients.accountsClient.NewListByResourceGroupPager,
		func(res armstorage.AccountsClientListResponse) []*armstorage.Account {
			return res.Value
		},
		func(res armstorage.AccountsClientListByResourceGroupResponse) []*armstorage.Account {
			return res.Value
		},
		func(account *armstorage.Account) error {
			// Discover backup policies
			backupVaults, err := d.discoverBackupVaults(account)
			if err != nil {
				log.Errorf("could not discover backup vaults: %v", err)
			}

			// Store voc.Backup for each entry in the backupMap
			d.handleBackupVaults(backupVaults)

			// Discover object storages
			objectStorages, err := d.discoverObjectStorages(account)
			if err != nil {
				return fmt.Errorf("could not handle object storages: %w", err)
			}

			// Discover file storages
			fileStorages, err := d.discoverFileStorages(account)
			if err != nil {
				return fmt.Errorf("could not handle file storages: %w", err)
			}

			storageResourcesList = append(storageResourcesList, objectStorages...)
			storageResourcesList = append(storageResourcesList, fileStorages...)

			// Create storage service for all storage account resources
			storageService, err := d.handleStorageAccount(account, storageResourcesList)
			if err != nil {
				return fmt.Errorf("could not create storage service: %w", err)
			}

			storageResourcesList = append(storageResourcesList, storageService)

			return nil
		})
	if err != nil {
		return nil, err
	}

	return storageResourcesList, nil
}

// handleBackupVaults creates a voc.Backup and stores it to the backupMap.
func (d *azureStorageDiscovery) handleBackupVaults(vaults []*armdataprotection.BackupVaultResource) {
	for _, vault := range vaults {
		// Discover backup instances for given vault
		instances, err := d.discoverBackupInstances(resourceGroupName(*vault.ID), *vault.Name)
		if err != nil {
			log.Errorf("could not discover backup instances: %v", err)
			continue
		}

		for _, instance := range instances {
			d.backupMap[idUpToStorageAccount(*instance.Properties.DataSourceInfo.ResourceID)] = &voc.Backup{
				Enabled: true,
				Policy:  *instance.Properties.PolicyInfo.PolicyID,
				// RetentionPeriod: , // TODO(all): Add retention period
				Storage: voc.ResourceID(*instance.ID),
				GeoLocation: voc.GeoLocation{
					Region: *vault.Location,
				},
			}
		}
	}
}

func (d *azureStorageDiscovery) discoverFileStorages(account *armstorage.Account) ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	// List all file shares in the specified resource group
	listPager := d.clients.fileStorageClient.NewListPager(resourceGroupName(util.Deref(account.ID)), util.Deref(account.Name), &armstorage.FileSharesClientListOptions{})
	for listPager.More() {
		pageResponse, err := listPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("%s: %v", ErrGettingNextPage, err)
			return nil, err
		}

		for _, value := range pageResponse.Value {
			fileStorages, err := d.handleFileStorage(account, value)
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

	// List all blob containers in the specified resource group
	listPager := d.clients.blobContainerClient.NewListPager(resourceGroupName(util.Deref(account.ID)), util.Deref(account.Name), &armstorage.BlobContainersClientListOptions{})
	for listPager.More() {
		pageResponse, err := listPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("%s: %v", ErrGettingNextPage, err)
			return nil, err
		}

		for _, value := range pageResponse.Value {
			objectStorages, err := d.handleObjectStorage(account, value)
			if err != nil {
				return nil, fmt.Errorf("could not handle object storage: %w", err)
			}
			log.Infof("Adding object storage '%s'", objectStorages.Name)

			list = append(list, objectStorages)
		}
	}

	return list, nil
}

func (d *azureStorageDiscovery) handleObjectStorage(account *armstorage.Account, container *armstorage.ListContainerItem) (*voc.ObjectStorage, error) {
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

	backup := d.backupMap[idUpToStorageAccount(*container.ID)]

	return &voc.ObjectStorage{
		Storage: &voc.Storage{
			Resource: discovery.NewResource(d,
				voc.ResourceID(util.Deref(container.ID)),
				util.Deref(container.Name),
				// We only have the creation time of the storage account the object storage belongs to
				account.Properties.CreationTime,
				voc.GeoLocation{
					// The location is the same as the storage account
					Region: util.Deref(account.Location),
				},
				// The storage account labels the object storage belongs to
				labels(account.Tags),
				voc.ObjectStorageType,
			),
			AtRestEncryption: enc,
			Immutability: &voc.Immutability{
				Enabled: util.Deref(container.Properties.HasImmutabilityPolicy),
			},
			ResourceLogging: d.createResourceLogging(),
			Backup:          backup,
		},
		PublicAccess: util.Deref(container.Properties.PublicAccess) != armstorage.PublicAccessNone,
	}, nil
}

func (d *azureStorageDiscovery) handleStorageAccount(account *armstorage.Account, storagesList []voc.IsCloudResource) (*voc.ObjectStorageService, error) {
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

	storageService := &voc.ObjectStorageService{
		StorageService: &voc.StorageService{
			Storage: storageResourceIDs,
			NetworkService: &voc.NetworkService{
				Networking: &voc.Networking{
					Resource: discovery.NewResource(d,
						voc.ResourceID(util.Deref(account.ID)),
						util.Deref(account.Name),
						account.Properties.CreationTime,
						voc.GeoLocation{
							Region: util.Deref(account.Location),
						},
						labels(account.Tags),
						voc.ObjectStorageServiceType,
					),
				},
				TransportEncryption: te,
			},
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

func (d *azureStorageDiscovery) handleFileStorage(account *armstorage.Account, fileshare *armstorage.FileShareItem) (*voc.FileStorage, error) {
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
			Resource: discovery.NewResource(d,
				voc.ResourceID(util.Deref(fileshare.ID)),
				util.Deref(fileshare.Name),
				// We only have the creation time of the storage account the file storage belongs to
				account.Properties.CreationTime,
				voc.GeoLocation{
					// The location is the same as the storage account
					Region: util.Deref(account.Location),
				},
				// The storage account labels the file storage belongs to
				labels(account.Tags),
				voc.FileStorageType,
			),
			ResourceLogging:  d.createResourceLogging(),
			AtRestEncryption: enc,
		},
	}, nil
}

// storageAtRestEncryption takes encryption properties of an armstorage.Account and converts it into our respective
// ontology object.
func storageAtRestEncryption(account *armstorage.Account) (enc voc.IsAtRestEncryption, err error) {
	if account == nil {
		return enc, ErrEmptyStorageAccount
	}

	if account.Properties == nil || account.Properties.Encryption.KeySource == nil {
		return enc, errors.New("keySource is empty")
	} else if *account.Properties.Encryption.KeySource == armstorage.KeySourceMicrosoftStorage {
		enc = &voc.ManagedKeyEncryption{
			AtRestEncryption: &voc.AtRestEncryption{
				Algorithm: "AES256",
				Enabled:   true,
			},
		}
	} else if *account.Properties.Encryption.KeySource == armstorage.KeySourceMicrosoftKeyvault {
		enc = &voc.CustomerKeyEncryption{
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

// TODO(all): Update to generic function or method
func (d *azureStorageDiscovery) createResourceLogging() (resourceLogging *voc.ResourceLogging) {
	if d.defenderProperties[DefenderStorageType] != nil {
		resourceLogging = &voc.ResourceLogging{
			MonitoringLogDataEnabled: d.defenderProperties[DefenderStorageType].monitoringLogDataEnabled,
			SecurityAlertsEnabled:    d.defenderProperties[DefenderStorageType].securityAlertsEnabled,
		}
	}

	return
}

// initAccountsClient creates the client if not already exists
func (d *azureStorageDiscovery) initAccountsClient() (err error) {
	d.clients.accountsClient, err = initClient(d.clients.accountsClient, d.azureDiscovery, armstorage.NewAccountsClient)
	return
}

// initBlobContainerClient creates the client if not already exists
func (d *azureStorageDiscovery) initBlobContainerClient() (err error) {
	d.clients.blobContainerClient, err = initClient(d.clients.blobContainerClient, d.azureDiscovery, armstorage.NewBlobContainersClient)
	return
}

// initFileStorageClient creates the client if not already exists
func (d *azureStorageDiscovery) initFileStorageClient() (err error) {
	d.clients.fileStorageClient, err = initClient(d.clients.fileStorageClient, d.azureDiscovery, armstorage.NewFileSharesClient)
	return
}

// initDataProtectionClient creates the client if not already exists
func (d *azureStorageDiscovery) initDataProtectionClient() (err error) {
	d.clients.dataProtectionClient, err = initClient(d.clients.dataProtectionClient, d.azureDiscovery, armdataprotection.NewBackupPoliciesClient)

	return
}

// initBackupVaultsClient creates the client if not already exists
func (d *azureStorageDiscovery) initBackupVaultsClient() (err error) {
	d.clients.backupVaultClient, err = initClient(d.clients.backupVaultClient, d.azureDiscovery, armdataprotection.NewBackupVaultsClient)

	return
}

// initBackupInstancesClient creates the client if not already exists
func (d *azureStorageDiscovery) initBackupInstancesClient() (err error) {
	d.clients.backupInstancesClient, err = initClient(d.clients.backupInstancesClient, d.azureDiscovery, armdataprotection.NewBackupInstancesClient)

	return
}

// discoverBackupVaults receives all backup vaults in the subscription.
func (d *azureStorageDiscovery) discoverBackupVaults(account *armstorage.Account) ([]*armdataprotection.BackupVaultResource, error) {
	var (
		list armdataprotection.BackupVaultsClientGetInResourceGroupResponse
		err  error
	)

	// List all backup vaults in the given resource group
	listPager := d.clients.backupVaultClient.NewGetInResourceGroupPager(resourceGroupName(*account.ID), &armdataprotection.BackupVaultsClientGetInResourceGroupOptions{})
	for listPager.More() {
		list, err = listPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("%s: %v", ErrGettingNextPage, err)
			return nil, err
		}
	}

	return list.Value, nil
}

// discoverBackupInstances retrieves the instances in a given backup vault.
func (d *azureStorageDiscovery) discoverBackupInstances(resourceGroup, vaultName string) ([]*armdataprotection.BackupInstanceResource, error) {
	var (
		list armdataprotection.BackupInstancesClientListResponse
		err  error
	)

	// List all instances in the given backupp vault
	listPager := d.clients.backupInstancesClient.NewListPager(resourceGroup, vaultName, &armdataprotection.BackupInstancesClientListOptions{})
	for listPager.More() {
		list, err = listPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("%s: %v", ErrGettingNextPage, err)
			return nil, err
		}
	}

	return list.Value, nil
}

// idUpToStorageAccount returns the resource ID cutting of the storage type information, e.g., '/subscriptions/XXXXXXXXXXXX-XXXX-XXXX-XXXXXXXXXXXX/resourceGroups/resourceGroupName/providers/Microsoft.Storage/storageAccounts/containerName'
func idUpToStorageAccount(id string) string {
	if id == "" {
		return ""
	}

	split := strings.Split(id, "/")

	if len(split) < 8 {
		return ""
	}

	return "/" + split[1] + "/" + split[2] + "/" + split[3] + "/" + split[4] + "/" + split[5] + "/" + split[6] + "/" + split[7] + "/" + split[8]
}
