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
	"clouditor.io/clouditor/internal/logging"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/persistence/gorm"
	"clouditor.io/clouditor/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (svc *Service) CreateTargetOfEvaluation(ctx context.Context, req *orchestrator.CreateTargetOfEvaluationRequest) (res *orchestrator.TargetOfEvaluation, err error) {
	var (
		c        *orchestrator.Catalog
		controls []*orchestrator.Control
	)

	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the cloud service according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessCreate, req) {
		return nil, service.ErrPermissionDenied
	}

	// We need to retrieve some additional meta-data about the security catalog, so we need to query it as well
	c, err = svc.GetCatalog(ctx, &orchestrator.GetCatalogRequest{CatalogId: req.TargetOfEvaluation.CatalogId})
	if err != nil {
		// The error is already a gRPC error, so we can just return it
		return nil, err
	}

	// Certain catalogs do not allow scoping, in this case we need to pre-populate all controls into the scope.
	if c.AllInScope {
		controls, err = api.ListAllPaginated(&orchestrator.ListControlsRequest{CatalogId: c.Id}, func(ctx context.Context, req *orchestrator.ListControlsRequest, opts ...grpc.CallOption) (*orchestrator.ListControlsResponse, error) {
			return svc.ListControls(ctx, req)
		}, func(res *orchestrator.ListControlsResponse) []*orchestrator.Control {
			return res.Controls
		})
		if err != nil {
			// The error is already a gRPC error, so we can just return it
			return nil, err
		}

		// If the catalog allows assurance levels, add only controls with the corresponsing assurance level.
		// Note: The upper assurance level includes the underlying assurance levels. Substantial includes basic and substantial and high include all control.
		if len(c.AssuranceLevels) > 0 {
			req.TargetOfEvaluation.ControlsInScope, err = getControls(controls, c.GetAssuranceLevels(), req.TargetOfEvaluation.GetAssuranceLevel())
			if err != nil {
				return nil, status.Errorf(codes.Internal, "assurance level error: %v", err)
			}
		} else {
			// The catalog does not allow assurance levels, add all controls.
			req.TargetOfEvaluation.ControlsInScope = controls
		}
	}

	// Create the ToE
	err = svc.storage.Create(&req.TargetOfEvaluation)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	go svc.informToeHooks(ctx, &orchestrator.TargetOfEvaluationChangeEvent{Type: orchestrator.TargetOfEvaluationChangeEvent_TYPE_TARGET_OF_EVALUATION_CREATED, TargetOfEvaluation: req.TargetOfEvaluation}, nil)

	res = req.TargetOfEvaluation

	logging.LogRequest(log, logrus.DebugLevel, logging.Create, req, fmt.Sprintf("and Catalog '%s'", req.TargetOfEvaluation.GetCatalogId()))

	return
}

// GetTargetOfEvaluation implements method for getting a TargetOfEvaluation, e.g. to show its state in the UI
func (svc *Service) GetTargetOfEvaluation(ctx context.Context, req *orchestrator.GetTargetOfEvaluationRequest) (response *orchestrator.TargetOfEvaluation, err error) {
	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the cloud service according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessRead, req) {
		return nil, service.ErrPermissionDenied
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
func (svc *Service) ListTargetsOfEvaluation(ctx context.Context, req *orchestrator.ListTargetsOfEvaluationRequest) (res *orchestrator.ListTargetsOfEvaluationResponse, err error) {
	var conds = []any{gorm.WithPreload("ControlsInScope")}

	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the cloud service according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessRead, req) {
		return nil, service.ErrPermissionDenied
	}

	// Either the cloud_service_id or the catalog_id is set and the conds are added accordingly.
	if req.GetCloudServiceId() != "" {
		conds = append(conds, "cloud_service_id = ?", req.CloudServiceId)
	} else if req.GetCatalogId() != "" {
		conds = append(conds, "catalog_id = ?", req.CatalogId)
	}

	res = new(orchestrator.ListTargetsOfEvaluationResponse)
	res.TargetOfEvaluation, res.NextPageToken, err = service.PaginateStorage[*orchestrator.TargetOfEvaluation](req, svc.storage, service.DefaultPaginationOpts, conds...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
	}
	return
}

// UpdateTargetOfEvaluation implements method for updating an existing TargetOfEvaluation
func (svc *Service) UpdateTargetOfEvaluation(ctx context.Context, req *orchestrator.UpdateTargetOfEvaluationRequest) (res *orchestrator.TargetOfEvaluation, err error) {
	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the cloud service according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessUpdate, req) {
		return nil, service.ErrPermissionDenied
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

	logging.LogRequest(log, logrus.DebugLevel, logging.Update, req, fmt.Sprintf("and Catalog '%s'", req.TargetOfEvaluation.GetCatalogId()))

	return
}

// RemoveTargetOfEvaluation implements method for removing a TargetOfEvaluation
func (svc *Service) RemoveTargetOfEvaluation(ctx context.Context, req *orchestrator.RemoveTargetOfEvaluationRequest) (response *emptypb.Empty, err error) {
	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the cloud service according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessDelete, req) {
		return nil, service.ErrPermissionDenied
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

	logging.LogRequest(log, logrus.DebugLevel, logging.Remove, req, fmt.Sprintf("and Catalog '%s'", req.GetCatalogId()))

	return &emptypb.Empty{}, nil
}

// RegisterToeHook registers the Target of Evaluation hook function
func (s *Service) RegisterToeHook(hook orchestrator.TargetOfEvaluationHookFunc) {
	s.hookMutex.Lock()
	defer s.hookMutex.Unlock()
	s.toeHooks = append(s.toeHooks, hook)
}

func (svc *Service) ListControlsInScope(ctx context.Context, req *orchestrator.ListControlsInScopeRequest) (res *orchestrator.ListControlsInScopeResponse, err error) {
	var conds = []any{gorm.WithoutPreload()}
	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the cloud service according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessRead, req) {
		return nil, service.ErrPermissionDenied
	}

	conds = append(conds, "target_of_evaluation_cloud_service_id = ? AND target_of_evaluation_catalog_id = ?", req.CloudServiceId, req.CatalogId)

	res = new(orchestrator.ListControlsInScopeResponse)
	res.ControlsInScope, res.NextPageToken, err = service.PaginateStorage[*orchestrator.ControlInScope](req, svc.storage, service.DefaultPaginationOpts, conds...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	return
}

func (svc *Service) AddControlToScope(ctx context.Context, req *orchestrator.AddControlToScopeRequest) (res *orchestrator.ControlInScope, err error) {
	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the cloud service according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessUpdate, req) {
		return nil, service.ErrPermissionDenied
	}

	err = svc.storage.Create(req.Scope)
	if err != nil && errors.Is(err, persistence.ErrUniqueConstraintFailed) {
		return nil, status.Error(codes.AlreadyExists, "entry already exists")
	} else if err != nil && errors.Is(err, persistence.ErrConstraintFailed) {
		return nil, status.Error(codes.NotFound, "ToE not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	res = req.Scope

	go svc.informToeHooks(ctx, &orchestrator.TargetOfEvaluationChangeEvent{Type: orchestrator.TargetOfEvaluationChangeEvent_TYPE_CONTROL_IN_SCOPE_ADDED, ControlInScope: req.GetScope()}, nil)

	logging.LogRequest(log, logrus.DebugLevel, logging.Add, req, fmt.Sprintf("with Control ID '%s'", req.Scope.GetControlId()))

	return
}

func (svc *Service) UpdateControlInScope(ctx context.Context, req *orchestrator.UpdateControlInScopeRequest) (res *orchestrator.ControlInScope, err error) {
	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the cloud service according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessUpdate, req) {
		return nil, service.ErrPermissionDenied
	}

	err = svc.storage.Update(req.Scope,
		"target_of_evaluation_cloud_service_id = ? AND "+
			"target_of_evaluation_catalog_id = ? AND "+
			"control_category_catalog_id = ? AND "+
			"control_category_name = ? AND "+
			"control_id = ?",
		req.Scope.TargetOfEvaluationCloudServiceId,
		req.Scope.TargetOfEvaluationCatalogId,
		req.Scope.ControlCategoryCatalogId,
		req.Scope.ControlCategoryName,
		req.Scope.ControlId)
	if err != nil && errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "ToE not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	res = req.Scope

	go svc.informToeHooks(ctx, &orchestrator.TargetOfEvaluationChangeEvent{Type: orchestrator.TargetOfEvaluationChangeEvent_TYPE_CONTROL_IN_SCOPE_UPDATED, ControlInScope: req.GetScope()}, nil)

	logging.LogRequest(log, logrus.DebugLevel, logging.Update, req, fmt.Sprintf("with Control ID '%s'", req.Scope.GetControlId()))

	return
}

func (svc *Service) RemoveControlFromScope(ctx context.Context, req *orchestrator.RemoveControlFromScopeRequest) (res *emptypb.Empty, err error) {
	// Validate request
	err = service.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the cloud service according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessDelete, req) {
		return nil, service.ErrPermissionDenied
	}

	err = svc.storage.Delete(orchestrator.ControlInScope{},
		"target_of_evaluation_cloud_service_id = ? AND "+
			"target_of_evaluation_catalog_id = ? AND "+
			"control_category_catalog_id = ? AND "+
			"control_category_name = ? AND "+
			"control_id = ?",
		req.CloudServiceId,
		req.CatalogId,
		req.CatalogId,
		req.ControlCategoryName,
		req.ControlId)
	if err != nil && errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "ToE not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	res = &emptypb.Empty{}

	go svc.informToeHooks(ctx, &orchestrator.TargetOfEvaluationChangeEvent{Type: orchestrator.TargetOfEvaluationChangeEvent_TYPE_CONTROL_IN_SCOPE_REMOVED, ControlInScope: nil}, nil)

	logging.LogRequest(log, logrus.DebugLevel, logging.Remove, req, fmt.Sprintf("with Control ID '%s'", req.GetControlId()))

	return
}

// getControls returns all controls based on the assurance level
func getControls(controls []*orchestrator.Control, levels []string, level string) ([]*orchestrator.Control, error) {
	var (
		low    []*orchestrator.Control
		medium []*orchestrator.Control
		high   []*orchestrator.Control
		c      = []*orchestrator.Control{}
	)

	// Check that levels and level is not empty
	if len(levels) < 3 {
		err := errors.New("assurance levels are empty")
		return c, err
	}

	if level == "" {
		err := errors.New("assurance level is empty")
		return c, err
	}

	// Add controls based on their assurance level to the lists low, medium and high. If a controls is not defined regarding the assurance level it is dropped.
	for i := range controls {
		switch controls[i].GetAssuranceLevel() {
		case levels[0]:
			low = append(low, controls[i])
		case levels[1]:
			medium = append(medium, controls[i])
		case levels[2]:
			high = append(high, controls[i])
		default:
			continue
		}
	}

	// Add all needed controls based on the assurance level and return
	switch level {
	case levels[0]:
		c = append(c, low...)
	case levels[1]:
		c = append(c, low...)
		c = append(c, medium...)
	case levels[2]:
		c = append(c, low...)
		c = append(c, medium...)
		c = append(c, high...)
	}

	return c, nil
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
