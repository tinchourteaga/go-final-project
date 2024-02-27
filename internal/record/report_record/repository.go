package report_record

import (
	"context"
	"database/sql"
	"errors"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
)

var (
	RepositoryErrNotFound = errors.New("product not found in database")
	RepositoryErrInternal = errors.New("database internal error")
)

const (
	GetAllReportRecords = "SELECT `products`.`id` AS `product_id`, `products`.`description` AS `description`, COUNT(`product_records`.`id`) AS `records_count` FROM `products` LEFT JOIN `product_records` ON `products`.`id` = `product_records`.`product_id` GROUP BY `product_records`.`product_id`, `products`.`id`, `products`.`description`"
	GetReportRecord     = "SELECT `products`.`id` AS `product_id`, `products`.`description` AS `description`, COUNT(`product_records`.`id`) AS `records_count` FROM `products` LEFT JOIN `product_records` ON `products`.`id` = `product_records`.`product_id` WHERE `products`.`id` = ? GROUP BY `product_records`.`product_id`, `products`.`id`, `products`.`description`;"
)

type Repository interface {
	GetAll(ctx context.Context) ([]domain.ReportRecord, error)
	Get(ctx context.Context, productID int) (domain.ReportRecord, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (repository *repository) GetAll(ctx context.Context) (reportRecords []domain.ReportRecord, err error) {
	rows, errQuery := repository.db.QueryContext(ctx, GetAllReportRecords)
	if errQuery != nil {
		logging.Log(errQuery)
		reportRecords = []domain.ReportRecord{}
		err = errQuery
		return
	}
	for rows.Next() {
		reportRecord := domain.ReportRecord{}
		_ = rows.Scan(&reportRecord.ProductID, &reportRecord.Description, &reportRecord.RecordsCount)
		reportRecords = append(reportRecords, reportRecord)
	}
	if len(reportRecords) == 0 {
		reportRecords = []domain.ReportRecord{}
	}
	return
}

func (repository *repository) Get(ctx context.Context, productID int) (reportRecord domain.ReportRecord, err error) {
	row := repository.db.QueryRowContext(ctx, GetReportRecord, productID)
	errScan := row.Scan(&reportRecord.ProductID, &reportRecord.Description, &reportRecord.RecordsCount)
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
