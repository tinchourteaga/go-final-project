package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/cmd/server/handler/requests"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/record/product_record"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type successfulProductRecordResponse struct {
	Data domain.ProductRecord `json:"data"`
}

type unsuccessfulProductRecordResponse struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func setupProductRecordHandlersEngineMock() (ctx *gin.Context, responseRecorder *httptest.ResponseRecorder) {
	logging.InitLog(nil)
	gin.SetMode(gin.TestMode)
	responseRecorder = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(responseRecorder)
	return
}

// TestProductRecord_Create_OK passes when data is correct (return 201 and new domain.ProductRecord)
func TestProductRecord_Create_OK(t *testing.T) {
	// Arrange
	productRecordRequest := requests.ProductRecordPOSTRequest{
		LastUpdateDate: "2022-12-24",
		PurchasePrice:  newFloatPointer(2.4),
		SalePrice:      newFloatPointer(5.6),
		ProductID:      newIntPointer(1),
	}
	id := 1
	expectedCode := http.StatusCreated
	expectedResponse, errConv := productRecordRequest.MapToDomain()
	assert.NoError(t, errConv)
	expectedResponse.ID = id

	// Act
	ctx, responseRecorder := setupProductRecordHandlersEngineMock()
	productRecordService := product_record.ServiceMock{ProductRecordRepository: []domain.ProductRecord{}, ExpectedID: id}
	productRecordHandler := NewProductRecord(&productRecordService)
	body, errMarshal := json.Marshal(&productRecordRequest)
	assert.Nil(t, errMarshal)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productRecordHandler.Create()(ctx)
	var response successfulProductRecordResponse
	errUnmarshal := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, errUnmarshal)
	assert.True(t, productRecordService.FlagSave)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedResponse, response.Data)
}

// TestProductRecord_Create_OKToday passes when data is correct and LasUpdateDate is today (return 201 and new domain.ProductRecord)
func TestProductRecord_Create_OKToday(t *testing.T) {
	// Setup
	now := time.Now()
	var day string
	if now.Day() < 10 {
		day = fmt.Sprintf("0%d", now.Day())
	} else {
		day = fmt.Sprintf("%d", now.Day())
	}
	var month string
	if now.Month() < 10 {
		month = fmt.Sprintf("0%d", now.Month())
	} else {
		month = fmt.Sprintf("%d", now.Month())
	}
	today := fmt.Sprintf("%d-%s-%s", now.Year(), month, day)

	// Arrange
	productRecordRequest := requests.ProductRecordPOSTRequest{
		LastUpdateDate: today,
		PurchasePrice:  newFloatPointer(2.4),
		SalePrice:      newFloatPointer(5.6),
		ProductID:      newIntPointer(1),
	}
	id := 1
	expectedCode := http.StatusCreated
	expectedResponse, errConv := productRecordRequest.MapToDomain()
	assert.NoError(t, errConv)
	expectedResponse.ID = id

	// Act
	ctx, responseRecorder := setupProductRecordHandlersEngineMock()
	productRecordService := product_record.ServiceMock{ProductRecordRepository: []domain.ProductRecord{}, ExpectedID: id}
	productRecordHandler := NewProductRecord(&productRecordService)
	body, errMarshal := json.Marshal(&productRecordRequest)
	assert.Nil(t, errMarshal)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productRecordHandler.Create()(ctx)
	var response successfulProductRecordResponse
	errUnmarshal := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, errUnmarshal)
	assert.True(t, productRecordService.FlagSave)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedResponse, response.Data)
}

// TestProductRecord_Create_Fail passes when data's format is incorrect (return 500 and empty error message)
func TestProductRecord_Create_Fail(t *testing.T) {
	// Arrange
	expectedErr := errors.New("")
	expectedCode := http.StatusInternalServerError

	// Act
	ctx, responseRecorder := setupProductRecordHandlersEngineMock()
	productRecordService := product_record.ServiceMock{}
	productRecordHandler := NewProductRecord(&productRecordService)
	body := []byte("{\"This\": \"has\" \"to\": \"fail\"}")
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productRecordHandler.Create()(ctx)
	var response unsuccessfulProductRecordResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, err)
	assert.False(t, productRecordService.FlagSave)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr.Error(), response.Message)
}

// TestProductRecord_Create_FailNecessaryFields passes when data doesn't have all necessary fields (return 422 and error message)
func TestProductRecord_Create_FailNecessaryFields(t *testing.T) {
	// Arrange
	productRecordRequest := requests.ProductRecordPOSTRequest{
		LastUpdateDate: "2022-12-24",
		PurchasePrice:  newFloatPointer(2.4),
		ProductID:      newIntPointer(1),
	}
	expectedCode := http.StatusUnprocessableEntity
	expectedErr := "Key: 'ProductRecordPOSTRequest.SalePrice' Error:Field validation for 'SalePrice' failed on the 'required' tag"

	// Act
	ctx, responseRecorder := setupProductRecordHandlersEngineMock()
	productRecordService := product_record.ServiceMock{}
	productRecordHandler := NewProductRecord(&productRecordService)
	body, errMarshal := json.Marshal(&productRecordRequest)
	assert.Nil(t, errMarshal)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productRecordHandler.Create()(ctx)
	var response unsuccessfulProductResponse
	errUnmarshal := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, errUnmarshal)
	assert.False(t, productRecordService.FlagSave)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr, response.Message)
}

// TestProductRecord_Create_FailCastError passes when data type doesn't align with struct definition (return 422 and error message)
func TestProductRecord_Create_FailCastError(t *testing.T) {
	// Arrange
	productRecordRequest := "This can not be casted to requests.ProductRecordPOSTRequest"
	expectedCode := http.StatusUnprocessableEntity
	expectedErr := "json: cannot unmarshal string into Go value of type requests.ProductRecordPOSTRequest"

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	productRecordService := product_record.ServiceMock{}
	productRecordHandler := NewProductRecord(&productRecordService)
	body, errMarshal := json.Marshal(&productRecordRequest)
	assert.Nil(t, errMarshal)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productRecordHandler.Create()(ctx)
	var response unsuccessfulProductRecordResponse
	errUnmarshal := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, errUnmarshal)
	assert.False(t, productRecordService.FlagSave)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr, response.Message)
}

// TestProductRecord_Create_Conflict passes when product not found (return 409 and error message product_record.ServiceErrForeignKeyNotFound)
func TestProductRecord_Create_Conflict(t *testing.T) {
	// Arrange
	productRecordRequest := requests.ProductRecordPOSTRequest{
		LastUpdateDate: "2022-12-24",
		PurchasePrice:  newFloatPointer(2.4),
		SalePrice:      newFloatPointer(5.6),
		ProductID:      newIntPointer(1),
	}
	expectedCode := http.StatusConflict
	expectedErr := product_record.ServiceErrForeignKeyNotFound

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	productRecordService := product_record.ServiceMock{ProductRecordRepository: []domain.ProductRecord{}, ForcedErrSave: expectedErr}
	productRecordHandler := NewProductRecord(&productRecordService)
	body, errMarshal := json.Marshal(&productRecordRequest)
	assert.Nil(t, errMarshal)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productRecordHandler.Create()(ctx)
	var response unsuccessfulProductRecordResponse
	errUnmarshal := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, errUnmarshal)
	assert.True(t, productRecordService.FlagSave)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr.Error(), response.Message)
}

// TestProductRecord_Create_OKButNotFound passes when product is created but can't be found in database (return 404 and error product_record.ServiceErrNotFound)
func TestProductRecord_Create_OKButNotFound(t *testing.T) {
	// Arrange
	productRecordRequest := requests.ProductRecordPOSTRequest{
		LastUpdateDate: "2022-12-24",
		PurchasePrice:  newFloatPointer(2.4),
		SalePrice:      newFloatPointer(5.6),
		ProductID:      newIntPointer(1),
	}
	forcedErr := product_record.ServiceErrNotFound
	expectedCode := http.StatusNotFound

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	productRecordService := product_record.ServiceMock{ProductRecordRepository: []domain.ProductRecord{}, ForcedErrSave: forcedErr}
	productRecordHandler := NewProductRecord(&productRecordService)
	body, errMarshal := json.Marshal(&productRecordRequest)
	assert.Nil(t, errMarshal)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productRecordHandler.Create()(ctx)
	var response unsuccessfulProductRecordResponse
	errUnmarshal := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, errUnmarshal)
	assert.True(t, productRecordService.FlagSave)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, ProductRecordErrCreatedButNotFound.Error(), response.Message)
}

// TestProductRecord_Create_InternalServerError passes when unexpected error occurs (return 500 and empty error message)
func TestProductRecord_Create_InternalServerError(t *testing.T) {
	// Arrange
	productRecordRequest := requests.ProductRecordPOSTRequest{
		LastUpdateDate: "2022-12-24",
		PurchasePrice:  newFloatPointer(2.4),
		SalePrice:      newFloatPointer(5.6),
		ProductID:      newIntPointer(1),
	}
	expectedCode := http.StatusInternalServerError
	expectedErr := errors.New("")

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	productRecordService := product_record.ServiceMock{ProductRecordRepository: []domain.ProductRecord{}, ForcedErrSave: expectedErr}
	productRecordHandler := NewProductRecord(&productRecordService)
	body, errMarshal := json.Marshal(&productRecordRequest)
	assert.Nil(t, errMarshal)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productRecordHandler.Create()(ctx)
	var response unsuccessfulProductRecordResponse
	errUnmarshal := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, errUnmarshal)
	assert.True(t, productRecordService.FlagSave)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr.Error(), response.Message)
}

// TestProductRecord_Create_BadDate passes when date is before today's date (return 409 and error ProductRecordErrDate)
func TestProductRecord_Create_BadDate(t *testing.T) {
	// Arrange
	productRecordRequest := requests.ProductRecordPOSTRequest{
		LastUpdateDate: "2021-12-24",
		PurchasePrice:  newFloatPointer(2.4),
		SalePrice:      newFloatPointer(5.6),
		ProductID:      newIntPointer(1),
	}
	expectedCode := http.StatusConflict
	expectedErr := ProductRecordErrDate

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	productRecordService := product_record.ServiceMock{ForcedErrSave: product_record.ServiceErrDate}
	productRecordHandler := NewProductRecord(&productRecordService)
	body, errMarshal := json.Marshal(&productRecordRequest)
	assert.Nil(t, errMarshal)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productRecordHandler.Create()(ctx)
	var response unsuccessfulProductRecordResponse
	errUnmarshal := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, errUnmarshal)
	assert.True(t, productRecordService.FlagSave)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr.Error(), response.Message)
}

// TestProductRecord_Create_BadDateNonsense passes when date is before today's date (return 409 and error ProductRecordErrInvalidDate)
func TestProductRecord_Create_BadDateNonsense(t *testing.T) {
	// Arrange
	productRecordRequest := requests.ProductRecordPOSTRequest{
		LastUpdateDate: "lkskasdfkp",
		PurchasePrice:  newFloatPointer(2.4),
		SalePrice:      newFloatPointer(5.6),
		ProductID:      newIntPointer(1),
	}
	expectedCode := http.StatusBadRequest

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	productRecordService := product_record.ServiceMock{}
	productRecordHandler := NewProductRecord(&productRecordService)
	body, errMarshal := json.Marshal(&productRecordRequest)
	assert.Nil(t, errMarshal)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productRecordHandler.Create()(ctx)
	var response unsuccessfulProductRecordResponse
	errUnmarshal := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, errUnmarshal)
	assert.False(t, productRecordService.FlagSave)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, ProductRecordErrInvalidDate.Error(), response.Message)
}
