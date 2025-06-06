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

package openstack

import (
	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/util"
	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/clusters"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// handleCluster creates a container resource based on the Clouditor Ontology
func (d *openstackDiscovery) handleCluster(cluster *clusters.Cluster) (ontology.IsResource, error) {
	r := &ontology.ContainerOrchestration{
		Id:           cluster.UUID,
		Name:         cluster.Name,
		CreationTime: timestamppb.New(cluster.CreatedAt),
		GeoLocation: &ontology.GeoLocation{
			Region: d.region,
		},
		Labels:   cluster.Labels,
		ParentId: util.Ref(cluster.ProjectID),
		Raw:      discovery.Raw(cluster),
	}

	log.Infof("Adding cluster '%s", cluster.Name)

	return r, nil
}
