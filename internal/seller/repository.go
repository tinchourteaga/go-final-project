package seller

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/go-sql-driver/mysql"
)

// Errors
var (
	ErrNotFound             = errors.New("seller not found")
	ErrLocalityNotExist     = errors.New("locality not exists")
	ErrAlreadyExists        = errors.New("cid already exists")
	ErrInternal             = errors.New("Database internal error")
	ErrForeignKeyConstraint = errors.New("a column table constraint fails")
)

// Repository encapsulates the storage of a Seller.
type Repository interface {
	GetAll(ctx context.Context) ([]domain.Seller, error)
	Get(ctx context.Context, id int) (domain.Seller, error)
	Exists(ctx context.Context, cid int) bool
	Save(ctx context.Context, s domain.Seller) (int, error)
	Update(ctx context.Context, s domain.Seller) error
	Delete(ctx context.Context, id int) error
}

type repository struct {
	db *sql.DB
}

const (
	GET_ALL_SELLERS                 = "SELECT * FROM sellers"
	GET_SELLER                      = "SELECT * FROM sellers WHERE id=?;"
	EXIST_SELLER                    = "SELECT cid FROM sellers WHERE cid=?;"
	SAVE_SELLER                     = "INSERT INTO sellers (cid, company_name, address, telephone, locality_id) VALUES (?, ?, ?, ?, ?)"
	UPDATE_SELLER                   = "UPDATE sellers SET cid=?, company_name=?, address=?, telephone=?, locality_id=? WHERE id=?"
	DELETE_SELLER                   = "DELETE FROM sellers WHERE id=?"
	MySqlNumberForeignKeyConstraint = 1452
)

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func CloseStmt(stmt *sql.Stmt) {
	err := stmt.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (r *repository) GetAll(ctx context.Context) ([]domain.Seller, error) {
	rows, err := r.db.Query(GET_ALL_SELLERS)
	if err != nil {
		return nil, err
	}

	var sellers []domain.Seller

	for rows.Next() {
		s := domain.Seller{}
		_ = rows.Scan(&s.ID, &s.CID, &s.CompanyName, &s.Address, &s.Telephone, &s.Locality_id)
		sellers = append(sellers, s)
	}

	return sellers, nil
}

func (r *repository) Get(ctx context.Context, id int) (domain.Seller, error) {
	row := r.db.QueryRow(GET_SELLER, id)
	s := domain.Seller{}
	err := row.Scan(&s.ID, &s.CID, &s.CompanyName, &s.Address, &s.Telephone, &s.Locality_id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return domain.Seller{}, ErrNotFound
		default:
			return domain.Seller{}, ErrInternal
		}
	}

	return s, nil
}

func (r *repository) Exists(ctx context.Context, cid int) bool {
	row := r.db.QueryRow(EXIST_SELLER, cid)
	err := row.Scan(&cid)
	return err == nil
}

func (r *repository) Save(ctx context.Context, s domain.Seller) (int, error) {
	stmt, err := r.db.Prepare(SAVE_SELLER)
	if err != nil {
		return 0, err
	}
	defer CloseStmt(stmt)

	res, err := stmt.Exec(s.CID, s.CompanyName, s.Address, s.Telephone, s.Locality_id)
	if err != nil {
		mysqlError, ok := err.(*mysql.MySQLError)
		if ok {
			switch mysqlError.Number {
			case MySqlNumberForeignKeyConstraint: // error on foreign key
				return 0, ErrForeignKeyConstraint
			}
		}
		return 0, ErrInternal
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *repository) Update(ctx context.Context, s domain.Seller) error {
	stmt, err := r.db.Prepare(UPDATE_SELLER)
	if err != nil {
		return err
	}
	defer CloseStmt(stmt)

	res, err := stmt.Exec(s.CID, s.CompanyName, s.Address, s.Telephone, s.Locality_id, s.ID)
	if err != nil {
		mysqlError, ok := err.(*mysql.MySQLError)
		if ok {
			switch mysqlError.Number {
			case MySqlNumberForeignKeyConstraint: // error on foreign key
				return ErrForeignKeyConstraint
			}
		}
		return ErrInternal
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	stmt, err := r.db.Prepare(DELETE_SELLER)
	if err != nil {
		return err
	}
	defer CloseStmt(stmt)

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affect < 1 {
		return ErrNotFound
	}

	return nil
}
