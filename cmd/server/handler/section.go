package handler

import (
	"net/http"
	"strconv"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/cmd/server/handler/requests"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/section"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/web"
	"github.com/gin-gonic/gin"
)

type Section struct {
	sectionService section.Service
}

func NewSection(s section.Service) *Section {
	return &Section{
		sectionService: s,
	}
}

// GetAll ListSections godoc
// @Summary     List sections
// @Tags        Sections
// @Description get sections
// @Produce     json
// @Success     200 {object} web.response
// @Failure     500 {object} web.errorResponse
// @Router      /sections [get]
func (s *Section) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := s.sectionService.GetAll(c)
		if err != nil {
			logging.Log(err)
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		//if data is empty replace null with empty slice
		if data == nil {
			web.Success(c, http.StatusOK, []domain.Section{})
			return
		}
		web.Success(c, http.StatusOK, data)
	}
}

// Get GetSectionByID godoc
// @Summary     Get section by ID
// @Tags        Sections
// @Description get section by ID
// @Produce     json
// @Param       id  path     int true "section id"
// @Success     200 {object} web.response
// @Failure     400 {object} web.errorResponse
// @Failure     404 {object} web.errorResponse
// @Failure     500 {object} web.errorResponse
// @Router      /sections/{id} [get]
func (s *Section) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		sectionId, err := strconv.Atoi(id)
		if err != nil {
			logging.Log(err)
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		data, err := s.sectionService.Get(c, sectionId)
		if err != nil {
			switch err {
			//check if error comes from the section not existing in the database
			case section.ErrNotFound:
				web.Error(c, http.StatusNotFound, "The section with id %d does not exists", sectionId)
			default:
				web.Error(c, http.StatusInternalServerError, err.Error())
			}
			logging.Log(err)
			return
		}
		web.Success(c, http.StatusOK, data)
	}
}

// Create CreateSection godoc
// @Summary     Create section
// @Tags        Sections
// @Description create section
// @Produce     json
// @Param       section body     requests.PostSection true "Section to store"
// @Success     201     {object} web.response
// @Failure     409     {object} web.errorResponse
// @Failure     422     {object} web.errorResponse
// @Failure     500     {object} web.errorResponse
// @Router      /sections [post]
func (s *Section) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req requests.PostSection
		if err := c.ShouldBindJSON(&req); err != nil {
			logging.Log(err)
			web.Error(c, http.StatusUnprocessableEntity, err.Error())
			return
		}
		sec, err := s.sectionService.Create(c, domain.Section{ID: 0, SectionNumber: req.SectionNumber, CurrentTemperature: *req.CurrentTemperature, MinimumTemperature: *req.MinimumTemperature, CurrentCapacity: req.CurrentCapacity, MinimumCapacity: req.MinimumCapacity, MaximumCapacity: req.MaximumCapacity, WarehouseID: req.WarehouseID, ProductTypeID: req.ProductTypeID})
		if err != nil {
			switch err {
			case section.ErrAlreadyExists:
				web.Error(c, http.StatusConflict, "a section with the section_number %d already exists", req.SectionNumber)
			default:
				web.Error(c, http.StatusInternalServerError, err.Error())
			}
			logging.Log(err)
			return
		}
		web.Success(c, http.StatusCreated, sec)
	}
}

// Update UpdateSection godoc
// @Summary     Update section
// @Tags        Sections
// @Description update section
// @Produce     json
// @Param       section body     requests.PatchSection true "Updated section"
// @Success     200     {object} web.response
// @Failure     400     {object} web.errorResponse
// @Failure     404     {object} web.errorResponse
// @Failure     409     {object} web.errorResponse
// @Failure     500     {object} web.errorResponse
// @Router      /sections/{id} [patch]
func (s *Section) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		sectionId, err := strconv.Atoi(id)
		if err != nil {
			logging.Log(err)
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		var req requests.PatchSection
		if err := c.ShouldBindJSON(&req); err != nil {
			logging.Log(err)
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		var currentTemperature, minimumTemperature int
		if req.CurrentTemperature == nil {
			currentTemperature = -273
		} else {
			currentTemperature = *req.CurrentTemperature
		}
		if req.MinimumTemperature == nil {
			minimumTemperature = -273
		} else {
			minimumTemperature = *req.MinimumTemperature
		}
		data, err := s.sectionService.Update(c, domain.Section{ID: sectionId, SectionNumber: req.SectionNumber, CurrentTemperature: currentTemperature, MinimumTemperature: minimumTemperature, CurrentCapacity: req.CurrentCapacity, MinimumCapacity: req.MinimumCapacity, MaximumCapacity: req.MaximumCapacity, WarehouseID: req.WarehouseID, ProductTypeID: req.ProductTypeID})
		if err != nil {
			switch err {
			case section.ErrNotFound:
				web.Error(c, http.StatusNotFound, "The section with id %d does not exists", sectionId)
			case section.ErrAlreadyExists:
				web.Error(c, http.StatusConflict, "a section with the section_number %d already exists", req.SectionNumber)
				return
			default:
				web.Error(c, http.StatusInternalServerError, err.Error())
			}
			logging.Log(err)
			return
		}
		web.Success(c, http.StatusOK, data)
	}
}

// Delete DeleteSection godoc
// @Summary     Delete section
// @Tags        Sections
// @Description delete section
// @Produce     json
// @Param       id path int true "section id"
// @Success     204
// @Failure     400 {object} web.errorResponse
// @Failure     404 {object} web.errorResponse
// @Failure     500 {object} web.errorResponse
// @Router      /sections/{id} [delete]
func (s *Section) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		sectionId, err := strconv.Atoi(id)
		if err != nil {
			logging.Log(err)
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		err = s.sectionService.Delete(c, sectionId)
		if err != nil {
			if err == section.ErrNotFound {
				logging.Log(err)
				web.Error(c, http.StatusNotFound, "The section with id %d does not exists", sectionId)
				return
			}
			logging.Log(err)
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		web.Success(c, http.StatusNoContent, "")
	}
}

// Get GetProductsBySection godoc
// @Summary     Get products by section
// @Tags        Sections
// @Description get products by section
// @Produce     json
// @Param       id  query    int false "section id"
// @Success     200 {object} web.response
// @Failure     400 {object} web.errorResponse
// @Failure     404 {object} web.errorResponse
// @Failure     500 {object} web.errorResponse
// @Router      /sections/reportProducts [get]
func (s *Section) GetSectionProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data []domain.ProductsBySection
		var err error
		var sectionID int

		id := c.Request.URL.Query().Get("id")
		if id != "" {
			sectionID, err = strconv.Atoi(id)
			if err != nil {
				logging.Log(err)
				web.Error(c, http.StatusBadRequest, err.Error())
				return
			}
			data, err = s.sectionService.GetSectionProducts(c, sectionID)
		} else {
			data, err = s.sectionService.GetSectionProducts(c, 0)
		}
		if err != nil {
			switch err {
			//check if error comes from the section not existing in the database
			case section.ErrNotFound:
				web.Error(c, http.StatusNotFound, "The section with id %d does not exists", sectionID)
			default:
				web.Error(c, http.StatusInternalServerError, err.Error())
			}
			logging.Log(err)
			return
		}
		web.Success(c, http.StatusOK, data)
	}
}
