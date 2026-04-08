// Copyright 2026 Fraunhofer AISEC
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

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// The code in this file is based on the Gophercloud library.
// Gophercloud is licensed under the Apache License 2.0. See the LICENSE file in the
// Gophercloud repository for the full license: https://github.com/gophercloud/gophercloud
// Source: https://raw.githubusercontent.com/gophercloud/gophercloud/refs/heads/main/openstack/networking/v2/ports/testing/fixtures.go (2026-03-05)
//
// Changes:
// - 2026-03-05: Rename const names, otherwise we have conflicts (anatheka)
// - 2026-03-05: Add HandleFuncs from https://github.com/gophercloud/gophercloud/blob/main/openstack/networking/v2/ports/testing/requests_test.go (anatheka)

const ListPortResponse = `
{
    "ports": [
        {
            "status": "ACTIVE",
            "binding:host_id": "devstack",
            "name": "",
            "admin_state_up": true,
            "network_id": "70c1db1f-b701-45bd-96e0-a313ee3430b3",
            "tenant_id": "",
            "device_owner": "network:router_gateway",
            "mac_address": "fa:16:3e:58:42:ed",
            "binding:vnic_type": "normal",
            "fixed_ips": [
                {
                    "subnet_id": "008ba151-0b8c-4a67-98b5-0d2b87666062",
                    "ip_address": "172.24.4.2"
                }
            ],
            "id": "d80b1a3b-4fc1-49f3-952e-1e2ab7081d8b",
            "security_groups": [],
            "dns_name": "test-port",
            "dns_assignment": [
              {
                "hostname": "test-port",
                "ip_address": "172.24.4.2",
                "fqdn": "test-port.openstack.local."
              }
            ],
            "device_id": "9ae135f4-b6e0-4dad-9e91-3c223e385824",
            "port_security_enabled": false,
            "created_at": "2019-06-30T04:15:37",
            "updated_at": "2019-06-30T05:18:49"
        }
    ]
}
`

const GetPortResponse = `
{
    "port": {
        "status": "ACTIVE",
        "binding:host_id": "devstack",
        "name": "",
        "allowed_address_pairs": [],
        "admin_state_up": true,
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "7e02058126cc4950b75f9970368ba177",
        "extra_dhcp_opts": [],
        "binding:vif_details": {
            "port_filter": true,
            "ovs_hybrid_plug": true
         },
        "binding:vif_type": "ovs",
        "device_owner": "network:router_interface",
        "port_security_enabled": false,
        "mac_address": "fa:16:3e:23:fd:d7",
        "binding:profile": {},
        "binding:vnic_type": "normal",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.1"
            }
        ],
        "id": "46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2",
        "security_groups": [],
        "dns_name": "test-port",
        "dns_assignment": [
          {
            "hostname": "test-port",
            "ip_address": "172.24.4.2",
            "fqdn": "test-port.openstack.local."
          }
        ],
        "device_id": "5e3898d7-11be-483e-9732-b2f5eccd2b2e",
        "created_at": "2019-06-30T04:15:37Z",
        "updated_at": "2019-06-30T05:18:49Z"
    }
}
`

const CreatePortRequest = `
{
    "port": {
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "name": "private-port",
        "admin_state_up": true,
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.2"
            }
        ],
        "security_groups": ["foo"],
        "allowed_address_pairs": [
          {
            "ip_address": "10.0.0.4",
            "mac_address": "fa:16:3e:c9:cb:f0"
          }
        ]
    }
}
`

const CreatePortResponse = `
{
    "port": {
        "status": "DOWN",
        "name": "private-port",
        "admin_state_up": true,
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "d6700c0c9ffa4f1cb322cd4a1f3906fa",
        "device_owner": "",
        "mac_address": "fa:16:3e:c9:cb:f0",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.2"
            }
        ],
        "id": "65c0ee9f-d634-4522-8954-51021b570b0d",
        "security_groups": [
            "f0ac4394-7e4a-4409-9701-ba8be283dbc3"
        ],
        "allowed_address_pairs": [
          {
            "ip_address": "10.0.0.4",
            "mac_address": "fa:16:3e:c9:cb:f0"
          }
        ],
        "device_id": ""
    }
}
`

const CreateOmitSecurityGroupsRequest = `
{
    "port": {
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "name": "private-port",
        "admin_state_up": true,
        "fixed_ips": [
          {
            "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
            "ip_address": "10.0.0.2"
          }
        ],
        "allowed_address_pairs": [
          {
            "ip_address": "10.0.0.4",
            "mac_address": "fa:16:3e:c9:cb:f0"
          }
        ]
    }
}
`

const CreateWithNoSecurityGroupsRequest = `
{
    "port": {
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "name": "private-port",
        "admin_state_up": true,
        "fixed_ips": [
          {
            "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
            "ip_address": "10.0.0.2"
          }
        ],
        "security_groups": [],
        "allowed_address_pairs": [
          {
            "ip_address": "10.0.0.4",
            "mac_address": "fa:16:3e:c9:cb:f0"
          }
        ]
    }
}
`

const CreateWithNoSecurityGroupsResponse = `
{
    "port": {
        "status": "DOWN",
        "name": "private-port",
        "admin_state_up": true,
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "d6700c0c9ffa4f1cb322cd4a1f3906fa",
        "device_owner": "",
        "mac_address": "fa:16:3e:c9:cb:f0",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.2"
            }
        ],
        "id": "65c0ee9f-d634-4522-8954-51021b570b0d",
        "allowed_address_pairs": [
          {
            "ip_address": "10.0.0.4",
            "mac_address": "fa:16:3e:c9:cb:f0"
          }
        ],
        "device_id": ""
    }
}
`

const CreateOmitSecurityGroupsResponse = `
{
    "port": {
        "status": "DOWN",
        "name": "private-port",
        "admin_state_up": true,
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "d6700c0c9ffa4f1cb322cd4a1f3906fa",
        "device_owner": "",
        "mac_address": "fa:16:3e:c9:cb:f0",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.2"
            }
        ],
        "id": "65c0ee9f-d634-4522-8954-51021b570b0d",
        "security_groups": [
            "f0ac4394-7e4a-4409-9701-ba8be283dbc3"
        ],
        "allowed_address_pairs": [
          {
            "ip_address": "10.0.0.4",
            "mac_address": "fa:16:3e:c9:cb:f0"
          }
        ],
        "device_id": ""
    }
}
`

const CreatePropagateUplinkStatusRequest = `
{
    "port": {
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "name": "private-port",
        "admin_state_up": true,
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.2"
            }
        ],
        "propagate_uplink_status": true
    }
}
`

const CreatePropagateUplinkStatusResponse = `
{
    "port": {
        "status": "DOWN",
        "name": "private-port",
        "admin_state_up": true,
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "d6700c0c9ffa4f1cb322cd4a1f3906fa",
        "device_owner": "",
        "mac_address": "fa:16:3e:c9:cb:f0",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.2"
            }
        ],
        "id": "65c0ee9f-d634-4522-8954-51021b570b0d",
        "propagate_uplink_status": true,
        "device_id": ""
    }
}
`

const CreateValueSpecRequest = `
{
    "port": {
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "name": "private-port",
        "admin_state_up": true,
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.2"
            }
        ],
        "security_groups": ["foo"],
        "allowed_address_pairs": [
          {
            "ip_address": "10.0.0.4",
            "mac_address": "fa:16:3e:c9:cb:f0"
          }
        ],
        "test": "value"
    }
}
`

const CreateValueSpecResponse = `
{
    "port": {
        "status": "DOWN",
        "name": "private-port",
        "admin_state_up": true,
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "d6700c0c9ffa4f1cb322cd4a1f3906fa",
        "device_owner": "",
        "mac_address": "fa:16:3e:c9:cb:f0",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.2"
            }
        ],
        "id": "65c0ee9f-d634-4522-8954-51021b570b0d",
        "security_groups": [
            "f0ac4394-7e4a-4409-9701-ba8be283dbc3"
        ],
        "allowed_address_pairs": [
          {
            "ip_address": "10.0.0.4",
            "mac_address": "fa:16:3e:c9:cb:f0"
          }
        ],
        "device_id": ""
    }
}
`

const CreatePortPortSecurityRequest = `
{
    "port": {
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "name": "private-port",
        "admin_state_up": true,
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.2"
            }
        ],
        "security_groups": ["foo"],
        "allowed_address_pairs": [
          {
            "ip_address": "10.0.0.4",
            "mac_address": "fa:16:3e:c9:cb:f0"
          }
        ],
        "port_security_enabled": false
    }
}
`

const CreatePortPortSecurityResponse = `
{
    "port": {
        "status": "DOWN",
        "name": "private-port",
        "admin_state_up": true,
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "d6700c0c9ffa4f1cb322cd4a1f3906fa",
        "device_owner": "",
        "mac_address": "fa:16:3e:c9:cb:f0",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.2"
            }
        ],
        "id": "65c0ee9f-d634-4522-8954-51021b570b0d",
        "security_groups": [
            "f0ac4394-7e4a-4409-9701-ba8be283dbc3"
        ],
        "allowed_address_pairs": [
          {
            "ip_address": "10.0.0.4",
            "mac_address": "fa:16:3e:c9:cb:f0"
          }
        ],
        "device_id": "",
        "port_security_enabled": false
    }
}
`

const UpdatePortRequest = `
{
    "port": {
        "name": "new_port_name",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.3"
            }
        ],
        "allowed_address_pairs": [
          {
            "ip_address": "10.0.0.4",
            "mac_address": "fa:16:3e:c9:cb:f0"
          }
        ],
        "security_groups": [
            "f0ac4394-7e4a-4409-9701-ba8be283dbc3"
        ],
        "mac_address": "fa:16:3e:c9:cb:f4"
    }
}
`

const UpdatePortResponse = `
{
    "port": {
        "status": "DOWN",
        "name": "new_port_name",
        "admin_state_up": true,
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "d6700c0c9ffa4f1cb322cd4a1f3906fa",
        "device_owner": "",
        "mac_address": "fa:16:3e:c9:cb:f4",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.3"
            }
        ],
        "allowed_address_pairs": [
          {
            "ip_address": "10.0.0.4",
            "mac_address": "fa:16:3e:c9:cb:f0"
          }
        ],
        "id": "65c0ee9f-d634-4522-8954-51021b570b0d",
        "security_groups": [
            "f0ac4394-7e4a-4409-9701-ba8be283dbc3"
        ],
        "device_id": ""
    }
}
`

const UpdateOmitSecurityGroupsRequest = `
{
  "port": {
      "name": "new_port_name",
      "fixed_ips": [
          {
            "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
            "ip_address": "10.0.0.3"
          }
        ],
        "allowed_address_pairs": [
          {
            "ip_address": "10.0.0.4",
            "mac_address": "fa:16:3e:c9:cb:f0"
          }
        ]
    }
}
`

const UpdateOmitSecurityGroupsResponse = `
{
    "port": {
        "status": "DOWN",
        "name": "new_port_name",
        "admin_state_up": true,
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "d6700c0c9ffa4f1cb322cd4a1f3906fa",
        "device_owner": "",
        "mac_address": "fa:16:3e:c9:cb:f0",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.3"
            }
        ],
        "allowed_address_pairs": [
          {
            "ip_address": "10.0.0.4",
            "mac_address": "fa:16:3e:c9:cb:f0"
          }
        ],
        "id": "65c0ee9f-d634-4522-8954-51021b570b0d",
        "security_groups": [
            "f0ac4394-7e4a-4409-9701-ba8be283dbc3"
        ],
        "device_id": ""
    }
}
`

const UpdatePropagateUplinkStatusRequest = `
{
    "port": {
        "propagate_uplink_status": true
    }
}
`

const UpdatePropagateUplinkStatusResponse = `
{
    "port": {
        "status": "DOWN",
        "name": "new_port_name",
        "admin_state_up": true,
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "d6700c0c9ffa4f1cb322cd4a1f3906fa",
        "device_owner": "",
        "mac_address": "fa:16:3e:c9:cb:f0",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.3"
            }
        ],
        "allowed_address_pairs": [
          {
            "ip_address": "10.0.0.4",
            "mac_address": "fa:16:3e:c9:cb:f0"
          }
        ],
        "id": "65c0ee9f-d634-4522-8954-51021b570b0d",
        "security_groups": [
            "f0ac4394-7e4a-4409-9701-ba8be283dbc3"
        ],
        "propagate_uplink_status": true,
        "device_id": ""
    }
}
`

const UpdateValueSpecsRequest = `
{
    "port": {
        "test": "update"
    }
}
`

const UpdateValueSpecsResponse = `
{
    "port": {
        "status": "DOWN",
        "name": "new_port_name",
        "admin_state_up": true,
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "d6700c0c9ffa4f1cb322cd4a1f3906fa",
        "device_owner": "",
        "mac_address": "fa:16:3e:c9:cb:f0",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.3"
            }
        ],
        "allowed_address_pairs": [
          {
            "ip_address": "10.0.0.4",
            "mac_address": "fa:16:3e:c9:cb:f0"
          }
        ],
        "id": "65c0ee9f-d634-4522-8954-51021b570b0d",
        "security_groups": [
            "f0ac4394-7e4a-4409-9701-ba8be283dbc3"
        ],
        "device_id": ""
    }
}
`

const UpdatePortPortSecurityRequest = `
{
    "port": {
        "port_security_enabled": false
    }
}
`

const UpdatePortPortSecurityResponse = `
{
    "port": {
        "status": "DOWN",
        "name": "private-port",
        "admin_state_up": true,
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "d6700c0c9ffa4f1cb322cd4a1f3906fa",
        "device_owner": "",
        "mac_address": "fa:16:3e:c9:cb:f0",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.2"
            }
        ],
        "id": "65c0ee9f-d634-4522-8954-51021b570b0d",
        "security_groups": [
            "f0ac4394-7e4a-4409-9701-ba8be283dbc3"
        ],
        "allowed_address_pairs": [
          {
            "ip_address": "10.0.0.4",
            "mac_address": "fa:16:3e:c9:cb:f0"
          }
        ],
        "device_id": "",
        "port_security_enabled": false
    }
}
`

const RemoveSecurityGroupRequest = `
{
    "port": {
      "name": "new_port_name",
      "fixed_ips": [
        {
          "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
          "ip_address": "10.0.0.3"
        }
      ],
      "allowed_address_pairs": [
        {
          "ip_address": "10.0.0.4",
          "mac_address": "fa:16:3e:c9:cb:f0"
        }
      ],
      "security_groups": []
    }
}
`

const RemoveSecurityGroupResponse = `
{
    "port": {
        "status": "DOWN",
        "name": "new_port_name",
        "admin_state_up": true,
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "d6700c0c9ffa4f1cb322cd4a1f3906fa",
        "device_owner": "",
        "mac_address": "fa:16:3e:c9:cb:f0",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.3"
            }
        ],
        "allowed_address_pairs": [
          {
            "ip_address": "10.0.0.4",
            "mac_address": "fa:16:3e:c9:cb:f0"
          }
        ],
        "id": "65c0ee9f-d634-4522-8954-51021b570b0d",
        "device_id": ""
    }
}
`

const RemoveAllowedAddressPairsRequest = `
{
    "port": {
      "name": "new_port_name",
      "fixed_ips": [
        {
          "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
          "ip_address": "10.0.0.3"
        }
      ],
      "allowed_address_pairs": [],
      "security_groups": [
        "f0ac4394-7e4a-4409-9701-ba8be283dbc3"
      ]
    }
}
`

const RemoveAllowedAddressPairsResponse = `
{
    "port": {
        "status": "DOWN",
        "name": "new_port_name",
        "admin_state_up": true,
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "d6700c0c9ffa4f1cb322cd4a1f3906fa",
        "device_owner": "",
        "mac_address": "fa:16:3e:c9:cb:f0",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.3"
            }
        ],
        "id": "65c0ee9f-d634-4522-8954-51021b570b0d",
        "security_groups": [
            "f0ac4394-7e4a-4409-9701-ba8be283dbc3"
        ],
        "device_id": ""
    }
}
`

const DontUpdateAllowedAddressPairsRequest = `
{
    "port": {
        "name": "new_port_name",
        "fixed_ips": [
          {
            "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
            "ip_address": "10.0.0.3"
          }
        ],
        "security_groups": [
            "f0ac4394-7e4a-4409-9701-ba8be283dbc3"
        ]
    }
}
`

const DontUpdateAllowedAddressPairsResponse = `
{
    "port": {
        "status": "DOWN",
        "name": "new_port_name",
        "admin_state_up": true,
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "d6700c0c9ffa4f1cb322cd4a1f3906fa",
        "device_owner": "",
        "mac_address": "fa:16:3e:c9:cb:f0",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.3"
            }
        ],
        "allowed_address_pairs": [
          {
            "ip_address": "10.0.0.4",
            "mac_address": "fa:16:3e:c9:cb:f0"
          }
        ],
        "id": "65c0ee9f-d634-4522-8954-51021b570b0d",
        "security_groups": [
            "f0ac4394-7e4a-4409-9701-ba8be283dbc3"
        ],
        "device_id": ""
    }
}
`

// GetWithExtraDHCPOptsResponse represents a raw port response with extra
// DHCP options.
const GetWithExtraDHCPOptsResponse = `
{
    "port": {
        "status": "ACTIVE",
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "d6700c0c9ffa4f1cb322cd4a1f3906fa",
        "extra_dhcp_opts": [
            {
                "opt_name": "option1",
                "opt_value": "value1",
                "ip_version": 4
            },
            {
                "opt_name": "option2",
                "opt_value": "value2",
                "ip_version": 4
            }
        ],
        "admin_state_up": true,
        "name": "port-with-extra-dhcp-opts",
        "device_owner": "",
        "mac_address": "fa:16:3e:c9:cb:f0",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.4"
            }
        ],
        "id": "65c0ee9f-d634-4522-8954-51021b570b0d",
        "device_id": ""
    }
}
`

// CreateWithExtraDHCPOptsRequest represents a raw port creation request
// with extra DHCP options.
const CreateWithExtraDHCPOptsRequest = `
{
    "port": {
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "name": "port-with-extra-dhcp-opts",
        "admin_state_up": true,
        "fixed_ips": [
            {
                "ip_address": "10.0.0.2",
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2"
            }
        ],
        "extra_dhcp_opts": [
            {
                "opt_name": "option1",
                "opt_value": "value1"
            }
        ]
    }
}
`

// CreateWithExtraDHCPOptsResponse represents a raw port creation response
// with extra DHCP options.
const CreateWithExtraDHCPOptsResponse = `
{
    "port": {
        "status": "DOWN",
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "d6700c0c9ffa4f1cb322cd4a1f3906fa",
        "extra_dhcp_opts": [
            {
                "opt_name": "option1",
                "opt_value": "value1",
                "ip_version": 4
            }
        ],
        "admin_state_up": true,
        "name": "port-with-extra-dhcp-opts",
        "device_owner": "",
        "mac_address": "fa:16:3e:c9:cb:f0",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.2"
            }
        ],
        "id": "65c0ee9f-d634-4522-8954-51021b570b0d",
        "device_id": ""
    }
}
`

// UpdateWithExtraDHCPOptsRequest represents a raw port update request with
// extra DHCP options.
const UpdateWithExtraDHCPOptsRequest = `
{
    "port": {
        "name": "updated-port-with-dhcp-opts",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.3"
            }
        ],
        "extra_dhcp_opts": [
            {
                "opt_name": "option1",
                "opt_value": null
            },
            {
                "opt_name": "option2",
                "opt_value": "value2"
            }
        ]
    }
}
`

// UpdateWithExtraDHCPOptsResponse represents a raw port update response with
// extra DHCP options.
const UpdateWithExtraDHCPOptsResponse = `
{
    "port": {
        "status": "DOWN",
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "d6700c0c9ffa4f1cb322cd4a1f3906fa",
        "extra_dhcp_opts": [
            {
                "opt_name": "option2",
                "opt_value": "value2",
                "ip_version": 4
            }
        ],
        "admin_state_up": true,
        "name": "updated-port-with-dhcp-opts",
        "device_owner": "",
        "mac_address": "fa:16:3e:c9:cb:f0",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.3"
            }
        ],
        "id": "65c0ee9f-d634-4522-8954-51021b570b0d",
        "device_id": ""
    }
}
`

func HandlePortsListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v2.0/ports", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ListPortResponse)
	})
}
