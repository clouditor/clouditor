package orchestrator

import (
	"context"
	"fmt"
	"os"
	"testing"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/internal/testutil/orchestratortest"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/persistence/gorm"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
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

	orchestratorService := NewService()
	// 1st case: Default catalogs stored
	listCatalogsResponse, err = orchestratorService.ListCatalogs(context.Background(), &orchestrator.ListCatalogsRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listCatalogsResponse.Catalogs)
	assert.Equal(t, len(listCatalogsResponse.Catalogs), 1)

	// 2nd case: Invalid request
	_, err = orchestratorService.ListCatalogs(context.Background(),
		&orchestrator.ListCatalogsRequest{OrderBy: "not a field"})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Contains(t, err.Error(), api.ErrInvalidColumnName.Error())
}

func TestService_UpdateCatalog(t *testing.T) {
	var (
		catalog *orchestrator.Catalog
		err     error
	)
	orchestratorService := NewService()

	// 1st case: Certificate is nil
	_, err = orchestratorService.UpdateCatalog(context.Background(), &orchestrator.UpdateCatalogRequest{})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// 2nd case: Certificate ID is nil
	_, err = orchestratorService.UpdateCatalog(context.Background(), &orchestrator.UpdateCatalogRequest{
		Catalog: catalog,
	})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// 3rd case: Certificate not found since there are no certificates yet
	_, err = orchestratorService.UpdateCatalog(context.Background(), &orchestrator.UpdateCatalogRequest{
		Catalog: &orchestrator.Catalog{
			Id:   testdata.MockCatalogID,
			Name: testdata.MockCatalogName,
		},
	})
	assert.Equal(t, codes.NotFound, status.Code(err))
	assert.ErrorContains(t, err, "catalog not found")

	// 4th case: Certificate updated successfully
	mockCatalog := orchestratortest.NewCatalog()
	err = orchestratorService.storage.Create(mockCatalog)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	// update the certificate's description and send the update request
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
	orchestratorService := NewService()

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

	// There are two catalogs in the db now (one default plus NewCatalog)
	listCatalogsResponse, err = orchestratorService.ListCatalogs(context.Background(), &orchestrator.ListCatalogsRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listCatalogsResponse.Catalogs)
	assert.Equal(t, 2, len(listCatalogsResponse.Catalogs))

	// Remove record
	_, err = orchestratorService.RemoveCatalog(context.Background(), &orchestrator.RemoveCatalogRequest{CatalogId: mockCatalog.Id})
	assert.NoError(t, err)

	// There are two records left in the DB
	listCatalogsResponse, err = orchestratorService.ListCatalogs(context.Background(), &orchestrator.ListCatalogsRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listCatalogsResponse.Catalogs)
	assert.Equal(t, 1, len(listCatalogsResponse.Catalogs))
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
		wantErr bool
	}{
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
				Controls: []*orchestrator.Control{{
					Id:                testdata.MockControlID,
					Name:              testdata.MockControlName,
					Description:       testdata.MockControlDescription,
					CategoryName:      testdata.MockCategoryName,
					CategoryCatalogId: testdata.MockCatalogID,
					// at this level, we will not have the metrics
					Metrics: []*assessment.Metric{},
					// at this level, we will not have the sub-controls
					Controls: []*orchestrator.Control{},
				}, {
					Id:                testdata.MockAnotherControlID,
					Name:              testdata.MockAnotherControlName,
					CategoryName:      testdata.MockCategoryName,
					CategoryCatalogId: testdata.MockCatalogID,
				}},
			},
			wantErr: false,
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
				catalogsFile:          tt.fields.catalogsFile,
				loadCatalogsFunc:      tt.fields.loadCatalogsFunc,
				events:                tt.fields.events,
			}
			gotRes, err := srv.GetCategory(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
		catalogsFile          string
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
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					// Create Catalog
					assert.NoError(t, s.Create(orchestratortest.NewCatalog()))
				})},
			args: args{req: &orchestrator.GetControlRequest{CatalogId: testdata.MockCatalogID, CategoryName: testdata.MockCategoryName, ControlId: testdata.MockControlID}},
			wantRes: &orchestrator.Control{
				Id:                testdata.MockControlID,
				CategoryName:      testdata.MockCategoryName,
				CategoryCatalogId: testdata.MockCatalogID,
				Name:              testdata.MockControlName,
				Description:       testdata.MockControlDescription,
				Metrics: []*assessment.Metric{{
					Id:          testdata.MockMetricID,
					Name:        testdata.MockMetricName,
					Description: testdata.MockMetricDescription,
					Scale:       assessment.Metric_ORDINAL,
					Range: &assessment.Range{
						Range: &assessment.Range_AllowedValues{AllowedValues: &assessment.AllowedValues{
							Values: []*structpb.Value{
								structpb.NewBoolValue(false),
								structpb.NewBoolValue(true),
							}}}},
				}},
				Controls: []*orchestrator.Control{{
					Id:                             testdata.MockSubControlID,
					Name:                           testdata.MockSubControlName,
					Description:                    testdata.MockSubControlDescription,
					Metrics:                        []*assessment.Metric{},
					CategoryName:                   testdata.MockCategoryName,
					CategoryCatalogId:              testdata.MockCatalogID,
					ParentControlId:                util.Ref(testdata.MockControlID),
					ParentControlCategoryCatalogId: util.Ref(testdata.MockCatalogID),
					ParentControlCategoryName:      util.Ref(testdata.MockCategoryName),
				}},
			},
			wantErr: false,
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
				catalogsFile:          tt.fields.catalogsFile,
				loadCatalogsFunc:      tt.fields.loadCatalogsFunc,
			}
			gotRes, err := srv.GetControl(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetControl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

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

	orchestratorService := NewService()
	// 1st case: No Controls stored
	listControlsResponse, err = orchestratorService.ListControls(context.Background(), &orchestrator.ListControlsRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listControlsResponse.Controls)
	assert.NotEmpty(t, listControlsResponse.Controls)

	// 2nd case: Two controls stored; note that we do not have to create an extra catalog/control since NewService above already loads the default catalogs/controls
	listControlsResponse, err = orchestratorService.ListControls(context.Background(), &orchestrator.ListControlsRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listControlsResponse.Controls)
	assert.NotEmpty(t, listControlsResponse.Controls)
	// there are two default controls
	assert.Equal(t, len(listControlsResponse.Controls), 2)

	// 3rd case: List controls for a specific catalog; first, create a new catalog with two controls, but only request controls for one of the existing catalogs
	orchestratortest.NewCatalog()
	listControlsResponse, err = orchestratorService.ListControls(context.Background(), &orchestrator.ListControlsRequest{
		CatalogId: "EUCS",
	})
	assert.NoError(t, err)
	assert.NotNil(t, listControlsResponse.Controls)
	assert.NotEmpty(t, listControlsResponse.Controls)
	// there are the two default controls
	assert.Equal(t, len(listControlsResponse.Controls), 2)

	// 4th case: Invalid request
	_, err = orchestratorService.ListControls(context.Background(),
		&orchestrator.ListControlsRequest{OrderBy: "not a field"})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Contains(t, err.Error(), api.ErrInvalidColumnName.Error())
}

func TestService_loadCatalogs(t *testing.T) {
	type fields struct {
		AssessmentResultHooks []func(result *assessment.AssessmentResult, err error)
		storage               persistence.Storage
		metricsFile           string
		loadMetricsFunc       func() ([]*assessment.Metric, error)
		catalogsFile          string
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
				catalogsFile: "demo_catalogs.json",
				storage:      &testutil.StorageWithError{SaveErr: ErrSomeError},
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
				catalogsFile: "demo_catalogs.json",
				storage:      testutil.NewInMemoryStorage(t),
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
				catalogsFile:          tt.fields.catalogsFile,
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
