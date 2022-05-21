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
	"context"
	"reflect"
	"testing"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/persistence"
)

func TestLoadRequirements(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name             string
		args             args
		wantRequirements []*orchestrator.Requirement
		wantErr          bool
	}{
		{
			name: "load",
			args: args{
				file: "requirements.json",
			},
			wantRequirements: []*orchestrator.Requirement{
				{
					Id:          "Req-1",
					Name:        "Make-it-Secure",
					Category:    "Some Category",
					Description: "You should make everything secure",
					Metrics: []*assessment.Metric{
						{Id: "TransportEncryptionEnabled"},
						{Id: "TransportEncryptionAlgorithm"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRequirements, err := LoadRequirements(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadRequirements() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRequirements, tt.wantRequirements) {
				t.Errorf("LoadRequirements() = %v, want %v", gotRequirements, tt.wantRequirements)
			}
		})
	}
}

func TestService_ListRequirements(t *testing.T) {
	type fields struct {
		metricConfigurations  map[string]map[string]*assessment.MetricConfiguration
		cloudServiceHooks     []orchestrator.CloudServiceHookFunc
		results               map[string]*assessment.AssessmentResult
		AssessmentResultHooks []func(result *assessment.AssessmentResult, err error)
		storage               persistence.Storage
		metricsFile           string
		loadMetricsFunc       func() ([]*assessment.Metric, error)
		requirements          []*orchestrator.Requirement
		events                chan *orchestrator.MetricChangeEvent
	}
	type args struct {
		in0 context.Context
		req *orchestrator.ListRequirementsRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *orchestrator.ListRequirementsResponse
		wantErr bool
	}{
		{
			name: "list requirements",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					_ = s.Save(&assessment.Metric{Id: "Metric1", Name: "Metric1"})
					_ = s.Save(&assessment.Metric{Id: "Metric2", Name: "Metric2"})
					_ = s.Save(&orchestrator.Requirement{Id: "Req1", Metrics: []*assessment.Metric{{Id: "Metric1"}}})
				}),
			},
			args: args{req: &orchestrator.ListRequirementsRequest{}},
			wantRes: &orchestrator.ListRequirementsResponse{
				Requirements: []*orchestrator.Requirement{
					{
						Id: "Req1",
						Metrics: []*assessment.Metric{
							{Id: "Metric1", Name: "Metric1"},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				metricConfigurations:  tt.fields.metricConfigurations,
				cloudServiceHooks:     tt.fields.cloudServiceHooks,
				results:               tt.fields.results,
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				metricsFile:           tt.fields.metricsFile,
				loadMetricsFunc:       tt.fields.loadMetricsFunc,
				requirements:          tt.fields.requirements,
				events:                tt.fields.events,
			}
			gotRes, err := svc.ListRequirements(tt.args.in0, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ListRequirements() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.ListRequirements() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_loadRequirements(t *testing.T) {
	type fields struct {
		metricConfigurations  map[string]map[string]*assessment.MetricConfiguration
		cloudServiceHooks     []orchestrator.CloudServiceHookFunc
		results               map[string]*assessment.AssessmentResult
		AssessmentResultHooks []func(result *assessment.AssessmentResult, err error)
		storage               persistence.Storage
		metricsFile           string
		loadMetricsFunc       func() ([]*assessment.Metric, error)
		requirements          []*orchestrator.Requirement
		events                chan *orchestrator.MetricChangeEvent
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "save error",
			fields: fields{
				storage: &testutil.StorageWithError{SaveErr: ErrSomeError},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				metricConfigurations:  tt.fields.metricConfigurations,
				cloudServiceHooks:     tt.fields.cloudServiceHooks,
				results:               tt.fields.results,
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				metricsFile:           tt.fields.metricsFile,
				loadMetricsFunc:       tt.fields.loadMetricsFunc,
				requirements:          tt.fields.requirements,
				events:                tt.fields.events,
			}
			if err := svc.loadRequirements(); (err != nil) != tt.wantErr {
				t.Errorf("Service.loadRequirements() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
