package buyer

import (
	"context"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
)

type Service interface {
	GetAll(ctx context.Context) ([]domain.Buyer, error)
	Save(ctx context.Context, b domain.Buyer) (domain.Buyer, error)
	Exists(ctx context.Context, cardNumberID string) bool
	Get(ctx context.Context, id int) (domain.Buyer, error)
	Update(ctx context.Context, b domain.Buyer) (domain.Buyer, error)
	Delete(ctx context.Context, id int) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

// GetAll returns a List of buyers if successful, or a error if it failed
func (s *service) GetAll(ctx context.Context) ([]domain.Buyer, error) {
	return s.repository.GetAll(ctx)
}

// Save returns the created a buyer if successful, or a error if it failed
// if the buyer is already exist, a error is returned
func (s *service) Save(ctx context.Context, b domain.Buyer) (domain.Buyer, error) {

	if s.repository.Exists(ctx, b.CardNumberID) {
		logging.Log(ErrAlreadyExists)
		return domain.Buyer{}, ErrAlreadyExists
	}

	buyer := domain.Buyer{
		CardNumberID: b.CardNumberID,
		FirstName:    b.FirstName,
		LastName:     b.LastName,
	}

	buyerId, errSave := s.repository.Save(ctx, b)
	if errSave != nil {
		logging.Log(errSave)
		return domain.Buyer{}, errSave
	}

	buyer.ID = buyerId

	return buyer, nil
}

// Exists returns true if the buyer is already exist, or false if it doesn´t
func (s *service) Exists(ctx context.Context, cardNumberID string) bool {
	return s.repository.Exists(ctx, cardNumberID)
}

// Get returns a buyer if successful, or a error if it failed
func (s *service) Get(ctx context.Context, id int) (domain.Buyer, error) {
	return s.repository.Get(ctx, id)
}

// Update returns the updated buyer if successful, or a error if it failed
// if a buyer with the given id doesn`t exist, an error is returned
// only the values not in a null state are updated
func (s *service) Update(ctx context.Context, b domain.Buyer) (domain.Buyer, error) {
	id := int(b.ID)
	data, err := s.repository.Get(ctx, id)
	if err != nil {
		logging.Log(err)
		return domain.Buyer{}, err
	}

	if b.FirstName != "" {
		data.FirstName = b.FirstName
	}

	if b.LastName != "" {
		data.LastName = b.LastName
	}

	if b.CardNumberID != "" {
		data.CardNumberID = b.CardNumberID
	}

	if err := s.repository.Update(ctx, data); err != nil {
		logging.Log(err)
		return domain.Buyer{}, err
	}

	return data, nil
}

// Delete returns an error if the deletion of the section failed
// if a section with the given id doesn`t exist, an error is returned
func (s *service) Delete(ctx context.Context, id int) error {
	data, err := s.repository.Get(ctx, id)

	if err != nil {
		logging.Log(err)
		return err
	}
	if data.ID == 0 {
		logging.Log(ErrNotFound)
		return ErrNotFound
	}
	return s.repository.Delete(ctx, id)
}
