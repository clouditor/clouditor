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
	"reflect"
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
		name          string
		args          args
		wantResp      string
		wantRespError error
	}{
		{
			name: "Missing resource",
			args: args{
				Evidence: &Evidence{
					Id:        "11111111-1111-1111-1111-111111111111",
					ToolId:    "mock",
					Timestamp: timestamppb.Now(),
				},
			},
			wantResp:      "",
			wantRespError: ErrNotValidResource,
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
				}},
			wantResp:      "",
			wantRespError: ErrResourceNotStruct,
		},
		{
			name: "Resource is missing the id field",
			args: args{
				Evidence: &Evidence{
					Id:        "11111111-1111-1111-1111-111111111111",
					Timestamp: timestamppb.Now(),
					ToolId:    "mock",
					Resource: toStruct(voc.VirtualMachine{
						Compute: &voc.Compute{
							CloudResource: &voc.CloudResource{
								Type: []string{"VirtualMachine"}},
						},
					}, t),
				}},
			wantRespError: ErrResourceIdMissing,
		},
		{
			name: "Id of resource is empty",
			args: args{
				Evidence: &Evidence{
					Id:        "11111111-1111-1111-1111-111111111111",
					Timestamp: timestamppb.Now(),
					ToolId:    "mock",
					Resource: toStruct(voc.VirtualMachine{
						Compute: &voc.Compute{
							CloudResource: &voc.CloudResource{
								ID:   "",
								Type: []string{"VirtualMachine"}},
						},
					}, t),
				}},
			wantRespError: ErrResourceIdMissing,
		},
		{
			name: "Type of resource is nil",
			args: args{
				Evidence: &Evidence{
					Id:        "11111111-1111-1111-1111-111111111111",
					Timestamp: timestamppb.Now(),
					ToolId:    "mock",
					Resource: toStruct(voc.VirtualMachine{
						Compute: &voc.Compute{
							CloudResource: &voc.CloudResource{
								ID:   "my-resource-id",
								Type: nil,
							},
						},
					}, t),
				}},
			// Since type is nil, the field exists, but it is not of array type
			wantRespError: ErrResourceTypeNotArrayOfStrings,
		},
		{
			name: "Array of resource types is empty",
			args: args{
				Evidence: &Evidence{
					Id:        "11111111-1111-1111-1111-111111111111",
					Timestamp: timestamppb.Now(),
					ToolId:    "mock",
					Resource: toStruct(voc.VirtualMachine{
						Compute: &voc.Compute{
							CloudResource: &voc.CloudResource{
								ID:   "my-resource-id",
								Type: []string{},
							},
						},
					}, t),
				}},
			wantRespError: ErrResourceTypeEmpty,
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
				}},
			wantRespError: ErrToolIdMissing,
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
			wantResp:      "",
			wantRespError: ErrTimestampMissing,
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
				}},
			wantResp:      "my-resource-id",
			wantRespError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resourceId, err := tt.args.Evidence.Validate()

			if isError := assert.ErrorIs(t, err, tt.wantRespError); isError {
				// If we have an error, we do not need to check other parts of the response (resourceId)
				return
			}

			if !reflect.DeepEqual(resourceId, tt.wantResp) {
				t.Errorf("Validate() gotResp = %v, want %v", resourceId, tt.wantResp)
			}

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
