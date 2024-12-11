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
	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/domains"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
)

// handleProject returns a [ontology.ResourceGroup] out of an existing [projects.Project].
func (d *openstackDiscovery) handleProject(project *projects.Project) (ontology.IsResource, error) {
	r := &ontology.ResourceGroup{
		Id:           project.ID,
		Name:         project.Name,
		CreationTime: nil, // project does not have a creation date
		GeoLocation: &ontology.GeoLocation{
			Region: "unknown", // TODO: Can we get the region?
		},
		Description: project.Description,
		Labels:      labels(util.Ref(project.Tags)),
		ParentId:    util.Ref(project.ParentID),
		Raw:         discovery.Raw(project),
	}

	log.Infof("Adding resource group '%s", project.Name)

	return r, nil
}

// handleDomain returns a [ontology.Account] out of an existing [domains.Domain].
func (d *openstackDiscovery) handleDomain(domain *domains.Domain) (ontology.IsResource, error) {
	r := &ontology.Account{
		Id:           domain.ID,
		Name:         domain.Name,
		Description:  domain.Description,
		CreationTime: nil, // domain do not have a creation date
		GeoLocation:  nil, // domain are global
		Labels:       nil, // domain do not have labels,
		ParentId:     nil, // domain are the top-most item and have no parent,
		Raw:          discovery.Raw(domain),
	}

	log.Infof("Adding domain '%s", domain.Name)

	return r, nil
}
