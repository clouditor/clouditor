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

package openstacktest

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// This code is based on the MockListResponse function from the Gophercloud library.
// Gophercloud is licensed under the Apache License 2.0. See the LICENSE file in the
// Gophercloud repository for the full license: https://github.com/gophercloud/gophercloud
// Source: https://github.com/gophercloud/gophercloud/blob/master/openstack/blockstorage/v3/volumes/testing/fixtures_test.go (2024-12-09)
//
// Changes:
// - 2025-03-12: Rename function name from MockListResponse() to MockStorageListResponse()(@anatheka)

func MockStorageListResponse(t *testing.T) {
	th.Mux.HandleFunc("/volumes/detail", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, `
  {
  "volumes": [
    {
      "volume_type": "lvmdriver-1",
      "created_at": "2015-09-17T03:35:03.000000",
      "bootable": "false",
      "name": "vol-001",
      "os-vol-mig-status-attr:name_id": null,
      "consistencygroup_id": null,
      "source_volid": null,
      "os-volume-replication:driver_data": null,
      "multiattach": false,
      "snapshot_id": null,
      "replication_status": "disabled",
      "os-volume-replication:extended_status": null,
      "encrypted": false,
      "os-vol-host-attr:host": "host-001",
      "availability_zone": "nova",
      "attachments": [{
        "server_id": "83ec2e3b-4321-422b-8706-a84185f52a0a",
        "attachment_id": "05551600-a936-4d4a-ba42-79a037c1-c91a",
        "attached_at": "2016-08-06T14:48:20.000000",
        "host_name": "foobar",
        "volume_id": "d6cacb1a-8b59-4c88-ad90-d70ebb82bb75",
        "device": "/dev/vdc",
        "id": "d6cacb1a-8b59-4c88-ad90-d70ebb82bb75"
      }],
      "id": "289da7f8-6440-407c-9fb4-7db01ec49164",
      "size": 75,
      "user_id": "ff1ce52c03ab433aaba9108c2e3ef541",
      "os-vol-tenant-attr:tenant_id": "304dc00909ac4d0da6c62d816bcb3459",
      "os-vol-mig-status-attr:migstat": null,
      "metadata": {"foo": "bar"},
      "status": "available",
      "description": null
    },
    {
      "volume_type": "lvmdriver-1",
      "created_at": "2015-09-17T03:32:29.000000",
      "bootable": "false",
      "name": "vol-002",
      "os-vol-mig-status-attr:name_id": null,
      "consistencygroup_id": null,
      "source_volid": null,
      "os-volume-replication:driver_data": null,
      "multiattach": false,
      "snapshot_id": null,
      "replication_status": "disabled",
      "os-volume-replication:extended_status": null,
      "encrypted": false,
      "os-vol-host-attr:host": null,
      "availability_zone": "nova",
      "attachments": [],
      "id": "96c3bda7-c82a-4f50-be73-ca7621794835",
      "size": 75,
      "user_id": "ff1ce52c03ab433aaba9108c2e3ef541",
      "os-vol-tenant-attr:tenant_id": "304dc00909ac4d0da6c62d816bcb3459",
      "os-vol-mig-status-attr:migstat": null,
      "metadata": {},
      "status": "available",
      "description": null
    }
  ],
	"volumes_links": [
	{
		"href": "%s/volumes/detail?marker=1",
		"rel": "next"
	}]
}
  `, th.Server.URL)
		case "1":
			fmt.Fprint(w, `{"volumes": []}`)
		default:
			t.Fatalf("Unexpected marker: [%s]", marker)
		}
	})
}

func MockGetResponse(t *testing.T) {
	th.Mux.HandleFunc("/volumes/d32019d3-bc6e-4319-9c1d-6722fc136a22", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `
{
  "volume": {
    "volume_type": "lvmdriver-1",
    "created_at": "2015-09-17T03:32:29.000000",
    "bootable": "false",
    "name": "vol-001",
    "os-vol-mig-status-attr:name_id": null,
    "consistencygroup_id": null,
    "source_volid": null,
    "os-volume-replication:driver_data": null,
    "multiattach": false,
    "snapshot_id": null,
    "replication_status": "disabled",
    "os-volume-replication:extended_status": null,
    "encrypted": false,
    "availability_zone": "nova",
    "attachments": [{
      "server_id": "83ec2e3b-4321-422b-8706-a84185f52a0a",
      "attachment_id": "05551600-a936-4d4a-ba42-79a037c1-c91a",
      "attached_at": "2016-08-06T14:48:20.000000",
      "host_name": "foobar",
      "volume_id": "d6cacb1a-8b59-4c88-ad90-d70ebb82bb75",
      "device": "/dev/vdc",
      "id": "d6cacb1a-8b59-4c88-ad90-d70ebb82bb75"
    }],
    "id": "d32019d3-bc6e-4319-9c1d-6722fc136a22",
    "size": 75,
    "user_id": "ff1ce52c03ab433aaba9108c2e3ef541",
    "os-vol-tenant-attr:tenant_id": "304dc00909ac4d0da6c62d816bcb3459",
    "os-vol-mig-status-attr:migstat": null,
    "metadata": {},
    "status": "available",
    "volume_image_metadata": {
      "container_format": "bare",
      "image_name": "centos"
    },
    "description": null
  }
}
      `)
	})
}

func MockCreateResponse(t *testing.T) {
	th.Mux.HandleFunc("/volumes", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "volume": {
    	"name": "vol-001",
        "size": 75
    }
}
      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

		fmt.Fprint(w, `
{
  "volume": {
    "size": 75,
    "id": "d32019d3-bc6e-4319-9c1d-6722fc136a22",
    "metadata": {},
    "created_at": "2015-09-17T03:32:29.044216",
    "encrypted": false,
    "bootable": "false",
    "availability_zone": "nova",
    "attachments": [],
    "user_id": "ff1ce52c03ab433aaba9108c2e3ef541",
    "status": "creating",
    "description": null,
    "volume_type": "lvmdriver-1",
    "name": "vol-001",
    "replication_status": "disabled",
    "consistencygroup_id": null,
    "source_volid": null,
    "snapshot_id": null,
    "multiattach": false
  }
}
    `)
	})
}

func MockDeleteResponse(t *testing.T) {
	th.Mux.HandleFunc("/volumes/d32019d3-bc6e-4319-9c1d-6722fc136a22", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusAccepted)
	})
}

func MockUpdateResponse(t *testing.T) {
	th.Mux.HandleFunc("/volumes/d32019d3-bc6e-4319-9c1d-6722fc136a22", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `
{
  "volume": {
    "name": "vol-002"
  }
}
        `)
	})
}

func MockCreateVolumeFromBackupResponse(t *testing.T) {
	th.Mux.HandleFunc("/volumes", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "volume": {
        "name": "vol-001",
        "backup_id": "20c792f0-bb03-434f-b653-06ef238e337e"
    }
}
`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

		fmt.Fprint(w, `
{
  "volume": {
    "size": 30,
    "id": "d32019d3-bc6e-4319-9c1d-6722fc136a22",
    "metadata": {},
    "created_at": "2015-09-17T03:32:29.044216",
    "encrypted": false,
    "bootable": "false",
    "availability_zone": "nova",
    "attachments": [],
    "user_id": "ff1ce52c03ab433aaba9108c2e3ef541",
    "status": "creating",
    "description": null,
    "volume_type": "lvmdriver-1",
    "name": "vol-001",
    "replication_status": "disabled",
    "consistencygroup_id": null,
    "source_volid": null,
    "snapshot_id": null,
    "backup_id": "20c792f0-bb03-434f-b653-06ef238e337e",
    "multiattach": false
  }
}`)
	})
}

func MockAttachResponse(t *testing.T) {
	th.Mux.HandleFunc("/volumes/cd281d77-8217-4830-be95-9528227c105c/action",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "POST")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, `
{
    "os-attach":
    {
        "mountpoint": "/mnt",
        "mode": "rw",
        "instance_uuid": "50902f4f-a974-46a0-85e9-7efc5e22dfdd"
    }
}
          `)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)

			fmt.Fprint(w, `{}`)
		})
}

func MockBeginDetachingResponse(t *testing.T) {
	th.Mux.HandleFunc("/volumes/cd281d77-8217-4830-be95-9528227c105c/action",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "POST")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, `
{
    "os-begin_detaching": {}
}
          `)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)

			fmt.Fprint(w, `{}`)
		})
}

func MockDetachResponse(t *testing.T) {
	th.Mux.HandleFunc("/volumes/cd281d77-8217-4830-be95-9528227c105c/action",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "POST")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, `
{
    "os-detach": {}
}
          `)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)

			fmt.Fprint(w, `{}`)
		})
}

func MockUploadImageResponse(t *testing.T) {
	th.Mux.HandleFunc("/volumes/cd281d77-8217-4830-be95-9528227c105c/action",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "POST")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, `
{
    "os-volume_upload_image": {
        "container_format": "bare",
        "force": true,
        "image_name": "test",
        "disk_format": "raw"
    }
}
          `)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)

			fmt.Fprint(w, `
{
    "os-volume_upload_image": {
        "container_format": "bare",
        "display_description": null,
        "id": "cd281d77-8217-4830-be95-9528227c105c",
        "image_id": "ecb92d98-de08-45db-8235-bbafe317269c",
        "image_name": "test",
        "disk_format": "raw",
        "size": 5,
        "status": "uploading",
        "updated_at": "2017-07-17T09:29:22.000000",
        "volume_type": {
            "created_at": "2016-05-04T08:54:14.000000",
            "deleted": false,
            "deleted_at": null,
            "description": null,
            "extra_specs": {
                "volume_backend_name": "basic.ru-2a"
            },
            "id": "b7133444-62f6-4433-8da3-70ac332229b7",
            "is_public": true,
            "name": "basic.ru-2a",
            "updated_at": "2016-05-04T09:15:33.000000"
        }
    }
}
          `)
		})
}

func MockReserveResponse(t *testing.T) {
	th.Mux.HandleFunc("/volumes/cd281d77-8217-4830-be95-9528227c105c/action",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "POST")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, `
{
    "os-reserve": {}
}
          `)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)

			fmt.Fprint(w, `{}`)
		})
}

func MockUnreserveResponse(t *testing.T) {
	th.Mux.HandleFunc("/volumes/cd281d77-8217-4830-be95-9528227c105c/action",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "POST")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, `
{
    "os-unreserve": {}
}
          `)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)

			fmt.Fprint(w, `{}`)
		})
}

func MockInitializeConnectionResponse(t *testing.T) {
	th.Mux.HandleFunc("/volumes/cd281d77-8217-4830-be95-9528227c105c/action",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "POST")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, `
{
    "os-initialize_connection":
    {
        "connector":
        {
        "ip":"127.0.0.1",
        "host":"stack",
        "initiator":"iqn.1994-05.com.redhat:17cf566367d2",
        "multipath": false,
        "platform": "x86_64",
        "os_type": "linux2"
        }
    }
}
          `)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)

			fmt.Fprint(w, `{
"connection_info": {
    "data": {
      "target_portals": [
        "172.31.17.48:3260"
      ],
      "auth_method": "CHAP",
      "auth_username": "5MLtcsTEmNN5jFVcT6ui",
      "access_mode": "rw",
      "target_lun": 0,
      "volume_id": "cd281d77-8217-4830-be95-9528227c105c",
      "target_luns": [
        0
      ],
      "target_iqns": [
        "iqn.2010-10.org.openstack:volume-cd281d77-8217-4830-be95-9528227c105c"
      ],
      "auth_password": "x854ZY5Re3aCkdNL",
      "target_discovered": false,
      "encrypted": false,
      "qos_specs": null,
      "target_iqn": "iqn.2010-10.org.openstack:volume-cd281d77-8217-4830-be95-9528227c105c",
      "target_portal": "172.31.17.48:3260"
    },
    "driver_volume_type": "iscsi"
  }
            }`)
		})
}

func MockTerminateConnectionResponse(t *testing.T) {
	th.Mux.HandleFunc("/volumes/cd281d77-8217-4830-be95-9528227c105c/action",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "POST")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, `
{
    "os-terminate_connection":
    {
        "connector":
        {
        "ip":"127.0.0.1",
        "host":"stack",
        "initiator":"iqn.1994-05.com.redhat:17cf566367d2",
        "multipath": true,
        "platform": "x86_64",
        "os_type": "linux2"
        }
    }
}
          `)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)

			fmt.Fprint(w, `{}`)
		})
}

func MockExtendSizeResponse(t *testing.T) {
	th.Mux.HandleFunc("/volumes/cd281d77-8217-4830-be95-9528227c105c/action",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "POST")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, `
{
    "os-extend":
    {
        "new_size": 3
    }
}
          `)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)

			fmt.Fprint(w, `{}`)
		})
}

func MockForceDeleteResponse(t *testing.T) {
	th.Mux.HandleFunc("/volumes/d32019d3-bc6e-4319-9c1d-6722fc136a22/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestBody(t, r, `{"os-force_delete":""}`)
		w.WriteHeader(http.StatusAccepted)
	})
}

func MockSetImageMetadataResponse(t *testing.T) {
	th.Mux.HandleFunc("/volumes/cd281d77-8217-4830-be95-9528227c105c/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
	"os-set_image_metadata": {
		"metadata": {
			"label": "test"
		}
	}
}
		`)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `{}`)
	})
}

func MockSetBootableResponse(t *testing.T) {
	th.Mux.HandleFunc("/volumes/cd281d77-8217-4830-be95-9528227c105c/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
	"os-set_bootable": {
		"bootable": true
	}
}
		`)
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Content-Length", "0")
		w.WriteHeader(http.StatusOK)
	})
}

func MockReImageResponse(t *testing.T) {
	th.Mux.HandleFunc("/volumes/cd281d77-8217-4830-be95-9528227c105c/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
	"os-reimage": {
		"image_id": "71543ced-a8af-45b6-a5c4-a46282108a90",
		"reimage_reserved": false
	}
}
		`)
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Content-Length", "0")
		w.WriteHeader(http.StatusAccepted)
	})
}

func MockChangeTypeResponse(t *testing.T) {
	th.Mux.HandleFunc("/volumes/cd281d77-8217-4830-be95-9528227c105c/action",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "POST")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, `
{
    "os-retype":
    {
		"new_type": "ssd",
		"migration_policy": "on-demand"
    }
}
          `)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)

			fmt.Fprint(w, `{}`)
		})
}

func MockResetStatusResponse(t *testing.T) {
	th.Mux.HandleFunc("/volumes/cd281d77-8217-4830-be95-9528227c105c/action",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "POST")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestJSONRequest(t, r, `
{
    "os-reset_status":
    {
		"status": "error",
		"attach_status": "detached",
		"migration_status": "migrating"
    }
}
          `)

			w.WriteHeader(http.StatusAccepted)
		})
}
