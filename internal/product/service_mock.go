package product

import (
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/gin-gonic/gin"
)

type ServiceMock struct {
	ProductRepository      []domain.Product
	ForcedErrGet           error
	ForcedErrGetAll        error
	ForcedErrSave          error
	ForcedErrPartialUpdate error
	ForcedErrDelete        error
	FlagGetAll             bool
	FlagGet                bool
	FlagSave               bool
	FlagPartialUpdate      bool
	FlagDelete             bool
	ExpectedID             int
}

// GetAll returns only weird SQL errors
func (service *ServiceMock) GetAll(_ *gin.Context) (products []domain.Product, err error) {
	service.FlagGetAll = true
	if service.ForcedErrGetAll == nil {
		if len(service.ProductRepository) == 0 {
			products = nil
		} else {
			products = service.ProductRepository
		}
	}
	err = service.ForcedErrGetAll
	return
}

// Get returns only weird SQL errors
func (service *ServiceMock) Get(_ *gin.Context, _ int) (product domain.Product, err error) {
	service.FlagGet = true
	if service.ForcedErrGet == nil {
		if len(service.ProductRepository) > 0 {
			product = service.ProductRepository[0]
			return
		}
		return
	}
	err = service.ForcedErrGet
	return
}

// Save returns ErrAlreadyExists, ErrNotFound and weird SQL errors
func (service *ServiceMock) Save(_ *gin.Context, product domain.Product) (p domain.Product, err error) {
	service.FlagSave = true
	if service.ForcedErrSave == nil {
		product.ID = service.ExpectedID
		p = product
		return
	}
	err = service.ForcedErrSave
	return
}

// PartialUpdate returns ErrAlreadyExists, ErrNotFound and weird SQL errors
func (service *ServiceMock) PartialUpdate(_ *gin.Context, id int, product domain.Product) (p domain.Product, err error) {
	service.FlagPartialUpdate = true
	if service.ForcedErrPartialUpdate == nil {
		product.ID = id
		if product.Description != "" {
			service.ProductRepository[0].Description = product.Description
		}
		if product.ExpirationRate != 0 {
			service.ProductRepository[0].ExpirationRate = product.ExpirationRate
		}
		if product.FreezingRate != 0 {
			service.ProductRepository[0].FreezingRate = product.FreezingRate
		}
		if product.Height != 0.0 {
			service.ProductRepository[0].Height = product.Height
		}
		if product.Length != 0.0 {
			service.ProductRepository[0].Length = product.Length
		}
		if product.NetWeight != 0.0 {
			service.ProductRepository[0].NetWeight = product.NetWeight
		}
		if product.RecommendedFreezingTemperature != 0.0 {
			service.ProductRepository[0].RecommendedFreezingTemperature = product.RecommendedFreezingTemperature
		}
		if product.Width != 0.0 {
			service.ProductRepository[0].Width = product.Width
		}
		if product.ProductTypeID != 0 {
			service.ProductRepository[0].ProductTypeID = product.ProductTypeID
		}
		if product.SellerID != nil {
			service.ProductRepository[0].SellerID = product.SellerID
		}
		if product.ProductCode != "" {
			service.ProductRepository[0].ProductCode = product.ProductCode
		}
		p = service.ProductRepository[0]
		return
	}
	err = service.ForcedErrPartialUpdate
	return
}

// Delete returns ErrNotFound and weird SQL errors
func (service *ServiceMock) Delete(_ *gin.Context, _ int) error {
	service.FlagDelete = true
	return service.ForcedErrDelete
}
