package product

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

var productTest = domain.Product{
	ID:                             1,
	Description:                    "Tomatoes",
	ExpirationRate:                 70,
	FreezingRate:                   20,
	Height:                         20.4,
	Length:                         10.3,
	NetWeight:                      0.0,
	ProductCode:                    "kasbkj8ats9aka9",
	RecommendedFreezingTemperature: -12.6,
	Width:                          6.7,
	ProductTypeID:                  3,
	SellerID:                       newIntPointer(5),
}

// TestRepository_Save_OK passes when seller_id exists and data is correct (return domain.Product.ID and nil error)
func TestRepository_Save_OK(t *testing.T) {
	// Arrange
	columns := []string{"id", "description", "expiration_rate", "freezing_rate", "height", "length", "net_weight", "product_code", "recommended_freezing_temperature", "width", "product_type_id", "seller_id"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(productTest.ID, productTest.Description, productTest.ExpirationRate, productTest.FreezingRate, productTest.Height, productTest.Length, productTest.NetWeight, productTest.ProductCode, productTest.RecommendedFreezingTemperature, productTest.Width, productTest.ProductTypeID, productTest.SellerID)

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	mock.ExpectPrepare(
		regexp.QuoteMeta(SaveProduct)).
		ExpectExec().
		WithArgs(productTest.Description, productTest.ExpirationRate, productTest.FreezingRate, productTest.Height, productTest.Length, productTest.NetWeight, productTest.ProductCode, productTest.RecommendedFreezingTemperature, productTest.Width, productTest.ProductTypeID, productTest.SellerID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	saveResultID, errSave := repo.Save(ctx, productTest)

	// Assert
	assert.NoError(t, errSave)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, saveResultID, productTest.ID)
}

// TestRepository_Save_FailSellerNonExistent passes when seller_id doesn't exist in the database (return 0 and error RepositoryErrForeignKeyConstraint)
func TestRepository_Save_FailSellerNonExistent(t *testing.T) {
	// Arrange
	expectedErr := RepositoryErrForeignKeyConstraint

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	mock.ExpectPrepare(regexp.QuoteMeta(SaveProduct)).
		ExpectExec().
		WithArgs(productTest.Description, productTest.ExpirationRate, productTest.FreezingRate, productTest.Height, productTest.Length, productTest.NetWeight, productTest.ProductCode, productTest.RecommendedFreezingTemperature, productTest.Width, productTest.ProductTypeID, productTest.SellerID).
		WillReturnError(&mysql.MySQLError{Number: MySqlNumberForeignKeyConstraint})
	saveResultID, errSave := repo.Save(ctx, productTest)

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
	mock.ExpectPrepare(regexp.QuoteMeta(SaveProduct)).
		ExpectExec().
		WillReturnResult(sqlmock.NewErrorResult(expectedErr))
	saveResultID, errSave := repo.Save(ctx, productTest)

	// Assert
	assert.EqualError(t, errSave, expectedErr.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, saveResultID, 0)
}

// TestRepository_Save_FailParsingError passes when query execution has an error that cant be cast to *mysql.MySQLError (return 0 and error message)
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
	mock.ExpectPrepare(regexp.QuoteMeta(SaveProduct)).
		ExpectExec().
		WithArgs(productTest.Description, productTest.ExpirationRate, productTest.FreezingRate, productTest.Height, productTest.Length, productTest.NetWeight, productTest.ProductCode, productTest.RecommendedFreezingTemperature, productTest.Width, productTest.ProductTypeID, productTest.SellerID).
		WillReturnError(expectedErr)
	saveResultID, errSave := repo.Save(ctx, productTest)

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
	mock.ExpectPrepare(regexp.QuoteMeta(SaveProduct)).WillReturnError(expectedErr)
	saveResultID, errSave := repo.Save(ctx, productTest)

	// Assert
	assert.EqualError(t, errSave, expectedErr.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, saveResultID, 0)
}

// TestRepository_GetAll_OK passes when query is successful (return slice of all domain.Product and nil error)
func TestRepository_GetAll_OK(t *testing.T) {
	// Arrange
	columns := []string{"id", "description", "expiration_rate", "freezing_rate", "height", "length", "net_weight", "product_code", "recommended_freezing_temperature", "width", "product_type_id", "seller_id"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(productTest.ID, productTest.Description, productTest.ExpirationRate, productTest.FreezingRate, productTest.Height, productTest.Length, productTest.NetWeight, productTest.ProductCode, productTest.RecommendedFreezingTemperature, productTest.Width, productTest.ProductTypeID, productTest.SellerID)

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mock.ExpectQuery(regexp.QuoteMeta(GetAllProducts)).WillReturnRows(rows)
	reportResult, errGetAll := repo.GetAll(ctx)

	// Assert
	assert.NoError(t, errGetAll)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, []domain.Product{productTest}, reportResult)
}

// TestRepository_GetAll_OKEmpty passes when query is successful and database is empty (return nil slice of domain.Product and nil error)
func TestRepository_GetAll_OKEmpty(t *testing.T) {
	// Arrange
	columns := []string{"id", "description", "expiration_rate", "freezing_rate", "height", "length", "net_weight", "product_code", "recommended_freezing_temperature", "width", "product_type_id", "seller_id"}
	rows := sqlmock.NewRows(columns)

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mock.ExpectQuery(regexp.QuoteMeta(GetAllProducts)).WillReturnRows(rows)
	reportResult, errGetAll := repo.GetAll(ctx)

	// Assert
	assert.NoError(t, errGetAll)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Nil(t, reportResult)
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
	mock.ExpectQuery(regexp.QuoteMeta(GetAllProducts)).WillReturnError(expectedErr)
	reportResult, errGetAll := repo.GetAll(ctx)

	// Assert
	assert.EqualError(t, errGetAll, expectedErr.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Nil(t, reportResult)
}

// TestRepository_Get_OK passes when query is successful (return domain.Product and nil error)
func TestRepository_Get_OK(t *testing.T) {
	// Arrange
	columns := []string{"id", "description", "expiration_rate", "freezing_rate", "height", "length", "net_weight", "product_code", "recommended_freezing_temperature", "width", "product_type_id", "seller_id"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(productTest.ID, productTest.Description, productTest.ExpirationRate, productTest.FreezingRate, productTest.Height, productTest.Length, productTest.NetWeight, productTest.ProductCode, productTest.RecommendedFreezingTemperature, productTest.Width, productTest.ProductTypeID, productTest.SellerID)

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mock.ExpectQuery(regexp.QuoteMeta(GetProduct)).WithArgs(productTest.ID).WillReturnRows(rows)
	productResult, errGet := repo.Get(ctx, productTest.ID)

	// Assert
	assert.NoError(t, errGet)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, productTest, productResult)
}

// TestRepository_Get_FailNotFound passes when the given id is not in database (return empty domain.Product and error RepositoryErrNotFound)
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
	mock.ExpectQuery(regexp.QuoteMeta(GetProduct)).WithArgs(productTest.ID).WillReturnError(sql.ErrNoRows)
	productResult, errGet := repo.Get(ctx, productTest.ID)

	// Assert
	assert.EqualError(t, errGet, expectedErr.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NotNil(t, productResult)
	assert.Empty(t, productResult)
}

// TestRepository_Get_FailInternalServerError passes when unexpected error occurs (return empty domain.Product and error RepositoryErrInternal)
func TestRepository_Get_FailInternalServerError(t *testing.T) {
	// Arrange
	expectedErr := RepositoryErrInternal

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mock.ExpectQuery(regexp.QuoteMeta(GetProduct)).WithArgs(productTest.ID).WillReturnError(sql.ErrConnDone)
	productResult, errGet := repo.Get(ctx, productTest.ID)

	// Assert
	assert.EqualError(t, errGet, expectedErr.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NotNil(t, productResult)
	assert.Empty(t, productResult)
}

// TestRepository_Exists_YES passes when given product_code exists on database (return true)
func TestRepository_Exists_YES(t *testing.T) {
	// Arrange
	columns := []string{"product_code"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(productTest.ProductCode)

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mock.ExpectQuery(regexp.QuoteMeta(ExistsProduct)).WithArgs(productTest.ProductCode).WillReturnRows(rows)
	exists := repo.Exists(ctx, productTest.ProductCode)

	// Assert
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.True(t, exists)
}

// TestRepository_Exists_NO passes when given product_code doesn't exist on database (return false)
func TestRepository_Exists_NO(t *testing.T) {
	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mock.ExpectQuery(regexp.QuoteMeta(ExistsProduct)).WithArgs(productTest.ProductCode).WillReturnError(sql.ErrNoRows)
	exists := repo.Exists(ctx, productTest.ProductCode)

	// Assert
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.False(t, exists)
}

// TestRepository_Update_OK passes when query executes with no errors (return nil error)
func TestRepository_Update_OK(t *testing.T) {
	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	mock.ExpectPrepare(
		regexp.QuoteMeta(UpdateProduct)).
		ExpectExec().
		WithArgs(productTest.Description, productTest.ExpirationRate, productTest.FreezingRate, productTest.Height, productTest.Length, productTest.NetWeight, productTest.ProductCode, productTest.RecommendedFreezingTemperature, productTest.Width, productTest.ProductTypeID, productTest.SellerID, productTest.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	errUpdate := repo.Update(ctx, productTest)

	// Assert
	assert.NoError(t, errUpdate)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepository_Update_FailParsingError passes when query execution has an error that cant be cast to *mysql.MySQLError (return error message)
func TestRepository_Update_FailParsingError(t *testing.T) {
	// Arrange
	expectedErr := errors.New("forced non mysql error")

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	mock.ExpectPrepare(regexp.QuoteMeta(UpdateProduct)).
		ExpectExec().
		WithArgs(productTest.Description, productTest.ExpirationRate, productTest.FreezingRate, productTest.Height, productTest.Length, productTest.NetWeight, productTest.ProductCode, productTest.RecommendedFreezingTemperature, productTest.Width, productTest.ProductTypeID, productTest.SellerID, productTest.ID).
		WillReturnError(expectedErr)
	errUpdate := repo.Update(ctx, productTest)

	// Assert
	assert.EqualError(t, errUpdate, expectedErr.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepository_Update_FailQuery passes when query fails to prepare (return error message)
func TestRepository_Update_FailQuery(t *testing.T) {
	// Arrange
	expectedErr := errors.New("forced query error")

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	mock.ExpectPrepare(regexp.QuoteMeta(UpdateProduct)).WillReturnError(expectedErr)
	errUpdate := repo.Update(ctx, productTest)

	// Assert
	assert.EqualError(t, errUpdate, expectedErr.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepository_Update_FailAffected passes when rowsAffected function returns an error (return error message)
func TestRepository_Update_FailAffected(t *testing.T) {
	// Arrange
	expectedErr := errors.New("forced affected error")

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	mock.ExpectPrepare(regexp.QuoteMeta(UpdateProduct)).
		ExpectExec().
		WillReturnResult(sqlmock.NewErrorResult(expectedErr))
	errUpdate := repo.Update(ctx, productTest)

	// Assert
	assert.EqualError(t, errUpdate, expectedErr.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepository_Update_FailSellerNonExistent passes when seller_id is not in database (return error RepositoryErrForeignKeyConstraint)
func TestRepository_Update_FailSellerNonExistent(t *testing.T) {
	// Arrange
	expectedErr := RepositoryErrForeignKeyConstraint
	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	mock.ExpectPrepare(regexp.QuoteMeta(UpdateProduct)).
		ExpectExec().
		WithArgs(productTest.Description, productTest.ExpirationRate, productTest.FreezingRate, productTest.Height, productTest.Length, productTest.NetWeight, productTest.ProductCode, productTest.RecommendedFreezingTemperature, productTest.Width, productTest.ProductTypeID, productTest.SellerID, productTest.ID).
		WillReturnError(&mysql.MySQLError{Number: MySqlNumberForeignKeyConstraint})
	errUpdate := repo.Update(ctx, productTest)

	// Assert
	assert.EqualError(t, errUpdate, expectedErr.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepository_Delete_OK passes when query executes with no errors (return nil error)
func TestRepository_Delete_OK(t *testing.T) {
	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	mock.ExpectPrepare(
		regexp.QuoteMeta(DeleteProduct)).
		ExpectExec().
		WithArgs(productTest.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	errDelete := repo.Delete(ctx, productTest.ID)

	// Assert
	assert.NoError(t, errDelete)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepository_Delete_FailExecError passes when query execution has an error (return error message)
func TestRepository_Delete_FailExecError(t *testing.T) {
	// Arrange
	expectedErr := errors.New("forced execution error")

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	mock.ExpectPrepare(regexp.QuoteMeta(DeleteProduct)).
		ExpectExec().
		WithArgs(productTest.ID).
		WillReturnError(expectedErr)
	errDelete := repo.Delete(ctx, productTest.ID)

	// Assert
	assert.EqualError(t, errDelete, expectedErr.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepository_Delete_FailQuery passes when query fails to prepare (return error message)
func TestRepository_Delete_FailQuery(t *testing.T) {
	// Arrange
	expectedErr := errors.New("forced query error")

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	mock.ExpectPrepare(regexp.QuoteMeta(DeleteProduct)).WillReturnError(expectedErr)
	errDelete := repo.Delete(ctx, productTest.ID)

	// Assert
	assert.EqualError(t, errDelete, expectedErr.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepository_Delete_FailAffected passes when rowsAffected function returns an error (return error message)
func TestRepository_Delete_FailAffected(t *testing.T) {
	// Arrange
	expectedErr := errors.New("forced affected error")

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	mock.ExpectPrepare(regexp.QuoteMeta(DeleteProduct)).
		ExpectExec().
		WillReturnResult(sqlmock.NewErrorResult(expectedErr))
	errDelete := repo.Delete(ctx, productTest.ID)

	// Assert
	assert.EqualError(t, errDelete, expectedErr.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepository_Delete_FailNoRowsAffected passes when rowsAffected function returns 0 rows affected (return error RepositoryErrNotFound)
func TestRepository_Delete_FailNoRowsAffected(t *testing.T) {
	// Arrange
	expectedErr := RepositoryErrNotFound

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	mock.ExpectPrepare(regexp.QuoteMeta(DeleteProduct)).
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(0, 0))
	errDelete := repo.Delete(ctx, productTest.ID)

	// Assert
	assert.EqualError(t, errDelete, expectedErr.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}
