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
	"clouditor.io/clouditor/api/evidence"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/ioutil"

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
	Results map[string]*orchestrator.AssessmentResult

	// Hooks
	AssessmentResultsHook func(result *assessment.Result, err error)
	EvidenceResultHook    func(result *evidence.Evidence, err error)

	db *gorm.DB
}

func init() {
	log = logrus.WithField("component", "orchestrator")
}

func NewService() *Service {
	s := Service{
		Results:              make(map[string]*orchestrator.AssessmentResult),
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
			log.Errorf("Could not retrieve default configuration for metric %s: %s. Ignoring metric", m.Id, err)
			continue
		}

		var config assessment.MetricConfiguration

		err = json.Unmarshal(b, &config)
		if err != nil {
			log.Errorf("Error in reading default configuration for metric %s: %s. Ignoring metric", m.Id, err)
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
		return nil, status.Errorf(codes.NotFound, "Could not find metric with id %s", request.MetricId)
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

	return nil, status.Errorf(codes.NotFound, "Could not find metric configuration for metric %s in service %s", req.MetricId, req.ServiceId)
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
func (s *Service) StoreAssessmentResult(_ context.Context, req *orchestrator.StoreAssessmentResultRequest) (response *orchestrator.StoreAssessmentResultResponse, err error) {

	response = &orchestrator.StoreAssessmentResultResponse{}

	// TODO(all): TBD
	// TODO(garuppel): Write validate function for assessmentResult
	//_, err = e.Validate()
	//if err != nil {
	//	return nil, status.Errorf(codes.InvalidArgument, "invalid evidence: %v", err)
	//}

	s.Results[req.Result.Id] = req.Result

	return
}

func (s *Service) RegisterAssessmentResultsHook(assessmentResultsHook func(result *assessment.Result, err error)) {
	s.AssessmentResultsHook = assessmentResultsHook
}

func (s *Service) RegisterEvidenceHook(evidenceHook func(result *evidence.Evidence, err error)) {
	s.EvidenceResultHook = evidenceHook
}

//// Tools
//
//// TODO Implement DeregisterAssessmentTool
//func ( *Service) RegisterAssessmentTool (ctx context.Context, request *orchestrator.RegisterAssessmentToolRequest) (tool *orchestrator.AssessmentTool, err error) {
//	// TBD
//	return tool, err
//}
//
//// TODO Implement UpdateAssessmentTool
//func ( *Service) UpdateAssessmentTool (ctx context.Context, request *orchestrator.UpdateAssessmentToolRequest) (tool *orchestrator.AssessmentTool, err error) {
//	// TBD
//	return tool, err
//}
//
//// TODO Implement DeregisterAssessmentTool
//func ( *Service) DeregisterAssessmentTool (ctx context.Context, request *orchestrator.DeregisterAssessmentToolRequest) (nil, err error) {
//	// TBD
//	return nil, err
//}
//
//
//// TODO Implement ListAssessmentTools
//func ( *Service) ListAssessmentTools (ctx context.Context, request *orchestrator.ListAssessmentToolsRequest) (tools *orchestrator.ListAssessmentToolsResponse, err error) {
//	// TBD
//	return tools, err
//}
//
//// TODO Implement GetAssessmentTool
//func ( *Service) GetAssessmentTool (ctx context.Context, request *orchestrator.GetAssessmentToolRequest) (tool *orchestrator.AssessmentTool, err error) {
//	// TBD
//	return tool, err
//}
