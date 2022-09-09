package orchestrator

import (
	"context"

	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/persistence/gorm"
)

func (srv *Service) GetCategory(ctx context.Context, req *orchestrator.GetCategoryRequest) (res *orchestrator.Category, err error) {
	res = new(orchestrator.Category)
	err = srv.storage.Get(&res, gorm.WithPreload("Controls", "parent_control_short_name IS NULL"), "name = ? AND catalog_id = ?", req.CategoryName, req.CatalogId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (srv *Service) GetControl(ctx context.Context, req *orchestrator.GetControlRequest) (res *orchestrator.Control, err error) {
	res = new(orchestrator.Control)
	err = srv.storage.Get(&res, "short_name = ? AND category_name = ? AND category_catalog_id = ?", req.ControlShortName, req.CategoryName, req.CatalogId)
	if err != nil {
		return nil, err
	}

	return res, nil
}
