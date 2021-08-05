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
	"fmt"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-02-01/storage"
	"github.com/Azure/go-autorest/autorest/to"
)

type azureStorageDiscovery struct {
	azureDiscovery
}

func NewAzureStorageDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &azureStorageDiscovery{}

	for _, opt := range opts {
		if auth, ok := opt.(*authorizerOption); ok {
			d.authOption = auth
		} else {
			d.options = append(d.options, opt)
		}
	}

	return d
}

func (d *azureStorageDiscovery) Name() string {
	return "Azure Storage Account"
}

func (d *azureStorageDiscovery) Description() string {
	return "Discovery Azure storage accounts."
}

func (d *azureStorageDiscovery) List() (list []voc.IsCloudResource, err error) {
	// make sure, we are authorized
	if err = d.authorize(); err != nil {
		return nil, fmt.Errorf("could not authorize Azure account: %w", err)
	}

	// Discover storage accounts
	storageAccounts, err := d.discoverStorageAccounts()
	if err != nil {
		return nil, fmt.Errorf("could not discover storage accounts: %w", err)
	}
	list = append(list, storageAccounts...)

	return
}

func (d *azureStorageDiscovery) discoverStorageAccounts() ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	client := storage.NewAccountsClient(to.String(d.sub.SubscriptionID))
	d.apply(&client.Client)

	result, err := client.List(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not list storage accounts: %w", err)
	}

	accounts := result.Values()
	for i := range accounts {
		s := handleStorageAccount(&accounts[i])

		log.Infof("Adding storage account %+v", s)

		list = append(list, s)
	}

	return list, err
}

func handleStorageAccount(account *storage.Account) voc.IsStorage {

	return &voc.ObjectStorage{
		Storage: &voc.Storage{
			CloudResource: &voc.CloudResource{
				ID:           voc.ResourceID(to.String(account.ID)),
				Name:         to.String(account.Name),
				CreationTime: account.CreationTime.Unix(),
				Type:         []string{"ObjectStorage", "Storage", "Resource"},
			},
			AtRestEncryption: &voc.AtRestEncryption{
				KeyManager: string(account.Encryption.KeySource),
				Algorithm:  "AES-265", // seems to be always AES-256
				Enabled:    to.Bool(account.Encryption.Services.Blob.Enabled),
			},
		},
		HttpEndpoint: &voc.HttpEndpoint{
			Url: to.String(account.PrimaryEndpoints.Blob),
			TransportEncryption: &voc.TransportEncryption{
				Enforced:   to.Bool(account.EnableHTTPSTrafficOnly),
				Enabled:    true, // cannot be disabled
				TlsVersion: string(account.MinimumTLSVersion),
				Algorithm:  "", // TBD
			},
		},
	}
}
