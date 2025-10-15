// Copyright 2025 Fraunhofer AISEC
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

package ionos

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

type mockSender struct {
}

// mockErrorSender is used to simulate errors in the RoundTrip method.
type mockErrorSender struct {
}

func newMockSender() *mockSender {
	m := &mockSender{}
	return m
}

func newMockErrorSender() *mockErrorSender {
	m := &mockErrorSender{}
	return m
}

// RoundTrip implements http.RoundTripper.
func (mockSender) RoundTrip(req *http.Request) (res *http.Response, err error) {
	// Check if the URL contains an empty segment
	// (e.g., "http://example.com//path") and return a 404 response if it does.
	if hasEmptySegmentinURL(req.URL.Path) {
		return createResponse(req, map[string]interface{}{}, 404)
	} else if strings.HasSuffix(req.URL.Path, "/labels") {
		// Mock response for labels endpoint
		return createResponse(req, map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"id": "label-1",
					"properties": map[string]interface{}{
						"key":   "label1",
						"value": "value1",
					},
				},
				{
					"id": "label-2",
					"properties": map[string]interface{}{
						"key":   "label2",
						"value": "value2",
					},
				},
			},
		}, 200)
	}

	switch req.URL.Path {
	case "/datacenters":
		return createResponse(req, map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"id": testdata.MockIonosDatacenterID1,
					"properties": map[string]interface{}{
						"name":        testdata.MockIonosDatacenterName1,
						"description": testdata.MockIonosDatacenterDescription1,
						"location":    testdata.MockIonosDatacenterLocation1,
					},
					"metadata": map[string]interface{}{},
				},
				{
					"id": testdata.MockIonosDatacenterID2,
					"properties": map[string]interface{}{
						"name":        testdata.MockIonosDatacenterName2,
						"description": testdata.MockIonosDatacenterDescription2,
						"location":    testdata.MockIonosDatacenterLocation2,
					},
					"metadata": map[string]interface{}{},
				},
			},
		}, 200)
	case "/datacenters/99d85e98-c3da-11ed-afa1-0242ac120002/servers":
		return createResponse(req, map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"id": testdata.MockIonosVMID1,
					"properties": map[string]interface{}{
						"name": testdata.MockIonosVMName1,
					},
					"metadata": map[string]interface{}{
						"createdDate": testdata.CreationTime,
					},
					"entities": map[string]interface{}{
						"volumes": map[string]interface{}{
							"items": []map[string]interface{}{
								{
									"id": testdata.MockIonosVolumeID1,
									"properties": map[string]interface{}{
										"name": testdata.MockIonosVolumeName1,
									},
								},
							},
						},
						"nics": map[string]interface{}{
							"items": []map[string]interface{}{
								{
									"id": testdata.MockIonosNicID1,
									"properties": map[string]interface{}{
										"name": testdata.MockIonosNicName1,
									},
								},
							},
						},
					},
				},
				{
					"id": testdata.MockIonosVMID2,
					"properties": map[string]interface{}{
						"name": testdata.MockIonosVMName2,
					},
					"metadata": map[string]interface{}{
						"createdDate": testdata.CreationTime,
					},
					"entities": map[string]interface{}{
						"volumes": map[string]interface{}{
							"items": []map[string]interface{}{
								{
									"id": testdata.MockIonosVolumeID2,
									"properties": map[string]interface{}{
										"name": testdata.MockIonosVolumeName2,
									},
								},
							},
						},
						"nics": map[string]interface{}{
							"items": []map[string]interface{}{
								{
									"id": testdata.MockIonosNicID2,
									"properties": map[string]interface{}{
										"name": testdata.MockIonosNicName2,
									},
								},
							},
						},
					},
				},
			},
		}, 200)
	case "/datacenters/99d85e98-c3da-11ed-afa1-0242ac120002/volumes":
		return createResponse(req, map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"id": testdata.MockIonosVolumeID1,
					"properties": map[string]interface{}{
						"name": testdata.MockIonosVolumeName1,
					},
					"metadata": map[string]interface{}{
						"createdDate": testdata.CreationTime,
					},
				},
				{
					"id": testdata.MockIonosVolumeID2,
					"properties": map[string]interface{}{
						"name": testdata.MockIonosVolumeName2,
					},
					"metadata": map[string]interface{}{
						"createdDate": testdata.CreationTime,
					},
				},
			},
		}, 200)
	case "/datacenters/99d85e98-c3da-11ed-afa1-0242ac120002/loadbalancers":
		return createResponse(req, map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"id": testdata.MockIonosLoadBalancerID1,
					"properties": map[string]interface{}{
						"name": testdata.MockIonosLoadBalancerName1,
					},
					"metadata": map[string]interface{}{
						"createdDate": testdata.CreationTime,
					},
				},
				{
					"id": testdata.MockIonosLoadBalancerID2,
					"properties": map[string]interface{}{
						"name": testdata.MockIonosLoadBalancerName2,
					},
					"metadata": map[string]interface{}{
						"createdDate": testdata.CreationTime,
					},
				},
			},
		}, 200)
	case "/datacenters/a9d85e98-c3da-11ed-afa1-0242ac120002/servers":
		return createResponse(req, map[string]interface{}{
			"items": []map[string]interface{}{},
		}, 200)
	case "/datacenters/a9d85e98-c3da-11ed-afa1-0242ac120002/volumes":
		return createResponse(req, map[string]interface{}{
			"items": []map[string]interface{}{},
		}, 200)
	case "/datacenters/a9d85e98-c3da-11ed-afa1-0242ac120002/loadbalancers":
		return createResponse(req, map[string]interface{}{
			"items": []map[string]interface{}{},
		}, 200)
	default:
		res, err = createResponse(req, map[string]interface{}{}, 404)
		log.Errorf("Not handling mock for %s yet", req.URL.Path)

	}
	return
}

func (mockErrorSender) RoundTrip(req *http.Request) (res *http.Response, err error) {
	return createResponse(req, map[string]interface{}{}, 404)
}

func createResponse(req *http.Request, body any, status int) (*http.Response, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, err
	}
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(buf),
		Request:    req,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}, nil
}

func NewMockIonosDiscovery(roundTrip http.RoundTripper) *ionosDiscovery {
	d := &ionosDiscovery{
		authConfig: &ionoscloud.Configuration{
			HTTPClient: &http.Client{
				Transport: roundTrip,
			},
		},
		ctID: config.DefaultTargetOfEvaluationID,
	}

	if _, ok := d.authConfig.HTTPClient.Transport.(*mockErrorSender); ok {
		d.clients.computeClient = computeClient(true)
	} else {
		d.clients.computeClient = computeClient(false)
	}

	return d
}

func computeClient(errClient bool) *ionoscloud.APIClient {
	client := ionoscloud.NewAPIClient(&ionoscloud.Configuration{
		HTTPClient: &http.Client{
			Transport: newMockSender(),
		},
		Servers: ionoscloud.ServerConfigurations{
			{
				URL: "https://mock"},
		},
	})

	if errClient {
		client.GetConfig().HTTPClient.Transport = &mockErrorSender{}
	}

	return client
}

func Test_ionosDiscovery_Name(t *testing.T) {
	d := NewIonosDiscovery()

	assert.Equal(t, "IONOS Cloud", d.Name())
}

func TestNewIonosDiscovery(t *testing.T) {
	type args struct {
		opts []DiscoveryOption
	}
	tests := []struct {
		name string
		args args
		want discovery.Discoverer
	}{
		{
			name: "Happy path",
			args: args{},
			want: &ionosDiscovery{
				ctID: config.DefaultTargetOfEvaluationID,
			},
		},
		{
			name: "Happy path: with target of evaluation id",
			args: args{
				opts: []DiscoveryOption{
					WithTargetOfEvaluationID(testdata.MockTargetOfEvaluationID1),
				},
			},
			want: &ionosDiscovery{
				ctID: testdata.MockTargetOfEvaluationID1,
			},
		},
		{
			name: "Happy path: with authorizer",
			args: args{
				opts: []DiscoveryOption{
					WithAuthorizer(&ionoscloud.Configuration{
						HTTPClient: http.DefaultClient,
					}),
				},
			},
			want: &ionosDiscovery{
				ctID: config.DefaultTargetOfEvaluationID,
				authConfig: &ionoscloud.Configuration{
					HTTPClient: http.DefaultClient,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewIonosDiscovery(tt.args.opts...)
			assert.Equal(t, tt.want, got, assert.CompareAllUnexported())
		})
	}
}

func Test_ionosDiscovery_TargetOfEvaluationID(t *testing.T) {
	type fields struct {
		authConfig *ionoscloud.Configuration
		clients    clients
		ctID       string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy path",
			fields: fields{
				ctID: config.DefaultTargetOfEvaluationID,
			},
			want: config.DefaultTargetOfEvaluationID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &ionosDiscovery{
				authConfig: tt.fields.authConfig,
				clients:    tt.fields.clients,
				ctID:       tt.fields.ctID,
			}
			if got := d.TargetOfEvaluationID(); got != tt.want {
				t.Errorf("ionosDiscovery.TargetOfEvaluationID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAuthorizer(t *testing.T) {
	type envVariable struct {
		hasEnvVariable   bool
		envVariableKey   string
		envVariableValue string
	}
	type fields struct {
		envVariables []envVariable
	}

	tests := []struct {
		name    string
		fields  fields
		want    assert.Want[*ionoscloud.Configuration]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path",
			fields: fields{
				envVariables: []envVariable{
					{
						hasEnvVariable:   true,
						envVariableKey:   "IONOS_USERNAME",
						envVariableValue: "test_username",
					},
					{
						hasEnvVariable:   true,
						envVariableKey:   "IONOS_PASSWORD",
						envVariableValue: "test_password",
					},
					{
						hasEnvVariable:   true,
						envVariableKey:   "IONOS_TOKEN",
						envVariableValue: "test_token",
					},
					{
						hasEnvVariable:   true,
						envVariableKey:   "IONOS_API_URL",
						envVariableValue: "test_api_url",
					},
				},
			},
			want: func(t *testing.T, got *ionoscloud.Configuration) bool {
				assert.Equal(t, "test_username", got.Username)
				assert.Equal(t, "test_password", got.Password)
				assert.Equal(t, "test_token", got.Token)
				return assert.Equal(t, "test_api_url", got.Host)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set env variables
			for _, env := range tt.fields.envVariables {
				if env.hasEnvVariable {
					t.Setenv(env.envVariableKey, env.envVariableValue)
				}
			}

			got, err := NewAuthorizer()

			tt.want(t, got)
			tt.wantErr(t, err)
		})
	}
}

func Test_ionosDiscovery_authorize(t *testing.T) {
	type fields struct {
		authConfig *ionoscloud.Configuration
		clients    clients
		ctID       string
	}
	tests := []struct {
		name   string
		fields fields
		want   assert.Want[*ionosDiscovery]
	}{
		{
			name: "Happy path",
			fields: fields{
				authConfig: &ionoscloud.Configuration{
					HTTPClient: http.DefaultClient,
				},
				clients: clients{},
			},
			want: func(t *testing.T, got *ionosDiscovery) bool {
				return assert.NotEmpty(t, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &ionosDiscovery{
				authConfig: tt.fields.authConfig,
				clients:    tt.fields.clients,
				ctID:       tt.fields.ctID,
			}
			d.authorize()

			tt.want(t, d)
		})
	}
}

func Test_ionosDiscovery_List(t *testing.T) {
	type fields struct {
		ionosDiscovery *ionosDiscovery
	}
	tests := []struct {
		name    string
		fields  fields
		want    assert.Want[[]ontology.IsResource]
		wantErr assert.WantErr
	}{
		{
			name: "error: discover datacenters",
			fields: fields{
				ionosDiscovery: NewMockIonosDiscovery(newMockErrorSender()),
			},
			want: assert.Nil[[]ontology.IsResource],
			wantErr: func(t *testing.T, err error) bool {
				return assert.ErrorContains(t, err, "could not discover datacenters:")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				ionosDiscovery: NewMockIonosDiscovery(newMockSender()),
			},
			want: func(t *testing.T, gotList []ontology.IsResource) bool {
				return assert.Equal(t, 6, len(gotList))
			},
			wantErr: assert.Nil[error],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.ionosDiscovery
			gotList, err := d.List()

			tt.want(t, gotList)
			tt.wantErr(t, err)
		})
	}
}
