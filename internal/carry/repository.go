package carry

import (
	"context"
	"database/sql"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/go-sql-driver/mysql"
)

// Errors
var (
	ErrAlreadyExists  = errors.New("carry cid already exists")
	ErrInternal       = errors.New("database internal error")
	ErrBodyValidation = errors.New("invalid request body")
	ErrFKConstraint   = errors.New("a column table constraint fails")
	ErrDataLong       = errors.New("a field exceeds the maximum length")
)

// Queries
const (
	EXIST_CARRY             = "SELECT cid FROM carries WHERE cid=?;"
	SAVE_CARRY              = "INSERT INTO carries (cid, company_name, address, telephone, locality_id) VALUES (?, ?, ?, ?, ?)"
	MySqlNumberFKConstraint = 1452
	MySqlNumberDataLong     = 1406
	MySqlNumberDuplicate    = 1062
)

// Repository encapsulates the storage of a carry.
type Repository interface {
	Save(ctx context.Context, carry domain.Carry) (int, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Save(ctx context.Context, carry domain.Carry) (int, error) {
	stmt, err := r.db.Prepare(SAVE_CARRY)
	if err != nil {
		logging.Log(ErrInternal)
		return 0, ErrInternal
	}

	res, err := stmt.Exec(&carry.CID, &carry.CompanyName, &carry.Address, &carry.Telephone, &carry.Locality_id)
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

	id, err := res.LastInsertId()
	if err != nil {
		logging.Log(ErrInternal)
		return 0, ErrInternal
	}

	return int(id), nil
}
