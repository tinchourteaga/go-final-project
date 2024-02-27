package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/locality"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockServiceLocality struct {
	ReportCarry []domain.ReportCarries
	Locality    domain.Locality
	DataMock    []domain.Locality
	Report      []domain.ReportSellers
	ErrorGet    error
	ErrorCreate error
	ErrorReport error
}

// *---------------------- Mock service functions -----------------*
func (l *MockServiceLocality) Get(ctx context.Context, id string) (locality domain.Locality, err error) {
	if l.ErrorGet != nil {
		err = l.ErrorGet
		return
	}
	locality = l.Locality
	return
}

func (l *MockServiceLocality) Create(ctx context.Context, locality domain.Locality) (loc domain.Locality, err error) {
	if l.ErrorCreate != nil {
		err = l.ErrorCreate
		return
	}
	loc = l.Locality
	return
}

func (l *MockServiceLocality) ReportSellers(ctx context.Context, locality_id *string) (report []domain.ReportSellers, err error) {
	if l.ErrorReport != nil {
		err = l.ErrorReport
		return
	}
	report = l.Report
	return
}

func (s *MockServiceLocality) ReportCarries(ctx context.Context, locality_id *string) (report []domain.ReportCarries, err error) {
	if s.ErrorReport != nil {
		return []domain.ReportCarries{}, s.ErrorReport
	}
	return s.ReportCarry, nil
}

// *---------------------- Others functions -----------------*
func createServerLocality() (ctx *gin.Context, recorder *httptest.ResponseRecorder) {
	logging.InitLog(nil)
	gin.SetMode(gin.TestMode)
	recorder = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(recorder)
	return
}

type responseLocality struct {
	Data domain.Locality `json:"data"`
}

type responseReportSellers struct {
	Data []domain.ReportSellers `json:"data"`
}

type responseErrorLocality struct {
	Message string `json:"message"`
}

type responseReportCarries struct {
	Data []domain.ReportCarries `json:"data"`
}

// *--------------------------- GetReportCarries ----------------------*
// TestGetReportCarries checks the correct operation of the GetReportCarries handler method when a non empty id is recieved
// Expected HTTP Status code: 200
func TestGetReportCarries(t *testing.T) {
	// Arrange
	testReport := domain.ReportCarries{
		LocalityID:   "0001",
		LocalityName: "Santiago",
		CarriesCount: 2,
	}
	expectedReport := []domain.ReportCarries{testReport}

	ctx, rr := createServerLocality()
	// add param id to request
	req := &http.Request{
		URL: &url.URL{},
	}
	query := req.URL.Query()
	query.Add("id", testReport.LocalityID)
	req.URL.RawQuery = query.Encode()
	ctx.Request = req

	service := MockServiceLocality{
		ReportCarry: expectedReport,
	}
	handler := NewLocality(&service)

	// Act
	handler.GetReportCarries()(ctx)

	/* Parse response body */
	response := rr.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseReportCarries
	err := json.Unmarshal(bytesBody, &body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, expectedReport, body.Data)
}

// TestGetReportCarries checks the correct operation of the GetReportCarries handler method when a empty id is recieved
// Expected HTTP Status code: 200
func TestGetReportCarriesEmptyID(t *testing.T) {
	// Arrange
	testReport := domain.ReportCarries{
		LocalityID:   "0001",
		LocalityName: "Santiago",
		CarriesCount: 2,
	}
	expectedReport := []domain.ReportCarries{testReport}

	ctx, rr := createServerLocality()

	service := MockServiceLocality{
		ReportCarry: expectedReport,
	}
	handler := NewLocality(&service)

	// Act
	handler.GetReportCarries()(ctx)

	/* Parse response body */
	response := rr.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseReportCarries
	err := json.Unmarshal(bytesBody, &body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, expectedReport, body.Data)
}

// TestGetReportCarriesFailNotFound is correct when service returns a NotFound error
// Expected HTTP Status code: 404
func TestGetReportCarriesFailNotFound(t *testing.T) {
	// Arrange
	expectedError := locality.ErrNotFound
	expectedStatus := http.StatusNotFound

	ctx, rr := createServerLocality()

	service := MockServiceLocality{
		ErrorReport: expectedError,
	}
	handler := NewLocality(&service)

	// Act
	handler.GetReportCarries()(ctx)

	/* Parse response body */
	response := rr.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseErrorLocality
	err := json.Unmarshal(bytesBody, &body)
	responseMessage := body.Message

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedError.Error(), responseMessage)
}

// TestGetReportCarriesFailInternal is correct when service returns an Internal error
// Expected HTTP Status code: 500
func TestGetReportCarriesFailInternal(t *testing.T) {
	// Arrange
	expectedError := locality.ErrInternal
	expectedStatus := http.StatusInternalServerError

	ctx, rr := createServerLocality()

	service := MockServiceLocality{
		ErrorReport: expectedError,
	}
	handler := NewLocality(&service)

	// Act
	handler.GetReportCarries()(ctx)

	/* Parse response body */
	response := rr.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseErrorLocality
	err := json.Unmarshal(bytesBody, &body)
	responseMessage := body.Message

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
	assert.Equal(t, expectedError.Error(), responseMessage)
}

// *--------------------------- Get -------------------------*
// TestGet_Locality_OK passes when return correct locality (status code 200)
func TestGet_Locality_OK(t *testing.T) {
	// Arrange
	expectedLocality := domain.Locality{
		ID:           "5700",
		LocalityName: "San Luis",
		ProvinceName: "San Luis",
		CountryName:  "Argentina",
	}
	id := "5700"
	ctx, rr := createServerLocality()
	ctx.AddParam("id", id)

	service := MockServiceLocality{
		Locality: expectedLocality,
	}
	handler := NewLocality(&service)

	// Act
	handler.Get()(ctx)

	/* Parse response body */
	response := rr.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseLocality
	err := json.Unmarshal(bytesBody, &body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, expectedLocality, body.Data)
}

// TestGetIdNotFound_Locality passes when return error id not found (status code 404)
func TestGetIdNotFound_Locality(t *testing.T) {
	// Arrange
	id := "1534"
	expectedError := fmt.Errorf("Id %s does not exist", id)
	ctx, rr := createServerLocality()
	ctx.AddParam("id", id)

	service := MockServiceLocality{
		ErrorGet: locality.ErrNotFound,
	}
	handler := NewLocality(&service)

	// Act
	handler.Get()(ctx)

	/* Parse response body */
	response := rr.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseErrorLocality
	err := json.Unmarshal(bytesBody, &body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.EqualError(t, expectedError, body.Message)
}

// TestGetInternalError_Locality passes when return an error (status code 500)
func TestGetInternalError_Locality(t *testing.T) {
	// Arrange
	expectedError := errors.New("Error on Get")
	ctx, rr := createServerLocality()

	service := MockServiceLocality{
		ErrorGet: expectedError,
	}
	handler := NewLocality(&service)

	// Act
	handler.Get()(ctx)

	/* Parse response body */
	response := rr.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseErrorLocality
	err := json.Unmarshal(bytesBody, &body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.EqualError(t, expectedError, body.Message)
}

// *--------------------------- Create ----------------------*
// TestCreate_Locality_OK passes when return locality created (status code 201)
func TestCreate_Locality_OK(t *testing.T) {
	// Arrange
	expectedCreateLocality := domain.Locality{
		ID:           "5700",
		LocalityName: "San Luis",
		ProvinceName: "San Luis",
		CountryName:  "Argentina",
	}

	id := "5700"
	localityName := "San Luis"
	provinceName := "San Luis"
	countryName := "Argentina"
	localityRequestBody := domain.Locality{
		ID:           id,
		LocalityName: localityName,
		ProvinceName: provinceName,
		CountryName:  countryName,
	}

	ctx, rr := createServerLocality()

	service := MockServiceLocality{
		Locality: expectedCreateLocality,
	}
	handler := NewLocality(&service)

	body, _ := json.Marshal(&localityRequestBody)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request

	// Act
	handler.Create()(ctx)

	/* Parse response body */
	var responseBody responseLocality
	err := json.Unmarshal(rr.Body.Bytes(), &responseBody)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, expectedCreateLocality, responseBody.Data)
}

// TestCreateErrorBadRequest_Seller passes when return an error for bad request (status code 400)
func TestCreateErrorBadRequest_Locality(t *testing.T) {
	// Arrange
	expectedError := errors.New("Bad Request, missing required fields")

	localityRequestBody := domain.Locality{
		ID:           "5700",
		LocalityName: "San Luis",
		CountryName:  "Argentina",
	}

	ctx, rr := createServerLocality()

	service := MockServiceLocality{
		ErrorCreate: expectedError,
	}
	handler := NewLocality(&service)

	body, _ := json.Marshal(&localityRequestBody)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request

	// Act
	handler.Create()(ctx)

	/* Parse response body */
	var responseBody responseErrorLocality
	err := json.Unmarshal(rr.Body.Bytes(), &responseBody)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.EqualError(t, expectedError, responseBody.Message)
}

// TestCreateErrorConflict_Locality passes when return an error for id already existe (status code 409)
func TestCreateErrorConflict_Locality(t *testing.T) {
	// Arrange
	expectedError := locality.ErrAlreadyExists

	localityRequestBody := domain.Locality{
		ID:           "5700",
		LocalityName: "San Luis",
		ProvinceName: "San Luis",
		CountryName:  "Argentina",
	}

	ctx, rr := createServerLocality()

	service := MockServiceLocality{
		ErrorCreate: expectedError,
	}
	handler := NewLocality(&service)

	body, _ := json.Marshal(&localityRequestBody)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request

	// Act
	handler.Create()(ctx)

	/* Parse response body */
	var responseBody responseErrorLocality
	err := json.Unmarshal(rr.Body.Bytes(), &responseBody)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.EqualError(t, expectedError, responseBody.Message)

}

// TestCreateUnprocessableEntity_Locality passes when return an error for unprocessable entity (status code 422)
func TestCreateUnprocessableEntity_Locality(t *testing.T) {
	// Arrange
	expectedError := errors.New("invalid request")
	ctx, rr := createServerLocality()

	service := MockServiceLocality{
		ErrorCreate: expectedError,
	}
	handler := NewLocality(&service)

	// Act
	handler.Create()(ctx)

	/* Parse response body */
	response := rr.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseErrorLocality
	err := json.Unmarshal(bytesBody, &body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	assert.EqualError(t, expectedError, body.Message)

}

// TestCreateInternalError_Locality passes when return an error for bad request (status code 500)
func TestCreateInternalError_Locality(t *testing.T) {
	// Arrange
	expectedError := locality.ErrInternal

	localityRequestBody := domain.Locality{
		ID:           "5700",
		LocalityName: "San Luis",
		ProvinceName: "San Luis",
		CountryName:  "Argentina",
	}

	ctx, rr := createServerLocality()

	service := MockServiceLocality{
		ErrorCreate: expectedError,
	}
	handler := NewLocality(&service)

	body, _ := json.Marshal(&localityRequestBody)
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(body)),
	}
	ctx.Request = request

	// Act
	handler.Create()(ctx)

	/* Parse response body */
	var responseBody responseErrorLocality
	err := json.Unmarshal(rr.Body.Bytes(), &responseBody)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.EqualError(t, expectedError, responseBody.Message)

}

// *--------------------------- Report ----------------------*
// TestGetReportSellers_OK passes when handler method recive a non empty id (status code 200)
func TestGetReportSellers_OK(t *testing.T) {
	// Arrange
	testReport := domain.ReportSellers{
		LocalityID:   "5700",
		LocalityName: "San Luis",
		SellersCount: 3,
	}
	expectedReport := []domain.ReportSellers{testReport}

	ctx, rr := createServerLocality()
	// add param id to request
	req := &http.Request{
		URL: &url.URL{},
	}
	query := req.URL.Query()
	query.Add("id", testReport.LocalityID)
	req.URL.RawQuery = query.Encode()
	ctx.Request = req

	service := MockServiceLocality{
		Report: expectedReport,
	}
	handler := NewLocality(&service)

	// Act
	handler.GetReportSellers()(ctx)

	/* Parse response body */
	response := rr.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseReportSellers
	err := json.Unmarshal(bytesBody, &body)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, expectedReport, body.Data)
}

// TestGetReportSellers_FailNotFound passes when service returns a NotFound error (status code 404)
func TestGetReportSellers_FailNotFound(t *testing.T) {
	// Arrange
	expectedError := locality.ErrNotFound

	ctx, rr := createServerLocality()

	service := MockServiceLocality{
		ErrorReport: locality.ErrNotFound,
	}
	handler := NewLocality(&service)

	// Act
	handler.GetReportSellers()(ctx)

	/* Parse response body */
	response := rr.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseErrorLocality
	err := json.Unmarshal(bytesBody, &body)
	responseMessage := body.Message

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
	assert.Equal(t, expectedError.Error(), responseMessage)
}

// TestGetReportSellers_FailInternal passes when service returns an Internal error (status code 500)
func TestGetReportSellers_FailInternal(t *testing.T) {
	// Arrange
	expectedError := locality.ErrInternal

	ctx, rr := createServerLocality()

	service := MockServiceLocality{
		ErrorReport: locality.ErrInternal,
	}
	handler := NewLocality(&service)

	// Act
	handler.GetReportSellers()(ctx)

	/* Parse response body */
	response := rr.Result()
	bytesBody, _ := io.ReadAll(response.Body)
	var body responseErrorLocality
	err := json.Unmarshal(bytesBody, &body)
	responseMessage := body.Message

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	assert.Equal(t, expectedError.Error(), responseMessage)
}
