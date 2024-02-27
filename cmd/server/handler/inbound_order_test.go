package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/inbound_order"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	logging.InitLog(nil)
}

type successfulResponseInboundOrders struct {
	Data domain.InboundOrder `json:"data"`
}

func createServerInboundOrders(mockRepository inbound_order.MockRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	service := inbound_order.NewService(&mockRepository)
	inboundOrderHandler := NewInboundOrder(service)
	router := gin.Default()

	inboundOrderRoutesGroup := router.Group("/api/v1/inboundOrders")

	inboundOrderRoutesGroup.POST("/", inboundOrderHandler.Create())

	return router
}

func createRequestTestInboundOrders(method, url, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")

	return req, httptest.NewRecorder()
}

func TestSaveInboundOrder_Ok(t *testing.T) {
	var response successfulResponseInboundOrders
	expectedInboundOrder := domain.InboundOrder{ID: 1, OrderDate: "01/01/2022", OrderNumber: "Test#1", EmployeeID: 1, ProductBatchID: 1, WarehouseID: 1}

	db := []domain.InboundOrder{}

	mockRepository := inbound_order.MockRepository{
		DataMockInboundOrders: db,
		ExpectedID:            1,
	}

	router := createServerInboundOrders(mockRepository)
	req, recorder := createRequestTestInboundOrders(http.MethodPost, "/api/v1/inboundOrders/", `{"order_date": "01/01/2022", "order_number": "Test#1", "employee_id": 1, "product_batch_id": 1, "warehouse_id": 1}`)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, 201, recorder.Code)
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, expectedInboundOrder, response.Data)
}

func TestSaveInboundOrder_Fail(t *testing.T) {
	db := []domain.InboundOrder{
		{ID: 1, OrderDate: "01/01/2022", OrderNumber: "Test#1", EmployeeID: 1, ProductBatchID: 1, WarehouseID: 1},
	}

	mockRepository := inbound_order.MockRepository{
		DataMockInboundOrders: db,
		ExpectedID:            1,
	}

	router := createServerInboundOrders(mockRepository)
	req, recorder := createRequestTestInboundOrders(http.MethodPost, "/api/v1/inboundOrders/", `{"order_number": "Test#1", "employee_id": 1, "product_batch_id": 1, "warehouse_id": 1}`)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, 422, recorder.Code)
}

func TestSaveInboundOrder_AlreadyExists(t *testing.T) {
	db := []domain.InboundOrder{
		{ID: 1, OrderDate: "01/01/2022", OrderNumber: "Test#1", EmployeeID: 1, ProductBatchID: 1, WarehouseID: 1},
	}

	mockRepository := inbound_order.MockRepository{
		DataMockInboundOrders:  db,
		ExpectedID:             1,
		MockErrorAlreadyExists: inbound_order.ErrInboundOrderAlreadyExists,
	}

	router := createServerInboundOrders(mockRepository)
	req, recorder := createRequestTestInboundOrders(http.MethodPost, "/api/v1/inboundOrders/", `{"order_date": "01/01/2022", "order_number": "Test#1", "employee_id": 1, "product_batch_id": 1, "warehouse_id": 1}`)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, 409, recorder.Code)
}

func TestSaveInboundOrder_EmptyOrderNumber(t *testing.T) {
	db := []domain.InboundOrder{}

	mockRepository := inbound_order.MockRepository{
		DataMockInboundOrders:     db,
		MockErrorEmptyOrderNumber: inbound_order.ErrEmptyOrderNumber,
	}

	router := createServerInboundOrders(mockRepository)
	req, recorder := createRequestTestInboundOrders(http.MethodPost, "/api/v1/inboundOrders/", `{"order_date": "01/01/2022", "order_number": "", "employee_id": 1, "product_batch_id": 1, "warehouse_id": 1}`)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, 409, recorder.Code)
}

func TestSaveInboundOrder_FailEmployeeFK(t *testing.T) {
	db := []domain.InboundOrder{}

	mockRepository := inbound_order.MockRepository{
		DataMockInboundOrders: db,
		MockErrorEmployeeFK:   inbound_order.ErrEmployeeNonExistent,
	}

	router := createServerInboundOrders(mockRepository)
	req, recorder := createRequestTestInboundOrders(http.MethodPost, "/api/v1/inboundOrders/", `{"order_date": "01/01/2022", "order_number": "Test#1", "employee_id": 10, "product_batch_id": 1, "warehouse_id": 1}`)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, 404, recorder.Code)
}

func TestSaveInboundOrder_FailWarehouseFK(t *testing.T) {
	db := []domain.InboundOrder{}

	mockRepository := inbound_order.MockRepository{
		DataMockInboundOrders: db,
		MockErrorWarehouseFK:  inbound_order.ErrWarehouseNonExistent,
	}

	router := createServerInboundOrders(mockRepository)
	req, recorder := createRequestTestInboundOrders(http.MethodPost, "/api/v1/inboundOrders/", `{"order_date": "01/01/2022", "order_number": "Test#1", "employee_id": 1, "product_batch_id": 1, "warehouse_id": 10}`)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, 404, recorder.Code)
}

func TestSaveInboundOrder_FailProductBatchFK(t *testing.T) {
	db := []domain.InboundOrder{}

	mockRepository := inbound_order.MockRepository{
		DataMockInboundOrders:   db,
		MockErrorProductBatchFK: inbound_order.ErrProductBatchNonExistent,
	}

	router := createServerInboundOrders(mockRepository)
	req, recorder := createRequestTestInboundOrders(http.MethodPost, "/api/v1/inboundOrders/", `{"order_date": "01/01/2022", "order_number": "Test#1", "employee_id": 1, "product_batch_id": 10, "warehouse_id": 1}`)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, 404, recorder.Code)
}

/* ============== Employee with inbound orders ==================== */
var (
	testEmployee = domain.Employee{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 1}
)

type successfulResponseEmployeeWithIO struct {
	Data domain.EmployeeWithInboundOrders `json:"data"`
}

type successfulSliceResponseEmployeeWithIO struct {
	Data []domain.EmployeeWithInboundOrders `json:"data"`
}

func createServerEmployeeWithIO(mockRepository inbound_order.MockRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	service := inbound_order.NewService(&mockRepository)
	inboundOrderHandler := NewInboundOrder(service)
	router := gin.Default()

	employeesRoutesGroup := router.Group("/api/v1/employees")
	employeesRoutesGroup.GET("/reportInboundOrders", inboundOrderHandler.GetAllEmployeesInboundOrders())
	employeesRoutesGroup.GET("/reportInboundOrders/:id", inboundOrderHandler.GetEmployeeInboundOrders())

	return router
}

func createRequestTestEmployeeWithIO(method, url, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")

	return req, httptest.NewRecorder()
}

func TestGetAllEmployeesInboundOrders_Ok(t *testing.T) {
	var sliceResponse successfulSliceResponseEmployeeWithIO
	db := []domain.EmployeeWithInboundOrders{
		{Employee: testEmployee, InboundOrders: 1},
	}

	mockRepository := inbound_order.MockRepository{
		DataMockEmployeesWithIO: db,
	}

	router := createServerEmployeeWithIO(mockRepository)
	req, recorder := createRequestTestEmployeeWithIO(http.MethodGet, "/api/v1/employees/reportInboundOrders", "")
	router.ServeHTTP(recorder, req)

	assert.Equal(t, 200, recorder.Code)
	err := json.Unmarshal(recorder.Body.Bytes(), &sliceResponse)
	assert.Nil(t, err)
	assert.Equal(t, mockRepository.DataMockEmployeesWithIO, sliceResponse.Data)
}

func TestGetAllEmployeesInboundOrders_Fail(t *testing.T) {
	db := []domain.EmployeeWithInboundOrders{
		{Employee: testEmployee, InboundOrders: 1},
	}

	mockRepository := inbound_order.MockRepository{
		DataMockEmployeesWithIO: db,
		MockErrorGetAll:         errors.New("error"),
	}

	router := createServerEmployeeWithIO(mockRepository)
	req, recorder := createRequestTestEmployeeWithIO(http.MethodGet, "/api/v1/employees/reportInboundOrders", "")
	router.ServeHTTP(recorder, req)

	assert.Equal(t, 500, recorder.Code)
}

func TestGetEmployeeInboundOrders_Ok(t *testing.T) {
	var response successfulResponseEmployeeWithIO
	db := []domain.EmployeeWithInboundOrders{
		{Employee: testEmployee, InboundOrders: 1},
	}

	mockRepository := inbound_order.MockRepository{
		DataMockEmployeesWithIO: db,
	}

	router := createServerEmployeeWithIO(mockRepository)
	req, recorder := createRequestTestEmployeeWithIO(http.MethodGet, "/api/v1/employees/reportInboundOrders/1", "")
	router.ServeHTTP(recorder, req)

	assert.Equal(t, 200, recorder.Code)
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, db[0], response.Data)
}

func TestGetEmployeeInboundOrders_NotFound(t *testing.T) {
	db := []domain.EmployeeWithInboundOrders{
		{Employee: testEmployee, InboundOrders: 1},
	}

	mockRepository := inbound_order.MockRepository{
		DataMockEmployeesWithIO: db,
		MockErrorGet:            inbound_order.ErrEmployeeWithInboundOrdersNotFound,
	}

	router := createServerEmployeeWithIO(mockRepository)
	req, recorder := createRequestTestEmployeeWithIO(http.MethodGet, "/api/v1/employees/reportInboundOrders/10", "")
	router.ServeHTTP(recorder, req)

	assert.Equal(t, 404, recorder.Code)
}

func TestGetEmployeeInboundOrders_InvalidID(t *testing.T) {
	db := []domain.EmployeeWithInboundOrders{
		{Employee: testEmployee, InboundOrders: 1},
	}

	mockRepository := inbound_order.MockRepository{
		DataMockEmployeesWithIO: db,
		MockErrorGet:            inbound_order.ErrEmployeeWithInboundOrdersNotFound,
	}

	router := createServerEmployeeWithIO(mockRepository)
	req, recorder := createRequestTestEmployeeWithIO(http.MethodGet, "/api/v1/employees/reportInboundOrders/thisIsAnInvalidID", "")
	router.ServeHTTP(recorder, req)

	assert.Equal(t, 400, recorder.Code)
}
