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
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
)

// handleResourceGroup returns a [voc.Account] out of an existing [armsubscription.Subscription].
func (d *azureDiscovery) handleResourceGroup(rg *armresources.ResourceGroup) voc.IsCloudResource {
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
func (d *azureDiscovery) handleSubscription(s *armsubscription.Subscription) *voc.Account {
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
