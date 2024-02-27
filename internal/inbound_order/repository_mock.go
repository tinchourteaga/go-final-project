package inbound_order

import (
	"context"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
)

type MockRepository struct {
	DataMockEmployeesWithIO   []domain.EmployeeWithInboundOrders
	DataMockInboundOrders     []domain.InboundOrder
	ExpectedID                int
	MockErrorGetAll           error
	MockErrorGet              error
	MockErrorSave             error
	MockErrorEmptyOrderNumber error
	MockErrorEmployeeFK       error
	MockErrorWarehouseFK      error
	MockErrorProductBatchFK   error
	MockErrorAlreadyExists    error
}

func (mockRepository *MockRepository) GetAllEmployeesInboundOrders(ctx context.Context) ([]domain.EmployeeWithInboundOrders, error) {
	if mockRepository.MockErrorGetAll != nil {
		return []domain.EmployeeWithInboundOrders{}, mockRepository.MockErrorGetAll
	}

	return mockRepository.DataMockEmployeesWithIO, nil
}

func (mockRepository *MockRepository) GetEmployeeInboundOrders(ctx context.Context, id int) (domain.EmployeeWithInboundOrders, error) {
	if mockRepository.MockErrorGet != nil {
		return domain.EmployeeWithInboundOrders{}, mockRepository.MockErrorGet
	}

	if len(mockRepository.DataMockEmployeesWithIO) > 0 {
		return mockRepository.DataMockEmployeesWithIO[0], nil
	}

	return domain.EmployeeWithInboundOrders{}, mockRepository.MockErrorGet
}

func (mockRepository *MockRepository) Save(ctx context.Context, newInboundOrder domain.InboundOrder) (int, error) {
	if mockRepository.MockErrorSave != nil {
		return 0, mockRepository.MockErrorSave
	}

	if mockRepository.MockErrorEmptyOrderNumber != nil {
		return 0, mockRepository.MockErrorEmptyOrderNumber
	}

	if mockRepository.MockErrorEmployeeFK != nil {
		return 0, mockRepository.MockErrorEmployeeFK
	}

	if mockRepository.MockErrorWarehouseFK != nil {
		return 0, mockRepository.MockErrorWarehouseFK
	}

	if mockRepository.MockErrorProductBatchFK != nil {
		return 0, mockRepository.MockErrorProductBatchFK
	}

	if mockRepository.MockErrorAlreadyExists != nil {
		return 0, mockRepository.MockErrorAlreadyExists
	}

	newInboundOrder.ID = mockRepository.ExpectedID
	mockRepository.DataMockInboundOrders = append(mockRepository.DataMockInboundOrders, newInboundOrder)
	return newInboundOrder.ID, nil
}
