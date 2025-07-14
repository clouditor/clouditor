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

func Test_labels(t *testing.T) {
	type args struct {
		labels ionoscloud.LabelResources
	}
	tests := []struct {
		name string
		args args
		want assert.Want[map[string]string]
	}{
		{
			name: "Happy path",
			args: args{
				labels: ionoscloud.LabelResources{
					Items: &[]ionoscloud.LabelResource{
						{
							Properties: &ionoscloud.LabelResourceProperties{
								Key:   util.Ref("key1"),
								Value: util.Ref("value1"),
							},
						},
						{
							Properties: &ionoscloud.LabelResourceProperties{
								Key:   util.Ref("key2"),
								Value: util.Ref("value2"),
							},
						},
					},
				},
			},
			want: func(t *testing.T, got map[string]string) bool {
				want := map[string]string{
					"key1": "value1",
					"key2": "value2",
				}
				return assert.Equal(t, want, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := labels(tt.args.labels)

			tt.want(t, got)
		})
	}
}
