package orchestrator_test

import (
	"context"
	"testing"

	"clouditor.io/clouditor/api/assessment"
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
