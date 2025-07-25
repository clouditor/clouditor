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

package assessment

import (
	"testing"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestMetricConfiguration_Validate(t *testing.T) {
	type fields struct {
		MetricConfiguration *MetricConfiguration
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "MetricConfiguration Operator is empty",
			fields: fields{
				MetricConfiguration: &MetricConfiguration{
					TargetValue: testdata.MockMetricConfigurationTargetValueString,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "metric_id: value length must be at least 1 characters")
			},
		},
		{
			name: "MetricConfiguration TargetValue is empty",
			fields: fields{
				MetricConfiguration: &MetricConfiguration{
					Operator: "==",
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "target_value: value is required")
			},
		},
		{
			name: "Successful Validation",
			fields: fields{
				MetricConfiguration: &MetricConfiguration{
					TargetValue:          testdata.MockMetricConfigurationTargetValueString,
					Operator:             "==",
					MetricId:             testdata.MockMetricID1,
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields.MetricConfiguration
			tt.wantErr(t, api.Validate(c))
		})
	}
}

func TestMetricConfiguration_Hash(t *testing.T) {
	type fields struct {
		sizeCache            protoimpl.SizeCache
		unknownFields        protoimpl.UnknownFields
		Operator             string
		TargetValue          *structpb.Value
		IsDefault            bool
		UpdatedAt            *timestamppb.Timestamp
		MetricId             string
		TargetOfEvaluationId string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy path",
			fields: fields{
				Operator:    "<",
				TargetValue: &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: 5}},
			},
			want: "PC1udW1iZXJfdmFsdWU6NQ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &MetricConfiguration{
				sizeCache:            tt.fields.sizeCache,
				unknownFields:        tt.fields.unknownFields,
				Operator:             tt.fields.Operator,
				TargetValue:          tt.fields.TargetValue,
				IsDefault:            tt.fields.IsDefault,
				UpdatedAt:            tt.fields.UpdatedAt,
				MetricId:             tt.fields.MetricId,
				TargetOfEvaluationId: tt.fields.TargetOfEvaluationId,
			}
			if got := x.Hash(); got != tt.want {
				t.Errorf("MetricConfiguration.Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetric_CategoryID(t *testing.T) {
	type fields struct {
		Category string
	}
	type args struct {
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantID string
	}{
		{
			name: "happy path",
			fields: fields{
				Category: "LoggingMonitoring",
			},
			wantID: "LoggingMonitoring",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Metric{
				Category: tt.fields.Category,
			}
			if gotID := c.Category; gotID != tt.wantID {
				t.Errorf("Metric.Category = %v, want %v", gotID, tt.wantID)
			}
		})
	}
}
