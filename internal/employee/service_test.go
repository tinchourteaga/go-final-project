package employee

import (
	"context"
	"errors"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/stretchr/testify/assert"
)

func init() {
	logging.InitLog(nil)
}

/* =============== GET ALL =============== */
func TestGetAll(t *testing.T) {
	var ctx context.Context
	expectedResult := 2
	db := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 3},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	mockRepository := MockRepository{
		DataMock:  db,
		MockError: nil,
	}

	service := NewService(&mockRepository)

	result, err := service.GetAll(ctx)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, len(result))
}

func TestGetAllFail(t *testing.T) {
	var ctx context.Context
	expectedErr := errors.New("error retrieving employees")
	db := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 3},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	mockRepository := MockRepository{
		DataMock:  db,
		MockError: errors.New("error retrieving employees"),
	}

	service := NewService(&mockRepository)

	result, err := service.GetAll(ctx)

	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, result)
}

/* =============== GET =============== */
func TestGet(t *testing.T) {
	var ctx context.Context
	expectedResult := domain.Employee{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7}
	id := 2
	db := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 3},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	mockRepository := MockRepository{
		DataMock:  db,
		MockError: nil,
	}

	service := NewService(&mockRepository)

	result, err := service.Get(ctx, id)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestGetFail(t *testing.T) {
	var ctx context.Context
	expectedErr := ErrEmployeeNotFound
	id := 5
	db := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 3},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	mockRepository := MockRepository{
		DataMock:  db,
		MockError: errors.New(""), // Service error overrides this msg
	}

	service := NewService(&mockRepository)

	result, err := service.Get(ctx, id)

	assert.EqualError(t, err, expectedErr.Error())
	assert.Empty(t, result) // Requirement indicates it should return nil, however, we will return an empty Employee
}

/* =============== SAVE =============== */
func TestSave(t *testing.T) {
	var ctx context.Context
	newEmployee := domain.Employee{CardNumberID: "111111", FirstName: "Martin", LastName: "Urteaga", WarehouseID: 10}
	expectedEmployee := domain.Employee{ID: 3, CardNumberID: "111111", FirstName: "Martin", LastName: "Urteaga", WarehouseID: 10}
	expectedResult := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 3},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
		{ID: 3, CardNumberID: "111111", FirstName: "Martin", LastName: "Urteaga", WarehouseID: 10},
	}

	db := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 3},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	mockRepository := MockRepository{
		DataMock:  db,
		MockError: nil,
	}

	service := NewService(&mockRepository)

	result, err := service.Save(ctx, newEmployee)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, mockRepository.DataMock)
	assert.Equal(t, expectedEmployee, result)
}

func TestSaveFail(t *testing.T) {
	var ctx context.Context
	newEmployee := domain.Employee{CardNumberID: "123456", FirstName: "Martin", LastName: "Urteaga", WarehouseID: 10}
	expectedErr := ErrEmployeeAlreadyExists

	db := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 3},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	mockRepository := MockRepository{
		DataMock:  db,
		MockError: errors.New(""), // Service error overrides this msg
	}

	service := NewService(&mockRepository)

	result, err := service.Save(ctx, newEmployee)

	assert.EqualError(t, err, expectedErr.Error())
	assert.Empty(t, result)
}

/* =============== UPDATE =============== */
func TestUpdate(t *testing.T) {
	var ctx context.Context
	employeeToUpdate := domain.Employee{ID: 1, FirstName: "Martin", LastName: "Urteaga", WarehouseID: 6}
	updatedEmployee := domain.Employee{ID: 1, CardNumberID: "123456", FirstName: "Martin", LastName: "Urteaga", WarehouseID: 6}

	expectedResult := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "Martin", LastName: "Urteaga", WarehouseID: 6},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	db := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 3},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	mockRepository := MockRepository{
		DataMock:  db,
		MockError: nil,
	}

	service := NewService(&mockRepository)

	result, err := service.Update(ctx, employeeToUpdate)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, mockRepository.DataMock)
	assert.Equal(t, updatedEmployee, result)
}

func TestUpdateFail(t *testing.T) {
	var ctx context.Context
	employeeToUpdate := domain.Employee{ID: 5, FirstName: "Martin", LastName: "Urteaga", WarehouseID: 6}
	expectedErr := ErrEmployeeNotFound

	db := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 3},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	mockRepository := MockRepository{
		DataMock:  db,
		MockError: errors.New(""), // Service error overrides this msg
	}

	service := NewService(&mockRepository)

	result, err := service.Update(ctx, employeeToUpdate)

	assert.EqualError(t, err, expectedErr.Error())
	assert.Empty(t, result) // Requirement indicates it should return nil, however, we will return an empty Employee
}

/* =============== DELETE =============== */
func TestDelete(t *testing.T) {
	var ctx context.Context
	id := 1

	expectedResult := []domain.Employee{
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	db := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 3},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	mockRepository := MockRepository{
		DataMock:  db,
		MockError: nil,
	}

	service := NewService(&mockRepository)

	err := service.Delete(ctx, id)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, mockRepository.DataMock)
}

func TestDeleteFail(t *testing.T) {
	var ctx context.Context
	id := 5
	expectedErr := ErrEmployeeNotFound

	db := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 3},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	mockRepository := MockRepository{
		DataMock:  db,
		MockError: ErrEmployeeNotFound,
	}

	service := NewService(&mockRepository)

	err := service.Delete(ctx, id)

	assert.EqualError(t, err, expectedErr.Error())
}
