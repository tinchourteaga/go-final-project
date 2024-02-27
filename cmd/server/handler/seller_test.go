package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/cmd/server/handler/requests"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/seller"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
)

type MockServiceSeller struct {
	Seller      domain.Seller
	DataMock    []domain.Seller
	ErrorGet    error
	ErrorGetAll error
	ErrorCreate error
	ErrorUpdate error
	ErrorDelete error
}

// *---------------------- Mock service functions -----------------*
func (s *MockServiceSeller) GetAll(ctx context.Context) (sellers []domain.Seller, err error) {
	if s.ErrorGetAll != nil {
		err = s.ErrorGetAll
		return
	}
	sellers = append(sellers, s.DataMock...)
	return
}

func (s *MockServiceSeller) Get(ctx context.Context, id int) (sell domain.Seller, err error) {
	if s.ErrorGet != nil {
		err = s.ErrorGet
		return
	}
	sell = s.Seller
	return
}

func (s *MockServiceSeller) Create(ctx context.Context, seller domain.Seller) (sell domain.Seller, err error) {
	if s.ErrorCreate != nil {
		err = s.ErrorCreate
		return
	}
	sell = s.Seller
	return
}

func (s *MockServiceSeller) Delete(ctx context.Context, id int) (err error) {
	if s.ErrorDelete != nil {
		err = s.ErrorDelete
		return
	}
	return
}

func (s *MockServiceSeller) Update(ctx context.Context, id int, cid *int, companyName, address, telephone, locality_id *string) (sell domain.Seller, err error) {
	if s.ErrorUpdate != nil {
		err = s.ErrorUpdate
		return
	}
	sell = s.Seller
	return
}

// *---------------------- Others functions -----------------*
func createServerSeller() (ctx *gin.Context, recorder *httptest.ResponseRecorder) {
	logging.InitLog(nil)
	gin.SetMode(gin.TestMode)
	recorder = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(recorder)
	return
}

type responseGetAllSellers struct {
	Data []domain.Seller `json:"data"`
}
type responseSeller struct {
	Data domain.Seller `json:"data"`
}

type responseError struct {
	Message string `json:"message"`
}

// *--------------------------- GetAll ----------------------*
// TestGetAll_Sellers passes when return correct sellers (status code 200)
func TestGetAll_Sellers(t *testing.T) {
	// Arrange
	expectedSellers := MockServiceSeller{
		DataMock: []domain.Seller{{
			ID:          1,
			CID:         1,
			CompanyName: "Kiosco 1",
			Address:     "Junin 323",
			Telephone:   "2664727336",
			Locality_id: "5700",
		}, {
			ID:          2,
			CID:         2,
			CompanyName: "Kiosco 2",
			Address:     "Colon 664",
			Telephone:   "2664233242",
			Locality_id: "5700",
		}},
	}

	ctx, rr := createServerSeller()

	service := MockServiceSeller{
		DataMock: expectedSellers.DataMock,
	}
	handler := NewSeller(&service)

	// Act
	handler.GetAll()(ctx)

	/* Parse response body */
	response := rr.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseGetAllSellers
	err := json.Unmarshal(bytesBody, &body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, expectedSellers.DataMock, body.Data)
}

// TestGetAllFail_Sellers passes when return an error in GetAll from db (status code 500)
func TestGetAllFail_Sellers(t *testing.T) {
	// Arrange
	expectedError := seller.ErrInternal

	ctx, rr := createServerSeller()

	service := MockServiceSeller{
		ErrorGetAll: seller.ErrInternal,
	}
	handler := NewSeller(&service)

	// Act
	handler.GetAll()(ctx)

	/* Parse response body */
	response := rr.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseError
	err := json.Unmarshal(bytesBody, &body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.EqualError(t, expectedError, body.Message)
}

// *--------------------------- Get -------------------------*
// TestGet_Seller passes when return correct seller (status code 200)
func TestGet_Seller(t *testing.T) {
	// Arrange
	expectedSeller := domain.Seller{
		ID:          1,
		CID:         1,
		CompanyName: "Kiosco 1",
		Address:     "Junin 323",
		Telephone:   "2664727336",
		Locality_id: "5700",
	}
	id := 1
	ctx, rr := createServerSeller()
	ctx.AddParam("id", fmt.Sprintf("%d", id))

	service := MockServiceSeller{
		Seller: expectedSeller,
	}
	handler := NewSeller(&service)

	// Act
	handler.Get()(ctx)

	/* Parse response body */
	response := rr.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseSeller
	err := json.Unmarshal(bytesBody, &body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, expectedSeller, body.Data)
}

// TestGetIvalidId_Seller passes when return error invalid id (status code 400)
func TestGetIvalidId_Seller(t *testing.T) {
	// Arrange
	expectedError := errors.New("Invalid ID")
	id := "asda"
	ctx, rr := createServerSeller()
	ctx.AddParam("id", id)

	service := MockServiceSeller{
		ErrorGet: expectedError,
	}
	handler := NewSeller(&service)

	// Act
	handler.Get()(ctx)

	/* Parse response body */
	response := rr.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseError
	err := json.Unmarshal(bytesBody, &body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.EqualError(t, expectedError, body.Message)
}

// TestGetIdNotFound_Seller passes when return error id not found (status code 404)
func TestGetIdNotFound_Seller(t *testing.T) {
	// Arrange
	id := 15
	expectedError := fmt.Errorf("Id %d does not exist", id)
	ctx, rr := createServerSeller()
	ctx.AddParam("id", fmt.Sprintf("%d", id))

	service := MockServiceSeller{
		ErrorGet: seller.ErrNotFound,
	}
	handler := NewSeller(&service)

	// Act
	handler.Get()(ctx)

	/* Parse response body */
	response := rr.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseError
	err := json.Unmarshal(bytesBody, &body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.EqualError(t, expectedError, body.Message)
}

// TestGetInternalServerError_Seller passes when return an internal error (status code 500)
func TestGetInternalServerError_Seller(t *testing.T) {
	// Arrange
	expectedError := seller.ErrInternal
	id := 1
	ctx, rr := createServerSeller()
	ctx.AddParam("id", fmt.Sprintf("%d", id))

	service := MockServiceSeller{
		ErrorGet: seller.ErrInternal,
	}
	handler := NewSeller(&service)

	// Act
	handler.Get()(ctx)

	/* Parse response body */
	response := rr.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseError
	err := json.Unmarshal(bytesBody, &body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.EqualError(t, expectedError, body.Message)
}

// *--------------------------- Create ----------------------*
// TestCreate_Seller passes when return seller created (status code 201)
func TestCreate_Seller(t *testing.T) {
	// Arrange
	expectedCreateSeller := domain.Seller{
		ID:          1,
		CID:         1,
		CompanyName: "Kiosco 1",
		Address:     "Junin 323",
		Telephone:   "2664727336",
		Locality_id: "5700",
	}

	cid := 1
	companyName := "Kiosco 1"
	address := "Junin 323"
	telephone := "2664727336"
	locality := "5700"
	sellerRequestBody := requests.SellerPostRequest{
		CID:         &cid,
		CompanyName: &companyName,
		Address:     &address,
		Telephone:   &telephone,
		Locality_id: &locality,
	}

	ctx, rr := createServerSeller()

	service := MockServiceSeller{
		Seller: expectedCreateSeller,
	}
	handler := NewSeller(&service)

	body, _ := json.Marshal(&sellerRequestBody)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request

	// Act
	handler.Create()(ctx)

	/* Parse response body */
	var responseBody responseSeller
	err := json.Unmarshal(rr.Body.Bytes(), &responseBody)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, expectedCreateSeller, responseBody.Data)
}

// TestCreateErrorBadRequest_Seller passes when return an error for bad request (status code 400)
func TestCreateErrorBadRequest_Seller(t *testing.T) {
	// Arrange
	expectedError := errors.New("Bad Request, missing required fields")

	cid := 1
	address := "Junin 323"
	telephone := "2664727336"
	locality := "5700"
	sellerRequestBody := requests.SellerPostRequest{
		CID:         &cid,
		Address:     &address,
		Telephone:   &telephone,
		Locality_id: &locality,
	}

	ctx, rr := createServerSeller()

	service := MockServiceSeller{
		ErrorCreate: expectedError,
	}
	handler := NewSeller(&service)

	body, _ := json.Marshal(&sellerRequestBody)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request

	// Act
	handler.Create()(ctx)

	/* Parse response body */
	var responseBody responseError
	err := json.Unmarshal(rr.Body.Bytes(), &responseBody)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.EqualError(t, expectedError, responseBody.Message)

}

// TestCreateErrorConflict_Seller passes when return an error for cid already existe (status code 409)
func TestCreateErrorConflict_Seller(t *testing.T) {
	// Arrange
	expectedError := seller.ErrAlreadyExists

	cid := 1
	companyName := "Kiosco 1"
	address := "Junin 323"
	telephone := "2664727336"
	locality := "5700"
	sellerRequestBody := requests.SellerPostRequest{
		CID:         &cid,
		CompanyName: &companyName,
		Address:     &address,
		Telephone:   &telephone,
		Locality_id: &locality,
	}

	ctx, rr := createServerSeller()

	service := MockServiceSeller{
		ErrorCreate: expectedError,
	}
	handler := NewSeller(&service)

	body, _ := json.Marshal(&sellerRequestBody)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request

	// Act
	handler.Create()(ctx)

	/* Parse response body */
	var responseBody responseError
	err := json.Unmarshal(rr.Body.Bytes(), &responseBody)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.EqualError(t, expectedError, responseBody.Message)

}

// TestCreateUnprocessableEntity_Seller passes when return an error for unprocessable entity (status code 422)
func TestCreateUnprocessableEntity_Seller(t *testing.T) {
	// Arrange
	expectedError := errors.New("invalid request")
	ctx, rr := createServerSeller()

	service := MockServiceSeller{
		ErrorCreate: expectedError,
	}
	handler := NewSeller(&service)

	// Act
	handler.Create()(ctx)

	/* Parse response body */
	response := rr.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseError
	err := json.Unmarshal(bytesBody, &body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	assert.EqualError(t, expectedError, body.Message)

}

// TestCreateFail_Seller passes when return an error for db (status code 500)
func TestCreateFail_Seller(t *testing.T) {
	// Arrange
	expectedError := seller.ErrInternal

	cid := 1
	companyName := "Kiosco 1"
	address := "Junin 323"
	telephone := "2664727336"
	locality := "5700"
	sellerRequestBody := requests.SellerPostRequest{
		CID:         &cid,
		CompanyName: &companyName,
		Address:     &address,
		Telephone:   &telephone,
		Locality_id: &locality,
	}

	ctx, rr := createServerSeller()

	service := MockServiceSeller{
		ErrorCreate: expectedError,
	}
	handler := NewSeller(&service)

	body, _ := json.Marshal(&sellerRequestBody)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request

	// Act
	handler.Create()(ctx)

	/* Parse response body */
	var responseBody responseError
	err := json.Unmarshal(rr.Body.Bytes(), &responseBody)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.EqualError(t, expectedError, responseBody.Message)

}

// *--------------------------- Update ----------------------*
// TestUpdate_Seller passes when return the updated seller (status code 200)
func TestUpdate_Seller(t *testing.T) {
	// Arrange
	expectedCreateSeller := domain.Seller{
		ID:          1,
		CID:         1,
		CompanyName: "Kiosco Patch",
		Address:     "Colon 323",
		Telephone:   "2664727336",
		Locality_id: "5700",
	}

	id := 1
	ctx, rr := createServerSeller()
	ctx.AddParam("id", fmt.Sprintf("%d", id))

	cid := 1
	companyName := "Kiosco Patch"
	address := "Colon 323"
	telephone := "2664727336"
	locality := "5700"
	sellerRequestBody := requests.SellerPatchRequest{
		CID:         &cid,
		CompanyName: &companyName,
		Address:     &address,
		Telephone:   &telephone,
		Locality_id: &locality,
	}

	service := MockServiceSeller{
		Seller: expectedCreateSeller,
	}
	handler := NewSeller(&service)

	body, _ := json.Marshal(&sellerRequestBody)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request

	// Act
	handler.Update()(ctx)

	/* Parse response body */
	var responseBody responseSeller
	err := json.Unmarshal(rr.Body.Bytes(), &responseBody)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, expectedCreateSeller, responseBody.Data)
}

// TestUpdateIvalidId_Seller passes when return error invalid id (status code 400)
func TestUpdateIvalidId_Seller(t *testing.T) {
	// Arrange
	expectedError := errors.New("Invalid ID")
	id := "asda"
	ctx, rr := createServerSeller()
	ctx.AddParam("id", id)

	service := MockServiceSeller{
		ErrorUpdate: expectedError,
	}
	handler := NewSeller(&service)

	// Act
	handler.Update()(ctx)

	/* Parse response body */
	response := rr.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseError
	err := json.Unmarshal(bytesBody, &body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.EqualError(t, expectedError, body.Message)
}

// TestUpdateIdNotFound_Seller passes when return error id not found (status code 404)
func TestUpdateIdNotFound_Seller(t *testing.T) {
	// Arrange
	id := 15
	expectedError := fmt.Errorf("Id %d does not exist", id)
	ctx, rr := createServerSeller()
	ctx.AddParam("id", fmt.Sprintf("%d", id))

	cid := 1
	companyName := "Kiosco 1"
	address := "Junin 323"
	telephone := "2664727336"
	locality := "5700"
	sellerRequestBody := requests.SellerPostRequest{
		CID:         &cid,
		CompanyName: &companyName,
		Address:     &address,
		Telephone:   &telephone,
		Locality_id: &locality,
	}

	service := MockServiceSeller{
		ErrorUpdate: seller.ErrNotFound,
	}
	handler := NewSeller(&service)
	body, _ := json.Marshal(&sellerRequestBody)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request

	// Act
	handler.Update()(ctx)

	/* Parse response body */
	response := rr.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var responseBody responseError
	err := json.Unmarshal(bytesBody, &responseBody)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.EqualError(t, expectedError, responseBody.Message)
}

// TestUpdateErrorConflict_Seller passes when return an error for cid already existe (status code 409)
func TestUpdateErrorConflict_Seller(t *testing.T) {
	// Arrange
	id := 1
	expectedError := seller.ErrAlreadyExists
	ctx, rr := createServerSeller()
	ctx.AddParam("id", fmt.Sprintf("%d", id))

	cid := 1
	companyName := "Kiosco 1"
	address := "Junin 323"
	telephone := "2664727336"
	locality := "5700"
	sellerRequestBody := requests.SellerPostRequest{
		CID:         &cid,
		CompanyName: &companyName,
		Address:     &address,
		Telephone:   &telephone,
		Locality_id: &locality,
	}

	service := MockServiceSeller{
		ErrorUpdate: seller.ErrAlreadyExists,
	}
	handler := NewSeller(&service)

	body, _ := json.Marshal(&sellerRequestBody)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request

	// Act
	handler.Update()(ctx)

	/* Parse response body */
	var responseBody responseError
	err := json.Unmarshal(rr.Body.Bytes(), &responseBody)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.EqualError(t, expectedError, responseBody.Message)

}

// TestCreateUnprocessablEentity_Seller passes when return an error for unprocessable entity (status code 422)
func TestUpdateUnprocessablEentity_Seller(t *testing.T) {
	// Arrange
	id := 1
	ctx, rr := createServerSeller()
	ctx.AddParam("id", fmt.Sprintf("%d", id))
	expectedError := errors.New("invalid request")

	service := MockServiceSeller{
		ErrorUpdate: expectedError,
	}
	handler := NewSeller(&service)

	// Act
	handler.Update()(ctx)

	/* Parse response body */
	response := rr.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseError
	err := json.Unmarshal(bytesBody, &body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	assert.EqualError(t, expectedError, body.Message)

}

// TestUpdateFail_Seller passes when return an error for db (status code 500)
func TestUpdateFail_Seller(t *testing.T) {
	// Arrange
	id := 1
	ctx, rr := createServerSeller()
	ctx.AddParam("id", fmt.Sprintf("%d", id))
	expectedError := seller.ErrInternal

	cid := 1
	companyName := "Kiosco 1"
	address := "Junin 323"
	telephone := "2664727336"
	locality := "5700"
	sellerRequestBody := requests.SellerPostRequest{
		CID:         &cid,
		CompanyName: &companyName,
		Address:     &address,
		Telephone:   &telephone,
		Locality_id: &locality,
	}

	service := MockServiceSeller{
		ErrorUpdate: seller.ErrInternal,
	}
	handler := NewSeller(&service)

	body, _ := json.Marshal(&sellerRequestBody)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request

	// Act
	handler.Update()(ctx)

	/* Parse response body */
	var responseBody responseError
	err := json.Unmarshal(rr.Body.Bytes(), &responseBody)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.EqualError(t, expectedError, responseBody.Message)

}

// *--------------------------- Delete ----------------------*
// TestDelete_Seller passes when delete function is succesfull (status code 204)
func TestDelete_Seller(t *testing.T) {
	// Arrange
	id := 1
	ctx, rr := createServerSeller()
	ctx.AddParam("id", fmt.Sprintf("%d", id))

	service := MockServiceSeller{
		DataMock: []domain.Seller{},
	}
	handler := NewSeller(&service)

	// Act
	handler.Delete()(ctx)

	// Assert
	assert.Equal(t, http.StatusNoContent, rr.Code)
}

// TestDeleteInvalidId_Seller passes when return error invalid id (status code 400)
func TestDeleteInvalidId_Seller(t *testing.T) {
	// Arrange
	expectedError := errors.New("Invalid ID")
	id := "asda"
	ctx, rr := createServerSeller()
	ctx.AddParam("id", id)

	service := MockServiceSeller{
		ErrorDelete: expectedError,
	}
	handler := NewSeller(&service)

	// Act
	handler.Delete()(ctx)

	/* Parse response body */
	response := rr.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseError
	err := json.Unmarshal(bytesBody, &body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.EqualError(t, expectedError, body.Message)

}

// TestDeleteIdNotFound_Seller passes when return error id not found (status code 404)
func TestDeleteIdNotFound_Seller(t *testing.T) {
	// Arrange
	id := 15
	expectedError := fmt.Errorf("Id %d does not exist", id)
	ctx, rr := createServerSeller()
	ctx.AddParam("id", fmt.Sprintf("%d", id))

	service := MockServiceSeller{
		ErrorDelete: seller.ErrNotFound,
	}
	handler := NewSeller(&service)

	// Act
	handler.Delete()(ctx)

	/* Parse response body */
	response := rr.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseError
	err := json.Unmarshal(bytesBody, &body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.EqualError(t, expectedError, body.Message)
}

// TestDeleteInternalServerError_Seller passes when return an internal error (status code 500)
func TestDeleteInternalServerError_Seller(t *testing.T) {
	// Arrange
	expectedError := seller.ErrInternal
	id := 1
	ctx, rr := createServerSeller()
	ctx.AddParam("id", fmt.Sprintf("%d", id))

	service := MockServiceSeller{
		ErrorDelete: seller.ErrInternal,
	}
	handler := NewSeller(&service)

	// Act
	handler.Delete()(ctx)

	/* Parse response body */
	response := rr.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseError
	err := json.Unmarshal(bytesBody, &body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.EqualError(t, expectedError, body.Message)
}
