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
	"testing"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/util"
	"clouditor.io/clouditor/v2/service/discovery/extra/csaf/providertest"

	"github.com/csaf-poc/csaf_distribution/v3/csaf"
)

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
				csID:   discovery.DefaultCloudServiceID,
				domain: "wid.cert-bund.de",
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
				domain: "wid.cert-bund.de",
				client: http.DefaultClient,
			},
		},
		{
			name: "Happy path: with domain",
			args: args{
				opts: []DiscoveryOption{WithProviderDomain("mock")},
			},
			want: &csafDiscovery{
				csID:   discovery.DefaultCloudServiceID,
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
	p := providertest.NewTrustedProvider(func(pmd *csaf.ProviderMetadata) {
		pmd.Publisher = &csaf.Publisher{
			Name:      util.Ref("Test Vendor"),
			Category:  util.Ref(csaf.CSAFCategoryVendor),
			Namespace: util.Ref("http://localhost"),
		}
	})
	defer p.Close()

	type fields struct {
		domain string
		csID   string
		client *http.Client
	}
	tests := []struct {
		name     string
		fields   fields
		wantList []ontology.IsResource
		wantErr  assert.WantErr
	}{
		{
			name: "fail",
			fields: fields{
				domain: "localhost:1234",
				client: http.DefaultClient,
				csID:   discovery.DefaultCloudServiceID,
			},
			wantErr: func(t *testing.T, err error) bool {
				return assert.ErrorContains(t, err, "could not load provider-metadata.json")
			},
		},
		{
			name: "happy path",
			fields: fields{
				domain: p.Domain(),
				client: p.Client(),
				csID:   discovery.DefaultCloudServiceID,
			},
			wantErr: func(t *testing.T, err error) bool {
				return assert.NoError(t, err)
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
			assert.Equal(t, tt.wantList, gotList)
		})
	}
}
