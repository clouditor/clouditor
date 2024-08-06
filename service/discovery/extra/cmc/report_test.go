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

package cmc

import (
	"reflect"
	"testing"

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
)

func Test_cmcDiscovery_discoverReports(t *testing.T) {
	type fields struct {
		csID       string
		cmcAddr    string
		campemPath string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []ontology.IsResource
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "error read CA from path",
			fields: fields{
				campemPath: "",
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not read certificate from path ")
			},
		},
		{
			name: "error getting attestation report",
			fields: fields{
				campemPath: "../../../../internal/testdata/cmc-discovery/certificate_remote_attestation-empty.pem",
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not add cert")
			},
		},
		{
			name: "error getting attestation report",
			fields: fields{
				campemPath: "../../../../internal/testdata/cmc-discovery/certificate_remote_attestation.pem",
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not get attestation report")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &cmcDiscovery{
				csID:      tt.fields.csID,
				cmcAddr:   tt.fields.cmcAddr,
				capemPath: tt.fields.campemPath,
			}
			got, err := d.discoverReports()

			tt.wantErr(t, err)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cmcDiscovery.discoverReports() = %v, want %v", got, tt.want)
			}
		})
	}
}
