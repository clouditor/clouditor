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

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/service"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestService_RegisterCloudService(t *testing.T) {
	tests := []struct {
		name    string
		req     *orchestrator.RegisterCloudServiceRequest
		res     *orchestrator.CloudService
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "missing request",
			req:  nil,
			res:  nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, api.ErrEmptyRequest.Error()) &&
					assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "missing service",
			req:  &orchestrator.RegisterCloudServiceRequest{},
			res:  nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "cloud_service: value is required") &&
					assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "missing service name",
			req:  &orchestrator.RegisterCloudServiceRequest{CloudService: &orchestrator.CloudService{}},
			res:  nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "cloud_service.name: value length must be at least 1 characters") &&
					assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Happy path: without metadata as input",
			req: &orchestrator.RegisterCloudServiceRequest{
				CloudService: &orchestrator.CloudService{
					Name:        "test",
					Description: "some",
				},
			},
			res: &orchestrator.CloudService{
				Name:        "test",
				Description: "some",
				Metadata:    &orchestrator.CloudService_Metadata{},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: with metadata as input",
			req: &orchestrator.RegisterCloudServiceRequest{
				CloudService: &orchestrator.CloudService{
					Name:        "test",
					Description: "some",
					Metadata: &orchestrator.CloudService_Metadata{
						Labels: map[string]string{
							"owner": "testOwner",
							"env":   "prod",
						},
					},
				},
			},
			res: &orchestrator.CloudService{
				Name:        "test",
				Description: "some",
				Metadata: &orchestrator.CloudService_Metadata{
					Labels: map[string]string{
						"owner": "testOwner",
						"env":   "prod",
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	orchestratorService := NewService()
	cloudService, err := orchestratorService.CreateDefaultTargetCloudService()
	assert.NoError(t, err)
	assert.NotNil(t, cloudService)
	assert.NoError(t, api.Validate(cloudService))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := orchestratorService.RegisterCloudService(context.Background(), tt.req)
			tt.wantErr(t, err)

			if tt.res != nil {
				assert.NotEmpty(t, res.Id)
			}

			// reset the IDs because we cannot compare them, since they are randomly generated
			if res != nil {
				assert.NoError(t, api.Validate(res))

				res.Id = ""
				// check creation/update time and reset
				assert.NotEmpty(t, res.CreatedAt)
				res.CreatedAt = nil

				assert.NotEmpty(t, res.UpdatedAt)
				res.UpdatedAt = nil
			}
			if tt.res != nil {
				tt.res.Id = ""
			}

			assert.Equal(t, tt.res, res)
		})
	}
}

func TestService_GetCloudService(t *testing.T) {
	tests := []struct {
		name    string
		svc     *Service
		ctx     context.Context
		req     *orchestrator.GetCloudServiceRequest
		res     *orchestrator.CloudService
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "invalid request",
			svc:  NewService(),
			ctx:  context.Background(),
			req:  nil,
			res:  nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, api.ErrEmptyRequest.Error()) &&
					assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "cloud service not found",
			svc:  NewService(),
			ctx:  context.Background(),
			req:  &orchestrator.GetCloudServiceRequest{CloudServiceId: testdata.MockCloudServiceID1},
			res:  nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "service not found") &&
					assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "valid",
			svc:  NewService(),
			ctx:  context.Background(),
			req:  &orchestrator.GetCloudServiceRequest{CloudServiceId: DefaultTargetCloudServiceId},
			res: &orchestrator.CloudService{
				Id:          DefaultTargetCloudServiceId,
				Name:        DefaultTargetCloudServiceName,
				Description: DefaultTargetCloudServiceDescription,
			},
			wantErr: assert.NoError,
		},
		{
			name: "permission denied",
			svc:  NewService(WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(false, testdata.MockCloudServiceID1))),
			ctx:  context.TODO(),
			req:  &orchestrator.GetCloudServiceRequest{CloudServiceId: DefaultTargetCloudServiceId},
			res:  nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, service.ErrPermissionDenied.Error()) &&
					assert.Equal(t, codes.PermissionDenied, status.Code(err))
			},
		},
		{
			name: "permission granted",
			svc: NewService(WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(false, testdata.MockCloudServiceID1)), WithStorage(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
				_ = s.Create(&orchestrator.CloudService{
					Id:        testdata.MockCloudServiceID1,
					Name:      "service1",
					CreatedAt: timestamppb.Now(),
					UpdatedAt: timestamppb.Now(),
				})
			}))),
			ctx: context.TODO(),
			req: &orchestrator.GetCloudServiceRequest{CloudServiceId: testdata.MockCloudServiceID1},
			res: &orchestrator.CloudService{
				Id:   testdata.MockCloudServiceID1,
				Name: "service1",
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.svc.CreateDefaultTargetCloudService()
			assert.NoError(t, err)

			res, err := tt.svc.GetCloudService(tt.ctx, tt.req)
			tt.wantErr(t, err)

			if tt.res != nil {
				assert.NoError(t, api.Validate(res))
				assert.NotEmpty(t, res.Id)
				// Check if timestamps are set and then delete for further checking
				assert.NotEmpty(t, res.CreatedAt)
				assert.NotEmpty(t, res.UpdatedAt)
				res.CreatedAt = nil
				res.UpdatedAt = nil
			}

			assert.Equal(t, tt.res, res)
		})
	}
}

func TestService_UpdateCloudService(t *testing.T) {
	var (
		cloudService *orchestrator.CloudService
		err          error
	)
	orchestratorService := NewService(WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(false, testdata.MockCloudServiceID1)))

	// 1st case: Service is nil
	_, err = orchestratorService.UpdateCloudService(context.TODO(), &orchestrator.UpdateCloudServiceRequest{})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// 2nd case: Service ID is nil
	_, err = orchestratorService.UpdateCloudService(context.TODO(), &orchestrator.UpdateCloudServiceRequest{
		CloudService: &orchestrator.CloudService{},
	})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// 3rd case: Service not found since there are no services yet
	_, err = orchestratorService.UpdateCloudService(context.TODO(), &orchestrator.UpdateCloudServiceRequest{
		CloudService: &orchestrator.CloudService{
			Id:          testdata.MockCloudServiceID1,
			Name:        DefaultTargetCloudServiceName,
			Description: DefaultTargetCloudServiceDescription,
		},
	})
	assert.Equal(t, codes.NotFound, status.Code(err))

	// 4th case: Service updated successfully
	err = orchestratorService.storage.Create(&orchestrator.CloudService{
		Id:          testdata.MockCloudServiceID1,
		Name:        DefaultTargetCloudServiceName,
		Description: DefaultTargetCloudServiceDescription,
	})
	assert.NoError(t, err)
	if err != nil {
		return
	}
	cloudService, err = orchestratorService.UpdateCloudService(context.TODO(), &orchestrator.UpdateCloudServiceRequest{
		CloudService: &orchestrator.CloudService{
			Id:          testdata.MockCloudServiceID1,
			Name:        "NewName",
			Description: "",
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, cloudService)
	assert.NoError(t, api.Validate(cloudService))
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
	// Check timestamps and delete it for further tests
	assert.NotEmpty(t, cloudServiceResponse.CreatedAt)
	assert.NotEmpty(t, cloudServiceResponse.UpdatedAt)
	cloudServiceResponse.CreatedAt = nil
	cloudServiceResponse.UpdatedAt = nil

	assert.Equal(t, &orchestrator.CloudService{
		Id:          DefaultTargetCloudServiceId,
		Name:        DefaultTargetCloudServiceName,
		Description: DefaultTargetCloudServiceDescription,
	}, cloudServiceResponse)

	// Check if CloudService is valid
	assert.NoError(t, api.Validate(cloudServiceResponse))

	// 2nd case: There is already a record for service (the default target service) -> Nothing added and no error
	cloudServiceResponse, err = orchestratorService.CreateDefaultTargetCloudService()
	assert.NoError(t, err)
	assert.Nil(t, cloudServiceResponse)
}

func TestService_ListCloudServices(t *testing.T) {
	type fields struct {
		AssessmentResultHooks []assessment.ResultHookFunc
		storage               persistence.Storage
		metricsFile           string
		loadMetricsFunc       func() ([]*assessment.Metric, error)
		catalogsFolder        string
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
				authz:   servicetest.NewAuthorizationStrategy(true),
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
				authz: servicetest.NewAuthorizationStrategy(true),
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
			name: "retrieve only allowed cloud services: no cloud service is allowed",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					// Store two cloud services, of which none we are allowed to retrieve in the test
					_ = s.Create(&orchestrator.CloudService{
						Id:   testdata.MockCloudServiceID1,
						Name: testdata.MockCloudServiceName1,
					})
					_ = s.Create(&orchestrator.CloudService{
						Id:   testdata.MockCloudServiceID2,
						Name: testdata.MockCloudServiceName2,
					})
				}),
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockCloudServiceID1),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.ListCloudServicesRequest{},
			},
			wantRes: &orchestrator.ListCloudServicesResponse{
				Services: []*orchestrator.CloudService{
					{
						Id:   testdata.MockCloudServiceID1,
						Name: testdata.MockCloudServiceName1,
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "retrieve only allowed cloud services",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					// Store two cloud services, of which only one we are allowed to retrieve in the test
					_ = s.Create(&orchestrator.CloudService{
						Id:   testdata.MockCloudServiceID1,
						Name: testdata.MockCloudServiceName1,
					})
					_ = s.Create(&orchestrator.CloudService{
						Id:   testdata.MockCloudServiceID2,
						Name: testdata.MockCloudServiceName1,
					})
				}),
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockCloudServiceID1),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.ListCloudServicesRequest{},
			},
			wantRes: &orchestrator.ListCloudServicesResponse{
				Services: []*orchestrator.CloudService{
					{
						Id:   testdata.MockCloudServiceID1,
						Name: testdata.MockCloudServiceName1,
					},
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				metricsFile:           tt.fields.metricsFile,
				loadMetricsFunc:       tt.fields.loadMetricsFunc,
				catalogsFolder:        tt.fields.catalogsFolder,
				loadCatalogsFunc:      tt.fields.loadCatalogsFunc,
				events:                tt.fields.events,
				authz:                 tt.fields.authz,
			}

			gotRes, err := svc.ListCloudServices(tt.args.ctx, tt.args.req)
			assert.NoError(t, api.Validate(gotRes))

			tt.wantErr(t, err, tt.args)
			assert.Equal(t, tt.wantRes, gotRes)
		})
	}
}

func TestService_GetCloudServiceStatistics(t *testing.T) {
	type fields struct {
		cloudServiceHooks     []orchestrator.CloudServiceHookFunc
		toeHooks              []orchestrator.TargetOfEvaluationHookFunc
		AssessmentResultHooks []assessment.ResultHookFunc
		storage               persistence.Storage
		metricsFile           string
		loadMetricsFunc       func() ([]*assessment.Metric, error)
		catalogsFolder        string
		loadCatalogsFunc      func() ([]*orchestrator.Catalog, error)
		events                chan *orchestrator.MetricChangeEvent
		authz                 service.AuthorizationStrategy
	}
	type args struct {
		ctx context.Context
		req *orchestrator.GetCloudServiceStatisticsRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *orchestrator.GetCloudServiceStatisticsResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Validate request error",
			args: args{
				ctx: context.TODO(),
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, api.ErrEmptyRequest.Error())
			},
		},
		{
			name: "Permission denied",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockCloudServiceID2),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.GetCloudServiceStatisticsRequest{
					CloudServiceId: testdata.MockCloudServiceID1,
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, service.ErrPermissionDenied.Error()) &&
					assert.Equal(t, codes.PermissionDenied, status.Code(err))
			},
		},
		{
			name: "Storage error: Get CloudService 'service not found'",
			fields: fields{
				authz:   servicetest.NewAuthorizationStrategy(false, testdata.MockCloudServiceID1),
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {}),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.GetCloudServiceStatisticsRequest{
					CloudServiceId: testdata.MockCloudServiceID1,
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "service not found") &&
					assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "Storage error: Get CloudService",
			fields: fields{
				authz:   servicetest.NewAuthorizationStrategy(false, testdata.MockCloudServiceID1),
				storage: &testutil.StorageWithError{GetErr: ErrSomeError},
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.GetCloudServiceStatisticsRequest{
					CloudServiceId: testdata.MockCloudServiceID1,
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "database error getting cloud service: some error") &&
					assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
		{
			name: "Happy path",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					// Store one cloud services
					_ = s.Create(&orchestrator.CloudService{
						Id:          testdata.MockCloudServiceID1,
						Name:        testdata.MockCloudServiceName1,
						Description: testdata.MockCloudServiceDescription1,
						CatalogsInScope: []*orchestrator.Catalog{
							{
								Id:          testdata.MockCatalogID,
								Name:        testdata.MockCatalogName,
								Description: testdata.MockCatalogDescription,
							},
						},
					})
					_ = s.Create(&orchestrator.CloudService{
						Id:   testdata.MockCloudServiceID2,
						Name: testdata.MockCloudServiceName2,
					})
					_ = s.Create(&evidence.Evidence{
						Id:             uuid.NewString(),
						CloudServiceId: testdata.MockCloudServiceID1,
					})
					_ = s.Create(&evidence.Evidence{
						Id:             uuid.NewString(),
						CloudServiceId: testdata.MockCloudServiceID2,
					})
					_ = s.Create(&assessment.AssessmentResult{
						Id:             uuid.NewString(),
						CloudServiceId: testdata.MockCloudServiceID1,
					})
					_ = s.Create(&assessment.AssessmentResult{
						Id:             uuid.NewString(),
						CloudServiceId: testdata.MockCloudServiceID1,
					})
					_ = s.Create(&assessment.AssessmentResult{
						Id:             uuid.NewString(),
						CloudServiceId: testdata.MockCloudServiceID2,
					})
					_ = s.Create(&discovery.Resource{
						Id:             uuid.NewString(),
						CloudServiceId: testdata.MockCloudServiceID1,
					})
					_ = s.Create(&discovery.Resource{
						Id:             uuid.NewString(),
						CloudServiceId: testdata.MockCloudServiceID2,
					})
				}),
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockCloudServiceID1),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.GetCloudServiceStatisticsRequest{
					CloudServiceId: testdata.MockCloudServiceID1,
				},
			},
			wantRes: &orchestrator.GetCloudServiceStatisticsResponse{
				NumberOfDiscoveredResources: 1,
				NumberOfAssessmentResults:   2,
				NumberOfEvidences:           1,
				NumberOfSelectedCatalogs:    1,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				cloudServiceHooks:     tt.fields.cloudServiceHooks,
				toeHooks:              tt.fields.toeHooks,
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				metricsFile:           tt.fields.metricsFile,
				loadMetricsFunc:       tt.fields.loadMetricsFunc,
				catalogsFolder:        tt.fields.catalogsFolder,
				loadCatalogsFunc:      tt.fields.loadCatalogsFunc,
				events:                tt.fields.events,
				authz:                 tt.fields.authz,
			}
			gotRes, err := s.GetCloudServiceStatistics(tt.args.ctx, tt.args.req)
			tt.wantErr(t, err)
			assert.Equal(t, tt.wantRes, gotRes)
		})
	}
}
