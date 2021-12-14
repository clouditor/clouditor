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

import "errors"

var (
	ErrTimestampMissing             = errors.New("timestamp in assessment result is missing")
	ErrMetricIdMissing              = errors.New("metric id in assessment result is missing")
	ErrMetricDataMissing            = errors.New("metric data in assessment result is missing")
	ErrEvidenceIdMissing            = errors.New("evidence id in assessment result is missing")
	ErrNonComplianceCommentsMissing = errors.New("non-compliance comments in assessment result is missing")
	ErrMetricDataOperatorMissing    = errors.New("operator in metric data is missing")
	ErrMetricDataTargetValueMissing = errors.New("target value in metric data is missing")
)

// Validate validates the assessment result according to several required fields
func (result *Result) Validate() (resourceId string, err error) {
	if result.Timestamp == nil {
		return "", ErrTimestampMissing
	}

	if result.MetricId == "" {
		return "", ErrMetricIdMissing
	}

	if result.MetricData == nil {
		return "", ErrMetricDataMissing
	}

	if result.MetricData.Operator == "" {
		return "", ErrMetricDataOperatorMissing
	}

	if result.MetricData.TargetValue == nil {
		return "", ErrMetricDataTargetValueMissing
	}

	// TODO(all): Do we have to check the target value type?
	//value := result.MetricData.GetTargetValue()
	//if value == nil {
	//	return "", ErrResourceNotStruct
	//}

	if result.EvidenceId == "" {
		return "", ErrEvidenceIdMissing
	}

	if result.NonComplianceComments == "" {
		return "", ErrNonComplianceCommentsMissing
	}

	return
}
