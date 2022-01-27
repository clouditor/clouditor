// Copyright 2021 Fraunhofer AISEC
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
	"embed"
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"io/ioutil"
	"sync"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/persistence"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

//go:embed metrics.json
var f embed.FS

var metrics []*assessment.Metric
var metricIndex map[string]*assessment.Metric
var defaultMetricConfigurations map[string]*assessment.MetricConfiguration
var log *logrus.Entry

var DefaultMetricsFile = "metrics.json"

// Service is an implementation of the Clouditor Orchestrator service
type Service struct {
	orchestrator.UnimplementedOrchestratorServer

	// metricConfigurations holds a double-map of metric configurations associated first by service ID and then metric ID
	metricConfigurations map[string]map[string]*assessment.MetricConfiguration

	// Currently only in-memory
	results map[string]*assessment.AssessmentResult

	// Hook
	AssessmentResultHooks []func(result *assessment.AssessmentResult, err error)
	// mu is used for (un)locking result hook calls
	mu sync.Mutex

	db *gorm.DB
}

func init() {
	log = logrus.WithField("component", "orchestrator")
}

func NewService() *Service {
	s := Service{
		results:              make(map[string]*assessment.AssessmentResult),
		metricConfigurations: make(map[string]map[string]*assessment.MetricConfiguration),
	}

	if err := LoadMetrics(DefaultMetricsFile); err != nil {
		log.Errorf("Could not load embedded metrics. Will continue with empty metric list: %v", err)
	}

	metricIndex = make(map[string]*assessment.Metric)
	defaultMetricConfigurations = make(map[string]*assessment.MetricConfiguration)

	for _, m := range metrics {
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

		metricIndex[m.Id] = m
		defaultMetricConfigurations[m.Id] = &config
	}

	s.db = persistence.GetDatabase()

	return &s
}

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

func (*Service) ListMetrics(_ context.Context, _ *orchestrator.ListMetricsRequest) (response *orchestrator.ListMetricsResponse, err error) {
	response = &orchestrator.ListMetricsResponse{
		Metrics: metrics,
	}

	return response, nil
}

func (*Service) GetMetric(_ context.Context, request *orchestrator.GetMetricsRequest) (response *assessment.Metric, err error) {
	var ok bool
	var metric *assessment.Metric

	if metric, ok = metricIndex[request.MetricId]; !ok {
		return nil, status.Errorf(codes.NotFound, "could not find metric with id %s", request.MetricId)
	}

	response = metric

	return response, nil
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

	return nil, status.Errorf(codes.NotFound, "could not find metric configuration for metric %s in service %s", req.MetricId, req.ServiceId)
}

// ListMetricConfigurations retrieves a list of MetricConfiguration objects for a particular target
// cloud service specified in req.
//
// The list MUST include a configuration for each known metric. If the user did not specify a custom
// configuration for a particular metric within the service, the default metric configuration is
// inserted into the list.
func (s *Service) ListMetricConfigurations(ctx context.Context, req *orchestrator.ListMetricConfigurationRequest) (response *orchestrator.ListMetricConfigurationResponse, err error) {
	response = &orchestrator.ListMetricConfigurationResponse{
		Configurations: make(map[string]*assessment.MetricConfiguration),
	}

	// TODO(oxisto): This is not very efficient, we should do this once at startup so that we can just return the map
	for metricId := range metricIndex {
		config, err := s.GetMetricConfiguration(ctx, &orchestrator.GetMetricConfigurationRequest{ServiceId: req.ServiceId, MetricId: metricId})

		if err != nil {
			return nil, err
		}

		response.Configurations[metricId] = config
	}

	return
}

// StoreAssessmentResult is a method implementation of the orchestrator interface: It receives an assessment result and stores it
func (s *Service) StoreAssessmentResult(_ context.Context, req *orchestrator.StoreAssessmentResultRequest) (resp *orchestrator.StoreAssessmentResultResponse, err error) {

	resp = &orchestrator.StoreAssessmentResultResponse{}

	_, err = req.Result.Validate()

	if err != nil {
		log.Errorf("Invalid assessment result: %v", err)
		newError := fmt.Errorf("invalid assessment result: %w", err)

		go s.informHook(nil, newError)

		return resp, status.Errorf(codes.InvalidArgument, "invalid req: %v", err)
	}

	// TODO(all): We do not check ID in Validate. Therefore, I assume that we have to set ID here?
	s.results[req.Result.Id] = req.Result

	go s.informHook(req.Result, nil)

	return
}

func (s *Service) StoreAssessmentResults(stream orchestrator.Orchestrator_StoreAssessmentResultsServer) (err error) {
	var result *assessment.AssessmentResult

	for {
		result, err = stream.Recv()

		if err != nil {
			// If no more input of the stream is available, return SendAndClose `error`
			if err == io.EOF {
				log.Infof("Stopped receiving streamed assessment results")
				return stream.SendAndClose(&emptypb.Empty{})
			}

			return err
		}

		// Call StoreAssessmentResult() for storing a single assessment
		storeAssessmentResultReq := &orchestrator.StoreAssessmentResultRequest{
			Result: result,
		}

		_, err = s.StoreAssessmentResult(context.Background(), storeAssessmentResultReq)
		if err != nil {
			return err
		}
	}

}

// ListAssessmentResults is a method implementation of the orchestrator interface
func (s *Service) ListAssessmentResults(_ context.Context, _ *assessment.ListAssessmentResultsRequest) (res *assessment.ListAssessmentResultsResponse, err error) {
	res = new(assessment.ListAssessmentResultsResponse)
	res.Results = []*assessment.AssessmentResult{}

	for _, result := range s.results {
		res.Results = append(res.Results, result)
	}

	return
}

func (s *Service) RegisterAssessmentResultHook(hook func(result *assessment.AssessmentResult, err error)) {
	s.AssessmentResultHooks = append(s.AssessmentResultHooks, hook)
}

func (s *Service) informHook(result *assessment.AssessmentResult, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Inform our hook, if we have any
	if s.AssessmentResultHooks != nil {
		for _, hook := range s.AssessmentResultHooks {
			// TODO(all): We could do hook concurrent again (assuming different hooks don't interfere with each other)
			hook(result, err)
		}
	}
}
