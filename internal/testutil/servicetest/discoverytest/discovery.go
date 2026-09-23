package discoverytest

import (
	"fmt"
	"math/rand/v2"
	"strconv"

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testdata"
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
				Id:       new("some-id-" + rand),
				Name:     new("some-name"),
				ParentId: new("some-storage-account-id"),
				Raw:      new("{}"),
			},
			&ontology.ObjectStorageService{
				Id:         new("some-storage-account-id-" + rand),
				Name:       new("some-storage-account-name"),
				StorageIds: []string{"some-id"},
				Raw:        new("{}"),
				HttpEndpoint: &ontology.HttpEndpoint{
					TransportEncryption: &ontology.TransportEncryption{
						Enforced:        new(false),
						Enabled:         new(true),
						ProtocolVersion: new(float32(1.2)),
					},
				},
			},
		}, nil
	default:
		return nil, nil
	}
}

func (TestDiscoverer) TargetOfEvaluationID() string {
	return testdata.MockTargetOfEvaluationID1
}
