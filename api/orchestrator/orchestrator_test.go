// Copyright 2022 Fraunhofer AISEC
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

package orchestrator

import (
	"database/sql/driver"
	"reflect"
	"testing"
)

func TestCloudService_Requirements_Scan(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    *CloudService_Requirements
	}{
		{
			name: "string type",
			args: args{
				value: "a,b",
			},
			wantErr: false,
			want: &CloudService_Requirements{
				RequirementIds: []string{"a", "b"},
			},
		},
		{
			name: "unsupported type",
			args: args{
				value: 0,
			},
			wantErr: true,
			want:    &CloudService_Requirements{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CloudService_Requirements{}

			if err := c.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("CloudService_Requirements.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(c, tt.want) {
				t.Errorf("CloudService_Requirements.Scan() = %v, want %v", c, tt.want)
			}
		})
	}
}

func TestCloudService_Requirements_Value(t *testing.T) {
	type fields struct {
		RequirementIds []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    driver.Value
		wantErr bool
	}{
		{
			name: "empty requirements",
			fields: fields{
				RequirementIds: nil,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "some requirements",
			fields: fields{
				RequirementIds: []string{"a", "b"},
			},
			want:    driver.Value(string("a,b")),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CloudService_Requirements{
				RequirementIds: tt.fields.RequirementIds,
			}

			got, err := c.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("CloudService_Requirements.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CloudService_Requirements.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
