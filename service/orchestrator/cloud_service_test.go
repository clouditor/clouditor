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

func TestService_RegisterCertificationTarget(t *testing.T) {
	tests := []struct {
		name    string
		req     *orchestrator.RegisterCertificationTargetRequest
		res     *orchestrator.CertificationTarget
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
			req:  &orchestrator.RegisterCertificationTargetRequest{},
			res:  nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "certification_target: value is required") &&
					assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "missing service name",
			req:  &orchestrator.RegisterCertificationTargetRequest{CertificationTarget: &orchestrator.CertificationTarget{}},
			res:  nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "certification_target.name: value length must be at least 1 characters") &&
					assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Happy path: without metadata as input",
			req: &orchestrator.RegisterCertificationTargetRequest{
				CertificationTarget: &orchestrator.CertificationTarget{
					Name:        "test",
					Description: "some",
				},
			},
			res: &orchestrator.CertificationTarget{
				Name:        "test",
				Description: "some",
				Metadata:    &orchestrator.CertificationTarget_Metadata{},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: with metadata as input",
			req: &orchestrator.RegisterCertificationTargetRequest{
				CertificationTarget: &orchestrator.CertificationTarget{
					Name:        "test",
					Description: "some",
					Metadata: &orchestrator.CertificationTarget_Metadata{
						Labels: map[string]string{
							"owner": "testOwner",
							"env":   "prod",
						},
					},
				},
			},
			res: &orchestrator.CertificationTarget{
				Name:        "test",
				Description: "some",
				Metadata: &orchestrator.CertificationTarget_Metadata{
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
	CertificationTarget, err := orchestratorService.CreateDefaultCertificationTarget()
	assert.NoError(t, err)
	assert.NotNil(t, CertificationTarget)
	assert.NoError(t, api.Validate(CertificationTarget))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := orchestratorService.RegisterCertificationTarget(context.Background(), tt.req)
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

func TestService_GetCertificationTarget(t *testing.T) {
	tests := []struct {
		name    string
		svc     *Service
		ctx     context.Context
		req     *orchestrator.GetCertificationTargetRequest
		res     *orchestrator.CertificationTarget
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
			name: "certification target not found",
			svc:  NewService(),
			ctx:  context.Background(),
			req:  &orchestrator.GetCertificationTargetRequest{CertificationTargetId: testdata.MockCertificationTargetID1},
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
			req:  &orchestrator.GetCertificationTargetRequest{CertificationTargetId: DefaultTargetCertificationTargetId},
			res: &orchestrator.CertificationTarget{
				Id:          DefaultTargetCertificationTargetId,
				Name:        DefaultTargetCertificationTargetName,
				Description: DefaultTargetCertificationTargetDescription,
			},
			wantErr: assert.NoError,
		},
		{
			name: "permission denied",
			svc:  NewService(WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(false, testdata.MockCertificationTargetID1))),
			ctx:  context.TODO(),
			req:  &orchestrator.GetCertificationTargetRequest{CertificationTargetId: DefaultTargetCertificationTargetId},
			res:  nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, service.ErrPermissionDenied.Error()) &&
					assert.Equal(t, codes.PermissionDenied, status.Code(err))
			},
		},
		{
			name: "permission granted",
			svc: NewService(WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(false, testdata.MockCertificationTargetID1)), WithStorage(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
				_ = s.Create(&orchestrator.CertificationTarget{
					Id:        testdata.MockCertificationTargetID1,
					Name:      "service1",
					CreatedAt: timestamppb.Now(),
					UpdatedAt: timestamppb.Now(),
				})
			}))),
			ctx: context.TODO(),
			req: &orchestrator.GetCertificationTargetRequest{CertificationTargetId: testdata.MockCertificationTargetID1},
			res: &orchestrator.CertificationTarget{
				Id:   testdata.MockCertificationTargetID1,
				Name: "service1",
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.svc.CreateDefaultCertificationTarget()
			assert.NoError(t, err)

			res, err := tt.svc.GetCertificationTarget(tt.ctx, tt.req)
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

func TestService_UpdateCertificationTarget(t *testing.T) {
	var (
		CertificationTarget *orchestrator.CertificationTarget
		err                 error
	)
	orchestratorService := NewService(WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(false, testdata.MockCertificationTargetID1)))

	// 1st case: Service is nil
	_, err = orchestratorService.UpdateCertificationTarget(context.TODO(), &orchestrator.UpdateCertificationTargetRequest{})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// 2nd case: Service ID is nil
	_, err = orchestratorService.UpdateCertificationTarget(context.TODO(), &orchestrator.UpdateCertificationTargetRequest{
		CertificationTarget: &orchestrator.CertificationTarget{},
	})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// 3rd case: Service not found since there are no services yet
	_, err = orchestratorService.UpdateCertificationTarget(context.TODO(), &orchestrator.UpdateCertificationTargetRequest{
		CertificationTarget: &orchestrator.CertificationTarget{
			Id:          testdata.MockCertificationTargetID1,
			Name:        DefaultTargetCertificationTargetName,
			Description: DefaultTargetCertificationTargetDescription,
		},
	})
	assert.Equal(t, codes.NotFound, status.Code(err))

	// 4th case: Service updated successfully
	err = orchestratorService.storage.Create(&orchestrator.CertificationTarget{
		Id:          testdata.MockCertificationTargetID1,
		Name:        DefaultTargetCertificationTargetName,
		Description: DefaultTargetCertificationTargetDescription,
	})
	assert.NoError(t, err)
	if err != nil {
		return
	}
	CertificationTarget, err = orchestratorService.UpdateCertificationTarget(context.TODO(), &orchestrator.UpdateCertificationTargetRequest{
		CertificationTarget: &orchestrator.CertificationTarget{
			Id:          testdata.MockCertificationTargetID1,
			Name:        "NewName",
			Description: "",
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, CertificationTarget)
	assert.NoError(t, api.Validate(CertificationTarget))
	assert.Equal(t, "NewName", CertificationTarget.Name)
	// Description should be overwritten with empty string
	assert.Equal(t, "", CertificationTarget.Description)
}

func TestService_RemoveCertificationTarget(t *testing.T) {
	var (
		CertificationTargetResponse      *orchestrator.CertificationTarget
		err                              error
		listCertificationTargetsResponse *orchestrator.ListCertificationTargetsResponse
	)
	orchestratorService := NewService()

	// 1st case: Empty service ID error
	_, err = orchestratorService.RemoveCertificationTarget(context.Background(), &orchestrator.RemoveCertificationTargetRequest{CertificationTargetId: ""})
	assert.Error(t, err)
	assert.Equal(t, status.Code(err), codes.InvalidArgument)

	// 2nd case: ErrRecordNotFound
	_, err = orchestratorService.RemoveCertificationTarget(context.Background(), &orchestrator.RemoveCertificationTargetRequest{CertificationTargetId: DefaultTargetCertificationTargetId})
	assert.Error(t, err)
	assert.Equal(t, status.Code(err), codes.NotFound)

	// 3rd case: Record removed successfully
	CertificationTargetResponse, err = orchestratorService.CreateDefaultCertificationTarget()
	assert.NoError(t, err)
	assert.NotNil(t, CertificationTargetResponse)

	// There is a record for certification targets in the DB (default one)
	listCertificationTargetsResponse, err = orchestratorService.ListCertificationTargets(context.Background(), &orchestrator.ListCertificationTargetsRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listCertificationTargetsResponse.Services)
	assert.NotEmpty(t, listCertificationTargetsResponse.Services)

	// Remove record
	_, err = orchestratorService.RemoveCertificationTarget(context.Background(), &orchestrator.RemoveCertificationTargetRequest{CertificationTargetId: DefaultTargetCertificationTargetId})
	assert.NoError(t, err)

	// There is a record for certification targets in the DB (default one)
	listCertificationTargetsResponse, err = orchestratorService.ListCertificationTargets(context.Background(), &orchestrator.ListCertificationTargetsRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listCertificationTargetsResponse.Services)
	assert.Empty(t, listCertificationTargetsResponse.Services)
}

func TestService_CreateDefaultTargetCertificationTarget(t *testing.T) {
	var (
		CertificationTargetResponse *orchestrator.CertificationTarget
		err                         error
	)
	orchestratorService := NewService()

	// 1st case: No records for certification targets -> Default target service is created
	CertificationTargetResponse, err = orchestratorService.CreateDefaultCertificationTarget()
	assert.NoError(t, err)
	// Check timestamps and delete it for further tests
	assert.NotEmpty(t, CertificationTargetResponse.CreatedAt)
	assert.NotEmpty(t, CertificationTargetResponse.UpdatedAt)
	CertificationTargetResponse.CreatedAt = nil
	CertificationTargetResponse.UpdatedAt = nil

	assert.Equal(t, &orchestrator.CertificationTarget{
		Id:          DefaultTargetCertificationTargetId,
		Name:        DefaultTargetCertificationTargetName,
		Description: DefaultTargetCertificationTargetDescription,
	}, CertificationTargetResponse)

	// Check if CertificationTarget is valid
	assert.NoError(t, api.Validate(CertificationTargetResponse))

	// 2nd case: There is already a record for service (the default target service) -> Nothing added and no error
	CertificationTargetResponse, err = orchestratorService.CreateDefaultCertificationTarget()
	assert.NoError(t, err)
	assert.Nil(t, CertificationTargetResponse)
}

func TestService_ListCertificationTargets(t *testing.T) {
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
		req *orchestrator.ListCertificationTargetsRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *orchestrator.ListCertificationTargetsResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "retrieve empty list",
			args: args{req: &orchestrator.ListCertificationTargetsRequest{}},
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			wantRes: &orchestrator.ListCertificationTargetsResponse{},
			wantErr: assert.NoError,
		},
		{
			name: "list with one item",
			args: args{req: &orchestrator.ListCertificationTargetsRequest{}},
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					service := &orchestrator.CertificationTarget{
						Id:          DefaultTargetCertificationTargetId,
						Name:        DefaultTargetCertificationTargetName,
						Description: DefaultTargetCertificationTargetDescription,
					}

					_ = s.Create(service)
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			wantRes: &orchestrator.ListCertificationTargetsResponse{
				Services: []*orchestrator.CertificationTarget{
					{
						Id:          DefaultTargetCertificationTargetId,
						Name:        DefaultTargetCertificationTargetName,
						Description: DefaultTargetCertificationTargetDescription,
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "retrieve only allowed certification targets: no certification target is allowed",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					// Store two certification targets, of which none we are allowed to retrieve in the test
					_ = s.Create(&orchestrator.CertificationTarget{
						Id:   testdata.MockCertificationTargetID1,
						Name: testdata.MockCertificationTargetName1,
					})
					_ = s.Create(&orchestrator.CertificationTarget{
						Id:   testdata.MockCertificationTargetID2,
						Name: testdata.MockCertificationTargetName2,
					})
				}),
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockCertificationTargetID1),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.ListCertificationTargetsRequest{},
			},
			wantRes: &orchestrator.ListCertificationTargetsResponse{
				Services: []*orchestrator.CertificationTarget{
					{
						Id:   testdata.MockCertificationTargetID1,
						Name: testdata.MockCertificationTargetName1,
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "retrieve only allowed certification targets",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					// Store two certification targets, of which only one we are allowed to retrieve in the test
					_ = s.Create(&orchestrator.CertificationTarget{
						Id:   testdata.MockCertificationTargetID1,
						Name: testdata.MockCertificationTargetName1,
					})
					_ = s.Create(&orchestrator.CertificationTarget{
						Id:   testdata.MockCertificationTargetID2,
						Name: testdata.MockCertificationTargetName1,
					})
				}),
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockCertificationTargetID1),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.ListCertificationTargetsRequest{},
			},
			wantRes: &orchestrator.ListCertificationTargetsResponse{
				Services: []*orchestrator.CertificationTarget{
					{
						Id:   testdata.MockCertificationTargetID1,
						Name: testdata.MockCertificationTargetName1,
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

			gotRes, err := svc.ListCertificationTargets(tt.args.ctx, tt.args.req)
			assert.NoError(t, api.Validate(gotRes))

			tt.wantErr(t, err, tt.args)
			assert.Equal(t, tt.wantRes, gotRes)
		})
	}
}

func TestService_GetCertificationTargetStatistics(t *testing.T) {
	type fields struct {
		CertificationTargetHooks []orchestrator.CertificationTargetHookFunc
		auditScopeHooks          []orchestrator.AuditScopeHookFunc
		AssessmentResultHooks    []assessment.ResultHookFunc
		storage                  persistence.Storage
		metricsFile              string
		loadMetricsFunc          func() ([]*assessment.Metric, error)
		catalogsFolder           string
		loadCatalogsFunc         func() ([]*orchestrator.Catalog, error)
		events                   chan *orchestrator.MetricChangeEvent
		authz                    service.AuthorizationStrategy
	}
	type args struct {
		ctx context.Context
		req *orchestrator.GetCertificationTargetStatisticsRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *orchestrator.GetCertificationTargetStatisticsResponse
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
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockCertificationTargetID2),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.GetCertificationTargetStatisticsRequest{
					CertificationTargetId: testdata.MockCertificationTargetID1,
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, service.ErrPermissionDenied.Error()) &&
					assert.Equal(t, codes.PermissionDenied, status.Code(err))
			},
		},
		{
			name: "Storage error: Get CertificationTarget 'service not found'",
			fields: fields{
				authz:   servicetest.NewAuthorizationStrategy(false, testdata.MockCertificationTargetID1),
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {}),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.GetCertificationTargetStatisticsRequest{
					CertificationTargetId: testdata.MockCertificationTargetID1,
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "service not found") &&
					assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "Storage error: Get CertificationTarget",
			fields: fields{
				authz:   servicetest.NewAuthorizationStrategy(false, testdata.MockCertificationTargetID1),
				storage: &testutil.StorageWithError{GetErr: ErrSomeError},
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.GetCertificationTargetStatisticsRequest{
					CertificationTargetId: testdata.MockCertificationTargetID1,
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "database error getting certification target: some error") &&
					assert.Equal(t, codes.Internal, status.Code(err))
			},
		},
		{
			name: "Happy path",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					// Store one certification targets
					_ = s.Create(&orchestrator.CertificationTarget{
						Id:          testdata.MockCertificationTargetID1,
						Name:        testdata.MockCertificationTargetName1,
						Description: testdata.MockCertificationTargetDescription1,
						CatalogsInScope: []*orchestrator.Catalog{
							{
								Id:          testdata.MockCatalogID,
								Name:        testdata.MockCatalogName,
								Description: testdata.MockCatalogDescription,
							},
						},
					})
					_ = s.Create(&orchestrator.CertificationTarget{
						Id:   testdata.MockCertificationTargetID2,
						Name: testdata.MockCertificationTargetName2,
					})
					_ = s.Create(&evidence.Evidence{
						Id:                    uuid.NewString(),
						CertificationTargetId: testdata.MockCertificationTargetID1,
					})
					_ = s.Create(&evidence.Evidence{
						Id:                    uuid.NewString(),
						CertificationTargetId: testdata.MockCertificationTargetID2,
					})
					_ = s.Create(&assessment.AssessmentResult{
						Id:                    uuid.NewString(),
						CertificationTargetId: testdata.MockCertificationTargetID1,
					})
					_ = s.Create(&assessment.AssessmentResult{
						Id:                    uuid.NewString(),
						CertificationTargetId: testdata.MockCertificationTargetID1,
					})
					_ = s.Create(&assessment.AssessmentResult{
						Id:                    uuid.NewString(),
						CertificationTargetId: testdata.MockCertificationTargetID2,
					})
					_ = s.Create(&discovery.Resource{
						Id:                    uuid.NewString(),
						CertificationTargetId: testdata.MockCertificationTargetID1,
					})
					_ = s.Create(&discovery.Resource{
						Id:                    uuid.NewString(),
						CertificationTargetId: testdata.MockCertificationTargetID2,
					})
				}),
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockCertificationTargetID1),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.GetCertificationTargetStatisticsRequest{
					CertificationTargetId: testdata.MockCertificationTargetID1,
				},
			},
			wantRes: &orchestrator.GetCertificationTargetStatisticsResponse{
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
				CertificationTargetHooks: tt.fields.CertificationTargetHooks,
				auditScopeHooks:          tt.fields.auditScopeHooks,
				AssessmentResultHooks:    tt.fields.AssessmentResultHooks,
				storage:                  tt.fields.storage,
				metricsFile:              tt.fields.metricsFile,
				loadMetricsFunc:          tt.fields.loadMetricsFunc,
				catalogsFolder:           tt.fields.catalogsFolder,
				loadCatalogsFunc:         tt.fields.loadCatalogsFunc,
				events:                   tt.fields.events,
				authz:                    tt.fields.authz,
			}
			gotRes, err := s.GetCertificationTargetStatistics(tt.args.ctx, tt.args.req)
			tt.wantErr(t, err)
			assert.Equal(t, tt.wantRes, gotRes)
		})
	}
}
