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
	"io/ioutil"
	"os"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const DefaultMetricPageSize = 50
const MaxMetricPageSize = 1000

// loadMetrics takes care of loading the metric definitions from the (embedded) metrics.json as
// well as the default metric implementations from the Rego files.
func (svc *Service) loadMetrics() (err error) {
	var (
		impl    *assessment.MetricImplementation
		metrics []*assessment.Metric
		b       []byte
	)

	b, err = f.ReadFile(svc.metricsFile)
	if err != nil {
		return fmt.Errorf("error while loading %s: %w", svc.metricsFile, err)
	}

	err = json.Unmarshal(b, &metrics)
	if err != nil {
		return fmt.Errorf("error in JSON marshal: %w", err)
	}

	svc.metrics = make(map[string]*assessment.Metric)
	defaultMetricConfigurations = make(map[string]*assessment.MetricConfiguration)

	for _, m := range metrics {
		// Load the Rego file
		file := fmt.Sprintf("policies/bundles/%s/metric.rego", m.Id)
		impl, err = loadMetricImplementation(m.Id, file)
		if err != nil {
			return fmt.Errorf("could not load metric implementation: %w", err)
		}

		// Save our metric implementation
		err = svc.storage.Save(impl, "metric_id = ?", m.Id)
		if err != nil {
			return fmt.Errorf("could not save metric implementation: %w", err)
		}

		// Look for the data.json to include default metric configurations
		fileName := fmt.Sprintf("policies/bundles/%s/data.json", m.Id)

		b, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Errorf("Could not retrieve default configuration for metric %s: %v. Ignoring metric", m.Id, err)
			continue
		}

		var config assessment.MetricConfiguration

		err = json.Unmarshal(b, &config)
		if err != nil {
			log.Errorf("Error in reading default configuration for metric %s: %v. Ignoring metric", m.Id, err)
			continue
		}

		config.IsDefault = true

		svc.metrics[m.Id] = m
		defaultMetricConfigurations[m.Id] = &config
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
	// Validate the metric request
	err = req.Metric.Validate(assessment.WithMetricRequiresId())
	if err != nil {
		newError := errors.New("validation of metric failed")
		log.Error(newError)
		return nil, status.Errorf(codes.InvalidArgument, "%v", newError)
	}

	// Check, if metric id already exists
	if _, ok := svc.metrics[req.Metric.Id]; ok {
		return nil, status.Error(codes.AlreadyExists, "metric already exists")
	}

	// Build a new metric out of the request
	metric = req.Metric

	// Append metric
	svc.metrics[req.Metric.Id] = metric

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
	var ok bool

	// Validate the metric request
	err = req.Metric.Validate()
	if err != nil {
		newError := fmt.Errorf("validation of metric failed: %w", err)
		log.Error(newError)
		return nil, status.Errorf(codes.InvalidArgument, "%v", newError)
	}

	// Check, if metric exists according to req.MetricId
	if metric, ok = svc.metrics[req.MetricId]; !ok {
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
		ok     bool
		metric *assessment.Metric
	)

	// TODO(oxisto): Validate the metric implementation request

	// Check, if metric exists according to req.MetricId
	if metric, ok = svc.metrics[req.MetricId]; !ok {
		err := fmt.Errorf("metric with identifier %s does not exist", req.MetricId)
		log.Error(err)
		return nil, status.Errorf(codes.NotFound, "%v", err)
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

	// Paginate the metrics according to the request
	res.Metrics, res.NextPageToken, err = service.PaginateMapValues(req, svc.metrics, func(a *assessment.Metric, b *assessment.Metric) bool {
		return a.Id < b.Id
	}, MaxMetricPageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate metrics: %v", err)
	}

	return res, nil
}

// GetMetric retrieves a metric specified by req.MetridId
func (svc *Service) GetMetric(_ context.Context, req *orchestrator.GetMetricRequest) (metric *assessment.Metric, err error) {
	var ok bool

	if metric, ok = svc.metrics[req.MetricId]; !ok {
		return nil, status.Errorf(codes.NotFound, "could not find metric with id %s", req.MetricId)
	}

	return
}

func (s *Service) GetMetricConfiguration(_ context.Context, req *orchestrator.GetMetricConfigurationRequest) (response *assessment.MetricConfiguration, err error) {
	// Check, if we have a specific configuration for the metric
	if config, ok := s.metricConfigurations[req.ServiceId][req.MetricId]; ok {
		return config, nil
	}

	// Otherwise, fall back to our default configuration
	if config, ok := defaultMetricConfigurations[req.MetricId]; ok {
		return config, nil
	}

	newError := fmt.Errorf("could not find metric configuration for metric %s in service %s", req.MetricId, req.ServiceId)
	log.Error(newError)

	return nil, status.Errorf(codes.NotFound, "%v", newError)
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

	// TODO(oxisto): This is not very efficient, we should do this once at startup so that we can just return the map
	for metricId := range svc.metrics {
		config, err := svc.GetMetricConfiguration(ctx, &orchestrator.GetMetricConfigurationRequest{ServiceId: req.ServiceId, MetricId: metricId})

		if err != nil {
			log.Errorf("Error getting metric configuration: %v", err)
			return nil, err
		}

		response.Configurations[metricId] = config
	}

	return
}

func (svc *Service) GetMetricImplementation(ctx context.Context, req *orchestrator.GetMetricImplementationRequest) (res *assessment.MetricImplementation, err error) {
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
