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
	"strings"

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
	d := &azureDefenderDiscovery{
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
	return "Azure Defender for Cloud"
}

func (*azureDefenderDiscovery) Description() string {
	return "Discovery Azure Defender for Clouds."
}

func (d *azureDefenderDiscovery) List() (list []voc.IsCloudResource, err error) {
	// make sure, we are authorized
	if err = d.authorize(); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrCouldNotAuthenticate, err)
	}

	log.Info("Discover Azure Defender for Cloud")

	// Discover defender
	defender, err := d.discoverDefender()
	if err != nil {
		return nil, fmt.Errorf("could not discover defender: %w", err)
	}
	list = append(list, defender...)

	return
}

// discoverDefender discovers the enabled status for the Denfender for X services.
func (d *azureDefenderDiscovery) discoverDefender() ([]voc.IsCloudResource, error) {
	var (
		list                     []voc.IsCloudResource
		monitoringLogDataEnabled = true
	)

	// initialize defender client: to get the properties set for defender for cloud, we have to get the pricing information
	if err := d.initDefenderClient(); err != nil {
		return nil, err
	}

	// List all pricings to get the enabled Defender for X
	pricingsList, err := d.clients.defenderClient.List(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not discover pricing")
	}

	for _, pricing := range pricingsList.Value {
		defender, err := d.handleDefender(pricing)
		if err != nil {
			return nil, fmt.Errorf("could not handle pricing: %w", err)
		}

		log.Infof("Adding defender '%s'", pricing.Name)
		list = append(list, defender)
	}

	// Get information if security monitoring is enabled for the subscription. The security monitoring is enabled for the subscription if all defender services are enabled.
	for _, defender := range list {
		defenderTest := defender.(*voc.Account)
		if !defenderTest.InventoryOfAssets {
			monitoringLogDataEnabled = false
			break
		}
	}

	list = append(list, &voc.Account{
		Resource: discovery.NewResource(d,
			voc.ResourceID(getSubscriptionID(string(list[0].GetID()))),
			getSubscriptionName(string(list[0].GetID())),
			nil,
			voc.GeoLocation{},
			nil,
			voc.AccountType,
		),
		MonitoringLogData: monitoringLogDataEnabled,
	})

	return list, nil
}

// TODO(anatheka): Which is the right Ontology resource for the defender?
func (d *azureDefenderDiscovery) handleDefender(pricing *armsecurity.Pricing) (*voc.Account, error) {
	var (
		inventoryOfAssetsEnabled bool
	)

	if pricing == nil {
		return nil, ErrEmptyPricing
	}

	// TODO(all): Maybe we have to check here which pricing tier is used and based on that the specific features for, e.g., Storage, IoT, ARM are used. For now, we add for each defener one voc.Account object.
	if *pricing.Properties.PricingTier == armsecurity.PricingTierFree {
		inventoryOfAssetsEnabled = false
	}

	// TODO(all): What should we use? Account or another Resource?
	return &voc.Account{
		Resource: discovery.NewResource(d,
			voc.ResourceID(*pricing.ID),
			*pricing.Name,
			nil,
			voc.GeoLocation{},
			nil,
			voc.AccountType,
		),
		InventoryOfAssets: inventoryOfAssetsEnabled,
	}, nil

}

// initDefenderClient creates the client if not already exists
func (d *azureDefenderDiscovery) initDefenderClient() (err error) {
	d.clients.defenderClient, err = initClient(d.clients.defenderClient, d.azureDiscovery, armsecurity.NewPricingsClient)
	return
}

// getSubscriptionName returns the Azure subscription number
func getSubscriptionName(id string) string {
	return strings.Split(id, "/")[2]
}

// getSubscriptionID returns the Azure subscription ID
func getSubscriptionID(id string) string {
	split := strings.Split(id, "/")

	return "/" + split[1] + "/" + split[2]
}
