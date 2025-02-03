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

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/v2/testhelper"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// The code in this file is based on the Gophercloud library.
// Gophercloud is licensed under the Apache License 2.0. See the LICENSE file in the
// Gophercloud repository for the full license: https://github.com/gophercloud/gophercloud
// Source: https://github.com/gophercloud/gophercloud/blob/master/openstack/compute/v2/servers/testing/fixtures_test.go (2024-12-09)
//
// Changes:
// - 2025-01-13: Added function HandleShowConsoleOutputSuccessfullyModified to get console output of server and delete `"length": 50` in TestJSONRequest, otherwise it does not work. (@anatheka)

// ServerListBody contains the canned body of a servers.List response.
const ServerListBody = `
{
	"servers": [
		{
			"status": "ACTIVE",
			"updated": "2014-09-25T13:10:10Z",
			"hostId": "29d3c8c896a45aa4c34e52247875d7fefc3d94bbcc9f622b5d204362",
			"OS-EXT-SRV-ATTR:host": "devstack",
			"addresses": {
				"private": [
					{
						"OS-EXT-IPS-MAC:mac_addr": "fa:16:3e:7c:1b:2b",
						"version": 4,
						"addr": "10.0.0.32",
						"OS-EXT-IPS:type": "fixed"
					}
				]
			},
			"links": [
				{
					"href": "http://104.130.131.164:8774/v2/fcad67a6189847c4aecfa3c81a05783b/servers/ef079b0c-e610-4dfb-b1aa-b49f07ac48e5",
					"rel": "self"
				},
				{
					"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/servers/ef079b0c-e610-4dfb-b1aa-b49f07ac48e5",
					"rel": "bookmark"
				}
			],
			"key_name": null,
			"image": {
				"id": "f90f6034-2570-4974-8351-6b49732ef2eb",
				"links": [
					{
						"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/images/f90f6034-2570-4974-8351-6b49732ef2eb",
						"rel": "bookmark"
					}
				]
			},
			"OS-EXT-STS:task_state": null,
			"OS-EXT-STS:vm_state": "active",
			"OS-EXT-SRV-ATTR:instance_name": "instance-0000001e",
			"OS-SRV-USG:launched_at": "2014-09-25T13:10:10.000000",
			"OS-EXT-SRV-ATTR:hypervisor_hostname": "devstack",
			"flavor": {
				"id": "1",
				"links": [
					{
						"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/flavors/1",
						"rel": "bookmark"
					}
				]
			},
			"id": "ef079b0c-e610-4dfb-b1aa-b49f07ac48e5",
			"security_groups": [
				{
					"name": "default"
				}
			],
			"OS-SRV-USG:terminated_at": null,
			"OS-EXT-AZ:availability_zone": "nova",
			"user_id": "9349aff8be7545ac9d2f1d00999a23cd",
			"name": "herp",
			"created": "2014-09-25T13:10:02Z",
			"tenant_id": "fcad67a6189847c4aecfa3c81a05783b",
			"OS-DCF:diskConfig": "MANUAL",
			"os-extended-volumes:volumes_attached": [
				{
					"id": "2bdbc40f-a277-45d4-94ac-d9881c777d33"
				}
			],
			"accessIPv4": "",
			"accessIPv6": "",
			"progress": 0,
			"OS-EXT-STS:power_state": 1,
			"config_drive": "",
			"metadata": {}
		},
		{
			"status": "ACTIVE",
			"updated": "2014-09-25T13:04:49Z",
			"hostId": "29d3c8c896a45aa4c34e52247875d7fefc3d94bbcc9f622b5d204362",
			"OS-EXT-SRV-ATTR:host": "devstack",
			"addresses": {
				"private": [
					{
						"OS-EXT-IPS-MAC:mac_addr": "fa:16:3e:9e:89:be",
						"version": 4,
						"addr": "10.0.0.31",
						"OS-EXT-IPS:type": "fixed"
					}
				]
			},
			"links": [
				{
					"href": "http://104.130.131.164:8774/v2/fcad67a6189847c4aecfa3c81a05783b/servers/9e5476bd-a4ec-4653-93d6-72c93aa682ba",
					"rel": "self"
				},
				{
					"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/servers/9e5476bd-a4ec-4653-93d6-72c93aa682ba",
					"rel": "bookmark"
				}
			],
			"key_name": null,
			"image": {
				"id": "f90f6034-2570-4974-8351-6b49732ef2eb",
				"links": [
					{
						"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/images/f90f6034-2570-4974-8351-6b49732ef2eb",
						"rel": "bookmark"
					}
				]
			},
			"OS-EXT-STS:task_state": null,
			"OS-EXT-STS:vm_state": "active",
			"OS-EXT-SRV-ATTR:instance_name": "instance-0000001d",
			"OS-SRV-USG:launched_at": "2014-09-25T13:04:49.000000",
			"OS-EXT-SRV-ATTR:hypervisor_hostname": "devstack",
			"flavor": {
				"id": "1",
				"links": [
					{
						"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/flavors/1",
						"rel": "bookmark"
					}
				]
			},
			"id": "9e5476bd-a4ec-4653-93d6-72c93aa682ba",
			"security_groups": [
				{
					"name": "default"
				}
			],
			"OS-SRV-USG:terminated_at": null,
			"OS-EXT-AZ:availability_zone": "nova",
			"user_id": "9349aff8be7545ac9d2f1d00999a23cd",
			"name": "derp",
			"created": "2014-09-25T13:04:41Z",
			"tenant_id": "fcad67a6189847c4aecfa3c81a05783b",
			"OS-DCF:diskConfig": "MANUAL",
			"os-extended-volumes:volumes_attached": [],
			"accessIPv4": "",
			"accessIPv6": "",
			"progress": 0,
			"OS-EXT-STS:power_state": 1,
			"config_drive": "",
			"metadata": {}
		},
		{
		"status": "ACTIVE",
		"updated": "2014-09-25T13:04:49Z",
		"hostId": "29d3c8c896a45aa4c34e52247875d7fefc3d94bbcc9f622b5d204362",
		"OS-EXT-SRV-ATTR:host": "devstack",
		"addresses": {
			"private": [
				{
					"OS-EXT-IPS-MAC:mac_addr": "fa:16:3e:9e:89:be",
					"version": 4,
					"addr": "10.0.0.31",
					"OS-EXT-IPS:type": "fixed"
				}
			]
		},
		"links": [
			{
				"href": "http://104.130.131.164:8774/v2/fcad67a6189847c4aecfa3c81a05783b/servers/9e5476bd-a4ec-4653-93d6-72c93aa682ba",
				"rel": "self"
			},
			{
				"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/servers/9e5476bd-a4ec-4653-93d6-72c93aa682ba",
				"rel": "bookmark"
			}
		],
		"key_name": null,
		"image": "",
		"OS-EXT-STS:task_state": null,
		"OS-EXT-STS:vm_state": "active",
		"OS-EXT-SRV-ATTR:instance_name": "instance-0000001d",
		"OS-SRV-USG:launched_at": "2014-09-25T13:04:49.000000",
		"OS-EXT-SRV-ATTR:hypervisor_hostname": "devstack",
		"flavor": {
			"id": "1",
			"links": [
				{
					"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/flavors/1",
					"rel": "bookmark"
				}
			]
		},
		"id": "9e5476bd-a4ec-4653-93d6-72c93aa682bb",
		"security_groups": [
			{
				"name": "default"
			}
		],
		"OS-SRV-USG:terminated_at": null,
		"OS-EXT-AZ:availability_zone": "nova",
		"user_id": "9349aff8be7545ac9d2f1d00999a23cd",
		"name": "merp",
		"created": "2014-09-25T13:04:41Z",
		"tenant_id": "fcad67a6189847c4aecfa3c81a05783b",
		"OS-DCF:diskConfig": "MANUAL",
		"os-extended-volumes:volumes_attached": [],
		"accessIPv4": "",
		"accessIPv6": "",
		"progress": 0,
		"OS-EXT-STS:power_state": 1,
		"config_drive": "",
		"metadata": {}
	}
	]
}
`

var (
	herpTimeCreated, _ = time.Parse(time.RFC3339, "2014-09-25T13:10:02Z")
	herpTimeUpdated, _ = time.Parse(time.RFC3339, "2014-09-25T13:10:10Z")
	// ServerHerp is a Server struct that should correspond to the first result in ServerListBody.
	ServerHerp = servers.Server{
		Status:  "ACTIVE",
		Updated: herpTimeUpdated,
		HostID:  "29d3c8c896a45aa4c34e52247875d7fefc3d94bbcc9f622b5d204362",
		Addresses: map[string]any{
			"private": []any{
				map[string]any{
					"OS-EXT-IPS-MAC:mac_addr": "fa:16:3e:7c:1b:2b",
					"version":                 float64(4),
					"addr":                    "10.0.0.32",
					"OS-EXT-IPS:type":         "fixed",
				},
			},
		},
		Links: []any{
			map[string]any{
				"href": "http://104.130.131.164:8774/v2/fcad67a6189847c4aecfa3c81a05783b/servers/ef079b0c-e610-4dfb-b1aa-b49f07ac48e5",
				"rel":  "self",
			},
			map[string]any{
				"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/servers/ef079b0c-e610-4dfb-b1aa-b49f07ac48e5",
				"rel":  "bookmark",
			},
		},
		Image: map[string]any{
			"id": "f90f6034-2570-4974-8351-6b49732ef2eb",
			"links": []any{
				map[string]any{
					"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/images/f90f6034-2570-4974-8351-6b49732ef2eb",
					"rel":  "bookmark",
				},
			},
		},
		Flavor: map[string]any{
			"id": "1",
			"links": []any{
				map[string]any{
					"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/flavors/1",
					"rel":  "bookmark",
				},
			},
		},
		ID:       "ef079b0c-e610-4dfb-b1aa-b49f07ac48e5",
		UserID:   "9349aff8be7545ac9d2f1d00999a23cd",
		Name:     "herp",
		Created:  herpTimeCreated,
		TenantID: "fcad67a6189847c4aecfa3c81a05783b",
		Metadata: map[string]string{},
		AttachedVolumes: []servers.AttachedVolume{
			{
				ID: "2bdbc40f-a277-45d4-94ac-d9881c777d33",
			},
		},
		SecurityGroups: []map[string]any{
			{
				"name": "default",
			},
		},
		Host:               "devstack",
		Hostname:           nil,
		HypervisorHostname: "devstack",
		InstanceName:       "instance-0000001e",
		LaunchIndex:        nil,
		ReservationID:      nil,
		RootDeviceName:     nil,
		Userdata:           nil,
		VmState:            "active",
		PowerState:         servers.RUNNING,
		LaunchedAt:         time.Date(2014, 9, 25, 13, 10, 10, 0, time.UTC),
		TerminatedAt:       time.Time{},
		DiskConfig:         servers.Manual,
		AvailabilityZone:   "nova",
	}

	derpTimeCreated, _ = time.Parse(time.RFC3339, "2014-09-25T13:04:41Z")
	derpTimeUpdated, _ = time.Parse(time.RFC3339, "2014-09-25T13:04:49Z")
	// ServerDerp is a Server struct that should correspond to the second server in ServerListBody.
	ServerDerp = servers.Server{
		Status:  "ACTIVE",
		Updated: derpTimeUpdated,
		HostID:  "29d3c8c896a45aa4c34e52247875d7fefc3d94bbcc9f622b5d204362",
		Addresses: map[string]any{
			"private": []any{
				map[string]any{
					"OS-EXT-IPS-MAC:mac_addr": "fa:16:3e:9e:89:be",
					"version":                 float64(4),
					"addr":                    "10.0.0.31",
					"OS-EXT-IPS:type":         "fixed",
				},
			},
		},
		Links: []any{
			map[string]any{
				"href": "http://104.130.131.164:8774/v2/fcad67a6189847c4aecfa3c81a05783b/servers/9e5476bd-a4ec-4653-93d6-72c93aa682ba",
				"rel":  "self",
			},
			map[string]any{
				"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/servers/9e5476bd-a4ec-4653-93d6-72c93aa682ba",
				"rel":  "bookmark",
			},
		},
		Image: map[string]any{
			"id": "f90f6034-2570-4974-8351-6b49732ef2eb",
			"links": []any{
				map[string]any{
					"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/images/f90f6034-2570-4974-8351-6b49732ef2eb",
					"rel":  "bookmark",
				},
			},
		},
		Flavor: map[string]any{
			"id": "1",
			"links": []any{
				map[string]any{
					"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/flavors/1",
					"rel":  "bookmark",
				},
			},
		},
		ID:              "9e5476bd-a4ec-4653-93d6-72c93aa682ba",
		UserID:          "9349aff8be7545ac9d2f1d00999a23cd",
		Name:            "derp",
		Created:         derpTimeCreated,
		TenantID:        "fcad67a6189847c4aecfa3c81a05783b",
		Metadata:        map[string]string{},
		AttachedVolumes: []servers.AttachedVolume{},
		SecurityGroups: []map[string]any{
			{
				"name": "default",
			},
		},
		Host:               "devstack",
		Hostname:           nil,
		HypervisorHostname: "devstack",
		InstanceName:       "instance-0000001d",
		LaunchIndex:        nil,
		ReservationID:      nil,
		RootDeviceName:     nil,
		Userdata:           nil,
		VmState:            "active",
		PowerState:         servers.RUNNING,
		LaunchedAt:         time.Date(2014, 9, 25, 13, 04, 49, 0, time.UTC),
		TerminatedAt:       time.Time{},
		DiskConfig:         servers.Manual,
		AvailabilityZone:   "nova",
	}

	ConsoleOutput = "abc"

	merpTimeCreated, _ = time.Parse(time.RFC3339, "2014-09-25T13:04:41Z")
	merpTimeUpdated, _ = time.Parse(time.RFC3339, "2014-09-25T13:04:49Z")
	// ServerMerp is a Server struct that should correspond to the second server in ServerListBody.
	ServerMerp = servers.Server{
		Status:  "ACTIVE",
		Updated: merpTimeUpdated,
		HostID:  "29d3c8c896a45aa4c34e52247875d7fefc3d94bbcc9f622b5d204362",
		Addresses: map[string]any{
			"private": []any{
				map[string]any{
					"OS-EXT-IPS-MAC:mac_addr": "fa:16:3e:9e:89:be",
					"version":                 float64(4),
					"addr":                    "10.0.0.31",
					"OS-EXT-IPS:type":         "fixed",
				},
			},
		},
		Links: []any{
			map[string]any{
				"href": "http://104.130.131.164:8774/v2/fcad67a6189847c4aecfa3c81a05783b/servers/9e5476bd-a4ec-4653-93d6-72c93aa682ba",
				"rel":  "self",
			},
			map[string]any{
				"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/servers/9e5476bd-a4ec-4653-93d6-72c93aa682ba",
				"rel":  "bookmark",
			},
		},
		Image: nil,
		Flavor: map[string]any{
			"id": "1",
			"links": []any{
				map[string]any{
					"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/flavors/1",
					"rel":  "bookmark",
				},
			},
		},
		ID:              "9e5476bd-a4ec-4653-93d6-72c93aa682bb",
		UserID:          "9349aff8be7545ac9d2f1d00999a23cd",
		Name:            "merp",
		Created:         merpTimeCreated,
		TenantID:        "fcad67a6189847c4aecfa3c81a05783b",
		Metadata:        map[string]string{},
		AttachedVolumes: []servers.AttachedVolume{},
		SecurityGroups: []map[string]any{
			{
				"name": "default",
			},
		},
		Host:               "devstack",
		Hostname:           nil,
		HypervisorHostname: "devstack",
		InstanceName:       "instance-0000001d",
		LaunchIndex:        nil,
		ReservationID:      nil,
		RootDeviceName:     nil,
		Userdata:           nil,
		VmState:            "active",
		PowerState:         servers.RUNNING,
		LaunchedAt:         time.Date(2014, 9, 25, 13, 04, 49, 0, time.UTC),
		TerminatedAt:       time.Time{},
		DiskConfig:         servers.Manual,
		AvailabilityZone:   "nova",
	}

	faultTimeCreated, _ = time.Parse(time.RFC3339, "2017-11-11T07:58:39Z")
	DerpFault           = servers.Fault{
		Code:    500,
		Created: faultTimeCreated,
		Details: "Stock details for test",
		Message: "Conflict updating instance c2ce4dea-b73f-4d01-8633-2c6032869281. " +
			"Expected: {'task_state': [u'spawning']}. Actual: {'task_state': None}",
	}
)

// HandleServerListSuccessfully sets up the test server to respond to a server detail List request.
func HandleServerListSuccessfully(t *testing.T) {
	testhelper.Mux.HandleFunc("/servers/detail", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprint(w, ServerListBody)
		case "9e5476bd-a4ec-4653-93d6-72c93aa682ba":
			fmt.Fprint(w, `{ "servers": [] }`)
		default:
			t.Fatalf("/servers/detail invoked with unexpected marker=[%s]", marker)
		}
	})
}

// HandleShowConsoleOutputSuccessfully sets up the test server to respond to a os-getConsoleOutput request with success.
func HandleShowConsoleOutputSuccessfully(t *testing.T, response string) {
	th.Mux.HandleFunc("/servers/ef079b0c-e610-4dfb-b1aa-b49f07ac48e5/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `{ "os-getConsoleOutput": {} }`)

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, response)
	})
}
