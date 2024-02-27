package carry

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func init() {
	logging.InitLog(nil)
}

// TestRepositorySave checks the correct operation of the Save repository method
func TestRepositorySave(t *testing.T) {
	// Arrange
	carry := domain.Carry{
		ID:          1,
		CID:         "CID#1",
		CompanyName: "some name",
		Address:     "corrientes 800",
		Telephone:   "4567-4567",
		Locality_id: "1",
	}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"id", "cid", "company_name", "address", "telephone", "locality_id"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(carry.ID, carry.CID, carry.CompanyName, carry.Address, carry.Telephone, carry.Locality_id)

	mock.ExpectPrepare(regexp.QuoteMeta(SAVE_CARRY))
	mock.ExpectExec(regexp.QuoteMeta(SAVE_CARRY)).WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	repository := NewRepository(db)

	newID, err := repository.Save(context.TODO(), carry)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, carry.ID, newID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepositorySaveFailPrepare is correct when function Prepare returns an error
func TestRepositorySaveFailPrepare(t *testing.T) {
	// Arrange
	carry := domain.Carry{
		ID:          1,
		CID:         "CID#1",
		CompanyName: "some name",
		Address:     "corrientes 800",
		Telephone:   "4567-4567",
		Locality_id: "1",
	}
	expectedError := ErrInternal

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(SAVE_CARRY)).WillReturnError(expectedError)

	// Act
	repository := NewRepository(db)

	newID, err := repository.Save(context.TODO(), carry)

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, newID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepositorySaveFailExecInternal is correct when function Exec returns an Internal error
func TestRepositorySaveFailExecInternal(t *testing.T) {
	// Arrange
	carry := domain.Carry{
		ID:          1,
		CID:         "CID#1",
		CompanyName: "some name",
		Address:     "corrientes 800",
		Telephone:   "4567-4567",
		Locality_id: "1",
	}
	expectedError := ErrInternal

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(SAVE_CARRY))
	mock.ExpectExec(regexp.QuoteMeta(SAVE_CARRY)).WillReturnError(expectedError)

	// Act
	repository := NewRepository(db)

	newID, err := repository.Save(context.TODO(), carry)

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, newID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepositorySaveFailExecMySQLFK is correct when function Exec returns an MySQL FK error
func TestRepositorySaveFailExecMySQLFK(t *testing.T) {
	// Arrange
	carry := domain.Carry{
		ID:          1,
		CID:         "CID#1",
		CompanyName: "some name",
		Address:     "corrientes 800",
		Telephone:   "4567-4567",
		Locality_id: "1",
	}
	expectedError := ErrFKConstraint

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(SAVE_CARRY))
	mock.ExpectExec(regexp.QuoteMeta(SAVE_CARRY)).WillReturnError(&mysql.MySQLError{Number: MySqlNumberFKConstraint})

	// Act
	repository := NewRepository(db)

	newID, err := repository.Save(context.TODO(), carry)

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, newID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepositorySaveFailExecMySQLDataLong is correct when function Exec returns an MySQL data long entry error
func TestRepositorySaveFailExecMySQLDataLong(t *testing.T) {
	// Arrange
	carry := domain.Carry{
		ID:          1,
		CID:         "CID#1",
		CompanyName: "some name",
		Address:     "corrientes 800",
		Telephone:   "4567-4567",
		Locality_id: "1",
	}
	expectedError := ErrDataLong

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(SAVE_CARRY))
	mock.ExpectExec(regexp.QuoteMeta(SAVE_CARRY)).WillReturnError(&mysql.MySQLError{Number: MySqlNumberDataLong})

	// Act
	repository := NewRepository(db)

	newID, err := repository.Save(context.TODO(), carry)

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, newID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepositorySaveFailExecMySQLAlreadyExists is correct when function Exec returns an MySQL duplicate entry error
func TestRepositorySaveFailExecMySQLAlreadyExists(t *testing.T) {
	// Arrange
	carry := domain.Carry{
		ID:          1,
		CID:         "CID#1",
		CompanyName: "some name",
		Address:     "corrientes 800",
		Telephone:   "4567-4567",
		Locality_id: "1",
	}
	expectedError := ErrAlreadyExists

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(SAVE_CARRY))
	mock.ExpectExec(regexp.QuoteMeta(SAVE_CARRY)).WillReturnError(&mysql.MySQLError{Number: MySqlNumberDuplicate})

	// Act
	repository := NewRepository(db)

	newID, err := repository.Save(context.TODO(), carry)

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, newID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepositorySaveFailLastID is correct when function LastInsertID returns an error
func TestRepositorySaveFailLastID(t *testing.T) {
	// Arrange
	expectedError := ErrInternal
	carry := domain.Carry{
		ID:          1,
		CID:         "CID#1",
		CompanyName: "some name",
		Address:     "corrientes 800",
		Telephone:   "4567-4567",
		Locality_id: "1",
	}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"id", "cid", "company_name", "address", "telephone", "locality_id"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(carry.ID, carry.CID, carry.CompanyName, carry.Address, carry.Telephone, carry.Locality_id)

	mock.ExpectPrepare(regexp.QuoteMeta(SAVE_CARRY))
	mock.ExpectExec(regexp.QuoteMeta(SAVE_CARRY)).WillReturnResult(sqlmock.NewErrorResult(ErrInternal))

	// Act
	repository := NewRepository(db)

	newID, err := repository.Save(context.TODO(), carry)

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, newID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
