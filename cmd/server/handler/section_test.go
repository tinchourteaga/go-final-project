package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/section"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type responseSection struct {
	Data domain.Section
}
type responseSections struct {
	Data []domain.Section
}
type responseProductsBySection struct {
	Data []domain.ProductsBySection
}
type responseErrorSection struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

var sectionService section.MockService = section.MockService{
	MockSections: []domain.Section{},
	MockError:    nil,
}

var s = createSectionServer()

// createSectionServer creates the mock server to be tested upon
func createSectionServer() *gin.Engine {

	gin.SetMode(gin.TestMode)

	p := NewSection(&sectionService)

	r := gin.Default()

	sec := r.Group("/sections")

	sec.GET("", p.GetAll())
	sec.GET("/:id", p.Get())
	sec.POST("", p.Create())
	sec.PATCH("/:id", p.Update())
	sec.DELETE("/:id", p.Delete())
	sec.GET("/reportProducts", p.GetSectionProducts())

	return r
}

// TestSectionCreateOk tests if the hanlder correctly calls the sectionService to create in storage and return the given section
func TestSectionCreateOk(t *testing.T) {
	sectionService = section.MockService{
		MockSections: []domain.Section{},
		MockError:    nil,
	}
	body := `{"section_number":1,"current_temperature":-1,"minimum_temperature":-5,"current_capacity":1,"minimum_capacity":1,"maximum_capacity":1,"warehouse_id":1,"product_type_id":1}`
	req, rw := createRequestTest(http.MethodPost, "/sections", body)
	s.ServeHTTP(rw, req)

	expected := domain.Section{
		ID:                 1,
		SectionNumber:      1,
		CurrentTemperature: -1,
		MinimumTemperature: -5,
		CurrentCapacity:    1,
		MinimumCapacity:    1,
		MaximumCapacity:    1,
		WarehouseID:        1,
		ProductTypeID:      1,
	}

	var objRes responseSection
	assert.Equal(t, 201, rw.Code)
	err := json.Unmarshal(rw.Body.Bytes(), &objRes)

	assert.Nil(t, err)
	assert.Equal(t, expected, objRes.Data)
}

// TestSectionCreateFail tests if the handler returns the correct error when the body given is incorrect or incomplete
func TestSectionCreateFail(t *testing.T) {
	req, rw := createRequestTest(http.MethodPost, "/sections", "")
	s.ServeHTTP(rw, req)

	assert.Equal(t, 422, rw.Code)
}

// TestSectionCreateConflict tests if the handler returns the correct error when a section with the given section number already exists
func TestSectionCreateConflict(t *testing.T) {
	sectionService.MockError = section.ErrAlreadyExists
	body := `{"section_number":1,"current_temperature":-1,"minimum_temperature":-5,"current_capacity":1,"minimum_capacity":1,"maximum_capacity":1,"warehouse_id":1,"product_type_id":1}`
	req, rw := createRequestTest(http.MethodPost, "/sections", body)
	s.ServeHTTP(rw, req)

	assert.Equal(t, 409, rw.Code)
}

func TestSectionCreateInternalErr(t *testing.T) {
	sectionService.MockError = section.ErrInternal
	body := `{"section_number":1,"current_temperature":-1,"minimum_temperature":-5,"current_capacity":1,"minimum_capacity":1,"maximum_capacity":1,"warehouse_id":1,"product_type_id":1}`
	req, rw := createRequestTest(http.MethodPost, "/sections", body)
	s.ServeHTTP(rw, req)

	assert.Equal(t, 500, rw.Code)
}

// TestSectionFindAll tests if the hanlder correctly calls the sectionService to return al the sections stored
func TestSectionFindAll(t *testing.T) {
	sectionService = section.MockService{
		MockSections: []domain.Section{
			{
				ID:                 1,
				SectionNumber:      1,
				CurrentTemperature: -1,
				MinimumTemperature: -5,
				CurrentCapacity:    1,
				MinimumCapacity:    1,
				MaximumCapacity:    1,
				WarehouseID:        1,
				ProductTypeID:      1,
			},
		},
		MockError: nil,
	}
	req, rw := createRequestTest(http.MethodGet, "/sections", "")
	s.ServeHTTP(rw, req)

	expected := 1

	var objRes responseSections
	assert.Equal(t, 200, rw.Code)
	err := json.Unmarshal(rw.Body.Bytes(), &objRes)

	assert.Nil(t, err)
	assert.Equal(t, expected, len(objRes.Data))
}

func TestSectionFindAllEmpty(t *testing.T) {
	sectionService = section.MockService{
		MockSections: nil,
		MockError: nil,
	}
	req, rw := createRequestTest(http.MethodGet, "/sections", "")
	s.ServeHTTP(rw, req)

	expected := []domain.Section{}

	var objRes responseSections
	assert.Equal(t, 200, rw.Code)
	err := json.Unmarshal(rw.Body.Bytes(), &objRes)

	assert.Nil(t, err)
	assert.Equal(t, expected, objRes.Data)
}

func TestSectionFindAllInternalErr(t *testing.T) {
	sectionService = section.MockService{
		MockError: section.ErrInternal,
	}
	req, rw := createRequestTest(http.MethodGet, "/sections", "")
	s.ServeHTTP(rw, req)

	expected := section.ErrInternal

	var objRes responseErrorSection
	assert.Equal(t, 500, rw.Code)
	err := json.Unmarshal(rw.Body.Bytes(), &objRes)

	assert.Nil(t, err)
	assert.EqualError(t, expected, objRes.Message)
}

// TestSectionFindByInvalidId tests if the handler returns the correct error when the given id isn´t a valid decimal number
func TestSectionFindByInvalidId(t *testing.T) {
	req, rw := createRequestTest(http.MethodGet, "/sections/a", "")
	s.ServeHTTP(rw, req)

	assert.Equal(t, 400, rw.Code)
}

// TestSectionFindByIdNonExistent tests if the handler returns the correct error when a section with the given id doesn´t exist
func TestSectionFindByIdNonExistent(t *testing.T) {
	sectionService.MockError = section.ErrNotFound
	req, rw := createRequestTest(http.MethodGet, "/sections/1", "")
	s.ServeHTTP(rw, req)

	expected := "The section with id 1 does not exists"

	var objRes responseErrorSection
	assert.Equal(t, 404, rw.Code)
	err := json.Unmarshal(rw.Body.Bytes(), &objRes)

	assert.Nil(t, err)
	assert.Equal(t, expected, objRes.Message)
}

// TestSectionFindByIdExistent tests if the handler correctly calls the sectionService to return the section with the given id
func TestSectionFindByIdExistent(t *testing.T) {
	sectionService = section.MockService{
		MockSections: []domain.Section{
			{
				ID:                 1,
				SectionNumber:      1,
				CurrentTemperature: -1,
				MinimumTemperature: -5,
				CurrentCapacity:    1,
				MinimumCapacity:    1,
				MaximumCapacity:    1,
				WarehouseID:        1,
				ProductTypeID:      1,
			},
		},
		MockError: nil,
	}
	req, rw := createRequestTest(http.MethodGet, "/sections/1", "")
	s.ServeHTTP(rw, req)

	expected := domain.Section{
		ID:                 1,
		SectionNumber:      1,
		CurrentTemperature: -1,
		MinimumTemperature: -5,
		CurrentCapacity:    1,
		MinimumCapacity:    1,
		MaximumCapacity:    1,
		WarehouseID:        1,
		ProductTypeID:      1,
	}

	var objRes responseSection
	assert.Equal(t, 200, rw.Code)
	err := json.Unmarshal(rw.Body.Bytes(), &objRes)

	assert.Nil(t, err)
	assert.Equal(t, expected, objRes.Data)
}

func TestSectionFindByIdInternalErr(t *testing.T) {
	sectionService.MockError = section.ErrInternal
	req, rw := createRequestTest(http.MethodGet, "/sections/1", "")
	s.ServeHTTP(rw, req)

	expected := section.ErrInternal

	var objRes responseErrorSection
	assert.Equal(t, 500, rw.Code)
	err := json.Unmarshal(rw.Body.Bytes(), &objRes)

	assert.Nil(t, err)
	assert.EqualError(t, expected, objRes.Message)
}

// TestSectionUpdateOk tests if the handler correctly calls the sectionService to update the section with the given id both when a value is given for each attribute and when is not
func TestSectionUpdateOk(t *testing.T) {
	sectionService = section.MockService{
		MockSections: []domain.Section{
			{
				ID:                 1,
				SectionNumber:      1,
				CurrentTemperature: -1,
				MinimumTemperature: -5,
				CurrentCapacity:    1,
				MinimumCapacity:    1,
				MaximumCapacity:    1,
				WarehouseID:        1,
				ProductTypeID:      1,
			},
		},
		MockError: nil,
	}
	req1, rw1 := createRequestTest(http.MethodPatch, "/sections/1", "{}")
	s.ServeHTTP(rw1, req1)

	expected1 := domain.Section{
		ID:                 1,
		SectionNumber:      1,
		CurrentTemperature: -1,
		MinimumTemperature: -5,
		CurrentCapacity:    1,
		MinimumCapacity:    1,
		MaximumCapacity:    1,
		WarehouseID:        1,
		ProductTypeID:      1,
	}

	var objRes1 responseSection
	assert.Equal(t, 200, rw1.Code)
	err1 := json.Unmarshal(rw1.Body.Bytes(), &objRes1)

	assert.Nil(t, err1)
	assert.Equal(t, expected1, objRes1.Data)

	body := `{"section_number":2,"current_temperature":1,"minimum_temperature":-10,"current_capacity":2,"minimum_capacity":2,"maximum_capacity":2,"warehouse_id":2,"product_type_id":2}`
	req2, rw2 := createRequestTest(http.MethodPatch, "/sections/1", body)
	s.ServeHTTP(rw2, req2)

	expected2 := domain.Section{
		ID:                 1,
		SectionNumber:      2,
		CurrentTemperature: 1,
		MinimumTemperature: -10,
		CurrentCapacity:    2,
		MinimumCapacity:    2,
		MaximumCapacity:    2,
		WarehouseID:        2,
		ProductTypeID:      2,
	}
	var objRes2 responseSection
	assert.Equal(t, 200, rw2.Code)
	err2 := json.Unmarshal(rw2.Body.Bytes(), &objRes2)

	assert.Nil(t, err2)
	assert.Equal(t, expected2, objRes2.Data)
}

// TestSectionUpdateInvalidId tests if the handler returns the correct error when the given id isn´t a valid decimal number
func TestSectionUpdateInvalidId(t *testing.T) {
	req, rw := createRequestTest(http.MethodPatch, "/sections/a", "")
	s.ServeHTTP(rw, req)

	assert.Equal(t, 400, rw.Code)
}

// TestSectionUpdateFail tests if the handler returns the correct error when the given is incorrect
func TestSectionUpdateFail(t *testing.T) {
	req, rw := createRequestTest(http.MethodPatch, "/sections/1", "")
	s.ServeHTTP(rw, req)

	assert.Equal(t, 400, rw.Code)
}

// TestSectionUpdateExistentSectionNumber tests if the handler returns the correct error when a section with the given section number already exists
func TestSectionUpdateExistentSectionNumber(t *testing.T) {
	sectionService.MockError = section.ErrAlreadyExists
	req, rw := createRequestTest(http.MethodPatch, "/sections/1", `{"section_number":2}`)
	s.ServeHTTP(rw, req)

	expected := "a section with the section_number 2 already exists"

	var objRes responseErrorSection
	assert.Equal(t, 409, rw.Code)
	err := json.Unmarshal(rw.Body.Bytes(), &objRes)

	assert.Nil(t, err)
	assert.Equal(t, expected, objRes.Message)
}

// TestSectionUpdateNonExistent tests if the handler returns the correct error when a section with the given id doesn´t exist
func TestSectionUpdateNonExistent(t *testing.T) {
	sectionService.MockError = section.ErrNotFound
	req, rw := createRequestTest(http.MethodPatch, "/sections/1", "{}")
	s.ServeHTTP(rw, req)

	expected := "The section with id 1 does not exists"

	var objRes responseErrorSection
	assert.Equal(t, 404, rw.Code)
	err := json.Unmarshal(rw.Body.Bytes(), &objRes)

	assert.Nil(t, err)
	assert.Equal(t, expected, objRes.Message)
}

func TestSectionUpdateInternalErr(t *testing.T) {
	sectionService.MockError = section.ErrInternal
	req, rw := createRequestTest(http.MethodPatch, "/sections/1", `{"section_number":2}`)
	s.ServeHTTP(rw, req)

	expected := section.ErrInternal

	var objRes responseErrorSection
	assert.Equal(t, 500, rw.Code)
	err := json.Unmarshal(rw.Body.Bytes(), &objRes)

	assert.Nil(t, err)
	assert.EqualError(t, expected, objRes.Message)
}

// TestSectionDeleteNonExistent tests if the handler returns the correct error when a section with the given id doesn´t exist
func TestSectionDeleteNonExistent(t *testing.T) {
	sectionService.MockError = section.ErrNotFound
	req, rw := createRequestTest(http.MethodDelete, "/sections/1", "")
	s.ServeHTTP(rw, req)

	expected := "The section with id 1 does not exists"

	var objRes responseErrorSection
	assert.Equal(t, 404, rw.Code)
	err := json.Unmarshal(rw.Body.Bytes(), &objRes)

	assert.Nil(t, err)
	assert.Equal(t, expected, objRes.Message)
}

func TestSectionDeleteInternalErr(t *testing.T) {
	sectionService.MockError = section.ErrInternal
	req, rw := createRequestTest(http.MethodDelete, "/sections/1", "")
	s.ServeHTTP(rw, req)

	expected := section.ErrInternal

	var objRes responseErrorSection
	assert.Equal(t, 500, rw.Code)
	err := json.Unmarshal(rw.Body.Bytes(), &objRes)

	assert.Nil(t, err)
	assert.EqualError(t, expected, objRes.Message)
}

// TestSectionDeleteExistent tests if the handler correctly calls the sectionService to delete the secion in storage with the given id
func TestSectionDeleteExistent(t *testing.T) {
	sectionService = section.MockService{
		MockSections: []domain.Section{
			{
				ID:                 1,
				SectionNumber:      1,
				CurrentTemperature: -1,
				MinimumTemperature: -5,
				CurrentCapacity:    1,
				MinimumCapacity:    1,
				MaximumCapacity:    1,
				WarehouseID:        1,
				ProductTypeID:      1,
			},
		},
		MockError: nil,
	}
	req, rw := createRequestTest(http.MethodDelete, "/sections/1", "")
	s.ServeHTTP(rw, req)

	expected := 0

	assert.Equal(t, 204, rw.Code)

	assert.Equal(t, expected, len(sectionService.MockSections))
}

// TestSectionDeleteInvalidId tests if the handler returns the correct error when the given id isn´t a valid decimal number
func TestSectionDeleteInvalidId(t *testing.T) {
	req, rw := createRequestTest(http.MethodDelete, "/sections/a", "")
	s.ServeHTTP(rw, req)

	assert.Equal(t, 400, rw.Code)
}

func TestSectionGetAllSectionProducts(t *testing.T) {
	sectionService = section.MockService{
		MockProductsBySection: []domain.ProductsBySection{
			{
				SectionID:     1,
				SectionNumber: 1,
				ProductsCount: 10,
			},
			{
				SectionID:     2,
				SectionNumber: 2,
				ProductsCount: 20,
			},
		},
		MockError: nil,
	}

	req, rw := createRequestTest(http.MethodGet, "/sections/reportProducts", "")
	s.ServeHTTP(rw, req)

	expected := []domain.ProductsBySection{
		{
			SectionID:     1,
			SectionNumber: 1,
			ProductsCount: 10,
		},
		{
			SectionID:     2,
			SectionNumber: 2,
			ProductsCount: 20,
		},
	}

	var objRes responseProductsBySection
	assert.Equal(t, 200, rw.Code)

	err := json.Unmarshal(rw.Body.Bytes(), &objRes)
	assert.Nil(t, err)
	assert.Equal(t, expected, objRes.Data)
}

func TestSectionGetSectionProducts(t *testing.T) {
	sectionService = section.MockService{
		MockProductsBySection: []domain.ProductsBySection{
			{
				SectionID:     1,
				SectionNumber: 1,
				ProductsCount: 10,
			},
		},
		MockError: nil,
	}

	req, rw := createRequestTest(http.MethodGet, "/sections/reportProducts?id=1", "")
	s.ServeHTTP(rw, req)

	expected := []domain.ProductsBySection{
		{
			SectionID:     1,
			SectionNumber: 1,
			ProductsCount: 10,
		},
	}

	var objRes responseProductsBySection
	assert.Equal(t, 200, rw.Code)

	err := json.Unmarshal(rw.Body.Bytes(), &objRes)
	assert.Nil(t, err)
	assert.Equal(t, expected, objRes.Data)
}

func TestSectionGetSectionProductsIdNotExistent(t *testing.T) {
	sectionService = section.MockService{
		MockError: section.ErrNotFound,
	}

	req, rw := createRequestTest(http.MethodGet, "/sections/reportProducts?id=1", "")
	s.ServeHTTP(rw, req)

	expected := "The section with id 1 does not exists"

	var objRes responseErrorSection
	assert.Equal(t, 404, rw.Code)

	err := json.Unmarshal(rw.Body.Bytes(), &objRes)
	assert.Nil(t, err)
	assert.Equal(t, expected, objRes.Message)
}

func TestSectionGetSectionProductsInvalidId(t *testing.T) {
	req, rw := createRequestTest(http.MethodGet, "/sections/reportProducts?id=a", "")
	s.ServeHTTP(rw, req)

	assert.Equal(t, 400, rw.Code)
}

func TestSectionGetSectionProductsInternalErr(t *testing.T) {
	sectionService = section.MockService{
		MockError: section.ErrInternal,
	}

	req, rw := createRequestTest(http.MethodGet, "/sections/reportProducts", "")
	s.ServeHTTP(rw, req)

	expected := section.ErrInternal

	var objRes responseErrorSection
	assert.Equal(t, 500, rw.Code)
	err := json.Unmarshal(rw.Body.Bytes(), &objRes)

	assert.Nil(t, err)
	assert.EqualError(t, expected, objRes.Message)
}