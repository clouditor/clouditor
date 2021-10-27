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

package policies_test

import (
	"clouditor.io/clouditor/voc"
	"testing"

	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/policies"
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

func TestRunEvidence(t *testing.T) {
	type fields struct {
		resource   voc.IsCloudResource
		evidenceID string
	}
	tests := []struct {
		name       string
		fields     fields
		applicable bool
		compliant  bool
		wantErr    bool
	}{
		{
			name: "ObjectStorage: Compliant Case",
			fields: fields{
				resource: voc.ObjectStorage{
					Storage: &voc.Storage{
						CloudResource: &voc.CloudResource{
							ID:           mockObjStorage1ResourceID,
							Name:         "aybazestorage",
							CreationTime: 1621086669,
							Type:         []string{"ObjectStorage", "Storage", "Resource"},
							GeoLocation:  voc.GeoLocation{},
						},
						AtRestEncryption: &voc.AtRestEncryption{
							Algorithm: "AES-256",
							Enabled:   true,
						},
					},
					HttpEndpoint: &voc.HttpEndpoint{
						Functionality: nil,
						Authenticity:  nil,
						TransportEncryption: &voc.TransportEncryption{
							Enforced:   true,
							Enabled:    true,
							TlsVersion: "1.3",
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
			applicable: true,
			compliant:  true,
			wantErr:    false,
		}, {
			name: "ObjectStorage: Non-Compliant Case",
			fields: fields{
				resource: voc.ObjectStorage{
					Storage: &voc.Storage{
						CloudResource: &voc.CloudResource{
							ID:           mockObjStorage2ResourceID,
							Name:         "aybazestorage",
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
			applicable: true,
			compliant:  false,
			wantErr:    false,
		},
		{
			name: "VM:Compliant Case",
			fields: fields{
				resource: voc.VirtualMachine{
					Compute: &voc.Compute{
						CloudResource: &voc.CloudResource{
							ID:   mockVM1ResourceID,
							Name: "aybazestorage",
							Type: []string{"Compute", "Virtual Machine", "Resource"},
						}},
					NetworkInterface: nil,
					BlockStorage:     nil,
					BootLog: &voc.BootLog{
						Log: &voc.Log{
							Output:          []voc.ResourceID{"SomeResourceId1", "SomeResourceId2"},
							Enabled:         true,
							RetentionPeriod: 36,
						},
					},
					OSLog: &voc.OSLog{
						Log: &voc.Log{
							Output:          []voc.ResourceID{"SomeResourceId2"},
							Enabled:         true,
							RetentionPeriod: 36,
						},
					},
				},
				evidenceID: mockVM1EvidenceID,
			},
			applicable: true,
			compliant:  true,
			wantErr:    false,
		},
		{
			name: "VM:Non-Compliant Case",
			fields: fields{
				resource: voc.VirtualMachine{
					Compute: &voc.Compute{
						CloudResource: &voc.CloudResource{
							ID:   mockVM2ResourceID,
							Name: "aybazestorage",
							Type: []string{"Compute", "Virtual Machine", "Resource"},
						}},
					NetworkInterface: nil,
					BlockStorage:     nil,
					BootLog: &voc.BootLog{
						Log: &voc.Log{
							Output:          []voc.ResourceID{},
							Enabled:         false,
							RetentionPeriod: 1,
						},
					},
					OSLog: &voc.OSLog{
						Log: &voc.Log{
							Output:          []voc.ResourceID{"SomeResourceId3"},
							Enabled:         false,
							RetentionPeriod: 1,
						},
					},
				},
				evidenceID: mockVM2EvidenceID,
			},
			applicable: true,
			compliant:  false,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		resource, err := voc.ToStruct(tt.fields.resource)
		assert.Nil(t, err)
		t.Run(tt.name, func(t *testing.T) {
			results, err := policies.RunEvidence(&evidence.Evidence{
				Id:         tt.fields.evidenceID,
				ResourceId: string(tt.fields.resource.GetID()),
				Resource:   resource,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("RunEvidence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotEmpty(t, results)
			for _, result := range results {
				assert.Equal(t, tt.applicable, result["applicable"].(bool))
				assert.Equal(t, tt.compliant, result["compliant"].(bool))
			}
		})
	}
}
