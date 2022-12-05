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

package orchestrator

import (
	"context"
	"testing"
	"time"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/service"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestService_ListAssessmentResults(t *testing.T) {
	type fields struct {
		results map[string]*assessment.AssessmentResult
		storage persistence.Storage
		authz   service.AuthorizationStrategy
	}
	type args struct {
		ctx context.Context
		req *assessment.ListAssessmentResultsRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *assessment.ListAssessmentResultsResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "request is missing",
			fields: fields{
				results: map[string]*assessment.AssessmentResult{},
				authz:   &service.AuthorizationStrategyAllowAll{},
			},
			args:    args{},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid request")
			},
		},
		{
			name: "request is empty",
			fields: fields{
				results: map[string]*assessment.AssessmentResult{},
				authz:   &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				req: &assessment.ListAssessmentResultsRequest{},
			},
			wantRes: &assessment.ListAssessmentResultsResponse{},
			wantErr: assert.NoError,
		},
		{
			name: "list all with allow all",
			fields: fields{
				results: map[string]*assessment.AssessmentResult{
					"1": {Id: "1", Timestamp: timestamppb.New(time.Unix(1, 0))},
					"2": {Id: "2", Timestamp: timestamppb.New(time.Unix(0, 0))},
				},
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{req: &assessment.ListAssessmentResultsRequest{}},
			wantRes: &assessment.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					{Id: "1", Timestamp: timestamppb.New(time.Unix(1, 0))},
					{Id: "2", Timestamp: timestamppb.New(time.Unix(0, 0))},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "list all denied",
			fields: fields{
				results: map[string]*assessment.AssessmentResult{
					"1": {Id: "1", Timestamp: timestamppb.New(time.Unix(1, 0))},
					"2": {Id: "2", Timestamp: timestamppb.New(time.Unix(0, 0))},
				},
				authz: &service.AuthorizationStrategyJWT{Key: testutil.TestCustomClaims},
			},
			args: args{
				ctx: testutil.TestContextOnlyService1,
				req: &assessment.ListAssessmentResultsRequest{},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, err, service.ErrPermissionDenied)
			},
		},
		{
			name: "specify filtered cloud service ID which is not allowed",
			fields: fields{
				results: map[string]*assessment.AssessmentResult{
					"1": {Id: "1", Timestamp: timestamppb.New(time.Unix(1, 0)), CloudServiceId: testutil.TestCloudService1},
					"2": {Id: "2", Timestamp: timestamppb.New(time.Unix(0, 0)), CloudServiceId: testutil.TestCloudService2},
				},
				authz: &service.AuthorizationStrategyJWT{Key: testutil.TestCustomClaims},
			},
			args: args{
				ctx: testutil.TestContextOnlyService1,
				req: &assessment.ListAssessmentResultsRequest{
					FilteredCloudServiceId: util.Ref(testutil.TestCloudService2),
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, err, service.ErrPermissionDenied)
			},
		},
		{
			name: "return filtered cloud service ID",
			fields: fields{
				results: map[string]*assessment.AssessmentResult{
					"1": {Id: "1", Timestamp: timestamppb.New(time.Unix(1, 0)), CloudServiceId: testutil.TestCloudService1},
					"2": {Id: "2", Timestamp: timestamppb.New(time.Unix(0, 0)), CloudServiceId: testutil.TestCloudService2},
				},
				authz: &service.AuthorizationStrategyJWT{Key: testutil.TestCustomClaims},
			},
			args: args{
				ctx: testutil.TestContextOnlyService1,
				req: &assessment.ListAssessmentResultsRequest{
					FilteredCloudServiceId: util.Ref(testutil.TestCloudService1),
				},
			},
			wantRes: &assessment.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					{Id: "1", Timestamp: timestamppb.New(time.Unix(1, 0)), CloudServiceId: testutil.TestCloudService1},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "return filtered cloud service ID and filtered compliant assessment results",
			fields: fields{
				results: getAssessmentResults(),
				authz:   &service.AuthorizationStrategyJWT{Key: testutil.TestCustomClaims},
			},
			args: args{
				ctx: testutil.TestContextOnlyService1,
				req: &assessment.ListAssessmentResultsRequest{
					FilteredCloudServiceId: util.Ref(testutil.TestCloudService1),
					FilteredCompliant:      util.Ref(true),
				},
			},
			wantRes: &assessment.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					{
						Id:             "1",
						Timestamp:      timestamppb.New(time.Unix(1, 0)),
						MetricId:       "TestMetricID-1",
						Compliant:      true,
						CloudServiceId: testutil.TestCloudService1,
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "return filtered cloud service ID and filtered non-compliant assessment results",
			fields: fields{
				results: getAssessmentResults(),
				authz:   &service.AuthorizationStrategyJWT{Key: testutil.TestCustomClaims},
			},
			args: args{
				ctx: testutil.TestContextOnlyService1,
				req: &assessment.ListAssessmentResultsRequest{
					FilteredCloudServiceId: util.Ref(testutil.TestCloudService1),
					FilteredCompliant:      util.Ref(false),
				},
			},
			wantRes: &assessment.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					{
						Id:             "3",
						Timestamp:      timestamppb.New(time.Unix(1, 0)),
						MetricId:       "TestMetricID-2",
						Compliant:      false,
						CloudServiceId: testutil.TestCloudService1,
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "return filtered compliant assessment results",
			fields: fields{
				results: getAssessmentResults(),
				authz:   &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				ctx: testutil.TestContextOnlyService1,
				req: &assessment.ListAssessmentResultsRequest{
					FilteredCompliant: util.Ref(true),
				},
			},
			wantRes: &assessment.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					{
						Id:             "1",
						Timestamp:      timestamppb.New(time.Unix(1, 0)),
						MetricId:       "TestMetricID-1",
						Compliant:      true,
						CloudServiceId: testutil.TestCloudService1,
					},
					{
						Id:             "2",
						Timestamp:      timestamppb.New(time.Unix(1, 0)),
						MetricId:       "TestMetricID-2",
						Compliant:      true,
						CloudServiceId: testutil.TestCloudService2,
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "return filtered non-compliant assessment results",
			fields: fields{
				results: getAssessmentResults(),
				authz:   &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				ctx: testutil.TestContextOnlyService1,
				req: &assessment.ListAssessmentResultsRequest{
					FilteredCompliant: util.Ref(false),
				},
			},
			wantRes: &assessment.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					{
						Id:             "3",
						Timestamp:      timestamppb.New(time.Unix(1, 0)),
						MetricId:       "TestMetricID-2",
						Compliant:      false,
						CloudServiceId: testutil.TestCloudService1,
					},
					{
						Id:             "4",
						Timestamp:      timestamppb.New(time.Unix(1, 0)),
						MetricId:       "TestMetricID-1",
						Compliant:      false,
						CloudServiceId: testutil.TestCloudService2,
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "return filtered cloud service ID and one filtered metric ID",
			fields: fields{
				results: getAssessmentResults(),
				authz:   &service.AuthorizationStrategyJWT{Key: testutil.TestCustomClaims},
			},
			args: args{
				ctx: testutil.TestContextOnlyService1,
				req: &assessment.ListAssessmentResultsRequest{
					FilteredCloudServiceId: util.Ref(testutil.TestCloudService1),
					FilteredMetricId:       []string{"TestMetricID-1"},
				},
			},
			wantRes: &assessment.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					{
						Id:             "1",
						Timestamp:      timestamppb.New(time.Unix(1, 0)),
						MetricId:       "TestMetricID-1",
						Compliant:      true,
						CloudServiceId: testutil.TestCloudService1,
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "return filtered cloud service ID and two filtered metric IDs",
			fields: fields{
				results: getAssessmentResults(),
				authz:   &service.AuthorizationStrategyJWT{Key: testutil.TestCustomClaims},
			},
			args: args{
				ctx: testutil.TestContextOnlyService1,
				req: &assessment.ListAssessmentResultsRequest{
					FilteredCloudServiceId: util.Ref(testutil.TestCloudService1),
					FilteredMetricId:       []string{"TestMetricID-1", "TestMetricID-2"},
				},
			},
			wantRes: &assessment.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					{
						Id:             "1",
						Timestamp:      timestamppb.New(time.Unix(1, 0)),
						MetricId:       "TestMetricID-1",
						Compliant:      true,
						CloudServiceId: testutil.TestCloudService1,
					},
					{
						Id:             "3",
						Timestamp:      timestamppb.New(time.Unix(1, 0)),
						MetricId:       "TestMetricID-2",
						Compliant:      false,
						CloudServiceId: testutil.TestCloudService1,
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "return one filtered metric ID",
			fields: fields{
				results: getAssessmentResults(),
				authz:   &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				ctx: testutil.TestContextOnlyService1,
				req: &assessment.ListAssessmentResultsRequest{
					FilteredMetricId: []string{"TestMetricID-1"},
				},
			},
			wantRes: &assessment.ListAssessmentResultsResponse{
				Results: []*assessment.AssessmentResult{
					{
						Id:             "1",
						Timestamp:      timestamppb.New(time.Unix(1, 0)),
						MetricId:       "TestMetricID-1",
						Compliant:      true,
						CloudServiceId: testutil.TestCloudService1,
					},
					{
						Id:             "4",
						Timestamp:      timestamppb.New(time.Unix(1, 0)),
						MetricId:       "TestMetricID-1",
						Compliant:      false,
						CloudServiceId: testutil.TestCloudService2,
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Invalid cloud service id request",
			fields: fields{
				results: getAssessmentResults(),
				authz:   &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				ctx: testutil.TestContextOnlyService1,
				req: &assessment.ListAssessmentResultsRequest{
					FilteredCloudServiceId: util.Ref("testCloudServiceID"),
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, assessment.ErrCloudServiceIDIsInvalid.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				results: tt.fields.results,
				storage: tt.fields.storage,
				authz:   tt.fields.authz,
			}
			gotRes, err := svc.ListAssessmentResults(tt.args.ctx, tt.args.req)
			tt.wantErr(t, err)

			if tt.wantRes == nil {
				assert.Nil(t, gotRes)
			} else {
				for _, elem := range gotRes.Results {
					assert.Contains(t, tt.wantRes.Results, elem)
				}

				assert.Equal(t, len(gotRes.Results), len(tt.wantRes.Results))

			}
		})
	}
}

func getAssessmentResults() (results map[string]*assessment.AssessmentResult) {
	results = map[string]*assessment.AssessmentResult{
		"1": {
			Id:             "1",
			Timestamp:      timestamppb.New(time.Unix(1, 0)),
			MetricId:       "TestMetricID-1",
			Compliant:      true,
			CloudServiceId: testutil.TestCloudService1,
		},
		"2": {
			Id:             "2",
			Timestamp:      timestamppb.New(time.Unix(1, 0)),
			MetricId:       "TestMetricID-2",
			Compliant:      true,
			CloudServiceId: testutil.TestCloudService2,
		},
		"3": {
			Id:             "3",
			Timestamp:      timestamppb.New(time.Unix(1, 0)),
			MetricId:       "TestMetricID-2",
			Compliant:      false,
			CloudServiceId: testutil.TestCloudService1,
		},
		"4": {
			Id:             "4",
			Timestamp:      timestamppb.New(time.Unix(1, 0)),
			MetricId:       "TestMetricID-1",
			Compliant:      false,
			CloudServiceId: testutil.TestCloudService2,
		},
	}

	return
}
