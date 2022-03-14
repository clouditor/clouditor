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
	"encoding/json"
	"errors"
	"fmt"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoadMetrics loads metrics definitions from a JSON file.
func LoadMetrics(metricsFile string) (err error) {
	var (
		b []byte
	)

	b, err = f.ReadFile(metricsFile)
	if err != nil {
		return fmt.Errorf("error while loading %s: %w", metricsFile, err)
	}

	err = json.Unmarshal(b, &metrics)
	if err != nil {
		return fmt.Errorf("error in JSON marshal: %w", err)
	}

	return nil
}

// CreateMetric creates a new metric in the database.
func (*Service) CreateMetric(_ context.Context, req *orchestrator.CreateMetricRequest) (metric *assessment.Metric, err error) {
	// Validate the metric request
	err = req.Metric.Validate(assessment.WithMetricRequiresId())
	if err != nil {
		newError := errors.New("validation of metric failed")
		log.Error(newError)
		return nil, status.Errorf(codes.InvalidArgument, "%v", newError)
	}

	// Check, if metric id already exists
	if _, ok := metricIndex[req.Metric.Id]; ok {
		return nil, status.Error(codes.AlreadyExists, "metric already exists")
	}

	// Build a new metric out of the request
	metric = req.Metric

	// Append metric
	metricIndex[req.Metric.Id] = metric
	metrics = append(metrics, metric)

	return
}

// UpdateMetric updates an existing metric, specified by the identifier in req.MetricId.
func (*Service) UpdateMetric(_ context.Context, req *orchestrator.UpdateMetricRequest) (metric *assessment.Metric, err error) {
	var ok bool

	// Validate the metric request
	err = req.Metric.Validate()
	if err != nil {
		newError := fmt.Errorf("validation of metric failed: %w", err)
		log.Error(newError)
		return nil, status.Errorf(codes.InvalidArgument, "%v", newError)
	}

	// Check, if metric exists according to req.MetricId
	if metric, ok = metricIndex[req.MetricId]; !ok {
		newError := fmt.Errorf("metric with identifier %s does not exist", req.MetricId)
		log.Error(newError)
		return nil, status.Errorf(codes.NotFound, "%v", newError)
	}

	// Update metric
	metric.Name = req.Metric.Name
	metric.Description = req.Metric.Description
	metric.Category = req.Metric.Category
	metric.Range = req.Metric.Range
	metric.Scale = req.Metric.Scale

	return
}

// ListMetrics lists all available metrics.
func (*Service) ListMetrics(_ context.Context, _ *orchestrator.ListMetricsRequest) (response *orchestrator.ListMetricsResponse, err error) {
	response = &orchestrator.ListMetricsResponse{
		Metrics: metrics,
	}

	return response, nil
}

// GetMetric retrieves a metric specified by req.MetridId
func (svc *Service) GetMetric(_ context.Context, req *orchestrator.GetMetricRequest) (metric *assessment.Metric, err error) {
	var ok bool

	if metric, ok = svc.metric(req.MetricId); !ok {
		return nil, status.Errorf(codes.NotFound, "could not find metric with id %s", req.MetricId)
	}

	return
}

func (svc *Service) metric(id string) (metric *assessment.Metric, ok bool) {
	metric, ok = metricIndex[id]
	return
}
