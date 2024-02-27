package handler

import (
	"net/http"
	"strconv"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/cmd/server/handler/requests"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/warehouse"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/web"
	"github.com/gin-gonic/gin"
)

type Warehouse struct {
	service warehouse.Service
}

func NewWarehouse(w warehouse.Service) *Warehouse {
	return &Warehouse{
		service: w,
	}
}

// Get GetWarehouseByID godoc
// @Summary     Get warehouse by ID
// @Tags        Warehouses
// @Description get warehouse by ID
// @Produce     json
// @Param       id  path     int true "warehouse id"
// @Success     200 {object} web.response
// @Failure     400 {object} web.errorResponse
// @Failure     404 {object} web.errorResponse
// @Failure     500 {object} web.errorResponse
// @Router      /api/v1/warehouses/{id} [get]
func (w *Warehouse) Get(ctx *gin.Context) {
	idString := ctx.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		logging.Log(warehouse.ErrBadRequest)
		web.Error(ctx, http.StatusBadRequest, warehouse.ErrBadRequest.Error())
		return
	}

	warehouseObtained, err := w.service.Get(ctx, id)
	if err != nil {
		switch err {
		case warehouse.ErrNotFound:
			logging.Log(warehouse.ErrNotFound)
			web.Error(ctx, http.StatusNotFound, warehouse.ErrNotFound.Error())
		default:
			logging.Log(err)
			web.Error(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}

	web.Success(ctx, http.StatusOK, warehouseObtained)
}

// GetAll ListWarehouses godoc
// @Summary     List warehouses
// @Tags        Warehouses
// @Description get warehouses
// @Produce     json
// @Success     200 {object} web.response
// @Failure     500 {object} web.errorResponse
// @Router      /api/v1/warehouses [get]
func (w *Warehouse) GetAll(ctx *gin.Context) {
	warehouses, err := w.service.GetAll(ctx)
	if err != nil {
		logging.Log(err)
		web.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	web.Success(ctx, http.StatusOK, warehouses)
}

// Create CreateWarehouse godoc
// @Summary     Create warehouse
// @Tags        Warehouses
// @Description create warehouse
// @Produce     json
// @Param       warehouse body     requests.WarehousePostRequest true "Warehouse to store"
// @Success     201       {object} web.response
// @Failure     409       {object} web.errorResponse
// @Failure     422       {object} web.errorResponse
// @Failure     500       {object} web.errorResponse
// @Router      /api/v1/warehouses [post]
func (w *Warehouse) Create(ctx *gin.Context) {
	var req requests.WarehousePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logging.Log(warehouse.ErrBodyValidation)
		web.Error(ctx, http.StatusUnprocessableEntity, warehouse.ErrBodyValidation.Error())
		return
	}

	warehouseCreated, err := w.service.Create(ctx, *req.Address, *req.Telephone, *req.WarehouseCode, *req.MinimumCapacity, *req.MinimumTemperature)
	if err != nil {
		switch err.Error() {
		case warehouse.ErrAlreadyExists.Error():
			logging.Log(warehouse.ErrAlreadyExists)
			web.Error(ctx, http.StatusConflict, err.Error())
		default:
			logging.Log(err)
			web.Error(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}

	web.Success(ctx, http.StatusCreated, warehouseCreated)
}

// Update UpdateWarehouse godoc
// @Summary     Update warehouse
// @Tags        Warehouses
// @Description update warehouse
// @Produce     json
// @Param       id        path     int                            true "warehouse id"
// @Param       warehouse body     requests.WarehousePatchRequest true "Warehouse to update"
// @Success     200       {object} web.response
// @Failure     400       {object} web.errorResponse
// @Failure     404       {object} web.errorResponse
// @Failure     409       {object} web.errorResponse
// @Failure     422       {object} web.errorResponse
// @Failure     500       {object} web.errorResponse
// @Router      /api/v1/warehouses/{id} [patch]
func (w *Warehouse) Update(ctx *gin.Context) {
	idString := ctx.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		logging.Log(warehouse.ErrBadRequest)
		web.Error(ctx, http.StatusBadRequest, warehouse.ErrBadRequest.Error())
		return
	}

	var req requests.WarehousePatchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logging.Log(warehouse.ErrBodyValidation)
		web.Error(ctx, http.StatusUnprocessableEntity, warehouse.ErrBodyValidation.Error())
		return
	}

	warehouseUpdated, err := w.service.Update(ctx, id, req.Address, req.Telephone, req.WarehouseCode, req.MinimumCapacity, req.MinimumTemperature)
	if err != nil {
		switch err {
		case warehouse.ErrAlreadyExists:
			logging.Log(warehouse.ErrAlreadyExists)
			web.Error(ctx, http.StatusConflict, err.Error())
		case warehouse.ErrNotFound:
			logging.Log(warehouse.ErrNotFound)
			web.Error(ctx, http.StatusNotFound, warehouse.ErrNotFound.Error())
		default:
			logging.Log(err)
			web.Error(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}
	web.Success(ctx, http.StatusOK, warehouseUpdated)
}

// Delete DeleteWarehouse godoc
// @Summary     Delete warehouse
// @Tags        Warehouses
// @Description delete warehouse
// @Produce     json
// @Param       id path int true "warehouse id"
// @Success     204
// @Failure     400 {object} web.errorResponse
// @Failure     404 {object} web.errorResponse
// @Failure     500 {object} web.errorResponse
// @Router      /api/v1/warehouses/{id} [delete]
func (w *Warehouse) Delete(ctx *gin.Context) {
	idString := ctx.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		logging.Log(warehouse.ErrBadRequest)
		web.Error(ctx, http.StatusBadRequest, warehouse.ErrBadRequest.Error())
		return
	}

	err = w.service.Delete(ctx, id)
	if err != nil {
		switch err {
		case warehouse.ErrNotFound:
			logging.Log(warehouse.ErrNotFound)
			web.Error(ctx, http.StatusNotFound, warehouse.ErrNotFound.Error())
		default:
			logging.Log(err)
			web.Error(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}

	web.Success(ctx, http.StatusNoContent, "")
}
