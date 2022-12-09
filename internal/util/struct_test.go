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
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"clouditor.io/clouditor/api/auth"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/internal/testutil/prototest"
)

// Not the most sophisticated tests. They will probably fail at some point when these proto messages change. Mocking
// a proto.Message (and thus implementing ProtoReflect() as well) would be ideal...
func TestGetFieldNames(t *testing.T) {
	var (
		fieldnames []string
	)

	// Successful
	fieldnames = GetFieldNames(&auth.User{})
	assert.Equal(t, []string{"username", "password", "email", "full_name", "shadow"}, fieldnames)
	assert.Equal(t, reflect.ValueOf(auth.User{}).NumField()-3, len(fieldnames))

	// Successful
	fieldnames = GetFieldNames(&prototest.TestStruct{})
	assert.Equal(t, reflect.ValueOf(prototest.TestStruct{}).NumField()-3, len(fieldnames))
	assert.Equal(t, []string{"test_name", "test_id", "test_description", "test_status"}, fieldnames)

	// Successful
	fieldnames = GetFieldNames(&evidence.Evidence{})
	assert.Equal(t, reflect.ValueOf(evidence.Evidence{}).NumField()-3, len(fieldnames))
	assert.Equal(t, []string{"id", "timestamp", "cloud_service_id", "tool_id", "raw", "resource"}, fieldnames)
}
