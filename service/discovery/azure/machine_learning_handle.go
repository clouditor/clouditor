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
	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (d *azureDiscovery) handleMLWorkspace(value *armmachinelearning.Workspace, computeList []string) (ontology.IsResource, error) {
	ml := &ontology.MachineLearningService{
		Id:                         resourceID(value.ID),
		Name:                       util.Deref(value.Name),
		CreationTime:               creationTime(value.SystemData.CreatedAt),
		GeoLocation:                location(value.Location),
		Labels:                     labels(value.Tags),
		ParentId:                   resourceGroupID(resourceIDPointer(value.ID)),
		Raw:                        discovery.Raw(value),
		InternetAccessibleEndpoint: getInternetAccessibleEndpoint(value.Properties.PublicNetworkAccess),
		StorageIds:                 []string{util.Deref(value.Properties.StorageAccount)},
		ComputeIds:                 computeList,
		Loggings: []*ontology.Logging{
			{
				Type: &ontology.Logging_ResourceLogging{
					ResourceLogging: getResourceLogging(value.Properties.ApplicationInsights),
				},
			},
		},
	}

	return ml, nil
}

// TODO(all): Should we move that to the compute file
func (d *azureDiscovery) handleMLCompute(value *armmachinelearning.ComputeResource, workspaceID *string) (ontology.IsResource, error) {
	var (
		compute   *ontology.VirtualMachine
		container *ontology.Container
		time      = &timestamppb.Timestamp{}
	)

	// Get properties vom ComputeResource
	if value.SystemData != nil && value.SystemData.CreatedAt != nil {
		time = creationTime(value.SystemData.CreatedAt)
	}

	// Get compute type specific properties for "VirtualMachine" or "ComputeInstance"
	switch c := value.Properties.(type) {
	case *armmachinelearning.ComputeInstance:
		container = &ontology.Container{
			Id:                  resourceID(value.ID),
			Name:                util.Deref(value.Name),
			CreationTime:        time,
			GeoLocation:         location(value.Location),
			Labels:              labels(value.Tags),
			ParentId:            resourceIDPointer(workspaceID),
			Raw:                 discovery.Raw(value, c.ComputeLocation),
			NetworkInterfaceIds: []string{},
		}
		return container, nil
	case *armmachinelearning.VirtualMachine:

		compute = &ontology.VirtualMachine{
			Id:                  resourceID(value.ID),
			Name:                util.Deref(value.Name),
			CreationTime:        time,
			GeoLocation:         location(value.Location),
			Labels:              labels(value.Tags),
			ParentId:            resourceIDPointer(workspaceID),
			Raw:                 discovery.Raw(value, c.ComputeLocation),
			NetworkInterfaceIds: []string{},
			BlockStorageIds:     []string{},
			MalwareProtection:   &ontology.MalwareProtection{},
		}

		return compute, nil
	}

	log.Debugf("Couldn't handle value '%s' because type '%s' is not supported.",
		util.Deref(value.Name), util.Deref(value.Type))
	return nil, nil
}
