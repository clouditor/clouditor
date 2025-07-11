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
	"fmt"

	"clouditor.io/clouditor/v2/internal/util"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/attachinterfaces"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// labels converts the resource tags to the ontology label
func labels(tags *[]string) map[string]string {
	l := make(map[string]string)

	for _, tag := range util.Deref(tags) {
		l[tag] = ""
	}

	return l
}

// getAttachedNetworkInterfaces gets the attached network interfaces to the given serverID.
func (d *openstackDiscovery) getAttachedNetworkInterfaces(serverID string) ([]string, error) {
	var (
		list []string
		err  error
	)

	if err = d.authorize(); err != nil {
		return nil, fmt.Errorf("could not authorize openstack: %w", err)
	}

	err = attachinterfaces.List(d.clients.computeClient, serverID).EachPage(context.Background(), func(_ context.Context, p pagination.Page) (bool, error) {
		ifc, err := attachinterfaces.ExtractInterfaces(p)
		if err != nil {
			return false, fmt.Errorf("could not extract network interface from page: %w", err)
		}

		for _, i := range ifc {
			list = append(list, i.NetID)
		}

		return true, nil
	})

	if err != nil {
		return nil, fmt.Errorf("could not list network interfaces: %w", err)
	}

	return list, nil
}

// setProjectInfo stores the project ID and name based on the given resource
func (d *openstackDiscovery) setProjectInfo(x interface{}) {

	switch v := x.(type) {
	case []volumes.Volume:
		d.project.projectID = v[0].TenantID
		d.project.projectName = v[0].TenantID // it is not possible to extract the project name
	case []servers.Server:
		d.project.projectID = v[0].TenantID
		d.project.projectName = v[0].TenantID // it is not possible to extract the project name
	case []networks.Network:
		d.project.projectID = v[0].TenantID
		d.project.projectName = v[0].TenantID // it is not possible to extract the project name
	default:
		log.Debugf("no known resource type found")
	}
}

// findOrCreateResourceGroupResource This method checks if a resourceGroup resource already is available in the list. Otherwise, it creates a new one.
func (d *openstackDiscovery) findOrCreateResourceGroupResource() {

}
