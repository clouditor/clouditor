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
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

// TODO(lebogg): Remove after testing
var log *logrus.Entry = logrus.WithField("component", "policies")

// applicableMetrics stores a list of applicable metrics per resourceType
var applicableMetrics = make(map[string][]string)

func RunEvidence(evidence *evidence.Evidence) ([]map[string]interface{}, error) {
	log.Println("Run evidence for resourceId", evidence.ResourceId, " and ID: ", evidence.Id)
	data := make([]map[string]interface{}, 0)
	var baseDir string = "."
	// check, if we are in the root of Clouditor
	if _, err := os.Stat("policies"); os.IsNotExist(err) {
		// in tests, we are relative to our current package
		baseDir = ".."
	}

	var m = evidence.Resource.GetStructValue().AsMap()

	var types []string

	if rawTypes, ok := m["type"].([]interface{}); ok {
		types = make([]string, len(rawTypes))
	} else {
		return nil, fmt.Errorf("got type '%T' but wanted '[]interface {}'. Check if resource types are specified ", rawTypes)
	}
	for i, v := range m["type"].([]interface{}) {
		// TODO(all): type assertion check good or unnecessary because we assume resoruceTypes to be always set as intended ([]string)?
		if t, ok := v.(string); !ok {
			return nil, fmt.Errorf("got type '%T' but wanted 'string'", t)
		} else {
			types[i] = t
		}
	}
	if key := strings.Join(types, "-"); applicableMetrics[key] == nil {
		files, err := scanBundleDir()
		if err != nil {
			return nil, fmt.Errorf("could not load metric bundles: %v", err)
		}

		for _, fileInfo := range files {
			file := fmt.Sprintf("%s/policies/policyBundles/%s", baseDir, fileInfo.Name())
			runMap, err := RunMap(file, m)
			if err != nil {
				return nil, err
			}
			if runMap != nil {
				data = append(data, runMap)

				if metric := applicableMetrics[key]; metric == nil {
					applicableMetrics[key] = []string{fileInfo.Name()}
				}
				applicableMetrics[key] = append(applicableMetrics[key], fileInfo.Name())
			}
		}
	} else {
		for _, metric := range applicableMetrics[key] {
			file := fmt.Sprintf("%s/policies/policyBundles/%s", baseDir, metric)
			runMap, err := RunMap(file, m)
			if err != nil {
				return nil, err
			}
			data = append(data, runMap)
		}
	}
	return data, nil
}

func scanBundleDir() ([]os.FileInfo, error) {

	dirname := "./policyBundles"

	f, err := os.Open(dirname)
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.Readdir(-1)
	_ = f.Close()
	if err != nil {
		log.Fatal(err)
	}

	return files, err
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
