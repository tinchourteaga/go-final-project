package product

import (
	"context"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
)

type RepositoryMock struct {
	db              []domain.Product
	ForcedErrGetAll error
	ForcedErrGet    error
	ForcedErrExists error
	ForcedErrSave   error
	ForcedErrUpdate error
	ForcedErrDelete error
	FlagGetAll      bool
	FlagGet         bool
	FlagExists      bool
	FlagSave        bool
	FlagUpdate      bool
	FlagDelete      bool
	ExpectedID      int
}

// GetAll returns only weird SQL errors
func (repository *RepositoryMock) GetAll(_ context.Context) (products []domain.Product, err error) {
	repository.FlagGetAll = true
	if repository.ForcedErrGetAll == nil {
		products = repository.db
	}
	err = repository.ForcedErrGetAll
	return
}

// Get returns ErrNotFound and ErrInternal
func (repository *RepositoryMock) Get(_ context.Context, _ int) (product domain.Product, err error) {
	repository.FlagGet = true
	if repository.ForcedErrGet == nil {
		if len(repository.db) > 0 {
			product = repository.db[0]
			return
		}
		return
	}
	err = repository.ForcedErrGet
	return
}

// Exists returns only boolean values
func (repository *RepositoryMock) Exists(_ context.Context, _ string) bool {
	repository.FlagExists = true
	return repository.ForcedErrExists != nil
}

// Save returns only weird SQL errors
func (repository *RepositoryMock) Save(_ context.Context, product domain.Product) (id int, err error) {
	repository.FlagSave = true
	if repository.ForcedErrSave == nil {
		id = repository.ExpectedID
		product.ID = id
		repository.db = append(repository.db, product)
		return
	}
	err = repository.ForcedErrSave
	return
}

// Update returns only weird SQL errors
func (repository *RepositoryMock) Update(_ context.Context, p domain.Product) error {
	repository.FlagUpdate = true
	if repository.ForcedErrUpdate == nil {
		if len(repository.db) > 0 {
			repository.db[0] = p
		}
	}
	return repository.ForcedErrUpdate
}

// Delete returns ErrNotFound and weird SQL errors
func (repository *RepositoryMock) Delete(_ context.Context, _ int) error {
	repository.FlagDelete = true
	return repository.ForcedErrDelete
}
