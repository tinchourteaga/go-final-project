package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/record/report_record"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type successfulReportRecordSliceResponse struct {
	Data []domain.ReportRecord `json:"data"`
}

type unsuccessfulReportRecordResponse struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func setupReportRecordHandlersEngineMock() (ctx *gin.Context, responseRecorder *httptest.ResponseRecorder) {
	logging.InitLog(nil)
	gin.SetMode(gin.TestMode)
	responseRecorder = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(responseRecorder)
	return
}

// TestReportRecord_GetReportRecords_All_OK passes when query is successful (return 200 and slice of all domain.ReportRecord)
func TestReportRecord_GetReportRecords_All_OK(t *testing.T) {
	reportRecordRepository := []domain.ReportRecord{
		{
			ProductID: 1,
		},
		{
			ProductID: 2,
		},
	}
	expectedCode := http.StatusOK
	ctx, responseRecorder := setupProductHandlersEngineMock()
	reportRecordService := report_record.ServiceMock{ReportRecordRepository: reportRecordRepository}
	reportRecordHandler := NewReportRecord(&reportRecordService)
	reportRecordHandler.GetReportRecords()(ctx)
	var response successfulReportRecordSliceResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.True(t, reportRecordService.FlagGet)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, reportRecordRepository, response.Data)
}

// TestReportRecord_GetReportRecords_All_OKWithEmptyDB passes when query is successful but database is empty (return 200 and empty slice of all domain.ReportRecord)
func TestReportRecord_GetReportRecords_All_OKWithEmptyDB(t *testing.T) {
	expectedCode := http.StatusOK
	ctx, responseRecorder := setupProductHandlersEngineMock()
	reportRecordService := report_record.ServiceMock{ReportRecordRepository: []domain.ReportRecord{}}
	reportRecordHandler := NewReportRecord(&reportRecordService)
	reportRecordHandler.GetReportRecords()(ctx)
	var response successfulReportRecordSliceResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.True(t, reportRecordService.FlagGet)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, 0, len(response.Data))
}

// TestReportRecord_GetReportRecords_All_InternalServerError passes when query is unsuccessful (return 500 and error message "")
func TestReportRecord_GetReportRecords_All_InternalServerError(t *testing.T) {
	expectedCode := http.StatusInternalServerError
	expectedErr := errors.New("")
	ctx, responseRecorder := setupReportRecordHandlersEngineMock()
	reportRecordService := report_record.ServiceMock{ReportRecordRepository: []domain.ReportRecord{}, ForcedErrGet: report_record.ServiceErrInternal}
	reportRecordHandler := NewReportRecord(&reportRecordService)
	reportRecordHandler.GetReportRecords()(ctx)
	var response unsuccessfulReportRecordResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.True(t, reportRecordService.FlagGet)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr.Error(), response.Message)
}

// TestReportRecord_GetReportRecords_One_OK passes when id exists (return 200 and domain.Product with given id)
func TestReportRecord_GetReportRecords_One_OK(t *testing.T) {
	reportRecordRepository := []domain.ReportRecord{
		{
			ProductID:    1,
			Description:  "A product's description",
			RecordsCount: 5,
		},
		{
			ProductID: 2,
		},
	}
	searchID := 1
	expectedCode := http.StatusOK
	ctx, responseRecorder := setupReportRecordHandlersEngineMock()
	req := &http.Request{URL: &url.URL{}}
	query := req.URL.Query()
	query.Add("id", fmt.Sprintf("%d", searchID))
	req.URL.RawQuery = query.Encode()
	ctx.Request = req
	reportRecordService := report_record.ServiceMock{ReportRecordRepository: reportRecordRepository}
	reportRecordHandler := NewReportRecord(&reportRecordService)
	reportRecordHandler.GetReportRecords()(ctx)
	var response successfulReportRecordSliceResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.True(t, reportRecordService.FlagGet)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, []domain.ReportRecord{reportRecordRepository[0]}, response.Data)
}

// TestReportRecord_GetReportRecords_One_IDNonExistent passes when the given id is not in database (return 404 and error message)
func TestReportRecord_GetReportRecords_One_IDNonExistent(t *testing.T) {
	searchID := 1
	forcedErr := report_record.ServiceErrNotFound
	expectedCode := http.StatusNotFound
	expectedErr := ReportRecordErrNotFound
	ctx, responseRecorder := setupReportRecordHandlersEngineMock()
	req := &http.Request{URL: &url.URL{}}
	query := req.URL.Query()
	query.Add("id", fmt.Sprintf("%d", searchID))
	req.URL.RawQuery = query.Encode()
	ctx.Request = req
	reportRecordService := report_record.ServiceMock{ReportRecordRepository: []domain.ReportRecord{}, ForcedErrGet: forcedErr}
	reportRecordHandler := NewReportRecord(&reportRecordService)
	reportRecordHandler.GetReportRecords()(ctx)
	var response unsuccessfulReportRecordResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.True(t, reportRecordService.FlagGet)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr.Error(), response.Message)
}

// TestReportRecord_GetReportRecords_One_InvalidID passes when the given id is not a number (return 400 and error message)
func TestReportRecord_GetReportRecords_One_InvalidID(t *testing.T) {
	searchID := "badID"
	expectedCode := http.StatusBadRequest
	expectedErr := ReportRecordErrInvalidID
	ctx, responseRecorder := setupReportRecordHandlersEngineMock()
	req := &http.Request{URL: &url.URL{}}
	query := req.URL.Query()
	query.Add("id", searchID)
	req.URL.RawQuery = query.Encode()
	ctx.Request = req
	reportRecordService := report_record.ServiceMock{ReportRecordRepository: []domain.ReportRecord{}}
	reportRecordHandler := NewReportRecord(&reportRecordService)
	reportRecordHandler.GetReportRecords()(ctx)
	var response unsuccessfulReportRecordResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.False(t, reportRecordService.FlagGet)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr.Error(), response.Message)
}

// TestReportRecord_GetReportRecords_One_InternalServerError passes when unexpected error occurs (return 500 and error message)
func TestReportRecord_GetReportRecords_One_InternalServerError(t *testing.T) {
	searchID := 1
	expectedCode := http.StatusInternalServerError
	expectedErr := errors.New("")
	ctx, responseRecorder := setupProductHandlersEngineMock()
	req := &http.Request{URL: &url.URL{}}
	query := req.URL.Query()
	query.Add("id", fmt.Sprintf("%d", searchID))
	req.URL.RawQuery = query.Encode()
	ctx.Request = req
	reportRecordService := report_record.ServiceMock{ReportRecordRepository: []domain.ReportRecord{}, ForcedErrGet: report_record.ServiceErrInternal}
	reportRecordHandler := NewReportRecord(&reportRecordService)
	reportRecordHandler.GetReportRecords()(ctx)
	var response unsuccessfulReportRecordResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.True(t, reportRecordService.FlagGet)
	assert.Equal(t, expectedCode, responseRecorder.Code)
	assert.Equal(t, expectedErr.Error(), response.Message)
}
