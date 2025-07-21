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
	"context"
	"fmt"

	"clouditor.io/clouditor/v2/internal/util"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/attachinterfaces"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/domains"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// labels converts the resource tags to the ontology label
func labels(tags *[]string) map[string]string {
	l := make(map[string]string)

	for _, tag := range util.Deref(tags) {
		l[tag] = ""
	}

	return l
}

// getAttachedNetworkInterfaces gets the attached network interfaces to the given serverID.
func (d *openstackDiscovery) getAttachedNetworkInterfaces(serverID string) ([]string, error) {
	var (
		list []string
		err  error
	)

	if err = d.authorize(); err != nil {
		return nil, fmt.Errorf("could not authorize openstack: %w", err)
	}

	err = attachinterfaces.List(d.clients.computeClient, serverID).EachPage(context.Background(), func(_ context.Context, p pagination.Page) (bool, error) {
		ifc, err := attachinterfaces.ExtractInterfaces(p)
		if err != nil {
			return false, fmt.Errorf("could not extract network interface from page: %w", err)
		}

		for _, i := range ifc {
			list = append(list, i.NetID)
		}

		return true, nil
	})

	if err != nil {
		return nil, fmt.Errorf("could not list network interfaces: %w", err)
	}

	return list, nil
}

// setProjectInfo stores the project ID and name based on the given resource
func (d *openstackDiscovery) setProjectInfo(x interface{}) error {
	var (
		projectID   string
		projectName string
		err         error
	)

	switch resource := x.(type) {
	case []servers.Server:
		projectID, err = getProjectID(resource[0])
		projectName = projectID // it is not possible to extract the project name
	case []networks.Network:
		projectID, err = getProjectID(resource[0])
		projectName = projectID // it is not possible to extract the project name
	case []volumes.Volume:
		projectID, err = getProjectID(resource[0])
		projectName = projectID // it is not possible to extract the project name
	case []projects.Project:
		projectID = resource[0].ID
		projectName = resource[0].Name
	case []domains.Domain:
		// Domain does not have a project ID or name, so we skip this
		return nil
	default:
		return fmt.Errorf("unknown resource type: %T", resource)
	}

	if err != nil {
		return fmt.Errorf("error getting project ID")
	}

	d.project.projectID = projectID
	d.project.projectName = projectName
	return nil
}

type resourceTypes interface {
	servers.Server | *servers.Server |
		networks.Network | *networks.Network |
		volumes.Volume | *volumes.Volume
}

// getProjectID returns the project/tenant ID from the given ionoscloud resouce object
func getProjectID[T resourceTypes](r T) (string, error) {
	switch v := any(r).(type) {
	case volumes.Volume:
		if v.TenantID != "" {
			return v.TenantID, nil
		}
	case *volumes.Volume:
		if v != nil && v.TenantID != "" {
			return v.TenantID, nil
		}

	case servers.Server:
		if v.TenantID != "" {
			return v.TenantID, nil
		}
	case *servers.Server:
		if v != nil && v.TenantID != "" {
			return v.TenantID, nil
		}

	case networks.Network:
		if v.TenantID != "" {
			return v.TenantID, nil
		}
		if v.ProjectID != "" {
			return v.ProjectID, nil
		}
	case *networks.Network:
		if v != nil {
			if v.TenantID != "" {
				return v.TenantID, nil
			}
			if v.ProjectID != "" {
				return v.ProjectID, nil
			}
		}
	default:
		return "", fmt.Errorf("unknown resource type: %T", r)
	}

	return "", fmt.Errorf("no tenant or project ID available")

}
