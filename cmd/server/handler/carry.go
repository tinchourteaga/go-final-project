package handler

import (
	"net/http"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/cmd/server/handler/requests"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/carry"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/web"
	"github.com/gin-gonic/gin"
)

type Carry struct {
	service carry.Service
}

func NewCarry(c carry.Service) *Carry {
	return &Carry{
		service: c,
	}
}

// Save SaveCarry godoc
// @Summary     Create carry
// @Tags        Carries
// @Description create carry
// @Produce     json
// @Param       carry body     requests.CarryPostRequest true "Carry to save"
// @Success     201   {object} web.response
// @Failure     409   {object} web.errorResponse
// @Failure     422   {object} web.errorResponse
// @Failure     500   {object} web.errorResponse
// @Router      /api/v1/carries [post]
func (c *Carry) Save(ctx *gin.Context) {
	var req requests.CarryPostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logging.Log(carry.ErrBodyValidation)
		web.Error(ctx, http.StatusUnprocessableEntity, carry.ErrBodyValidation.Error())
		return
	}

	carryCreated, err := c.service.Save(ctx, *req.CID, *req.CompanyName, *req.Address, *req.Telephone, *req.Locality_id)
	if err != nil {
		switch err.Error() {
		case carry.ErrAlreadyExists.Error():
			logging.Log(carry.ErrAlreadyExists)
			web.Error(ctx, http.StatusConflict, err.Error())
		case carry.ErrFKConstraint.Error():
			logging.Log(carry.ErrFKConstraint)
			web.Error(ctx, http.StatusConflict, err.Error())
		case carry.ErrDataLong.Error():
			logging.Log(carry.ErrDataLong)
			web.Error(ctx, http.StatusUnprocessableEntity, err.Error())
		default:
			logging.Log(err)
			web.Error(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}

	web.Success(ctx, http.StatusCreated, carryCreated)
}
