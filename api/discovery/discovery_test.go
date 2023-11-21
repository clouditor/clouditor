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
	"reflect"
	"testing"

	"clouditor.io/clouditor/voc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

func TestResource_ToVocResource(t *testing.T) {
	type fields struct {
		Id             string
		CloudServiceId string
		ResourceType   string
		Properties     *structpb.Value
	}
	tests := []struct {
		name    string
		fields  fields
		want    voc.IsCloudResource
		wantErr bool
	}{
		{
			name: "happy path VM",
			fields: fields{
				Id:             "vm1",
				CloudServiceId: "service1",
				ResourceType:   "VirtualMachine,Compute",
				Properties: func() *structpb.Value {
					var s structpb.Struct
					raw := []byte(`{"blockStorage": ["bs1"]}`)
					err := protojson.Unmarshal(raw, &s)
					assert.NoError(t, err)
					return structpb.NewStructValue(&s)
				}(),
			},
			want: &voc.VirtualMachine{
				BlockStorage: []voc.ResourceID{"bs1"},
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
			got, err := r.ToVocResource()
			if (err != nil) != tt.wantErr {
				t.Errorf("Resource.ToVocResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resource.ToVocResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
