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
	"testing"

	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
)

func TestValidate(t *testing.T) {
	var nilReq *orchestrator.CreateTargetOfEvaluationRequest = nil

	type args struct {
		req IncomingRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Missing request",
			args: args{
				req: nil,
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrEmptyRequest.Error())
			},
		},
		{
			name: "Missing request",
			args: args{
				req: nilReq,
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrEmptyRequest.Error())
			},
		},
		{
			name: "Invalid request",
			args: args{
				req: &orchestrator.CreateTargetOfEvaluationRequest{},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Contains(t, err.Error(), "invalid request")
			},
		},
		{
			name: "Happy path",
			args: args{
				req: &orchestrator.CreateTargetOfEvaluationRequest{
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CertificationTargetId: "11111111-1111-1111-1111-111111111111",
						CatalogId:             "0000",
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.args.req)
			tt.wantErr(t, err)
		})
	}
}
