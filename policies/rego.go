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
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/util"
	"github.com/open-policy-agent/opa/v1/rego"
	"github.com/open-policy-agent/opa/v1/storage"
	"github.com/open-policy-agent/opa/v1/storage/inmem"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		mrtc: &metricsCache{m: make(map[string][]*assessment.Metric)},
		qc:   newQueryCache(),
		pkg:  DefaultRegoPackage,
	}

	for _, o := range opts {
		o(&re)
	}

	return &re
}

// Eval evaluates a given evidence against all available Rego policies and returns the result of all policies that were
// considered to be applicable. In order to avoid multiple unwrapping, the callee will already supply an unwrapped
// ontology resource in r.
func (re *regoEval) Eval(evidence *evidence.Evidence, r ontology.IsResource, related map[string]ontology.IsResource, src MetricsSource) (data []*CombinedResult, err error) {
	var (
		baseDir string
		m       map[string]any
		mm      map[string]any
		types   []string
	)

	baseDir = "."

	m, err = ontology.ResourceMap(r)
	if err != nil {
		return nil, err
	}

	if related != nil {
		am := make(map[string]interface{})
		for key, value := range related {
			mm, err = ontology.ResourceMap(value)
			if err != nil {
				return nil, err
			}
			am[key] = mm
		}

		m["related"] = am
	}

	types = ontology.ResourceTypes(r)
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
		cached = []*assessment.Metric{}
		for _, metric := range metrics {
			// Try to evaluate it and check, if the metric is applicable (in which case we are getting a result). We
			// need to differentiate here between an execution error (which might be temporary) and an error if the
			// metric configuration or implementation is not found. The latter case happens if the metric is not
			// assessed within the Clouditor toolset but we need to know that the metric exists, e.g., because it is
			// evaluated by an external tool. In this case, we can just pretend that the metric is not applicable for us
			// and continue.
			runMap, err := re.evalMap(baseDir, evidence.CertificationTargetId, metric, m, src)
			if err != nil {
				// Try to retrieve the gRPC status from the error, to check if the metric implementation just does not exist.
				status, ok := status.FromError(err)
				if ok && status.Code() == codes.NotFound &&
					(strings.Contains(status.Message(), "implementation for metric not found") ||
						strings.Contains(status.Message(), "metric configuration not found for metric")) {
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
				cached = append(cached, metric)

				data = append(data, runMap)
			}
		}

		// Set it and unlock
		re.mrtc.m[key] = cached
		log.Infof("Resource type %v has the following %v applicable metric(s): %v", key, len(re.mrtc.m[key]), idsOf(re.mrtc.m[key]))

		re.mrtc.Unlock()
	} else {
		for _, metric := range cached {
			runMap, err := re.evalMap(baseDir, evidence.CertificationTargetId, metric, m, src)
			if err != nil {
				return nil, err
			}
			// Add runMap to data only if metric was applicable. runMap=nil and err=nil means the metric was not
			// applicable.
			// This shouldn't happen in theory since it was tested above when the metric cache got initialized. But when
			// there is new evidence which has set the resource types and tool id correctly (their combination builds
			// the key for the cache), all metrics are applied due to the cache - even when all corresponding resource
			// fields are not set properly.
			if runMap != nil {
				data = append(data, runMap)
			}
		}
	}

	return data, nil
}

// HandleMetricEvent takes care of handling metric events, such as evicting cache entries for the
// appropriate metrics.
func (re *regoEval) HandleMetricEvent(event *orchestrator.MetricChangeEvent) (err error) {
	if event.Type == orchestrator.MetricChangeEvent_TYPE_IMPLEMENTATION_CHANGED {
		log.Infof("Implementation of %s has changed. Clearing cache for this metric", event.MetricId)
	} else if event.Type == orchestrator.MetricChangeEvent_TYPE_CONFIG_CHANGED {
		log.Infof("Configuration of %s has changed. Clearing cache for this metric", event.MetricId)
	}

	// Evict the cache for the given metric
	re.qc.Evict(event.MetricId)

	return nil
}

func (re *regoEval) evalMap(baseDir string, targetID string, metric *assessment.Metric, m map[string]interface{}, src MetricsSource) (result *CombinedResult, err error) {
	var (
		query  *rego.PreparedEvalQuery
		key    string
		pkg    string
		prefix string
	)

	// We need to check, if the metric configuration has been changed.
	config, err := src.MetricConfiguration(targetID, metric)
	if err != nil {
		return nil, fmt.Errorf("could not fetch metric configuration for metric %s: %w", metric.Id, err)
	}

	// We build a key out of the metric and its configuration, so we are creating a new Rego implementation
	// if the metric configuration (i.e. its hash) for a particular target of evaluation has changed.
	key = fmt.Sprintf("%s-%s-%s", metric.Id, targetID, config.Hash())

	// Try to fetch a cached prepared query for the specified key. If the key is not found, we create a new query with
	// the function specified as the second parameter
	query, err = re.qc.Get(key, func(key string) (*rego.PreparedEvalQuery, error) {
		var (
			tx   storage.Transaction
			impl *assessment.MetricImplementation
		)

		// Create paths for bundle directory and utility functions file
		bundle := fmt.Sprintf("%s/policies/bundles/%s/%s/", baseDir, metric.CategoryID(), metric.Id)
		operators := fmt.Sprintf("%s/policies/operators.rego", baseDir)

		// The contents of the data map is available as the data variable within the Rego evaluation
		data := map[string]interface{}{
			"target_value": config.TargetValue.AsInterface(),
			"operator":     config.Operator,
			"config":       config,
		}

		// Create a new in-memory Rego store based on our data map
		store := inmem.NewFromObject(data)
		ctx := context.Background()

		// Create a new transaction in the store
		tx, err = store.NewTransaction(ctx, storage.WriteParams)
		if err != nil {
			return nil, fmt.Errorf("could not create transaction: %w", err)
		}

		prefix = re.pkg

		// Convert camelCase metric in under_score_style for package name
		pkg = util.CamelCaseToSnakeCase(metric.Id)

		// Fetch the metric implementation, i.e., the Rego code from the metric source
		impl, err = src.MetricImplementation(assessment.MetricImplementation_LANGUAGE_REGO, metric)
		if err != nil {
			return nil, fmt.Errorf("could not fetch policy for metric %s: %w", metric.Id, err)
		}

		// Insert/Update the policy. The bundle path depends on the metric ID
		err = store.UpsertPolicy(context.Background(), tx, bundle+"metric.rego", []byte(impl.Code))
		if err != nil {
			return nil, fmt.Errorf("could not upsert policy: %w", err)
		}

		// Create a new Rego prepared query evaluation, which can later be used to query the metric on any object (input)
		query, err := rego.New(
			rego.Query(fmt.Sprintf(`
			output = data.%s.%s;
			applicable = data.%s.%s.applicable;
			compliant = data.%s.%s.compliant;
			operator = data.clouditor.operator;
			target_value = data.clouditor.target_value;
			config = data.clouditor.config`, prefix, pkg, prefix, pkg, prefix, pkg)),
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
			return nil, fmt.Errorf("could not prepare rego evaluation for metric %s: %w", metric.Id, err)
		}

		// Commit the transaction into the store
		err = store.Commit(ctx, tx)
		if err != nil {
			return nil, fmt.Errorf("could not commit transaction: %w", err)
		}

		return &query, nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not fetch cached query for metric %s: %w", metric.Id, err)
	}

	results, err := query.Eval(context.Background(), rego.EvalInput(m))
	if err != nil {
		return nil, fmt.Errorf("could not evaluate rego policy: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results. probably the package name of metric %s is wrong", metric.Id)
	}

	result = &CombinedResult{
		Applicable:  results[0].Bindings["applicable"].(bool),
		Compliant:   results[0].Bindings["compliant"].(bool),
		Operator:    results[0].Bindings["operator"].(string),
		TargetValue: results[0].Bindings["target_value"],
		MetricID:    metric.Id,
	}

	// A little trick to convert the map-based metric configuration back to a real object
	result.Config = new(assessment.MetricConfiguration)
	if err = reencode(results[0].Bindings["config"], result.Config); err != nil {
		return nil, err
	}

	// Enable the new results
	output := results[0].Bindings["output"]
	if results, ok := output.(map[string]interface{})["results"]; ok {
		result.ComparisonResult = make([]*assessment.ComparisonResult, 0)
		if err = reencode(results, &result.ComparisonResult); err != nil {
			return nil, err
		}
	}

	// Check, if the metric supplies an additional message
	if msg, ok := output.(map[string]interface{})["message"]; ok {
		// Also append a short comment that details can be found in the ... details, if we have any
		if len(result.ComparisonResult) > 0 {
			result.Message = fmt.Sprintf("%s %s", msg, assessment.AdditionalDetailsMessage)
		} else {
			result.Message = assessment.AdditionalDetailsMessage
		}
	} else if result.Compliant {
		result.Message = assessment.DefaultCompliantMessage
	} else if !result.Compliant {
		result.Message = assessment.DefaultNonCompliantMessage
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

func reencode[T any](in any, out *T) (err error) {
	var b []byte
	if b, err = json.Marshal(in); err != nil {
		return fmt.Errorf("JSON marshal failed: %w", err)
	}

	if err = json.Unmarshal(b, out); err != nil {
		return fmt.Errorf("JSON unmarshal failed: %w", err)
	}

	return
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

func idsOf(metrics []*assessment.Metric) (ids []string) {
	for _, metric := range metrics {
		ids = append(ids, metric.Id)
	}
	return
}
