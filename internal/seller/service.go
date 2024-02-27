package seller

import (
	"context"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
)

var (
	ServiceErrNotFound           = errors.New("seller not found")
	ServiceErrInternal           = errors.New("internal error")
	ServiceErrAlreadyExists      = errors.New("seller code already exists")
	ServiceErrForeignKeyNotFound = errors.New("locality not found")
)

// Service represents a service layer for Seller
type Service interface {
	GetAll(ctx context.Context) ([]domain.Seller, error)
	Create(ctx context.Context, sell domain.Seller) (domain.Seller, error)
	Get(ctx context.Context, id int) (domain.Seller, error)
	Delete(ctx context.Context, id int) error
	Update(context.Context, int, *int, *string, *string, *string, *string) (domain.Seller, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

// GetAll returns all sellers from db
func (s *service) GetAll(ctx context.Context) (sellers []domain.Seller, err error) {
	sellers, err = s.repository.GetAll(ctx)
	if err != nil {
		logging.Log(err)
		return
	}
	return
}

// Create save a seller in the db and returns that seller
func (s *service) Create(ctx context.Context, sell domain.Seller) (seller domain.Seller, err error) {
	if s.repository.Exists(ctx, sell.CID) {
		err = ErrAlreadyExists
		return
	}

	seller = domain.Seller{
		CID:         sell.CID,
		CompanyName: sell.CompanyName,
		Address:     sell.Address,
		Telephone:   sell.Telephone,
		Locality_id: sell.Locality_id,
	}

	sellerID, err := s.repository.Save(ctx, seller)
	if err != nil {
		logging.Log(err)
		return domain.Seller{}, err
	}

	seller.ID = sellerID

	return
}

// Get returns a seller by id
func (s *service) Get(ctx context.Context, id int) (seller domain.Seller, err error) {
	seller, err = s.repository.Get(ctx, id)
	if err != nil {
		logging.Log(err)
		return
	}
	return
}

// Delete receive an id and delete from db
func (s *service) Delete(ctx context.Context, id int) (err error) {
	err = s.repository.Delete(ctx, id)
	if err != nil {
		logging.Log(err)
		return
	}
	return
}

// Update receive an id and update all fields from a seller
func (s *service) Update(ctx context.Context, id int, cid *int, companyName, address, telephone *string, locality_id *string) (sellerToUpdate domain.Seller, err error) {
	sellerToUpdate, err = s.repository.Get(ctx, id)
	if err != nil {
		return domain.Seller{}, err
	}

	if cid != nil {
		// if to compare if already exist cid and return conflict (409)
		if sellerToUpdate.CID != *cid {
			if s.repository.Exists(ctx, *cid) {
				return domain.Seller{}, ErrAlreadyExists
			}
		}

		sellerToUpdate.CID = *cid
	}

	if companyName != nil {
		sellerToUpdate.CompanyName = *companyName
	}

	if address != nil {
		sellerToUpdate.Address = *address
	}

	if telephone != nil {
		sellerToUpdate.Telephone = *telephone
	}

	if locality_id != nil {
		sellerToUpdate.Locality_id = *locality_id
	}

	err = s.repository.Update(ctx, sellerToUpdate)
	if err != nil {
		switch err {
		case ErrForeignKeyConstraint:
			logging.Log(ServiceErrForeignKeyNotFound)
			return domain.Seller{}, ServiceErrForeignKeyNotFound
		case ErrAlreadyExists:
			// This is in case we implement unique with product_code (not happening on Sprint III)
			return domain.Seller{}, ServiceErrAlreadyExists
		default:
			logging.Log(ServiceErrInternal)
			return domain.Seller{}, ServiceErrInternal
		}
	}

	return
}
