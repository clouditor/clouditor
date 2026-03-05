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
	"strings"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/v2/pagination"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// handleNetworkInterfaces creates a network interface resource based on the Clouditor Ontology
func (d *openstackDiscovery) handleNetworkInterfaces(network *networks.Network) (ontology.IsResource, error) {
	var (
		l3FirewallEnabled   = false
		restrictedPortsList = []string{}
	)

	// Check if any port associated with the network has security groups enabled.
	// If at least one port has security groups, we consider the L3 firewall to be enabled for the entire network.
	pager := ports.List(d.clients.networkClient, ports.ListOpts{
		NetworkID: network.ID,
	})
	err := pager.EachPage(context.Background(), func(ctx context.Context, page pagination.Page) (bool, error) {
		portList, err := ports.ExtractPorts(page)
		if err != nil {
			return false, err
		}

		for _, port := range portList {
			if len(port.SecurityGroups) > 0 {
				l3FirewallEnabled = true
				restrictedPortsList = append(restrictedPortsList, d.getRestrictedPorts(port.SecurityGroups)...)
				return true, nil
			}
		}
		return true, nil
	})
	if err != nil {
		log.Errorf("Error listing ports for network %s: %s", network.ID, err)
	}

	r := &ontology.NetworkInterface{
		Id:           network.ID,
		Name:         network.Name,
		Description:  network.Description,
		CreationTime: timestamppb.New(network.CreatedAt),
		GeoLocation: &ontology.GeoLocation{
			Region: d.region,
		},
		Labels:   labels(util.Ref(network.Tags)),
		ParentId: util.Ref(network.ProjectID),
		Raw:      discovery.Raw(network),
		AccessRestriction: &ontology.AccessRestriction{
			Type: &ontology.AccessRestriction_L3Firewall{
				L3Firewall: &ontology.L3Firewall{
					Enabled:         l3FirewallEnabled,
					RestrictedPorts: strings.Join(restrictedPortsList, ", "),
				},
			},
		},
	}

	log.Infof("Adding network interface '%s", network.Name)

	return r, nil
}

func (d *openstackDiscovery) handlePorts(network *networks.Network) (ontology.IsResource, error) {

	return nil, nil
}
