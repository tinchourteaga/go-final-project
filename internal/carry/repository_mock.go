package carry

import (
	"context"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
)

type MockRepo struct {
	mockError error
}

func (r *MockRepo) Save(ctx context.Context, carry domain.Carry) (int, error) {
	if r.mockError != nil {
		return 0, r.mockError
	}
	return 1, nil
}
