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
	"clouditor.io/clouditor/internal/testvariables"
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
					MetricId: testvariables.MockMetricID,
					Configuration: &assessment.MetricConfiguration{
						Operator:    "<",
						TargetValue: testvariables.MockMetricConfigurationTargetValueString,
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
					MetricId:       testvariables.MockMetricID,
					CloudServiceId: "00000000000000000000",
					Configuration: &assessment.MetricConfiguration{
						Operator:    "<",
						TargetValue: testvariables.MockMetricConfigurationTargetValueString,
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
					CloudServiceId: testvariables.MockCloudServiceID,
					Configuration: &assessment.MetricConfiguration{
						Operator:    "<",
						TargetValue: testvariables.MockMetricConfigurationTargetValueString,
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
					MetricId:       testvariables.MockMetricID,
					CloudServiceId: testvariables.MockCloudServiceID,
					Configuration: &assessment.MetricConfiguration{
						Operator:       "<",
						TargetValue:    testvariables.MockMetricConfigurationTargetValueString,
						MetricId:       testvariables.MockMetricID,
						CloudServiceId: testvariables.MockCloudServiceID,
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
					MetricId: testvariables.MockMetricID,
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
					MetricId:       testvariables.MockMetricID,
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
					CloudServiceId: testvariables.MockCloudServiceID,
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
					MetricId:       testvariables.MockMetricID,
					CloudServiceId: testvariables.MockCloudServiceID,
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
