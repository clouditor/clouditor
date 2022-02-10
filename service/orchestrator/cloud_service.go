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
	"clouditor.io/clouditor/persistence"
	"context"
	"errors"

	"clouditor.io/clouditor/api/orchestrator"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const DefaultTargetCloudServiceId = "00000000-0000-0000-000000000000"
const DefaultTargetCloudServiceName = "default"
const DefaultTargetCloudServiceDescription = "The default target cloud service"

func (s *Service) RegisterCloudService(_ context.Context, req *orchestrator.RegisterCloudServiceRequest) (service *orchestrator.CloudService, err error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, orchestrator.ErrRequestIsNil.Error())
	}
	if req.Service == nil {
		return nil, status.Errorf(codes.InvalidArgument, orchestrator.ErrServiceIsNil.Error())
	}
	if req.Service.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, orchestrator.ErrNameIsMissing.Error())
	}

	service = new(orchestrator.CloudService)

	// Generate a new ID
	service.Id = uuid.NewString()
	service.Name = req.Service.Name
	service.Description = req.Service.Description

	// Persist the service in our database
	err = s.db.Create(&service)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not add cloud service to the database: %v", err)
	}

	return
}

// ListCloudServices implements method for OrchestratorServer interface for listing all cloud services
func (s *Service) ListCloudServices(_ context.Context, _ *orchestrator.ListCloudServicesRequest) (response *orchestrator.ListCloudServicesResponse, err error) {
	response = new(orchestrator.ListCloudServicesResponse)
	response.Services = make([]*orchestrator.CloudService, 0)

	err = s.db.List(&response.Services)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	return response, nil
}

// GetCloudService implements method for OrchestratorServer interface for getting a cloud service with provided id
func (s *Service) GetCloudService(_ context.Context, req *orchestrator.GetCloudServiceRequest) (response *orchestrator.CloudService, err error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, orchestrator.ErrRequestIsNil.Error())
	}
	if req.ServiceId == "" {
		return nil, status.Errorf(codes.InvalidArgument, orchestrator.ErrIDIsMissing.Error())
	}

	response = new(orchestrator.CloudService)
	err = s.db.Get(&response, "Id = ?", req.ServiceId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "service not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	return response, nil
}

// UpdateCloudService implements method for OrchestratorServer interface for updating a cloud service
func (s *Service) UpdateCloudService(_ context.Context, req *orchestrator.UpdateCloudServiceRequest) (response *orchestrator.CloudService, err error) {
	if req.Service == nil {
		return nil, status.Errorf(codes.InvalidArgument, "service is empty")
	}

	if req.ServiceId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "service id is empty")
	}

	count, err := s.db.Count(req.Service, "Id = ?", req.ServiceId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}
	if count == 0 {
		return nil, status.Error(codes.NotFound, "service not found")
	}

	// Add id to response because otherwise it will overwrite ID with empty string
	response = req.Service
	response.Id = req.ServiceId

	err = s.db.Update(response, "Id = ?", req.ServiceId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	return
}

// RemoveCloudService implements method for OrchestratorServer interface for removing a cloud service
func (s *Service) RemoveCloudService(_ context.Context, req *orchestrator.RemoveCloudServiceRequest) (response *emptypb.Empty, err error) {
	if req.ServiceId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "service id is empty")
	}

	err = s.db.Delete(&orchestrator.CloudService{}, req.ServiceId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "service not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	return &emptypb.Empty{}, nil
}

// CreateDefaultTargetCloudService creates a new "default" target cloud services,
// if no target service exists in the database.
//
// If a new target cloud service was created, it will be returned.
func (s *Service) CreateDefaultTargetCloudService() (service *orchestrator.CloudService, err error) {
	count, err := s.db.Count(service)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		// Create a default target cloud service
		service =
			&orchestrator.CloudService{
				Id:          DefaultTargetCloudServiceId,
				Name:        DefaultTargetCloudServiceName,
				Description: DefaultTargetCloudServiceDescription,
			}
		// Save it directly into the database, so that we can set the ID
		err = s.db.Create(&service)
		if err != nil {
			log.Infof("Could not create default target cloud service %s", service.Id)
		}

	}
	return
}
