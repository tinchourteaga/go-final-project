package product_record

import (
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
	"time"
)

func setupProductRecordServiceTest() (ctx *gin.Context) {
	logging.InitLog(nil)
	gin.SetMode("test")
	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	return
}

// TestService_Get_OK passes when id exists (return domain.ProductRecord and nil error)
func TestService_Get_OK(t *testing.T) {
	// Arrange
	searchID := 1

	// Act
	productRecordRepository := RepositoryMock{db: []domain.ProductRecord{productRecordTest}}
	productRecordService := NewService(&productRecordRepository)
	ctx := setupProductRecordServiceTest()
	productRecord, errGet := productRecordService.Get(ctx, searchID)

	// Assert
	assert.NoError(t, errGet)
	assert.Equal(t, productRecordTest, productRecord)
	assert.True(t, productRecordRepository.FlagGet)
}

// TestService_Get_FailNotFound passes when the given id is not in database (return empty domain.ProductRecord and error ServiceErrNotFound)
func TestService_Get_FailNotFound(t *testing.T) {
	// Arrange
	searchID := 1
	expectedErr := ServiceErrNotFound

	// Act
	productRecordRepository := RepositoryMock{ForcedErrGet: RepositoryErrNotFound}
	productRecordService := NewService(&productRecordRepository)
	ctx := setupProductRecordServiceTest()
	productRecord, errGet := productRecordService.Get(ctx, searchID)

	// Assert
	assert.EqualError(t, errGet, expectedErr.Error())
	assert.NotNil(t, productRecord)
	assert.Empty(t, productRecord)
	assert.True(t, productRecordRepository.FlagGet)
}

// TestService_Get_FailInternal passes when there is a problem in repository layer (return empty domain.ProductRecord and error ServiceErrInternal)
func TestService_Get_FailInternal(t *testing.T) {
	// Arrange
	searchID := 1
	expectedErr := ServiceErrInternal

	// Act
	productRecordRepository := RepositoryMock{ForcedErrGet: RepositoryErrInternal}
	productRecordService := NewService(&productRecordRepository)
	ctx := setupProductRecordServiceTest()
	productRecord, errGet := productRecordService.Get(ctx, searchID)

	// Assert
	assert.EqualError(t, errGet, expectedErr.Error())
	assert.NotNil(t, productRecord)
	assert.Empty(t, productRecord)
	assert.True(t, productRecordRepository.FlagGet)
}

// TestService_Save_OK passes when data is correct (return domain.ProductRecord and nil error)
func TestService_Save_OK(t *testing.T) {
	// Act
	productRecordRepository := RepositoryMock{ExpectedID: 1}
	productRecordService := NewService(&productRecordRepository)
	ctx := setupProductRecordServiceTest()
	productRecord, errSave := productRecordService.Save(ctx, productRecordTest)

	// Assert
	assert.NoError(t, errSave)
	assert.Equal(t, productRecordTest, productRecord)
	assert.True(t, productRecordRepository.FlagSave)
	assert.True(t, productRecordRepository.FlagGet)
}

// TestService_Save_FailNotFound passes when product record is created but can't be found in database (return empty domain.ProductRecord and error ServiceErrNotFound)
func TestService_Save_FailNotFound(t *testing.T) {
	// Arrange
	expectedErr := ServiceErrNotFound

	// Act
	productRecordRepository := RepositoryMock{ForcedErrGet: RepositoryErrNotFound}
	productRecordService := NewService(&productRecordRepository)
	ctx := setupProductRecordServiceTest()
	productRecord, errSave := productRecordService.Save(ctx, productRecordTest)

	// Assert
	assert.EqualError(t, errSave, expectedErr.Error())
	assert.NotNil(t, productRecord)
	assert.Empty(t, productRecord)
	assert.True(t, productRecordRepository.FlagSave)
	assert.True(t, productRecordRepository.FlagGet)
}

// TestService_Save_FailForeignKey passes when Save fails because product is not found in database (return empty domain.ProductRecord and error ServiceErrForeignKeyNotFound)
func TestService_Save_FailForeignKey(t *testing.T) {
	// Arrange
	expectedErr := ServiceErrForeignKeyNotFound

	// Act
	productRecordRepository := RepositoryMock{ForcedErrSave: RepositoryErrForeignKeyConstraint}
	productRecordService := NewService(&productRecordRepository)
	ctx := setupProductRecordServiceTest()
	productRecord, errSave := productRecordService.Save(ctx, productRecordTest)

	// Assert
	assert.EqualError(t, errSave, expectedErr.Error())
	assert.NotNil(t, productRecord)
	assert.Empty(t, productRecord)
	assert.True(t, productRecordRepository.FlagSave)
	assert.False(t, productRecordRepository.FlagGet)
}

// TestService_Save_FailInternal passes when there is a problem in repository layer (return empty domain.ProductRecord and error ServiceErrInternal)
func TestService_Save_FailInternal(t *testing.T) {
	// Arrange
	expectedErr := ServiceErrInternal

	// Act
	productRecordRepository := RepositoryMock{ForcedErrSave: RepositoryErrInternal}
	productRecordService := NewService(&productRecordRepository)
	ctx := setupProductRecordServiceTest()
	productRecord, errSave := productRecordService.Save(ctx, productRecordTest)

	// Assert
	assert.EqualError(t, errSave, expectedErr.Error())
	assert.NotNil(t, productRecord)
	assert.Empty(t, productRecord)
	assert.True(t, productRecordRepository.FlagSave)
	assert.False(t, productRecordRepository.FlagGet)
}

// TestService_Save_FailBadDate passes when date is before today's date (return empty domain.ProductRecord and error ServiceErrDate)
func TestService_Save_FailBadDate(t *testing.T) {
	// Arrange
	date, errDate := time.Parse(domain.ISO8601, "2021-12-24")
	assert.NoError(t, errDate)
	productRecordTest.LastUpdateDate = domain.MySqlTime{Time: date}
	expectedErr := ServiceErrDate

	// Act
	productRecordRepository := RepositoryMock{}
	productRecordService := NewService(&productRecordRepository)
	ctx := setupProductRecordServiceTest()
	productRecord, errSave := productRecordService.Save(ctx, productRecordTest)

	// Assert
	assert.EqualError(t, errSave, expectedErr.Error())
	assert.NotNil(t, productRecord)
	assert.Empty(t, productRecord)
	assert.False(t, productRecordRepository.FlagSave)
}
