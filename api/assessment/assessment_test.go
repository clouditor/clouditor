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
	"reflect"
	"testing"

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
	tests := []struct {
		name          string
		args          args
		wantResp      string
		wantRespError error
		wantErr       bool
	}{
		{
			name: "Missing assessment result id",
			args: args{
				&AssessmentResult{
					// Empty id
					Id:       "",
					MetricId: "MockMetricID",
					MetricConfiguration: &MetricConfiguration{
						Operator: "MockOperator",
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
			// cannot access unexported invalidLengthError of uuid package. Use the error string directly
			wantRespError: ErrIdInvalidFormat,
			wantErr:       true,
		},
		{
			name: "Wrong length of assessment result id",
			args: args{
				&AssessmentResult{
					// Only 4 characters
					Id:       "1234",
					MetricId: "MockMetricID",
					MetricConfiguration: &MetricConfiguration{
						Operator: "MockOperator",
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
			// cannot access unexported invalidLengthError of uuid package. Use the error string directly
			wantRespError: ErrIdInvalidFormat,
			wantErr:       true,
		},
		{
			name: "Wrong format of assessment result id",
			args: args{
				&AssessmentResult{
					// Wrong format: 'x' not allowed (no hexadecimal character)
					Id:       "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
					MetricId: "MockMetricID",
					MetricConfiguration: &MetricConfiguration{
						Operator: "MockOperator",
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
			// Copied error of uuid package.
			wantRespError: ErrIdInvalidFormat,
			wantErr:       true,
		},
		{
			name: "Missing assessment result timestamp",
			args: args{
				&AssessmentResult{
					Id:       assessmentResultID,
					MetricId: "MockMetricID",
					MetricConfiguration: &MetricConfiguration{
						Operator: "MockOperator",
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
			wantResp:      "",
			wantRespError: ErrTimestampMissing,
			wantErr:       true,
		},
		{
			name: "Missing assessment result metric id",
			args: args{
				&AssessmentResult{
					Id:        assessmentResultID,
					Timestamp: timestamppb.Now(),
					MetricConfiguration: &MetricConfiguration{
						Operator: "MockOperator",
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
			wantResp:      "",
			wantRespError: ErrMetricIdMissing,
			wantErr:       true,
		},
		{
			name: "Missing assessment result resource types",
			args: args{
				&AssessmentResult{
					Id:        assessmentResultID,
					Timestamp: timestamppb.Now(),
					MetricId:  "MockMetricID",
					MetricConfiguration: &MetricConfiguration{
						Operator: "MockOperator",
						TargetValue: &structpb.Value{
							Kind: &structpb.Value_StringValue{
								StringValue: "MockTargetValue",
							},
						},
					},
					ResourceId: "myResource",
					EvidenceId: mockEvidenceID,
				},
			},
			wantResp:      "",
			wantRespError: ErrResourceTypesMissing,
			wantErr:       true,
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
			wantResp:      "",
			wantRespError: ErrMetricConfigurationMissing,
			wantErr:       true,
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
					},
					EvidenceId:    mockEvidenceID,
					ResourceTypes: []string{"Resource"},
				},
			},
			wantResp:      "",
			wantRespError: ErrMetricConfigurationOperatorMissing,
			wantErr:       true,
		},
		{
			name: "Missing assessment result metric configuration target value",
			args: args{
				&AssessmentResult{
					Id:        assessmentResultID,
					Timestamp: timestamppb.Now(),
					MetricId:  "MockMetricID",
					MetricConfiguration: &MetricConfiguration{
						Operator: "MockOperator",
					},
					EvidenceId:    mockEvidenceID,
					ResourceTypes: []string{"Resource"},
				},
			},
			wantResp:      "",
			wantRespError: ErrMetricConfigurationTargetValueMissing,
			wantErr:       true,
		},
		{
			name: "Missing assessment result evidence id",
			args: args{
				&AssessmentResult{
					Id:        assessmentResultID,
					Timestamp: timestamppb.Now(),
					MetricId:  "MockMetricID",
					MetricConfiguration: &MetricConfiguration{
						Operator: "MockOperator",
						TargetValue: &structpb.Value{
							Kind: &structpb.Value_StringValue{
								StringValue: "MockTargetValue",
							},
						},
					},
					ResourceTypes: []string{"Resource"},
				},
			},
			wantResp:      "",
			wantRespError: ErrEvidenceIdInvalidFormat,
			wantErr:       true,
		},
		{
			name: "Valid assessment result",
			args: args{
				&AssessmentResult{
					Id:        assessmentResultID,
					Timestamp: timestamppb.Now(),
					MetricId:  "MockMetricID",
					MetricConfiguration: &MetricConfiguration{
						Operator: "MockOperator",
						TargetValue: &structpb.Value{
							Kind: &structpb.Value_StringValue{
								StringValue: "MockTargetValue",
							},
						},
					},
					EvidenceId:    mockEvidenceID,
					ResourceId:    "myResource",
					ResourceTypes: []string{"Resource"},
				},
			},
			wantResp:      "",
			wantRespError: nil,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resourceId, err := tt.args.AssessmentResult.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(resourceId, tt.wantResp) {
				t.Errorf("Validate() gotResp = %v, want %v", resourceId, tt.wantResp)
			}
			assert.Equal(t, resourceId, tt.wantResp)

			if err != nil {
				assert.ErrorIs(t, err, tt.wantRespError)
			}
		})
	}
}
