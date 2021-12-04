// Copyright 2021 Fraunhofer AISEC
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

	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/persistence"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (*Service) RegisterCloudService(_ context.Context, req *orchestrator.RegisterCloudServiceRequest) (service *orchestrator.CloudService, err error) {
	if req.Service == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Service is empty")
	}

	if req.Service.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Service name is empty")
	}

	db := persistence.GetDatabase()

	service = new(orchestrator.CloudService)

	// Generate a new ID
	service.Id = uuid.NewString()
	service.Name = req.Service.Name

	// Persist the service in our database
	db.Create(&service)

	return
}

func (s *Service) ListCloudServices(_ context.Context, _ *orchestrator.ListCloudServicesRequest) (response *orchestrator.ListCloudServicesResponse, err error) {
	response = new(orchestrator.ListCloudServicesResponse)
	response.Services = make([]*orchestrator.CloudService, 0)

	err = s.db.Find(&response.Services).Error
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	return response, nil
}

func (s *Service) GetCloudService(_ context.Context, req *orchestrator.GetCloudServiceRequest) (response *orchestrator.CloudService, err error) {
	if req.ServiceId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Service id is empty")
	}

	response = new(orchestrator.CloudService)
	err = s.db.Find(&response).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "Service not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	return response, nil
}

func (s *Service) UpdateCloudService(_ context.Context, req *orchestrator.UpdateCloudServiceRequest) (response *orchestrator.CloudService, err error) {
	if req.Service == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Service is empty")
	}

	if req.ServiceId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Service id is empty")
	}

	var count int64
	err = s.db.Model(&orchestrator.CloudService{}).Count(&count).Error
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	if count == 0 {
		return nil, status.Error(codes.NotFound, "Service not found")
	}

	req.Service.Id = req.ServiceId
	err = s.db.Save(&req.Service).Error
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %s", err)
	}

	return req.Service, nil
}

// CreateDefaultTargetCloudService creates a new "default" target cloud services,
// if no target service exists in the database. If a new cloud service was created, it will be returned.
func (s *Service) CreateDefaultTargetCloudService() (service *orchestrator.CloudService, err error) {
	var count int64
	s.db.Model(&orchestrator.CloudService{}).Count(&count)

	if count == 0 {
		// Create a default target cloud service
		service, err = s.RegisterCloudService(context.Background(),
			&orchestrator.RegisterCloudServiceRequest{
				Service: &orchestrator.CloudService{
					Name:        "default",
					Description: "The default target cloud service",
				}})

		log.Infof("Created new default target cloud service %s", service.Id)

		return service, err
	}

	return
}
