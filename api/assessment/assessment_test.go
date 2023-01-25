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

package assessment

import (
	"testing"

	"clouditor.io/clouditor/internal/util"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_ValidateAssessmentResult(t *testing.T) {
	type args struct {
		AssessmentResult *AssessmentResult
	}

	const assessmentResultID = "11111111-1111-1111-1111-111111111111"
	const mockEvidenceID = "11111111-1111-1111-1111-111111111111"
	const mockServiceID = "11111111-1111-1111-1111-111111111111"
	tests := []struct {
		name     string
		args     args
		wantResp string
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "Missing assessment result id",
			args: args{
				&AssessmentResult{
					// Empty id
					Id:       "",
					MetricId: "MockMetricID",
					MetricConfiguration: &MetricConfiguration{
						Operator: "==",
						TargetValue: &structpb.Value{
							Kind: &structpb.Value_StringValue{
								StringValue: "MockTargetValue",
							},
						},
					},
					EvidenceId:    mockEvidenceID,
					ResourceTypes: []string{"Resource"},
				},
			},
			wantResp: "",
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid uuid format")
			},
		},
		{
			name: "Wrong length of assessment result id",
			args: args{
				&AssessmentResult{
					// Only 4 characters
					Id:       "1234",
					MetricId: "MockMetricID",
					MetricConfiguration: &MetricConfiguration{
						Operator: "==",
						TargetValue: &structpb.Value{
							Kind: &structpb.Value_StringValue{
								StringValue: "MockTargetValue",
							},
						},
					},
					EvidenceId:    mockEvidenceID,
					ResourceTypes: []string{"Resource"},
				},
			},
			wantResp: "",
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid uuid format")
			},
		},
		{
			name: "Wrong format of assessment result id",
			args: args{
				&AssessmentResult{
					// Wrong format: 'x' not allowed (no hexadecimal character)
					Id:       "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
					MetricId: "MockMetricID",
					MetricConfiguration: &MetricConfiguration{
						Operator: "==",
						TargetValue: &structpb.Value{
							Kind: &structpb.Value_StringValue{
								StringValue: "MockTargetValue",
							},
						},
					},
					EvidenceId:    mockEvidenceID,
					ResourceTypes: []string{"Resource"},
				},
			},
			wantResp: "",
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid uuid format")
			},
		},
		{
			name: "Missing assessment result timestamp",
			args: args{
				&AssessmentResult{
					Id:       assessmentResultID,
					MetricId: "MockMetricID",
					MetricConfiguration: &MetricConfiguration{
						Operator: ">",
						TargetValue: &structpb.Value{
							Kind: &structpb.Value_StringValue{
								StringValue: "MockTargetValue",
							},
						},
					},
					EvidenceId:    mockEvidenceID,
					ResourceTypes: []string{"Resource"},
				},
			},
			wantResp: "",
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "value is required")
			},
		},
		{
			name: "Missing assessment result metric id",
			args: args{
				&AssessmentResult{
					Id:        assessmentResultID,
					Timestamp: timestamppb.Now(),
					MetricConfiguration: &MetricConfiguration{
						Operator: "<",
						TargetValue: &structpb.Value{
							Kind: &structpb.Value_StringValue{
								StringValue: "MockTargetValue",
							},
						},
					},
					ResourceTypes: []string{"Resource"},
					EvidenceId:    mockEvidenceID,
				},
			},
			wantResp: "",
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "value length must be at least 1 runes")
			},
		},
		{
			name: "Missing assessment result resource types",
			args: args{
				&AssessmentResult{
					Id:        assessmentResultID,
					Timestamp: timestamppb.Now(),
					MetricId:  "MockMetricID",
					MetricConfiguration: &MetricConfiguration{
						MetricId: "MockMetricID",
						Operator: "==",
						TargetValue: &structpb.Value{
							Kind: &structpb.Value_StringValue{
								StringValue: "MockTargetValue",
							},
						},
						CloudServiceId: mockServiceID,
					},
					ResourceId: "myResource",
					EvidenceId: mockEvidenceID,
				},
			},
			wantResp: "",
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "ResourceTypes: value must contain at least 1 item(s")
			},
		},
		{
			name: "Missing assessment result metric configuration",
			args: args{
				&AssessmentResult{
					Id:            assessmentResultID,
					Timestamp:     timestamppb.Now(),
					MetricId:      "MockMetricID",
					EvidenceId:    mockEvidenceID,
					ResourceTypes: []string{"Resource"},
				},
			},
			wantResp: "",
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "MetricConfiguration: value is required")
			},
		},
		{
			name: "Missing assessment result metric configuration operator",
			args: args{
				&AssessmentResult{
					Id:        assessmentResultID,
					Timestamp: timestamppb.Now(),
					MetricId:  "MockMetricID",
					MetricConfiguration: &MetricConfiguration{
						TargetValue: &structpb.Value{
							Kind: &structpb.Value_StringValue{
								StringValue: "MockTargetValue",
							},
						},
						MetricId: "MockMetricID",
					},
					EvidenceId:    mockEvidenceID,
					ResourceTypes: []string{"Resource"},
				},
			},
			wantResp: "",
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid uuid format")
			},
		},
		{
			name: "Missing assessment result metric configuration target value",
			args: args{
				&AssessmentResult{
					Id:        assessmentResultID,
					Timestamp: timestamppb.Now(),
					MetricId:  "MockMetricID",
					MetricConfiguration: &MetricConfiguration{
						Operator: "<",
					},
					EvidenceId:    mockEvidenceID,
					ResourceTypes: []string{"Resource"},
				},
			},
			wantResp: "",
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "TargetValue: value is required")
			},
		},
		{
			name: "Missing assessment result evidence id",
			args: args{
				&AssessmentResult{
					Id:        assessmentResultID,
					Timestamp: timestamppb.Now(),
					MetricId:  "MockMetricID",
					MetricConfiguration: &MetricConfiguration{
						Operator: ">",
						MetricId: "MockMetricID",
						TargetValue: &structpb.Value{
							Kind: &structpb.Value_StringValue{
								StringValue: "MockTargetValue",
							},
						},
					},
					ResourceTypes: []string{"Resource"},
				},
			},
			wantResp: "",
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid uuid format")
			},
		},
		{
			name: "Valid assessment result",
			args: args{
				&AssessmentResult{
					Id:        assessmentResultID,
					Timestamp: timestamppb.Now(),
					MetricId:  "MockMetricID",
					MetricConfiguration: &MetricConfiguration{
						Operator: "==",
						MetricId: "MockMetricID",
						TargetValue: &structpb.Value{
							Kind: &structpb.Value_StringValue{
								StringValue: "MockTargetValue",
							},
						},
						CloudServiceId: mockServiceID,
					},
					EvidenceId:     mockEvidenceID,
					ResourceId:     "myResource",
					ResourceTypes:  []string{"Resource"},
					CloudServiceId: mockServiceID,
				},
			},
			wantResp: "",
			wantErr:  assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.AssessmentResult.Validate()
			if tt.wantErr != nil {
				tt.wantErr(t, err, tt.args)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListAssessmentResultsRequest_Validate(t *testing.T) {
	type fields struct {
		req *ListAssessmentResultsRequest
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Request is empty",
			fields: fields{
				&ListAssessmentResultsRequest{},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Invalid cloud service id",
			fields: fields{
				req: &ListAssessmentResultsRequest{
					FilteredCloudServiceId: util.Ref("invalidCloudServiceId"),
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "FilteredCloudServiceId: value must be a valid UUID")
			},
		},
		{
			name: "No filtered cloud service id",
			fields: fields{
				req: &ListAssessmentResultsRequest{},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path",
			fields: fields{
				req: &ListAssessmentResultsRequest{
					FilteredCloudServiceId: util.Ref("00000000-0000-0000-0000-000000000000"),
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.req.Validate()
			tt.wantErr(t, err)
		})
	}
}
