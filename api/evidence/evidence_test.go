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

package evidence

import (
	"testing"

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/prototest"
	"clouditor.io/clouditor/v2/internal/util"
	anypb "google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestEvidence_GetOntologyResource(t *testing.T) {
	type fields struct {
		Id                             string
		Timestamp                      *timestamppb.Timestamp
		TargetOfEvaluationId           string
		ToolId                         string
		Resource                       *ontology.Resource
		ExperimentalRelatedResourceIds []string
	}
	tests := []struct {
		name   string
		fields fields
		want   ontology.IsResource
	}{
		{
			name: "happy path",
			fields: fields{
				Resource: &ontology.Resource{
					Type: &ontology.Resource_VirtualMachine{
						VirtualMachine: &ontology.VirtualMachine{
							Id: "vm-1",
						},
					},
				},
			},
			want: &ontology.VirtualMachine{
				Id: "vm-1",
			},
		},
		{
			name: "resource is nil",
			fields: fields{
				Resource: nil,
			},
			want: nil,
		},
		{
			name: "resource is empty",
			fields: fields{
				Resource: &ontology.Resource{},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev := &Evidence{
				Id:                             tt.fields.Id,
				Timestamp:                      tt.fields.Timestamp,
				TargetOfEvaluationId:           tt.fields.TargetOfEvaluationId,
				ToolId:                         tt.fields.ToolId,
				Resource:                       tt.fields.Resource,
				ExperimentalRelatedResourceIds: tt.fields.ExperimentalRelatedResourceIds,
			}

			assert.Equal(t, tt.want, ev.GetOntologyResource())
		})
	}
}

func TestResource_ToOntologyResource(t *testing.T) {
	type fields struct {
		Id                   string
		TargetOfEvaluationId string
		ResourceType         string
		Properties           *anypb.Any
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
				Id:                   "vm1",
				TargetOfEvaluationId: "target1",
				ResourceType:         "VirtualMachine",
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
				Id:                   "vm1",
				TargetOfEvaluationId: "target1",
				ResourceType:         "Something",
				Properties:           prototest.NewAny(t, &emptypb.Empty{}),
			},
			want: nil,
			wantErr: func(t *testing.T, err error) bool {
				return assert.ErrorContains(t, err, ontology.ErrNotOntologyResource.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resource{
				Id:                   tt.fields.Id,
				TargetOfEvaluationId: tt.fields.TargetOfEvaluationId,
				ResourceType:         tt.fields.ResourceType,
				Properties:           tt.fields.Properties,
			}
			got, err := r.ToOntologyResource()

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestToEvidenceResource(t *testing.T) {
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
				ctID:        testdata.MockTargetOfEvaluationID1,
				collectorID: testdata.MockEvidenceToolID1,
			},
			want: &Resource{
				Id:                   "my-block-storage",
				TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
				ToolId:               testdata.MockEvidenceToolID1,
				ResourceType:         "BlockStorage,Storage,Infrastructure,Resource",
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
			gotR, err := ToEvidenceResource(tt.args.resource, tt.args.ctID, tt.args.collectorID)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, gotR)
		})
	}
}
