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
	"testing"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/orchestratortest"
	"clouditor.io/clouditor/v2/internal/util"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/persistence/gorm"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestService_CreateCatalog(t *testing.T) {
	// Mock catalogs
	mockCatalogWithoutMetadata := orchestratortest.NewCatalog()
	mockCatalogWithMetadata := orchestratortest.NewCatalog()
	mockCatalogWithMetadata.Metadata = &orchestrator.Catalog_Metadata{
		Color: util.Ref("#007FC3"),
	}
	mockCatalogWithoutID := orchestratortest.NewCatalog()
	mockCatalogWithoutID.Id = ""

	type args struct {
		in0 context.Context
		req *orchestrator.CreateCatalogRequest
	}
	tests := []struct {
		name    string
		args    args
		wantRes *orchestrator.Catalog
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "missing request",
			args: args{
				context.Background(),
				nil,
			},
			wantRes: nil,
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
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.ErrorContains(t, err, "catalog: value is required")
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
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.ErrorContains(t, err, "catalog.id: value length must be at least 1 characters")
				return assert.Equal(t, status.Code(err), codes.InvalidArgument)
			},
		},
		{
			name: "Happy path: without metadata",
			args: args{
				context.Background(),
				&orchestrator.CreateCatalogRequest{
					Catalog: mockCatalogWithoutMetadata,
				},
			},
			wantRes: mockCatalogWithoutMetadata,
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: with metadata",
			args: args{
				context.Background(),
				&orchestrator.CreateCatalogRequest{
					Catalog: mockCatalogWithMetadata,
				},
			},
			wantRes: mockCatalogWithMetadata,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			gotRes, err := s.CreateCatalog(tt.args.in0, tt.args.req)
			tt.wantErr(t, err)

			// If no error is wanted, check response
			if !assert.Equal(t, tt.wantRes, gotRes) {
				// Check catalog structure with validation method
				assert.NoError(t, api.Validate(gotRes))
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
		wantResponse assert.Want[*orchestrator.Catalog]
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "invalid request",
			fields: fields{storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
				// Create Catalog
				assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
			})},
			wantResponse: assert.Nil[*orchestrator.Catalog],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "empty request")
			},
		},
		{
			name:         "catalog ID empty",
			fields:       fields{storage: testutil.NewInMemoryStorage(t)},
			args:         args{req: &orchestrator.GetCatalogRequest{CatalogId: ""}},
			wantResponse: assert.Nil[*orchestrator.Catalog],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "catalog_id: value length must be at least 1 characters")
			},
		},
		{
			name: "catalog not found",
			fields: fields{storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
				// Create Catalog
				assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
			})},
			args:         args{req: &orchestrator.GetCatalogRequest{CatalogId: "a"}},
			wantResponse: assert.Nil[*orchestrator.Catalog],
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
			args: args{req: &orchestrator.GetCatalogRequest{CatalogId: testdata.MockCatalogID1}},
			wantResponse: func(t *testing.T, got *orchestrator.Catalog) bool {
				want := orchestratortest.NewCatalog()
				assert.Equal(t, 1, len(got.Categories))
				return assert.Equal(t, want.Id, got.Id)
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
			tt.wantErr(t, err)
			tt.wantResponse(t, res)
		})
	}
}

func TestService_ListCatalogs(t *testing.T) {
	var (
		listCatalogsResponse *orchestrator.ListCatalogsResponse
		err                  error
	)

	orchestratorService := NewService(WithCatalogsFolder("internal/testdata/catalogs"))
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
			Id:              testdata.MockCatalogID1,
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
	orchestratorService := NewService(WithCatalogsFolder("internal/testdata/empty_catalogs"))

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
		TargetOfEvaluationHooks []orchestrator.TargetOfEvaluationHookFunc
		AssessmentResultHooks   []assessment.ResultHookFunc
		storage                 persistence.Storage
		catalogsFile            string
		loadCatalogsFunc        func() ([]*orchestrator.Catalog, error)
		events                  chan *orchestrator.MetricChangeEvent
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
				req: &orchestrator.GetCategoryRequest{CatalogId: testdata.MockCatalogID1, CategoryName: testdata.MockCategoryName},
			},
			wantRes: &orchestrator.Category{
				Name:        testdata.MockCategoryName,
				Description: testdata.MockCategoryDescription,
				CatalogId:   testdata.MockCatalogID1,
				Controls: []*orchestrator.Control{
					{
						Id:                testdata.MockControlID1,
						Name:              testdata.MockControlName,
						CategoryName:      testdata.MockCategoryName,
						CategoryCatalogId: testdata.MockCatalogID1,
						Description:       testdata.MockControlDescription,
						Controls:          []*orchestrator.Control{},
					},
					{
						Id:                testdata.MockControlID2,
						Name:              testdata.MockControlName,
						CategoryName:      testdata.MockCategoryName,
						CategoryCatalogId: testdata.MockCatalogID1,
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
				TargetOfEvaluationHooks: tt.fields.TargetOfEvaluationHooks,
				AssessmentResultHooks:   tt.fields.AssessmentResultHooks,
				storage:                 tt.fields.storage,
				catalogsFolder:          tt.fields.catalogsFile,
				loadCatalogsFunc:        tt.fields.loadCatalogsFunc,
				events:                  tt.fields.events,
			}
			gotRes, err := srv.GetCategory(tt.args.ctx, tt.args.req)

			tt.wantErr(t, err)
			assert.Equal(t, tt.wantRes, gotRes)
		})
	}
}

func TestService_GetControl(t *testing.T) {
	type fields struct {
		TargetOfEvaluationHooks []orchestrator.TargetOfEvaluationHookFunc
		AssessmentResultHooks   []assessment.ResultHookFunc
		storage                 persistence.Storage
		catalogsFolder          string
		loadCatalogsFunc        func() ([]*orchestrator.Catalog, error)
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
			args:    args{req: &orchestrator.GetControlRequest{CatalogId: testdata.MockCatalogID1, CategoryName: testdata.MockCategoryName, ControlId: "WrongControlID"}},
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
			args: args{req: &orchestrator.GetControlRequest{CatalogId: testdata.MockCatalogID1, CategoryName: testdata.MockCategoryName, ControlId: testdata.MockControlID1}},
			wantRes: &orchestrator.Control{
				Id:                testdata.MockControlID1,
				CategoryName:      testdata.MockCategoryName,
				CategoryCatalogId: testdata.MockCatalogID1,
				Name:              testdata.MockControlName,
				Description:       testdata.MockControlDescription,
				Controls: []*orchestrator.Control{{
					Id:                             testdata.MockSubControlID11,
					Name:                           testdata.MockSubControlName,
					Description:                    testdata.MockSubControlDescription,
					Metrics:                        []*assessment.Metric{}, // metrics on sub-controls are not returned
					CategoryName:                   testdata.MockCategoryName,
					CategoryCatalogId:              testdata.MockCatalogID1,
					AssuranceLevel:                 &testdata.AssuranceLevelBasic,
					ParentControlId:                util.Ref(testdata.MockControlID1),
					ParentControlCategoryCatalogId: util.Ref(testdata.MockCatalogID1),
					ParentControlCategoryName:      util.Ref(testdata.MockCategoryName),
				}},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Service{
				TargetOfEvaluationHooks: tt.fields.TargetOfEvaluationHooks,
				AssessmentResultHooks:   tt.fields.AssessmentResultHooks,
				storage:                 tt.fields.storage,
				catalogsFolder:          tt.fields.catalogsFolder,
				loadCatalogsFunc:        tt.fields.loadCatalogsFunc,
			}
			gotRes, err := srv.GetControl(tt.args.ctx, tt.args.req)

			tt.wantErr(t, err)
			assert.Equal(t, tt.wantRes, gotRes)
		})
	}
}

func TestService_ListControls(t *testing.T) {
	var (
		res *orchestrator.ListControlsResponse
		err error
		c   *orchestrator.Control
		sub *orchestrator.Control
	)

	orchestratorService := NewService(WithCatalogsFolder("internal/testdata/empty_catalogs"))
	// 1st case: No Controls stored
	res, err = orchestratorService.ListControls(context.Background(), &orchestrator.ListControlsRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, res.Controls)
	assert.Empty(t, res.Controls)

	// 2nd case: 3 controls stored; note that we do not have to create an extra catalog/control since NewService above already loads the default catalogs/controls
	orchestratorService = NewService(WithCatalogsFolder("internal/testdata/catalogs"))
	res, err = orchestratorService.ListControls(context.Background(), &orchestrator.ListControlsRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, res.Controls)
	assert.NotEmpty(t, res.Controls)
	// there are 3 default controls
	assert.Equal(t, 3, len(res.Controls))

	// 3th case: List controls for a specific catalog and category.
	res, err = orchestratorService.ListControls(context.Background(), &orchestrator.ListControlsRequest{
		CatalogId:    "TestCatalog",
		CategoryName: "Secure Category",
	})
	assert.NoError(t, err)
	assert.NotNil(t, res.Controls)
	assert.NotEmpty(t, res.Controls)
	assert.Equal(t, 3, len(res.Controls))

	// Make sure, that control parent information is set correctly
	c = res.Controls[0]
	assert.Equal(t, 2, len(c.Controls))

	sub = c.Controls[0]
	assert.Equal(t, c.Id, util.Deref(sub.ParentControlId))
	assert.Equal(t, c.CategoryName, util.Deref(sub.ParentControlCategoryName))
	assert.Equal(t, c.CategoryCatalogId, util.Deref(sub.ParentControlCategoryCatalogId))

	// 4th case: Filter by assurance level
	res, err = orchestratorService.ListControls(context.Background(), &orchestrator.ListControlsRequest{
		CatalogId:    "TestCatalog",
		CategoryName: "Secure Category",
		Filter: &orchestrator.ListControlsRequest_Filter{
			AssuranceLevels: []string{"substantial"},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, res.Controls)
	assert.NotEmpty(t, res.Controls)
	assert.Equal(t, 2, len(res.Controls))

	res, err = orchestratorService.ListControls(context.Background(), &orchestrator.ListControlsRequest{
		CatalogId:    "TestCatalog",
		CategoryName: "Secure Category",
		Filter: &orchestrator.ListControlsRequest_Filter{
			AssuranceLevels: []string{"substantial", "high"},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, res.Controls)
	assert.NotEmpty(t, res.Controls)
	assert.Equal(t, 3, len(res.Controls))
}

func TestService_loadCatalogs(t *testing.T) {
	type fields struct {
		AssessmentResultHooks []assessment.ResultHookFunc
		storage               persistence.Storage
		catalogsFolder        string
		loadCatalogsFunc      func() ([]*orchestrator.Catalog, error)
		events                chan *orchestrator.MetricChangeEvent
	}
	tests := []struct {
		name    string
		fields  fields
		wantSvc assert.Want[*Service]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "storage error",
			fields: fields{
				catalogsFolder: "internal/testdata/empty_catalogs",
				storage:        &testutil.StorageWithError{SaveErr: ErrSomeError},
			},
			wantSvc: assert.NotNil[*Service],
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
			wantSvc: assert.NotNil[*Service],
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
			wantSvc: func(t *testing.T, got *Service) bool {
				catalog := new(orchestrator.Catalog)
				err := got.storage.Get(catalog, gorm.WithPreload("Categories.Controls", "parent_control_id IS NULL"), "Id = ?", "DemoCatalog")
				assert.NoError(t, err)

				err = api.Validate(catalog)
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
				catalogsFolder:        tt.fields.catalogsFolder,
				loadCatalogsFunc:      tt.fields.loadCatalogsFunc,
				events:                tt.fields.events,
			}

			err := svc.loadCatalogs()
			tt.wantErr(t, err)
			tt.wantSvc(t, svc)
		})
	}
}
