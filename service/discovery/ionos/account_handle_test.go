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

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/util"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_ionosDiscovery_handleDatacenter(t *testing.T) {
	type fields struct {
		ionosDiscovery *ionosDiscovery
	}
	type args struct {
		dc ionoscloud.Datacenter
	}
	tests := []struct {
		name    string // description of this test case
		fields  fields
		args    args
		want    assert.Want[ontology.IsResource]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path",
			fields: fields{
				ionosDiscovery: NewMockIonosDiscovery(newMockSender()),
			},
			args: args{
				dc: ionoscloud.Datacenter{
					Id: util.Ref(testdata.MockIonosDatacenterID1),
					Properties: &ionoscloud.DatacenterProperties{
						Name:        util.Ref(testdata.MockIonosDatacenterName1),
						Description: util.Ref(testdata.MockIonosDatacenterDescription1),
						Location:    util.Ref(testdata.MockIonosDatacenterLocation1),
					},
					Metadata: &ionoscloud.DatacenterElementMetadata{
						CreatedDate: &ionoscloud.IonosTime{Time: testdata.CreationTime},
					},
				},
			},
			want: func(t *testing.T, got ontology.IsResource) bool {
				want := &ontology.Account{
					Id:           testdata.MockIonosDatacenterID1,
					Name:         testdata.MockIonosDatacenterName1,
					Description:  testdata.MockIonosDatacenterDescription1,
					CreationTime: timestamppb.New(testdata.CreationTime),
					GeoLocation: &ontology.GeoLocation{
						Region: testdata.MockIonosDatacenterLocation1,
					},
					Labels: map[string]string{
						"label1": "value1",
						"label2": "value2",
					},
				}

				got0 := got.(*ontology.Account)
				// Check if raw field is not empty and remove the content for comparison
				assert.NotEmpty(t, got0.Raw)
				got0.Raw = ""

				return assert.Equal(t, want, got0)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.ionosDiscovery
			got, gotErr := d.handleDatacenter(tt.args.dc)

			tt.want(t, got)
			tt.wantErr(t, gotErr)
		})
	}
}
