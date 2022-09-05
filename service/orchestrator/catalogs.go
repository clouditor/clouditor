package orchestrator

import (
	"context"

	"clouditor.io/clouditor/api/orchestrator"
)

func (srv *Service) GetCategory(ctx context.Context, req *orchestrator.GetCategoryRequest) (res *orchestrator.Category, err error) {
	res = new(orchestrator.Category)
	err = srv.storage.Get(&res, "name = ? AND catalog_id = ?", req.CategoryName, req.CatalogId)
	if err != nil {
		return nil, err
	}

	return res, nil
}
