package buyer

import (
	"context"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
)

// Errors
var (
	NotFound       = errors.New("buyer not found")
	OneRowAffected = errors.New("one row Affected")
)

type MockRepository struct {
	Data []domain.Buyer
	Err  error
}

// GetAll returns a list od buyers or weird SQL errors
func (m *MockRepository) GetAll(ctx context.Context) ([]domain.Buyer, error) {
	if len(m.Data) == 0 {
		return nil, NotFound
	}
	return m.Data, nil
}

// Get returns a buyer or NotFound
func (m *MockRepository) Get(ctx context.Context, id int) (domain.Buyer, error) {
	result := []domain.Buyer{}

	if m.Err != nil {
		return domain.Buyer{}, m.Err
	}

	if len(m.Data) == 0 {
		return domain.Buyer{}, NotFound
	}

	for _, buyer := range m.Data {
		if buyer.ID == id {
			result = append(result, buyer)
		}
	}

	if len(result) == 0 {
		return domain.Buyer{}, NotFound
	}
	return result[0], nil
}

// Exists returns only boolean values
func (m *MockRepository) Exists(ctx context.Context, cardNumberID string) bool {
	flag := false
	if len(m.Data) > 0 {
		for _, buyer := range m.Data {
			if buyer.CardNumberID == cardNumberID {
				flag = true
			}
		}
	}
	return flag
}

// Save returns buyer Id or weird SQL errors
func (m *MockRepository) Save(ctx context.Context, b domain.Buyer) (int, error) {
	if m.Err != nil {
		return 0, m.Err
	}
	return b.ID, nil
}

// Update returns only weird SQL errors or nil
func (m *MockRepository) Update(ctx context.Context, b domain.Buyer) error {
	if m.Err != nil {
		return m.Err
	}
	return nil
}

// Delete returns only weird SQL errors or nil
func (m *MockRepository) Delete(ctx context.Context, id int) error {
	if m.Err != nil {
		return m.Err
	}
	return nil
}
