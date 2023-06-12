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

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/internal/testutil/servicetest"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/persistence/gorm"
	"clouditor.io/clouditor/service"
)

var ErrSomeError = errors.New("some error")

func TestService_loadMetrics(t *testing.T) {
	type fields struct {
		AssessmentResultHooks []func(result *assessment.AssessmentResult, err error)
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
			name: "storage error",
			fields: fields{
				metricsFile: "metrics.json",
				storage:     &testutil.StorageWithError{SaveErr: ErrSomeError},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrSomeError)
			},
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
						Id:    "TransportEncryptionEnabled",
						Name:  "A slightly updated metric",
						Scale: assessment.Metric_NOMINAL,
						Range: &assessment.Range{Range: &assessment.Range_AllowedValues{}},
					},
				},
			},
			wantMetric: &assessment.Metric{
				Id:    "TransportEncryptionEnabled",
				Name:  "A slightly updated metric",
				Scale: assessment.Metric_NOMINAL,
				Range: &assessment.Range{Range: &assessment.Range_AllowedValues{}},
			},
			wantErr: assert.NoError,
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
						Id:    "DoesProbablyNotExist",
						Name:  "UpdateMetricName",
						Scale: assessment.Metric_NOMINAL,
						Range: &assessment.Range{Range: &assessment.Range_AllowedValues{}},
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
		wantErr    bool
	}{
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
			wantErr: false,
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
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				storage: tt.fields.storage,
			}
			gotMetric, err := svc.GetMetric(tt.args.in0, tt.args.req)

			assert.NoError(t, gotMetric.Validate())

			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetMetric() error = %v, wantErrMessage %v", err, tt.wantErr)
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
	// TODO(anatheka): Test failes because of incorrect metrics.json file
	// assert.NoError(t, response.Validate())
}

func TestService_GetMetricImplementation(t *testing.T) {
	type fields struct {
		AssessmentResultHooks []func(result *assessment.AssessmentResult, err error)
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
		wantErr bool
	}{
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
			wantErr: true,
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
			wantErr: true,
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

			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetMetricImplementation() error = %v, wantErrMessage %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.GetMetricImplementation() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_UpdateMetricImplementation(t *testing.T) {
	type fields struct {
		AssessmentResultHooks []func(result *assessment.AssessmentResult, err error)
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
					Implementation: &assessment.MetricImplementation{MetricId: "notfound"},
				},
			},
			wantImpl: assert.Empty,
			wantErr:  true,
		},
		{
			name: "storage error",
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
			wantErr:  true,
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
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.UpdateMetricImplementation() error = %v, wantErrMessage %v", err, tt.wantErr)
				return
			}

			tt.wantImpl(t, gotImpl)
		})
	}
}

func TestService_GetMetricConfiguration(t *testing.T) {
	type fields struct {
		AssessmentResultHooks []func(result *assessment.AssessmentResult, err error)
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
			want: assert.Empty,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				gotStatus, ok := status.FromError(err)
				if !ok {
					return false
				}
				return assert.Equal(t, gotStatus.Code(), codes.NotFound)
			},
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
		AssessmentResultHooks []func(result *assessment.AssessmentResult, err error)
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
		wantErr      bool
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
			wantErr: true,
		},
		{
			name: "no error",
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
					})
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
					testdata.MockMetricID1: {
						MetricId:       testdata.MockMetricID1,
						CloudServiceId: testdata.MockCloudServiceID1,
						Operator:       "==",
					},
				},
			},
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

			// TODO(anatheka): Test failes because of incorrect metrics.json file
			// assert.NoError(t, gotResponse.Validate())

			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ListMetricConfigurations() error = %v, wantErrMessage %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("Service.ListMetricConfigurations() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestService_UpdateMetricConfiguration(t *testing.T) {
	type fields struct {
		AssessmentResultHooks []func(result *assessment.AssessmentResult, err error)
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
				err := svc.storage.Get(&config, gorm.WithoutPreload(), "cloud_service_id = ? AND metric_id = ?", testdata.MockCloudServiceID1, testdata.MockMetricID1)
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
				err := svc.storage.Get(&config, gorm.WithoutPreload(), "cloud_service_id = ? AND metric_id = ?", testdata.MockCloudServiceID1, testdata.MockMetricID1)
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
