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

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/api/orchestrator"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithField("component", "policies")
)

// metricsCache holds all cached metrics for different combinations of Tools with resource types
type metricsCache struct {
	sync.RWMutex
	// Metrics cached in a map. Key is composed of tool id and resource types concatenation
	m map[string][]string
}

// TODO(oxisto): Rename to AssessmentEngine or something?
type PolicyEval interface {
	Eval(evidence *evidence.Evidence, src MetricsSource) (data []*Result, err error)
	HandleMetricEvent(event *orchestrator.MetricChangeEvent) (err error)
}

type Result struct {
	Applicable  bool
	Compliant   bool
	TargetValue interface{} `mapstructure:"target_value"`
	Operator    string
	MetricId    string
}

// MetricsSource is used to retrieve a list of metrics and to retrieve a metric
// configuration as well as implementation for a particular metric (and target service)
type MetricsSource interface {
	Metrics() ([]*assessment.Metric, error)
	MetricConfiguration(serviceId, metricId string) (*assessment.MetricConfiguration, error)
	MetricImplementation(lang assessment.MetricImplementation_Language, metric string) (*assessment.MetricImplementation, error)
}

// RequirementsSource is used to retrieve a list of requirements
type RequirementsSource interface {
	Requirements() ([]*orchestrator.Requirement, error)
}

// createKey creates a key by concatenating toolID and all types
func createKey(evidence *evidence.Evidence, types []string) (key string) {
	// Merge toolID and types to one slice and concatenate all its elements
	key = strings.Join(append(types, evidence.ToolId), "-")
	key = strings.ReplaceAll(key, " ", "")
	return
}
