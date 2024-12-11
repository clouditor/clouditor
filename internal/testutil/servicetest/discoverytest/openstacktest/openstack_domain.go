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

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/domains"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// The code in this file is based on the Gophercloud library.
// Gophercloud is licensed under the Apache License 2.0. See the LICENSE file in the
// Gophercloud repository for the full license: https://github.com/gophercloud/gophercloud
// Source: https://github.com/gophercloud/gophercloud/blob/v0.25.0/openstack/identity/v3/domains/testing/fixtures.go(2024-12-10)
//
// Changes:
// - 2024-12-10: Changed variable from ListOutput to ListDomainOutput (@anatheka)
// - 2024-12-10: Changed variable from GetOutput to GetDomainOutput (@anatheka)
// - 2024-12-10: Changed variable from CreateRequest to ListDomainOutput (@anatheka)
// - 2024-12-10: Changed variable from CreateRequest to ListDomainOutput (@anatheka)
// - 2024-12-10: Changed variable from UpdateRequest to UpdateDomainRequest (@anatheka)
// - 2024-12-10: Changed variable from UpdateOutput to UpdateDomainOutput (@anatheka)

// ListDomainOutput provides a single page of Domain results.
const ListDomainOutput = `
{
    "links": {
        "next": null,
        "previous": null,
        "self": "http://example.com/identity/v3/domains"
    },
    "domains": [
        {
            "enabled": true,
            "id": "2844b2a08be147a08ef58317d6471f1f",
            "links": {
                "self": "http://example.com/identity/v3/domains/2844b2a08be147a08ef58317d6471f1f"
            },
            "name": "domain one",
            "description": "some description"
        },
        {
            "enabled": true,
            "id": "9fe1d3",
            "links": {
                "self": "https://example.com/identity/v3/domains/9fe1d3"
            },
            "name": "domain two"
        }
    ]
}
`

// GetDomainOutput provides a Get result.
const GetDomainOutput = `
{
    "domain": {
        "enabled": true,
        "id": "9fe1d3",
        "links": {
            "self": "https://example.com/identity/v3/domains/9fe1d3"
        },
        "name": "domain two"
    }
}
`

// CreateDomainRequest provides the input to a Create request.
const CreateDomainRequest = `
{
    "domain": {
        "name": "domain two"
    }
}
`

// UpdateDomainRequest provides the input to as Update request.
const UpdateDomainRequest = `
{
    "domain": {
        "description": "Staging Domain"
    }
}
`

// UpdateDomainOutput provides an update result.
const UpdateDomainOutput = `
{
    "domain": {
		"enabled": true,
        "id": "9fe1d3",
        "links": {
            "self": "https://example.com/identity/v3/domains/9fe1d3"
        },
        "name": "domain two",
        "description": "Staging Domain"
    }
}
`

// FirstDomain is the first domain in the List request.
var FirstDomain = domains.Domain{
	Enabled: true,
	ID:      "2844b2a08be147a08ef58317d6471f1f",
	Links: map[string]interface{}{
		"self": "http://example.com/identity/v3/domains/2844b2a08be147a08ef58317d6471f1f",
	},
	Name:        "domain one",
	Description: "some description",
}

// SecondDomain is the second domain in the List request.
var SecondDomain = domains.Domain{
	Enabled: true,
	ID:      "9fe1d3",
	Links: map[string]interface{}{
		"self": "https://example.com/identity/v3/domains/9fe1d3",
	},
	Name: "domain two",
}

// SecondDomainUpdated is how SecondDomain should look after an Update.
var SecondDomainUpdated = domains.Domain{
	Enabled: true,
	ID:      "9fe1d3",
	Links: map[string]interface{}{
		"self": "https://example.com/identity/v3/domains/9fe1d3",
	},
	Name:        "domain two",
	Description: "Staging Domain",
}

// ExpectedDomainsSlice is the slice of domains expected to be returned from ListDomainOutput.
var ExpectedDomainsSlice = []domains.Domain{FirstDomain, SecondDomain}

// HandleListDomainsSuccessfully creates an HTTP handler at `/domains` on the
// test handler mux that responds with a list of two domains.
func HandleListDomainsSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/domains", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, ListDomainOutput)
	})
}

// HandleGetDomainSuccessfully creates an HTTP handler at `/domains` on the
// test handler mux that responds with a single domain.
func HandleGetDomainSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/domains/9fe1d3", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, GetDomainOutput)
	})
}

// HandleCreateDomainSuccessfully creates an HTTP handler at `/domains` on the
// test handler mux that tests domain creation.
func HandleCreateDomainSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/domains", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, CreateDomainRequest)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, GetDomainOutput)
	})
}

// HandleDeleteDomainSuccessfully creates an HTTP handler at `/domains` on the
// test handler mux that tests domain deletion.
func HandleDeleteDomainSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/domains/9fe1d3", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleUpdateDomainSuccessfully creates an HTTP handler at `/domains` on the
// test handler mux that tests domain update.
func HandleUpdateDomainSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/domains/9fe1d3", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, UpdateDomainRequest)

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, UpdateDomainOutput)
	})
}
