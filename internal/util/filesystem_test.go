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

package util

import (
	"io"
	"testing"

	"clouditor.io/clouditor/v2/internal/testutil/assert"
)

func TestGetJSONFilenames(t *testing.T) {
	type args struct {
		folder string
	}
	tests := []struct {
		name    string
		args    args
		want    assert.Want[[]string]
		wantErr assert.WantErr
	}{
		{
			name: "Empty input folder",
			args: args{
				folder: "",
			},
			want: assert.Nil[[]string],
			wantErr: func(t *testing.T, gotErr error) bool {
				return assert.ErrorContains(t, gotErr, "open : no such file or directory")
			},
		},
		{
			name: "Happy path",
			args: args{
				folder: "../testdata/catalogs",
			},
			want: func(t *testing.T, got []string) bool {
				return assert.Equal[[]string](t, []string{"../testdata/catalogs/test_catalog.json"}, got)
			},
			wantErr: assert.Nil[error],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetJSONFilenames(tt.args.folder)
			tt.wantErr(t, err)
			tt.want(t, got)
		})
	}
}

func TestExpandPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name            string
		userHomeDirFunc func() (string, error)
		args            args
		want            string
		wantErr         assert.WantErr
	}{
		{
			name: "fail",
			userHomeDirFunc: func() (string, error) {
				return "", io.EOF
			},
			args: args{
				path: "~",
			},
			want: "",
			wantErr: func(t *testing.T, gotErr error) bool {
				return assert.ErrorContains(t, gotErr, "could not find retrieve current user: EOF")
			},
		},
		{
			name: "happy path with home",
			userHomeDirFunc: func() (string, error) {
				return "/home/test", nil
			},
			args: args{
				path: "~/test",
			},
			want: "/home/test/test",

			wantErr: assert.Nil[error],
		},
		{
			name: "happy path relative",
			userHomeDirFunc: func() (string, error) {
				return "/home/test", nil
			},
			args: args{
				path: "test",
			},
			wantErr: assert.Nil[error],
			want:    "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := userHomeDirFunc
			userHomeDirFunc = tt.userHomeDirFunc
			defer func() {
				userHomeDirFunc = old
			}()

			got, err := ExpandPath(tt.args.path)

			assert.Equal(t, tt.want, got)
			tt.wantErr(t, err)
		})
	}
}
