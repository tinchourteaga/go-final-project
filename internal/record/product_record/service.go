package product_record

import (
	"errors"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/gin-gonic/gin"
	"time"
)

var (
	ServiceErrNotFound           = errors.New("product not found")
	ServiceErrInternal           = errors.New("internal error")
	ServiceErrForeignKeyNotFound = errors.New("product not found")
	ServiceErrDate               = errors.New("invalid date")
)

type Service interface {
	Get(ctx *gin.Context, id int) (domain.ProductRecord, error)
	Save(ctx *gin.Context, record domain.ProductRecord) (domain.ProductRecord, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (service *service) Get(ctx *gin.Context, id int) (domain.ProductRecord, error) {
	productRecord, errGet := service.repository.Get(ctx, id)
	if errGet != nil {
		logging.Log(errGet)
		switch errGet {
		case RepositoryErrNotFound:
			return domain.ProductRecord{}, ServiceErrNotFound
		default:
			return domain.ProductRecord{}, ServiceErrInternal
		}
	}
	return productRecord, nil
}

func (service *service) Save(ctx *gin.Context, record domain.ProductRecord) (domain.ProductRecord, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	if record.LastUpdateDate.Before(today) {
		logging.Log(ServiceErrDate)
		return domain.ProductRecord{}, ServiceErrDate
	}
	id, errSave := service.repository.Save(ctx, record)
	if errSave != nil {
		logging.Log(errSave)
		switch errSave {
		case RepositoryErrForeignKeyConstraint:
			return domain.ProductRecord{}, ServiceErrForeignKeyNotFound
		default:
			return domain.ProductRecord{}, ServiceErrInternal
		}
	}
	return service.Get(ctx, id)
}
