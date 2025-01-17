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

package openstack

import (
	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// handleNetworkInterfaces creates a network interface resource based on the Clouditor Ontology
func (d *openstackDiscovery) handleNetworkInterfaces(network *networks.Network) (ontology.IsResource, error) {
	var parentId string

	// Check if a network interface is assigned to the discovered project. If it is not, do not add the parent ID, allowing the network interface to remain unconnected to the project while still being associated with a virtual machine resource and displayed in the UI graph.
	if d.projectID == network.TenantID {
		parentId = network.TenantID
	}

	r := &ontology.NetworkInterface{
		Id:           network.ID,
		Name:         network.Name,
		Description:  network.Description,
		CreationTime: timestamppb.New(network.CreatedAt),
		GeoLocation: &ontology.GeoLocation{
			Region: d.region,
		},
		Labels:   labels(util.Ref(network.Tags)),
		ParentId: util.Ref(parentId),
		Raw:      discovery.Raw(network),
	}

	log.Infof("Adding network interface '%s", network.Name)

	return r, nil
}
