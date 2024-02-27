package employee

import (
	"context"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
)

type MockRepository struct {
	DataMock  []domain.Employee
	MockError error
}

func (mockRepository *MockRepository) GetAll(ctx context.Context) ([]domain.Employee, error) {
	if mockRepository.MockError != nil {
		return []domain.Employee{}, mockRepository.MockError
	}
	return mockRepository.DataMock, nil
}

func (mockRepository *MockRepository) Get(ctx context.Context, id int) (domain.Employee, error) {
	for _, employee := range mockRepository.DataMock {
		if employee.ID == id {
			return employee, nil
		}
	}
	return domain.Employee{}, mockRepository.MockError
}

func (mockRepository *MockRepository) Exists(ctx context.Context, cardNumberID string) bool {
	for _, employee := range mockRepository.DataMock {
		if employee.CardNumberID == cardNumberID {
			return true
		}
	}
	return false
}

func (mockRepository *MockRepository) Save(ctx context.Context, newEmployee domain.Employee) (int, error) {
	if mockRepository.Exists(ctx, newEmployee.CardNumberID) {
		return 0, mockRepository.MockError
	}

	lastID := getLastID(mockRepository)
	newEmployee.ID = lastID
	mockRepository.DataMock = append(mockRepository.DataMock, newEmployee)
	return newEmployee.ID, nil
}

func (mockRepository *MockRepository) Update(ctx context.Context, updatedEmployee domain.Employee) error {
	for i, employee := range mockRepository.DataMock {
		if employee.ID == updatedEmployee.ID {
			mockRepository.DataMock[i] = updatedEmployee
			return nil
		}
	}
	return mockRepository.MockError
}

func (mockRepository *MockRepository) Delete(ctx context.Context, id int) error {
	for i, employee := range mockRepository.DataMock {
		if employee.ID == id {
			mockRepository.DataMock = append(mockRepository.DataMock[:i], mockRepository.DataMock[i+1:]...)
			return nil
		}
	}
	return mockRepository.MockError
}

func getLastID(mockRepository *MockRepository) int {
	lastID := 0

	if len(mockRepository.DataMock) > 0 {
		lastID = mockRepository.DataMock[len(mockRepository.DataMock)-1].ID
	}

	lastID++
	return lastID
}
