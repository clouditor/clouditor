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
	"clouditor.io/clouditor/internal/constants"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dataprotection/armdataprotection"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/security/armsecurity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

var (
	ErrEmptyStorageAccount        = errors.New("storage account is empty")
	ErrMissingDiskEncryptionSetID = errors.New("no disk encryption set ID was specified")
	ErrBackupStorageNotAvailable  = errors.New("backup storages not available")
)

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
		return nil, fmt.Errorf("could not discover sql databases: %w", err)
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
			cosmos, err := d.handleCosmosDB(dbAccount)
			if err != nil {
				return fmt.Errorf("could not cosmos db accounts: %w", err)
			}
			log.Infof("Adding Cosmos DB account '%s", *dbAccount.Name)
			list = append(list, cosmos)

			return nil
		})
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (d *azureStorageDiscovery) handleCosmosDB(account *armcosmos.DatabaseAccountGetResults) (voc.IsCloudResource, error) {
	var (
		enc voc.IsAtRestEncryption
		err error
	)

	// initialize Cosmos DB client
	if err = d.initCosmosDBClient(); err != nil {
		return nil, err
	}

	// Check if KeyVaultURI is set
	// By default the Cosmos DB account is encrypted by Azure managed keys. Optionally, it is possible to add a second encryption layer with customer key encryption. (see https://learn.microsoft.com/en-us/azure/cosmos-db/how-to-setup-customer-managed-keys?tabs=azure-portal)
	if account.Properties.KeyVaultKeyURI != nil {
		enc = &voc.CustomerKeyEncryption{
			AtRestEncryption: &voc.AtRestEncryption{
				Enabled: true,
				// Algorithm: algorithm, //TODO(anatheka): How do we get the algorithm? Are we available to do it by the related resources?
			},
			KeyUrl: util.Deref(account.Properties.KeyVaultKeyURI),
		}
	} else {
		enc = &voc.ManagedKeyEncryption{
			AtRestEncryption: &voc.AtRestEncryption{
				Enabled:   true,
				Algorithm: AES256,
			},
		}
	}

	// Create Cosmos DB database account voc object
	dbStorage := &voc.DatabaseStorage{
		Storage: &voc.Storage{
			Resource: discovery.NewResource(d,
				voc.ResourceID(*account.ID),
				util.Deref(account.Name),
				account.SystemData.CreatedAt,
				voc.GeoLocation{
					Region: *account.Location,
				},
				labels(account.Tags),
				resourceGroupID(account.ID),
				voc.DatabaseStorageType,
				account),

			AtRestEncryption: enc,
		},
	}

	return dbStorage, nil
}

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
		dbStorage voc.IsCloudResource
		dbService voc.IsCloudResource
		list      []voc.IsCloudResource
		err       error
	)

	// initialize SQL databases client
	if err = d.initDatabasesClient(); err != nil {
		return nil, err
	}

	// Get databases for given server
	serverlistPager := d.clients.databasesClient.NewListByServerPager(resourceGroupName(util.Deref(server.ID)), *server.Name, &armsql.DatabasesClientListByServerOptions{})
	for serverlistPager.More() {
		pageResponse, err := serverlistPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("%s: %v", ErrGettingNextPage, err)
			return nil, err
		}

		for _, value := range pageResponse.Value {
			// Getting anomaly detection status
			anomalyDetectionEnabeld, err := d.anomalyDetectionEnabled(server, value)
			if err != nil {
				log.Errorf("error getting anomaly detection info for database '%s': %v", *value.Name, err)
			}

			// Create database service voc object
			//
			// TODO(oxisto): This is not 100 % accurate. According to our ontology definition, the SQL server would be
			// the database service and individual databases would be DatabaseStorage objects. However, the problem is
			// that azure defines anomaly detection on a per-database level and we currently have anomaly detection as
			// part of the service
			dbService = &voc.DatabaseService{
				StorageService: &voc.StorageService{
					NetworkService: &voc.NetworkService{
						Networking: &voc.Networking{
							Resource: discovery.NewResource(d,
								voc.ResourceID(*value.ID),
								*value.Name,
								value.Properties.CreationDate,
								voc.GeoLocation{
									Region: *value.Location,
								},
								labels(value.Tags),
								resourceGroupID(value.ID),
								voc.DatabaseServiceType,
								server,
								value,
							),
						},
						// TODO(all): TransportEncryption, HttpEndpoint
					},
				},
				AnomalyDetection: &voc.AnomalyDetection{
					Enabled: anomalyDetectionEnabeld,
				},
			}

			list = append(list, dbService)

			// Create database storage voc object
			dbStorage = &voc.DatabaseStorage{
				Storage: &voc.Storage{
					Resource: discovery.NewResource(d,
						voc.ResourceID(*value.ID),
						*value.Name,
						value.Properties.CreationDate,
						voc.GeoLocation{
							Region: *value.Location,
						},
						labels(value.Tags),
						// the DB service is our parent
						dbService.GetID(),
						voc.DatabaseStorageType,
						value),
					AtRestEncryption: &voc.AtRestEncryption{
						Enabled:   *value.Properties.IsInfraEncryptionEnabled,
						Algorithm: AES256,
					},
					// TODO(all): Backups
				},
			}

			list = append(list, dbStorage)
		}
	}
	return list, nil
}

func (d *azureStorageDiscovery) discoverStorageAccounts() ([]voc.IsCloudResource, error) {
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

	// Add backuped storage account objects
	if d.backupMap[DataSourceTypeStorageAccountObject] != nil && d.backupMap[DataSourceTypeStorageAccountObject].backupStorages != nil {
		storageResourcesList = append(storageResourcesList, d.backupMap[DataSourceTypeStorageAccountObject].backupStorages...)
	}

	return storageResourcesList, nil
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

func (d *azureStorageDiscovery) handleStorageAccount(account *armstorage.Account, storagesList []voc.IsCloudResource) (*voc.ObjectStorageService, error) {
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
						voc.ResourceID(util.Deref(account.ID)),
						util.Deref(account.Name),
						account.Properties.CreationTime,
						voc.GeoLocation{
							Region: util.Deref(account.Location),
						},
						labels(account.Tags),
						resourceGroupID(account.ID),
						voc.ObjectStorageServiceType,
						account,
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

func (d *azureStorageDiscovery) handleFileStorage(account *armstorage.Account, fileshare *armstorage.FileShareItem) (*voc.FileStorage, error) {
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
				// the storage account is our parent
				voc.ResourceID(util.Deref(account.ID)),
				voc.FileStorageType,
				account, fileshare,
			),
			ResourceLogging: &voc.ResourceLogging{
				Logging: &voc.Logging{
					MonitoringLogDataEnabled: monitoringLogDataEnabled,
					SecurityAlertsEnabled:    securityAlertsEnabled,
				},
			},
			AtRestEncryption: enc,
		},
	}, nil
}

func (d *azureStorageDiscovery) handleObjectStorage(account *armstorage.Account, container *armstorage.ListContainerItem) (*voc.ObjectStorage, error) {
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
				// the storage account is our parent
				voc.ResourceID(util.Deref(account.ID)),
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
			Backups: backups,
		},
		PublicAccess: util.Deref(container.Properties.PublicAccess) != armstorage.PublicAccessNone,
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
