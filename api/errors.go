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
)

var (
	ErrInvalidColumnName             = errors.New("column name is invalid")
	ErrEmptyRequest                  = errors.New("empty request")
	ErrInvalidRequest                = errors.New("invalid request")
	ErrResourceNotStruct             = errors.New("resource in evidence is not struct value")
	ErrResourceNotMap                = errors.New("resource in evidence is not a map")
	ErrResourceIdIsEmpty             = errors.New("resource id in evidence is empty")
	ErrResourceIdNotString           = errors.New("resource id in evidence is not a string")
	ErrResourceIdFieldMissing        = errors.New("resource in evidence is missing the id field")
	ErrResourceTypeFieldMissing      = errors.New("field type in evidence is missing")
	ErrResourceTypeNotArrayOfStrings = errors.New("resource type in evidence is not an array of strings")
	ErrResourceTypeEmpty             = errors.New("resource type (array) in evidence is empty")
)
