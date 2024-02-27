package locality

import (
	"context"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
)

// Service represents a service layer for Locality
type Service interface {
	ReportCarries(ctx context.Context, locality_ID *string) ([]domain.ReportCarries, error)
	Create(ctx context.Context, locality domain.Locality) (domain.Locality, error)
	Get(ctx context.Context, id string) (domain.Locality, error)
	ReportSellers(ctx context.Context, locality_ID *string) ([]domain.ReportSellers, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) ReportCarries(ctx context.Context, locality_ID *string) ([]domain.ReportCarries, error) {
	if locality_ID != nil {
		return s.repository.ReportCarriesByLocationID(ctx, *locality_ID)
	}
	return s.repository.ReportCarries(ctx)
}

func (s *service) ReportSellers(ctx context.Context, locality_ID *string) (report []domain.ReportSellers, err error) {
	if locality_ID != nil {
		report, err = s.repository.ReportSellersByLocationID(ctx, *locality_ID)
		return
	}

	report, err = s.repository.ReportSellers(ctx)
	return
}

// Create save a locality in the db and returns that locality
func (s *service) Create(ctx context.Context, locality domain.Locality) (l domain.Locality, err error) {
	if s.repository.Exists(ctx, locality.ID) {
		err = ErrAlreadyExists
		return
	}

	l = domain.Locality{
		ID:           locality.ID,
		LocalityName: locality.LocalityName,
		ProvinceName: locality.ProvinceName,
		CountryName:  locality.CountryName,
	}

	_, err = s.repository.Save(ctx, l)
	if err != nil {
		logging.Log(err)
		return domain.Locality{}, err
	}

	return
}

// Get returns a locality by id
func (s *service) Get(ctx context.Context, id string) (l domain.Locality, err error) {
	l, err = s.repository.Get(ctx, id)
	if err != nil {
		logging.Log(err)
	}
	return
}
