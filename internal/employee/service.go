package employee

import (
	"context"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
)

// Errors
var (
	//ErrEmployeeNotFound is returned when the requested result are not found
	ErrEmployeeNotFound = errors.New("employee not found")
	//ErrEmployeeAlreadyExists is returned when the requested result already exists
	ErrEmployeeAlreadyExists = errors.New("employee already exists")
	//ErrEmployeeNotUpdated is returned when the requested result cannot be updated in db
	ErrEmployeeNotUpdated = errors.New("cannot update employee")
	//ErrEmployeedNotSaved is returned when the requested result cannot be saved in db
	ErrEmployeeNotSaved = errors.New("cannot save employee")
)

type Service interface {
	// GetAll returns all the employees that exist and are inside the repository
	GetAll(ctx context.Context) ([]domain.Employee, error)
	// Get returns employee with the specified ID if it exists inside the repository
	Get(ctx context.Context, id int) (domain.Employee, error)
	// Save creates a new employee with the specified data inside the repository
	Save(ctx context.Context, employee domain.Employee) (domain.Employee, error)
	// Update updates the employee data inside the repository
	Update(ctx context.Context, employee domain.Employee) (domain.Employee, error)
	// Delete removes the employee with the specified ID from the repository
	Delete(ctx context.Context, id int) error
}

type service struct {
	repository Repository
}

func NewService(repo Repository) Service {
	return &service{
		repository: repo,
	}
}

func (service *service) GetAll(ctx context.Context) ([]domain.Employee, error) {
	employees, err := service.repository.GetAll(ctx)

	if err != nil {
		logging.Log(err)
		return nil, err
	}

	return employees, nil
}

func (service *service) Get(ctx context.Context, id int) (domain.Employee, error) {
	var emptyEmployee domain.Employee
	employee, err := service.repository.Get(ctx, id)

	if err != nil && employee == emptyEmployee {
		logging.Log(ErrEmployeeNotFound)
		return domain.Employee{}, ErrEmployeeNotFound
	}
	if err != nil {
		logging.Log(err)
		return domain.Employee{}, err
	}

	return employee, nil
}

func (service *service) Save(ctx context.Context, employee domain.Employee) (domain.Employee, error) {
	if service.repository.Exists(ctx, employee.CardNumberID) {
		logging.Log(ErrEmployeeAlreadyExists)
		return domain.Employee{}, ErrEmployeeAlreadyExists
	}

	id, err := service.repository.Save(ctx, employee)

	employee.ID = id

	if err != nil {
		switch err.Error() {
		case ErrWarehouseNonExistent.Error():
			logging.Log(ErrWarehouseNonExistent)
			return domain.Employee{}, ErrWarehouseNonExistent
		default:
			logging.Log(ErrEmployeeNotSaved)
			return domain.Employee{}, ErrEmployeeNotSaved
		}
	}

	return employee, nil
}

func (service *service) Update(ctx context.Context, employee domain.Employee) (domain.Employee, error) {
	updatedEmployee, err := service.Get(ctx, employee.ID)

	if err != nil {
		logging.Log(err)
		return domain.Employee{}, err
	}

	if employee.FirstName != "" {
		updatedEmployee.FirstName = employee.FirstName
	}
	if employee.LastName != "" {
		updatedEmployee.LastName = employee.LastName
	}
	if employee.WarehouseID > 0 {
		updatedEmployee.WarehouseID = employee.WarehouseID
	}

	err = service.repository.Update(ctx, updatedEmployee)

	if err != nil {
		logging.Log(ErrEmployeeNotUpdated)
		return domain.Employee{}, ErrEmployeeNotUpdated
	}

	return updatedEmployee, nil
}

func (service *service) Delete(ctx context.Context, id int) error {
	err := service.repository.Delete(ctx, id)

	if err != nil {
		if err.Error() == ErrEmployeeNotFound.Error() {
			logging.Log(ErrEmployeeNotFound)
			return ErrEmployeeNotFound
		} else {
			logging.Log(err)
			return err
		}
	}

	return nil
}
