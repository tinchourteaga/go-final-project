package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/cmd/server/handler/requests"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/buyer"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/web"
	"github.com/gin-gonic/gin"
)

type Buyer struct {
	buyerService buyer.Service
}

func NewBuyer(s buyer.Service) *Buyer {
	return &Buyer{
		buyerService: s,
	}
}

// Get GetBuyerByID godoc
// @Summary     Get buyer by ID
// @Tags        Buyers
// @Description get buyer by ID
// @Produce     json
// @Param       id  path     int true "buyer id"
// @Success     200 {object} web.response
// @Failure     400 {object} web.errorResponse
// @Failure     404 {object} web.errorResponse
// @Failure     500 {object} web.errorResponse
// @Router      /api/v1/buyers/{id} [get]
func (b *Buyer) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			logging.Log("invalid Id")
			web.Error(c, http.StatusBadRequest, "invalid Id")
			return
		}
		data, errGet := b.buyerService.Get(c, id)
		if errGet != nil {
			switch errGet {
			case buyer.ErrNotFound:
				logging.Log(errors.New(fmt.Sprintf("buyer with %d id not found", id)))
				web.Error(c, http.StatusNotFound, "buyer with %d id not found", id)
			default:
				logging.Log(errGet.Error())
				web.Error(c, http.StatusInternalServerError, errGet.Error())
			}
			return
		}
		web.Success(c, http.StatusOK, data)
	}
}

// GetAll ListBuyers godoc
// @Summary     List buyers
// @Tags        Buyers
// @Description get buyers
// @Produce     json
// @Success     200 {object} web.response
// @Failure     404 {object} web.errorResponse
// @Router      /api/v1/buyers [get]
func (b *Buyer) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := b.buyerService.GetAll(c)
		if err != nil {
			logging.Log(err.Error())
			web.Error(c, http.StatusInternalServerError, err.Error())
		}
		if data == nil {
			web.Success(c, http.StatusOK, []domain.Buyer{})
			return
		}
		web.Success(c, http.StatusOK, data)
	}
}

// Create CreateBuyer godoc
// @Summary     Create buyer
// @Tags        Buyers
// @Description create buyer
// @Produce     json
// @Param       buyer body     requests.RequestBuyerPost true "buyer to store"
// @Success     201   {object} web.response
// @Failure     404   {object} web.errorResponse
// @Failure     409   {object} web.errorResponse
// @Failure     422   {object} web.errorResponse
// @Failure     500   {object} web.errorResponse
// @Router      /api/v1/buyers [post]
func (b *Buyer) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req requests.RequestBuyerPost
		if err := c.ShouldBindJSON(&req); err != nil {
			logging.Log(err.Error())
			web.Error(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		b, errCreate := b.buyerService.Save(c, domain.Buyer(req))
		if errCreate != nil {
			switch errCreate {
			case buyer.ErrAlreadyExists:
				logging.Log(errCreate.Error())
				web.Error(c, http.StatusConflict, errCreate.Error())
			case buyer.ErrNotFound:
				logging.Log("buyer created with no errors but not found in database")
				web.Error(c, http.StatusNotFound, "buyer created with no errors but not found in database")
			default:
				logging.Log(errCreate.Error())
				web.Error(c, http.StatusInternalServerError, errCreate.Error())
			}
			return
		}
		web.Success(c, http.StatusCreated, b)
	}
}

// Update UpdateBuyer godoc
// @Summary     Update buyer
// @Tags        Buyers
// @Description update buyer
// @Produce     json
// @Param       id        path     int                        true "buyer id"
// @Param       warehouse body     requests.RequestBuyerPatch true "buyer to update"
// @Success     200       {object} web.response
// @Failure     404       {object} web.errorResponse
// @Failure     422       {object} web.errorResponse
// @Failure     500       {object} web.errorResponse
// @Router      /api/v1/buyers/{id} [patch]
func (b *Buyer) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			logging.Log("invalid Id")
			web.Error(c, http.StatusBadRequest, "invalid Id")
			return
		}

		var req requests.RequestBuyerPatch
		if err := c.ShouldBindJSON(&req); err != nil {
			logging.Log(err.Error())
			web.Error(c, http.StatusUnprocessableEntity, err.Error())
			return
		}
		req.ID = id

		dataUpdate, errUpdate := b.buyerService.Update(c, domain.Buyer(req))
		if errUpdate != nil {
			switch errUpdate {
			case buyer.ErrAlreadyExists:
				logging.Log(errUpdate.Error())
				web.Error(c, http.StatusConflict, errUpdate.Error())
			case buyer.ErrNotFound:
				logging.Log(errors.New(fmt.Sprintf("buyer with %d id not found", id)))
				web.Error(c, http.StatusNotFound, "buyer with id %d not found", id)
			default:
				logging.Log(errUpdate.Error())
				web.Error(c, http.StatusInternalServerError, errUpdate.Error())
			}
			return
		}
		web.Success(c, http.StatusOK, dataUpdate)
	}
}

// Delete DeleteBuyer godoc
// @Summary     Delete buyer
// @Tags        Buyers
// @Description delete buyer
// @Produce     json
// @Param       id path int true "buyer id"
// @Success     204
// @Failure     404 {object} web.errorResponse
// @Failure     422 {object} web.errorResponse
// @Failure     500 {object} web.errorResponse
// @Router      /api/v1/buyers/{id} [delete]
func (b *Buyer) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			logging.Log("invalid Id")
			web.Error(c, http.StatusBadRequest, "invalid Id")
			return
		}
		errDelete := b.buyerService.Delete(c, id)
		if errDelete != nil {
			switch errDelete {
			case buyer.ErrNotFound:
				logging.Log(errors.New(fmt.Sprintf("buyer with %d id not found", id)))
				web.Error(c, http.StatusNotFound, "buyer with id %d not found", id)
			default:
				logging.Log(errDelete.Error())
				web.Error(c, http.StatusInternalServerError, errDelete.Error())
			}
			return
		}
		web.Success(c, http.StatusNoContent, "Deleted ok")
	}
}
