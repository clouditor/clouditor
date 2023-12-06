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
	"context"
	"errors"

	"google.golang.org/protobuf/proto"
)

type ResultHookFunc func(ctx context.Context, result *AssessmentResult, err error)

var (
	ErrMetricConfigurationMissing            = errors.New("metric configuration in assessment result is missing")
	ErrMetricConfigurationOperatorMissing    = errors.New("operator in metric data is missing")
	ErrMetricConfigurationTargetValueMissing = errors.New("target value in metric data is missing")
)

const AssessmentToolId = "Clouditor Assessment"

func (req *AssessEvidenceRequest) GetPayload() proto.Message {
	return req.Evidence
}

// GetCloudServiceId is a shortcut to implement CloudServiceRequest. It returns the cloud service ID of the inner
// object.
func (req *AssessEvidenceRequest) GetCloudServiceId() string {
	return req.GetEvidence().GetCloudServiceId()
}
