package seller

import (
	"context"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
)

type MockRepositorySeller struct {
	Seller        domain.Seller
	DataMock      []domain.Seller
	ErrorMock     error
	ErrorCidExist error
}

func (r *MockRepositorySeller) GetAll(ctx context.Context) (sellers []domain.Seller, err error) {
	if r.ErrorMock != nil {
		err = r.ErrorMock
		return
	}
	sellers = append(sellers, r.DataMock...)
	return
}

func (r *MockRepositorySeller) Get(ctx context.Context, id int) (s domain.Seller, err error) {
	if r.ErrorMock != nil {
		err = r.ErrorMock
		return
	}
	s = r.Seller
	return
}

func (r *MockRepositorySeller) Save(ctx context.Context, s domain.Seller) (id int, err error) {
	if r.ErrorMock != nil {
		err = r.ErrorMock
		return
	}
	id = 1
	return
}

func (r *MockRepositorySeller) Exists(ctx context.Context, cid int) (exist bool) {
	exist = r.ErrorCidExist != nil
	return
}

func (r *MockRepositorySeller) Update(ctx context.Context, s domain.Seller) (err error) {
	if r.ErrorMock != nil {
		err = r.ErrorMock
		return
	}
	return
}

func (r *MockRepositorySeller) Delete(ctx context.Context, id int) (err error) {
	if r.ErrorMock != nil {
		err = r.ErrorMock
		return
	}
	return
}
