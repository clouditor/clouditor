package orchestrator_test

import (
	"context"
	"testing"

	"clouditor.io/clouditor/api/orchestrator"
	service_orchestrator "clouditor.io/clouditor/service/orchestrator"
	"github.com/stretchr/testify/assert"
)

func TestListMetrics(t *testing.T) {
	service := service_orchestrator.Service{}

	response, err := service.ListMetrics(context.TODO(), &orchestrator.ListMetricsRequest{})

	assert.Nil(t, err)
	assert.Greater(t, len(response.Metrics), 0)
}
