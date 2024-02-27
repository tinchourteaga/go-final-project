package handler

import (
	"net/http"
	"strconv"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/cmd/server/handler/requests"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/seller"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Seller struct {
	sellerService seller.Service
}

func NewSeller(s seller.Service) *Seller {
	return &Seller{
		sellerService: s,
	}
}

// GetAll ListSellers godoc
// @Summary     List sellers
// @Tags        Sellers
// @Description get sellers
// @Produce     json
// @Success     200 {object} web.response      "Get sellers"
// @Failure     500 {object} web.errorResponse "Internal server error"
// @Router      /api/v1/sellers [get]
func (s *Seller) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		seller, err := s.sellerService.GetAll(c)
		if err != nil {
			logging.Log(err)
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		web.Success(c, http.StatusOK, seller)
	}
}

// Get Seller by id godoc
// @Summary     Seller by id
// @Tags        Sellers
// @Description get seller
// @Produce     json
// @Param       id  path     int               true "seller id"
// @Success     200 {object} web.response      "Get seller"
// @Failure     400 {object} web.errorResponse "BadRequest"
// @Failure     404 {object} web.errorResponse "Not found"
// @Failure     500 {object} web.errorResponse "Internal server error"
// @Router      /api/v1/sellers/{id} [get]
func (s *Seller) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		sellerId, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			logging.Log(err)
			web.Error(c, http.StatusBadRequest, "Invalid ID")
			return
		}

		sellerObtained, err := s.sellerService.Get(c, int(sellerId))
		if err != nil {
			logging.Log(err)
			switch err.Error() {
			case seller.ErrNotFound.Error():
				web.Error(c, http.StatusNotFound, "Id %d does not exist", sellerId)
			default:
				web.Error(c, http.StatusInternalServerError, err.Error())
			}
			return
		}

		web.Success(c, http.StatusOK, sellerObtained)
	}
}

// Create seller godoc
// @Summary Create Seller
// @Tags    Sellers
// @Accept  json
// @Produce json
// @Param   seller body     requests.SellerPostRequest true "Seller to Create"
// @Success 201    {object} web.response               "New seller"
// @Failure 400    {object} web.errorResponse          "BadRequest"
// @Failure 409    {object} web.errorResponse          "Conflict"
// @Failure 422    {object} web.errorResponse          "UnprocessableEntity"
// @Failure 500    {object} web.errorResponse          "Internal server error"
// @Router  /api/v1/sellers    [POST]
func (s *Seller) Create() gin.HandlerFunc {
	return func(c *gin.Context) {

		var req requests.SellerPostRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logging.Log(err)
			switch err.(type) {
			case validator.ValidationErrors:
				web.Error(c, http.StatusBadRequest, "Bad Request, missing required fields")
			default:
				web.Error(c, http.StatusUnprocessableEntity, err.Error())
			}

			return
		}

		sellerCreated, err := s.sellerService.Create(c, domain.Seller{CID: *req.CID, CompanyName: *req.CompanyName, Address: *req.Address, Telephone: *req.Telephone, Locality_id: *req.Locality_id})
		if err != nil {
			logging.Log(err)
			switch err {
			case seller.ErrAlreadyExists:
				web.Error(c, http.StatusConflict, err.Error())
			case seller.ErrForeignKeyConstraint:
				web.Error(c, http.StatusConflict, err.Error())
			default:
				web.Error(c, http.StatusInternalServerError, err.Error())
			}
			return
		}

		web.Success(c, http.StatusCreated, sellerCreated)
	}
}

// Update godoc
// @Summary Update
// @Tags    Sellers
// @Accept  json
// @Produce json
// @Param   id     path     int                         true "seller id"
// @Param   seller body     requests.SellerPatchRequest true "seller"
// @Success 200    {object} web.response                "Seller"
// @Failure 400    {object} web.errorResponse           "BadRequest"
// @Failure 404    {object} web.errorResponse           "NotFound"
// @Failure 409    {object} web.errorResponse           "Conflict"
// @Failure 422    {object} web.errorResponse           "UnprocessableEntity"
// @Failure 500    {object} web.errorResponse           "Internal server error"
// @Router  /api/v1/sellers/{id} [PATCH]
func (s *Seller) Update() gin.HandlerFunc {
	return func(c *gin.Context) {

		sellerId, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			logging.Log(err)
			web.Error(c, http.StatusBadRequest, "Invalid ID")
			return
		}

		var req requests.SellerPatchRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logging.Log(err)
			web.Error(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		sellerUpdated, err := s.sellerService.Update(c, int(sellerId), req.CID, req.CompanyName, req.Address, req.Telephone, req.Locality_id)
		if err != nil {
			logging.Log(err)
			switch err {
			case seller.ErrNotFound:
				web.Error(c, http.StatusNotFound, "Id %d does not exist", sellerId)
			case seller.ErrAlreadyExists:
				web.Error(c, http.StatusConflict, err.Error())
			case seller.ErrForeignKeyConstraint:
				web.Error(c, http.StatusConflict, err.Error())
			default:
				web.Error(c, http.StatusInternalServerError, err.Error())
			}
			return
		}

		web.Success(c, http.StatusOK, sellerUpdated)
	}
}

// Delete seller
// @Summary Delete seller
// @Tags    Sellers
// @Param   id path int true "seller id"
// @Success 204
// @Failure 400 {object} web.errorResponse "BadRequest"
// @Failure 404 {object} web.errorResponse "Not found"
// @Failure 500 {object} web.errorResponse "Internal server error"
// @Router  /api/v1/sellers/{id}   [DELETE]
func (s *Seller) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		sellerId, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			logging.Log(err)
			web.Error(c, http.StatusBadRequest, "Invalid ID")
			return
		}

		err = s.sellerService.Delete(c, int(sellerId))
		if err != nil {
			logging.Log(err)
			switch err {
			case seller.ErrNotFound:
				web.Error(c, http.StatusNotFound, "Id %d does not exist", sellerId)
			default:
				web.Error(c, http.StatusInternalServerError, err.Error())
			}
			return
		}

		web.Success(c, http.StatusNoContent, "")
	}
}
