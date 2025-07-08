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

package ionos

import (
	"clouditor.io/clouditor/v2/internal/util"
	"github.com/ionos-cloud/sdk-go-bundle/products/compute"
)

// getBlockStorageIds lists the block storage IDs attached to the given server
func getBlockStorageIds(server compute.Server) (blockStorageIds []string) {
	for _, volume := range util.Deref(server.Entities.Volumes.Items) {
		blockStorageIds = append(blockStorageIds, util.Deref(volume.GetId()))
	}

	return blockStorageIds
}

// getNetworkInterfaceIds lists the network interface IDs attached to the given server
func getNetworkInterfaceIds(server compute.Server) (NetworkInterfaceIds []string) {
	for _, nic := range util.Deref(server.Entities.Nics.Items) {
		NetworkInterfaceIds = append(NetworkInterfaceIds, util.Deref(nic.GetId()))
	}

	return NetworkInterfaceIds
}
