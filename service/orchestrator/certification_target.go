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
	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/logging"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/service"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	DefaultCertificationTargetId          = "00000000-0000-0000-0000-000000000000"
	DefaultCertificationTargetName        = "default"
	DefaultCertificationTargetDescription = "The default certification target"
	DefaultCertificationTargetType        = orchestrator.CertificationTarget_TARGET_TYPE_CLOUD
)

func (s *Service) RegisterCertificationTarget(ctx context.Context, req *orchestrator.RegisterCertificationTargetRequest) (res *orchestrator.CertificationTarget, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	res = new(orchestrator.CertificationTarget)

	// Generate a new ID
	res.Id = uuid.NewString()
	res.Name = req.CertificationTarget.Name
	res.Description = req.CertificationTarget.Description

	now := timestamppb.Now()

	res.CreatedAt = now
	res.UpdatedAt = now

	res.Metadata = &orchestrator.CertificationTarget_Metadata{}
	if req.CertificationTarget.Metadata != nil {
		res.Metadata.Labels = req.CertificationTarget.Metadata.Labels
		res.Metadata.Icon = req.CertificationTarget.Metadata.Icon
	}

	// Persist the service in our database
	err = s.storage.Create(res)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not add certification target to the database: %v", err)
	}

	go s.informHooks(ctx, res, nil)

	logging.LogRequest(log, logrus.DebugLevel, logging.Register, req)

	return
}

// ListCertificationTargets implements method for OrchestratorServer interface for listing all certification targets
func (svc *Service) ListCertificationTargets(ctx context.Context, req *orchestrator.ListCertificationTargetsRequest) (
	res *orchestrator.ListCertificationTargetsResponse, err error) {
	var conds []any
	var allowed []string
	var all bool

	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	res = new(orchestrator.ListCertificationTargetsResponse)

	// Retrieve list of allowed certification target according to our authorization strategy. No need to specify any conditions
	// to our storage request, if we are allowed to see all certification targets.
	all, allowed = svc.authz.AllowedCertificationTargets(ctx)
	if !all {
		conds = append(conds, allowed)
	}

	// Paginate the certification targets according to the request
	res.Services, res.NextPageToken, err = service.PaginateStorage[*orchestrator.CertificationTarget](req, svc.storage,
		service.DefaultPaginationOpts, conds...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
	}

	return
}

// GetCertificationTarget implements method for OrchestratorServer interface for getting a certification target with provided id
func (s *Service) GetCertificationTarget(ctx context.Context, req *orchestrator.GetCertificationTargetRequest) (response *orchestrator.CertificationTarget, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the certification target according to our authorization strategy.
	if !s.authz.CheckAccess(ctx, service.AccessRead, req) {
		return nil, service.ErrPermissionDenied
	}

	response = new(orchestrator.CertificationTarget)

	err = s.storage.Get(response, "Id = ?", req.CertificationTargetId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "service not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	return response, nil
}

// UpdateCertificationTarget implements method for OrchestratorServer interface for updating a certification target
func (s *Service) UpdateCertificationTarget(ctx context.Context, req *orchestrator.UpdateCertificationTargetRequest) (res *orchestrator.CertificationTarget, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the certification target according to our authorization strategy.
	if !s.authz.CheckAccess(ctx, service.AccessUpdate, req) {
		return nil, service.ErrPermissionDenied
	}

	count, err := s.storage.Count(req.CertificationTarget, "id = ?", req.CertificationTarget.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	if count == 0 {
		return nil, status.Error(codes.NotFound, "service not found")
	}

	// Add id to response because otherwise it will overwrite ID with empty string
	res = req.CertificationTarget
	res.UpdatedAt = timestamppb.Now()

	// Since UpdateCertificationTarget is a PUT method, we use storage.Save
	err = s.storage.Save(res, "Id = ?", req.CertificationTarget.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	go s.informHooks(ctx, res, nil)

	logging.LogRequest(log, logrus.DebugLevel, logging.Update, req)

	return
}

// RemoveCertificationTarget implements method for OrchestratorServer interface for removing a certification target
func (s *Service) RemoveCertificationTarget(ctx context.Context, req *orchestrator.RemoveCertificationTargetRequest) (response *emptypb.Empty, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the certification target according to our authorization strategy.
	if !s.authz.CheckAccess(ctx, service.AccessDelete, req) {
		return nil, service.ErrPermissionDenied
	}

	err = s.storage.Delete(&orchestrator.CertificationTarget{Id: req.CertificationTargetId})
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "service not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	go s.informHooks(ctx, nil, nil)

	logging.LogRequest(log, logrus.DebugLevel, logging.Remove, req)

	return &emptypb.Empty{}, nil
}

// GetCertificationTargetStatistics implements method for OrchestratorServer interface for retrieving certification target statistics
func (s *Service) GetCertificationTargetStatistics(ctx context.Context, req *orchestrator.GetCertificationTargetStatisticsRequest) (response *orchestrator.GetCertificationTargetStatisticsResponse, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the certification target according to our authorization strategy.
	if !s.authz.CheckAccess(ctx, service.AccessRead, req) {
		return nil, service.ErrPermissionDenied
	}

	response = &orchestrator.GetCertificationTargetStatisticsResponse{}

	// Get number of selected catalogs
	CertificationTarget := new(orchestrator.CertificationTarget)
	err = s.storage.Get(CertificationTarget, "Id = ?", req.CertificationTargetId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "service not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error getting certification target: %s", err)
	}
	response.NumberOfSelectedCatalogs = int64(len(CertificationTarget.CatalogsInScope))

	// Get number of discovered resources
	resources := new(discovery.Resource)
	count, err := s.storage.Count(resources, "certification_target_id = ?", req.CertificationTargetId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error counting resources: %s", err)
	}
	response.NumberOfDiscoveredResources = count

	// Get number of evidences
	ev := new(evidence.Evidence)
	count, err = s.storage.Count(ev, "certification_target_id = ?", req.CertificationTargetId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error counting evidences: %s", err)
	}
	response.NumberOfEvidences = count

	// Get number of assessment results
	res := new(assessment.AssessmentResult)
	count, err = s.storage.Count(res, "certification_target_id = ?", req.CertificationTargetId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error counting assessment results: %s", err)
	}
	response.NumberOfAssessmentResults = count

	return response, nil
}

// CreateDefaultCertificationTarget creates a new "default" certification target,
// if no certification target exists in the database.
//
// If a new certification target was created, it will be returned.
func (s *Service) CreateDefaultCertificationTarget() (service *orchestrator.CertificationTarget, err error) {
	log.Infof("Trying to create new default certification target...")

	count, err := s.storage.Count(service)
	if err != nil {
		return nil, fmt.Errorf("storage error: %w", err)
	}

	if count == 0 {
		now := timestamppb.Now()

		// Create a default certification target
		service =
			&orchestrator.CertificationTarget{
				Id:          DefaultCertificationTargetId,
				Name:        DefaultCertificationTargetName,
				Description: DefaultCertificationTargetDescription,
				CreatedAt:   now,
				UpdatedAt:   now,
				TargetType:  DefaultCertificationTargetType,
			}

		// Save it in the database
		err = s.storage.Create(service)
		if err != nil {
			return nil, fmt.Errorf("storage error: %w", err)
		} else {
			log.Infof("Created new default target certification target: %s", service.Id)
		}
	} else {
		log.Infof("Default target certification target already exist.")
	}

	return
}
