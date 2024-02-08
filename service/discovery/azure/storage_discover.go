//go:build exclude

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
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

var (
	ErrEmptyStorageAccount        = errors.New("storage account is empty")
	ErrMissingDiskEncryptionSetID = errors.New("no disk encryption set ID was specified")
	ErrBackupStorageNotAvailable  = errors.New("backup storages not available")
)

// discoverCosmosDB discovers Cosmos DB accounts
func (d *azureDiscovery) discoverCosmosDB() ([]voc.IsCloudResource, error) {
	var (
		list []voc.IsCloudResource
		err  error
	)

	// initialize Cosmos DB client
	if err := d.initCosmosDBClient(); err != nil {
		return nil, err
	}

	// Discover Cosmos DB
	err = listPager(d,
		d.clients.cosmosDBClient.NewListPager,
		d.clients.cosmosDBClient.NewListByResourceGroupPager,
		func(res armcosmos.DatabaseAccountsClientListResponse) []*armcosmos.DatabaseAccountGetResults {
			return res.Value
		},
		func(res armcosmos.DatabaseAccountsClientListByResourceGroupResponse) []*armcosmos.DatabaseAccountGetResults {
			return res.Value
		},
		func(dbAccount *armcosmos.DatabaseAccountGetResults) error {
			cosmos, err := d.handleCosmosDB(dbAccount)
			if err != nil {
				return fmt.Errorf("could not cosmos db accounts: %w", err)
			}
			log.Infof("Adding Cosmos DB account '%s", *dbAccount.Name)
			list = append(list, cosmos...)

			return nil
		})
	if err != nil {
		return nil, err
	}

	return list, nil
}

// discoverMongoDBDatabases returns a list of Mongo DB databases for a specific Mongo DB account
func (d *azureDiscovery) discoverMongoDBDatabases(account *armcosmos.DatabaseAccountGetResults, atRestEnc voc.IsAtRestEncryption) []voc.IsCloudResource {
	var (
		list []voc.IsCloudResource
		err  error
	)

	// initialize Mongo DB resources client
	if err = d.initMongoDResourcesBClient(); err != nil {
		log.Errorf("error initializing Mongo DB resource client: %v", err)
		return list
	}

	// Discover Mongo DB databases
	serverlistPager := d.clients.mongoDBResourcesClient.NewListMongoDBDatabasesPager(resourceGroupName(util.Deref(account.ID)), *account.Name, &armcosmos.MongoDBResourcesClientListMongoDBDatabasesOptions{})
	for serverlistPager.More() {
		pageResponse, err := serverlistPager.NextPage(context.TODO())
		if err != nil {
			log.Errorf("%s: %v", ErrGettingNextPage, err)
			return list
		}

		for _, value := range pageResponse.Value {
			// Create Cosmos DB database storage voc object
			mongoDB := &voc.DatabaseStorage{
				Storage: &voc.Storage{
					Resource: discovery.NewResource(d,
						voc.ResourceID(*value.ID),
						util.Deref(value.Name),
						nil, // creation time of database not available
						voc.GeoLocation{
							Region: *value.Location,
						},
						labels(value.Tags),
						voc.ResourceID(*account.ID),
						voc.DatabaseStorageType,
						account,
						value),

					AtRestEncryption: atRestEnc,
				},
			}
			list = append(list, mongoDB)
		}
	}

	return list
}

// discoverSqlServers discovers the sql server and databases
func (d *azureDiscovery) discoverSqlServers() ([]voc.IsCloudResource, error) {
	var (
		list []voc.IsCloudResource
		err  error
	)

	// initialize SQL server client
	if err := d.initSQLServersClient(); err != nil {
		return nil, err
	}

	// Discover sql server
	err = listPager(d,
		d.clients.sqlServersClient.NewListPager,
		d.clients.sqlServersClient.NewListByResourceGroupPager,
		func(res armsql.ServersClientListResponse) []*armsql.Server {
			return res.Value
		},
		func(res armsql.ServersClientListByResourceGroupResponse) []*armsql.Server {
			return res.Value
		},
		func(server *armsql.Server) error {
			db, err := d.handleSqlServer(server)
			if err != nil {
				return fmt.Errorf("could not handle sql database: %w", err)
			}
			log.Infof("Adding sql database '%s", *server.Name)
			list = append(list, db...)

			return nil
		})
	if err != nil {
		return nil, err
	}

	return list, nil
}

// getSqlDBs returns a list of SQL databases for a specific SQL account
func (d *azureDiscovery) getSqlDBs(server *armsql.Server) ([]voc.IsCloudResource, []voc.IsAnomalyDetection) {
	var (
		list                 []voc.IsCloudResource
		anomalyDetectionList []voc.IsAnomalyDetection
		err                  error
	)

	// initialize SQL databases client
	if err = d.initDatabasesClient(); err != nil {
		log.Errorf("error initializing database client: %v", err)
		return list, anomalyDetectionList
	}

	// Get databases for given server
	serverlistPager := d.clients.databasesClient.NewListByServerPager(resourceGroupName(util.Deref(server.ID)), *server.Name, &armsql.DatabasesClientListByServerOptions{})
	for serverlistPager.More() {
		pageResponse, err := serverlistPager.NextPage(context.TODO())
		if err != nil {
			log.Errorf("%s: %v", ErrGettingNextPage, err)
			return list, anomalyDetectionList
		}

		for _, value := range pageResponse.Value {
			// Create anomaly detection property
			// Get anomaly detection status
			anomalyDetectionEnabled, err := d.anomalyDetectionEnabled(server, value)
			if err != nil {
				log.Errorf("error getting anomaly detection info for database '%s': %v", *value.Name, err)
			}

			a := &voc.AnomalyDetection{
				Scope:   voc.ResourceID(*value.ID),
				Enabled: anomalyDetectionEnabled,
			}

			anomalyDetectionList = append(anomalyDetectionList, a)

			// Create database storage voc object
			sqlDB := &voc.DatabaseStorage{
				Storage: &voc.Storage{
					Resource: discovery.NewResource(d,
						voc.ResourceID(*value.ID),
						*value.Name,
						value.Properties.CreationDate,
						voc.GeoLocation{
							Region: *value.Location,
						},
						labels(value.Tags),
						voc.ResourceID(*server.ID),
						voc.DatabaseStorageType,
						value),
					AtRestEncryption: &voc.AtRestEncryption{
						Enabled:   *value.Properties.IsInfraEncryptionEnabled,
						Algorithm: AES256,
					},
					// TODO(all): Backups
				},
			}
			list = append(list, sqlDB)
		}
	}

	return list, anomalyDetectionList
}

func (d *azureDiscovery) discoverStorageAccounts() ([]voc.IsCloudResource, error) {
	var storageResourcesList []voc.IsCloudResource

	// initialize backup policies client
	if err := d.initBackupPoliciesClient(); err != nil {
		return nil, err
	}

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

	// Discover backup vaults
	err := d.discoverBackupVaults()
	if err != nil {
		log.Errorf("could not discover backup vaults: %v", err)
	}

	// Discover object and file storages
	err = listPager(d,
		d.clients.accountsClient.NewListPager,
		d.clients.accountsClient.NewListByResourceGroupPager,
		func(res armstorage.AccountsClientListResponse) []*armstorage.Account {
			return res.Value
		},
		func(res armstorage.AccountsClientListByResourceGroupResponse) []*armstorage.Account {
			return res.Value
		},
		func(account *armstorage.Account) error {
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

	// Add backup storage account objects
	if d.backupMap[DataSourceTypeStorageAccountObject] != nil && d.backupMap[DataSourceTypeStorageAccountObject].backupStorages != nil {
		storageResourcesList = append(storageResourcesList, d.backupMap[DataSourceTypeStorageAccountObject].backupStorages...)
	}

	return storageResourcesList, nil
}

func (d *azureDiscovery) discoverFileStorages(account *armstorage.Account) ([]voc.IsCloudResource, error) {
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

func (d *azureDiscovery) discoverObjectStorages(account *armstorage.Account) ([]voc.IsCloudResource, error) {
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
