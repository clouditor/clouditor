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

package azure

import (
	"fmt"
	"strings"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/constants"
	"clouditor.io/clouditor/v2/internal/util"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

func (d *azureDiscovery) handleCosmosDB(account *armcosmos.DatabaseAccountGetResults) ([]ontology.IsResource, error) {
	var (
		atRestEnc *ontology.AtRestEncryption
		err       error
		list      []ontology.IsResource
	)

	// initialize Cosmos DB client
	if err = d.initCosmosDBClient(); err != nil {
		return nil, err
	}

	// Check if KeyVaultURI is set for Cosmos DB account By default the Cosmos DB account is encrypted by Azure managed
	// keys. Optionally, it is possible to add a second encryption layer with customer key encryption. (see
	// https://learn.microsoft.com/en-us/azure/cosmos-db/how-to-setup-customer-managed-keys?tabs=azure-portal)
	if account.Properties.KeyVaultKeyURI != nil {
		atRestEnc = &ontology.AtRestEncryption{
			Type: &ontology.AtRestEncryption_CustomerKeyEncryption{
				CustomerKeyEncryption: &ontology.CustomerKeyEncryption{
					Enabled: true,
					// Algorithm: algorithm, //TODO(anatheka): How do we get the algorithm? Are we available to do it by
					// the related resources?
					KeyUrl: util.Deref(account.Properties.KeyVaultKeyURI),
				},
			},
		}
	} else {
		atRestEnc = &ontology.AtRestEncryption{
			Type: &ontology.AtRestEncryption_ManagedKeyEncryption{
				ManagedKeyEncryption: &ontology.ManagedKeyEncryption{
					Enabled:   true,
					Algorithm: constants.AES256,
				},
			},
		}
	}

	// Create Cosmos DB database service ontology object for the database account
	//
	// TODO(oxisto): Actually, CosmosDB is a multi-model database, but for now we just model this as a document
	// database.
	dbService := &ontology.DocumentDatabaseService{
		Id:           resourceID(account.ID),
		Name:         util.Deref(account.Name),
		CreationTime: creationTime(account.SystemData.CreatedAt),
		GeoLocation:  location(account.Location),
		Labels:       labels(account.Tags),
		ParentId:     resourceGroupID(account.ID),
		Raw:          discovery.Raw(account),
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

func (d *azureDiscovery) handleSqlServer(server *armsql.Server) ([]ontology.IsResource, error) {
	var (
		dbList               []ontology.IsResource
		anomalyDetectionList []*ontology.AnomalyDetection
		dbService            ontology.IsResource
		list                 []ontology.IsResource
	)

	// Get SQL database storages and the corresponding anomaly detection property
	dbList, anomalyDetectionList = d.getSqlDBs(server)

	// Create SQL database service voc object for SQL server
	dbService = &ontology.RelationalDatabaseService{
		Id:           resourceID(server.ID),
		Name:         util.Deref(server.Name),
		CreationTime: nil,
		GeoLocation:  location(server.Location),
		Labels:       labels(server.Tags),
		ParentId:     resourceGroupID(server.ID),
		Raw:          discovery.Raw(server),
		// TODO(all): HttpEndpoint
		TransportEncryption: &ontology.TransportEncryption{
			Enabled:         true,
			Enforced:        true,
			Protocol:        constants.TLS,
			ProtocolVersion: tlsVersion((*string)(server.Properties.MinimalTLSVersion)),
		},
		AnomalyDetections: anomalyDetectionList,
	}

	// Add SQL database service
	list = append(list, dbService)

	// Add SQL database storages
	list = append(list, dbList...)

	return list, nil
}

func (d *azureDiscovery) handleStorageAccount(account *armstorage.Account, storagesList []ontology.IsResource, activityLogging *ontology.ActivityLogging, rawActivityLogging string) (*ontology.ObjectStorageService, error) {
	var (
		storageResourceIDs []string
	)

	if account == nil {
		return nil, ErrEmptyStorageAccount
	}

	// Get all object storage IDs
	for _, storage := range storagesList {
		if strings.Contains(string(storage.GetId()), accountName(util.Deref(account.ID))) {
			storageResourceIDs = append(storageResourceIDs, storage.GetId())
		}
	}

	te := &ontology.TransportEncryption{
		Enforced:        util.Deref(account.Properties.EnableHTTPSTrafficOnly),
		Enabled:         true, // cannot be disabled
		Protocol:        constants.TLS,
		ProtocolVersion: tlsVersion((*string)(account.Properties.MinimumTLSVersion)),
	}

	storageService := &ontology.ObjectStorageService{
		Id:                  resourceID(account.ID),
		Name:                util.Deref(account.Name),
		StorageIds:          storageResourceIDs,
		CreationTime:        creationTime(account.Properties.CreationTime),
		GeoLocation:         location(account.Location),
		Labels:              labels(account.Tags),
		ParentId:            resourceGroupID(account.ID),
		Raw:                 discovery.Raw(account, rawActivityLogging),
		TransportEncryption: te,
		HttpEndpoint: &ontology.HttpEndpoint{
			Url:                 generalizeURL(util.Deref(account.Properties.PrimaryEndpoints.Blob)),
			TransportEncryption: te,
		},
		ActivityLogging: activityLogging,
	}

	return storageService, nil
}

func (d *azureDiscovery) handleFileStorage(account *armstorage.Account, fileshare *armstorage.FileShareItem, activityLogging *ontology.ActivityLogging, rawActivityLogging string) (*ontology.FileStorage, error) {
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

	return &ontology.FileStorage{
		Id:           resourceID(fileshare.ID),
		Name:         util.Deref(fileshare.Name),
		CreationTime: creationTime(account.Properties.CreationTime), // We only have the creation time of the storage account the file storage belongs to
		GeoLocation:  location(account.Location),                    // The location is the same as the storage account
		Labels:       labels(account.Tags),                          // The storage account labels the file storage belongs to
		ParentId:     resourceIDPointer(account.ID),                 // the storage account is our parent
		Raw:          discovery.Raw(account, fileshare),
		ResourceLogging: &ontology.ResourceLogging{
			MonitoringLogDataEnabled: monitoringLogDataEnabled,
			SecurityAlertsEnabled:    securityAlertsEnabled,
		},
		ActivityLogging:  activityLogging,
		AtRestEncryption: enc,
	}, nil
}

func (d *azureDiscovery) handleObjectStorage(account *armstorage.Account, container *armstorage.ListContainerItem, activityLogging *ontology.ActivityLogging) (*ontology.ObjectStorage, error) {
	var (
		backups                  []*ontology.Backup
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

	return &ontology.ObjectStorage{
		Id:               resourceID(container.ID),
		Name:             util.Deref(container.Name),
		CreationTime:     creationTime(account.Properties.CreationTime), // We only have the creation time of the storage account the file storage belongs to
		GeoLocation:      location(account.Location),                    // The location is the same as the storage account
		Labels:           labels(account.Tags),                          // The storage account labels the file storage belongs to
		ParentId:         resourceIDPointer(account.ID),                 // the storage account is our parent
		Raw:              discovery.Raw(account, container),
		AtRestEncryption: enc,
		Immutability: &ontology.Immutability{
			Enabled: util.Deref(container.Properties.HasImmutabilityPolicy),
		},
		ResourceLogging: &ontology.ResourceLogging{
			MonitoringLogDataEnabled: monitoringLogDataEnabled,
			SecurityAlertsEnabled:    securityAlertsEnabled,
		},
		ActivityLogging: activityLogging,
		Backups:         backups,
		PublicAccess:    util.Deref(container.Properties.PublicAccess) != armstorage.PublicAccessNone,
	}, nil
}
