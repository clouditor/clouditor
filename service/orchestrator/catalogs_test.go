// Copyright 2023 Fraunhofer AISEC
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
	"fmt"
	"os"
	"testing"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/internal/testutil/servicetest/orchestratortest"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/persistence/gorm"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func TestService_CreateCatalog(t *testing.T) {
	// Mock catalogs
	mockCatalog := orchestratortest.NewCatalog()
	mockCatalogWithoutID := orchestratortest.NewCatalog()
	mockCatalogWithoutID.Id = ""

	type args struct {
		in0 context.Context
		req *orchestrator.CreateCatalogRequest
	}
	tests := []struct {
		name         string
		args         args
		wantResponse *orchestrator.Catalog
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "missing request",
			args: args{
				context.Background(),
				nil,
			},
			wantResponse: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.ErrorContains(t, err, "empty request")
				return assert.Equal(t, status.Code(err), codes.InvalidArgument)
			},
		},
		{
			name: "missing catalog",
			args: args{
				context.Background(),
				&orchestrator.CreateCatalogRequest{},
			},
			wantResponse: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.ErrorContains(t, err, "Catalog: value is required")
				return assert.Equal(t, status.Code(err), codes.InvalidArgument)
			},
		},
		{
			name: "missing catalog id",
			args: args{
				context.Background(),
				&orchestrator.CreateCatalogRequest{
					Catalog: mockCatalogWithoutID,
				},
			},
			wantResponse: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.ErrorContains(t, err, "Catalog.Id: value length must be at least 1 runes")
				return assert.Equal(t, status.Code(err), codes.InvalidArgument)
			},
		},
		{
			name: "valid catalog",
			args: args{
				context.Background(),
				&orchestrator.CreateCatalogRequest{
					Catalog: mockCatalog,
				},
			},
			wantResponse: mockCatalog,
			wantErr:      assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			gotResponse, err := s.CreateCatalog(tt.args.in0, tt.args.req)
			tt.wantErr(t, err)

			// If no error is wanted, check response
			if !proto.Equal(gotResponse, tt.wantResponse) {
				t.Errorf("Service.CreateCatalog() = %v, want %v", gotResponse, tt.wantResponse)

				// Check catalog structure with validation method
				assert.NoError(t, gotResponse.Validate())
			}
		})
	}
}

func TestService_GetCatalog(t *testing.T) {
	type fields struct {
		storage persistence.Storage
	}
	type args struct {
		req *orchestrator.GetCatalogRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantResponse assert.ValueAssertionFunc
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "invalid request",
			fields: fields{storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
				// Create Catalog
				assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
			})},
			wantResponse: assert.Nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "empty request")
			},
		},
		{
			name:         "catalog ID empty",
			fields:       fields{storage: testutil.NewInMemoryStorage(t)},
			args:         args{req: &orchestrator.GetCatalogRequest{CatalogId: ""}},
			wantResponse: assert.Nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "CatalogId: value length must be at least 1 runes")
			},
		},
		{
			name: "catalog not found",
			fields: fields{storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
				// Create Catalog
				assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
			})},
			args:         args{req: &orchestrator.GetCatalogRequest{CatalogId: "a"}},
			wantResponse: assert.Nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.NotFound, status.Code(err))
				return assert.ErrorContains(t, err, "catalog not found")
			},
		},
		{
			name: "valid",
			fields: fields{storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
				// Create Catalog
				assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
			})},
			args: args{req: &orchestrator.GetCatalogRequest{CatalogId: testdata.MockCatalogID}},
			wantResponse: func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				res, ok := i.(*orchestrator.Catalog)
				want := orchestratortest.NewCatalog()
				assert.True(t, ok)
				fmt.Println(res)
				assert.Equal(t, 1, len(res.Categories))
				return assert.Equal(t, want.Id, res.Id)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orchestratorService := Service{
				storage: tt.fields.storage,
			}
			res, err := orchestratorService.GetCatalog(context.Background(), tt.args.req)

			// Validate the error via the ErrorAssertionFunc function
			tt.wantErr(t, err)

			// Validate the response via the ValueAssertionFunc function
			tt.wantResponse(t, res)
		})
	}
}

func TestService_ListCatalogs(t *testing.T) {
	var (
		listCatalogsResponse *orchestrator.ListCatalogsResponse
		err                  error
	)

	orchestratorService := NewService(WithCatalogsFolder("catalogs"))
	// 1st case: Default catalogs stored
	listCatalogsResponse, err = orchestratorService.ListCatalogs(context.Background(), &orchestrator.ListCatalogsRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listCatalogsResponse.Catalogs)
	assert.Equal(t, 1, len(listCatalogsResponse.Catalogs))
}

func TestService_UpdateCatalog(t *testing.T) {
	var (
		catalog *orchestrator.Catalog
		err     error
	)
	orchestratorService := NewService()

	// 1st case: Catalog is nil
	_, err = orchestratorService.UpdateCatalog(context.Background(), &orchestrator.UpdateCatalogRequest{})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// 2nd case: Catalog ID is nil
	_, err = orchestratorService.UpdateCatalog(context.Background(), &orchestrator.UpdateCatalogRequest{
		Catalog: catalog,
	})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// 3rd case: Catalog not found since there are no catalogs yet
	_, err = orchestratorService.UpdateCatalog(context.Background(), &orchestrator.UpdateCatalogRequest{
		Catalog: &orchestrator.Catalog{
			Id:              testdata.MockCatalogID,
			Name:            testdata.MockCatalogName,
			AssuranceLevels: []string{testdata.AssuranceLevelBasic, testdata.AssuranceLevelSubstantial, testdata.AssuranceLevelHigh},
		},
	})
	assert.Equal(t, codes.NotFound, status.Code(err))
	assert.ErrorContains(t, err, "catalog not found")

	// 4th case: Catalog updated successfully
	mockCatalog := orchestratortest.NewCatalog()
	err = orchestratorService.storage.Create(mockCatalog)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	// update the Catalog's description and send the update request
	mockCatalog.Description = "new description"
	catalog, err = orchestratorService.UpdateCatalog(context.Background(), &orchestrator.UpdateCatalogRequest{
		Catalog: mockCatalog,
	})
	assert.NoError(t, err)
	assert.NotNil(t, catalog)
	assert.Equal(t, "new description", catalog.Description)
}

func TestService_RemoveCatalog(t *testing.T) {
	var (
		err                  error
		listCatalogsResponse *orchestrator.ListCatalogsResponse
	)
	orchestratorService := NewService(WithCatalogsFolder("internal/testcatalogs/emptyTestFolder"))

	// 1st case: Empty catalog ID error
	_, err = orchestratorService.RemoveCatalog(context.Background(), &orchestrator.RemoveCatalogRequest{CatalogId: ""})
	assert.Error(t, err)
	assert.Equal(t, status.Code(err), codes.InvalidArgument)

	// 2nd case: ErrRecordNotFound
	_, err = orchestratorService.RemoveCatalog(context.Background(), &orchestrator.RemoveCatalogRequest{CatalogId: "0000"})
	assert.Error(t, err)
	assert.Equal(t, status.Code(err), codes.NotFound)

	// 3rd case: Record removed successfully
	mockCatalog := orchestratortest.NewCatalog()
	err = orchestratorService.storage.Create(mockCatalog)
	assert.NoError(t, err)

	// There is one catalog in the db now (one default plus NewCatalog)
	listCatalogsResponse, err = orchestratorService.ListCatalogs(context.Background(), &orchestrator.ListCatalogsRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listCatalogsResponse.Catalogs)
	assert.Equal(t, 1, len(listCatalogsResponse.Catalogs))

	// Remove record
	_, err = orchestratorService.RemoveCatalog(context.Background(), &orchestrator.RemoveCatalogRequest{CatalogId: mockCatalog.Id})
	assert.NoError(t, err)

	// There is no record left in the DB
	listCatalogsResponse, err = orchestratorService.ListCatalogs(context.Background(), &orchestrator.ListCatalogsRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listCatalogsResponse.Catalogs)
	assert.Equal(t, 0, len(listCatalogsResponse.Catalogs))
}

func TestService_GetCategory(t *testing.T) {
	type fields struct {
		cloudServiceHooks     []orchestrator.CloudServiceHookFunc
		AssessmentResultHooks []func(result *assessment.AssessmentResult, err error)
		storage               persistence.Storage
		metricsFile           string
		loadMetricsFunc       func() ([]*assessment.Metric, error)
		catalogsFile          string
		loadCatalogsFunc      func() ([]*orchestrator.Catalog, error)
		events                chan *orchestrator.MetricChangeEvent
	}
	type args struct {
		ctx context.Context
		req *orchestrator.GetCategoryRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *orchestrator.Category
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Category not found",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					// Create Catalog
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
				})},
			args: args{
				req: &orchestrator.GetCategoryRequest{CatalogId: "WrongCatalogID", CategoryName: testdata.MockCategoryName},
			},
			wantRes: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.NotFound, status.Code(err))
				return assert.ErrorContains(t, err, "category not found")
			},
		},
		{
			name: "valid",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					// Create Catalog
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
				})},
			args: args{
				req: &orchestrator.GetCategoryRequest{CatalogId: testdata.MockCatalogID, CategoryName: testdata.MockCategoryName},
			},
			wantRes: &orchestrator.Category{
				Name:        testdata.MockCategoryName,
				Description: testdata.MockCategoryDescription,
				CatalogId:   testdata.MockCatalogID,
				Controls: []*orchestrator.Control{
					{
						Id:                testdata.MockControlID1,
						Name:              testdata.MockControlName,
						CategoryName:      testdata.MockCategoryName,
						CategoryCatalogId: testdata.MockCatalogID,
						Description:       testdata.MockControlDescription,
						Controls:          []*orchestrator.Control{},
					},
					{
						Id:                testdata.MockControlID2,
						Name:              testdata.MockControlName,
						CategoryName:      testdata.MockCategoryName,
						CategoryCatalogId: testdata.MockCatalogID,
						Description:       testdata.MockControlDescription,
						Controls:          []*orchestrator.Control{},
					},
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Service{
				cloudServiceHooks:     tt.fields.cloudServiceHooks,
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				metricsFile:           tt.fields.metricsFile,
				loadMetricsFunc:       tt.fields.loadMetricsFunc,
				catalogsFolder:        tt.fields.catalogsFile,
				loadCatalogsFunc:      tt.fields.loadCatalogsFunc,
				events:                tt.fields.events,
			}
			gotRes, err := srv.GetCategory(tt.args.ctx, tt.args.req)

			tt.wantErr(t, err)
			if !proto.Equal(gotRes, tt.wantRes) {
				t.Errorf("Service.GetCategory() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_GetControl(t *testing.T) {
	type fields struct {
		cloudServiceHooks     []orchestrator.CloudServiceHookFunc
		AssessmentResultHooks []func(result *assessment.AssessmentResult, err error)
		storage               persistence.Storage
		metricsFile           string
		loadMetricsFunc       func() ([]*assessment.Metric, error)
		catalogsFolder        string
		loadCatalogsFunc      func() ([]*orchestrator.Catalog, error)
	}
	type args struct {
		ctx context.Context
		req *orchestrator.GetControlRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *orchestrator.Control
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Control not found",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					// Create Catalog
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
				})},
			args:    args{req: &orchestrator.GetControlRequest{CatalogId: testdata.MockCatalogID, CategoryName: testdata.MockCategoryName, ControlId: "WrongControlID"}},
			wantRes: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.NotFound, status.Code(err))
				return assert.ErrorContains(t, err, "control not found")
			},
		},
		{
			name: "valid",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					// Create Catalog
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
				})},
			args: args{req: &orchestrator.GetControlRequest{CatalogId: testdata.MockCatalogID, CategoryName: testdata.MockCategoryName, ControlId: testdata.MockControlID1}},
			wantRes: &orchestrator.Control{
				Id:                testdata.MockControlID1,
				CategoryName:      testdata.MockCategoryName,
				CategoryCatalogId: testdata.MockCatalogID,
				Name:              testdata.MockControlName,
				Description:       testdata.MockControlDescription,
				Controls: []*orchestrator.Control{{
					Id:                             testdata.MockSubControlID11,
					Name:                           testdata.MockSubControlName,
					Description:                    testdata.MockSubControlDescription,
					Metrics:                        []*assessment.Metric{}, // metrics on sub-controls are not returned
					CategoryName:                   testdata.MockCategoryName,
					CategoryCatalogId:              testdata.MockCatalogID,
					AssuranceLevel:                 &testdata.AssuranceLevelBasic,
					ParentControlId:                util.Ref(testdata.MockControlID1),
					ParentControlCategoryCatalogId: util.Ref(testdata.MockCatalogID),
					ParentControlCategoryName:      util.Ref(testdata.MockCategoryName),
				}},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Service{
				cloudServiceHooks:     tt.fields.cloudServiceHooks,
				AssessmentResultHooks: tt.fields.AssessmentResultHooks,
				storage:               tt.fields.storage,
				metricsFile:           tt.fields.metricsFile,
				loadMetricsFunc:       tt.fields.loadMetricsFunc,
				catalogsFolder:        tt.fields.catalogsFolder,
				loadCatalogsFunc:      tt.fields.loadCatalogsFunc,
			}
			gotRes, err := srv.GetControl(tt.args.ctx, tt.args.req)

			tt.wantErr(t, err)
			if !proto.Equal(gotRes, tt.wantRes) {
				t.Errorf("Service.GetControl() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_ListControls(t *testing.T) {
	var (
		listControlsResponse *orchestrator.ListControlsResponse
		err                  error
	)

	orchestratorService := NewService(WithCatalogsFolder("internal/testcatalogs/emptyTestFolder"))
	// 1st case: No Controls stored
	listControlsResponse, err = orchestratorService.ListControls(context.Background(), &orchestrator.ListControlsRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listControlsResponse.Controls)
	assert.Empty(t, listControlsResponse.Controls)

	// 2nd case: 30 controls stored; note that we do not have to create an extra catalog/control since NewService above already loads the default catalogs/controls
	orchestratorService = NewService(WithCatalogsFolder("catalogs"))
	listControlsResponse, err = orchestratorService.ListControls(context.Background(), &orchestrator.ListControlsRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listControlsResponse.Controls)
	assert.NotEmpty(t, listControlsResponse.Controls)
	// there are 30 default controls
	assert.Equal(t, 30, len(listControlsResponse.Controls))

	// 3th case: List controls for a specific catalog and category.
	listControlsResponse, err = orchestratorService.ListControls(context.Background(), &orchestrator.ListControlsRequest{
		CatalogId:    "DemoCatalog",
		CategoryName: "Communication Security",
	})
	assert.NoError(t, err)
	assert.NotNil(t, listControlsResponse.Controls)
	assert.NotEmpty(t, listControlsResponse.Controls)
	assert.Equal(t, 4, len(listControlsResponse.Controls))
}

func TestService_loadCatalogs(t *testing.T) {
	type fields struct {
		AssessmentResultHooks []func(result *assessment.AssessmentResult, err error)
		storage               persistence.Storage
		metricsFile           string
		loadMetricsFunc       func() ([]*assessment.Metric, error)
		catalogsFolder        string
		loadCatalogsFunc      func() ([]*orchestrator.Catalog, error)
		events                chan *orchestrator.MetricChangeEvent
	}
	tests := []struct {
		name       string
		fields     fields
		wantResult assert.ValueAssertionFunc
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "json not found",
			fields: fields{
				metricsFile: "notfound.json",
			},
			wantResult: assert.NotEmpty,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, os.ErrNotExist)
			},
		},
		{
			name: "storage error",
			fields: fields{
				catalogsFolder: "internal/testcatalogs/emptyTestFolder",
				storage:        &testutil.StorageWithError{SaveErr: ErrSomeError},
			},
			wantResult: assert.NotEmpty,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrSomeError)
			},
		},
		{
			name: "custom loading function with error",
			fields: fields{
				loadCatalogsFunc: func() ([]*orchestrator.Catalog, error) {
					return nil, ErrSomeError
				},
			},
			wantResult: assert.NotEmpty,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrSomeError)
			},
		},
		{
			name: "Happy path",
			fields: fields{
				catalogsFolder: "catalogs",
				storage:        testutil.NewInMemoryStorage(t),
			},
			wantResult: func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				svc, ok := i.(*Service)
				assert.True(t, ok)

				catalog := new(orchestrator.Catalog)
				err := svc.storage.Get(catalog, gorm.WithPreload("Categories.Controls", "parent_control_id IS NULL"), "Id = ?", "DemoCatalog")
				assert.NoError(t, err)

				err = catalog.Validate()
				return assert.NoError(t, err)
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
			}

			err := svc.loadCatalogs()

			// Validate the error via the ErrorAssertionFunc function
			tt.wantErr(t, err)

			// Validate the result via the ValueAssertionFunc function
			tt.wantResult(t, svc)
		})
	}
}
