//
// Copyright 2016-2023 Fraunhofer AISEC
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

package evidence

import "google.golang.org/protobuf/proto"

// GetTargetOfEvaluationId is a shortcut to implement TargetOfEvaluationRequest. It returns
// the target of evaluation ID of the inner object.
func (req *StoreEvidenceRequest) GetTargetOfEvaluationId() string {
	return req.GetEvidence().GetTargetOfEvaluationId()
}

// GetPayload is a shortcut to implement EvidenceRequest. It returns the Evidence of the request.
func (req *StoreEvidenceRequest) GetPayload() proto.Message {
	return req.Evidence
}

// GetTargetOfEvaluationId is a shortcut to implement TargetOfEvaluationRequest. It returns
// the target of evaluation ID of the inner object.
func (req *UpdateResourceRequest) GetTargetOfEvaluationId() string {
	return req.Resource.GetTargetOfEvaluationId()
}
