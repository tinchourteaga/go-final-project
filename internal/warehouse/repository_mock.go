package warehouse

import (
	"context"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
)

type MockRepo struct {
	mockWarehouse     domain.Warehouse
	mockWarehouses    []domain.Warehouse
	mockErrorInternal error
	mockErrorExists   error
	mockErrorUpdate   error
}

func (r *MockRepo) GetAll(ctx context.Context) ([]domain.Warehouse, error) {
	if r.mockErrorInternal != nil {
		return []domain.Warehouse{}, r.mockErrorInternal
	}
	return r.mockWarehouses, nil
}

func (r *MockRepo) Get(ctx context.Context, id int) (domain.Warehouse, error) {
	if r.mockErrorInternal != nil {
		return domain.Warehouse{}, r.mockErrorInternal
	}
	return r.mockWarehouse, nil
}

func (r *MockRepo) Exists(ctx context.Context, warehouseCode string) bool {
	return r.mockErrorExists != nil
}

func (r *MockRepo) Save(ctx context.Context, w domain.Warehouse) (int, error) {
	if r.mockErrorInternal != nil {
		return 0, r.mockErrorInternal
	}
	return 1, nil
}

func (r *MockRepo) Update(ctx context.Context, w domain.Warehouse) error {
	if r.mockErrorUpdate != nil {
		return r.mockErrorUpdate
	}
	return nil
}

func (r *MockRepo) Delete(ctx context.Context, id int) error {
	if r.mockErrorInternal != nil {
		return r.mockErrorInternal
	}
	return nil
}
