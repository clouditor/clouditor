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

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/util"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// handleDatacenter creates a datacenter resource based on the Clouditor Ontology
func (d *ionosDiscovery) handleDatacenter(dc ionoscloud.Datacenter) (ontology.IsResource, error) {
	// Getting labels
	l, _, err := d.clients.computeClient.LabelsApi.
		DatacentersLabelsGet(context.Background(), util.Deref(dc.GetId())).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("error getting labels for datacenter %s: %s", util.Deref(dc.Id), err)
	}

	r := &ontology.Account{
		Id:           util.Deref(dc.Id),
		Name:         util.Deref(dc.Properties.Name),
		CreationTime: timestamppb.New(util.Deref(dc.Metadata.GetCreatedDate())),
		GeoLocation: &ontology.GeoLocation{
			Region: util.Deref(dc.Properties.Location),
		},
		Labels:      labels(l),
		Raw:         discovery.Raw(dc),
		Description: util.Deref(dc.Properties.Description),
	}

	return r, nil
}
