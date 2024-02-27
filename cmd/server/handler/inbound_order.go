package handler

import (
	"net/http"
	"strconv"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/cmd/server/handler/requests"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/inbound_order"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/web"
	"github.com/gin-gonic/gin"
)

type InboundOrder struct {
	inboundOrderService inbound_order.Service
}

func NewInboundOrder(ib inbound_order.Service) *InboundOrder {
	return &InboundOrder{
		inboundOrderService: ib,
	}
}

// GetAllEmployeesInboundOrders godoc
// @Summary     List employees and their inbound orders
// @Tags        Inbound Orders
// @Description Lists all existing employees and their inbound orders from database
// @Produce     json
// @Success     200 {object} web.response      "List of employees and their inbound orders"
// @Failure     500 {object} web.errorResponse "Connection to database error"
// @Router      /api/v1/employees/reportInboundOrders [get]
func (inboundOrder *InboundOrder) GetAllEmployeesInboundOrders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		employees, err := inboundOrder.inboundOrderService.GetAllEmployeesInboundOrders(ctx)

		if err != nil {
			logging.Log(err)
			web.Error(ctx, http.StatusInternalServerError, err.Error())
			return
		}

		if employees == nil {
			web.Success(ctx, http.StatusOK, []domain.EmployeeWithInboundOrders{})
			return
		}

		web.Success(ctx, http.StatusOK, employees)
	}
}

// GetEmployeeInboundOrders godoc
// @Summary     Get employee by ID
// @Tags        InboundOrders
// @Description Retrieves existing employee by ID and its inbound orders from database
// @Produce     json
// @Param       id  path     int               true "Employee id"
// @Success     200 {object} web.response      "Employee with inbound orders"
// @Failure     400 {object} web.errorResponse "Invalid id type"
// @Failure     404 {object} web.errorResponse "Employee not found"
// @Failure     500 {object} web.errorResponse "Connection to database error"
// @Router      /api/v1/employees/reportInboundOrders/{id} [get]
func (inboundOrder *InboundOrder) GetEmployeeInboundOrders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			logging.Log("invalid ip")
			web.Error(ctx, http.StatusBadRequest, "invalid id")
			return
		}

		obtainedEmployee, err := inboundOrder.inboundOrderService.GetEmployeeInboundOrders(ctx, id)

		if err != nil {
			logging.Log(err)
			switch err {
			case inbound_order.ErrEmployeeWithInboundOrdersNotFound:
				web.Error(ctx, http.StatusNotFound, inbound_order.ErrEmployeeWithInboundOrdersNotFound.Error())
			default:
				web.Error(ctx, http.StatusInternalServerError, err.Error())
			}
			return
		}

		web.Success(ctx, http.StatusOK, obtainedEmployee)
	}
}

// Create godoc
// @Summary     Create inbound order
// @Tags        InboundOrders
// @Description Creates a new inbound order in database
// @Accept      json
// @Produce     json
// @Param       inboundOrder body     requests.InboundOrderDTOPOST true "Inbound order to be stored"
// @Success     201          {object} web.response                 "Inbound order created"
// @Failure     409          {object} web.errorResponse            "Inbound order with order number already exists error"
// @Failure     422          {object} web.errorResponse            "Missing field or type casting error"
// @Failure     500          {object} web.errorResponse            "Connection to database error"
// @Router      /api/v1/inboundOrders [post]
func (inboundOrder *InboundOrder) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req requests.InboundOrderDTOPOST

		if err := ctx.ShouldBindJSON(&req); err != nil {
			logging.Log(err)
			web.Error(ctx, http.StatusUnprocessableEntity, err.Error())
			return
		}

		newInboundOrder := domain.InboundOrder{OrderDate: *req.OrderDate, OrderNumber: *req.OrderNumber, EmployeeID: *req.EmployeeID, ProductBatchID: *req.ProductBatchID, WarehouseID: *req.WarehouseID}
		newInboundOrder, err := inboundOrder.inboundOrderService.Save(ctx, newInboundOrder)

		if err != nil {
			logging.Log(err)
			switch err {
			case inbound_order.ErrInboundOrderAlreadyExists:
				web.Error(ctx, http.StatusConflict, inbound_order.ErrInboundOrderAlreadyExists.Error())
			case inbound_order.ErrInboundOrderNotSaved:
				web.Error(ctx, http.StatusInternalServerError, inbound_order.ErrInboundOrderNotSaved.Error())
			case inbound_order.ErrEmptyOrderNumber:
				web.Error(ctx, http.StatusConflict, inbound_order.ErrEmptyOrderNumber.Error())
			case inbound_order.ErrEmployeeNonExistent:
				web.Error(ctx, http.StatusNotFound, inbound_order.ErrEmployeeNonExistent.Error())
			case inbound_order.ErrWarehouseNonExistent:
				web.Error(ctx, http.StatusNotFound, inbound_order.ErrWarehouseNonExistent.Error())
			case inbound_order.ErrProductBatchNonExistent:
				web.Error(ctx, http.StatusNotFound, inbound_order.ErrProductBatchNonExistent.Error())
			default:
				web.Error(ctx, http.StatusInternalServerError, err.Error())
			}
			return
		}

		web.Success(ctx, http.StatusCreated, newInboundOrder)
	}
}
