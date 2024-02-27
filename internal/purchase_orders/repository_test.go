package purchaseorders

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func init() {
	logging.InitLog(nil)
}

// TestSaveOrderSuccess passes when return Purchase_orders created
func TestSaveOrderSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(INSERT_ORDER_QUERY)).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
	orderId := 1

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	order := domain.Purchase_orders{
		ID:              orderId,
		OrderNumber:     "001",
		OrderDate:       "25-05-2022",
		TrackingCode:    "001",
		BuyerId:         2,
		ProductRecordId: 5,
		OrderStatusId:   1,
	}
	repo := NewRepository(db)
	o, err := repo.SaveOrder(ctx, order)
	assert.NoError(t, err)
	assert.NotZero(t, o)
	assert.Equal(t, o, orderId)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSaveOrderFailPrepare passes when return an error for Internal Server Error
func TestSaveOrderFailPrepare(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	orderId := 1
	mock.ExpectPrepare(regexp.QuoteMeta(INSERT_ORDER_QUERY)).WillReturnError(ErrInternal)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	order := domain.Purchase_orders{
		ID:              orderId,
		OrderNumber:     "001",
		OrderDate:       "25-05-2022",
		TrackingCode:    "001",
		BuyerId:         2,
		ProductRecordId: 5,
		OrderStatusId:   1,
	}
	repo := NewRepository(db)
	o, err := repo.SaveOrder(ctx, order)

	assert.Empty(t, o)
	assert.EqualError(t, ErrInternal, err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSaveOrderFailInternalError passes when return an error for Internal Server Error
func TestSaveOrderFailInternalError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	orderId := 1
	mock.ExpectPrepare(regexp.QuoteMeta(INSERT_ORDER_QUERY)).ExpectExec().WillReturnError(ErrInternal)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	order := domain.Purchase_orders{
		ID:              orderId,
		OrderNumber:     "001",
		OrderDate:       "25-05-2022",
		TrackingCode:    "001",
		BuyerId:         2,
		ProductRecordId: 5,
		OrderStatusId:   1,
	}
	repo := NewRepository(db)
	o, err := repo.SaveOrder(ctx, order)

	assert.Empty(t, o)
	assert.EqualError(t, ErrInternal, err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSaveOrderFailExecMySQLFK passes when return an error for fail an MySQL FK error
func TestSaveOrderFailExecMySQLFK(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	orderId := 1
	mock.ExpectPrepare(regexp.QuoteMeta(INSERT_ORDER_QUERY)).ExpectExec().WillReturnError(&mysql.MySQLError{Number: MySqlNumberFKConstraint})
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	order := domain.Purchase_orders{
		ID:              orderId,
		OrderNumber:     "001",
		OrderDate:       "25-05-2022",
		TrackingCode:    "001",
		BuyerId:         2,
		ProductRecordId: 5,
		OrderStatusId:   1,
	}
	repo := NewRepository(db)
	o, err := repo.SaveOrder(ctx, order)

	assert.Empty(t, o)
	assert.EqualError(t, ErrFKConstraint, err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSaveOrderFailExecMySQLDataLong passes when return an MySQL data long entry error
func TestSaveOrderFailExecMySQLDataLong(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	orderId := 1
	mock.ExpectPrepare(regexp.QuoteMeta(INSERT_ORDER_QUERY)).ExpectExec().WillReturnError(&mysql.MySQLError{Number: MySqlNumberDataLong})
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	order := domain.Purchase_orders{
		ID:              orderId,
		OrderNumber:     "001",
		OrderDate:       "25-05-2022",
		TrackingCode:    "001",
		BuyerId:         2,
		ProductRecordId: 5,
		OrderStatusId:   1,
	}
	repo := NewRepository(db)
	o, err := repo.SaveOrder(ctx, order)

	assert.Empty(t, o)
	assert.EqualError(t, ErrDataLong, err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSaveOrderFailExecMySQLAlreadyExists passes when return an an MySQL duplicate entry error
func TestSaveOrderFailExecMySQLAlreadyExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	orderId := 1
	mock.ExpectPrepare(regexp.QuoteMeta(INSERT_ORDER_QUERY)).ExpectExec().WillReturnError(&mysql.MySQLError{Number: MySqlNumberDuplicate})
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	order := domain.Purchase_orders{
		ID:              orderId,
		OrderNumber:     "001",
		OrderDate:       "25-05-2022",
		TrackingCode:    "001",
		BuyerId:         2,
		ProductRecordId: 5,
		OrderStatusId:   1,
	}
	repo := NewRepository(db)
	o, err := repo.SaveOrder(ctx, order)

	assert.Empty(t, o)
	assert.EqualError(t, ErrAlreadyExists, err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSaveOrderFailLastID passes when function LastInsertID returns an error
func TestSaveOrderFailLastID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	orderId := 1
	mock.ExpectPrepare(regexp.QuoteMeta(INSERT_ORDER_QUERY)).ExpectExec().WillReturnResult(sqlmock.NewErrorResult(ErrInternal))
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	order := domain.Purchase_orders{
		ID:              orderId,
		OrderNumber:     "001",
		OrderDate:       "25-05-2022",
		TrackingCode:    "001",
		BuyerId:         2,
		ProductRecordId: 5,
		OrderStatusId:   1,
	}
	repo := NewRepository(db)
	o, err := repo.SaveOrder(ctx, order)

	assert.Empty(t, o)
	assert.EqualError(t, ErrInternal, err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestExistsSuccess passes when return true
func TestExistsSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id"})
	orderNumber := "001"
	orderId := 1

	rows.AddRow(orderId)
	mock.ExpectQuery(regexp.QuoteMeta(EXISTS_ORDER_QUERY)).WithArgs(orderNumber).WillReturnRows(rows)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	repo := NewRepository(db)
	ok := repo.Exists(ctx, orderNumber)

	assert.NotEmpty(t, ok)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestExistsFail passes when return false
func TestExistsFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	orderNumber := "001"
	mock.ExpectQuery(regexp.QuoteMeta(EXISTS_ORDER_QUERY)).WithArgs(orderNumber).WillReturnError(sql.ErrNoRows)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	repo := NewRepository(db)
	ok := repo.Exists(ctx, orderNumber)

	assert.Equal(t, false, ok)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetByBuyerIdSuccess passes when return list of purchase_orders with same buyer_id
func TestGetByBuyerIdSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"buyer_id", "card_number_id", "first_name", "last_name", "orders_count"})
	buyerId := 1
	orders := domain.Purchase_orders_buyer{
		ID:           buyerId,
		CardNumberId: "001",
		FirstName:    "Comprador 1",
		LastName:     "Vendedor 1",
		OrdersCount:  1,
	}
	rows.AddRow(orders.ID, orders.CardNumberId, orders.FirstName, orders.LastName, orders.OrdersCount)
	mock.ExpectQuery(regexp.QuoteMeta(GET_ORDERS_BY_BUYERID_QUERY)).WillReturnRows(rows)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	repo := NewRepository(db)
	countList, err := repo.GetByBuyerId(ctx, buyerId)

	assert.NoError(t, err)
	assert.NotEmpty(t, countList)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetByBuyerIdFailInternal passes when function Scan cant parsed the values
func TestGetByBuyerIdFailInternal(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	buyerId := 1
	rows := sqlmock.NewRows([]string{"buyer_id", "card_number_id", "first_name", "last_name", "orders_count"})
	rows.AddRow(nil, nil, nil, nil, nil)

	mock.ExpectQuery(regexp.QuoteMeta(GET_ORDERS_BY_BUYERID_QUERY)).WillReturnRows(rows)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	repo := NewRepository(db)
	countList, err := repo.GetByBuyerId(ctx, buyerId)

	assert.Empty(t, countList)
	assert.EqualError(t, ErrInternal, err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetByBuyerIdFailNoRows passes when there is no rows to Scan
func TestGetByBuyerIdFailNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	buyerId := 1
	mock.ExpectQuery(regexp.QuoteMeta(GET_ORDERS_BY_BUYERID_QUERY)).WillReturnError(sql.ErrNoRows)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	repo := NewRepository(db)
	countList, err := repo.GetByBuyerId(ctx, buyerId)
	assert.Empty(t, countList)
	assert.EqualError(t, ErrNotFound, err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetByBuyerIdFailInternalError passes when return an Internal Error (Status 500)
func TestGetByBuyerIdFailInternalError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	buyerId := 1
	mock.ExpectQuery(regexp.QuoteMeta(GET_ORDERS_BY_BUYERID_QUERY)).WillReturnError(sql.ErrConnDone)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	repo := NewRepository(db)
	countList, err := repo.GetByBuyerId(ctx, buyerId)
	fmt.Println("Debugger Agus 8")
	assert.Empty(t, countList)
	assert.EqualError(t, ErrInternal, err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetAllByBuyerSuccess passes when return list of purchase_orders
func TestGetAllByBuyerSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"buyer_id", "card_number_id", "first_name", "last_name", "orders_count"})
	buyerId := 1
	orders := domain.Purchase_orders_buyer{
		ID:           buyerId,
		CardNumberId: "001",
		FirstName:    "Comprador 1",
		LastName:     "Vendedor 1",
		OrdersCount:  2,
	}
	rows.AddRow(orders.ID, orders.CardNumberId, orders.FirstName, orders.LastName, orders.OrdersCount)
	mock.ExpectQuery(regexp.QuoteMeta(GETALL_ORDERS_BY_BUYERID_QUERY)).WillReturnRows(rows)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	repo := NewRepository(db)
	countList, err := repo.GetAllByBuyer(ctx)

	assert.NoError(t, err)
	assert.NotEmpty(t, countList)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetAllByBuyerFailQuery passes when return an error
func TestGetAllByBuyerFailQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(GETALL_ORDERS_BY_BUYERID_QUERY)).WillReturnError(ErrInternal)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	repo := NewRepository(db)
	countList, err := repo.GetAllByBuyer(ctx)

	assert.Empty(t, countList)
	assert.EqualError(t, ErrInternal, err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetAllByBuyerFailScan passes when function Scan returns an error
func TestGetAllByBuyerFailScan(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"buyer_id", "card_number_id", "first_name", "last_name", "orders_count"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(nil, nil, nil, nil, nil)

	mock.ExpectQuery(regexp.QuoteMeta(GETALL_ORDERS_BY_BUYERID_QUERY)).WillReturnRows(rows)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	repo := NewRepository(db)
	countList, err := repo.GetAllByBuyer(ctx)

	assert.Empty(t, countList)
	assert.EqualError(t, ErrInternal, err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}
