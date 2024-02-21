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
	"testing"

	"clouditor.io/clouditor/v2/internal/testutil/assert"

	"github.com/sirupsen/logrus"
)

func TestCapitalizeFormatter_Format(t *testing.T) {
	type fields struct {
		Formatter logrus.Formatter
	}
	type args struct {
		entry *logrus.Entry
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantMessage string
	}{
		{
			name: "First letter is already capitalized",
			fields: fields{
				Formatter: &logrus.TextFormatter{},
			},
			args: args{
				entry: &logrus.Entry{Message: "Capitalized"},
			},
			wantMessage: "Capitalized",
		},
		{
			name: "First letter is not capitalized",
			fields: fields{
				Formatter: &logrus.TextFormatter{},
			},
			args: args{
				entry: &logrus.Entry{Message: "capitalized"},
			},
			wantMessage: "Capitalized",
		},
		{
			name: "First letter is apostrophe",
			fields: fields{
				Formatter: &logrus.TextFormatter{},
			},
			args: args{
				entry: &logrus.Entry{Message: "'capitalized'"},
			},
			wantMessage: "'capitalized'",
		},
		{
			name: "Message is empty",
			fields: fields{
				Formatter: &logrus.TextFormatter{},
			},
			args: args{
				entry: &logrus.Entry{},
			},
			wantMessage: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := CapitalizeFormatter{
				Formatter: tt.fields.Formatter,
			}
			// Format (incl. TextFormatter's Format) doesn't throw any error
			gotEntry, _ := f.Format(tt.args.entry)

			formatter := logrus.TextFormatter{}

			// Format wantedMessage with TextFormatter and test against gotEntry.
			// TextFormatter's Format doesn't throw any error
			wantEntry, _ := formatter.Format(&logrus.Entry{Message: tt.wantMessage})
			assert.Equal(t, wantEntry, gotEntry)
		})
	}
}
