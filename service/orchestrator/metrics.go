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
	"io"
	"os"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/persistence/gorm"
	"clouditor.io/clouditor/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// loadMetrics takes care of loading the metric definitions from the (embedded) metrics.json as
// well as the default metric implementations from the Rego files.
func (svc *Service) loadMetrics() (err error) {
	var (
		metrics []*assessment.Metric
	)

	// Default to loading metrics from our embedded file system
	if svc.loadMetricsFunc == nil {
		svc.loadMetricsFunc = svc.loadEmbeddedMetrics
	}

	// Execute our metric loading function
	metrics, err = svc.loadMetricsFunc()
	if err != nil {
		return fmt.Errorf("could not load metrics: %w", err)
	}

	defaultMetricConfigurations = make(map[string]*assessment.MetricConfiguration)

	// Try to prepare the (initial) metric implementations and configurations. We can still continue if they should
	// fail, since they can still be updated later during runtime. Also, this makes it possible to load metrics of which
	// we intentionally do not have the implementation, because they are assess outside the Clouditor toolset, but we
	// still need to be aware of the particular metric.
	for _, m := range metrics {
		err = svc.prepareMetric(m)
		if err != nil {
			log.Warnf("Could not prepare implementation or default configuration for metric %s: %v", m.Id, err)
		}
	}

	err = svc.storage.Save(metrics)
	if err != nil {
		log.Errorf("Error while saving metrics: %v", err)
	}

	return
}

// prepareMetric takes care of the heavy lifting of loading the default implementation and configuration of a particular
// metric and storing them into the service.
func (svc *Service) prepareMetric(m *assessment.Metric) (err error) {
	var (
		config *assessment.MetricConfiguration
	)

	// Load the Rego file
	file := fmt.Sprintf("policies/bundles/%s/metric.rego", m.Id)
	m.Implementation, err = loadMetricImplementation(m.Id, file)
	if err != nil {
		return fmt.Errorf("could not load metric implementation: %w", err)
	}

	// Look for the data.json to include default metric configurations
	fileName := fmt.Sprintf("policies/bundles/%s/data.json", m.Id)

	// Load the default configuration file
	b, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("could not retrieve default configuration for metric %s: %w", m.Id, err)
	}

	err = json.Unmarshal(b, &config)
	if err != nil {
		return fmt.Errorf("error in reading default configuration for metric %s: %w", m.Id, err)
	}

	config.IsDefault = true

	defaultMetricConfigurations[m.Id] = config

	return
}

// loadEmbeddedMetrics loads the metric definitions from the embedded file system using the path specified in
// the service's metricsFile.
func (svc *Service) loadEmbeddedMetrics() (metrics []*assessment.Metric, err error) {
	var b []byte

	b, err = f.ReadFile(svc.metricsFile)
	if err != nil {
		return nil, fmt.Errorf("error while loading %s: %w", svc.metricsFile, err)
	}

	err = json.Unmarshal(b, &metrics)
	if err != nil {
		return nil, fmt.Errorf("error in JSON marshal: %w", err)
	}

	return
}

// loadMetricImplementation loads a metric implementation from a Rego file on a filesystem.
func loadMetricImplementation(metricID, file string) (impl *assessment.MetricImplementation, err error) {
	// Fetch the metric implementation directly from our file
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	impl = &assessment.MetricImplementation{
		MetricId: metricID,
		Lang:     assessment.MetricImplementation_REGO,
		Code:     string(b),
	}

	return
}

// CreateMetric creates a new metric in the database.
func (svc *Service) CreateMetric(_ context.Context, req *orchestrator.CreateMetricRequest) (metric *assessment.Metric, err error) {
	var count int64

	// Validate the metric request
	err = req.Metric.Validate(assessment.WithMetricRequiresId())
	if err != nil {
		newError := errors.New("validation of metric failed")
		log.Error(newError)
		return nil, status.Errorf(codes.InvalidArgument, "%v", newError)
	}

	// Check, if metric id already exists
	count, err = svc.storage.Count(metric, "id = ?", req.Metric.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	if count > 0 {
		return nil, status.Error(codes.AlreadyExists, "metric already exists")
	}

	// Build a new metric out of the request
	metric = req.Metric

	// Append metric
	err = svc.storage.Create(metric)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	// Notify event listeners
	go func() {
		svc.events <- &orchestrator.MetricChangeEvent{
			Type:     orchestrator.MetricChangeEvent_METADATA_CHANGED,
			MetricId: metric.Id,
		}
	}()

	return
}

// UpdateMetric updates an existing metric, specified by the identifier in req.MetricId.
func (svc *Service) UpdateMetric(_ context.Context, req *orchestrator.UpdateMetricRequest) (metric *assessment.Metric, err error) {
	// Validate the metric request
	err = req.Metric.Validate()
	if err != nil {
		newError := fmt.Errorf("validation of metric failed: %w", err)
		log.Error(newError)
		return nil, status.Errorf(codes.InvalidArgument, "%v", newError)
	}

	// Check, if metric exists according to req.MetricId
	err = svc.storage.Get(&metric, "id = ?", req.MetricId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "metric not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	// Update metric
	metric.Name = req.Metric.Name
	metric.Description = req.Metric.Description
	metric.Category = req.Metric.Category
	metric.Range = req.Metric.Range
	metric.Scale = req.Metric.Scale

	err = svc.storage.Save(metric, "id = ? ", metric.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	// Notify event listeners
	go func() {
		svc.events <- &orchestrator.MetricChangeEvent{
			Type:     orchestrator.MetricChangeEvent_METADATA_CHANGED,
			MetricId: metric.Id,
		}
	}()

	return
}

// UpdateMetricImplementation updates an existing metric implementation, specified by the identifier in req.MetricId.
func (svc *Service) UpdateMetricImplementation(_ context.Context, req *orchestrator.UpdateMetricImplementationRequest) (impl *assessment.MetricImplementation, err error) {
	var (
		metric *assessment.Metric
	)

	// TODO(oxisto): Validate the metric implementation request

	// Check, if metric exists according to req.MetricId
	err = svc.storage.Get(&metric, "id = ?", req.MetricId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "metric not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	// Update implementation
	impl = req.Implementation
	impl.MetricId = req.MetricId

	// Store it in the database
	err = svc.storage.Save(impl, "metric_id = ?", impl.MetricId)
	if err != nil {
		return nil, fmt.Errorf("could not save metric implementation: %w", err)
	}

	// Notify event listeners
	go func() {
		svc.events <- &orchestrator.MetricChangeEvent{
			Type:     orchestrator.MetricChangeEvent_IMPLEMENTATION_CHANGED,
			MetricId: metric.Id,
		}
	}()

	return
}

// ListMetrics lists all available metrics.
func (svc *Service) ListMetrics(_ context.Context, req *orchestrator.ListMetricsRequest) (res *orchestrator.ListMetricsResponse, err error) {
	res = new(orchestrator.ListMetricsResponse)

	// Validate the request
	if err = api.ValidateListRequest[*assessment.Metric](req); err != nil {
		err = fmt.Errorf("invalid request: %w", err)
		log.Error(err)
		err = status.Errorf(codes.InvalidArgument, "%v", err)
		return
	}

	// Paginate the metrics according to the request
	res.Metrics, res.NextPageToken, err = service.PaginateStorage[*assessment.Metric](req, svc.storage,
		service.DefaultPaginationOpts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate metrics: %v", err)
	}

	return res, nil
}

// GetMetric retrieves a metric specified by req.MetridId
func (svc *Service) GetMetric(_ context.Context, req *orchestrator.GetMetricRequest) (metric *assessment.Metric, err error) {
	err = svc.storage.Get(&metric, "id = ?", req.MetricId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "metric not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	return
}

func (svc *Service) GetMetricConfiguration(_ context.Context, req *orchestrator.GetMetricConfigurationRequest) (res *assessment.MetricConfiguration, err error) {
	res = new(assessment.MetricConfiguration)
	err = svc.storage.Get(res, gorm.WithoutPreload(), "cloud_service_id = ? AND metric_id = ?", req.CloudServiceId, req.MetricId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		// Otherwise, fall back to our default configuration
		if config, ok := defaultMetricConfigurations[req.MetricId]; ok {
			return config, nil
		}

		return nil, status.Errorf(codes.NotFound, "metric configuration not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	return
}

// UpdateMetricConfiguration updates the configuration for a metric, specified by the identifier in req.MetricId.
func (svc *Service) UpdateMetricConfiguration(_ context.Context, req *orchestrator.UpdateMetricConfigurationRequest) (res *assessment.MetricConfiguration, err error) {
	// TODO(oxisto): Validate the request

	// Make sure that the configuration also has metric/service ID set
	req.Configuration.CloudServiceId = req.CloudServiceId
	req.Configuration.MetricId = req.MetricId

	err = svc.storage.Save(&req.Configuration)
	if err != nil && errors.Is(err, persistence.ErrConstraintFailed) {
		return nil, status.Errorf(codes.NotFound, "metric or service does not exist")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	// Notify event listeners
	go func() {
		svc.events <- &orchestrator.MetricChangeEvent{
			Type:           orchestrator.MetricChangeEvent_CONFIG_CHANGED,
			CloudServiceId: req.CloudServiceId,
			MetricId:       req.MetricId,
		}
	}()

	// Update response
	res = req.Configuration

	return
}

// ListMetricConfigurations retrieves a list of MetricConfiguration objects for a particular target
// cloud service specified in req.
//
// The list MUST include a configuration for each known metric. If the user did not specify a custom
// configuration for a particular metric within the service, the default metric configuration is
// inserted into the list.
func (svc *Service) ListMetricConfigurations(ctx context.Context, req *orchestrator.ListMetricConfigurationRequest) (response *orchestrator.ListMetricConfigurationResponse, err error) {
	response = &orchestrator.ListMetricConfigurationResponse{
		Configurations: make(map[string]*assessment.MetricConfiguration),
	}

	var metrics []*assessment.Metric
	err = svc.storage.List(&metrics, "", true, 0, -1)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	// TODO(oxisto): This is not very efficient, we should do this once at startup so that we can just return the map
	for _, metric := range metrics {
		config, err := svc.GetMetricConfiguration(ctx, &orchestrator.GetMetricConfigurationRequest{CloudServiceId: req.CloudServiceId, MetricId: metric.Id})
		if err == nil {
			response.Configurations[metric.Id] = config
		}
	}

	return
}

func (svc *Service) GetMetricImplementation(_ context.Context, req *orchestrator.GetMetricImplementationRequest) (res *assessment.MetricImplementation, err error) {
	res = new(assessment.MetricImplementation)

	// TODO(oxisto): Validate GetMetricImplementationRequest
	err = svc.storage.Get(res, "metric_id = ?", req.MetricId)
	if err == persistence.ErrRecordNotFound {
		return nil, status.Error(codes.NotFound, "implementation for metric not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "could not retrieve metric implementation: %v", err)
	}

	return
}

// SubscribeMetricChangeEvents implements a stream of metric events to the subscribed client.
func (svc *Service) SubscribeMetricChangeEvents(_ *orchestrator.SubscribeMetricChangeEventRequest, stream orchestrator.Orchestrator_SubscribeMetricChangeEventsServer) (err error) {
	var (
		event *orchestrator.MetricChangeEvent
	)

	// TODO(oxisto): Do we also need a (empty) recv func again?

	for {
		// TODO(oxisto): Does this work for multiple subcribers/readers or do we need a channel each?
		// Wait for a new event in our event channel
		event = <-svc.events

		err = stream.Send(event)

		// Check for send errors
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			err = fmt.Errorf("cannot stream response to the client: %w", err)
			log.Error(err)

			return status.Errorf(codes.Unknown, "%v", err)
		}
	}
}
