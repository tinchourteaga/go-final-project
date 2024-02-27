package report_record

import (
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/gin-gonic/gin"
)

type ServiceMock struct {
	ReportRecordRepository []domain.ReportRecord
	ForcedErrGet           error
	FlagGet                bool
}

func (service *ServiceMock) Get(_ *gin.Context, id *int) (reports []domain.ReportRecord, err error) {
	service.FlagGet = true
	if service.ForcedErrGet == nil {
		if id == nil {
			if len(service.ReportRecordRepository) == 0 {
				reports = nil
			} else {
				reports = service.ReportRecordRepository
			}
		} else {
			if len(service.ReportRecordRepository) > 0 {
				reports = []domain.ReportRecord{service.ReportRecordRepository[0]}
			}
		}
	}
	err = service.ForcedErrGet
	return
}
