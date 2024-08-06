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
	"time"

	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_SafeTimestamp(t *testing.T) {

	testTime := time.Date(2000, 01, 20, 9, 20, 12, 123, time.UTC)
	testTimeUnix := testTime.Unix()

	type args struct {
		t *time.Time
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "Empty time",
			args: args{
				t: &time.Time{},
			},
			want: 0,
		},
		{
			name: "Time is nil",
			args: args{},
			want: 0,
		},
		{
			name: "Valid time",
			args: args{
				t: &testTime,
			},
			want: testTimeUnix,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, SafeTimestamp(tt.args.t))
		})
	}
}

func TestTimestamp(t *testing.T) {
	type args struct {
		t string
	}
	tests := []struct {
		name string
		args args
		want assert.Want[*timestamppb.Timestamp]
	}{
		{
			name: "empty input",
			args: args{},
			want: func(t *testing.T, got *timestamppb.Timestamp) bool {
				return assert.Equal(t, &timestamppb.Timestamp{}, got)
			},
		},
		{
			name: "Happy path",
			args: args{
				t: "2024-08-06T09:39:25Z",
			},
			want: func(t *testing.T, got *timestamppb.Timestamp) bool {
				time, err := time.Parse(time.RFC3339, "2024-08-06T09:39:25Z")
				assert.NoError(t, err)

				return assert.Equal(t, timestamppb.New(time), got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Timestamp(tt.args.t)

			tt.want(t, got)
		})
	}
}
