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

package openstacktest

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// The code in this file is based on the Gophercloud library.
// Gophercloud is licensed under the Apache License 2.0. See the LICENSE file in the
// Gophercloud repository for the full license: https://github.com/gophercloud/gophercloud
// Source: https://github.com/gophercloud/gophercloud/blob/master/openstack/networking/v2/networks/testing/fixtures.go (2024-12-09)
//
// Changes:
// - 2024-12-10: Add function HandleNetworkListSuccessfully() and add second path for "/networks" based on the HandleFunc in https://github.com/gophercloud/gophercloud/blob/5770765aa037e1572cbaa9474113010a1397e822/openstack/networking/v2/networks/testing/requests_test.go (anatheka)
// - 2025-01-28: Add ProjectID property to all network objects (anatheka)

const ListResponse = `
{
    "networks": [
        {
            "status": "ACTIVE",
            "subnets": [
                "54d6f61d-db07-451c-9ab3-b9609b6b6f0b"
            ],
            "name": "public",
            "admin_state_up": true,
            "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
            "project_id": "4fd44f30292945e481c7b8a0c8908869",
            "created_at": "2019-06-30T04:15:37",
            "updated_at": "2019-06-30T05:18:49",
            "shared": true,
            "id": "d32019d3-bc6e-4319-9c1d-6722fc136a22",
            "provider:segmentation_id": 9876543210,
            "provider:physical_network": null,
            "provider:network_type": "local",
            "router:external": true,
            "port_security_enabled": true,
            "dns_domain": "local.",
            "mtu": 1500
        },
        {
            "status": "ACTIVE",
            "subnets": [
                "08eae331-0402-425a-923c-34f7cfe39c1b"
            ],
            "name": "private",
            "admin_state_up": true,
            "tenant_id": "26a7980765d0414dbc1fc1f88cdb7e6e",
            "created_at": "2019-06-30T04:15:37Z",
            "updated_at": "2019-06-30T05:18:49Z",
            "shared": false,
            "id": "db193ab3-96e3-4cb3-8fc5-05f4296d0324",
            "provider:segmentation_id": 1234567890,
            "provider:physical_network": null,
            "provider:network_type": "local",
            "router:external": false,
            "port_security_enabled": false,
            "dns_domain": "",
            "mtu": 1500
        }
    ]
}`

const GetResponse = `
{
    "network": {
        "status": "ACTIVE",
        "subnets": [
            "54d6f61d-db07-451c-9ab3-b9609b6b6f0b"
        ],
        "name": "public",
        "admin_state_up": true,
        "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
        "project_id": "4fd44f30292945e481c7b8a0c8908869",
        "created_at": "2019-06-30T04:15:37",
        "updated_at": "2019-06-30T05:18:49",
        "shared": true,
        "id": "d32019d3-bc6e-4319-9c1d-6722fc136a22",
        "provider:segmentation_id": 9876543210,
        "provider:physical_network": null,
        "provider:network_type": "local",
        "router:external": true,
        "port_security_enabled": true,
        "dns_domain": "local.",
        "mtu": 1500
    }
}`

const CreateRequest = `
{
    "network": {
        "name": "private",
        "admin_state_up": true
    }
}`

const CreateResponse = `
{
    "network": {
        "status": "ACTIVE",
        "subnets": ["08eae331-0402-425a-923c-34f7cfe39c1b"],
        "name": "private",
        "admin_state_up": true,
        "tenant_id": "26a7980765d0414dbc1fc1f88cdb7e6e",
        "created_at": "2019-06-30T04:15:37Z",
        "updated_at": "2019-06-30T05:18:49Z",
        "shared": false,
        "id": "db193ab3-96e3-4cb3-8fc5-05f4296d0324",
        "provider:segmentation_id": 9876543210,
        "provider:physical_network": null,
        "provider:network_type": "local",
        "dns_domain": ""
    }
}`

const CreatePortSecurityRequest = `
{
    "network": {
        "name": "private",
        "admin_state_up": true,
        "port_security_enabled": false
    }
}`

const CreatePortSecurityResponse = `
{
    "network": {
        "status": "ACTIVE",
        "subnets": ["08eae331-0402-425a-923c-34f7cfe39c1b"],
        "name": "private",
        "admin_state_up": true,
        "tenant_id": "26a7980765d0414dbc1fc1f88cdb7e6e",
        "created_at": "2019-06-30T04:15:37Z",
        "updated_at": "2019-06-30T05:18:49Z",
        "shared": false,
        "id": "db193ab3-96e3-4cb3-8fc5-05f4296d0324",
        "provider:segmentation_id": 9876543210,
        "provider:physical_network": null,
        "provider:network_type": "local",
        "port_security_enabled": false
    }
}`

const CreateOptionalFieldsRequest = `
{
  "network": {
      "name": "public",
      "admin_state_up": true,
      "shared": true,
      "tenant_id": "12345",
      "availability_zone_hints": ["zone1", "zone2"]
  }
}`

const UpdateRequest = `
{
    "network": {
        "name": "new_network_name",
        "admin_state_up": false,
        "shared": true
    }
}`

const UpdateResponse = `
{
    "network": {
        "status": "ACTIVE",
        "subnets": [],
        "name": "new_network_name",
        "admin_state_up": false,
        "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
        "project_id": "4fd44f30292945e481c7b8a0c8908869",
        "created_at": "2019-06-30T04:15:37Z",
        "updated_at": "2019-06-30T05:18:49Z",
        "shared": true,
        "id": "4e8e5957-649f-477b-9e5b-f1f75b21c03c",
        "provider:segmentation_id": 1234567890,
        "provider:physical_network": null,
        "provider:network_type": "local"
    }
}`

const UpdatePortSecurityRequest = `
{
    "network": {
        "port_security_enabled": false
    }
}`

const UpdatePortSecurityResponse = `
{
    "network": {
        "status": "ACTIVE",
        "subnets": ["08eae331-0402-425a-923c-34f7cfe39c1b"],
        "name": "private",
        "admin_state_up": true,
        "tenant_id": "26a7980765d0414dbc1fc1f88cdb7e6e",
        "created_at": "2019-06-30T04:15:37Z",
        "updated_at": "2019-06-30T05:18:49Z",
        "shared": false,
        "id": "4e8e5957-649f-477b-9e5b-f1f75b21c03c",
        "provider:segmentation_id": 9876543210,
        "provider:physical_network": null,
        "provider:network_type": "local",
        "port_security_enabled": false
    }
}`

var createdTime, _ = time.Parse(time.RFC3339, "2019-06-30T04:15:37Z")
var updatedTime, _ = time.Parse(time.RFC3339, "2019-06-30T05:18:49Z")

var (
	Network1 = networks.Network{
		Status:       "ACTIVE",
		Subnets:      []string{"54d6f61d-db07-451c-9ab3-b9609b6b6f0b"},
		Name:         "public",
		AdminStateUp: true,
		TenantID:     "4fd44f30292945e481c7b8a0c8908869",
		ProjectID:    "4fd44f30292945e481c7b8a0c8908869",
		CreatedAt:    createdTime,
		UpdatedAt:    updatedTime,
		Shared:       true,
		ID:           "d32019d3-bc6e-4319-9c1d-6722fc136a22",
	}
)

var (
	Network2 = networks.Network{
		Status:       "ACTIVE",
		Subnets:      []string{"08eae331-0402-425a-923c-34f7cfe39c1b"},
		Name:         "private",
		AdminStateUp: true,
		TenantID:     "26a7980765d0414dbc1fc1f88cdb7e6e",
		ProjectID:    "26a7980765d0414dbc1fc1f88cdb7e6e",
		CreatedAt:    createdTime,
		UpdatedAt:    updatedTime,
		Shared:       false,
		ID:           "db193ab3-96e3-4cb3-8fc5-05f4296d0324",
	}
)

var ExpectedNetworkSlice = []networks.Network{Network1, Network2}

func HandleNetworkListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/networks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ListResponse)
	})

	th.Mux.HandleFunc("/v2.0/networks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ListResponse)
	})
}
