package product

import (
	"errors"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/gin-gonic/gin"
)

var (
	ServiceErrNotFound           = errors.New("product not found")
	ServiceErrInternal           = errors.New("internal error")
	ServiceErrAlreadyExists      = errors.New("product code already exists")
	ServiceErrForeignKeyNotFound = errors.New("seller not found")
)

type Service interface {
	GetAll(ctx *gin.Context) ([]domain.Product, error)
	Get(ctx *gin.Context, id int) (domain.Product, error)
	Save(ctx *gin.Context, product domain.Product) (domain.Product, error)
	PartialUpdate(ctx *gin.Context, id int, product domain.Product) (domain.Product, error)
	Delete(ctx *gin.Context, id int) error
}

type service struct {
	productRepository Repository
}

func NewService(repository Repository) Service {
	return &service{
		productRepository: repository,
	}
}

// GetAll returns a list with all the Products from database or nil if it's empty.
// If there is any error it is returned to the controller layer to be handled.
func (s *service) GetAll(ctx *gin.Context) ([]domain.Product, error) {
	products, errGetAll := s.productRepository.GetAll(ctx)
	if errGetAll != nil {
		logging.Log(errGetAll)
		return []domain.Product{}, ServiceErrInternal
	}
	return products, nil
}

// Get returns a Product from database or error if not found.
// If there is any error it is returned to the controller layer to be handled.
func (s *service) Get(ctx *gin.Context, id int) (domain.Product, error) {
	product, errGet := s.productRepository.Get(ctx, id)
	if errGet != nil {
		logging.Log(errGet)
		switch errGet {
		case RepositoryErrNotFound:
			return domain.Product{}, ServiceErrNotFound
		default:
			return domain.Product{}, ServiceErrInternal
		}
	}
	return product, nil
}

// Save stores the given values in a new Product in te database.
// sellerID is optional and productCode should be unique.
// After storing, Save retrieves the new Product from the database and returns it.
// If there is any error it is returned to the controller layer to be handled.
func (s *service) Save(ctx *gin.Context, product domain.Product) (domain.Product, error) {
	if s.productRepository.Exists(ctx, product.ProductCode) {
		return domain.Product{}, ServiceErrAlreadyExists
	}
	prodID, errSave := s.productRepository.Save(ctx, product)
	if errSave != nil {
		logging.Log(errSave)
		switch errSave {
		case RepositoryErrForeignKeyConstraint:
			return domain.Product{}, ServiceErrForeignKeyNotFound
		case RepositoryErrAlreadyExists:
			// This is in case we implement unique with product_code (not happening on Sprint III)
			return domain.Product{}, ServiceErrAlreadyExists
		default:
			return domain.Product{}, ServiceErrInternal
		}
	}
	return s.Get(ctx, prodID)
}

// PartialUpdate retrieves a Product from database and checks for values != nil to update.
// productCode should be unique.
// After updating, PartialUpdate retrieves the updated Product from the database and returns it.
// If there is any error it is returned to the controller layer to be handled.
// This method needs this many attributes because there is no other way of knowing the number of attributes to update that accepts 'zero' values
func (s *service) PartialUpdate(ctx *gin.Context, id int, product domain.Product) (domain.Product, error) {
	productOriginal, errGetOriginal := s.Get(ctx, id)
	if errGetOriginal != nil {
		return domain.Product{}, errGetOriginal
	}
	// This is the only way to do this with Vanilla Go.
	// We're losing the ability to use 'zero' values but at least we can do partial updates
	// Comparison 'product.x != productOriginal.x' is useless, if they are the same, and we change it... it stays the same!
	if product.Description != "" {
		productOriginal.Description = product.Description
	}
	if product.ExpirationRate != 0 {
		productOriginal.ExpirationRate = product.ExpirationRate
	}
	if product.FreezingRate != 0 {
		productOriginal.FreezingRate = product.FreezingRate
	}
	if product.Height != 0.0 {
		productOriginal.Height = product.Height
	}
	if product.Length != 0.0 {
		productOriginal.Length = product.Length
	}
	if product.NetWeight != 0.0 {
		productOriginal.NetWeight = product.NetWeight
	}
	if product.RecommendedFreezingTemperature != 0.0 {
		productOriginal.RecommendedFreezingTemperature = product.RecommendedFreezingTemperature
	}
	if product.Width != 0.0 {
		productOriginal.Width = product.Width
	}
	if product.ProductTypeID != 0 {
		productOriginal.ProductTypeID = product.ProductTypeID
	}
	if product.SellerID != nil {
		productOriginal.SellerID = product.SellerID
	}
	// checks product_code to be updated is not the same as the one in the struct:
	//     - if it is, does nothing
	//     - if it's not, updates product_code and checks if it's unique on the database
	if product.ProductCode != "" && product.ProductCode != productOriginal.ProductCode {
		productOriginal.ProductCode = product.ProductCode
		if s.productRepository.Exists(ctx, productOriginal.ProductCode) {
			return domain.Product{}, ServiceErrAlreadyExists
		}
	}
	errUpdate := s.productRepository.Update(ctx, productOriginal)
	if errUpdate != nil {
		logging.Log(errUpdate)
		switch errUpdate {
		case RepositoryErrForeignKeyConstraint:
			return domain.Product{}, ServiceErrForeignKeyNotFound
		case RepositoryErrAlreadyExists:
			// This is in case we implement unique with product_code (not happening on Sprint III)
			return domain.Product{}, ServiceErrAlreadyExists
		default:
			return domain.Product{}, ServiceErrInternal
		}
	}
	return s.Get(ctx, id)
}

// Delete removes a Product from database.
// If there is any error it is returned to the controller layer to be handled.
func (s *service) Delete(ctx *gin.Context, id int) error {
	errDelete := s.productRepository.Delete(ctx, id)
	if errDelete != nil {
		logging.Log(errDelete)
		switch errDelete {
		case RepositoryErrNotFound:
			return ServiceErrNotFound
		default:
			return ServiceErrInternal
		}
	}
	return nil
}
