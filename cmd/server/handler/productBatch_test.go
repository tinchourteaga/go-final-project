package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/productBatch"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type responseProductBatch struct {
	Data domain.ProductBatch
}

var productBatchService productbatch.MockService = productbatch.MockService{
	MockProductBatches: []domain.ProductBatch{},
	MockError:          nil,
}

var pbS = createProductBatchServer()

// createSectionServer creates the mock server to be tested upon
func createProductBatchServer() *gin.Engine {

	gin.SetMode(gin.TestMode)

	p := NewProductBatch(&productBatchService)

	r := gin.Default()

	sec := r.Group("/productBatches")

	sec.POST("", p.Create())

	return r
}

func TestProductBatchCreateOk(t *testing.T) {
	productBatchService = productbatch.MockService{
		MockProductBatches: []domain.ProductBatch{},
		MockError:          nil,
	}
	body := `{"batch_number":1,"current_quantity": 1,"current_temperature": 1,"due_date": "1999-12-12","initial_quantity": 1,"manufacturing_date": "1999-12-12","manufacturing_hour": 1,"minimum_temperature": 1,"product_id": 1,"section_id": 1}`

	req, rw := createRequestTest(http.MethodPost, "/productBatches", body)
	pbS.ServeHTTP(rw, req)

	expected := domain.ProductBatch{
		ID:                 1,
		BatchNumber:        1,
		CurrentQuantity:    1,
		CurrentTemperature: 1,
		DueDate:            "1999-12-12",
		InitialQuantity:    1,
		ManufacturingDate:  "1999-12-12",
		ManufacturingHour:  1,
		MinimumTemperature: 1,
		ProductID:          1,
		SectionID:          1,
	}

	var objRes responseProductBatch
	assert.Equal(t, 201, rw.Code)
	err := json.Unmarshal(rw.Body.Bytes(), &objRes)

	assert.Nil(t, err)
	assert.Equal(t, expected, objRes.Data)
}

func TestProductBatchCreateFail(t *testing.T) {
	req, rw := createRequestTest(http.MethodPost, "/productBatches", "")
	pbS.ServeHTTP(rw, req)

	assert.Equal(t, 422, rw.Code)
}

func TestProductBatchCreateInvalidDate(t *testing.T) {
	productBatchService.MockError = productbatch.ErrDateValue
	body := `{"batch_number":1,"current_quantity": 1,"current_temperature": 1,"due_date": "kadmfdkls","initial_quantity": 1,"manufacturing_date": "llcsmd","manufacturing_hour": 1,"minimum_temperature": 1,"product_id": 1,"section_id": 1}`
	req, rw := createRequestTest(http.MethodPost, "/productBatches", body)
	pbS.ServeHTTP(rw, req)

	assert.Equal(t, 400, rw.Code)
}

func TestProductBatchCreateConflict(t *testing.T) {
	productBatchService.MockError = productbatch.ErrAlreadyExists
	body := `{"batch_number":4,"current_quantity": 1,"current_temperature": 1,"due_date": "1999-12-12","initial_quantity": 1,"manufacturing_date": "1999-12-12","manufacturing_hour": 1,"minimum_temperature": 1,"product_id": 1,"section_id": 1}`
	req, rw := createRequestTest(http.MethodPost, "/productBatches", body)
	pbS.ServeHTTP(rw, req)

	assert.Equal(t, 409, rw.Code)
}

func TestProductBatchCreateInternal(t *testing.T) {
	productBatchService.MockError = productbatch.ErrInternal
	body := `{"batch_number":4,"current_quantity": 1,"current_temperature": 1,"due_date": "1999-12-12","initial_quantity": 1,"manufacturing_date": "1999-12-12","manufacturing_hour": 1,"minimum_temperature": 1,"product_id": 1,"section_id": 1}`
	req, rw := createRequestTest(http.MethodPost, "/productBatches", body)
	pbS.ServeHTTP(rw, req)

	assert.Equal(t, 500, rw.Code)
}

func TestProductBatchCreateMissingForeignProduct(t *testing.T) {
	productBatchService.MockError = productbatch.ErrForeignProductNotFound
	body := `{"batch_number":4,"current_quantity": 1,"current_temperature": 1,"due_date": "1999-12-12","initial_quantity": 1,"manufacturing_date": "1999-12-12","manufacturing_hour": 1,"minimum_temperature": 1,"product_id": 2,"section_id": 1}`
	req, rw := createRequestTest(http.MethodPost, "/productBatches", body)
	pbS.ServeHTTP(rw, req)

	assert.Equal(t, 409, rw.Code)
}

func TestProductBatchCreateMissingForeignSection(t *testing.T) {
	productBatchService.MockError = productbatch.ErrForeignSectionNotFound
	body := `{"batch_number":4,"current_quantity": 1,"current_temperature": 1,"due_date": "1999-12-12","initial_quantity": 1,"manufacturing_date": "1999-12-12","manufacturing_hour": 1,"minimum_temperature": 1,"product_id": 1,"section_id": 2}`
	req, rw := createRequestTest(http.MethodPost, "/productBatches", body)
	pbS.ServeHTTP(rw, req)

	assert.Equal(t, 409, rw.Code)
}
