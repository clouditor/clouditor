package discovery_test

import (
	"context"
	"testing"

	"clouditor.io/clouditor/api/discovery"
	service_discovery "clouditor.io/clouditor/service/discovery"
	"clouditor.io/clouditor/voc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
)

var service *service_discovery.Service

type mockDiscoverer struct {
}

func (m mockDiscoverer) Name() string { return "just mocking" }

func (m mockDiscoverer) List() ([]voc.IsResource, error) {
	return []voc.IsResource{
		&voc.ObjectStorageResource{
			StorageResource: voc.StorageResource{
				Resource: voc.Resource{
					ID:   "some-id",
					Name: "some-name",
				},
			},
		},
	}, nil
}

func TestQuery(t *testing.T) {
	var (
		discoverer discovery.Discoverer
		response   *discovery.QueryResponse
		err        error
	)

	service = service_discovery.NewService()

	// use our mock discoverer
	discoverer = mockDiscoverer{}

	// discover some resources
	service.StartDiscovery(discoverer)

	// query them
	response, err = service.Query(context.Background(), &emptypb.Empty{})

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Result.Values)

	m := response.Result.Values[0].GetStructValue().AsMap()

	assert.NotNil(t, m)
	assert.Equal(t, "some-id", m["id"])
	assert.Equal(t, "some-name", m["name"])
}
