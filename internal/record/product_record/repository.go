package product_record

import (
	"context"
	"database/sql"
	"errors"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/product"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/go-sql-driver/mysql"
)

var (
	RepositoryErrNotFound             = errors.New("product record not found in database")
	RepositoryErrInternal             = errors.New("database internal error")
	RepositoryErrForeignKeyConstraint = errors.New("a foreign key constraint fails")
)

const (
	SaveProductRecord               = "INSERT INTO `product_records`(`last_update_date`, `purchase_price`, `sale_price`, `product_id`) VALUES (?, ?, ?, ?);"
	GetProductRecord                = "SELECT `id`, `last_update_date`, `purchase_price`, `sale_price`, `product_id` FROM `product_records` WHERE `id` = ?;"
	MySqlNumberForeignKeyConstraint = 1452
)

type Repository interface {
	Get(ctx context.Context, id int) (domain.ProductRecord, error)
	Save(ctx context.Context, record domain.ProductRecord) (int, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (repository *repository) Get(ctx context.Context, id int) (productRecord domain.ProductRecord, err error) {
	row := repository.db.QueryRowContext(ctx, GetProductRecord, id)
	errScan := row.Scan(&productRecord.ID, &productRecord.LastUpdateDate.Time, &productRecord.PurchasePrice, &productRecord.SalePrice, &productRecord.ProductID)
	if errScan != nil {
		logging.Log(errScan)
		switch errScan {
		case sql.ErrNoRows:
			err = RepositoryErrNotFound
		default:
			err = RepositoryErrInternal
		}
	}
	return
}

func (repository *repository) Save(ctx context.Context, record domain.ProductRecord) (savedID int, err error) {
	stmt, errPrepare := repository.db.PrepareContext(ctx, SaveProductRecord)
	if errPrepare != nil {
		logging.Log(errPrepare)
		err = errPrepare
		return
	}
	defer product.CloseStmt(stmt)
	result, errExec := stmt.ExecContext(ctx, record.LastUpdateDate.Time, record.PurchasePrice, record.SalePrice, record.ProductID)
	if errExec != nil {
		logging.Log(errExec)
		message, ok := errExec.(*mysql.MySQLError)
		if ok {
			switch message.Number {
			case MySqlNumberForeignKeyConstraint:
				err = RepositoryErrForeignKeyConstraint
			}
			return
		}
		err = errExec
		return
	}
	id, errID := result.LastInsertId()
	if errID != nil {
		logging.Log(errID)
		err = errID
		return
	}
	savedID = int(id)
	return
}
