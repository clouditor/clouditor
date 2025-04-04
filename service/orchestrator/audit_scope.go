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
	"slices"
	"strings"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/evaluation"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/logging"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/service"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (svc *Service) CreateAuditScope(ctx context.Context, req *orchestrator.CreateAuditScopeRequest) (res *orchestrator.AuditScope, err error) {
	if req.AuditScope == nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", api.ErrEmptyRequest)
	}

	// We want to add the UUID without retrieving it from the request, so we need to include it first and then perform the validation check.
	req.AuditScope.Id = uuid.NewString()

	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Check, if this request has access to the target of evaluation according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessCreate, req) {
		return nil, service.ErrPermissionDenied
	}

	// Create the Audit Scope
	err = svc.storage.Create(&req.AuditScope)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v: %v", persistence.ErrDatabase, err)
	}

	go svc.informToeHooks(ctx, &orchestrator.AuditScopeChangeEvent{Type: orchestrator.AuditScopeChangeEvent_TYPE_AUDIT_SCOPE_CREATED, AuditScope: req.AuditScope}, nil)

	res = req.AuditScope

	logging.LogRequest(log, logrus.DebugLevel, logging.Create, req, fmt.Sprintf("and Catalog '%s'", req.AuditScope.GetCatalogId()))

	return
}

// GetAuditScope implements method for getting an Audit Scope, e.g. to show its state in the UI
func (svc *Service) GetAuditScope(ctx context.Context, req *orchestrator.GetAuditScopeRequest) (res *orchestrator.AuditScope, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	res = new(orchestrator.AuditScope)
	err = svc.storage.Get(res, "id = ?", req.GetAuditScopeId())
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "%v", api.ErrAuditScopeNotFound)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "%v: %v", persistence.ErrDatabase, err)
	}

	// Check if client is allowed to access the audit scope
	all, allowed := svc.authz.AllowedTargetOfEvaluations(ctx)
	if !all && !slices.Contains(allowed, res.TargetOfEvaluationId) {
		// Important to nil the response since it is set already
		return nil, status.Error(codes.PermissionDenied, service.ErrPermissionDenied.Error())
	}

	return res, nil
}

// ListAuditScopes implements a method for listing the Audit Scopes
func (svc *Service) ListAuditScopes(ctx context.Context, req *orchestrator.ListAuditScopesRequest) (res *orchestrator.ListAuditScopesResponse, err error) {
	var allowed []string
	var all bool

	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Retrieve a list of allowed target of evaluations according to our
	// authorization strategy. No need to specify any conditions to our storage
	// request, if we are allowed to see all target of evaluations.
	all, allowed = svc.authz.AllowedTargetOfEvaluations(ctx)

	// The content of the filtered target of evaluation ID must be in the list of allowed target of evaluation IDs,
	// unless one can access *all* the target of evaluations.
	if !all && req.Filter != nil && req.Filter.TargetOfEvaluationId != nil && !slices.Contains(allowed, req.Filter.GetTargetOfEvaluationId()) {
		return nil, service.ErrPermissionDenied
	}

	var query []string
	var args []any

	// Filtering the audit scopes by
	// * target of evaluation ID
	// * catalog ID
	if req.Filter != nil {
		if req.Filter.TargetOfEvaluationId != nil {
			query = append(query, "target_of_evaluation_id = ?")
			args = append(args, req.Filter.GetTargetOfEvaluationId())
			// conds = append(conds, "target_of_evaluation_id = ?", req.TargetOfEvaluationId)
		}
		if req.Filter.CatalogId != nil {
			query = append(query, "catalog_id = ?")
			args = append(args, req.Filter.GetCatalogId())
			// conds = append(conds, "catalog_id = ?", req.CatalogId)
		}
	}

	res = new(orchestrator.ListAuditScopesResponse)

	// In any case, we need to make sure that we only select audit scopes that we
	// have access to (if we do not have access to all)
	if !all {
		query = append(query, "target_of_evaluation_id IN ?")
		args = append(args, allowed)
	}

	// Join query with AND and prepend the query
	args = append([]any{strings.Join(query, " AND ")}, args...)

	// Paginate the audit scopes according to the request
	res.AuditScopes, res.NextPageToken, err = service.PaginateStorage[*orchestrator.AuditScope](req, svc.storage, service.DefaultPaginationOpts, args...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate audit scopes: %v", err)
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

	// Check, if this request has access to the target of evaluation according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessUpdate, req) {
		return nil, service.ErrPermissionDenied
	}

	res = req.AuditScope

	err = svc.storage.Update(res, "Id = ?", req.AuditScope.GetId())
	if err != nil && errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "%v", api.ErrAuditScopeNotFound)
	} else if err != nil && errors.Is(err, persistence.ErrConstraintFailed) {
		return nil, status.Errorf(codes.NotFound, "%v", persistence.ErrConstraintFailed)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "%v: %v", persistence.ErrDatabase, err)
	}

	go svc.informToeHooks(ctx, &orchestrator.AuditScopeChangeEvent{Type: orchestrator.AuditScopeChangeEvent_TYPE_AUDIT_SCOPE_UPDATED, AuditScope: req.AuditScope}, nil)

	logging.LogRequest(log, logrus.DebugLevel, logging.Update, req, fmt.Sprintf("and Catalog '%s'", req.AuditScope.GetCatalogId()))

	return
}

// RemoveAuditScope implements method for removing an Audit Scope
func (svc *Service) RemoveAuditScope(ctx context.Context, req *orchestrator.RemoveAuditScopeRequest) (response *emptypb.Empty, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Get audit scope
	auditScope, err := svc.GetAuditScope(context.Background(), &orchestrator.GetAuditScopeRequest{
		AuditScopeId: req.GetAuditScopeId(),
	})
	if err != nil {
		// the GetAuditScope method already checks the errors and we can return the error as it is
		return nil, err
	}

	// TODO(all): Currently we do not need that check as it is already done when getting the audit scope. But we will need it if we will check the request type as well.
	// Check, if this request has access to the target of evaluation according to our authorization strategy.
	if !svc.authz.CheckAccess(ctx, service.AccessDelete, auditScope) {
		return nil, service.ErrPermissionDenied
	}

	// Delete entry
	err = svc.storage.Delete(&orchestrator.AuditScope{}, "Id = ?", req.GetAuditScopeId())
	if err != nil { // Only internal errors left since others (Permission and NotFound) are already covered
		return nil, status.Errorf(codes.Internal, "%v: %v", persistence.ErrDatabase, err)
	}

	go svc.informToeHooks(ctx, &orchestrator.AuditScopeChangeEvent{Type: orchestrator.AuditScopeChangeEvent_TYPE_AUDIT_SCOPE_REMOVED, AuditScope: auditScope}, nil)

	if req.RemoveEvaluationResults {
		err := svc.storage.Delete(&evaluation.EvaluationResult{}, "audit_scope_id = ?", req.AuditScopeId)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "database error: %v", err)
		}
	}

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
