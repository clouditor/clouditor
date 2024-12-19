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
	"testing"

	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/volumes"
)

func Test_getParentID(t *testing.T) {
	type args struct {
		volume *volumes.Volume
	}
	tests := []struct {
		name string
		args args
		want assert.Want[string]
	}{
		{
			name: "Happy path: no attached server available",
			args: args{
				&volumes.Volume{
					TenantID: testdata.MockVolumeTenantID,
				},
			},
			want: func(t *testing.T, got string) bool {
				return assert.Equal(t, testdata.MockVolumeTenantID, got)
			},
		},
		{
			name: "Happy path: attached serverID",
			args: args{
				&volumes.Volume{
					TenantID: testdata.MockVolumeTenantID,
					Attachments: []volumes.Attachment{
						{
							ServerID: testdata.MockServerID1,
						},
					},
				},
			},
			want: func(t *testing.T, got string) bool {
				return assert.Equal(t, testdata.MockServerID1, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getParentID(tt.args.volume)

			tt.want(t, got)
		})
	}
}
