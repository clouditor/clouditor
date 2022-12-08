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

package orchestrator

import (
	"context"
	"errors"
	"fmt"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/persistence/gorm"
	"clouditor.io/clouditor/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (svc *Service) CreateTargetOfEvaluation(ctx context.Context, req *orchestrator.CreateTargetOfEvaluationRequest) (res *orchestrator.TargetOfEvaluation, err error) {
	var (
		c     *orchestrator.Catalog
		lcres *orchestrator.ListControlsResponse
	)

	// We need to retrieve some additional meta-data about the security catalog, so we need to query it as well
	c, err = svc.GetCatalog(ctx, &orchestrator.GetCatalogRequest{CatalogId: req.TargetOfEvaluation.CatalogId})
	if err != nil {
		// The error is already a gRPC error, so we can just return it
		return nil, err
	}

	// Certain catalogs do not allow scoping, in this case we need to pre-populate all controls into the scope.
	if c.AllInScope {
		lcres, err = svc.ListControls(ctx, &orchestrator.ListControlsRequest{CatalogId: req.TargetOfEvaluation.CatalogId})
		if err != nil {
			// The error is already a gRPC error, so we can just return it
			return nil, err
		}

		// Make all controls in scope
		// TODO: Certain catalogs differentiate between assurance levels
		req.TargetOfEvaluation.ControlsInScope = lcres.Controls
	}

	// Create the ToE
	err = svc.storage.Create(&req.TargetOfEvaluation)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	go svc.informToeHooks(ctx, &orchestrator.TargetOfEvaluationChangeEvent{Type: orchestrator.TargetOfEvaluationChangeEvent_TYPE_TARGET_OF_EVALUATION_CREATED, TargetOfEvaluation: req.TargetOfEvaluation}, nil)

	res = req.TargetOfEvaluation

	return
}

// GetTargetOfEvaluation implements method for getting a TargetOfEvaluation, e.g. to show its state in the UI
func (svc *Service) GetTargetOfEvaluation(_ context.Context, req *orchestrator.GetTargetOfEvaluationRequest) (response *orchestrator.TargetOfEvaluation, err error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, api.ErrRequestIsNil.Error())
	}
	if req.CloudServiceId == "" || req.CatalogId == "" {
		return nil, status.Errorf(codes.NotFound, orchestrator.ErrToEIDIsMissing.Error())
	}

	response = new(orchestrator.TargetOfEvaluation)
	err = svc.storage.Get(response, gorm.WithoutPreload(), "cloud_service_id = ? AND catalog_id = ?", req.CloudServiceId, req.CatalogId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "ToE not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}
	return response, nil
}

// ListTargetsOfEvaluation implements method for getting a TargetOfEvaluation
func (svc *Service) ListTargetsOfEvaluation(_ context.Context, req *orchestrator.ListTargetsOfEvaluationRequest) (res *orchestrator.ListTargetsOfEvaluationResponse, err error) {
	// Validate the request
	if err = api.ValidateListRequest[*orchestrator.TargetOfEvaluation](req); err != nil {
		err = fmt.Errorf("invalid request: %w", err)
		log.Error(err)
		err = status.Errorf(codes.InvalidArgument, "%v", err)
		return
	}

	res = new(orchestrator.ListTargetsOfEvaluationResponse)
	res.TargetOfEvaluation, res.NextPageToken, err = service.PaginateStorage[*orchestrator.TargetOfEvaluation](req, svc.storage, service.DefaultPaginationOpts, gorm.WithoutPreload())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
	}
	return
}

// UpdateTargetOfEvaluation implements method for updating an existing TargetOfEvaluation
func (svc *Service) UpdateTargetOfEvaluation(ctx context.Context, req *orchestrator.UpdateTargetOfEvaluationRequest) (res *orchestrator.TargetOfEvaluation, err error) {
	if req.TargetOfEvaluation == nil {
		return nil, status.Errorf(codes.InvalidArgument, "ToE is empty")
	}

	if req.TargetOfEvaluation.GetCloudServiceId() == "" || req.TargetOfEvaluation.GetCatalogId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is empty")
	}

	res = req.TargetOfEvaluation

	err = svc.storage.Update(res, "cloud_service_id = ? AND catalog_id = ?", req.TargetOfEvaluation.GetCloudServiceId(), req.TargetOfEvaluation.GetCatalogId())

	if err != nil && errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "ToE not found")
	} else if err != nil && errors.Is(err, persistence.ErrConstraintFailed) {
		return nil, status.Error(codes.NotFound, "ToE not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	go svc.informToeHooks(ctx, &orchestrator.TargetOfEvaluationChangeEvent{Type: orchestrator.TargetOfEvaluationChangeEvent_TYPE_TARGET_OF_EVALUATION_UPDATED, TargetOfEvaluation: req.TargetOfEvaluation}, nil)

	return
}

// RemoveTargetOfEvaluation implements method for removing a TargetOfEvaluation
func (svc *Service) RemoveTargetOfEvaluation(ctx context.Context, req *orchestrator.RemoveTargetOfEvaluationRequest) (response *emptypb.Empty, err error) {
	if req.CloudServiceId == "" || req.CatalogId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "ToE id is empty")
	}

	err = svc.storage.Delete(&orchestrator.TargetOfEvaluation{}, "cloud_service_id = ? AND catalog_id = ?", req.CloudServiceId, req.CatalogId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "ToE not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	// Since we don't have a TargetOfEvaluation object, we create one to be able to inform the hook about the deleted TargetOfEvaluation.
	toe := &orchestrator.TargetOfEvaluation{
		CloudServiceId: req.GetCloudServiceId(),
		CatalogId:      req.GetCatalogId(),
	}
	go svc.informToeHooks(ctx, &orchestrator.TargetOfEvaluationChangeEvent{Type: orchestrator.TargetOfEvaluationChangeEvent_TYPE_TARGET_OF_EVALUATION_REMOVED, TargetOfEvaluation: toe}, nil)
	return &emptypb.Empty{}, nil
}

// informToeHooks informs the registered hook function either of a event change for the Target of Evaluation or Control Monitoring Status
func (s *Service) informToeHooks(ctx context.Context, event *orchestrator.TargetOfEvaluationChangeEvent, err error) {
	s.hookMutex.RLock()
	hooks := s.toeHooks
	defer s.hookMutex.RUnlock()

	// Inform our hook, if we have any
	if len(hooks) > 0 {
		for _, hook := range hooks {
			// We could do hook concurrent again (assuming different hooks don't interfere with each other)
			hook(ctx, event, err)
		}
	}
}

// RegisterToeHook registers the Target of Evaluation hook function
func (s *Service) RegisterToeHook(hook orchestrator.TargetOfEvaluationHookFunc) {
	s.hookMutex.Lock()
	defer s.hookMutex.Unlock()
	s.toeHooks = append(s.toeHooks, hook)
}

func (svc *Service) ListControlMonitoringStatus(ctx context.Context, req *orchestrator.ListControlMonitoringStatusRequest) (res *orchestrator.ListControlMonitoringStatusResponse, err error) {
	// Check, if this request has access to the cloud service according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessRead, req) {
		return nil, service.ErrPermissionDenied
	}

	res = new(orchestrator.ListControlMonitoringStatusResponse)
	res.Status, res.NextPageToken, err = service.PaginateStorage[*orchestrator.ControlMonitoringStatus](req, svc.storage, service.DefaultPaginationOpts, gorm.WithoutPreload())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	return
}

func (svc *Service) UpdateControlMonitoringStatus(ctx context.Context, req *orchestrator.UpdateControlMonitoringStatusRequest) (res *orchestrator.ControlMonitoringStatus, err error) {
	// Check, if this request has access to the cloud service according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessRead, req) {
		return nil, service.ErrPermissionDenied
	}

	err = svc.storage.Update(req.Status,
		"target_of_evaluation_cloud_service_id = ? AND "+
			"target_of_evaluation_catalog_id = ? AND "+
			"control_category_catalog_id = ? AND "+
			"control_category_name = ? AND "+
			"control_id = ?",
		req.Status.TargetOfEvaluationCloudServiceId,
		req.Status.TargetOfEvaluationCatalogId,
		req.Status.ControlCategoryCatalogId,
		req.Status.ControlCategoryName,
		req.Status.ControlId)
	if err != nil && errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "ToE not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	res = req.Status

	return
}
