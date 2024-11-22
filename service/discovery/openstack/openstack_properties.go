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

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/attachinterfaces"
	"github.com/gophercloud/gophercloud/pagination"
)

// TODO(anatheka): TBD: How is the string formatted?
// labels converts the resource tags to the ontology label
func labels(tags *[]string) map[string]string {
	l := make(map[string]string)

	// for i := range util.Deref(tags) {
	// }

	return l
}

func (d *openstackDiscovery) getAttachedNetworkInterfaces(serverID string) ([]string, error) {
	var (
		list []string
		err  error
	)

	if err = d.authorize(); err != nil {
		return nil, fmt.Errorf("could not authorize openstack: %w", err)
	}

	err = attachinterfaces.List(d.clients.computeClient, serverID).EachPage(func(p pagination.Page) (bool, error) {
		ifc, err := attachinterfaces.ExtractInterfaces(p)
		if err != nil {
			return false, fmt.Errorf("could not extract network interface from page: %w", err)
		}

		for _, i := range ifc {
			list = append(list, i.PortID)
		}

		return true, nil
	})

	if err != nil {
		return nil, fmt.Errorf("could not list network interfaces: %w", err)
	}

	return list, nil
}
