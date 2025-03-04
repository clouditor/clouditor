// Copyright 2020-2024 Fraunhofer AISEC
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
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning"
)

func getComputeStringList(values []ontology.IsResource) []string {
	var list []string

	for _, value := range values {
		list = append(list, value.GetId())

	}

	return list
}

func getInternetAccessibleEndpoint(status *armmachinelearning.PublicNetworkAccess) bool {
	// Check if status is empty
	if status == nil {
		return false
	}

	return util.Deref(status) == armmachinelearning.PublicNetworkAccessEnabled
}

// getResourceLogging returns true if application insights contains a string (applicationInsights enabled), otherwise it returns false
func getResourceLogging(log *string) *ontology.ResourceLogging {
	// Check if logging service storage is available
	if util.Deref(log) == "" {
		return &ontology.ResourceLogging{
			Enabled: false,
		}
	}

	return &ontology.ResourceLogging{
		Enabled:           true,
		LoggingServiceIds: []string{resourceID(log)},
	}
}

func getAtRestEncryption(enc *armmachinelearning.EncryptionProperty) (atRestEnc *ontology.AtRestEncryption) {

	// If the encryption property is nil, the ML workspace has managed key encryption in use
	if enc == nil {
		return &ontology.AtRestEncryption{
			Type: &ontology.AtRestEncryption_ManagedKeyEncryption{
				ManagedKeyEncryption: &ontology.ManagedKeyEncryption{
					Enabled:   true,
					Algorithm: AES256,
				},
			},
		}
	}

	if util.Deref(enc.KeyVaultProperties.KeyVaultArmID) == "" {
		atRestEnc = &ontology.AtRestEncryption{
			Type: &ontology.AtRestEncryption_ManagedKeyEncryption{
				ManagedKeyEncryption: &ontology.ManagedKeyEncryption{
					Enabled:   getEncryptionStatus(enc.Status),
					Algorithm: AES256,
				},
			},
		}
	} else {
		atRestEnc = &ontology.AtRestEncryption{
			Type: &ontology.AtRestEncryption_CustomerKeyEncryption{
				CustomerKeyEncryption: &ontology.CustomerKeyEncryption{
					Enabled: getEncryptionStatus(enc.Status),
					KeyUrl:  resourceID(enc.KeyVaultProperties.KeyVaultArmID),
				},
			},
		}
	}

	return atRestEnc

}

func getEncryptionStatus(enc *armmachinelearning.EncryptionStatus) bool {
	return util.Deref(enc) == armmachinelearning.EncryptionStatusEnabled
}
