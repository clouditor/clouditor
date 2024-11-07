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
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
)

// computeClient returns the compute client if initialized
func (d *openstackDiscovery) computeClient() (err error) {
	if d.clients.computeClient == nil {
		d.clients.computeClient, err = openstack.NewComputeV2(d.provider, gophercloud.EndpointOpts{
			Region: os.Getenv(RegionName),
		})
		if err != nil {
			return fmt.Errorf("could not create compute client: %w", err)
		}
	}
	return nil
}

// storageClient returns the compute client if initialized
func (d *openstackDiscovery) storageClient() (err error) {
	if d.clients.storageClient == nil {
		d.clients.storageClient, err = openstack.NewBlockStorageV3(d.provider, gophercloud.EndpointOpts{
			Region: os.Getenv(RegionName),
		})
		if err != nil {
			return fmt.Errorf("could not create block storage client: %w", err)
		}
	}

	return nil
}
