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
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/orchestratortest"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/service"
	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func init() {
	viper.Set(config.DefaultCertificationTargetNameFlag, config.DefaultCertificationTargetName)
	viper.Set(config.DefaultCertificationTargetDescriptionFlag, config.DefaultCertificationTargetDescription)
	viper.Set(config.DefaultCertificationTargetTypeFlag, int32(config.DefaultCertificationTargetType))
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
				return assert.ErrorContains(t, err, "certification target not found") &&
					assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "valid",
			svc:  NewService(),
			ctx:  context.Background(),
			req:  &orchestrator.GetCertificationTargetRequest{CertificationTargetId: DefaultCertificationTargetId},
			res: &orchestrator.CertificationTarget{
				Id:          DefaultCertificationTargetId,
				Name:        config.DefaultCertificationTargetName,
				Description: config.DefaultCertificationTargetDescription,
				TargetType:  config.DefaultCertificationTargetType,
			},
			wantErr: assert.NoError,
		},
		{
			name: "permission denied",
			svc:  NewService(WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(false, testdata.MockCertificationTargetID1))),
			ctx:  context.TODO(),
			req:  &orchestrator.GetCertificationTargetRequest{CertificationTargetId: DefaultCertificationTargetId},
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
					Name:      "target1",
					CreatedAt: timestamppb.Now(),
					UpdatedAt: timestamppb.Now(),
				})
			}))),
			ctx: context.TODO(),
			req: &orchestrator.GetCertificationTargetRequest{CertificationTargetId: testdata.MockCertificationTargetID1},
			res: &orchestrator.CertificationTarget{
				Id:   testdata.MockCertificationTargetID1,
				Name: "target1",
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
		target *orchestrator.CertificationTarget
		err    error
	)
	orchestratorService := NewService(WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(false, testdata.MockCertificationTargetID1)))

	// 1st case: Certification Target is nil
	_, err = orchestratorService.UpdateCertificationTarget(context.TODO(), &orchestrator.UpdateCertificationTargetRequest{})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// 2nd case: Certification Target ID is nil
	_, err = orchestratorService.UpdateCertificationTarget(context.TODO(), &orchestrator.UpdateCertificationTargetRequest{
		CertificationTarget: &orchestrator.CertificationTarget{},
	})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// 3rd case: Certification Target not found since there are no certification targets yet
	_, err = orchestratorService.UpdateCertificationTarget(context.TODO(), &orchestrator.UpdateCertificationTargetRequest{
		CertificationTarget: &orchestrator.CertificationTarget{
			Id:          testdata.MockCertificationTargetID1,
			Name:        config.DefaultCertificationTargetName,
			Description: config.DefaultCertificationTargetDescription,
		},
	})
	assert.Equal(t, codes.NotFound, status.Code(err))

	// 4th case: Certification Target updated successfully
	err = orchestratorService.storage.Create(&orchestrator.CertificationTarget{
		Id:          testdata.MockCertificationTargetID1,
		Name:        config.DefaultCertificationTargetName,
		Description: config.DefaultCertificationTargetDescription,
	})
	assert.NoError(t, err)
	if err != nil {
		return
	}
	target, err = orchestratorService.UpdateCertificationTarget(context.TODO(), &orchestrator.UpdateCertificationTargetRequest{
		CertificationTarget: &orchestrator.CertificationTarget{
			Id:          testdata.MockCertificationTargetID1,
			Name:        "NewName",
			Description: "",
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, target)
	assert.NoError(t, api.Validate(target))
	assert.Equal(t, "NewName", target.Name)
	// Description should be overwritten with empty string
	assert.Equal(t, "", target.Description)
}

func TestService_RemoveCertificationTarget(t *testing.T) {
	var (
		CertificationTargetResponse      *orchestrator.CertificationTarget
		err                              error
		listCertificationTargetsResponse *orchestrator.ListCertificationTargetsResponse
	)
	orchestratorService := NewService()

	// 1st case: Empty certification target ID error
	_, err = orchestratorService.RemoveCertificationTarget(context.Background(), &orchestrator.RemoveCertificationTargetRequest{CertificationTargetId: ""})
	assert.Error(t, err)
	assert.Equal(t, status.Code(err), codes.InvalidArgument)

	// 2nd case: ErrRecordNotFound
	_, err = orchestratorService.RemoveCertificationTarget(context.Background(), &orchestrator.RemoveCertificationTargetRequest{CertificationTargetId: DefaultCertificationTargetId})
	assert.Error(t, err)
	assert.Equal(t, status.Code(err), codes.NotFound)

	// 3rd case: Record removed successfully
	CertificationTargetResponse, err = orchestratorService.CreateDefaultCertificationTarget()
	assert.NoError(t, err)
	assert.NotNil(t, CertificationTargetResponse)

	// There is a record for certification targets in the DB (default one)
	listCertificationTargetsResponse, err = orchestratorService.ListCertificationTargets(context.Background(), &orchestrator.ListCertificationTargetsRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listCertificationTargetsResponse.Targets)
	assert.NotEmpty(t, listCertificationTargetsResponse.Targets)

	// Remove record
	_, err = orchestratorService.RemoveCertificationTarget(context.Background(), &orchestrator.RemoveCertificationTargetRequest{CertificationTargetId: DefaultCertificationTargetId})
	assert.NoError(t, err)

	// There is a record for certification targets in the DB (default one)
	listCertificationTargetsResponse, err = orchestratorService.ListCertificationTargets(context.Background(), &orchestrator.ListCertificationTargetsRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listCertificationTargetsResponse.Targets)
	assert.Empty(t, listCertificationTargetsResponse.Targets)
}

func TestService_CreateDefaultCertificationTarget(t *testing.T) {
	var (
		CertificationTargetResponse *orchestrator.CertificationTarget
		err                         error
	)
	orchestratorService := NewService()

	// 1st case: No records for certification targets -> Default certification target is created
	CertificationTargetResponse, err = orchestratorService.CreateDefaultCertificationTarget()
	assert.NoError(t, err)
	// Check timestamps and delete it for further tests
	assert.NotEmpty(t, CertificationTargetResponse.CreatedAt)
	assert.NotEmpty(t, CertificationTargetResponse.UpdatedAt)
	CertificationTargetResponse.CreatedAt = nil
	CertificationTargetResponse.UpdatedAt = nil

	assert.Equal(t, &orchestrator.CertificationTarget{
		Id:          DefaultCertificationTargetId,
		Name:        config.DefaultCertificationTargetName,
		Description: config.DefaultCertificationTargetDescription,
		TargetType:  config.DefaultCertificationTargetType,
	}, CertificationTargetResponse)

	// Check if CertificationTarget is valid
	assert.NoError(t, api.Validate(CertificationTargetResponse))

	// 2nd case: There is already a record for the certification target (the default certification target) -> Nothing added and no error
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
					target := &orchestrator.CertificationTarget{
						Id:          DefaultCertificationTargetId,
						Name:        config.DefaultCertificationTargetName,
						Description: config.DefaultCertificationTargetDescription,
					}

					_ = s.Create(target)
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			wantRes: &orchestrator.ListCertificationTargetsResponse{
				Targets: []*orchestrator.CertificationTarget{
					{
						Id:          DefaultCertificationTargetId,
						Name:        config.DefaultCertificationTargetName,
						Description: config.DefaultCertificationTargetDescription,
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
				Targets: []*orchestrator.CertificationTarget{
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
				Targets: []*orchestrator.CertificationTarget{
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
			name: "Storage error: getting resources",
			fields: fields{
				storage: &testutil.StorageWithError{CountErr: gorm.ErrInvalidDB},
				authz:   servicetest.NewAuthorizationStrategy(false, testdata.MockCertificationTargetID1),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.GetCertificationTargetStatisticsRequest{
					CertificationTargetId: testdata.MockCertificationTargetID1,
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.Internal, status.Code(err))
				return assert.ErrorContains(t, err, "database error counting resources:")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					// Store one certification target
					assert.NoError(t, s.Create(&orchestrator.CertificationTarget{
						Id:          testdata.MockCertificationTargetID1,
						Name:        testdata.MockCertificationTargetName1,
						Description: testdata.MockCertificationTargetDescription1,
					}))

					// Store evidences for certification target
					assert.NoError(t, s.Create(&evidence.Evidence{
						Id:                    uuid.NewString(),
						CertificationTargetId: testdata.MockCertificationTargetID1,
					}))
					assert.NoError(t, s.Create(&evidence.Evidence{
						Id:                    uuid.NewString(),
						CertificationTargetId: testdata.MockCertificationTargetID2,
					}))

					// Store assessment results for certification target
					assert.NoError(t, s.Create(&assessment.AssessmentResult{
						Id:                    uuid.NewString(),
						CertificationTargetId: testdata.MockCertificationTargetID1,
					}))
					assert.NoError(t, s.Create(&assessment.AssessmentResult{
						Id:                    uuid.NewString(),
						CertificationTargetId: testdata.MockCertificationTargetID1,
					}))
					assert.NoError(t, s.Create(&assessment.AssessmentResult{
						Id:                    uuid.NewString(),
						CertificationTargetId: testdata.MockCertificationTargetID2,
					}))

					// Store audit scopes
					assert.NoError(t, s.Create(orchestratortest.NewAuditScope("", "", testdata.MockCertificationTargetID1)))
					assert.NoError(t, s.Create(orchestratortest.NewAuditScope("", "", testdata.MockCertificationTargetID1)))

					// Store resources
					assert.NoError(t, s.Create(&discovery.Resource{
						Id:                    uuid.NewString(),
						CertificationTargetId: testdata.MockCertificationTargetID1,
					}))
					assert.NoError(t, s.Create(&discovery.Resource{
						Id:                    uuid.NewString(),
						CertificationTargetId: testdata.MockCertificationTargetID2,
					}))
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
				NumberOfSelectedCatalogs:    2,
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

func TestService_RegisterCertificationTarget(t *testing.T) {
	type fields struct {
		UnimplementedOrchestratorServer orchestrator.UnimplementedOrchestratorServer
		CertificationTargetHooks        []orchestrator.CertificationTargetHookFunc
		auditScopeHooks                 []orchestrator.AuditScopeHookFunc
		AssessmentResultHooks           []assessment.ResultHookFunc
		storage                         persistence.Storage
		metricsFile                     string
		loadMetricsFunc                 func() ([]*assessment.Metric, error)
		catalogsFolder                  string
		loadCatalogsFunc                func() ([]*orchestrator.Catalog, error)
		events                          chan *orchestrator.MetricChangeEvent
		authz                           service.AuthorizationStrategy
	}
	type args struct {
		ctx context.Context
		req *orchestrator.RegisterCertificationTargetRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes assert.Want[*orchestrator.CertificationTarget]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Request validation error",
			args: args{
				req: &orchestrator.RegisterCertificationTargetRequest{
					CertificationTarget: &orchestrator.CertificationTarget{},
				},
			},
			wantRes: assert.Nil[*orchestrator.CertificationTarget],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, " validation error:\n - certification_target.name: value length must be at least 1 characters [string.min_len]")
			},
		},
		{
			name: "Database error",
			fields: fields{
				storage: &testutil.StorageWithError{CreateErr: gorm.ErrInvalidDB},
			},
			args: args{
				req: &orchestrator.RegisterCertificationTargetRequest{
					CertificationTarget: orchestratortest.NewCertificationTarget(),
				},
			},
			wantRes: assert.Nil[*orchestrator.CertificationTarget],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not add certification target to the database:")
			},
		},
		{
			name: "Happy path: with metadata as input",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
			},
			args: args{
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
			},
			wantRes: func(t *testing.T, got *orchestrator.CertificationTarget) bool {
				want := &orchestrator.CertificationTarget{
					Name:        "test",
					Description: "some",
					CreatedAt:   &timestamppb.Timestamp{},
					UpdatedAt:   &timestamppb.Timestamp{},
					Metadata: &orchestrator.CertificationTarget_Metadata{
						Labels: map[string]string{
							"owner": "testOwner",
							"env":   "prod",
						},
					},
				}

				// Check if ID is set and delete it for the comparison
				assert.NotEmpty(t, got.GetId())
				got.Id = ""

				// Check if timestamp is set and delete it for the comparison
				assert.NotEmpty(t, got.CreatedAt)
				got.CreatedAt = &timestamppb.Timestamp{}

				// Check if updated_at is set and delete it for the comparison
				assert.NotEmpty(t, got.UpdatedAt)
				got.UpdatedAt = &timestamppb.Timestamp{}

				return assert.Equal(t, want, got)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: without metadata as input",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
			},
			args: args{
				req: &orchestrator.RegisterCertificationTargetRequest{
					CertificationTarget: &orchestrator.CertificationTarget{
						Name:        "test",
						Description: "some",
					},
				},
			},
			wantRes: func(t *testing.T, got *orchestrator.CertificationTarget) bool {
				want := &orchestrator.CertificationTarget{
					Name:        "test",
					Description: "some",
					CreatedAt:   &timestamppb.Timestamp{},
					UpdatedAt:   &timestamppb.Timestamp{},
					Metadata:    &orchestrator.CertificationTarget_Metadata{},
				}

				// Check if ID is set and delete it for the comparison
				assert.NotEmpty(t, got.GetId())
				got.Id = ""

				// Check if timestamp is set and delete it for the comparison
				assert.NotEmpty(t, got.CreatedAt)
				got.CreatedAt = &timestamppb.Timestamp{}

				// Check if updated_at is set and delete it for the comparison
				assert.NotEmpty(t, got.UpdatedAt)
				got.UpdatedAt = &timestamppb.Timestamp{}

				return assert.Equal(t, want, got)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				UnimplementedOrchestratorServer: tt.fields.UnimplementedOrchestratorServer,
				CertificationTargetHooks:        tt.fields.CertificationTargetHooks,
				auditScopeHooks:                 tt.fields.auditScopeHooks,
				AssessmentResultHooks:           tt.fields.AssessmentResultHooks,
				storage:                         tt.fields.storage,
				metricsFile:                     tt.fields.metricsFile,
				loadMetricsFunc:                 tt.fields.loadMetricsFunc,
				catalogsFolder:                  tt.fields.catalogsFolder,
				loadCatalogsFunc:                tt.fields.loadCatalogsFunc,
				events:                          tt.fields.events,
				authz:                           tt.fields.authz,
			}
			gotRes, err := s.RegisterCertificationTarget(tt.args.ctx, tt.args.req)

			tt.wantErr(t, err)
			tt.wantRes(t, gotRes)
		})
	}
}
