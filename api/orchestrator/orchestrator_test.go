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
	"testing"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/internal/testdata"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetricConfigurationRequest_Validate(t *testing.T) {
	type fields struct {
		Request *UpdateMetricConfigurationRequest
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Missing CloudServiceId",
			fields: fields{
				Request: &UpdateMetricConfigurationRequest{
					MetricId: testdata.MockMetricID,
					Configuration: &assessment.MetricConfiguration{
						Operator:    "<",
						TargetValue: testdata.MockMetricConfigurationTargetValueString,
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "CloudServiceId: value must be a valid UUID")
			},
		},
		{
			name: "Wrong CloudServiceId",
			fields: fields{
				Request: &UpdateMetricConfigurationRequest{
					MetricId:       testdata.MockMetricID,
					CloudServiceId: "00000000000000000000",
					Configuration: &assessment.MetricConfiguration{
						Operator:    "<",
						TargetValue: testdata.MockMetricConfigurationTargetValueString,
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "CloudServiceId: value must be a valid UUID")
			},
		},
		{
			name: "Missing MetricId",
			fields: fields{
				Request: &UpdateMetricConfigurationRequest{
					CloudServiceId: testdata.MockCloudServiceID,
					Configuration: &assessment.MetricConfiguration{
						Operator:    "<",
						TargetValue: testdata.MockMetricConfigurationTargetValueString,
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "MetricId: value length must be at least 1 runes")
			},
		},
		{
			name: "Without error",
			fields: fields{
				Request: &UpdateMetricConfigurationRequest{
					MetricId:       testdata.MockMetricID,
					CloudServiceId: testdata.MockCloudServiceID,
					Configuration: &assessment.MetricConfiguration{
						Operator:       "<",
						TargetValue:    testdata.MockMetricConfigurationTargetValueString,
						MetricId:       testdata.MockMetricID,
						CloudServiceId: testdata.MockCloudServiceID,
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.fields.Request
			tt.wantErr(t, req.Validate())
		})
	}
}

func TestGetMetricConfigurationRequest_Validate(t *testing.T) {
	type fields struct {
		Request *GetMetricConfigurationRequest
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Missing CloudServiceId",
			fields: fields{
				Request: &GetMetricConfigurationRequest{
					MetricId: testdata.MockMetricID,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "CloudServiceId: value must be a valid UUID")
			},
		},
		{
			name: "Wrong CloudServiceId",
			fields: fields{
				Request: &GetMetricConfigurationRequest{
					MetricId:       testdata.MockMetricID,
					CloudServiceId: "00000000000000000000",
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "CloudServiceId: value must be a valid UUID")
			},
		},
		{
			name: "Missing MetricId",
			fields: fields{
				Request: &GetMetricConfigurationRequest{
					CloudServiceId: testdata.MockCloudServiceID,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "MetricId: value length must be at least 1 runes")
			},
		},
		{
			name: "Without error",
			fields: fields{
				Request: &GetMetricConfigurationRequest{
					MetricId:       testdata.MockMetricID,
					CloudServiceId: testdata.MockCloudServiceID,
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.fields.Request
			tt.wantErr(t, req.Validate())
		})
	}
}

func TestAddControlToScopeRequest_GetCloudServiceId(t *testing.T) {
	type fields struct {
		Scope *ControlInScope
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy path",
			fields: fields{
				&ControlInScope{
					TargetOfEvaluationCloudServiceId: testdata.MockCloudServiceID,
				},
			},
			want: testdata.MockCloudServiceID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &AddControlToScopeRequest{
				Scope: tt.fields.Scope,
			}
			if got := req.GetCloudServiceId(); got != tt.want {
				t.Errorf("AddControlToScopeRequest.GetCloudServiceId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateControlInScopeRequest_GetCloudServiceId(t *testing.T) {
	type fields struct {
		Scope *ControlInScope
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy path",
			fields: fields{
				&ControlInScope{
					TargetOfEvaluationCloudServiceId: testdata.MockCloudServiceID,
				},
			},
			want: testdata.MockCloudServiceID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &UpdateControlInScopeRequest{
				Scope: tt.fields.Scope,
			}
			if got := req.GetCloudServiceId(); got != tt.want {
				t.Errorf("AddControlToScopeRequest.GetCloudServiceId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateCloudServiceRequest_GetCloudServiceId(t *testing.T) {
	type fields struct {
		CloudService *CloudService
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy path",
			fields: fields{
				&CloudService{
					Id: testdata.MockCloudServiceID,
				},
			},
			want: testdata.MockCloudServiceID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &UpdateCloudServiceRequest{
				CloudService: tt.fields.CloudService,
			}
			if got := req.GetCloudServiceId(); got != tt.want {
				t.Errorf("UpdateCloudServiceRequest.GetCloudServiceId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreAssessmentResultRequest_GetCloudServiceId(t *testing.T) {
	type fields struct {
		Result *assessment.AssessmentResult
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy path",
			fields: fields{
				Result: &assessment.AssessmentResult{
					CloudServiceId: testdata.MockCloudServiceID,
				},
			},
			want: testdata.MockCloudServiceID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &StoreAssessmentResultRequest{
				Result: tt.fields.Result,
			}
			if got := req.GetCloudServiceId(); got != tt.want {
				t.Errorf("StoreAssessmentResultRequest.GetCloudServiceId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateTargetOfEvaluationRequest_GetCloudServiceId(t *testing.T) {
	type fields struct {
		TargetOfEvaluation *TargetOfEvaluation
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy path",
			fields: fields{
				&TargetOfEvaluation{
					CloudServiceId: testdata.MockCloudServiceID,
				},
			},
			want: testdata.MockCloudServiceID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &CreateTargetOfEvaluationRequest{
				TargetOfEvaluation: tt.fields.TargetOfEvaluation,
			}
			if got := req.GetCloudServiceId(); got != tt.want {
				t.Errorf("CreateTargetOfEvaluationRequest.GetCloudServiceId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestControlInScope_TableName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Happy path",
			want: "controls_in_scope",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ControlInScope{}
			if got := c.TableName(); got != tt.want {
				t.Errorf("ControlInScope.TableName() = %v, want %v", got, tt.want)
			}
		})
	}
}
