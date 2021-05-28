/*
 * Copyright 2021 Fraunhofer AISEC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *           $$\                           $$\ $$\   $$\
 *           $$ |                          $$ |\__|  $$ |
 *  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
 * $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
 * $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
 * $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
 * \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
 *  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
 *
 * This file is part of Clouditor Community Edition.
 */

package azure

import (
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-02-01/storage"
)

type azureStorageDiscovery struct{}

func NewAzureStorageDiscovery() discovery.Discoverer {
	return &azureStorageDiscovery{}
}

func (d *azureStorageDiscovery) Name() string {
	return "Azure Storage Account"
}

func (d *azureStorageDiscovery) Description() string {
	return "Discovery Azure storage accounts."
}

func (d *azureStorageDiscovery) List() (list []voc.IsResource, err error) {

	client := storage.NewAccountsClient(*azureAuthorizer.sub.SubscriptionID)
	client.Authorizer = azureAuthorizer.authorizer

	result, _ := client.List(azureAuthorizer.ctx)

	for _, v := range result.Values() {
		s := handleStorageAccount(v)

		log.Infof("Adding storage account %+v", s)

		list = append(list, s)
	}

	return
}

func handleStorageAccount(account storage.Account) voc.IsStorage {
	return &voc.ObjectStorageResource{StorageResource: voc.StorageResource{
		Resource: voc.Resource{
			ID:           *account.ID,
			Name:         *account.Name,
			CreationTime: account.CreationTime.Unix(),
		},
		AtRestEncryption: voc.NewAtRestEncryption(
			*account.Encryption.Services.Blob.Enabled,
			"AES-265", // seems to be always AES-256
			string(account.Encryption.KeySource),
		)},
		HttpEndpoint: &voc.HttpEndpoint{
			URL: *account.PrimaryEndpoints.Blob,
			TransportEncryption: voc.NewTransportEncryption(
				true, // cannot be disabled
				*account.EnableHTTPSTrafficOnly,
				string(account.MinimumTLSVersion),
			),
		},
	}
}
