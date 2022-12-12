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
	"time"

	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/voc"
	"github.com/stretchr/testify/assert"
)

func Test_regoEval_Eval(t *testing.T) {
	type fields struct {
		// TODO(oxisto): move to args
		resource voc.IsCloudResource
		// TODO(oxisto): move to args
		evidenceID string

		qc      *queryCache
		mrtc    *metricsCache
		storage persistence.Storage
		pkg     string
	}
	type args struct {
		src MetricsSource
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		applicable bool
		compliant  map[string]bool
		wantErr    bool
	}{
		/*{
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
				mrtc:       &metricsCache{m: make(map[string][]string)},
				storage:    testutil.NewInMemoryStorage(t),
				pkg:        DefaultRegoPackage,
			},
			args:       args{src: &mockMetricsSource{t: t}},
			applicable: true,
			compliant: map[string]bool{
				"AtRestEncryptionAlgorithm":         true,
				"AtRestEncryptionEnabled":           true,
				"CustomerKeyEncryption":             true,
				"ObjectStoragePublicAccessDisabled": true,
				"ResourceInventory":                 true,
			},
			wantErr: false,
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
					PublicAccess: true,
				},
				evidenceID: mockObjStorage2EvidenceID,
				qc:         newQueryCache(),
				mrtc:       &metricsCache{m: make(map[string][]string)},
				storage:    testutil.NewInMemoryStorage(t),
				pkg:        DefaultRegoPackage,
			},
			args:       args{src: &mockMetricsSource{t: t}},
			applicable: true,
			compliant: map[string]bool{
				"AtRestEncryptionAlgorithm":         false,
				"AtRestEncryptionEnabled":           false,
				"CustomerKeyEncryption":             false,
				"ObjectStoragePublicAccessDisabled": false,
				"ResourceInventory":                 true,
			},
			wantErr: false,
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
					PublicAccess: true,
				},
				evidenceID: mockObjStorage2EvidenceID,
				qc:         newQueryCache(),
				mrtc:       &metricsCache{m: make(map[string][]string)},
				storage:    testutil.NewInMemoryStorage(t),
				pkg:        DefaultRegoPackage,
			},
			args:       args{src: &mockMetricsSource{t: t}},
			applicable: true,
			compliant: map[string]bool{
				"AtRestEncryptionAlgorithm":         false,
				"AtRestEncryptionEnabled":           false,
				"CustomerKeyEncryption":             false,
				"ObjectStoragePublicAccessDisabled": false,
				"ResourceInventory":                 true,
			},
			wantErr: false,
		},
		{
			name: "VM: Compliant Case",
			fields: fields{
				resource: voc.VirtualMachine{
					AutomaticUpdates: &voc.AutomaticUpdates{
						Enabled:      true,
						Interval:     time.Hour * 24,
						SecurityOnly: true,
					},
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
							RetentionPeriod: 36 * time.Hour * 24,
						},
					},
					OsLogging: &voc.OSLogging{
						Logging: &voc.Logging{
							LoggingService:  []voc.ResourceID{"SomeResourceId2"},
							Enabled:         true,
							RetentionPeriod: 36 * time.Hour * 24,
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
				mrtc:       &metricsCache{m: make(map[string][]string)},
				storage:    testutil.NewInMemoryStorage(t),
				pkg:        DefaultRegoPackage,
			},
			args:       args{src: &mockMetricsSource{t: t}},
			applicable: true,
			compliant: map[string]bool{
				"AutomaticUpdatesEnabled":      true,
				"AutomaticUpdatesInterval":     true,
				"AutomaticUpdatesSecurityOnly": true,
				"BootLoggingEnabled":           true,
				"BootLoggingRetention":         true,
				"MalwareProtectionEnabled":     true,
				"OSLoggingRetention":           true,
				"OSLoggingEnabled":             true,
				"ResourceInventory":            true,
			},
			wantErr: false,
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
							RetentionPeriod: 1 * time.Hour * 24,
						},
					},
					OsLogging: &voc.OSLogging{
						Logging: &voc.Logging{
							LoggingService:  []voc.ResourceID{"SomeResourceId3"},
							Enabled:         false,
							RetentionPeriod: 1 * time.Hour * 24,
						},
					},
				},
				evidenceID: mockVM2EvidenceID,
				qc:         newQueryCache(),
				mrtc:       &metricsCache{m: make(map[string][]string)},
				storage:    testutil.NewInMemoryStorage(t),
				pkg:        DefaultRegoPackage,
			},
			args:       args{src: &mockMetricsSource{t: t}},
			applicable: true,
			compliant: map[string]bool{
				"AutomaticUpdatesEnabled":      false,
				"AutomaticUpdatesInterval":     false,
				"AutomaticUpdatesSecurityOnly": false,
				"BootLoggingEnabled":           false,
				"BootLoggingRetention":         false,
				"MalwareProtectionEnabled":     false,
				"OSLoggingEnabled":             false,
				"OSLoggingRetention":           false,
				"ResourceInventory":            true,
			},
			wantErr: false,
		},*/
		{
			name: "Updated metric configuration",
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
							RetentionPeriod: 1 * time.Hour * 24,
						},
					},
					OsLogging: &voc.OSLogging{
						Logging: &voc.Logging{
							LoggingService:  []voc.ResourceID{"SomeResourceId3"},
							Enabled:         false,
							RetentionPeriod: 1 * time.Hour * 24,
						},
					},
				},
				evidenceID: mockVM2EvidenceID,
				qc:         newQueryCache(),
				mrtc:       &metricsCache{m: make(map[string][]string)},
				storage:    testutil.NewInMemoryStorage(t),
				pkg:        DefaultRegoPackage,
			},
			args:       args{src: &updatedMockMetricsSource{mockMetricsSource{t: t}}},
			applicable: true,
			compliant: map[string]bool{
				"AutomaticUpdatesEnabled":      false,
				"AutomaticUpdatesInterval":     false,
				"AutomaticUpdatesSecurityOnly": false,
				"BootLoggingEnabled":           false,
				"BootLoggingRetention":         false,
				"MalwareProtectionEnabled":     false,
				"OSLoggingEnabled":             false,
				"OSLoggingRetention":           false,
				"ResourceInventory":            true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		resource, err := voc.ToStruct(tt.fields.resource)
		assert.NoError(t, err)
		t.Run(tt.name, func(t *testing.T) {
			pe := regoEval{
				qc:   tt.fields.qc,
				mrtc: tt.fields.mrtc,
				pkg:  tt.fields.pkg,
			}
			results, err := pe.Eval(&evidence.Evidence{
				Id:       tt.fields.evidenceID,
				Resource: resource,
			}, tt.args.src)

			if (err != nil) != tt.wantErr {
				t.Errorf("RunEvidence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.NotEmpty(t, results)

			var compliants = map[string]bool{}

			for _, result := range results {
				if result.Applicable {
					compliants[result.MetricID] = result.Compliant
				}
			}

			assert.Equal(t, compliants, tt.compliant)
		})
	}
}
