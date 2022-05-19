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

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/structpb"
)

type EvidenceHookFunc func(result *Evidence, err error)

var (
	ErrEvidenceIdInvalidFormat = errors.New("evidence id not in expected format (UUID) or missing")
	ErrNotValidResource        = errors.New("resource in evidence is missing")
	ErrResourceNotStruct       = errors.New("resource in evidence is not struct value")
	ErrResourceNotMap          = errors.New("resource in evidence is not a map")
	ErrResourceIdMissing       = errors.New("resource in evidence is missing the id field")
	ErrResourceIdNotString     = errors.New("resource id in evidence is not a string")
	ErrToolIdMissing           = errors.New("tool id in evidence is missing")
	ErrTimestampMissing        = errors.New("timestamp in evidence is missing")
	ErrResourceIdFieldMissing  = errors.New("field id is missing")
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
		return "", ErrResourceIdFieldMissing
	} else if field == "" {
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

	if _, err = uuid.Parse(evidence.Id); err != nil {
		return "", ErrEvidenceIdInvalidFormat
	}

	return
}

// ResourceTypes parses the embedded resource of this evidence and returns its types according to the ontology.
func (evidence *Evidence) ResourceTypes() (types []string, err error) {
	var (
		m     map[string]interface{}
		value *structpb.Value
	)

	value = evidence.Resource
	if value == nil {
		return
	}

	m = value.GetStructValue().AsMap()

	if rawTypes, ok := m["type"].([]interface{}); ok {
		if len(rawTypes) != 0 {
			types = make([]string, len(rawTypes))
		} else {
			return nil, fmt.Errorf("list of types is empty")
		}
	} else {
		return nil, fmt.Errorf("got type '%T' but wanted '[]interface {}'. Check if resource types are specified ", rawTypes)
	}
	for i, v := range m["type"].([]interface{}) {
		if t, ok := v.(string); !ok {
			return nil, fmt.Errorf("got type '%T' but wanted 'string'", t)
		} else {
			types[i] = t
		}
	}

	return
}
