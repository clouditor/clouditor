// Copyright 2016-2022 Fraunhofer AISEC
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
	"reflect"
	"runtime"
	"sync"
	"testing"
	"time"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/evaluation"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/evaluationtest"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/orchestratortest"
	"clouditor.io/clouditor/v2/internal/util"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func TestService_CreateAuditScope(t *testing.T) {

	catalogWithoutAssuranceLevelList := orchestratortest.NewCatalog()
	catalogWithoutAssuranceLevelList.AssuranceLevels = []string{}

	type fields struct {
		TargetOfEvaluationHooks []orchestrator.TargetOfEvaluationHookFunc
		AssessmentResultHooks   []assessment.ResultHookFunc
		storage                 persistence.Storage
		catalogsFolder          string
		loadCatalogsFunc        func() ([]*orchestrator.Catalog, error)
		events                  chan *orchestrator.MetricChangeEvent
		authz                   service.AuthorizationStrategy
	}
	type args struct {
		ctx context.Context
		req *orchestrator.CreateAuditScopeRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.Want[*orchestrator.AuditScope]
		wantSvc assert.Want[*Service]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "error: empty request",
			args: args{
				ctx: context.Background(),
				req: &orchestrator.CreateAuditScopeRequest{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, api.ErrEmptyRequest.Error())
			},
			wantSvc: assert.NotNil[*Service],
			want:    assert.Nil[*orchestrator.AuditScope],
		},
		{
			name: "error: invalid request",
			args: args{
				ctx: context.Background(),
				req: &orchestrator.CreateAuditScopeRequest{
					AuditScope: &orchestrator.AuditScope{
						CatalogId: testdata.MockCatalogID1,
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "audit_scope.target_of_evaluation_id: value is empty")
			},
			wantSvc: assert.NotNil[*Service],
			want:    assert.Nil[*orchestrator.AuditScope],
		},
		{
			name: "error: permission denied",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(&orchestrator.TargetOfEvaluation{Id: testdata.MockTargetOfEvaluationID1}))
				}),
				authz: servicetest.NewAuthorizationStrategy(false),
			},
			args: args{
				ctx: context.Background(),
				req: &orchestrator.CreateAuditScopeRequest{
					AuditScope: &orchestrator.AuditScope{
						CatalogId:            testdata.MockCatalogID1,
						Name:                 testdata.MockAuditScopeName1,
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.PermissionDenied, status.Code(err))
				return assert.ErrorContains(t, err, service.ErrPermissionDenied.Error())
			},
			wantSvc: assert.NotNil[*Service],
			want:    assert.Nil[*orchestrator.AuditScope],
		},
		{
			name: "error: database error",
			fields: fields{
				storage: &testutil.StorageWithError{CreateErr: gorm.ErrInvalidDB},
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				ctx: context.Background(),
				req: &orchestrator.CreateAuditScopeRequest{
					AuditScope: &orchestrator.AuditScope{
						CatalogId:            testdata.MockCatalogID1,
						Name:                 testdata.MockAuditScopeName1,
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.Internal, status.Code(err))
				return assert.ErrorContains(t, err, persistence.ErrDatabase.Error())
			},
			wantSvc: assert.NotNil[*Service],
			want:    assert.Nil[*orchestrator.AuditScope],
		},
		{
			name: "valid and assurance level not set",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{req: &orchestrator.CreateAuditScopeRequest{
				AuditScope: orchestratortest.NewAuditScope("", testdata.MockAuditScopeID1, "", testdata.MockAuditScopeName1),
			}},
			wantSvc: func(t *testing.T, svc *Service) bool {
				auditScope := []*orchestrator.AuditScope{}
				// Check if audit scope is created in the DB
				err := svc.storage.List(&auditScope, "", true, 0, -1)
				assert.NoError(t, err)
				assert.Equal(t, 1, len(auditScope))

				want := &orchestrator.AuditScope{
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					Name:                 testdata.MockAuditScopeName1,
					CatalogId:            testdata.MockCatalogID1,
				}

				// We check if the ID is set and remove it, as we are unable to verify it.
				assert.NotEmpty(t, auditScope[0].GetId())
				auditScope[0].Id = ""
				return assert.Equal(t, want, auditScope[0])
			},
			want:    assert.AnyValue[*orchestrator.AuditScope],
			wantErr: assert.NoError,
		},
		{
			name: "valid and assurance level set",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Create(&orchestrator.TargetOfEvaluation{Id: testdata.MockTargetOfEvaluationID1})
					assert.NoError(t, err)
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{req: &orchestrator.CreateAuditScopeRequest{
				AuditScope: orchestratortest.NewAuditScope(testdata.AssuranceLevelHigh, testdata.MockAuditScopeID1, "", testdata.MockAuditScopeName1),
			}},
			wantSvc: func(t *testing.T, svc *Service) bool {
				auditScope := []*orchestrator.AuditScope{}
				// Check if audit scope is created in the DB
				err := svc.storage.List(&auditScope, "", true, 0, -1)
				assert.NoError(t, err)
				assert.Equal(t, 1, len(auditScope))

				want := &orchestrator.AuditScope{
					TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					CatalogId:            testdata.MockCatalogID1,
					Name:                 testdata.MockAuditScopeName1,
					AssuranceLevel:       &testdata.AssuranceLevelHigh,
				}

				// We check if the ID is set and remove it, as we are unable to verify it.
				assert.NotEmpty(t, auditScope[0].GetId())
				auditScope[0].Id = ""
				return assert.Equal(t, want, auditScope[0])
			},
			want:    assert.AnyValue[*orchestrator.AuditScope],
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				TargetOfEvaluationHooks: tt.fields.TargetOfEvaluationHooks,
				AssessmentResultHooks:   tt.fields.AssessmentResultHooks,
				storage:                 tt.fields.storage,
				catalogsFolder:          tt.fields.catalogsFolder,
				loadCatalogsFunc:        tt.fields.loadCatalogsFunc,
				events:                  tt.fields.events,
				authz:                   tt.fields.authz,
			}

			gotRes, err := svc.CreateAuditScope(tt.args.ctx, tt.args.req)
			tt.wantErr(t, err)
			tt.wantSvc(t, svc)
			tt.want(t, gotRes)
		})
	}
}

func TestService_GetAuditScope(t *testing.T) {
	type fields struct {
		storage persistence.Storage
		authz   service.AuthorizationStrategy
	}
	type args struct {
		req *orchestrator.GetAuditScopeRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantResponse assert.Want[*orchestrator.AuditScope]
		wantErr      assert.ErrorAssertionFunc
	}{

		{
			name: "Error: invalid request",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{req: &orchestrator.GetAuditScopeRequest{
				AuditScopeId: "",
			}},
			wantResponse: assert.Nil[*orchestrator.AuditScope],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "audit_scope_id: value is empty, which is not a valid UUID ")
			},
		},
		{
			name: "Error: auditScope not found",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Create(orchestratortest.NewAuditScope(testdata.AssuranceLevelBasic, "", "", testdata.MockAuditScopeName1))
					assert.NoError(t, err)
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{req: &orchestrator.GetAuditScopeRequest{
				AuditScopeId: testdata.MockAuditScopeID2,
			}},
			wantResponse: assert.Nil[*orchestrator.AuditScope],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.NotFound, status.Code(err))
				return assert.ErrorContains(t, err, api.ErrAuditScopeNotFound.Error())
			},
		},
		{
			name: "Error: permission denied",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Create(orchestratortest.NewAuditScope(testdata.AssuranceLevelBasic, testdata.MockAuditScopeID1, testdata.MockTargetOfEvaluationID1, testdata.MockAuditScopeName1))
					assert.NoError(t, err)
				}),
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID2),
			},
			args: args{req: &orchestrator.GetAuditScopeRequest{
				AuditScopeId: testdata.MockAuditScopeID1,
			}},
			wantResponse: assert.Nil[*orchestrator.AuditScope],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.PermissionDenied, status.Code(err))
				return assert.ErrorContains(t, err, "access denied")
			},
		},
		{
			name: "Error: database error",
			fields: fields{
				storage: &testutil.StorageWithError{GetErr: gorm.ErrInvalidDB},
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{req: &orchestrator.GetAuditScopeRequest{
				AuditScopeId: testdata.MockAuditScopeID1,
			}},
			wantResponse: assert.Nil[*orchestrator.AuditScope],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.Internal, status.Code(err))
				return assert.ErrorContains(t, err, persistence.ErrDatabase.Error())
			},
		},
		{
			name: "Happy path",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {

					err := s.Create(orchestratortest.NewAuditScope(testdata.AssuranceLevelBasic, testdata.MockAuditScopeID1, "", testdata.MockAuditScopeName1))
					assert.NoError(t, err)
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{req: &orchestrator.GetAuditScopeRequest{
				AuditScopeId: testdata.MockAuditScopeID1,
			}},
			wantResponse: func(t *testing.T, got *orchestrator.AuditScope) bool {
				want := orchestratortest.NewAuditScope(testdata.AssuranceLevelBasic, testdata.MockAuditScopeID1, "", testdata.MockAuditScopeName1)

				return assert.NoError(t, api.Validate(got)) &&
					assert.Equal(t, want.TargetOfEvaluationId, got.TargetOfEvaluationId) &&
					assert.Equal(t, want.CatalogId, got.CatalogId)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orchestratorService := Service{
				storage: tt.fields.storage,
				authz:   tt.fields.authz,
			}
			res, err := orchestratorService.GetAuditScope(context.Background(), tt.args.req)
			tt.wantErr(t, err)
			tt.wantResponse(t, res)
		})
	}
}

func TestToeHook(t *testing.T) {
	var (
		hookCallCounter = 0
		wg              sync.WaitGroup
	)
	wg.Add(2)

	firstHookFunction := func(ctx context.Context, event *orchestrator.AuditScopeChangeEvent, err error) {
		hookCallCounter++
		log.Println("Hello from inside the first auditScope hook function")

		// Checking the status
		// UpdateAuditScope is called, so the status must be AUDIT_SCOPE_UPDATE
		if *event.GetType().Enum() != orchestrator.AuditScopeChangeEvent_TYPE_AUDIT_SCOPE_UPDATED {
			return
		}

		wg.Done()
	}

	secondHookFunction := func(ctx context.Context, event *orchestrator.AuditScopeChangeEvent, err error) {
		hookCallCounter++
		log.Println("Hello from inside the second auditScope hook function")

		wg.Done()
	}

	svc := NewService()
	svc.RegisterToeHook(firstHookFunction)
	svc.RegisterToeHook(secondHookFunction)

	// Check if first hook is registered
	funcName1 := runtime.FuncForPC(reflect.ValueOf(svc.auditScopeHooks[0]).Pointer()).Name()
	funcName2 := runtime.FuncForPC(reflect.ValueOf(firstHookFunction).Pointer()).Name()
	assert.Equal(t, funcName1, funcName2)

	// Check if second hook is registered
	funcName1 = runtime.FuncForPC(reflect.ValueOf(svc.auditScopeHooks[1]).Pointer()).Name()
	funcName2 = runtime.FuncForPC(reflect.ValueOf(secondHookFunction).Pointer()).Name()
	assert.Equal(t, funcName1, funcName2)

	type args struct {
		ctx context.Context
		req *orchestrator.UpdateAuditScopeRequest
	}
	tests := []struct {
		name    string
		args    args
		wantRes *orchestrator.AuditScope
		wantErr assert.WantErr
	}{
		{
			name: "Store first audit scope to the map",
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.UpdateAuditScopeRequest{
					AuditScope: &orchestrator.AuditScope{
						Id:                   testdata.MockAuditScopeID1,
						Name:                 testdata.MockAuditScopeName1,
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
						CatalogId:            testdata.MockCatalogID1,
						AssuranceLevel:       &testdata.AssuranceLevelSubstantial,
					},
				},
			},
			wantErr: assert.Nil[error],
			wantRes: &orchestrator.AuditScope{
				Id:                   testdata.MockAuditScopeID1,
				Name:                 testdata.MockAuditScopeName1,
				TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
				CatalogId:            testdata.MockCatalogID1,
				AssuranceLevel:       &testdata.AssuranceLevelSubstantial,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hookCallCounter = 0

			// Create service
			s := svc

			// Create Audit Scope in DB
			err := s.storage.Create(orchestratortest.NewAuditScope(testdata.AssuranceLevelBasic, testdata.MockAuditScopeID1, "", testdata.MockAuditScopeName1))
			assert.NoError(t, err)

			gotRes, err := s.UpdateAuditScope(tt.args.ctx, tt.args.req)

			// wait for all hooks (2 hooks)
			wg.Wait()

			assert.NoError(t, api.Validate(gotRes))

			tt.wantErr(t, err)
			assert.Equal(t, tt.wantRes, gotRes)
			assert.Equal(t, 2, hookCallCounter)
		})
	}
}

func TestService_ListAuditScopes(t *testing.T) {
	type fields struct {
		storage persistence.Storage
		authz   service.AuthorizationStrategy
	}
	type args struct {
		ctx context.Context
		req *orchestrator.ListAuditScopesRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes assert.Want[*orchestrator.ListAuditScopesResponse]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Error: validation error",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args:    args{},
			wantRes: assert.Nil[*orchestrator.ListAuditScopesResponse],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "empty request")
			},
		},
		{
			name: "Permission denied",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Create(orchestratortest.MockAuditScopeCertTargetID1)
					assert.NoError(t, err)
				}),
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID1),
			},
			args: args{
				ctx: context.Background(),
				req: &orchestrator.ListAuditScopesRequest{
					Filter: &orchestrator.ListAuditScopesRequest_Filter{
						TargetOfEvaluationId: util.Ref(testdata.MockTargetOfEvaluationID2),
					},
				},
			},
			wantRes: assert.Nil[*orchestrator.ListAuditScopesResponse],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.PermissionDenied, status.Code(err))
				return assert.ErrorContains(t, err, "access denied")
			},
		},
		{
			name: "TargetOfEvaluationId filter and no access rights",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Create(orchestratortest.MockAuditScopeCertTargetID1)
					assert.NoError(t, err)
					err = s.Create(orchestratortest.MockAuditScopeCertTargetID2)
					assert.NoError(t, err)

				}),
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID2),
			},
			args: args{
				ctx: context.Background(),
				req: &orchestrator.ListAuditScopesRequest{
					Filter: &orchestrator.ListAuditScopesRequest_Filter{
						TargetOfEvaluationId: util.Ref(testdata.MockTargetOfEvaluationID2),
					},
				},
			},
			wantRes: func(t *testing.T, got *orchestrator.ListAuditScopesResponse) bool {
				assert.Equal(t, 0, len(got.AuditScopes))
				return assert.Empty(t, got.AuditScopes)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: with targetOfEvaluationId filter and access rights",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Create(orchestratortest.MockAuditScopeCertTargetID1)
					assert.NoError(t, err)
					err = s.Create(orchestratortest.MockAuditScopeCertTargetID2)
					assert.NoError(t, err)
					err = s.Create(orchestratortest.MockAuditScopeCertTargetID3)
					assert.NoError(t, err)
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				ctx: context.Background(),
				req: &orchestrator.ListAuditScopesRequest{
					Filter: &orchestrator.ListAuditScopesRequest_Filter{
						TargetOfEvaluationId: util.Ref(testdata.MockTargetOfEvaluationID1),
					},
				},
			},
			wantRes: func(t *testing.T, got *orchestrator.ListAuditScopesResponse) bool {
				assert.Equal(t, 2, len(got.AuditScopes))
				want := []*orchestrator.AuditScope{
					orchestratortest.MockAuditScopeCertTargetID1,
					orchestratortest.MockAuditScopeCertTargetID2,
				}
				return assert.Equal(t, want, got.AuditScopes)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: with catalogID filter and access rights",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Create(orchestratortest.MockAuditScopeCertTargetID1)
					assert.NoError(t, err)
					err = s.Create(orchestratortest.MockAuditScopeCertTargetID2)
					assert.NoError(t, err)
					err = s.Create(orchestratortest.MockAuditScopeCertTargetID3)
					assert.NoError(t, err)
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				ctx: context.Background(),
				req: &orchestrator.ListAuditScopesRequest{
					Filter: &orchestrator.ListAuditScopesRequest_Filter{
						CatalogId: util.Ref(testdata.MockCatalogID1),
					},
				},
			},
			wantRes: func(t *testing.T, got *orchestrator.ListAuditScopesResponse) bool {
				assert.Equal(t, 2, len(got.AuditScopes))
				want := []*orchestrator.AuditScope{
					orchestratortest.MockAuditScopeCertTargetID1,
					orchestratortest.MockAuditScopeCertTargetID3,
				}
				return assert.Equal(t, want, got.AuditScopes)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: without filter",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Create(orchestratortest.MockAuditScopeCertTargetID1)
					assert.NoError(t, err)
					err = s.Create(orchestratortest.MockAuditScopeCertTargetID2)
					assert.NoError(t, err)
					err = s.Create(orchestratortest.MockAuditScopeCertTargetID3)
					assert.NoError(t, err)
				}),
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID1),
			},
			args: args{
				ctx: context.Background(),
				req: &orchestrator.ListAuditScopesRequest{},
			},
			wantRes: func(t *testing.T, got *orchestrator.ListAuditScopesResponse) bool {
				assert.Equal(t, 2, len(got.AuditScopes))
				want := []*orchestrator.AuditScope{
					orchestratortest.MockAuditScopeCertTargetID1,
					orchestratortest.MockAuditScopeCertTargetID2,
				}
				return assert.Equal(t, want, got.AuditScopes)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				storage: tt.fields.storage,
				authz:   tt.fields.authz,
			}
			gotRes, err := svc.ListAuditScopes(tt.args.ctx, tt.args.req)
			tt.wantRes(t, gotRes)
			tt.wantErr(t, err)
		})
	}
}

func TestService_RemoveAuditScope(t *testing.T) {
	type fields struct {
		svc *Service
	}
	type args struct {
		ctx context.Context
		req *orchestrator.RemoveAuditScopeRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantResponse assert.Want[*emptypb.Empty]
		wantSvc      assert.Want[*Service]
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "Error: validation error",
			fields: fields{
				svc: NewService(),
			},
			args: args{
				ctx: nil,
				req: nil,
			},
			wantResponse: assert.Nil[*emptypb.Empty],
			wantSvc: func(t *testing.T, got *Service) bool {
				return true
			},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, api.ErrEmptyRequest.Error())
			},
		},
		{
			name: "Error: audit scope not found",
			fields: fields{
				svc: NewService(
					WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(
						true)),
					// Create empty storage => No audit scope can be found
					WithStorage(testutil.NewInMemoryStorage(t))),
			},
			args: args{
				ctx: nil,
				req: &orchestrator.RemoveAuditScopeRequest{
					AuditScopeId: testdata.MockAuditScopeID1,
				},
			},
			wantResponse: assert.Nil[*emptypb.Empty],
			wantSvc: func(t *testing.T, got *Service) bool {
				return true
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.NotFound, status.Code(err))
				return assert.ErrorContains(t, err, api.ErrAuditScopeNotFound.Error())
			},
		},
		// {
		// TODO(all): Currently we are not able to check that, as it is already checked by the GetAuditScope call in the method. As soon as we check the request type as well, we need this check.
		// 	name: "Error: permission denied",
		// 	fields: fields{
		// 		svc: NewService(
		// 			WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID1)),
		// 			WithStorage(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
		// 				assert.NoError(t, s.Create(orchestratortest.NewAuditScope("", testdata.MockAuditScopeID1, testdata.MockTargetOfEvaluationID2)))
		// 			})),
		// 		),
		// 	},
		// 	args: args{
		// 		ctx: nil,
		// 		req: &orchestrator.RemoveAuditScopeRequest{
		// 			AuditScopeId: testdata.MockAuditScopeID1,
		// 		},
		// 	},
		// 	wantResponse: assert.Nil[*emptypb.Empty],
		// 	wantSvc: func(t *testing.T, got *Service) bool {
		// 		n, err := got.storage.Count(&orchestrator.AuditScope{})
		// 		assert.NoError(t, err)
		// 		return assert.Equal(t, 1, n)
		// 	},
		// 	wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
		// 		assert.Equal(t, codes.PermissionDenied, status.Code(err))
		// 		return assert.ErrorContains(t, err, service.ErrPermissionDenied.Error())
		// 	},
		// },
		{
			name: "Error: remove evaluation results",
			fields: fields{
				svc: NewService(
					WithAuthorizationStrategy(&service.AuthorizationStrategyAllowAll{}),
					WithStorage(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						assert.NoError(t, s.Create(orchestratortest.NewAuditScope("", testdata.MockAuditScopeID1, "", testdata.MockAuditScopeName1)))
						assert.NoError(t, s.Create(orchestratortest.NewAuditScope("", testdata.MockAuditScopeID2, "", testdata.MockAuditScopeName2)))
					})),
				),
			},
			args: args{
				ctx: nil,
				req: &orchestrator.RemoveAuditScopeRequest{
					AuditScopeId:            testdata.MockAuditScopeID1,
					RemoveEvaluationResults: true,
				},
			},
			wantResponse: assert.Nil[*emptypb.Empty],
			wantSvc: func(t *testing.T, got *Service) bool {
				return true
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.Internal, status.Code(err))
				return assert.ErrorContains(t, err, "database error")
			},
		},
		{
			name: "Happy path: remove evaluation results",
			fields: fields{
				svc: NewService(
					WithAuthorizationStrategy(&service.AuthorizationStrategyAllowAll{}),
					WithStorage(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						assert.NoError(t, s.Create(orchestratortest.NewAuditScope("", testdata.MockAuditScopeID1, "", testdata.MockAuditScopeName1)))
						assert.NoError(t, s.Create(orchestratortest.NewAuditScope("", testdata.MockAuditScopeID2, "", testdata.MockAuditScopeName2)))
						assert.NoError(t, s.Create(evaluationtest.MockEvaluationResults))
						assert.NoError(t, s.Create(&evaluation.EvaluationResult{
							Id:                   testdata.MockEvaluationResult10ID,
							Timestamp:            timestamppb.New(time.Unix(3, 0)),
							TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
							AuditScopeId:         testdata.MockAuditScopeID2,
							ControlCategoryName:  testdata.MockCategoryName,
							ControlCatalogId:     testdata.MockCatalogID1,
							Status:               evaluation.EvaluationStatus_EVALUATION_STATUS_COMPLIANT,
							ControlId:            testdata.MockSubControlID21,
							ParentControlId:      util.Ref(testdata.MockControlID2),
							AssessmentResultIds:  []string{},
						}))
					})),
				),
			},
			args: args{
				ctx: nil,
				req: &orchestrator.RemoveAuditScopeRequest{
					AuditScopeId:            testdata.MockAuditScopeID1,
					RemoveEvaluationResults: true,
				},
			},
			wantResponse: assert.NotNil[*emptypb.Empty],
			wantSvc: func(t *testing.T, got *Service) bool {
				// Verify that audit scope with ID 2 is still in the DB (by counting the number of occurrences = 1)
				n, err := got.storage.Count(&orchestrator.AuditScope{}, "id=?", testdata.MockAuditScopeID2)
				assert.NoError(t, err)
				assert.Equal(t, 1, n)

				// Verify that audit scope with ID 1 is not in the DB anymore (by counting the number of occurrences = 0)
				n, err = got.storage.Count(&orchestrator.AuditScope{}, "id=?", testdata.MockAuditScopeID1)
				assert.NoError(t, err)
				assert.Equal(t, 0, n)

				// Verify that evaluation results with audit scope ID 1 is not in the DB anymore (by counting the number of occurrences = 0)
				n, err = got.storage.Count(&evaluation.EvaluationResult{}, "audit_scope_id=?", testdata.MockAuditScopeID1)
				assert.NoError(t, err)
				return assert.Equal(t, 0, n)

			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: with authorization allAllowed",
			fields: fields{
				svc: NewService(
					WithAuthorizationStrategy(&service.AuthorizationStrategyAllowAll{}),
					WithStorage(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						assert.NoError(t, s.Create(orchestratortest.NewAuditScope("", testdata.MockAuditScopeID1, "", testdata.MockAuditScopeName1)))
						assert.NoError(t, s.Create(orchestratortest.NewAuditScope("", testdata.MockAuditScopeID2, "", testdata.MockAuditScopeName2)))
					})),
				),
			},
			args: args{
				ctx: nil,
				req: &orchestrator.RemoveAuditScopeRequest{
					AuditScopeId: testdata.MockAuditScopeID1,
				},
			},
			wantResponse: assert.NotNil[*emptypb.Empty],
			wantSvc: func(t *testing.T, got *Service) bool {
				// Verify that audit scope with ID 2 is still in the DB (by counting the number of occurrences = 1)
				n, err := got.storage.Count(&orchestrator.AuditScope{}, "id=?", testdata.MockAuditScopeID2)
				assert.NoError(t, err)
				assert.Equal(t, 1, n)

				// Verify that the default audit scope isn't in the DB anymore (occurrences = 0)
				x, err := got.storage.Count(&orchestrator.AuditScope{}, "id=?", testdata.MockAuditScopeID1)
				assert.NoError(t, err)
				return assert.Equal(t, 0, x)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: with authorization for audit scopes with a certain specific target of evaluation",
			fields: fields{
				svc: NewService(
					WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID2)),
					WithStorage(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						assert.NoError(t, s.Create(orchestratortest.NewAuditScope("", testdata.MockAuditScopeID1, testdata.MockTargetOfEvaluationID1, testdata.MockAuditScopeName1)))
						assert.NoError(t, s.Create(orchestratortest.NewAuditScope("", testdata.MockAuditScopeID2, testdata.MockTargetOfEvaluationID2, testdata.MockAuditScopeName2)))
					})),
				),
			},
			args: args{
				ctx: nil,
				req: &orchestrator.RemoveAuditScopeRequest{
					AuditScopeId: testdata.MockAuditScopeID2,
				},
			},
			wantResponse: assert.NotNil[*emptypb.Empty],
			wantSvc: func(t *testing.T, got *Service) bool {
				// Verify that audit scope with ID 2 is still in the DB (by counting the number of occurrences = 1)
				n, err := got.storage.Count(&orchestrator.AuditScope{}, "id=?", testdata.MockAuditScopeID1)
				assert.NoError(t, err)
				assert.Equal(t, 1, n)

				// Verify that the default audit scope isn't in the DB anymore (occurrences = 0)
				x, err := got.storage.Count(&orchestrator.AuditScope{}, "id=?", testdata.MockAuditScopeID2)
				assert.NoError(t, err)
				return assert.Equal(t, 0, x)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := tt.fields.svc.RemoveAuditScope(tt.args.ctx, tt.args.req)

			tt.wantResponse(t, res)
			tt.wantSvc(t, tt.fields.svc)
			tt.wantErr(t, err)
		})
	}
}

func TestService_UpdateAuditScope(t *testing.T) {
	type fields struct {
		svc *Service
	}
	type args struct {
		ctx context.Context
		req *orchestrator.UpdateAuditScopeRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes assert.Want[*orchestrator.AuditScope]
		wantSvc assert.Want[*Service]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Error: validation error",
			fields: fields{
				svc: NewService(),
			},
			args: args{
				ctx: nil,
				req: nil,
			},
			wantRes: assert.Nil[*orchestrator.AuditScope],
			wantSvc: func(t *testing.T, got *Service) bool {
				return true
			},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, api.ErrEmptyRequest.Error())
			},
		},
		{
			name: "Error: permission denied",
			fields: fields{
				svc: NewService(
					WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(false)),
				),
			},
			args: args{
				ctx: nil,
				req: &orchestrator.UpdateAuditScopeRequest{AuditScope: orchestratortest.MockAuditScopeCertTargetID1},
			},
			wantRes: assert.Nil[*orchestrator.AuditScope],
			wantSvc: func(t *testing.T, got *Service) bool {
				return true
			},
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				assert.Equal(t, codes.PermissionDenied, status.Code(err))
				return assert.ErrorContains(t, err, service.ErrPermissionDenied.Error())
			},
		},
		{
			name: "Error: permission denied: with authorization for audit scopes with a certain specific target of evaluation",
			fields: fields{
				svc: NewService(
					WithStorage(testutil.NewInMemoryStorage(t)),
					WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID2)),
				),
			},
			args: args{
				ctx: context.Background(),
				req: &orchestrator.UpdateAuditScopeRequest{
					AuditScope: orchestratortest.MockAuditScopeCertTargetID1,
				},
			},
			wantRes: assert.Nil[*orchestrator.AuditScope],
			wantSvc: func(t *testing.T, got *Service) bool {
				return true
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.PermissionDenied, status.Code(err))
				return assert.ErrorContains(t, err, service.ErrPermissionDenied.Error())
			},
		},
		{
			name: "Error: audit scope not found",
			fields: fields{
				svc: NewService(
					WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(
						true)),
					// Create empty storage => No audit scope can be found
					WithStorage(testutil.NewInMemoryStorage(t))),
			},
			args: args{
				ctx: nil,
				req: &orchestrator.UpdateAuditScopeRequest{
					AuditScope: orchestratortest.MockAuditScopeCertTargetID1,
				},
			},
			wantRes: assert.Nil[*orchestrator.AuditScope],
			wantSvc: func(t *testing.T, got *Service) bool {
				return true
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.NotFound, status.Code(err))
				return assert.ErrorContains(t, err, api.ErrAuditScopeNotFound.Error())
			},
		},
		{
			name: "Error: constraint failed",
			fields: fields{
				svc: NewService(
					WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(true)),
					WithStorage(&testutil.StorageWithError{UpdateErr: persistence.ErrConstraintFailed}),
				),
			},
			args: args{req: &orchestrator.UpdateAuditScopeRequest{
				AuditScope: orchestratortest.MockAuditScopeCertTargetID1,
			}},
			wantRes: assert.Nil[*orchestrator.AuditScope],
			wantSvc: func(t *testing.T, svc *Service) bool {
				n, err := svc.storage.Count(&orchestrator.AuditScope{})
				assert.NoError(t, err)
				return assert.Equal(t, 0, n)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.NotFound, status.Code(err))
				return assert.ErrorContains(t, err, persistence.ErrConstraintFailed.Error())
			},
		},
		{
			name: "Error: database error",
			fields: fields{
				svc: NewService(
					WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(true)),
					WithStorage(&testutil.StorageWithError{UpdateErr: gorm.ErrInvalidDB}),
				),
			},
			args: args{req: &orchestrator.UpdateAuditScopeRequest{
				AuditScope: orchestratortest.MockAuditScopeCertTargetID1,
			}},
			wantRes: assert.Nil[*orchestrator.AuditScope],
			wantSvc: func(t *testing.T, svc *Service) bool {
				n, err := svc.storage.Count(&orchestrator.AuditScope{})
				assert.NoError(t, err)
				return assert.Equal(t, 0, n)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.Internal, status.Code(err))
				return assert.ErrorContains(t, err, persistence.ErrDatabase.Error())
			},
		},
		{
			name: "Happy path: with authorization for audit scopes with a certain specific target of evaluation",
			fields: fields{
				svc: NewService(
					WithStorage(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						err := s.Create(&orchestrator.AuditScope{
							Id:                   testdata.MockAuditScopeID1,
							TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
							CatalogId:            testdata.MockCatalogID1,
							AssuranceLevel:       &testdata.AssuranceLevelHigh,
						})
						assert.NoError(t, err)
					})),
					WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(false, testdata.MockTargetOfEvaluationID1)),
				),
			},
			args: args{
				ctx: context.Background(),
				req: &orchestrator.UpdateAuditScopeRequest{
					AuditScope: orchestratortest.MockAuditScopeCertTargetID1,
				},
			},
			wantRes: func(t *testing.T, got *orchestrator.AuditScope) bool {
				return assert.Equal(t, orchestratortest.MockAuditScopeCertTargetID1, got)
			},
			wantSvc: func(t *testing.T, svc *Service) bool {
				res := &orchestrator.AuditScope{}

				// Check if audit scope is updated in the DB
				err := svc.storage.Get(res, "id = ?", testdata.MockAuditScopeID1)
				assert.NoError(t, err)
				return assert.Equal(t, orchestratortest.MockAuditScopeCertTargetID1, res)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: with authorization allAllowed",
			fields: fields{
				svc: NewService(
					WithStorage(testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
						err := s.Create(&orchestrator.AuditScope{
							Id:                   testdata.MockAuditScopeID1,
							TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
							CatalogId:            testdata.MockCatalogID1,
							AssuranceLevel:       &testdata.AssuranceLevelHigh,
						})
						assert.NoError(t, err)
					})),
					WithAuthorizationStrategy(servicetest.NewAuthorizationStrategy(true)),
				),
			},
			args: args{
				ctx: context.Background(),
				req: &orchestrator.UpdateAuditScopeRequest{
					AuditScope: orchestratortest.MockAuditScopeCertTargetID1,
				},
			},
			wantRes: func(t *testing.T, got *orchestrator.AuditScope) bool {
				return assert.Equal(t, orchestratortest.MockAuditScopeCertTargetID1, got)
			},
			wantSvc: func(t *testing.T, svc *Service) bool {
				res := &orchestrator.AuditScope{}

				// Check if audit scope is updated in the DB
				err := svc.storage.Get(res, "id = ?", testdata.MockAuditScopeID1)
				assert.NoError(t, err)
				return assert.Equal(t, orchestratortest.MockAuditScopeCertTargetID1, res)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.fields.svc.UpdateAuditScope(tt.args.ctx, tt.args.req)

			tt.wantRes(t, gotRes)
			tt.wantSvc(t, tt.fields.svc)
			tt.wantErr(t, err)
		})
	}
}
