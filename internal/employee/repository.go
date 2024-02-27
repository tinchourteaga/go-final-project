package employee

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/go-sql-driver/mysql"
)

const (
	GetAllEmployees = "SELECT * FROM employees"
	GetEmployeeByID = "SELECT * FROM employees WHERE id=?;"
	EmployeeExists  = "SELECT card_number_id FROM employees WHERE card_number_id=?;"
	SaveEmployee    = "INSERT INTO employees(card_number_id,first_name,last_name,warehouse_id) VALUES (?,?,?,?)"
	UpdateEmployee  = "UPDATE employees SET first_name=?, last_name=?, warehouse_id=?  WHERE id=?"
	DeleteEmployee  = "DELETE FROM employees WHERE id=?"
)

const (
	ForeignKeyConstraint = 1452
)

// Errors
var (
	ErrWarehouseNonExistent = errors.New("the associated warehouse does not exist")
)

// Repository encapsulates the storage of a employee.
type Repository interface {
	GetAll(ctx context.Context) ([]domain.Employee, error)
	Get(ctx context.Context, id int) (domain.Employee, error)
	Exists(ctx context.Context, cardNumberID string) bool
	Save(ctx context.Context, e domain.Employee) (int, error)
	Update(ctx context.Context, e domain.Employee) error
	Delete(ctx context.Context, id int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAll(ctx context.Context) ([]domain.Employee, error) {
	rows, err := r.db.Query(GetAllEmployees)
	if err != nil {
		logging.Log(err)
		return nil, err
	}

	var employees []domain.Employee

	for rows.Next() {
		e := domain.Employee{}
		_ = rows.Scan(&e.ID, &e.CardNumberID, &e.FirstName, &e.LastName, &e.WarehouseID)
		employees = append(employees, e)
	}

	return employees, nil
}

func (r *repository) Get(ctx context.Context, id int) (domain.Employee, error) {
	row := r.db.QueryRow(GetEmployeeByID, id)
	e := domain.Employee{}
	err := row.Scan(&e.ID, &e.CardNumberID, &e.FirstName, &e.LastName, &e.WarehouseID)
	if err != nil {
		logging.Log(err)
		return domain.Employee{}, err
	}

	return e, nil
}

func (r *repository) Exists(ctx context.Context, cardNumberID string) bool {
	row := r.db.QueryRow(EmployeeExists, cardNumberID)
	err := row.Scan(&cardNumberID)
	return err == nil
}

func (r *repository) Save(ctx context.Context, e domain.Employee) (int, error) {
	stmt, err := r.db.Prepare(SaveEmployee)
	if err != nil {
		logging.Log(err)
		return 0, err
	}

	res, err := stmt.Exec(&e.CardNumberID, &e.FirstName, &e.LastName, &e.WarehouseID)

	if err != nil {
		message, ok := err.(*mysql.MySQLError)
		if ok {
			switch message.Number {
			case ForeignKeyConstraint:
				if strings.Contains(message.Message, "warehouses") {
					logging.Log(ErrWarehouseNonExistent)
					return 0, ErrWarehouseNonExistent
				}
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

func (r *repository) Update(ctx context.Context, e domain.Employee) error {
	stmt, err := r.db.Prepare(UpdateEmployee)
	if err != nil {
		logging.Log(err)
		return err
	}

	res, err := stmt.Exec(&e.FirstName, &e.LastName, &e.WarehouseID, &e.ID)
	if err != nil {
		logging.Log(err)
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		logging.Log(err)
		return err
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	stmt, err := r.db.Prepare(DeleteEmployee)
	if err != nil {
		logging.Log(err)
		return err
	}

	res, err := stmt.Exec(id)
	if err != nil {
		logging.Log(err)
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		logging.Log(err)
		return err
	}

	if affect < 1 {
		logging.Log(ErrEmployeeNotFound)
		return ErrEmployeeNotFound
	}

	return nil
}
