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

package evidence

import "errors"

type EvidenceHookFunc func(result *Evidence, err error)

var (
	ErrNotValidResource    = errors.New("resource in evidence is not a map")
	ErrResourceNotStruct   = errors.New("resource in evidence is not struct value")
	ErrResourceNotMap      = errors.New("resource in evidence is not a map")
	ErrResourceIdMissing   = errors.New("resource in evidence is missing the id field")
	ErrResourceIdNotString = errors.New("resource id in evidence is not a string")
	ErrToolIdMissing       = errors.New("tool id in evidence is missing")
	ErrTimestampMissing    = errors.New("timestamp in evidence is missing")
)

// Validate validates the evidence according to several required fields
func (evidence *Evidence) Validate() (resourceId string, err error) {
	if evidence.Resource == nil {
		return "", ErrNotValidResource
	}

	value := evidence.Resource.GetStructValue()
	if value == nil {
		return "", ErrResourceNotStruct
	}

	m := evidence.Resource.GetStructValue().AsMap()
	if m == nil {
		return "", ErrResourceNotMap
	}

	field, ok := m["id"]
	if !ok {
		return "", ErrResourceIdMissing
	}

	resourceId, ok = field.(string)
	if !ok {
		return "", ErrResourceIdNotString
	}

	if evidence.ToolId == "" {
		return "", ErrToolIdMissing
	}

	if evidence.Timestamp == nil {
		return "", ErrTimestampMissing
	}

	return
}
