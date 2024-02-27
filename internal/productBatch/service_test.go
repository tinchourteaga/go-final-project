package productbatch

import (
	"context"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreateOk(t *testing.T) {
	// ARRANGE
	repository := MockRepository{mockProductBatches: []domain.ProductBatch{}}
	service := NewService(&repository)
	ctx := new(context.Context)

	expected := domain.ProductBatch{
		ID:                 1,
		BatchNumber:        1,
		CurrentQuantity:    1,
		CurrentTemperature: 1,
		DueDate:            "1999-12-12",
		InitialQuantity:    1,
		ManufacturingDate:  "1999-12-12",
		ManufacturingHour:  1,
		MinimumTemperature: 1,
		ProductID:          1,
		SectionID:          1,
	}

	// ACT
	result, err := service.Create(*ctx, domain.ProductBatch{
		BatchNumber:        1,
		CurrentQuantity:    1,
		CurrentTemperature: 1,
		DueDate:            "1999-12-12",
		InitialQuantity:    1,
		ManufacturingDate:  "1999-12-12",
		ManufacturingHour:  1,
		MinimumTemperature: 1,
		ProductID:          1,
		SectionID:          1,
	})

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestCreateAlreadyExistsConflict(t *testing.T) {
	// ARRANGE
	repository := MockRepository{mockError: ErrAlreadyExists}
	service := NewService(&repository)
	ctx := new(context.Context)

	expected := ErrAlreadyExists

	// ACT
	result, err := service.Create(*ctx, domain.ProductBatch{})

	// ASSERT
	assert.Empty(t, result)
	assert.EqualError(t, expected, err.Error())
}

func TestCreateMissingSectionConflict(t *testing.T) {
	// ARRANGE
	repository := MockRepository{mockError: ErrForeignSectionNotFound}
	service := NewService(&repository)
	ctx := new(context.Context)

	expected := ErrForeignSectionNotFound

	// ACT
	result, err := service.Create(*ctx, domain.ProductBatch{})

	// ASSERT
	assert.Empty(t, result)
	assert.EqualError(t, expected, err.Error())
}

func TestCreateMissingProductConflict(t *testing.T) {
	// ARRANGE
	repository := MockRepository{mockError: ErrForeignProductNotFound}
	service := NewService(&repository)
	ctx := new(context.Context)

	expected := ErrForeignProductNotFound

	// ACT
	result, err := service.Create(*ctx, domain.ProductBatch{})

	// ASSERT
	assert.Empty(t, result)
	assert.EqualError(t, expected, err.Error())
}
