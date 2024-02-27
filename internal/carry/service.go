package carry

import (
	"context"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
)

// Service provides the public methods of a carry service.
type Service interface {
	Save(ctx context.Context, CID string, CompanyName string, Address string, Telephone string, Locality_id string) (domain.Carry, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

// Save returns the created carry provided by the repository if succesful.
// if the cid is not unique, a error is returned.
// any other error encountered is also returned.
func (s *service) Save(ctx context.Context, cid string, companyName string, address string, telephone string, locality_id string) (domain.Carry, error) {
	carry := domain.Carry{
		CID:         cid,
		CompanyName: companyName,
		Address:     address,
		Telephone:   telephone,
		Locality_id: locality_id,
	}
	carryID, err := s.repository.Save(ctx, carry)
	if err != nil {
		logging.Log(err)
		return domain.Carry{}, err
	}

	carry.ID = carryID
	return carry, nil
}
