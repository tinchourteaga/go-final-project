package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/employee"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	logging.InitLog(nil)
}

type successfulResponseEmployee struct {
	Data domain.Employee `json:"data"`
}

type successfulSliceResponseEmployee struct {
	Data []domain.Employee `json:"data"`
}

func createServerEmployee(mockRepository employee.MockRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	service := employee.NewService(&mockRepository)
	employeeHandler := NewEmployee(service)
	router := gin.Default()

	employeesRoutesGroup := router.Group("/api/v1/employees")

	employeesRoutesGroup.GET("/", employeeHandler.GetAll())
	employeesRoutesGroup.GET("/:id", employeeHandler.Get())
	employeesRoutesGroup.POST("/", employeeHandler.Create())
	employeesRoutesGroup.PATCH("/:id", employeeHandler.Update())
	employeesRoutesGroup.DELETE("/:id", employeeHandler.Delete())

	return router
}

func createRequestTestEmployee(method, url, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")

	return req, httptest.NewRecorder()
}

/* =============== GET ALL =============== */
func TestGetAllEmployee(t *testing.T) {
	var sliceResponse successfulSliceResponseEmployee
	db := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 3},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	mockRepository := employee.MockRepository{
		DataMock:  db,
		MockError: nil,
	}

	router := createServerEmployee(mockRepository)
	req, recorder := createRequestTestEmployee(http.MethodGet, "/api/v1/employees/", "")
	router.ServeHTTP(recorder, req)

	assert.Equal(t, 200, recorder.Code)
	err := json.Unmarshal(recorder.Body.Bytes(), &sliceResponse)
	assert.Nil(t, err)
	assert.Equal(t, mockRepository.DataMock, sliceResponse.Data)
}

func TestGetAllEmployeeFail(t *testing.T) {
	db := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 3},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	mockRepository := employee.MockRepository{
		DataMock:  db,
		MockError: errors.New("error"),
	}

	router := createServerEmployee(mockRepository)
	req, recorder := createRequestTestEmployee(http.MethodGet, "/api/v1/employees/", "")
	router.ServeHTTP(recorder, req)
	assert.Equal(t, 500, recorder.Code)
}

/* =============== GET =============== */
func TestGetEmployee(t *testing.T) {
	var response successfulResponseEmployee
	db := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 3},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	mockRepository := employee.MockRepository{
		DataMock:  db,
		MockError: nil,
	}

	router := createServerEmployee(mockRepository)
	req, recorder := createRequestTestEmployee(http.MethodGet, "/api/v1/employees/2", "")
	router.ServeHTTP(recorder, req)

	assert.Equal(t, 200, recorder.Code)
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, db[1], response.Data)
}

func TestGetEmployeeFail(t *testing.T) {
	db := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 3},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	mockRepository := employee.MockRepository{
		DataMock:  db,
		MockError: errors.New(""),
	}

	router := createServerEmployee(mockRepository)
	req, recorder := createRequestTestEmployee(http.MethodGet, "/api/v1/employees/5", "")
	router.ServeHTTP(recorder, req)
	assert.Equal(t, 404, recorder.Code)
}

func TestGetEmployeeInvalidID(t *testing.T) {
	db := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 3},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	mockRepository := employee.MockRepository{
		DataMock:  db,
		MockError: errors.New(""),
	}

	router := createServerEmployee(mockRepository)
	req, recorder := createRequestTestEmployee(http.MethodGet, "/api/v1/employees/thisIsAnInvalidID", "")
	router.ServeHTTP(recorder, req)
	assert.Equal(t, 400, recorder.Code)
}

/* =============== SAVE =============== */

func TestSaveEmployee(t *testing.T) {
	var response successfulResponseEmployee
	expectedEmployee := domain.Employee{ID: 3, CardNumberID: "111111", FirstName: "Martin", LastName: "Urteaga", WarehouseID: 6}
	db := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 3},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	mockRepository := employee.MockRepository{
		DataMock:  db,
		MockError: nil,
	}

	router := createServerEmployee(mockRepository)
	req, recorder := createRequestTestEmployee(http.MethodPost, "/api/v1/employees/", `{"card_number_id": "111111", "first_name": "Martin", "last_name": "Urteaga", "warehouse_id": 6}`)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, 201, recorder.Code)
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, expectedEmployee, response.Data)
}

func TestSaveEmployeeFail(t *testing.T) {
	db := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 3},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	mockRepository := employee.MockRepository{
		DataMock:  db,
		MockError: errors.New(""),
	}

	router := createServerEmployee(mockRepository)
	req, recorder := createRequestTestEmployee(http.MethodPost, "/api/v1/employees/", `{"last_name": "Urteaga", "warehouse_id": 6}`)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, 422, recorder.Code)
}

func TestSaveEmployeeConflict(t *testing.T) {
	db := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 3},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	mockRepository := employee.MockRepository{
		DataMock:  db,
		MockError: errors.New(""),
	}

	router := createServerEmployee(mockRepository)
	req, recorder := createRequestTestEmployee(http.MethodPost, "/api/v1/employees/", `{"card_number_id": "123456", "first_name": "Martin", "last_name": "Urteaga", "warehouse_id": 6}`)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, 409, recorder.Code)
}

/* =============== UPDATE =============== */
func TestUpdateEmployee(t *testing.T) {
	var response successfulResponseEmployee
	db := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 3},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	mockRepository := employee.MockRepository{
		DataMock:  db,
		MockError: nil,
	}

	router := createServerEmployee(mockRepository)
	req, recorder := createRequestTestEmployee(http.MethodPatch, "/api/v1/employees/1", `{"first_name": "Martin", "last_name": "Urteaga", "warehouse_id": 6}`)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, 200, recorder.Code)
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, mockRepository.DataMock[0], response.Data)
}

func TestUpdateEmployeeFail(t *testing.T) {
	db := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 3},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	mockRepository := employee.MockRepository{
		DataMock:  db,
		MockError: errors.New(""),
	}

	router := createServerEmployee(mockRepository)
	req, recorder := createRequestTestEmployee(http.MethodPatch, "/api/v1/employees/5", `{"first_name": "Martin", "last_name": "Urteaga", "warehouse_id": 6}`)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, 404, recorder.Code)
}

func TestUpdateEmployeeInvalidID(t *testing.T) {
	db := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 3},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	mockRepository := employee.MockRepository{
		DataMock:  db,
		MockError: errors.New(""),
	}

	router := createServerEmployee(mockRepository)
	req, recorder := createRequestTestEmployee(http.MethodPatch, "/api/v1/employees/thisIsAnInvalidID", `{"first_name": "Martin", "last_name": "Urteaga", "warehouse_id": 6}`)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, 400, recorder.Code)
}

/* =============== DELETE =============== */
func TestDeleteEmployee(t *testing.T) {
	db := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 3},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	mockRepository := employee.MockRepository{
		DataMock:  db,
		MockError: nil,
	}

	router := createServerEmployee(mockRepository)
	req, recorder := createRequestTestEmployee(http.MethodDelete, "/api/v1/employees/1", "")
	router.ServeHTTP(recorder, req)
	assert.Equal(t, 204, recorder.Code)
}

func TestDeleteEmployeeFail(t *testing.T) {
	db := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 3},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	mockRepository := employee.MockRepository{
		DataMock:  db,
		MockError: employee.ErrEmployeeNotFound,
	}

	router := createServerEmployee(mockRepository)
	req, recorder := createRequestTestEmployee(http.MethodDelete, "/api/v1/employees/5", "")
	router.ServeHTTP(recorder, req)
	assert.Equal(t, 404, recorder.Code)
}

func TestDeleteEmployeeInvalidID(t *testing.T) {
	db := []domain.Employee{
		{ID: 1, CardNumberID: "123456", FirstName: "John", LastName: "Doe", WarehouseID: 3},
		{ID: 2, CardNumberID: "654321", FirstName: "Jane", LastName: "Doe", WarehouseID: 7},
	}

	mockRepository := employee.MockRepository{
		DataMock:  db,
		MockError: employee.ErrEmployeeNotFound,
	}

	router := createServerEmployee(mockRepository)
	req, recorder := createRequestTestEmployee(http.MethodDelete, "/api/v1/employees/thisIsAnInvalidID", "")
	router.ServeHTTP(recorder, req)
	assert.Equal(t, 400, recorder.Code)
}
