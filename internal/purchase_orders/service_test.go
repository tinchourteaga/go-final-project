package purchaseorders

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	logging.InitLog(nil)
}

// TestSaveOrdersSuccess passes when return Purchase_orders and not errors
func TestSaveOrdersSuccess(t *testing.T) {

	expectedOrder := domain.Purchase_orders{
		ID:              1,
		OrderNumber:     "001",
		OrderDate:       "2022-05-05",
		TrackingCode:    "sdaksdjf664387",
		BuyerId:         1,
		ProductRecordId: 5,
		OrderStatusId:   200,
	}
	result := []domain.Purchase_orders{}
	result = append(result, expectedOrder)
	mockRepo := MockRepository{
		Data: result,
	}
	s := NewService(&mockRepo)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	p, err := s.SaveOrder(ctx, expectedOrder)

	assert.NoError(t, err)
	assert.NotEmpty(t, p)
}

// TestSaveOrdersFailureSave passes when return an error because OrderNumber already exists
func TestSaveOrdersFailureSave(t *testing.T) {

	expectedOrder := domain.Purchase_orders{
		ID:              1,
		OrderNumber:     "001",
		OrderDate:       "2022-05-05",
		TrackingCode:    "sdaksdjf664387",
		BuyerId:         1,
		ProductRecordId: 5,
		OrderStatusId:   200,
	}
	mockRepo := MockRepository{
		ErrExist: ErrAlreadyExists,
	}
	s := NewService(&mockRepo)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	p, err := s.SaveOrder(ctx, expectedOrder)
	assert.EqualError(t, ErrAlreadyExists, err.Error())
	assert.Empty(t, p)
}

// TestSaveOrdersFailureSaveInternalError passes when return an error (Status 500)
func TestSaveOrdersFailureSaveInternalError(t *testing.T) {

	expectedOrder := domain.Purchase_orders{
		ID:              1,
		OrderNumber:     "001",
		OrderDate:       "2022-05-05",
		TrackingCode:    "sdaksdjf664387",
		BuyerId:         1,
		ProductRecordId: 5,
		OrderStatusId:   200,
	}
	mockRepo := MockRepository{
		Err: ErrInternal,
	}
	s := NewService(&mockRepo)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	p, err := s.SaveOrder(ctx, expectedOrder)
	assert.EqualError(t, ErrInternal, err.Error())
	assert.Empty(t, p)
}

// TestExistsSuccessOrders passes when return true
func TestExistsSuccessOrders(t *testing.T) {
	orderId := "001"
	mockRepo := MockRepository{
		ErrExist: errors.New("order_number already exists"),
	}
	s := NewService(&mockRepo)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	err := s.Exists(ctx, orderId)

	assert.Equal(t, true, err)
}

// TestExistsFailOrders passes when return false
func TestExistsFailOrders(t *testing.T) {
	orderId := "001"
	mockRepo := MockRepository{}
	s := NewService(&mockRepo)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	err := s.Exists(ctx, orderId)

	assert.Equal(t, false, err)
}

// TestGetAllByBuyerSuccessWithId passes when return false
func TestGetAllByBuyerSuccessWithId(t *testing.T) {
	expectedOrder := domain.Purchase_orders_buyer{
		ID:           1,
		CardNumberId: "001",
		FirstName:    "Comprador 1",
		LastName:     "Vendedor 1",
		OrdersCount:  2,
	}
	orderId := 1
	result := []domain.Purchase_orders_buyer{}
	result = append(result, expectedOrder)
	mockRepo := MockRepository{
		DataOrdersByBuyers: result,
	}
	s := NewService(&mockRepo)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	p, err := s.GetAllByBuyer(ctx, orderId)

	assert.Nil(t, err)
	assert.NotEmpty(t, p)
	assert.Equal(t, result, p)
}

func TestGetAllByBuyerSuccessNotId(t *testing.T) {
	expectedOrder := domain.Purchase_orders_buyer{
		ID:           1,
		CardNumberId: "001",
		FirstName:    "Comprador 1",
		LastName:     "Vendedor 1",
		OrdersCount:  2,
	}
	orderId := 0
	result := []domain.Purchase_orders_buyer{}
	result = append(result, expectedOrder)
	mockRepo := MockRepository{
		DataOrdersByBuyers: result,
	}
	s := NewService(&mockRepo)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	p, err := s.GetAllByBuyer(ctx, orderId)

	assert.Nil(t, err)
	assert.NotEmpty(t, p)
	assert.Equal(t, result, p)
}

func TestGetAllByBuyerFailWithId(t *testing.T) {
	orderId := 1
	mockRepo := MockRepository{
		Err: ErrInternal,
	}
	s := NewService(&mockRepo)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	p, err := s.GetAllByBuyer(ctx, orderId)
	assert.EqualError(t, ErrInternal, err.Error())
	assert.Empty(t, p)
}

func TestGetAllByBuyerFailNotId(t *testing.T) {
	orderId := 0
	mockRepo := MockRepository{
		Err: ErrInternal,
	}
	s := NewService(&mockRepo)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	p, err := s.GetAllByBuyer(ctx, orderId)

	assert.EqualError(t, ErrInternal, err.Error())
	assert.Empty(t, p)
}
