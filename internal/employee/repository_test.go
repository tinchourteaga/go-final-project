package employee

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/stretchr/testify/assert"
)

func init() {
	logging.InitLog(nil)
}

var (
	ErrInternalError = errors.New("internal error")
	ErrNoRows        = errors.New("sql: no rows in result set")
)

var (
	employeeTest = domain.Employee{
		ID:           1,
		CardNumberID: "asd321",
		FirstName:    "Martin",
		LastName:     "Urteaga Naya",
		WarehouseID:  1,
	}
	employeeTestWarehouseFK = domain.Employee{
		ID:           1,
		CardNumberID: "asd321",
		FirstName:    "Martin",
		LastName:     "Urteaga Naya",
		WarehouseID:  35,
	}
)

func TestRepositoryEmployeeSave_Ok(t *testing.T) {
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	columns := []string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(employeeTest.ID, employeeTest.CardNumberID, employeeTest.FirstName, employeeTest.LastName, employeeTest.WarehouseID)
	mock.ExpectPrepare(
		regexp.QuoteMeta(SaveEmployee)).
		ExpectExec().
		WithArgs(employeeTest.CardNumberID, employeeTest.FirstName, employeeTest.LastName, employeeTest.WarehouseID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	resultID, errSave := repo.Save(ctx, employeeTest)
	assert.NoError(t, errSave)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, resultID, employeeTest.ID)
}

func TestRepositoryEmployeeSave_FailWarehouseFK(t *testing.T) {
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	columns := []string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(employeeTestWarehouseFK.ID, employeeTestWarehouseFK.CardNumberID, employeeTestWarehouseFK.FirstName, employeeTestWarehouseFK.LastName, employeeTestWarehouseFK.WarehouseID)
	mock.ExpectPrepare(
		regexp.QuoteMeta(SaveEmployee)).
		ExpectExec().
		WithArgs(employeeTestWarehouseFK.CardNumberID, employeeTestWarehouseFK.FirstName, employeeTestWarehouseFK.LastName, employeeTestWarehouseFK.WarehouseID).
		WillReturnError(ErrWarehouseNonExistent)
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	result, err := repo.Save(ctx, employeeTestWarehouseFK)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrWarehouseNonExistent.Error())
	assert.Equal(t, 0, result)
}

func TestRepositoryEmployeeSave_FailPrepare(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectPrepare(regexp.QuoteMeta(SaveEmployee)).WillReturnError(ErrInternalError)
	repo := NewRepository(db)
	result, err := repo.Save(context.TODO(), employeeTest)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrInternalError.Error())
	assert.Empty(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryEmployeeSave_FailExec(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectPrepare(regexp.QuoteMeta(SaveEmployee))
	mock.ExpectExec(regexp.QuoteMeta(SaveEmployee)).WillReturnError(ErrInternalError)
	repo := NewRepository(db)
	result, err := repo.Save(context.TODO(), employeeTest)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrInternalError.Error())
	assert.Empty(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryEmployeeSave_NotExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mock.ExpectQuery(regexp.QuoteMeta(EmployeeExists)).WithArgs(employeeTest.CardNumberID).WillReturnError(ErrNoRows)
	result := repo.Exists(ctx, employeeTest.CardNumberID)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, false, result)
}

/*
func TestRepositoryEmployeeSave_Exists(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	columns := []string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(employeeTest.ID, employeeTest.CardNumberID, employeeTest.FirstName, employeeTest.LastName, employeeTest.WarehouseID)
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mock.ExpectQuery(regexp.QuoteMeta(EmployeeExists)).WithArgs(employeeTest.CardNumberID).WillReturnRows(rows)
	result := repo.Exists(ctx, employeeTest.CardNumberID)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, true, result)
}
*/

func TestRepositoryEmployeesGetAll_Ok(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	columns := []string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(employeeTest.ID, employeeTest.CardNumberID, employeeTest.FirstName, employeeTest.LastName, employeeTest.WarehouseID)
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mock.ExpectQuery(regexp.QuoteMeta(GetAllEmployees)).WillReturnRows(rows)
	result, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, []domain.Employee{employeeTest}, result)
}

func TestRepositoryEmployeeGet_Ok(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	columns := []string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(employeeTest.ID, employeeTest.CardNumberID, employeeTest.FirstName, employeeTest.LastName, employeeTest.WarehouseID)
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mock.ExpectQuery(regexp.QuoteMeta(GetEmployeeByID)).WithArgs(employeeTest.ID).WillReturnRows(rows)
	result, err := repo.Get(ctx, employeeTest.ID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, employeeTest, result)
}

func TestRepositoryEmployeeGet_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mock.ExpectQuery(regexp.QuoteMeta(GetEmployeeByID)).WithArgs(employeeTest.ID).WillReturnError(ErrEmployeeNotFound)
	result, err := repo.Get(ctx, employeeTest.ID)
	assert.EqualError(t, err, ErrEmployeeNotFound.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Empty(t, result)
}

func TestRepositoryEmployeeDelete_Ok(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	columns := []string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(employeeTest.ID, employeeTest.CardNumberID, employeeTest.FirstName, employeeTest.LastName, employeeTest.WarehouseID)
	mock.ExpectPrepare(regexp.QuoteMeta(DeleteEmployee))
	mock.ExpectExec(regexp.QuoteMeta(DeleteEmployee)).WillReturnResult(sqlmock.NewResult(1, 1))
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = repo.Delete(ctx, employeeTest.ID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryEmployeeDelete_FailPrepare(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectPrepare(regexp.QuoteMeta(DeleteEmployee)).WillReturnError(ErrInternalError)
	repo := NewRepository(db)
	err = repo.Delete(context.TODO(), employeeTest.ID)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrInternalError.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryEmployeeDelete_FailExec(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectPrepare(regexp.QuoteMeta(DeleteEmployee))
	mock.ExpectExec(regexp.QuoteMeta(DeleteEmployee)).WillReturnError(ErrInternalError)
	repo := NewRepository(db)
	err = repo.Delete(context.TODO(), employeeTest.ID)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrInternalError.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryEmployeeDelete_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectPrepare(regexp.QuoteMeta(DeleteEmployee))
	mock.ExpectExec(regexp.QuoteMeta(DeleteEmployee)).WillReturnError(ErrEmployeeNotFound)
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = repo.Delete(ctx, employeeTest.ID)
	assert.EqualError(t, err, ErrEmployeeNotFound.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryEmployeeUpdate_Ok(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	columns := []string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(employeeTest.ID, employeeTest.CardNumberID, employeeTest.FirstName, employeeTest.LastName, employeeTest.WarehouseID)
	mock.ExpectPrepare(
		regexp.QuoteMeta(UpdateEmployee)).
		ExpectExec().
		WithArgs(employeeTest.FirstName, employeeTest.LastName, employeeTest.WarehouseID, employeeTest.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	err = repo.Update(ctx, employeeTest)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryEmployeeUpdate_FailPrepare(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectPrepare(regexp.QuoteMeta(UpdateEmployee)).WillReturnError(ErrInternalError)
	repo := NewRepository(db)
	err = repo.Update(context.TODO(), employeeTest)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrInternalError.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryEmployeeUpdate_FailExec(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectPrepare(regexp.QuoteMeta(UpdateEmployee))
	mock.ExpectExec(regexp.QuoteMeta(UpdateEmployee)).WillReturnError(ErrInternalError)
	repo := NewRepository(db)
	err = repo.Update(context.TODO(), employeeTest)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrInternalError.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}
