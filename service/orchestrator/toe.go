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
	err = svc.storage.Create(&req.Toe)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	res = req.Toe

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

// ListTargetOfEvaluation implements method for getting a TargetOfEvaluation
func (svc *Service) ListTargetsOfEvaluation(_ context.Context, req *orchestrator.ListTargetsOfEvaluationRequest) (res *orchestrator.ListTargetsOfEvaluationResponse, err error) {
	// Validate the request
	if err = api.ValidateListRequest[*orchestrator.TargetOfEvaluation](req); err != nil {
		err = fmt.Errorf("invalid request: %w", err)
		log.Error(err)
		err = status.Errorf(codes.InvalidArgument, "%v", err)
		return
	}

	res = new(orchestrator.ListTargetsOfEvaluationResponse)
	res.Toes, res.NextPageToken, err = service.PaginateStorage[*orchestrator.TargetOfEvaluation](req, svc.storage, service.DefaultPaginationOpts, gorm.WithoutPreload())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
	}
	return
}

// UpdateTargetOfEvaluation implements method for updating an existing TargetOfEvaluation
func (svc *Service) UpdateTargetOfEvaluation(_ context.Context, req *orchestrator.UpdateTargetOfEvaluationRequest) (res *orchestrator.TargetOfEvaluation, err error) {
	if req.CloudServiceId == "" || req.CatalogId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is empty")
	}

	if req.Toe == nil {
		return nil, status.Errorf(codes.InvalidArgument, "ToE is empty")
	}

	res = req.Toe
	res.CloudServiceId = req.Toe.CloudServiceId
	res.CatalogId = req.Toe.CatalogId

	err = svc.storage.Update(res, "cloud_service_id = ? AND catalog_id = ?", res.CloudServiceId, res.CatalogId)

	if err != nil && errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "ToE not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}
	return
}

// RemoveTargetOfEvaluation implements method for removing a TargetOfEvaluation
func (svc *Service) RemoveTargetOfEvaluation(_ context.Context, req *orchestrator.RemoveTargetOfEvaluationRequest) (response *emptypb.Empty, err error) {
	if req.CloudServiceId == "" || req.CatalogId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "ToE id is empty")
	}

	err = svc.storage.Delete(&orchestrator.TargetOfEvaluation{}, "cloud_service_id = ? AND catalog_id = ?", req.CloudServiceId, req.CatalogId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "ToE not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	return &emptypb.Empty{}, nil
}
