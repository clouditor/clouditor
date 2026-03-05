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
	"context"
	"strings"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/containers"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// handleBlockStorage creates a block storage resource based on the Clouditor Ontology
func (d *openstackDiscovery) handleBlockStorage(volume *volumes.Volume) (ontology.IsResource, error) {
	// Get Name, if exits, otherwise take the ID
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
			Region: d.region,
		},
		ParentId: util.Ref(getParentID(volume)),
		Labels:   map[string]string{}, // Not available
		Raw:      discovery.Raw(volume),
		AtRestEncryption: &ontology.AtRestEncryption{
			Type: &ontology.AtRestEncryption_CustomerKeyEncryption{
				CustomerKeyEncryption: &ontology.CustomerKeyEncryption{
					Enabled: volume.Encrypted,
					// Algorithm: , // Not available
				},
			},
		},
	}
	log.Infof("Adding block storage '%s'", volume.Name)

	return r, nil
}

// handleObjectStorage creates an object storage resource based on the Clouditor Ontology
func (d *openstackDiscovery) handleObjectStorage(container *containers.Container) (ontology.IsResource, error) {
	var (
		isPublic bool
	)
	res := containers.Get(context.Background(), d.clients.storageClient, container.Name, containers.GetOpts{})

	header, err := res.Extract()
	if err != nil {
		log.Errorf("Error extracting container details for container '%s': %s", container.Name, err)
	} else {
		// Check if the container is public by looking for the "X-Container-Read" header and checking if it contains ".r:*"
		for _, acl := range header.Read {
			if strings.Contains(acl, ".r:*") {
				isPublic = true
				break
			}
		}
	}

	r := &ontology.ObjectStorage{
		Id:           container.Name,
		Name:         container.Name,
		Description:  "", // Not available
		PublicAccess: isPublic,
	}

	return r, nil
}
