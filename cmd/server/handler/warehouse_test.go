package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/cmd/server/handler/requests"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/warehouse"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MOCK SERVICE
type MockWarehouseService struct {
	mockWarehouse     domain.Warehouse
	mockWarehouses    []domain.Warehouse
	mockErrorInternal error
	mockErrorUpdate   error
}

func (s *MockWarehouseService) Get(ctx context.Context, id int) (domain.Warehouse, error) {
	if s.mockErrorInternal != nil {
		return domain.Warehouse{}, s.mockErrorInternal
	}
	return s.mockWarehouse, nil
}

func (s *MockWarehouseService) GetAll(ctx context.Context) ([]domain.Warehouse, error) {
	if s.mockErrorInternal != nil {
		return []domain.Warehouse{}, s.mockErrorInternal
	}
	return s.mockWarehouses, nil
}

func (s *MockWarehouseService) Create(ctx context.Context, address string, telephone string, warehouseCode string, minimumCapacity int, minimumTemperature int) (domain.Warehouse, error) {
	if s.mockErrorInternal != nil {
		return domain.Warehouse{}, s.mockErrorInternal
	}
	return s.mockWarehouse, nil
}

func (s *MockWarehouseService) Delete(ctx context.Context, id int) error {
	if s.mockErrorInternal != nil {
		return s.mockErrorInternal
	}
	return nil
}

func (s *MockWarehouseService) Update(ctx context.Context, id int, address *string, telephone *string, warehouseCode *string, minimumCapacity *int, minimumTemperature *int) (domain.Warehouse, error) {
	if s.mockErrorUpdate != nil {
		return domain.Warehouse{}, s.mockErrorUpdate
	}
	return s.mockWarehouse, nil
}

// MOCK GIN
func mockWarehouseGin(warehouseID string, structBody interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	logging.InitLog(nil)
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.AddParam("id", warehouseID)
	body, _ := json.Marshal(&structBody)
	req := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = req
	return ctx, recorder
}

// RESPONSE STRUCTS

type responseDataWarehouse struct {
	Data domain.Warehouse `json:"data"`
}

type responseDataWarehouses struct {
	Data []domain.Warehouse `json:"data"`
}

type warehouseErrorResponse struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// TESTS

// TestWarehouseGet checks the correct operation of the Get handler method
// Expected HTTP Status code: 200
func TestWarehouseGet(t *testing.T) {
	// arrange
	expectedWarehouse := domain.Warehouse{
		ID:                 1,
		Address:            "Monroe 1230",
		Telephone:          "47470000",
		WarehouseCode:      "DHM1",
		MinimumCapacity:    10,
		MinimumTemperature: 0,
	}
	expectedStatus := http.StatusOK

	mockService := MockWarehouseService{mockWarehouse: expectedWarehouse}
	handler := NewWarehouse(&mockService)

	ctx, recorder := mockWarehouseGin("1", "")

	// act
	handler.Get(ctx)

	// parse response body
	response := recorder.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseDataWarehouse
	err := json.Unmarshal(bytesBody, &body)
	result := body.Data

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedWarehouse, result)
}

// TestWarehouseGetFailureStrConv is correct when strconv function returns an error
// Expected HTTP Status code: 400
func TestWarehouseGetFailureStrConv(t *testing.T) {
	// arrange
	expectedStatus := http.StatusBadRequest
	expectedError := warehouse.ErrBadRequest

	mockService := MockWarehouseService{mockErrorInternal: expectedError}
	handler := NewWarehouse(&mockService)

	warehouseID := "a"
	ctx, recorder := mockWarehouseGin(warehouseID, "")

	// act
	handler.Get(ctx)

	// parse response body
	response := recorder.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body warehouseErrorResponse
	err := json.Unmarshal(bytesBody, &body)
	responseMessage := body.Message

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedError.Error(), responseMessage)
}

// TestWarehouseGetFailureNotFound is correct when the id does not match any warehouse
// Expected HTTP Status code: 404
func TestWarehouseGetFailureNotFound(t *testing.T) {
	// arrange
	expectedStatus := http.StatusNotFound
	expectedError := warehouse.ErrNotFound

	mockService := MockWarehouseService{mockErrorInternal: expectedError}
	handler := NewWarehouse(&mockService)

	ctx, recorder := mockWarehouseGin("1", "")

	// act
	handler.Get(ctx)

	// parse response body
	response := recorder.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body warehouseErrorResponse
	err := json.Unmarshal(bytesBody, &body)
	responseMessage := body.Message

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedError.Error(), responseMessage)
}

// TestWarehouseGetFailureInternal is correct when an internal error is encountered
// Expected HTTP Status code: 500
func TestWarehouseGetFailureInternal(t *testing.T) {
	// arrange
	expectedStatus := http.StatusInternalServerError
	expectedError := warehouse.ErrInternal

	mockService := MockWarehouseService{mockErrorInternal: expectedError}
	handler := NewWarehouse(&mockService)

	ctx, recorder := mockWarehouseGin("1", "")

	// act
	handler.Get(ctx)

	// parse response body
	response := recorder.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body warehouseErrorResponse
	err := json.Unmarshal(bytesBody, &body)
	responseMessage := body.Message

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedError.Error(), responseMessage)
}

// TestWarehouseGetAll checks the correct operation of the GetAll handler method
// Expected HTTP Status code: 200
func TestWarehouseGetAll(t *testing.T) {
	// arrange
	expectedWarehouses := []domain.Warehouse{
		{
			ID:                 1,
			Address:            "Monroe 1230",
			Telephone:          "47470000",
			WarehouseCode:      "DHM1",
			MinimumCapacity:    10,
			MinimumTemperature: 0,
		},
	}
	expectedStatus := http.StatusOK

	mockService := MockWarehouseService{mockWarehouses: expectedWarehouses}
	handler := NewWarehouse(&mockService)

	ctx, recorder := mockWarehouseGin("", "")

	// act
	handler.GetAll(ctx)

	// parse response body
	response := recorder.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseDataWarehouses
	err := json.Unmarshal(bytesBody, &body)
	result := body.Data

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedWarehouses, result)
}

// TestWarehouseGetAllFailureInternal is correct when an internal error is encountered
// Expected HTTP Status code: 500
func TestWarehouseGetAllFailureInternal(t *testing.T) {
	// arrange
	expectedStatus := http.StatusInternalServerError
	expectedError := warehouse.ErrInternal

	mockService := MockWarehouseService{mockErrorInternal: expectedError}
	handler := NewWarehouse(&mockService)

	ctx, recorder := mockWarehouseGin("", "")

	// act
	handler.GetAll(ctx)

	// parse response body
	response := recorder.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body warehouseErrorResponse
	err := json.Unmarshal(bytesBody, &body)
	responseMessage := body.Message

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedError.Error(), responseMessage)
}

// TestWarehouseCreate checks the correct operation of the Create handler method
// Expected HTTP Status code: 201
func TestWarehouseCreate(t *testing.T) {
	// arrange
	address := "Monroe 1230"
	telephone := "47470000"
	warehouseCode := "DHM1"
	minimumCapacity := 10
	minimumTemperature := 0
	requestWarehouse := requests.WarehousePostRequest{
		Address:            &address,
		Telephone:          &telephone,
		WarehouseCode:      &warehouseCode,
		MinimumCapacity:    &minimumCapacity,
		MinimumTemperature: &minimumTemperature,
	}
	expectedWarehouse := domain.Warehouse{
		ID:                 1,
		Address:            address,
		Telephone:          telephone,
		WarehouseCode:      warehouseCode,
		MinimumCapacity:    minimumCapacity,
		MinimumTemperature: minimumTemperature,
	}
	expectedStatus := http.StatusCreated

	mockService := MockWarehouseService{mockWarehouse: expectedWarehouse}
	handler := NewWarehouse(&mockService)

	ctx, recorder := mockWarehouseGin("", requestWarehouse)

	// act
	handler.Create(ctx)

	// parse response body
	response := recorder.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseDataWarehouse
	err := json.Unmarshal(bytesBody, &body)
	result := body.Data

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedWarehouse, result)
}

// TestWarehouseCreateFailureUnprocessable is correct when the sended entity is unprocessable
// Expected HTTP Status code: 422
func TestWarehouseCreateFailureUnprocessable(t *testing.T) {
	// arrange
	address := "Monroe 1230"
	requestWarehouse := requests.WarehousePostRequest{
		Address: &address,
	}
	expectedError := warehouse.ErrBodyValidation
	expectedStatus := http.StatusUnprocessableEntity

	mockService := MockWarehouseService{}
	handler := NewWarehouse(&mockService)

	ctx, recorder := mockWarehouseGin("", requestWarehouse)

	// act
	handler.Create(ctx)

	// parse response body
	response := recorder.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body warehouseErrorResponse
	err := json.Unmarshal(bytesBody, &body)
	responseMessage := body.Message

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedError.Error(), responseMessage)
}

// TestWarehouseCreateFailureConflict is correct when the sended entity has a conflict with another entity
// Expected HTTP Status code: 409
func TestWarehouseCreateFailureConflict(t *testing.T) {
	// arrange
	address := "Monroe 1230"
	telephone := "47470000"
	warehouseCode := "DHM1"
	minimumCapacity := 10
	minimumTemperature := 0
	requestWarehouse := requests.WarehousePostRequest{
		Address:            &address,
		Telephone:          &telephone,
		WarehouseCode:      &warehouseCode,
		MinimumCapacity:    &minimumCapacity,
		MinimumTemperature: &minimumTemperature,
	}
	expectedStatus := http.StatusConflict
	expectedError := warehouse.ErrAlreadyExists

	mockService := MockWarehouseService{mockErrorInternal: expectedError}
	handler := NewWarehouse(&mockService)

	ctx, recorder := mockWarehouseGin("", requestWarehouse)

	// act
	handler.Create(ctx)

	// parse response body
	response := recorder.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body warehouseErrorResponse
	err := json.Unmarshal(bytesBody, &body)
	responseMessage := body.Message

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedError.Error(), responseMessage)
}

// TestWarehouseCreateFailureInternal is correct when an internal error is encountered
// Expected HTTP Status code: 500
func TestWarehouseCreateFailureInternal(t *testing.T) {
	// arrange
	address := "Monroe 1230"
	telephone := "47470000"
	warehouseCode := "DHM1"
	minimumCapacity := 10
	minimumTemperature := 0
	requestWarehouse := requests.WarehousePostRequest{
		Address:            &address,
		Telephone:          &telephone,
		WarehouseCode:      &warehouseCode,
		MinimumCapacity:    &minimumCapacity,
		MinimumTemperature: &minimumTemperature,
	}
	expectedStatus := http.StatusInternalServerError
	expectedError := warehouse.ErrInternal

	mockService := MockWarehouseService{mockErrorInternal: expectedError}
	handler := NewWarehouse(&mockService)

	ctx, recorder := mockWarehouseGin("", requestWarehouse)

	// act
	handler.Create(ctx)

	// parse response body
	response := recorder.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body warehouseErrorResponse
	err := json.Unmarshal(bytesBody, &body)
	responseMessage := body.Message

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedError.Error(), responseMessage)
}

// TestWarehouseUpdate checks the correct operation of the Update handler method
// Expected HTTP Status code: 200
func TestWarehouseUpdate(t *testing.T) {
	// arrange
	address := "Monroe 1230"
	telephone := "47470000"
	warehouseCode := "DHM1"
	minimumCapacity := 10
	minimumTemperature := 0
	requestWarehouse := requests.WarehousePatchRequest{
		Address:            &address,
		Telephone:          &telephone,
		WarehouseCode:      &warehouseCode,
		MinimumCapacity:    &minimumCapacity,
		MinimumTemperature: &minimumTemperature,
	}
	expectedWarehouse := domain.Warehouse{
		ID:                 1,
		Address:            address,
		Telephone:          telephone,
		WarehouseCode:      warehouseCode,
		MinimumCapacity:    minimumCapacity,
		MinimumTemperature: minimumTemperature,
	}
	expectedStatus := http.StatusOK

	mockService := MockWarehouseService{mockWarehouse: expectedWarehouse}
	handler := NewWarehouse(&mockService)

	ctx, recorder := mockWarehouseGin("1", requestWarehouse)

	// act
	handler.Update(ctx)

	// parse response body
	response := recorder.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseDataWarehouse
	err := json.Unmarshal(bytesBody, &body)
	result := body.Data

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedWarehouse, result)
}

// TestWarehouseUpdateFailureStrConv is correct when strconv function returns an error
// Expected HTTP Status code: 400
func TestWarehouseUpdateFailureStrConv(t *testing.T) {
	// arrange
	expectedError := warehouse.ErrBadRequest
	expectedStatus := http.StatusBadRequest

	mockService := MockWarehouseService{}
	handler := NewWarehouse(&mockService)

	warehouseID := "a"
	ctx, recorder := mockWarehouseGin(warehouseID, "")

	// act
	handler.Update(ctx)

	// parse response body
	response := recorder.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body warehouseErrorResponse
	err := json.Unmarshal(bytesBody, &body)
	responseMessage := body.Message

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedError.Error(), responseMessage)
}

// TestWarehouseUpdateFailureUnprocessable is correct when the sended entity is unprocessable
// Expected HTTP Status code: 422
func TestWarehouseUpdateFailureUnprocessable(t *testing.T) {
	// arrange
	requestWarehouse := ""
	expectedError := warehouse.ErrBodyValidation
	expectedStatus := http.StatusUnprocessableEntity

	mockService := MockWarehouseService{}
	handler := NewWarehouse(&mockService)

	ctx, recorder := mockWarehouseGin("1", requestWarehouse)

	// act
	handler.Update(ctx)

	// parse response body
	response := recorder.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body warehouseErrorResponse
	err := json.Unmarshal(bytesBody, &body)
	responseMessage := body.Message

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedError.Error(), responseMessage)
}

// TestWarehouseUpdateFailureConflict is correct when the sended entity has a conflict with another entity
// Expected HTTP Status code: 409
func TestWarehouseUpdateFailureConflict(t *testing.T) {
	// arrange
	warehouseCode := "DHM1"
	requestWarehouse := requests.WarehousePostRequest{
		WarehouseCode: &warehouseCode,
	}
	expectedStatus := http.StatusConflict
	expectedError := warehouse.ErrAlreadyExists

	mockService := MockWarehouseService{mockErrorUpdate: expectedError}
	handler := NewWarehouse(&mockService)

	ctx, recorder := mockWarehouseGin("1", requestWarehouse)

	// act
	handler.Update(ctx)

	// parse response body
	response := recorder.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body warehouseErrorResponse
	err := json.Unmarshal(bytesBody, &body)
	responseMessage := body.Message

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedError.Error(), responseMessage)
}

// TestWarehouseUpdateFailureNotFound is correct when the id does not match any warehouse
// Expected HTTP Status code: 404
func TestWarehouseUpdateFailureNotFound(t *testing.T) {
	// arrange
	address := "Monroe 1230"
	requestWarehouse := requests.WarehousePostRequest{
		Address: &address,
	}
	expectedStatus := http.StatusNotFound
	expectedError := warehouse.ErrNotFound

	mockService := MockWarehouseService{mockErrorUpdate: expectedError}
	handler := NewWarehouse(&mockService)

	ctx, recorder := mockWarehouseGin("1", requestWarehouse)

	// act
	handler.Update(ctx)

	// parse response body
	response := recorder.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body warehouseErrorResponse
	err := json.Unmarshal(bytesBody, &body)
	responseMessage := body.Message

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedError.Error(), responseMessage)
}

// TestWarehouseUpdateFailureInternal is correct when an internal error is encountered
// Expected HTTP Status code: 500
func TestWarehouseUpdateFailureInternal(t *testing.T) {
	// arrange
	address := "Monroe 1230"
	requestWarehouse := requests.WarehousePostRequest{
		Address: &address,
	}
	expectedStatus := http.StatusInternalServerError
	expectedError := warehouse.ErrInternal

	mockService := MockWarehouseService{mockErrorUpdate: expectedError}
	handler := NewWarehouse(&mockService)

	ctx, recorder := mockWarehouseGin("1", requestWarehouse)

	// act
	handler.Update(ctx)

	// parse response body
	response := recorder.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body warehouseErrorResponse
	err := json.Unmarshal(bytesBody, &body)
	responseMessage := body.Message

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedError.Error(), responseMessage)
}

// TestWarehouseDelete checks the correct operation of the Delete handler method
// Expected HTTP Status code: 204
func TestWarehouseDelete(t *testing.T) {
	// arrange
	expectedStatus := http.StatusNoContent

	mockService := MockWarehouseService{}
	handler := NewWarehouse(&mockService)

	ctx, recorder := mockWarehouseGin("1", "")

	// act
	handler.Delete(ctx)

	// parse response
	response := recorder.Result()

	// assert
	assert.Equal(t, expectedStatus, response.StatusCode)
}

// TestWarehouseDeleteFailureStrConv is correct when strconv function returns an error
// Expected HTTP Status code: 400
func TestWarehouseDeleteFailureStrConv(t *testing.T) {
	// arrange
	expectedStatus := http.StatusBadRequest
	expectedError := warehouse.ErrBadRequest

	mockService := MockWarehouseService{}
	handler := NewWarehouse(&mockService)

	warehouseID := "a"
	ctx, recorder := mockWarehouseGin(warehouseID, "")

	// act
	handler.Delete(ctx)

	// parse response body
	response := recorder.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body warehouseErrorResponse
	err := json.Unmarshal(bytesBody, &body)
	responseMessage := body.Message

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedError.Error(), responseMessage)
}

// TestWarehouseDeleteFailureNotFound is correct when the id does not match any warehouse
// Expected HTTP Status code: 404
func TestWarehouseDeleteFailureNotFound(t *testing.T) {
	// arrange
	expectedStatus := http.StatusNotFound
	expectedError := warehouse.ErrNotFound

	mockService := MockWarehouseService{mockErrorInternal: expectedError}
	handler := NewWarehouse(&mockService)

	ctx, recorder := mockWarehouseGin("1", "")

	// act
	handler.Delete(ctx)

	// parse response body
	response := recorder.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body warehouseErrorResponse
	err := json.Unmarshal(bytesBody, &body)
	responseMessage := body.Message

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedError.Error(), responseMessage)
}

// TestWarehouseDeleteFailureInternal is correct when an internal error is encountered
// Expected HTTP Status code: 500
func TestWarehouseDeleteFailureInternal(t *testing.T) {
	// arrange
	expectedStatus := http.StatusInternalServerError
	expectedError := warehouse.ErrInternal

	mockService := MockWarehouseService{mockErrorInternal: expectedError}
	handler := NewWarehouse(&mockService)

	ctx, recorder := mockWarehouseGin("1", "")

	// act
	handler.Delete(ctx)

	// parse response body
	response := recorder.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body warehouseErrorResponse
	err := json.Unmarshal(bytesBody, &body)
	responseMessage := body.Message

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedError.Error(), responseMessage)
}
