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

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/orchestratortest"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/persistence/gorm"
	"clouditor.io/clouditor/v2/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestService_CreateAuditScope(t *testing.T) {

	catalogWithoutAssuranceLevelList := orchestratortest.NewCatalog()
	catalogWithoutAssuranceLevelList.AssuranceLevels = []string{}

	type fields struct {
		CertificationTargetHooks []orchestrator.CertificationTargetHookFunc
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
			name: "Empty request",
			args: args{
				ctx: context.Background(),
				req: &orchestrator.CreateAuditScopeRequest{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "empty request")
			},
			want: assert.Nil[*orchestrator.AuditScope],
		},
		{
			name: "Invalid request",
			args: args{
				ctx: context.Background(),
				req: &orchestrator.CreateAuditScopeRequest{
					AuditScope: &orchestrator.AuditScope{
						CatalogId: "testcatalog",
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "audit_scope.id: value length must be at least 1 characters")
			},
			want: assert.Nil[*orchestrator.AuditScope],
		},
		{
			name: "Error getting catalog",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {}),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{req: &orchestrator.CreateAuditScopeRequest{
				AuditScope: orchestratortest.NewAuditScope(testdata.AssuranceLevelBasic),
			}},
			want: assert.Nil[*orchestrator.AuditScope],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "invalid catalog or certification target")
			},
		},
		{
			name: "valid and assurance level not set",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Create(&orchestrator.CertificationTarget{Id: testdata.MockCertificationTargetID1})
					assert.NoError(t, err)

					err = s.Create(catalogWithoutAssuranceLevelList)
					assert.NoError(t, err)
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{req: &orchestrator.CreateAuditScopeRequest{
				AuditScope: orchestratortest.NewAuditScope(""),
			}},
			wantSvc: func(t *testing.T, got *Service) bool {
				// We want to assert that certain things happened in our database
				var auditScopes []*orchestrator.AuditScope
				// for join tables, do not use preload (which is the default)
				err := got.storage.List(&auditScopes, "", false, 0, -1, gorm.WithoutPreload())
				if !assert.NoError(t, err) {
					return false
				}
				if !assert.Equal(t, 1, len(auditScopes)) {
					return false
				}

				var service orchestrator.CertificationTarget
				err = got.storage.Get(&service, "id = ?", testdata.MockCertificationTargetID1)
				if !assert.NoError(t, err) {
					return false
				}

				return assert.Equal(t, 1, len(service.CatalogsInScope))

			},
			want:    assert.AnyValue[*orchestrator.AuditScope],
			wantErr: assert.NoError,
		},
		{
			name: "valid and assurance level set",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Create(&orchestrator.CertificationTarget{Id: testdata.MockCertificationTargetID1})
					assert.NoError(t, err)

					err = s.Create(orchestratortest.NewCatalog())
					assert.NoError(t, err)
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{req: &orchestrator.CreateAuditScopeRequest{
				AuditScope: orchestratortest.NewAuditScope(testdata.AssuranceLevelBasic),
			}},
			wantSvc: func(t *testing.T, got *Service) bool {
				// We want to assert that certain things happened in our database
				var auditScopes []*orchestrator.AuditScope
				// for join tables, do not use preload (which is the default)
				err := got.storage.List(&auditScopes, "", false, 0, -1, gorm.WithoutPreload())
				if !assert.NoError(t, err) {
					return false
				}
				if !assert.Equal(t, 1, len(auditScopes)) {
					return false
				}

				var service orchestrator.CertificationTarget
				err = got.storage.Get(&service, "id = ?", testdata.MockCertificationTargetID1)
				if !assert.NoError(t, err) {
					return false
				}

				return assert.Equal(t, 1, len(service.CatalogsInScope))
			},
			want:    assert.AnyValue[*orchestrator.AuditScope],
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				CertificationTargetHooks: tt.fields.CertificationTargetHooks,
				AssessmentResultHooks:    tt.fields.AssessmentResultHooks,
				storage:                  tt.fields.storage,
				metricsFile:              tt.fields.metricsFile,
				loadMetricsFunc:          tt.fields.loadMetricsFunc,
				catalogsFolder:           tt.fields.catalogsFolder,
				loadCatalogsFunc:         tt.fields.loadCatalogsFunc,
				events:                   tt.fields.events,
				authz:                    tt.fields.authz,
			}

			gotRes, err := svc.CreateAuditScope(tt.args.ctx, tt.args.req)
			tt.wantErr(t, err)
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
			name: "empty request",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args:         args{req: nil},
			wantResponse: assert.Nil[*orchestrator.AuditScope],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "empty request")
			},
		},
		{
			name: "invalid request",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{req: &orchestrator.GetAuditScopeRequest{
				CertificationTargetId: "",
			}},
			wantResponse: assert.Nil[*orchestrator.AuditScope],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "certification_target_id: value is empty, which is not a valid UUID ")
			},
		},
		{
			name: "auditScope not found",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Create(&orchestrator.CertificationTarget{Id: testdata.MockCertificationTargetID1})
					assert.NoError(t, err)

					err = s.Create(orchestratortest.NewCatalog())
					assert.NoError(t, err)

					err = s.Create(orchestratortest.NewAuditScope(testdata.AssuranceLevelBasic))
					assert.NoError(t, err)
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{req: &orchestrator.GetAuditScopeRequest{
				CertificationTargetId: testdata.MockCertificationTargetID2,
				CatalogId:             testdata.MockCatalogID,
			}},
			wantResponse: assert.Nil[*orchestrator.AuditScope],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.NotFound, status.Code(err))
				return assert.ErrorContains(t, err, "Audit Scope not found")
			},
		},
		{
			name: "valid auditScope",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Create(&orchestrator.CertificationTarget{Id: testdata.MockCertificationTargetID1})
					assert.NoError(t, err)

					err = s.Create(orchestratortest.NewCatalog())
					assert.NoError(t, err)

					err = s.Create(orchestratortest.NewAuditScope(testdata.AssuranceLevelBasic))
					assert.NoError(t, err)
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{req: &orchestrator.GetAuditScopeRequest{
				CertificationTargetId: testdata.MockCertificationTargetID1,
				CatalogId:             testdata.MockCatalogID,
			}},
			wantResponse: func(t *testing.T, got *orchestrator.AuditScope) bool {
				want := orchestratortest.NewAuditScope(testdata.AssuranceLevelBasic)

				return assert.NoError(t, api.Validate(got)) &&
					assert.Equal(t, want.CertificationTargetId, got.CertificationTargetId) &&
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

func TestService_ListAuditScopes(t *testing.T) {
	var (
		listAuditScopesResponse *orchestrator.ListAuditScopesResponse
		err                     error
	)

	orchestratorService := NewService()
	err = orchestratorService.storage.Create(&orchestrator.CertificationTarget{Id: testdata.MockCertificationTargetID1})
	assert.NoError(t, err)
	err = orchestratorService.storage.Create(orchestratortest.NewCatalog())
	assert.NoError(t, err)

	// 1st case: No Audit Scopes stored
	listAuditScopesResponse, err = orchestratorService.ListAuditScopes(context.Background(), &orchestrator.ListAuditScopesRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listAuditScopesResponse.AuditScopes)
	assert.Empty(t, listAuditScopesResponse.AuditScopes)

	// 2nd case: One Audit Scope stored
	err = orchestratorService.storage.Create(orchestratortest.NewAuditScope(testdata.AssuranceLevelBasic))
	assert.NoError(t, err)

	listAuditScopesResponse, err = orchestratorService.ListAuditScopes(context.Background(), &orchestrator.ListAuditScopesRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listAuditScopesResponse.AuditScopes)
	assert.NotEmpty(t, listAuditScopesResponse.AuditScopes)
	assert.NoError(t, api.Validate(listAuditScopesResponse.AuditScopes[0]))
	assert.Equal(t, 1, len(listAuditScopesResponse.AuditScopes))
}

func TestService_UpdateAuditScope(t *testing.T) {
	var (
		auditScope *orchestrator.AuditScope
		err        error
	)
	orchestratorService := NewService()
	err = orchestratorService.storage.Create(&orchestrator.CertificationTarget{Id: testdata.MockCertificationTargetID1})
	assert.NoError(t, err)
	err = orchestratorService.storage.Create(orchestratortest.NewCatalog())
	assert.NoError(t, err)

	// 1st case: Audit Scope is nil
	_, err = orchestratorService.UpdateAuditScope(context.Background(), &orchestrator.UpdateAuditScopeRequest{
		AuditScope: nil,
	})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// 2nd case: Ids are empty
	_, err = orchestratorService.UpdateAuditScope(context.Background(), &orchestrator.UpdateAuditScopeRequest{
		AuditScope: &orchestrator.AuditScope{},
	})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.ErrorContains(t, err, "audit_scope.certification_target_id: value is empty, which is not a valid UUID")

	// 3rd case: Audit Scope not found since there are no Audit Scopes
	_, err = orchestratorService.UpdateAuditScope(context.Background(), &orchestrator.UpdateAuditScopeRequest{
		AuditScope: orchestratortest.NewAuditScope(testdata.AssuranceLevelBasic),
	})
	assert.Equal(t, codes.NotFound, status.Code(err))

	// 4th case: Audit Scope updated successfully
	err = orchestratorService.storage.Create(orchestratortest.NewAuditScope(testdata.AssuranceLevelBasic))
	assert.NoError(t, err)

	// update the auditScope's assurance level and send the update request
	auditScope, err = orchestratorService.UpdateAuditScope(context.Background(), &orchestrator.UpdateAuditScopeRequest{
		AuditScope: &orchestrator.AuditScope{
			CertificationTargetId: testdata.MockCertificationTargetID1,
			CatalogId:             testdata.MockCatalogID,
			AssuranceLevel:        &testdata.AssuranceLevelBasic,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, auditScope)
	assert.NoError(t, api.Validate(auditScope))
	assert.Equal(t, &testdata.AssuranceLevelBasic, auditScope.AssuranceLevel)
}

func TestService_RemoveAuditScope(t *testing.T) {
	var (
		err                     error
		listAuditScopesResponse *orchestrator.ListAuditScopesResponse
	)
	orchestratorService := NewService()
	err = orchestratorService.storage.Create(&orchestrator.CertificationTarget{Id: testdata.MockCertificationTargetID1})
	assert.NoError(t, err)
	err = orchestratorService.storage.Create(orchestratortest.NewCatalog())
	assert.NoError(t, err)

	// 1st case: Empty ID error
	_, err = orchestratorService.RemoveAuditScope(context.Background(), &orchestrator.RemoveAuditScopeRequest{
		CertificationTargetId: "",
		CatalogId:             "",
	})
	assert.Error(t, err)
	assert.Equal(t, status.Code(err), codes.InvalidArgument)

	// 2nd case: ErrRecordNotFound
	_, err = orchestratorService.RemoveAuditScope(context.Background(), &orchestrator.RemoveAuditScopeRequest{
		CertificationTargetId: testdata.MockCertificationTargetID1,
		CatalogId:             "0000",
	})
	assert.Error(t, err)
	assert.Equal(t, status.Code(err), codes.NotFound)

	// 3rd case: Record removed successfully
	err = orchestratorService.storage.Create(orchestratortest.NewAuditScope(testdata.AssuranceLevelBasic))
	assert.NoError(t, err)

	// Verify that there is a record for Audit Scope in the DB
	listAuditScopesResponse, err = orchestratorService.ListAuditScopes(context.Background(), &orchestrator.ListAuditScopesRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listAuditScopesResponse.AuditScopes)
	assert.NoError(t, api.Validate(listAuditScopesResponse.AuditScopes[0]))
	assert.Equal(t, 1, len(listAuditScopesResponse.AuditScopes))

	// Remove record
	_, err = orchestratorService.RemoveAuditScope(context.Background(), &orchestrator.RemoveAuditScopeRequest{
		CertificationTargetId: testdata.MockCertificationTargetID1,
		CatalogId:             testdata.MockCatalogID,
	})
	assert.NoError(t, err)

	// There is no record for Audit Scope in the DB
	listAuditScopesResponse, err = orchestratorService.ListAuditScopes(context.Background(), &orchestrator.ListAuditScopesRequest{})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(listAuditScopesResponse.AuditScopes))
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
			name: "Store first assessment result to the map",
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.UpdateAuditScopeRequest{
					AuditScope: &orchestrator.AuditScope{
						CertificationTargetId: testdata.MockCertificationTargetID1,
						CatalogId:             testdata.MockCatalogID,
						AssuranceLevel:        &testdata.AssuranceLevelSubstantial,
					},
				},
			},
			wantErr: assert.Nil[error],
			wantRes: &orchestrator.AuditScope{
				CertificationTargetId: testdata.MockCertificationTargetID1,
				CatalogId:             testdata.MockCatalogID,
				AssuranceLevel:        &testdata.AssuranceLevelSubstantial,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hookCallCounter = 0

			// Create service
			s := svc
			err := s.storage.Create(&orchestrator.CertificationTarget{Id: testdata.MockCertificationTargetID1})
			assert.NoError(t, err)

			// Create catalog
			err = s.storage.Create(orchestratortest.NewCatalog())
			assert.NoError(t, err)

			// Create new Audit Scope
			err = s.storage.Create(orchestratortest.NewAuditScope(testdata.AssuranceLevelBasic))
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
