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

package orchestrator_test

import (
	"clouditor.io/clouditor/api/assessment"
	"context"
	"io/fs"
	"testing"

	"clouditor.io/clouditor/api/orchestrator"
	service_orchestrator "clouditor.io/clouditor/service/orchestrator"
	"github.com/stretchr/testify/assert"
)

var service service_orchestrator.Service

func TestListMetrics(t *testing.T) {
	var (
		response *orchestrator.ListMetricsResponse
		err      error
	)

	response, err = service.ListMetrics(context.TODO(), &orchestrator.ListMetricsRequest{})

	assert.Nil(t, err)
	assert.NotEmpty(t, response.Metrics)
}

func TestGetMetric(t *testing.T) {
	var (
		request *orchestrator.GetMetricsRequest
		metric  *assessment.Metric
		err     error
	)

	request = &orchestrator.GetMetricsRequest{
		MetricId: 1,
	}

	metric, err = service.GetMetric(context.TODO(), request)

	assert.Nil(t, err)
	assert.NotNil(t, metric)
	assert.Equal(t, request.MetricId, metric.Id)
}

func TestLoad(t *testing.T) {
	var err = service_orchestrator.LoadMetrics("notfound.json")

	assert.ErrorIs(t, err, fs.ErrNotExist)

	err = service_orchestrator.LoadMetrics("metrics.json")

	assert.Nil(t, err)
}
