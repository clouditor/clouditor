// Copyright 2023 Fraunhofer AISEC
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

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/security/armsecurity"
)

var (
	ErrEmptyPricing = errors.New("pricing is empty")
)

type azureDefenderDiscovery struct {
	*azureDiscovery
}

func NewAzureDefenderDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &azureStorageDiscovery{
		&azureDiscovery{
			discovererComponent: StorageComponent,
			csID:                discovery.DefaultCloudServiceID,
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(d.azureDiscovery)
	}

	return d
}

func (*azureDefenderDiscovery) Name() string {
	return "Azure Storage Account"
}

func (*azureDefenderDiscovery) Description() string {
	return "Discovery Azure Defender for Clouds."
}

func (d *azureDefenderDiscovery) List() (list []voc.IsCloudResource, err error) {
	// make sure, we are authorized
	if err = d.authorize(); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrCouldNotAuthenticate, err)
	}

	log.Info("Discover Azure storage resources")

	// Discover defender
	defender, err := d.discoverDefender()
	if err != nil {
		return nil, fmt.Errorf("could not discover defender: %w", err)
	}
	list = append(list, defender...)

	return
}

func (d *azureDefenderDiscovery) discoverDefender() ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	// initialize defender client: to get the properties set for defender for cloud, we have to get the pricing information
	if err := d.initDefenderClient(); err != nil {
		return nil, err
	}

	// TODO(anatheka): List all storage accounts across all resource groups
	pricingsList, err := d.clients.defenderClient.List(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not discover pricing")
	}

	for _, pricing := range pricingsList.Value {
		defender, err := d.handleDefender(pricing)
		if err != nil {
			return nil, fmt.Errorf("could not handle pricing: %w", err)
		}
		log.Infof("Adding defender '%s'", defender.Name)

		list = append(list, defender)
	}

	return list, nil
}

func (d *azureDefenderDiscovery) handleDefender(pricing *armsecurity.Pricing) (*voc.Account, error) {
	var inventoryOfAssetsEnabled bool

	if pricing == nil {
		return nil, ErrEmptyPricing
	}

	if *pricing.Properties.PricingTier == armsecurity.PricingTierFree {
		inventoryOfAssetsEnabled = true
	}

	return &voc.Account{
		Resource: discovery.NewResource(d,
			voc.ResourceID(*pricing.ID),
			*pricing.Name,
			nil,
			voc.GeoLocation{},
			nil,
			voc.AccountType,
		),
		InventoryOfAssetsEnabled: inventoryOfAssetsEnabled,
	}, nil
}

// initDefenderClient creates the client if not already exists
func (d *azureDefenderDiscovery) initDefenderClient() (err error) {
	d.clients.defenderClient, err = initClient(d.clients.defenderClient, d.azureDiscovery, armsecurity.NewPricingsClient)
	return
}
