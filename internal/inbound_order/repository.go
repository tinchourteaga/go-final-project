package inbound_order

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/go-sql-driver/mysql"
)

// Queries
const (
	GetAllEmployeesInboundOrders = `SELECT e.id, e.card_number_id, e.first_name, e.last_name, e.warehouse_id, COUNT(ib.id) AS inbound_orders_count
	FROM inbound_orders AS ib
	RIGHT JOIN employees AS e ON ib.employee_id = e.id
	GROUP BY e.id;`
	GetEmployeeInboundOrders = `SELECT e.id, e.card_number_id, e.first_name, e.last_name, e.warehouse_id, COUNT(ib.id) AS inbound_orders_count
	FROM inbound_orders AS ib
	INNER JOIN employees AS e ON ib.employee_id = e.id
	GROUP BY e.id
	HAVING e.id = ?;`
	SaveInboundOrder = `INSERT INTO inbound_orders (order_date, order_number, employee_id, product_batch_id, warehouse_id)
	VALUES (?, ?, ?, ?, ?);`
	InboundOrderExists = `SELECT order_number FROM inbound_orders WHERE order_number = ?;`
)

const (
	ForeignKeyConstraint   = 1452
	DuplicateKeyConstraint = 1062
)

// Errors
var (
	ErrEmployeeNonExistent     = errors.New("the associated employee does not exist")
	ErrWarehouseNonExistent    = errors.New("the associated warehouse does not exist")
	ErrProductBatchNonExistent = errors.New("the associated product batch does not exist")
)

// Repository encapsulates the storage of an inbound order.
type Repository interface {
	GetAllEmployeesInboundOrders(ctx context.Context) ([]domain.EmployeeWithInboundOrders, error)
	GetEmployeeInboundOrders(ctx context.Context, id int) (domain.EmployeeWithInboundOrders, error)
	Save(ctx context.Context, inboundOrder domain.InboundOrder) (int, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAllEmployeesInboundOrders(ctx context.Context) ([]domain.EmployeeWithInboundOrders, error) {
	rows, err := r.db.Query(GetAllEmployeesInboundOrders)

	if err != nil {
		logging.Log(err)
		return nil, err
	}

	var employees []domain.EmployeeWithInboundOrders

	for rows.Next() {
		employee := domain.EmployeeWithInboundOrders{}
		_ = rows.Scan(&employee.ID, &employee.CardNumberID, &employee.FirstName,
			&employee.LastName, &employee.WarehouseID, &employee.InboundOrders)
		employees = append(employees, employee)
	}

	return employees, nil
}

func (r *repository) GetEmployeeInboundOrders(ctx context.Context, id int) (domain.EmployeeWithInboundOrders, error) {
	row := r.db.QueryRow(GetEmployeeInboundOrders, id)
	employee := domain.EmployeeWithInboundOrders{}
	err := row.Scan(&employee.ID, &employee.CardNumberID, &employee.FirstName,
		&employee.LastName, &employee.WarehouseID, &employee.InboundOrders)

	if err != nil {
		logging.Log(err)
		return domain.EmployeeWithInboundOrders{}, err
	}

	return employee, nil
}

func (r *repository) Save(ctx context.Context, inboundOrder domain.InboundOrder) (int, error) {
	stmt, err := r.db.Prepare(SaveInboundOrder)

	if err != nil {
		logging.Log(err)
		return 0, err
	}

	if inboundOrder.OrderNumber == "" {
		logging.Log(ErrEmptyOrderNumber)
		return 0, ErrEmptyOrderNumber
	}

	res, err := stmt.Exec(&inboundOrder.OrderDate, &inboundOrder.OrderNumber,
		&inboundOrder.EmployeeID, &inboundOrder.ProductBatchID, &inboundOrder.WarehouseID)

	if err != nil {
		message, ok := err.(*mysql.MySQLError)
		if ok {
			switch message.Number {
			case ForeignKeyConstraint:
				if strings.Contains(message.Message, "employees") {
					logging.Log(ErrEmployeeNonExistent)
					return 0, ErrEmployeeNonExistent
				}
				if strings.Contains(message.Message, "warehouses") {
					logging.Log(ErrWarehouseNonExistent)
					return 0, ErrWarehouseNonExistent
				}
				if strings.Contains(message.Message, "product_batches") {
					logging.Log(ErrProductBatchNonExistent)
					return 0, ErrProductBatchNonExistent
				}
			case DuplicateKeyConstraint:
				logging.Log(ErrInboundOrderAlreadyExists)
				return 0, ErrInboundOrderAlreadyExists
			}
		}
		logging.Log(err)
		return 0, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		logging.Log(err)
		return 0, err
	}

	return int(id), nil
}
