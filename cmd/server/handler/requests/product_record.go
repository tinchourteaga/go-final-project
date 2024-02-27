package requests

import (
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"time"
)

// A ProductRecordPOSTRequest
//   - uses pointers to allow 'zero' values on database
//   - gin.context.ShouldBindJSON() validates 'required' on specified fields
type ProductRecordPOSTRequest struct {
	LastUpdateDate string   `json:"last_update_date" binding:"required"`
	PurchasePrice  *float32 `json:"purchase_price" binding:"required"`
	SalePrice      *float32 `json:"sale_price" binding:"required"`
	ProductID      *int     `json:"product_id" binding:"required"`
}

func (request *ProductRecordPOSTRequest) MapToDomain() (domain.ProductRecord, error) {
	date, errConv := StringToMySQLDate(request.LastUpdateDate)
	if errConv != nil {
		return domain.ProductRecord{}, errConv
	}
	productRecord := domain.ProductRecord{
		LastUpdateDate: domain.MySqlTime{Time: date},
		PurchasePrice:  *request.PurchasePrice,
		SalePrice:      *request.SalePrice,
		ProductID:      *request.ProductID,
	}
	return productRecord, nil
}

func StringToMySQLDate(date string) (time.Time, error) {
	return time.Parse(domain.ISO8601, date)
}
