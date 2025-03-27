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
	"clouditor.io/clouditor/v2/internal/testutil/assert"
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

func TestEvidence_GetResourceId(t *testing.T) {
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
		want   string
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
			want: "vm-1",
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

			assert.Equal(t, tt.want, ev.GetResourceId())
		})
	}
}
