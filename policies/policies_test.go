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
	"encoding/json"
	"testing"

	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/policies"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestRun(t *testing.T) {
	var (
		m    map[string]interface{}
		data []map[string]interface{}
		s    *structpb.Struct
		err  error
	)

	j := `{
		"atRestEncryption": {
			"algorithm": "AES-256",
			"enabled": true,
			"keyManager": "Microsoft.Storage"
		},
		"creationTime": 1621086669,
		"httpEndpoint": {
			"transportEncryption": {
				"algorithm": "TLS",
				"enabled": true,
				"enforced": true,
				"tlsVersion": 1.3
			},
			"url": "https://aybazestorage.blob.core.windows.net/"
		},
		"id": "/subscriptions/e3ed0e96-57bc-4d81-9594-f239540cd77a/resourceGroups/titan/providers/Microsoft.Storage/storageAccounts/aybazestorage",
		"name": "aybazestorage"
	}`

	err = json.Unmarshal([]byte(j), &m)

	assert.Nil(t, err)

	s, err = structpb.NewStruct(m)

	assert.Nil(t, err)

	data, err = policies.RunEvidence(&evidence.Evidence{
		Resource:   structpb.NewStructValue(s),
		ResourceId: "/subscriptions/e3ed0e96-57bc-4d81-9594-f239540cd77a/resourceGroups/titan/providers/Microsoft.Storage/storageAccounts/aybazestorage",
	})

	assert.Nil(t, err)

	// Test metric TransportEncryptionAlgorithm
	assert.NotNil(t, data[0])
	assert.Equal(t, true, data[0]["compliant"])
	assert.Equal(t, true, data[0]["applicable"])

	// Test metric TLSVersion
	assert.NotNil(t, data[1])
	assert.Equal(t, true, data[1]["compliant"])
	assert.Equal(t, true, data[1]["applicable"])

	// Test metric TransportEncryptionEnabled
	assert.NotNil(t, data[2])
	assert.Equal(t, true, data[2]["compliant"])
	assert.Equal(t, true, data[2]["applicable"])

	// Test metric TransportEncryptionEnforced
	assert.NotNil(t, data[3])
	assert.Equal(t, true, data[3]["compliant"])
	assert.Equal(t, true, data[3]["applicable"])

	// Test metric EncryptionAtRestEnabled
	assert.NotNil(t, data[4])
	assert.Equal(t, true, data[4]["compliant"])
	assert.Equal(t, true, data[4]["applicable"])

	// Test metric EncryptionAtRestAlgorithm
	assert.NotNil(t, data[5])
	assert.Equal(t, true, data[5]["compliant"])
	assert.Equal(t, true, data[5]["applicable"])

	// Testing VM
	j = `{
		"bootLog" : {
			"enabled" : true,
			"retentionPeriod" : 36
		},
		"OSLog" : {
			"enabled" : true,
			"retentionPeriod" : 90
		},
		"id": "/subscriptions/e3ed0e96-57bc-4d81-9594-f239540cd77a/resourceGroups/titan/providers/Microsoft.Storage/virtualMachine/mockvm",
		"name": "aybazestorage"
	}`

	m = make(map[string]interface{})
	err = json.Unmarshal([]byte(j), &m)

	assert.Nil(t, err)

	s, err = structpb.NewStruct(m)

	assert.Nil(t, err)

	data, err = policies.RunEvidence(&evidence.Evidence{
		Resource:   structpb.NewStructValue(s),
		ResourceId: "/subscriptions/e3ed0e96-57bc-4d81-9594-f239540cd77a/resourceGroups/titan/providers/Microsoft.Storage/virtualMachine/mockvm",
	})
	assert.Nil(t, err)
	assert.NotNil(t, data)

	// Test metric BootLoggingEnabled
	assert.NotNil(t, data[0])
	assert.Equal(t, true, data[0]["compliant"])
	assert.Equal(t, true, data[0]["applicable"])

	// Test metric BootLoggingRetention
	assert.NotNil(t, data[1])
	assert.Equal(t, true, data[1]["compliant"])
	assert.Equal(t, true, data[1]["applicable"])

	// Test metric OSLoggingEnabled
	assert.NotNil(t, data[2])
	assert.Equal(t, true, data[2]["compliant"])
	assert.Equal(t, true, data[2]["applicable"])

	// Test metricOSLoggingRetention
	assert.NotNil(t, data[3])
	assert.Equal(t, true, data[3]["compliant"])
	assert.Equal(t, true, data[3]["applicable"])

	// Repeat to check if only the metrics are evaluated that are needed
	data, err = policies.RunEvidence(&evidence.Evidence{
		Resource:   structpb.NewStructValue(s),
		ResourceId: "/subscriptions/e3ed0e96-57bc-4d81-9594-f239540cd77a/resourceGroups/titan/providers/Microsoft.Storage/virtualMachine/mockvm",
	})
	assert.Nil(t, err)

}
