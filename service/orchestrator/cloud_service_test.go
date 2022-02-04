// Copyright 2022 Fraunhofer AISEC
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
	"testing"

	"clouditor.io/clouditor/api/orchestrator"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func TestRegisterCloudService(t *testing.T) {
	tests := []struct {
		name string
		req  *orchestrator.RegisterCloudServiceRequest
		res  *orchestrator.CloudService
		err  error
	}{
		{
			"missing service",
			&orchestrator.RegisterCloudServiceRequest{},
			nil,
			status.Error(codes.InvalidArgument, "service is empty"),
		},
		{
			"missing service name",
			&orchestrator.RegisterCloudServiceRequest{Service: &orchestrator.CloudService{}},
			nil,
			status.Error(codes.InvalidArgument, "service name is empty"),
		},
		{
			"valid",
			&orchestrator.RegisterCloudServiceRequest{Service: &orchestrator.CloudService{Name: "test", Description: "some"}},
			&orchestrator.CloudService{Name: "test", Description: "some"},
			nil,
		},
	}

	cloudService, err := service.CreateDefaultTargetCloudService()
	assert.Nil(t, err)
	assert.NotNil(t, cloudService)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := service.RegisterCloudService(context.Background(), tt.req)

			if tt.err == nil {
				assert.Equal(t, err, tt.err)
			} else {
				assert.EqualError(t, err, tt.err.Error())
			}

			if tt.res != nil {
				assert.NotEmpty(t, res.Id)
			}

			// reset the IDs because we cannot compare them, since they are randomly generated
			if res != nil {
				res.Id = ""
			}

			if tt.res != nil {
				tt.res.Id = ""
			}

			assert.True(t, proto.Equal(res, tt.res), "%v != %v", res, tt.res)
		})
	}

	// Reset DB
	assert.Nil(t, gormX.Reset())
}

func TestService_ListCloudServices(t *testing.T) {
	var (
		listCloudServicesResponse *orchestrator.ListCloudServicesResponse
		cloudService              *orchestrator.CloudService
		err                       error
	)

	// 1st case: No services stored
	listCloudServicesResponse, err = service.ListCloudServices(context.Background(), &orchestrator.ListCloudServicesRequest{})
	assert.Nil(t, err)
	assert.NotNil(t, listCloudServicesResponse.Services)
	assert.Empty(t, listCloudServicesResponse.Services)

	// 2nd case: One service stored
	cloudService, err = service.CreateDefaultTargetCloudService()
	assert.Nil(t, err)
	assert.NotNil(t, cloudService)

	listCloudServicesResponse, err = service.ListCloudServices(context.Background(), &orchestrator.ListCloudServicesRequest{})
	assert.Nil(t, err)
	assert.NotNil(t, listCloudServicesResponse.Services)
	assert.NotEmpty(t, listCloudServicesResponse.Services)
	assert.Equal(t, len(listCloudServicesResponse.Services), 1)

	// Reset DB so that other tests work correctly
	err = service.db.Delete(&orchestrator.CloudService{}, "")
	assert.Nil(t, err)

	// Reset DB
	assert.Nil(t, gormX.Reset())
}

func TestGetCloudService(t *testing.T) {
	tests := []struct {
		name string
		req  *orchestrator.GetCloudServiceRequest
		res  *orchestrator.CloudService
		err  error
	}{
		{
			"missing request",
			nil,
			nil,
			status.Error(codes.InvalidArgument, "service id is empty"),
		},
		{
			"missing service id",
			&orchestrator.GetCloudServiceRequest{},
			nil,
			status.Error(codes.InvalidArgument, "service id is empty"),
		},
		{
			"invalid service id",
			&orchestrator.GetCloudServiceRequest{ServiceId: "does-not-exist"},
			nil,
			status.Error(codes.NotFound, "service not found"),
		},
		{
			"valid",
			&orchestrator.GetCloudServiceRequest{ServiceId: DefaultTargetCloudServiceId},
			&orchestrator.CloudService{
				Id:          DefaultTargetCloudServiceId,
				Name:        DefaultTargetCloudServiceName,
				Description: DefaultTargetCloudServiceDescription,
			},
			nil,
		},
	}

	cloudService, err := service.CreateDefaultTargetCloudService()
	assert.Nil(t, err)
	assert.NotNil(t, cloudService)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := service.GetCloudService(context.Background(), tt.req)

			if tt.err == nil {
				assert.Equal(t, tt.err, err)
			} else {
				assert.EqualError(t, err, tt.err.Error())
			}

			if tt.res != nil {
				assert.NotEmpty(t, res.Id)
			}

			assert.True(t, proto.Equal(res, tt.res), "%v != %v", res, tt.res)
		})
	}

	// Reset DB so that other tests work correctly
	err = service.db.Delete(&orchestrator.CloudService{}, "")
	assert.Nil(t, err)

	// Reset DB
	assert.Nil(t, gormX.Reset())
}

func TestService_UpdateCloudService(t *testing.T) {
	var (
		cloudService *orchestrator.CloudService
		err          error
	)

	// 1st case: Service is nil
	_, err = service.UpdateCloudService(context.Background(), &orchestrator.UpdateCloudServiceRequest{})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// 2nd case: Service ID is nil
	_, err = service.UpdateCloudService(context.Background(), &orchestrator.UpdateCloudServiceRequest{
		Service: &orchestrator.CloudService{},
	})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// 3rd case: Service not found since there are no services yet
	_, err = service.UpdateCloudService(context.Background(), &orchestrator.UpdateCloudServiceRequest{
		Service: &orchestrator.CloudService{
			Id: DefaultTargetCloudServiceId,
		},
		ServiceId: DefaultTargetCloudServiceId,
	})
	assert.Equal(t, codes.NotFound, status.Code(err))
	// 4th case: Service updated successfully
	_, err = service.CreateDefaultTargetCloudService()
	assert.Nil(t, err)
	if err != nil {
		return
	}
	cloudService, err = service.UpdateCloudService(context.Background(), &orchestrator.UpdateCloudServiceRequest{
		Service: &orchestrator.CloudService{
			Id:          DefaultTargetCloudServiceId,
			Name:        "NewName",
			Description: "NewDescription",
		},
		ServiceId: DefaultTargetCloudServiceId,
	})
	assert.Nil(t, err)
	assert.NotNil(t, cloudService)
	assert.Equal(t, "NewName", cloudService.Name)
	assert.Equal(t, "NewDescription", cloudService.Description)

	cloudService, err = service.UpdateCloudService(context.Background(), &orchestrator.UpdateCloudServiceRequest{
		ServiceId: DefaultTargetCloudServiceId,
		Service: &orchestrator.CloudService{
			Id:          "",
			Name:        "",
			Description: "",
		},
	})

	// Reset DB
	assert.Nil(t, gormX.Reset())
}

func TestService_RemoveCloudService(t *testing.T) {
	var (
		cloudServiceResponse      *orchestrator.CloudService
		err                       error
		listCloudServicesResponse *orchestrator.ListCloudServicesResponse
	)

	// 1st case: Empty service ID error
	_, err = service.RemoveCloudService(context.Background(), &orchestrator.RemoveCloudServiceRequest{ServiceId: ""})
	assert.NotNil(t, err)
	assert.Equal(t, status.Code(err), codes.InvalidArgument)

	// 2nd case: ErrRecordNotFound
	_, err = service.RemoveCloudService(context.Background(), &orchestrator.RemoveCloudServiceRequest{ServiceId: DefaultTargetCloudServiceId})
	assert.NotNil(t, err)
	assert.Equal(t, status.Code(err), codes.NotFound)

	// 3rd case: Record removed successfully
	cloudServiceResponse, err = service.CreateDefaultTargetCloudService()
	assert.Nil(t, err)
	assert.NotNil(t, cloudServiceResponse)

	// There is a record for cloud services in the DB (default one)
	listCloudServicesResponse, err = service.ListCloudServices(context.Background(), &orchestrator.ListCloudServicesRequest{})
	assert.Nil(t, err)
	assert.NotNil(t, listCloudServicesResponse.Services)
	assert.NotEmpty(t, listCloudServicesResponse.Services)

	// Remove record
	_, err = service.RemoveCloudService(context.Background(), &orchestrator.RemoveCloudServiceRequest{ServiceId: DefaultTargetCloudServiceId})
	assert.Nil(t, err)

	// There is a record for cloud services in the DB (default one)
	listCloudServicesResponse, err = service.ListCloudServices(context.Background(), &orchestrator.ListCloudServicesRequest{})
	assert.Nil(t, err)
	assert.NotNil(t, listCloudServicesResponse.Services)
	assert.Empty(t, listCloudServicesResponse.Services)

	// Reset DB so that other tests work correctly
	err = service.db.Delete(&orchestrator.CloudService{}, "")
	assert.Nil(t, err)

	// Reset DB
	assert.Nil(t, gormX.Reset())
}

func TestService_CreateDefaultTargetCloudService(t *testing.T) {
	var (
		cloudServiceResponse *orchestrator.CloudService
		err                  error
	)

	// 1st case: No records for cloud services -> Default target service is created
	cloudServiceResponse, err = service.CreateDefaultTargetCloudService()
	assert.Nil(t, err)
	assert.Equal(t, &orchestrator.CloudService{
		Id:          DefaultTargetCloudServiceId,
		Name:        DefaultTargetCloudServiceName,
		Description: DefaultTargetCloudServiceDescription,
	}, cloudServiceResponse)

	// 2nd case: There is already a record for service (the default target service) -> Nothing added and no error
	cloudServiceResponse, err = service.CreateDefaultTargetCloudService()
	assert.Nil(t, err)
	assert.Nil(t, cloudServiceResponse)

	// Reset DB so that other tests work correctly
	err = service.db.Delete(&orchestrator.CloudService{}, "")
	assert.Nil(t, err)

	// Reset DB
	assert.Nil(t, gormX.Reset())
}
