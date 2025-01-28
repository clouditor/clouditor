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
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/domains"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
)

// discoverDomains discovers domains.
func (d *openstackDiscovery) discoverDomains() (list []ontology.IsResource, err error) {
	var opts domains.ListOptsBuilder = &domains.ListOpts{}
	list, err = genericList(d, d.identityClient, domains.List, d.handleDomain, domains.ExtractDomains, opts)

	if err != nil {
		// if we cannot retrieve the domain information by calling the API or from the environment variables, we will add the information manually if we already got the domain ID and/or domain name
		log.Debugf("could not discover domains due to insufficient permissions, but we can proceed with less domain information: %v", err)

		if d.domain.domainID == "" || d.domain.domainName == "" {
			err := fmt.Errorf("neither the domain ID nor the domain name are available: %v", err)
			return nil, err
		}

		r := &ontology.Account{
			Id:   d.domain.domainID,
			Name: d.domain.domainName,
			Raw:  discovery.Raw("Domain information manually added."),
		}

		list = append(list, r)
	}

	return list, nil
}

// discoverProjects discovers projects/tenants. OpenStack project and tenant are interchangeable.
func (d *openstackDiscovery) discoverProjects() (list []ontology.IsResource, err error) {
	var opts projects.ListOptsBuilder = &projects.ListOpts{}
	list, err = genericList(d, d.identityClient, projects.List, d.handleProject, projects.ExtractProjects, opts)

	if err != nil {
		// if we cannot retrieve the project information by calling the API or from the environment variables, we will add the information manually if we already got the project ID and/or project name
		log.Debugf("could not discover projects/tenants due to insufficient permissions, but we can proceed with less project/tenant information: %v", err)

		if d.project.projectID == "" || d.project.projectName == "" {
			err := fmt.Errorf("neither the project ID nor the project name are available: %v", err)
			return nil, err
		}

		r := &ontology.ResourceGroup{
			Id:       d.project.projectID,
			Name:     d.project.projectName,
			ParentId: &d.domain.domainID,
			Raw:      discovery.Raw("Project/Tenant information manually added."),
		}

		list = append(list, r)
	}

	return list, nil
}
