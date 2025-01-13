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

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// handleServer creates a virtual machine resource based on the Clouditor Ontology
func (d *openstackDiscovery) handleServer(server *servers.Server) (ontology.IsResource, error) {
	var (
		err         error
		bootLogging *ontology.BootLogging
		osLogging   *ontology.OSLogging
	)

	// TODO(anatheka): Check again!
	// boot and os logging are logged together in the console log
	consoleOutput := servers.ShowConsoleOutput(context.Background(), d.clients.computeClient, server.ID, servers.ShowConsoleOutputOpts{})
	if consoleOutput.Result.Err == nil {
		bootLogging = &ontology.BootLogging{
			Enabled: true,
		}
		// osLogging = &ontology.OSLogging{
		// 	Enabled: true,
		// }
	} else {
		log.Errorf("Error getting boot logging: %s", consoleOutput.Err)
		// When an error occurs, we assume that boot logging is disabled.
		bootLogging = &ontology.BootLogging{
			Enabled: false,
		}
	}

	r := &ontology.VirtualMachine{
		Id:           server.ID,
		Name:         server.Name,
		CreationTime: timestamppb.New(server.Created),
		GeoLocation: &ontology.GeoLocation{
			Region: d.region, // TODO: Can we get the region?
		},
		Labels:            labels(server.Tags),
		ParentId:          util.Ref(server.TenantID),
		Raw:               discovery.Raw(server),
		MalwareProtection: &ontology.MalwareProtection{},
		BootLogging:       bootLogging,
		OsLogging:         osLogging,
		ActivityLogging:   &ontology.ActivityLogging{},
		AutomaticUpdates:  &ontology.AutomaticUpdates{},
	}

	// Get attached block storage IDs
	for _, v := range server.AttachedVolumes {
		r.BlockStorageIds = append(r.BlockStorageIds, v.ID)
	}

	// Get attached network interface IDs
	r.NetworkInterfaceIds, err = d.getAttachedNetworkInterfaces(server.ID)
	if err != nil {
		return nil, fmt.Errorf("could not discover attached network interfaces: %w", err)
	}

	log.Infof("Adding server '%s", server.Name)

	return r, nil
}
