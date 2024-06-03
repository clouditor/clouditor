package discoverytest

import (
	"fmt"
	"math/rand/v2"
	"strconv"

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/util"
)

// TestDiscoverer implements Discoverer and mocks the API to cloud resources
type TestDiscoverer struct {
	// testCase allows for different implementations for table tests in TestStartDiscovery
	TestCase  int
	ServiceId string
}

func (TestDiscoverer) Name() string { return "just mocking" }

func (m *TestDiscoverer) List() ([]ontology.IsResource, error) {
	// random number is used to get different resource IDs if more than one discoverer is used in the tests
	// the number should be a 2 digit number, so it is easier to cut it off if needed
	rand := strconv.Itoa(rand.IntN(99-10) + 10)
	switch m.TestCase {
	case 0:
		return nil, fmt.Errorf("mock error in List()")
	case 2:
		return []ontology.IsResource{
			&ontology.ObjectStorage{
				Id:       "some-id-" + rand,
				Name:     "some-name",
				ParentId: util.Ref("some-storage-account-id"),
				Raw:      "{}",
			},
			&ontology.ObjectStorageService{
				Id:         "some-storage-account-id-" + rand,
				Name:       "some-storage-account-name",
				StorageIds: []string{"some-id"},
				Raw:        "{}",
				HttpEndpoint: &ontology.HttpEndpoint{
					TransportEncryption: &ontology.TransportEncryption{
						Enforced:        false,
						Enabled:         true,
						ProtocolVersion: 1.2,
					},
				},
			},
		}, nil
	default:
		return nil, nil
	}
}

func (TestDiscoverer) CloudServiceID() string {
	return testdata.MockCloudServiceID1
}
