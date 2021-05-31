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
	"time"

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

// just a shortcut
type m map[string]interface{}

func TestRunRule(t *testing.T) {
	var err error

	var testData = []struct {
		name            string
		json            map[string]interface{}
		ccl             string
		expectedError   error
		expectedSuccess bool
	}{
		{
			name:            "string compare",
			json:            m{"stringField": 1},
			ccl:             "Object has stringField == 1",
			expectedSuccess: true,
		},
		{
			name:            "int compare",
			json:            m{"intField": 1},
			ccl:             "Object has intField == 1",
			expectedSuccess: true,
		},
		{
			name:            "float compare",
			json:            m{"floatField": 1.5},
			ccl:             "Object has floatField == 1.5",
			expectedSuccess: true,
		},
		{
			name:            "bool compare",
			json:            m{"boolField": false},
			ccl:             "Object has boolField == false",
			expectedSuccess: true,
		},
		{
			name:            "not expression",
			json:            m{"boolField": false},
			ccl:             "Object has not boolField == false",
			expectedSuccess: false,
		},
		{
			name:            "isempty expression with nonexisting field",
			json:            m{},
			ccl:             "Object has empty field",
			expectedSuccess: true,
		},
		{
			name:            "isempty expression with empty string",
			json:            m{"field": ""},
			ccl:             "Object has empty field",
			expectedSuccess: true,
		},
		{
			name:            "isempty expression with boolean false",
			json:            m{"field": false},
			ccl:             "Object has empty field",
			expectedSuccess: true,
		},
		{
			name:            "isempty expression with integer zero",
			json:            m{"field": 0},
			ccl:             "Object has empty field",
			expectedSuccess: true,
		},
		{
			name:          "field does not exist",
			json:          m{},
			ccl:           "Object has (field == 1)",
			expectedError: ccl.ErrFieldNameNotFound,
		},
		{
			name:          "field not a map",
			json:          m{"nested": 1},
			ccl:           "Object has nested.field == 1",
			expectedError: ccl.ErrFieldNoMap,
		},
		{
			name:            "nested field",
			json:            m{"nested": map[string]interface{}{"field": 1}},
			ccl:             "Object has nested.field == 1",
			expectedSuccess: true,
		},
		{
			name:          "nested field does not exist",
			json:          m{"nested": map[string]interface{}{"field": 1}},
			ccl:           "Object has nested.someotherfield == 1",
			expectedError: ccl.ErrFieldNameNotFound,
		},
		{
			name:            "not equals operator int",
			json:            m{"field": 2},
			ccl:             "Object has field != 1",
			expectedSuccess: true,
		},
		{
			name:            "not equals operator string",
			json:            m{"field": "value"},
			ccl:             "Object has field != \"some_value\"",
			expectedSuccess: true,
		},
		{
			name:            "not equals operator float",
			json:            m{"field": 1.5},
			ccl:             "Object has field != 2",
			expectedSuccess: true,
		},
		{
			name:            "not equals operator bool",
			json:            m{"field": false},
			ccl:             "Object has field != true",
			expectedSuccess: true,
		},
		{
			name:            "less operator",
			json:            m{"field": 1},
			ccl:             "Object has field < 2",
			expectedSuccess: true,
		},
		{
			name:            "less equals operator",
			json:            m{"field": 2},
			ccl:             "Object has field <= 2",
			expectedSuccess: true,
		},
		{
			name:            "greater operator",
			json:            m{"field": 2},
			ccl:             "Object has field > 1",
			expectedSuccess: true,
		},
		{
			name:            "greater equals operator",
			json:            m{"field": 1},
			ccl:             "Object has field >= 1",
			expectedSuccess: true,
		},
		{
			name:            "contains operator",
			json:            m{"field": "myvalue"},
			ccl:             "Object has field contains \"my\"",
			expectedSuccess: true,
		},
		{
			name:            "within expression",
			json:            m{"field": 4},
			ccl:             "Object has field within 1,2,3,4",
			expectedSuccess: true,
		},
		{
			name:            "time comparison before",
			json:            m{"field": time.Now().Add(1 * time.Hour * 24)},
			ccl:             "Object has field before 2 days",
			expectedSuccess: true,
		},
		{
			name:            "time comparison after",
			json:            m{"field": time.Now().Add(10 * time.Second)},
			ccl:             "Object has field after 2 seconds",
			expectedSuccess: true,
		},
		{
			name:            "time comparison younger",
			json:            m{"field": time.Now().Add(-2 * time.Hour * 24 * 30)},
			ccl:             "Object has field younger 3 months",
			expectedSuccess: true,
		},
		{
			name:            "time comparison older",
			json:            m{"field": time.Now().Add(-10 * time.Hour * 24 * 30)},
			ccl:             "Object has field older now",
			expectedSuccess: true,
		},
		{
			name: "in any expression",
			json: m{"array": []map[string]interface{}{
				{"field": 1},
				{"field": 2},
			}},
			ccl:             "Object has field == 1 in any array",
			expectedSuccess: true,
		},
		{
			name: "in all expression",
			json: m{"array": []map[string]interface{}{
				{"field": 1},
				{"field": 2},
			}},
			ccl:             "Object has field == 1 in all array",
			expectedSuccess: false,
		},
		{
			name:          "syntax error",
			json:          m{},
			ccl:           "Object has nonsense",
			expectedError: ccl.ErrUnexpectedExpression,
		},
		{
			name:          "syntax error in any expression",
			json:          m{},
			ccl:           "Object has (field == in any array",
			expectedError: ccl.ErrUnexpectedExpression,
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			success, err := ccl.RunRule(data.ccl, data.json)

			assert.ErrorIs(t, err, data.expectedError, "did not match error")
			assert.Equal(t, data.expectedSuccess, success, "did not match expected success outcome")
		})
	}

	assert.Nil(t, err)
}
