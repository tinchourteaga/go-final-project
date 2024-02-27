package purchaseorders

import (
	"context"
	"database/sql"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/go-sql-driver/mysql"
)

var (
	//ErrOrderAlreadyExists = errors.New("order_number already exists")
	ErrInternal       = errors.New("database internal error")
	ErrBodyValidation = errors.New("invalid request body")
	ErrNotFound       = errors.New("Purchase Order not found")
	ErrAlreadyExists  = errors.New("order_number already exists")
	ErrFKConstraint   = errors.New("a column table constraint fails")
	ErrDataLong       = errors.New("a field exceeds the maximum length")
)

const (
	INSERT_ORDER_QUERY             = "INSERT INTO purchase_orders (order_number, order_date, tracking_code, buyer_id, product_record_id,order_status_id) VALUES (?, ?, ?, ?, ?, ?);"
	EXISTS_ORDER_QUERY             = "SELECT id FROM purchase_orders WHERE order_number = ?;"
	GET_ORDERS_BY_BUYERID_QUERY    = "SELECT b.id 'buyer_id', b.card_number_id, b.first_name, b.last_name , count(p.id) 'orders_count' FROM purchase_orders p RIGHT JOIN buyers b ON p.buyer_id = b.id WHERE b.id = ? GROUP BY b.id, b.card_number_id, b.first_name, b.last_name;"
	GETALL_ORDERS_BY_BUYERID_QUERY = "SELECT b.id 'buyer_id', b.card_number_id, b.first_name, b.last_name , count(p.id) 'orders_count' FROM purchase_orders p RIGHT JOIN buyers b ON p.buyer_id = b.id GROUP BY b.id, b.card_number_id, b.first_name, b.last_name;"
	MySqlNumberFKConstraint        = 1452
	MySqlNumberDataLong            = 1406
	MySqlNumberDuplicate           = 1062
)

type Repository interface {
	SaveOrder(ctx context.Context, p domain.Purchase_orders) (int, error)
	Exists(ctx context.Context, orderID string) bool
	GetByBuyerId(ctx context.Context, buyerId int) ([]domain.Purchase_orders_buyer, error)
	GetAllByBuyer(ctx context.Context) ([]domain.Purchase_orders_buyer, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) SaveOrder(ctx context.Context, p domain.Purchase_orders) (int, error) {
	stmt, err := r.db.PrepareContext(ctx, INSERT_ORDER_QUERY)
	if err != nil {
		logging.Log(err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, &p.OrderNumber, &p.OrderDate, &p.TrackingCode, &p.BuyerId, &p.ProductRecordId, &p.OrderStatusId)
	if err != nil {
		mysqlError, ok := err.(*mysql.MySQLError)
		if ok {
			switch mysqlError.Number {
			case MySqlNumberFKConstraint:
				logging.Log(ErrFKConstraint)
				return 0, ErrFKConstraint
			case MySqlNumberDataLong:
				logging.Log(ErrDataLong)
				return 0, ErrDataLong
			case MySqlNumberDuplicate:
				logging.Log(ErrAlreadyExists)
				return 0, ErrAlreadyExists
			}
		}
		logging.Log(ErrInternal)
		return 0, ErrInternal
	}

	id, err := result.LastInsertId()
	if err != nil {
		logging.Log(err)
		return 0, err
	}

	return int(id), nil
}

func (r *repository) Exists(ctx context.Context, orderNumber string) bool {
	row := r.db.QueryRowContext(ctx, EXISTS_ORDER_QUERY, orderNumber)
	err := row.Scan(&orderNumber)
	return err == nil
}

func (r *repository) GetByBuyerId(ctx context.Context, buyerId int) ([]domain.Purchase_orders_buyer, error) {
	var o domain.Purchase_orders_buyer
	var result []domain.Purchase_orders_buyer

	rows := r.db.QueryRowContext(ctx, GET_ORDERS_BY_BUYERID_QUERY, buyerId)

	if err := rows.Scan(&o.ID, &o.CardNumberId, &o.FirstName, &o.LastName, &o.OrdersCount); err != nil {
		switch err {
		case sql.ErrNoRows:
			logging.Log(ErrNotFound)
			return nil, ErrNotFound
		default:
			logging.Log(ErrInternal)
			return nil, ErrInternal
		}
	}
	result = append(result, o)
	return result, nil
}

func (r *repository) GetAllByBuyer(ctx context.Context) ([]domain.Purchase_orders_buyer, error) {
	var result []domain.Purchase_orders_buyer

	rows, err := r.db.QueryContext(ctx, GETALL_ORDERS_BY_BUYERID_QUERY)
	if err != nil {
		logging.Log(ErrInternal)
		return nil, ErrInternal
	}

	for rows.Next() {
		var o domain.Purchase_orders_buyer
		if err := rows.Scan(&o.ID, &o.CardNumberId, &o.FirstName, &o.LastName, &o.OrdersCount); err != nil {
			logging.Log(ErrInternal)
			return nil, ErrInternal
		}
		result = append(result, o)
	}

	return result, nil
}
