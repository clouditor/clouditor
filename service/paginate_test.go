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
	"reflect"
	"testing"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/testutil"
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
		wantPage []int
		wantNbt  string
		wantErr  bool
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
			wantPage: []int{1, 2},
			wantNbt:  "CAIQAg==",
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
			wantPage: []int{3, 4},
			wantNbt:  "CAQQAg==",
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
			wantPage: []int{5},
			wantNbt:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPage, gotNbt, err := PaginateSlice(tt.args.req, tt.args.values, func(a int, b int) bool { return a < b }, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("PaginateSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPage, tt.wantPage) {
				t.Errorf("PaginateSlice() gotPage = %v, want %v", gotPage, tt.wantPage)
			}
			if gotNbt != tt.wantNbt {
				t.Errorf("PaginateSlice() gotNbt = %v, want %v", gotNbt, tt.wantNbt)
			}
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
		wantPage []orchestrator.CloudService
		wantNbt  string
		wantErr  bool
	}{
		{
			name: "first page",
			args: args{
				req: &orchestrator.ListCloudServicesRequest{
					PageSize:  2,
					PageToken: "",
				},
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					_ = s.Save(&orchestrator.CloudService{Id: "1"})
					_ = s.Save(&orchestrator.CloudService{Id: "2"})
					_ = s.Save(&orchestrator.CloudService{Id: "3"})
					_ = s.Save(&orchestrator.CloudService{Id: "4"})
					_ = s.Save(&orchestrator.CloudService{Id: "5"})
				}),
				opts: PaginationOpts{10, 10},
			},
			wantPage: []orchestrator.CloudService{
				{Id: "1", ConfiguredMetrics: []*assessment.Metric{}, CatalogsInScope: []*orchestrator.Catalog{}},
				{Id: "2", ConfiguredMetrics: []*assessment.Metric{}, CatalogsInScope: []*orchestrator.Catalog{}},
			},
			wantNbt: "CAIQAg==",
		},
		{
			name: "next page",
			args: args{
				req: &orchestrator.ListCloudServicesRequest{
					PageSize:  2,
					PageToken: "CAIQAg==",
				},
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					_ = s.Save(&orchestrator.CloudService{Id: "1"})
					_ = s.Save(&orchestrator.CloudService{Id: "2"})
					_ = s.Save(&orchestrator.CloudService{Id: "3"})
					_ = s.Save(&orchestrator.CloudService{Id: "4"})
					_ = s.Save(&orchestrator.CloudService{Id: "5"})
				}),
				opts: PaginationOpts{10, 10},
			},
			wantPage: []orchestrator.CloudService{
				{Id: "3", ConfiguredMetrics: []*assessment.Metric{}, CatalogsInScope: []*orchestrator.Catalog{}},
				{Id: "4", ConfiguredMetrics: []*assessment.Metric{}, CatalogsInScope: []*orchestrator.Catalog{}},
			},
			wantNbt: "CAQQAg==",
		},
		{
			name: "last page",
			args: args{
				req: &orchestrator.ListCloudServicesRequest{
					PageSize:  2,
					PageToken: "CAQQAg==",
				},
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					_ = s.Save(&orchestrator.CloudService{Id: "1"})
					_ = s.Save(&orchestrator.CloudService{Id: "2"})
					_ = s.Save(&orchestrator.CloudService{Id: "3"})
					_ = s.Save(&orchestrator.CloudService{Id: "4"})
					_ = s.Save(&orchestrator.CloudService{Id: "5"})
				}),
				opts: PaginationOpts{10, 10},
			},
			wantPage: []orchestrator.CloudService{{Id: "5", ConfiguredMetrics: []*assessment.Metric{}, CatalogsInScope: []*orchestrator.Catalog{}}},
			wantNbt:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPage, gotNbt, err := PaginateStorage[orchestrator.CloudService](tt.args.req, tt.args.storage,
				tt.args.opts, tt.args.conds...)
			if (err != nil) != tt.wantErr {
				t.Errorf("PaginateStorage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPage, tt.wantPage) {
				t.Errorf("PaginateStorage() gotPage = %v, want %v", gotPage, tt.wantPage)
			}
			if gotNbt != tt.wantNbt {
				t.Errorf("PaginateStorage() gotNbt = %v, want %v", gotNbt, tt.wantNbt)
			}
		})
	}
}
