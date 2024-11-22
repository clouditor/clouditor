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

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/volumes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// handleBlockStorage creates a block storage resource based on the Clouditor Ontology
func (d *openstackDiscovery) handleBlockStorage(volume *volumes.Volume) (ontology.IsResource, error) {
	name := volume.Name

	if volume.Name == "" {
		name = volume.ID
	}

	r := &ontology.BlockStorage{
		Id:           volume.ID,
		Name:         name,
		Description:  volume.Description,
		CreationTime: timestamppb.New(volume.CreatedAt),
		GeoLocation: &ontology.GeoLocation{
			Region: "unknown", // TODO: Can we get the region?
		},
		Labels: map[string]string{}, // Not available
		// ParentId: util.Ref(volume.TenantID),
		Raw: discovery.Raw(volume),
	}

	return r, nil
}
