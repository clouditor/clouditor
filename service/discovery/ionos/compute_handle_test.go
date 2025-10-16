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

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/util"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_ionosDiscovery_handleServer(t *testing.T) {
	type fields struct {
		ionosDiscovery *ionosDiscovery
	}
	type args struct {
		server ionoscloud.Server
		dc     ionoscloud.Datacenter
	}
	tests := []struct {
		name    string
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
				server: ionoscloud.Server{
					Id: util.Ref(testdata.MockIonosVMID1),
					Properties: &ionoscloud.ServerProperties{
						Name: util.Ref(testdata.MockIonosVMName1),
					},
					Metadata: &ionoscloud.DatacenterElementMetadata{
						CreatedDate: &ionoscloud.IonosTime{Time: testdata.CreationTime},
					},
					Entities: &ionoscloud.ServerEntities{
						Volumes: &ionoscloud.AttachedVolumes{
							Items: &[]ionoscloud.Volume{
								{
									Id: util.Ref(testdata.MockIonosVolumeID1),
									Properties: &ionoscloud.VolumeProperties{
										Name: util.Ref(testdata.MockIonosVolumeName1),
									},
								},
							},
						},
						Nics: &ionoscloud.Nics{
							Items: &[]ionoscloud.Nic{
								{
									Id: util.Ref(testdata.MockIonosNicID1),
									Properties: &ionoscloud.NicProperties{
										Name: util.Ref(testdata.MockIonosNicName1),
									},
								},
							},
						},
					},
				},
				dc: ionoscloud.Datacenter{
					Id: util.Ref(testdata.MockIonosDatacenterID1),
					Properties: &ionoscloud.DatacenterProperties{
						Location: util.Ref(testdata.MockIonosDatacenterLocation1),
					},
				},
			},
			want: func(t *testing.T, got ontology.IsResource) bool {
				want := &ontology.VirtualMachine{
					Id:   testdata.MockIonosVMID1,
					Name: testdata.MockIonosVMName1,
					Labels: map[string]string{
						"label1": "value1",
						"label2": "value2",
					},
					CreationTime: timestamppb.New(testdata.CreationTime),
					GeoLocation: &ontology.GeoLocation{
						Region: testdata.MockIonosDatacenterLocation1,
					},
					ParentId:        util.Ref(testdata.MockIonosDatacenterID1),
					BlockStorageIds: []string{testdata.MockIonosVolumeID1},
					NetworkInterfaceIds: []string{
						testdata.MockIonosNicID1,
					},
					ActivityLogging: &ontology.ActivityLogging{
						Enabled: true,
					},
				}
				got0 := got.(*ontology.VirtualMachine)
				// Check if raw field is not empty and then remove it for comparison
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
			got, err := d.handleServer(tt.args.server, tt.args.dc)

			tt.want(t, got)
			tt.wantErr(t, err)
		})
	}
}

func Test_ionosDiscovery_handleBlockStorage(t *testing.T) {
	type fields struct {
		ionosDiscovery *ionosDiscovery
	}
	type args struct {
		blockStorage ionoscloud.Volume
		dc           ionoscloud.Datacenter
	}
	tests := []struct {
		name    string
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
				blockStorage: ionoscloud.Volume{
					Id: util.Ref(testdata.MockIonosVolumeID1),
					Properties: &ionoscloud.VolumeProperties{
						Name: util.Ref(testdata.MockIonosVolumeName1),
					},
					Metadata: &ionoscloud.DatacenterElementMetadata{
						CreatedDate: &ionoscloud.IonosTime{Time: testdata.CreationTime},
					},
				},
				dc: ionoscloud.Datacenter{
					Id: util.Ref(testdata.MockIonosDatacenterID1),
					Properties: &ionoscloud.DatacenterProperties{
						Location: util.Ref(testdata.MockIonosDatacenterLocation2),
					},
				},
			},
			want: func(t *testing.T, got ontology.IsResource) bool {
				want := &ontology.BlockStorage{
					Id:   testdata.MockIonosVolumeID1,
					Name: testdata.MockIonosVolumeName1,
					Labels: map[string]string{
						"label1": "value1",
						"label2": "value2",
					},
					CreationTime: timestamppb.New(testdata.CreationTime),
					GeoLocation: &ontology.GeoLocation{
						Region: testdata.MockIonosDatacenterLocation2,
					},
					ParentId: util.Ref(testdata.MockIonosDatacenterID1),
				}
				got0 := got.(*ontology.BlockStorage)

				// Check if raw field is not empty and then remove it for comparison
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
			got, err := d.handleBlockStorage(tt.args.blockStorage, tt.args.dc)

			tt.want(t, got)
			tt.wantErr(t, err)
		})
	}
}

func Test_ionosDiscovery_handleLoadBalancer(t *testing.T) {
	type fields struct {
		ionosDiscovery *ionosDiscovery
	}
	type args struct {
		loadBalancer ionoscloud.Loadbalancer
		dc           ionoscloud.Datacenter
	}
	tests := []struct {
		name    string
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
				loadBalancer: ionoscloud.Loadbalancer{
					Id: util.Ref(string(testdata.MockIonosLoadBalancerID1)),
					Properties: &ionoscloud.LoadbalancerProperties{
						Name: util.Ref(string(testdata.MockIonosLoadBalancerName1)),
					},
					Metadata: &ionoscloud.DatacenterElementMetadata{
						CreatedDate: &ionoscloud.IonosTime{Time: testdata.CreationTime},
					},
				},
				dc: ionoscloud.Datacenter{
					Id: util.Ref(testdata.MockIonosDatacenterID1),
					Properties: &ionoscloud.DatacenterProperties{
						Location: util.Ref(testdata.MockIonosDatacenterLocation1),
					},
				},
			},
			want: func(t *testing.T, got ontology.IsResource) bool {
				want := &ontology.LoadBalancer{
					Id:   testdata.MockIonosLoadBalancerID1,
					Name: testdata.MockIonosLoadBalancerName1,
					Labels: map[string]string{
						"label1": "value1",
						"label2": "value2",
					},
					CreationTime: timestamppb.New(testdata.CreationTime),
					GeoLocation: &ontology.GeoLocation{
						Region: testdata.MockIonosDatacenterLocation1,
					},
					ParentId: util.Ref(testdata.MockIonosDatacenterID1),
				}
				got0 := got.(*ontology.LoadBalancer)

				// Check if raw field is not empty and then remove it for comparison
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
			got, err := d.handleLoadBalancer(tt.args.loadBalancer, tt.args.dc)

			tt.want(t, got)
			tt.wantErr(t, err)
		})
	}
}
