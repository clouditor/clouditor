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

	"clouditor.io/clouditor/voc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_ValidateEvidence(t *testing.T) {
	type args struct {
		Evidence *Evidence
	}

	tests := []struct {
		name     string
		args     args
		wantResp string
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "Missing resource",
			args: args{
				Evidence: &Evidence{
					Id:             "11111111-1111-1111-1111-111111111111",
					ToolId:         "mock",
					Timestamp:      timestamppb.Now(),
					CloudServiceId: "11111111-1111-1111-1111-111111111111",
				},
			},
			wantResp: "",
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "Resource: value is required")
			},
		},
		{
			name: "Resource is not a struct",
			args: args{
				Evidence: &Evidence{
					Id:        "11111111-1111-1111-1111-111111111111",
					Timestamp: timestamppb.Now(),
					ToolId:    "mock",
					Resource: &structpb.Value{
						Kind: &structpb.Value_StringValue{
							StringValue: "MockTargetValue",
						},
					},
					CloudServiceId: "11111111-1111-1111-1111-111111111111",
				}},
			wantResp: "",
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrResourceNotStruct)
			},
		},
		{
			name: "Missing toolId",
			args: args{
				Evidence: &Evidence{
					Id:        "11111111-1111-1111-1111-111111111111",
					Timestamp: timestamppb.Now(),
					Resource: toStruct(voc.VirtualMachine{
						Compute: &voc.Compute{
							Resource: &voc.Resource{
								ID:   "my-resource-id",
								Type: []string{"VirtualMachine"}},
						},
					}, t),
					CloudServiceId: "11111111-1111-1111-1111-111111111111",
				}},
			wantResp: "",
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "ToolId: value length must be at least 1 runes")
			},
		},
		{
			name: "Missing timestamp",
			args: args{
				Evidence: &Evidence{
					Id:     "11111111-1111-1111-1111-111111111111",
					ToolId: "mock",
					Resource: toStruct(voc.VirtualMachine{
						Compute: &voc.Compute{
							Resource: &voc.Resource{
								ID:   "my-resource-id",
								Type: []string{"VirtualMachine"}},
						},
					}, t),
				}},
			wantResp: "",
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "Timestamp: value is required")
			},
		},
		{
			name: "Valid evidence",
			args: args{
				Evidence: &Evidence{
					Id:        "11111111-1111-1111-1111-111111111111",
					Timestamp: timestamppb.Now(),
					ToolId:    "mock",
					Resource: toStruct(voc.VirtualMachine{
						Compute: &voc.Compute{
							Resource: &voc.Resource{
								ID:   "my-resource-id",
								Type: []string{"VirtualMachine"}},
						},
					}, t),
					CloudServiceId: "11111111-1111-1111-1111-111111111111",
				}},
			wantResp: "my-resource-id",
			wantErr:  assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := tt.args.Evidence.ValidateResource()
			tt.wantErr(t, err)
		})
	}
}

// toStruct transforms r to a struct and asserts if it was successful
func toStruct(r voc.IsCloudResource, t *testing.T) (s *structpb.Value) {
	s, err := voc.ToStruct(r)
	if err != nil {
		assert.Error(t, err)
	}

	return
}
