package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/cmd/server/handler/requests"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/product"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type successfulProductResponse struct {
	Data domain.Product `json:"data"`
}

type successfulProductSliceResponse struct {
	Data []domain.Product `json:"data"`
}

type unsuccessfulProductResponse struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func newStringPointer(value string) *string {
	return &value
}

func newFloatPointer(value float32) *float32 {
	return &value
}

func newIntPointer(value int) *int {
	return &value
}

func setupProductHandlersEngineMock() (ctx *gin.Context, responseRecorder *httptest.ResponseRecorder) {
	logging.InitLog(nil)
	gin.SetMode(gin.TestMode)
	responseRecorder = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(responseRecorder)
	return
}

// TestProduct_Create_OK passes when data is correct (return 201 and new domain.Product)
func TestProduct_Create_OK(t *testing.T) {
	// Arrange
	productRequest := requests.ProductPOSTRequest{
		Description:                    newStringPointer("hola"),
		ExpirationRate:                 newIntPointer(2),
		FreezingRate:                   newIntPointer(3),
		Height:                         newFloatPointer(3.4),
		Length:                         newFloatPointer(4.3),
		NetWeight:                      newFloatPointer(3.0),
		ProductCode:                    newStringPointer("pcode"),
		RecommendedFreezingTemperature: newFloatPointer(34.5),
		Width:                          newFloatPointer(3.1),
		ProductTypeID:                  newIntPointer(4),
	}
	id := 1
	expectedCode := http.StatusCreated
	expectedResponse := productRequest.MapToDomain()
	expectedResponse.ID = id

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	productService := product.ServiceMock{ProductRepository: []domain.Product{}, ExpectedID: id}
	productHandler := NewProduct(&productService)
	body, errMarshal := json.Marshal(&productRequest)
	assert.Nil(t, errMarshal)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productHandler.Create()(ctx)
	var response successfulProductResponse
	errUnmarshal := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, errUnmarshal)
	assert.True(t, productService.FlagSave)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedResponse, response.Data)
}

// TestProduct_Create_Fail passes when data's format is incorrect (return 500 and empty error message)
func TestProduct_Create_Fail(t *testing.T) {
	// Arrange
	expectedErr := errors.New("")
	expectedCode := http.StatusInternalServerError

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	productService := product.ServiceMock{ProductRepository: []domain.Product{}}
	productHandler := NewProduct(&productService)
	body := []byte("{\"This\": \"has\" \"to\": \"fail\"}")
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productHandler.Create()(ctx)
	var response unsuccessfulProductResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, err)
	assert.False(t, productService.FlagSave)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr.Error(), response.Message)
}

// TestProduct_Create_FailNecessaryFields passes when data doesn't have all necessary fields (return 400 and error message)
func TestProduct_Create_FailNecessaryFields(t *testing.T) {
	// Arrange
	productRequest := requests.ProductPOSTRequest{
		Description:                    newStringPointer("hola"),
		ExpirationRate:                 newIntPointer(2),
		FreezingRate:                   newIntPointer(3),
		Height:                         newFloatPointer(3.4),
		Length:                         newFloatPointer(4.3),
		NetWeight:                      newFloatPointer(3.0),
		ProductCode:                    newStringPointer("pcode"),
		RecommendedFreezingTemperature: newFloatPointer(34.5),
		ProductTypeID:                  newIntPointer(4),
	}
	expectedCode := http.StatusBadRequest
	expectedErr := "Key: 'ProductPOSTRequest.Width' Error:Field validation for 'Width' failed on the 'required' tag"

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	productService := product.ServiceMock{ProductRepository: []domain.Product{}}
	productHandler := NewProduct(&productService)
	body, errMarshal := json.Marshal(&productRequest)
	assert.Nil(t, errMarshal)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productHandler.Create()(ctx)
	var response unsuccessfulProductResponse
	errUnmarshal := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, errUnmarshal)
	assert.False(t, productService.FlagSave)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr, response.Message)
}

// TestProduct_Create_FailCastError passes when data type doesn't align with struct definition (return 422 and error message)
func TestProduct_Create_FailCastError(t *testing.T) {
	// Arrange
	productRequest := "This can not be casted to requests.ProductPOSTRequest"
	expectedCode := http.StatusUnprocessableEntity
	expectedErr := "json: cannot unmarshal string into Go value of type requests.ProductPOSTRequest"

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	productService := product.ServiceMock{ProductRepository: []domain.Product{}}
	productHandler := NewProduct(&productService)
	body, errMarshal := json.Marshal(&productRequest)
	assert.Nil(t, errMarshal)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productHandler.Create()(ctx)
	var response unsuccessfulProductResponse
	errUnmarshal := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, errUnmarshal)
	assert.False(t, productService.FlagSave)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr, response.Message)
}

// TestProduct_Create_Conflict passes when product_code already exists (return 409 and error message product.ServiceErrAlreadyExists)
func TestProduct_Create_Conflict(t *testing.T) {
	// Arrange
	productRequest := requests.ProductPOSTRequest{
		Description:                    newStringPointer("hola"),
		ExpirationRate:                 newIntPointer(2),
		FreezingRate:                   newIntPointer(3),
		Height:                         newFloatPointer(3.4),
		Length:                         newFloatPointer(4.3),
		NetWeight:                      newFloatPointer(3.0),
		ProductCode:                    newStringPointer("pcode"),
		RecommendedFreezingTemperature: newFloatPointer(34.5),
		Width:                          newFloatPointer(3.1),
		ProductTypeID:                  newIntPointer(4),
	}
	expectedCode := http.StatusConflict
	expectedErr := product.ServiceErrAlreadyExists

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	productService := product.ServiceMock{ProductRepository: []domain.Product{}, ForcedErrSave: expectedErr}
	productHandler := NewProduct(&productService)
	body, errMarshal := json.Marshal(&productRequest)
	assert.Nil(t, errMarshal)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productHandler.Create()(ctx)
	var response unsuccessfulProductResponse
	errUnmarshal := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, errUnmarshal)
	assert.True(t, productService.FlagSave)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr.Error(), response.Message)
}

// TestProduct_Create_OKButNotFound passes when product is created but can't be found in database (return 404 and error ProductErrCreatedButNotFound)
func TestProduct_Create_OKButNotFound(t *testing.T) {
	// Arrange
	productRequest := requests.ProductPOSTRequest{
		Description:                    newStringPointer("hola"),
		ExpirationRate:                 newIntPointer(2),
		FreezingRate:                   newIntPointer(3),
		Height:                         newFloatPointer(3.4),
		Length:                         newFloatPointer(4.3),
		NetWeight:                      newFloatPointer(3.0),
		ProductCode:                    newStringPointer("pcode"),
		RecommendedFreezingTemperature: newFloatPointer(34.5),
		Width:                          newFloatPointer(3.1),
		ProductTypeID:                  newIntPointer(4),
		SellerID:                       newIntPointer(3),
	}
	forcedErr := product.ServiceErrNotFound
	expectedCode := http.StatusNotFound

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	productService := product.ServiceMock{ProductRepository: []domain.Product{}, ForcedErrSave: forcedErr}
	productHandler := NewProduct(&productService)
	body, errMarshal := json.Marshal(&productRequest)
	assert.Nil(t, errMarshal)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productHandler.Create()(ctx)
	var response unsuccessfulProductResponse
	errUnmarshal := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, errUnmarshal)
	assert.True(t, productService.FlagSave)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, ProductErrCreatedButNotFound.Error(), response.Message)
}

// TestProduct_Create_InternalServerError passes when unexpected error occurs (return 500 and empty error message)
func TestProduct_Create_InternalServerError(t *testing.T) {
	// Arrange
	productRequest := requests.ProductPOSTRequest{
		Description:                    newStringPointer("hola"),
		ExpirationRate:                 newIntPointer(2),
		FreezingRate:                   newIntPointer(3),
		Height:                         newFloatPointer(3.4),
		Length:                         newFloatPointer(4.3),
		NetWeight:                      newFloatPointer(3.0),
		ProductCode:                    newStringPointer("pcode"),
		RecommendedFreezingTemperature: newFloatPointer(34.5),
		Width:                          newFloatPointer(3.1),
		ProductTypeID:                  newIntPointer(4),
	}
	expectedCode := http.StatusInternalServerError
	expectedErr := errors.New("")

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	productService := product.ServiceMock{ProductRepository: []domain.Product{}, ForcedErrSave: expectedErr}
	productHandler := NewProduct(&productService)
	body, errMarshal := json.Marshal(&productRequest)
	assert.Nil(t, errMarshal)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productHandler.Create()(ctx)
	var response unsuccessfulProductResponse
	errUnmarshal := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, errUnmarshal)
	assert.True(t, productService.FlagSave)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr.Error(), response.Message)
}

// TestProduct_Create_ForeignKeyNotFound passes when seller_id is not in database (return 404 and error product.ServiceErrForeignKeyNotFound)
func TestProduct_Create_ForeignKeyNotFound(t *testing.T) {
	// Arrange
	productRequest := requests.ProductPOSTRequest{
		Description:                    newStringPointer("hola"),
		ExpirationRate:                 newIntPointer(2),
		FreezingRate:                   newIntPointer(3),
		Height:                         newFloatPointer(3.4),
		Length:                         newFloatPointer(4.3),
		NetWeight:                      newFloatPointer(3.0),
		ProductCode:                    newStringPointer("pcode"),
		RecommendedFreezingTemperature: newFloatPointer(34.5),
		Width:                          newFloatPointer(3.1),
		ProductTypeID:                  newIntPointer(4),
	}
	expectedCode := http.StatusNotFound
	expectedErr := product.ServiceErrForeignKeyNotFound

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	productService := product.ServiceMock{ProductRepository: []domain.Product{}, ForcedErrSave: expectedErr}
	productHandler := NewProduct(&productService)
	body, errMarshal := json.Marshal(&productRequest)
	assert.Nil(t, errMarshal)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productHandler.Create()(ctx)
	var response unsuccessfulProductResponse
	errUnmarshal := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, errUnmarshal)
	assert.True(t, productService.FlagSave)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr.Error(), response.Message)
}

// TestProduct_GetAll_OK passes when query is successful (return 200 and slice of all domain.Product)
func TestProduct_GetAll_OK(t *testing.T) {
	// Arrange
	productRepository := []domain.Product{
		{
			ID: 1,
		},
		{
			ID: 2,
		},
	}
	expectedCode := http.StatusOK

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	productService := product.ServiceMock{ProductRepository: productRepository}
	productHandler := NewProduct(&productService)
	productHandler.GetAll()(ctx)
	var response successfulProductSliceResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, err)
	assert.True(t, productService.FlagGetAll)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, productRepository, response.Data)
}

// TestProduct_GetAll_OKWithEmptyDB passes when query is successful but database is empty (return 200 and empty slice of all domain.Product)
func TestProduct_GetAll_OKWithEmptyDB(t *testing.T) {
	// Arrange
	expectedCode := http.StatusOK

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	productService := product.ServiceMock{ProductRepository: []domain.Product{}}
	productHandler := NewProduct(&productService)
	productHandler.GetAll()(ctx)
	var response successfulProductSliceResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, err)
	assert.True(t, productService.FlagGetAll)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, 0, len(response.Data))
}

// TestProduct_GetAll_InternalServerError passes when query is unsuccessful (return 500 and empty error message)
func TestProduct_GetAll_InternalServerError(t *testing.T) {
	// Arrange
	expectedCode := http.StatusInternalServerError
	expectedErr := errors.New("")

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	productService := product.ServiceMock{ProductRepository: []domain.Product{}, ForcedErrGetAll: product.ServiceErrInternal}
	productHandler := NewProduct(&productService)
	productHandler.GetAll()(ctx)
	var response unsuccessfulProductResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, err)
	assert.True(t, productService.FlagGetAll)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr.Error(), response.Message)
}

// TestProduct_Get_OK passes when id exists (return 200 and domain.Product with given id)
func TestProduct_Get_OK(t *testing.T) {
	// Arrange
	productRepository := []domain.Product{
		{
			ID:                             1,
			Description:                    "Tomatoes",
			ExpirationRate:                 70,
			FreezingRate:                   20,
			Height:                         20.4,
			Length:                         10.3,
			NetWeight:                      0.0,
			ProductCode:                    "kasbkj8ats9aka9",
			RecommendedFreezingTemperature: -12.6,
			Width:                          6.7,
			ProductTypeID:                  3,
			SellerID:                       newIntPointer(5),
		},
		{
			ID: 2,
		},
	}
	searchID := 1
	expectedCode := http.StatusOK

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	ctx.AddParam("id", fmt.Sprintf("%d", searchID))
	productService := product.ServiceMock{ProductRepository: productRepository}
	productHandler := NewProduct(&productService)
	productHandler.Get()(ctx)
	var response successfulProductResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, err)
	assert.True(t, productService.FlagGet)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, productRepository[0], response.Data)
}

// TestProduct_Get_IDNonExistent passes when the given id is not in database (return 404 and error ProductErrNotFound)
func TestProduct_Get_IDNonExistent(t *testing.T) {
	// Arrange
	searchID := 1
	forcedErr := product.ServiceErrNotFound
	expectedCode := http.StatusNotFound
	expectedErr := ProductErrNotFound

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	ctx.AddParam("id", fmt.Sprintf("%d", searchID))
	productService := product.ServiceMock{ProductRepository: []domain.Product{}, ForcedErrGet: forcedErr}
	productHandler := NewProduct(&productService)
	productHandler.Get()(ctx)
	var response unsuccessfulProductResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, err)
	assert.True(t, productService.FlagGet)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr.Error(), response.Message)
}

// TestProduct_Get_InvalidID passes when the given id is not a number (return 400 and error ProductErrInvalidID)
func TestProduct_Get_InvalidID(t *testing.T) {
	// Arrange
	searchID := "badID"
	expectedCode := http.StatusBadRequest
	expectedErr := ProductErrInvalidID

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	ctx.AddParam("id", searchID)
	productService := product.ServiceMock{ProductRepository: []domain.Product{}}
	productHandler := NewProduct(&productService)
	productHandler.Get()(ctx)
	var response unsuccessfulProductResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, err)
	assert.False(t, productService.FlagGet)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr.Error(), response.Message)
}

// TestProduct_Get_InternalServerError passes when unexpected error occurs (return 500 and empty error message)
func TestProduct_Get_InternalServerError(t *testing.T) {
	// Arrange
	searchID := 1
	expectedCode := http.StatusInternalServerError
	expectedErr := errors.New("")

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	ctx.AddParam("id", fmt.Sprintf("%d", searchID))
	productService := product.ServiceMock{ProductRepository: []domain.Product{}, ForcedErrGet: product.ServiceErrInternal}
	productHandler := NewProduct(&productService)
	productHandler.Get()(ctx)
	var response unsuccessfulProductResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, err)
	assert.True(t, productService.FlagGet)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr.Error(), response.Message)
}

// TestProduct_PartialUpdate_OK passes when data is correct (return 200 and updated domain.Product)
func TestProduct_PartialUpdate_OK(t *testing.T) {
	// Arrange
	productRequest := requests.ProductPATCHRequest{
		Description:    newStringPointer("hola"),
		ExpirationRate: newIntPointer(2),
	}
	productRepository := []domain.Product{
		{
			ID:                             1,
			Description:                    "Tomatoes",
			ExpirationRate:                 70,
			FreezingRate:                   20,
			Height:                         20.4,
			Length:                         10.3,
			NetWeight:                      0.0,
			ProductCode:                    "kasbkj8ats9aka9",
			RecommendedFreezingTemperature: -12.6,
			Width:                          6.7,
			ProductTypeID:                  3,
			SellerID:                       newIntPointer(5),
		},
	}
	searchID := 1
	expectedCode := http.StatusOK
	expectedResponse := productRepository[0]
	expectedResponse.Description = "hola"
	expectedResponse.ExpirationRate = 2

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	ctx.AddParam("id", fmt.Sprintf("%d", searchID))
	productService := product.ServiceMock{ProductRepository: productRepository}
	productHandler := NewProduct(&productService)
	body, errMarshal := json.Marshal(&productRequest)
	assert.Nil(t, errMarshal)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productHandler.PartialUpdate()(ctx)
	var response successfulProductResponse
	errUnmarshal := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, errUnmarshal)
	assert.True(t, productService.FlagPartialUpdate)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedResponse, response.Data)
}

// TestProduct_PartialUpdate_IDNonExistent passes when the given id is not in database (return 404 and error ProductErrNotFound)
func TestProduct_PartialUpdate_IDNonExistent(t *testing.T) {
	// Arrange
	productRequest := requests.ProductPATCHRequest{}
	searchID := 1
	forcedErr := product.ServiceErrNotFound
	expectedCode := http.StatusNotFound
	expectedErr := ProductErrNotFound

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	ctx.AddParam("id", fmt.Sprintf("%d", searchID))
	productService := product.ServiceMock{ProductRepository: []domain.Product{}, ForcedErrPartialUpdate: forcedErr}
	productHandler := NewProduct(&productService)
	body, errMarshal := json.Marshal(&productRequest)
	assert.Nil(t, errMarshal)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productHandler.PartialUpdate()(ctx)
	var response unsuccessfulProductResponse
	errUnmarshal := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, errUnmarshal)
	assert.True(t, productService.FlagPartialUpdate)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr.Error(), response.Message)
}

// TestProduct_PartialUpdate_InvalidID passes when the given id is not a number (return 400 and error ProductErrInvalidID)
func TestProduct_PartialUpdate_InvalidID(t *testing.T) {
	// Arrange
	productRequest := requests.ProductPATCHRequest{}
	searchID := "badID"
	expectedCode := http.StatusBadRequest
	expectedErr := ProductErrInvalidID

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	ctx.AddParam("id", searchID)
	productService := product.ServiceMock{ProductRepository: []domain.Product{}}
	productHandler := NewProduct(&productService)
	body, errMarshal := json.Marshal(&productRequest)
	assert.Nil(t, errMarshal)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productHandler.PartialUpdate()(ctx)
	var response unsuccessfulProductResponse
	errUnmarshal := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, errUnmarshal)
	assert.False(t, productService.FlagPartialUpdate)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr.Error(), response.Message)
}

// TestProduct_PartialUpdate_FailCastError passes when data type doesn't align with struct definition (return 422 and error message)
func TestProduct_PartialUpdate_FailCastError(t *testing.T) {
	// Arrange
	productRequest := "This can not be casted to requests.ProductPATCHRequest"
	searchID := 1
	expectedCode := http.StatusUnprocessableEntity
	expectedErr := "json: cannot unmarshal string into Go value of type requests.ProductPATCHRequest"

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	ctx.AddParam("id", fmt.Sprintf("%d", searchID))
	productService := product.ServiceMock{ProductRepository: []domain.Product{}}
	productHandler := NewProduct(&productService)
	body, errMarshal := json.Marshal(&productRequest)
	assert.Nil(t, errMarshal)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productHandler.PartialUpdate()(ctx)
	var response unsuccessfulProductResponse
	errUnmarshal := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, errUnmarshal)
	assert.False(t, productService.FlagPartialUpdate)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr, response.Message)
}

// TestProduct_PartialUpdate_Conflict passes when product_code already exists (return 409 and error message product.ServiceErrAlreadyExists)
func TestProduct_PartialUpdate_Conflict(t *testing.T) {
	// Arrange
	productRequest := requests.ProductPATCHRequest{}
	searchID := 1
	expectedCode := http.StatusConflict
	expectedErr := product.ServiceErrAlreadyExists

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	ctx.AddParam("id", fmt.Sprintf("%d", searchID))
	productService := product.ServiceMock{ProductRepository: []domain.Product{}, ForcedErrPartialUpdate: expectedErr}
	productHandler := NewProduct(&productService)
	body, errMarshal := json.Marshal(&productRequest)
	assert.Nil(t, errMarshal)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productHandler.PartialUpdate()(ctx)
	var response unsuccessfulProductResponse
	errUnmarshal := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, errUnmarshal)
	assert.True(t, productService.FlagPartialUpdate)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr.Error(), response.Message)
}

// TestProduct_PartialUpdate_InternalServerError passes when unexpected error occurs (return 500 and empty error message)
func TestProduct_PartialUpdate_InternalServerError(t *testing.T) {
	// Arrange
	productRequest := requests.ProductPATCHRequest{}
	searchID := 1
	expectedCode := http.StatusInternalServerError
	expectedErr := errors.New("")

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	ctx.AddParam("id", fmt.Sprintf("%d", searchID))
	productService := product.ServiceMock{ProductRepository: []domain.Product{}, ForcedErrPartialUpdate: product.ServiceErrInternal}
	productHandler := NewProduct(&productService)
	body, errMarshal := json.Marshal(&productRequest)
	assert.Nil(t, errMarshal)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productHandler.PartialUpdate()(ctx)
	var response unsuccessfulProductResponse
	errUnmarshal := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, errUnmarshal)
	assert.True(t, productService.FlagPartialUpdate)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr.Error(), response.Message)
}

// TestProduct_PartialUpdate_ForeignKeyNotFound passes when seller_id is not in database (return 404 and error product.ServiceErrForeignKeyNotFound)
func TestProduct_PartialUpdate_ForeignKeyNotFound(t *testing.T) {
	// Arrange
	productRequest := requests.ProductPATCHRequest{}
	searchID := 1
	expectedCode := http.StatusNotFound
	expectedErr := product.ServiceErrForeignKeyNotFound

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	ctx.AddParam("id", fmt.Sprintf("%d", searchID))
	productService := product.ServiceMock{ProductRepository: []domain.Product{}, ForcedErrPartialUpdate: expectedErr}
	productHandler := NewProduct(&productService)
	body, errMarshal := json.Marshal(&productRequest)
	assert.Nil(t, errMarshal)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productHandler.PartialUpdate()(ctx)
	var response unsuccessfulProductResponse
	errUnmarshal := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, errUnmarshal)
	assert.True(t, productService.FlagPartialUpdate)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr.Error(), response.Message)
}

// TestProduct_PartialUpdate_Fail passes when data's format is incorrect (return 500 and empty error message)
func TestProduct_PartialUpdate_Fail(t *testing.T) {
	// Arrange
	expectedErr := errors.New("")
	searchID := 1
	expectedCode := http.StatusInternalServerError

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	ctx.AddParam("id", fmt.Sprintf("%d", searchID))
	productService := product.ServiceMock{ProductRepository: []domain.Product{}}
	productHandler := NewProduct(&productService)
	body := []byte("{\"This\": \"has\" \"to\": \"fail\"}")
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request
	productHandler.PartialUpdate()(ctx)
	var response unsuccessfulProductResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, err)
	assert.False(t, productService.FlagSave)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr.Error(), response.Message)
}

// TestProduct_Delete_OK passes when id exists and deletion is successful (return 204)
func TestProduct_Delete_OK(t *testing.T) {
	// Arrange
	searchID := 1
	expectedCode := http.StatusNoContent

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	ctx.AddParam("id", fmt.Sprintf("%d", searchID))
	productService := product.ServiceMock{ProductRepository: []domain.Product{}}
	productHandler := NewProduct(&productService)
	productHandler.Delete()(ctx)

	// Assert
	assert.True(t, productService.FlagDelete)
	assert.Equal(t, expectedCode, responseRecorder.Code)
}

// TestProduct_Delete_IDNonExistent passes when the given id is not in database (return 404 and error ProductErrNotFound)
func TestProduct_Delete_IDNonExistent(t *testing.T) {
	// Arrange
	searchID := 1
	forcedErr := product.ServiceErrNotFound
	expectedCode := http.StatusNotFound
	expectedErr := ProductErrNotFound

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	ctx.AddParam("id", fmt.Sprintf("%d", searchID))
	productService := product.ServiceMock{ProductRepository: []domain.Product{}, ForcedErrDelete: forcedErr}
	productHandler := NewProduct(&productService)
	productHandler.Delete()(ctx)
	var response unsuccessfulProductResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, err)
	assert.True(t, productService.FlagDelete)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr.Error(), response.Message)
}

// TestProduct_Delete_InvalidID passes when the given id is not a number (return 400 and error ProductErrInvalidID)
func TestProduct_Delete_InvalidID(t *testing.T) {
	// Arrange
	searchID := "badID"
	expectedCode := http.StatusBadRequest
	expectedErr := ProductErrInvalidID

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	ctx.AddParam("id", searchID)
	productService := product.ServiceMock{ProductRepository: []domain.Product{}}
	productHandler := NewProduct(&productService)
	productHandler.Delete()(ctx)
	var response unsuccessfulProductResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, err)
	assert.False(t, productService.FlagDelete)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr.Error(), response.Message)
}

// TestProduct_Delete_InternalServerError passes when unexpected error occurs (return 500 and empty error message)
func TestProduct_Delete_InternalServerError(t *testing.T) {
	// Arrange
	searchID := 1
	expectedCode := http.StatusInternalServerError
	expectedErr := errors.New("")

	// Act
	ctx, responseRecorder := setupProductHandlersEngineMock()
	ctx.AddParam("id", fmt.Sprintf("%d", searchID))
	productService := product.ServiceMock{ProductRepository: []domain.Product{}, ForcedErrDelete: product.ServiceErrInternal}
	productHandler := NewProduct(&productService)
	productHandler.Delete()(ctx)
	var response unsuccessfulProductResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)

	// Assert
	assert.Nil(t, err)
	assert.True(t, productService.FlagDelete)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr.Error(), response.Message)
}
