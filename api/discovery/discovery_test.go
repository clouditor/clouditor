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
	"testing"

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/prototest"
	"clouditor.io/clouditor/v2/internal/util"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestResource_ToOntologyResource(t *testing.T) {
	type fields struct {
		Id                    string
		CertificationTargetId string
		ResourceType          string
		Properties            *anypb.Any
	}
	tests := []struct {
		name    string
		fields  fields
		want    ontology.IsResource
		wantErr assert.WantErr
	}{
		{
			name: "happy path VM",
			fields: fields{
				Id:                    "vm1",
				CertificationTargetId: "target1",
				ResourceType:          "VirtualMachine",
				Properties: prototest.NewAny(t, &ontology.VirtualMachine{
					Id:              "vm1",
					BlockStorageIds: []string{"bs1"},
				}),
			},
			want: &ontology.VirtualMachine{
				Id:              "vm1",
				BlockStorageIds: []string{"bs1"},
			},
			wantErr: assert.Nil[error],
		},
		{
			name: "not an ontology resource",
			fields: fields{
				Id:                    "vm1",
				CertificationTargetId: "target1",
				ResourceType:          "Something",
				Properties:            prototest.NewAny(t, &emptypb.Empty{}),
			},
			want: nil,
			wantErr: func(t *testing.T, err error) bool {
				return assert.ErrorContains(t, err, ErrNotOntologyResource.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resource{
				Id:                    tt.fields.Id,
				CertificationTargetId: tt.fields.CertificationTargetId,
				ResourceType:          tt.fields.ResourceType,
				Properties:            tt.fields.Properties,
			}
			got, err := r.ToOntologyResource()

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestToDiscoveryResource(t *testing.T) {
	type args struct {
		resource    ontology.IsResource
		ctID        string
		collectorID string
	}
	tests := []struct {
		name    string
		args    args
		want    *Resource
		wantErr assert.WantErr
	}{
		{
			name: "happy path",
			args: args{
				resource: &ontology.BlockStorage{
					Id:   "my-block-storage",
					Name: "My Block Storage",
					Backups: []*ontology.Backup{
						{
							Enabled:   true,
							StorageId: util.Ref("my-offsite-backup-id"),
						},
					},
				},
				ctID:        testdata.MockCertificationTargetID1,
				collectorID: testdata.MockEvidenceToolID1,
			},
			want: &Resource{
				Id:                    "my-block-storage",
				CertificationTargetId: testdata.MockCertificationTargetID1,
				ToolId:                testdata.MockEvidenceToolID1,
				ResourceType:          "BlockStorage,Storage,Infrastructure,Resource",
				Properties: prototest.NewAny(t, &ontology.BlockStorage{
					Id:   "my-block-storage",
					Name: "My Block Storage",
					Backups: []*ontology.Backup{
						{
							Enabled:   true,
							StorageId: util.Ref("my-offsite-backup-id"),
						},
					},
				}),
			},
			wantErr: assert.Nil[error],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := ToDiscoveryResource(tt.args.resource, tt.args.ctID, tt.args.collectorID)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, gotR)
		})
	}
}
