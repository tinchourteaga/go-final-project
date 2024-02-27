package report_record

import (
	"errors"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func setupReportRecordServiceTest() (ctx *gin.Context) {
	logging.InitLog(nil)
	gin.SetMode("test")
	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	return
}

// TestService_Get_OK passes when id exists (return domain.ReportRecord and nil error)
func TestService_Get_OK(t *testing.T) {
	// Arrange
	searchID := 1

	// Act
	reportRecordRepository := RepositoryMock{db: []domain.ReportRecord{reportRecordTest}}
	reportRecordService := NewService(&reportRecordRepository)
	ctx := setupReportRecordServiceTest()
	reportRecord, errGet := reportRecordService.Get(ctx, &searchID)

	// Assert
	assert.NoError(t, errGet)
	assert.Equal(t, []domain.ReportRecord{reportRecordTest}, reportRecord)
	assert.True(t, reportRecordRepository.FlagGet)
}

// TestService_Get_FailNotFound passes when the given id is not in database (return empty domain.ReportRecord and error ServiceErrNotFound)
func TestService_Get_FailNotFound(t *testing.T) {
	// Arrange
	searchID := 1
	expectedErr := ServiceErrNotFound

	// Act
	reportRecordRepository := RepositoryMock{ForcedErrGet: RepositoryErrNotFound}
	reportRecordService := NewService(&reportRecordRepository)
	ctx := setupReportRecordServiceTest()
	reportRecord, errGet := reportRecordService.Get(ctx, &searchID)

	// Assert
	assert.EqualError(t, errGet, expectedErr.Error())
	assert.NotNil(t, reportRecord)
	assert.Empty(t, reportRecord)
	assert.True(t, reportRecordRepository.FlagGet)
}

// TestService_Get_FailInternal passes when there is a problem in repository layer (return empty domain.ReportRecord and error ServiceErrInternal)
func TestService_Get_FailInternal(t *testing.T) {
	// Arrange
	searchID := 1
	expectedErr := ServiceErrInternal

	// Act
	reportRecordRepository := RepositoryMock{ForcedErrGet: RepositoryErrInternal}
	reportRecordService := NewService(&reportRecordRepository)
	ctx := setupReportRecordServiceTest()
	reportRecord, errGet := reportRecordService.Get(ctx, &searchID)

	// Assert
	assert.EqualError(t, errGet, expectedErr.Error())
	assert.NotNil(t, reportRecord)
	assert.Empty(t, reportRecord)
	assert.True(t, reportRecordRepository.FlagGet)
}

// TestService_GetAll_OK passes when there are no errors on repository layer (return slice of all domain.ReportRecord and nil error)
func TestService_GetAll_OK(t *testing.T) {
	// Act
	reportRecordRepository := RepositoryMock{db: []domain.ReportRecord{reportRecordTest}}
	reportRecordService := NewService(&reportRecordRepository)
	ctx := setupReportRecordServiceTest()
	reportRecords, errGetAll := reportRecordService.Get(ctx, nil)

	// Assert
	assert.NoError(t, errGetAll)
	assert.Equal(t, reportRecordRepository.db, reportRecords)
	assert.True(t, reportRecordRepository.FlagGetAll)
}

// TestService_GetAll_OKEmpty passes when there are no errors on repository layer and database is empty (return empty slice of all domain.ReportRecord and nil error)
func TestService_GetAll_OKEmpty(t *testing.T) {
	// Act
	reportRecordRepository := RepositoryMock{}
	reportRecordService := NewService(&reportRecordRepository)
	ctx := setupReportRecordServiceTest()
	reportRecords, errGetAll := reportRecordService.Get(ctx, nil)

	// Assert
	assert.NoError(t, errGetAll)
	assert.NotNil(t, reportRecords)
	assert.Empty(t, reportRecords)
	assert.True(t, reportRecordRepository.FlagGetAll)
}

// TestService_GetAll_Fail passes when there is a problem in repository layer (return an empty slice of domain.ReportRecord and error ServiceErrInternal)
func TestService_GetAll_Fail(t *testing.T) {
	// Arrange
	expectedErr := ServiceErrInternal

	// Act
	reportRecordRepository := RepositoryMock{ForcedErrGetAll: errors.New("forced query error")}
	reportRecordService := NewService(&reportRecordRepository)
	ctx := setupReportRecordServiceTest()
	reportRecords, errGetAll := reportRecordService.Get(ctx, nil)

	// Assert
	assert.EqualError(t, errGetAll, expectedErr.Error())
	assert.NotNil(t, reportRecords)
	assert.Empty(t, reportRecords)
	assert.True(t, reportRecordRepository.FlagGetAll)
}
