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
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/carry"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MOCK SERVICE
type MockCarryService struct {
	mockCarry domain.Carry
	mockError error
}

func (s *MockCarryService) Save(ctx context.Context, CID string, CompanyName string, Address string, Telephone string, Locality_id string) (domain.Carry, error) {
	if s.mockError != nil {
		return domain.Carry{}, s.mockError
	}
	return s.mockCarry, nil
}

// MOCK GIN
func mockCarryGin(structBody interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	logging.InitLog(nil)
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	body, _ := json.Marshal(&structBody)
	req := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = req
	return ctx, recorder
}

// RESPONSE STRUCTS

type responseDataCarry struct {
	Data domain.Carry `json:"data"`
}

type carryErrorResponse struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// TESTS

// TestCarrySave checks the correct operation of the Save handler method
// Expected HTTP Status code: 201
func TestCarrySave(t *testing.T) {
	// arrange
	CID := "CID#1"
	company_name := "some name"
	address := "corrientes 800"
	telephone := "4567-4567"
	locality_id := "1"
	requestCarry := requests.CarryPostRequest{
		CID:         &CID,
		CompanyName: &company_name,
		Address:     &address,
		Telephone:   &telephone,
		Locality_id: &locality_id,
	}
	expectedCarry := domain.Carry{
		ID:          1,
		CID:         CID,
		CompanyName: company_name,
		Address:     address,
		Telephone:   telephone,
		Locality_id: locality_id,
	}
	expectedStatus := http.StatusCreated

	mockService := MockCarryService{mockCarry: expectedCarry}
	handler := NewCarry(&mockService)

	ctx, recorder := mockCarryGin(requestCarry)

	// act
	handler.Save(ctx)

	// parse response body
	response := recorder.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseDataCarry
	err := json.Unmarshal(bytesBody, &body)
	result := body.Data

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedCarry, result)
}

// TestCarrySaveFailureUnprocessable is correct when the sended entity is unprocessable
// Expected HTTP Status code: 422
func TestCarrySaveFailureUnprocessable(t *testing.T) {
	// arrange
	address := "Monroe 1230"
	requestCarry := requests.CarryPostRequest{
		Address: &address,
	}
	expectedError := carry.ErrBodyValidation
	expectedStatus := http.StatusUnprocessableEntity

	mockService := MockCarryService{}
	handler := NewCarry(&mockService)

	ctx, recorder := mockCarryGin(requestCarry)

	// act
	handler.Save(ctx)

	// parse response body
	response := recorder.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body carryErrorResponse
	err := json.Unmarshal(bytesBody, &body)
	responseMessage := body.Message

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedError.Error(), responseMessage)
}

// TestCarrySaveFailureConflict is correct when the sended entity has a conflict with another entity
// Expected HTTP Status code: 409
func TestCarrySaveFailureConflict(t *testing.T) {
	// arrange
	CID := "CID#1"
	company_name := "some name"
	address := "corrientes 800"
	telephone := "4567-4567"
	locality_id := "1"
	requestCarry := requests.CarryPostRequest{
		CID:         &CID,
		CompanyName: &company_name,
		Address:     &address,
		Telephone:   &telephone,
		Locality_id: &locality_id,
	}
	expectedStatus := http.StatusConflict
	expectedError := carry.ErrAlreadyExists

	mockService := MockCarryService{mockError: expectedError}
	handler := NewCarry(&mockService)

	ctx, recorder := mockCarryGin(requestCarry)

	// act
	handler.Save(ctx)

	// parse response body
	response := recorder.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body carryErrorResponse
	err := json.Unmarshal(bytesBody, &body)
	responseMessage := body.Message

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedError.Error(), responseMessage)
}

// TestCarrySaveFailureConstraint is correct when the sended entity has a foreign key constraint
// Expected HTTP Status code: 409
func TestCarrySaveFailureConstraint(t *testing.T) {
	// arrange
	CID := "CID#1"
	company_name := "some name"
	address := "corrientes 800"
	telephone := "4567-4567"
	locality_id := "1"
	requestCarry := requests.CarryPostRequest{
		CID:         &CID,
		CompanyName: &company_name,
		Address:     &address,
		Telephone:   &telephone,
		Locality_id: &locality_id,
	}
	expectedStatus := http.StatusConflict
	expectedError := carry.ErrFKConstraint

	mockService := MockCarryService{mockError: expectedError}
	handler := NewCarry(&mockService)

	ctx, recorder := mockCarryGin(requestCarry)

	// act
	handler.Save(ctx)

	// parse response body
	response := recorder.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body carryErrorResponse
	err := json.Unmarshal(bytesBody, &body)
	responseMessage := body.Message

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedError.Error(), responseMessage)
}

// TestCarrySaveFailureDataLong is correct when the sended entity exceeds the CID maximum length
// Expected HTTP Status code: 422
func TestCarrySaveFailureDataLong(t *testing.T) {
	// arrange
	CID := "CID#11111111"
	company_name := "some name"
	address := "corrientes 800"
	telephone := "4567-4567"
	locality_id := "1"
	requestCarry := requests.CarryPostRequest{
		CID:         &CID,
		CompanyName: &company_name,
		Address:     &address,
		Telephone:   &telephone,
		Locality_id: &locality_id,
	}
	expectedStatus := http.StatusUnprocessableEntity
	expectedError := carry.ErrDataLong

	mockService := MockCarryService{mockError: expectedError}
	handler := NewCarry(&mockService)

	ctx, recorder := mockCarryGin(requestCarry)

	// act
	handler.Save(ctx)

	// parse response body
	response := recorder.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body carryErrorResponse
	err := json.Unmarshal(bytesBody, &body)
	responseMessage := body.Message

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedError.Error(), responseMessage)
}

// TestCarrySaveFailureInternal is correct when an internal error is encountered
// Expected HTTP Status code: 500
func TestCarrySaveFailureInternal(t *testing.T) {
	// arrange
	CID := "CID#1"
	company_name := "some name"
	address := "corrientes 800"
	telephone := "4567-4567"
	locality_id := "1"
	requestCarry := requests.CarryPostRequest{
		CID:         &CID,
		CompanyName: &company_name,
		Address:     &address,
		Telephone:   &telephone,
		Locality_id: &locality_id,
	}
	expectedStatus := http.StatusInternalServerError
	expectedError := carry.ErrInternal

	mockService := MockCarryService{mockError: expectedError}
	handler := NewCarry(&mockService)

	ctx, recorder := mockCarryGin(requestCarry)

	// act
	handler.Save(ctx)

	// parse response body
	response := recorder.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body carryErrorResponse
	err := json.Unmarshal(bytesBody, &body)
	responseMessage := body.Message

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedError.Error(), responseMessage)
}
