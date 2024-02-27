package productbatch

import (
	"context"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
)

type Service interface {
	Create(c context.Context, pb domain.ProductBatch) (domain.ProductBatch, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) Create(c context.Context, pb domain.ProductBatch) (domain.ProductBatch, error) {
	id, err := s.repository.Save(c, pb)
	if err != nil {
		logging.Log(err)
		return domain.ProductBatch{}, err
	}
	pb.ID = id
	return pb, nil
}
