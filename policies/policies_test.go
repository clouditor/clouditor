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
	"time"

	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/clitest"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	mockBlockStorage1ID       = "/mockresources/storage/storage1"
	mockBlockStorage2ID       = "/mockresources/storage/storage2"
)

func TestMain(m *testing.M) {
	clitest.AutoChdir()

	os.Exit(m.Run())
}

type mockMetricsSource struct {
	t *testing.T
}

func (*mockMetricsSource) Metrics() (metrics []*assessment.Metric, err error) {
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
	for _, m := range metrics {
		if m.Scale == assessment.Metric_SCALE_UNSPECIFIED {
			return nil, fmt.Errorf("no scale specified for metric '%s'", m.Name)
		}
		if m.Range == nil {
			return nil, fmt.Errorf("no range specified for metric '%s'", m.Name)
		}
	}

	return
}

func (m *mockMetricsSource) MetricConfiguration(targetID string, metric *assessment.Metric) (*assessment.MetricConfiguration, error) {
	// Fetch the metric configuration directly from our file
	bundle := fmt.Sprintf("policies/metrics/metrics/%s/%s/data.json", metric.CategoryID(), metric.Id)

	b, err := os.ReadFile(bundle)
	assert.NoError(m.t, err)

	var config assessment.MetricConfiguration
	err = protojson.Unmarshal(b, &config)
	assert.NoError(m.t, err)

	config.IsDefault = true
	config.MetricId = metric.Id
	config.CertificationTargetId = targetID

	return &config, nil
}

func (m *mockMetricsSource) MetricImplementation(_ assessment.MetricImplementation_Language, metric *assessment.Metric) (*assessment.MetricImplementation, error) {
	// Fetch the metric implementation directly from our file
	bundle := fmt.Sprintf("policies/metrics/metrics/%s/%s/metric.rego", metric.CategoryID(), metric.Id)

	b, err := os.ReadFile(bundle)
	assert.NoError(m.t, err)

	var impl = &assessment.MetricImplementation{
		MetricId: metric.Id,
		Lang:     assessment.MetricImplementation_LANGUAGE_REGO,
		Code:     string(b),
	}

	return impl, nil
}

type updatedMockMetricsSource struct {
	mockMetricsSource
}

func (*updatedMockMetricsSource) MetricConfiguration(targetID string, metric *assessment.Metric) (*assessment.MetricConfiguration, error) {
	return &assessment.MetricConfiguration{
		Operator:              "==",
		TargetValue:           structpb.NewBoolValue(false),
		IsDefault:             false,
		UpdatedAt:             timestamppb.New(time.Date(2022, 12, 1, 0, 0, 0, 0, time.Local)),
		MetricId:              metric.Id,
		CertificationTargetId: targetID,
	}, nil
}
