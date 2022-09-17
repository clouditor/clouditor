package orchestrator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/persistence/gorm"
	"clouditor.io/clouditor/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateCatalog implements method for creating a new catalog
func (svc *Service) CreateCatalog(_ context.Context, req *orchestrator.CreateCatalogRequest) (
	*orchestrator.Catalog, error) {
	// Validate request
	if req == nil {
		return nil,
			status.Errorf(codes.InvalidArgument, api.ErrRequestIsNil.Error())
	}
	if req.Catalog == nil {
		return nil,
			status.Errorf(codes.InvalidArgument, orchestrator.ErrCatalogIsNil.Error())
	}
	if req.Catalog.Id == "" {
		return nil,
			status.Errorf(codes.InvalidArgument, orchestrator.ErrCatalogIDIsMissing.Error())
	}

	// Persist the new catalog in our database
	err := svc.storage.Create(req.Catalog)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	// Return catalog
	return req.Catalog, nil
}

// GetCatalog implements method for getting a catalog, e.g. to show its state in the UI
func (svc *Service) GetCatalog(_ context.Context, req *orchestrator.GetCatalogRequest) (response *orchestrator.Catalog, err error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, api.ErrRequestIsNil.Error())
	}
	if req.CatalogId == "" {
		return nil, status.Errorf(codes.NotFound, orchestrator.ErrCatalogIDIsMissing.Error())
	}

	response = new(orchestrator.Catalog)
	err = svc.storage.Get(response, gorm.WithPreload("Categories.Controls", "parent_control_id IS NULL"), "Id = ?", req.CatalogId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "catalog not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}
	return response, nil
}

// ListCatalogs implements method for getting a catalog, e.g. to show its state in the UI
func (svc *Service) ListCatalogs(_ context.Context, req *orchestrator.ListCatalogsRequest) (res *orchestrator.ListCatalogsResponse, err error) {
	// Validate the request
	if err = api.ValidateListRequest[*orchestrator.Catalog](req); err != nil {
		err = fmt.Errorf("invalid request: %w", err)
		log.Error(err)
		err = status.Errorf(codes.InvalidArgument, "%v", err)
		return
	}

	res = new(orchestrator.ListCatalogsResponse)

	res.Catalogs, res.NextPageToken, err = service.PaginateStorage[*orchestrator.Catalog](req, svc.storage,
		service.DefaultPaginationOpts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
	}
	return
}

// UpdateCatalog implements method for updating an existing catalog
func (svc *Service) UpdateCatalog(_ context.Context, req *orchestrator.UpdateCatalogRequest) (res *orchestrator.Catalog, err error) {
	if req.CatalogId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "catalog id is empty")
	}

	if req.Catalog == nil {
		return nil, status.Errorf(codes.InvalidArgument, "catalog is empty")
	}

	res = req.Catalog
	res.Id = req.CatalogId

	err = svc.storage.Update(res, "id = ?", res.Id)

	if err != nil && errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "catalog not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}
	return
}

// RemoveCatalog implements method for removing a catalog
func (svc *Service) RemoveCatalog(_ context.Context, req *orchestrator.RemoveCatalogRequest) (response *emptypb.Empty, err error) {
	if req.CatalogId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "catalog id is empty")
	}

	err = svc.storage.Delete(&orchestrator.Catalog{}, "Id = ?", req.CatalogId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "catalog not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (srv *Service) GetCategory(ctx context.Context, req *orchestrator.GetCategoryRequest) (res *orchestrator.Category, err error) {
	res = new(orchestrator.Category)
	err = srv.storage.Get(&res, gorm.WithPreload("Controls", "parent_control_id IS NULL"), "name = ? AND catalog_id = ?", req.CategoryName, req.CatalogId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "category not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	return res, nil
}

func (srv *Service) GetControl(ctx context.Context, req *orchestrator.GetControlRequest) (res *orchestrator.Control, err error) {
	res = new(orchestrator.Control)
	err = srv.storage.Get(&res, "Id = ? AND category_name = ? AND category_catalog_id = ?", req.ControlId, req.CategoryName, req.CatalogId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "control not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	return res, nil
}

func (srv *Service) ListControls(ctx context.Context, req *orchestrator.ListControlsRequest) (res *orchestrator.ListControlsResponse, err error) {
	// Validate the request
	if err = api.ValidateListRequest[*orchestrator.Control](req); err != nil {
		err = fmt.Errorf("invalid request: %w", err)
		log.Error(err)
		err = status.Errorf(codes.InvalidArgument, "%v", err)
		return
	}

	res = new(orchestrator.ListControlsResponse)

	// If the category name is set (additional binding), forward it as a condition to the pagination method
	if req.CategoryName != "" && req.CatalogId != "" {
		res.Controls, res.NextPageToken, err = service.PaginateStorage[*orchestrator.Control](req, srv.storage,
			service.DefaultPaginationOpts, "category_name = ? AND catalog_id = ?", req.CategoryName, req.CatalogId)
	} else {
		res.Controls, res.NextPageToken, err = service.PaginateStorage[*orchestrator.Control](req, srv.storage,
			service.DefaultPaginationOpts)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
	}
	return
}

// LoadCatalogs loads catalog definitions from a JSON file.
func (svc *Service) loadCatalogs() (err error) {
	var catalogs []*orchestrator.Catalog

	log.Infof("Loading catalogs from %s", svc.catalogsFile)

	// Default to loading catalogs from our embedded file system
	if svc.loadCatalogsFunc == nil {
		svc.loadCatalogsFunc = svc.loadEmbeddedCatalogs
	}

	// Execute our catalogs loading function
	catalogs, err = svc.loadCatalogsFunc()
	if err != nil {
		return fmt.Errorf("could not load catalogs: %w", err)
	}

	err = svc.storage.Save(catalogs)
	if err != nil {
		log.Errorf("Error while saving catalog %v", err)
	}

	return
}

func (svc *Service) loadEmbeddedCatalogs() (catalogs []*orchestrator.Catalog, err error) {
	var b []byte

	b, err = f.ReadFile(svc.catalogsFile)
	if err != nil {
		return nil, fmt.Errorf("error while loading %s: %w", svc.catalogsFile, err)
	}

	err = json.Unmarshal(b, &catalogs)
	if err != nil {
		return nil, fmt.Errorf("error in JSON marshal: %w", err)
	}

	// We need to make sure that sub-controls have the category_name and category_catalog_id of their parents set, otherwise we are failing a constraint.
	for _, catalog := range catalogs {
		for _, category := range catalog.Categories {
			for _, control := range category.Controls {
				for _, sub := range control.Controls {
					sub.CategoryName = category.Name
					sub.CategoryCatalogId = catalog.Id
				}
			}
		}
	}
	return
}
