package product_record

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func init() {
	logging.InitLog(nil)
}

var productRecordTest = domain.ProductRecord{
	ID:             1,
	LastUpdateDate: domain.MySqlTime{Time: time.Now()},
	PurchasePrice:  2.5,
	SalePrice:      3.5,
	ProductID:      1,
}

// TestRepository_Save_OK passes when product id exists and data is correct (return domain.ProductRecord.ID, nil error)
func TestRepository_Save_OK(t *testing.T) {
	// Arrange
	columns := []string{"id", "last_update_date", "purchase_price", "sale_price", "product_id"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(productRecordTest.ID, productRecordTest.LastUpdateDate.Time, productRecordTest.PurchasePrice, productRecordTest.SalePrice, productRecordTest.ProductID)

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	mock.ExpectPrepare(
		regexp.QuoteMeta(SaveProductRecord)).
		ExpectExec().
		WithArgs(productRecordTest.LastUpdateDate.Time, productRecordTest.PurchasePrice, productRecordTest.SalePrice, productRecordTest.ProductID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	saveResultID, errSave := repo.Save(ctx, productRecordTest)

	// Assert
	assert.NoError(t, errSave)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, saveResultID, productRecordTest.ID)
}

// TestRepository_Save_FailProductNonExistent passes when product_id doesn't exist in the database (return 0, error message RepositoryErrForeignKeyConstraint)
func TestRepository_Save_FailProductNonExistent(t *testing.T) {
	// Arrange
	expectedErr := RepositoryErrForeignKeyConstraint

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	mock.ExpectPrepare(regexp.QuoteMeta(SaveProductRecord)).
		ExpectExec().
		WithArgs(productRecordTest.LastUpdateDate.Time, productRecordTest.PurchasePrice, productRecordTest.SalePrice, productRecordTest.ProductID).
		WillReturnError(&mysql.MySQLError{Number: MySqlNumberForeignKeyConstraint})
	saveResultID, errSave := repo.Save(ctx, productRecordTest)

	// Assert
	assert.EqualError(t, errSave, expectedErr.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, saveResultID, 0)
}

// TestRepository_Save_FailQuery passes when query fails to prepare (return 0 and error message)
func TestRepository_Save_FailQuery(t *testing.T) {
	// Arrange
	expectedErr := errors.New("forced query error")

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	mock.ExpectPrepare(regexp.QuoteMeta(SaveProductRecord)).WillReturnError(expectedErr)
	saveResultID, errSave := repo.Save(ctx, productRecordTest)

	// Assert
	assert.EqualError(t, errSave, expectedErr.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, saveResultID, 0)
}

// TestRepository_Save_FailParsingError passes when query execution has an error that can't be cast to *mysql.MySQLError (return 0 and error message)
func TestRepository_Save_FailParsingError(t *testing.T) {
	// Arrange
	expectedErr := errors.New("forced non mysql error")

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	mock.ExpectPrepare(regexp.QuoteMeta(SaveProductRecord)).
		ExpectExec().
		WithArgs(productRecordTest.LastUpdateDate.Time, productRecordTest.PurchasePrice, productRecordTest.SalePrice, productRecordTest.ProductID).
		WillReturnError(expectedErr)
	saveResultID, errSave := repo.Save(ctx, productRecordTest)

	// Assert
	assert.EqualError(t, errSave, expectedErr.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, saveResultID, 0)
}

// TestRepository_Save_FailLastID passes when lastID function returns an error (return 0 and error message)
func TestRepository_Save_FailLastID(t *testing.T) {
	// Arrange
	expectedErr := errors.New("forced last id error")

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	mock.ExpectPrepare(regexp.QuoteMeta(SaveProductRecord)).
		ExpectExec().
		WillReturnResult(sqlmock.NewErrorResult(expectedErr))
	saveResultID, errSave := repo.Save(ctx, productRecordTest)

	// Assert
	assert.EqualError(t, errSave, expectedErr.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, saveResultID, 0)
}

// TestRepository_Get_OK passes when query is successful (return domain.ProductRecord and nil error)
func TestRepository_Get_OK(t *testing.T) {
	// Arrange
	columns := []string{"id", "last_update_date", "purchase_price", "sale_price", "product_id"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(productRecordTest.ID, productRecordTest.LastUpdateDate.Time, productRecordTest.PurchasePrice, productRecordTest.SalePrice, productRecordTest.ProductID)

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mock.ExpectQuery(regexp.QuoteMeta(GetProductRecord)).WithArgs(productRecordTest.ProductID).WillReturnRows(rows)
	productResult, errGet := repo.Get(ctx, productRecordTest.ProductID)

	// Assert
	assert.NoError(t, errGet)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, productRecordTest, productResult)
}

// TestRepository_Get_FailNotFound passes when the given id is not in database (return empty domain.ProductRecord and error RepositoryErrNotFound)
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
	mock.ExpectQuery(regexp.QuoteMeta(GetProductRecord)).WithArgs(productRecordTest.ProductID).WillReturnError(sql.ErrNoRows)
	productResult, errGet := repo.Get(ctx, productRecordTest.ProductID)

	// Assert
	assert.EqualError(t, errGet, expectedErr.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NotNil(t, productResult)
	assert.Empty(t, productResult)
}

// TestRepository_Get_FailInternalErr passes when unexpected error occurs (return empty domain.ProductRecord and error RepositoryErrInternal)
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
	mock.ExpectQuery(regexp.QuoteMeta(GetProductRecord)).WithArgs(productRecordTest.ProductID).WillReturnError(sql.ErrConnDone)
	productResult, errGet := repo.Get(ctx, productRecordTest.ProductID)

	// Assert
	assert.EqualError(t, errGet, expectedErr.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NotNil(t, productResult)
	assert.Empty(t, productResult)
}
