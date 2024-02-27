package seller

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

var (
	ErrInternalTest = errors.New("Internal error")
)

var seller_test = domain.Seller{
	ID:          1,
	CID:         1,
	CompanyName: "Kiosco 1",
	Address:     "Mitre 1323",
	Telephone:   "26647273999",
	Locality_id: "5700",
}

func TestExist_OK(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"id"}
	rows := sqlmock.NewRows(columns)

	rows.AddRow(seller_test.ID)
	mock.ExpectQuery(regexp.QuoteMeta(EXIST_SELLER)).WithArgs(seller_test.ID).WillReturnRows(rows)

	// Act
	repo := NewRepository(db)
	result := repo.Exists(context.TODO(), seller_test.ID)

	// Assert
	assert.True(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// * ---------------------- GetAll ----------------------------
// TestGetAll_Seller_OK passes when return all corrects sellers
func TestGetAll_Seller_OK(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	column := []string{"id", "cid", "company_name", "address", "telephone", "locality_id"}
	rows := sqlmock.NewRows(column)
	sellers := []domain.Seller{{ID: seller_test.ID, CID: seller_test.CID, CompanyName: seller_test.CompanyName, Address: seller_test.Address, Telephone: seller_test.Telephone, Locality_id: seller_test.Locality_id}}

	for _, s := range sellers {
		rows.AddRow(s.ID, s.CID, s.CompanyName, s.Address, s.Telephone, s.Locality_id)
	}

	mock.ExpectQuery(regexp.QuoteMeta(GET_ALL_SELLERS)).WillReturnRows(rows)

	// Act
	repo := NewRepository(db)
	result, err := repo.GetAll(context.TODO())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, sellers, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetAll_Seller_Fail passes when return an internal error
func TestGetAll_Seller_Fail(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	column := []string{"id", "cid", "company_name", "address", "telephone", "locality_id"}
	rows := sqlmock.NewRows(column)
	sellers := []domain.Seller{{ID: seller_test.ID, CID: seller_test.CID, CompanyName: seller_test.CompanyName, Address: seller_test.Address, Telephone: seller_test.Telephone, Locality_id: seller_test.Locality_id}}

	for _, s := range sellers {
		rows.AddRow(s.ID, s.CID, s.CompanyName, s.Address, s.Telephone, s.Locality_id)
	}

	mock.ExpectQuery(regexp.QuoteMeta(GET_ALL_SELLERS)).WillReturnError(ErrInternalTest)

	// Act
	repo := NewRepository(db)
	result, err := repo.GetAll(context.TODO())

	// Assert
	assert.EqualError(t, err, ErrInternalTest.Error())
	assert.Empty(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// * ------------------------ Get ------------------------------
// TestGet_Seller_OK passes when return the correct seller by id
func TestGet_Seller_OK(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	column := []string{"id", "cid", "company_name", "address", "telephone", "locality_id"}
	rows := sqlmock.NewRows(column)
	seller := domain.Seller{ID: seller_test.ID, CID: seller_test.CID, CompanyName: seller_test.CompanyName, Address: seller_test.Address, Telephone: seller_test.Telephone, Locality_id: seller_test.Locality_id}

	rows.AddRow(seller.ID, seller.CID, seller.CompanyName, seller.Address, seller.Telephone, seller.Locality_id)

	mock.ExpectQuery(regexp.QuoteMeta(GET_SELLER)).WillReturnRows(rows)

	// Act
	repo := NewRepository(db)
	result, err := repo.Get(context.TODO(), seller.ID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, seller, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGet_Seller_FailErrNotFound passes when return an error not found
func TestGet_Seller_FailErrNotFound(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	column := []string{"id", "cid", "company_name", "address", "telephone", "locality_id"}
	rows := sqlmock.NewRows(column)
	seller := domain.Seller{ID: seller_test.ID, CID: seller_test.CID, CompanyName: seller_test.CompanyName, Address: seller_test.Address, Telephone: seller_test.Telephone, Locality_id: seller_test.Locality_id}

	rows.AddRow(seller.ID, seller.CID, seller.CompanyName, seller.Address, seller.Telephone, seller.Locality_id)

	mock.ExpectQuery(regexp.QuoteMeta(GET_SELLER)).WithArgs(seller_test.ID).WillReturnError(sql.ErrNoRows)

	// Act
	repo := NewRepository(db)
	result, err := repo.Get(context.TODO(), seller.ID)

	// Assert
	assert.EqualError(t, err, ErrNotFound.Error())
	assert.Empty(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGet_Seller_FailErrInternal passes when return an internal error
func TestGet_Seller_FailErrInternal(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	column := []string{"id", "cid", "company_name", "address", "telephone", "locality_id"}
	rows := sqlmock.NewRows(column)
	seller := domain.Seller{ID: seller_test.ID, CID: seller_test.CID, CompanyName: seller_test.CompanyName, Address: seller_test.Address, Telephone: seller_test.Telephone, Locality_id: seller_test.Locality_id}

	rows.AddRow(seller.ID, seller.CID, seller.CompanyName, seller.Address, seller.Telephone, seller.Locality_id)

	mock.ExpectQuery(regexp.QuoteMeta(GET_SELLER)).WithArgs(seller_test.ID).WillReturnError(sql.ErrConnDone)

	// Act
	repo := NewRepository(db)
	result, err := repo.Get(context.TODO(), seller.ID)

	// Assert
	assert.EqualError(t, err, ErrInternal.Error())
	assert.Empty(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// * ---------------------- Save --------------------------
// TestSave_Seller_OK passes when return seller created
func TestSave_Seller_OK(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(SAVE_SELLER))
	mock.ExpectExec(regexp.QuoteMeta(SAVE_SELLER)).WillReturnResult(sqlmock.NewResult(1, 1))

	column := []string{"id", "cid", "company_name", "address", "telephone", "locality_id"}
	rows := sqlmock.NewRows(column)
	rows.AddRow(seller_test.ID, seller_test.CID, seller_test.CompanyName, seller_test.Address, seller_test.Telephone, seller_test.Locality_id)

	// Act
	repository := NewRepository(db)

	newID, err := repository.Save(context.TODO(), seller_test)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, seller_test.ID, newID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSave_Seller_FailForeignKeyConstraint passes when seller id doesn't exist in the database
func TestSave_Seller_FailForeignKeyConstraint(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(SAVE_SELLER))
	mock.ExpectExec(regexp.QuoteMeta(SAVE_SELLER)).WillReturnError(&mysql.MySQLError{Number: MySqlNumberForeignKeyConstraint})

	// Act
	repository := NewRepository(db)

	newID, err := repository.Save(context.TODO(), seller_test)

	// Assert
	assert.EqualError(t, err, ErrForeignKeyConstraint.Error())
	assert.Equal(t, 0, newID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSave_Seller_FailPrepare passes when function Prepare returns an error
func TestSave_Seller_FailPrepare(t *testing.T) {
	// Arrange
	seller := domain.Seller{
		ID:          1,
		CID:         1,
		CompanyName: "Kiosco 1",
		Address:     "Mitre 1323",
		Telephone:   "26647273999",
		Locality_id: "5700",
	}
	expectedError := ErrInternal

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(SAVE_SELLER)).WillReturnError(expectedError)

	// Act
	repository := NewRepository(db)

	newID, err := repository.Save(context.TODO(), seller)

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, newID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSave_Seller_FailErrInternal passes when returns an internal error
func TestSave_Seller_FailErrInternal(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(SAVE_SELLER))
	mock.ExpectExec(regexp.QuoteMeta(SAVE_SELLER)).WillReturnError(sql.ErrConnDone)

	// Act
	repository := NewRepository(db)

	newID, err := repository.Save(context.TODO(), seller_test)

	// Assert
	assert.EqualError(t, err, ErrInternal.Error())
	assert.Equal(t, 0, newID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSave_Seller_FailLastID passes when function LastInsertID returns an error
func TestSave_Seller_FailLastID(t *testing.T) {
	// Arrange
	expectedError := ErrInternal
	seller := domain.Seller{
		ID:          1,
		CID:         1,
		CompanyName: "Kiosco 1",
		Address:     "Mitre 1323",
		Telephone:   "26647273999",
		Locality_id: "5700",
	}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	column := []string{"id", "cid", "company_name", "address", "telephone", "locality_id"}
	rows := sqlmock.NewRows(column)
	rows.AddRow(seller_test.ID, seller_test.CID, seller_test.CompanyName, seller_test.Address, seller_test.Telephone, seller_test.Locality_id)

	mock.ExpectPrepare(regexp.QuoteMeta(SAVE_SELLER))
	mock.ExpectExec(regexp.QuoteMeta(SAVE_SELLER)).WillReturnResult(sqlmock.NewErrorResult(ErrInternal))

	// Act
	repository := NewRepository(db)

	newID, err := repository.Save(context.TODO(), seller)

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, newID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// * ---------------------- Update --------------------------
// TestUpdate_Seller_OK passes when return seller updated successfully
func TestUpdate_Seller_OK(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	seller := domain.Seller{
		ID:          1,
		CID:         1,
		CompanyName: "Kiosco 1",
		Address:     "Mitre 1323",
		Telephone:   "26647273999",
		Locality_id: "5700",
	}

	mock.ExpectPrepare(regexp.QuoteMeta(UPDATE_SELLER))
	mock.ExpectExec(regexp.QuoteMeta(UPDATE_SELLER)).WithArgs(seller.CID, seller.CompanyName, seller.Address, seller.Telephone, seller.Locality_id, seller.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	column := []string{"id", "cid", "company_name", "address", "telephone", "locality_id"}
	rows := sqlmock.NewRows(column)
	rows.AddRow(seller_test.ID, seller_test.CID, seller_test.CompanyName, seller_test.Address, seller_test.Telephone, seller_test.Locality_id)

	repository := NewRepository(db)

	// Act
	result := repository.Update(context.TODO(), seller)

	// Assert
	assert.NoError(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUpdate_Seller_FailPrepare passes when function Prepare returns an error
func TestUpdate_Seller_FailPrepare(t *testing.T) {
	// Arrange
	seller := domain.Seller{
		ID:          1,
		CID:         1,
		CompanyName: "Kiosco 1",
		Address:     "Mitre 1323",
		Telephone:   "26647273999",
		Locality_id: "5700",
	}
	expectedError := ErrInternal

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(UPDATE_SELLER)).WillReturnError(expectedError)

	// Act
	repository := NewRepository(db)

	result := repository.Update(context.TODO(), seller)

	// Assert
	assert.EqualError(t, result, expectedError.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUpdate_Seller_FailForeignKeyConstraint passes when seller id doesn't exist in the database
func TestUpdate_Seller_FailForeignKeyConstraint(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(UPDATE_SELLER))
	mock.ExpectExec(regexp.QuoteMeta(UPDATE_SELLER)).WillReturnError(&mysql.MySQLError{Number: MySqlNumberForeignKeyConstraint})

	// Act
	repository := NewRepository(db)

	result := repository.Update(context.TODO(), seller_test)

	// Assert
	assert.EqualError(t, result, ErrForeignKeyConstraint.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUpdate_Seller_FailErrInternal passes when returns an internal error
func TestUpdate_Seller_FailErrInternal(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(UPDATE_SELLER))
	mock.ExpectExec(regexp.QuoteMeta(UPDATE_SELLER)).WillReturnError(sql.ErrConnDone)

	// Act
	repository := NewRepository(db)

	result := repository.Update(context.TODO(), seller_test)

	// Assert
	assert.EqualError(t, result, ErrInternal.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// * ---------------------- Delete --------------------------
// TestDelete_Selle_OK passes when delete a seller is successful
func TestDelete_Selle_OK(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(DELETE_SELLER))
	mock.ExpectExec(regexp.QuoteMeta(DELETE_SELLER)).WithArgs(seller_test.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	repository := NewRepository(db)

	// Act
	result := repository.Delete(context.TODO(), seller_test.ID)

	// Assert
	assert.NoError(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestDelete_Seller_FailPrepare passes when function Prepare returns an error
func TestDelete_Seller_FailPrepare(t *testing.T) {
	// Arrange
	expectedError := ErrInternal

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(DELETE_SELLER)).WillReturnError(expectedError)

	// Act
	repository := NewRepository(db)

	result := repository.Delete(context.TODO(), seller_test.ID)

	// Assert
	assert.EqualError(t, result, expectedError.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}
