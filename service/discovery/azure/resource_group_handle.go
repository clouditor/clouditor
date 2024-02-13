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
	"clouditor.io/clouditor/api/ontology"
	"clouditor.io/clouditor/internal/util"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
)

// handleResourceGroup returns a [voc.Account] out of an existing [armsubscription.Subscription].
func (d *azureDiscovery) handleResourceGroup(rg *armresources.ResourceGroup) ontology.IsResource {
	return &ontology.ResourceGroup{
		Id:   util.Deref(rg.ID),
		Name: util.Deref(rg.Name),
		GeoLocation: &ontology.GeoLocation{
			Region: util.Deref(rg.Location),
		},
		Labels:   labels(rg.Tags),
		ParentId: util.Deref(d.sub.ID),
	}
}

// handleSubscription returns a [voc.Account] out of an existing [armsubscription.Subscription].
func (d *azureDiscovery) handleSubscription(s *armsubscription.Subscription) *ontology.Account {
	return &ontology.Account{
		Id:           util.Deref(s.ID),
		Name:         util.Deref(s.DisplayName),
		CreationTime: nil, // subscriptions do not have a creation date
		GeoLocation:  nil, // subscriptions are global
		Labels:       nil, // subscriptions do not have labels,
		// TODO(oxisto): we should make this an optional string
		ParentId: "", // subscriptions are the top-most item and have no parent,
	}
}
