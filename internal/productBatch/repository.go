package productbatch

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/go-sql-driver/mysql"
)

var (
	ErrNotFound               = errors.New("product batch not found")
	ErrAlreadyExistsCode      = 1062
	ErrAlreadyExists          = errors.New("batch code already exists")
	ErrDateValueCode          = 1292
	ErrDateValue              = errors.New("the string provided does not match the valid date format yyyy-mm-dd")
	ErrForeignNotFoundCode    = 1452
	ErrForeignProductNotFound = errors.New("the given id does not have a product atached to it")
	ErrForeignSectionNotFound = errors.New("the given id does not have a section atached to it")
	ErrInternal               = errors.New("database internal error")
)

const (
	SaveProductBatch = "INSERT INTO product_batches (batch_number,current_quantity,current_temperature,due_date,initial_quantity,manufacturing_date,manufacturing_hour,minimum_temperature,product_id,section_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
)

func init() {
	logging.InitLog(nil)
}

// Repository encapsulates the storage of a section.
type Repository interface {
	Save(ctx context.Context, pb domain.ProductBatch) (int, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Save(ctx context.Context, pb domain.ProductBatch) (int, error) {

	stmt, err := r.db.Prepare(SaveProductBatch)
	if err != nil {
		logging.Log(err)
		return 0, ErrInternal
	}

	res, err := stmt.Exec(&pb.BatchNumber, &pb.CurrentQuantity, &pb.CurrentTemperature, &pb.DueDate, &pb.InitialQuantity, &pb.ManufacturingDate, &pb.ManufacturingHour, &pb.MinimumTemperature, &pb.ProductID, &pb.SectionID)
	if err != nil {
		message, ok := err.(*mysql.MySQLError)
		if ok {
			switch int(message.Number) {
			case ErrDateValueCode:
				logging.Log(err)
				return 0, ErrDateValue
			case ErrForeignNotFoundCode:
				if strings.Contains(message.Message, "product_id") {
					logging.Log(err)
					return 0, ErrForeignProductNotFound
				} else {
					logging.Log(err)
					return 0, ErrForeignSectionNotFound
				}
			case ErrAlreadyExistsCode:
				logging.Log(err)
				return 0, ErrAlreadyExists
			}
		}
		logging.Log(err)
		return 0, ErrInternal
	}

	id, err := res.LastInsertId()
	if err != nil {
		logging.Log(err)
		return 0, ErrInternal
	}

	return int(id), nil
}
