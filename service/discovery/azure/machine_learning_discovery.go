// Copyright 2020-2024 Fraunhofer AISEC
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

package azure

import (
	"context"
	"fmt"

	"clouditor.io/clouditor/v2/api/ontology"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning"
)

// Discover machine learning workspace
func (d *azureDiscovery) discoverMLWorkspaces() ([]ontology.IsResource, error) {
	var list []ontology.IsResource

	// initialize machine learning client
	if err := d.initMachineLearningClient(); err != nil {
		return nil, err
	}

	// List all ML workspaces
	serverListPager := d.clients.mlWorkspaceClient.NewListBySubscriptionPager(&armmachinelearning.WorkspacesClientListBySubscriptionOptions{})
	for serverListPager.More() {
		pageResponse, err := serverListPager.NextPage(context.TODO())
		if err != nil {
			log.Errorf("%s: %v", ErrGettingNextPage, err)
			return list, err
		}

		// TODO(anatheka): New Resource to Cloud Ontology, but where?
		// Add storage, atRestEncryption (keyVault), ...?
		for _, value := range pageResponse.Value {
			workspace, err := handleMLWorkspace(value)

			if err != nil {
				return nil, fmt.Errorf("could not handle ML workspace: %w", err)
			}

			log.Infof("Adding ML workspace '%s'", workspace.GetName())

			list = append(list, workspace)
		}
	}

	return list, nil
}