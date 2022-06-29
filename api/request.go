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
	"errors"

	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/proto"

	"clouditor.io/clouditor/internal/util"
)

var (
	ErrInvalidColumnName = errors.New("column name is invalid")
	ErrRequestIsNil      = errors.New("request is empty")
)

// ListRequest indicates a proto message being a request for ListXXX RPCs
type ListRequest interface {
	GetOrderBy() string
	GetAsc() bool
	proto.Message
}

func ValidateListReq(req ListRequest, responseType any) (err error) {
	// req must be non-nil
	if req == nil {
		err = ErrRequestIsNil
		return
	}

	// Avoid DB injections by whitelisting the valid orderBy statements
	whitelist, err := util.GetFieldNames(responseType)
	// Add empty string indicating no explicit ordering
	whitelist = append(whitelist, "")
	if !slices.Contains(whitelist, req.GetOrderBy()) {
		err = ErrInvalidColumnName
		return
	}

	return

}
