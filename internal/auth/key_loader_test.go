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

package auth

import (
	"io"
	"testing"
)

func TestExpandPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name            string
		userHomeDirFunc func() (string, error)
		args            args
		wantOut         string
		wantErr         bool
	}{
		{
			name: "fail",
			userHomeDirFunc: func() (string, error) {
				return "", io.EOF
			},
			args: args{
				path: "~",
			},
			wantErr: true,
		},
		{
			name: "happy path with home",
			userHomeDirFunc: func() (string, error) {
				return "/home/test", nil
			},
			args: args{
				path: "~/test",
			},
			wantOut: "/home/test/test",
		},
		{
			name: "happy path relative",
			userHomeDirFunc: func() (string, error) {
				return "/home/test", nil
			},
			args: args{
				path: "test",
			},
			wantOut: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := userHomeDirFunc
			userHomeDirFunc = tt.userHomeDirFunc
			defer func() {
				userHomeDirFunc = old
			}()

			gotOut, err := ExpandPath(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExpandPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut != tt.wantOut {
				t.Errorf("ExpandPath() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
