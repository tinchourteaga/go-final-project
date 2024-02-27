package purchaseorders

import (
	"context"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
)

type MockRepository struct {
	Data               []domain.Purchase_orders
	DataOrdersByBuyers []domain.Purchase_orders_buyer
	Err                error
	ErrExist           error
}

func (m *MockRepository) Exists(ctx context.Context, orderNumber string) bool {
	return m.ErrExist != nil
}

func (m *MockRepository) SaveOrder(ctx context.Context, p domain.Purchase_orders) (int, error) {
	if m.Err != nil {
		return 0, m.Err
	}
	return p.ID, nil
}

func (m *MockRepository) GetByBuyerId(ctx context.Context, buyerId int) ([]domain.Purchase_orders_buyer, error) {
	var result []domain.Purchase_orders_buyer
	if m.Err != nil {
		return nil, m.Err
	}

	if m.DataOrdersByBuyers != nil && buyerId != 0 {
		result = m.DataOrdersByBuyers
	}
	return result, nil
}

func (m *MockRepository) GetAllByBuyer(ctx context.Context) ([]domain.Purchase_orders_buyer, error) {
	var result []domain.Purchase_orders_buyer
	if m.Err != nil {
		return nil, m.Err
	}

	if m.DataOrdersByBuyers != nil {
		result = m.DataOrdersByBuyers
	}
	return result, nil
}
