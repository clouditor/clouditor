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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (svc *Service) CreateTargetOfEvaluation(ctx context.Context, req *orchestrator.CreateTargetOfEvaluationRequest) (res *orchestrator.TargetOfEvaluation, err error) {
	// Validate request
	err = api.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the cloud service according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessCreate, req) {
		return nil, service.ErrPermissionDenied
	}

	// Create the ToE
	err = svc.storage.Create(&req.TargetOfEvaluation)
	if err != nil && errors.Is(err, persistence.ErrConstraintFailed) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid catalog or cloud service")
	}
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
	err = api.ValidateRequest(req)
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
	var conds = []any{gorm.WithoutPreload()}

	// Validate request
	err = api.ValidateRequest(req)
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
	err = api.ValidateRequest(req)
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
	err = api.ValidateRequest(req)
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
