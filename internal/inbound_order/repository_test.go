package inbound_order

import (
	"context"
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
	inboundOrderTest = domain.InboundOrder{
		ID:             1,
		OrderDate:      "10/10/20",
		OrderNumber:    "Order#Test",
		EmployeeID:     1,
		ProductBatchID: 1,
		WarehouseID:    1,
	}
	emptyInboundOrderTest = domain.InboundOrder{
		ID:             1,
		OrderDate:      "10/10/20",
		OrderNumber:    "",
		EmployeeID:     1,
		ProductBatchID: 1,
		WarehouseID:    1,
	}
	inboundOrderTestEmployeeFK = domain.InboundOrder{
		ID:             1,
		OrderDate:      "10/10/20",
		OrderNumber:    "Order#Test",
		EmployeeID:     10,
		ProductBatchID: 1,
		WarehouseID:    1,
	}
	inboundOrderTestWarehouseFK = domain.InboundOrder{
		ID:             1,
		OrderDate:      "10/10/20",
		OrderNumber:    "Order#Test",
		EmployeeID:     1,
		ProductBatchID: 1,
		WarehouseID:    10,
	}
	inboundOrderTestProductBatchFK = domain.InboundOrder{
		ID:             1,
		OrderDate:      "10/10/20",
		OrderNumber:    "Order#Test",
		EmployeeID:     1,
		ProductBatchID: 10,
		WarehouseID:    1,
	}
	employee = domain.Employee{
		ID:           1,
		CardNumberID: "asd321",
		FirstName:    "Ignacio",
		LastName:     "Naya",
		WarehouseID:  1,
	}
	employeeWithInboundOrders = domain.EmployeeWithInboundOrders{
		Employee:      employee,
		InboundOrders: 5,
	}
)

func TestRepositoryInboundOrdersSave_Ok(t *testing.T) {
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	columns := []string{"id", "order_date", "order_number", "employee_id", "product_batch_id", "warehouse_id"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(inboundOrderTest.ID, inboundOrderTest.OrderDate, inboundOrderTest.OrderNumber, inboundOrderTest.EmployeeID, inboundOrderTest.ProductBatchID, inboundOrderTest.WarehouseID)
	mock.ExpectPrepare(
		regexp.QuoteMeta(SaveInboundOrder)).
		ExpectExec().
		WithArgs(inboundOrderTest.OrderDate, inboundOrderTest.OrderNumber, inboundOrderTest.EmployeeID, inboundOrderTest.ProductBatchID, inboundOrderTest.WarehouseID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	resultID, errSave := repo.Save(ctx, inboundOrderTest)
	assert.NoError(t, errSave)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, resultID, inboundOrderTest.ID)
}

func TestRepositoryInboundOrdersSave_FailEmployeeFK(t *testing.T) {
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	columns := []string{"id", "order_date", "order_number", "employee_id", "product_batch_id", "warehouse_id"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(inboundOrderTestEmployeeFK.ID, inboundOrderTestEmployeeFK.OrderDate, inboundOrderTestEmployeeFK.OrderNumber, inboundOrderTestEmployeeFK.EmployeeID, inboundOrderTestEmployeeFK.ProductBatchID, inboundOrderTestEmployeeFK.WarehouseID)
	mock.ExpectPrepare(
		regexp.QuoteMeta(SaveInboundOrder)).
		ExpectExec().
		WithArgs(inboundOrderTestEmployeeFK.OrderDate, inboundOrderTestEmployeeFK.OrderNumber, inboundOrderTestEmployeeFK.EmployeeID, inboundOrderTestEmployeeFK.ProductBatchID, inboundOrderTestEmployeeFK.WarehouseID).
		WillReturnError(ErrEmployeeNonExistent)
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	result, err := repo.Save(ctx, inboundOrderTestEmployeeFK)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrEmployeeNonExistent.Error())
	assert.Equal(t, 0, result)
}

func TestRepositoryInboundOrdersSave_FailWarehouseFK(t *testing.T) {
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	columns := []string{"id", "order_date", "order_number", "employee_id", "product_batch_id", "warehouse_id"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(inboundOrderTestWarehouseFK.ID, inboundOrderTestWarehouseFK.OrderDate, inboundOrderTestWarehouseFK.OrderNumber, inboundOrderTestWarehouseFK.EmployeeID, inboundOrderTestWarehouseFK.ProductBatchID, inboundOrderTestWarehouseFK.WarehouseID)
	mock.ExpectPrepare(
		regexp.QuoteMeta(SaveInboundOrder)).
		ExpectExec().
		WithArgs(inboundOrderTestWarehouseFK.OrderDate, inboundOrderTestWarehouseFK.OrderNumber, inboundOrderTestWarehouseFK.EmployeeID, inboundOrderTestWarehouseFK.ProductBatchID, inboundOrderTestWarehouseFK.WarehouseID).
		WillReturnError(ErrWarehouseNonExistent)
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	result, err := repo.Save(ctx, inboundOrderTestWarehouseFK)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrWarehouseNonExistent.Error())
	assert.Equal(t, 0, result)
}

func TestRepositoryInboundOrdersSave_FailProductBatchFK(t *testing.T) {
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	columns := []string{"id", "order_date", "order_number", "employee_id", "product_batch_id", "warehouse_id"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(inboundOrderTestProductBatchFK.ID, inboundOrderTestProductBatchFK.OrderDate, inboundOrderTestProductBatchFK.OrderNumber, inboundOrderTestProductBatchFK.EmployeeID, inboundOrderTestProductBatchFK.ProductBatchID, inboundOrderTestProductBatchFK.WarehouseID)
	mock.ExpectPrepare(
		regexp.QuoteMeta(SaveInboundOrder)).
		ExpectExec().
		WithArgs(inboundOrderTestProductBatchFK.OrderDate, inboundOrderTestProductBatchFK.OrderNumber, inboundOrderTestProductBatchFK.EmployeeID, inboundOrderTestProductBatchFK.ProductBatchID, inboundOrderTestProductBatchFK.WarehouseID).
		WillReturnError(ErrProductBatchNonExistent)
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	result, err := repo.Save(ctx, inboundOrderTestProductBatchFK)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrProductBatchNonExistent.Error())
	assert.Equal(t, 0, result)
}

func TestRepositoryInboundOrdersSave_Conflict(t *testing.T) {
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	columns := []string{"id", "order_date", "order_number", "employee_id", "product_batch_id", "warehouse_id"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(emptyInboundOrderTest.ID, emptyInboundOrderTest.OrderDate, emptyInboundOrderTest.OrderNumber, emptyInboundOrderTest.EmployeeID, emptyInboundOrderTest.ProductBatchID, emptyInboundOrderTest.WarehouseID)
	mock.ExpectPrepare(
		regexp.QuoteMeta(SaveInboundOrder)).
		ExpectExec().
		WithArgs(emptyInboundOrderTest.OrderDate, emptyInboundOrderTest.OrderNumber, emptyInboundOrderTest.EmployeeID, emptyInboundOrderTest.ProductBatchID, emptyInboundOrderTest.WarehouseID).
		WillReturnError(ErrEmptyOrderNumber)
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	result, err := repo.Save(ctx, emptyInboundOrderTest)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrEmptyOrderNumber.Error())
	assert.Equal(t, 0, result)
}

func TestRepositoryInboundOrdersGetAllEmployees_Ok(t *testing.T) {
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	columns := []string{"id", "card_number_id", "first_name", "last_name", "warehouse_id", "inbound_orders_count"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(employeeWithInboundOrders.ID, employeeWithInboundOrders.CardNumberID, employeeWithInboundOrders.FirstName,
		employeeWithInboundOrders.LastName, employeeWithInboundOrders.WarehouseID, employeeWithInboundOrders.InboundOrders)
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mock.ExpectQuery(regexp.QuoteMeta(GetAllEmployeesInboundOrders)).WillReturnRows(rows)
	result, err := repo.GetAllEmployeesInboundOrders(ctx)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, []domain.EmployeeWithInboundOrders{employeeWithInboundOrders}, result)
}

func TestRepositoryInboundOrdersGetEmployee_Ok(t *testing.T) {
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	columns := []string{"id", "card_number_id", "first_name", "last_name", "warehouse_id", "inbound_orders_count"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(employeeWithInboundOrders.ID, employeeWithInboundOrders.CardNumberID, employeeWithInboundOrders.FirstName,
		employeeWithInboundOrders.LastName, employeeWithInboundOrders.WarehouseID, employeeWithInboundOrders.InboundOrders)
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mock.ExpectQuery(regexp.QuoteMeta(GetEmployeeInboundOrders)).WithArgs(employeeWithInboundOrders.ID).WillReturnRows(rows)
	result, err := repo.GetEmployeeInboundOrders(ctx, emptyInboundOrderTest.EmployeeID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, employeeWithInboundOrders, result)
}

func TestRepositoryInboundOrdersGetEmployee_NotFound(t *testing.T) {
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	repo := NewRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mock.ExpectQuery(regexp.QuoteMeta(GetEmployeeInboundOrders)).WithArgs(employeeWithInboundOrders.ID).WillReturnError(ErrEmployeeWithInboundOrdersNotFound)
	result, err := repo.GetEmployeeInboundOrders(ctx, employeeWithInboundOrders.ID)
	assert.EqualError(t, err, ErrEmployeeWithInboundOrdersNotFound.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Empty(t, result)
}
