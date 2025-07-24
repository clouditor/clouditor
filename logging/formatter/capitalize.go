// Copyright 2016-2022 Fraunhofer AISEC
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

package formatter

import (
	"unicode"

	"github.com/sirupsen/logrus"
)

// CapitalizeFormatter implements the Formatter interface for customizing the logging to our needs
type CapitalizeFormatter struct {
	logrus.Formatter
}

// Format capitalizes the first letter of entry's message
func (f CapitalizeFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// If message is empty, nothing to format. Prevents `out of range` error
	if entry.Message != "" {
		// Convert message to runes, capitalize first rune and convert back to string again
		msgAsRune := []rune(entry.Message)
		msgAsRune[0] = unicode.ToUpper(msgAsRune[0])
		entry.Message = string(msgAsRune)
	}

	return f.Formatter.Format(entry)
}
