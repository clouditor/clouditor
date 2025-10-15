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

package ionos

import (
	"testing"

	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/util"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func Test_ionosDiscovery_discoverDatacenters(t *testing.T) {
	type fields struct {
		ionosDiscovery *ionosDiscovery
	}
	tests := []struct {
		name    string
		fields  fields
		want    assert.Want[*ionoscloud.Datacenters]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "error: could not list datacenters",
			fields: fields{
				ionosDiscovery: NewMockIonosDiscovery(newMockErrorSender()),
			},
			want: assert.Empty[*ionoscloud.Datacenters],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not list datacenters:")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				ionosDiscovery: NewMockIonosDiscovery(newMockSender()),
			},
			want: func(t *testing.T, got *ionoscloud.Datacenters) bool {
				return assert.Equal(t, 2, len(util.Deref(got.Items)))
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.ionosDiscovery
			got, _, err := d.discoverDatacenters()
			// TODO(all): Check list as well

			tt.want(t, got)
			tt.wantErr(t, err)
		})
	}
}
