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

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/util"
)

func (d *ionosDiscovery) discoverServer() (list []ontology.IsResource, err error) {
	// List all datacenters
	datacenters, _, err := d.clients.computeClient.DataCentersApi.DatacentersGet(context.Background()).Execute()
	if err != nil {
		return nil, fmt.Errorf("could not list datacenters: %w", err)
	}

	// List all servers from the given datacenters
	for _, dc := range *datacenters.Items {
		// List all servers in the datacenter
		servers, _, err := d.clients.computeClient.ServersApi.DatacentersServersGet(context.Background(), util.Deref(dc.Id)).Depth(2).Execute() // Depth(3) to include the volumes and NICs
		if err != nil {
			return nil, fmt.Errorf("could not list servers for datacenter %s: %w", util.Deref(dc.Id), err)
		}

		for _, server := range *servers.Items {
			r, err := d.handleServer(server, dc)
			if err != nil {
				return nil, fmt.Errorf("could not handle server %s: %w", util.Deref(server.Id), err)
			}

			log.Debug("Adding resource %+w", r.GetId())

			list = append(list, r)
		}

	}

	return list, nil
}
