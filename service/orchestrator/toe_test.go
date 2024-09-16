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

func TestService_CreateTargetOfEvaluation(t *testing.T) {

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
		req *orchestrator.CreateTargetOfEvaluationRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.Want[*orchestrator.TargetOfEvaluation]
		wantSvc assert.Want[*Service]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Invalid request",
			args: args{
				ctx: context.Background(),
				req: &orchestrator.CreateTargetOfEvaluationRequest{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "target_of_evaluation: value is required")
			},
			want: assert.Nil[*orchestrator.TargetOfEvaluation],
		},
		{
			name: "Error getting catalog",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {}),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args: args{req: &orchestrator.CreateTargetOfEvaluationRequest{
				TargetOfEvaluation: orchestratortest.NewTargetOfEvaluation(testdata.AssuranceLevelBasic),
			}},
			want: assert.Nil[*orchestrator.TargetOfEvaluation],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "invalid catalog or cloud service")
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
			args: args{req: &orchestrator.CreateTargetOfEvaluationRequest{
				TargetOfEvaluation: orchestratortest.NewTargetOfEvaluation(""),
			}},
			wantSvc: func(t *testing.T, got *Service) bool {
				// We want to assert that certain things happened in our database
				var toes []*orchestrator.TargetOfEvaluation
				// for join tables, do not use preload (which is the default)
				err := got.storage.List(&toes, "", false, 0, -1, gorm.WithoutPreload())
				if !assert.NoError(t, err) {
					return false
				}
				if !assert.Equal(t, 1, len(toes)) {
					return false
				}

				var service orchestrator.CertificationTarget
				err = got.storage.Get(&service, "id = ?", testdata.MockCertificationTargetID1)
				if !assert.NoError(t, err) {
					return false
				}

				return assert.Equal(t, 1, len(service.CatalogsInScope))

			},
			want:    assert.AnyValue[*orchestrator.TargetOfEvaluation],
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
			args: args{req: &orchestrator.CreateTargetOfEvaluationRequest{
				TargetOfEvaluation: orchestratortest.NewTargetOfEvaluation(testdata.AssuranceLevelBasic),
			}},
			wantSvc: func(t *testing.T, got *Service) bool {
				// We want to assert that certain things happened in our database
				var toes []*orchestrator.TargetOfEvaluation
				// for join tables, do not use preload (which is the default)
				err := got.storage.List(&toes, "", false, 0, -1, gorm.WithoutPreload())
				if !assert.NoError(t, err) {
					return false
				}
				if !assert.Equal(t, 1, len(toes)) {
					return false
				}

				var service orchestrator.CertificationTarget
				err = got.storage.Get(&service, "id = ?", testdata.MockCertificationTargetID1)
				if !assert.NoError(t, err) {
					return false
				}

				return assert.Equal(t, 1, len(service.CatalogsInScope))
			},
			want:    assert.AnyValue[*orchestrator.TargetOfEvaluation],
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

			gotRes, err := svc.CreateTargetOfEvaluation(tt.args.ctx, tt.args.req)
			tt.wantErr(t, err)
			tt.want(t, gotRes)
		})
	}
}

func TestService_GetTargetOfEvaluation(t *testing.T) {
	type fields struct {
		storage persistence.Storage
		authz   service.AuthorizationStrategy
	}
	type args struct {
		req *orchestrator.GetTargetOfEvaluationRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantResponse assert.Want[*orchestrator.TargetOfEvaluation]
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "empty request",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   servicetest.NewAuthorizationStrategy(true),
			},
			args:         args{req: nil},
			wantResponse: assert.Nil[*orchestrator.TargetOfEvaluation],
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
			args: args{req: &orchestrator.GetTargetOfEvaluationRequest{
				CertificationTargetId: "",
			}},
			wantResponse: assert.Nil[*orchestrator.TargetOfEvaluation],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "certification_target_id: value is empty, which is not a valid UUID ")
			},
		},
		{
			name: "toe not found",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Create(&orchestrator.CertificationTarget{Id: testdata.MockCertificationTargetID1})
					assert.NoError(t, err)

					err = s.Create(orchestratortest.NewCatalog())
					assert.NoError(t, err)

					err = s.Create(orchestratortest.NewTargetOfEvaluation(testdata.AssuranceLevelBasic))
					assert.NoError(t, err)
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{req: &orchestrator.GetTargetOfEvaluationRequest{
				CertificationTargetId: testdata.MockCertificationTargetID2,
				CatalogId:             testdata.MockCatalogID,
			}},
			wantResponse: assert.Nil[*orchestrator.TargetOfEvaluation],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.NotFound, status.Code(err))
				return assert.ErrorContains(t, err, "ToE not found")
			},
		},
		{
			name: "valid toe",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Create(&orchestrator.CertificationTarget{Id: testdata.MockCertificationTargetID1})
					assert.NoError(t, err)

					err = s.Create(orchestratortest.NewCatalog())
					assert.NoError(t, err)

					err = s.Create(orchestratortest.NewTargetOfEvaluation(testdata.AssuranceLevelBasic))
					assert.NoError(t, err)
				}),
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{req: &orchestrator.GetTargetOfEvaluationRequest{
				CertificationTargetId: testdata.MockCertificationTargetID1,
				CatalogId:             testdata.MockCatalogID,
			}},
			wantResponse: func(t *testing.T, got *orchestrator.TargetOfEvaluation) bool {
				want := orchestratortest.NewTargetOfEvaluation(testdata.AssuranceLevelBasic)

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
			res, err := orchestratorService.GetTargetOfEvaluation(context.Background(), tt.args.req)
			tt.wantErr(t, err)
			tt.wantResponse(t, res)
		})
	}
}

func TestService_ListTargetsOfEvaluation(t *testing.T) {
	var (
		listTargetsOfEvaluationResponse *orchestrator.ListTargetsOfEvaluationResponse
		err                             error
	)

	orchestratorService := NewService()
	err = orchestratorService.storage.Create(&orchestrator.CertificationTarget{Id: testdata.MockCertificationTargetID1})
	assert.NoError(t, err)
	err = orchestratorService.storage.Create(orchestratortest.NewCatalog())
	assert.NoError(t, err)

	// 1st case: No ToEs stored
	listTargetsOfEvaluationResponse, err = orchestratorService.ListTargetsOfEvaluation(context.Background(), &orchestrator.ListTargetsOfEvaluationRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listTargetsOfEvaluationResponse.TargetOfEvaluation)
	assert.Empty(t, listTargetsOfEvaluationResponse.TargetOfEvaluation)

	// 2nd case: One ToE stored
	err = orchestratorService.storage.Create(orchestratortest.NewTargetOfEvaluation(testdata.AssuranceLevelBasic))
	assert.NoError(t, err)

	listTargetsOfEvaluationResponse, err = orchestratorService.ListTargetsOfEvaluation(context.Background(), &orchestrator.ListTargetsOfEvaluationRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listTargetsOfEvaluationResponse.TargetOfEvaluation)
	assert.NotEmpty(t, listTargetsOfEvaluationResponse.TargetOfEvaluation)
	assert.NoError(t, api.Validate(listTargetsOfEvaluationResponse.TargetOfEvaluation[0]))
	assert.Equal(t, 1, len(listTargetsOfEvaluationResponse.TargetOfEvaluation))
}

func TestService_UpdateTargetOfEvaluation(t *testing.T) {
	var (
		toe *orchestrator.TargetOfEvaluation
		err error
	)
	orchestratorService := NewService()
	err = orchestratorService.storage.Create(&orchestrator.CertificationTarget{Id: testdata.MockCertificationTargetID1})
	assert.NoError(t, err)
	err = orchestratorService.storage.Create(orchestratortest.NewCatalog())
	assert.NoError(t, err)

	// 1st case: ToE is nil
	_, err = orchestratorService.UpdateTargetOfEvaluation(context.Background(), &orchestrator.UpdateTargetOfEvaluationRequest{
		TargetOfEvaluation: nil,
	})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// 2nd case: Ids are empty
	_, err = orchestratorService.UpdateTargetOfEvaluation(context.Background(), &orchestrator.UpdateTargetOfEvaluationRequest{
		TargetOfEvaluation: &orchestrator.TargetOfEvaluation{},
	})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.ErrorContains(t, err, "target_of_evaluation.certification_target_id: value is empty, which is not a valid UUID")

	// 3rd case: ToE not found since there are no ToEs
	_, err = orchestratorService.UpdateTargetOfEvaluation(context.Background(), &orchestrator.UpdateTargetOfEvaluationRequest{
		TargetOfEvaluation: orchestratortest.NewTargetOfEvaluation(testdata.AssuranceLevelBasic),
	})
	assert.Equal(t, codes.NotFound, status.Code(err))

	// 4th case: ToE updated successfully
	err = orchestratorService.storage.Create(orchestratortest.NewTargetOfEvaluation(testdata.AssuranceLevelBasic))
	assert.NoError(t, err)

	// update the toe's assurance level and send the update request
	toe, err = orchestratorService.UpdateTargetOfEvaluation(context.Background(), &orchestrator.UpdateTargetOfEvaluationRequest{
		TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
			CertificationTargetId: testdata.MockCertificationTargetID1,
			CatalogId:             testdata.MockCatalogID,
			AssuranceLevel:        &testdata.AssuranceLevelBasic,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, toe)
	assert.NoError(t, api.Validate(toe))
	assert.Equal(t, &testdata.AssuranceLevelBasic, toe.AssuranceLevel)
}

func TestService_RemoveTargetOfEvaluation(t *testing.T) {
	var (
		err                             error
		listTargetsOfEvaluationResponse *orchestrator.ListTargetsOfEvaluationResponse
	)
	orchestratorService := NewService()
	err = orchestratorService.storage.Create(&orchestrator.CertificationTarget{Id: testdata.MockCertificationTargetID1})
	assert.NoError(t, err)
	err = orchestratorService.storage.Create(orchestratortest.NewCatalog())
	assert.NoError(t, err)

	// 1st case: Empty ID error
	_, err = orchestratorService.RemoveTargetOfEvaluation(context.Background(), &orchestrator.RemoveTargetOfEvaluationRequest{
		CertificationTargetId: "",
		CatalogId:             "",
	})
	assert.Error(t, err)
	assert.Equal(t, status.Code(err), codes.InvalidArgument)

	// 2nd case: ErrRecordNotFound
	_, err = orchestratorService.RemoveTargetOfEvaluation(context.Background(), &orchestrator.RemoveTargetOfEvaluationRequest{
		CertificationTargetId: testdata.MockCertificationTargetID1,
		CatalogId:             "0000",
	})
	assert.Error(t, err)
	assert.Equal(t, status.Code(err), codes.NotFound)

	// 3rd case: Record removed successfully
	err = orchestratorService.storage.Create(orchestratortest.NewTargetOfEvaluation(testdata.AssuranceLevelBasic))
	assert.NoError(t, err)

	// Verify that there is a record for ToE in the DB
	listTargetsOfEvaluationResponse, err = orchestratorService.ListTargetsOfEvaluation(context.Background(), &orchestrator.ListTargetsOfEvaluationRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listTargetsOfEvaluationResponse.TargetOfEvaluation)
	assert.NoError(t, api.Validate(listTargetsOfEvaluationResponse.TargetOfEvaluation[0]))
	assert.Equal(t, 1, len(listTargetsOfEvaluationResponse.TargetOfEvaluation))

	// Remove record
	_, err = orchestratorService.RemoveTargetOfEvaluation(context.Background(), &orchestrator.RemoveTargetOfEvaluationRequest{
		CertificationTargetId: testdata.MockCertificationTargetID1,
		CatalogId:             testdata.MockCatalogID,
	})
	assert.NoError(t, err)

	// There is no record for ToE in the DB
	listTargetsOfEvaluationResponse, err = orchestratorService.ListTargetsOfEvaluation(context.Background(), &orchestrator.ListTargetsOfEvaluationRequest{})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(listTargetsOfEvaluationResponse.TargetOfEvaluation))
}

func TestToeHook(t *testing.T) {
	var (
		hookCallCounter = 0
		wg              sync.WaitGroup
	)
	wg.Add(2)

	firstHookFunction := func(ctx context.Context, event *orchestrator.TargetOfEvaluationChangeEvent, err error) {
		hookCallCounter++
		log.Println("Hello from inside the first toe hook function")

		// Checking the status
		// UpdateTargetOfEvaluation is called, so the status must be TOE_UPDATE
		if *event.GetType().Enum() != orchestrator.TargetOfEvaluationChangeEvent_TYPE_TARGET_OF_EVALUATION_UPDATED {
			return
		}

		wg.Done()
	}

	secondHookFunction := func(ctx context.Context, event *orchestrator.TargetOfEvaluationChangeEvent, err error) {
		hookCallCounter++
		log.Println("Hello from inside the second toe hook function")

		wg.Done()
	}

	svc := NewService()
	svc.RegisterToeHook(firstHookFunction)
	svc.RegisterToeHook(secondHookFunction)

	// Check if first hook is registered
	funcName1 := runtime.FuncForPC(reflect.ValueOf(svc.toeHooks[0]).Pointer()).Name()
	funcName2 := runtime.FuncForPC(reflect.ValueOf(firstHookFunction).Pointer()).Name()
	assert.Equal(t, funcName1, funcName2)

	// Check if second hook is registered
	funcName1 = runtime.FuncForPC(reflect.ValueOf(svc.toeHooks[1]).Pointer()).Name()
	funcName2 = runtime.FuncForPC(reflect.ValueOf(secondHookFunction).Pointer()).Name()
	assert.Equal(t, funcName1, funcName2)

	type args struct {
		ctx context.Context
		req *orchestrator.UpdateTargetOfEvaluationRequest
	}
	tests := []struct {
		name    string
		args    args
		wantRes *orchestrator.TargetOfEvaluation
		wantErr assert.WantErr
	}{
		{
			name: "Store first assessment result to the map",
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.UpdateTargetOfEvaluationRequest{
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CertificationTargetId: testdata.MockCertificationTargetID1,
						CatalogId:             testdata.MockCatalogID,
						AssuranceLevel:        &testdata.AssuranceLevelSubstantial,
					},
				},
			},
			wantErr: assert.Nil[error],
			wantRes: &orchestrator.TargetOfEvaluation{
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

			// Create new ToE
			err = s.storage.Create(orchestratortest.NewTargetOfEvaluation(testdata.AssuranceLevelBasic))
			assert.NoError(t, err)

			gotRes, err := s.UpdateTargetOfEvaluation(tt.args.ctx, tt.args.req)

			// wait for all hooks (2 hooks)
			wg.Wait()

			assert.NoError(t, api.Validate(gotRes))

			tt.wantErr(t, err)
			assert.Equal(t, tt.wantRes, gotRes)
			assert.Equal(t, 2, hookCallCounter)
		})
	}
}
