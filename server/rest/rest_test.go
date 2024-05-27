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

package rest

import (
	"testing"
)

func Test_corsConfig_OriginAllowed(t *testing.T) {
	type fields struct {
		allowedOrigins []string
		allowedHeaders []string
		allowedMethods []string
	}
	type args struct {
		origin string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "Allow non-browser origin",
			fields: fields{},
			args: args{
				origin: "", // origin is only explicitly set by a browser
			},
			want: true,
		},
		{
			name: "Allowed origin",
			fields: fields{
				allowedOrigins: []string{"clouditor.io", "localhost"},
			},
			args: args{
				origin: "clouditor.io",
			},
			want: true,
		},
		{
			name: "Disallowed origin",
			fields: fields{
				allowedOrigins: []string{"clouditor.io", "localhost"},
			},
			args: args{
				origin: "clouditor.com",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cors := &corsConfig{
				allowedOrigins: tt.fields.allowedOrigins,
				allowedHeaders: tt.fields.allowedHeaders,
				allowedMethods: tt.fields.allowedMethods,
			}
			if got := cors.OriginAllowed(tt.args.origin); got != tt.want {
				t.Errorf("corsConfig.OriginAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}
