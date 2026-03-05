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

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

func (d *openstackDiscovery) getRestrictedPorts(portSecurityGroup []string) []string {
	var (
		restrictedPortsList []string
	)

	for _, sgID := range portSecurityGroup {
		// Get security group details
		pager := groups.List(d.clients.networkClient, groups.ListOpts{
			ID: sgID,
		})

		err := pager.EachPage(context.Background(), func(ctx context.Context, page pagination.Page) (bool, error) {
			sgList, err := groups.ExtractGroups(page)
			if err != nil {
				return false, err
			}

			for _, sg := range sgList {
				for _, rule := range sg.Rules {
					if rule.Direction == "ingress" {
						if rule.Protocol == "tcp" || rule.Protocol == "udp" {
							// Check if the port range is specified
							if rule.PortRangeMin == 0 && rule.PortRangeMax == 0 {
								// If no port range is specified, it means all ports are allowed
								restrictedPortsList = append(restrictedPortsList, "all")
							} else if rule.PortRangeMin == rule.PortRangeMax {
								// If the port range is a single port, add that port to the list
								restrictedPortsList = append(restrictedPortsList, string(rule.PortRangeMin))
							} else {
								// If the port range is specified, add each port in the range to the list
								for port := rule.PortRangeMin; port <= rule.PortRangeMax; port++ {
									restrictedPortsList = append(restrictedPortsList, string(port))
								}

							}
						}
					}
				}
			}

			return true, nil
		})
		if err != nil {
			continue
		}
	}

	return restrictedPortsList
}
