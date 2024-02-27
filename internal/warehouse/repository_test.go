package warehouse

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/stretchr/testify/assert"
)

func init() {
	logging.InitLog(nil)
}

// * ---------------------- GetAll --------------------------
// TestRepositoryGetAll checks the correct operation of the GetAll repository method
func TestRepositoryGetAll(t *testing.T) {
	// Arrange
	warehouse := domain.Warehouse{
		ID:                 1,
		Address:            "avenida siempre viva",
		Telephone:          "+5699911021",
		WarehouseCode:      "W001",
		MinimumCapacity:    100,
		MinimumTemperature: 0,
	}
	warehouses := []domain.Warehouse{warehouse}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"id", "address", "telephone", "warehouse_code", "minimum_capacity", "minimum_temperature"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(warehouse.ID, warehouse.Address, warehouse.Telephone, warehouse.WarehouseCode, warehouse.MinimumCapacity, warehouse.MinimumTemperature)

	mock.ExpectQuery(regexp.QuoteMeta(GET_ALL_WAREHOUSES)).WillReturnRows(rows)

	// Act
	repository := NewRepository(db)
	result, err := repository.GetAll(context.TODO())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, warehouses, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepositoryGetAllFail is correct when function Query returns an error
func TestRepositoryGetAllFail(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expectedError := ErrInternal

	mock.ExpectQuery(regexp.QuoteMeta(GET_ALL_WAREHOUSES)).WillReturnError(expectedError)

	// Act
	repository := NewRepository(db)
	result, err := repository.GetAll(context.TODO())

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// * ---------------------- Get --------------------------
// TestRepositoryGetAll checks the correct operation of the GetAll repository method
func TestRepositoryGet(t *testing.T) {
	// Arrange
	warehouse := domain.Warehouse{
		ID:                 1,
		Address:            "avenida siempre viva",
		Telephone:          "+5699911021",
		WarehouseCode:      "W001",
		MinimumCapacity:    100,
		MinimumTemperature: 0,
	}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"id", "address", "telephone", "warehouse_code", "minimum_capacity", "minimum_temperature"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(warehouse.ID, warehouse.Address, warehouse.Telephone, warehouse.WarehouseCode, warehouse.MinimumCapacity, warehouse.MinimumTemperature)

	mock.ExpectQuery(regexp.QuoteMeta(GET_WAREHOUSE)).WillReturnRows(rows)

	// Act
	repository := NewRepository(db)
	result, err := repository.Get(context.TODO(), warehouse.ID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, warehouse, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepositoryGetFailScanNoRows is correct when function Scan returns a NoRows error
func TestRepositoryGetFailScanNoRows(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	warehouseID := 1
	expectedError := ErrNotFound

	columns := []string{"id", "address", "telephone", "warehouse_code", "minimum_capacity", "minimum_temperature"}
	rows := sqlmock.NewRows(columns)
	mock.ExpectQuery(regexp.QuoteMeta(GET_WAREHOUSE)).WillReturnRows(rows)

	// Act
	repository := NewRepository(db)
	result, err := repository.Get(context.TODO(), warehouseID)

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepositoryGetFailScanNoRows is correct when function Scan returns an Internal error
func TestRepositoryGetFailScanInternal(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	warehouseID := 1
	expectedError := ErrInternal

	columns := []string{"id", "address", "telephone", "warehouse_code", "minimum_capacity", "minimum_temperature"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(nil, nil, nil, nil, nil, nil) // this can't be parsed by Scan function
	mock.ExpectQuery(regexp.QuoteMeta(GET_WAREHOUSE)).WillReturnRows(rows)

	// Act
	repository := NewRepository(db)
	result, err := repository.Get(context.TODO(), warehouseID)

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// * ---------------------- Exists --------------------------
// TestRepositoryExistsFail is correct when there is no warehouse with the given code
func TestRepositoryExistsFail(t *testing.T) {
	// Arrange
	expectedResult := false

	warehouseCode := "W001"
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"id", "address", "telephone", "warehouse_code", "minimum_capacity", "minimum_temperature"}
	rows := sqlmock.NewRows(columns)

	mock.ExpectQuery(regexp.QuoteMeta(EXISTS)).WillReturnRows(rows)

	// Act
	repository := NewRepository(db)
	result := repository.Exists(context.TODO(), warehouseCode)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// * ---------------------- Save --------------------------
// TestRepositorySave checks the correct operation of the Save repository method
func TestRepositorySave(t *testing.T) {
	// Arrange
	warehouse := domain.Warehouse{
		ID:                 1,
		Address:            "avenida siempre viva",
		Telephone:          "+5699911021",
		WarehouseCode:      "W001",
		MinimumCapacity:    100,
		MinimumTemperature: 0,
	}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"id", "address", "telephone", "warehouse_code", "minimum_capacity", "minimum_temperature"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(warehouse.ID, warehouse.Address, warehouse.Telephone, warehouse.WarehouseCode, warehouse.MinimumCapacity, warehouse.MinimumTemperature)

	mock.ExpectPrepare(regexp.QuoteMeta(SAVE_WAREHOUSE))
	mock.ExpectExec(regexp.QuoteMeta(SAVE_WAREHOUSE)).WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	repository := NewRepository(db)
	newID, err := repository.Save(context.TODO(), warehouse)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, warehouse.ID, newID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepositorySaveFailPrepare is correct when function Prepare returns an error
func TestRepositorySaveFailPrepare(t *testing.T) {
	// Arrange
	warehouse := domain.Warehouse{
		ID:                 1,
		Address:            "avenida siempre viva",
		Telephone:          "+5699911021",
		WarehouseCode:      "W001",
		MinimumCapacity:    100,
		MinimumTemperature: 0,
	}
	expectedError := ErrInternal

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(SAVE_WAREHOUSE)).WillReturnError(expectedError)

	// Act
	repository := NewRepository(db)

	newID, err := repository.Save(context.TODO(), warehouse)

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, newID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepositorySaveFailExec is correct when function Exec returns an error
func TestRepositorySaveFailExec(t *testing.T) {
	// Arrange
	warehouse := domain.Warehouse{
		ID:                 1,
		Address:            "avenida siempre viva",
		Telephone:          "+5699911021",
		WarehouseCode:      "W001",
		MinimumCapacity:    100,
		MinimumTemperature: 0,
	}
	expectedError := ErrInternal

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(SAVE_WAREHOUSE))
	mock.ExpectExec(regexp.QuoteMeta(SAVE_WAREHOUSE)).WillReturnError(expectedError)

	// Act
	repository := NewRepository(db)

	newID, err := repository.Save(context.TODO(), warehouse)

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, newID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepositorySaveFailLastID is correct when function LastInsertID returns an error
func TestRepositorySaveFailLastID(t *testing.T) {
	// Arrange
	warehouse := domain.Warehouse{
		ID:                 1,
		Address:            "avenida siempre viva",
		Telephone:          "+5699911021",
		WarehouseCode:      "W001",
		MinimumCapacity:    100,
		MinimumTemperature: 0,
	}
	expectedError := ErrInternal
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"id", "address", "telephone", "warehouse_code", "minimum_capacity", "minimum_temperature"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(warehouse.ID, warehouse.Address, warehouse.Telephone, warehouse.WarehouseCode, warehouse.MinimumCapacity, warehouse.MinimumTemperature)

	mock.ExpectPrepare(regexp.QuoteMeta(SAVE_WAREHOUSE))
	mock.ExpectExec(regexp.QuoteMeta(SAVE_WAREHOUSE)).WillReturnResult(sqlmock.NewErrorResult(ErrInternal))

	// Act
	repository := NewRepository(db)

	newID, err := repository.Save(context.TODO(), warehouse)

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, newID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// * ---------------------- Update --------------------------
// TestRepositoryUpdate checks the correct operation of the Update repository method
func TestRepositoryUpdate(t *testing.T) {
	// Arrange
	warehouse := domain.Warehouse{
		ID:                 1,
		Address:            "avenida siempre viva",
		Telephone:          "+5699911021",
		WarehouseCode:      "W001",
		MinimumCapacity:    100,
		MinimumTemperature: 0,
	}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"id", "address", "telephone", "warehouse_code", "minimum_capacity", "minimum_temperature"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(warehouse.ID, warehouse.Address, warehouse.Telephone, warehouse.WarehouseCode, warehouse.MinimumCapacity, warehouse.MinimumTemperature)

	mock.ExpectPrepare(regexp.QuoteMeta(UPDATE_WAREHOUSE))
	mock.ExpectExec(regexp.QuoteMeta(UPDATE_WAREHOUSE)).WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	repository := NewRepository(db)
	err = repository.Update(context.TODO(), warehouse)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepositoryUpdateFailPrepare is correct when function Prepare returns an error
func TestRepositoryUpdateFailPrepare(t *testing.T) {
	// Arrange
	warehouse := domain.Warehouse{
		ID:                 1,
		Address:            "avenida siempre viva",
		Telephone:          "+5699911021",
		WarehouseCode:      "W001",
		MinimumCapacity:    100,
		MinimumTemperature: 0,
	}
	expectedError := ErrInternal

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(UPDATE_WAREHOUSE)).WillReturnError(expectedError)

	// Act
	repository := NewRepository(db)

	err = repository.Update(context.TODO(), warehouse)

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepositoryUpdateFailExec is correct when function Exec returns an error
func TestRepositoryUpdateFailExec(t *testing.T) {
	// Arrange
	warehouse := domain.Warehouse{
		ID:                 1,
		Address:            "avenida siempre viva",
		Telephone:          "+5699911021",
		WarehouseCode:      "W001",
		MinimumCapacity:    100,
		MinimumTemperature: 0,
	}
	expectedError := ErrInternal

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(UPDATE_WAREHOUSE))
	mock.ExpectExec(regexp.QuoteMeta(UPDATE_WAREHOUSE)).WillReturnError(expectedError)

	// Act
	repository := NewRepository(db)

	err = repository.Update(context.TODO(), warehouse)

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepositoryUpdateFailRowsAffected is correct when function RowsAffected returns an error
func TestRepositoryUpdateFailRowsAffected(t *testing.T) {
	// Arrange
	warehouse := domain.Warehouse{
		ID:                 1,
		Address:            "avenida siempre viva",
		Telephone:          "+5699911021",
		WarehouseCode:      "W001",
		MinimumCapacity:    100,
		MinimumTemperature: 0,
	}
	expectedError := ErrInternal
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"id", "address", "telephone", "warehouse_code", "minimum_capacity", "minimum_temperature"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(warehouse.ID, warehouse.Address, warehouse.Telephone, warehouse.WarehouseCode, warehouse.MinimumCapacity, warehouse.MinimumTemperature)

	mock.ExpectPrepare(regexp.QuoteMeta(UPDATE_WAREHOUSE))
	mock.ExpectExec(regexp.QuoteMeta(UPDATE_WAREHOUSE)).WillReturnResult(sqlmock.NewErrorResult(ErrInternal))

	// Act
	repository := NewRepository(db)

	err = repository.Update(context.TODO(), warehouse)

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// * ---------------------- Delete --------------------------
// TestRepositoryDelete checks the correct operation of the Delete repository method
func TestRepositoryDelete(t *testing.T) {
	// Arrange
	warehouseID := 1
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(DELETE_WAREHOUSE))
	mock.ExpectExec(regexp.QuoteMeta(DELETE_WAREHOUSE)).WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	repository := NewRepository(db)
	err = repository.Delete(context.TODO(), warehouseID)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepositoryDeleteFailPrepare is correct when function Prepare returns an error
func TestRepositoryDeleteFailPrepare(t *testing.T) {
	// Arrange
	warehouseID := 1
	expectedError := ErrInternal

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(DELETE_WAREHOUSE)).WillReturnError(expectedError)

	// Act
	repository := NewRepository(db)

	err = repository.Delete(context.TODO(), warehouseID)

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepositoryDeleteFailExec is correct when function Exec returns an error
func TestRepositoryDeleteFailExec(t *testing.T) {
	// Arrange
	warehouseID := 1
	expectedError := ErrInternal

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(DELETE_WAREHOUSE))
	mock.ExpectExec(regexp.QuoteMeta(DELETE_WAREHOUSE)).WillReturnError(expectedError)

	// Act
	repository := NewRepository(db)

	err = repository.Delete(context.TODO(), warehouseID)

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepositoryDeleteFailRowsAffected is correct when function RowsAffected returns an error
func TestRepositoryDeleteFailRowsAffected(t *testing.T) {
	// Arrange
	warehouseID := 1
	expectedError := ErrInternal
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(DELETE_WAREHOUSE))
	mock.ExpectExec(regexp.QuoteMeta(DELETE_WAREHOUSE)).WillReturnResult(sqlmock.NewErrorResult(ErrInternal))

	// Act
	repository := NewRepository(db)

	err = repository.Delete(context.TODO(), warehouseID)

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepositoryDeleteFailNotFound is correct when there is no wharehouse with the given id
func TestRepositoryDeleteFailNotFound(t *testing.T) {
	// Arrange
	warehouseID := 1
	expectedError := ErrNotFound
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(DELETE_WAREHOUSE))
	mock.ExpectExec(regexp.QuoteMeta(DELETE_WAREHOUSE)).WillReturnResult(sqlmock.NewResult(1, 0)) // forced 0 affected rows

	// Act
	repository := NewRepository(db)
	err = repository.Delete(context.TODO(), warehouseID)

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}
