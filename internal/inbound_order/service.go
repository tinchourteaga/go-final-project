package inbound_order

import (
	"context"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
)

var (
	//ErrEmployeeWithInboundOrdersNotFound is returned when the requested result are not found
	ErrEmployeeWithInboundOrdersNotFound = errors.New("employee not found")
	//ErrInboundOrderAlreadyExists is returned when the requested result already exists
	ErrInboundOrderAlreadyExists = errors.New("inbound order already exists")
	//ErrInboundOrderNotSaved is returned when the requested result cannot be saved in db
	ErrInboundOrderNotSaved = errors.New("cannot save inbound order")
	// ErrEmptyOrderNumber is returned when the request does not have an orderNumber
	ErrEmptyOrderNumber = errors.New("orderNumber field cannot be empty")
)

type Service interface {
	// GetAllEmployeesInboundOrders returns all the employees and their inbound orders that exist and are inside the repository
	GetAllEmployeesInboundOrders(ctx context.Context) ([]domain.EmployeeWithInboundOrders, error)
	// GetEmployeeInboundOrders returns employee and its inbound orders if it exists inside the repository
	GetEmployeeInboundOrders(ctx context.Context, id int) (domain.EmployeeWithInboundOrders, error)
	// Save creates a new inbound order with the specified data inside the repository
	Save(ctx context.Context, inboundOrder domain.InboundOrder) (domain.InboundOrder, error)
}

type service struct {
	repository Repository
}

func NewService(repo Repository) Service {
	return &service{
		repository: repo,
	}
}

func (service *service) GetAllEmployeesInboundOrders(ctx context.Context) ([]domain.EmployeeWithInboundOrders, error) {
	employees, err := service.repository.GetAllEmployeesInboundOrders(ctx)

	if err != nil {
		logging.Log(err)
		return nil, err
	}

	return employees, nil
}

func (service *service) GetEmployeeInboundOrders(ctx context.Context, id int) (domain.EmployeeWithInboundOrders, error) {
	var emptyEmployee domain.EmployeeWithInboundOrders
	employee, err := service.repository.GetEmployeeInboundOrders(ctx, id)

	if err != nil && employee == emptyEmployee {
		logging.Log(ErrEmployeeWithInboundOrdersNotFound)
		return domain.EmployeeWithInboundOrders{}, ErrEmployeeWithInboundOrdersNotFound
	}
	if err != nil {
		logging.Log(err)
		return domain.EmployeeWithInboundOrders{}, err
	}

	return employee, nil
}

func (service *service) Save(ctx context.Context, inboundOrder domain.InboundOrder) (domain.InboundOrder, error) {
	id, err := service.repository.Save(ctx, inboundOrder)

	inboundOrder.ID = id

	if err != nil {
		logging.Log(err)
		switch err.Error() {
		case ErrEmptyOrderNumber.Error():
			return domain.InboundOrder{}, ErrEmptyOrderNumber
		case ErrEmployeeNonExistent.Error():
			return domain.InboundOrder{}, ErrEmployeeNonExistent
		case ErrWarehouseNonExistent.Error():
			return domain.InboundOrder{}, ErrWarehouseNonExistent
		case ErrProductBatchNonExistent.Error():
			return domain.InboundOrder{}, ErrProductBatchNonExistent
		case ErrInboundOrderAlreadyExists.Error():
			return domain.InboundOrder{}, ErrInboundOrderAlreadyExists
		default:
			return domain.InboundOrder{}, ErrInboundOrderNotSaved
		}
	}

	return inboundOrder, nil
}
