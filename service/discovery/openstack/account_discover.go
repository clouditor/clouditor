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
	"clouditor.io/clouditor/v2/api/ontology"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/domains"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
)

// discoverDomains discovers domains.
func (d *openstackDiscovery) discoverDomains() (list []ontology.IsResource, err error) {
	var opts domains.ListOptsBuilder = &domains.ListOpts{}
	list, err = genericList(d, d.identityClient, domains.List, d.handleDomain, domains.ExtractDomains, opts)

	if err != nil {
		// if we cannot retrieve the domain information by calling the API or from the environment variables, we will not be able to succeed
		if d.domainID == "" || d.domainName == "" {
			return nil, err
		}

		// TODO(all): Or should it be Errorf?
		log.Debugf("could not discover domains due to insufficient permissions, but we can proceed with less domain information: %v", err)
		r := &ontology.Account{
			Id:   d.domainID,
			Name: d.domainName,
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
		// TODO(all): Or should it be Errorf?
		log.Debugf("could not discover projects/tenants due to insufficient permissions, but we can proceed with less project/tenant information: %v", err)

		r := &ontology.ResourceGroup{
			Id:       "6b715d5b91964beaa14100f011dc6339", //TODO(anatheka): Add information
			Name:     "testName",                         //TODO(anatheka): Add information
			ParentId: &d.domainID,
		}

		// Set projectID/projectName
		d.projectID = ""   //TODO(anatheka): TBD
		d.projectName = "" //TODO(anatheka): TBD

		list = append(list, r)
	}

	return list, nil
}
