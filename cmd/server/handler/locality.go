package handler

import (
	"net/http"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/locality"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Locality struct {
	localityService locality.Service
}

func NewLocality(l locality.Service) *Locality {
	return &Locality{
		localityService: l,
	}
}

// Get GetReportSellers godoc
// @Summary     Get total sellers number
// @Tags        Localities
// @Description Get total sellers number of a given locality
// @Produce     json
// @Param       id  query    string false "locality id"
// @Success     200 {object} web.response
// @Failure     404 {object} web.errorResponse
// @Failure     500 {object} web.errorResponse
// @Router      /api/v1/localities/reportCarries [get]
func (l *Locality) GetReportCarries() gin.HandlerFunc {
	return func(c *gin.Context) {
		var localityID *string
		queryLocalityID := c.Query("id")
		if queryLocalityID == "" {
			localityID = nil
		} else {
			localityID = &queryLocalityID
		}
		report, err := l.localityService.ReportCarries(c, localityID)
		if err != nil {
			switch err {
			case locality.ErrNotFound:
				web.Error(c, http.StatusNotFound, err.Error())
			default:
				web.Error(c, http.StatusInternalServerError, err.Error())
			}
			return
		}
		web.Success(c, http.StatusOK, report)
	}
}

// @Router /api/v1/localities/reportSellers [get]
func (l *Locality) GetReportSellers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var localityId *string

		queryLocalityID := c.Query("id")
		if queryLocalityID == "" {
			localityId = nil
		} else {
			localityId = &queryLocalityID
		}

		reportSellers, err := l.localityService.ReportSellers(c, localityId)
		if err != nil {
			logging.Log(err)
			switch err {
			case locality.ErrNotFound:
				web.Error(c, http.StatusNotFound, err.Error())
			default:
				web.Error(c, http.StatusInternalServerError, err.Error())
			}
			return
		}
		web.Success(c, http.StatusOK, reportSellers)
	}
}

// Get Locality by id godoc
// @Summary     Locality by id
// @Tags        Localities
// @Description get locality
// @Produce     json
// @Param       id  path     int               true "locality id"
// @Success     200 {object} web.response      "Get locality"
// @Failure     400 {object} web.errorResponse "BadRequest"
// @Failure     404 {object} web.errorResponse "Not found"
// @Failure     500 {object} web.errorResponse "Internal server error"
// @Router      /api/v1/localities/{id} [get]
func (l *Locality) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		localityId := c.Param("id")

		localityObtained, err := l.localityService.Get(c, localityId)
		if err != nil {
			logging.Log(err)
			switch err.Error() {
			case locality.ErrNotFound.Error():
				web.Error(c, http.StatusNotFound, "Id %s does not exist", localityId)
			default:
				web.Error(c, http.StatusInternalServerError, err.Error())
			}
			return
		}

		web.Success(c, http.StatusOK, localityObtained)
	}
}

// Create locality godoc
// @Summary Create Locality
// @Tags    Localities
// @Accept  json
// @Produce json
// @Param   locality body     domain.Locality   true "Locality to Create"
// @Success 201      {object} web.response      "New locality"
// @Failure 400      {object} web.errorResponse "BadRequest"
// @Failure 409      {object} web.errorResponse "Conflict"
// @Failure 422      {object} web.errorResponse "UnprocessableEntity"
// @Failure 500      {object} web.errorResponse "Internal server error"
// @Router  /api/v1/localities    [POST]
func (l *Locality) Create() gin.HandlerFunc {
	return func(c *gin.Context) {

		var req domain.Locality
		if err := c.ShouldBindJSON(&req); err != nil {
			switch err.(type) {
			case validator.ValidationErrors:
				web.Error(c, http.StatusBadRequest, "Bad Request, missing required fields")
			default:
				web.Error(c, http.StatusUnprocessableEntity, err.Error())
			}

			return
		}

		localityCreated, err := l.localityService.Create(c, domain.Locality{ID: req.ID, LocalityName: req.LocalityName, ProvinceName: req.ProvinceName, CountryName: req.CountryName})
		if err != nil {
			switch err {
			case locality.ErrAlreadyExists:
				web.Error(c, http.StatusConflict, err.Error())
			default:
				web.Error(c, http.StatusInternalServerError, err.Error())
			}
			return
		}

		web.Success(c, http.StatusCreated, localityCreated)
	}
}
