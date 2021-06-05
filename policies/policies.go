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
	"context"
	"encoding/json"
	"fmt"

	"github.com/open-policy-agent/opa/rego"
)

func Run(file string) (err error) {
	var m map[string]interface{}
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
	if err != nil {
		return fmt.Errorf("could not unmarshal JSON: %w", err)
	}

	ctx := context.TODO()
	r, err := rego.New(
		rego.Query("data.clouditor"),
		rego.Load([]string{"tls.rego"}, nil),
	).PrepareForEval(ctx)
	if err != nil {
		return fmt.Errorf("could not prepare rego evaluation: %w", err)
	}

	results, err := r.Eval(ctx, rego.EvalInput(m))
	if err != nil {
		return fmt.Errorf("could not evaluate rego policy: %w", err)
	}

	fmt.Printf("%+v", results[0].Expressions[0].Value)

	return nil
}
