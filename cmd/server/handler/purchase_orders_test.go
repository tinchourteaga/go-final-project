package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	purchaseorders "github.com/extmatperez/meli_bootcamp_go_w6-2/internal/purchase_orders"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func createServerPurchaseOrders(mockRepository purchaseorders.MockRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	service := purchaseorders.NewService(&mockRepository)
	handler := NewPurchaseOrders(service)

	r := gin.Default()

	pr := r.Group("/api/v1/purchase_orders")
	pr.POST("/", handler.CreateOrder())

	pr2 := r.Group("/api/v1/reportPurchaseOrder")
	pr2.GET("", handler.GetAllOrdersByBuyers())

	return r
}

func createRequestTestPurchaseOrders(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")

	return req, httptest.NewRecorder()
}

// TestCreateOrderSuccess passes when return Purchase_orders created (status code 201)
func TestCreateOrderSuccess(t *testing.T) {
	//Arrange and Act
	repo := purchaseorders.MockRepository{}
	r := createServerPurchaseOrders(repo)
	req, recorder := createRequestTestPurchaseOrders(http.MethodPost, "/api/v1/purchase_orders/", `{"id":1, "order_number":"002", "order_date":"2022-10-10", "tracking_code":"asd233501", "buyer_id": 1, "product_record_id":1, "order_status_id":1}`)
	r.ServeHTTP(recorder, req)

	//asserts
	assert.Equal(t, http.StatusCreated, recorder.Code)
}

// TestCreateOrderFailInternalServerError passes when return an error for internal server error (status code 500)
func TestCreateOrderFailInternalServerError(t *testing.T) {

	//Arrange and Act
	repo := purchaseorders.MockRepository{
		Err: purchaseorders.ErrInternal,
	}
	r := createServerPurchaseOrders(repo)
	req, recorder := createRequestTestPurchaseOrders(http.MethodPost, "/api/v1/purchase_orders/", `{"id":1, "order_number":"003", "order_date":"2022-10-10", "tracking_code":"asd233501", "buyer_id": 1, "product_record_id":1, "order_status_id":1}`)
	r.ServeHTTP(recorder, req)

	//asserts
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

// TestCreateOrderFailInternalServerError passes when return an error for FK (status code 1452)
func TestCreateOrderFailErrorFKConstraint(t *testing.T) {

	//Arrange and Act
	repo := purchaseorders.MockRepository{
		Err: purchaseorders.ErrFKConstraint,
	}
	r := createServerPurchaseOrders(repo)
	req, recorder := createRequestTestPurchaseOrders(http.MethodPost, "/api/v1/purchase_orders/", `{"id":1, "order_number":"004", "order_date":"2022-10-10", "tracking_code":"asd233501", "buyer_id": 1, "product_record_id":1, "order_status_id":5}`)
	r.ServeHTTP(recorder, req)

	//asserts
	assert.Equal(t, http.StatusConflict, recorder.Code)
}

// TestCreateOrderFailErrAlreadyExists passes when return an error (status code 409)
func TestCreateOrderFailErrAlreadyExists(t *testing.T) {

	//Arrange and Act
	repo := purchaseorders.MockRepository{
		Err: purchaseorders.ErrAlreadyExists,
	}
	r := createServerPurchaseOrders(repo)
	req, recorder := createRequestTestPurchaseOrders(http.MethodPost, "/api/v1/purchase_orders/", `{"id":1, "order_number":"004", "order_date":"2022-10-10", "tracking_code":"asd233501", "buyer_id": 1, "product_record_id":1, "order_status_id":5}`)
	r.ServeHTTP(recorder, req)

	//asserts
	assert.Equal(t, http.StatusConflict, recorder.Code)
}

// TestCreateOrderFailErrDataLong passes when return an error (status code 422)
func TestCreateOrderFailErrDataLong(t *testing.T) {

	//Arrange and Act
	repo := purchaseorders.MockRepository{
		Err: purchaseorders.ErrDataLong,
	}
	r := createServerPurchaseOrders(repo)
	req, recorder := createRequestTestPurchaseOrders(http.MethodPost, "/api/v1/purchase_orders/", `{"id":1, "order_number":"004", "order_date":"2022-10-10", "tracking_code":"asd233501", "buyer_id": 1, "product_record_id":1, "order_status_id":5}`)
	r.ServeHTTP(recorder, req)

	//asserts
	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
}

// TestCreateOrderFailStatusUnprocessableEntity passes when return an error for bad request (status code 422)
func TestCreateOrderFailStatusUnprocessableEntity(t *testing.T) {

	//Arrange and Act
	repo := purchaseorders.MockRepository{}
	r := createServerPurchaseOrders(repo)
	req, recorder := createRequestTestPurchaseOrders(http.MethodPost, "/api/v1/purchase_orders/", `{"id":1, "order_number":"005", "order_date":"2022-10-10", "tracking_code":"asd233501", "buyer_id": 1, "product_record_id":1 "order_status_id":1}`)
	r.ServeHTTP(recorder, req)

	//asserts
	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
}

// TestGetAllOrdersByBuyers passes when return a list of Purchase_orders (status code 200)
func TestGetAllOrdersByBuyers(t *testing.T) {
	//Arrange
	ListPurchaseOrders := []domain.Purchase_orders{
		{
			ID:              1,
			OrderNumber:     "002",
			OrderDate:       "2022-10-10",
			TrackingCode:    "asd233501",
			BuyerId:         1,
			ProductRecordId: 1,
			OrderStatusId:   1,
		},
	}
	//Act
	repo := purchaseorders.MockRepository{
		Data: ListPurchaseOrders,
	}
	r := createServerPurchaseOrders(repo)
	req, recorder := createRequestTestPurchaseOrders(http.MethodGet, "/api/v1/reportPurchaseOrder", "")
	r.ServeHTTP(recorder, req)

	//arrange
	assert.Equal(t, http.StatusOK, recorder.Code)
}

// TestGetAllOrdersByBuyersSuccessById passes when return Purchase_orders (status code 200)
func TestGetAllOrdersByBuyersSuccessById(t *testing.T) {
	//Arrange
	ListPurchaseOrders := []domain.Purchase_orders{
		{
			ID:              1,
			OrderNumber:     "002",
			OrderDate:       "2022-10-10",
			TrackingCode:    "asd233501",
			BuyerId:         1,
			ProductRecordId: 1,
			OrderStatusId:   1,
		},
	}
	//Act
	repo := purchaseorders.MockRepository{
		Data: ListPurchaseOrders,
	}
	r := createServerPurchaseOrders(repo)
	req, recorder := createRequestTestPurchaseOrders(http.MethodGet, "/api/v1/reportPurchaseOrder?id=1", "")
	r.ServeHTTP(recorder, req)

	//arrange
	assert.Equal(t, http.StatusOK, recorder.Code)
}

// TestGetAllOrdersByBuyersFailInvalidId passes when an error for invalid Id (status code 400)
func TestGetAllOrdersByBuyersFailNullId(t *testing.T) {
	//Arrange
	ListPurchaseOrders := []domain.Purchase_orders{
		{
			ID:              1,
			OrderNumber:     "002",
			OrderDate:       "2022-10-10",
			TrackingCode:    "asd233501",
			BuyerId:         1,
			ProductRecordId: 1,
			OrderStatusId:   1,
		},
	}
	//Act
	repo := purchaseorders.MockRepository{
		Data: ListPurchaseOrders,
		//Err:  purchaseorders.ErrNotFound,
	}
	r := createServerPurchaseOrders(repo)
	req, recorder := createRequestTestPurchaseOrders(http.MethodGet, "/api/v1/reportPurchaseOrder?id=aaa", "")
	r.ServeHTTP(recorder, req)

	//arrange
	fmt.Println("Debugger Agus 9")
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

// TestGetAllOrdersByBuyersFailStatusInternalServerError passes when return an error for Internal Server Error (status code 500)
func TestGetAllOrdersByBuyersFailStatusInternalServerError(t *testing.T) {
	//Arrange
	ListPurchaseOrders := []domain.Purchase_orders{
		{
			ID:              1,
			OrderNumber:     "002",
			OrderDate:       "2022-10-10",
			TrackingCode:    "asd233501",
			BuyerId:         1,
			ProductRecordId: 1,
			OrderStatusId:   1,
		},
	}
	//Act
	repo := purchaseorders.MockRepository{
		Data: ListPurchaseOrders,
		Err:  purchaseorders.ErrInternal,
	}
	r := createServerPurchaseOrders(repo)
	req, recorder := createRequestTestPurchaseOrders(http.MethodGet, "/api/v1/reportPurchaseOrder", "")
	r.ServeHTTP(recorder, req)

	//arrange
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}
