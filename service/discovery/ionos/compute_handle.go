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
	"context"
	"fmt"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/util"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// handleServer creates a virtual machine resource based on the Clouditor Ontology
func (d *ionosDiscovery) handleServer(server ionoscloud.Server, dc ionoscloud.Datacenter) (ontology.IsResource, error) {
	// Getting labels
	l, _, err := d.clients.computeClient.LabelsApi.
		DatacentersServersLabelsGet(context.Background(), util.Deref(dc.GetId()), util.Deref(server.GetId())).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("error getting labels for server %s: %s", util.Deref(server.Id), err)
	}

	r := &ontology.VirtualMachine{
		Id:           util.Deref(server.Id),
		Name:         util.Deref(server.Properties.Name),
		CreationTime: timestamppb.New(*server.Metadata.GetCreatedDate()),
		GeoLocation: &ontology.GeoLocation{
			Region: util.Deref(dc.Properties.GetLocation()),
		},
		Labels:              labels(l),
		ParentId:            dc.GetId(),
		Raw:                 discovery.Raw(server, dc),
		BlockStorageIds:     getBlockStorageIds(server),
		NetworkInterfaceIds: getNetworkInterfaceIds(server),
		// BootLogging: , // Not available
		// OsLogging: , // Not available
		// MalwareProtection: , // Not available
		// AutomaticUpdates: , // Not available
		// Description: , // Not available
		ActivityLogging: &ontology.ActivityLogging{
			Enabled: true, // is always enabled
			// RetentionPeriod: , // TODO(all): TBD
			// LoggingServiceIds: , // Not available, the entries are only available in the requests service
		},
	}

	return r, nil
}

// handleBlockStorage creates a block storage resource based on the Clouditor Ontology
func (d *ionosDiscovery) handleBlockStorage(blockStorage ionoscloud.Volume, dc ionoscloud.Datacenter) (ontology.IsResource, error) {
	// Getting labels
	l, _, err := d.clients.computeClient.LabelsApi.
		DatacentersServersLabelsGet(context.Background(), util.Deref(dc.GetId()), util.Deref(blockStorage.GetId())).
		Execute()
	if err != nil {
		log.Errorf("error getting labels for server %s: %s", util.Deref(blockStorage.Id), err)
	}

	r := &ontology.BlockStorage{
		Id:           util.Deref(blockStorage.Id),
		Name:         util.Deref(blockStorage.Properties.Name),
		CreationTime: timestamppb.New(*blockStorage.Metadata.GetCreatedDate()),
		GeoLocation: &ontology.GeoLocation{
			Region: util.Deref(dc.Properties.Location),
		},
		Labels:   labels(l),
		ParentId: dc.GetId(),
		Raw:      discovery.Raw(blockStorage, dc),
		// AtRestEncryption: ,
		// Backups: ,

	}

	return r, nil
}

// handleLoadBalancer creates a load balancer resource based on the Clouditor Ontology
func (d *ionosDiscovery) handleLoadBalancer(loadBalancer ionoscloud.Loadbalancer, dc ionoscloud.Datacenter) (ontology.IsResource, error) {
	// Getting labels
	l, _, err := d.clients.computeClient.LabelsApi.
		DatacentersServersLabelsGet(context.Background(), util.Deref(dc.GetId()), util.Deref(loadBalancer.GetId())).
		Execute()
	if err != nil {
		log.Errorf("error getting labels for load balancer %s: %s", util.Deref(loadBalancer.Id), err)
	}
	r := &ontology.LoadBalancer{
		Id:           util.Deref(loadBalancer.Id),
		Name:         util.Deref(loadBalancer.Properties.Name),
		CreationTime: timestamppb.New(*loadBalancer.Metadata.GetCreatedDate()),
		GeoLocation: &ontology.GeoLocation{
			Region: util.Deref(dc.Properties.Location),
		},
		Labels:   labels(l),
		ParentId: dc.GetId(),
		Raw:      discovery.Raw(loadBalancer, dc),
		// Description: , // Not available
	}

	return r, nil
}
