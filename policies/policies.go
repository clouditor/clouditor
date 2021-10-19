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
	"clouditor.io/clouditor/api/evidence"
	"context"
	"errors"
	"fmt"
	"github.com/open-policy-agent/opa/rego"
	"os"
)

// applicableMetrics stores a list of applicable metrics per resourceType
var applicableMetrics = make(map[string][]int)

func RunEvidence(evidence *evidence.Evidence) ([]map[string]interface{}, error) {
	// TODO(lebogg): Think about magic number 10
	data := make([]map[string]interface{}, 0)
	var baseDir string = "."
	// check, if we are in the root of Clouditor
	if _, err := os.Stat("policies"); os.IsNotExist(err) {
		// in tests, we are relative to our current package
		// TODO(lebogg): Is it possible to do go to the root dir with a function?
		baseDir = ".."
	}

	var m = evidence.Resource.GetStructValue().AsMap()

	// TODO(lebogg): For now check via resourceIDs. Later with (concatenation of) resourceTypes (more efficient hash)
	if resID := evidence.GetResourceId(); applicableMetrics[resID] == nil {
		// TODO(lebogg): Replace magic number of 10 (current metrics)
		for i := 1; i < 11; i++ {
			// ToDo(lebogg): Test if direction to each bundle works
			file := fmt.Sprintf("%s/policies/bundle%d", baseDir, i)
			runMap, err := RunMap(file, m)
			if err != nil {
				return nil, err
			}
			if runMap != nil {
				data = append(data, runMap)

				if metric := applicableMetrics[resID]; metric == nil {
					applicableMetrics[resID] = []int{i}
				}
				applicableMetrics[resID] = append(applicableMetrics[resID], i)
			}
		}
	} else {
		for _, metric := range applicableMetrics[resID] {
			file := fmt.Sprintf("%s/policies/bundle%d", baseDir, metric)
			runMap, err := RunMap(file, m)
			if err != nil {
				return nil, err
			}
			data = append(data, runMap)
		}
	}
	return data, nil
}

func RunMap(bundle string, m map[string]interface{}) (data map[string]interface{}, err error) {
	var (
		ok bool
	)

	ctx := context.TODO()
	r, err := rego.New(
		rego.Query("data.clouditor"),
		rego.LoadBundle(bundle),
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
	} else if data["applicable"] == false {
		return nil, nil
	} else {
		return data, nil
	}
}
