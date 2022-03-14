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
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidence"
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/storage/inmem"
	"os"
	"strings"
	"unicode"
)

// applicableMetrics stores a list of applicable metrics per resourceType
var applicableMetrics = make(map[string][]string)

type Result struct {
	Applicable  bool
	Compliant   bool
	TargetValue interface{} `mapstructure:"target_value"`
	Operator    string
	MetricId    string
}

// MetricConfigurationSource can be used to retrieve a metric configuration for a particular metric (and target service)
type MetricConfigurationSource interface {
	MetricConfiguration(metric string) (*assessment.MetricConfiguration, error)
}

func RunEvidence(evidence *evidence.Evidence, holder MetricConfigurationSource) ([]*Result, error) {
	data := make([]*Result, 0)
	var baseDir = "."

	var m = evidence.Resource.GetStructValue().AsMap()

	var types []string

	if rawTypes, ok := m["type"].([]interface{}); ok {
		if len(rawTypes) != 0 {
			types = make([]string, len(rawTypes))
		} else {
			return nil, fmt.Errorf("list of types is empty")
		}
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
			runMap, err := RunMap(baseDir, fileInfo.Name(), m, holder)
			if err != nil {
				return nil, err
			}

			if runMap != nil {
				metricId := fileInfo.Name()
				applicableMetrics[key] = append(applicableMetrics[key], metricId)
				runMap.MetricId = metricId

				data = append(data, runMap)
			}
		}
	} else {
		for _, metric := range applicableMetrics[key] {
			runMap, err := RunMap(baseDir, metric, m, holder)
			if err != nil {
				return nil, err
			}

			runMap.MetricId = metric
			data = append(data, runMap)
		}
	}

	return data, nil
}

func RunMap(baseDir string, metric string, m map[string]interface{}, holder MetricConfigurationSource) (result *Result, err error) {
	var (
		tx storage.Transaction
	)

	// Create paths for bundle directory and utility functions file
	bundle := fmt.Sprintf("%s/policies/bundles/%s/", baseDir, metric)
	operators := fmt.Sprintf("%s/policies/operators.rego", baseDir)

	config, err := holder.MetricConfiguration(metric)
	if err != nil {
		return nil, fmt.Errorf("could not fetch metric configuration: %w", err)
	}

	c := map[string]interface{}{
		"target_value": config.TargetValue.AsInterface(),
		"operator":     config.Operator,
	}

	// Convert camelCase metric in under_score_style for package name
	metric = camelCaseToSnakeCase(metric)

	store := inmem.NewFromObject(c)
	ctx := context.Background()

	tx, err = store.NewTransaction(ctx, storage.WriteParams)
	if err != nil {
		return nil, fmt.Errorf("could not create transaction: %w", err)
	}

	r, err := rego.New(
		rego.Query(fmt.Sprintf(`x = {
			"applicable": data.clouditor.metrics.%s.applicable, 
			"compliant": data.clouditor.metrics.%s.compliant, 
			"operator": data.clouditor.operator,
			"target_value": data.clouditor.target_value,
		}`, metric, metric)),
		rego.Package("clouditor.metrics"),
		rego.Store(store),
		rego.Transaction(tx),
		rego.Input(m),
		rego.Load(
			[]string{
				bundle + "metric.rego",
				operators,
			},
			nil),
	).PrepareForEval(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not prepare rego evaluation: %w", err)
	}

	err = store.Commit(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	results, err := r.Eval(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not evaluate rego policy: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results. probably the package name of the metric is wrong")
	}

	result = new(Result)
	err = mapstructure.Decode(results[0].Bindings["x"], result)

	if err != nil {
		return nil, fmt.Errorf("expected data is not a map[string]interface{}: %w", err)
	}

	if !result.Applicable {
		return nil, nil
	} else {
		return result, nil
	}
}

// camelCaseToSnakeCase converts a `camelCase` string to `snake_case`
func camelCaseToSnakeCase(input string) string {
	if input == "" {
		return ""
	}

	snakeCase := make([]rune, 0, len(input))
	runeArray := []rune(input)

	for i := range runeArray {
		if i > 0 && marksNewWord(i, runeArray) {
			snakeCase = append(snakeCase, '_', unicode.ToLower(runeArray[i]))
		} else {
			snakeCase = append(snakeCase, unicode.ToLower(runeArray[i]))
		}
	}

	return string(snakeCase)
}

// marksNewWord checks if the current character starts a new word excluding the first word
func marksNewWord(i int, input []rune) bool {

	if i >= len(input) {
		return false
	}

	// If previous or following rune/character is lowercase or rune is a number than it is a new word
	if i < len(input)-1 && unicode.IsUpper(input[i]) && unicode.IsLower(input[i+1]) {
		return true
	} else if i > 0 && unicode.IsLower(input[i-1]) && unicode.IsUpper(input[i]) {
		return true
	} else if unicode.IsDigit(input[i]) {
		return true
	}

	return false
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
