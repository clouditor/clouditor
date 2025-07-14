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
	"context"
	"encoding/json"
	"io"
	"net/http"
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

func newMockSender() *mockSender {
	m := &mockSender{}
	return m
}
func (mockSender) Do(req *http.Request) (res *http.Response, err error) {
	switch req.URL.Path {
	case "/subscriptions":
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":             "/subscriptions/00000000-0000-0000-0000-000000000000",
					"subscriptionId": "00000000-0000-0000-0000-000000000000",
					"name":           "sub1",
					"displayName":    "displayName",
				},
			},
		}, 200)
	default:
		res, err = createResponse(req, map[string]interface{}{}, 404)
		log.Errorf("Not handling mock for %s yet", req.URL.Path)

	}
	return
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

type mockAuthorizer struct{}

func (*mockAuthorizer) GetConfiguration(_ context.Context) (ionoscloud.Configuration, error) {
	var config ionoscloud.Configuration

	return config, nil
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
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path",
			fields: fields{
				authConfig: &ionoscloud.Configuration{
					HTTPClient: http.DefaultClient,
				},
				clients: clients{},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &ionosDiscovery{
				authConfig: tt.fields.authConfig,
				clients:    tt.fields.clients,
				ctID:       tt.fields.ctID,
			}
			err := d.authorize()

			tt.wantErr(t, err)
		})
	}
}

// TODO(anatheka): Add test
func Test_ionosDiscovery_List(t *testing.T) {
	type fields struct {
		authConfig *ionoscloud.Configuration
		clients    clients
		ctID       string
	}
	tests := []struct {
		name     string
		fields   fields
		wantList assert.Want[[]ontology.IsResource]
		wantErr  assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &ionosDiscovery{
				authConfig: tt.fields.authConfig,
				clients:    tt.fields.clients,
				ctID:       tt.fields.ctID,
			}
			gotList, err := d.List()

			tt.wantList(t, gotList)
			tt.wantErr(t, err)
		})
	}
}
