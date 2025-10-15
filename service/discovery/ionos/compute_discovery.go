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

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

// discoverDatacenters lists all datacenters in the IONOS cloud and returns them as a list of ontology resources
func (d *ionosDiscovery) discoverDatacenters() (*ionoscloud.Datacenters, error) {
	// List all datacenters
	dc, _, err := d.clients.computeClient.DataCentersApi.DatacentersGet(context.Background()).Depth(1).Execute()
	if err != nil {
		return util.Ref(dc), fmt.Errorf("could not list datacenters: %w", err)
	}

	// for _, datacenter := range util.Deref(dc.Items) {
	// }

	return util.Ref(dc), err
}

// discoverResources lists all resources in the given datacenter and returns them as a list of ontology resources
func (d *ionosDiscovery) discoverServers(dc ionoscloud.Datacenter) (list []ontology.IsResource, err error) {
	// List all servers in the datacenter
	servers, _, err := d.clients.computeClient.ServersApi.DatacentersServersGet(context.Background(), util.Deref(dc.Id)).Depth(3).Execute() // Depth(3) to include the volumes and NICs
	if err != nil {
		return nil, fmt.Errorf("could not list servers for datacenter %s: %w", util.Deref(dc.Id), err)
	}

	for _, server := range util.Deref(servers.Items) {
		r, err := d.handleServer(server, dc)
		if err != nil {
			return nil, fmt.Errorf("could not handle server %s: %w", util.Deref(server.Id), err)
		}

		log.Debug("Adding server %+w", r.GetId())

		list = append(list, r)
	}

	return list, nil
}

// discoverNetworks lists all networks in the given datacenter and returns them as a list of ontology resources
func (d *ionosDiscovery) discoverBlockStorages(dc ionoscloud.Datacenter) (list []ontology.IsResource, err error) {
	blockStorages, _, err := d.clients.computeClient.VolumesApi.DatacentersVolumesGet(context.Background(), util.Deref(dc.Id)).Depth(1).Execute()
	if err != nil {
		return nil, fmt.Errorf("could not list block storages for datacenter %s: %w", util.Deref(dc.Id), err)
	}
	for _, blockStorage := range util.Deref(blockStorages.Items) {
		r, err := d.handleBlockStorage(blockStorage, dc)
		if err != nil {
			return nil, fmt.Errorf("could not handle block storage %s: %w", util.Deref(blockStorage.Id), err)
		}

		log.Debug("Adding block storage %+w", r.GetId())

		list = append(list, r)
	}
	return list, nil
}

// discoverLoadBalancers lists all load balancers in the given datacenter and returns them as a list of ontology resources
func (d *ionosDiscovery) discoverLoadBalancers(dc ionoscloud.Datacenter) (list []ontology.IsResource, err error) {
	loadBalancers, _, err := d.clients.computeClient.LoadBalancersApi.DatacentersLoadbalancersGet(context.Background(), util.Deref(dc.Id)).Depth(1).Execute()
	if err != nil {
		return nil, fmt.Errorf("could not list load balancers for datacenter %s: %w", util.Deref(dc.Id), err)
	}
	for _, loadBalancer := range util.Deref(loadBalancers.Items) {
		r, err := d.handleLoadBalancer(loadBalancer, dc)
		if err != nil {
			return nil, fmt.Errorf("could not handle load balancer %s: %w", util.Deref(loadBalancer.Id), err)
		}

		log.Debug("Adding load balancer %+w", r.GetId())

		list = append(list, r)
	}
	return list, nil
}
