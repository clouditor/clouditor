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
	"clouditor.io/clouditor/internal/defaults"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
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
					MetricId: "TestMetric",
					Configuration: &assessment.MetricConfiguration{
						Operator:    "<",
						TargetValue: &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: "11111"}},
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
					MetricId:       "TestMetric",
					CloudServiceId: "00000000000000000000",
					Configuration: &assessment.MetricConfiguration{
						Operator:    "<",
						TargetValue: &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: "11111"}},
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
					CloudServiceId: "00000000-0000-0000-0000-000000000000",
					Configuration: &assessment.MetricConfiguration{
						Operator:    "<",
						TargetValue: &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: "11111"}},
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
					MetricId:       "TestMetric",
					CloudServiceId: "00000000-0000-0000-0000-000000000000",
					Configuration: &assessment.MetricConfiguration{
						Operator:       "<",
						TargetValue:    &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: "11111"}},
						MetricId:       "TestMetric",
						CloudServiceId: "00000000-0000-0000-0000-000000000000",
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
					MetricId: "TestMetric",
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
					MetricId:       "TestMetric",
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
					CloudServiceId: "00000000-0000-0000-0000-000000000000",
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
					MetricId:       "TestMetric",
					CloudServiceId: "00000000-0000-0000-0000-000000000000",
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

func TestTargetOfEvaluation_Validate(t *testing.T) {
	type fields struct {
		toe *TargetOfEvaluation
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "ToE missing",
			fields: fields{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrToEIsMissing.Error())
			},
		},
		{
			name: "AssuranceLevel missing",
			fields: fields{
				toe: &TargetOfEvaluation{
					CloudServiceId: defaults.DefaultTargetCloudServiceID,
					CatalogId:      defaults.DefaultCatalogID,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrAssuranceLevelIsMissing.Error())
			},
		},
		{
			name: "CatalogId missing",
			fields: fields{
				toe: &TargetOfEvaluation{
					CloudServiceId: defaults.DefaultTargetCloudServiceID,
					AssuranceLevel: &defaults.AssuranceLevelHigh,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrCatalogIDIsMissing.Error())
			},
		},
		{
			name: "CloudServiceId is missing",
			fields: fields{
				toe: &TargetOfEvaluation{
					CatalogId:      defaults.DefaultCatalogID,
					AssuranceLevel: &defaults.AssuranceLevelHigh,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, assessment.ErrCloudServiceIDIsMissing.Error())
			},
		},
		{
			name: "Happy path",
			fields: fields{
				toe: &TargetOfEvaluation{
					CloudServiceId: defaults.DefaultTargetCloudServiceID,
					CatalogId:      defaults.DefaultCatalogID,
					AssuranceLevel: &defaults.AssuranceLevelHigh,
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.fields.toe
			tt.wantErr(t, req.Validate())
		})
	}
}
