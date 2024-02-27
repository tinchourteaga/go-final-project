package handler

import (
	"encoding/json"
	"errors"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/cmd/server/handler/requests"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/record/product_record"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

var (
	ProductRecordErrCreatedButNotFound = errors.New("product created with no errors but not found in database")
	ProductRecordErrDate               = errors.New("input date cannot be less than today")
	ProductRecordErrInvalidDate        = errors.New("invalid input date")
)

type ProductRecord struct {
	productRecordService product_record.Service
}

func NewProductRecord(service product_record.Service) *ProductRecord {
	return &ProductRecord{
		productRecordService: service,
	}
}

// Create
// @Summary     POST "ProductRecord"
// @Description "Creates a new ProductRecord on the database"
// @Tags        ProductRecords
// @Accept      json
// @Produce     json
// @Param       product_record body     requests.ProductRecordPOSTRequest true "ProductRecord to be created"
// @Success     201            {object} web.response                      "ProductRecord created"
// @Failure     404            {object} web.errorResponse                 "ProductRecord created but not found in database"
// @Failure     409            {object} web.errorResponse                 "Product not found"
// @Failure     422            {object} web.errorResponse                 "Missing Field or Type casting error"
// @Failure     500            {object} web.errorResponse                 "Unknown or unhandled error"
// @Router      /api/v1/productRecords [post]
func (pr *ProductRecord) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var productRecordPOSTRequest requests.ProductRecordPOSTRequest
		if err := ctx.ShouldBindJSON(&productRecordPOSTRequest); err != nil {
			logging.Log(err)
			switch err.(type) {
			case validator.ValidationErrors:
				web.Error(ctx, http.StatusUnprocessableEntity, err.Error())
			case *json.UnmarshalTypeError:
				web.Error(ctx, http.StatusUnprocessableEntity, err.Error())
			default:
				// errorMessage = "" for security reasons (we don't want to expose internal data to the outside)
				web.Error(ctx, http.StatusInternalServerError, "")
			}
			return
		}
		productRecordRequested, errMap := productRecordPOSTRequest.MapToDomain()
		if errMap != nil {
			logging.Log(errMap)
			web.Error(ctx, http.StatusBadRequest, ProductRecordErrInvalidDate.Error())
			return
		}
		prod, errSave := pr.productRecordService.Save(ctx, productRecordRequested)
		if errSave != nil {
			logging.Log(errSave)
			switch errSave {
			case product_record.ServiceErrDate:
				web.Error(ctx, http.StatusConflict, ProductRecordErrDate.Error())
			case product_record.ServiceErrForeignKeyNotFound:
				web.Error(ctx, http.StatusConflict, errSave.Error())
			case product_record.ServiceErrNotFound:
				web.Error(ctx, http.StatusNotFound, ProductRecordErrCreatedButNotFound.Error())
			default:
				// errorMessage = "" for security reasons (we don't want to expose internal data to the outside)
				web.Error(ctx, http.StatusInternalServerError, "")
			}
			return
		}
		web.Success(ctx, http.StatusCreated, prod)
	}
}
