package product_record

import (
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/gin-gonic/gin"
)

type ServiceMock struct {
	ProductRecordRepository []domain.ProductRecord
	ForcedErrGet            error
	ForcedErrSave           error
	FlagGet                 bool
	FlagSave                bool
	ExpectedID              int
}

func (service *ServiceMock) Get(_ *gin.Context, _ int) (productRecord domain.ProductRecord, err error) {
	service.FlagGet = true
	if service.ForcedErrGet == nil {
		if len(service.ProductRecordRepository) > 0 {
			productRecord = service.ProductRecordRepository[0]
			return
		}
		return
	}
	err = service.ForcedErrGet
	return
}

func (service *ServiceMock) Save(_ *gin.Context, productRecord domain.ProductRecord) (pr domain.ProductRecord, err error) {
	service.FlagSave = true
	if service.ForcedErrSave == nil {
		productRecord.ID = service.ExpectedID
		pr = productRecord
		return
	}
	err = service.ForcedErrSave
	return
}
