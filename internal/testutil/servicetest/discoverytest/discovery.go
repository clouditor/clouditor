package discoverytest

import (
	"fmt"

	"clouditor.io/clouditor/api/ontology"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/util"
)

// TestDiscoverer implements Discoverer and mocks the API to cloud resources
type TestDiscoverer struct {
	// testCase allows for different implementations for table tests in TestStartDiscovery
	TestCase  int
	ServiceId string
}

func (TestDiscoverer) Name() string { return "just mocking" }

func (m *TestDiscoverer) List() ([]ontology.IsResource, error) {
	switch m.TestCase {
	case 0:
		return nil, fmt.Errorf("mock error in List()")
	case 2:
		return []ontology.IsResource{
			&ontology.ObjectStorage{
				Id:       "some-id",
				Name:     "some-name",
				ParentId: util.Ref("some-storage-account-id"),
				Raw:      "{}",
			},
			&ontology.ObjectStorageService{
				Id:         "some-storage-account-id",
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
