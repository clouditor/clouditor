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
	"context"
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type EvidenceHookFunc func(ctx context.Context, evidence *Evidence, err error)

var (
	ErrResourceNotStruct             = errors.New("resource in evidence is not struct value")
	ErrResourceNotMap                = errors.New("resource in evidence is not a map")
	ErrResourceIdIsEmpty             = errors.New("resource id in evidence is empty")
	ErrResourceIdNotString           = errors.New("resource id in evidence is not a string")
	ErrResourceIdFieldMissing        = errors.New("resource in evidence is missing the id field")
	ErrResourceTypeFieldMissing      = errors.New("field type in evidence is missing")
	ErrResourceTypeNotArrayOfStrings = errors.New("resource type in evidence is not an array of strings")
	ErrResourceTypeEmpty             = errors.New("resource type (array) in evidence is empty")
)

// ValidateWithResource validates the evidence according to its resource
func (evidence *Evidence) ValidateWithResource() (resourceId string, err error) {
	err = evidence.Validate()
	if err != nil {
		return "", err
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
		return "", ErrResourceIdIsEmpty
	}

	resourceId, ok = field.(string)
	if !ok {
		return "", ErrResourceIdNotString
	}

	_, ok = m["type"]
	if !ok {
		return "", ErrResourceTypeFieldMissing
	}

	// Check if resource is a slice
	fieldType, ok := m["type"].([]interface{})
	if !ok {
		// Resource is not a slice
		return "", ErrResourceTypeNotArrayOfStrings
	} else if len(fieldType) == 0 {
		// Resource slice is empty
		return "", ErrResourceTypeEmpty
	} else {
		if _, ok := fieldType[0].(string); !ok {
			// Resource slice does not contain string values
			return "", ErrResourceTypeNotArrayOfStrings
		}
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

func (req *StoreEvidenceRequest) GetPayload() proto.Message {
	return req.Evidence
}
