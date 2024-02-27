package productbatch

import (
	"context"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
)

type MockService struct {
	MockProductBatches []domain.ProductBatch
	MockError          error
}

func (s *MockService) Create(c context.Context, pb domain.ProductBatch) (domain.ProductBatch, error) {
	if s.MockError != nil {
		return domain.ProductBatch{}, s.MockError
	}
	id := len(s.MockProductBatches) + 1
	pb.ID = id
	s.MockProductBatches = append(s.MockProductBatches, pb)
	return s.MockProductBatches[id-1], nil
}
func (s *MockService) Exists(c context.Context, productBatchNumber int) error {
	return nil
}
