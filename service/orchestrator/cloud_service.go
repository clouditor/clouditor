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

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/service"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	DefaultTargetCloudServiceId          = "00000000-0000-0000-0000-000000000000"
	DefaultTargetCloudServiceName        = "default"
	DefaultTargetCloudServiceDescription = "The default target cloud service"
)

func (s *Service) RegisterCloudService(ctx context.Context, req *orchestrator.RegisterCloudServiceRequest) (service *orchestrator.CloudService, err error) {
	if req == nil {
		go s.informHooks(ctx, service, api.ErrRequestIsNil)
		return nil, status.Errorf(codes.InvalidArgument, api.ErrRequestIsNil.Error())
	}
	if req.Service == nil {
		go s.informHooks(ctx, service, orchestrator.ErrServiceIsNil)
		return nil, status.Errorf(codes.InvalidArgument, orchestrator.ErrServiceIsNil.Error())
	}
	if req.Service.Name == "" {
		go s.informHooks(ctx, service, orchestrator.ErrNameIsMissing)
		return nil, status.Errorf(codes.InvalidArgument, orchestrator.ErrNameIsMissing.Error())
	}

	service = new(orchestrator.CloudService)

	// Generate a new ID
	service.Id = uuid.NewString()
	service.Name = req.Service.Name
	service.Description = req.Service.Description

	// Persist the service in our database
	err = s.storage.Create(service)
	if err != nil {
		go s.informHooks(ctx, service, err)
		return nil, status.Errorf(codes.Internal, "could not add cloud service to the database: %v", err)
	}

	go s.informHooks(ctx, service, nil)
	return
}

// ListCloudServices implements method for OrchestratorServer interface for listing all cloud services
func (svc *Service) ListCloudServices(ctx context.Context, req *orchestrator.ListCloudServicesRequest) (
	res *orchestrator.ListCloudServicesResponse, err error) {
	var conds []any
	var allowed []string
	var all bool

	// Validate tne request
	if err = api.ValidateListRequest[*orchestrator.CloudService](req); err != nil {
		err = fmt.Errorf("invalid request: %w", err)
		log.Error(err)
		err = status.Errorf(codes.InvalidArgument, "%v", err)
		return
	}

	res = new(orchestrator.ListCloudServicesResponse)

	// Retrieve list of allowed cloud service according to our authorization strategy. No need to specify any conditions
	// to our storage request, if we are allowed to see all cloud services.
	all, allowed = svc.authz.AllowedCloudServices(ctx)
	if !all {
		conds = append(conds, allowed)
	}

	// Paginate the cloud services according to the request
	res.Services, res.NextPageToken, err = service.PaginateStorage[*orchestrator.CloudService](req, svc.storage,
		service.DefaultPaginationOpts, conds...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
	}

	return
}

// GetCloudService implements method for OrchestratorServer interface for getting a cloud service with provided id
func (s *Service) GetCloudService(ctx context.Context, req *orchestrator.GetCloudServiceRequest) (response *orchestrator.CloudService, err error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, api.ErrRequestIsNil.Error())
	}
	if req.CloudServiceId == "" {
		return nil, status.Errorf(codes.InvalidArgument, orchestrator.ErrIDIsMissing.Error())
	}

	// Check, if this request has access to the cloud service according to our authorization strategy.
	if !s.authz.CheckAccess(ctx, service.AccessRead, req) {
		return nil, service.ErrPermissionDenied
	}

	response = new(orchestrator.CloudService)
	err = s.storage.Get(response, "Id = ?", req.CloudServiceId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "service not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	return response, nil
}

// UpdateCloudService implements method for OrchestratorServer interface for updating a cloud service
func (s *Service) UpdateCloudService(ctx context.Context, req *orchestrator.UpdateCloudServiceRequest) (response *orchestrator.CloudService, err error) {
	if req.Service == nil {
		go s.informHooks(ctx, response, errors.New(codes.InvalidArgument.String()))
		return nil, status.Errorf(codes.InvalidArgument, "service is empty")
	}

	if req.CloudServiceId == "" {
		go s.informHooks(ctx, response, errors.New(codes.InvalidArgument.String()))
		return nil, status.Errorf(codes.InvalidArgument, "service id is empty")
	}

	// Check, if this request has access to the cloud service according to our authorization strategy.
	if !s.authz.CheckAccess(ctx, service.AccessUpdate, req) {
		go s.informHooks(ctx, response, service.ErrPermissionDenied)
		return nil, service.ErrPermissionDenied
	}

	count, err := s.storage.Count(req.Service, "id = ?", req.CloudServiceId)
	if err != nil {
		go s.informHooks(ctx, response, err)
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	if count == 0 {
		go s.informHooks(ctx, response, errors.New(codes.NotFound.String()))
		return nil, status.Error(codes.NotFound, "service not found")
	}

	// Add id to response because otherwise it will overwrite ID with empty string
	response = req.Service
	response.Id = req.CloudServiceId

	// Since UpdateCloudService is a PUT method, we use storage.Save
	err = s.storage.Save(response, "Id = ?", req.CloudServiceId)
	if err != nil {
		go s.informHooks(ctx, response, err)
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}
	go s.informHooks(ctx, response, nil)
	return
}

// RemoveCloudService implements method for OrchestratorServer interface for removing a cloud service
func (s *Service) RemoveCloudService(ctx context.Context, req *orchestrator.RemoveCloudServiceRequest) (response *emptypb.Empty, err error) {
	if req.CloudServiceId == "" {
		go s.informHooks(ctx, nil, errors.New(codes.InvalidArgument.String()))
		return nil, status.Errorf(codes.InvalidArgument, "service id is empty")
	}

	// Check, if this request has access to the cloud service according to our authorization strategy.
	if !s.authz.CheckAccess(ctx, service.AccessDelete, req) {
		go s.informHooks(ctx, nil, service.ErrPermissionDenied)
		return nil, service.ErrPermissionDenied
	}

	err = s.storage.Delete(&orchestrator.CloudService{Id: req.CloudServiceId})
	if errors.Is(err, persistence.ErrRecordNotFound) {
		go s.informHooks(ctx, nil, err)
		return nil, status.Errorf(codes.NotFound, "service not found")
	} else if err != nil {
		go s.informHooks(ctx, nil, err)
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	go s.informHooks(ctx, nil, nil)
	return &emptypb.Empty{}, nil
}

// CreateDefaultTargetCloudService creates a new "default" target cloud services,
// if no target service exists in the database.
//
// If a new target cloud service was created, it will be returned.
func (s *Service) CreateDefaultTargetCloudService() (service *orchestrator.CloudService, err error) {
	log.Infof("Trying to create new default target cloud service...")

	count, err := s.storage.Count(service)
	if err != nil {
		return nil, fmt.Errorf("storage error: %w", err)
	}

	if count == 0 {
		// Create a default target cloud service
		service =
			&orchestrator.CloudService{
				Id:          DefaultTargetCloudServiceId,
				Name:        DefaultTargetCloudServiceName,
				Description: DefaultTargetCloudServiceDescription,
			}

		// Save it in the database
		err = s.storage.Create(service)
		if err != nil {
			return nil, fmt.Errorf("storage error: %w", err)
		} else {
			log.Infof("Created new default target cloud service: %s", service.Id)
		}
	}

	return
}
