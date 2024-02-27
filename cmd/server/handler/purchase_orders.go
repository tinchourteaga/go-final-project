package handler

import (
	"net/http"
	"strconv"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/cmd/server/handler/requests"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	purchaseorders "github.com/extmatperez/meli_bootcamp_go_w6-2/internal/purchase_orders"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/web"
	"github.com/gin-gonic/gin"
)

type Purchase_Order struct {
	service purchaseorders.Service
}

func NewPurchaseOrders(s purchaseorders.Service) *Purchase_Order {
	return &Purchase_Order{
		service: s,
	}
}

// Create CreateOrder godoc
// @Summary     Create Purchase_Order
// @Tags        Purchase_Order
// @Description create Purchase_Order
// @Produce     json
// @Param       purchase_order body     requests.RequestPurchaseOrdersPost true "Purchase_Order to store"
// @Success     201            {object} web.response
// @Failure     404            {object} web.errorResponse
// @Failure     409            {object} web.errorResponse
// @Failure     422            {object} web.errorResponse
// @Failure     500            {object} web.errorResponse
// @Router      /api/v1/purchase_orders [post]
func (o *Purchase_Order) CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req requests.RequestPurchaseOrdersPost
		if err := c.ShouldBindJSON(&req); err != nil {
			logging.Log(purchaseorders.ErrBodyValidation.Error())
			web.Error(c, http.StatusUnprocessableEntity, purchaseorders.ErrBodyValidation.Error())
			return
		}

		po, errCreate := o.service.SaveOrder(c, domain.Purchase_orders(req))
		if errCreate != nil {
			switch errCreate.Error() {
			case purchaseorders.ErrAlreadyExists.Error():
				logging.Log(errCreate.Error())
				web.Error(c, http.StatusConflict, errCreate.Error())
			case purchaseorders.ErrFKConstraint.Error():
				logging.Log(errCreate.Error())
				web.Error(c, http.StatusConflict, errCreate.Error())
			case purchaseorders.ErrDataLong.Error():
				logging.Log(errCreate.Error())
				web.Error(c, http.StatusUnprocessableEntity, errCreate.Error())
			default:
				logging.Log(errCreate.Error())
				web.Error(c, http.StatusInternalServerError, errCreate.Error())
			}
			return
		}
		web.Success(c, http.StatusCreated, po)
	}
}

// GetAllOrdersByBuyers List Purchase_Order godoc
// @Summary     List Purchase_Order
// @Tags        Create Purchase_Order
// @Description get all Purchase_Order by buyer_id
// @Produce     json
// @Success     200 {object} web.response
// @Failure     404 {object} web.errorResponse
// @Router      /api/v1/reportPurchaseOrder [get]
func (o *Purchase_Order) GetAllOrdersByBuyers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data []domain.Purchase_orders_buyer
		var errGet error
		id := c.Query("id")
		if id != "" {
			idNum, errId := strconv.Atoi(id)
			if errId != nil {
				logging.Log("invalid Id")
				web.Error(c, http.StatusBadRequest, "invalid Id")
				return
			}
			data, errGet = o.service.GetAllByBuyer(c, idNum)

		} else {
			data, errGet = o.service.GetAllByBuyer(c, 0)
		}
		if errGet != nil {
			switch errGet {
			case purchaseorders.ErrNotFound:
				logging.Log(errGet.Error())
				web.Error(c, http.StatusNotFound, errGet.Error())
				return
			default:
				logging.Log(errGet.Error())
				web.Error(c, http.StatusInternalServerError, errGet.Error())
				return
			}
		}
		web.Success(c, http.StatusOK, data)
	}
}
