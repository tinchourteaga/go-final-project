package section

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

var section_test = domain.Section{
	ID:                 1,
	SectionNumber:      2,
	CurrentTemperature: 3.0,
	MinimumTemperature: -2.0,
	CurrentCapacity:    90,
	MaximumCapacity:    1100,
	MinimumCapacity:    20,
	WarehouseID:        2,
	ProductTypeID:      2,
}

var productsBySection_test = domain.ProductsBySection{
	SectionID:     1,
	SectionNumber: 1,
	ProductsCount: 100,
}

func TestExist_OK(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"section_number"}
	rows := sqlmock.NewRows(columns)

	rows.AddRow(section_test.ID)
	mock.ExpectQuery(regexp.QuoteMeta(ExistsSection)).WithArgs(section_test.SectionNumber).WillReturnRows(rows)

	// ACT
	repo := NewRepository(db)
	result := repo.Exists(context.TODO(), section_test.SectionNumber)

	// ASSERT
	assert.True(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProductsBySections_Ok(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := []domain.ProductsBySection{productsBySection_test}

	columns := []string{
		"section_id",
		"section_number",
		"products_count",
	}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(
		productsBySection_test.SectionID,
		productsBySection_test.SectionNumber,
		productsBySection_test.ProductsCount,
	)

	mock.ExpectQuery(regexp.QuoteMeta(ProductsBySections)).WillReturnRows(rows)

	// ACT
	repo := NewRepository(db)

	productsBySections, err := repo.GetProductsBySections(context.TODO())

	// ASSERT
	assert.NoError(t, err)
	assert.Equal(t, expected, productsBySections)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProductsBySections_InternalErr(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := ErrInternal

	mock.ExpectQuery(regexp.QuoteMeta(ProductsBySections)).WillReturnError(&mysql.MySQLError{})

	// ACT
	repo := NewRepository(db)

	productsBySections, err := repo.GetProductsBySections(context.TODO())

	// ASSERT
	assert.Empty(t, productsBySections)
	assert.EqualError(t, err, expected.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProductsBySections_ScanErr(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{
		"test_column_1",
		"test_column_2",
		"test_column_3",
		"test_column_4",
	}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(
		"test_row_value_1",
		"test_row_value_2",
		"test_row_value_3",
		"test_row_value_4",
	)

	expected := ErrInternal

	mock.ExpectQuery(regexp.QuoteMeta(ProductsBySections)).WillReturnRows(rows)

	// ACT
	repo := NewRepository(db)

	productsBySections, err := repo.GetProductsBySections(context.TODO())

	// ASSERT
	assert.Empty(t, productsBySections)
	assert.EqualError(t, err, expected.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProductsBySection_Ok(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := []domain.ProductsBySection{productsBySection_test}

	columns := []string{
		"section_id",
		"section_number",
		"products_count",
	}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(
		productsBySection_test.SectionID,
		productsBySection_test.SectionNumber,
		productsBySection_test.ProductsCount,
	)

	mock.ExpectQuery(regexp.QuoteMeta(ProductsBySection)).WillReturnRows(rows)

	// ACT
	repo := NewRepository(db)

	productsBySections, err := repo.GetProductsBySection(context.TODO(), productsBySection_test.SectionID)

	// ASSERT
	assert.NoError(t, err)
	assert.Equal(t, expected, productsBySections)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProductsBySection_NotFound(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := ErrNotFound

	mock.ExpectQuery(regexp.QuoteMeta(ProductsBySection)).WillReturnError(sql.ErrNoRows)

	// ACT
	repo := NewRepository(db)

	id, err := repo.GetProductsBySection(context.TODO(), productsBySection_test.SectionID)

	// ASSERT
	assert.Empty(t, id)
	assert.EqualError(t, err, expected.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProductsBySection_InternalErr(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := ErrInternal

	mock.ExpectQuery(regexp.QuoteMeta(ProductsBySection)).WillReturnError(errors.New("error"))

	// ACT
	repo := NewRepository(db)

	id, err := repo.GetProductsBySection(context.TODO(), productsBySection_test.SectionID)

	// ASSERT
	assert.Empty(t, id)
	assert.EqualError(t, err, expected.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAll_Ok(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := []domain.Section{section_test}

	columns := []string{
		"id",
		"section_number",
		"current_temperature",
		"minimum_temperature",
		"current_capacity",
		"maximum_capacity",
		"minimum_capacity",
		"warehouse_id",
		"product_type_id",
	}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(
		section_test.ID,
		section_test.SectionNumber,
		section_test.CurrentTemperature,
		section_test.MinimumTemperature,
		section_test.CurrentCapacity,
		section_test.MinimumCapacity,
		section_test.MaximumCapacity,
		section_test.WarehouseID,
		section_test.ProductTypeID,
	)

	mock.ExpectQuery(regexp.QuoteMeta(GetAllSections)).WillReturnRows(rows)

	// ACT
	repo := NewRepository(db)

	sections, err := repo.GetAll(context.TODO())

	// ASSERT
	assert.NoError(t, err)
	assert.Equal(t, expected, sections)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAll_InternalErr(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := ErrInternal

	mock.ExpectQuery(regexp.QuoteMeta(GetAllSections)).WillReturnError(&mysql.MySQLError{})

	// ACT
	repo := NewRepository(db)

	sections, err := repo.GetAll(context.TODO())

	// ASSERT
	assert.Empty(t, sections)
	assert.EqualError(t, err, expected.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGet_Ok(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{
		"id",
		"section_number",
		"current_temperature",
		"minimum_temperature",
		"current_capacity",
		"maximum_capacity",
		"minimum_capacity",
		"warehouse_id",
		"product_type_id",
	}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(
		section_test.ID,
		section_test.SectionNumber,
		section_test.CurrentTemperature,
		section_test.MinimumTemperature,
		section_test.CurrentCapacity,
		section_test.MinimumCapacity,
		section_test.MaximumCapacity,
		section_test.WarehouseID,
		section_test.ProductTypeID,
	)

	mock.ExpectQuery(regexp.QuoteMeta(GetSection)).WithArgs(section_test.ID).WillReturnRows(rows)

	// ACT
	repo := NewRepository(db)

	section, err := repo.Get(context.TODO(), section_test.ID)

	// ASSERT
	assert.NoError(t, err)
	assert.Equal(t, section_test, section)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGet_NotFound(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := ErrNotFound

	mock.ExpectQuery(regexp.QuoteMeta(GetSection)).WillReturnError(sql.ErrNoRows)

	// ACT
	repo := NewRepository(db)

	section, err := repo.Get(context.TODO(), section_test.ID)

	// ASSERT
	assert.Empty(t, section)
	assert.EqualError(t, err, expected.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGet_InternalErr(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := ErrInternal

	mock.ExpectQuery(regexp.QuoteMeta(GetSection)).WillReturnError(errors.New("error"))

	// ACT
	repo := NewRepository(db)

	section, err := repo.Get(context.TODO(), section_test.ID)

	// ASSERT
	assert.Empty(t, section)
	assert.EqualError(t, err, expected.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_Ok(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(SaveSection))
	mock.ExpectExec(regexp.QuoteMeta(SaveSection)).WillReturnResult(sqlmock.NewResult(1, 1))

	columns := []string{
		"id",
		"section_number",
		"current_temperature",
		"minimum_temperature",
		"current_capacity",
		"maximum_capacity",
		"minimum_capacity",
		"warehouse_id",
		"product_type_id",
	}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(
		section_test.ID,
		section_test.SectionNumber,
		section_test.CurrentTemperature,
		section_test.MinimumTemperature,
		section_test.CurrentCapacity,
		section_test.MinimumCapacity,
		section_test.MaximumCapacity,
		section_test.WarehouseID,
		section_test.ProductTypeID,
	)

	// ACT
	repo := NewRepository(db)

	id, err := repo.Save(context.TODO(), section_test)

	// ASSERT
	assert.NoError(t, err)
	assert.Equal(t, section_test.ID, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_MissingForeign(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := ErrForeignNotFound

	mock.ExpectPrepare(regexp.QuoteMeta(SaveSection))
	mock.ExpectExec(regexp.QuoteMeta(SaveSection)).WillReturnError(&mysql.MySQLError{Number: uint16(ErrForeignNotFoundCode)})

	// ACT
	repo := NewRepository(db)

	id, err := repo.Save(context.TODO(), section_test)

	// ASSERT
	assert.Empty(t, id)
	assert.EqualError(t, err, expected.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_InternalErr(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := ErrInternal

	mock.ExpectPrepare(regexp.QuoteMeta(SaveSection))
	mock.ExpectExec(regexp.QuoteMeta(SaveSection)).WillReturnError(&mysql.MySQLError{})

	// ACT
	repo := NewRepository(db)

	id, err := repo.Save(context.TODO(), section_test)

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

	id, err := repo.Save(context.TODO(), section_test)

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

	mock.ExpectPrepare(regexp.QuoteMeta(SaveSection))
	mock.ExpectExec(regexp.QuoteMeta(SaveSection)).WillReturnResult(sqlmock.NewErrorResult(ErrInternal))

	// ACT
	repo := NewRepository(db)

	id, err := repo.Save(context.TODO(), section_test)

	// ASSERT
	assert.Empty(t, id)
	assert.EqualError(t, err, expected.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_Ok(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(UpdateSection))
	mock.ExpectExec(regexp.QuoteMeta(UpdateSection)).WillReturnResult(sqlmock.NewResult(1, 1))

	// ACT
	repo := NewRepository(db)

	err = repo.Update(context.TODO(), section_test)

	// ASSERT
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_RowsAffectedErr(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := ErrInternal

	mock.ExpectPrepare(regexp.QuoteMeta(UpdateSection))
	mock.ExpectExec(regexp.QuoteMeta(UpdateSection)).WillReturnResult(sqlmock.NewErrorResult(ErrInternal))

	// ACT
	repo := NewRepository(db)

	err = repo.Update(context.TODO(), section_test)

	// ASSERT
	assert.EqualError(t, err, expected.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_ExecErr(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := ErrInternal

	mock.ExpectPrepare(regexp.QuoteMeta(UpdateSection))
	mock.ExpectExec(regexp.QuoteMeta(UpdateSection)).WillReturnError(&mysql.MySQLError{})

	// ACT
	repo := NewRepository(db)

	err = repo.Update(context.TODO(), section_test)

	// ASSERT
	assert.EqualError(t, err, expected.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_PrepareErr(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := ErrInternal

	// ACT
	repo := NewRepository(db)

	err = repo.Update(context.TODO(), section_test)

	// ASSERT
	assert.EqualError(t, err, expected.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_Ok(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(DeleteSection))
	mock.ExpectExec(regexp.QuoteMeta(DeleteSection)).WillReturnResult(sqlmock.NewResult(1, 1))

	// ACT
	repo := NewRepository(db)

	err = repo.Delete(context.TODO(),1)

	//ASSERT
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_RowsAffectedErr(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := ErrInternal

	mock.ExpectPrepare(regexp.QuoteMeta(DeleteSection))
	mock.ExpectExec(regexp.QuoteMeta(DeleteSection)).WillReturnResult(sqlmock.NewErrorResult(ErrInternal))

	// ACT
	repo := NewRepository(db)

	err = repo.Delete(context.TODO(), 1)

	// ASSERT
	assert.EqualError(t, err, expected.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_ExecErr(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := ErrInternal

	mock.ExpectPrepare(regexp.QuoteMeta(DeleteSection))
	mock.ExpectExec(regexp.QuoteMeta(DeleteSection)).WillReturnError(&mysql.MySQLError{})

	// ACT
	repo := NewRepository(db)

	err = repo.Delete(context.TODO(), 1)

	// ASSERT
	assert.EqualError(t, err, expected.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_PrepareErr(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := ErrInternal

	// ACT
	repo := NewRepository(db)

	err = repo.Delete(context.TODO(), 1)

	// ASSERT
	assert.EqualError(t, err, expected.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_NotFoundErr(t *testing.T) {
	// ARRANGE
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := ErrNotFound

	mock.ExpectPrepare(regexp.QuoteMeta(DeleteSection))
	mock.ExpectExec(regexp.QuoteMeta(DeleteSection)).WillReturnResult(sqlmock.NewResult(0, 0))

	// ACT
	repo := NewRepository(db)

	err = repo.Delete(context.TODO(), 1)

	// ASSERT
	assert.EqualError(t, err, expected.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}