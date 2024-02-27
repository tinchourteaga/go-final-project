package warehouse

import (
	"context"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
)

// Service provides the public methods of a warehouse service.
type Service interface {
	Get(ctx context.Context, id int) (domain.Warehouse, error)
	GetAll(ctx context.Context) ([]domain.Warehouse, error)
	Create(ctx context.Context, address string, telephone string, warehouseCode string, minimumCapacity int, minimumTemperature int) (domain.Warehouse, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, address *string, telephone *string, warehouseCode *string, minimumCapacity *int, minimumTemperature *int) (domain.Warehouse, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

// Get returns the collection of warehouses provided by the repository.
// any error encountered is also returned.
func (s *service) Get(ctx context.Context, id int) (domain.Warehouse, error) {
	warehouse, err := s.repository.Get(ctx, id)
	if err != nil {
		logging.Log(err)
		return domain.Warehouse{}, err
	}
	return warehouse, nil
}

// GetAll returns a warehouse provided by the repository if succesful.
// any error encountered is also returned.
func (s *service) GetAll(ctx context.Context) ([]domain.Warehouse, error) {
	warehouses, err := s.repository.GetAll(ctx)
	if err != nil {
		logging.Log(err)
		return nil, err
	}
	return warehouses, nil
}

// Create returns the created warehouse provided by the repository if succesful.
// if the warehouseCode is not unique, a error is returned.
// any other error encountered is also returned.
func (s *service) Create(ctx context.Context, address string, telephone string, warehouseCode string, minimumCapacity int, minimumTemperature int) (domain.Warehouse, error) {
	warehouseExists := s.repository.Exists(ctx, warehouseCode)
	if warehouseExists {
		logging.Log(ErrAlreadyExists)
		return domain.Warehouse{}, ErrAlreadyExists
	}

	warehouse := domain.Warehouse{
		Address:            address,
		Telephone:          telephone,
		WarehouseCode:      warehouseCode,
		MinimumCapacity:    minimumCapacity,
		MinimumTemperature: minimumTemperature,
	}
	warehouseID, err := s.repository.Save(ctx, warehouse)
	if err != nil {
		logging.Log(err)
		return domain.Warehouse{}, err
	}

	warehouse.ID = warehouseID
	return warehouse, nil
}

// Delete returns an error if the deletion of the warehouse failed.
// if a warehouse with the given id doesn`t exist, an error is returned.
// any other error encountered is also returned.
func (s *service) Delete(ctx context.Context, id int) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		logging.Log(err)
		return err
	}
	return nil
}

// Update returns the updated warehouse if successful.
// if a warehouse with the given id doesn't exist, an error is returned.
// if the warehouseCode is not unique (with exception to the warehouse currently updating), an error is returned.
// any other error encountered is also returned.
// only the values not in a null state are updated.
func (s *service) Update(ctx context.Context, id int, address *string, telephone *string, warehouseCode *string, minimumCapacity *int, minimumTemperature *int) (domain.Warehouse, error) {
	// Get Original Warehouse
	warehouse, err := s.repository.Get(ctx, id)
	if err != nil {
		logging.Log(err)
		return domain.Warehouse{}, err
	}

	// Valid warehouseCode not exists in another Warehouse
	if warehouseCode != nil {
		if warehouse.WarehouseCode != *warehouseCode {
			warehouseExists := s.repository.Exists(ctx, *warehouseCode)
			if warehouseExists {
				logging.Log(ErrAlreadyExists)
				return domain.Warehouse{}, ErrAlreadyExists
			}
		}
		warehouse.WarehouseCode = *warehouseCode // update with a validated warehouseCode
	}

	// Change all other valid entries
	if address != nil {
		warehouse.Address = *address
	}
	if telephone != nil {
		warehouse.Telephone = *telephone
	}
	if minimumCapacity != nil {
		warehouse.MinimumCapacity = *minimumCapacity
	}
	if minimumTemperature != nil {
		warehouse.MinimumTemperature = *minimumTemperature
	}

	// Update only valid entries
	err = s.repository.Update(ctx, warehouse)
	if err != nil {
		logging.Log(err)
		return domain.Warehouse{}, err
	}
	return warehouse, nil
}
