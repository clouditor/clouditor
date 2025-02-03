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
	"strings"
	"sync"

	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithField("component", "policies")
)

// metricsCache holds all cached metrics for different combinations of Tools with resource types
type metricsCache struct {
	sync.RWMutex
	// Metrics cached in a map. Key is composed of tool id and resource types concatenation
	m map[string][]*assessment.Metric
}

// PolicyEval is an interface for the policy evaluation engine
type PolicyEval interface {
	// Eval evaluates a given evidence against a metric coming from the metrics source. In order to avoid unnecessarily
	// unwrapping, the callee of this function needs to supply the unwrapped ontology resource, since they most likely
	// unwrapped the resource already, e.g. to check for validation.
	Eval(evidence *evidence.Evidence, r ontology.IsResource, related map[string]ontology.IsResource, src MetricsSource) (data []*CombinedResult, err error)
	HandleMetricEvent(event *orchestrator.MetricChangeEvent) (err error)
}

type CombinedResult struct {
	Applicable bool
	Compliant  bool
	// TODO(oxisto): They are now part of the individual comparison results
	TargetValue interface{}
	// TODO(oxisto): They are now part of the individual comparison results
	Operator string
	MetricID string
	Config   *assessment.MetricConfiguration

	// ComparisonResult is an optional feature to get more infos about the comparisons
	ComparisonResult []*assessment.ComparisonResult

	// Message contains an optional string that the metric can supply to provide a human readable representation of the result
	Message string
}

// MetricsSource is used to retrieve a list of metrics and to retrieve a metric
// configuration as well as implementation for a particular metric (and certification target)
type MetricsSource interface {
	Metrics() ([]*assessment.Metric, error)
	MetricConfiguration(targetID string, metric *assessment.Metric) (*assessment.MetricConfiguration, error)
	MetricImplementation(lang assessment.MetricImplementation_Language, metric *assessment.Metric) (*assessment.MetricImplementation, error)
}

// ControlsSource is used to retrieve a list of controls
type ControlsSource interface {
	Controls() ([]*orchestrator.Control, error)
}

// createKey creates a key by concatenating toolID and all types
func createKey(evidence *evidence.Evidence, types []string) (key string) {
	// Merge toolID and types to one slice and concatenate all its elements
	key = strings.Join(append(types, evidence.ToolId), "-")
	key = strings.ReplaceAll(key, " ", "")
	return
}
