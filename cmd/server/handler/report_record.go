package handler

import (
	"errors"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/record/report_record"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/web"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var (
	ReportRecordErrInvalidID = errors.New("invalid ID")
	ReportRecordErrNotFound  = errors.New("product not found")
)

type ReportRecord struct {
	reportRecordService report_record.Service
}

func NewReportRecord(service report_record.Service) *ReportRecord {
	return &ReportRecord{
		reportRecordService: service,
	}
}

// GetReportRecords
// @Summary     GET all ReportRecord or one ReportRecord by ID
// @Description Retrieves all ReportRecord or one ReportRecord by ID from database
// @Tags        ReportRecords
// @Produce     json
// @Param       product_id query    int               false "Product ID"
// @Success     200        {object} web.response      "Report Records"
// @Failure     400        {object} web.errorResponse "Invalid ID"
// @Failure     404        {object} web.errorResponse "Product not found"
// @Failure     500        {object} web.errorResponse "Unknown or unhandled error"
// @Router      /api/v1/products/reportRecords [get]
func (rr *ReportRecord) GetReportRecords() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idString := ctx.Query("id")
		var id *int
		if idString != "" {
			idConv, errStrConv := strconv.Atoi(idString)
			if errStrConv != nil {
				logging.Log(errStrConv)
				web.Error(ctx, http.StatusBadRequest, ReportRecordErrInvalidID.Error())
				return
			}
			id = &idConv
		}
		reports, errGet := rr.reportRecordService.Get(ctx, id)
		if errGet != nil {
			logging.Log(errGet)
			switch errGet {
			case report_record.ServiceErrNotFound:
				web.Error(ctx, http.StatusNotFound, ReportRecordErrNotFound.Error())
			default:
				// errorMessage = "" for security reasons (we don't want to expose internal data to the outside)
				web.Error(ctx, http.StatusInternalServerError, "")
			}
			return
		}
		web.Success(ctx, http.StatusOK, reports)
	}
}
