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
	"unicode"

	"clouditor.io/clouditor/v2/api/ontology"
)

// CamelCaseToSnakeCase converts a `camelCase` string to `snake_case`
func CamelCaseToSnakeCase(input string) string {
	if input == "" {
		return ""
	}

	snakeCase := make([]rune, 0, len(input))
	runeArray := []rune(input)

	for i := range runeArray {
		if i > 0 && marksNewWord(i, runeArray) {
			snakeCase = append(snakeCase, '_', unicode.ToLower(runeArray[i]))
		} else {
			snakeCase = append(snakeCase, unicode.ToLower(runeArray[i]))
		}
	}

	return string(snakeCase)
}

// marksNewWord checks if the current character starts a new word excluding the first word
func marksNewWord(i int, input []rune) bool {

	if i >= len(input) {
		return false
	}

	// If previous or following rune/character is lowercase or a number then it is a new word
	if i < len(input)-1 && unicode.IsUpper(input[i]) && unicode.IsLower(input[i+1]) {
		return true
	} else if i > 0 && unicode.IsLower(input[i-1]) && unicode.IsUpper(input[i]) {
		return true
	} else if i < len(input)-1 && unicode.IsDigit(input[i]) && unicode.IsLower(input[i+1]) {
		return true
	} else if i > 0 && unicode.IsUpper(input[i-1]) && unicode.IsDigit(input[i]) {
		return true
	}

	return false
}

// ListResourceIDs return a list of the given resource IDs
func ListResourceIDs(r []ontology.IsResource) []string {
	var a = []string{}

	for _, v := range r {
		a = append(a, v.GetId())
	}

	return a
}
