package section

import (
	"context"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/stretchr/testify/assert"
)

// TestCreateOk tests if the service correctly calls the repository to create and return the given section
func TestCreateOk(t *testing.T) {
	// ARANGE
	repository := MockRepository{}
	service := NewService(&repository)
	ctx := new(context.Context)

	expected := domain.Section{
		ID:                 1,
		SectionNumber:      1,
		CurrentTemperature: 2.0,
		MinimumTemperature: -1.0,
		CurrentCapacity:    100,
		MaximumCapacity:    1000,
		MinimumCapacity:    10,
		WarehouseID:        1,
		ProductTypeID:      1,
	}

	//ACT
	result, err := service.Create(*ctx,
		domain.Section{
			ID:                 1,
			SectionNumber:      1,
			CurrentTemperature: 2.0,
			MinimumTemperature: -1.0,
			CurrentCapacity:    100,
			MaximumCapacity:    1000,
			MinimumCapacity:    10,
			WarehouseID:        1,
			ProductTypeID:      1,
		})

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

// TestCreateConflict tests if the service returns the correct error if the repository already exists
func TestCreateConflict(t *testing.T) {
	// ARRANGE
	repository := MockRepository{mockError: ErrAlreadyExists,
		mockSections: []domain.Section{
			{
				ID:                 1,
				SectionNumber:      1,
				CurrentTemperature: 2.0,
				MinimumTemperature: -1.0,
				CurrentCapacity:    100,
				MaximumCapacity:    1000,
				MinimumCapacity:    10,
				WarehouseID:        1,
				ProductTypeID:      1,
			},
		},
	}
	service := NewService(&repository)
	ctx := new(context.Context)

	expected := ErrAlreadyExists

	// ACT
	result, err := service.Create(*ctx,
		domain.Section{
			ID:                 1,
			SectionNumber:      1,
			CurrentTemperature: 2.0,
			MinimumTemperature: -1.0,
			CurrentCapacity:    100,
			MaximumCapacity:    1000,
			MinimumCapacity:    10,
			WarehouseID:        1,
			ProductTypeID:      1,
		})

	// ASSERT
	assert.Empty(t, result)
	assert.EqualError(t, expected, err.Error())

}

func TestCreateInternalErr(t *testing.T) {
	// ARRANGE
	repository := MockRepository{mockError: ErrInternal}
	service := NewService(&repository)
	ctx := new(context.Context)

	expected := ErrInternal

	// ACT
	result, err := service.Create(*ctx, domain.Section{})

	// ASSERT
	assert.Empty(t, result)
	assert.EqualError(t, expected, err.Error())
}

// TestFindAll tests if the service correctly calls the repository to return all sections stored
func TestFindAll(t *testing.T) {
	// ARRANGE
	repository := MockRepository{
		mockSections: []domain.Section{
			{
				ID:                 1,
				SectionNumber:      1,
				CurrentTemperature: 2.0,
				MinimumTemperature: -1.0,
				CurrentCapacity:    100,
				MaximumCapacity:    1000,
				MinimumCapacity:    10,
				WarehouseID:        1,
				ProductTypeID:      1,
			},
			{
				ID:                 2,
				SectionNumber:      2,
				CurrentTemperature: 3.0,
				MinimumTemperature: -1.0,
				CurrentCapacity:    100,
				MaximumCapacity:    1000,
				MinimumCapacity:    10,
				WarehouseID:        1,
				ProductTypeID:      2,
			},
		}}
	service := NewService(&repository)
	ctx := new(context.Context)

	expected := []domain.Section{
		{
			ID:                 1,
			SectionNumber:      1,
			CurrentTemperature: 2.0,
			MinimumTemperature: -1.0,
			CurrentCapacity:    100,
			MaximumCapacity:    1000,
			MinimumCapacity:    10,
			WarehouseID:        1,
			ProductTypeID:      1,
		},
		{
			ID:                 2,
			SectionNumber:      2,
			CurrentTemperature: 3.0,
			MinimumTemperature: -1.0,
			CurrentCapacity:    100,
			MaximumCapacity:    1000,
			MinimumCapacity:    10,
			WarehouseID:        1,
			ProductTypeID:      2,
		},
	}

	// ACT
	result, err := service.GetAll(*ctx)

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

// TestFindByIdNon test if the service returns the correct error when a section with the given id doesn´t exists
func TestFindByIdNonExistent(t *testing.T) {
	// ARRANGE
	repository := MockRepository{mockGetError: ErrNotFound}
	service := NewService(&repository)
	ctx := new(context.Context)

	expected := ErrNotFound

	// ACT
	result, err := service.Get(*ctx, 1)

	// ASSERT
	assert.Empty(t, result)
	assert.EqualError(t, expected, err.Error())
}

// TestFindByIdExistent tests if the service correctly calls the repository to return the section with the given id
func TestFindByIdExistent(t *testing.T) {
	// ARRANGE
	repository := MockRepository{
		mockSections: []domain.Section{
			{
				ID:                 1,
				SectionNumber:      1,
				CurrentTemperature: 2.0,
				MinimumTemperature: -1.0,
				CurrentCapacity:    100,
				MaximumCapacity:    1000,
				MinimumCapacity:    10,
				WarehouseID:        1,
				ProductTypeID:      1,
			},
		},
		mockGetError: nil,
	}
	service := NewService(&repository)
	ctx := new(context.Context)

	expected := domain.Section{
		ID:                 1,
		SectionNumber:      1,
		CurrentTemperature: 2.0,
		MinimumTemperature: -1.0,
		CurrentCapacity:    100,
		MaximumCapacity:    1000,
		MinimumCapacity:    10,
		WarehouseID:        1,
		ProductTypeID:      1,
	}

	// ACT
	result, err := service.Get(*ctx, 1)

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

// TestUpdateExistent tests if the service correctly calls the repository to update the section with the given id both when a value is given for each attribute and when is not
func TestUpdateExistent(t *testing.T) {
	// ARRANGE
	repository := MockRepository{
		mockSections: []domain.Section{
			{
				ID:                 1,
				SectionNumber:      1,
				CurrentTemperature: 2.0,
				MinimumTemperature: -1.0,
				CurrentCapacity:    100,
				MaximumCapacity:    1000,
				MinimumCapacity:    10,
				WarehouseID:        1,
				ProductTypeID:      1,
			},
		},
		mockError: nil,
		mockGetError: nil,
	}
	service := NewService(&repository)
	ctx := new(context.Context)

	expected1 := domain.Section{
		ID:                 1,
		SectionNumber:      1,
		CurrentTemperature: 2.0,
		MinimumTemperature: -1.0,
		CurrentCapacity:    100,
		MaximumCapacity:    1000,
		MinimumCapacity:    10,
		WarehouseID:        1,
		ProductTypeID:      1,
	}
	expected2 := domain.Section{
		ID:                 1,
		SectionNumber:      2,
		CurrentTemperature: 3.0,
		MinimumTemperature: -2.0,
		CurrentCapacity:    90,
		MaximumCapacity:    1100,
		MinimumCapacity:    20,
		WarehouseID:        2,
		ProductTypeID:      2,
	}

	// ACT
	// Check if values of section remain the same after empty values are passed
	// Empty value of CurrentTemperature and MinimumTemperature is absolute zero = -273, which is different from empty value of float64
	result1, err1 := service.Update(*ctx, domain.Section{
		CurrentTemperature: -273,
		MinimumTemperature: -273,
	})
	// Check if values of section change after not-empty values are passed
	result2, err2 := service.Update(*ctx, domain.Section{
		ID:                 1,
		SectionNumber:      2,
		CurrentTemperature: 3.0,
		MinimumTemperature: -2.0,
		CurrentCapacity:    90,
		MaximumCapacity:    1100,
		MinimumCapacity:    20,
		WarehouseID:        2,
		ProductTypeID:      2,
	})

	// ASSERT
	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Equal(t, expected1, result1)
	assert.Equal(t, expected2, result2)
}

// TestUpdateNonExistent tests if the service returns the correct error when a section with the given id doesn´t exists
func TestUpdateNonExistent(t *testing.T) {
	// ARRANGE
	repository := MockRepository{
		mockGetError: ErrNotFound,
	}
	service := NewService(&repository)
	ctx := new(context.Context)

	expected := ErrNotFound

	// ACT
	result, err := service.Update(*ctx, domain.Section{})

	//ASSERT
	assert.Empty(t, result)
	assert.EqualError(t, expected, err.Error())
}

func TestUpdateExistentSectionNumber(t *testing.T) {
	// ARRANGE
	repository := MockRepository{
		mockSections: []domain.Section{
			{
				ID:                 1,
				SectionNumber:      1,
				CurrentTemperature: 2.0,
				MinimumTemperature: -1.0,
				CurrentCapacity:    100,
				MaximumCapacity:    1000,
				MinimumCapacity:    10,
				WarehouseID:        1,
				ProductTypeID:      1,
			},
		},
		mockError: ErrAlreadyExists,
		mockGetError: nil,
	}
	service := NewService(&repository)
	ctx := new(context.Context)

	expected := ErrAlreadyExists

	// ACT
	result, err := service.Update(*ctx, domain.Section{
			ID:                 1,
			SectionNumber:      2,
			CurrentTemperature: 2.0,
			MinimumTemperature: -1.0,
			CurrentCapacity:    100,
			MaximumCapacity:    1000,
			MinimumCapacity:    10,
			WarehouseID:        1,
			ProductTypeID:      1,
		},
	)

	//ASSERT
	assert.Empty(t, result)
	assert.EqualError(t, err, expected.Error())
}

func TestUpdateInternalErr(t *testing.T) {
	// ARRANGE
	repository := MockRepository{
		mockSections: []domain.Section{
			{
				ID:                 1,
				SectionNumber:      1,
				CurrentTemperature: 2.0,
				MinimumTemperature: -1.0,
				CurrentCapacity:    100,
				MaximumCapacity:    1000,
				MinimumCapacity:    10,
				WarehouseID:        1,
				ProductTypeID:      1,
			},
		},
		mockError: ErrInternal,
		mockGetError: nil,
	}
	service := NewService(&repository)
	ctx := new(context.Context)

	expected := ErrInternal

	// ACT
	result, err := service.Update(*ctx, domain.Section{})

	//ASSERT
	assert.Empty(t, result)
	assert.EqualError(t, expected, err.Error())
}

// TestDeleteNonExistent test if the service returns the correct error when a section with the given id doesn´t exists
func TestDeleteNonExistent(t *testing.T) {
	// ARRANGE
	repository := MockRepository{
		mockGetError: ErrNotFound,
	}
	service := NewService(&repository)
	ctx := new(context.Context)

	expected := ErrNotFound

	// ACT
	err := service.Delete(*ctx, 1)

	// ASSERT
	assert.EqualError(t, expected, err.Error())
}

// TestDeleteOk tests if the service correctly calls the repository to delete the section with the given id from storage
func TestDeleteOk(t *testing.T) {
	// ARRANGE
	repository := MockRepository{
		mockSections: []domain.Section{
			{
				ID:                 1,
				SectionNumber:      1,
				CurrentTemperature: 2.0,
				MinimumTemperature: -1.0,
				CurrentCapacity:    100,
				MaximumCapacity:    1000,
				MinimumCapacity:    10,
				WarehouseID:        1,
				ProductTypeID:      1,
			},
		},
	}
	service := NewService(&repository)
	ctx := new(context.Context)

	// ACT
	err1 := service.Delete(*ctx, 1)
	result, err2 := service.GetAll(*ctx)

	// ASSERT
	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Empty(t, result)
}

func TestGetAllSectionProductsOk(t *testing.T) {
	//	ARRANGE
	repository := MockRepository{
		mockProductsBySection: []domain.ProductsBySection{
			{
				SectionID:     1,
				SectionNumber: 1,
				ProductsCount: 100,
			},
			{
				SectionID:     2,
				SectionNumber: 2,
				ProductsCount: 90,
			},
		},
	}
	service := NewService(&repository)
	ctx := new(context.Context)

	expected := []domain.ProductsBySection{
		{
			SectionID:     1,
			SectionNumber: 1,
			ProductsCount: 100,
		},
		{
			SectionID:     2,
			SectionNumber: 2,
			ProductsCount: 90,
		},
	}

	//	ACT
	result, err := service.GetSectionProducts(*ctx, 0)

	//	ASSERT
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestGetSectionProductsOk(t *testing.T) {
	//	ARRANGE
	repository := MockRepository{
		mockProductsBySection: []domain.ProductsBySection{
			{
				SectionID:     1,
				SectionNumber: 1,
				ProductsCount: 100,
			},
			{
				SectionID:     2,
				SectionNumber: 2,
				ProductsCount: 90,
			},
		},
	}
	service := NewService(&repository)
	ctx := new(context.Context)

	expected := []domain.ProductsBySection{
		{
			SectionID:     1,
			SectionNumber: 1,
			ProductsCount: 100,
		},
	}
	//	ACT
	result, err := service.GetSectionProducts(*ctx, 1)

	//	ASSERT
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestGetSectionProductsIdNotExistent(t *testing.T) {
	// ARRANGE
	repository := MockRepository{mockError: ErrNotFound}
	service := NewService(&repository)
	ctx := new(context.Context)

	expected := ErrNotFound

	// ACT
	result, err := service.GetSectionProducts(*ctx, 3)

	// ASSERT
	assert.Empty(t, result)
	assert.EqualError(t, expected, err.Error())
}
