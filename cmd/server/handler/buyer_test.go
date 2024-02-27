package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/buyer"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func createServer(mockRepository buyer.MockRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	service := buyer.NewService(&mockRepository)
	handler := NewBuyer(service)

	r := gin.Default()

	pr := r.Group("/api/v1/buyers")

	pr.GET("", handler.GetAll())
	pr.GET("/:id", handler.Get())
	pr.POST("/", handler.Create())
	pr.DELETE("/:id", handler.Delete())
	pr.PATCH("/:id", handler.Update())

	return r
}

// createRequestTest returns a request and a response recorder
func createRequestTest(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")

	//guarda la response que obtiene el servidor
	return req, httptest.NewRecorder()
}

func TestGetAllSuccess(t *testing.T) {
	//arrange
	ListBuyers := []domain.Buyer{
		{ID: 1, CardNumberID: "001", FirstName: "Comprador 1", LastName: "Vendedor 1"},
		{ID: 2, CardNumberID: "002", FirstName: "Comprador 2", LastName: "Vendedor 2"},
		{ID: 3, CardNumberID: "003", FirstName: "Comprador 3", LastName: "Vendedor 3"},
	}

	//Act
	repo := buyer.MockRepository{
		Data: ListBuyers,
	}

	r := createServer(repo)
	req, recorder := createRequestTest(http.MethodGet, "/api/v1/buyers", "")
	r.ServeHTTP(recorder, req)
	////arrange
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestGetAllFail(t *testing.T) {

	repo := buyer.MockRepository{
		Data: nil,
	}

	r := createServer(repo)
	req, recorder := createRequestTest(http.MethodGet, "/api/v1/buyers", "")
	r.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestGetSuccess(t *testing.T) {
	//arrange
	ListBuyers := []domain.Buyer{
		{ID: 1, CardNumberID: "001", FirstName: "Comprador 1", LastName: "Vendedor 1"},
		{ID: 2, CardNumberID: "002", FirstName: "Comprador 2", LastName: "Vendedor 2"},
		{ID: 3, CardNumberID: "003", FirstName: "Comprador 3", LastName: "Vendedor 3"},
	}

	//Act
	repo := buyer.MockRepository{
		Data: ListBuyers,
	}

	r := createServer(repo)
	req, recorder := createRequestTest(http.MethodGet, "/api/v1/buyers/2", "")
	r.ServeHTTP(recorder, req)
	//arrange
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestGetFailBadRequest(t *testing.T) {
	//arrange
	ListBuyers := []domain.Buyer{
		{ID: 1, CardNumberID: "001", FirstName: "Comprador 1", LastName: "Vendedor 1"},
		{ID: 2, CardNumberID: "002", FirstName: "Comprador 2", LastName: "Vendedor 2"},
		{ID: 3, CardNumberID: "003", FirstName: "Comprador 3", LastName: "Vendedor 3"},
	}

	//Act
	repo := buyer.MockRepository{
		Data: ListBuyers,
	}

	r := createServer(repo)
	req, recorder := createRequestTest(http.MethodGet, "/api/v1/buyers/aaa", "")
	r.ServeHTTP(recorder, req)
	//arrange
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestGetFailNotId(t *testing.T) {
	//arrange
	ListBuyers := []domain.Buyer{
		{ID: 1, CardNumberID: "001", FirstName: "Comprador 1", LastName: "Vendedor 1"},
		{ID: 2, CardNumberID: "002", FirstName: "Comprador 2", LastName: "Vendedor 2"},
		{ID: 3, CardNumberID: "003", FirstName: "Comprador 3", LastName: "Vendedor 3"},
	}

	//Act
	repo := buyer.MockRepository{
		Data: ListBuyers,
	}

	r := createServer(repo)
	req, recorder := createRequestTest(http.MethodGet, "/api/v1/buyers", "")
	r.ServeHTTP(recorder, req)
	//arrange
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestGetFailStatusInternalServerError(t *testing.T) {
	//arrange

	//Act
	repo := buyer.MockRepository{
		Data: nil,
	}

	r := createServer(repo)
	req, recorder := createRequestTest(http.MethodGet, "/api/v1/buyers/1", "")
	r.ServeHTTP(recorder, req)
	//arrange
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestGetFailNotFound(t *testing.T) {

	//arrange
	//id := 15
	//arrange
	ListBuyers := []domain.Buyer{
		{ID: 1, CardNumberID: "001", FirstName: "Comprador 1", LastName: "Vendedor 1"},
		{ID: 2, CardNumberID: "002", FirstName: "Comprador 2", LastName: "Vendedor 2"},
		{ID: 3, CardNumberID: "003", FirstName: "Comprador 3", LastName: "Vendedor 3"},
	}

	//Act
	repo := buyer.MockRepository{
		Data: ListBuyers,
		Err:  buyer.ErrNotFound,
	}

	r := createServer(repo)
	req, recorder := createRequestTest(http.MethodGet, "/api/v1/buyers/15", "")
	r.ServeHTTP(recorder, req)
	//arrange
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestCreateSuccess(t *testing.T) {
	//arrange
	ListBuyers := []domain.Buyer{
		{ID: 1, CardNumberID: "001", FirstName: "Comprador 1", LastName: "Vendedor 1"},
		{ID: 2, CardNumberID: "002", FirstName: "Comprador 2", LastName: "Vendedor 2"},
		{ID: 3, CardNumberID: "003", FirstName: "Comprador 3", LastName: "Vendedor 3"},
	}

	//Act
	repo := buyer.MockRepository{
		Data: ListBuyers,
	}

	r := createServer(repo)
	req, recorder := createRequestTest(http.MethodPost, "/api/v1/buyers/", `{"card_number_id": "004", "first_name": "Comprador 4", "last_name": "Vendedor 4"}`)
	r.ServeHTTP(recorder, req)
	//arrange
	assert.Equal(t, http.StatusCreated, recorder.Code)
}

func TestCreateFailBuyerNil(t *testing.T) {
	//arrange
	ListBuyers := []domain.Buyer{
		{ID: 1, CardNumberID: "001", FirstName: "Comprador 1", LastName: "Vendedor 1"},
		{ID: 2, CardNumberID: "002", FirstName: "Comprador 2", LastName: "Vendedor 2"},
		{ID: 3, CardNumberID: "003", FirstName: "Comprador 3", LastName: "Vendedor 3"},
	}

	//Act
	repo := buyer.MockRepository{
		Data: ListBuyers,
	}

	r := createServer(repo)
	req, recorder := createRequestTest(http.MethodPost, "/api/v1/buyers/", "")
	r.ServeHTTP(recorder, req)
	//arrange
	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
}

func TestCreateFailIdAlreadyExists(t *testing.T) {
	//arrange
	ListBuyers := []domain.Buyer{
		{ID: 1, CardNumberID: "001", FirstName: "Comprador 1", LastName: "Vendedor 1"},
		{ID: 2, CardNumberID: "002", FirstName: "Comprador 2", LastName: "Vendedor 2"},
		{ID: 3, CardNumberID: "003", FirstName: "Comprador 3", LastName: "Vendedor 3"},
	}

	//Act
	repo := buyer.MockRepository{
		Data: ListBuyers,
	}

	r := createServer(repo)
	req, recorder := createRequestTest(http.MethodPost, "/api/v1/buyers/", `{"card_number_id": "002", "first_name": "Comprador 2", "last_name": "Vendedor 2"}`)
	r.ServeHTTP(recorder, req)
	//arrange
	assert.Equal(t, http.StatusConflict, recorder.Code)
}

func TestCreateFailStatusNotFound(t *testing.T) {
	//arrange
	ListBuyers := []domain.Buyer{
		{ID: 1, CardNumberID: "001", FirstName: "Comprador 1", LastName: "Vendedor 1"},
		{ID: 2, CardNumberID: "002", FirstName: "Comprador 2", LastName: "Vendedor 2"},
		{ID: 3, CardNumberID: "003", FirstName: "Comprador 3", LastName: "Vendedor 3"},
	}

	//Act
	repo := buyer.MockRepository{
		Data: ListBuyers,
		Err:  buyer.ErrNotFound,
	}

	r := createServer(repo)
	req, recorder := createRequestTest(http.MethodPost, "/api/v1/buyers/", `{"card_number_id": "004", "first_name": "Comprador 4", "last_name": "Vendedor 4"}`)
	r.ServeHTTP(recorder, req)
	//arrange
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestCreateFailInternalServerError(t *testing.T) {
	//arrange
	ListBuyers := []domain.Buyer{
		{ID: 1, CardNumberID: "001", FirstName: "Comprador 1", LastName: "Vendedor 1"},
		{ID: 2, CardNumberID: "002", FirstName: "Comprador 2", LastName: "Vendedor 2"},
		{ID: 3, CardNumberID: "003", FirstName: "Comprador 3", LastName: "Vendedor 3"},
	}

	//Act
	repo := buyer.MockRepository{
		Data: ListBuyers,
		Err:  buyer.ErrInternal,
	}

	r := createServer(repo)
	req, recorder := createRequestTest(http.MethodPost, "/api/v1/buyers/", `{"card_number_id": "004", "first_name": "Comprador 4", "last_name": "Vendedor 4"}`)
	r.ServeHTTP(recorder, req)
	//arrange
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestDeleteSuccess(t *testing.T) {
	//arrange
	ListBuyers := []domain.Buyer{
		{ID: 1, CardNumberID: "001", FirstName: "Comprador 1", LastName: "Vendedor 1"},
		{ID: 2, CardNumberID: "002", FirstName: "Comprador 2", LastName: "Vendedor 2"},
		{ID: 3, CardNumberID: "003", FirstName: "Comprador 3", LastName: "Vendedor 3"},
	}

	//Act
	repo := buyer.MockRepository{
		Data: ListBuyers,
	}

	r := createServer(repo)
	req, recorder := createRequestTest(http.MethodDelete, "/api/v1/buyers/2", "")
	r.ServeHTTP(recorder, req)
	//arrange
	assert.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestDeleteFailBadRequest(t *testing.T) {
	//arrange
	ListBuyers := []domain.Buyer{
		{ID: 1, CardNumberID: "001", FirstName: "Comprador 1", LastName: "Vendedor 1"},
		{ID: 2, CardNumberID: "002", FirstName: "Comprador 2", LastName: "Vendedor 2"},
		{ID: 3, CardNumberID: "003", FirstName: "Comprador 3", LastName: "Vendedor 3"},
	}

	//Act
	repo := buyer.MockRepository{
		Data: ListBuyers,
	}

	r := createServer(repo)
	req, recorder := createRequestTest(http.MethodDelete, "/api/v1/buyers/aaa", "")
	r.ServeHTTP(recorder, req)
	//arrange
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestDeleteFailStatusNotFoundNotId(t *testing.T) {
	//arrange
	ListBuyers := []domain.Buyer{
		{ID: 1, CardNumberID: "001", FirstName: "Comprador 1", LastName: "Vendedor 1"},
		{ID: 2, CardNumberID: "002", FirstName: "Comprador 2", LastName: "Vendedor 2"},
		{ID: 3, CardNumberID: "003", FirstName: "Comprador 3", LastName: "Vendedor 3"},
	}

	//Act
	repo := buyer.MockRepository{
		Data: ListBuyers,
	}

	r := createServer(repo)
	req, recorder := createRequestTest(http.MethodDelete, "/api/v1/buyers", "")
	r.ServeHTTP(recorder, req)
	//arrange
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestDeleteFailStatusNotFound(t *testing.T) {
	//arrange
	ListBuyers := []domain.Buyer{
		{ID: 1, CardNumberID: "001", FirstName: "Comprador 1", LastName: "Vendedor 1"},
		{ID: 2, CardNumberID: "002", FirstName: "Comprador 2", LastName: "Vendedor 2"},
		{ID: 3, CardNumberID: "003", FirstName: "Comprador 3", LastName: "Vendedor 3"},
	}

	//Act
	repo := buyer.MockRepository{
		Data: ListBuyers,
		Err:  buyer.ErrNotFound,
	}

	r := createServer(repo)
	req, recorder := createRequestTest(http.MethodDelete, "/api/v1/buyers/1", "")
	r.ServeHTTP(recorder, req)
	//arrange
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestDeleteFailInternalServerError(t *testing.T) {
	//arrange
	ListBuyers := []domain.Buyer{
		{ID: 1, CardNumberID: "001", FirstName: "Comprador 1", LastName: "Vendedor 1"},
		{ID: 2, CardNumberID: "002", FirstName: "Comprador 2", LastName: "Vendedor 2"},
		{ID: 3, CardNumberID: "003", FirstName: "Comprador 3", LastName: "Vendedor 3"},
	}

	//Act
	repo := buyer.MockRepository{
		Data: ListBuyers,
		Err:  buyer.ErrInternal,
	}

	r := createServer(repo)
	req, recorder := createRequestTest(http.MethodDelete, "/api/v1/buyers/2", "")
	r.ServeHTTP(recorder, req)
	//arrange
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestUpdateSuccess(t *testing.T) {
	//arrange
	ListBuyers := []domain.Buyer{
		{ID: 1, CardNumberID: "001", FirstName: "Comprador 1", LastName: "Vendedor 1"},
		{ID: 2, CardNumberID: "002", FirstName: "Comprador 2", LastName: "Vendedor 2"},
		{ID: 3, CardNumberID: "003", FirstName: "Comprador 3", LastName: "Vendedor 3"},
	}

	//Act
	repo := buyer.MockRepository{
		Data: ListBuyers,
	}

	r := createServer(repo)
	req, recorder := createRequestTest(http.MethodPatch, "/api/v1/buyers/2", `{"card_number_id": "004", "first_name": "Comprador 4", "last_name": "Vendedor 4"}`)
	r.ServeHTTP(recorder, req)
	//arrange
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestUpdateFailStatusBadRequest(t *testing.T) {
	//arrange
	ListBuyers := []domain.Buyer{
		{ID: 1, CardNumberID: "001", FirstName: "Comprador 1", LastName: "Vendedor 1"},
		{ID: 2, CardNumberID: "002", FirstName: "Comprador 2", LastName: "Vendedor 2"},
		{ID: 3, CardNumberID: "003", FirstName: "Comprador 3", LastName: "Vendedor 3"},
	}

	//Act
	repo := buyer.MockRepository{
		Data: ListBuyers,
	}

	r := createServer(repo)
	req, recorder := createRequestTest(http.MethodPatch, "/api/v1/buyers/aa", `{"card_number_id": "004", "first_name": "Comprador 4", "last_name": "Vendedor 4"}`)
	r.ServeHTTP(recorder, req)
	//arrange
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestUpdateFailStatusUnprocessableEntity(t *testing.T) {
	//arrange
	ListBuyers := []domain.Buyer{
		{ID: 1, CardNumberID: "001", FirstName: "Comprador 1", LastName: "Vendedor 1"},
		{ID: 2, CardNumberID: "002", FirstName: "Comprador 2", LastName: "Vendedor 2"},
		{ID: 3, CardNumberID: "003", FirstName: "Comprador 3", LastName: "Vendedor 3"},
	}

	//Act
	repo := buyer.MockRepository{
		Data: ListBuyers,
	}

	r := createServer(repo)
	req, recorder := createRequestTest(http.MethodPatch, "/api/v1/buyers/1", `{"card_number_id": 004}`)
	r.ServeHTTP(recorder, req)
	//arrange
	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
}

func TestUpdateFailStatusConflict(t *testing.T) {
	//arrange
	ListBuyers := []domain.Buyer{
		{ID: 1, CardNumberID: "001", FirstName: "Comprador 1", LastName: "Vendedor 1"},
		{ID: 2, CardNumberID: "002", FirstName: "Comprador 2", LastName: "Vendedor 2"},
		{ID: 3, CardNumberID: "003", FirstName: "Comprador 3", LastName: "Vendedor 3"},
	}

	//Act
	repo := buyer.MockRepository{
		Data: ListBuyers,
		Err:  buyer.ErrAlreadyExists,
	}

	r := createServer(repo)
	req, recorder := createRequestTest(http.MethodPatch, "/api/v1/buyers/1", `{"card_number_id": "004", "first_name": "Comprador 4", "last_name": "Vendedor 4"}`)
	r.ServeHTTP(recorder, req)
	//arrange
	assert.Equal(t, http.StatusConflict, recorder.Code)
}

func TestUpdateFailStatusNotFound(t *testing.T) {
	//arrange
	ListBuyers := []domain.Buyer{
		{ID: 1, CardNumberID: "001", FirstName: "Comprador 1", LastName: "Vendedor 1"},
		{ID: 2, CardNumberID: "002", FirstName: "Comprador 2", LastName: "Vendedor 2"},
		{ID: 3, CardNumberID: "003", FirstName: "Comprador 3", LastName: "Vendedor 3"},
	}

	//Act
	repo := buyer.MockRepository{
		Data: ListBuyers,
		Err:  buyer.ErrNotFound,
	}

	r := createServer(repo)
	req, recorder := createRequestTest(http.MethodPatch, "/api/v1/buyers/1", `{"card_number_id": "004", "first_name": "Comprador 4", "last_name": "Vendedor 4"}`)
	r.ServeHTTP(recorder, req)
	//arrange
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestUpdateFailStatusInternalServerError(t *testing.T) {
	//arrange
	ListBuyers := []domain.Buyer{
		{ID: 1, CardNumberID: "001", FirstName: "Comprador 1", LastName: "Vendedor 1"},
		{ID: 2, CardNumberID: "002", FirstName: "Comprador 2", LastName: "Vendedor 2"},
		{ID: 3, CardNumberID: "003", FirstName: "Comprador 3", LastName: "Vendedor 3"},
	}

	//Act
	repo := buyer.MockRepository{
		Data: ListBuyers,
		Err:  errors.New(""),
	}

	r := createServer(repo)
	req, recorder := createRequestTest(http.MethodPatch, "/api/v1/buyers/1", `{"card_number_id": "004", "first_name": "Comprador 4", "last_name": "Vendedor 4"}`)
	r.ServeHTTP(recorder, req)
	//arrange
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}
