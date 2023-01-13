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

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/service"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func TestService_RegisterCloudService(t *testing.T) {
	tests := []struct {
		name string
		req  *orchestrator.RegisterCloudServiceRequest
		res  *orchestrator.CloudService
		err  error
	}{
		{
			"missing request",
			nil,
			nil,
			status.Error(codes.InvalidArgument, "empty request"),
		},
		{
			"missing service",
			&orchestrator.RegisterCloudServiceRequest{},
			nil,
			status.Error(codes.InvalidArgument, "invalid request: invalid RegisterCloudServiceRequest.CloudService: value is required"),
		},
		{
			"missing service name",
			&orchestrator.RegisterCloudServiceRequest{CloudService: &orchestrator.CloudService{}},
			nil,
			status.Error(codes.InvalidArgument, "invalid request: invalid RegisterCloudServiceRequest.CloudService: embedded message failed validation | caused by: invalid CloudService.Name: value length must be at least 1 runes"),
		},
		{
			"valid",
			&orchestrator.RegisterCloudServiceRequest{CloudService: &orchestrator.CloudService{Name: "test", Description: "some"}},
			&orchestrator.CloudService{Name: "test", Description: "some"},
			nil,
		},
	}
	orchestratorService := NewService()
	cloudService, err := orchestratorService.CreateDefaultTargetCloudService()
	assert.NoError(t, err)
	assert.NotNil(t, cloudService)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := orchestratorService.RegisterCloudService(context.Background(), tt.req)

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
}

func TestService_GetCloudService(t *testing.T) {
	tests := []struct {
		name string
		svc  *Service
		ctx  context.Context
		req  *orchestrator.GetCloudServiceRequest
		res  *orchestrator.CloudService
		err  error
	}{
		{
			"invalid request",
			NewService(),
			context.Background(),
			nil,
			nil,
			status.Error(codes.InvalidArgument, "empty request"),
		},
		{
			"cloud service not found",
			NewService(),
			context.Background(),
			&orchestrator.GetCloudServiceRequest{CloudServiceId: testdata.MockCloudServiceID},
			nil,
			status.Error(codes.NotFound, "service not found"),
		},
		{
			"valid",
			NewService(),
			context.Background(),
			&orchestrator.GetCloudServiceRequest{CloudServiceId: DefaultTargetCloudServiceId},
			&orchestrator.CloudService{
				Id:          DefaultTargetCloudServiceId,
				Name:        DefaultTargetCloudServiceName,
				Description: DefaultTargetCloudServiceDescription,
			},
			nil,
		},
		{
			"permission denied",
			NewService(WithAuthorizationStrategyJWT(testutil.TestCustomClaims)),
			testutil.TestContextOnlyService1,
			&orchestrator.GetCloudServiceRequest{CloudServiceId: DefaultTargetCloudServiceId},
			nil,
			service.ErrPermissionDenied,
		},
		{
			"permission granted",
			NewService(WithAuthorizationStrategyJWT(testutil.TestCustomClaims), WithStorage(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
				_ = s.Create(&orchestrator.CloudService{
					Id:   testutil.TestCloudService1,
					Name: "service1",
				})
			}))),
			testutil.TestContextOnlyService1,
			&orchestrator.GetCloudServiceRequest{CloudServiceId: testutil.TestCloudService1},
			&orchestrator.CloudService{
				Id:   testutil.TestCloudService1,
				Name: "service1",
			},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.svc.CreateDefaultTargetCloudService()
			assert.NoError(t, err)

			res, err := tt.svc.GetCloudService(tt.ctx, tt.req)

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
}

func TestService_UpdateCloudService(t *testing.T) {
	var (
		cloudService *orchestrator.CloudService
		err          error
	)
	orchestratorService := NewService(WithAuthorizationStrategyJWT(testutil.TestCustomClaims))

	// 1st case: Service is nil
	_, err = orchestratorService.UpdateCloudService(testutil.TestContextOnlyService1, &orchestrator.UpdateCloudServiceRequest{})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// 2nd case: Service ID is nil
	_, err = orchestratorService.UpdateCloudService(testutil.TestContextOnlyService1, &orchestrator.UpdateCloudServiceRequest{
		CloudService: &orchestrator.CloudService{},
	})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// 3rd case: Service not found since there are no services yet
	_, err = orchestratorService.UpdateCloudService(testutil.TestContextOnlyService1, &orchestrator.UpdateCloudServiceRequest{
		CloudService: &orchestrator.CloudService{
			Id:          testutil.TestCloudService1,
			Name:        DefaultTargetCloudServiceName,
			Description: DefaultTargetCloudServiceDescription,
		},
	})
	assert.Equal(t, codes.NotFound, status.Code(err))

	// 4th case: Service updated successfully
	err = orchestratorService.storage.Create(&orchestrator.CloudService{
		Id:          testutil.TestCloudService1,
		Name:        DefaultTargetCloudServiceName,
		Description: DefaultTargetCloudServiceDescription,
	})
	assert.NoError(t, err)
	if err != nil {
		return
	}
	cloudService, err = orchestratorService.UpdateCloudService(testutil.TestContextOnlyService1, &orchestrator.UpdateCloudServiceRequest{
		CloudService: &orchestrator.CloudService{
			Id:          testutil.TestCloudService1,
			Name:        "NewName",
			Description: "",
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, cloudService)
	assert.Equal(t, "NewName", cloudService.Name)
	// Description should be overwritten with empty string
	assert.Equal(t, "", cloudService.Description)
}

func TestService_RemoveCloudService(t *testing.T) {
	var (
		cloudServiceResponse      *orchestrator.CloudService
		err                       error
		listCloudServicesResponse *orchestrator.ListCloudServicesResponse
	)
	orchestratorService := NewService()

	// 1st case: Empty service ID error
	_, err = orchestratorService.RemoveCloudService(context.Background(), &orchestrator.RemoveCloudServiceRequest{CloudServiceId: ""})
	assert.Error(t, err)
	assert.Equal(t, status.Code(err), codes.InvalidArgument)

	// 2nd case: ErrRecordNotFound
	_, err = orchestratorService.RemoveCloudService(context.Background(), &orchestrator.RemoveCloudServiceRequest{CloudServiceId: DefaultTargetCloudServiceId})
	assert.Error(t, err)
	assert.Equal(t, status.Code(err), codes.NotFound)

	// 3rd case: Record removed successfully
	cloudServiceResponse, err = orchestratorService.CreateDefaultTargetCloudService()
	assert.NoError(t, err)
	assert.NotNil(t, cloudServiceResponse)

	// There is a record for cloud services in the DB (default one)
	listCloudServicesResponse, err = orchestratorService.ListCloudServices(context.Background(), &orchestrator.ListCloudServicesRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listCloudServicesResponse.Services)
	assert.NotEmpty(t, listCloudServicesResponse.Services)

	// Remove record
	_, err = orchestratorService.RemoveCloudService(context.Background(), &orchestrator.RemoveCloudServiceRequest{CloudServiceId: DefaultTargetCloudServiceId})
	assert.NoError(t, err)

	// There is a record for cloud services in the DB (default one)
	listCloudServicesResponse, err = orchestratorService.ListCloudServices(context.Background(), &orchestrator.ListCloudServicesRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listCloudServicesResponse.Services)
	assert.Empty(t, listCloudServicesResponse.Services)
}

func TestService_CreateDefaultTargetCloudService(t *testing.T) {
	var (
		cloudServiceResponse *orchestrator.CloudService
		err                  error
	)
	orchestratorService := NewService()

	// 1st case: No records for cloud services -> Default target service is created
	cloudServiceResponse, err = orchestratorService.CreateDefaultTargetCloudService()
	assert.NoError(t, err)
	assert.Equal(t, &orchestrator.CloudService{
		Id:          DefaultTargetCloudServiceId,
		Name:        DefaultTargetCloudServiceName,
		Description: DefaultTargetCloudServiceDescription,
	}, cloudServiceResponse)

	// 2nd case: There is already a record for service (the default target service) -> Nothing added and no error
	cloudServiceResponse, err = orchestratorService.CreateDefaultTargetCloudService()
	assert.NoError(t, err)
	assert.Nil(t, cloudServiceResponse)
}

func TestService_ListCloudServices(t *testing.T) {
	type fields struct {
		results               map[string]*assessment.AssessmentResult
		AssessmentResultHooks []func(result *assessment.AssessmentResult, err error)
		storage               persistence.Storage
		metricsFile           string
		loadMetricsFunc       func() ([]*assessment.Metric, error)
		catalogsFile          string
		loadCatalogsFunc      func() ([]*orchestrator.Catalog, error)
		events                chan *orchestrator.MetricChangeEvent
		authz                 service.AuthorizationStrategy
	}
	type args struct {
		ctx context.Context
		req *orchestrator.ListCloudServicesRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *orchestrator.ListCloudServicesResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "retrieve empty list",
			args: args{req: &orchestrator.ListCloudServicesRequest{}},
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   &service.AuthorizationStrategyAllowAll{},
			},
			wantRes: &orchestrator.ListCloudServicesResponse{},
			wantErr: assert.NoError,
		},
		{
			name: "list with one item",
			args: args{req: &orchestrator.ListCloudServicesRequest{}},
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					service := &orchestrator.CloudService{
						Id:          DefaultTargetCloudServiceId,
						Name:        DefaultTargetCloudServiceName,
						Description: DefaultTargetCloudServiceDescription,
					}

					_ = s.Create(service)
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			wantRes: &orchestrator.ListCloudServicesResponse{
				Services: []*orchestrator.CloudService{
					{
						Id:          DefaultTargetCloudServiceId,
						Name:        DefaultTargetCloudServiceName,
						Description: DefaultTargetCloudServiceDescription,
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "invalid request",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				req: &orchestrator.ListCloudServicesRequest{
					OrderBy: "not a field",
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, api.ErrInvalidColumnName.Error())
			},
		},
		{
			name: "retrieve only allowed cloud services",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					// Store two cloud services, of which only one we are allowed to retrieve in the test
					_ = s.Create(&orchestrator.CloudService{Id: testdata.MockCloudServiceID})
					_ = s.Create(&orchestrator.CloudService{Id: testdata.MockAnotherCloudServiceID})
				}),
				authz: &service.AuthorizationStrategyJWT{Key: testutil.TestCustomClaims},
			},
			args: args{
				ctx: testutil.TestContextOnlyService1,
				req: &orchestrator.ListCloudServicesRequest{},
			},
			wantRes: &orchestrator.ListCloudServicesResponse{
				Services: []*orchestrator.CloudService{
					{Id: testdata.MockCloudServiceID},
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				results:               tt.fields.results,
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				metricsFile:           tt.fields.metricsFile,
				loadMetricsFunc:       tt.fields.loadMetricsFunc,
				catalogsFile:          tt.fields.catalogsFile,
				loadCatalogsFunc:      tt.fields.loadCatalogsFunc,
				events:                tt.fields.events,
				authz:                 tt.fields.authz,
			}

			gotRes, err := svc.ListCloudServices(tt.args.ctx, tt.args.req)
			tt.wantErr(t, err, tt.args)

			if !proto.Equal(gotRes, tt.wantRes) {
				t.Errorf("Service.ListCloudServices() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
