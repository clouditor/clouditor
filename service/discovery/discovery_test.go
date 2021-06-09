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

package discovery_test

import (
	"context"
	"testing"
	"time"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/discovery"
	service_assessment "clouditor.io/clouditor/service/assessment"
	service_discovery "clouditor.io/clouditor/service/discovery"
	"clouditor.io/clouditor/service/standalone"
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
					Type: []string{"ObjectStorage", "Storage", "Resource"},
				},
			},
			HttpEndpoint: &voc.HttpEndpoint{
				TransportEncryption: voc.NewTransportEncryption(true, false, "TLS1_2"),
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

	service = service_discovery.NewService()
	service.AssessmentStream, _ = client.StreamEvidences(context.Background())

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

	// make the test wait for streaming envidence
	select {
	case <-ready:
		return
	case <-time.After(10 * time.Second):
		assert.Fail(t, "Timeout while waiting for evidence assessment result to be ready")
	}
}
