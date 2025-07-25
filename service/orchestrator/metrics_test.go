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
	filepath "path/filepath"
	"testing"
	"time"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/orchestratortest"
	"clouditor.io/clouditor/v2/persistence"
	persistence_gorm "clouditor.io/clouditor/v2/persistence/gorm"
	"clouditor.io/clouditor/v2/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

var (
	ErrSomeError = errors.New("some error")

	// We need to define the following vars here because we could get import cycle errors in ./internal/testdata/testdata.go
	MockMetricConfiguration1 = &assessment.MetricConfiguration{
		MetricId:             testdata.MockMetricID1,
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
		Operator:             "==",
		TargetValue:          testdata.MockMetricConfigurationTargetValueString,
	}

	MockMetricConfiguration2 = &assessment.MetricConfiguration{
		MetricId:             testdata.MockMetricID2,
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID2,
		Operator:             "==",
		TargetValue:          testdata.MockMetricConfigurationTargetValueString,
	}

	MockMetric1 = &assessment.Metric{
		Id:          testdata.MockMetricID1,
		Description: testdata.MockMetricDescription1,
		Category:    testdata.MockMetricCategory1,
		Version:     "1.0",
		Comments:    "Test comments",
	}
)

func Test_loadMetricsFromMetricsRepository(t *testing.T) {
	// Create a temporary directory structure for test files
	validDir := filepath.Join("internal", "testdata", "mock_metrics", "valid_metrics")
	invalidDir := filepath.Join("internal", "testdata", "mock_metrics", "invalid_metrics")
	unreadableDir := filepath.Join("internal", "testdata", "mock_metrics", "unreadable_metrics")
	unreadableFile := filepath.Join(unreadableDir, "unreadable_metric.yaml")

	// Change the file permissions to 0000, so it becomes unreadable
	if err := os.Chmod(unreadableFile, 0000); err != nil {
		log.Fatalf("Error changing file permissions: %v", err)
	}

	tests := []struct {
		name       string
		path       string
		wantMetric []*assessment.Metric
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "Invalid path",
			path: "doesnotexist",
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "error accessing path")
			},
		},
		{
			name: "Invalid file",
			path: unreadableDir,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "error reading file")
			},
		},
		{
			name: "Invalid metric",
			path: invalidDir,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "error decoding metric")
			},
		},
		{
			name: "Happy path",
			path: validDir,
			wantMetric: []*assessment.Metric{
				{
					Id:          "TestMetric",
					Description: "Test Metric 1",
					Category:    "TestCategory",
					Version:     "1.0",
					Comments:    "Test comments",
				},
				{
					Id:          "TestMetric",
					Description: "Test Metric 2",
					Category:    "TestCategory",
					Version:     "1.0",
					Comments:    "Test comments",
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			svc := &Service{}
			defaultMetricsPath = tt.path
			gotMetrics, err := svc.loadMetricsFromMetricsRepository()

			if tt.wantMetric != nil {
				assert.NoError(t, api.Validate(gotMetrics[0]))
				assert.NoError(t, api.Validate(gotMetrics[1]))
			}

			tt.wantErr(t, err)
		})
	}
	// Change the file permissions back
	if err := os.Chmod(unreadableFile, 0644); err != nil {
		log.Fatalf("Error changing file permissions: %v", err)
	}
}

func Test_loadMetricImplementation(t *testing.T) {
	metric, err := os.ReadFile("internal/testutil/metrictest/metric.rego")
	if err != nil {
		t.Fatalf("Failed to read metric file: %v", err)
	}

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
				metricID: testdata.MockMetricID1,
				file:     "doesnotexist.rego",
			},
			wantImpl: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, os.ErrNotExist)
			},
		},
		{
			name: "Happy path",
			args: args{
				metricID: testdata.MockMetricID1,
				file:     "internal/testutil/metrictest/metric.rego",
			},
			wantImpl: &assessment.MetricImplementation{
				MetricId: testdata.MockMetricID1,
				Lang:     assessment.MetricImplementation_LANGUAGE_REGO,
				Code:     string(metric),
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotImpl, err := loadMetricImplementation(tt.args.metricID, tt.args.file)
			if tt.wantErr(t, err) && err == nil {
				assert.NoError(t, api.Validate(gotImpl))
			}

			if tt.wantImpl != nil {
				assert.NotEmpty(t, gotImpl)

				// Check if time is set
				assert.NotEmpty(t, gotImpl.UpdatedAt)
			}

			// Ignore updated_at
			assert.Equal(t, tt.wantImpl, gotImpl, protocmp.IgnoreFields(&assessment.MetricImplementation{}, "updated_at"))
		})
	}
}

func TestService_CreateMetric(t *testing.T) {
	type fields struct {
		storage persistence.Storage
	}
	type args struct {
		in0 context.Context
		req *orchestrator.CreateMetricRequest
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantMetric *assessment.Metric
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "Create valid metric",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
			},
			args: args{
				context.TODO(),
				&orchestrator.CreateMetricRequest{
					Metric: &assessment.Metric{
						Id:          testdata.MockMetricID1,
						Description: testdata.MockMetricDescription1,
						Category:    testdata.MockMetricCategory1,
						Version:     "1.0",
						Comments:    "Test comments",
					},
				},
			},
			wantMetric: &assessment.Metric{
				Id:          testdata.MockMetricID1,
				Description: testdata.MockMetricDescription1,
				Category:    testdata.MockMetricCategory1,
				Version:     "1.0",
				Comments:    "Test comments",
			},
			wantErr: assert.NoError,
		},
		{
			name: "Create invalid metric",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
			},
			args: args{
				context.TODO(),
				&orchestrator.CreateMetricRequest{
					Metric: &assessment.Metric{},
				},
			},
			wantMetric: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "metric.id: value length must be at least 1 characters")
			},
		},
		{
			name: "Create metric and set to deprecated",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
			},
			args: args{
				context.TODO(),
				&orchestrator.CreateMetricRequest{
					Metric: &assessment.Metric{
						Id:              "TLSVersion",
						Description:     testdata.MockMetricDescription1,
						Category:        testdata.MockMetricCategory1,
						Version:         "1.0",
						Comments:        "Test comments",
						DeprecatedSince: timestamppb.Now(),
					},
				},
			},
			wantMetric: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "the metric shouldn't be set to deprecated at creation time")
			},
		},
		{
			name: "Create metric which already exists",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(&assessment.Metric{Id: "TLSVersion"}))
				}),
			},
			args: args{
				context.TODO(),
				&orchestrator.CreateMetricRequest{
					Metric: &assessment.Metric{
						Id:          "TLSVersion",
						Description: testdata.MockMetricDescription1,
						Category:    testdata.MockMetricCategory1,
						Version:     "1.0",
						Comments:    "Test comments",
					},
				},
			},
			wantMetric: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "rpc error: code = AlreadyExists desc = metric already exists")
			},
		},
		{
			name: "Error while counting",
			fields: fields{
				storage: &testutil.StorageWithError{CountErr: ErrSomeError},
			},
			args: args{
				context.TODO(),
				&orchestrator.CreateMetricRequest{
					Metric: &assessment.Metric{
						Id:          testdata.MockMetricID1,
						Description: testdata.MockMetricDescription1,
						Category:    testdata.MockMetricCategory1,
						Version:     "1.0",
						Comments:    "Test comments",
					},
				},
			},
			wantMetric: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "rpc error: code = Internal desc = database error: some error")
			},
		},
		{
			name: "Error while creating",
			fields: fields{
				storage: &testutil.StorageWithError{CreateErr: ErrSomeError},
			},
			args: args{
				context.TODO(),
				&orchestrator.CreateMetricRequest{
					Metric: &assessment.Metric{
						Id:          testdata.MockMetricID1,
						Description: testdata.MockMetricDescription1,
						Category:    testdata.MockMetricCategory1,
						Version:     "1.0",
						Comments:    "Test comments",
					},
				},
			},
			wantMetric: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "rpc error: code = Internal desc = database error: some error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				storage: tt.fields.storage,
			}
			gotMetric, err := svc.CreateMetric(tt.args.in0, tt.args.req)
			if tt.wantErr(t, err) && err == nil {
				assert.NoError(t, api.Validate(gotMetric))
			}
			assert.Equal(t, tt.wantMetric, gotMetric)
		})
	}
}

func TestService_UpdateMetric(t *testing.T) {
	timestamp := timestamppb.Now()

	type fields struct {
		storage persistence.Storage
	}
	type args struct {
		in0 context.Context
		req *orchestrator.UpdateMetricRequest
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantMetric *assessment.Metric
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "Update existing metric",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					_ = s.Create(&assessment.Metric{Id: "TransportEncryptionEnabled"})
					_ = s.Create(&orchestrator.TargetOfEvaluation{Id: testdata.MockTargetOfEvaluationID1})
				}),
			},
			args: args{
				context.TODO(),
				&orchestrator.UpdateMetricRequest{
					Metric: &assessment.Metric{
						Id:              "TransportEncryptionEnabled",
						Description:     testdata.MockMetricDescription1,
						Category:        testdata.MockMetricCategory1,
						Version:         "1.0",
						Comments:        "Test comments",
						DeprecatedSince: timestamp,
					},
				},
			},
			wantMetric: &assessment.Metric{
				Id:              "TransportEncryptionEnabled",
				Description:     testdata.MockMetricDescription1,
				Category:        testdata.MockMetricCategory1,
				Version:         "1.0",
				Comments:        "Test comments",
				DeprecatedSince: timestamp,
			},
			wantErr: assert.NoError,
		},
		{
			name: "storage error: Get",
			fields: fields{
				storage: &testutil.StorageWithError{GetErr: ErrSomeError},
			},
			args: args{
				context.TODO(),
				&orchestrator.UpdateMetricRequest{
					Metric: &assessment.Metric{
						Id:              "TransportEncryptionEnabled",
						Description:     testdata.MockMetricDescription1,
						Category:        testdata.MockMetricCategory1,
						Version:         "1.0",
						Comments:        "Test comments",
						DeprecatedSince: timestamp,
					},
				},
			},
			wantMetric: nil,
			wantErr:    wantStatusCode(codes.Internal),
		},
		{
			name: "storage error: Save",
			fields: fields{
				storage: &testutil.StorageWithError{SaveErr: ErrSomeError},
			},
			args: args{
				context.TODO(),
				&orchestrator.UpdateMetricRequest{
					Metric: &assessment.Metric{
						Id:              "TransportEncryptionEnabled",
						Description:     testdata.MockMetricDescription1,
						Category:        testdata.MockMetricCategory1,
						Version:         "1.0",
						Comments:        "Test comments",
						DeprecatedSince: timestamp,
					},
				},
			},
			wantMetric: nil,
			wantErr:    wantStatusCode(codes.Internal),
		},
		{
			name: "Update non-existing metric",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
			},
			args: args{
				context.TODO(),
				&orchestrator.UpdateMetricRequest{
					Metric: &assessment.Metric{
						Id:              "DoesProbablyNotExist",
						Description:     testdata.MockMetricDescription1,
						Category:        testdata.MockMetricCategory1,
						Version:         "1.0",
						Comments:        "Test comments",
						DeprecatedSince: timestamp,
					},
				},
			},
			wantMetric: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "rpc error: code = NotFound desc = metric not found")
			},
		},
		{
			name: "Updating invalid metric",
			args: args{
				context.TODO(),
				&orchestrator.UpdateMetricRequest{
					Metric: &assessment.Metric{
						Id:      "DoesProbablyNotExist",
						Version: "",
					},
				},
			},
			wantMetric: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "metric.version: value length must be at least 1 characters")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				storage: tt.fields.storage,
			}
			gotMetric, err := svc.UpdateMetric(tt.args.in0, tt.args.req)
			if tt.wantErr(t, err) && err == nil {
				assert.NoError(t, api.Validate(gotMetric))
			}
			assert.Equal(t, tt.wantMetric, gotMetric)
		})
	}
}

func TestService_GetMetric(t *testing.T) {
	type fields struct {
		storage persistence.Storage
	}
	type args struct {
		in0 context.Context
		req *orchestrator.GetMetricRequest
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantMetric *assessment.Metric
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "Invalid request",
			args:       args{},
			wantMetric: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, api.ErrEmptyRequest.Error())
			},
		},
		{
			name: "storage error",
			fields: fields{
				storage: &testutil.StorageWithError{GetErr: ErrSomeError},
			},
			args: args{
				context.TODO(),
				&orchestrator.GetMetricRequest{
					MetricId: "TransportEncryptionEnabled",
				},
			},
			wantMetric: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrSomeError.Error())
			},
		},
		{
			name: "Get existing metric",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					_ = s.Create(&orchestrator.TargetOfEvaluation{Id: testdata.MockTargetOfEvaluationID1})
					_ = s.Create(&assessment.Metric{
						Id:          "TransportEncryptionEnabled",
						Description: "This metric describes, whether transport encryption is turned on or not",
						Category:    testdata.MockMetricCategory1,
						Version:     "1.0",
						Comments:    "Test comments",
					})
				}),
			},
			args: args{
				context.TODO(),
				&orchestrator.GetMetricRequest{
					MetricId: "TransportEncryptionEnabled",
				},
			},
			wantMetric: &assessment.Metric{
				Id:          "TransportEncryptionEnabled",
				Description: "This metric describes, whether transport encryption is turned on or not",
				Category:    testdata.MockMetricCategory1,
				Version:     "1.0",
				Comments:    "Test comments",
			},
			wantErr: assert.NoError,
		},
		{
			name:   "Get non-existing metric",
			fields: fields{storage: testutil.NewInMemoryStorage(t)},
			args: args{
				context.TODO(),
				&orchestrator.GetMetricRequest{
					MetricId: "DoesProbablyNotExist",
				},
			},
			wantMetric: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "metric not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				storage: tt.fields.storage,
			}
			gotMetric, err := svc.GetMetric(tt.args.in0, tt.args.req)
			if tt.wantErr(t, err) && err == nil {
				assert.NoError(t, api.Validate(gotMetric))
			}
			assert.Equal(t, tt.wantMetric, gotMetric)
		})
	}
}

func TestService_ListMetrics(t *testing.T) {
	timestamp := timestamppb.Now()

	type fields struct {
		TargetOfEvaluationHooks []orchestrator.TargetOfEvaluationHookFunc
		auditScopeHooks         []orchestrator.AuditScopeHookFunc
		AssessmentResultHooks   []assessment.ResultHookFunc
		storage                 persistence.Storage
		loadMetricsFunc         func() ([]*assessment.Metric, error)
		catalogsFolder          string
		loadCatalogsFunc        func() ([]*orchestrator.Catalog, error)
		events                  chan *orchestrator.MetricChangeEvent
		authz                   service.AuthorizationStrategy
	}
	type args struct {
		in0 context.Context
		req *orchestrator.ListMetricsRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *orchestrator.ListMetricsResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Invalid input",
			args: args{},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, api.ErrEmptyRequest.Error())
			},
		},
		{
			name: "Happy path: all active metrics",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					_ = s.Create(&assessment.Metric{
						Id:          testdata.MockMetricID1,
						Description: testdata.MockMetricDescription1,
						Category:    testdata.MockMetricCategory1,
						Version:     "1.0",
						Comments:    "Test comments",
					})
					_ = s.Create(&assessment.Metric{
						Id:              testdata.MockMetricID2,
						Description:     testdata.MockMetricDescription2,
						Category:        testdata.MockMetricCategory1,
						Version:         "1.0",
						Comments:        "Test comments",
						DeprecatedSince: timestamp,
					})
				}),
			},
			args: args{
				req: &orchestrator.ListMetricsRequest{},
			},
			wantRes: &orchestrator.ListMetricsResponse{
				Metrics: []*assessment.Metric{MockMetric1},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: including deprecated metrics",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					_ = s.Create(&assessment.Metric{
						Id:          testdata.MockMetricID1,
						Description: testdata.MockMetricDescription1,
						Category:    testdata.MockMetricCategory1,
						Version:     "1.0",
						Comments:    "Test comments",
					})
					_ = s.Create(&assessment.Metric{
						Id:              testdata.MockMetricID2,
						Description:     testdata.MockMetricDescription2,
						Category:        testdata.MockMetricCategory1,
						Version:         "1.0",
						Comments:        "Test comments",
						DeprecatedSince: timestamp,
					})
				}),
			},
			args: args{
				req: &orchestrator.ListMetricsRequest{
					Filter: &orchestrator.ListMetricsRequest_Filter{
						IncludeDeprecated: proto.Bool(true),
					},
				},
			},
			wantRes: &orchestrator.ListMetricsResponse{
				Metrics: []*assessment.Metric{
					{
						Id:          testdata.MockMetricID1,
						Description: testdata.MockMetricDescription1,
						Category:    testdata.MockMetricCategory1,
						Version:     "1.0",
						Comments:    "Test comments",
					},
					{
						Id:              testdata.MockMetricID2,
						Description:     testdata.MockMetricDescription2,
						Category:        testdata.MockMetricCategory1,
						Version:         "1.0",
						Comments:        "Test comments",
						DeprecatedSince: timestamp,
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				TargetOfEvaluationHooks: tt.fields.TargetOfEvaluationHooks,
				auditScopeHooks:         tt.fields.auditScopeHooks,
				AssessmentResultHooks:   tt.fields.AssessmentResultHooks,
				loadMetricsFunc:         tt.fields.loadMetricsFunc,
				storage:                 tt.fields.storage,
				catalogsFolder:          tt.fields.catalogsFolder,
				loadCatalogsFunc:        tt.fields.loadCatalogsFunc,
				events:                  tt.fields.events,
				authz:                   tt.fields.authz,
			}
			gotRes, err := svc.ListMetrics(tt.args.in0, tt.args.req)
			tt.wantErr(t, err)

			if tt.wantRes != nil {
				assert.NoError(t, api.Validate(gotRes))
			}
			assert.Equal(t, tt.wantRes, gotRes)
		})
	}
}

func TestService_GetMetricImplementation(t *testing.T) {
	type fields struct {
		AssessmentResultHooks []assessment.ResultHookFunc
		storage               persistence.Storage
		events                chan *orchestrator.MetricChangeEvent
		catalogsFolder        string
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
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "Invalid input",
			wantRes: nil,
			wantErr: wantStatusCode(codes.InvalidArgument),
		},
		{
			name: "metric not found",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
			},
			args: args{
				req: &orchestrator.GetMetricImplementationRequest{
					MetricId: testdata.MockMetricID1,
				},
			},
			wantErr: wantStatusCode(codes.NotFound),
		},
		{
			name: "storage error",
			fields: fields{
				storage: &testutil.StorageWithError{GetErr: ErrSomeError},
			},
			args: args{
				req: &orchestrator.GetMetricImplementationRequest{
					MetricId: testdata.MockMetricID1,
				},
			},
			wantErr: wantStatusCode(codes.Internal),
		},
		{
			name: "metric found",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Save(&assessment.Metric{Id: testdata.MockMetricID1})
					assert.NoError(t, err)
					err = s.Save(&assessment.MetricImplementation{
						MetricId: testdata.MockMetricID1,
						Lang:     assessment.MetricImplementation_LANGUAGE_REGO,
						Code:     "package test",
					}, "metric_id = ?", testdata.MockMetricID1)
					assert.NoError(t, err)
				}),
			},
			args: args{
				req: &orchestrator.GetMetricImplementationRequest{
					MetricId: testdata.MockMetricID1,
				},
			},
			wantRes: &assessment.MetricImplementation{
				MetricId: testdata.MockMetricID1,
				Lang:     assessment.MetricImplementation_LANGUAGE_REGO,
				Code:     "package test",
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				catalogsFolder:        tt.fields.catalogsFolder,
				events:                tt.fields.events,
			}

			gotRes, err := svc.GetMetricImplementation(tt.args.ctx, tt.args.req)
			if tt.wantErr(t, err) && err == nil {
				assert.NoError(t, api.Validate(gotRes))
			}

			assert.Equal(t, tt.wantRes, gotRes)
		})
	}
}

func TestService_UpdateMetricImplementation(t *testing.T) {
	type fields struct {
		AssessmentResultHooks []assessment.ResultHookFunc
		storage               persistence.Storage
		catalogsFolder        string
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
		wantImpl assert.Want[*assessment.MetricImplementation]
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "Invalid input",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
			},
			args: args{
				req: &orchestrator.UpdateMetricImplementationRequest{
					Implementation: &assessment.MetricImplementation{MetricId: "notfound"},
				},
			},
			wantImpl: assert.Nil[*assessment.MetricImplementation],
			wantErr:  wantStatusCode(codes.InvalidArgument),
		},
		{
			name: "Metric not found",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
			},
			args: args{
				req: &orchestrator.UpdateMetricImplementationRequest{
					Implementation: &assessment.MetricImplementation{
						MetricId: "TransportEncryptionEnabled",
						Lang:     assessment.MetricImplementation_LANGUAGE_REGO,
						Code:     "package example",
					},
				},
			},
			wantImpl: assert.Nil[*assessment.MetricImplementation],
			wantErr:  wantStatusCode(codes.NotFound),
		},
		{
			name: "storage error: Get",
			fields: fields{
				storage: &testutil.StorageWithError{GetErr: ErrSomeError},
			},
			args: args{
				req: &orchestrator.UpdateMetricImplementationRequest{
					Implementation: &assessment.MetricImplementation{
						MetricId: "TransportEncryptionEnabled",
						Lang:     assessment.MetricImplementation_LANGUAGE_REGO,
						Code:     "package example",
					},
				},
			},
			wantImpl: assert.Nil[*assessment.MetricImplementation],
			wantErr:  wantStatusCode(codes.Internal),
		},
		{
			name: "storage error: Save",
			fields: fields{
				storage: &testutil.StorageWithError{SaveErr: ErrSomeError},
			},
			args: args{
				req: &orchestrator.UpdateMetricImplementationRequest{
					Implementation: &assessment.MetricImplementation{
						MetricId: "TransportEncryptionEnabled",
						Lang:     assessment.MetricImplementation_LANGUAGE_REGO,
						Code:     "package example",
					},
				},
			},
			wantImpl: assert.Nil[*assessment.MetricImplementation],
			wantErr:  wantStatusCode(codes.Internal),
		},
		{
			name: "update",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					_ = s.Create(&assessment.Metric{Id: "TransportEncryptionEnabled"})
				}),
			},
			args: args{
				req: &orchestrator.UpdateMetricImplementationRequest{
					Implementation: &assessment.MetricImplementation{
						MetricId: "TransportEncryptionEnabled",
						Lang:     assessment.MetricImplementation_LANGUAGE_REGO,
						Code:     "package example",
					},
				},
			},
			wantImpl: func(t *testing.T, got *assessment.MetricImplementation) bool {
				return assert.Equal(t, "TransportEncryptionEnabled", got.MetricId) &&
					assert.Equal(t, assessment.MetricImplementation_LANGUAGE_REGO, got.Lang) &&
					assert.Equal(t, "package example", got.Code) &&
					assert.True(t, got.UpdatedAt.AsTime().Before(time.Now())) &&
					assert.NoError(t, api.Validate(got))
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				catalogsFolder:        tt.fields.catalogsFolder,
				events:                tt.fields.events,
			}
			gotImpl, err := svc.UpdateMetricImplementation(tt.args.in0, tt.args.req)

			tt.wantErr(t, err)
			tt.wantImpl(t, gotImpl)
		})
	}
}

func TestService_GetMetricConfiguration(t *testing.T) {
	type fields struct {
		AssessmentResultHooks []assessment.ResultHookFunc
		storage               persistence.Storage
		catalogsFolder        string
		events                chan *orchestrator.MetricChangeEvent
		authz                 service.AuthorizationStrategy
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
		want         assert.Want[*assessment.MetricConfiguration]
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "Invalid input",
			args: args{
				req: &orchestrator.GetMetricConfigurationRequest{
					TargetOfEvaluationId: "InvalidTargetOfEvaluationID",
				},
			},
			want:    assert.Nil[*assessment.MetricConfiguration],
			wantErr: wantStatusCode(codes.InvalidArgument),
		},
		{
			name: "Permission denied",
			args: args{
				req: &orchestrator.GetMetricConfigurationRequest{
					MetricId:             testdata.MockMetricID1,
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
				},
			},
			fields: fields{
				authz: &servicetest.AuthorizationStrategyMock{},
			},
			want:    assert.Nil[*assessment.MetricConfiguration],
			wantErr: wantStatusCode(codes.PermissionDenied),
		},
		{
			name: "storage error",
			fields: fields{
				authz:   &service.AuthorizationStrategyAllowAll{},
				storage: &testutil.StorageWithError{GetErr: ErrSomeError},
			},
			args: args{
				req: &orchestrator.GetMetricConfigurationRequest{
					MetricId:             testdata.MockMetricID1,
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
				},
			},
			want: assert.Nil[*assessment.MetricConfiguration],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrSomeError.Error())
			},
		},
		{
			name: "metric found",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					_ = s.Create(&assessment.Metric{Id: testdata.MockMetricID1})
					_ = s.Create(&orchestrator.TargetOfEvaluation{
						Id: testdata.MockTargetOfEvaluationID1,
					})
					_ = s.Create(&assessment.MetricConfiguration{
						MetricId:             testdata.MockMetricID1,
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
						Operator:             "==",
						TargetValue:          &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: "1111"}},
					})
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.GetMetricConfigurationRequest{
					MetricId:             testdata.MockMetricID1,
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
				},
			},
			wantResponse: &assessment.MetricConfiguration{
				MetricId:             testdata.MockMetricID1,
				TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
				Operator:             "==",
				TargetValue:          &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: "1111"}},
			},
			want: func(t *testing.T, got *assessment.MetricConfiguration) bool {
				wantResp := &assessment.MetricConfiguration{
					MetricId:             testdata.MockMetricID1,
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					Operator:             "==",
					TargetValue:          &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: "1111"}},
				}

				return assert.Equal(t, wantResp, got)
			},
			wantErr: assert.NoError,
		},
		{
			name: "metric not found",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.GetMetricConfigurationRequest{
					MetricId:             "NotExists",
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
				},
			},
			want:    assert.Nil[*assessment.MetricConfiguration],
			wantErr: wantStatusCode(codes.NotFound),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				catalogsFolder:        tt.fields.catalogsFolder,
				events:                tt.fields.events,
				authz:                 tt.fields.authz,
			}
			gotResponse, err := s.GetMetricConfiguration(tt.args.in0, tt.args.req)
			if tt.wantErr(t, err, tt.args) && err == nil {
				assert.NoError(t, api.Validate(gotResponse))
			}

			tt.want(t, gotResponse)
		})
	}
}

func TestService_ListMetricConfigurations(t *testing.T) {
	type fields struct {
		AssessmentResultHooks []assessment.ResultHookFunc
		storage               persistence.Storage
		catalogsFolder        string
		events                chan *orchestrator.MetricChangeEvent
		authz                 service.AuthorizationStrategy
	}
	type args struct {
		ctx context.Context
		req *orchestrator.ListMetricConfigurationRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *orchestrator.ListMetricConfigurationResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "error",
			fields: fields{
				storage: &testutil.StorageWithError{ListErr: ErrSomeError},
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.ListMetricConfigurationRequest{},
			},
			wantErr: wantStatusCode(codes.InvalidArgument),
		},
		{
			name: "Permission denied",
			fields: fields{
				authz: &servicetest.AuthorizationStrategyMock{},
			},
			args: args{
				req: &orchestrator.ListMetricConfigurationRequest{
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
				},
			},
			wantRes: nil,
			wantErr: wantStatusCode(codes.PermissionDenied),
		},
		{
			name: "storage error",
			fields: fields{
				storage: &testutil.StorageWithError{ListErr: ErrSomeError},
				authz:   &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				ctx: context.Background(),
				req: &orchestrator.ListMetricConfigurationRequest{
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
				},
			},
			wantErr: wantStatusCode(codes.Internal),
		},
		{
			name: "no error",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					_ = s.Create(&assessment.Metric{Id: testdata.MockMetricID1})
					_ = s.Create(&orchestrator.TargetOfEvaluation{
						Id: testdata.MockTargetOfEvaluationID1,
					})
					_ = s.Create(MockMetricConfiguration1)
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.ListMetricConfigurationRequest{
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
				},
			},
			wantRes: &orchestrator.ListMetricConfigurationResponse{
				Configurations: map[string]*assessment.MetricConfiguration{
					testdata.MockMetricID1: MockMetricConfiguration1,
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				catalogsFolder:        tt.fields.catalogsFolder,
				events:                tt.fields.events,
				authz:                 tt.fields.authz,
			}
			gotRes, err := svc.ListMetricConfigurations(tt.args.ctx, tt.args.req)
			tt.wantErr(t, err)

			if tt.wantRes != nil {
				// Check if response validation succeeds
				assert.NoError(t, api.Validate(gotRes))
			}
			assert.Equal(t, tt.wantRes, gotRes)
		})
	}
}

func TestService_UpdateMetricConfiguration(t *testing.T) {
	type fields struct {
		AssessmentResultHooks []assessment.ResultHookFunc
		storage               persistence.Storage
		loadMetricsFunc       func() ([]*assessment.Metric, error)
		catalogsFolder        string
		loadCatalogsFunc      func() ([]*orchestrator.Catalog, error)
		events                chan *orchestrator.MetricChangeEvent
		authz                 service.AuthorizationStrategy
	}
	type args struct {
		in0 context.Context
		req *orchestrator.UpdateMetricConfigurationRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.Want[*assessment.MetricConfiguration]
		wantSvc assert.Want[*Service]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Metric configuration invalid",
			args: args{
				req: &orchestrator.UpdateMetricConfigurationRequest{
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					MetricId:             testdata.MockMetricID1,
					Configuration: &assessment.MetricConfiguration{
						Operator:             "invalidOperator",
						TargetValue:          testdata.MockMetricConfigurationTargetValueString,
						MetricId:             testdata.MockMetricID1,
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					},
				},
			},
			want: assert.Nil[*assessment.MetricConfiguration],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "value does not match regex pattern")
			},
		},
		{
			name: "Permission denied",
			fields: fields{
				authz: &servicetest.AuthorizationStrategyMock{},
			},
			args: args{
				req: &orchestrator.UpdateMetricConfigurationRequest{
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					MetricId:             testdata.MockMetricID1,
					Configuration: &assessment.MetricConfiguration{
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
						MetricId:             testdata.MockMetricID1,
						Operator:             "<",
						TargetValue:          testdata.MockMetricConfigurationTargetValueString,
					},
				},
			},
			want: assert.Nil[*assessment.MetricConfiguration],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, service.ErrPermissionDenied.Error())
			},
		},
		{
			name: "storage error",
			fields: fields{
				authz:   &service.AuthorizationStrategyAllowAll{},
				storage: &testutil.StorageWithError{SaveErr: ErrSomeError},
			},
			args: args{
				req: &orchestrator.UpdateMetricConfigurationRequest{
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					MetricId:             testdata.MockMetricID1,
					Configuration: &assessment.MetricConfiguration{
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
						MetricId:             testdata.MockMetricID1,
						Operator:             "<",
						TargetValue:          testdata.MockMetricConfigurationTargetValueString,
					},
				},
			},
			want: assert.Nil[*assessment.MetricConfiguration],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrSomeError.Error())
			},
		},
		{
			name: "metricId is missing in request",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.UpdateMetricConfigurationRequest{
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					Configuration: &assessment.MetricConfiguration{
						Operator:             "<",
						TargetValue:          testdata.MockMetricConfigurationTargetValueString,
						MetricId:             testdata.MockMetricID1,
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					},
				},
			},
			want: assert.Nil[*assessment.MetricConfiguration],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "at least")
			},
		},
		{
			name: "TargetOfEvaluationID is missing in request",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.UpdateMetricConfigurationRequest{
					MetricId: testdata.MockMetricID1,
					Configuration: &assessment.MetricConfiguration{
						Operator:             "<",
						TargetValue:          testdata.MockMetricConfigurationTargetValueString,
						MetricId:             testdata.MockMetricID1,
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					},
				},
			},
			want: assert.Nil[*assessment.MetricConfiguration],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "target_of_evaluation_id: value is empty, which is not a valid UUI")
			},
		},
		{
			name: "TargetOfEvaluationID is invalid in request",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.UpdateMetricConfigurationRequest{
					TargetOfEvaluationId: "00000000-000000000000",
					MetricId:             testdata.MockMetricID1,
					Configuration: &assessment.MetricConfiguration{
						Operator:             "<",
						TargetValue:          testdata.MockMetricConfigurationTargetValueString,
						MetricId:             testdata.MockMetricID1,
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					},
				},
			},
			want: assert.Nil[*assessment.MetricConfiguration],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "target_of_evaluation_id: value must be a valid UUID")
			},
		},
		{
			name: "configuration is missing in request",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.UpdateMetricConfigurationRequest{
					MetricId:             testdata.MockMetricID1,
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
				},
			},
			want: assert.Nil[*assessment.MetricConfiguration],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "configuration: value is required")
			},
		},
		{
			name: "metricId is missing in configuration",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.UpdateMetricConfigurationRequest{
					MetricId:             testdata.MockMetricID1,
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					Configuration: &assessment.MetricConfiguration{
						Operator:             "<",
						TargetValue:          testdata.MockMetricConfigurationTargetValueString,
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					},
				},
			},
			want: assert.Nil[*assessment.MetricConfiguration],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "at least")
			},
		},
		{
			name: "TargetOfEvaluationId is missing in configuration",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.UpdateMetricConfigurationRequest{
					MetricId:             testdata.MockMetricID1,
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					Configuration: &assessment.MetricConfiguration{
						Operator:    "<",
						TargetValue: testdata.MockMetricConfigurationTargetValueString,
						MetricId:    testdata.MockMetricID1,
					},
				},
			},
			want: assert.Nil[*assessment.MetricConfiguration],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "configuration.target_of_evaluation_id: value is empty, which is not a valid UUID")
			},
		},
		{
			name: "TargetOfEvaluationId is invalid in configuration",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.UpdateMetricConfigurationRequest{
					MetricId:             testdata.MockMetricID1,
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					Configuration: &assessment.MetricConfiguration{
						Operator:             "<",
						TargetValue:          testdata.MockMetricConfigurationTargetValueString,
						MetricId:             testdata.MockMetricID1,
						TargetOfEvaluationId: "00000000-000000000000",
					},
				},
			},
			want: assert.Nil[*assessment.MetricConfiguration],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "target_of_evaluation_id: value must be a valid UUID")
			},
		},
		{
			name: "metric does not exist in storage",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.UpdateMetricConfigurationRequest{
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					MetricId:             testdata.MockMetricID1,
					Configuration: &assessment.MetricConfiguration{
						Operator:             "<",
						TargetValue:          testdata.MockMetricConfigurationTargetValueString,
						MetricId:             testdata.MockMetricID1,
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					},
				},
			},
			want: assert.Nil[*assessment.MetricConfiguration],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "metric or service does not exist")
			},
		},
		{
			name: "TargetOfEvaluation does not exist in storage",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					_ = s.Create(&assessment.Metric{Id: testdata.MockMetricID1})
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.UpdateMetricConfigurationRequest{
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					MetricId:             testdata.MockMetricID1,
					Configuration: &assessment.MetricConfiguration{
						Operator:             "<",
						TargetValue:          testdata.MockMetricConfigurationTargetValueString,
						MetricId:             testdata.MockMetricID1,
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					},
				},
			},
			want: assert.Nil[*assessment.MetricConfiguration],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "metric or service does not exist")
			},
		},
		{
			name: "append metric configuration",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					_ = s.Create(&assessment.Metric{Id: testdata.MockMetricID1})
					_ = s.Create(&orchestrator.TargetOfEvaluation{Id: testdata.MockTargetOfEvaluationID1})
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.UpdateMetricConfigurationRequest{
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					MetricId:             testdata.MockMetricID1,
					Configuration: &assessment.MetricConfiguration{
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
						MetricId:             testdata.MockMetricID1,
						Operator:             "<",
						TargetValue:          testdata.MockMetricConfigurationTargetValueString,
					},
				},
			},
			wantErr: assert.NoError,
			wantSvc: func(t *testing.T, got *Service) bool {
				var config *assessment.MetricConfiguration
				err := got.storage.Get(&config, persistence_gorm.WithoutPreload(), "target_of_evaluation_id = ? AND metric_id = ?", testdata.MockTargetOfEvaluationID1, testdata.MockMetricID1)
				if !assert.NoError(t, err) {
					return false
				}

				return assert.Equal(t, config.Operator, "<")
			},
			want: func(t *testing.T, got *assessment.MetricConfiguration) bool {
				return assert.NoError(t, api.Validate(got))
			},
		},
		{
			name: "replace metric configuration",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					_ = s.Create(&assessment.Metric{Id: testdata.MockMetricID1})
					_ = s.Create(&orchestrator.TargetOfEvaluation{
						Id: testdata.MockTargetOfEvaluationID1,
					})
					_ = s.Create(&assessment.MetricConfiguration{
						MetricId:             testdata.MockMetricID1,
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
						Operator:             ">",
					})
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.UpdateMetricConfigurationRequest{
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					MetricId:             testdata.MockMetricID1,
					Configuration: &assessment.MetricConfiguration{
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
						MetricId:             testdata.MockMetricID1,
						Operator:             "<",
						TargetValue:          testdata.MockMetricConfigurationTargetValueString,
					},
				},
			},
			wantErr: assert.NoError,
			wantSvc: func(t *testing.T, got *Service) bool {
				var config *assessment.MetricConfiguration
				err := got.storage.Get(&config, persistence_gorm.WithoutPreload(), "target_of_evaluation_id = ? AND metric_id = ?", testdata.MockTargetOfEvaluationID1, testdata.MockMetricID1)
				if !assert.NoError(t, err) {
					return false
				}

				return assert.Equal(t, config.Operator, "<")
			},
			want: func(t *testing.T, got *assessment.MetricConfiguration) bool {
				return assert.NoError(t, api.Validate(got))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				loadMetricsFunc:       tt.fields.loadMetricsFunc,
				catalogsFolder:        tt.fields.catalogsFolder,
				loadCatalogsFunc:      tt.fields.loadCatalogsFunc,
				events:                tt.fields.events,
				authz:                 tt.fields.authz,
			}
			gotRes, err := svc.UpdateMetricConfiguration(tt.args.in0, tt.args.req)
			tt.wantErr(t, err)
			tt.want(t, gotRes)
			assert.Optional(t, tt.wantSvc, svc)
		})
	}
}

func wantStatusCode(code codes.Code) assert.ErrorAssertionFunc {
	return func(t assert.TestingT, err error, i ...interface{}) bool {
		gotStatus, ok := status.FromError(err)
		if !ok {
			return false
		}
		return assert.Equal(t, gotStatus.Code(), code)
	}
}

func TestService_RemoveMetric(t *testing.T) {
	timestamp := timestamppb.New(time.Date(2017, 12, 1, 0, 0, 0, 0, time.Local))
	// timestamp := timestamppb.New(time.Time{})
	// timestamp := time.Date(2011, 7, 1, 0, 0, 0, 0, time.UTC)
	type fields struct {
		storage persistence.Storage
	}
	type args struct {
		ctx context.Context
		req *orchestrator.RemoveMetricRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes assert.Want[*emptypb.Empty]
		wantSvc assert.Want[*Service]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Validation Error - Request is nil",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
				}),
			},
			args: args{
				ctx: nil,
				req: nil,
			},
			wantRes: assert.Nil[*emptypb.Empty],
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, api.ErrEmptyRequest.Error())
			},
		},
		{
			name: "Validation Error - metric id is empty",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
				}),
			},
			args: args{
				ctx: nil,
				req: &orchestrator.RemoveMetricRequest{MetricId: ""},
			},
			wantRes: assert.Nil[*emptypb.Empty],
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, api.ErrInvalidRequest.Error())
			},
		},
		{
			name: "Error - Internal (Get)",
			fields: fields{
				storage: &testutil.StorageWithError{GetErr: gorm.ErrInvalidDB},
			},
			args: args{
				ctx: nil,
				req: &orchestrator.RemoveMetricRequest{MetricId: testdata.MockMetricID1},
			},
			wantRes: assert.Nil[*emptypb.Empty],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.Internal, status.Code(err))
				return assert.ErrorContains(t, err, gorm.ErrInvalidDB.Error())
			},
		},
		{
			name: "Error - Not Found",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
				}),
			},
			args: args{
				ctx: nil,
				req: &orchestrator.RemoveMetricRequest{MetricId: testdata.MockMetricID1},
			},
			wantRes: assert.Nil[*emptypb.Empty],
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				assert.Equal(t, codes.NotFound, status.Code(err))
				return assert.ErrorContains(t, err, ErrMetricNotFound.Error())
			},
		},
		{
			name: "Happy path",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					_ = s.Create(orchestratortest.NewMetric())
				}),
			},
			args: args{
				context.TODO(),
				&orchestrator.RemoveMetricRequest{
					MetricId: testdata.MockMetricID1,
				},
			},
			wantRes: assert.NotNil[*emptypb.Empty],
			wantSvc: func(t *testing.T, got *Service) bool {
				var gotMetric *assessment.Metric

				err := got.storage.Get(&gotMetric, "id = ?", testdata.MockMetricID1)
				assert.NoError(t, err)

				return assert.NotEmpty(t, gotMetric.DeprecatedSince)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: metric already removed in the past",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					_ = s.Create(&orchestrator.TargetOfEvaluation{Id: testdata.MockTargetOfEvaluationID1})
					_ = s.Create(&assessment.Metric{
						Id:              testdata.MockMetricID1,
						Description:     testdata.MockMetricDescription1,
						Category:        testdata.MockMetricCategory1,
						Version:         "1.0",
						Comments:        "comments",
						DeprecatedSince: timestamp,
					},
					)
				}),
			},
			args: args{
				context.TODO(),
				&orchestrator.RemoveMetricRequest{
					MetricId: testdata.MockMetricID1,
				},
			},
			wantSvc: func(t *testing.T, got *Service) bool {
				var gotMetric *assessment.Metric

				err := got.storage.Get(&gotMetric, "id = ?", testdata.MockMetricID1)
				assert.NoError(t, err)

				assert.Equal(t, timestamp, gotMetric.DeprecatedSince)
				return assert.NotEmpty(t, gotMetric.DeprecatedSince)
			},
			wantRes: assert.NotNil[*emptypb.Empty],
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				storage: tt.fields.storage,
			}

			res, err := svc.RemoveMetric(context.TODO(), tt.args.req)
			tt.wantErr(t, err)
			tt.wantRes(t, res)
			assert.Optional(t, tt.wantSvc, svc)
		})
	}
}
