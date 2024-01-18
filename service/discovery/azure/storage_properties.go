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

	"clouditor.io/clouditor/internal/constants"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

// checkTlsVersion returns Clouditor's TLS version constants for the given TLS version
func checkTlsVersion(version string) string {
	// Check TLS version
	switch version {
	case constants.TLS1_0, constants.TLS1_1, constants.TLS1_2:
		return version
	case "1.0", "1_0":
		return constants.TLS1_0
	case "1.1", "1_1":
		return constants.TLS1_1
	case "1.2", "1_2":
		return constants.TLS1_2
	default:
		log.Warningf("'%s' is no implemented TLS version.", version)
		return ""
	}
}

// storageAtRestEncryption takes encryption properties of an armstorage.Account and converts it into our respective ontology object.
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
