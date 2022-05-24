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

package api

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestStatusFromWrappedError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		args   args
		wantS  *status.Status
		wantOk bool
	}{
		{
			name: "err is status.Status",
			args: args{
				err: status.Errorf(codes.NotFound, "not found"),
			},
			wantS:  status.New(codes.NotFound, "not found"),
			wantOk: true,
		},
		{
			name: "wrapped in fmt.Errorf",
			args: args{
				err: fmt.Errorf("some error: %w", status.Errorf(codes.NotFound, "not found")),
			},
			wantS:  status.New(codes.NotFound, "not found"),
			wantOk: true,
		},
		{
			name: "no status",
			args: args{
				err: fmt.Errorf("some error: %w", errors.New("some other error")),
			},
			wantS:  nil,
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotS, gotOk := StatusFromWrappedError(tt.args.err)
			if !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("StatusFromWrappedError() gotS = %v, want %v", gotS, tt.wantS)
			}
			if gotOk != tt.wantOk {
				t.Errorf("StatusFromWrappedError() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
