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

package evaluation

import (
	"testing"

	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/defaults"
	"github.com/stretchr/testify/assert"
)

func TestStartEvaluationRequest_Validate(t *testing.T) {
	type fields struct {
		Request *StartEvaluationRequest
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "Request is empty",
			fields: fields{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrRequestIsEmpty.Error())
			},
		},
		{
			name: "ToE is missing in request",
			fields: fields{
				Request: &StartEvaluationRequest{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "TODO")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				Request: &StartEvaluationRequest{
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CloudServiceId: defaults.DefaultTargetCloudServiceID,
						CatalogId:      defaults.DefaultCatalogID,
						AssuranceLevel: &defaults.AssuranceLevelHigh,
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

func TestStopEvaluationRequest_Validate(t *testing.T) {
	type fields struct {
		Request *StopEvaluationRequest
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "Request is empty",
			fields: fields{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrRequestIsEmpty.Error())
			},
		},
		{
			name: "Missing ControlID in request",
			fields: fields{
				Request: &StopEvaluationRequest{
					CategoryName: defaults.DefaultEUCSCategoryName,
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CloudServiceId: defaults.DefaultTargetCloudServiceID,
						CatalogId:      defaults.DefaultCatalogID,
						AssuranceLevel: &defaults.AssuranceLevelHigh,
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrControlIDIsMissing.Error())
			},
		},
		{
			name: "Missing CategoryName in request",
			fields: fields{
				Request: &StopEvaluationRequest{
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CloudServiceId: defaults.DefaultTargetCloudServiceID,
						CatalogId:      defaults.DefaultCatalogID,
						AssuranceLevel: &defaults.AssuranceLevelHigh,
					},
					ControlId: defaults.DefaultEUCSLowerLevelControlID137,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrCategoryNameIsMissing.Error())
			},
		},
		{
			name: "ToE is missing in request",
			fields: fields{
				Request: &StopEvaluationRequest{
					ControlId:    defaults.DefaultEUCSLowerLevelControlID137,
					CategoryName: defaults.DefaultEUCSCategoryName,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "TODO")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				Request: &StopEvaluationRequest{
					ControlId: defaults.DefaultEUCSLowerLevelControlID137,
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CloudServiceId: defaults.DefaultTargetCloudServiceID,
						CatalogId:      defaults.DefaultCatalogID,
						AssuranceLevel: &defaults.AssuranceLevelHigh,
					},
					CategoryName: defaults.DefaultEUCSCategoryName,
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
