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
func (d *ionosDiscovery) discoverDatacenters() (*ionoscloud.Datacenters, []ontology.IsResource, error) {
	var (
		list []ontology.IsResource
	)
	// List all datacenters
	dc, _, err := d.clients.client.DataCentersApi.DatacentersGet(context.Background()).Depth(1).Execute()
	if err != nil {
		return nil, nil, fmt.Errorf("could not list datacenters: %w", err)
	}

	for _, datacenter := range util.Deref(dc.Items) {
		r, err := d.handleDatacenter(datacenter)
		if err != nil {
			return nil, nil, fmt.Errorf("could not handle datacenter %s: %w", util.Deref(datacenter.Id), err)
		}

		log.Infof("Adding datacenter '%s'", r.GetId())

		list = append(list, r)
	}

	return util.Ref(dc), list, err
}
