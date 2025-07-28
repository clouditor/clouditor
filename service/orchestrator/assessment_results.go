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
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/logging"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/service"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetAssessmentResult gets one assessment result by id
func (svc *Service) GetAssessmentResult(ctx context.Context, req *orchestrator.GetAssessmentResultRequest) (res *assessment.AssessmentResult, err error) {
	var (
		all     bool
		allowed []string
	)

	// Validate request
	if err = api.Validate(req); err != nil {
		return
	}

	// Fetch result
	res = new(assessment.AssessmentResult)
	err = svc.storage.Get(res, "Id = ?", req.Id)

	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "assessment result not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	// Check if target_of_evaluation_id in assessment_result is within allowed or one can access *all* the target of evaluations
	all, allowed = svc.authz.AllowedTargetOfEvaluations(ctx)

	// The content of the filtered target of evaluation ID must be in the list of allowed target of evaluation IDs,
	// unless one can access *all* the target of evaluations.
	if !all && !slices.Contains(allowed, res.GetTargetOfEvaluationId()) {
		return nil, service.ErrPermissionDenied
	}

	return
}

// ListAssessmentResults is a method implementation of the orchestrator interface
func (svc *Service) ListAssessmentResults(ctx context.Context, req *orchestrator.ListAssessmentResultsRequest) (res *orchestrator.ListAssessmentResultsResponse, err error) {
	var allowed []string
	var all bool

	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Retrieve list of allowed target of evaluation according to our authorization strategy. No need to specify any conditions
	// to our storage request, if we are allowed to see all target of evaluations.
	all, allowed = svc.authz.AllowedTargetOfEvaluations(ctx)

	// The content of the filtered target of evaluation ID must be in the list of allowed target of evaluation IDs,
	// unless one can access *all* the target of evaluations.
	if !all && req.Filter != nil && req.Filter.TargetOfEvaluationId != nil && !slices.Contains(allowed, req.Filter.GetTargetOfEvaluationId()) {
		return nil, service.ErrPermissionDenied
	}

	res = new(orchestrator.ListAssessmentResultsResponse)

	var query []string
	var args []any

	// Filtering the assessment results by
	// * target of evaluation ID
	// * compliant status
	// * metric ID
	// * tool ID
	// * assessment result ID(s)
	// * evidence ID in the history field
	if req.Filter != nil {
		if req.Filter.TargetOfEvaluationId != nil {
			query = append(query, "target_of_evaluation_id = ?")
			args = append(args, req.Filter.GetTargetOfEvaluationId())
		}
		if req.Filter.Compliant != nil {
			query = append(query, "compliant = ?")
			args = append(args, req.Filter.GetCompliant())
		}
		if req.Filter.MetricIds != nil {
			query = append(query, "metric_id IN ?")
			args = append(args, req.Filter.GetMetricIds())
		}
		if req.Filter.AssessmentResultIds != nil {
			query = append(query, "id IN ?")
			args = append(args, req.Filter.GetAssessmentResultIds())
		}
		if req.Filter.ToolId != nil {
			query = append(query, "tool_id = ?")
			args = append(args, req.Filter.ToolId)
		}
		if req.Filter.HistoryEvidenceId != nil {
			query = append(query,
				`EXISTS (
					SELECT 1
					FROM json_each(history)
					WHERE json_extract(json_each.value, '$.evidence_id') = ?
				)`,
			)
			args = append(args, req.Filter.GetHistoryEvidenceId())
		}
	}

	// In any case, we need to make sure that we only select assessment results of target of evaluations that we have access to
	// (if we do not have access to all)
	if !all {
		query = append(query, "target_of_evaluation_id IN ?")
		args = append(args, allowed)
	}

	// If we want to have it grouped by resource ID (and metric ID), we need to do a raw query
	if req.GetLatestByResourceId() {
		// In the raw SQL, we need to build the whole WHERE statement
		var where string

		if len(query) > 0 {
			where = "WHERE " + strings.Join(query, " AND ")
		}

		// Execute the raw SQL statement
		err = svc.storage.Raw(&res.Results,
			fmt.Sprintf(`WITH sorted_results AS (
				SELECT *, ROW_NUMBER() OVER (PARTITION BY resource_id, metric_id ORDER BY created_at DESC) AS row_number
				FROM assessment_results
				%s
		  	)
		  	SELECT * FROM sorted_results WHERE row_number = 1;`, where), args...)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "database error: %v", err)
		}
	} else {
		// Join query with AND and prepend the query
		args = append([]any{strings.Join(query, " AND ")}, args...)

		// Paginate the results according to the request
		res.Results, res.NextPageToken, err = service.PaginateStorage[*assessment.AssessmentResult](req, svc.storage, service.DefaultPaginationOpts, args...)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
		}
	}

	return
}

// StoreAssessmentResult is a method implementation of the orchestrator interface: It receives an assessment result and stores it
func (svc *Service) StoreAssessmentResult(ctx context.Context, req *orchestrator.StoreAssessmentResultRequest) (res *orchestrator.StoreAssessmentResultResponse, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the target of evaluation according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessRead, req) {
		return nil, service.ErrPermissionDenied
	}

	err = svc.storage.Create(req.Result)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	go svc.informHook(ctx, req.Result, nil)

	res = &orchestrator.StoreAssessmentResultResponse{}

	logging.LogRequest(log, logrus.DebugLevel, logging.Store, req)

	return res, nil
}

func (s *Service) StoreAssessmentResults(stream orchestrator.Orchestrator_StoreAssessmentResultsServer) (err error) {
	var (
		result *orchestrator.StoreAssessmentResultRequest
		res    *orchestrator.StoreAssessmentResultsResponse
	)

	for {
		result, err = stream.Recv()

		// If no more input of the stream is available, return
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			newError := fmt.Errorf("cannot receive stream request: %w", err)
			log.Error(newError)
			return status.Errorf(codes.Unknown, "%v", newError)
		}

		// Call StoreAssessmentResult() for storing a single assessment
		storeAssessmentResultReq := &orchestrator.StoreAssessmentResultRequest{
			Result: result.Result,
		}
		_, err = s.StoreAssessmentResult(stream.Context(), storeAssessmentResultReq)
		if err != nil {
			// Create response message. The StoreAssessmentResult method does not need that message, so we have to create it here for the stream response.
			res = &orchestrator.StoreAssessmentResultsResponse{
				Status:        false,
				StatusMessage: err.Error(),
			}
		} else {
			res = &orchestrator.StoreAssessmentResultsResponse{
				Status: true,
			}
		}

		log.Debugf("Assessment result received (%v)", result.Result.Id)

		err = stream.Send(res)

		// Check for send errors
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			newError := fmt.Errorf("cannot stream response to the client: %w", err)
			log.Error(newError)
			return status.Errorf(codes.Unknown, "%v", newError.Error())
		}
	}
}

func (s *Service) RegisterAssessmentResultHook(hook assessment.ResultHookFunc) {
	s.hookMutex.Lock()
	defer s.hookMutex.Unlock()
	s.AssessmentResultHooks = append(s.AssessmentResultHooks, hook)
}

func (s *Service) informHook(ctx context.Context, result *assessment.AssessmentResult, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Inform our hook, if we have any
	if s.AssessmentResultHooks != nil {
		for _, hook := range s.AssessmentResultHooks {
			// TODO(all): We could do hook concurrent again (assuming different hooks don't interfere with each other)
			hook(ctx, result, err)
		}
	}
}
