// Copyright 2021-2022 Fraunhofer AISEC
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
	"fmt"
	"sync"

	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/internal/util"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/storage/inmem"
)

type regoEval struct {
	// qc contains cached Rego queries
	qc *queryCache

	// mrtc stores a list of applicable metrics per resourceType
	mrtc *metricsResourceTypeCache
}

type queryCache struct {
	sync.Mutex
	cache map[string]*rego.PreparedEvalQuery
}

type orElseFunc func(key string) (query *rego.PreparedEvalQuery, err error)

func NewRegoEval() PolicyEval {
	return &regoEval{
		mrtc: &metricsResourceTypeCache{m: make(map[string][]string)},
		qc:   newQueryCache(),
	}
}

// Eval evaluates a given evidence against all available Rego policies and returns the result of all policies that were
// considered to be applicable.
func (re *regoEval) Eval(evidence *evidence.Evidence, holder MetricConfigurationSource) (data []*Result, err error) {
	var (
		baseDir string
		m       map[string]interface{}
		types   []string
	)

	baseDir = "."
	m = evidence.Resource.GetStructValue().AsMap()

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

	key := createKey(types)

	re.mrtc.RLock()
	metrics := re.mrtc.m[key]
	re.mrtc.RUnlock()

	// TODO(lebogg): Try to optimize duplicated code
	if metrics == nil {
		files, err := scanBundleDir(baseDir)
		if err != nil {
			return nil, fmt.Errorf("could not load metric bundles: %w", err)
		}

		// Lock until we looped through all files
		re.mrtc.Lock()

		// Start with an empty list, otherwise we might copy metrics into the list
		// that are added by a parallel execution - which might occur if both goroutines
		// start at the exactly same time.
		metrics = []string{}
		for _, fileInfo := range files {
			runMap, err := re.evalMap(baseDir, fileInfo.Name(), m, holder)
			if err != nil {
				return nil, err
			}

			if runMap != nil {
				metricId := fileInfo.Name()
				metrics = append(metrics, metricId)
				runMap.MetricId = metricId

				data = append(data, runMap)
			}
		}

		// Set it and unlock
		re.mrtc.m[key] = metrics
		log.Infof("Resource type %v has the following %v applicable metric(s): %v", key, len(re.mrtc.m[key]), re.mrtc.m[key])

		re.mrtc.Unlock()
	} else {
		for _, metric := range metrics {
			runMap, err := re.evalMap(baseDir, metric, m, holder)
			if err != nil {
				return nil, err
			}

			runMap.MetricId = metric
			data = append(data, runMap)
		}
	}

	return data, nil
}

func (re *regoEval) evalMap(baseDir string, metric string, m map[string]interface{}, holder MetricConfigurationSource) (result *Result, err error) {
	var (
		query *rego.PreparedEvalQuery
		key   string
	)

	// We need to check, if the metric configuration has been changed. Any caching of this
	// configuration will be done by the MetricConfigurationSource.
	config, err := holder.MetricConfiguration(metric)
	if err != nil {
		return nil, fmt.Errorf("could not fetch metric configuration: %w", err)
	}

	// We build a key out of the metric and its configuration, so we are creating a new Rego implementation
	// if the metric configuration (i.e. its hash) has changed.
	key = fmt.Sprintf("%s-%s", metric, config.Hash())

	query, err = re.qc.Get(key, func(key string) (*rego.PreparedEvalQuery, error) {
		var (
			tx storage.Transaction
		)

		// Create paths for bundle directory and utility functions file
		bundle := fmt.Sprintf("%s/policies/bundles/%s/", baseDir, metric)
		operators := fmt.Sprintf("%s/policies/operators.rego", baseDir)

		c := map[string]interface{}{
			"target_value": config.TargetValue.AsInterface(),
			"operator":     config.Operator,
		}

		// Convert camelCase metric in under_score_style for package name
		metric = util.CamelCaseToSnakeCase(metric)

		store := inmem.NewFromObject(c)
		ctx := context.Background()

		tx, err = store.NewTransaction(ctx, storage.WriteParams)
		if err != nil {
			return nil, fmt.Errorf("could not create transaction: %w", err)
		}

		query, err := rego.New(
			rego.Query(fmt.Sprintf(`
			applicable = data.clouditor.metrics.%s.applicable;
			compliant = data.clouditor.metrics.%s.compliant;
			operator = data.clouditor.operator;
			target_value = data.clouditor.target_value`, metric, metric)),
			rego.Package("clouditor.metrics"),
			rego.Store(store),
			rego.Transaction(tx),
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

		return &query, nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not fetch cached query: %w", err)
	}

	results, err := query.Eval(context.Background(), rego.EvalInput(m))
	if err != nil {
		return nil, fmt.Errorf("could not evaluate rego policy: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results. probably the package name of the metric is wrong")
	}

	result = &Result{
		Applicable:  results[0].Bindings["applicable"].(bool),
		Compliant:   results[0].Bindings["compliant"].(bool),
		Operator:    results[0].Bindings["operator"].(string),
		TargetValue: results[0].Bindings["target_value"],
	}

	if !result.Applicable {
		return nil, nil
	} else {
		return result, nil
	}
}

func newQueryCache() *queryCache {
	return &queryCache{
		cache: make(map[string]*rego.PreparedEvalQuery),
	}
}

// Get returns the prepared query for the given key. If the key was not found in the cache,
// the orElse function is executed to populate the cache.
func (qc *queryCache) Get(key string, orElse orElseFunc) (query *rego.PreparedEvalQuery, err error) {
	var (
		ok bool
	)

	// Lock the cache
	qc.Lock()
	// And defer the unlock
	defer qc.Unlock()

	// Check, if query is contained in the cache
	query, ok = qc.cache[key]
	if ok {
		return
	}

	// Otherwise, the orElse function is executed to fetch the query
	query, err = orElse(key)
	if err != nil {
		return nil, err
	}

	// Update the cache
	qc.cache[key] = query
	return
}
