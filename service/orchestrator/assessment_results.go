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

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/service"
)

// ListAssessmentResults is a method implementation of the orchestrator interface
func (svc *Service) ListAssessmentResults(ctx context.Context, req *assessment.ListAssessmentResultsRequest) (res *assessment.ListAssessmentResultsResponse, err error) {
	var values = maps.Values(svc.results)
	var filtered_values []*assessment.AssessmentResult
	var allowed []string
	var all bool

	// Check, if the current filter is valid according to the authorization strategy. Omitting the cloud service ID is
	// only allowed, if one can access *all* the cloud services.
	all, allowed = svc.authz.AllowedCloudServices(ctx)
	if !all && util.Deref(req.FilteredCloudServiceId) == "" {
		return nil, service.ErrPermissionDenied
	}

	// Furthermore, the content of the filtered cloud service ID must be in the list of allowed cloud service IDs,
	// unless one can access *all* the cloud services.
	if !all && util.Deref(req.FilteredCloudServiceId) != "" && !slices.Contains(allowed, util.Deref(req.FilteredCloudServiceId)) {
		return nil, service.ErrPermissionDenied
	}

	for _, v := range values {
		// Check for filtered cloud service ID
		if util.Deref(req.FilteredCloudServiceId) != "" && v.CloudServiceId != util.Deref(req.FilteredCloudServiceId) {
			continue
		}

		// Check for filter compliant
		if req.FilteredCompliant != nil && util.Deref(req.FilteredCompliant) != v.Compliant {
			continue
		}

		// Check for filter metric ID
		if req.FilteredMetricId != nil && util.Deref(req.FilteredMetricId) != v.MetricId {
			continue
		}

		filtered_values = append(filtered_values, v)
	}

	res = new(assessment.ListAssessmentResultsResponse)

	// Paginate the results according to the request
	res.Results, res.NextPageToken, err = service.PaginateSlice(req, filtered_values, func(a *assessment.AssessmentResult, b *assessment.AssessmentResult) bool {
		return a.Timestamp.AsTime().After(b.Timestamp.AsTime())
	}, service.DefaultPaginationOpts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
	}

	return
}
