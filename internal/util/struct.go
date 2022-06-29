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

package util

import (
	"errors"
	"reflect"
)

var (
	// ErrNoStruct indicates the passed argument is not a struct
	ErrNoStruct = errors.New("no struct")
)

// GetFieldNames extracts all field names of struct T. Returns error if T is no struct.
// TODO(lebogg): Only take exported fields
func GetFieldNames[T any]() (fieldNames []string, err error) {
	var aStruct T
	// Check aStruct is a struct
	if reflect.TypeOf(aStruct).Kind() != reflect.Struct {
		err = ErrNoStruct
		return
	}
	// Get all fields of aStruct and add their names to fieldNames
	t := reflect.TypeOf(aStruct)
	fields := reflect.VisibleFields(t)
	for _, f := range fields {
		fieldNames = append(fieldNames, f.Name)
	}
	return
}
