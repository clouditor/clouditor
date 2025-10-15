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
)

func Test_ionosDiscovery_discoverServer(t *testing.T) {
	type fields struct {
		ionosDiscovery *ionosDiscovery
	}
	type args struct {
		dc ionoscloud.Datacenter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.Want[[]ontology.IsResource]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "error: listing servers",
			fields: fields{
				ionosDiscovery: NewMockIonosDiscovery(newMockErrorSender()),
			},
			want: assert.Nil[[]ontology.IsResource],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not list servers for datacenter")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				ionosDiscovery: NewMockIonosDiscovery(newMockSender()),
			},
			args: args{
				dc: ionoscloud.Datacenter{
					Id: util.Ref(testdata.MockIonosDatacenterID1),
					Properties: &ionoscloud.DatacenterProperties{
						Name: util.Ref(testdata.MockIonosDatacenterName1),
					},
					Metadata: &ionoscloud.DatacenterElementMetadata{
						CreatedDate: &ionoscloud.IonosTime{Time: testdata.CreationTime},
					},
				},
			},
			want: func(t *testing.T, got []ontology.IsResource) bool {
				assert.Equal(t, 2, len(got))
				_, ok := got[0].(*ontology.VirtualMachine)
				_, ok1 := got[1].(*ontology.VirtualMachine)
				if !ok || !ok1 {
					return assert.Fail(t, "expected both resources to be VirtualMachine")
				}
				return true
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.ionosDiscovery
			got, err := d.discoverServers(tt.args.dc)

			tt.want(t, got)
			tt.wantErr(t, err)
		})
	}
}

func Test_ionosDiscovery_discoverBlockStorages(t *testing.T) {
	type fields struct {
		ionosDiscovery *ionosDiscovery
	}
	type args struct {
		dc ionoscloud.Datacenter
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantList assert.Want[[]ontology.IsResource]
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "error: listing servers",
			fields: fields{
				ionosDiscovery: NewMockIonosDiscovery(newMockErrorSender()),
			},
			wantList: assert.Nil[[]ontology.IsResource],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not list block storages for datacenter")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				ionosDiscovery: NewMockIonosDiscovery(newMockSender()),
			},
			args: args{
				dc: ionoscloud.Datacenter{
					Id: util.Ref(testdata.MockIonosDatacenterID1),
					Properties: &ionoscloud.DatacenterProperties{
						Name: util.Ref(testdata.MockIonosDatacenterName1),
					},
					Metadata: &ionoscloud.DatacenterElementMetadata{
						CreatedDate: &ionoscloud.IonosTime{Time: testdata.CreationTime},
					},
				},
			},
			wantList: func(t *testing.T, gotList []ontology.IsResource) bool {
				return assert.Equal(t, 2, len(gotList))
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.ionosDiscovery
			gotList, err := d.discoverBlockStorages(tt.args.dc)

			tt.wantList(t, gotList)
			tt.wantErr(t, err)
		})
	}
}

func Test_ionosDiscovery_discoverLoadBalancers(t *testing.T) {
	type fields struct {
		ionosDiscovery *ionosDiscovery
	}
	type args struct {
		dc ionoscloud.Datacenter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.Want[[]ontology.IsResource]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "error: listing servers",
			fields: fields{
				ionosDiscovery: NewMockIonosDiscovery(newMockErrorSender()),
			},
			want: assert.Nil[[]ontology.IsResource],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not list load balancers for datacenter")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				ionosDiscovery: NewMockIonosDiscovery(newMockSender()),
			},
			args: args{
				dc: ionoscloud.Datacenter{
					Id: util.Ref(testdata.MockIonosDatacenterID1),
					Properties: &ionoscloud.DatacenterProperties{
						Name: util.Ref(testdata.MockIonosDatacenterName1),
					},
					Metadata: &ionoscloud.DatacenterElementMetadata{
						CreatedDate: &ionoscloud.IonosTime{Time: testdata.CreationTime},
					},
				},
			},
			want: func(t *testing.T, got []ontology.IsResource) bool {
				assert.Equal(t, 2, len(got))
				_, ok := got[0].(*ontology.LoadBalancer)
				_, ok1 := got[1].(*ontology.LoadBalancer)
				if !ok || !ok1 {
					return assert.Fail(t, "expected both resources to be LoadBalancer")
				}
				return true
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.ionosDiscovery
			gotList, err := d.discoverLoadBalancers(tt.args.dc)

			tt.want(t, gotList)
			tt.wantErr(t, err)
		})
	}
}
