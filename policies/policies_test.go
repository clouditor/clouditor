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

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/policies"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestRun(t *testing.T) {
	var (
		m    map[string]interface{}
		data map[string]interface{}
		s    *structpb.Struct
		err  error
	)

	j := `{
		"atRestEncryption": {
			"algorithm": "AES-265",
			"enabled": true,
			"keyManager": "Microsoft.Storage"
		},
		"creationTime": 1621086669,
		"httpEndpoint": {
			"transportEncryption": {
				"enabled": true,
				"enforced": true,
				"tlsVersion": "TLS1_2"
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

	data, err = policies.RunEvidence("metric1.rego", &assessment.Evidence{
		Resource: structpb.NewStructValue(s),
	})

	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, true, data["compliant"])
}
