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
	"io/fs"
	"testing"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestLoadMetrics(t *testing.T) {
	var err = LoadMetrics("notfound.json")

	assert.ErrorIs(t, err, fs.ErrNotExist)

	err = LoadMetrics("metrics.json")

	assert.NoError(t, err)
}

func TestService_CreateMetric(t *testing.T) {
	type args struct {
		in0 context.Context
		req *orchestrator.CreateMetricRequest
	}
	tests := []struct {
		name       string
		args       args
		wantMetric *assessment.Metric
		wantErr    bool
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
			wantMetric: nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMetric, err := service.CreateMetric(tt.args.in0, tt.args.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("Service.CreateMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
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
		name       string
		args       args
		wantMetric *assessment.Metric
		wantErr    bool
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
				},
			},
			wantMetric: nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service = NewService()
			gotMetric, err := service.UpdateMetric(tt.args.in0, tt.args.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("Service.UpdateMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
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
			service = NewService()
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

	response, err = service.ListMetrics(context.TODO(), &orchestrator.ListMetricsRequest{})

	assert.NoError(t, err)
	assert.NotEmpty(t, response.Metrics)
}
