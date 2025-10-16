// Copyright 2025 Fraunhofer AISEC
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

package ionos

import (
	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/util"
	"google.golang.org/protobuf/types/known/timestamppb"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

// handleNetworkInterfaces creates a network interface resource based on the Clouditor Ontology
func (d *ionosDiscovery) handleNetworkInterfaces(nic ionoscloud.Nic, dc ionoscloud.Datacenter) (ontology.IsResource, error) {
	r := &ontology.NetworkInterface{
		Id:           util.Deref(nic.Id),
		Name:         util.Deref(nic.Properties.Name),
		CreationTime: timestamppb.New(util.Deref(nic.Metadata.GetCreatedDate())),
		GeoLocation: &ontology.GeoLocation{
			Region: util.Deref(dc.Properties.GetLocation()),
		},
		// Labels:   , // Not available
		ParentId: dc.GetId(),
		Raw:      discovery.Raw(nic, dc),
		AccessRestriction: &ontology.AccessRestriction{
			Type: &ontology.AccessRestriction_L3Firewall{
				L3Firewall: &ontology.L3Firewall{
					Enabled:         util.Deref(nic.Properties.GetFirewallActive()),
					RestrictedPorts: "",
				},
			},
		},
	}

	return r, nil
}
