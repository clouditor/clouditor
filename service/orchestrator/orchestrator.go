/*
 * Copyright 2021 Fraunhofer AISEC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *           $$\                           $$\ $$\   $$\
 *           $$ |                          $$ |\__|  $$ |
 *  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
 * $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
 * $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
 * $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
 * \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
 *  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
 *
 * This file is part of Clouditor Community Edition.
 */

package discovery

import (
	"context"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate protoc -I ../../proto -I ../../third_party orchestrator.proto --go_out=../.. --go-grpc_out=../.. --go_opt=Mevidence.proto=clouditor.io/clouditor/api/assessment  --openapi_out=../../openapi/orchestrator

var metrics []*assessment.Metric
var metricIndex map[int32]*assessment.Metric

// Service is an implementation of the Clouditor Orchestrator service
type Service struct {
	orchestrator.UnimplementedOrchestratorServer
}

func init() {
	metrics = []*assessment.Metric{
		{
			Id:          1,
			Name:        "Transport Encryption",
			Description: "This metric describes, whether transport encryption is turned on or not",
			Category:    "",
			Scale:       assessment.Metric_ORDINAL,
			Range: &assessment.Range{
				Range: &assessment.Range_Order{
					Order: &assessment.Order{Values: []string{"false", "true"}},
				},
			},
			TargetValue: &assessment.Metric_TargetValueString{TargetValueString: "true"},
		},
	}

	for _, v := range metrics {
		metricIndex[v.Id] = v
	}
}

func (s *Service) ListMetrics(ctx context.Context, request *orchestrator.ListMetricsRequest) (response *orchestrator.ListMetricsResponse, err error) {
	response = &orchestrator.ListMetricsResponse{
		Metrics: metrics,
	}

	return response, nil
}

func (s *Service) GetMetric(ctx context.Context, request *orchestrator.GetMetricsRequest) (metric *assessment.Metric, err error) {
	var ok bool

	if metric, ok = metricIndex[request.MetricId]; !ok {
		return nil, status.Errorf(codes.NotFound, "Could not find metric with id %d", request.MetricId)
	}

	return metric, nil
}
