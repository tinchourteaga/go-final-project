package report_record

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func init() {
	logging.InitLog(nil)
}

var reportRecordTest = domain.ReportRecord{
	ProductID:    1,
	Description:  "Samsung S21 FE",
	RecordsCount: 4,
}

// TestRepository_GetAll_OK passes when query is successful (return slice of all domain.ReportRecord and nil error)
func TestRepository_GetAll_OK(t *testing.T) {
	// Arrange
	columns := []string{"product_id", "description", "records_count"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(reportRecordTest.ProductID, reportRecordTest.Description, reportRecordTest.RecordsCount)

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mock.ExpectQuery(regexp.QuoteMeta(GetAllReportRecords)).WillReturnRows(rows)
	reportResult, errGetAll := repo.GetAll(ctx)

	// Assert
	assert.NoError(t, errGetAll)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, []domain.ReportRecord{reportRecordTest}, reportResult)
}

// TestRepository_GetAll_OKEmpty passes when query is successful and database is empty (return nil slice of domain.ReportRecord and nil error)
func TestRepository_GetAll_OKEmpty(t *testing.T) {
	// Arrange
	columns := []string{"product_id", "description", "records_count"}
	rows := sqlmock.NewRows(columns)

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mock.ExpectQuery(regexp.QuoteMeta(GetAllReportRecords)).WillReturnRows(rows)
	reportResult, errGetAll := repo.GetAll(ctx)

	// Assert
	assert.NoError(t, errGetAll)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NotNil(t, reportResult)
	assert.Empty(t, reportResult)
}

// TestRepository_GetAll_Fail passes when query fails (return nil and error message)
func TestRepository_GetAll_Fail(t *testing.T) {
	// Arrange
	expectedErr := errors.New("forced error on query")

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mock.ExpectQuery(regexp.QuoteMeta(GetAllReportRecords)).WillReturnError(expectedErr)
	reportResult, errGetAll := repo.GetAll(ctx)

	// Assert
	assert.EqualError(t, errGetAll, expectedErr.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NotNil(t, reportResult)
	assert.Empty(t, reportResult)
}

// TestRepository_Get_OK passes when query is successful (return domain.ReportRecord and nil error)
func TestRepository_Get_OK(t *testing.T) {
	// Arrange
	columns := []string{"product_id", "description", "records_count"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(reportRecordTest.ProductID, reportRecordTest.Description, reportRecordTest.RecordsCount)

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mock.ExpectQuery(regexp.QuoteMeta(GetReportRecord)).WithArgs(reportRecordTest.ProductID).WillReturnRows(rows)
	reportResult, errGet := repo.Get(ctx, reportRecordTest.ProductID)

	// Assert
	assert.NoError(t, errGet)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, reportRecordTest, reportResult)
}

// TestRepository_Get_FailNotFound passes when the given id is not in database (return empty domain.ReportRecord and error RepositoryErrNotFound)
func TestRepository_Get_FailNotFound(t *testing.T) {
	// Arrange
	expectedErr := RepositoryErrNotFound

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mock.ExpectQuery(regexp.QuoteMeta(GetReportRecord)).WithArgs(reportRecordTest.ProductID).WillReturnError(sql.ErrNoRows)
	reportResult, errGet := repo.Get(ctx, reportRecordTest.ProductID)

	// Assert
	assert.EqualError(t, errGet, expectedErr.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NotNil(t, reportResult)
	assert.Empty(t, reportResult)
}

// TestRepository_Get_FailInternalErr passes when unexpected error occurs (return empty domain.ReportRecord and error RepositoryErrInternal)
func TestRepository_Get_FailInternalErr(t *testing.T) {
	// Arrange
	expectedErr := RepositoryErrInternal

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mock.ExpectQuery(regexp.QuoteMeta(GetReportRecord)).WithArgs(reportRecordTest.ProductID).WillReturnError(sql.ErrConnDone)
	reportResult, errGet := repo.Get(ctx, reportRecordTest.ProductID)

	// Assert
	assert.EqualError(t, errGet, expectedErr.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NotNil(t, reportResult)
	assert.Empty(t, reportResult)
}
