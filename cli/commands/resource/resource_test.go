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
//go:build exclude

package resource

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/cli"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/testutil/clitest"
	"clouditor.io/clouditor/server"
	service_discovery "clouditor.io/clouditor/service/discovery"
	"clouditor.io/clouditor/voc"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestMain(m *testing.M) {
	svc := service_discovery.NewService()
	svc.StartDiscovery(mockDiscoverer{testCase: 2})

	os.Exit(clitest.RunCLITest(m, server.WithDiscovery(svc)))
}

func TestAddCommands(t *testing.T) {
	cmd := NewResourceCommand()

	// Check if sub commands were added
	assert.True(t, cmd.HasSubCommands())

	// Check if NewListCommand was added
	for _, v := range cmd.Commands() {
		if v.Use == "list" {
			return
		}
	}
	t.Errorf("No list command was added")
}

func TestNewListCommand(t *testing.T) {
	var err error
	var b bytes.Buffer

	cli.Output = &b

	cmd := NewListResourcesCommand()
	err = cmd.RunE(nil, []string{})
	assert.NoError(t, err)

	var response = &discovery.ListResourcesResponse{}
	err = protojson.Unmarshal(b.Bytes(), response)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Results)
}

// Mocking code below is copied from clouditor.io/service/discovery

// mockDiscoverer implements Discoverer and mocks the API to cloud resources
type mockDiscoverer struct {
	// testCase allows for different implementations for table tests in TestStartDiscovery
	testCase int
}

func (mockDiscoverer) Name() string { return "just mocking" }

func (m mockDiscoverer) List() ([]voc.IsCloudResource, error) {
	switch m.testCase {
	case 0:
		return nil, fmt.Errorf("mock error in List()")
	case 1:
		return []voc.IsCloudResource{wrongFormattedResource()}, nil
	case 2:
		return []voc.IsCloudResource{
			&voc.ObjectStorage{
				Storage: &voc.Storage{
					Resource: &voc.Resource{
						ID:   testdata.MockResourceID1,
						Name: testdata.MockResourceName1,
						Type: []string{"ObjectStorage", "Storage", "Resource"},
					},
				},
			},
			&voc.ObjectStorageService{
				StorageService: &voc.StorageService{
					Storage: []voc.ResourceID{"some-id"},
					NetworkService: &voc.NetworkService{
						Networking: &voc.Networking{
							Resource: &voc.Resource{
								ID:   testdata.MockResourceStorageID,
								Name: testdata.MockResourceStorageName,
								Type: []string{"StorageService", "NetworkService", "Networking", "Resource"},
							},
						},
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

func (mockDiscoverer) CloudServiceID() string {
	return testdata.MockCloudServiceID1
}

func wrongFormattedResource() voc.IsCloudResource {
	res1 := mockIsCloudResource{Another: nil}
	res2 := mockIsCloudResource{Another: &res1}
	res1.Another = &res2
	return res1
}

// mockIsCloudResource implements mockIsCloudResource interface.
// It is used for json.marshal to fail since it contains circular dependency
type mockIsCloudResource struct {
	Another *mockIsCloudResource `json:"Another"`
}

func (mockIsCloudResource) GetID() voc.ResourceID {
	return "MockResourceID"
}

func (mockIsCloudResource) GetServiceID() string {
	return "MockServiceID"
}

func (mockIsCloudResource) SetServiceID(_ string) {

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

func (mockIsCloudResource) GetRaw() string {
	return ""
}

func (mockIsCloudResource) SetRaw(_ string) {
}

func (mockIsCloudResource) Related() []string {
	return []string{}
}
