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

package ionos

import (
	"testing"

	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/util"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func Test_getBlockStorageIds(t *testing.T) {
	type args struct {
		server ionoscloud.Server
	}
	tests := []struct {
		name string
		args args
		want assert.Want[[]string]
	}{
		{
			name: "no volumes available",
			args: args{
				server: ionoscloud.Server{
					Entities: &ionoscloud.ServerEntities{
						Volumes: &ionoscloud.AttachedVolumes{},
					},
				},
			},
			want: func(t *testing.T, got []string) bool {
				return assert.Empty(t, got)
			},
		},
		{
			name: "Happy path",
			args: args{
				server: ionoscloud.Server{
					Entities: &ionoscloud.ServerEntities{
						Volumes: &ionoscloud.AttachedVolumes{
							Items: &[]ionoscloud.Volume{
								{
									Id: util.Ref(testdata.MockIonosVolumeID1),
								},
								{
									Id: util.Ref(testdata.MockIonosVolumeID2),
								},
							},
						},
					},
				},
			},
			want: func(t *testing.T, got []string) bool {
				want := []string{testdata.MockIonosVolumeID1, testdata.MockIonosVolumeID2}
				return assert.Equal(t, want, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBlockStorageIds := getBlockStorageIds(tt.args.server)
			tt.want(t, gotBlockStorageIds)
		})
	}
}

func Test_getNetworkInterfaceIds(t *testing.T) {
	type args struct {
		server ionoscloud.Server
	}
	tests := []struct {
		name string
		args args
		want assert.Want[[]string]
	}{
		{
			name: "Happy path",
			args: args{
				server: ionoscloud.Server{
					Entities: &ionoscloud.ServerEntities{
						Nics: &ionoscloud.Nics{
							Items: &[]ionoscloud.Nic{
								{
									Id: util.Ref(testdata.MockIonosNicID1),
								},
								{
									Id: util.Ref(testdata.MockIonosNicID2),
								},
							},
						},
					},
				},
			},
			want: func(t *testing.T, got []string) bool {
				want := []string{testdata.MockIonosNicID1, testdata.MockIonosNicID2}
				return assert.Equal(t, want, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNetworkInterfaceIds := getNetworkInterfaceIds(tt.args.server)

			tt.want(t, gotNetworkInterfaceIds)
		})
	}
}
