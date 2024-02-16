// Copyright 2023 Fraunhofer AISEC
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

package runtime

import (
	"testing"
	"time"

	"clouditor.io/clouditor/v2/internal/util"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestRuntime_VersionString(t *testing.T) {
	type fields struct {
		ReleaseVersion *string
		CommitHash     string
		CommitTime     *timestamppb.Timestamp
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "release",
			fields: fields{
				ReleaseVersion: util.Ref("v2.0.0"),
			},
			want: "v2.0.0",
		},
		{
			name: "pseudo-version",
			fields: fields{
				ReleaseVersion: nil,
				CommitHash:     "1234567890ab",
				CommitTime:     timestamppb.New(time.Unix(0, 0)),
			},
			want: "v0.0.0-19700101000000-1234567890ab",
		},
		{
			name:   "no-version",
			fields: fields{},
			want:   "v0.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Runtime{
				ReleaseVersion: tt.fields.ReleaseVersion,
				CommitHash:     tt.fields.CommitHash,
				CommitTime:     tt.fields.CommitTime,
			}
			if got := r.VersionString(); got != tt.want {
				t.Errorf("Runtime.VersionString() = %v, want %v", got, tt.want)
			}
		})
	}
}
