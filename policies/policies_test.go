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
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"clouditor.io/clouditor/internal/testutil/clitest"
	"google.golang.org/protobuf/encoding/protojson"

	"clouditor.io/clouditor/api/assessment"
	"github.com/stretchr/testify/assert"
)

const (
	mockObjStorage1EvidenceID = "1"
	mockObjStorage1ResourceID = "/mockresources/storages/object1"
	mockObjStorage2EvidenceID = "2"
	mockObjStorage2ResourceID = "/mockresources/storages/object2"
	mockVM1EvidenceID         = "3"
	mockVM1ResourceID         = "/mockresources/compute/vm1"
	mockVM2EvidenceID         = "4"
	mockVM2ResourceID         = "/mockresources/compute/vm2"
)

func TestMain(m *testing.M) {
	clitest.AutoChdir()

	os.Exit(m.Run())
}

type mockMetricsSource struct {
	t *testing.T
}

func (m *mockMetricsSource) Metrics() (metrics []*assessment.Metric, err error) {
	var (
		b           []byte
		metricsFile = "service/orchestrator/metrics.json"
	)

	b, err = os.ReadFile(metricsFile)
	if err != nil {
		return nil, fmt.Errorf("error while loading %s: %w", metricsFile, err)
	}

	err = json.Unmarshal(b, &metrics)
	if err != nil {
		return nil, fmt.Errorf("error in JSON marshal: %w", err)
	}

	return
}

func (m *mockMetricsSource) MetricConfiguration(serviceId, metricId string) (*assessment.MetricConfiguration, error) {
	// Fetch the metric configuration directly from our file
	bundle := fmt.Sprintf("policies/bundles/%s/data.json", metricId)

	b, err := os.ReadFile(bundle)
	assert.NoError(m.t, err)

	var config assessment.MetricConfiguration
	err = protojson.Unmarshal(b, &config)
	assert.NoError(m.t, err)

	return &config, nil
}

func (m *mockMetricsSource) MetricImplementation(lang assessment.MetricImplementation_Language, metric string) (*assessment.MetricImplementation, error) {
	// Fetch the metric implementation directly from our file
	bundle := fmt.Sprintf("policies/bundles/%s/metric.rego", metric)

	b, err := os.ReadFile(bundle)
	assert.NoError(m.t, err)

	var impl = &assessment.MetricImplementation{
		MetricId: metric,
		Lang:     assessment.MetricImplementation_REGO,
		Code:     string(b),
	}

	return impl, nil
}
