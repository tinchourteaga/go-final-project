package handler

import (
	"net/http"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/cmd/server/handler/requests"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	productbatch "github.com/extmatperez/meli_bootcamp_go_w6-2/internal/productBatch"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/web"
	"github.com/gin-gonic/gin"
)

type ProductBatch struct {
	productBatchService productbatch.Service
}

func NewProductBatch(s productbatch.Service) *ProductBatch {
	return &ProductBatch{
		productBatchService: s,
	}
}

// Create CreateProductBatch godoc
// @Summary     Create product batch
// @Tags        Sections
// @Description create product batch
// @Produce     json
// @Param       productBatch body     requests.PostProductBatch true "Product batch to store"
// @Success     201          {object} web.response
// @Failure     409          {object} web.errorResponse
// @Failure     422          {object} web.errorResponse
// @Failure     500          {object} web.errorResponse
// @Router      /productBatches [post]
func (pb *ProductBatch) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req requests.PostProductBatch
		if err := c.ShouldBindJSON(&req); err != nil {
			logging.Log(err)
			web.Error(c, http.StatusUnprocessableEntity, err.Error())
			return
		}
		prodBatch, err := pb.productBatchService.Create(c, domain.ProductBatch{ID: 0, BatchNumber: req.BatchNumber, CurrentQuantity: req.CurrentQuantity, CurrentTemperature: req.CurrentTemperature, DueDate: req.DueDate, InitialQuantity: req.InitialQuantity, ManufacturingDate: req.ManufacturingDate, ManufacturingHour: req.ManufacturingHour, MinimumTemperature: req.MinimumTemperature, ProductID: req.ProductID, SectionID: req.SectionID})
		if err != nil {
			switch err {
			case productbatch.ErrAlreadyExists:
				web.Error(c, http.StatusConflict, "a product batch with the batch_number %d already exists", req.BatchNumber)
			case productbatch.ErrForeignProductNotFound, productbatch.ErrForeignSectionNotFound:
				web.Error(c, http.StatusConflict, err.Error())
			case productbatch.ErrDateValue:
				web.Error(c, http.StatusBadRequest, err.Error())
			default:
				web.Error(c, http.StatusInternalServerError, err.Error())
			}
			logging.Log(err)
			return
		}
		web.Success(c, http.StatusCreated, prodBatch)
	}
}
