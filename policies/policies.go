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
	"os"
	"strings"

	"clouditor.io/clouditor/api/evidence"
	"github.com/open-policy-agent/opa/rego"
)

// applicableMetrics stores a list of applicable metrics per resourceType
var applicableMetrics = make(map[string][]string)

func RunEvidence(evidence *evidence.Evidence) ([]map[string]interface{}, error) {
	data := make([]map[string]interface{}, 0)
	var baseDir string = "."

	var m = evidence.Resource.GetStructValue().AsMap()

	var types []string

	if rawTypes, ok := m["type"].([]interface{}); ok {
		types = make([]string, len(rawTypes))
	} else {
		return nil, fmt.Errorf("got type '%T' but wanted '[]interface {}'. Check if resource types are specified ", rawTypes)
	}
	for i, v := range m["type"].([]interface{}) {
		if t, ok := v.(string); !ok {
			return nil, fmt.Errorf("got type '%T' but wanted 'string'", t)
		} else {
			types[i] = t
		}
	}

	// TODO(lebogg): Try to optimize duplicated code
	if key := createKey(types); applicableMetrics[key] == nil {
		files, err := scanBundleDir(baseDir)
		if err != nil {
			return nil, fmt.Errorf("could not load metric bundles: %w", err)
		}

		for _, fileInfo := range files {
			runMap, err := RunMap(baseDir, fileInfo.Name(), m)
			if err != nil {
				return nil, err
			}

			if runMap != nil {
				metricId := fileInfo.Name()
				applicableMetrics[key] = append(applicableMetrics[key], metricId)
				runMap["metricId"] = metricId

				data = append(data, runMap)
			}
		}
	} else {
		for _, metric := range applicableMetrics[key] {
			runMap, err := RunMap(baseDir, metric, m)
			if err != nil {
				return nil, err
			}

			runMap["metricId"] = metric
			data = append(data, runMap)
		}
	}

	return data, nil
}

func RunMap(baseDir string, metric string, m map[string]interface{}) (data map[string]interface{}, err error) {
	var (
		ok bool
	)

	// Create paths for bundle directory and utility functions file
	bundle := fmt.Sprintf("%s/policies/bundles/%s/", baseDir, metric)
	operators := fmt.Sprintf("%s/policies/operators.rego", baseDir)

	ctx := context.TODO()
	r, err := rego.New(
		rego.Query("data.clouditor"),
		rego.Load(
			[]string{
				bundle + "metric.rego",
				bundle + "data.json",
				operators,
			},
			nil),
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
		fmt.Println(data)
		return data, nil
	}
}

func scanBundleDir(baseDir string) ([]os.FileInfo, error) {
	dirname := baseDir + "/policies/bundles"

	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	files, err := f.Readdir(-1)
	_ = f.Close()
	if err != nil {
		return nil, err
	}

	return files, err
}

func createKey(types []string) (key string) {
	key = strings.Join(types, "-")
	key = strings.ReplaceAll(key, " ", "")
	return
}
