package report_record

import (
	"errors"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/gin-gonic/gin"
)

var (
	ServiceErrNotFound = errors.New("product record not found")
	ServiceErrInternal = errors.New("database internal error")
)

type Service interface {
	Get(ctx *gin.Context, productID *int) ([]domain.ReportRecord, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (service *service) Get(ctx *gin.Context, productID *int) ([]domain.ReportRecord, error) {
	var reports []domain.ReportRecord
	if productID == nil {
		var errGetAll error
		reports, errGetAll = service.repository.GetAll(ctx)
		if errGetAll != nil {
			logging.Log(errGetAll)
			return []domain.ReportRecord{}, ServiceErrInternal
		}
		if reports == nil {
			return []domain.ReportRecord{}, nil
		}
	} else {
		report, errGet := service.repository.Get(ctx, *productID)
		if errGet != nil {
			logging.Log(errGet)
			switch errGet {
			case RepositoryErrNotFound:
				return []domain.ReportRecord{}, ServiceErrNotFound
			case RepositoryErrInternal:
				return []domain.ReportRecord{}, ServiceErrInternal
			}
		}
		reports = append(reports, report)
	}
	return reports, nil
}
