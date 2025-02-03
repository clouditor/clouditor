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

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/attachinterfaces"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// The code in this file is based on the Gophercloud library.
// Gophercloud is licensed under the Apache License 2.0. See the LICENSE file in the
// Gophercloud repository for the full license: https://github.com/gophercloud/gophercloud
// Source: https://github.com/gophercloud/gophercloud/blob/master/openstack/compute/v2/attachinterfaces/testing/fixtures_test.go (2024-12-09)
//
// Changes:
// - 2024-12-10: Extended function HandleInterfaceListSuccessfully() for all 3 virtual machines given in /internal/testuitl/servicetes/openstacktest/openstack_compute.go (@anatheka)
// - 2025-01-07:  (@anatheka)

// ListInterfacesExpected represents an expected response from a ListInterfaces request.
var ListInterfacesExpected = []attachinterfaces.Interface{
	{
		PortState: "ACTIVE",
		FixedIPs: []attachinterfaces.FixedIP{
			{
				SubnetID:  "d7906db4-a566-4546-b1f4-5c7fa70f0bf3",
				IPAddress: "10.0.0.7",
			},
			{
				SubnetID:  "45906d64-a548-4276-h1f8-kcffa80fjbnl",
				IPAddress: "10.0.0.8",
			},
		},
		PortID:  "0dde1598-b374-474e-986f-5b8dd1df1d4e",
		NetID:   "8a5fe506-7e9f-4091-899b-96336909d93c",
		MACAddr: "fa:16:3e:38:2d:80",
	},
}

// GetInterfaceExpected represents an expected response from a GetInterface request.
var GetInterfaceExpected = attachinterfaces.Interface{
	PortState: "ACTIVE",
	FixedIPs: []attachinterfaces.FixedIP{
		{
			SubnetID:  "d7906db4-a566-4546-b1f4-5c7fa70f0bf3",
			IPAddress: "10.0.0.7",
		},
		{
			SubnetID:  "45906d64-a548-4276-h1f8-kcffa80fjbnl",
			IPAddress: "10.0.0.8",
		},
	},
	PortID:  "0dde1598-b374-474e-986f-5b8dd1df1d4e",
	NetID:   "8a5fe506-7e9f-4091-899b-96336909d93c",
	MACAddr: "fa:16:3e:38:2d:80",
}

// CreateInterfacesExpected represents an expected response from a CreateInterface request.
var CreateInterfacesExpected = attachinterfaces.Interface{
	PortState: "ACTIVE",
	FixedIPs: []attachinterfaces.FixedIP{
		{
			SubnetID:  "d7906db4-a566-4546-b1f4-5c7fa70f0bf3",
			IPAddress: "10.0.0.7",
		},
	},
	PortID:  "0dde1598-b374-474e-986f-5b8dd1df1d4e",
	NetID:   "8a5fe506-7e9f-4091-899b-96336909d93c",
	MACAddr: "fa:16:3e:38:2d:80",
}

// HandleInterfaceListSuccessfully sets up the test server to respond to a ListInterfaces request.
func HandleInterfaceListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/servers/ef079b0c-e610-4dfb-b1aa-b49f07ac48e5/os-interface", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"interfaceAttachments": [
				{
					"port_state":"ACTIVE",
					"fixed_ips": [
						{
							"subnet_id": "d7906db4-a566-4546-b1f4-5c7fa70f0bf3",
							"ip_address": "10.0.0.7"
						},
						{
							"subnet_id": "45906d64-a548-4276-h1f8-kcffa80fjbnl",
							"ip_address": "10.0.0.8"
						}
					],
					"port_id": "0dde1598-b374-474e-986f-5b8dd1df1d4e",
					"net_id": "8a5fe506-7e9f-4091-899b-96336909d93c",
					"mac_addr": "fa:16:3e:38:2d:80"
				}
			]
		}`)
	})

	th.Mux.HandleFunc("/servers/9e5476bd-a4ec-4653-93d6-72c93aa682ba/os-interface", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"interfaceAttachments": [
				{
					"port_state":"ACTIVE",
					"fixed_ips": [
						{
							"subnet_id": "d7906db4-a566-4546-b1f4-5c7fa70f0bf3",
							"ip_address": "10.0.0.7"
						},
						{
							"subnet_id": "45906d64-a548-4276-h1f8-kcffa80fjbnl",
							"ip_address": "10.0.0.8"
						}
					],
					"port_id": "0dde1598-b374-474e-986f-5b8dd1df1d4e",
					"net_id": "8a5fe506-7e9f-4091-899b-96336909d93c",
					"mac_addr": "fa:16:3e:38:2d:80"
				}
			]
		}`)
	})

	th.Mux.HandleFunc("/servers/9e5476bd-a4ec-4653-93d6-72c93aa682bb/os-interface", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"interfaceAttachments": [
				{
					"port_state":"ACTIVE",
					"fixed_ips": [
						{
							"subnet_id": "d7906db4-a566-4546-b1f4-5c7fa70f0bf3",
							"ip_address": "10.0.0.7"
						},
						{
							"subnet_id": "45906d64-a548-4276-h1f8-kcffa80fjbnl",
							"ip_address": "10.0.0.8"
						}
					],
					"port_id": "0dde1598-b374-474e-986f-5b8dd1df1d4e",
					"net_id": "8a5fe506-7e9f-4091-899b-96336909d93c",
					"mac_addr": "fa:16:3e:38:2d:80"
				}
			]
		}`)
	})
}

// HandleInterfaceGetSuccessfully sets up the test server to respond to a GetInterface request.
func HandleInterfaceGetSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/servers/ef079b0c-e610-4dfb-b1aa-b49f07ac48e5/os-interface/0dde1598-b374-474e-986f-5b8dd1df1d4e", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"interfaceAttachment":
				{
					"port_state":"ACTIVE",
					"fixed_ips": [
						{
							"subnet_id": "d7906db4-a566-4546-b1f4-5c7fa70f0bf3",
							"ip_address": "10.0.0.7"
						},
						{
							"subnet_id": "45906d64-a548-4276-h1f8-kcffa80fjbnl",
							"ip_address": "10.0.0.8"
						}
					],
					"port_id": "0dde1598-b374-474e-986f-5b8dd1df1d4e",
					"net_id": "8a5fe506-7e9f-4091-899b-96336909d93c",
					"mac_addr": "fa:16:3e:38:2d:80"
				}
		}`)
	})
}

// HandleInterfaceCreateSuccessfully sets up the test server to respond to a CreateInterface request.
func HandleInterfaceCreateSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/servers/ef079b0c-e610-4dfb-b1aa-b49f07ac48e5/os-interface", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `{
			  "interfaceAttachment": {
				"net_id": "8a5fe506-7e9f-4091-899b-96336909d93c"
			  }
		}`)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"interfaceAttachment":
				{
					"port_state":"ACTIVE",
					"fixed_ips": [
						{
							"subnet_id": "d7906db4-a566-4546-b1f4-5c7fa70f0bf3",
							"ip_address": "10.0.0.7"
						}
					],
					"port_id": "0dde1598-b374-474e-986f-5b8dd1df1d4e",
					"net_id": "8a5fe506-7e9f-4091-899b-96336909d93c",
					"mac_addr": "fa:16:3e:38:2d:80"
				}
		}`)
	})
}

// HandleInterfaceDeleteSuccessfully sets up the test server to respond to a DeleteInterface request.
func HandleInterfaceDeleteSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/servers/ef079b0c-e610-4dfb-b1aa-b49f07ac48e5/os-interface/0dde1598-b374-474e-986f-5b8dd1df1d4e", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusAccepted)
	})
}
