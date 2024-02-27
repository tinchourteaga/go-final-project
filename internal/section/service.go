package section

import (
	"context"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
)

type Service interface {
	// GetAll returns all the sections that exist inside the repository
	GetAll(context.Context) ([]domain.Section, error)
	// Get returns the section with the specified ID in the repository, if it exists
	Get(c context.Context, id int) (domain.Section, error)
	// Create saves the specified section inside the repository
	Create(c context.Context, section domain.Section) (domain.Section, error)
	// Update updates the section with the specified data in the repository, if it exists
	Update(c context.Context, section domain.Section) (domain.Section, error)
	// Delete deletes a section with the specified ID from the repository
	Delete(c context.Context, id int) error
	// Exists checks if a section with the specified section number exists inside the repository
	Exists(c context.Context, sectionNumber int) error
	// GetSectionProducts returns the amount of products in each section, if the section id given = 0, or only of the section with the same id, if one is given
	GetSectionProducts(c context.Context, sectionID int) ([]domain.ProductsBySection, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) GetAll(c context.Context) ([]domain.Section, error) {
	return s.repository.GetAll(c)
}

func (s *service) Get(c context.Context, id int) (domain.Section, error) {
	return s.repository.Get(c, id)
}

// Create returns the created section if successful, or a error if it failed
func (s *service) Create(c context.Context, section domain.Section) (domain.Section, error) {
	if err := s.Exists(c, section.SectionNumber); err != nil {
		logging.Log(err)
		return domain.Section{}, err
	}
	id, err := s.repository.Save(c, section)
	if err != nil {
		logging.Log(err)
		return domain.Section{}, err
	}
	section.ID = id
	return section, nil
}

// Update returns the updated section if successful, or a error if it failed
// if a section with the given id doesn`t exist, an error is returned
// if the sectionNumber is not unique (with exception to the section currently updating), a error is returned
// only the values not in a null state are updated
func (s *service) Update(c context.Context, newSection domain.Section) (domain.Section, error) {
	section, err := s.Get(c, newSection.ID)
	if err != nil {
		logging.Log(err)
		return domain.Section{}, err
	}
	if newSection.SectionNumber != 0 && newSection.SectionNumber != section.SectionNumber {
		if err := s.Exists(c, newSection.SectionNumber); err != nil {
			logging.Log(err)
			return domain.Section{}, err
		}
		section.SectionNumber = newSection.SectionNumber
	}
	if newSection.CurrentTemperature > -273 {
		section.CurrentTemperature = newSection.CurrentTemperature
	}
	if newSection.MinimumTemperature > -273 {
		section.MinimumTemperature = newSection.MinimumTemperature
	}
	if newSection.CurrentCapacity != 0 {
		section.CurrentCapacity = newSection.CurrentCapacity
	}
	if newSection.MinimumCapacity != 0 {
		section.MinimumCapacity = newSection.MinimumCapacity
	}
	if newSection.MaximumCapacity != 0 {
		section.MaximumCapacity = newSection.MaximumCapacity
	}
	if newSection.WarehouseID != 0 {
		section.WarehouseID = newSection.WarehouseID
	}
	if newSection.ProductTypeID != 0 {
		section.ProductTypeID = newSection.ProductTypeID
	}
	err = s.repository.Update(c, section)
	if err != nil {
		logging.Log(err)
		return domain.Section{}, err
	}
	return section, nil
}

// Delete returns an error if the deletion of the section failed
// if a section with the given id doesn`t exist, an error is returned
func (s *service) Delete(c context.Context, id int) error {
	_, err := s.repository.Get(c, id)
	if err != nil {
		logging.Log(err)
		return err
	}
	err = s.repository.Delete(c, id)
	logging.Log(err)
	return err
}

func (s *service) Exists(c context.Context, sectionNumber int) error {
	if s.repository.Exists(c, sectionNumber) {
		return ErrAlreadyExists
	}
	return nil
}

func (s *service) GetSectionProducts(c context.Context, sectionID int) ([]domain.ProductsBySection, error) {
	if sectionID == 0 {
		return s.repository.GetProductsBySections(c)
	}
	return s.repository.GetProductsBySection(c, sectionID)
}
