// Copyright 2021 Fraunhofer AISEC
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

package discovery

import (
	reflect "reflect"
	"testing"

	"clouditor.io/clouditor/api/ontology"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/testutil/prototest"
	"clouditor.io/clouditor/internal/util"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestResource_ToOntologyResource(t *testing.T) {
	type fields struct {
		Id             string
		CloudServiceId string
		ResourceType   string
		Properties     *anypb.Any
	}
	tests := []struct {
		name    string
		fields  fields
		want    ontology.IsResource
		wantErr bool
	}{
		{
			name: "happy path VM",
			fields: fields{
				Id:             "vm1",
				CloudServiceId: "service1",
				ResourceType:   "VirtualMachine",
				Properties: prototest.NewAny(t, &ontology.VirtualMachine{
					Id:              "vm1",
					BlockStorageIds: []string{"bs1"},
				}),
			},
			want: &ontology.VirtualMachine{
				Id:              "vm1",
				BlockStorageIds: []string{"bs1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resource{
				Id:             tt.fields.Id,
				CloudServiceId: tt.fields.CloudServiceId,
				ResourceType:   tt.fields.ResourceType,
				Properties:     tt.fields.Properties,
			}
			got, err := r.ToOntologyResource()
			if (err != nil) != tt.wantErr {
				t.Errorf("Resource.ToOntologyResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !proto.Equal(got, tt.want) {
				t.Errorf("Resource.ToOntologyResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToDiscoveryResource(t *testing.T) {
	type args struct {
		resource ontology.IsResource
		csID     string
	}
	tests := []struct {
		name    string
		args    args
		wantR   *Resource
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				resource: &ontology.BlockStorage{
					Id:   "my-block-storage",
					Name: "My Block Storage",
					Backup: []*ontology.Backup{
						{
							Enabled:   true,
							StorageId: util.Ref("my-offsite-backup-id"),
						},
					},
				},
				csID: testdata.MockCloudServiceID1,
			},
			wantR: &Resource{
				Id:             "my-block-storage",
				CloudServiceId: testdata.MockCloudServiceID1,
				ResourceType:   "BlockStorage,Storage,CloudResource,Resource",
				Properties: prototest.NewAny(t, &ontology.BlockStorage{
					Id:   "my-block-storage",
					Name: "My Block Storage",
					Backup: []*ontology.Backup{
						{
							Enabled:   true,
							StorageId: util.Ref("my-offsite-backup-id"),
						},
					},
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := ToDiscoveryResource(tt.args.resource, tt.args.csID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToDiscoveryResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("ToDiscoveryResource() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}
