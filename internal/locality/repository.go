package locality

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
	ErrBadRequest           = errors.New("bad request")
	ErrNotFound             = errors.New("locality not found")
	ErrInternal             = errors.New("database internal error")
	ErrForeignKeyConstraint = errors.New("a column table constraint fails")
	ErrAlreadyExists        = errors.New("id already exists")
)

// Repository encapsulates the storage of a Locality.
type Repository interface {
	ReportCarries(ctx context.Context) ([]domain.ReportCarries, error)
	ReportCarriesByLocationID(ctx context.Context, id string) ([]domain.ReportCarries, error)
	Exists(ctx context.Context, id string) bool
	Save(ctx context.Context, l domain.Locality) (string, error)
	Get(ctx context.Context, id string) (domain.Locality, error)
	ReportSellers(ctx context.Context) ([]domain.ReportSellers, error)
	ReportSellersByLocationID(ctx context.Context, id string) ([]domain.ReportSellers, error)
}

type repository struct {
	db *sql.DB
}

const (
	GET_CARRIES                     = "SELECT localities.id, localities.locality_name, COUNT(carries.locality_id) AS carries_count FROM carries RIGHT JOIN localities ON carries.locality_id = localities.id GROUP BY carries.locality_id, localities.locality_name, localities.id;"
	GET_CARRIES_BY_LOCATION_ID      = "SELECT localities.id, localities.locality_name, COUNT(carries.locality_id) AS carries_count FROM carries RIGHT JOIN localities ON carries.locality_id = localities.id WHERE localities.id = ? GROUP BY carries.locality_id, localities.locality_name, localities.id;"
	GET_SELLERS                     = "SELECT localities.id, localities.locality_name, COUNT(sellers.locality_id) AS sellers_count FROM sellers RIGHT JOIN localities ON sellers.locality_id = localities.id GROUP BY sellers.locality_id, localities.locality_name, localities.id;"
	GET_SELLERS_BY_LOCATION_ID      = "SELECT localities.id, localities.locality_name, COUNT(sellers.locality_id) AS sellers_count FROM sellers RIGHT JOIN localities ON sellers.locality_id = localities.id WHERE localities.id = ? GROUP BY sellers.locality_id, localities.locality_name, localities.id;"
	SAVE_LOCALITY                   = "INSERT INTO localities (id, locality_name, province_name, country_name) VALUES (?, ?, ?, ?)"
	EXIST_LOCALITY                  = "SELECT id FROM localities WHERE id=?"
	GET_LOCALITY                    = "SELECT id, locality_name, province_name, country_name FROM localities WHERE id=?;"
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

func (r *repository) ReportCarries(ctx context.Context) ([]domain.ReportCarries, error) {
	rows, err := r.db.Query(GET_CARRIES)
	if err != nil {
		return nil, ErrInternal
	}

	var report []domain.ReportCarries

	for rows.Next() {
		row := domain.ReportCarries{}
		err = rows.Scan(&row.LocalityID, &row.LocalityName, &row.CarriesCount)
		if err != nil {
			return nil, ErrInternal
		}
		report = append(report, row)
	}
	return report, nil
}

func (r *repository) ReportCarriesByLocationID(ctx context.Context, id string) ([]domain.ReportCarries, error) {
	row := r.db.QueryRow(GET_CARRIES_BY_LOCATION_ID, id)
	report := domain.ReportCarries{}
	err := row.Scan(&report.LocalityID, &report.LocalityName, &report.CarriesCount)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, ErrInternal
		}
	}

	return []domain.ReportCarries{report}, nil
}

func (r *repository) ReportSellers(ctx context.Context) ([]domain.ReportSellers, error) {
	rows, err := r.db.Query(GET_SELLERS)
	if err != nil {
		return nil, ErrInternal
	}

	var report []domain.ReportSellers

	for rows.Next() {
		row := domain.ReportSellers{}
		err = rows.Scan(&row.LocalityID, &row.LocalityName, &row.SellersCount)
		if err != nil {
			return nil, ErrInternal
		}
		report = append(report, row)
	}
	return report, nil
}

func (r *repository) ReportSellersByLocationID(ctx context.Context, id string) ([]domain.ReportSellers, error) {
	row := r.db.QueryRow(GET_SELLERS_BY_LOCATION_ID, id)
	report := domain.ReportSellers{}

	err := row.Scan(&report.LocalityID, &report.LocalityName, &report.SellersCount)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, ErrInternal
		}
	}

	return []domain.ReportSellers{report}, nil
}

func (r *repository) Exists(ctx context.Context, id string) bool {
	rows := r.db.QueryRow(EXIST_LOCALITY, id)
	err := rows.Scan(&id)
	return err == nil
}

func (r *repository) Save(ctx context.Context, l domain.Locality) (string, error) {
	stmt, err := r.db.Prepare(SAVE_LOCALITY)
	if err != nil {
		return "0", err
	}

	_, err = stmt.Exec(l.ID, l.LocalityName, l.ProvinceName, l.CountryName)
	if err != nil {
		mysqlError, ok := err.(*mysql.MySQLError)
		if ok {
			switch mysqlError.Number {
			case MySqlNumberForeignKeyConstraint: // error on foreign key
				return "0", ErrForeignKeyConstraint
			}
		}
		return "0", ErrInternal
	}

	return l.ID, nil
}

func (r *repository) Get(ctx context.Context, id string) (domain.Locality, error) {
	row := r.db.QueryRow(GET_LOCALITY, id)
	l := domain.Locality{}
	err := row.Scan(&l.ID, &l.LocalityName, &l.ProvinceName, &l.CountryName)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return domain.Locality{}, ErrNotFound
		default:
			return domain.Locality{}, ErrInternal
		}
	}

	return l, nil
}
