// Copyright 2021-2022 Fraunhofer AISEC
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
	"testing"

	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/voc"
	"github.com/stretchr/testify/assert"
)

func Test_regoEval_Eval(t *testing.T) {
	type fields struct {
		// TODO(oxisto): move to args
		resource voc.IsCloudResource
		// TODO(oxisto): move to args
		evidenceID string

		qc   *queryCache
		mrtc *metricsResourceTypeCache
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
				},
				evidenceID: mockObjStorage1EvidenceID,
				qc:         newQueryCache(),
				mrtc:       &metricsResourceTypeCache{m: make(map[string][]string)},
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
				},
				evidenceID: mockObjStorage2EvidenceID,
				qc:         newQueryCache(),
				mrtc:       &metricsResourceTypeCache{m: make(map[string][]string)},
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
				},
				evidenceID: mockObjStorage2EvidenceID,
				qc:         newQueryCache(),
				mrtc:       &metricsResourceTypeCache{m: make(map[string][]string)},
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
							Type: []string{"Virtual Machine", "Compute", "Resource"},
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
				qc:         newQueryCache(),
				mrtc:       &metricsResourceTypeCache{m: make(map[string][]string)},
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
				qc:         newQueryCache(),
				mrtc:       &metricsResourceTypeCache{m: make(map[string][]string)},
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
			pe := regoEval{
				qc:   tt.fields.qc,
				mrtc: tt.fields.mrtc,
			}
			results, err := pe.Eval(&evidence.Evidence{
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
