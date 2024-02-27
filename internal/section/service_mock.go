package section

import (
	"context"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
)

type MockService struct {
	MockSections          []domain.Section
	MockProductsBySection []domain.ProductsBySection
	MockError             error
}

func (s *MockService) GetAll(context.Context) ([]domain.Section, error) {
	if s.MockError != nil {
		return nil, s.MockError
	}
	return s.MockSections, nil
}

func (s *MockService) Get(c context.Context, id int) (domain.Section, error) {
	if s.MockError != nil {
		return domain.Section{}, s.MockError
	}
	return s.MockSections[0], nil
}

func (s *MockService) Create(c context.Context, section domain.Section) (domain.Section, error) {
	if s.MockError != nil {
		return domain.Section{}, s.MockError
	}
	id := len(s.MockSections) + 1
	section.ID = id
	s.MockSections = append(s.MockSections, section)
	return s.MockSections[id-1], nil
}

func (s *MockService) Update(c context.Context, section domain.Section) (domain.Section, error) {
	if s.MockError != nil {
		return domain.Section{}, s.MockError
	}
	if section.SectionNumber != 0 {
		s.MockSections[0].SectionNumber = section.SectionNumber
	}
	if section.CurrentTemperature > -273 {
		s.MockSections[0].CurrentTemperature = section.CurrentTemperature
	}
	if section.MinimumTemperature > -273 {
		s.MockSections[0].MinimumTemperature = section.MinimumTemperature
	}
	if section.CurrentCapacity != 0 {
		s.MockSections[0].CurrentCapacity = section.CurrentCapacity
	}
	if section.MinimumCapacity != 0 {
		s.MockSections[0].MinimumCapacity = section.MinimumCapacity
	}
	if section.MaximumCapacity != 0 {
		s.MockSections[0].MaximumCapacity = section.MaximumCapacity
	}
	if section.WarehouseID != 0 {
		s.MockSections[0].WarehouseID = section.WarehouseID
	}
	if section.ProductTypeID != 0 {
		s.MockSections[0].ProductTypeID = section.ProductTypeID
	}
	return s.MockSections[0], nil
}

func (s *MockService) Delete(c context.Context, id int) error {
	if s.MockError != nil {
		return s.MockError
	}
	s.MockSections = s.MockSections[1:]
	return nil
}

func (s *MockService) Exists(c context.Context, sectionNumber int) error {
	return nil
}

func (s *MockService) GetSectionProducts(c context.Context, sectionID int) ([]domain.ProductsBySection, error) {
	if s.MockError != nil {
		return nil, s.MockError
	}
	if sectionID != 0 {
		return []domain.ProductsBySection{s.MockProductsBySection[0]}, nil
	}
	return s.MockProductsBySection, nil
}
