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

package policies

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"clouditor.io/clouditor/internal/testutil/clitest"
	"clouditor.io/clouditor/voc"
	"google.golang.org/protobuf/encoding/protojson"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidence"
	"github.com/stretchr/testify/assert"
)

const (
	mockObjStorage1EvidenceID = "1"
	mockObjStorage1ResourceID = "/mockresources/storages/object1"
	mockObjStorage2EvidenceID = "2"
	mockObjStorage2ResourceID = "/mockresources/storages/object2"
	mockVM1EvidenceID         = "3"
	mockVM1ResourceID         = "/mockresources/compute/vm1"
	mockVM2EvidenceID         = "4"
	mockVM2ResourceID         = "/mockresources/compute/vm2"
)

func TestMain(m *testing.M) {
	clitest.AutoChdir()

	os.Exit(m.Run())
}

func TestRunEvidence(t *testing.T) {
	type fields struct {
		resource   voc.IsCloudResource
		evidenceID string
	}
	type args struct {
		source MetricConfigurationSource
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		applicable bool
		compliant  bool
		wantErr    bool
	}{
		{
			name: "ObjectStorage: Compliant Case",
			fields: fields{
				resource: voc.ObjectStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           mockObjStorage1ResourceID,
							CreationTime: 1621086669,
							Type:         []string{"ObjectStorage", "Storage", "Resource"},
							GeoLocation:  voc.GeoLocation{},
						},
						AtRestEncryption: &voc.CustomerKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
							KeyUrl: "SomeUrl",
						},
					},
					HttpEndpoint: &voc.HttpEndpoint{
						Functionality: nil,
						Authenticity:  nil,
						TransportEncryption: &voc.TransportEncryption{
							Enforced:   true,
							Enabled:    true,
							TlsVersion: "TLS1.3",
							Algorithm:  "TLS",
						},
						Url:     "https://aybazestorage.blob.core.windows.net/",
						Method:  "",
						Handler: "",
						Path:    "",
					},
				},
				evidenceID: mockObjStorage1EvidenceID,
			},
			args:       args{source: &mockMetricConfigurationSource{t: t}},
			applicable: true,
			compliant:  true,
			wantErr:    false,
		}, {
			name: "ObjectStorage: Non-Compliant Case with no Encryption at rest",
			fields: fields{
				resource: voc.ObjectStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           mockObjStorage2ResourceID,
							CreationTime: 1621086669,
							Type:         []string{"ObjectStorage", "Storage", "Resource"},
							GeoLocation:  voc.GeoLocation{},
						},
						AtRestEncryption: &voc.AtRestEncryption{
							Algorithm: "NoGoodAlg",
							Enabled:   false,
						},
					},
					HttpEndpoint: &voc.HttpEndpoint{
						Functionality: nil,
						Authenticity:  nil,
						TransportEncryption: &voc.TransportEncryption{
							Enforced:   false,
							Enabled:    false,
							TlsVersion: "1.0",
							Algorithm:  "NoTLS",
						},
						Url:     "https://aybazestorage.blob.core.windows.net/",
						Method:  "",
						Handler: "",
						Path:    "",
					},
				},
				evidenceID: mockObjStorage2EvidenceID,
			},
			args:       args{source: &mockMetricConfigurationSource{t: t}},
			applicable: true,
			compliant:  false,
			wantErr:    false,
		},
		{
			name: "ObjectStorage: Non-Compliant Case 2 with no customer managed key",
			fields: fields{
				resource: voc.ObjectStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           mockObjStorage2ResourceID,
							CreationTime: 1621086669,
							Type:         []string{"ObjectStorage", "Storage", "Resource"},
							GeoLocation:  voc.GeoLocation{},
						},
						AtRestEncryption: &voc.ManagedKeyEncryption{
							AtRestEncryption: &voc.AtRestEncryption{
								// Normally given but for test case purpose only check that no key URL is given
								Algorithm: "",
								Enabled:   false,
							},
						},
					},
					HttpEndpoint: &voc.HttpEndpoint{
						Functionality: nil,
						Authenticity:  nil,
						TransportEncryption: &voc.TransportEncryption{
							Enforced:   false,
							Enabled:    false,
							TlsVersion: "1.0",
							Algorithm:  "NoTLS",
						},
						Url:     "https://aybazestorage.blob.core.windows.net/",
						Method:  "",
						Handler: "",
						Path:    "",
					},
				},
				evidenceID: mockObjStorage2EvidenceID,
			},
			args:       args{source: &mockMetricConfigurationSource{t: t}},
			applicable: true,
			compliant:  false,
			wantErr:    false,
		},
		{
			name: "VM: Compliant Case",
			fields: fields{
				resource: voc.VirtualMachine{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:   mockVM1ResourceID,
							Type: []string{"Compute", "Virtual Machine", "Resource"},
						}},
					BlockStorage: nil,
					BootLogging: &voc.BootLogging{
						Logging: &voc.Logging{
							LoggingService:  []voc.ResourceID{"SomeResourceId1", "SomeResourceId2"},
							Enabled:         true,
							RetentionPeriod: 36,
						},
					},
					OSLogging: &voc.OSLogging{
						Logging: &voc.Logging{
							LoggingService:  []voc.ResourceID{"SomeResourceId2"},
							Enabled:         true,
							RetentionPeriod: 36,
						},
					},
				},
				evidenceID: mockVM1EvidenceID,
			},
			args:       args{source: &mockMetricConfigurationSource{t: t}},
			applicable: true,
			compliant:  true,
			wantErr:    false,
		},
		{
			name: "VM: Non-Compliant Case",
			fields: fields{
				resource: voc.VirtualMachine{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:   mockVM2ResourceID,
							Type: []string{"Compute", "Virtual Machine", "Resource"},
						}},
					BlockStorage: nil,
					BootLogging: &voc.BootLogging{
						Logging: &voc.Logging{
							LoggingService:  []voc.ResourceID{},
							Enabled:         false,
							RetentionPeriod: 1,
						},
					},
					OSLogging: &voc.OSLogging{
						Logging: &voc.Logging{
							LoggingService:  []voc.ResourceID{"SomeResourceId3"},
							Enabled:         false,
							RetentionPeriod: 1,
						},
					},
				},
				evidenceID: mockVM2EvidenceID,
			},
			args:       args{source: &mockMetricConfigurationSource{t: t}},
			applicable: true,
			compliant:  false,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		resource, err := voc.ToStruct(tt.fields.resource)
		assert.NoError(t, err)
		t.Run(tt.name, func(t *testing.T) {
			results, err := RunEvidence(&evidence.Evidence{
				Id:       tt.fields.evidenceID,
				Resource: resource,
			}, tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunEvidence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotEmpty(t, results)
			for _, result := range results {
				assert.Equal(t, tt.applicable, result.Applicable)
				assert.Equal(t, tt.compliant, result.Compliant)
			}
		})
	}
}

type mockMetricConfigurationSource struct {
	t *testing.T
}

func (m *mockMetricConfigurationSource) MetricConfiguration(metric string) (*assessment.MetricConfiguration, error) {
	// Fetch the metric configuration directly from our file

	bundle := fmt.Sprintf("policies/bundles/%s/data.json", metric)
	file, err := os.OpenFile(bundle, os.O_RDONLY, 0600)
	assert.NoError(m.t, err)

	b, err := ioutil.ReadAll(file)
	assert.NoError(m.t, err)

	var config assessment.MetricConfiguration
	err = protojson.Unmarshal(b, &config)
	assert.NoError(m.t, err)

	return &config, nil
}
