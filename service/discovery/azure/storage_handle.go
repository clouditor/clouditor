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
	"fmt"
	"strings"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/constants"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

func (d *azureDiscovery) handleCosmosDB(account *armcosmos.DatabaseAccountGetResults) ([]voc.IsCloudResource, error) {
	var (
		atRestEnc voc.IsAtRestEncryption
		err       error
		list      []voc.IsCloudResource
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
	} else {
		atRestEnc = &voc.ManagedKeyEncryption{
			AtRestEncryption: &voc.AtRestEncryption{
				Enabled:   true,
				Algorithm: AES256,
			},
		}
	}

	// Create Cosmos DB database service voc object for the database account
	dbService := &voc.DatabaseService{
		StorageService: &voc.StorageService{
			NetworkService: &voc.NetworkService{
				Networking: &voc.Networking{
					Resource: discovery.NewResource(d,
						voc.ResourceID(*account.ID),
						*account.Name,
						account.SystemData.CreatedAt,
						voc.GeoLocation{
							Region: *account.Location,
						},
						labels(account.Tags),
						resourceGroupID(account.ID),
						voc.DatabaseServiceType,
						account,
					),
				},
			},
		},
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

func (d *azureDiscovery) handleSqlServer(server *armsql.Server) ([]voc.IsCloudResource, error) {
	var (
		dbList               []voc.IsCloudResource
		anomalyDetectionList []voc.IsAnomalyDetection
		dbService            voc.IsCloudResource
		list                 []voc.IsCloudResource
	)

	// Get SQL database storages and the corresponding anomaly detection property
	dbList, anomalyDetectionList = d.getSqlDBs(server)

	// Create SQL database service voc object for SQL server
	dbService = &voc.DatabaseService{
		StorageService: &voc.StorageService{
			NetworkService: &voc.NetworkService{
				Networking: &voc.Networking{
					Resource: discovery.NewResource(d,
						voc.ResourceID(*server.ID),
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
					TlsVersion: tlsVersion(util.Deref(server.Properties.MinimalTLSVersion)),
				},
			},
		},
		AnomalyDetection: anomalyDetectionList,
	}

	// Add SQL database service
	list = append(list, dbService)

	// Add SQL database storages
	list = append(list, dbList...)

	return list, nil
}

func (d *azureDiscovery) handleStorageAccount(account *armstorage.Account, storagesList []voc.IsCloudResource) (*voc.ObjectStorageService, error) {
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
		TlsVersion: tlsVersion((string(util.Deref(account.Properties.MinimumTLSVersion)))),
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

func (d *azureDiscovery) handleFileStorage(account *armstorage.Account, fileshare *armstorage.FileShareItem) (*voc.FileStorage, error) {
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

func (d *azureDiscovery) handleObjectStorage(account *armstorage.Account, container *armstorage.ListContainerItem) (*voc.ObjectStorage, error) {
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
