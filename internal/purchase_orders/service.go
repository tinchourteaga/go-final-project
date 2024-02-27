package purchaseorders

import (
	"context"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
)

type Service interface {
	SaveOrder(ctx context.Context, p domain.Purchase_orders) (domain.Purchase_orders, error)
	Exists(ctx context.Context, orderID string) bool
	GetAllByBuyer(ctx context.Context, orderId int) ([]domain.Purchase_orders_buyer, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

// Save returns the created a purchase_order if successful, or a error if it failed
// if the purchase_order is already exist, a error is returned
func (s *service) SaveOrder(ctx context.Context, p domain.Purchase_orders) (domain.Purchase_orders, error) {
	if s.repository.Exists(ctx, p.OrderNumber) {
		return domain.Purchase_orders{}, ErrAlreadyExists
	}

	order := domain.Purchase_orders{
		OrderNumber:     p.OrderNumber,
		OrderDate:       p.OrderDate,
		TrackingCode:    p.TrackingCode,
		BuyerId:         p.BuyerId,
		ProductRecordId: p.ProductRecordId,
		OrderStatusId:   p.OrderStatusId,
	}

	id, err := s.repository.SaveOrder(ctx, order)
	if err != nil {
		logging.Log(err)
		return domain.Purchase_orders{}, err
	}

	order.ID = id
	return order, nil
}

// Exists returns true if the Purchase_order is already exist, or false if it doesn´t
func (s *service) Exists(ctx context.Context, orderID string) bool {
	return s.repository.Exists(ctx, orderID)
}

// GetAll returns a List of Purchase_order if successful, or a error if it failed
func (s *service) GetAllByBuyer(ctx context.Context, orderId int) ([]domain.Purchase_orders_buyer, error) {

	if orderId != 0 {
		return s.repository.GetByBuyerId(ctx, orderId)
	}
	return s.repository.GetAllByBuyer(ctx)
}
