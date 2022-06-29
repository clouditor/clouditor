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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFieldNames(t *testing.T) {
	var (
		fieldnames []string
		err        error
	)
	type someStruct struct {
		Name string
		// TODO(lebogg): When unexported, linter throws err 'field `secret` is unused (unused)'
		Secret int
	}
	// Avoid linter msg: "type `someStruct` is unused (unused)"
	_ = someStruct{}

	// Type has to be struct, not string
	_, err = GetFieldNames[string]()
	assert.ErrorIs(t, err, ErrNoStruct)

	// Type has to be struct, not int
	_, err = GetFieldNames[int]()
	assert.ErrorIs(t, err, ErrNoStruct)

	// Successful
	fieldnames, err = GetFieldNames[someStruct]()
	assert.NoError(t, err)
	assert.Equal(t, fieldnames, []string{"Name", "Secret"})

}
