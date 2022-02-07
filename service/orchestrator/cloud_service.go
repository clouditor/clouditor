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
	"clouditor.io/clouditor/api/orchestrator"
	"context"
	"errors"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

const DefaultTargetCloudServiceId = "00000000-0000-0000-000000000000"
const DefaultTargetCloudServiceName = "default"
const DefaultTargetCloudServiceDescription = "The default target cloud service"

// RegisterCloudService implements method for OrchestratorServer interface for registering a cloud service
func (s *Service) RegisterCloudService(_ context.Context, req *orchestrator.RegisterCloudServiceRequest) (service *orchestrator.CloudService, err error) {
	if err = req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
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
	response.Services = make([]*orchestrator.CloudService, 0, 10)

	err = s.db.Read(&response.Services)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	return response, nil
}

// GetCloudService implements method for OrchestratorServer interface for getting a cloud service with provided id
func (s *Service) GetCloudService(_ context.Context, req *orchestrator.GetCloudServiceRequest) (response *orchestrator.CloudService, err error) {
	if err = req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	response = new(orchestrator.CloudService)
	err = s.db.Read(response, "Id = ?", req.ServiceId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "service not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	return response, nil
}

// UpdateCloudService implements method for OrchestratorServer interface for updating a cloud service
// TODO(all): remove ServiceId from request since it is accessible in service already
func (s *Service) UpdateCloudService(_ context.Context, req *orchestrator.UpdateCloudServiceRequest) (response *orchestrator.CloudService, err error) {
	if err = req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// Check if cloud service is in the DB
	_, err = s.GetCloudService(context.Background(), &orchestrator.GetCloudServiceRequest{ServiceId: req.ServiceId})

	// Update cloud service
	// TODO(lebogg): See questions above: Copying serviceId to the struct makes it unnecessary in the first place IMO
	req.Service.Id = req.ServiceId
	if err != nil {
		return nil, err
	}
	err = s.db.Update(req.Service)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	return req.Service, nil
}

// RemoveCloudService implements method for OrchestratorServer interface for removing a cloud service
func (s *Service) RemoveCloudService(_ context.Context, req *orchestrator.RemoveCloudServiceRequest) (response *emptypb.Empty, err error) {
	if err = req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = s.db.Delete(&orchestrator.CloudService{}, req.ServiceId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "service not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	return &emptypb.Empty{}, nil
}

// CreateDefaultTargetCloudService implements method for OrchestratorServer interface for creating a new "default"
// target cloud services, if no target service exists in the database.
// Returns new created target cloud service. Otherwise, returns nil (if a record exists) or error if sth went wrong
func (s *Service) CreateDefaultTargetCloudService() (*orchestrator.CloudService, error) {
	var currentServices []*orchestrator.CloudService
	err := s.db.Read(&currentServices)
	// Error while reading DB
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, status.Errorf(codes.Internal, "db error: %v", err)
	} else if len(currentServices) == 0 {
		// There are no cloud services yet -> Create a default target cloud service and return it
		service :=
			&orchestrator.CloudService{
				Id:          DefaultTargetCloudServiceId,
				Name:        DefaultTargetCloudServiceName,
				Description: DefaultTargetCloudServiceDescription,
			}
		// Save it directly into the database, so that we can set the ID
		err = s.db.Create(&service)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error while writing new target cloud service to db: %v", err)
		}

		return service, nil
	} else {
		// There is at least one cloud service -> Return no new default target service and no error
		return nil, nil
	}
}
