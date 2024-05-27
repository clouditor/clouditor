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

// csaf contains a discover that discovery security advisory information from a CSAF trusted provider
package csaf

import (
	"net/http"
	"os"
	"testing"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/discoverytest/csaf/providertest"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/csaf-poc/csaf_distribution/v3/csaf"
)

// validAdvisory contains the structure of a valid CSAF Advisory that validates against the JSON schema
var validAdvisory = &csaf.Advisory{
	Document: &csaf.Document{
		Category:    util.Ref(csaf.DocumentCategory("csaf_security_advisory")),
		CSAFVersion: util.Ref(csaf.CSAFVersion20),
		Title:       util.Ref("Buffer overflow in Test Product"),
		Publisher: &csaf.DocumentPublisher{
			Name:      util.Ref("Test Vendor"),
			Category:  util.Ref(csaf.CSAFCategoryVendor),
			Namespace: util.Ref("http://localhost"),
		},
		Tracking: &csaf.Tracking{
			ID:                 util.Ref(csaf.TrackingID("some-id")),
			CurrentReleaseDate: util.Ref("2020-07-01T10:09:07Z"),
			InitialReleaseDate: util.Ref("2020-07-01T10:09:07Z"),
			Generator: &csaf.Generator{
				Date: util.Ref("2020-07-01T10:09:07Z"),
				Engine: &csaf.Engine{
					Name:    util.Ref("test"),
					Version: util.Ref("1.0"),
				},
			},
			Status:  util.Ref(csaf.CSAFTrackingStatusFinal),
			Version: util.Ref(csaf.RevisionNumber("1")),
			RevisionHistory: csaf.Revisions{
				&csaf.Revision{
					Date:    util.Ref("2020-07-01T10:09:07Z"),
					Number:  util.Ref(csaf.RevisionNumber("1")),
					Summary: util.Ref("First and final version"),
				},
			},
		},
	},
	ProductTree: &csaf.ProductTree{
		Branches: csaf.Branches{
			&csaf.Branch{
				Category: util.Ref(csaf.CSAFBranchCategoryVendor),
				Name:     util.Ref("Test Vendor"),
				Product: &csaf.FullProductName{
					Name:      util.Ref("Test Product"),
					ProductID: util.Ref(csaf.ProductID("CSAFPID-0001")),
				},
			},
		},
	},
}

var goodProvider *providertest.TrustedProvider

func TestMain(m *testing.M) {
	var advisories = map[csaf.TLPLabel][]*csaf.Advisory{
		csaf.TLPLabelWhite: {
			validAdvisory,
		},
	}

	goodProvider = providertest.NewTrustedProvider(
		advisories,
		providertest.NewGoodIndexTxtWriter(),
		func(pmd *csaf.ProviderMetadata) {
			pmd.Publisher = &csaf.Publisher{
				Name:      util.Ref("Test Vendor"),
				Category:  util.Ref(csaf.CSAFCategoryVendor),
				Namespace: util.Ref("http://localhost"),
			}
		})
	defer goodProvider.Close()

	code := m.Run()
	os.Exit(code)
}

func TestNewTrustedProviderDiscovery(t *testing.T) {
	type args struct {
		opts []DiscoveryOption
	}
	tests := []struct {
		name string
		args args
		want discovery.Discoverer
	}{
		{
			name: "Happy path",
			args: args{},
			want: &csafDiscovery{
				csID:   config.DefaultCloudServiceID,
				domain: "clouditor.io",
				client: http.DefaultClient,
			},
		},
		{
			name: "Happy path: with cloud service id",
			args: args{
				opts: []DiscoveryOption{WithCloudServiceID(testdata.MockCloudServiceID1)},
			},
			want: &csafDiscovery{
				csID:   testdata.MockCloudServiceID1,
				domain: "clouditor.io",
				client: http.DefaultClient,
			},
		},
		{
			name: "Happy path: with domain",
			args: args{
				opts: []DiscoveryOption{WithProviderDomain("mock")},
			},
			want: &csafDiscovery{
				csID:   config.DefaultCloudServiceID,
				client: http.DefaultClient,
				domain: "mock",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTrustedProviderDiscovery(tt.args.opts...)
			assert.Equal(t, tt.want, got, assert.CompareAllUnexported())
		})
	}
}

func Test_csafDiscovery_List(t *testing.T) {
	type fields struct {
		domain string
		csID   string
		client *http.Client
	}
	tests := []struct {
		name     string
		fields   fields
		wantList assert.Want[[]ontology.IsResource]
		wantErr  assert.WantErr
	}{
		{
			name: "fail",
			fields: fields{
				domain: "localhost:1234",
				client: http.DefaultClient,
				csID:   config.DefaultCloudServiceID,
			},
			wantErr: func(t *testing.T, err error) bool {
				return assert.ErrorContains(t, err, "could not load provider-metadata.json")
			},
			wantList: assert.Empty[[]ontology.IsResource],
		},
		{
			name: "happy path",
			fields: fields{
				domain: goodProvider.Domain(),
				client: goodProvider.Client(),
				csID:   config.DefaultCloudServiceID,
			},
			wantErr: func(t *testing.T, err error) bool {
				return assert.NoError(t, err)
			},
			wantList: func(t *testing.T, got []ontology.IsResource) bool {
				return assert.NotEmpty(t, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &csafDiscovery{
				domain: tt.fields.domain,
				client: tt.fields.client,
				csID:   tt.fields.csID,
			}
			gotList, err := d.List()
			tt.wantErr(t, err)
			tt.wantList(t, gotList)
		})
	}
}
