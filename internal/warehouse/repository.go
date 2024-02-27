package warehouse

import (
	"context"
	"database/sql"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
)

// Errors
var (
	ErrNotFound       = errors.New("warehouse not found")
	ErrAlreadyExists  = errors.New("warehouse code already exists")
	ErrInternal       = errors.New("database internal error")
	ErrBadRequest     = errors.New("bad request")
	ErrBodyValidation = errors.New("invalid request body")
)

// Queries
const (
	GET_ALL_WAREHOUSES = "SELECT * FROM warehouses"
	GET_WAREHOUSE      = "SELECT * FROM warehouses WHERE id=?;"
	EXISTS             = "SELECT warehouse_code FROM warehouses WHERE warehouse_code=?;"
	SAVE_WAREHOUSE     = "INSERT INTO warehouses (address, telephone, warehouse_code, minimum_capacity, minimum_temperature) VALUES (?, ?, ?, ?, ?)"
	UPDATE_WAREHOUSE   = "UPDATE warehouses SET address=?, telephone=?, warehouse_code=?, minimum_capacity=?, minimum_temperature=? WHERE id=?"
	DELETE_WAREHOUSE   = "DELETE FROM warehouses WHERE id=?"
)

// Repository encapsulates the storage of a warehouse.
type Repository interface {
	GetAll(ctx context.Context) ([]domain.Warehouse, error)
	Get(ctx context.Context, id int) (domain.Warehouse, error)
	Exists(ctx context.Context, warehouseCode string) bool
	Save(ctx context.Context, w domain.Warehouse) (int, error)
	Update(ctx context.Context, w domain.Warehouse) error
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

func (r *repository) GetAll(ctx context.Context) ([]domain.Warehouse, error) {
	rows, err := r.db.Query(GET_ALL_WAREHOUSES)
	if err != nil {
		logging.Log(err)
		return nil, err
	}

	var warehouses []domain.Warehouse

	for rows.Next() {
		w := domain.Warehouse{}
		_ = rows.Scan(&w.ID, &w.Address, &w.Telephone, &w.WarehouseCode, &w.MinimumCapacity, &w.MinimumTemperature)
		warehouses = append(warehouses, w)
	}

	return warehouses, nil
}

func (r *repository) Get(ctx context.Context, id int) (domain.Warehouse, error) {
	row := r.db.QueryRow(GET_WAREHOUSE, id)
	w := domain.Warehouse{}
	err := row.Scan(&w.ID, &w.Address, &w.Telephone, &w.WarehouseCode, &w.MinimumCapacity, &w.MinimumTemperature)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			logging.Log(ErrNotFound)
			return domain.Warehouse{}, ErrNotFound
		default:
			logging.Log(ErrInternal)
			return domain.Warehouse{}, ErrInternal
		}
	}

	return w, nil
}

func (r *repository) Exists(ctx context.Context, warehouseCode string) bool {
	row := r.db.QueryRow(EXISTS, warehouseCode)
	err := row.Scan(&warehouseCode)
	return err == nil
}

func (r *repository) Save(ctx context.Context, w domain.Warehouse) (int, error) {
	stmt, err := r.db.Prepare(SAVE_WAREHOUSE)
	if err != nil {
		logging.Log(err)
		return 0, err
	}

	res, err := stmt.Exec(&w.Address, &w.Telephone, &w.WarehouseCode, &w.MinimumCapacity, &w.MinimumTemperature)
	if err != nil {
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

func (r *repository) Update(ctx context.Context, w domain.Warehouse) error {
	stmt, err := r.db.Prepare(UPDATE_WAREHOUSE)
	if err != nil {
		logging.Log(err)
		return err
	}

	res, err := stmt.Exec(&w.Address, &w.Telephone, &w.WarehouseCode, &w.MinimumCapacity, &w.MinimumTemperature, &w.ID)
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
	stmt, err := r.db.Prepare(DELETE_WAREHOUSE)
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
		logging.Log(ErrNotFound)
		return ErrNotFound
	}

	return nil
}
