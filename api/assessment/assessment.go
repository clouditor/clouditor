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

package assessment

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type ResultHookFunc func(result *AssessmentResult, err error)

var (
	ErrTimestampMissing                      = errors.New("timestamp in assessment result is missing")
	ErrMetricIdMissing                       = errors.New("metric id in assessment result is missing")
	ErrMetricConfigurationMissing            = errors.New("metric configuration in assessment result is missing")
	ErrEvidenceIdMissing                     = errors.New("evidence id in assessment result is missing")
	ErrMetricConfigurationOperatorMissing    = errors.New("operator in metric data is missing")
	ErrMetricConfigurationTargetValueMissing = errors.New("target value in metric data is missing")
)

// Validate validates the assessment result according to several required fields
func (result *AssessmentResult) Validate() (resourceId string, err error) {
	if _, err = uuid.Parse(result.Id); err != nil {
		fmt.Println(uuid.IsInvalidLengthError(err))
		return "", err
	}

	if result.Timestamp == nil {
		return "", ErrTimestampMissing
	}

	if result.MetricId == "" {
		return "", ErrMetricIdMissing
	}

	if result.MetricConfiguration == nil {
		return "", ErrMetricConfigurationMissing
	}

	if result.MetricConfiguration.Operator == "" {
		return "", ErrMetricConfigurationOperatorMissing
	}

	if result.MetricConfiguration.TargetValue == nil {
		return "", ErrMetricConfigurationTargetValueMissing
	}

	if result.EvidenceId == "" {
		return "", ErrEvidenceIdMissing
	}

	// TODO(garuppel): Missing resource ID?

	return
}
