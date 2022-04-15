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
	"errors"
	"os"
	"reflect"
	"testing"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/persistence"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

var ErrSomeError = errors.New("some error")

const MockMetricID = "SomeMetric"

func TestService_CreateMetric(t *testing.T) {
	type args struct {
		in0 context.Context
		req *orchestrator.CreateMetricRequest
	}
	tests := []struct {
		name           string
		args           args
		wantMetric     *assessment.Metric
		wantErr        bool
		wantErrMessage string
	}{
		{
			name: "Create valid metric",
			args: args{
				context.TODO(),
				&orchestrator.CreateMetricRequest{
					Metric: &assessment.Metric{
						Id:   "MyTransportEncryptionEnabled",
						Name: "A very good metric",
					},
				},
			},
			wantMetric: &assessment.Metric{
				Id:   "MyTransportEncryptionEnabled",
				Name: "A very good metric",
			},
			wantErr: false,
		},
		{
			name: "Create invalid metric",
			args: args{
				context.TODO(),
				&orchestrator.CreateMetricRequest{
					Metric: &assessment.Metric{},
				},
			},
			wantMetric:     nil,
			wantErr:        true,
			wantErrMessage: "rpc error: code = InvalidArgument desc = validation of metric failed",
		},
		{
			name: "Create metric which already exists",
			args: args{
				context.TODO(),
				&orchestrator.CreateMetricRequest{
					Metric: &assessment.Metric{
						Id:       "TLSVersion",
						Name:     "TLSMetricMockName",
						Category: "",
						Scale:    assessment.Metric_NOMINAL,
					},
				},
			},
			wantMetric:     nil,
			wantErr:        true,
			wantErrMessage: "rpc error: code = AlreadyExists desc = metric already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService()
			gotMetric, err := service.CreateMetric(tt.args.in0, tt.args.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("Service.CreateMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				assert.Equal(t, tt.wantErrMessage, err.Error())
			}

			if !proto.Equal(gotMetric, tt.wantMetric) {
				t.Errorf("Service.CreateMetric() = %v, want %v", gotMetric, tt.wantMetric)
			}

		})
	}
}

func TestService_UpdateMetric(t *testing.T) {
	type args struct {
		in0 context.Context
		req *orchestrator.UpdateMetricRequest
	}
	tests := []struct {
		name           string
		args           args
		wantMetric     *assessment.Metric
		wantErr        bool
		wantErrMessage string
	}{
		{
			name: "Update existing metric",
			args: args{
				context.TODO(),
				&orchestrator.UpdateMetricRequest{
					MetricId: "TransportEncryptionEnabled",
					Metric: &assessment.Metric{
						Name: "A slightly updated metric",
					},
				},
			},
			wantMetric: &assessment.Metric{
				Id:   "TransportEncryptionEnabled",
				Name: "A slightly updated metric",
			},
			wantErr: false,
		},
		{
			name: "Update non-existing metric",
			args: args{
				context.TODO(),
				&orchestrator.UpdateMetricRequest{
					MetricId: "DoesProbablyNotExist",
					Metric: &assessment.Metric{
						Id:   "UpdateMetricID",
						Name: "UpdateMetricName",
					},
				},
			},
			wantMetric:     nil,
			wantErr:        true,
			wantErrMessage: "rpc error: code = NotFound desc = metric with identifier DoesProbablyNotExist does not exist",
		},
		{
			name: "Updating invalid metric",
			args: args{
				context.TODO(),
				&orchestrator.UpdateMetricRequest{
					MetricId: "DoesProbablyNotExist",
					Metric:   &assessment.Metric{},
				},
			},
			wantMetric:     nil,
			wantErr:        true,
			wantErrMessage: "rpc error: code = InvalidArgument desc = validation of metric failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService()
			gotMetric, err := service.UpdateMetric(tt.args.in0, tt.args.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("Service.UpdateMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				assert.Contains(t, err.Error(), tt.wantErrMessage)
			}

			if !proto.Equal(gotMetric, tt.wantMetric) {
				t.Errorf("Service.UpdateMetric() = %v, want %v", gotMetric, tt.wantMetric)
			}
		})
	}
}

func TestService_GetMetric(t *testing.T) {
	type args struct {
		in0 context.Context
		req *orchestrator.GetMetricRequest
	}
	tests := []struct {
		name       string
		args       args
		wantMetric *assessment.Metric
		wantErr    bool
	}{
		{
			name: "Get existing metric",
			args: args{
				context.TODO(),
				&orchestrator.GetMetricRequest{
					MetricId: "TransportEncryptionEnabled",
				},
			},
			wantMetric: &assessment.Metric{
				Id:          "TransportEncryptionEnabled",
				Name:        "Transport Encryption: Enabled",
				Description: "This metric describes, whether transport encryption is turned on or not",
				Scale:       assessment.Metric_ORDINAL,
				Range: &assessment.Range{
					Range: &assessment.Range_AllowedValues{AllowedValues: &assessment.AllowedValues{
						Values: []*structpb.Value{
							structpb.NewBoolValue(false),
							structpb.NewBoolValue(true),
						}}}},
			},
			wantErr: false,
		},
		{
			name: "Get non-existing metric",
			args: args{
				context.TODO(),
				&orchestrator.GetMetricRequest{
					MetricId: "DoesProbablyNotExist",
				},
			},
			wantMetric: nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService()
			gotMetric, err := service.GetMetric(tt.args.in0, tt.args.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !proto.Equal(gotMetric, tt.wantMetric) {
				t.Errorf("Service.GetMetric() = %v, want %v", gotMetric, tt.wantMetric)
			}
		})
	}
}

func TestService_ListMetrics(t *testing.T) {
	var (
		response *orchestrator.ListMetricsResponse
		err      error
	)
	service := NewService()

	response, err = service.ListMetrics(context.TODO(), &orchestrator.ListMetricsRequest{})

	assert.NoError(t, err)
	assert.NotEmpty(t, response.Metrics)
}

func TestService_GetMetricImplementation(t *testing.T) {
	type fields struct {
		metricConfigurations  map[string]map[string]*assessment.MetricConfiguration
		results               map[string]*assessment.AssessmentResult
		AssessmentResultHooks []func(result *assessment.AssessmentResult, err error)
		storage               persistence.Storage
		metricsFile           string
		requirements          []*orchestrator.Requirement
		events                chan *orchestrator.MetricChangeEvent
	}
	type args struct {
		ctx context.Context
		req *orchestrator.GetMetricImplementationRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *assessment.MetricImplementation
		wantErr bool
	}{
		{
			name: "metric not found",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
			},
			args: args{
				req: &orchestrator.GetMetricImplementationRequest{
					MetricId: MockMetricID,
				},
			},
			wantErr: true,
		},
		{
			name: "storage error",
			fields: fields{
				storage: &testutil.StorageWithError{GetErr: ErrSomeError},
			},
			args: args{
				req: &orchestrator.GetMetricImplementationRequest{
					MetricId: MockMetricID,
				},
			},
			wantErr: true,
		},
		{
			name: "metric found",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Save(&assessment.MetricImplementation{
						MetricId: MockMetricID,
						Lang:     assessment.MetricImplementation_REGO,
						Code:     "package test",
					}, "metric_id = ?", MockMetricID)
					if err != nil {
						t.Errorf("Could not save implementation: %v", err)
					}
				}),
			},
			args: args{
				req: &orchestrator.GetMetricImplementationRequest{
					MetricId: MockMetricID,
				},
			},
			wantRes: &assessment.MetricImplementation{
				MetricId: MockMetricID,
				Lang:     assessment.MetricImplementation_REGO,
				Code:     "package test",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				metricConfigurations:  tt.fields.metricConfigurations,
				results:               tt.fields.results,
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				metricsFile:           tt.fields.metricsFile,
				requirements:          tt.fields.requirements,
				events:                tt.fields.events,
			}

			gotRes, err := svc.GetMetricImplementation(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetMetricImplementation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.GetMetricImplementation() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_loadMetrics(t *testing.T) {
	type fields struct {
		metricConfigurations  map[string]map[string]*assessment.MetricConfiguration
		results               map[string]*assessment.AssessmentResult
		AssessmentResultHooks []func(result *assessment.AssessmentResult, err error)
		storage               persistence.Storage
		metricsFile           string
		requirements          []*orchestrator.Requirement
		events                chan *orchestrator.MetricChangeEvent
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "json not found",
			fields: fields{
				metricsFile: "notfound.json",
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, os.ErrNotExist)
			},
		},
		{
			name: "storage error",
			fields: fields{
				metricsFile: "metrics.json",
				storage:     &testutil.StorageWithError{SaveErr: ErrSomeError},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrSomeError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				metricConfigurations:  tt.fields.metricConfigurations,
				results:               tt.fields.results,
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				metricsFile:           tt.fields.metricsFile,
				requirements:          tt.fields.requirements,
				events:                tt.fields.events,
			}

			err := svc.loadMetrics()
			if tt.wantErr != nil {
				tt.wantErr(t, err)
			}
		})
	}
}

func TestService_UpdateMetricImplementation(t *testing.T) {
	type fields struct {
		metricConfigurations  map[string]map[string]*assessment.MetricConfiguration
		results               map[string]*assessment.AssessmentResult
		AssessmentResultHooks []func(result *assessment.AssessmentResult, err error)
		storage               persistence.Storage
		metricsFile           string
		metrics               map[string]*assessment.Metric
		requirements          []*orchestrator.Requirement
		events                chan *orchestrator.MetricChangeEvent
	}
	type args struct {
		in0 context.Context
		req *orchestrator.UpdateMetricImplementationRequest
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantImpl *assessment.MetricImplementation
		wantErr  bool
	}{
		{
			name: "metric not found",
			fields: fields{
				storage:     testutil.NewInMemoryStorage(t),
				metricsFile: "metrics.json",
			},
			args: args{
				req: &orchestrator.UpdateMetricImplementationRequest{
					MetricId: "notfound",
				},
			},
			wantErr: true,
		},
		{
			name: "storage error",
			fields: fields{
				storage: &testutil.StorageWithError{SaveErr: ErrSomeError},
			},
			args: args{
				req: &orchestrator.UpdateMetricImplementationRequest{
					MetricId: "TransportEncryptionEnabled",
					Implementation: &assessment.MetricImplementation{
						Lang: assessment.MetricImplementation_REGO,
						Code: "package example",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "update",
			fields: fields{
				storage:     testutil.NewInMemoryStorage(t),
				metricsFile: "metrics.json",
				metrics: map[string]*assessment.Metric{
					"TransportEncryptionEnabled": {},
				},
			},
			args: args{
				req: &orchestrator.UpdateMetricImplementationRequest{
					MetricId: "TransportEncryptionEnabled",
					Implementation: &assessment.MetricImplementation{
						Lang: assessment.MetricImplementation_REGO,
						Code: "package example",
					},
				},
			},
			wantImpl: &assessment.MetricImplementation{
				MetricId: "TransportEncryptionEnabled",
				Lang:     assessment.MetricImplementation_REGO,
				Code:     "package example",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				metricConfigurations:  tt.fields.metricConfigurations,
				results:               tt.fields.results,
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				metricsFile:           tt.fields.metricsFile,
				metrics:               tt.fields.metrics,
				requirements:          tt.fields.requirements,
				events:                tt.fields.events,
			}
			gotImpl, err := svc.UpdateMetricImplementation(tt.args.in0, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.UpdateMetricImplementation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotImpl, tt.wantImpl) {
				t.Errorf("Service.UpdateMetricImplementation() = %v, want %v", gotImpl, tt.wantImpl)
			}
		})
	}
}

func Test_loadMetricImplementation(t *testing.T) {
	type args struct {
		metricID string
		file     string
	}
	tests := []struct {
		name     string
		args     args
		wantImpl *assessment.MetricImplementation
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "not existing",
			args: args{
				metricID: MockMetricID,
				file:     "doesnotexist.rego",
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, os.ErrNotExist)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotImpl, err := loadMetricImplementation(tt.args.metricID, tt.args.file)
			if tt.wantErr != nil {
				tt.wantErr(t, err, tt.args)
			}
			if !reflect.DeepEqual(gotImpl, tt.wantImpl) {
				t.Errorf("loadMetricImplementation() = %v, want %v", gotImpl, tt.wantImpl)
			}
		})
	}
}

func TestService_GetMetricConfiguration(t *testing.T) {
	type fields struct {
		metricConfigurations  map[string]map[string]*assessment.MetricConfiguration
		results               map[string]*assessment.AssessmentResult
		AssessmentResultHooks []func(result *assessment.AssessmentResult, err error)
		storage               persistence.Storage
		metrics               map[string]*assessment.Metric
		metricsFile           string
		requirements          []*orchestrator.Requirement
		events                chan *orchestrator.MetricChangeEvent
	}
	type args struct {
		in0 context.Context
		req *orchestrator.GetMetricConfigurationRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantResponse *assessment.MetricConfiguration
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name:   "metric not found",
			fields: fields{},
			args: args{
				req: &orchestrator.GetMetricConfigurationRequest{
					MetricId: "NotExists",
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				status, ok := status.FromError(err)
				if !ok {
					return false
				}
				return assert.Equal(t, status.Code(), codes.NotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				metricConfigurations:  tt.fields.metricConfigurations,
				results:               tt.fields.results,
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				metrics:               tt.fields.metrics,
				metricsFile:           tt.fields.metricsFile,
				requirements:          tt.fields.requirements,
				events:                tt.fields.events,
			}
			gotResponse, err := s.GetMetricConfiguration(tt.args.in0, tt.args.req)
			if tt.wantErr != nil {
				tt.wantErr(t, err, tt.args)
			}
			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("Service.GetMetricConfiguration() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestService_ListMetricConfigurations(t *testing.T) {
	type fields struct {
		metricConfigurations  map[string]map[string]*assessment.MetricConfiguration
		results               map[string]*assessment.AssessmentResult
		AssessmentResultHooks []func(result *assessment.AssessmentResult, err error)
		storage               persistence.Storage
		metrics               map[string]*assessment.Metric
		metricsFile           string
		requirements          []*orchestrator.Requirement
		events                chan *orchestrator.MetricChangeEvent
	}
	type args struct {
		ctx context.Context
		req *orchestrator.ListMetricConfigurationRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantResponse *orchestrator.ListMetricConfigurationResponse
		wantErr      bool
	}{
		{
			name: "error",
			fields: fields{
				metrics: map[string]*assessment.Metric{
					MockMetricID: {},
				},
			},
			args: args{
				req: &orchestrator.ListMetricConfigurationRequest{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				metricConfigurations:  tt.fields.metricConfigurations,
				results:               tt.fields.results,
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				metrics:               tt.fields.metrics,
				metricsFile:           tt.fields.metricsFile,
				requirements:          tt.fields.requirements,
				events:                tt.fields.events,
			}
			gotResponse, err := svc.ListMetricConfigurations(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ListMetricConfigurations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("Service.ListMetricConfigurations() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}
