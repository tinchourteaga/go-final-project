package report_record

import (
	"context"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
)

type RepositoryMock struct {
	db              []domain.ReportRecord
	ForcedErrGetAll error
	ForcedErrGet    error
	FlagGetAll      bool
	FlagGet         bool
}

func (repository *RepositoryMock) GetAll(_ context.Context) (reports []domain.ReportRecord, err error) {
	repository.FlagGetAll = true
	if repository.ForcedErrGetAll == nil {
		reports = repository.db
	}
	err = repository.ForcedErrGetAll
	return
}

func (repository *RepositoryMock) Get(_ context.Context, _ int) (report domain.ReportRecord, err error) {
	repository.FlagGet = true
	if repository.ForcedErrGet == nil {
		if len(repository.db) > 0 {
			report = repository.db[0]
			return
		}
		return
	}
	err = repository.ForcedErrGet
	return
}
