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

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate protoc -I ../../proto -I ../../third_party orchestrator.proto --go_out=../.. --go-grpc_out=../.. --go_opt=Mevidence.proto=clouditor.io/clouditor/api/assessment --go-grpc_opt=Mevidence.proto=clouditor.io/clouditor/api/assessment --openapi_out=../../openapi/orchestrator

//go:embed metrics.json
var f embed.FS

var metrics []*assessment.Metric
var metricIndex map[int32]*assessment.Metric
var log *logrus.Entry

var DefaultMetricsFile = "metrics.json"

// Service is an implementation of the Clouditor Orchestrator service
type Service struct {
	orchestrator.UnimplementedOrchestratorServer
}

func init() {
	log = logrus.WithField("component", "orchestrator")

	if err := LoadMetrics(DefaultMetricsFile); err != nil {
		log.Errorf("Could not load embedded metrics. Will continue with empty metric list: %v", err)
	}

	metricIndex = make(map[int32]*assessment.Metric)
	for _, v := range metrics {
		metricIndex[v.Id] = v
	}
}

// Metrics
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

func (_ *Service) ListMetrics(_ context.Context, _ *orchestrator.ListMetricsRequest) (response *orchestrator.ListMetricsResponse, err error) {
	response = &orchestrator.ListMetricsResponse{
		Metrics: metrics,
	}

	return response, nil
}

func (_ *Service) GetMetric(_ context.Context, request *orchestrator.GetMetricsRequest) (response *orchestrator.GetMetricResponse, err error) {
	var ok bool
	var metric *assessment.Metric

	if metric, ok = metricIndex[request.MetricId]; !ok {
		return nil, status.Errorf(codes.NotFound, "Could not find metric with id %d", request.MetricId)
	}

	response = &orchestrator.GetMetricResponse{
		Metric: metric,
	}

	return response, nil
}

//// Tools
//
//// TODO Implement DeregisterAssessmentTool
//func (_ *Service) RegisterAssessmentTool (ctx context.Context, request *orchestrator.RegisterAssessmentToolRequest) (tool *orchestrator.AssessmentTool, err error) {
//	// TBD
//	return tool, err
//}
//
//// TODO Implement UpdateAssessmentTool
//func (_ *Service) UpdateAssessmentTool (ctx context.Context, request *orchestrator.UpdateAssessmentToolRequest) (tool *orchestrator.AssessmentTool, err error) {
//	// TBD
//	return tool, err
//}
//
//// TODO Implement DeregisterAssessmentTool
//func (_ *Service) DeregisterAssessmentTool (ctx context.Context, request *orchestrator.DeregisterAssessmentToolRequest) (nil, err error) {
//	// TBD
//	return nil, err
//}
//
//
//// TODO Implement ListAssessmentTools
//func (_ *Service) ListAssessmentTools (ctx context.Context, request *orchestrator.ListAssessmentToolsRequest) (tools *orchestrator.ListAssessmentToolsResponse, err error) {
//	// TBD
//	return tools, err
//}
//
//// TODO Implement GetAssessmentTool
//func (_ *Service) GetAssessmentTool (ctx context.Context, request *orchestrator.GetAssessmentToolRequest) (tool *orchestrator.AssessmentTool, err error) {
//	// TBD
//	return tool, err
//}
