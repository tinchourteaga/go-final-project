package product

import (
	"context"
	"database/sql"
	"errors"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/go-sql-driver/mysql"
	"log"
)

var (
	RepositoryErrNotFound             = errors.New("product not found in database")
	RepositoryErrInternal             = errors.New("database internal error")
	RepositoryErrForeignKeyConstraint = errors.New("a foreign key constraint fails")
	RepositoryErrAlreadyExists        = errors.New("product code already exists")
)

const (
	SaveProduct                     = "INSERT INTO products(description, expiration_rate, freezing_rate, height, lenght, netweight, product_code, recommended_freezing_temperature, width, id_product_type, id_seller) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	GetProduct                      = "SELECT id, description, expiration_rate, freezing_rate, height, lenght, netweight, product_code, recommended_freezing_temperature, width, id_product_type, id_seller FROM products WHERE id = ?;"
	GetAllProducts                  = "SELECT id, description, expiration_rate, freezing_rate, height, lenght, netweight, product_code, recommended_freezing_temperature, width, id_product_type, id_seller FROM products;"
	UpdateProduct                   = "UPDATE products SET description = ?, expiration_rate = ?, freezing_rate = ?, height = ?, lenght = ?, netweight = ?, product_code = ?, recommended_freezing_temperature = ?, width = ?, id_product_type = ?, id_seller = ? WHERE id = ?"
	DeleteProduct                   = "DELETE FROM products WHERE id = ?"
	ExistsProduct                   = "SELECT product_code FROM products WHERE product_code = ?;"
	MySqlNumberForeignKeyConstraint = 1452
	MySqlNumberDuplicateEntry       = 1062
)

// Repository encapsulates the storage of a Product.
type Repository interface {
	GetAll(ctx context.Context) ([]domain.Product, error)
	Get(ctx context.Context, id int) (domain.Product, error)
	Exists(ctx context.Context, productCode string) bool
	Save(ctx context.Context, p domain.Product) (int, error)
	Update(ctx context.Context, p domain.Product) error
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

func CloseStmt(stmt *sql.Stmt) {
	err := stmt.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (r *repository) Exists(ctx context.Context, productCode string) bool {
	row := r.db.QueryRowContext(ctx, ExistsProduct, productCode)
	errScan := row.Scan(&productCode)
	return errScan == nil
}

func (r *repository) GetAll(ctx context.Context) ([]domain.Product, error) {
	rows, errQuery := r.db.QueryContext(ctx, GetAllProducts)
	if errQuery != nil {
		logging.Log(errQuery)
		return nil, errQuery
	}
	var products []domain.Product
	for rows.Next() {
		p := domain.Product{}
		_ = rows.Scan(&p.ID, &p.Description, &p.ExpirationRate, &p.FreezingRate, &p.Height, &p.Length, &p.NetWeight, &p.ProductCode, &p.RecommendedFreezingTemperature, &p.Width, &p.ProductTypeID, &p.SellerID)
		products = append(products, p)
	}
	return products, nil
}

func (r *repository) Get(ctx context.Context, id int) (domain.Product, error) {
	row := r.db.QueryRowContext(ctx, GetProduct, id)
	p := domain.Product{}
	errScan := row.Scan(&p.ID, &p.Description, &p.ExpirationRate, &p.FreezingRate, &p.Height, &p.Length, &p.NetWeight, &p.ProductCode, &p.RecommendedFreezingTemperature, &p.Width, &p.ProductTypeID, &p.SellerID)
	if errScan != nil {
		logging.Log(errScan)
		switch errScan {
		case sql.ErrNoRows:
			return domain.Product{}, RepositoryErrNotFound
		default:
			return domain.Product{}, RepositoryErrInternal
		}
	}
	return p, nil
}

func (r *repository) Save(ctx context.Context, p domain.Product) (int, error) {
	stmt, errPrepare := r.db.PrepareContext(ctx, SaveProduct)
	if errPrepare != nil {
		logging.Log(errPrepare)
		return 0, errPrepare
	}
	defer CloseStmt(stmt)
	result, errExec := stmt.ExecContext(ctx, p.Description, p.ExpirationRate, p.FreezingRate, p.Height, p.Length, p.NetWeight, p.ProductCode, p.RecommendedFreezingTemperature, p.Width, p.ProductTypeID, p.SellerID)
	if errExec != nil {
		logging.Log(errExec)
		message, ok := errExec.(*mysql.MySQLError)
		if ok {
			switch message.Number {
			case MySqlNumberForeignKeyConstraint:
				return 0, RepositoryErrForeignKeyConstraint
			case MySqlNumberDuplicateEntry:
				// This is in case we implement unique with product_code (not happening on Sprint III)
				return 0, RepositoryErrAlreadyExists
			}
		}
		return 0, errExec
	}
	id, errID := result.LastInsertId()
	if errID != nil {
		logging.Log(errID)
		return 0, errID
	}
	return int(id), nil
}

func (r *repository) Update(ctx context.Context, p domain.Product) error {
	stmt, errPrepare := r.db.PrepareContext(ctx, UpdateProduct)
	if errPrepare != nil {
		logging.Log(errPrepare)
		return errPrepare
	}
	defer CloseStmt(stmt)
	result, errExec := stmt.ExecContext(ctx, p.Description, p.ExpirationRate, p.FreezingRate, p.Height, p.Length, p.NetWeight, p.ProductCode, p.RecommendedFreezingTemperature, p.Width, p.ProductTypeID, p.SellerID, p.ID)
	if errExec != nil {
		logging.Log(errExec)
		message, ok := errExec.(*mysql.MySQLError)
		if ok {
			switch message.Number {
			case MySqlNumberForeignKeyConstraint:
				return RepositoryErrForeignKeyConstraint
			case MySqlNumberDuplicateEntry:
				// This is in case we implement unique with product_code (not happening on Sprint III)
				return RepositoryErrAlreadyExists
			}
		}
		return errExec
	}
	_, errAffection := result.RowsAffected()
	if errAffection != nil {
		logging.Log(errAffection)
		return errAffection
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	stmt, errPrepare := r.db.PrepareContext(ctx, DeleteProduct)
	if errPrepare != nil {
		logging.Log(errPrepare)
		return errPrepare
	}
	defer CloseStmt(stmt)
	result, errExec := stmt.ExecContext(ctx, id)
	if errExec != nil {
		logging.Log(errExec)
		return errExec
	}
	affectedRows, errAffection := result.RowsAffected()
	if errAffection != nil {
		logging.Log(errAffection)
		return errAffection
	}
	if affectedRows < 1 {
		logging.Log(RepositoryErrNotFound)
		return RepositoryErrNotFound
	}
	return nil
}
