package buyer

import (
	"context"
	"database/sql"
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

// TestSaveBuyerSuccess passes when return Buyer created
func TestSaveBuyerSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(INSERT_QUERY)).ExpectExec().WillReturnResult(sqlmock.NewResult(4, 1))
	orderId := 4

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	order := domain.Buyer{
		ID:           4,
		CardNumberID: "004",
		FirstName:    "Fernando",
		LastName:     "Perez",
	}

	repo := NewRepository(db)
	o, err := repo.Save(ctx, order)

	assert.NoError(t, err)
	assert.NotZero(t, o)
	assert.Equal(t, o, orderId)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSaveBuyerFailPrepare passes when return an error for Internal Server Error
func TestSaveBuyerFailPrepare(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(INSERT_QUERY)).WillReturnError(ErrInternal)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	order := domain.Buyer{
		ID:           4,
		CardNumberID: "004",
		FirstName:    "Fernando",
		LastName:     "Perez",
	}
	repo := NewRepository(db)
	o, err := repo.Save(ctx, order)

	assert.Empty(t, o)
	assert.EqualError(t, ErrInternal, err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSaveBuyerFailInternalError passes when return an error for Internal Server Error
func TestSaveBuyerFailInternalError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(INSERT_QUERY)).ExpectExec().WillReturnError(ErrInternal)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	order := domain.Buyer{
		ID:           4,
		CardNumberID: "004",
		FirstName:    "Fernando",
		LastName:     "Perez",
	}
	repo := NewRepository(db)
	o, err := repo.Save(ctx, order)

	assert.Empty(t, o)
	assert.EqualError(t, ErrInternal, err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSaveBuyerFailLastID passes when function LastInsertID returns an error
func TestSaveBuyerFailLastID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(INSERT_QUERY)).ExpectExec().WillReturnResult(sqlmock.NewErrorResult(ErrInternal))
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	order := domain.Buyer{
		ID:           4,
		CardNumberID: "004",
		FirstName:    "Fernando",
		LastName:     "Perez",
	}
	repo := NewRepository(db)
	o, err := repo.Save(ctx, order)

	assert.Empty(t, o)
	assert.EqualError(t, ErrInternal, err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestExistsBuyerSuccess passes when return true
func TestExistsBuyerSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id"})
	CardNumberId := "001"
	Id := 1

	rows.AddRow(Id)
	mock.ExpectQuery(regexp.QuoteMeta(EXISTS_QUERY)).WithArgs(CardNumberId).WillReturnRows(rows)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	repo := NewRepository(db)
	ok := repo.Exists(ctx, CardNumberId)

	assert.NotEmpty(t, ok)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestExistsBuyerFail passes when return false
func TestExistsBuyerFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	CardNumberId := "001"
	mock.ExpectQuery(regexp.QuoteMeta(EXISTS_QUERY)).WithArgs(CardNumberId).WillReturnError(sql.ErrNoRows)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	repo := NewRepository(db)
	ok := repo.Exists(ctx, CardNumberId)

	assert.Equal(t, false, ok)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetBuyerSuccess passes when return list of purchase_orders with same buyer_id
func TestGetBuyerSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name"})
	buyerId := 1
	orders := domain.Buyer{
		ID:           1,
		CardNumberID: "001",
		FirstName:    "Fernando",
		LastName:     "Perez",
	}
	rows.AddRow(orders.ID, orders.CardNumberID, orders.FirstName, orders.LastName)
	mock.ExpectQuery(regexp.QuoteMeta(GET_BY_ID_QUERY)).WillReturnRows(rows)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	repo := NewRepository(db)
	countList, err := repo.Get(ctx, buyerId)

	assert.NoError(t, err)
	assert.NotEmpty(t, countList)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetBuyerFailInternal passes when function Scan cant parsed the values
func TestGetBuyerFailInternal(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	buyerId := 1
	rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name"})
	rows.AddRow(nil, nil, nil, nil)

	mock.ExpectQuery(regexp.QuoteMeta(GET_BY_ID_QUERY)).WillReturnRows(rows)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	repo := NewRepository(db)
	countList, err := repo.Get(ctx, buyerId)

	assert.Empty(t, countList)
	assert.EqualError(t, ErrInternal, err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetBuyerFailNoRows passes when there is no rows to Scan
func TestGetBuyerFailNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	buyerId := 1
	mock.ExpectQuery(regexp.QuoteMeta(GET_BY_ID_QUERY)).WillReturnError(sql.ErrNoRows)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	repo := NewRepository(db)
	countList, err := repo.Get(ctx, buyerId)
	assert.Empty(t, countList)
	assert.EqualError(t, ErrNotFound, err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetAllBuyersSuccess passes when return list of buyers
func TestGetAllBuyersSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name"})
	orders := domain.Buyer{
		ID:           1,
		CardNumberID: "001",
		FirstName:    "Fernando",
		LastName:     "Perez",
	}
	rows.AddRow(orders.ID, orders.CardNumberID, orders.FirstName, orders.LastName)
	mock.ExpectQuery(regexp.QuoteMeta(GET_ALL_QUERY)).WillReturnRows(rows)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	repo := NewRepository(db)
	countList, err := repo.GetAll(ctx)

	assert.NoError(t, err)
	assert.NotEmpty(t, countList)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetAllBuyersFailQuery passes when return an error
func TestGetAllBuyersFailQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(GET_ALL_QUERY)).WillReturnError(ErrInternal)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	repo := NewRepository(db)
	countList, err := repo.GetAll(ctx)

	assert.Empty(t, countList)
	assert.EqualError(t, ErrInternal, err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSaveBuyerSuccess passes when return nil
func TestUpdateBuyerSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	Buyer_Id := 4
	rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name"})
	order := domain.Buyer{
		ID:           Buyer_Id,
		CardNumberID: "004",
		FirstName:    "Fernando",
		LastName:     "Perez",
	}
	rows.AddRow(order.ID, order.CardNumberID, order.FirstName, order.LastName)
	mock.ExpectPrepare(regexp.QuoteMeta(UPDATE_QUERY)).ExpectExec().WillReturnResult(sqlmock.NewResult(4, 1))
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	repo := NewRepository(db)
	err = repo.Update(ctx, order)

	assert.Empty(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUpdateBuyerFailPrepare passes when return an error "data base internal error"
func TestUpdateBuyerFailPrepare(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	Buyer_Id := 4
	rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name"})
	order := domain.Buyer{
		ID:           Buyer_Id,
		CardNumberID: "004",
		FirstName:    "Fernando",
		LastName:     "Perez",
	}
	rows.AddRow(order.ID, order.CardNumberID, order.FirstName, order.LastName)
	mock.ExpectPrepare(regexp.QuoteMeta(UPDATE_QUERY)).WillReturnError(ErrInternal)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	repo := NewRepository(db)
	err = repo.Update(ctx, order)

	assert.NotEmpty(t, err)
	assert.EqualError(t, ErrInternal, err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUpdateBuyerFailExecQuery passes when return an error "data base internal error"
func TestUpdateBuyerFailExecQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	Buyer_Id := 4
	rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name"})
	order := domain.Buyer{
		ID:           Buyer_Id,
		CardNumberID: "004",
		FirstName:    "Fernando",
		LastName:     "Perez",
	}
	rows.AddRow(order.ID, order.CardNumberID, order.FirstName, order.LastName)
	mock.ExpectPrepare(regexp.QuoteMeta(UPDATE_QUERY)).ExpectExec().WillReturnError(ErrInternal)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	repo := NewRepository(db)
	err = repo.Update(ctx, order)

	assert.NotEmpty(t, err)
	assert.EqualError(t, ErrInternal, err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUpdateBuyerFailRowsAffected passes when function RowsAffected returns an error"
func TestUpdateBuyerFailRowsAffected(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	Buyer_Id := 4
	rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name"})
	order := domain.Buyer{
		ID:           Buyer_Id,
		CardNumberID: "004",
		FirstName:    "Fernando",
		LastName:     "Perez",
	}
	rows.AddRow(order.ID, order.CardNumberID, order.FirstName, order.LastName)
	mock.ExpectPrepare(regexp.QuoteMeta(UPDATE_QUERY)).ExpectExec().WillReturnResult(sqlmock.NewErrorResult(ErrInternal))
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	repo := NewRepository(db)
	err = repo.Update(ctx, order)

	assert.NotEmpty(t, err)
	assert.EqualError(t, ErrInternal, err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSaveBuyerSuccess passes when return nil and buyer´s deleted
func TestDeleteBuyerSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	Buyer_Id := 1
	mock.ExpectPrepare(regexp.QuoteMeta(DELETE_QUERY)).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	repo := NewRepository(db)
	err = repo.Delete(ctx, Buyer_Id)

	assert.Empty(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestDeleteBuyerFailPrepare passes when function Prepare returns an error
func TestDeleteBuyerFailPrepare(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	Buyer_Id := 1
	mock.ExpectPrepare(regexp.QuoteMeta(DELETE_QUERY)).WillReturnError(ErrInternal)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	repo := NewRepository(db)
	err = repo.Delete(ctx, Buyer_Id)

	assert.NotEmpty(t, err)
	assert.EqualError(t, ErrInternal, err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestDeleteBuyerFailExecQuery passes when return an error "data base internal error"
func TestDeleteBuyerFailExecQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	Buyer_Id := 1
	mock.ExpectPrepare(regexp.QuoteMeta(DELETE_QUERY)).ExpectExec().WillReturnError(ErrInternal)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	repo := NewRepository(db)
	err = repo.Delete(ctx, Buyer_Id)

	assert.NotEmpty(t, err)
	assert.EqualError(t, ErrInternal, err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestDeleteBuyerFailRowsAffected passes when function RowsAffected returns an error"
func TestDeleteBuyerFailRowsAffected(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	Buyer_Id := 1
	mock.ExpectPrepare(regexp.QuoteMeta(DELETE_QUERY)).ExpectExec().WillReturnResult(sqlmock.NewErrorResult(ErrInternal))
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	repo := NewRepository(db)
	err = repo.Delete(ctx, Buyer_Id)

	assert.NotEmpty(t, err)
	assert.EqualError(t, ErrInternal, err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}
