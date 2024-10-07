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

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/logging"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/persistence/gorm"
	"clouditor.io/clouditor/v2/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (svc *Service) CreateAuditScope(ctx context.Context, req *orchestrator.CreateAuditScopeRequest) (res *orchestrator.AuditScope, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the certification target according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessCreate, req) {
		return nil, service.ErrPermissionDenied
	}

	// Create the Audit Scope
	err = svc.storage.Create(&req.AuditScope)
	if err != nil && errors.Is(err, persistence.ErrConstraintFailed) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid catalog or certification target")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	go svc.informToeHooks(ctx, &orchestrator.AuditScopeChangeEvent{Type: orchestrator.AuditScopeChangeEvent_TYPE_AUDIT_SCOPE_CREATED, AuditScope: req.AuditScope}, nil)

	res = req.AuditScope

	logging.LogRequest(log, logrus.DebugLevel, logging.Create, req, fmt.Sprintf("and Catalog '%s'", req.AuditScope.GetCatalogId()))

	return
}

// GetAuditScope implements method for getting a AuditScope, e.g. to show its state in the UI
func (svc *Service) GetAuditScope(ctx context.Context, req *orchestrator.GetAuditScopeRequest) (response *orchestrator.AuditScope, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the certification target according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessRead, req) {
		return nil, service.ErrPermissionDenied
	}

	response = new(orchestrator.AuditScope)
	err = svc.storage.Get(response, gorm.WithoutPreload(), "certification_target_id = ? AND catalog_id = ?", req.CertificationTargetId, req.CatalogId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "Audit Scope not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}
	return response, nil
}

// ListTargetsOfEvaluation implements method for getting a AuditScope
func (svc *Service) ListTargetsOfEvaluation(ctx context.Context, req *orchestrator.ListAuditScopesRequest) (res *orchestrator.ListAuditScopesResponse, err error) {
	var conds = []any{gorm.WithoutPreload()}

	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the certification target according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessRead, req) {
		return nil, service.ErrPermissionDenied
	}

	// Either the certification_target_id or the catalog_id is set and the conds are added accordingly.
	if req.GetCertificationTargetId() != "" {
		conds = append(conds, "certification_target_id = ?", req.CertificationTargetId)
	} else if req.GetCatalogId() != "" {
		conds = append(conds, "catalog_id = ?", req.CatalogId)
	}

	res = new(orchestrator.ListAuditScopesResponse)
	res.AuditScope, res.NextPageToken, err = service.PaginateStorage[*orchestrator.AuditScope](req, svc.storage, service.DefaultPaginationOpts, conds...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
	}
	return
}

// UpdateAuditScope implements method for updating an existing AuditScope
func (svc *Service) UpdateAuditScope(ctx context.Context, req *orchestrator.UpdateAuditScopeRequest) (res *orchestrator.AuditScope, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the certification target according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessUpdate, req) {
		return nil, service.ErrPermissionDenied
	}

	res = req.AuditScope

	err = svc.storage.Update(res, "certification_target_id = ? AND catalog_id = ?", req.AuditScope.GetCertificationTargetId(), req.AuditScope.GetCatalogId())

	if err != nil && errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Audit Scope not found")
	} else if err != nil && errors.Is(err, persistence.ErrConstraintFailed) {
		return nil, status.Error(codes.NotFound, "Audit Scope not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	go svc.informToeHooks(ctx, &orchestrator.AuditScopeChangeEvent{Type: orchestrator.AuditScopeChangeEvent_TYPE_AUDIT_SCOPE_UPDATED, AuditScope: req.AuditScope}, nil)

	logging.LogRequest(log, logrus.DebugLevel, logging.Update, req, fmt.Sprintf("and Catalog '%s'", req.AuditScope.GetCatalogId()))

	return
}

// RemoveAuditScope implements method for removing a AuditScope
func (svc *Service) RemoveAuditScope(ctx context.Context, req *orchestrator.RemoveAuditScopeRequest) (response *emptypb.Empty, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the certification target according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessDelete, req) {
		return nil, service.ErrPermissionDenied
	}

	err = svc.storage.Delete(&orchestrator.AuditScope{}, "certification_target_id = ? AND catalog_id = ?", req.CertificationTargetId, req.CatalogId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "Audit Scope not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	// Since we don't have a AuditScope object, we create one to be able to inform the hook about the deleted AuditScope.
	auditScope := &orchestrator.AuditScope{
		CertificationTargetId: req.GetCertificationTargetId(),
		CatalogId:             req.GetCatalogId(),
	}
	go svc.informToeHooks(ctx, &orchestrator.AuditScopeChangeEvent{Type: orchestrator.AuditScopeChangeEvent_TYPE_AUDIT_SCOPE_REMOVED, AuditScope: auditScope}, nil)

	logging.LogRequest(log, logrus.DebugLevel, logging.Remove, req, fmt.Sprintf("and Catalog '%s'", req.GetCatalogId()))

	return &emptypb.Empty{}, nil
}

// informToeHooks informs the registered hook function either of a event change for the Audit Scope or Control Monitoring Status
func (s *Service) informToeHooks(ctx context.Context, event *orchestrator.AuditScopeChangeEvent, err error) {
	s.hookMutex.RLock()
	hooks := s.auditScopeHooks
	defer s.hookMutex.RUnlock()

	// Inform our hook, if we have any
	if len(hooks) > 0 {
		for _, hook := range hooks {
			// We could do hook concurrent again (assuming different hooks don't interfere with each other)
			hook(ctx, event, err)
		}
	}
}

// RegisterToeHook registers the Audit Scope hook function
func (s *Service) RegisterToeHook(hook orchestrator.AuditScopeHookFunc) {
	s.hookMutex.Lock()
	defer s.hookMutex.Unlock()
	s.auditScopeHooks = append(s.auditScopeHooks, hook)
}
