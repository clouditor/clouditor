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

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/util"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/storage/inmem"
	"google.golang.org/grpc/codes"
)

// DefaultRegoPackage is the default package name for the Rego files
const DefaultRegoPackage = "clouditor.metrics"

type regoEval struct {
	// qc contains cached Rego queries
	qc *queryCache

	// mrtc stores a list of applicable metrics per toolID and resourceType
	mrtc *metricsCache

	// pkg is the base package name that is used in the Rego files
	pkg string
}

type queryCache struct {
	sync.Mutex
	cache map[string]*rego.PreparedEvalQuery
}

type orElseFunc func(key string) (query *rego.PreparedEvalQuery, err error)

type RegoEvalOption func(re *regoEval)

// WithPackageName is an option to configure the package name
func WithPackageName(pkg string) RegoEvalOption {
	return func(re *regoEval) {
		re.pkg = pkg
	}
}

func NewRegoEval(opts ...RegoEvalOption) PolicyEval {
	re := regoEval{
		mrtc: &metricsCache{m: make(map[string][]string)},
		qc:   newQueryCache(),
		pkg:  DefaultRegoPackage,
	}

	for _, o := range opts {
		o(&re)
	}

	return &re
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

	types, err = evidence.ResourceTypes()
	if err != nil {
		return nil, fmt.Errorf("could not extract resource types from evidence: %w", err)
	}

	key := createKey(evidence, types)

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
			// Try to evaluate it and check, if the metric is applicable (in which case we are getting a result). We
			// need to differentiate here between an execution error (which might be temporary) and an error if the
			// metric configuration or implementation is not found. The latter case happens if the metric is not
			// assessed within the Clouditor toolset but we need to know that the metric exists, e.g., because it is
			// evaluated by an external tool. In this case, we can just pretend that the metric is not applicable for us
			// and continue.
			runMap, err := re.evalMap(baseDir, evidence.CloudServiceId, metric.Id, m, src)
			if err != nil {
				// Try to retrieve the gRPC status from the error, to check if the metric implementation just does not exist.
				status, ok := api.StatusFromWrappedError(err)
				if ok && status.Code() == codes.NotFound {
					log.Warnf("Resource type %v ignored metric %v because of its missing implementation or default configuration ", key, metric.Id)
					// In this case, we can continue
					continue
				}

				// Otherwise, we are not really in a state where our cache is valid, so we mark it as not cached at all.
				re.mrtc.m[key] = nil

				// Unlock, to avoid deadlock and return from here with the error
				re.mrtc.Unlock()
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
			runMap, err := re.evalMap(baseDir, evidence.CloudServiceId, metric, m, src)
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

func (re *regoEval) evalMap(baseDir string, serviceID, metricID string, m map[string]interface{}, src MetricsSource) (result *Result, err error) {
	var (
		query  *rego.PreparedEvalQuery
		key    string
		pkg    string
		prefix string
	)

	// We need to check, if the metric configuration has been changed.
	config, err := src.MetricConfiguration(serviceID, metricID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch metric configuration for metric %s: %w", metricID, err)
	}

	// We build a key out of the metric and its configuration, so we are creating a new Rego implementation
	// if the metric configuration (i.e. its hash) for a particular service has changed.
	key = fmt.Sprintf("%s-%s-%s", metricID, serviceID, config.Hash())

	query, err = re.qc.Get(key, func(key string) (*rego.PreparedEvalQuery, error) {
		var (
			tx   storage.Transaction
			impl *assessment.MetricImplementation
		)

		// Create paths for bundle directory and utility functions file
		bundle := fmt.Sprintf("%s/policies/bundles/%s/", baseDir, metricID)
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

		prefix = re.pkg

		// Convert camelCase metric in under_score_style for package name
		pkg = util.CamelCaseToSnakeCase(metricID)

		impl, err = src.MetricImplementation(assessment.MetricImplementation_REGO, metricID)
		if err != nil {
			return nil, fmt.Errorf("could not fetch policy for metric %s: %w", metricID, err)
		}

		err = store.UpsertPolicy(context.Background(), tx, bundle+"metric.rego", []byte(impl.Code))
		if err != nil {
			return nil, fmt.Errorf("could not upsert policy: %w", err)
		}

		query, err := rego.New(
			rego.Query(fmt.Sprintf(`
			applicable = data.%s.%s.applicable;
			compliant = data.%s.%s.compliant;
			operator = data.clouditor.operator;
			target_value = data.clouditor.target_value`, prefix, pkg, prefix, pkg)),
			rego.Package(prefix),
			rego.Store(store),
			rego.Transaction(tx),
			rego.Load(
				[]string{
					operators,
				},
				nil),
		).PrepareForEval(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not prepare rego evaluation for metric %s: %w", metricID, err)
		}

		err = store.Commit(ctx, tx)
		if err != nil {
			return nil, fmt.Errorf("could not commit transaction: %w", err)
		}

		return &query, nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not fetch cached query for metric %s: %w", metricID, err)
	}

	results, err := query.Eval(context.Background(), rego.EvalInput(m))
	if err != nil {
		return nil, fmt.Errorf("could not evaluate rego policy: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results. probably the package name of metric %s is wrong", metricID)
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
