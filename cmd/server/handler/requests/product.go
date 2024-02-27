package requests

import "github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"

var (
	DefaultIntValue            = 0
	DefaultFloatValue  float32 = 0.0
	DefaultStringValue         = ""
)

// A ProductPOSTRequest
//   - uses pointers to allow 'zero' values on database
//   - gin.context.ShouldBindJSON() validates 'required' on specified fields
type ProductPOSTRequest struct {
	Description                    *string  `json:"description" binding:"required"`
	ExpirationRate                 *int     `json:"expiration_rate" binding:"required"`
	FreezingRate                   *int     `json:"freezing_rate" binding:"required"`
	Height                         *float32 `json:"height" binding:"required"`
	Length                         *float32 `json:"length" binding:"required"`
	NetWeight                      *float32 `json:"net_weight" binding:"required"`
	ProductCode                    *string  `json:"product_code" binding:"required"`
	RecommendedFreezingTemperature *float32 `json:"recommended_freezing_temperature" binding:"required"`
	Width                          *float32 `json:"width" binding:"required"`
	ProductTypeID                  *int     `json:"product_type_id" binding:"required"`
	SellerID                       *int     `json:"seller_id"`
}

// A ProductPATCHRequest
//   - uses pointers to allow 'zero' values on database -> This actually changed on Service implementation, we had to sacrifice 'zeros' for code simplicity
//   - uses pointers to allow partial updates
type ProductPATCHRequest struct {
	Description                    *string  `json:"description"`
	ExpirationRate                 *int     `json:"expiration_rate"`
	FreezingRate                   *int     `json:"freezing_rate"`
	Height                         *float32 `json:"height"`
	Length                         *float32 `json:"length"`
	NetWeight                      *float32 `json:"net_weight"`
	ProductCode                    *string  `json:"product_code"`
	RecommendedFreezingTemperature *float32 `json:"recommended_freezing_temperature"`
	Width                          *float32 `json:"width"`
	ProductTypeID                  *int     `json:"product_type_id"`
	SellerID                       *int     `json:"seller_id"`
}

func (request *ProductPOSTRequest) MapToDomain() domain.Product {
	return domain.Product{
		Description:                    *request.Description,
		ExpirationRate:                 *request.ExpirationRate,
		FreezingRate:                   *request.FreezingRate,
		Height:                         *request.Height,
		Length:                         *request.Length,
		NetWeight:                      *request.NetWeight,
		ProductCode:                    *request.ProductCode,
		RecommendedFreezingTemperature: *request.RecommendedFreezingTemperature,
		Width:                          *request.Width,
		ProductTypeID:                  *request.ProductTypeID,
		SellerID:                       request.SellerID,
	}
}

func (request *ProductPATCHRequest) MapToDomain() domain.Product {
	// If any attribute is nil, it will be replaced by its 'zero' value.
	// domain.Product has to have a value, database doesn't allow null values (Except sellerID).
	if request.Description == nil {
		request.Description = &DefaultStringValue
	}
	if request.ExpirationRate == nil {
		request.ExpirationRate = &DefaultIntValue
	}
	if request.FreezingRate == nil {
		request.FreezingRate = &DefaultIntValue
	}
	if request.Height == nil {
		request.Height = &DefaultFloatValue
	}
	if request.Length == nil {
		request.Length = &DefaultFloatValue
	}
	if request.NetWeight == nil {
		request.NetWeight = &DefaultFloatValue
	}
	if request.RecommendedFreezingTemperature == nil {
		request.RecommendedFreezingTemperature = &DefaultFloatValue
	}
	if request.Width == nil {
		request.Width = &DefaultFloatValue
	}
	if request.ProductTypeID == nil {
		request.ProductTypeID = &DefaultIntValue
	}
	if request.ProductCode == nil {
		request.ProductCode = &DefaultStringValue
	}
	return domain.Product{
		Description:                    *request.Description,
		ExpirationRate:                 *request.ExpirationRate,
		FreezingRate:                   *request.FreezingRate,
		Height:                         *request.Height,
		Length:                         *request.Length,
		NetWeight:                      *request.NetWeight,
		ProductCode:                    *request.ProductCode,
		RecommendedFreezingTemperature: *request.RecommendedFreezingTemperature,
		Width:                          *request.Width,
		ProductTypeID:                  *request.ProductTypeID,
		SellerID:                       request.SellerID,
	}
}
