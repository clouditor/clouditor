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
	"context"
	"errors"
	"fmt"

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/constants"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

// storageAtRestEncryption takes encryption properties of an armstorage.Account and converts it into our respective ontology object.
func storageAtRestEncryption(account *armstorage.Account) (enc *ontology.AtRestEncryption, err error) {
	if account == nil {
		return enc, ErrEmptyStorageAccount
	}

	if account.Properties == nil || account.Properties.Encryption.KeySource == nil {
		return enc, errors.New("keySource is empty")
	} else if util.Deref(account.Properties.Encryption.KeySource) == armstorage.KeySourceMicrosoftStorage {
		enc = &ontology.AtRestEncryption{
			Type: &ontology.AtRestEncryption_ManagedKeyEncryption{
				ManagedKeyEncryption: &ontology.ManagedKeyEncryption{
					Algorithm: constants.AES256,
					Enabled:   true,
				},
			},
		}
	} else if util.Deref(account.Properties.Encryption.KeySource) == armstorage.KeySourceMicrosoftKeyvault {
		enc = &ontology.AtRestEncryption{
			Type: &ontology.AtRestEncryption_CustomerKeyEncryption{
				CustomerKeyEncryption: &ontology.CustomerKeyEncryption{
					Algorithm: "", // TODO(all): TBD
					Enabled:   true,
					// TODO(oxisto): This should also include the key!
					KeyUrl: util.Deref(account.Properties.Encryption.KeyVaultProperties.KeyVaultURI),
				},
			},
		}
	}

	return enc, nil
}

// anomalyDetectionEnabled returns true if Azure Advanced Threat Protection is enabled for the database.
func (d *azureDiscovery) anomalyDetectionEnabled(server *armsql.Server, db *armsql.Database) (bool, error) {
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

// getActivityLogging returns the activity logging information for the storage account, blob, table and file storage including their raw information
func (d *azureDiscovery) getActivityLogging(account *armstorage.Account) (activityLoggingAccount, activityLoggingBlob, activityLoggingTable, activityLoggingFile *ontology.ActivityLogging, rawAccount, rawBlob, rawTable, rawFile string) {

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
