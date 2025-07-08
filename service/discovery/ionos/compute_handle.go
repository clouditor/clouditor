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

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/util"
	"github.com/ionos-cloud/sdk-go-bundle/products/compute"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// handleServer creates a virtual machine resource based on the Clouditor Ontology
func (d *ionosDiscovery) handleServer(server compute.Server, dc compute.Datacenter) (ontology.IsResource, error) {
	// Getting labels
	l, _, err := d.clients.computeClient.LabelsApi.
		DatacentersServersLabelsGet(context.Background(), util.Deref(dc.GetId()), util.Deref(server.GetId())).
		Execute()
	if err != nil {
		log.Errorf("error getting labels for server %s: %s", util.Deref(server.Id), err)
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
