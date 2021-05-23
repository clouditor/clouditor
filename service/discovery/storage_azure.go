// Copyright 2021 Fraunhofer AISEC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package discovery

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/resources/mgmt/subscriptions"
	azure_storage "github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-02-01/storage"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

type azureStorageDiscovery struct{}

func (d *azureStorageDiscovery) List() (list []Storage, err error) {
	// create an authorizer from env vars or Azure Managed Service Idenity
	authorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		log.Errorf("Could not authenticate to Azure: %s", err)
		return
	}

	subClient := subscriptions.NewClient()
	subClient.Authorizer = authorizer

	// get first subription
	page, _ := subClient.List(context.Background())
	sub := page.Values()[0]

	log.Infof("Using %s as subscription", *sub.SubscriptionID)

	client := azure_storage.NewAccountsClient(*sub.SubscriptionID)
	client.Authorizer = authorizer

	ctx := context.Background()

	result, _ := client.List(ctx)

	for _, v := range result.Values() {
		endpoint := &HttpEndpoint{&TransportEncryption{
			Enforced:   *v.EnableHTTPSTrafficOnly,
			Enabled:    true, // cannot be disabled
			TlsVersion: string(v.MinimumTLSVersion),
		}}

		s := &objectStorage{storage{resource{
			id:           *v.ID,
			name:         *v.Name,
			creationTime: &v.CreationTime.Time,
		}, &AtRestEncryption{}}, endpoint}

		log.Infof("Adding storage account %+v", s)

		list = append(list, s)
	}

	return
}
