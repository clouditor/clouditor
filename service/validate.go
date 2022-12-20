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

package service

import (
	"reflect"
	"strings"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/internal/util"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// ValidateRequest validates an incoming request according to different criteria:
//   - If the request is nil, [api.ErrEmptyRequest] is returned
//   - The request is validated according to the generated validation method
//   - Lastly, if the request is a [api.PaginatedRequest], an additional check is performed to ensure only valid columns are listed
func ValidateRequest(req IncomingRequest) (err error) {
	// Check, if request is nil. We need to check whether the interface itself is nil (untyped nil); this happens if
	// someone is directly setting nil to a variable of the interface IncomingRequest. Furthermore, we need to check,
	// whether the *value* of the interface is nil. This can happen if nil is first assigned to a variable of a struct
	// (pointer) that implements the interface. If this variable is then passed to the validate function, the req
	// paramater is not nil, but the value of the interface representing it is.
	if req == nil || reflect.ValueOf(req).IsNil() {
		return status.Errorf(codes.InvalidArgument, "%s", api.ErrEmptyRequest)
	}

	// Validate request
	err = req.Validate()
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	// Check, if request is a list request
	if preq, ok := req.(api.PaginatedRequest); ok {
		whitelist := util.GetFieldNames(preq)
		// Add empty string indicating no explicit ordering
		whitelist = append(whitelist, "")

		normalizedReq := strings.ToLower(preq.GetOrderBy())
		if !slices.Contains(whitelist, normalizedReq) {
			return status.Errorf(codes.InvalidArgument, "invalid request: %v", api.ErrInvalidColumnName)
		}
	}

	return nil
}

type IncomingRequest interface {
	Validate() error
	proto.Message
}
