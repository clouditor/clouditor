package assessment_test

import (
	"context"
	"os"
	"testing"
	"time"

	"clouditor.io/clouditor/api/assessment"
	service_assessment "clouditor.io/clouditor/service/assessment"
	"clouditor.io/clouditor/service/standalone"
	"clouditor.io/clouditor/voc"
	"github.com/stretchr/testify/assert"
)

func TestListResults(t *testing.T) {
	var (
		response *assessment.ListAssessmentResultsResponse
		err      error
	)

	// make sure, that we are in the clouditor root folder to find the policies
	err = os.Chdir("../../..")

	assert.Nil(t, err)

	var ready chan bool = make(chan bool)

	assessmentServer := standalone.NewAssessmentServer().(*service_assessment.Service)
	assessmentServer.ResultHook = func(result *assessment.Result, err error) {
		assert.Nil(t, err)
		assert.NotNil(t, result)

		assert.Equal(t, "some-id", result.ResourceId)
		assert.Equal(t, true, result.Compliant)

		ready <- true
	}

	client := standalone.NewAssessmentClient()

	resource := &voc.ObjectStorageResource{
		StorageResource: voc.StorageResource{
			Resource: voc.Resource{
				ID:   "some-id",
				Type: []string{"ObjectStorage", "Storage", "Resource"},
			},
		},
		HttpEndpoint: &voc.HttpEndpoint{
			TransportEncryption: voc.NewTransportEncryption(true, true, "TLS1_2"),
		},
	}

	s, err := voc.ToStruct(resource)

	assert.Nil(t, err)

	evidence := &assessment.Evidence{
		ResourceId:        "some-id",
		ApplicableMetrics: []int32{1},
		Resource:          s,
	}

	_, err = client.StoreEvidence(context.Background(), &assessment.StoreEvidenceRequest{
		Evidence: evidence,
	})

	assert.Nil(t, err)

	// make the test wait for envidence to be stored
	select {
	case <-ready:
		break
	case <-time.After(10 * time.Second):
		assert.Fail(t, "Timeout while waiting for evidence assessment result to be ready")
	}

	// query them
	response, err = client.ListAssessmentResults(context.Background(), &assessment.ListAssessmentResultsRequest{})

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Results)

	r := response.Results[0]

	assert.NotNil(t, r)
	assert.Equal(t, "some-id", r.ResourceId)
	assert.Equal(t, int32(1), r.MetricId)
}
