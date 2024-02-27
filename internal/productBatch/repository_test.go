package productbatch

import (
	"context"

	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

var productBatch_test = domain.ProductBatch{
	ID:                 1,
	BatchNumber:        6,
	CurrentQuantity:    1,
	CurrentTemperature: 1,
	DueDate:            "1999-12-12",
	InitialQuantity:    1,
	ManufacturingDate:  "1999-12-12",
	ManufacturingHour:  1,
	MinimumTemperature: 1,
	ProductID:          2,
	SectionID:          1,
}

func TestCreate_Ok(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(SaveProductBatch))
	mock.ExpectExec(regexp.QuoteMeta(SaveProductBatch)).WillReturnResult(sqlmock.NewResult(1, 1))

	columns := []string{
		"id",
		"batch_number",
		"current_quantity",
		"current_temperature",
		"due_date",
		"initial_quantity",
		"manufacturing_date",
		"manufacturing_hour",
		"minimum_temperature",
		"product_id",
		"section_id",
	}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(
		productBatch_test.ID,
		productBatch_test.BatchNumber,
		productBatch_test.CurrentQuantity,
		productBatch_test.CurrentTemperature,
		productBatch_test.DueDate,
		productBatch_test.InitialQuantity,
		productBatch_test.ManufacturingDate,
		productBatch_test.ManufacturingHour,
		productBatch_test.MinimumTemperature,
		productBatch_test.ProductID,
		productBatch_test.SectionID,
	)

	// ACT
	repo := NewRepository(db)

	id, err := repo.Save(context.TODO(), productBatch_test)

	// ASSERT
	assert.NoError(t, err)
	assert.Equal(t, productBatch_test.ID, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_Conflict(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := ErrAlreadyExists

	mock.ExpectPrepare(regexp.QuoteMeta(SaveProductBatch))
	mock.ExpectExec(regexp.QuoteMeta(SaveProductBatch)).WillReturnError(&mysql.MySQLError{Number: uint16(ErrAlreadyExistsCode)})

	// ACT
	repo := NewRepository(db)

	id, err := repo.Save(context.TODO(), productBatch_test)

	// ASSERT
	assert.Empty(t, id)
	assert.EqualError(t, err, expected.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_MissingForeignProductId(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := ErrForeignProductNotFound

	mock.ExpectPrepare(regexp.QuoteMeta(SaveProductBatch))
	mock.ExpectExec(regexp.QuoteMeta(SaveProductBatch)).WillReturnError(&mysql.MySQLError{Number: uint16(ErrForeignNotFoundCode), Message: "product_id"})

	// ACT
	repo := NewRepository(db)

	id, err := repo.Save(context.TODO(), productBatch_test)

	// ASSERT
	assert.Empty(t, id)
	assert.EqualError(t, err, expected.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_MissingForeignSectionId(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := ErrForeignSectionNotFound

	mock.ExpectPrepare(regexp.QuoteMeta(SaveProductBatch))
	mock.ExpectExec(regexp.QuoteMeta(SaveProductBatch)).WillReturnError(&mysql.MySQLError{Number: uint16(ErrForeignNotFoundCode), Message: "section_id"})

	// ACT
	repo := NewRepository(db)

	id, err := repo.Save(context.TODO(), productBatch_test)

	// ASSERT
	assert.Empty(t, id)
	assert.EqualError(t, err, expected.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_InvalidDate(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := ErrDateValue

	mock.ExpectPrepare(regexp.QuoteMeta(SaveProductBatch))
	mock.ExpectExec(regexp.QuoteMeta(SaveProductBatch)).WillReturnError(&mysql.MySQLError{Number: uint16(ErrDateValueCode)})

	// ACT
	repo := NewRepository(db)

	id, err := repo.Save(context.TODO(), productBatch_test)

	// ASSERT
	assert.Empty(t, id)
	assert.EqualError(t, err, expected.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_InternalError(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := ErrInternal

	mock.ExpectPrepare(regexp.QuoteMeta(SaveProductBatch))
	mock.ExpectExec(regexp.QuoteMeta(SaveProductBatch)).WillReturnError(&mysql.MySQLError{})

	// ACT
	repo := NewRepository(db)

	id, err := repo.Save(context.TODO(), productBatch_test)

	// ASSERT
	assert.Empty(t, id)
	assert.EqualError(t, err, expected.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_LastIdError(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := ErrInternal

	mock.ExpectPrepare(regexp.QuoteMeta(SaveProductBatch))
	mock.ExpectExec(regexp.QuoteMeta(SaveProductBatch)).WillReturnResult(sqlmock.NewErrorResult(ErrInternal))

	// ACT
	repo := NewRepository(db)

	id, err := repo.Save(context.TODO(), productBatch_test)

	// ASSERT
	assert.Empty(t, id)
	assert.EqualError(t, err, expected.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_PrepareErr(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := ErrInternal

	// ACT
	repo := NewRepository(db)

	id, err := repo.Save(context.TODO(), productBatch_test)

	// ASSERT
	assert.Empty(t, id)
	assert.EqualError(t, err, expected.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}