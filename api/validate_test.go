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

package api

import (
	"testing"

	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/voc"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestValidate(t *testing.T) {
	var nilReq *orchestrator.CreateTargetOfEvaluationRequest = nil

	type args struct {
		req IncomingRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Missing request",
			args: args{
				req: nil,
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrEmptyRequest.Error())
			},
		},
		{
			name: "Missing request",
			args: args{
				req: nilReq,
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrEmptyRequest.Error())
			},
		},
		{
			name: "Invalid request",
			args: args{
				req: &orchestrator.CreateTargetOfEvaluationRequest{},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Contains(t, err.Error(), "invalid request")
			},
		},
		{
			name: "Happy path",
			args: args{
				req: &orchestrator.CreateTargetOfEvaluationRequest{
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CloudServiceId: "11111111-1111-1111-1111-111111111111",
						CatalogId:      "0000",
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.args.req)
			tt.wantErr(t, err)
		})
	}
}

func TestValidateWithResource(t *testing.T) {
	type args struct {
		Evidence *evidence.Evidence
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
				Evidence: &evidence.Evidence{
					Id:             testdata.MockCloudServiceID1,
					ToolId:         testdata.MockEvidenceToolID1,
					Timestamp:      timestamppb.Now(),
					CloudServiceId: testdata.MockCloudServiceID1,
				},
			},
			wantResp: "",
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "resource: value is required ")
			},
		},
		{
			name: "Resource is not a struct",
			args: args{
				Evidence: &evidence.Evidence{
					Id:        testdata.MockCloudServiceID1,
					Timestamp: timestamppb.Now(),
					ToolId:    testdata.MockEvidenceToolID1,
					Resource: &structpb.Value{
						Kind: &structpb.Value_StringValue{
							StringValue: "MockTargetValue",
						},
					},
					CloudServiceId: testdata.MockCloudServiceID1,
				}},
			wantResp: "",
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrResourceNotStruct)
			},
		},
		{
			name: "Missing resource Id field",
			args: args{
				Evidence: &evidence.Evidence{
					Id:        testdata.MockCloudServiceID1,
					Timestamp: timestamppb.Now(),
					ToolId:    testdata.MockEvidenceToolID1,
					Resource: toStruct(voc.VirtualMachine{
						Compute: &voc.Compute{
							Resource: &voc.Resource{
								Type: []string{"VirtualMachine"}},
						},
					}, t),
					CloudServiceId: testdata.MockCloudServiceID1,
				}},
			wantResp: "",
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrResourceIdIsEmpty.Error())
			},
		},
		{
			name: "Missing resource type field",
			args: args{
				Evidence: &evidence.Evidence{
					Id:        testdata.MockCloudServiceID1,
					Timestamp: timestamppb.Now(),
					ToolId:    testdata.MockEvidenceToolID1,
					Resource: toStruct(voc.VirtualMachine{
						Compute: &voc.Compute{
							Resource: &voc.Resource{
								ID: testdata.MockResourceID1,
							},
						},
					}, t),
					CloudServiceId: testdata.MockCloudServiceID1,
				}},
			wantResp: "",
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrResourceTypeNotArrayOfStrings.Error())
			},
		},
		{
			name: "Missing resource type field is empty",
			args: args{
				Evidence: &evidence.Evidence{
					Id:        testdata.MockCloudServiceID1,
					Timestamp: timestamppb.Now(),
					ToolId:    testdata.MockEvidenceToolID1,
					Resource: toStruct(voc.VirtualMachine{
						Compute: &voc.Compute{
							Resource: &voc.Resource{
								ID:   testdata.MockResourceID1,
								Type: []string{},
							},
						},
					}, t),
					CloudServiceId: testdata.MockCloudServiceID1,
				}},
			wantResp: "",
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrResourceTypeEmpty.Error())
			},
		},
		{
			name: "Missing toolId",
			args: args{
				Evidence: &evidence.Evidence{
					Id:        testdata.MockCloudServiceID1,
					Timestamp: timestamppb.Now(),
					Resource: toStruct(voc.VirtualMachine{
						Compute: &voc.Compute{
							Resource: &voc.Resource{
								ID:   testdata.MockResourceID1,
								Type: []string{"VirtualMachine"}},
						},
					}, t),
					CloudServiceId: testdata.MockCloudServiceID1,
				}},
			wantResp: "",
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "tool_id: value length must be at least 1 characters")
			},
		},
		{
			name: "Missing timestamp",
			args: args{
				Evidence: &evidence.Evidence{
					Id:     testdata.MockCloudServiceID1,
					ToolId: testdata.MockEvidenceToolID1,
					Resource: toStruct(voc.VirtualMachine{
						Compute: &voc.Compute{
							Resource: &voc.Resource{
								ID:   testdata.MockResourceID1,
								Type: []string{"VirtualMachine"}},
						},
					}, t),
				}},
			wantResp: "",
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "timestamp: value is required")
			},
		},
		{
			name: "Valid evidence",
			args: args{
				Evidence: &evidence.Evidence{
					Id:        testdata.MockCloudServiceID1,
					Timestamp: timestamppb.Now(),
					ToolId:    testdata.MockEvidenceToolID1,
					Resource: toStruct(voc.VirtualMachine{
						Compute: &voc.Compute{
							Resource: &voc.Resource{
								ID:   testdata.MockResourceID1,
								Type: []string{"VirtualMachine"}},
						},
					}, t),
					CloudServiceId: testdata.MockCloudServiceID1,
				}},
			wantResp: string(testdata.MockResourceID1),
			wantErr:  assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := ValidateWithResource(tt.args.Evidence)
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
