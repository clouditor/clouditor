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
	"sort"
	"sync"
	"testing"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/internal/testutil/orchestratortest"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/persistence/gorm"
	"clouditor.io/clouditor/service"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func TestService_CreateTargetOfEvaluation(t *testing.T) {

	catalogWithoutAssuranceLevelList := orchestratortest.NewCatalog()
	catalogWithoutAssuranceLevelList.AssuranceLevels = []string{}

	type fields struct {
		cloudServiceHooks     []orchestrator.CloudServiceHookFunc
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
		req *orchestrator.CreateTargetOfEvaluationRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.ValueAssertionFunc
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
				return assert.ErrorContains(t, err, "invalid CreateTargetOfEvaluationRequest.TargetOfEvaluation: value is required")
			},
			want: assert.Empty,
		},
		{
			name: "Error getting catalog",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {}),
				authz:   &service.AuthorizationStrategyAllowAll{},
			},
			args: args{req: &orchestrator.CreateTargetOfEvaluationRequest{
				TargetOfEvaluation: orchestratortest.NewTargetOfEvaluation(testdata.AssuranceLevelBasic),
			}},
			want: assert.Empty,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.NotFound, status.Code(err))
				return assert.ErrorContains(t, err, "catalog not found")
			},
		},
		{
			name: "valid and assurance level not set",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Create(&orchestrator.CloudService{Id: testdata.MockCloudServiceID})
					assert.NoError(t, err)

					err = s.Create(catalogWithoutAssuranceLevelList)
					assert.NoError(t, err)
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{req: &orchestrator.CreateTargetOfEvaluationRequest{
				TargetOfEvaluation: orchestratortest.NewTargetOfEvaluation(""),
			}},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				svc := i2[0].(*Service)

				// We want to assert that certain things happened in our database
				var toes []*orchestrator.TargetOfEvaluation
				// for join tables, do not use preload (which is the default)
				err := svc.storage.List(&toes, "", false, 0, -1, gorm.WithoutPreload())
				if !assert.NoError(t, err) {
					return false
				}
				if !assert.Equal(t, 1, len(toes)) {
					return false
				}

				var service orchestrator.CloudService
				err = svc.storage.Get(&service, "id = ?", testdata.MockCloudServiceID)
				if !assert.NoError(t, err) {
					return false
				}

				return assert.Equal(t, 1, len(service.CatalogsInScope))
			},
			wantErr: assert.NoError,
		},
		{
			name: "valid and assurance level set",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Create(&orchestrator.CloudService{Id: testdata.MockCloudServiceID})
					assert.NoError(t, err)

					err = s.Create(orchestratortest.NewCatalog())
					assert.NoError(t, err)
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{req: &orchestrator.CreateTargetOfEvaluationRequest{
				TargetOfEvaluation: orchestratortest.NewTargetOfEvaluation(testdata.AssuranceLevelBasic),
			}},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				svc := i2[0].(*Service)

				// We want to assert that certain things happened in our database
				var toes []*orchestrator.TargetOfEvaluation
				// for join tables, do not use preload (which is the default)
				err := svc.storage.List(&toes, "", false, 0, -1, gorm.WithoutPreload())
				if !assert.NoError(t, err) {
					return false
				}
				if !assert.Equal(t, 1, len(toes)) {
					return false
				}

				var service orchestrator.CloudService
				err = svc.storage.Get(&service, "id = ?", testdata.MockCloudServiceID)
				if !assert.NoError(t, err) {
					return false
				}

				return assert.Equal(t, 1, len(service.CatalogsInScope))
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				cloudServiceHooks:     tt.fields.cloudServiceHooks,
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				metricsFile:           tt.fields.metricsFile,
				loadMetricsFunc:       tt.fields.loadMetricsFunc,
				catalogsFile:          tt.fields.catalogsFile,
				loadCatalogsFunc:      tt.fields.loadCatalogsFunc,
				events:                tt.fields.events,
				authz:                 tt.fields.authz,
			}

			gotRes, err := svc.CreateTargetOfEvaluation(tt.args.ctx, tt.args.req)
			tt.wantErr(t, err)
			tt.want(t, gotRes, svc)

			if tt.want != nil {
				assert.NoError(t, gotRes.Validate())

			}
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
		wantResponse assert.ValueAssertionFunc
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "empty request",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   &service.AuthorizationStrategyAllowAll{},
			},
			args:         args{req: nil},
			wantResponse: assert.Nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "empty request")
			},
		},
		{
			name: "invalid request",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   &service.AuthorizationStrategyAllowAll{},
			},
			args: args{req: &orchestrator.GetTargetOfEvaluationRequest{
				CloudServiceId: "",
			}},
			wantResponse: assert.Nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "CloudServiceId: value must be a valid UUID")
			},
		},
		{
			name: "toe not found",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Create(&orchestrator.CloudService{Id: testdata.MockCloudServiceID})
					assert.NoError(t, err)

					err = s.Create(orchestratortest.NewCatalog())
					assert.NoError(t, err)

					err = s.Create(orchestratortest.NewTargetOfEvaluation(testdata.AssuranceLevelBasic))
					assert.NoError(t, err)
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{req: &orchestrator.GetTargetOfEvaluationRequest{
				CloudServiceId: testutil.TestCloudService2,
				CatalogId:      testdata.MockCatalogID,
			}},
			wantResponse: assert.Nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.NotFound, status.Code(err))
				return assert.ErrorContains(t, err, "ToE not found")
			},
		},
		{
			name: "valid toe",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					err := s.Create(&orchestrator.CloudService{Id: testdata.MockCloudServiceID})
					assert.NoError(t, err)

					err = s.Create(orchestratortest.NewCatalog())
					assert.NoError(t, err)

					err = s.Create(orchestratortest.NewTargetOfEvaluation(testdata.AssuranceLevelBasic))
					assert.NoError(t, err)
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{req: &orchestrator.GetTargetOfEvaluationRequest{
				CloudServiceId: testdata.MockCloudServiceID,
				CatalogId:      "Cat1234",
			}},
			wantResponse: func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				res, ok := i.(*orchestrator.TargetOfEvaluation)
				want := orchestratortest.NewTargetOfEvaluation(testdata.AssuranceLevelBasic)

				return assert.True(t, ok) &&
					assert.NoError(t, res.Validate()) &&
					assert.Equal(t, want.CloudServiceId, res.CloudServiceId) &&
					assert.Equal(t, want.CatalogId, res.CatalogId)
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

			// Validate the error via the ErrorAssertionFunc function
			tt.wantErr(t, err)

			// Validate the response via the ValueAssertionFunc function
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
	err = orchestratorService.storage.Create(&orchestrator.CloudService{Id: testdata.MockCloudServiceID})
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
	assert.NoError(t, listTargetsOfEvaluationResponse.TargetOfEvaluation[0].Validate())
	assert.Equal(t, 1, len(listTargetsOfEvaluationResponse.TargetOfEvaluation))
}

func TestService_UpdateTargetOfEvaluation(t *testing.T) {
	var (
		toe *orchestrator.TargetOfEvaluation
		err error
	)
	orchestratorService := NewService()
	err = orchestratorService.storage.Create(&orchestrator.CloudService{Id: testdata.MockCloudServiceID})
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
	assert.ErrorContains(t, err, "CloudServiceId: value must be a valid UUID")

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
			CloudServiceId: testdata.MockCloudServiceID,
			CatalogId:      testdata.MockCatalogID,
			AssuranceLevel: &testdata.AssuranceLevelBasic,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, toe)
	assert.NoError(t, toe.Validate())
	assert.Equal(t, &testdata.AssuranceLevelBasic, toe.AssuranceLevel)
}

func TestService_RemoveTargetOfEvaluation(t *testing.T) {
	var (
		err                             error
		listTargetsOfEvaluationResponse *orchestrator.ListTargetsOfEvaluationResponse
	)
	orchestratorService := NewService()
	err = orchestratorService.storage.Create(&orchestrator.CloudService{Id: testdata.MockCloudServiceID})
	assert.NoError(t, err)
	err = orchestratorService.storage.Create(orchestratortest.NewCatalog())
	assert.NoError(t, err)

	// 1st case: Empty ID error
	_, err = orchestratorService.RemoveTargetOfEvaluation(context.Background(), &orchestrator.RemoveTargetOfEvaluationRequest{
		CloudServiceId: "",
		CatalogId:      "",
	})
	assert.Error(t, err)
	assert.Equal(t, status.Code(err), codes.InvalidArgument)

	// 2nd case: ErrRecordNotFound
	_, err = orchestratorService.RemoveTargetOfEvaluation(context.Background(), &orchestrator.RemoveTargetOfEvaluationRequest{
		CloudServiceId: testdata.MockCloudServiceID,
		CatalogId:      "0000",
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
	assert.NoError(t, listTargetsOfEvaluationResponse.TargetOfEvaluation[0].Validate())
	assert.Equal(t, 1, len(listTargetsOfEvaluationResponse.TargetOfEvaluation))

	// Remove record
	_, err = orchestratorService.RemoveTargetOfEvaluation(context.Background(), &orchestrator.RemoveTargetOfEvaluationRequest{
		CloudServiceId: testdata.MockCloudServiceID,
		CatalogId:      testdata.MockCatalogID,
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

	service := NewService()
	service.RegisterToeHook(firstHookFunction)
	service.RegisterToeHook(secondHookFunction)

	// Check if first hook is registered
	funcName1 := runtime.FuncForPC(reflect.ValueOf(service.toeHooks[0]).Pointer()).Name()
	funcName2 := runtime.FuncForPC(reflect.ValueOf(firstHookFunction).Pointer()).Name()
	assert.Equal(t, funcName1, funcName2)

	// Check if second hook is registered
	funcName1 = runtime.FuncForPC(reflect.ValueOf(service.toeHooks[1]).Pointer()).Name()
	funcName2 = runtime.FuncForPC(reflect.ValueOf(secondHookFunction).Pointer()).Name()
	assert.Equal(t, funcName1, funcName2)

	type args struct {
		ctx context.Context
		req *orchestrator.UpdateTargetOfEvaluationRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *orchestrator.TargetOfEvaluation
		wantErr  bool
	}{
		{
			name: "Store first assessment result to the map",
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.UpdateTargetOfEvaluationRequest{
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
						CloudServiceId: testdata.MockCloudServiceID,
						CatalogId:      testdata.MockCatalogID,
						AssuranceLevel: &testdata.AssuranceLevelSubstantial,
					},
				},
			},
			wantErr: false,
			wantResp: &orchestrator.TargetOfEvaluation{
				CloudServiceId: testdata.MockCloudServiceID,
				CatalogId:      testdata.MockCatalogID,
				AssuranceLevel: &testdata.AssuranceLevelSubstantial,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hookCallCounter = 0

			// Create service
			s := service
			err := s.storage.Create(&orchestrator.CloudService{Id: testdata.MockCloudServiceID})
			assert.NoError(t, err)

			// Create catalog
			err = s.storage.Create(orchestratortest.NewCatalog())
			assert.NoError(t, err)

			// Create new ToE
			err = s.storage.Create(orchestratortest.NewTargetOfEvaluation(testdata.AssuranceLevelBasic))
			assert.NoError(t, err)

			gotResp, err := s.UpdateTargetOfEvaluation(tt.args.ctx, tt.args.req)

			// wait for all hooks (2 hooks)
			wg.Wait()

			assert.NoError(t, gotResp.Validate())

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateTargetOfEvaluation() error = %v, wantErrMessage %v", err, tt.wantErr)
				return
			}
			if !proto.Equal(gotResp, tt.wantResp) {
				t.Errorf("UpdateTargetOfEvaluation() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
			assert.Equal(t, 2, hookCallCounter)
		})
	}
}
func TestService_ListControlsInScope(t *testing.T) {
	type fields struct {
		storage persistence.Storage
		authz   service.AuthorizationStrategy
	}
	type args struct {
		ctx context.Context
		req *orchestrator.ListControlsInScopeRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *orchestrator.ListControlsInScopeResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Empty request",
			args: args{
				ctx: context.Background(),
				req: nil,
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "Missing request property",
			args: args{
				ctx: context.Background(),
				req: &orchestrator.ListControlsInScopeRequest{
					CloudServiceId: testdata.MockCloudServiceID,
				},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "invalid ListControlsInScopeRequest.CatalogId: value length must be at least 1 runes")
			},
		},
		{
			name: "no controls explicitly selected - all controls status unspecified",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
					assert.NoError(t, s.Create(&orchestrator.CloudService{Id: testdata.MockCloudServiceID}))
					assert.NoError(t, s.Create(orchestratortest.NewTargetOfEvaluation(testdata.AssuranceLevelBasic)))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.ListControlsInScopeRequest{
					CloudServiceId: testdata.MockCloudServiceID,
					CatalogId:      testdata.MockCatalogID,
				},
			},
			wantRes: &orchestrator.ListControlsInScopeResponse{
				ControlsInScope: []*orchestrator.ControlInScope{
					{
						ControlId:                        testdata.MockControlID1,
						ControlCategoryName:              testdata.MockCategoryName,
						ControlCategoryCatalogId:         testdata.MockCatalogID,
						TargetOfEvaluationCloudServiceId: testdata.MockCloudServiceID,
						TargetOfEvaluationCatalogId:      testdata.MockCatalogID,
						MonitoringStatus:                 orchestrator.MonitoringStatus_MONITORING_STATUS_UNSPECIFIED,
					},
					{
						ControlId:                        testdata.MockSubControlID11,
						ControlCategoryName:              testdata.MockCategoryName,
						ControlCategoryCatalogId:         testdata.MockCatalogID,
						TargetOfEvaluationCloudServiceId: testdata.MockCloudServiceID,
						TargetOfEvaluationCatalogId:      testdata.MockCatalogID,
						MonitoringStatus:                 orchestrator.MonitoringStatus_MONITORING_STATUS_UNSPECIFIED,
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "one control explicitly set to continuously monitored",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
					assert.NoError(t, s.Create(&orchestrator.CloudService{Id: testdata.MockCloudServiceID}))
					assert.NoError(t, s.Create(orchestratortest.NewTargetOfEvaluation(testdata.AssuranceLevelBasic)))
					assert.NoError(t, s.Update(&orchestrator.ControlInScope{
						ControlId:                        testdata.MockControlID1,
						ControlCategoryName:              testdata.MockCategoryName,
						ControlCategoryCatalogId:         testdata.MockCatalogID,
						TargetOfEvaluationCloudServiceId: testdata.MockCloudServiceID,
						TargetOfEvaluationCatalogId:      testdata.MockCatalogID,
						MonitoringStatus:                 orchestrator.MonitoringStatus_MONITORING_STATUS_AUTOMATICALLY_MONITORED,
					}, orchestrator.ControlInScope{
						ControlId:                        testdata.MockControlID1,
						ControlCategoryName:              testdata.MockCategoryName,
						ControlCategoryCatalogId:         testdata.MockCatalogID,
						TargetOfEvaluationCloudServiceId: testdata.MockCloudServiceID,
						TargetOfEvaluationCatalogId:      testdata.MockCatalogID,
					}))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				ctx: context.TODO(),
				req: &orchestrator.ListControlsInScopeRequest{
					CloudServiceId: testdata.MockCloudServiceID,
					CatalogId:      testdata.MockCatalogID,
				},
			},
			wantRes: &orchestrator.ListControlsInScopeResponse{
				ControlsInScope: []*orchestrator.ControlInScope{
					{
						ControlId:                        testdata.MockControlID1,
						ControlCategoryName:              testdata.MockCategoryName,
						ControlCategoryCatalogId:         testdata.MockCatalogID,
						TargetOfEvaluationCloudServiceId: testdata.MockCloudServiceID,
						TargetOfEvaluationCatalogId:      testdata.MockCatalogID,
						MonitoringStatus:                 orchestrator.MonitoringStatus_MONITORING_STATUS_AUTOMATICALLY_MONITORED,
					},
					{
						ControlId:                        testdata.MockSubControlID11,
						ControlCategoryName:              testdata.MockCategoryName,
						ControlCategoryCatalogId:         testdata.MockCatalogID,
						TargetOfEvaluationCloudServiceId: testdata.MockCloudServiceID,
						TargetOfEvaluationCatalogId:      testdata.MockCatalogID,
						MonitoringStatus:                 orchestrator.MonitoringStatus_MONITORING_STATUS_UNSPECIFIED,
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "permission denied",
			fields: fields{
				authz: &service.AuthorizationStrategyJWT{CloudServicesKey: testutil.TestCustomClaims},
			},
			args: args{
				ctx: testutil.TestContextOnlyService1,
				req: &orchestrator.ListControlsInScopeRequest{
					CloudServiceId: testutil.TestCloudService2,
					CatalogId:      testdata.MockCatalogID,
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, service.ErrPermissionDenied)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				storage: tt.fields.storage,
				authz:   tt.fields.authz,
			}

			gotRes, err := svc.ListControlsInScope(tt.args.ctx, tt.args.req)
			tt.wantErr(t, err, tt.args)

			if !proto.Equal(gotRes, tt.wantRes) {
				t.Errorf("Service.ListControlInScope() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_AddControlToScope(t *testing.T) {
	type fields struct {
		storage persistence.Storage
		authz   service.AuthorizationStrategy
	}
	type args struct {
		in0 context.Context
		req *orchestrator.AddControlToScopeRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *orchestrator.ControlInScope
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Invalid request",
			args: args{
				req: nil,
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "already exists",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
					assert.NoError(t, s.Create(&orchestrator.CloudService{Id: testdata.MockCloudServiceID}))
					assert.NoError(t, s.Create(orchestratortest.NewTargetOfEvaluation(testdata.AssuranceLevelBasic)))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				in0: context.TODO(),
				req: &orchestrator.AddControlToScopeRequest{
					Scope: &orchestrator.ControlInScope{
						ControlId:                        testdata.MockControlID1,
						ControlCategoryName:              testdata.MockCategoryName,
						ControlCategoryCatalogId:         testdata.MockCatalogID,
						TargetOfEvaluationCloudServiceId: testdata.MockCloudServiceID,
						TargetOfEvaluationCatalogId:      testdata.MockCatalogID,
						MonitoringStatus:                 orchestrator.MonitoringStatus_MONITORING_STATUS_AUTOMATICALLY_MONITORED,
					},
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, codes.AlreadyExists, status.Code(err))
			},
		},
		{
			name: "ToE not found",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				in0: context.TODO(),
				req: &orchestrator.AddControlToScopeRequest{
					Scope: &orchestrator.ControlInScope{
						ControlId:                        testdata.MockControlID1,
						ControlCategoryName:              testdata.MockCategoryName,
						ControlCategoryCatalogId:         testdata.MockCatalogID,
						TargetOfEvaluationCloudServiceId: testdata.MockCloudServiceID,
						TargetOfEvaluationCatalogId:      testdata.MockCatalogID,
						MonitoringStatus:                 orchestrator.MonitoringStatus_MONITORING_STATUS_AUTOMATICALLY_MONITORED,
					},
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "valid update",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
					assert.NoError(t, s.Create(&orchestrator.CloudService{Id: testdata.MockCloudServiceID}))
					assert.NoError(t, s.Create(orchestratortest.NewTargetOfEvaluation(testdata.AssuranceLevelBasic)))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				in0: context.TODO(),
				req: &orchestrator.AddControlToScopeRequest{
					Scope: &orchestrator.ControlInScope{
						ControlId:                        testdata.MockControlID2,
						ControlCategoryName:              testdata.MockCategoryName,
						ControlCategoryCatalogId:         testdata.MockCatalogID,
						TargetOfEvaluationCloudServiceId: testdata.MockCloudServiceID,
						TargetOfEvaluationCatalogId:      testdata.MockCatalogID,
						MonitoringStatus:                 orchestrator.MonitoringStatus_MONITORING_STATUS_AUTOMATICALLY_MONITORED,
					},
				},
			},
			wantRes: &orchestrator.ControlInScope{
				ControlId:                        testdata.MockControlID2,
				ControlCategoryName:              testdata.MockCategoryName,
				ControlCategoryCatalogId:         testdata.MockCatalogID,
				TargetOfEvaluationCloudServiceId: testdata.MockCloudServiceID,
				TargetOfEvaluationCatalogId:      testdata.MockCatalogID,
				MonitoringStatus:                 orchestrator.MonitoringStatus_MONITORING_STATUS_AUTOMATICALLY_MONITORED,
			},
			wantErr: assert.NoError,
		},
		{
			name: "permission denied",
			fields: fields{
				authz: &service.AuthorizationStrategyJWT{CloudServicesKey: testutil.TestCustomClaims},
			},
			args: args{
				in0: testutil.TestContextOnlyService1,
				req: &orchestrator.AddControlToScopeRequest{
					Scope: &orchestrator.ControlInScope{
						TargetOfEvaluationCloudServiceId: testutil.TestCloudService2,
						ControlId:                        testdata.MockControlID1,
						ControlCategoryName:              testdata.MockCategoryName,
						ControlCategoryCatalogId:         testdata.MockCatalogID,
						TargetOfEvaluationCatalogId:      testdata.MockCatalogID,
						MonitoringStatus:                 orchestrator.MonitoringStatus_MONITORING_STATUS_AUTOMATICALLY_MONITORED,
					},
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, service.ErrPermissionDenied)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				storage: tt.fields.storage,
				authz:   tt.fields.authz,
			}

			gotRes, err := svc.AddControlToScope(tt.args.in0, tt.args.req)
			tt.wantErr(t, err)

			if !proto.Equal(gotRes, tt.wantRes) {
				t.Errorf("Service.AddControlToScope() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_UpdateControlInScope(t *testing.T) {
	type fields struct {
		storage persistence.Storage
		authz   service.AuthorizationStrategy
	}
	type args struct {
		in0 context.Context
		req *orchestrator.UpdateControlInScopeRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *orchestrator.ControlInScope
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Invalid request",
			args: args{
				req: nil,
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "valid update",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
					assert.NoError(t, s.Create(&orchestrator.CloudService{Id: testdata.MockCloudServiceID}))
					assert.NoError(t, s.Create(orchestratortest.NewTargetOfEvaluation(testdata.AssuranceLevelBasic)))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				in0: context.TODO(),
				req: &orchestrator.UpdateControlInScopeRequest{
					Scope: &orchestrator.ControlInScope{
						ControlId:                        testdata.MockControlID1,
						ControlCategoryName:              testdata.MockCategoryName,
						ControlCategoryCatalogId:         testdata.MockCatalogID,
						TargetOfEvaluationCloudServiceId: testdata.MockCloudServiceID,
						TargetOfEvaluationCatalogId:      testdata.MockCatalogID,
						MonitoringStatus:                 orchestrator.MonitoringStatus_MONITORING_STATUS_AUTOMATICALLY_MONITORED,
					},
				},
			},
			wantRes: &orchestrator.ControlInScope{
				ControlId:                        testdata.MockControlID1,
				ControlCategoryName:              testdata.MockCategoryName,
				ControlCategoryCatalogId:         testdata.MockCatalogID,
				TargetOfEvaluationCloudServiceId: testdata.MockCloudServiceID,
				TargetOfEvaluationCatalogId:      testdata.MockCatalogID,
				MonitoringStatus:                 orchestrator.MonitoringStatus_MONITORING_STATUS_AUTOMATICALLY_MONITORED,
			},
			wantErr: assert.NoError,
		},
		{
			name: "ToE not found",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				in0: context.TODO(),
				req: &orchestrator.UpdateControlInScopeRequest{
					Scope: &orchestrator.ControlInScope{
						ControlId:                        testdata.MockControlID1,
						ControlCategoryName:              testdata.MockCategoryName,
						ControlCategoryCatalogId:         testdata.MockCatalogID,
						TargetOfEvaluationCloudServiceId: testdata.MockCloudServiceID,
						TargetOfEvaluationCatalogId:      testdata.MockCatalogID,
						MonitoringStatus:                 orchestrator.MonitoringStatus_MONITORING_STATUS_AUTOMATICALLY_MONITORED,
					},
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "permission denied",
			fields: fields{
				authz: &service.AuthorizationStrategyJWT{CloudServicesKey: testutil.TestCustomClaims},
			},
			args: args{
				in0: testutil.TestContextOnlyService1,
				req: &orchestrator.UpdateControlInScopeRequest{
					Scope: &orchestrator.ControlInScope{
						TargetOfEvaluationCloudServiceId: testutil.TestCloudService2,
						ControlId:                        testdata.MockControlID1,
						ControlCategoryName:              testdata.MockCategoryName,
						ControlCategoryCatalogId:         testdata.MockCatalogID,
						TargetOfEvaluationCatalogId:      testdata.MockCatalogID,
						MonitoringStatus:                 orchestrator.MonitoringStatus_MONITORING_STATUS_AUTOMATICALLY_MONITORED,
					},
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, service.ErrPermissionDenied)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				storage: tt.fields.storage,
				authz:   tt.fields.authz,
			}
			gotRes, err := svc.UpdateControlInScope(tt.args.in0, tt.args.req)
			assert.NoError(t, gotRes.Validate())
			tt.wantErr(t, err, tt.args)

			if !proto.Equal(gotRes, tt.wantRes) {
				t.Errorf("Service.UpdateControlInScope() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_RemoveControlFromScope(t *testing.T) {
	type fields struct {
		storage persistence.Storage
		authz   service.AuthorizationStrategy
	}
	type args struct {
		in0 context.Context
		req *orchestrator.RemoveControlFromScopeRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Invalid request",
			args: args{
				req: nil,
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, codes.InvalidArgument, status.Code(err))
			},
		},
		{
			name: "valid remove",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
					assert.NoError(t, s.Create(&orchestrator.CloudService{Id: testdata.MockCloudServiceID}))
					assert.NoError(t, s.Create(orchestratortest.NewTargetOfEvaluation(testdata.AssuranceLevelBasic)))
				}),
				authz: &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				in0: context.TODO(),
				req: &orchestrator.RemoveControlFromScopeRequest{
					ControlId:           testdata.MockControlID1,
					ControlCategoryName: testdata.MockCategoryName,
					CloudServiceId:      testdata.MockCloudServiceID,
					CatalogId:           testdata.MockCatalogID,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "ToE not found",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t),
				authz:   &service.AuthorizationStrategyAllowAll{},
			},
			args: args{
				in0: context.TODO(),
				req: &orchestrator.RemoveControlFromScopeRequest{
					ControlId:           testdata.MockControlID1,
					ControlCategoryName: testdata.MockCategoryName,
					CloudServiceId:      testdata.MockCloudServiceID,
					CatalogId:           testdata.MockCatalogID,
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, codes.NotFound, status.Code(err))
			},
		},
		{
			name: "permission denied",
			fields: fields{
				authz: &service.AuthorizationStrategyJWT{CloudServicesKey: testutil.TestCustomClaims},
			},
			args: args{
				in0: testutil.TestContextOnlyService1,
				req: &orchestrator.RemoveControlFromScopeRequest{
					CloudServiceId:      testutil.TestCloudService2,
					ControlId:           testdata.MockControlID1,
					ControlCategoryName: testdata.MockCategoryName,
					CatalogId:           testdata.MockCatalogID,
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, service.ErrPermissionDenied)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				storage: tt.fields.storage,
				authz:   tt.fields.authz,
			}
			_, err := svc.RemoveControlFromScope(tt.args.in0, tt.args.req)
			tt.wantErr(t, err, tt.args)
		})
	}
}

func Test_getControls(t *testing.T) {
	type args struct {
		controls        []*orchestrator.Control
		level           string
		assuranceLevels []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*orchestrator.Control
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Missing controls",
			args: args{
				controls:        []*orchestrator.Control{},
				level:           testdata.AssuranceLevelBasic,
				assuranceLevels: []string{testdata.AssuranceLevelBasic, testdata.AssuranceLevelSubstantial, testdata.AssuranceLevelHigh},
			},
			want:    []*orchestrator.Control{},
			wantErr: assert.NoError,
		},
		{
			name: "List of assurance levels unspecified",
			args: args{
				controls:        orchestratortest.MockControlsInScope,
				assuranceLevels: []string{},
			},
			want: []*orchestrator.Control{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "assurance levels are empty")
			},
		},
		{
			name: "Assurance level unspecified",
			args: args{
				controls:        orchestratortest.MockControlsInScope,
				level:           "",
				assuranceLevels: []string{testdata.AssuranceLevelBasic, testdata.AssuranceLevelSubstantial, testdata.AssuranceLevelHigh},
			},
			want: []*orchestrator.Control{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "assurance level is empty")
			},
		},
		{
			name: "Happy path with assurance level basic (1 assurance level)",
			args: args{
				controls:        orchestratortest.MockControlsInScope,
				level:           testdata.AssuranceLevelHigh,
				assuranceLevels: []string{testdata.AssuranceLevelHigh},
			},
			want:    []*orchestrator.Control{orchestratortest.MockControlsInScope3, orchestratortest.MockControlsInScopeSubControl32},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path with assurance level basic (3 assurance level)",
			args: args{
				controls:        orchestratortest.MockControlsInScope,
				level:           testdata.AssuranceLevelBasic,
				assuranceLevels: []string{testdata.AssuranceLevelBasic, testdata.AssuranceLevelSubstantial, testdata.AssuranceLevelHigh},
			},
			want:    []*orchestrator.Control{orchestratortest.MockControlsInScope1, orchestratortest.MockControlsInScopeSubControl11, orchestratortest.MockControlsInScope2, orchestratortest.MockControlsInScopeSubControl21},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path with assurance level substantial (3 assurance level)",
			args: args{
				controls:        orchestratortest.MockControlsInScope,
				level:           testdata.AssuranceLevelSubstantial,
				assuranceLevels: []string{testdata.AssuranceLevelBasic, testdata.AssuranceLevelSubstantial, testdata.AssuranceLevelHigh},
			},
			want:    []*orchestrator.Control{orchestratortest.MockControlsInScope1, orchestratortest.MockControlsInScopeSubControl11, orchestratortest.MockControlsInScope2, orchestratortest.MockControlsInScopeSubControl21, orchestratortest.MockControlsInScope3, orchestratortest.MockControlsInScopeSubControl31},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path with assurance level high (3 assurance level)",
			args: args{
				controls:        orchestratortest.MockControlsInScope,
				level:           testdata.AssuranceLevelHigh,
				assuranceLevels: []string{testdata.AssuranceLevelBasic, testdata.AssuranceLevelSubstantial, testdata.AssuranceLevelHigh},
			},
			want:    []*orchestrator.Control{orchestratortest.MockControlsInScope1, orchestratortest.MockControlsInScopeSubControl11, orchestratortest.MockControlsInScope2, orchestratortest.MockControlsInScopeSubControl21, orchestratortest.MockControlsInScope3, orchestratortest.MockControlsInScopeSubControl31, orchestratortest.MockControlsInScopeSubControl32},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getControls(tt.args.controls, tt.args.assuranceLevels, tt.args.level)

			// Sort slices
			sort.Slice(tt.want, func(i, j int) bool {
				return tt.want[i].GetId() > tt.want[j].GetId()
			})
			sort.Slice(got, func(i, j int) bool {
				return got[i].GetId() > got[j].GetId()
			})

			assert.Equal(t, tt.want, got)
			tt.wantErr(t, err)
		})
	}
}
