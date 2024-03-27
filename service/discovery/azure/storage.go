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
	"slices"
	"strings"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/constants"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dataprotection/armdataprotection"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/security/armsecurity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

var (
	ErrEmptyStorageAccount        = errors.New("storage account is empty")
	ErrMissingDiskEncryptionSetID = errors.New("no disk encryption set ID was specified")
)

// Currently supports only one backup. There could be more and even a metric that may check multiple backups
var backupOf = make(map[string]string)

type azureStorageDiscovery struct {
	*azureDiscovery
	defenderProperties map[string]*defenderProperties
}

func NewAzureStorageDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &azureStorageDiscovery{
		&azureDiscovery{
			discovererComponent: StorageComponent,
			csID:                discovery.DefaultCloudServiceID,
			backupMap:           make(map[string]*backup),
		},
		make(map[string]*defenderProperties),
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

	// initialize defender client
	if err := d.initDefenderClient(); err != nil {
		return nil, fmt.Errorf("could not initialize defender client: %w", err)
	}

	// Discover Defender for X properties to add it to the required resource properties
	d.defenderProperties, err = d.discoverDefender()
	if err != nil {
		return nil, fmt.Errorf("could not discover Defender for X: %w", err)
	}

	// Discover storage accounts
	storageAccounts, err := d.discoverStorageAccounts()
	if err != nil {
		return nil, fmt.Errorf("could not discover storage accounts: %w", err)
	}
	list = append(list, storageAccounts...)

	// Discover sql databases
	dbs, err := d.discoverSqlServers()
	if err != nil {
		return nil, fmt.Errorf("could not discover sql databases: %w", err)
	}
	list = append(list, dbs...)

	// Discover Cosmos DB
	cosmosDB, err := d.discoverCosmosDB()
	if err != nil {
		return nil, fmt.Errorf("could not discover cosmos db accounts: %w", err)
	}
	list = append(list, cosmosDB...)

	return
}

// discoverCosmosDB discovers Cosmos DB accounts
func (d *azureStorageDiscovery) discoverCosmosDB() ([]voc.IsCloudResource, error) {
	var (
		list []voc.IsCloudResource
		err  error
	)

	// initialize Cosmos DB client
	if err := d.initCosmosDBClient(); err != nil {
		return nil, err
	}

	// Discover Cosmos DB
	err = listPager(d.azureDiscovery,
		d.clients.cosmosDBClient.NewListPager,
		d.clients.cosmosDBClient.NewListByResourceGroupPager,
		func(res armcosmos.DatabaseAccountsClientListResponse) []*armcosmos.DatabaseAccountGetResults {
			return res.Value
		},
		func(res armcosmos.DatabaseAccountsClientListByResourceGroupResponse) []*armcosmos.DatabaseAccountGetResults {
			return res.Value
		},
		func(dbAccount *armcosmos.DatabaseAccountGetResults) error {
			// Get ActivityLogging for the CosmosDB
			activityLogging, raw, err := d.discoverDiagnosticSettings(util.Deref(dbAccount.ID))
			if err != nil {
				log.Error("could not discover diagnostic settings for the storage account: %w", err)
			}

			cosmos, err := d.handleCosmosDB(dbAccount, activityLogging, raw)
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

func (d *azureStorageDiscovery) handleCosmosDB(account *armcosmos.DatabaseAccountGetResults, activityLogging *voc.ActivityLogging, raw string) ([]voc.IsCloudResource, error) {
	var (
		atRestEnc           voc.IsAtRestEncryption
		err                 error
		list                []voc.IsCloudResource
		publicNetworkAccess = false
	)

	// initialize Cosmos DB client
	if err = d.initCosmosDBClient(); err != nil {
		return nil, err
	}

	// Check if KeyVaultURI is set for Cosmos DB account
	// By default the Cosmos DB account is encrypted by Azure managed keys. Optionally, it is possible to add a second encryption layer with customer key encryption. (see https://learn.microsoft.com/en-us/azure/cosmos-db/how-to-setup-customer-managed-keys?tabs=azure-portal)
	if account.Properties.KeyVaultKeyURI != nil {
		atRestEnc = &voc.CustomerKeyEncryption{
			AtRestEncryption: &voc.AtRestEncryption{
				Enabled: true,
				// Algorithm: algorithm, //TODO(anatheka): How do we get the algorithm? Are we available to do it by the related resources?
			},
			KeyUrl: util.Deref(account.Properties.KeyVaultKeyURI),
		}
		// Add key URI to keyUsage for tracking usages of single keys.
		// Hacky, but we haven't related Evidences yet.
		// We add the ID of this resource to the list of usages for the given key. But only if it is not there already
		if !slices.Contains(keyUsage[util.Deref(account.Properties.KeyVaultKeyURI)], util.Deref(account.ID)) {
			keyUsage[util.Deref(account.Properties.KeyVaultKeyURI)] =
				append(keyUsage[util.Deref(account.Properties.KeyVaultKeyURI)], util.Deref(account.ID))
		}
	} else {
		atRestEnc = &voc.ManagedKeyEncryption{
			AtRestEncryption: &voc.AtRestEncryption{
				Enabled:   true,
				Algorithm: AES256,
			},
		}
	}

	// Check if resource is public available
	if util.Deref(account.Properties.PublicNetworkAccess) == "Enabled" {
		publicNetworkAccess = true
	}

	// Create Cosmos DB database service voc object for the database account
	dbService := &voc.DatabaseService{
		StorageService: &voc.StorageService{
			NetworkService: &voc.NetworkService{
				Networking: &voc.Networking{
					Resource: discovery.NewResource(d,
						voc.ResourceID(resourceID(account.ID)),
						*account.Name,
						account.SystemData.CreatedAt,
						voc.GeoLocation{
							Region: *account.Location,
						},
						labels(account.Tags),
						resourceGroupID(account.ID),
						voc.DatabaseServiceType,
						account,
						raw,
					),
				},
			},
			ActivityLogging: activityLogging,
			Redundancy:      getCosmosDBRedundancy(account),
		},
		PublicAccess: publicNetworkAccess,
	}

	// Add Mongo DB database service
	list = append(list, dbService)

	// Check account kind and add Mongo DB databases storages
	switch util.Deref(account.Kind) {
	case armcosmos.DatabaseAccountKindMongoDB:
		// Get Mongo databases
		list = append(list, d.discoverMongoDBDatabases(account, atRestEnc)...)
	case armcosmos.DatabaseAccountKindGlobalDocumentDB:
		log.Infof("%s not yet implemented", armcosmos.DatabaseAccountKindGlobalDocumentDB)
	case armcosmos.DatabaseAccountKindParse:
		log.Infof("%s not yet implemented", armcosmos.DatabaseAccountKindParse)
	default:
		log.Warningf("Account kind '%s' not yet implemented", util.Deref(account.Kind))
	}

	return list, nil
}

// func getCosmosDBRedundancy(acc *armcosmos.DatabaseAccountGetResults) *voc.Redundancy {
// 	r := &voc.Redundancy{}
// 	locations := acc.Properties.Locations
// 	// If one location has zone redundancy enabled, we define the resource as zone redundant
// 	for _, l := range locations {
// 		if util.Deref(l.IsZoneRedundant) {
// 			r.Zone = true
// 		}
// 	}
// 	// If there are more than 1 region that means data is replicated geo-redundantly
// 	if len(locations) > 1 {
// 		r.Geo = true
// 	}
// 	return r
// }

// discoverSqlServers discovers the sql server and databases
func (d *azureStorageDiscovery) discoverSqlServers() ([]voc.IsCloudResource, error) {
	var (
		list []voc.IsCloudResource
		err  error
	)

	// initialize SQL server client
	if err := d.initSQLServersClient(); err != nil {
		return nil, err
	}

	// Discover sql server
	err = listPager(d.azureDiscovery,
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

func (d *azureStorageDiscovery) handleSqlServer(server *armsql.Server) ([]voc.IsCloudResource, error) {
	var (
		dbList []voc.IsCloudResource
		// anomalyDetectionList []voc.IsAnomalyDetection
		dbService           voc.IsCloudResource
		list                []voc.IsCloudResource
		publicNetworkAccess = false
	)

	// Get SQL database storages and the corresponding anomaly detection property
	dbList, _ = d.getSqlDBs(server)

	// Check if resource is public available
	if util.Deref(server.Properties.PublicNetworkAccess) == "Enabled" {
		publicNetworkAccess = true
	}

	// Create SQL database service voc object for SQL server
	dbService = &voc.DatabaseService{
		StorageService: &voc.StorageService{
			NetworkService: &voc.NetworkService{
				Networking: &voc.Networking{
					Resource: discovery.NewResource(d,
						voc.ResourceID(resourceID(server.ID)),
						*server.Name,
						nil, // creation time not available
						voc.GeoLocation{
							Region: *server.Location,
						},
						labels(server.Tags),
						resourceGroupID(server.ID),
						voc.DatabaseServiceType,
						server,
					),
				},
				// TODO(all): HttpEndpoint
				TransportEncryption: &voc.TransportEncryption{
					Enabled:    true,
					Enforced:   true,
					TlsVersion: tlsVersion(server.Properties.MinimalTLSVersion),
				},
			},
		},
		PublicAccess: publicNetworkAccess,
		// AnomalyDetection: anomalyDetectionList, //TODO(all): The anomaly detection changed, we have to update that.
	}

	// Add SQL database service
	list = append(list, dbService)

	// Add SQL database storages
	list = append(list, dbList...)

	return list, nil
}
func (d *azureStorageDiscovery) discoverStorageAccounts() ([]voc.IsCloudResource, error) {
	var (
		storageResourcesList []voc.IsCloudResource
	)

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

	// initialize table client
	if err := d.initTableStorageClient(); err != nil {
		return nil, err
	}

	// initialize diagnostic settings client
	if err := d.initDiagnosticsSettingsClient(); err != nil {
		return nil, err
	}

	// Discover backup vaults
	err := d.azureDiscovery.discoverBackupVaults()
	if err != nil {
		log.Errorf("could not discover backup vaults: %v", err)
	}

	// Discover object and file storages
	err = listPager(d.azureDiscovery,
		d.clients.accountsClient.NewListPager,
		d.clients.accountsClient.NewListByResourceGroupPager,
		func(res armstorage.AccountsClientListResponse) []*armstorage.Account {
			return res.Value
		},
		func(res armstorage.AccountsClientListByResourceGroupResponse) []*armstorage.Account {
			return res.Value
		},
		func(account *armstorage.Account) error {
			activityLoggingAccount, activityLoggingBlob, activityLoggingFile, activityLoggingTable, rawAccount, rawBlob, rawTable, rawFile := d.getActivityLogging(account)

			// Discover object storages
			objectStorages, err := d.discoverObjectStorages(account, activityLoggingBlob, rawBlob)
			if err != nil {
				return fmt.Errorf("could not handle object storages: %w", err)
			}

			// Discover file storages
			fileStorages, err := d.discoverFileStorages(account, activityLoggingFile, rawFile)
			if err != nil {
				return fmt.Errorf("could not handle file storages: %w", err)
			}

			// Discover file storages
			tableStorages, err := d.discoverTableStorages(account, activityLoggingTable, rawTable)
			if err != nil {
				return fmt.Errorf("could not handle table storages: %w", err)
			}

			storageResourcesList = append(storageResourcesList, objectStorages...)
			storageResourcesList = append(storageResourcesList, fileStorages...)
			storageResourcesList = append(storageResourcesList, tableStorages...)

			// Create storage service for all storage account resources
			storageService, err := d.handleStorageAccount(account, storageResourcesList, activityLoggingAccount, rawAccount)
			if err != nil {
				return fmt.Errorf("could not create storage service: %w", err)
			}

			storageResourcesList = append(storageResourcesList, storageService)

			return nil
		})
	if err != nil {
		return nil, err
	}

	// Add backuped storage account objects
	if d.backupMap[DataSourceTypeStorageAccountObject] != nil && d.backupMap[DataSourceTypeStorageAccountObject].backupStorages != nil {
		storageResourcesList = append(storageResourcesList, d.backupMap[DataSourceTypeStorageAccountObject].backupStorages...)
	}

	return storageResourcesList, nil
}

func (d *azureStorageDiscovery) getActivityLogging(account *armstorage.Account) (activityLoggingAccount, activityLoggingBlob, activityLoggingTable, activityLoggingFile *voc.ActivityLogging, rawAccount, rawBlob, rawTable, rawFile string) {

	var err error

	// Get ActivityLogging for the storage account
	activityLoggingAccount, rawAccount, err = d.discoverDiagnosticSettings(util.Deref(account.ID))
	if err != nil {
		log.Errorf("could not discover diagnostic settings for the storage account: %v", err)
	}

	// Get ActivityLogging for the blob service
	activityLoggingBlob, rawBlob, err = d.discoverDiagnosticSettings(util.Deref(account.ID) + "/blobServices/default")
	if err != nil {
		log.Errorf("could not discover diagnostic settings for the blob service: %v", err)
	}

	// Get ActivityLogging for the table service
	activityLoggingTable, rawTable, err = d.discoverDiagnosticSettings(util.Deref(account.ID) + "/tableServices/default")
	if err != nil {
		log.Errorf("could not discover diagnostic settings for the table service: %v", err)
	}

	// Get ActivityLogging for the file service
	activityLoggingFile, rawFile, err = d.discoverDiagnosticSettings(util.Deref(account.ID) + "/fileServices/default")
	if err != nil {
		log.Errorf("could not discover diagnostic settings for the file service: %v", err)
	}

	return

}

func (d *azureStorageDiscovery) discoverFileStorages(account *armstorage.Account, activityLogging *voc.ActivityLogging, raw string) ([]voc.IsCloudResource, error) {
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
			fileStorages, err := d.handleFileStorage(account, value, activityLogging, raw)
			if err != nil {
				return nil, fmt.Errorf("could not handle file storage: %w", err)
			}

			log.Infof("Adding file storage '%s", fileStorages.Name)

			list = append(list, fileStorages)
		}
	}

	return list, nil
}

func (d *azureStorageDiscovery) discoverObjectStorages(account *armstorage.Account, activityLogging *voc.ActivityLogging, raw string) ([]voc.IsCloudResource, error) {
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
			objectStorages, err := d.handleObjectStorage(account, value, activityLogging)
			if err != nil {
				return nil, fmt.Errorf("could not handle object storage: %w", err)
			}
			log.Infof("Adding object storage '%s'", objectStorages.Name)
			list = append(list, objectStorages)

			// Add objects as well (mainly for UI rendering and connections, to avoid UI bugs)
			objects, err := d.handleObjects(account, value, raw)
			if err != nil {
				// We don't quit here since it is not crucial to have objects.
				log.Warnf("could not handle objects of object storage: %s", err.Error())
				objects = []*voc.Object{}
			}
			for _, o := range objects {
				list = append(list, o)
			}
		}
	}

	return list, nil
}

func (d *azureStorageDiscovery) discoverTableStorages(account *armstorage.Account, activityLogging *voc.ActivityLogging, raw string) ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	// List all blob containers in the specified resource group
	listPager := d.clients.tableStorageClient.NewListPager(resourceGroupName(util.Deref(account.ID)), util.Deref(account.Name), &armstorage.TableClientListOptions{})
	for listPager.More() {
		pageResponse, err := listPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("%s: %v", ErrGettingNextPage, err)
			return nil, err
		}

		for _, value := range pageResponse.Value {
			tableStorage, err := d.handleTableStorage(account, value, activityLogging, raw)
			if err != nil {
				return nil, fmt.Errorf("could not handle table storage: %w", err)
			}
			log.Infof("Adding table storage '%s'", tableStorage.Name)

			list = append(list, tableStorage)

		}
	}

	return list, nil
}

func (d *azureStorageDiscovery) handleStorageAccount(account *armstorage.Account, storagesList []voc.IsCloudResource, activityLogging *voc.ActivityLogging, raw string) (*voc.ObjectStorageService, error) {
	var (
		storageResourceIDs []voc.ResourceID
	)

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
		TlsVersion: string(util.Deref(account.Properties.MinimumTLSVersion)),
		Algorithm:  constants.TLS,
	}

	storageService := &voc.ObjectStorageService{
		StorageService: &voc.StorageService{
			Storage: storageResourceIDs,
			NetworkService: &voc.NetworkService{
				Networking: &voc.Networking{
					Resource: discovery.NewResource(d,
						voc.ResourceID(resourceID(account.ID)),
						util.Deref(account.Name),
						account.Properties.CreationTime,
						voc.GeoLocation{
							Region: util.Deref(account.Location),
						},
						labels(account.Tags),
						resourceGroupID(account.ID),
						voc.ObjectStorageServiceType,
						account, raw,
					),
				},
				TransportEncryption: te,
			},
			Redundancy:      getStorageAccountRedundancy(account),
			ActivityLogging: activityLogging,
		},

		HttpEndpoint: &voc.HttpEndpoint{
			Url:                 generalizeURL(util.Deref(account.Properties.PrimaryEndpoints.Blob)),
			TransportEncryption: te,
		},
		PublicAccess: getPublicAccessOfStorageAccount(account),
	}

	return storageService, nil
}

func getPublicAccessOfStorageAccount(acc *armstorage.Account) bool {
	return util.Deref(acc.Properties.PublicNetworkAccess) == "Enabled"
}

func (d *azureStorageDiscovery) handleFileStorage(account *armstorage.Account, fileshare *armstorage.FileShareItem, activityLogging *voc.ActivityLogging, raw string) (*voc.FileStorage, error) {
	var (
		monitoringLogDataEnabled bool
		securityAlertsEnabled    bool
	)

	if account == nil {
		return nil, ErrEmptyStorageAccount
	}

	// It is possible that the fileshare is not empty. In that case we have to check if a mandatory field is empty, so the whole disk is empty
	if fileshare == nil || fileshare.ID == nil {
		return nil, fmt.Errorf("fileshare is nil")
	}

	// Get atRestEncryptionEnabled
	enc, err := storageAtRestEncryption(account)
	if err != nil {
		return nil, fmt.Errorf("could not get file storage properties for the atRestEncryption: %w", err)
	}

	// Get monitoringLogDataEnabled and securityAlertsEnabled
	if d.defenderProperties[DefenderStorageType] != nil {
		monitoringLogDataEnabled = d.defenderProperties[DefenderVirtualMachineType].monitoringLogDataEnabled
		securityAlertsEnabled = d.defenderProperties[DefenderVirtualMachineType].securityAlertsEnabled
	}

	return &voc.FileStorage{
		Storage: &voc.Storage{
			Resource: discovery.NewResource(d,
				voc.ResourceID(resourceID(fileshare.ID)),
				util.Deref(fileshare.Name),
				// We only have the creation time of the storage account the file storage belongs to
				account.Properties.CreationTime,
				voc.GeoLocation{
					// The location is the same as the storage account
					Region: util.Deref(account.Location),
				},
				// The storage account labels the file storage belongs to
				labels(account.Tags),
				// the storage account is our parent
				voc.ResourceID(resourceID(account.ID)),
				voc.FileStorageType,
				account, fileshare, raw,
			),
			ResourceLogging: &voc.ResourceLogging{
				Logging: &voc.Logging{
					MonitoringLogDataEnabled: monitoringLogDataEnabled,
					SecurityAlertsEnabled:    securityAlertsEnabled,
				},
			},
			ActivityLogging:  activityLogging,
			AtRestEncryption: enc,
			Redundancy:       getStorageAccountRedundancy(account),
		},
	}, nil
}

func (d *azureStorageDiscovery) handleObjectStorage(account *armstorage.Account, container *armstorage.ListContainerItem, activityLogging *voc.ActivityLogging) (*voc.ObjectStorage, error) {
	var (
		backups                  []*voc.Backup
		monitoringLogDataEnabled bool
		securityAlertsEnabled    bool
	)

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

	if d.backupMap[DataSourceTypeStorageAccountObject] != nil && d.backupMap[DataSourceTypeStorageAccountObject].backup[util.Deref(account.ID)] != nil {
		backups = d.backupMap[DataSourceTypeStorageAccountObject].backup[util.Deref(account.ID)]
	}
	backups = backupsEmptyCheck(backups)

	if d.defenderProperties[DefenderStorageType] != nil {
		monitoringLogDataEnabled = d.defenderProperties[DefenderVirtualMachineType].monitoringLogDataEnabled
		securityAlertsEnabled = d.defenderProperties[DefenderVirtualMachineType].securityAlertsEnabled
	}

	// Check if container is acting as a backups. If so, they are also added to backupOf
	isBackup := d.isBackup(account, container)

	return &voc.ObjectStorage{
		Storage: &voc.Storage{
			Resource: discovery.NewResource(d,
				voc.ResourceID(resourceID(container.ID)),
				util.Deref(container.Name),
				// We only have the creation time of the storage account the object storage belongs to
				account.Properties.CreationTime,
				voc.GeoLocation{
					// The location is the same as the storage account
					Region: util.Deref(account.Location),
				},
				// The storage account labels the object storage belongs to
				labels(account.Tags),
				// the storage account is our parent
				voc.ResourceID(resourceID(account.ID)),
				voc.ObjectStorageType,
				account, container,
			),
			AtRestEncryption: enc,
			Immutability: &voc.Immutability{
				Enabled: util.Deref(container.Properties.HasImmutabilityPolicy),
			},
			ResourceLogging: &voc.ResourceLogging{
				Logging: &voc.Logging{
					MonitoringLogDataEnabled: monitoringLogDataEnabled,
					SecurityAlertsEnabled:    securityAlertsEnabled,
				},
			},
			ActivityLogging: activityLogging,
			Backups:         backups,
			// Todo(lebogg): Add tests
			Redundancy: getStorageAccountRedundancy(account),
		},
		ContainerPublicAccess: util.Deref(container.Properties.PublicAccess) != armstorage.PublicAccessNone, // This is not the public network access https://learn.microsoft.com/en-us/azure/storage/common/storage-network-security?tabs=azure-portal, but the container public access
		IsBackup:              isBackup,
	}, nil
}

// isBackup checks if container is used as a backup - and metadata.  If so, it is
// added to backupOf
func (d *azureStorageDiscovery) isBackup(account *armstorage.Account, container *armstorage.ListContainerItem) bool {
	// Get specific container with mor details, e.g. meta data
	res, err := d.clients.blobContainerClient.Get(context.Background(), resourceGroupName(util.Deref(account.ID)),
		util.Deref(account.Name), util.Deref(container.Name), &armstorage.BlobContainersClientGetOptions{})
	if err != nil {
		log.Warnf("Error while retrieving container '%s' to find out if it is a backup", util.Deref(container.Name))
		return false
	}

	// Check if the container serves as backup
	if b, ok := res.ContainerProperties.Metadata["backupOf"]; ok {
		backupOf[resourceID(b)] = resourceID(container.ID)
		return true
	}
	return false
}

func (d *azureStorageDiscovery) handleTableStorage(account *armstorage.Account, table *armstorage.Table, activityLogging *voc.ActivityLogging, raw string) (*voc.DatabaseStorage, error) {
	var (
		backups                  []*voc.Backup
		monitoringLogDataEnabled bool
		securityAlertsEnabled    bool
	)

	if account == nil {
		return nil, ErrEmptyStorageAccount
	}

	// It is possible that the table is empty. In that case we have to check if a mandatory field is empty, so the whole disk is empty
	if table == nil || table.ID == nil {
		return nil, fmt.Errorf("table is nil")
	}

	enc, err := storageAtRestEncryption(account)
	if err != nil {
		return nil, fmt.Errorf("could not get object storage properties for the atRestEncryption: %w", err)
	}
	if d.backupMap[DataSourceTypeStorageAccountObject] != nil && d.backupMap[DataSourceTypeStorageAccountObject].backup[util.Deref(account.ID)] != nil {
		backups = d.backupMap[DataSourceTypeStorageAccountObject].backup[util.Deref(account.ID)]
	} else { // approach with Tagging
		if backupLocation, ok := backupOf["https://"+resourceID(account.Name)+".table.core.windows.net/"+resourceID(table.Name)]; ok {
			backups = []*voc.Backup{
				{
					Availability:        nil,
					TransportEncryption: nil,
					Storage:             voc.ResourceID(backupLocation),
					Enabled:             true,
					RetentionPeriod:     0,
					Interval:            0,
				},
			}
		}
	}
	backups = backupsEmptyCheck(backups)

	if d.defenderProperties[DefenderStorageType] != nil {
		monitoringLogDataEnabled = d.defenderProperties[DefenderVirtualMachineType].monitoringLogDataEnabled
		securityAlertsEnabled = d.defenderProperties[DefenderVirtualMachineType].securityAlertsEnabled
	}

	return &voc.DatabaseStorage{
		Storage: &voc.Storage{
			Resource: discovery.NewResource(d,
				voc.ResourceID(resourceID(table.ID)),
				util.Deref(table.Name),
				// We only have the creation time of the storage account the object storage belongs to
				account.Properties.CreationTime,
				voc.GeoLocation{
					// The location is the same as the storage account
					Region: util.Deref(account.Location),
				},
				// The storage account labels the object storage belongs to
				labels(account.Tags),
				// the storage account is our parent
				voc.ResourceID(resourceID(account.ID)),
				voc.DatabaseStorageType,
				account, table, raw,
			),
			AtRestEncryption: enc,
			Backups:          backups,
			Immutability:     nil, // TODO
			ResourceLogging: &voc.ResourceLogging{
				Logging: &voc.Logging{
					MonitoringLogDataEnabled: monitoringLogDataEnabled,
					SecurityAlertsEnabled:    securityAlertsEnabled,
				},
			},
			ActivityLogging: activityLogging,
			Redundancy:      getStorageAccountRedundancy(account),
		},
	}, nil
}

// TODO(lebogg): Add tests
func getStorageAccountRedundancy(account *armstorage.Account) (r *voc.Redundancy) {
	r = new(voc.Redundancy)
	name := util.Deref(account.SKU.Name)
	switch name {
	// LRS denotes local redundancy
	case armstorage.SKUNameStandardLRS, armstorage.SKUNamePremiumLRS:
		r.Local = true
	// ZRS denotes zone redundancy
	case armstorage.SKUNameStandardZRS, armstorage.SKUNamePremiumZRS:
		r.Zone = true
	// GRS denotes geo redundancy which also includes local redundancy in Azure
	case armstorage.SKUNameStandardGRS, armstorage.SKUNameStandardRAGRS:
		r.Local = true
		r.Geo = true
	// GZRS denotes geo redundancy + zone redundancy
	case armstorage.SKUNameStandardGZRS, armstorage.SKUNameStandardRAGZRS:
		// r.Local = true // local redundancy only in secondary location. TODO(all): Discuss all options
		r.Zone = true
		r.Geo = true
	// When there are new SKU types in the future we will probably miss it. Print out a warning if there is a name we
	// don't consider so far.
	default:
		log.Warnf("Unknown redundancy model (via SKU) for storage account '%s': '%s'. Probably, we should add it.",
			util.Deref(account.SKU.Name), name)
		// consideredAccountTypes shows how many account types (SKUs) we consider so far. It has to be a "magic" number.
		consideredAccountTypes := 8
		log.Warnf("Currently there are %d different SKU types. We consider %d types so far",
			len(armstorage.PossibleSKUNameValues()), consideredAccountTypes)
	}
	return
}

// storageAtRestEncryption takes encryption properties of an armstorage.Account and converts it into our respective
// ontology object.
func storageAtRestEncryption(account *armstorage.Account) (enc voc.IsAtRestEncryption, err error) {
	if account == nil {
		return enc, ErrEmptyStorageAccount
	}

	if account.Properties == nil || account.Properties.Encryption.KeySource == nil {
		return enc, errors.New("keySource is empty")
	} else if util.Deref(account.Properties.Encryption.KeySource) == armstorage.KeySourceMicrosoftStorage {
		enc = &voc.ManagedKeyEncryption{
			AtRestEncryption: &voc.AtRestEncryption{
				Algorithm: AES256,
				Enabled:   true,
			},
		}
	} else if util.Deref(account.Properties.Encryption.KeySource) == armstorage.KeySourceMicrosoftKeyvault {
		enc = &voc.CustomerKeyEncryption{
			AtRestEncryption: &voc.AtRestEncryption{
				Algorithm: "", // TODO(all): TBD
				Enabled:   true,
			},
			KeyUrl: util.Deref(account.Properties.Encryption.KeyVaultProperties.KeyVaultURI),
		}
	}

	return enc, nil
}

// anomalyDetectionEnabled returns true if Azure Advanced Threat Protection is enabled for the database.
func (d *azureStorageDiscovery) anomalyDetectionEnabled(server *armsql.Server, db *armsql.Database) (bool, error) {
	// initialize threat protection client
	if err := d.initThreatProtectionClient(); err != nil {
		return false, err
	}

	listPager := d.clients.threatProtectionClient.NewListByDatabasePager(resourceGroupName(util.Deref(db.ID)), *server.Name, *db.Name, &armsql.DatabaseAdvancedThreatProtectionSettingsClientListByDatabaseOptions{})
	for listPager.More() {
		pageResponse, err := listPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("%s: %v", ErrGettingNextPage, err)
			return false, err
		}

		for _, value := range pageResponse.Value {
			if *value.Properties.State == armsql.AdvancedThreatProtectionStateEnabled {
				return true, nil
			}
		}
	}
	return false, nil
}

// getSqlDBs returns a list of SQL databases for a specific SQL account
func (d *azureStorageDiscovery) getSqlDBs(server *armsql.Server) ([]voc.IsCloudResource, []voc.IsAnomalyDetection) {
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
				Scope:   voc.ResourceID(resourceID(value.ID)),
				Enabled: anomalyDetectionEnabled,
			}

			anomalyDetectionList = append(anomalyDetectionList, a)

			// Create database storage voc object
			sqlDB := &voc.DatabaseStorage{
				Storage: &voc.Storage{
					Resource: discovery.NewResource(d,
						voc.ResourceID(resourceID(value.ID)),
						*value.Name,
						value.Properties.CreationDate,
						voc.GeoLocation{
							Region: *value.Location,
						},
						labels(value.Tags),
						voc.ResourceID(resourceID(server.ID)),
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

func (d *azureStorageDiscovery) initTableStorageClient() (err error) {
	d.clients.tableStorageClient, err = initClient(d.clients.tableStorageClient, d.azureDiscovery, armstorage.NewTableClient)
	return
}

// initDefenderClient creates the client if not already exists
func (d *azureStorageDiscovery) initDefenderClient() (err error) {
	d.clients.defenderClient, err = initClient(d.clients.defenderClient, d.azureDiscovery, armsecurity.NewPricingsClient)

	return
}

// initBackupPoliciesClient creates the client if not already exists
func (d *azureStorageDiscovery) initBackupPoliciesClient() (err error) {
	d.clients.backupPoliciesClient, err = initClient(d.clients.backupPoliciesClient, d.azureDiscovery, armdataprotection.NewBackupPoliciesClient)

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

// initDatabasesClient creates the client if not already exists
func (d *azureStorageDiscovery) initDatabasesClient() (err error) {
	d.clients.databasesClient, err = initClient(d.clients.databasesClient, d.azureDiscovery, armsql.NewDatabasesClient)

	return
}

// initSQLServersClient creates the client if not already exists
func (d *azureStorageDiscovery) initSQLServersClient() (err error) {
	d.clients.sqlServersClient, err = initClient(d.clients.sqlServersClient, d.azureDiscovery, armsql.NewServersClient)

	return
}

// initCosmosDBClient creates the client if not already exists
func (d *azureStorageDiscovery) initCosmosDBClient() (err error) {
	d.clients.cosmosDBClient, err = initClient(d.clients.cosmosDBClient, d.azureDiscovery, armcosmos.NewDatabaseAccountsClient)

	return
}

// initThreatProtectionClient creates the client if not already exists
func (d *azureStorageDiscovery) initThreatProtectionClient() (err error) {
	d.clients.threatProtectionClient, err = initClient(d.clients.threatProtectionClient, d.azureDiscovery, armsql.NewDatabaseAdvancedThreatProtectionSettingsClient)

	return
}

// initMongoDBResourcesClient creates the client if not already exists
func (d *azureDiscovery) initMongoDBResourcesClient() (err error) {
	d.clients.mongoDBResourcesClient, err = initClient(d.clients.mongoDBResourcesClient, d, armcosmos.NewMongoDBResourcesClient)

	return
}

// handleObjects returns all objects of a container. It also checks if single Objects are backups and add these to the
// backupOf map accordingly.
func (d *azureStorageDiscovery) handleObjects(acc *armstorage.Account, container *armstorage.ListContainerItem, raw string) (objects []*voc.Object, err error) {
	// Get blobs and check their tags + metadata to check if there are backups
	var (
		client *azblob.Client
		// Determines if the given blob is a backup
		isBackup bool
	)

	client, err = azblob.NewClient(util.Deref(acc.Properties.PrimaryEndpoints.Blob), d.cred, nil)
	if err != nil {
		return nil, fmt.Errorf("could not creat azblob client: %v", err)
	}
	pager := client.NewListBlobsFlatPager(util.Deref(container.Name), &azblob.ListBlobsFlatOptions{
		Include: azblob.ListBlobsInclude{Tags: true, Metadata: true},
	})
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("could not load next page, probably you do not have the right "+
				"permissions ('Storage Blob Data Contributor' and 'Reader' roles are needed at least): %v", err)
		}
		// If there is no segment (although it is required) continue with next page
		if page.Segment == nil {
			continue
		}
		for _, blobItem := range page.Segment.BlobItems {

			// Label and Backup Stuff
			blobLabels := make(map[string]string)
			if blobItem.BlobTags != nil {
				for _, t := range blobItem.BlobTags.BlobTagSet {
					k := util.Deref(t.Key)
					v := util.Deref(t.Value)
					// Add to blob labels
					blobLabels[k] = v
					// Check backup label
					if k == "backupOf" {
						isBackup = true
						backupOf[v] =
							"https://" + util.Deref(acc.Name) + ".blob.core.windows.net/" +
								util.Deref(container.Name) + "/" + util.Deref(blobItem.Name)
					}
				}

			}
			// This can potentially overwrite a 'backupOf' defined in Tags, but they should be the same.
			if blobItem.Metadata != nil {
				for k, value := range blobItem.Metadata {
					v := util.Deref(value)
					// Add to blob labels
					blobLabels[k] = v
					// Check backup lable
					if k == "backupof" { // only lowercase
						isBackup = true
						backupOf[v] =
							"https://" + util.Deref(acc.Name) + ".blob.core.windows.net/" +
								util.Deref(container.Name) + "/" + util.Deref(blobItem.Name)
					}
				}
			}

			// Add resource to list
			objects = append(objects, &voc.Object{
				Resource: discovery.NewResource(d,
					voc.ResourceID("https://"+resourceID(acc.Name)+".blob.core.windows.net/"+
						resourceID(container.Name)+"/"+resourceID(blobItem.Name)),
					util.Deref(blobItem.Name),
					// We only have the creation time of the storage account the object storage belongs to
					acc.Properties.CreationTime,
					voc.GeoLocation{
						// The location is the same as the storage account
						Region: util.Deref(acc.Location),
					},
					blobLabels,
					// the storage account is our parent
					voc.ResourceID(resourceID(container.ID)),
					voc.ObjectType,
					container, blobItem, raw,
				),
				IsBackup: isBackup,
			})
		}
	}
	return
}

// discoverMongoDBDatabases returns a list of Mongo DB databases for a specific Mongo DB account
func (d *azureStorageDiscovery) discoverMongoDBDatabases(account *armcosmos.DatabaseAccountGetResults, atRestEnc voc.IsAtRestEncryption) []voc.IsCloudResource {
	var (
		list []voc.IsCloudResource
		err  error
	)

	// initialize Mongo DB resources client
	if err = d.initMongoDBResourcesClient(); err != nil {
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
						voc.ResourceID(resourceID(value.ID)),
						util.Deref(value.Name),
						nil, // creation time of database not available
						voc.GeoLocation{
							Region: *value.Location,
						},
						labels(value.Tags),
						voc.ResourceID(resourceID(account.ID)),
						voc.DatabaseStorageType,
						account,
						value),

					AtRestEncryption: atRestEnc,
					Redundancy:       nil, // Redundancy is done over database service (Cosmos DB)
					ActivityLogging:  nil, // ActivityLogging is done over database service (Cosmos DB)
				},
			}
			list = append(list, mongoDB)
		}
	}

	return list
}

// getCosmosDBRedundancy returns for a given cosmos DB account the redundancy object in the voc format. Currently, only
// zone redundancy is supported
func getCosmosDBRedundancy(account *armcosmos.DatabaseAccountGetResults) (r *voc.Redundancy) {
	r = &voc.Redundancy{}
	locations := account.Properties.Locations
	for _, l := range locations {
		r.Zone = util.Deref(l.IsZoneRedundant)
	}
	if len(locations) > 1 {
		r.Geo = true
		r.Zone = true // Probably, we don't want that in main this way
	}
	return
}
