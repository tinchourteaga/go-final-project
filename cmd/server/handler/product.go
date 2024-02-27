package handler

import (
	"encoding/json"
	"errors"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/cmd/server/handler/requests"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/product"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/web"
	"github.com/gin-gonic/gin"
)

var (
	ProductErrInvalidID          = errors.New("invalid ID")
	ProductErrNotFound           = errors.New("product not found")
	ProductErrCreatedButNotFound = errors.New("product created with no errors but not found in database")
)

type Product struct {
	productService product.Service
}

func NewProduct(service product.Service) *Product {
	return &Product{
		productService: service,
	}
}

// GetAll
// @Summary     GET []Product
// @Description List of all Product from database
// @Tags        Products
// @Produce     json
// @Success     200 {object} web.response      "List of Products"
// @Failure     500 {object} web.errorResponse "Problems with the database"
// @Router      /api/v1/products [get]
func (p *Product) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		products, err := p.productService.GetAll(ctx)
		if err != nil {
			logging.Log(err)
			// errorMessage = "" for security reasons (we don't want to expose internal data to the outside)
			web.Error(ctx, http.StatusInternalServerError, "")
			return
		}
		if products == nil {
			web.Success(ctx, http.StatusOK, []domain.Product{})
			return
		}
		web.Success(ctx, http.StatusOK, products)
	}
}

// Get
// @Summary     GET Product by ID
// @Description Retrieves one Product from database by ID
// @Tags        Products
// @Produce     json
// @Param       id  path     int               true "Product ID"
// @Success     200 {object} web.response      "Product"
// @Failure     400 {object} web.errorResponse "Invalid ID"
// @Failure     404 {object} web.errorResponse "Product not found"
// @Failure     500 {object} web.errorResponse "Unknown or unhandled error"
// @Router      /api/v1/products/{id} [get]
func (p *Product) Get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idString := ctx.Param("id")
		id, errStrConv := strconv.Atoi(idString)
		if errStrConv != nil {
			logging.Log(errStrConv)
			web.Error(ctx, http.StatusBadRequest, ProductErrInvalidID.Error())
			return
		}
		prod, errGet := p.productService.Get(ctx, id)
		if errGet != nil {
			logging.Log(errGet)
			switch errGet {
			case product.ServiceErrNotFound:
				web.Error(ctx, http.StatusNotFound, ProductErrNotFound.Error())
			default:
				// errorMessage = "" for security reasons (we don't want to expose internal data to the outside)
				web.Error(ctx, http.StatusInternalServerError, "")
			}
			return
		}
		web.Success(ctx, http.StatusOK, prod)
	}
}

// Create
// @Summary     POST Product
// @Description Creates a new Product on the database
// @Tags        Products
// @Accept      json
// @Produce     json
// @Param       product body     requests.ProductPOSTRequest true "Product to be created"
// @Success     201     {object} web.response                "Product created"
// @Failure     400     {object} web.errorResponse           "Missing field"
// @Failure     404     {object} web.errorResponse           "Product created but not found in database or Seller not found"
// @Failure     409     {object} web.errorResponse           "Product code already exists"
// @Failure     422     {object} web.errorResponse           "Type casting error"
// @Failure     500     {object} web.errorResponse           "Unknown or unhandled error"
// @Router      /api/v1/products [post]
func (p *Product) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var productPOSTRequest requests.ProductPOSTRequest
		if err := ctx.ShouldBindJSON(&productPOSTRequest); err != nil {
			logging.Log(err)
			switch err.(type) {
			case validator.ValidationErrors:
				web.Error(ctx, http.StatusBadRequest, err.Error())
			case *json.UnmarshalTypeError:
				web.Error(ctx, http.StatusUnprocessableEntity, err.Error())
			default:
				// errorMessage = "" for security reasons (we don't want to expose internal data to the outside)
				web.Error(ctx, http.StatusInternalServerError, "")
			}
			return
		}
		productRequested := productPOSTRequest.MapToDomain()
		prod, errSave := p.productService.Save(ctx, productRequested)
		if errSave != nil {
			logging.Log(errSave)
			switch errSave {
			case product.ServiceErrAlreadyExists:
				web.Error(ctx, http.StatusConflict, errSave.Error())
			case product.ServiceErrForeignKeyNotFound:
				web.Error(ctx, http.StatusNotFound, errSave.Error())
			case product.ServiceErrNotFound:
				web.Error(ctx, http.StatusNotFound, ProductErrCreatedButNotFound.Error())
			default:
				// errorMessage = "" for security reasons (we don't want to expose internal data to the outside)
				web.Error(ctx, http.StatusInternalServerError, "")
			}
			return
		}
		web.Success(ctx, http.StatusCreated, prod)
	}
}

// PartialUpdate
// @Summary     PATCH Product by ID
// @Description Partially updates an existing Product on the database by ID
// @Tags        Products
// @Accept      json
// @Produce     json
// @Param       id      path     int                          true "Product ID"
// @Param       product body     requests.ProductPATCHRequest true "Product to be updated"
// @Success     200     {object} web.response                 "Product updated"
// @Failure     400     {object} web.errorResponse            "Invalid field or ID"
// @Failure     404     {object} web.errorResponse            "Product not found or Seller not found"
// @Failure     409     {object} web.errorResponse            "Product code already exists"
// @Failure     422     {object} web.errorResponse            "Type casting error"
// @Failure     500     {object} web.errorResponse            "Unknown or unhandled error"
// @Router      /api/v1/products/{id} [patch]
func (p *Product) PartialUpdate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idString := ctx.Param("id")
		id, errStrConv := strconv.Atoi(idString)
		if errStrConv != nil {
			logging.Log(errStrConv)
			web.Error(ctx, http.StatusBadRequest, ProductErrInvalidID.Error())
			return
		}
		var productPATCHRequest requests.ProductPATCHRequest
		if err := ctx.ShouldBindJSON(&productPATCHRequest); err != nil {
			logging.Log(err)
			switch err.(type) {
			case *json.UnmarshalTypeError:
				web.Error(ctx, http.StatusUnprocessableEntity, err.Error())
			default:
				// errorMessage = "" for security reasons (we don't want to expose internal data to the outside)
				web.Error(ctx, http.StatusInternalServerError, "")
			}
			return
		}
		productModificationRequested := productPATCHRequest.MapToDomain()
		prod, errPartialUpdate := p.productService.PartialUpdate(ctx, id, productModificationRequested)
		if errPartialUpdate != nil {
			logging.Log(errPartialUpdate)
			switch errPartialUpdate {
			case product.ServiceErrAlreadyExists:
				web.Error(ctx, http.StatusConflict, errPartialUpdate.Error())
			case product.ServiceErrForeignKeyNotFound:
				web.Error(ctx, http.StatusNotFound, errPartialUpdate.Error())
			case product.ServiceErrNotFound:
				web.Error(ctx, http.StatusNotFound, ProductErrNotFound.Error())
			default:
				// errorMessage = "" for security reasons (we don't want to expose internal data to the outside)
				web.Error(ctx, http.StatusInternalServerError, "")
			}
			return
		}
		web.Success(ctx, http.StatusOK, prod)
	}
}

// Delete
// @Summary     DELETE Product by ID
// @Description Remove a Product from the database by ID
// @Tags        Products
// @Produce     json
// @Success     204
// @Failure     400 {object} web.errorResponse "Invalid ID"
// @Failure     404 {object} web.errorResponse "Product not found"
// @Failure     500 {object} web.errorResponse "Unknown or unhandled error"
// @Router      /api/v1/products/{id} [delete]
func (p *Product) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idString := ctx.Param("id")
		id, errStrConv := strconv.Atoi(idString)
		if errStrConv != nil {
			logging.Log(errStrConv)
			web.Error(ctx, http.StatusBadRequest, ProductErrInvalidID.Error())
			return
		}
		errDelete := p.productService.Delete(ctx, id)
		if errDelete != nil {
			logging.Log(errDelete)
			switch errDelete {
			case product.ServiceErrNotFound:
				web.Error(ctx, http.StatusNotFound, ProductErrNotFound.Error())
			default:
				// errorMessage = "" for security reasons (we don't want to expose internal data to the outside)
				web.Error(ctx, http.StatusInternalServerError, "")
			}
			return
		}
		// Using ctx.Status(http.StatusNoContent) works on release mode, but on testing mode 204 is changed to 200 latter on
		web.Success(ctx, http.StatusNoContent, "")
	}
}
