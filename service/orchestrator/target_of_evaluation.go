// Copyright 2021-2022 Fraunhofer AISEC
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

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/logging"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/service"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	DefaultTargetOfEvaluationId = "00000000-0000-0000-0000-000000000000"
)

func (s *Service) CreateTargetOfEvaluation(ctx context.Context, req *orchestrator.CreateTargetOfEvaluationRequest) (res *orchestrator.TargetOfEvaluation, err error) {
	// A new target of evaluation typically does not contain a UUID; therefore, we will add one here. This must be done before the validation check to prevent validation failure.
	req.TargetOfEvaluation.Id = uuid.NewString()

	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	res = new(orchestrator.TargetOfEvaluation)

	res.Id = req.TargetOfEvaluation.Id
	res.Name = req.TargetOfEvaluation.Name
	res.Description = req.TargetOfEvaluation.Description

	now := timestamppb.Now()

	res.CreatedAt = now
	res.UpdatedAt = now

	res.Metadata = &orchestrator.TargetOfEvaluation_Metadata{}
	if req.TargetOfEvaluation.Metadata != nil {
		res.Metadata.Labels = req.TargetOfEvaluation.Metadata.Labels
		res.Metadata.Icon = req.TargetOfEvaluation.Metadata.Icon
	}

	// Persist the target of evaluation in our database
	err = s.storage.Create(res)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not add target of evaluation to the database: %v", err)
	}

	go s.informHooks(ctx, res, nil)

	logging.LogRequest(log, logrus.DebugLevel, logging.Register, req)

	return
}

// ListTargetsOfEvaluation implements method for OrchestratorServer interface for listing all target of evaluations
func (svc *Service) ListTargetsOfEvaluation(ctx context.Context, req *orchestrator.ListTargetsOfEvaluationRequest) (
	res *orchestrator.ListTargetsOfEvaluationResponse, err error) {
	var conds []any
	var allowed []string
	var all bool

	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	res = new(orchestrator.ListTargetsOfEvaluationResponse)

	// Retrieve list of allowed target of evaluation according to our authorization strategy. No need to specify any conditions
	// to our storage request, if we are allowed to see all target of evaluations.
	all, allowed = svc.authz.AllowedTargetOfEvaluations(ctx)
	if !all {
		conds = append(conds, allowed)
	}

	// Paginate the target of evaluations according to the request
	res.Targets, res.NextPageToken, err = service.PaginateStorage[*orchestrator.TargetOfEvaluation](req, svc.storage,
		service.DefaultPaginationOpts, conds...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
	}

	return
}

// GetTargetOfEvaluation implements method for OrchestratorServer interface for getting a target of evaluation with provided id
func (s *Service) GetTargetOfEvaluation(ctx context.Context, req *orchestrator.GetTargetOfEvaluationRequest) (response *orchestrator.TargetOfEvaluation, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the target of evaluation according to our authorization strategy.
	if !s.authz.CheckAccess(ctx, service.AccessRead, req) {
		return nil, service.ErrPermissionDenied
	}

	response = new(orchestrator.TargetOfEvaluation)

	err = s.storage.Get(response, "Id = ?", req.TargetOfEvaluationId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "target of evaluation not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	return response, nil
}

// UpdateTargetOfEvaluation implements method for OrchestratorServer interface for updating a target of evaluation
func (s *Service) UpdateTargetOfEvaluation(ctx context.Context, req *orchestrator.UpdateTargetOfEvaluationRequest) (res *orchestrator.TargetOfEvaluation, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the target of evaluation according to our authorization strategy.
	if !s.authz.CheckAccess(ctx, service.AccessUpdate, req) {
		return nil, service.ErrPermissionDenied
	}

	count, err := s.storage.Count(req.TargetOfEvaluation, "id = ?", req.TargetOfEvaluation.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	if count == 0 {
		return nil, status.Error(codes.NotFound, "service not found")
	}

	// Add id to response because otherwise it will overwrite ID with empty string
	res = req.TargetOfEvaluation
	res.UpdatedAt = timestamppb.Now()

	// Since UpdateTargetOfEvaluation is a PUT method, we use storage.Save
	err = s.storage.Save(res, "Id = ?", req.TargetOfEvaluation.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	go s.informHooks(ctx, res, nil)

	logging.LogRequest(log, logrus.DebugLevel, logging.Update, req)

	return
}

// RemoveTargetOfEvaluation implements method for OrchestratorServer interface for removing a target of evaluation
func (s *Service) RemoveTargetOfEvaluation(ctx context.Context, req *orchestrator.RemoveTargetOfEvaluationRequest) (response *emptypb.Empty, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the target of evaluation according to our authorization strategy.
	if !s.authz.CheckAccess(ctx, service.AccessDelete, req) {
		return nil, service.ErrPermissionDenied
	}

	err = s.storage.Delete(&orchestrator.TargetOfEvaluation{Id: req.TargetOfEvaluationId})
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "target of evaluation not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	go s.informHooks(ctx, nil, nil)

	logging.LogRequest(log, logrus.DebugLevel, logging.Remove, req)

	return &emptypb.Empty{}, nil
}

// GetTargetOfEvaluationStatistics implements method for OrchestratorServer interface for retrieving target of evaluation statistics
func (s *Service) GetTargetOfEvaluationStatistics(ctx context.Context, req *orchestrator.GetTargetOfEvaluationStatisticsRequest) (response *orchestrator.GetTargetOfEvaluationStatisticsResponse, err error) {
	var (
		auditScopes *orchestrator.ListAuditScopesResponse
	)

	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the target of evaluation according to our authorization strategy.
	if !s.authz.CheckAccess(ctx, service.AccessRead, req) {
		return nil, service.ErrPermissionDenied
	}

	response = &orchestrator.GetTargetOfEvaluationStatisticsResponse{}

	// Get number of selected catalogs
	auditScopes, err = s.ListAuditScopes(ctx,
		&orchestrator.ListAuditScopesRequest{
			Filter: &orchestrator.ListAuditScopesRequest_Filter{
				TargetOfEvaluationId: &req.TargetOfEvaluationId,
			}})
	if err != nil {
		return nil, err
	}
	response.NumberOfSelectedCatalogs = int64(len(auditScopes.AuditScopes))

	// Get number of discovered resources
	resources := new(evidence.Resource)
	count, err := s.storage.Count(resources, "target_of_evaluation_id = ?", req.TargetOfEvaluationId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error counting resources: %s", err)
	}
	response.NumberOfDiscoveredResources = count

	// Get number of evidences
	ev := new(evidence.Evidence)
	count, err = s.storage.Count(ev, "target_of_evaluation_id = ?", req.TargetOfEvaluationId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error counting evidences: %s", err)
	}
	response.NumberOfEvidences = count

	// Get number of assessment results
	res := new(assessment.AssessmentResult)
	count, err = s.storage.Count(res, "target_of_evaluation_id = ?", req.TargetOfEvaluationId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error counting assessment results: %s", err)
	}
	response.NumberOfAssessmentResults = count

	return response, nil
}

// CreateDefaultTargetOfEvaluation creates a new "default" target of evaluation,
// if no target of evaluation exists in the database.
//
// If a new target of evaluation was created, it will be returned.
func (s *Service) CreateDefaultTargetOfEvaluation() (target *orchestrator.TargetOfEvaluation, err error) {
	log.Infof("Trying to create new default target of evaluation...")

	count, err := s.storage.Count(target)
	if err != nil {
		return nil, fmt.Errorf("storage error: %w", err)
	}

	if count == 0 {
		now := timestamppb.Now()

		// Create a default target of evaluation
		target =
			&orchestrator.TargetOfEvaluation{
				Id:          DefaultTargetOfEvaluationId,
				Name:        viper.GetString(config.DefaultTargetOfEvaluationNameFlag),
				Description: viper.GetString(config.DefaultTargetOfEvaluationDescriptionFlag),
				CreatedAt:   now,
				UpdatedAt:   now,
				TargetType:  orchestrator.TargetOfEvaluation_TargetType(viper.GetInt32(config.DefaultTargetOfEvaluationTypeFlag)),
			}

		// Save it in the database
		err = s.storage.Create(target)
		if err != nil {
			return nil, fmt.Errorf("storage error: %w", err)
		} else {
			log.Infof("Created new default target of evaluation: %s", target.Id)
		}
	} else {
		log.Infof("Default target of evaluation already exist.")
	}

	return
}
