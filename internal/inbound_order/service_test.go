package inbound_order

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

var (
	testEmployee = domain.Employee{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 1}
	ctx          context.Context
)

func TestGetAllEmployeesInboundOrders_Ok(t *testing.T) {
	expectedResult := 1
	db := []domain.EmployeeWithInboundOrders{
		{Employee: testEmployee, InboundOrders: 1},
	}

	mockRepository := MockRepository{
		DataMockEmployeesWithIO: db,
	}

	service := NewService(&mockRepository)

	result, err := service.GetAllEmployeesInboundOrders(ctx)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, len(result))
}

func TestGetAllEmployeesInboundOrders_Fail(t *testing.T) {
	expectedErr := errors.New("error")

	mockRepository := MockRepository{
		MockErrorGetAll: errors.New("error"),
	}

	service := NewService(&mockRepository)

	result, err := service.GetAllEmployeesInboundOrders(ctx)

	assert.Error(t, err)
	assert.EqualError(t, err, expectedErr.Error())
	assert.Empty(t, result)
}

func TestGetEmployeeInboundOrders_Ok(t *testing.T) {
	expectedResult := domain.EmployeeWithInboundOrders{Employee: testEmployee, InboundOrders: 1}
	id := 1
	db := []domain.EmployeeWithInboundOrders{
		{Employee: testEmployee, InboundOrders: 1},
	}

	mockRepository := MockRepository{
		DataMockEmployeesWithIO: db,
	}

	service := NewService(&mockRepository)

	result, err := service.GetEmployeeInboundOrders(ctx, id)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestGetEmployeeInboundOrders_NotFound(t *testing.T) {
	expectedErr := ErrEmployeeWithInboundOrdersNotFound
	id := 1

	mockRepository := MockRepository{
		MockErrorGet: errors.New("error"),
	}

	service := NewService(&mockRepository)

	result, err := service.GetEmployeeInboundOrders(ctx, id)

	assert.Error(t, err)
	assert.EqualError(t, err, expectedErr.Error())
	assert.Empty(t, result)
}

func TestSaveInboundOrder_Ok(t *testing.T) {
	newInboundOrder := domain.InboundOrder{OrderDate: "01/01/2022", OrderNumber: "Test#1", EmployeeID: 1, ProductBatchID: 1, WarehouseID: 1}
	expectedInboundOrder := domain.InboundOrder{ID: 1, OrderDate: "01/01/2022", OrderNumber: "Test#1", EmployeeID: 1, ProductBatchID: 1, WarehouseID: 1}
	expectedResult := []domain.InboundOrder{
		expectedInboundOrder,
	}

	db := []domain.InboundOrder{}

	mockRepository := MockRepository{
		DataMockInboundOrders: db,
		ExpectedID:            1,
	}

	service := NewService(&mockRepository)

	result, err := service.Save(ctx, newInboundOrder)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, mockRepository.DataMockInboundOrders)
	assert.Equal(t, expectedInboundOrder, result)
}

func TestSaveInboundOrder_Fail(t *testing.T) {
	newInboundOrder := domain.InboundOrder{OrderDate: "01/01/2022", OrderNumber: "Test#1", EmployeeID: 1, ProductBatchID: 1, WarehouseID: 1}
	expectedErr := ErrInboundOrderNotSaved

	db := []domain.InboundOrder{}

	mockRepository := MockRepository{
		DataMockInboundOrders: db,
		MockErrorSave:         ErrInboundOrderNotSaved,
	}

	service := NewService(&mockRepository)

	result, err := service.Save(ctx, newInboundOrder)

	assert.Error(t, err)
	assert.EqualError(t, err, expectedErr.Error())
	assert.Empty(t, result)
}

func TestSaveInboundOrder_AlreadyExists(t *testing.T) {
	newInboundOrder := domain.InboundOrder{OrderDate: "01/01/2022", OrderNumber: "Test#1", EmployeeID: 1, ProductBatchID: 1, WarehouseID: 1}
	expectedErr := ErrInboundOrderAlreadyExists

	db := []domain.InboundOrder{
		{OrderDate: "10/10/2022", OrderNumber: "Test#1", EmployeeID: 2, ProductBatchID: 2, WarehouseID: 2},
	}

	mockRepository := MockRepository{
		DataMockInboundOrders:  db,
		MockErrorAlreadyExists: ErrInboundOrderAlreadyExists,
	}

	service := NewService(&mockRepository)

	result, err := service.Save(ctx, newInboundOrder)

	assert.Error(t, err)
	assert.EqualError(t, err, expectedErr.Error())
	assert.Empty(t, result)
}

func TestSaveInboundOrder_EmptyOrderNumber(t *testing.T) {
	newInboundOrder := domain.InboundOrder{OrderDate: "01/01/2022", OrderNumber: "", EmployeeID: 1, ProductBatchID: 1, WarehouseID: 1}
	expectedErr := ErrEmptyOrderNumber

	db := []domain.InboundOrder{}

	mockRepository := MockRepository{
		DataMockInboundOrders:     db,
		MockErrorEmptyOrderNumber: ErrEmptyOrderNumber,
	}

	service := NewService(&mockRepository)

	result, err := service.Save(ctx, newInboundOrder)

	assert.Error(t, err)
	assert.EqualError(t, err, expectedErr.Error())
	assert.Empty(t, result)
}

func TestSaveInboundOrder_FailEmployeeFK(t *testing.T) {
	newInboundOrder := domain.InboundOrder{OrderDate: "01/01/2022", OrderNumber: "Test#1", EmployeeID: 1, ProductBatchID: 1, WarehouseID: 1}
	expectedErr := ErrEmployeeNonExistent

	db := []domain.InboundOrder{}

	mockRepository := MockRepository{
		DataMockInboundOrders: db,
		MockErrorEmployeeFK:   ErrEmployeeNonExistent,
	}

	service := NewService(&mockRepository)

	result, err := service.Save(ctx, newInboundOrder)

	assert.Error(t, err)
	assert.EqualError(t, err, expectedErr.Error())
	assert.Empty(t, result)
}

func TestSaveInboundOrder_FailWarehouseFK(t *testing.T) {
	newInboundOrder := domain.InboundOrder{OrderDate: "01/01/2022", OrderNumber: "Test#1", EmployeeID: 1, ProductBatchID: 1, WarehouseID: 1}
	expectedErr := ErrWarehouseNonExistent

	db := []domain.InboundOrder{}

	mockRepository := MockRepository{
		DataMockInboundOrders: db,
		MockErrorEmployeeFK:   ErrWarehouseNonExistent,
	}

	service := NewService(&mockRepository)

	result, err := service.Save(ctx, newInboundOrder)

	assert.Error(t, err)
	assert.EqualError(t, err, expectedErr.Error())
	assert.Empty(t, result)
}

func TestSaveInboundOrder_FailProductBatchFK(t *testing.T) {
	newInboundOrder := domain.InboundOrder{OrderDate: "01/01/2022", OrderNumber: "Test#1", EmployeeID: 1, ProductBatchID: 1, WarehouseID: 1}
	expectedErr := ErrProductBatchNonExistent

	db := []domain.InboundOrder{}

	mockRepository := MockRepository{
		DataMockInboundOrders: db,
		MockErrorEmployeeFK:   ErrProductBatchNonExistent,
	}

	service := NewService(&mockRepository)

	result, err := service.Save(ctx, newInboundOrder)

	assert.Error(t, err)
	assert.EqualError(t, err, expectedErr.Error())
	assert.Empty(t, result)
}
