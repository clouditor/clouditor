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
	"fmt"

	"context"
	"strings"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/constants"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumetypes"
	"github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/containers"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// handleBlockStorage creates a block storage resource based on the Clouditor Ontology
func (d *openstackDiscovery) handleBlockStorage(volume *volumes.Volume) (ontology.IsResource, error) {
	var are *ontology.AtRestEncryption

	// Get Name, if exits, otherwise take the ID
	name := volume.Name
	if volume.Name == "" {
		name = volume.ID
	}

	// Get encryption information
	vType, err := volumetypes.Get(context.Background(), d.clients.computeClient, volume.ID).Extract()
	if err != nil {
		log.Errorf("Error getting volume type information for volume '%s': %s", volume.Name, err)
	} else {
		enc, err := volumetypes.GetEncryption(context.Background(), d.clients.computeClient, vType.ID).Extract()
		if err != nil {
			log.Errorf("Error getting encryption information for volume '%s': %s", volume.Name, err)
		}
		if enc.EncryptionID != "" {
			are = &ontology.AtRestEncryption{
				Type: &ontology.AtRestEncryption_CustomerKeyEncryption{
					CustomerKeyEncryption: &ontology.CustomerKeyEncryption{
						Enabled:   true,
						Algorithm: enc.Cipher,
					},
				},
			}
		}

	}

	r := &ontology.BlockStorage{
		Id:           volume.ID,
		Name:         name,
		Description:  volume.Description,
		CreationTime: timestamppb.New(volume.CreatedAt),
		GeoLocation: &ontology.GeoLocation{
			Region: d.region,
		},
		ParentId:         util.Ref(getParentID(volume)),
		Labels:           map[string]string{}, // Not available
		Raw:              discovery.Raw(volume),
		AtRestEncryption: are,
	}
	log.Infof("Adding block storage '%s'", volume.Name)

	// Create project resource for the parentId if not available
	err = d.addProjectIfMissing(volume.TenantID, volume.TenantID, d.domain.domainID)
	if err != nil {
		return nil, fmt.Errorf("could not handle project for block storage %s: %w", volume.ID, err)
	}

	log.Infof("Adding block storage '%s", r.Name)
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
		Description:  "",                                         // Not available
		ParentId:     util.Ref(d.clients.storageClient.Endpoint), // parent is the object storage service stored in the endpoint
		PublicAccess: isPublic,
	}

	return r, nil
}

// handleObjectStorageService creates an object storage service resource based on the Clouditor Ontology
func (d *openstackDiscovery) handleObjectStorageService() (ontology.IsResource, error) {
	var (
		te *ontology.TransportEncryption
	)
	if strings.HasPrefix(d.clients.storageClient.Endpoint, "https://") {
		te = &ontology.TransportEncryption{
			Enabled:  true,
			Enforced: true,
			Protocol: constants.TLS,
			// ProtocolVersion: 1.2, // information not available
		}
	}

	r := &ontology.ObjectStorageService{
		Id:   d.clients.storageClient.Endpoint,
		Name: "Swift Object Storage Service",
		GeoLocation: &ontology.GeoLocation{
			Region: d.region,
		},
		ParentId: util.Ref(d.configuredProject.projectID),
		HttpEndpoint: &ontology.HttpEndpoint{
			Url:                 d.clients.storageClient.Endpoint,
			TransportEncryption: te,
		},
	}

	return r, nil
}
