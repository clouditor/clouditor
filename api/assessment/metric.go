// Copyright 2022 Fraunhofer AISEC
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
	"encoding/base64"
	"errors"
	"fmt"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	ErrMetricNameMissing             = errors.New("metric name is missing")
	ErrMetricEmpty                   = errors.New("metric is missing or empty")
	ErrTargetOfEvaluationIDIsMissing = errors.New("target of evaluation id is missing")
	ErrTargetOfEvaluationIDIsInvalid = errors.New("target of evaluation id is invalid")
)

// Hash provides a simple string based hash for this metric configuration. It can be used
// to provide a key for a map or a cache.
func (x *MetricConfiguration) Hash() string {
	return base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf("%v-%v", x.Operator, x.TargetValue)))
}

func (x *MetricConfiguration) MarshalJSON() (b []byte, err error) {
	return protojson.Marshal(x)
}

func (x *MetricConfiguration) UnmarshalJSON(b []byte) (err error) {
	return protojson.Unmarshal(b, x)
}
