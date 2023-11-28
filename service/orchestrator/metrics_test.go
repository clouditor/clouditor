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
	"time"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/internal/testutil/servicetest"
	"clouditor.io/clouditor/internal/testutil/servicetest/orchestratortest"
	"clouditor.io/clouditor/persistence"
	persistence_gorm "clouditor.io/clouditor/persistence/gorm"
	"clouditor.io/clouditor/service"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	// "gorm.io/gorm"
)

var (
	ErrSomeError = errors.New("some error")

	// We need to define the following vars here because we could get import cycle errors in ./internal/testdata/testdata.go
	MockMetricConfiguration1 = &assessment.MetricConfiguration{
		MetricId:       testdata.MockMetricID1,
		CloudServiceId: testdata.MockCloudServiceID1,
		Operator:       "==",
		TargetValue:    testdata.MockMetricConfigurationTargetValueString,
	}

	MockMetricRange1 = &assessment.Range{
		Range: &assessment.Range_AllowedValues{
			AllowedValues: &assessment.AllowedValues{
				Values: []*structpb.Value{
					structpb.NewBoolValue(false),
					structpb.NewBoolValue(true),
				}}}}

	MockMetric1 = &assessment.Metric{
		Id:          testdata.MockMetricID1,
		Name:        testdata.MockMetricName1,
		Description: testdata.MockMetricDescription1,
		Scale:       assessment.Metric_ORDINAL,
		Range:       MockMetricRange1,
	}
)

func TestService_loadMetrics(t *testing.T) {
	type fields struct {
		AssessmentResultHooks []assessment.ResultHookFunc
		storage               persistence.Storage
		metricsFile           string
		loadMetricsFunc       func() ([]*assessment.Metric, error)
		catalogsFolder        string
		loadCatalogsFunc      func() ([]*orchestrator.Catalog, error)
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
			// There is a storage error for creating a single metric which is not forwarded/returned
			name: "storage error (but not returned)",
			fields: fields{
				metricsFile: "metrics.json",
				storage:     &testutil.StorageWithError{CreateErr: ErrSomeError},
			},
			wantErr: assert.NoError,
		},
		{
			name: "custom loading function with error",
			fields: fields{
				loadMetricsFunc: func() ([]*assessment.Metric, error) {
					return nil, ErrSomeError
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrSomeError)
			},
		},
		{
			// To test more accurately, we would need, e.g. an integration test since this function only returns error
			name: "happy path",
			fields: fields{
				// nil enforce the local embedded metric function to be used
				loadMetricsFunc: nil,
				// empty string enforces the DefaultValue to be used
				metricsFile: DefaultMetricsFile,
				storage:     testutil.NewInMemoryStorage(t),
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				metricsFile:           tt.fields.metricsFile,
				loadMetricsFunc:       tt.fields.loadMetricsFunc,
				catalogsFolder:        tt.fields.catalogsFolder,
				loadCatalogsFunc:      tt.fields.loadCatalogsFunc,
				events:                tt.fields.events,
			}

			err := svc.loadMetrics()
			tt.wantErr(t, err)
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
				Code: `package clouditor.metrics.admin_mfa_enabled

import data.clouditor.compare
import future.keywords.every
import input as identity

default applicable = false

default compliant = false

applicable {
	# we are only interested in some kind of privileged user    
	identity.privileged
}

compliant {
	# count the number of "factors"
	compare(data.operator, data.target_value, count(identity.authenticity))

	# also make sure, that we do not have any "NoAuthentication" in the factor and all are activated
	every factor in identity.authenticity {
		# TODO(oxisto): we do not have this type property (yet)
		not factor.type == "NoAuthentication"

		factor.activated == true
	}
}
`,
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotImpl, err := loadMetricImplementation(tt.args.metricID, tt.args.file)

			assert.NoError(t, gotImpl.Validate())
			tt.wantErr(t, err)

			if tt.wantImpl != nil {
				assert.NotEmpty(t, gotImpl)

				// Check if time is set and than delete if for the deepEqual check
				assert.NotEmpty(t, gotImpl.UpdatedAt)
				gotImpl.UpdatedAt = nil
			}

			if !reflect.DeepEqual(gotImpl, tt.wantImpl) {
				t.Errorf("loadMetricImplementation() = %v, want %v", gotImpl, tt.wantImpl)
			}
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
						Id:    testdata.MockMetricID1,
						Name:  testdata.MockMetricName1,
						Scale: assessment.Metric_ORDINAL,
						Range: &assessment.Range{Range: &assessment.Range_MinMax{}},
					},
				},
			},
			wantMetric: &assessment.Metric{
				Id:    testdata.MockMetricID1,
				Name:  testdata.MockMetricName1,
				Scale: assessment.Metric_ORDINAL,
				Range: &assessment.Range{Range: &assessment.Range_MinMax{}},
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
				return assert.ErrorContains(t, err, "invalid Metric.Id: value length must be at least 1 runes")
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
						Name:            "TLSMetricMockName",
						Category:        "",
						Scale:           assessment.Metric_NOMINAL,
						Range:           &assessment.Range{},
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
						Id:       "TLSVersion",
						Name:     "TLSMetricMockName",
						Category: "",
						Scale:    assessment.Metric_NOMINAL,
						Range:    &assessment.Range{},
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
						Id:    testdata.MockMetricName1,
						Name:  testdata.MockMetricName1,
						Scale: assessment.Metric_NOMINAL,
						Range: &assessment.Range{},
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
						Id:    testdata.MockMetricID1,
						Name:  testdata.MockMetricName1,
						Scale: assessment.Metric_NOMINAL,
						Range: &assessment.Range{},
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

			assert.NoError(t, gotMetric.Validate())
			tt.wantErr(t, err)

			if !proto.Equal(gotMetric, tt.wantMetric) {
				t.Errorf("Service.CreateMetric() = %v, want %v", gotMetric, tt.wantMetric)
			}
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
				}),
			},
			args: args{
				context.TODO(),
				&orchestrator.UpdateMetricRequest{
					Metric: &assessment.Metric{
						Id:              "TransportEncryptionEnabled",
						Name:            "A slightly updated metric",
						Scale:           assessment.Metric_NOMINAL,
						Range:           &assessment.Range{Range: &assessment.Range_AllowedValues{}},
						DeprecatedSince: timestamp,
					},
				},
			},
			wantMetric: &assessment.Metric{
				Id:              "TransportEncryptionEnabled",
				Name:            "A slightly updated metric",
				Scale:           assessment.Metric_NOMINAL,
				Range:           &assessment.Range{Range: &assessment.Range_AllowedValues{}},
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
						Name:            "A slightly updated metric",
						Scale:           assessment.Metric_NOMINAL,
						Range:           &assessment.Range{Range: &assessment.Range_AllowedValues{}},
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
						Name:            "TransportEncryptionEnabled",
						Description:     testdata.MockMetricDescription1,
						Category:        testdata.MockMetricCategory1,
						Scale:           assessment.Metric_NOMINAL,
						Range:           &assessment.Range{Range: &assessment.Range_AllowedValues{}},
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
						Name:            "UpdateMetricName",
						Scale:           assessment.Metric_NOMINAL,
						Range:           &assessment.Range{Range: &assessment.Range_AllowedValues{}},
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
						Id: "DoesProbablyNotExist",
					},
				},
			},
			wantMetric: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "Name: value length must be at least 1 runes")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				storage: tt.fields.storage,
			}
			gotMetric, err := svc.UpdateMetric(tt.args.in0, tt.args.req)

			assert.NoError(t, gotMetric.Validate())
			tt.wantErr(t, err)

			if !proto.Equal(gotMetric, tt.wantMetric) {
				t.Errorf("Service.UpdateMetric() = %v, want %v", gotMetric, tt.wantMetric)
			}
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
					_ = s.Create(&assessment.Metric{
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

			assert.NoError(t, gotMetric.Validate())
			tt.wantErr(t, err)
			if !proto.Equal(gotMetric, tt.wantMetric) {
				t.Errorf("Service.GetMetric() = %v, want %v", gotMetric, tt.wantMetric)
			}
		})
	}
}

func TestService_ListMetrics(t *testing.T) {
	timestamp := timestamppb.Now()

	type fields struct {
		cloudServiceHooks     []orchestrator.CloudServiceHookFunc
		toeHooks              []orchestrator.TargetOfEvaluationHookFunc
		AssessmentResultHooks []assessment.ResultHookFunc
		storage               persistence.Storage
		metricsFile           string
		loadMetricsFunc       func() ([]*assessment.Metric, error)
		catalogsFolder        string
		loadCatalogsFunc      func() ([]*orchestrator.Catalog, error)
		events                chan *orchestrator.MetricChangeEvent
		authz                 service.AuthorizationStrategy
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
						Name:        testdata.MockMetricName1,
						Description: testdata.MockMetricDescription1,
						Scale:       assessment.Metric_ORDINAL,
						Range:       MockMetricRange1,
					})
					_ = s.Create(&assessment.Metric{
						Id:              testdata.MockMetricID2,
						Name:            testdata.MockMetricName2,
						Description:     testdata.MockMetricDescription2,
						Scale:           assessment.Metric_ORDINAL,
						Range:           MockMetricRange1,
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
						Name:        testdata.MockMetricName1,
						Description: testdata.MockMetricDescription1,
						Scale:       assessment.Metric_ORDINAL,
						Range:       MockMetricRange1,
					})
					_ = s.Create(&assessment.Metric{
						Id:              testdata.MockMetricID2,
						Name:            testdata.MockMetricName2,
						Description:     testdata.MockMetricDescription2,
						Scale:           assessment.Metric_ORDINAL,
						Range:           MockMetricRange1,
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
						Name:        testdata.MockMetricName1,
						Description: testdata.MockMetricDescription1,
						Scale:       assessment.Metric_ORDINAL,
						Range:       MockMetricRange1,
					},
					{
						Id:              testdata.MockMetricID2,
						Name:            testdata.MockMetricName2,
						Description:     testdata.MockMetricDescription2,
						Scale:           assessment.Metric_ORDINAL,
						Range:           MockMetricRange1,
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
				cloudServiceHooks:     tt.fields.cloudServiceHooks,
				toeHooks:              tt.fields.toeHooks,
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				metricsFile:           tt.fields.metricsFile,
				loadMetricsFunc:       tt.fields.loadMetricsFunc,
				catalogsFolder:        tt.fields.catalogsFolder,
				loadCatalogsFunc:      tt.fields.loadCatalogsFunc,
				events:                tt.fields.events,
				authz:                 tt.fields.authz,
			}
			gotRes, err := svc.ListMetrics(tt.args.in0, tt.args.req)

			assert.NoError(t, gotRes.Validate())

			tt.wantErr(t, err)

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.ListMetrics() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_GetMetricImplementation(t *testing.T) {
	type fields struct {
		AssessmentResultHooks []assessment.ResultHookFunc
		storage               persistence.Storage
		metricsFile           string
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
				metricsFile:           tt.fields.metricsFile,
				catalogsFolder:        tt.fields.catalogsFolder,
				events:                tt.fields.events,
			}

			gotRes, err := svc.GetMetricImplementation(tt.args.ctx, tt.args.req)
			assert.NoError(t, gotRes.Validate())

			tt.wantErr(t, err)
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.GetMetricImplementation() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_UpdateMetricImplementation(t *testing.T) {
	type fields struct {
		AssessmentResultHooks []assessment.ResultHookFunc
		storage               persistence.Storage
		metricsFile           string
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
		wantImpl assert.ValueAssertionFunc
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "Invalid input",
			fields: fields{
				storage:     testutil.NewInMemoryStorage(t),
				metricsFile: "metrics.json",
			},
			args: args{
				req: &orchestrator.UpdateMetricImplementationRequest{
					Implementation: &assessment.MetricImplementation{MetricId: "notfound"},
				},
			},
			wantImpl: assert.Empty,
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
			wantImpl: assert.Empty,
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
			wantImpl: assert.Empty,
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
			wantImpl: assert.Empty,
			wantErr:  wantStatusCode(codes.Internal),
		},
		{
			name: "update",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					_ = s.Create(&assessment.Metric{Id: "TransportEncryptionEnabled"})
				}),
				metricsFile: "metrics.json",
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
			wantImpl: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				var impl = i1.(*assessment.MetricImplementation)

				return assert.Equal(t, "TransportEncryptionEnabled", impl.MetricId) &&
					assert.Equal(t, assessment.MetricImplementation_LANGUAGE_REGO, impl.Lang) &&
					assert.Equal(t, "package example", impl.Code) &&
					assert.True(t, impl.UpdatedAt.AsTime().Before(time.Now())) &&
					assert.NoError(t, impl.Validate())
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				metricsFile:           tt.fields.metricsFile,
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
		metricsFile           string
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
		want         assert.ValueAssertionFunc
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "Invalid input",
			args: args{
				req: &orchestrator.GetMetricConfigurationRequest{
					CloudServiceId: "InvalidCloudServiceID",
				},
			},
			want:    assert.Empty,
			wantErr: wantStatusCode(codes.InvalidArgument),
		},
		{
			name: "Permission denied",
			args: args{
				req: &orchestrator.GetMetricConfigurationRequest{
					MetricId:       testdata.MockMetricID1,
					CloudServiceId: testdata.MockCloudServiceID1,
				},
			},
			fields: fields{
				authz: &servicetest.AuthorizationStrategyMock{},
			},
			want:    assert.Empty,
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
					MetricId:       testdata.MockMetricID1,
					CloudServiceId: testdata.MockCloudServiceID1,
				},
			},
			want: assert.Empty,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrSomeError.Error())
			},
		},
		{
			name: "metric found",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					_ = s.Create(&assessment.Metric{Id: testdata.MockMetricID1})
					_ = s.Create(&orchestrator.CloudService{
						Id: testdata.MockCloudServiceID1,
					})
					_ = s.Create(&assessment.MetricConfiguration{
						MetricId:       testdata.MockMetricID1,
						CloudServiceId: testdata.MockCloudServiceID1,
						Operator:       "==",
						TargetValue:    &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: "1111"}},
					})
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.GetMetricConfigurationRequest{
					MetricId:       testdata.MockMetricID1,
					CloudServiceId: testdata.MockCloudServiceID1,
				},
			},
			wantResponse: &assessment.MetricConfiguration{
				MetricId:       testdata.MockMetricID1,
				CloudServiceId: testdata.MockCloudServiceID1,
				Operator:       "==",
				TargetValue:    &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: "1111"}},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				resp := i1.(*assessment.MetricConfiguration)
				wantResp := i2[0].(*assessment.MetricConfiguration)

				// We have to check the TargetValue independently and delete it, because otherwise DeepEqual throws an error.
				assert.Equal(t, resp.TargetValue.GetStringValue(), wantResp.TargetValue.GetStringValue())
				resp.TargetValue = nil
				wantResp.TargetValue = nil

				if !reflect.DeepEqual(resp, wantResp) {
					t.Errorf("Service.GetMetricConfiguration() = %v, want %v", resp, wantResp)
					return false
				}

				return true
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
					MetricId:       "NotExists",
					CloudServiceId: testdata.MockCloudServiceID1,
				},
			},
			want:    assert.Empty,
			wantErr: wantStatusCode(codes.NotFound),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				metricsFile:           tt.fields.metricsFile,
				catalogsFolder:        tt.fields.catalogsFolder,
				events:                tt.fields.events,
				authz:                 tt.fields.authz,
			}
			gotResponse, err := s.GetMetricConfiguration(tt.args.in0, tt.args.req)
			assert.NoError(t, gotResponse.Validate())

			tt.wantErr(t, err, tt.args)

			tt.want(t, gotResponse, tt.wantResponse)
		})
	}
}

func TestService_ListMetricConfigurations(t *testing.T) {
	type fields struct {
		AssessmentResultHooks []assessment.ResultHookFunc
		storage               persistence.Storage
		metricsFile           string
		catalogsFolder        string
		events                chan *orchestrator.MetricChangeEvent
		authz                 service.AuthorizationStrategy
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
		wantErr      assert.ErrorAssertionFunc
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
					CloudServiceId: testdata.MockCloudServiceID1,
				},
			},
			wantResponse: nil,
			wantErr:      wantStatusCode(codes.PermissionDenied),
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
					CloudServiceId: testdata.MockCloudServiceID1,
				},
			},
			wantErr: wantStatusCode(codes.Internal),
		},
		{
			name: "no error",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					_ = s.Create(&assessment.Metric{Id: testdata.MockMetricID1})
					_ = s.Create(&orchestrator.CloudService{
						Id: testdata.MockCloudServiceID1,
					})
					_ = s.Create(MockMetricConfiguration1)
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.ListMetricConfigurationRequest{
					CloudServiceId: testdata.MockCloudServiceID1,
				},
			},
			wantResponse: &orchestrator.ListMetricConfigurationResponse{
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
				metricsFile:           tt.fields.metricsFile,
				catalogsFolder:        tt.fields.catalogsFolder,
				events:                tt.fields.events,
				authz:                 tt.fields.authz,
			}
			gotResponse, err := svc.ListMetricConfigurations(tt.args.ctx, tt.args.req)

			// Check if response validation succeds
			assert.NoError(t, gotResponse.Validate())

			tt.wantErr(t, err)

			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("Service.ListMetricConfigurations() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestService_UpdateMetricConfiguration(t *testing.T) {
	type fields struct {
		AssessmentResultHooks []assessment.ResultHookFunc
		storage               persistence.Storage
		metricsFile           string
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
		want    assert.ValueAssertionFunc
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Metric configuration invalid",
			args: args{
				req: &orchestrator.UpdateMetricConfigurationRequest{
					CloudServiceId: testdata.MockCloudServiceID1,
					MetricId:       testdata.MockMetricID1,
					Configuration: &assessment.MetricConfiguration{
						Operator:       "invalidOperator",
						TargetValue:    testdata.MockMetricConfigurationTargetValueString,
						MetricId:       testdata.MockMetricID1,
						CloudServiceId: testdata.MockCloudServiceID1,
					},
				},
			},
			want: assert.Empty,
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
					CloudServiceId: testdata.MockCloudServiceID1,
					MetricId:       testdata.MockMetricID1,
					Configuration: &assessment.MetricConfiguration{
						CloudServiceId: testdata.MockCloudServiceID1,
						MetricId:       testdata.MockMetricID1,
						Operator:       "<",
						TargetValue:    testdata.MockMetricConfigurationTargetValueString,
					},
				},
			},
			want: assert.Empty,
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
					CloudServiceId: testdata.MockCloudServiceID1,
					MetricId:       testdata.MockMetricID1,
					Configuration: &assessment.MetricConfiguration{
						CloudServiceId: testdata.MockCloudServiceID1,
						MetricId:       testdata.MockMetricID1,
						Operator:       "<",
						TargetValue:    testdata.MockMetricConfigurationTargetValueString,
					},
				},
			},
			want: assert.Empty,
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
					CloudServiceId: testdata.MockCloudServiceID1,
					Configuration: &assessment.MetricConfiguration{
						Operator:       "<",
						TargetValue:    testdata.MockMetricConfigurationTargetValueString,
						MetricId:       testdata.MockMetricID1,
						CloudServiceId: testdata.MockCloudServiceID1,
					},
				},
			},
			want: assert.Empty,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "at least")
			},
		},
		{
			name: "cloudServiceID is missing in request",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.UpdateMetricConfigurationRequest{
					MetricId: testdata.MockMetricID1,
					Configuration: &assessment.MetricConfiguration{
						Operator:       "<",
						TargetValue:    testdata.MockMetricConfigurationTargetValueString,
						MetricId:       testdata.MockMetricID1,
						CloudServiceId: testdata.MockCloudServiceID1,
					},
				},
			},
			want: assert.Empty,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "CloudServiceId: value must be a valid UUID")
			},
		},
		{
			name: "cloudServiceID is invalid in request",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.UpdateMetricConfigurationRequest{
					CloudServiceId: "00000000-000000000000",
					MetricId:       testdata.MockMetricID1,
					Configuration: &assessment.MetricConfiguration{
						Operator:       "<",
						TargetValue:    testdata.MockMetricConfigurationTargetValueString,
						MetricId:       testdata.MockMetricID1,
						CloudServiceId: testdata.MockCloudServiceID1,
					},
				},
			},
			want: assert.Empty,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "CloudServiceId: value must be a valid UUID")
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
					MetricId:       testdata.MockMetricID1,
					CloudServiceId: testdata.MockCloudServiceID1,
				},
			},
			want: assert.Empty,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "Configuration: value is required")
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
					MetricId:       testdata.MockMetricID1,
					CloudServiceId: testdata.MockCloudServiceID1,
					Configuration: &assessment.MetricConfiguration{
						Operator:       "<",
						TargetValue:    testdata.MockMetricConfigurationTargetValueString,
						CloudServiceId: testdata.MockCloudServiceID1,
					},
				},
			},
			want: assert.Empty,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "at least")
			},
		},
		{
			name: "cloudServiceId is missing in configuration",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.UpdateMetricConfigurationRequest{
					MetricId:       testdata.MockMetricID1,
					CloudServiceId: testdata.MockCloudServiceID1,
					Configuration: &assessment.MetricConfiguration{
						Operator:    "<",
						TargetValue: testdata.MockMetricConfigurationTargetValueString,
						MetricId:    testdata.MockMetricID1,
					},
				},
			},
			want: assert.Empty,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "MetricConfiguration.CloudServiceId: value must be a valid UUID")
			},
		},
		{
			name: "cloudServiceId is invalid in configuration",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.UpdateMetricConfigurationRequest{
					MetricId:       testdata.MockMetricID1,
					CloudServiceId: testdata.MockCloudServiceID1,
					Configuration: &assessment.MetricConfiguration{
						Operator:       "<",
						TargetValue:    testdata.MockMetricConfigurationTargetValueString,
						MetricId:       testdata.MockMetricID1,
						CloudServiceId: "00000000-000000000000",
					},
				},
			},
			want: assert.Empty,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "MetricConfiguration.CloudServiceId: value must be a valid UUID")
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
					CloudServiceId: testdata.MockCloudServiceID1,
					MetricId:       testdata.MockMetricID1,
					Configuration: &assessment.MetricConfiguration{
						Operator:       "<",
						TargetValue:    testdata.MockMetricConfigurationTargetValueString,
						MetricId:       testdata.MockMetricID1,
						CloudServiceId: testdata.MockCloudServiceID1,
					},
				},
			},
			want: assert.Empty,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "metric or service does not exist")
			},
		},
		{
			name: "cloudService does not exist in storage",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					_ = s.Create(&assessment.Metric{Id: testdata.MockMetricID1})
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.UpdateMetricConfigurationRequest{
					CloudServiceId: testdata.MockCloudServiceID1,
					MetricId:       testdata.MockMetricID1,
					Configuration: &assessment.MetricConfiguration{
						Operator:       "<",
						TargetValue:    testdata.MockMetricConfigurationTargetValueString,
						MetricId:       testdata.MockMetricID1,
						CloudServiceId: testdata.MockCloudServiceID1,
					},
				},
			},
			want: assert.Empty,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "metric or service does not exist")
			},
		},
		{
			name: "append metric configuration",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					_ = s.Create(&assessment.Metric{Id: testdata.MockMetricID1})
					_ = s.Create(&orchestrator.CloudService{Id: testdata.MockCloudServiceID1})
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.UpdateMetricConfigurationRequest{
					CloudServiceId: testdata.MockCloudServiceID1,
					MetricId:       testdata.MockMetricID1,
					Configuration: &assessment.MetricConfiguration{
						CloudServiceId: testdata.MockCloudServiceID1,
						MetricId:       testdata.MockMetricID1,
						Operator:       "<",
						TargetValue:    testdata.MockMetricConfigurationTargetValueString,
					},
				},
			},
			wantErr: assert.NoError,
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				svc := i2[0].(*Service)

				var config *assessment.MetricConfiguration
				err := svc.storage.Get(&config, persistence_gorm.WithoutPreload(), "cloud_service_id = ? AND metric_id = ?", testdata.MockCloudServiceID1, testdata.MockMetricID1)
				if !assert.NoError(t, err) {
					return false
				}

				return assert.True(t, proto.Equal(config, i1.(proto.Message)))
			},
		},
		{
			name: "replace metric configuration",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					_ = s.Create(&assessment.Metric{Id: testdata.MockMetricID1})
					_ = s.Create(&orchestrator.CloudService{
						Id: testdata.MockCloudServiceID1,
					})
					_ = s.Create(&assessment.MetricConfiguration{
						MetricId:       testdata.MockMetricID1,
						CloudServiceId: testdata.MockCloudServiceID1,
						Operator:       ">",
					})
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				req: &orchestrator.UpdateMetricConfigurationRequest{
					CloudServiceId: testdata.MockCloudServiceID1,
					MetricId:       testdata.MockMetricID1,
					Configuration: &assessment.MetricConfiguration{
						CloudServiceId: testdata.MockCloudServiceID1,
						MetricId:       testdata.MockMetricID1,
						Operator:       "<",
						TargetValue:    testdata.MockMetricConfigurationTargetValueString,
					},
				},
			},
			wantErr: assert.NoError,
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				svc := i2[0].(*Service)

				var config *assessment.MetricConfiguration
				err := svc.storage.Get(&config, persistence_gorm.WithoutPreload(), "cloud_service_id = ? AND metric_id = ?", testdata.MockCloudServiceID1, testdata.MockMetricID1)
				if !assert.NoError(t, err) {
					return false
				}

				return assert.True(t, proto.Equal(config, i1.(proto.Message)))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				metricsFile:           tt.fields.metricsFile,
				loadMetricsFunc:       tt.fields.loadMetricsFunc,
				catalogsFolder:        tt.fields.catalogsFolder,
				loadCatalogsFunc:      tt.fields.loadCatalogsFunc,
				events:                tt.fields.events,
				authz:                 tt.fields.authz,
			}
			gotRes, err := svc.UpdateMetricConfiguration(tt.args.in0, tt.args.req)
			assert.NoError(t, gotRes.Validate())

			tt.wantErr(t, err)

			tt.want(t, gotRes, svc)
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
		wantRes assert.ValueAssertionFunc
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
			wantRes: assert.Nil,
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
			wantRes: assert.Nil,
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
			wantRes: assert.Nil,
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
			wantRes: assert.Nil,
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
			wantRes: func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				assert.NotNil(t, i)
				_, ok := i.(*emptypb.Empty)
				assert.True(t, ok)

				assert.NotNil(t, i2)
				s, ok := i2[0].(persistence.Storage)
				assert.True(t, ok)

				var gotMetric *assessment.Metric

				err := s.Get(&gotMetric, "id = ?", testdata.MockMetricID1)
				assert.NoError(t, err)

				return assert.NotEmpty(t, gotMetric.DeprecatedSince)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				storage: tt.fields.storage,
			}

			res, err := svc.RemoveMetric(context.TODO(), tt.args.req)

			// Run ErrorAssertionFunc
			tt.wantErr(t, err)

			tt.wantRes(t, res, tt.fields.storage)
		})
	}
}
