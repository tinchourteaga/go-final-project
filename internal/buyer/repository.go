package buyer

import (
	"context"
	"database/sql"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
)

// Errors
var (
	ErrNotFound      = errors.New("buyer not found")
	ErrAlreadyExists = errors.New("card_number_id already exists")
	ErrInternal      = errors.New("database internal error")
	ErrDataLong      = errors.New("a field exceeds the maximum length")
)

const (
	GET_ALL_QUERY        = "SELECT * FROM buyers;"
	GET_BY_ID_QUERY      = "SELECT * FROM buyers WHERE id = ?;"
	EXISTS_QUERY         = "SELECT card_number_id FROM buyers WHERE card_number_id=?;"
	INSERT_QUERY         = "INSERT INTO buyers(card_number_id,first_name,last_name) VALUES (?,?,?);"
	UPDATE_QUERY         = "UPDATE buyers SET first_name=?, last_name=?, card_number_id=?  WHERE id=?;"
	DELETE_QUERY         = "DELETE FROM buyers WHERE id = ?;"
	MySqlNumberDataLong  = 1406
	MySqlNumberDuplicate = 1062
)

// Repository encapsulates the storage of a buyer.
type Repository interface {
	GetAll(ctx context.Context) ([]domain.Buyer, error)
	Get(ctx context.Context, id int) (domain.Buyer, error)
	Exists(ctx context.Context, cardNumberID string) bool
	Save(ctx context.Context, b domain.Buyer) (int, error)
	Update(ctx context.Context, b domain.Buyer) error
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

func (r *repository) GetAll(ctx context.Context) ([]domain.Buyer, error) {
	rows, err := r.db.Query(GET_ALL_QUERY)
	if err != nil {
		logging.Log(err)
		return nil, err
	}

	var buyers []domain.Buyer

	for rows.Next() {
		b := domain.Buyer{}
		_ = rows.Scan(&b.ID, &b.CardNumberID, &b.FirstName, &b.LastName)
		buyers = append(buyers, b)
	}

	return buyers, nil
}

func (r *repository) Get(ctx context.Context, id int) (domain.Buyer, error) {
	row := r.db.QueryRow(GET_BY_ID_QUERY, id)
	b := domain.Buyer{}
	err := row.Scan(&b.ID, &b.CardNumberID, &b.FirstName, &b.LastName)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			logging.Log(ErrNotFound)
			return domain.Buyer{}, ErrNotFound
		default:
			logging.Log(ErrInternal)
			return domain.Buyer{}, ErrInternal
		}
	}

	return b, nil
}

func (r *repository) Exists(ctx context.Context, cardNumberID string) bool {
	row := r.db.QueryRow(EXISTS_QUERY, cardNumberID)
	err := row.Scan(&cardNumberID)
	return err == nil
}

func (r *repository) Save(ctx context.Context, b domain.Buyer) (int, error) {
	stmt, err := r.db.Prepare(INSERT_QUERY)
	if err != nil {
		logging.Log(err)
		return 0, err
	}

	res, err := stmt.Exec(&b.CardNumberID, &b.FirstName, &b.LastName)
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

func (r *repository) Update(ctx context.Context, b domain.Buyer) error {
	stmt, err := r.db.Prepare(UPDATE_QUERY)
	if err != nil {
		logging.Log(err)
		return err
	}

	res, err := stmt.Exec(&b.FirstName, &b.LastName, &b.CardNumberID, &b.ID)
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
	stmt, err := r.db.Prepare(DELETE_QUERY)
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
		logging.Log(err)
		return ErrNotFound
	}

	return nil
}
