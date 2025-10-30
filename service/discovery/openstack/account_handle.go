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

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/domains"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
)

// handleDomain returns a [ontology.Account] out of an existing [domains.Domain].
func (d *openstackDiscovery) handleDomain(domain *domains.Domain) (ontology.IsResource, error) {
	r := &ontology.Account{
		Id:           domain.ID,
		Name:         domain.Name,
		Description:  domain.Description,
		CreationTime: nil, // domain does not have a creation date
		GeoLocation:  nil, // domain is global
		Labels:       nil, // domain does not have labels,
		ParentId:     nil, // domain is the top-most item and have no parent,
		Raw:          discovery.Raw(domain),
	}

	log.Infof("Adding domain '%s", r.Name)

	return r, nil
}

// handleProject returns a [ontology.ResourceGroup] out of an existing [projects.Project].
func (d *openstackDiscovery) handleProject(project *projects.Project) (ontology.IsResource, error) {
	r := &ontology.ResourceGroup{
		Id:          project.ID,
		Name:        project.Name,
		Description: project.Description,

		CreationTime: nil, // project does not have a creation date
		GeoLocation: &ontology.GeoLocation{
			Region: d.region,
		},
		Labels:   labels(util.Ref(project.Tags)),
		ParentId: util.Ref(project.ParentID),
		Raw:      discovery.Raw(project),
	}

	log.Infof("Adding project '%s", r.Name)

	return r, nil
}

// addProjectIfMissing checks if the project information is available and adds it to the list of projects.
// If the project is already in the list, it will skip the creation.
func (d *openstackDiscovery) addProjectIfMissing(projectID, projectName, domainID string) error {
	if projectID == "" || projectName == "" || domainID == "" {
		return fmt.Errorf("cannot create project resource: project ID, project name, or domain ID is empty")
	}

	if _, ok := d.discoveredProjects[projectID]; ok {
		log.Debugf("Project with ID '%s' already exists, skipping creation", projectID)
		return nil
	}

	r := &ontology.ResourceGroup{
		Id:       projectID,
		Name:     projectName,
		ParentId: util.Ref(domainID),
		Raw:      discovery.Raw("Project/Tenant information manually added."),
	}

	// Add project to the list of projects
	d.discoveredProjects[r.GetId()] = r

	return nil
}
