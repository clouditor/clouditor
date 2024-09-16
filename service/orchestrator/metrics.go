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

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/logging"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/persistence/gorm"
	"clouditor.io/clouditor/v2/service"
	"github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ErrMetricNotFound indicates the certification was not found
var ErrMetricNotFound = status.Error(codes.NotFound, "metric not found")

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
		// Somehow, we first need to create the metric, otherwise we are running into weird constraint issues.
		// We use Save() instead of Create()because it may be that the metrics already exist in the database and just need to be updated.
		err = svc.storage.Save(m, "id = ? ", m.Id)
		if err != nil {
			log.Errorf("Error while saving metric `%s`: %v", m.Id, err)
			continue
		}

		err = prepareMetric(m)
		if err != nil {
			log.Warnf("Could not prepare implementation or default configuration for metric %s: %v", m.Id, err)
			continue
		}

		err = svc.storage.Save(m.Implementation, "metric_id = ? ", m.Id)
		if err != nil {
			log.Errorf("Error while saving metric implementation for '%s': %v", m.Id, err)
			continue
		}

		log.Debugf("Metric loaded with id '%s'.", m.GetId())

	}

	// Here we have return nil, as the previous errors are only a warning and not a real error for the calling function.
	return nil
}

// prepareMetric takes care of the heavy lifting of loading the default implementation and configuration of a particular
// metric and storing them into the service.
func prepareMetric(m *assessment.Metric) (err error) {
	var (
		config *assessment.MetricConfiguration
	)

	// Load the Rego file
	file := fmt.Sprintf("policies/bundles/%s/%s/metric.rego", m.CategoryID(), m.Id)
	m.Implementation, err = loadMetricImplementation(m.Id, file)
	if err != nil {
		return fmt.Errorf("could not load metric implementation: %w", err)
	}

	// Look for the data.json to include default metric configurations
	fileName := fmt.Sprintf("policies/bundles/%s/%s/data.json", m.CategoryID(), m.Id)

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
	config.MetricId = m.GetId()

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
		MetricId:  metricID,
		Lang:      assessment.MetricImplementation_LANGUAGE_REGO,
		Code:      string(b),
		UpdatedAt: timestamppb.Now(),
	}

	return
}

// CreateMetric creates a new metric in the database.
func (svc *Service) CreateMetric(_ context.Context, req *orchestrator.CreateMetricRequest) (metric *assessment.Metric, err error) {
	var count int64

	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
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

	if metric.DeprecatedSince != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "the metric shouldn't be set to deprecated at creation time")
	}

	// Append metric
	err = svc.storage.Create(metric)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	// Notify event listeners
	go func() {
		svc.events <- &orchestrator.MetricChangeEvent{
			Type:     orchestrator.MetricChangeEvent_TYPE_METADATA_CHANGED,
			MetricId: metric.Id,
		}
	}()

	logging.LogRequest(log, logrus.DebugLevel, logging.Create, req)

	return
}

// UpdateMetric updates an existing metric, specified by the identifier in req.MetricId.
func (svc *Service) UpdateMetric(_ context.Context, req *orchestrator.UpdateMetricRequest) (metric *assessment.Metric, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	metric = new(assessment.Metric)

	// Check, if metric exists according to req.Metric.Id
	err = svc.storage.Get(&metric, "id = ?", req.Metric.Id)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, ErrMetricNotFound
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	// Update metric
	metric.Name = req.Metric.Name
	metric.Description = req.Metric.Description
	metric.Category = req.Metric.Category
	metric.Range = req.Metric.Range
	metric.Scale = req.Metric.Scale
	metric.DeprecatedSince = req.Metric.DeprecatedSince

	err = svc.storage.Save(metric, "id = ? ", metric.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	// Notify event listeners
	go func() {
		svc.events <- &orchestrator.MetricChangeEvent{
			Type:     orchestrator.MetricChangeEvent_TYPE_METADATA_CHANGED,
			MetricId: metric.Id,
		}
	}()

	logging.LogRequest(log, logrus.DebugLevel, logging.Update, req)

	return
}

// UpdateMetricImplementation updates an existing metric implementation, specified by the identifier in req.MetricId.
func (svc *Service) UpdateMetricImplementation(_ context.Context, req *orchestrator.UpdateMetricImplementationRequest) (impl *assessment.MetricImplementation, err error) {
	var (
		metric *assessment.Metric
	)

	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check, if metric exists according to the metric ID
	err = svc.storage.Get(&metric, "id = ?", req.Implementation.MetricId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, ErrMetricNotFound
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	// Update implementation
	impl = req.Implementation
	impl.UpdatedAt = timestamppb.Now()

	// Store it in the database
	err = svc.storage.Save(impl, "metric_id = ?", impl.MetricId)
	if err != nil {
		return nil, fmt.Errorf("could not save metric implementation: %w", err)
	}

	// Notify event listeners
	go func() {
		svc.events <- &orchestrator.MetricChangeEvent{
			Type:     orchestrator.MetricChangeEvent_TYPE_IMPLEMENTATION_CHANGED,
			MetricId: metric.Id,
		}
	}()

	logging.LogRequest(log, logrus.DebugLevel, logging.Update, req)

	return
}

// ListMetrics lists all available metrics.
func (svc *Service) ListMetrics(_ context.Context, req *orchestrator.ListMetricsRequest) (res *orchestrator.ListMetricsResponse, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	res = new(orchestrator.ListMetricsResponse)
	var conds []any

	// Add the deprecated metrics as well if requested
	if !req.Filter.GetIncludeDeprecated() {
		conds = append(conds, "deprecated_since IS NULL")
	}

	// Paginate the metrics according to the request
	res.Metrics, res.NextPageToken, err = service.PaginateStorage[*assessment.Metric](req, svc.storage,
		service.DefaultPaginationOpts, conds...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate metrics: %v", err)
	}

	return res, nil
}

// RemoveMetric removes a metric specified by req.MetricId. The metric is not deleted, but the property deprecated is set to true for backward compatibility reasons.
func (svc *Service) RemoveMetric(ctx context.Context, req *orchestrator.RemoveMetricRequest) (res *emptypb.Empty, err error) {
	var (
		metric *assessment.Metric
	)

	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check, if metric exists according to the metric ID
	err = svc.storage.Get(&metric, "id = ?", req.MetricId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, ErrMetricNotFound
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	// Set timestamp if not already set
	if metric.DeprecatedSince == nil {
		metric.DeprecatedSince = timestamppb.Now()
	}

	// Update metric with property deprecated is true
	err = svc.storage.Update(metric, "Id = ?", req.MetricId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	logging.LogRequest(log, logrus.DebugLevel, logging.Remove, req)

	return &emptypb.Empty{}, nil
}

// GetMetric retrieves a metric specified by req.MetricId.
func (svc *Service) GetMetric(_ context.Context, req *orchestrator.GetMetricRequest) (metric *assessment.Metric, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	err = svc.storage.Get(&metric, "id = ?", req.MetricId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, ErrMetricNotFound
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	return
}

func (svc *Service) GetMetricConfiguration(ctx context.Context, req *orchestrator.GetMetricConfigurationRequest) (res *assessment.MetricConfiguration, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the cloud service according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessRead, req) {
		return nil, service.ErrPermissionDenied
	}

	res = new(assessment.MetricConfiguration)

	err = svc.storage.Get(res, gorm.WithoutPreload(), "cloud_service_id = ? AND metric_id = ?", req.CertificationTargetId, req.MetricId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		// Otherwise, fall back to our default configuration
		if config, ok := defaultMetricConfigurations[req.MetricId]; ok {
			// Copy the metric configuration and set the cloud service id
			newConfig := &assessment.MetricConfiguration{
				Operator:              config.GetOperator(),
				TargetValue:           config.GetTargetValue(),
				IsDefault:             config.GetIsDefault(),
				UpdatedAt:             config.GetUpdatedAt(),
				MetricId:              config.GetMetricId(),
				CertificationTargetId: req.GetCertificationTargetId(),
			}

			return newConfig, nil
		}

		return nil, status.Errorf(codes.NotFound, "metric configuration not found for metric with id '%s'", req.GetMetricId())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	return
}

// UpdateMetricConfiguration updates the configuration for a metric, specified by the identifier in req.MetricId.
func (svc *Service) UpdateMetricConfiguration(ctx context.Context, req *orchestrator.UpdateMetricConfigurationRequest) (res *assessment.MetricConfiguration, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the cloud service according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessRead, req) {
		return nil, service.ErrPermissionDenied
	}

	// Make sure that the configuration also has updatedAt and isDefault set
	req.Configuration.UpdatedAt = timestamppb.Now()
	req.Configuration.IsDefault = false

	err = svc.storage.Save(&req.Configuration, "metric_id = ? AND cloud_service_id = ?", req.GetMetricId(), req.GetCertificationTargetId())
	if err != nil && errors.Is(err, persistence.ErrConstraintFailed) {
		return nil, status.Errorf(codes.NotFound, "metric or service does not exist")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	// Notify event listeners
	go func() {
		svc.events <- &orchestrator.MetricChangeEvent{
			Type:                  orchestrator.MetricChangeEvent_TYPE_CONFIG_CHANGED,
			CertificationTargetId: req.CertificationTargetId,
			MetricId:              req.MetricId,
		}
	}()

	// Update response
	res = req.Configuration

	logging.LogRequest(log, logrus.DebugLevel, logging.Update, req)

	return
}

// ListMetricConfigurations retrieves a list of MetricConfiguration objects for a particular target
// cloud service specified in req.
//
// The list MUST include a configuration for each known metric. If the user did not specify a custom
// configuration for a particular metric within the service, the default metric configuration is
// inserted into the list.
func (svc *Service) ListMetricConfigurations(ctx context.Context, req *orchestrator.ListMetricConfigurationRequest) (response *orchestrator.ListMetricConfigurationResponse, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the cloud service according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessRead, req) {
		return nil, service.ErrPermissionDenied
	}

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
		config, err := svc.GetMetricConfiguration(ctx, &orchestrator.GetMetricConfigurationRequest{CertificationTargetId: req.CertificationTargetId, MetricId: metric.Id})
		if err == nil {
			response.Configurations[metric.Id] = config
		}
	}

	return
}

// GetMetricImplementation retrieves a metric implementation specified by req.MetricId.
func (svc *Service) GetMetricImplementation(_ context.Context, req *orchestrator.GetMetricImplementationRequest) (res *assessment.MetricImplementation, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	res = new(assessment.MetricImplementation)

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
