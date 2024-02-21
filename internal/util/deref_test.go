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

	"clouditor.io/clouditor/v2/internal/testutil/assert"
)

type myStruct struct {
	Test string
}

func TestDeref(t *testing.T) {
	var testValue string
	assert.Equal(t, testValue, Deref(&testValue))

	testValue = "testString"
	assert.Equal(t, testValue, Deref(&testValue))

	var testInt32 int32 = 12
	assert.Equal(t, testInt32, Deref(&testInt32))

	var testInt64 int64 = 12
	assert.Equal(t, testInt64, Deref(&testInt64))

	var testFloat32 float32 = 1.5
	assert.Equal(t, testFloat32, Deref(&testFloat32))

	var testFloat64 float32 = 1.5
	assert.Equal(t, testFloat64, Deref(&testFloat64))

	var testBool = true
	assert.Equal(t, testBool, Deref(&testBool))

	testStruct := myStruct{
		Test: "test",
	}
	assert.Equal(t, testStruct, Deref(&testStruct))

	testByteArray := []byte("testByteArray")
	assert.Equal(t, testByteArray, Deref(&testByteArray))
}

func TestRef(t *testing.T) {
	var testValue string
	assert.Equal(t, &testValue, Ref(testValue))

	testValue = "testString"
	assert.Equal(t, &testValue, Ref(testValue))

	var testInt32 int32 = 12
	assert.Equal(t, &testInt32, Ref(testInt32))

	var testInt64 int64 = 12
	assert.Equal(t, &testInt64, Ref(testInt64))

	var testFloat32 float32 = 1.5
	assert.Equal(t, &testFloat32, Ref(testFloat32))

	var testFloat64 float32 = 1.5
	assert.Equal(t, &testFloat64, Ref(testFloat64))

	var testBool = true
	assert.Equal(t, &testBool, Ref(testBool))

	testStruct := myStruct{
		Test: "test",
	}
	assert.Equal(t, &testStruct, Ref(testStruct))

	testByteArray := []byte("testByteArray")
	assert.Equal(t, &testByteArray, Ref(testByteArray))
}
