package orchestrator

import (
	"context"
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
	err = svc.storage.Get(response, gorm.WithPreload("Categories.Controls", "parent_control_short_name IS NULL"), "Id = ?", req.CatalogId)
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
	err = srv.storage.Get(&res, gorm.WithPreload("Controls", "parent_control_short_name IS NULL"), "name = ? AND catalog_id = ?", req.CategoryName, req.CatalogId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "category not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	return res, nil
}

func (srv *Service) GetControl(ctx context.Context, req *orchestrator.GetControlRequest) (res *orchestrator.Control, err error) {
	res = new(orchestrator.Control)
	err = srv.storage.Get(&res, "short_name = ? AND category_name = ? AND category_catalog_id = ?", req.ControlShortName, req.CategoryName, req.CatalogId)
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
	if req.CategoryName != "" {
		res.Controls, res.NextPageToken, err = service.PaginateStorage[*orchestrator.Control](req, srv.storage,
			service.DefaultPaginationOpts, "category_name = ?", req.CategoryName)
	} else {
		res.Controls, res.NextPageToken, err = service.PaginateStorage[*orchestrator.Control](req, srv.storage,
			service.DefaultPaginationOpts)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
	}
	return
}
