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
	"strings"
	"sync"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/api/orchestrator"
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
func (re *regoEval) Eval(evidence *evidence.Evidence, src MetricsSource) (data []*Result, err error) {
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
	cached := re.mrtc.m[key]
	re.mrtc.RUnlock()

	// TODO(lebogg): Try to optimize duplicated code
	if cached == nil {
		metrics, err := src.Metrics()
		if err != nil {
			return nil, fmt.Errorf("could not retrieve metric definitions: %w", err)
		}

		// Lock until we looped through all files
		re.mrtc.Lock()

		// Start with an empty list, otherwise we might copy metrics into the list
		// that are added by a parallel execution - which might occur if both goroutines
		// start at the exactly same time.
		cached = []string{}
		for _, metric := range metrics {
			runMap, err := re.evalMap(baseDir, metric.Id, m, src)
			if err != nil {
				return nil, err
			}

			if runMap != nil {
				cached = append(cached, metric.Id)
				runMap.MetricId = metric.Id

				data = append(data, runMap)
			}
		}

		// Set it and unlock
		re.mrtc.m[key] = cached
		log.Infof("Resource type %v has the following %v applicable metric(s): %v", key, len(re.mrtc.m[key]), re.mrtc.m[key])

		re.mrtc.Unlock()
	} else {
		for _, metric := range cached {
			runMap, err := re.evalMap(baseDir, metric, m, src)
			if err != nil {
				return nil, err
			}

			runMap.MetricId = metric
			data = append(data, runMap)
		}
	}

	return data, nil
}

// HandleMetricEvent takes care of handling metric events, such as evicting cache entries for the
// appropriate metrics.
func (re *regoEval) HandleMetricEvent(event *orchestrator.MetricChangeEvent) (err error) {
	if event.Type == orchestrator.MetricChangeEvent_IMPLEMENTATION_CHANGED {
		log.Infof("Implementation of %s has changed. Clearing cache for this metric", event.MetricId)
	} else if event.Type == orchestrator.MetricChangeEvent_CONFIG_CHANGED {
		log.Infof("Configuration of %s has changed. Clearing cache for this metric", event.MetricId)
	}

	// Evict the cache for the given metric
	re.qc.Evict(event.MetricId)

	return nil
}

func (re *regoEval) evalMap(baseDir string, metric string, m map[string]interface{}, src MetricsSource) (result *Result, err error) {
	var (
		query *rego.PreparedEvalQuery
		key   string
		pkg   string
	)

	// We need to check, if the metric configuration has been changed.
	config, err := src.MetricConfiguration(metric)
	if err != nil {
		return nil, fmt.Errorf("could not fetch metric configuration: %w", err)
	}

	// We build a key out of the metric and its configuration, so we are creating a new Rego implementation
	// if the metric configuration (i.e. its hash) has changed.
	key = fmt.Sprintf("%s-%s", metric, config.Hash())

	query, err = re.qc.Get(key, func(key string) (*rego.PreparedEvalQuery, error) {
		var (
			tx   storage.Transaction
			impl *assessment.MetricImplementation
		)

		// Create paths for bundle directory and utility functions file
		bundle := fmt.Sprintf("%s/policies/bundles/%s/", baseDir, metric)
		operators := fmt.Sprintf("%s/policies/operators.rego", baseDir)

		c := map[string]interface{}{
			"target_value": config.TargetValue.AsInterface(),
			"operator":     config.Operator,
		}

		store := inmem.NewFromObject(c)
		ctx := context.Background()

		tx, err = store.NewTransaction(ctx, storage.WriteParams)
		if err != nil {
			return nil, fmt.Errorf("could not create transaction: %w", err)
		}

		// Convert camelCase metric in under_score_style for package name
		pkg = util.CamelCaseToSnakeCase(metric)

		// TODO (oxisto): we should probably do this using some Rego store implementation

		// Check, if an override in our database exists
		/*var impl assessment.MetricImplementation
		err = re.storage.Get(&impl, assessment.MetricImplementation{
			MetricId: metric,
			Lang:     assessment.MetricImplementation_REGO,
		})
		if err == persistence.ErrRecordNotFound {
			// Load from file if it is not in the database
			data, _ = os.ReadFile(bundle + "metric.rego")
		} else if err != nil {
			return nil, fmt.Errorf("could fetch policy from database: %w", err)
		} else {
			// Take the implementation from the DB
			data = []byte(impl.Code)
		}*/
		impl, err = src.MetricImplementation(assessment.MetricImplementation_REGO, metric)
		if err != nil {
			return nil, fmt.Errorf("could not fetch policy: %w", err)
		}

		err = store.UpsertPolicy(context.Background(), tx, bundle+"metric.rego", []byte(impl.Code))
		if err != nil {
			return nil, fmt.Errorf("could not upsert policy: %w", err)
		}

		query, err := rego.New(
			rego.Query(fmt.Sprintf(`
			applicable = data.clouditor.metrics.%s.applicable;
			compliant = data.clouditor.metrics.%s.compliant;
			operator = data.clouditor.operator;
			target_value = data.clouditor.target_value`, pkg, pkg)),
			rego.Package("clouditor.metrics"),
			rego.Store(store),
			rego.Transaction(tx),
			rego.Load(
				[]string{
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

func (qc *queryCache) Empty() {
	qc.Lock()
	defer qc.Unlock()

	for k := range qc.cache {
		delete(qc.cache, k)
	}
}

// Evict deletes all keys from the cache that belong to the given metric.
func (qc *queryCache) Evict(metric string) {
	qc.Lock()
	defer qc.Unlock()

	// Look for keys that begin with the metric
	for k := range qc.cache {
		if strings.HasPrefix(k, metric) {
			delete(qc.cache, k)
		}
	}
}
