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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/internal/testutil/clitest"
	"clouditor.io/clouditor/voc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
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
	mockBlockStorage1ID       = "/mockresources/storage/storage1"
	mockBlockStorage2ID       = "/mockresources/storage/storage2"
)

func TestMain(m *testing.M) {
	clitest.AutoChdir()

	os.Exit(m.Run())
}

// TODO(oxisto): Move these to Rego unit tests
func TestRunEvidence(t *testing.T) {
	type fields struct {
		resource   voc.IsCloudResource
		evidenceID string
	}
	type args struct {
		source  MetricConfigurationSource
		related map[string]*structpb.Value
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		applicable bool
		compliant  bool
		wantErr    bool
		want       []*Result
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
				},
				evidenceID: mockObjStorage1EvidenceID,
			},
			args:       args{source: &mockMetricConfigurationSource{t: t}},
			applicable: true,
			compliant:  true,
			wantErr:    false,
			want: []*Result{
				{
					Applicable:  true,
					Compliant:   true,
					TargetValue: "AES256",
					Operator:    "==",
					MetricId:    "AtRestEncryptionAlgorithm",
				},
				{
					Applicable:  true,
					Compliant:   true,
					TargetValue: true,
					Operator:    "==",
					MetricId:    "AtRestEncryptionEnabled",
				},
				{
					Applicable:  true,
					Compliant:   true,
					TargetValue: true,
					Operator:    "==",
					MetricId:    "CustomerKeyEncryption",
				},
			},
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
				},
				evidenceID: mockObjStorage2EvidenceID,
			},
			args:       args{source: &mockMetricConfigurationSource{t: t}},
			applicable: true,
			compliant:  false,
			wantErr:    false,
			want: []*Result{
				{
					Applicable:  true,
					Compliant:   false,
					TargetValue: "AES256",
					Operator:    "==",
					MetricId:    "AtRestEncryptionAlgorithm",
				},
				{
					Applicable:  true,
					Compliant:   false,
					TargetValue: true,
					Operator:    "==",
					MetricId:    "AtRestEncryptionEnabled",
				},
				{
					Applicable:  true,
					Compliant:   false,
					TargetValue: true,
					Operator:    "==",
					MetricId:    "CustomerKeyEncryption",
				},
			},
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
				},
				evidenceID: mockObjStorage2EvidenceID,
			},
			args:       args{source: &mockMetricConfigurationSource{t: t}},
			applicable: true,
			compliant:  false,
			wantErr:    false,
			want: []*Result{
				{
					Applicable:  true,
					Compliant:   false,
					TargetValue: "AES256",
					Operator:    "==",
					MetricId:    "AtRestEncryptionAlgorithm",
				},
				{
					Applicable:  true,
					Compliant:   false,
					TargetValue: true,
					Operator:    "==",
					MetricId:    "AtRestEncryptionEnabled",
				},
				{
					Applicable:  true,
					Compliant:   false,
					TargetValue: true,
					Operator:    "==",
					MetricId:    "CustomerKeyEncryption",
				},
			},
		},
		{
			name: "VM: Compliant Case",
			fields: fields{
				resource: voc.VirtualMachine{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:   mockVM1ResourceID,
							Type: []string{"VirtualMachine", "Compute", "Resource"},
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
					MalwareProtection: &voc.MalwareProtection{
						Enabled:              true,
						DaysSinceActive:      5,
						NumberOfThreatsFound: 5,
						ApplicationLogging: &voc.ApplicationLogging{
							Logging: &voc.Logging{
								Enabled:         true,
								RetentionPeriod: 36,
								LoggingService:  []voc.ResourceID{"SomeAnalyticsService?"},
							},
						},
					},
				},
				evidenceID: mockVM1EvidenceID,
			},
			args:       args{source: &mockMetricConfigurationSource{t: t}},
			applicable: true,
			compliant:  true,
			wantErr:    false,
			want: []*Result{
				{
					Applicable:  true,
					Compliant:   true,
					TargetValue: true,
					Operator:    "==",
					MetricId:    "BootLoggingEnabled",
				},
				{
					Applicable:  true,
					Compliant:   true,
					TargetValue: []interface{}{"SomeResourceId1", "SomeResourceId2"},
					Operator:    "==",
					MetricId:    "BootLoggingSecureTransport",
				},
				{
					Applicable:  true,
					Compliant:   true,
					TargetValue: json.Number("35"),
					Operator:    ">=",
					MetricId:    "BootLoggingRetention",
				},
				{
					Applicable:  true,
					Compliant:   true,
					TargetValue: true,
					Operator:    "==",
					MetricId:    "MalwareProtectionEnabled",
				},
				{
					Applicable:  true,
					Compliant:   true,
					TargetValue: []interface{}{"SomeAnalyticsService?", "?"},
					Operator:    "==",
					MetricId:    "MalwareProtectionSecureTransport",
				},
				{
					Applicable:  true,
					Compliant:   true,
					TargetValue: true,
					Operator:    "==",
					MetricId:    "OSLoggingEnabled",
				},
				{
					Applicable:  true,
					Compliant:   true,
					TargetValue: []interface{}{"SomeResourceId1", "SomeResourceId2"},
					Operator:    "==",
					MetricId:    "OSLoggingSecureTransport",
				},
				{
					Applicable:  true,
					Compliant:   true,
					TargetValue: json.Number("35"),
					Operator:    ">=",
					MetricId:    "OSLoggingRetention",
				},
			},
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
			want: []*Result{
				{
					Applicable:  true,
					Compliant:   false,
					TargetValue: true,
					Operator:    "==",
					MetricId:    "BootLoggingEnabled",
				},
				{
					Applicable:  true,
					Compliant:   false,
					TargetValue: []interface{}{"SomeResourceId1", "SomeResourceId2"},
					Operator:    "==",
					MetricId:    "BootLoggingSecureTransport",
				},
				{
					Applicable:  true,
					Compliant:   false,
					TargetValue: json.Number("35"),
					Operator:    ">=",
					MetricId:    "BootLoggingRetention",
				},
				{
					Applicable:  true,
					Compliant:   false,
					TargetValue: true,
					Operator:    "==",
					MetricId:    "MalwareProtectionEnabled",
				},
				{
					Applicable:  true,
					Compliant:   false,
					TargetValue: true,
					Operator:    "==",
					MetricId:    "OSLoggingEnabled",
				},
				{
					Applicable:  true,
					Compliant:   false,
					TargetValue: []interface{}{"SomeResourceId1", "SomeResourceId2"},
					Operator:    "==",
					MetricId:    "OSLoggingSecureTransport",
				},
				{
					Applicable:  true,
					Compliant:   false,
					TargetValue: json.Number("35"),
					Operator:    ">=",
					MetricId:    "OSLoggingRetention",
				},
			},
		},
		{
			name: "VM: Related Evidence",
			fields: fields{
				resource: voc.VirtualMachine{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:   mockVM1ResourceID,
							Type: []string{"VirtualMachine", "Compute", "Resource"},
						}},
					BlockStorage: []voc.ResourceID{mockBlockStorage1ID},
				},
				evidenceID: mockVM1EvidenceID,
			},
			args: args{
				source: &mockMetricConfigurationSource{t: t},
				related: map[string]*structpb.Value{
					mockBlockStorage1ID: testutil.ToStruct(&voc.BlockStorage{
						Storage: &voc.Storage{
							Resource: &voc.Resource{
								ID:   mockBlockStorage1ID,
								Type: []string{"BlockStorage", "Storage", "Resource"},
							},
							AtRestEncryption: voc.AtRestEncryption{
								Enabled:   true,
								Algorithm: "AES256",
							},
						},
					}, t),
				},
			},
			wantErr: false,
			want: []*Result{
				{
					Applicable:  true,
					Compliant:   false,
					TargetValue: true,
					Operator:    "==",
					MetricId:    "BootLoggingEnabled",
				},
				{
					Applicable:  true,
					Compliant:   false,
					TargetValue: []interface{}{"SomeResourceId1", "SomeResourceId2"},
					Operator:    "==",
					MetricId:    "BootLoggingSecureTransport",
				},
				{
					Applicable:  true,
					Compliant:   false,
					TargetValue: json.Number("35"),
					Operator:    ">=",
					MetricId:    "BootLoggingRetention",
				},
				{
					Applicable:  true,
					Compliant:   false,
					TargetValue: true,
					Operator:    "==",
					MetricId:    "MalwareProtectionEnabled",
				},
				{
					Applicable:  true,
					Compliant:   false,
					TargetValue: true,
					Operator:    "==",
					MetricId:    "OSLoggingEnabled",
				},
				{
					Applicable:  true,
					Compliant:   false,
					TargetValue: []interface{}{"SomeResourceId1", "SomeResourceId2"},
					Operator:    "==",
					MetricId:    "OSLoggingSecureTransport",
				},
				{
					Applicable:  true,
					Compliant:   false,
					TargetValue: json.Number("35"),
					Operator:    ">=",
					MetricId:    "OSLoggingRetention",
				},
				{
					Applicable:  true,
					Compliant:   true,
					TargetValue: true,
					Operator:    "==",
					MetricId:    "VirtualMachineDiskEncryptionEnabled",
				},
			},
		},
	}

	for _, tt := range tests {
		resource, err := voc.ToStruct(tt.fields.resource)
		assert.NoError(t, err)
		t.Run(tt.name, func(t *testing.T) {
			applicableMetrics.Clear()

			results, err := RunEvidence(&evidence.Evidence{
				Id:       tt.fields.evidenceID,
				Resource: resource,
			}, tt.args.source, tt.args.related)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunEvidence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, results)
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
