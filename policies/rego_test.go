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

	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/prototest"
	"clouditor.io/clouditor/v2/persistence"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_regoEval_Eval(t *testing.T) {
	type fields struct {
		qc      *queryCache
		mrtc    *metricsCache
		storage persistence.Storage
		pkg     string
	}
	type args struct {
		resource   ontology.IsResource
		related    map[string]ontology.IsResource
		evidenceID string
		src        MetricsSource
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		applicable bool
		compliant  map[string]bool
		wantErr    assert.WantErr
	}{
		{
			name: "ObjectStorage: Compliant Case",
			fields: fields{
				qc:      newQueryCache(),
				mrtc:    &metricsCache{m: make(map[string][]*assessment.Metric)},
				storage: testutil.NewInMemoryStorage(t),
				pkg:     DefaultRegoPackage,
			},
			args: args{
				resource: &ontology.ObjectStorage{
					Id:           mockObjStorage1ResourceID,
					CreationTime: timestamppb.New(time.Unix(1621086669, 0)),
					AtRestEncryption: &ontology.AtRestEncryption{
						Type: &ontology.AtRestEncryption_CustomerKeyEncryption{
							CustomerKeyEncryption: &ontology.CustomerKeyEncryption{
								Algorithm: "AES256",
								Enabled:   true,
								KeyUrl:    "SomeUrl",
							},
						},
					},
				},
				evidenceID: mockObjStorage1EvidenceID,
				src:        &mockMetricsSource{t: t},
			},
			applicable: true,
			compliant: map[string]bool{
				"AtRestEncryptionAlgorithm":         true,
				"AtRestEncryptionEnabled":           true,
				"CustomerKeyEncryption":             true,
				"ObjectStoragePublicAccessDisabled": true,
				"ResourceInventory":                 true,
			},
			wantErr: assert.Nil[error],
		},
		{
			name: "ObjectStorage: Non-Compliant Case with no Encryption at rest",
			fields: fields{
				qc:      newQueryCache(),
				mrtc:    &metricsCache{m: make(map[string][]*assessment.Metric)},
				storage: testutil.NewInMemoryStorage(t),
				pkg:     DefaultRegoPackage,
			},
			args: args{
				resource: &ontology.ObjectStorage{
					Id:           mockObjStorage1ResourceID,
					CreationTime: timestamppb.New(time.Unix(1621086669, 0)),
					AtRestEncryption: &ontology.AtRestEncryption{
						Type: &ontology.AtRestEncryption_CustomerKeyEncryption{
							CustomerKeyEncryption: &ontology.CustomerKeyEncryption{
								Algorithm: "NoGoodAlg",
								Enabled:   false,
							},
						},
					},
					PublicAccess: true,
				},
				evidenceID: mockObjStorage2EvidenceID,
				src:        &mockMetricsSource{t: t},
			},
			applicable: true,
			compliant: map[string]bool{
				"AtRestEncryptionAlgorithm":         false,
				"AtRestEncryptionEnabled":           false,
				"CustomerKeyEncryption":             false,
				"ObjectStoragePublicAccessDisabled": false,
				"ResourceInventory":                 true,
			},
			wantErr: assert.Nil[error],
		},
		{
			name: "ObjectStorage: Non-Compliant Case 2 with no customer managed key",
			fields: fields{
				qc:      newQueryCache(),
				mrtc:    &metricsCache{m: make(map[string][]*assessment.Metric)},
				storage: testutil.NewInMemoryStorage(t),
				pkg:     DefaultRegoPackage,
			},
			args: args{
				resource: &ontology.ObjectStorage{
					Id:           mockObjStorage1ResourceID,
					CreationTime: timestamppb.New(time.Unix(1621086669, 0)),
					AtRestEncryption: &ontology.AtRestEncryption{
						Type: &ontology.AtRestEncryption_CustomerKeyEncryption{
							CustomerKeyEncryption: &ontology.CustomerKeyEncryption{
								// Normally given but for test case purpose only check that no key URL is given
								Algorithm: "",
								Enabled:   false,
							},
						},
					},
					PublicAccess: true,
				},
				evidenceID: mockObjStorage2EvidenceID,
				src:        &mockMetricsSource{t: t},
			},
			applicable: true,
			compliant: map[string]bool{
				"AtRestEncryptionAlgorithm":         false,
				"AtRestEncryptionEnabled":           false,
				"CustomerKeyEncryption":             false,
				"ObjectStoragePublicAccessDisabled": false,
				"ResourceInventory":                 true,
			},
			wantErr: assert.Nil[error],
		},
		{
			name: "VM: Compliant Case",
			fields: fields{
				qc:      newQueryCache(),
				mrtc:    &metricsCache{m: make(map[string][]*assessment.Metric)},
				storage: testutil.NewInMemoryStorage(t),
				pkg:     DefaultRegoPackage,
			},
			args: args{
				src: &mockMetricsSource{t: t},
				resource: &ontology.VirtualMachine{
					Id: mockVM1ResourceID,
					AutomaticUpdates: &ontology.AutomaticUpdates{
						Enabled:      true,
						Interval:     durationpb.New(time.Hour * 24 * 30),
						SecurityOnly: true,
					},
					BootLogging: &ontology.BootLogging{
						LoggingServiceIds: []string{"SomeResourceId1", "SomeResourceId2"},
						Enabled:           true,
						RetentionPeriod:   durationpb.New(36 * time.Hour * 24),
					},
					OsLogging: &ontology.OSLogging{
						LoggingServiceIds: []string{"SomeResourceId2"},
						Enabled:           true,
						RetentionPeriod:   durationpb.New(36 * time.Hour * 24),
					},
					MalwareProtection: &ontology.MalwareProtection{
						Enabled:              true,
						DaysSinceActive:      durationpb.New(time.Hour * 24 * 5),
						NumberOfThreatsFound: 5,
						ApplicationLogging: &ontology.ApplicationLogging{
							Enabled:           true,
							RetentionPeriod:   durationpb.New(time.Hour * 24 * 36),
							LoggingServiceIds: []string{"SomeAnalyticsService?"},
						},
					},
				},
				evidenceID: mockVM1EvidenceID,
			},
			applicable: true,
			compliant: map[string]bool{
				"AutomaticUpdatesEnabled":      true,
				"AutomaticUpdatesInterval":     true,
				"AutomaticUpdatesSecurityOnly": true,
				"BootLoggingEnabled":           true,
				"BootLoggingOutput":            true,
				"BootLoggingRetention":         true,
				"MalwareProtectionEnabled":     true,
				"MalwareProtectionOutput":      true,
				"OSLoggingRetention":           true,
				"OSLoggingOutput":              true,
				"OSLoggingEnabled":             true,
				"ResourceInventory":            true,
			},
			wantErr: assert.Nil[error],
		},
		{
			name: "VM: Non-Compliant Case",
			fields: fields{
				qc:      newQueryCache(),
				mrtc:    &metricsCache{m: make(map[string][]*assessment.Metric)},
				storage: testutil.NewInMemoryStorage(t),
				pkg:     DefaultRegoPackage,
			},
			args: args{
				resource: &ontology.VirtualMachine{
					Id: mockVM2ResourceID,
					BootLogging: &ontology.BootLogging{
						LoggingServiceIds: nil,
						Enabled:           false,
						RetentionPeriod:   durationpb.New(1 * time.Hour * 24),
					},
					OsLogging: &ontology.OSLogging{
						LoggingServiceIds: []string{"SomeResourceId3"},
						Enabled:           false,
						RetentionPeriod:   durationpb.New(1 * time.Hour * 24),
					},
				},
				evidenceID: mockVM2EvidenceID,
				src:        &mockMetricsSource{t: t},
			},
			applicable: true,
			compliant: map[string]bool{
				"AutomaticUpdatesEnabled":      false,
				"AutomaticUpdatesInterval":     false,
				"AutomaticUpdatesSecurityOnly": false,
				"BootLoggingEnabled":           false,
				"BootLoggingOutput":            false,
				"BootLoggingRetention":         false,
				"MalwareProtectionEnabled":     false,
				"OSLoggingEnabled":             false,
				"OSLoggingOutput":              true,
				"OSLoggingRetention":           false,
				"ResourceInventory":            true,
			},
			wantErr: assert.Nil[error],
		},
		{
			name: "VM: Related Evidence: non-compliant VMDiskEncryptionEnabled",
			fields: fields{
				qc:      newQueryCache(),
				mrtc:    &metricsCache{m: make(map[string][]*assessment.Metric)},
				storage: testutil.NewInMemoryStorage(t),
				pkg:     DefaultRegoPackage,
			},
			args: args{
				resource: &ontology.VirtualMachine{
					Id:              mockVM2ResourceID,
					BlockStorageIds: []string{mockBlockStorage1ID},
				},
				evidenceID: mockVM1EvidenceID,
				src:        &mockMetricsSource{t: t},
				related: map[string]ontology.IsResource{
					mockBlockStorage1ID: &ontology.BlockStorage{
						Id: mockBlockStorage1ID,
						AtRestEncryption: &ontology.AtRestEncryption{
							Type: &ontology.AtRestEncryption_CustomerKeyEncryption{
								CustomerKeyEncryption: &ontology.CustomerKeyEncryption{
									Enabled:   false,
									Algorithm: "AES256",
								},
							},
						},
					},
				},
			},
			compliant: map[string]bool{
				"AutomaticUpdatesEnabled":             false,
				"AutomaticUpdatesInterval":            false,
				"AutomaticUpdatesSecurityOnly":        false,
				"BootLoggingEnabled":                  false,
				"BootLoggingOutput":                   false,
				"BootLoggingRetention":                false,
				"MalwareProtectionEnabled":            false,
				"OSLoggingEnabled":                    false,
				"OSLoggingOutput":                     false,
				"OSLoggingRetention":                  false,
				"ResourceInventory":                   true,
				"VirtualMachineDiskEncryptionEnabled": false,
			},
			wantErr: assert.Nil[error],
		},
		{
			name: "VM: Related Evidence",
			fields: fields{
				qc:      newQueryCache(),
				mrtc:    &metricsCache{m: make(map[string][]*assessment.Metric)},
				storage: testutil.NewInMemoryStorage(t),
				pkg:     DefaultRegoPackage,
			},
			args: args{
				resource: &ontology.VirtualMachine{
					Id:              mockVM2ResourceID,
					BlockStorageIds: []string{mockBlockStorage1ID},
				},
				evidenceID: mockVM1EvidenceID,
				src:        &mockMetricsSource{t: t},
				related: map[string]ontology.IsResource{
					mockBlockStorage1ID: &ontology.BlockStorage{
						Id: mockBlockStorage1ID,
						AtRestEncryption: &ontology.AtRestEncryption{
							Type: &ontology.AtRestEncryption_CustomerKeyEncryption{
								CustomerKeyEncryption: &ontology.CustomerKeyEncryption{
									Enabled:   true,
									Algorithm: "AES256",
								},
							},
						},
					},
				},
			},
			compliant: map[string]bool{
				"AutomaticUpdatesEnabled":             false,
				"AutomaticUpdatesInterval":            false,
				"AutomaticUpdatesSecurityOnly":        false,
				"BootLoggingEnabled":                  false,
				"BootLoggingOutput":                   false,
				"BootLoggingRetention":                false,
				"MalwareProtectionEnabled":            false,
				"OSLoggingEnabled":                    false,
				"OSLoggingOutput":                     false,
				"OSLoggingRetention":                  false,
				"ResourceInventory":                   true,
				"VirtualMachineDiskEncryptionEnabled": true,
			},
			wantErr: assert.Nil[error],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pe := regoEval{
				qc:   tt.fields.qc,
				mrtc: tt.fields.mrtc,
				pkg:  tt.fields.pkg,
			}
			results, err := pe.Eval(&evidence.Evidence{
				Id:       tt.args.evidenceID,
				Resource: prototest.NewAny(t, tt.args.resource),
			}, tt.args.resource, tt.args.related, tt.args.src)

			tt.wantErr(t, err)

			assert.NotEmpty(t, results)

			var compliants = map[string]bool{}

			for _, result := range results {
				if result.Applicable {
					compliants[result.MetricID] = result.Compliant
				}
			}

			assert.Equal(t, tt.compliant, compliants)
		})
	}
}

func Test_regoEval_evalMap(t *testing.T) {
	type fields struct {
		qc   *queryCache
		mrtc *metricsCache
		pkg  string
	}
	type args struct {
		baseDir    string
		serviceID  string
		metric     *assessment.Metric
		categoryID string
		m          map[string]interface{}
		src        MetricsSource
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult assert.Want[*Result]
		wantErr    assert.WantErr
	}{
		{
			name: "default metric configuration",
			fields: fields{
				qc:   newQueryCache(),
				mrtc: &metricsCache{m: make(map[string][]*assessment.Metric)},
				pkg:  DefaultRegoPackage,
			},
			args: args{
				serviceID: testdata.MockCloudServiceID1,
				metric: &assessment.Metric{
					Id:       "AutomaticUpdatesEnabled",
					Category: "Endpoint Security",
				},
				baseDir: ".",
				m: map[string]interface{}{
					"automaticUpdates": map[string]interface{}{
						"enabled": true,
					},
				},
				src: &mockMetricsSource{t: t},
			},
			wantResult: func(t *testing.T, got *Result) bool {
				want := &Result{
					Applicable:  true,
					Compliant:   true,
					TargetValue: true,
					Operator:    "==",
					MetricID:    "AutomaticUpdatesEnabled",
					Config: &assessment.MetricConfiguration{
						Operator:       "==",
						TargetValue:    structpb.NewBoolValue(true),
						IsDefault:      true,
						UpdatedAt:      nil,
						MetricId:       "AutomaticUpdatesEnabled",
						CloudServiceId: testdata.MockCloudServiceID1,
					},
				}

				return assert.Equal(t, want, got)
			},
			wantErr: assert.Nil[error],
		},
		{
			name: "updated metric configuration",
			fields: fields{
				qc:   newQueryCache(),
				mrtc: &metricsCache{m: make(map[string][]*assessment.Metric)},
				pkg:  DefaultRegoPackage,
			},
			args: args{
				serviceID: testdata.MockCloudServiceID1,
				metric: &assessment.Metric{
					Id:       "AutomaticUpdatesEnabled",
					Category: "Endpoint Security",
				},
				baseDir: ".",
				m: map[string]interface{}{
					"automaticUpdates": map[string]interface{}{
						"enabled": true,
					},
				},
				src: &updatedMockMetricsSource{mockMetricsSource{t: t}},
			},
			wantResult: func(t *testing.T, got *Result) bool {
				want := &Result{
					Applicable:  true,
					Compliant:   false,
					TargetValue: false,
					Operator:    "==",
					MetricID:    "AutomaticUpdatesEnabled",
					Config: &assessment.MetricConfiguration{
						Operator:       "==",
						TargetValue:    structpb.NewBoolValue(false),
						IsDefault:      false,
						UpdatedAt:      timestamppb.New(time.Date(2022, 12, 1, 0, 0, 0, 0, time.Local)),
						MetricId:       "AutomaticUpdatesEnabled",
						CloudServiceId: testdata.MockCloudServiceID1,
					},
				}

				return assert.Equal(t, want, got)
			},
			wantErr: assert.Nil[error],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := &regoEval{
				qc:   tt.fields.qc,
				mrtc: tt.fields.mrtc,
				pkg:  tt.fields.pkg,
			}
			gotResult, err := re.evalMap(tt.args.baseDir, tt.args.serviceID, tt.args.metric, tt.args.m, tt.args.src)

			tt.wantErr(t, err)
			tt.wantResult(t, gotResult)
		})
	}
}
