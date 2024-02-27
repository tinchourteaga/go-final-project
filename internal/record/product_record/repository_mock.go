package product_record

import (
	"context"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
)

type RepositoryMock struct {
	db            []domain.ProductRecord
	ForcedErrGet  error
	ForcedErrSave error
	FlagGet       bool
	FlagSave      bool
	ExpectedID    int
}

func (repository *RepositoryMock) Get(_ context.Context, _ int) (productRecord domain.ProductRecord, err error) {
	repository.FlagGet = true
	if repository.ForcedErrGet == nil {
		if len(repository.db) > 0 {
			productRecord = repository.db[0]
			return
		}
	}
	err = repository.ForcedErrGet
	return
}

func (repository *RepositoryMock) Save(_ context.Context, record domain.ProductRecord) (id int, err error) {
	repository.FlagSave = true
	if repository.ForcedErrSave == nil {
		id = repository.ExpectedID
		record.ID = id
		repository.db = append(repository.db, record)
		return
	}
	err = repository.ForcedErrSave
	return
}
