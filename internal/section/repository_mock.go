package section

import (
	"context"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
)

type MockRepository struct {
	mockSections			[]domain.Section
	mockProductsBySection	[]domain.ProductsBySection
	mockError				error
	mockGetError			error
}

func (r *MockRepository) GetAll(ctx context.Context) ([]domain.Section, error) {
	return r.mockSections, nil
}

func (r *MockRepository) Get(ctx context.Context, id int) (domain.Section, error) {
	if r.mockGetError != nil {
		return domain.Section{}, r.mockGetError
	}
	return r.mockSections[0], nil
}

func (r *MockRepository) Exists(ctx context.Context, cid int) bool {
	return r.mockError == ErrAlreadyExists
}

func (r *MockRepository) Save(ctx context.Context, s domain.Section) (int, error) {
	if r.mockError != nil {
		return 0, r.mockError
	}
	id := len(r.mockSections)
	s.ID = id + 1
	r.mockSections = append(r.mockSections, s)
	return s.ID, nil
}

func (r *MockRepository) Update(ctx context.Context, s domain.Section) error {
	if r.mockError != nil {
		return r.mockError
	}
	r.mockSections[0] = s
	return nil
}

func (r *MockRepository) Delete(ctx context.Context, id int) error {
	if r.mockError != nil {
		return r.mockError
	}
	r.mockSections = r.mockSections[1:]
	return nil
}

func (r *MockRepository) GetProductsBySections(ctx context.Context) ([]domain.ProductsBySection, error) {
	if r.mockError != nil {
		return nil, r.mockError
	}
	return r.mockProductsBySection, nil
}

func (r *MockRepository) GetProductsBySection(ctx context.Context, sectionID int) ([]domain.ProductsBySection, error) {
	if r.mockError != nil {
		return nil, r.mockError
	}
	return []domain.ProductsBySection{r.mockProductsBySection[0]}, nil
}
