package section

import (
	"context"
	"database/sql"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/go-sql-driver/mysql"
)

var (
	ErrNotFound            = errors.New("section not found")
	ErrAlreadyExists       = errors.New("section code already exists")
	ErrInternal            = errors.New("database internal error")
	ErrForeignNotFoundCode = 1452
	ErrForeignNotFound     = errors.New("the given id does not have a warehouse atached to it")
)

const (
	GetAllSections     = `SELECT * FROM sections;`
	GetSection         = `SELECT * FROM sections WHERE id=?;`
	ExistsSection      = `SELECT section_number FROM sections WHERE section_number=?;`
	SaveSection        = `INSERT INTO sections (section_number, current_temperature, minimum_temperature, current_capacity, minimum_capacity, maximum_capacity, warehouse_id, id_product_type) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`
	UpdateSection      = `UPDATE sections SET section_number=?, current_temperature=?, minimum_temperature=?, current_capacity=?, minimum_capacity=?, maximum_capacity=?, warehouse_id=?, id_product_type=? WHERE id=?;`
	DeleteSection      = `DELETE FROM sections WHERE id=?;`
	ProductsBySections = `SELECT s.id, s.section_number, IFNULL(sum(pb.current_quantity), 0) as products_count FROM product_batches as pb
							RIGHT JOIN sections as s ON s.id = pb.section_id
							GROUP BY s.id;`
	ProductsBySection = `SELECT s.id, s.section_number, IFNULL(sum(pb.current_quantity), 0) as products_count FROM product_batches as pb
							RIGHT JOIN sections as s ON s.id = pb.section_id
							WHERE s.id = ?
							GROUP BY s.id;`
)

func init() {
	logging.InitLog(nil)
}

// Repository encapsulates the storage of a section.
type Repository interface {
	GetAll(ctx context.Context) ([]domain.Section, error)
	Get(ctx context.Context, id int) (domain.Section, error)
	Exists(ctx context.Context, cid int) bool
	Save(ctx context.Context, s domain.Section) (int, error)
	Update(ctx context.Context, s domain.Section) error
	Delete(ctx context.Context, id int) error
	GetProductsBySections(ctx context.Context) ([]domain.ProductsBySection, error)
	GetProductsBySection(ctx context.Context, sectionID int) ([]domain.ProductsBySection, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAll(ctx context.Context) ([]domain.Section, error) {
	rows, err := r.db.Query(GetAllSections)
	if err != nil {
		logging.Log(err)
		return nil, ErrInternal
	}

	var sections []domain.Section

	for rows.Next() {
		s := domain.Section{}
		_ = rows.Scan(&s.ID, &s.SectionNumber, &s.CurrentTemperature, &s.MinimumTemperature, &s.CurrentCapacity, &s.MinimumCapacity, &s.MaximumCapacity, &s.WarehouseID, &s.ProductTypeID)
		sections = append(sections, s)
	}

	return sections, nil
}

func (r *repository) Get(ctx context.Context, id int) (domain.Section, error) {
	row := r.db.QueryRow(GetSection, id)
	s := domain.Section{}
	err := row.Scan(&s.ID, &s.SectionNumber, &s.CurrentTemperature, &s.MinimumTemperature, &s.CurrentCapacity, &s.MinimumCapacity, &s.MaximumCapacity, &s.WarehouseID, &s.ProductTypeID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			logging.Log(err)
			return domain.Section{}, ErrNotFound
		default:
			logging.Log(err)
			return domain.Section{}, ErrInternal
		}
	}

	return s, nil
}

func (r *repository) Exists(ctx context.Context, sectionNumber int) bool {
	row := r.db.QueryRow(ExistsSection, sectionNumber)
	err := row.Scan(&sectionNumber)
	return err == nil
}

func (r *repository) Save(ctx context.Context, s domain.Section) (int, error) {
	stmt, err := r.db.Prepare(SaveSection)
	if err != nil {
		logging.Log(err)
		return 0, ErrInternal
	}

	res, err := stmt.Exec(&s.SectionNumber, &s.CurrentTemperature, &s.MinimumTemperature, &s.CurrentCapacity, &s.MinimumCapacity, &s.MaximumCapacity, &s.WarehouseID, &s.ProductTypeID)
	if err != nil {
		message, ok := err.(*mysql.MySQLError)
		if ok {
			switch int(message.Number) {
			case ErrForeignNotFoundCode:
				logging.Log(err)
				return 0, ErrForeignNotFound
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

func (r *repository) Update(ctx context.Context, s domain.Section) error {
	stmt, err := r.db.Prepare(UpdateSection)
	if err != nil {
		logging.Log(err)
		return ErrInternal
	}

	res, err := stmt.Exec(&s.SectionNumber, &s.CurrentTemperature, &s.MinimumTemperature, &s.CurrentCapacity, &s.MinimumCapacity, &s.MaximumCapacity, &s.WarehouseID, &s.ProductTypeID, &s.ID)
	if err != nil {
		logging.Log(err)
		return ErrInternal
	}

	_, err = res.RowsAffected()
	if err != nil {
		logging.Log(err)
		return ErrInternal
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	stmt, err := r.db.Prepare(DeleteSection)
	if err != nil {
		logging.Log(err)
		return ErrInternal
	}

	res, err := stmt.Exec(id)
	if err != nil {
		logging.Log(err)
		return ErrInternal
	}

	affect, err := res.RowsAffected()
	if err != nil {
		logging.Log(err)
		return ErrInternal
	}

	if affect < 1 {
		logging.Log(err)
		return ErrNotFound
	}

	return nil
}

func (r *repository) GetProductsBySections(ctx context.Context) ([]domain.ProductsBySection, error) {
	rows, err := r.db.Query(ProductsBySections)
	if err != nil {
		logging.Log(err)
		return nil, ErrInternal
	}

	var productsBySections []domain.ProductsBySection

	for rows.Next() {
		pByS := domain.ProductsBySection{}
		err = rows.Scan(&pByS.SectionID, &pByS.SectionNumber, &pByS.ProductsCount)
		if err != nil {
			logging.Log(err)
			return nil, ErrInternal
		}
		productsBySections = append(productsBySections, pByS)
	}

	return productsBySections, nil
}

func (r *repository) GetProductsBySection(ctx context.Context, sectionID int) ([]domain.ProductsBySection, error) {
	row := r.db.QueryRow(ProductsBySection, sectionID)

	var productsBySections []domain.ProductsBySection

	pByS := domain.ProductsBySection{}
	err := row.Scan(&pByS.SectionID, &pByS.SectionNumber, &pByS.ProductsCount)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			logging.Log(err)
			return nil, ErrNotFound
		default:
			logging.Log(err)
			return nil, ErrInternal
		}
	}
	productsBySections = append(productsBySections, pByS)

	return productsBySections, nil
}
