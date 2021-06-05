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
