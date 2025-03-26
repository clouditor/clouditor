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

package service

import (
	"testing"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/testutil"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/persistence"
)

func TestPaginateSlice(t *testing.T) {
	type args struct {
		req    api.PaginatedRequest
		values []int
		opts   PaginationOpts
	}
	tests := []struct {
		name     string
		args     args
		wantPage assert.Want[[]int]
		wantNbt  assert.Want[string]
		wantErr  assert.WantErr
	}{
		{
			name: "first page",
			args: args{
				req: &orchestrator.ListAssessmentResultsRequest{
					PageSize:  2,
					PageToken: "",
				},
				values: []int{1, 2, 3, 4, 5},
				opts:   PaginationOpts{10, 10},
			},
			wantPage: func(t *testing.T, got []int) bool {
				return assert.Equal(t, []int{1, 2}, got)
			},
			wantNbt: func(t *testing.T, got string) bool {
				return assert.Equal(t, "CAIQAg==", got)
			},
			wantErr: assert.Nil[error],
		},
		{
			name: "next page",
			args: args{
				req: &orchestrator.ListAssessmentResultsRequest{
					PageSize:  2,
					PageToken: "CAIQAg==",
				},
				values: []int{1, 2, 3, 4, 5},
				opts:   PaginationOpts{10, 10},
			},
			wantPage: func(t *testing.T, got []int) bool {
				return assert.Equal(t, []int{3, 4}, got)
			},
			wantNbt: func(t *testing.T, got string) bool {
				return assert.Equal(t, "CAQQAg==", got)
			},
			wantErr: assert.Nil[error],
		},
		{
			name: "last page",
			args: args{
				req: &orchestrator.ListAssessmentResultsRequest{
					PageSize:  2,
					PageToken: "CAQQAg==",
				},
				values: []int{1, 2, 3, 4, 5},
				opts:   PaginationOpts{10, 10},
			},
			wantPage: func(t *testing.T, got []int) bool {
				return assert.Equal(t, []int{5}, got)
			},
			wantNbt: func(t *testing.T, got string) bool {
				return assert.Equal(t, "", got)
			},
			wantErr: assert.Nil[error],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPage, gotNbt, err := PaginateSlice(tt.args.req, tt.args.values, func(a int, b int) bool { return a < b }, tt.args.opts)

			tt.wantErr(t, err)
			tt.wantNbt(t, gotNbt)
			tt.wantPage(t, gotPage)
		})
	}
}

func TestPaginateStorage(t *testing.T) {
	type args struct {
		req     api.PaginatedRequest
		storage persistence.Storage
		opts    PaginationOpts
		conds   []interface{}
	}
	tests := []struct {
		name     string
		args     args
		wantPage assert.Want[[]orchestrator.TargetOfEvaluation]
		wantNbt  assert.Want[string]
		wantErr  assert.WantErr
	}{
		{
			name: "first page",
			args: args{
				req: &orchestrator.ListTargetOfEvaluationsRequest{
					PageSize:  2,
					PageToken: "",
				},
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Save(&orchestrator.TargetOfEvaluation{Id: "1"}))
					assert.NoError(t, s.Save(&orchestrator.TargetOfEvaluation{Id: "2"}))
					assert.NoError(t, s.Save(&orchestrator.TargetOfEvaluation{Id: "3"}))
					assert.NoError(t, s.Save(&orchestrator.TargetOfEvaluation{Id: "4"}))
					assert.NoError(t, s.Save(&orchestrator.TargetOfEvaluation{Id: "5"}))
				}),
				opts: PaginationOpts{10, 10},
			},
			wantPage: func(t *testing.T, got []orchestrator.TargetOfEvaluation) bool {
				want := []orchestrator.TargetOfEvaluation{
					{Id: "1", ConfiguredMetrics: []*assessment.Metric{}},
					{Id: "2", ConfiguredMetrics: []*assessment.Metric{}},
				}
				return assert.Equal(t, want, got)
			},
			wantNbt: func(t *testing.T, got string) bool {
				return assert.Equal(t, "CAIQAg==", got)
			},
			wantErr: assert.Nil[error],
		},
		{
			name: "next page",
			args: args{
				req: &orchestrator.ListTargetOfEvaluationsRequest{
					PageSize:  2,
					PageToken: "CAIQAg==",
				},
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Save(&orchestrator.TargetOfEvaluation{Id: "1"}))
					assert.NoError(t, s.Save(&orchestrator.TargetOfEvaluation{Id: "2"}))
					assert.NoError(t, s.Save(&orchestrator.TargetOfEvaluation{Id: "3"}))
					assert.NoError(t, s.Save(&orchestrator.TargetOfEvaluation{Id: "4"}))
					assert.NoError(t, s.Save(&orchestrator.TargetOfEvaluation{Id: "5"}))
				}),
				opts: PaginationOpts{10, 10},
			},
			wantPage: func(t *testing.T, got []orchestrator.TargetOfEvaluation) bool {
				want := []orchestrator.TargetOfEvaluation{
					{Id: "3", ConfiguredMetrics: []*assessment.Metric{}},
					{Id: "4", ConfiguredMetrics: []*assessment.Metric{}},
				}
				return assert.Equal(t, want, got)
			},
			wantNbt: func(t *testing.T, got string) bool {
				return assert.Equal(t, "CAQQAg==", got)
			},
			wantErr: assert.Nil[error],
		},
		{
			name: "last page",
			args: args{
				req: &orchestrator.ListTargetOfEvaluationsRequest{
					PageSize:  2,
					PageToken: "CAQQAg==",
				},
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Save(&orchestrator.TargetOfEvaluation{Id: "1"}))
					assert.NoError(t, s.Save(&orchestrator.TargetOfEvaluation{Id: "2"}))
					assert.NoError(t, s.Save(&orchestrator.TargetOfEvaluation{Id: "3"}))
					assert.NoError(t, s.Save(&orchestrator.TargetOfEvaluation{Id: "4"}))
					assert.NoError(t, s.Save(&orchestrator.TargetOfEvaluation{Id: "5"}))
				}),
				opts: PaginationOpts{10, 10},
			},
			wantPage: func(t *testing.T, got []orchestrator.TargetOfEvaluation) bool {
				want := []orchestrator.TargetOfEvaluation{{Id: "5", ConfiguredMetrics: []*assessment.Metric{}}}

				return assert.Equal(t, want, got)
			},
			wantNbt: func(t *testing.T, got string) bool {
				return assert.Equal(t, "", got)
			},
			wantErr: assert.Nil[error],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPage, gotNbt, err := PaginateStorage[orchestrator.TargetOfEvaluation](tt.args.req, tt.args.storage,
				tt.args.opts, tt.args.conds...)

			tt.wantErr(t, err)
			tt.wantNbt(t, gotNbt)
			tt.wantPage(t, gotPage)
		})
	}
}
