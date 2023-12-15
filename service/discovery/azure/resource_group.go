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
	"fmt"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
)

type azureResourceGroupDiscovery struct {
	*azureDiscovery
}

func NewAzureResourceGroupDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &azureResourceGroupDiscovery{
		&azureDiscovery{
			discovererComponent: ComputeComponent,
			csID:                discovery.DefaultCloudServiceID,
			backupMap:           make(map[string]*backup),
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(d.azureDiscovery)
	}

	return d
}

func (*azureResourceGroupDiscovery) Name() string {
	return "Azure Resource Group"
}

func (*azureResourceGroupDiscovery) Description() string {
	return "Discovery of Azure resource groups."
}

// List resource groups and cloud account
func (d *azureResourceGroupDiscovery) List() (list []voc.IsCloudResource, err error) {
	// make sure we are authorized
	if err = d.authorize(); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrCouldNotAuthenticate, err)
	}

	// initialize client
	if err := d.initResourceGroupsClient(); err != nil {
		return nil, err
	}

	// Build an account as the most top-level item. Our subscription will serve as the account
	acc := d.handleSubscription(d.sub)

	list = append(list, acc)

	listPager := d.clients.rgClient.NewListPager(&armresources.ResourceGroupsClientListOptions{})
	for listPager.More() {
		page, err := listPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("%s: %v", ErrGettingNextPage, err)
			return nil, err
		}

		for _, rg := range page.Value {
			// If we are scoped to one resource group, we can skip the rest of the groups
			if d.azureDiscovery.rg != nil && util.Deref(rg.Name) != util.Deref(d.azureDiscovery.rg) {
				continue
			}

			r := d.handleResourceGroup(rg)

			log.Infof("Adding resource group '%s'", r.GetName())

			list = append(list, r)
		}
	}

	return
}

// handleResourceGroup returns a [voc.Account] out of an existing [armsubscription.Subscription].
func (d *azureResourceGroupDiscovery) handleResourceGroup(rg *armresources.ResourceGroup) voc.IsCloudResource {
	return &voc.ResourceGroup{
		Resource: discovery.NewResource(
			d,
			voc.ResourceID(util.Deref(rg.ID)),
			util.Deref(rg.Name),
			nil,
			voc.GeoLocation{
				Region: util.Deref(rg.Location),
			},
			labels(rg.Tags),
			voc.ResourceID(*d.sub.ID),
			voc.ResourceGroupType,
		),
	}
}

// handleSubscription returns a [voc.Account] out of an existing [armsubscription.Subscription].
func (d *azureResourceGroupDiscovery) handleSubscription(s *armsubscription.Subscription) *voc.Account {
	return &voc.Account{
		Resource: discovery.NewResource(
			d,
			voc.ResourceID(util.Deref(s.ID)),
			util.Deref(s.DisplayName),
			// subscriptions do not have a creation date
			nil,
			// subscriptions are global
			voc.GeoLocation{},
			// subscriptions do not have labels,
			nil,
			// subscriptions are the top-most item and have no parent,
			"",
			voc.AccountType,
		),
	}
}

// azureResourceGroupDiscovery creates the client if not already exists
func (d *azureResourceGroupDiscovery) initResourceGroupsClient() (err error) {
	d.clients.rgClient, err = initClient(d.clients.rgClient, d.azureDiscovery, armresources.NewResourceGroupsClient)
	return
}
