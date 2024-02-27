package productbatch

import (
	"context"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
)

type MockRepository struct {
	mockProductBatches []domain.ProductBatch
	mockError          error
}

func (r *MockRepository) Save(ctx context.Context, pb domain.ProductBatch) (int, error) {
	if r.mockError != nil {
		return 0, r.mockError
	}
	id := len(r.mockProductBatches)
	pb.ID = id + 1
	r.mockProductBatches = append(r.mockProductBatches, pb)
	return pb.ID, nil
}
func (r *MockRepository) Exists(ctx context.Context, cid int) bool {
	return r.mockError == ErrAlreadyExists
}
