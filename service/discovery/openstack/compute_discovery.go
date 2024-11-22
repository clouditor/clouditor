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

	"clouditor.io/clouditor/v2/api/ontology"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

// List lists OpenStack servers (compute resources) and translates them into the Clouditor ontology
func (d *openstackDiscovery) List() (list []ontology.IsResource, err error) {

	// TODO(anatheka): Add Discover network interfaces

	// Discover servers
	servers, err := d.discoverServers()
	if err != nil {
		return nil, fmt.Errorf("could not discover servers: %w", err)
	}
	list = append(list, servers...)

	return
}

func (d *openstackDiscovery) discoverServers() (list []ontology.IsResource, err error) {
	// TODO(oxisto): Limit the list to a specific tenant?
	var opts servers.ListOptsBuilder = &servers.ListOpts{}
	list, err = genericList(d, d.computeClient, servers.List, d.handleServer, servers.ExtractServers, opts)

	return
}
