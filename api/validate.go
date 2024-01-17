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

package api

import (
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/internal/util"

	"github.com/bufbuild/protovalidate-go"
	"github.com/bufbuild/protovalidate-go/legacy"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Validate validates an incoming request according to different criteria:
//   - If the request is nil, [api.ErrEmptyRequest] is returned
//   - The request is validated according to the generated validation method
//   - Lastly, if the request is a [api.PaginatedRequest], an additional check is performed to ensure only valid columns are listed
//
// Note: This function already returns a gRPC error, so the error can be returned directly without any wrapping in a
// request function.
func Validate(req IncomingRequest) (err error) {
	// Check, if request is nil. We need to check whether the interface itself is nil (untyped nil); this happens if
	// someone is directly setting nil to a variable of the interface IncomingRequest. Furthermore, we need to check,
	// whether the *value* of the interface is nil. This can happen if nil is first assigned to a variable of a struct
	// (pointer) that implements the interface. If this variable is then passed to the validate function, the req
	// parameter is not nil, but the value of the interface representing it is.
	if util.IsNil(req) {
		return status.Errorf(codes.InvalidArgument, "%s", ErrEmptyRequest)
	}

	// TODO(oxisto): Re-use validator?
	v, err := protovalidate.New(
		legacy.WithLegacySupport(legacy.ModeMerge),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to initialize validator: %s", err)
	}

	// Validate request
	err = v.Validate(req)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "%v: %v", ErrInvalidRequest, err)
	}

	return nil
}

// ValidateWithResource validates the evidence according to its resource
// TODO(oxisto): Replace with CEL?
func ValidateWithResource(ev *evidence.Evidence) (resourceId string, err error) {
	err = Validate(ev)
	if err != nil {
		return "", err
	}

	value := ev.Resource.GetStructValue()
	if value == nil {
		return "", ErrResourceNotStruct
	}

	m := ev.Resource.GetStructValue().AsMap()
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

type IncomingRequest interface {
	proto.Message
}
