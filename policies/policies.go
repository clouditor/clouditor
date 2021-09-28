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
	"errors"
	"fmt"

	"clouditor.io/clouditor/api/evidence"
	"github.com/open-policy-agent/opa/rego"
)

func RunEvidence(file string, evidence *evidence.Evidence) (data map[string]interface{}, err error) {
	var m = evidence.Resource.GetStructValue().AsMap()

	return RunMap(file, m)
}

func RunMap(file string, m map[string]interface{}) (data map[string]interface{}, err error) {
	var (
		ok bool
	)

	ctx := context.TODO()
	r, err := rego.New(
		rego.Query("data.clouditor"),
		rego.Load([]string{file}, nil),
	).PrepareForEval(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not prepare rego evaluation: %w", err)
	}

	results, err := r.Eval(ctx, rego.EvalInput(m))
	if err != nil {
		return nil, fmt.Errorf("could not evaluate rego policy: %w", err)
	}

	if data, ok = results[0].Expressions[0].Value.(map[string]interface{}); !ok {
		return nil, errors.New("expected data is not a map[string]interface{}")
	}

	return
}
