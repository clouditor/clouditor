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
	"context"
	"fmt"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-02-01/storage"
	"github.com/Azure/go-autorest/autorest/to"
)

type azureStorageDiscovery struct {
}

func NewAzureStorageDiscovery() discovery.Discoverer {
	return &azureStorageDiscovery{}
}

func (d *azureStorageDiscovery) Name() string {
	return "Azure Storage Account"
}

func (d *azureStorageDiscovery) Description() string {
	return "Discovery Azure storage accounts."
}

var StorageListFunc ListFunc = storage.AccountsClient.List
var StorageClientFunc ClientFunc = createClient

type ListFunc func(storage.AccountsClient, context.Context) (storage.AccountListResultPage, error)
type ClientFunc func() (client storage.AccountsClient, err error)

func (d *azureStorageDiscovery) List() (list []voc.IsResource, err error) {
	client, _ := StorageClientFunc()

	result, _ := StorageListFunc(client, azureAuthorizer.ctx)

	for _, v := range result.Values() {
		s := handleStorageAccount(v)

		log.Infof("Adding storage account %+v", s)

		list = append(list, s)
	}

	return
}

func createClient() (client storage.AccountsClient, err error) {
	if err = azureAuthorizer.Authorize(); err != nil {
		return storage.AccountsClient{}, fmt.Errorf("could not authorize Azure account: %w", err)
	}

	client = storage.NewAccountsClient(*azureAuthorizer.sub.SubscriptionID)
	client.Authorizer = azureAuthorizer.authorizer

	return
}

func handleStorageAccount(account storage.Account) voc.IsStorage {
	return &voc.ObjectStorageResource{StorageResource: voc.StorageResource{
		Resource: voc.Resource{
			ID:           to.String(account.ID),
			Name:         to.String(account.Name),
			CreationTime: account.CreationTime.Unix(),
		},
		AtRestEncryption: voc.NewAtRestEncryption(
			to.Bool(account.Encryption.Services.Blob.Enabled),
			"AES-265", // seems to be always AES-256
			string(account.Encryption.KeySource),
		)},
		HttpEndpoint: &voc.HttpEndpoint{
			URL: to.String(account.PrimaryEndpoints.Blob),
			TransportEncryption: voc.NewTransportEncryption(
				true, // cannot be disabled
				to.Bool(account.EnableHTTPSTrafficOnly),
				string(account.MinimumTLSVersion),
			),
		},
	}
}
