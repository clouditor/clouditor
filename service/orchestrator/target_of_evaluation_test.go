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

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func init() {
	viper.Set(config.DefaultTargetOfEvaluationNameFlag, config.DefaultTargetOfEvaluationName)
	viper.Set(config.DefaultTargetOfEvaluationDescriptionFlag, config.DefaultTargetOfEvaluationDescription)
	viper.Set(config.DefaultTargetOfEvaluationTypeFlag, int32(config.DefaultTargetOfEvaluationType))
}

func TestService_GetTargetOfEvaluation(t *testing.T) {
	tests := []struct {
		name    string
		svc     *Service
		ctx     context.Context
		req     *orchestrator.GetTargetOfEvaluationRequest
		res     *orchestrator.TargetOfEvaluation
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
			name: "target of evaluation not found",
			svc:  NewService(),
			ctx:  context.Background(),
			req:  &orchestrator.GetTargetOfEvaluationRequest{TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1},
			res:  nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "target of evaluation not found") &&
					assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "valid",
			svc:  NewService(),
			ctx:  context.Background(),
			req:  &orchestrator.GetTargetOfEvaluationRequest{TargetOfEvaluationId: DefaultTargetOfEvaluationId},
			res: &orchestrator.TargetOfEvaluation{
				Id:          DefaultTargetOfEvaluationId,
				Name:        config.DefaultTargetOfEvaluationName,
				Description: config.DefaultTargetOfEvaluationDescription,
				TargetType:  config.DefaultTargetOfEvaluationType,
			},
			wantErr: assert.NoError,
		},
		{
			name: "permission denied",
			svc:  NewService(WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID1))),
			ctx:  context.TODO(),
			req:  &orchestrator.GetTargetOfEvaluationRequest{TargetOfEvaluationId: DefaultTargetOfEvaluationId},
			res:  nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, service.ErrPermissionDenied.Error()) &&
					assert.Equal(t, codes.PermissionDenied, status.Code(err))
			},
		},
		{
			name: "permission granted",
			svc: NewService(WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID1)), WithStorage(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
				_ = s.Create(&orchestrator.TargetOfEvaluation{
					Id:        testdata.MockTargetOfEvaluationID1,
					Name:      "target1",
					CreatedAt: timestamppb.Now(),
					UpdatedAt: timestamppb.Now(),
				})
			}))),
			ctx: context.TODO(),
			req: &orchestrator.GetTargetOfEvaluationRequest{TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1},
			res: &orchestrator.TargetOfEvaluation{
				Id:   testdata.MockTargetOfEvaluationID1,
				Name: "target1",
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.svc.CreateDefaultTargetOfEvaluation()
			assert.NoError(t, err)

			res, err := tt.svc.GetTargetOfEvaluation(tt.ctx, tt.req)
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

func TestService_UpdateTargetOfEvaluation(t *testing.T) {
	var (
		target *orchestrator.TargetOfEvaluation
		err    error
	)
	orchestratorService := NewService(WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID1)))

	// 1st case: Target of Evaluation is nil
	_, err = orchestratorService.UpdateTargetOfEvaluation(context.TODO(), &orchestrator.UpdateTargetOfEvaluationRequest{})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// 2nd case: Target of Evaluation ID is nil
	_, err = orchestratorService.UpdateTargetOfEvaluation(context.TODO(), &orchestrator.UpdateTargetOfEvaluationRequest{
		TargetOfEvaluation: &orchestrator.TargetOfEvaluation{},
	})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// 3rd case: Target of Evaluation not found since there are no target of evaluations yet
	_, err = orchestratorService.UpdateTargetOfEvaluation(context.TODO(), &orchestrator.UpdateTargetOfEvaluationRequest{
		TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
			Id:          testdata.MockTargetOfEvaluationID1,
			Name:        config.DefaultTargetOfEvaluationName,
			Description: config.DefaultTargetOfEvaluationDescription,
		},
	})
	assert.Equal(t, codes.NotFound, status.Code(err))

	// 4th case: Target of Evaluation updated successfully
	err = orchestratorService.storage.Create(&orchestrator.TargetOfEvaluation{
		Id:          testdata.MockTargetOfEvaluationID1,
		Name:        config.DefaultTargetOfEvaluationName,
		Description: config.DefaultTargetOfEvaluationDescription,
	})
	assert.NoError(t, err)
	if err != nil {
		return
	}
	target, err = orchestratorService.UpdateTargetOfEvaluation(context.TODO(), &orchestrator.UpdateTargetOfEvaluationRequest{
		TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
			Id:          testdata.MockTargetOfEvaluationID1,
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

func TestService_RemoveTargetOfEvaluation(t *testing.T) {
	var (
		TargetOfEvaluationResponse      *orchestrator.TargetOfEvaluation
		err                             error
		listTargetsOfEvaluationResponse *orchestrator.ListTargetsOfEvaluationResponse
	)
	orchestratorService := NewService()

	// 1st case: Empty target of evaluation ID error
	_, err = orchestratorService.RemoveTargetOfEvaluation(context.Background(), &orchestrator.RemoveTargetOfEvaluationRequest{TargetOfEvaluationId: ""})
	assert.Error(t, err)
	assert.Equal(t, status.Code(err), codes.InvalidArgument)

	// 2nd case: ErrRecordNotFound
	_, err = orchestratorService.RemoveTargetOfEvaluation(context.Background(), &orchestrator.RemoveTargetOfEvaluationRequest{TargetOfEvaluationId: DefaultTargetOfEvaluationId})
	assert.Error(t, err)
	assert.Equal(t, status.Code(err), codes.NotFound)

	// 3rd case: Record removed successfully
	TargetOfEvaluationResponse, err = orchestratorService.CreateDefaultTargetOfEvaluation()
	assert.NoError(t, err)
	assert.NotNil(t, TargetOfEvaluationResponse)

	// There is a record for target of evaluations in the DB (default one)
	listTargetsOfEvaluationResponse, err = orchestratorService.ListTargetsOfEvaluation(context.Background(), &orchestrator.ListTargetsOfEvaluationRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listTargetsOfEvaluationResponse.Targets)
	assert.NotEmpty(t, listTargetsOfEvaluationResponse.Targets)

	// Remove record
	_, err = orchestratorService.RemoveTargetOfEvaluation(context.Background(), &orchestrator.RemoveTargetOfEvaluationRequest{TargetOfEvaluationId: DefaultTargetOfEvaluationId})
	assert.NoError(t, err)

	// There is a record for target of evaluations in the DB (default one)
	listTargetsOfEvaluationResponse, err = orchestratorService.ListTargetsOfEvaluation(context.Background(), &orchestrator.ListTargetsOfEvaluationRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listTargetsOfEvaluationResponse.Targets)
	assert.Empty(t, listTargetsOfEvaluationResponse.Targets)
}

func TestService_CreateDefaultTargetOfEvaluation(t *testing.T) {
	var (
		TargetOfEvaluationResponse *orchestrator.TargetOfEvaluation
		err                        error
	)
	orchestratorService := NewService()

	// 1st case: No records for target of evaluations -> Default target of evaluation is created
	TargetOfEvaluationResponse, err = orchestratorService.CreateDefaultTargetOfEvaluation()
	assert.NoError(t, err)
	// Check timestamps and delete it for further tests
	assert.NotEmpty(t, TargetOfEvaluationResponse.CreatedAt)
	assert.NotEmpty(t, TargetOfEvaluationResponse.UpdatedAt)
	TargetOfEvaluationResponse.CreatedAt = nil
	TargetOfEvaluationResponse.UpdatedAt = nil

	assert.Equal(t, &orchestrator.TargetOfEvaluation{
		Id:          DefaultTargetOfEvaluationId,
		Name:        config.DefaultTargetOfEvaluationName,
		Description: config.DefaultTargetOfEvaluationDescription,
		TargetType:  config.DefaultTargetOfEvaluationType,
	}, TargetOfEvaluationResponse)

	// Check if TargetOfEvaluation is valid
	assert.NoError(t, api.Validate(TargetOfEvaluationResponse))

	// 2nd case: There is already a record for the target of evaluation (the default target of evaluation) -> Nothing added and no error
	TargetOfEvaluationResponse, err = orchestratorService.CreateDefaultTargetOfEvaluation()
	assert.NoError(t, err)
	assert.Nil(t, TargetOfEvaluationResponse)
}

func TestService_ListTargetsOfEvaluation(t *testing.T) {
	type fields struct {
		AssessmentResultHooks []assessment.ResultHookFunc
		storage               persistence.Storage
		catalogsFolder        string
		loadCatalogsFunc      func() ([]*orchestrator.Catalog, error)
		events                chan *orchestrator.MetricChangeEvent
		authz                 service.AuthorizationStrategy
	}
	type args struct {
		ctx context.Context
		req *orchestrator.ListTargetsOfEvaluationRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *orchestrator.ListTargetsOfEvaluationResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "retrieve empty list",
			args: args{req: &orchestrator.ListTargetsOfEvaluationRequest{}},
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			wantRes: &orchestrator.ListTargetsOfEvaluationResponse{},
			wantErr: assert.NoError,
		},
		{
			name: "list with one item",
			args: args{req: &orchestrator.ListTargetsOfEvaluationRequest{}},
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					target := &orchestrator.TargetOfEvaluation{
						Id:          DefaultTargetOfEvaluationId,
						Name:        config.DefaultTargetOfEvaluationName,
						Description: config.DefaultTargetOfEvaluationDescription,
					}

					_ = s.Create(target)
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			wantRes: &orchestrator.ListTargetsOfEvaluationResponse{
				Targets: []*orchestrator.TargetOfEvaluation{
					{
						Id:          DefaultTargetOfEvaluationId,
						Name:        config.DefaultTargetOfEvaluationName,
						Description: config.DefaultTargetOfEvaluationDescription,
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "retrieve only allowed target of evaluations: no target of evaluation is allowed",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					// Store two target of evaluations, of which none we are allowed to retrieve in the test
					_ = s.Create(&orchestrator.TargetOfEvaluation{
						Id:   testdata.MockTargetOfEvaluationID1,
						Name: testdata.MockTargetOfEvaluationName1,
					})
					_ = s.Create(&orchestrator.TargetOfEvaluation{
						Id:   testdata.MockTargetOfEvaluationID2,
						Name: testdata.MockTargetOfEvaluationName2,
					})
				}),
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID1),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.ListTargetsOfEvaluationRequest{},
			},
			wantRes: &orchestrator.ListTargetsOfEvaluationResponse{
				Targets: []*orchestrator.TargetOfEvaluation{
					{
						Id:   testdata.MockTargetOfEvaluationID1,
						Name: testdata.MockTargetOfEvaluationName1,
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "retrieve only allowed target of evaluations",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					// Store two target of evaluations, of which only one we are allowed to retrieve in the test
					_ = s.Create(&orchestrator.TargetOfEvaluation{
						Id:   testdata.MockTargetOfEvaluationID1,
						Name: testdata.MockTargetOfEvaluationName1,
					})
					_ = s.Create(&orchestrator.TargetOfEvaluation{
						Id:   testdata.MockTargetOfEvaluationID2,
						Name: testdata.MockTargetOfEvaluationName1,
					})
				}),
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID1),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.ListTargetsOfEvaluationRequest{},
			},
			wantRes: &orchestrator.ListTargetsOfEvaluationResponse{
				Targets: []*orchestrator.TargetOfEvaluation{
					{
						Id:   testdata.MockTargetOfEvaluationID1,
						Name: testdata.MockTargetOfEvaluationName1,
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
				catalogsFolder:        tt.fields.catalogsFolder,
				loadCatalogsFunc:      tt.fields.loadCatalogsFunc,
				events:                tt.fields.events,
				authz:                 tt.fields.authz,
			}

			gotRes, err := svc.ListTargetsOfEvaluation(tt.args.ctx, tt.args.req)
			assert.NoError(t, api.Validate(gotRes))

			tt.wantErr(t, err, tt.args)
			assert.Equal(t, tt.wantRes, gotRes)
		})
	}
}

func TestService_GetTargetOfEvaluationStatistics(t *testing.T) {
	type fields struct {
		auditScopeHooks       []orchestrator.AuditScopeHookFunc
		AssessmentResultHooks []assessment.ResultHookFunc
		storage               persistence.Storage
		loadMetricsFunc       func() ([]*assessment.Metric, error)
		catalogsFolder        string
		loadCatalogsFunc      func() ([]*orchestrator.Catalog, error)
		events                chan *orchestrator.MetricChangeEvent
		authz                 service.AuthorizationStrategy
	}
	type args struct {
		ctx context.Context
		req *orchestrator.GetTargetOfEvaluationStatisticsRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *orchestrator.GetTargetOfEvaluationStatisticsResponse
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
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID2),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.GetTargetOfEvaluationStatisticsRequest{
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
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
				authz:   servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID1),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.GetTargetOfEvaluationStatisticsRequest{
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
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
					// Store one target of evaluation
					assert.NoError(t, s.Create(&orchestrator.TargetOfEvaluation{
						Id:          testdata.MockTargetOfEvaluationID1,
						Name:        testdata.MockTargetOfEvaluationName1,
						Description: testdata.MockTargetOfEvaluationDescription1,
					}))

					// Store evidences for target of evaluation
					assert.NoError(t, s.Create(&evidence.Evidence{
						Id:                   uuid.NewString(),
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					}))
					assert.NoError(t, s.Create(&evidence.Evidence{
						Id:                   uuid.NewString(),
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID2,
					}))

					// Store assessment results for target of evaluation
					assert.NoError(t, s.Create(&assessment.AssessmentResult{
						Id:                   uuid.NewString(),
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					}))
					assert.NoError(t, s.Create(&assessment.AssessmentResult{
						Id:                   uuid.NewString(),
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					}))
					assert.NoError(t, s.Create(&assessment.AssessmentResult{
						Id:                   uuid.NewString(),
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID2,
					}))

					// Store audit scopes
					assert.NoError(t, s.Create(orchestratortest.NewAuditScope("", "", testdata.MockTargetOfEvaluationID1)))
					assert.NoError(t, s.Create(orchestratortest.NewAuditScope("", "", testdata.MockTargetOfEvaluationID1)))

					// Store resources
					assert.NoError(t, s.Create(&discovery.Resource{
						Id:                   uuid.NewString(),
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					}))
					assert.NoError(t, s.Create(&discovery.Resource{
						Id:                   uuid.NewString(),
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID2,
					}))
				}),
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID1),
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.GetTargetOfEvaluationStatisticsRequest{
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
				},
			},
			wantRes: &orchestrator.GetTargetOfEvaluationStatisticsResponse{
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
				auditScopeHooks:       tt.fields.auditScopeHooks,
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				catalogsFolder:        tt.fields.catalogsFolder,
				loadCatalogsFunc:      tt.fields.loadCatalogsFunc,
				events:                tt.fields.events,
				authz:                 tt.fields.authz,
			}
			gotRes, err := s.GetTargetOfEvaluationStatistics(tt.args.ctx, tt.args.req)
			tt.wantErr(t, err)
			assert.Equal(t, tt.wantRes, gotRes)
		})
	}
}

func TestService_CreateTargetOfEvaluation(t *testing.T) {
	type fields struct {
		UnimplementedOrchestratorServer orchestrator.UnimplementedOrchestratorServer
		TargetOfEvaluationHooks         []orchestrator.TargetOfEvaluationHookFunc
		auditScopeHooks                 []orchestrator.AuditScopeHookFunc
		AssessmentResultHooks           []assessment.ResultHookFunc
		storage                         persistence.Storage
		loadMetricsFunc                 func() ([]*assessment.Metric, error)
		catalogsFolder                  string
		loadCatalogsFunc                func() ([]*orchestrator.Catalog, error)
		events                          chan *orchestrator.MetricChangeEvent
		authz                           service.AuthorizationStrategy
	}
	type args struct {
		ctx context.Context
		req *orchestrator.CreateTargetOfEvaluationRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes assert.Want[*orchestrator.TargetOfEvaluation]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Request validation error",
			args: args{
				req: &orchestrator.CreateTargetOfEvaluationRequest{
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{},
				},
			},
			wantRes: assert.Nil[*orchestrator.TargetOfEvaluation],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, " validation error:\n - target_of_evaluation.name: value length must be at least 1 characters [string.min_len]")
			},
		},
		{
			name: "Database error",
			fields: fields{
				storage: &testutil.StorageWithError{CreateErr: gorm.ErrInvalidDB},
			},
			args: args{
				req: &orchestrator.CreateTargetOfEvaluationRequest{
					TargetOfEvaluation: orchestratortest.NewTargetOfEvaluation(),
				},
			},
			wantRes: assert.Nil[*orchestrator.TargetOfEvaluation],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not add target of evaluation to the database:")
			},
		},
		{
			name: "Happy path: with metadata as input",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
			},
			args: args{
				req: &orchestrator.CreateTargetOfEvaluationRequest{
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						Name:        "test",
						Description: "some",
						Metadata: &orchestrator.TargetOfEvaluation_Metadata{
							Labels: map[string]string{
								"owner": "testOwner",
								"env":   "prod",
							},
						},
					},
				},
			},
			wantRes: func(t *testing.T, got *orchestrator.TargetOfEvaluation) bool {
				want := &orchestrator.TargetOfEvaluation{
					Name:        "test",
					Description: "some",
					CreatedAt:   &timestamppb.Timestamp{},
					UpdatedAt:   &timestamppb.Timestamp{},
					Metadata: &orchestrator.TargetOfEvaluation_Metadata{
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
				req: &orchestrator.CreateTargetOfEvaluationRequest{
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						Name:        "test",
						Description: "some",
					},
				},
			},
			wantRes: func(t *testing.T, got *orchestrator.TargetOfEvaluation) bool {
				want := &orchestrator.TargetOfEvaluation{
					Name:        "test",
					Description: "some",
					CreatedAt:   &timestamppb.Timestamp{},
					UpdatedAt:   &timestamppb.Timestamp{},
					Metadata:    &orchestrator.TargetOfEvaluation_Metadata{},
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
				TargetOfEvaluationHooks:         tt.fields.TargetOfEvaluationHooks,
				auditScopeHooks:                 tt.fields.auditScopeHooks,
				AssessmentResultHooks:           tt.fields.AssessmentResultHooks,
				storage:                         tt.fields.storage,
				catalogsFolder:                  tt.fields.catalogsFolder,
				loadCatalogsFunc:                tt.fields.loadCatalogsFunc,
				events:                          tt.fields.events,
				authz:                           tt.fields.authz,
			}
			gotRes, err := s.CreateTargetOfEvaluation(tt.args.ctx, tt.args.req)

			tt.wantErr(t, err)
			tt.wantRes(t, gotRes)
		})
	}
}
