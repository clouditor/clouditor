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

package discovery

import (
	"context"
	"fmt"
	"testing"
	"time"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/voc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestStartDiscovery(t *testing.T) {
	discoveryService := NewService()

	type fields struct {
		assessmentStream    assessment.Assessment_AssessEvidencesClient
		evidenceStoreStream evidence.EvidenceStore_StoreEvidencesClient
		discoverer          discovery.Discoverer
	}

	tests := []struct {
		name          string
		fields        fields
		checkEvidence bool
	}{
		{
			name: "Err in discoverer",
			fields: fields{
				discoverer: mockDiscoverer{testCase: 0}},
		},
		{
			name: "Err in marshaling the resource containing circular dependencies",
			fields: fields{
				discoverer: mockDiscoverer{testCase: 1}},
		},
		{
			name: "No err in discoverer but no evidence stream to assessment",
			fields: fields{
				discoverer: mockDiscoverer{testCase: 2},
			},
		},
		{
			name: "No err in discoverer but no evidence stream to evidence store available",
			fields: fields{
				assessmentStream: &mockAssessmentStream{},
				discoverer:       mockDiscoverer{testCase: 2}},
		},
		{
			name: "No err in discoverer but streaming to assessment fails",
			fields: fields{
				assessmentStream:    &mockAssessmentStream{},
				evidenceStoreStream: mockEvidenceStoreStream{},
				discoverer:          mockDiscoverer{testCase: 2}},
		},
		{
			name: "No err",
			fields: fields{
				assessmentStream:    &mockAssessmentStream{connectionEstablished: true},
				evidenceStoreStream: mockEvidenceStoreStream{},
				discoverer:          mockDiscoverer{testCase: 2}},
			checkEvidence: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			discoveryService.AssessmentStream = tt.fields.assessmentStream
			discoveryService.EvidenceStoreStream = tt.fields.evidenceStoreStream
			discoveryService.StartDiscovery(tt.fields.discoverer)

			// APIs for assessment and evidence store both send the same evidence. Thus, testing one is enough.
			if tt.checkEvidence {
				e := discoveryService.AssessmentStream.(*mockAssessmentStream).sentEvidence
				// Check if UUID has been created
				assert.NotEmpty(t, e.Id)
				// Check if cloud resources / properties are there
				assert.NotEmpty(t, e.Resource)
				// Check if ID of mockDiscovery's resource is mapped to resource id of the evidence
				list, _ := tt.fields.discoverer.List()
				assert.Equal(t, string(list[0].GetID()), e.Resource.GetStructValue().AsMap()["id"].(string))
			}
		})
	}

}

func TestQuery(t *testing.T) {
	s := NewService()
	s.StartDiscovery(mockDiscoverer{testCase: 2})

	type fields struct {
		resources    map[string]voc.IsCloudResource
		queryRequest *discovery.QueryRequest
	}
	tests := []struct {
		name string
		fields
		numberOfQueriedResources int
		wantErr                  bool
	}{
		{
			name: "Err when unmarshalling",
			fields: fields{
				queryRequest: &discovery.QueryRequest{},
				resources: map[string]voc.IsCloudResource{
					"MockResourceId": wrongFormattedResource(),
				},
			},
			numberOfQueriedResources: 1,
			wantErr:                  true,
		},
		{
			name:                     "Filter type",
			fields:                   fields{queryRequest: &discovery.QueryRequest{FilteredType: "Compute"}},
			numberOfQueriedResources: 0,
		},
		{
			name:                     "No filtering",
			fields:                   fields{queryRequest: &discovery.QueryRequest{}},
			numberOfQueriedResources: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Add resources (bad formatted ones)
			if res := tt.fields.resources; res != nil {
				for k, v := range res {
					s.resources[k] = v
				}
			}
			response, err := s.Query(context.TODO(), &discovery.QueryRequest{FilteredType: tt.fields.queryRequest.FilteredType})
			assert.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				return
			}
			assert.Equal(t, tt.numberOfQueriedResources, len(response.Results.Values))
		})
		// If a bad resource was added it will be removed. Otherwise no-op
		delete(s.resources, "MockResourceId")
	}

}

// mockDiscoverer implements Discoverer and mocks the API to cloud resources
type mockDiscoverer struct {
	// testCase allows for different implementations for table tests in TestStartDiscovery
	testCase int
}

func (m mockDiscoverer) Name() string { return "just mocking" }

func (m mockDiscoverer) List() ([]voc.IsCloudResource, error) {
	switch m.testCase {
	case 0:
		return nil, fmt.Errorf("mock Error in List()")
	case 1:
		return []voc.IsCloudResource{wrongFormattedResource()}, nil
	case 2:
		return []voc.IsCloudResource{
			&voc.ObjectStorage{
				Storage: &voc.Storage{
					CloudResource: &voc.CloudResource{
						ID:   "some-id",
						Name: "some-name",
						Type: []string{"ObjectStorage", "Storage", "Resource"},
					},
				},
				HttpEndpoint: &voc.HttpEndpoint{
					TransportEncryption: &voc.TransportEncryption{
						Enforced:   false,
						Enabled:    true,
						TlsVersion: "TLS1_2",
					},
				},
			},
		}, nil
	default:
		return nil, nil
	}
}

func wrongFormattedResource() voc.IsCloudResource {
	res1 := mockIsCloudResource{Another: nil}
	res2 := mockIsCloudResource{Another: &res1}
	res1.Another = &res2
	return res1
}

// mockAssessmentStream implements Assessment_AssessEvidencesClient interface
type mockAssessmentStream struct {
	// We add sentEvidence field to test the evidence that would be sent over gRPC
	sentEvidence *evidence.Evidence
	// We add connectionEstablished to differentiate between the case where evidences can be sent and not
	connectionEstablished bool
}

func (m *mockAssessmentStream) Send(req *assessment.AssessEvidenceRequest) (err error) {
	e := req.Evidence
	if m.connectionEstablished {
		m.sentEvidence = e
	} else {
		err = fmt.Errorf("MocK Send error")
	}
	return
}

func (*mockAssessmentStream) CloseAndRecv() (*emptypb.Empty, error) {
	return nil, nil
}

func (*mockAssessmentStream) Header() (metadata.MD, error) {
	return nil, nil
}

func (*mockAssessmentStream) Trailer() metadata.MD {
	return nil
}

func (*mockAssessmentStream) CloseSend() error {
	return nil
}

func (*mockAssessmentStream) Context() context.Context {
	return nil
}

func (*mockAssessmentStream) SendMsg(_ interface{}) error {
	return nil
}

func (*mockAssessmentStream) RecvMsg(_ interface{}) error {
	return nil
}

// mockEvidenceStoreStream implements EvidenceStore_StoreEvidencesClient interface
type mockEvidenceStoreStream struct {
}

func (mockEvidenceStoreStream) Send(_ *evidence.StoreEvidenceRequest) error {
	return fmt.Errorf("MocK Send error")
}

func (mockEvidenceStoreStream) CloseAndRecv() (*emptypb.Empty, error) {
	return nil, nil
}

func (mockEvidenceStoreStream) Header() (metadata.MD, error) {
	return nil, nil
}

func (mockEvidenceStoreStream) Trailer() metadata.MD {
	return nil
}

func (mockEvidenceStoreStream) CloseSend() error {
	return nil
}

func (mockEvidenceStoreStream) Context() context.Context {
	return nil
}

func (mockEvidenceStoreStream) SendMsg(_ interface{}) error {
	return nil
}

func (mockEvidenceStoreStream) RecvMsg(_ interface{}) error {
	return nil
}

// mockIsCloudResource implements mockIsCloudResource interface.
// It is used for json.marshal to fail since it contains circular dependency
type mockIsCloudResource struct {
	Another *mockIsCloudResource `json:"Another"`
}

func (mockIsCloudResource) GetID() voc.ResourceID {
	return "MockResourceId"
}

func (mockIsCloudResource) GetName() string {
	return ""
}

func (mockIsCloudResource) GetType() []string {
	return nil
}

func (mockIsCloudResource) HasType(_ string) bool {
	return false
}

func (mockIsCloudResource) GetCreationTime() *time.Time {
	return nil
}

//// ToDo: Adapt TestQuery at a later stage when we fully implement the standalone version
//func TestQuery(t *testing.T) {
//	var (
//		discoverer discovery.Discoverer
//		response   *discovery.QueryResponse
//		err        error
//	)
//
//	var ready = make(chan bool)
//
//	assessmentServer := standalone.NewAssessmentServer().(*service_assessment.Service)
//	assessmentServer.resultHook = func(result *assessment.Result, err error) {
//		if result.MetricId == 1 {
//			assert.Nil(t, err)
//			assert.NotNil(t, result)
//
//			assert.Equal(t, "some-id", result.ResourceId)
//			assert.Equal(t, true, result.Compliant)
//		}
//
//		ready <- true
//	}
//
//	client := standalone.NewAssessmentClient()
//
//	service = service_discovery.NewService()
//	service.AssessmentStream, _ = client.AssessEvidences(context.Background())
//
//	// use our mock discoverer
//	discoverer = mockDiscoverer{}
//
//	// discover some resources
//	service.StartDiscovery(discoverer)
//
//	// make the test wait for streaming evidence
//	select {
//	case <-ready:
//		break
//	case <-time.After(10 * time.Second):
//		assert.Fail(t, "Timeout while waiting for evidence assessment result to be ready")
//	}
//
//	// query them
//	response, err = service.Query(context.Background(), &discovery.QueryRequest{
//		// this should only result 1 resource, and not the compute resource
//		FilteredType: "ObjectStorage",
//	})
//
//	assert.Nil(t, err)
//	assert.NotNil(t, response)
//	assert.NotEmpty(t, response.Result.Values)
//
//	m := response.Result.Values[0].GetStructValue().AsMap()
//
//	assert.NotNil(t, m)
//	assert.Equal(t, "some-id", m["id"])
//	assert.Equal(t, "some-name", m["name"])
//}
