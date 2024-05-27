package orchestrator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/logging"
	"clouditor.io/clouditor/v2/internal/util"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/persistence/gorm"
	"clouditor.io/clouditor/v2/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateCatalog implements a method for creating a new catalog.
func (svc *Service) CreateCatalog(_ context.Context, req *orchestrator.CreateCatalogRequest) (
	*orchestrator.Catalog, error) {
	// Validate request
	err := api.Validate(req)
	if err != nil {
		return nil, err
	}

	// Persist the new catalog in our database
	err = svc.storage.Create(req.Catalog)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	logging.LogRequest(log, logrus.DebugLevel, logging.Create, req)

	// Return catalog
	return req.Catalog, nil
}

// GetCatalog retrieves a control specified by the catalog ID, the control's category
// name and the control ID. If present, it also includes a list of sub-controls and any metrics associated to any
// controls.
func (svc *Service) GetCatalog(_ context.Context, req *orchestrator.GetCatalogRequest) (response *orchestrator.Catalog, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	response = new(orchestrator.Catalog)
	err = svc.storage.Get(response,
		// Preload fills in associated entities, in this case controls. We want to only select those controls which do
		// not have a parent, e.g., the top-level
		gorm.WithPreload("Categories.Controls", "parent_control_id IS NULL"),
		// Select catalog by ID
		"Id = ?", req.CatalogId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "catalog not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}
	return response, nil
}

// ListCatalogs Lists all security controls catalogs. Each catalog includes a list of its
// categories but no additional sub-resources.
func (svc *Service) ListCatalogs(_ context.Context, req *orchestrator.ListCatalogsRequest) (res *orchestrator.ListCatalogsResponse, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	res = new(orchestrator.ListCatalogsResponse)
	res.Catalogs, res.NextPageToken, err = service.PaginateStorage[*orchestrator.Catalog](req, svc.storage,
		service.DefaultPaginationOpts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
	}
	return
}

// UpdateCatalog implements a method for updating an existing catalog
func (svc *Service) UpdateCatalog(_ context.Context, req *orchestrator.UpdateCatalogRequest) (res *orchestrator.Catalog, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	res = req.Catalog

	err = svc.storage.Update(res, "id = ?", res.Id)

	if err != nil && errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "catalog not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	logging.LogRequest(log, logrus.DebugLevel, logging.Update, req)

	return
}

// RemoveCatalog implements a method for removing a catalog
func (svc *Service) RemoveCatalog(_ context.Context, req *orchestrator.RemoveCatalogRequest) (response *emptypb.Empty, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	err = svc.storage.Delete(&orchestrator.Catalog{}, "Id = ?", req.CatalogId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "catalog not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	logging.LogRequest(log, logrus.DebugLevel, logging.Remove, req)

	return &emptypb.Empty{}, nil
}

// GetCategory retrieves a category of a catalog specified by the catalog ID and the category name. It includes the
// first level of controls within each category.
func (srv *Service) GetCategory(_ context.Context, req *orchestrator.GetCategoryRequest) (res *orchestrator.Category, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	res = new(orchestrator.Category)
	err = srv.storage.Get(&res,
		// Preload fills in associated entities, in this case controls. We want to only select those controls which do
		// not have a parent, e.g., the top-level
		gorm.WithPreload("Controls", "parent_control_id IS NULL"),
		// Select the category by name and catalog ID
		"name = ? AND catalog_id = ?", req.CategoryName, req.CatalogId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "category not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	return res, nil
}

// GetControl retrieves a control specified by the catalog ID, the control's category name and the control ID. If
// present, it also includes a list of sub-controls and any metrics associated to the control.
func (srv *Service) GetControl(_ context.Context, req *orchestrator.GetControlRequest) (res *orchestrator.Control, err error) {
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	res = new(orchestrator.Control)
	err = srv.storage.Get(&res,
		// We only want to select controls for the specified category and catalog
		"Id = ? AND category_name = ? AND category_catalog_id = ?", req.ControlId, req.CategoryName, req.CatalogId)
	if errors.Is(err, persistence.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "control not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	return res, nil
}

// ListControls lists controls. If no additional parameters are specified, this lists all controls. If a catalog ID and
// a category name is specified, then only controls containing in this category are returned.
func (srv *Service) ListControls(_ context.Context, req *orchestrator.ListControlsRequest) (res *orchestrator.ListControlsResponse, err error) {
	var (
		args  []any
		query []string
	)
	// Validate request
	err = api.Validate(req)
	if err != nil {
		return nil, err
	}

	res = new(orchestrator.ListControlsResponse)

	// If the category name is set (additional binding), forward it as a condition to the pagination method
	if req.CategoryName != "" && req.CatalogId != "" {
		query = append(query, "category_name = ?")
		args = append(args, req.CategoryName)
	}
	if req.CatalogId != "" {
		query = append(query, "category_catalog_id = ?")
		args = append(args, req.CatalogId)
	}

	// Apply additional filter
	if req.Filter != nil {
		if len(req.Filter.AssuranceLevels) > 0 {
			query = append(query, "(assurance_level IN ? OR assurance_level IS NULL)")
			args = append(args, req.Filter.AssuranceLevels)
		}
	}

	// Join query with AND and prepend the query
	args = append([]any{strings.Join(query, " AND ")}, args...)

	res.Controls, res.NextPageToken, err = service.PaginateStorage[*orchestrator.Control](req, srv.storage,
		service.DefaultPaginationOpts, args...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not paginate results: %v", err)
	}
	return
}

// loadCatalogs loads catalog definitions from a JSON file.
func (svc *Service) loadCatalogs() (err error) {
	var catalogs []*orchestrator.Catalog

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

	log.Debug("Catalogs successfully stored.")

	return
}

func (svc *Service) loadEmbeddedCatalogs() (catalogs []*orchestrator.Catalog, err error) {
	var (
		b        []byte
		fileList []string
	)

	// Get all filenames
	fileList, err = util.GetJSONFilenames(svc.catalogsFolder)

	log.Infof("Loading catalogs from files %s.", strings.Join(fileList, ", "))

	// Get catalog for each file
	for i := range fileList {
		var catalogsFromFile []*orchestrator.Catalog

		b, err = os.ReadFile(filepath.Join(fileList[i]))
		if err != nil {
			log.Errorf("error while loading %s: %v", fileList[i], err)
			continue
		}

		err = json.Unmarshal(b, &catalogsFromFile)
		if err != nil {
			log.Errorf("error in JSON marshal for file %s: %v", fileList[i], err)
			continue
		}

		catalogs = append(catalogs, catalogsFromFile...)
	}

	// We need to make sure that sub-controls have the category_name and category_catalog_id of their parents set,
	// otherwise we are failing a constraint.
	for _, catalog := range catalogs {
		for _, category := range catalog.Categories {
			for _, control := range category.Controls {
				for _, sub := range control.Controls {
					sub.CategoryName = category.Name
					sub.CategoryCatalogId = catalog.Id

					// Make sure we are dealing with a copy of control when we
					// take an address of its property and not the loop var,
					// which gets overridden in each loop.
					control := control

					// Also set the parent information, so we do not need to set
					// it in the original file to make it easier
					sub.ParentControlCategoryCatalogId = &control.CategoryCatalogId
					sub.ParentControlCategoryName = &control.CategoryName
					sub.ParentControlId = &control.Id
				}
			}
		}
	}

	return
}
