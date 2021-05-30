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

package ccl_test

import (
	"encoding/json"
	"testing"

	"clouditor.io/clouditor/ccl"
	"github.com/stretchr/testify/assert"
)

func TestRuleFromFile(t *testing.T) {
	var (
		err     error
		success bool
		j       = `{"field": "value"}`
		o       map[string]interface{}
	)

	err = json.Unmarshal([]byte(j), &o)

	assert.Nil(t, err)

	success, err = ccl.RunRuleFromFile("test_rule.ccl", o)

	assert.Nil(t, err)
	assert.True(t, success)
}

func TestRunRule_intCompare(t *testing.T) {
	var (
		err     error
		success bool
		j       = `{"intField": 1}`
		o       map[string]interface{}
	)

	err = json.Unmarshal([]byte(j), &o)

	assert.Nil(t, err)

	success, err = ccl.RunRule("Object has intField == 1", o)

	assert.Nil(t, err)
	assert.True(t, success)
}

func TestRunRule_floatCompare(t *testing.T) {
	var (
		err     error
		success bool
		j       = `{"floatField": 1.5}`
		o       map[string]interface{}
	)

	err = json.Unmarshal([]byte(j), &o)

	assert.Nil(t, err)

	success, err = ccl.RunRule("Object has floatField == 1.5", o)

	assert.Nil(t, err)
	assert.True(t, success)
}

func TestRunRule_boolCompare(t *testing.T) {
	var (
		err     error
		success bool
		j       = `{"boolField": false}`
		o       map[string]interface{}
	)

	err = json.Unmarshal([]byte(j), &o)

	assert.Nil(t, err)

	success, err = ccl.RunRule("Object has boolField == false", o)

	assert.Nil(t, err)
	assert.True(t, success)
}
