package buyer

import (
	"context"
	"errors"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/stretchr/testify/assert"
)

func init() {
	logging.InitLog(nil)
}

// list of buyers for unit test
var ListBuyers = []domain.Buyer{
	{
		ID:           1,
		CardNumberID: "001",
		FirstName:    "Comprador 1",
		LastName:     "Vendedor 1",
	}, {
		ID:           2,
		CardNumberID: "002",
		FirstName:    "Comprador 2",
		LastName:     "Vendedor 2",
	}, {
		ID:           3,
		CardNumberID: "003",
		FirstName:    "Comprador 3",
		LastName:     "Vendedor 3",
	}, {
		ID:           4,
		CardNumberID: "004",
		FirstName:    "Comprador 4",
		LastName:     "Vendedor 4",
	}, {
		ID:           5,
		CardNumberID: "005",
		FirstName:    "Comprador 5",
		LastName:     "Vendedor 5",
	},
}

// TestGetAllSuccess passes when there are no errors on repository layer
// (return slice of all domain.Buyers and nil error)
func TestGetAllSuccess(t *testing.T) {
	//arrange
	Expectedresult := ListBuyers
	//Act
	MockRepo := MockRepository{
		Data: Expectedresult,
	}
	serv := NewService(&MockRepo)
	ctx := new(context.Context)
	result, err := serv.GetAll(*ctx)

	//arrange
	assert.Nil(t, err)
	assert.Equal(t, Expectedresult, result)
}

// TestGetAllFail passes when there are error on repository
// (return "buyer not found", http.StatusNotFound)
func TestGetAllFail(t *testing.T) {
	//arrange
	errExpected := errors.New("buyer not found")
	//Act
	MockRepo := MockRepository{
		Data: nil,
	}
	serv := NewService(&MockRepo)
	ctx := new(context.Context)
	result, err := serv.GetAll(*ctx)

	//arrange
	assert.Nil(t, result)
	assert.EqualError(t, errExpected, err.Error())
}

// TestGetSuccess passes when there are no errors on repository layer
// (return domain.Buyer and nil error)
func TestGetSuccess(t *testing.T) {
	//arrange
	id := 4
	Expectedresult := domain.Buyer{
		ID:           4,
		CardNumberID: "004",
		FirstName:    "Comprador 4",
		LastName:     "Vendedor 4",
	}

	//Act
	MockRepo := MockRepository{
		Data: ListBuyers,
	}
	serv := NewService(&MockRepo)
	ctx := new(context.Context)
	result, err := serv.Get(*ctx, id)

	//arrange
	assert.Nil(t, err)
	assert.Equal(t, Expectedresult, result)
}

// TestGetAllFail passes when there are error on repository
// (return "buyer not found", http.StatusNotFound)
func TestGetFailNullDataBase(t *testing.T) {
	//arrange
	id := 4
	errExpected := errors.New("buyer not found")
	ExpectedResult := domain.Buyer{}
	//Act
	MockRepo := MockRepository{
		Data: []domain.Buyer{},
	}
	serv := NewService(&MockRepo)
	ctx := new(context.Context)
	result, err := serv.Get(*ctx, id)

	//arrange
	assert.Equal(t, ExpectedResult, result)
	assert.EqualError(t, errExpected, err.Error())
}

// TestGetFailIncorrectId passes when the given id is not in database
// (return domain.buyer{} and error "buyer not found", http.StatusNotFound)
func TestGetFailIncorrectId(t *testing.T) {
	//arrange
	id := 10
	errExpected := errors.New("buyer not found")
	ExpectedResult := domain.Buyer{}

	//Act
	MockRepo := MockRepository{
		Data: ListBuyers,
	}
	serv := NewService(&MockRepo)
	ctx := new(context.Context)
	result, err := serv.Get(*ctx, id)

	//arrange
	assert.Equal(t, ExpectedResult, result)
	assert.EqualError(t, errExpected, err.Error())
}

// TestExistsSuccess passes when the given id is in database
// (return true)
func TestExistsSuccess(t *testing.T) {
	//arrange
	id := "005"
	ExpectedResult := true

	//Act
	MockRepo := MockRepository{
		Data: ListBuyers,
	}
	serv := NewService(&MockRepo)
	ctx := new(context.Context)
	result := serv.Exists(*ctx, id)

	//arrange
	assert.Equal(t, ExpectedResult, result)
}

// TestExistsFail passes when the given id is not in database
// (return false)
func TestExistsFail(t *testing.T) {
	//arrange
	id := "288"
	ExpectedResult := false

	//Act
	MockRepo := MockRepository{
		Data: ListBuyers,
	}
	serv := NewService(&MockRepo)
	ctx := new(context.Context)
	result := serv.Exists(*ctx, id)

	//arrange
	assert.Equal(t, ExpectedResult, result)
}

// TestSaveSuccess passes when data is correct (return domain.Buyer{} and nil error)
func TestSaveSuccess(t *testing.T) {
	//arrange
	newBuyer := domain.Buyer{
		ID:           6,
		CardNumberID: "006",
		FirstName:    "Comprador 6",
		LastName:     "Vendedor 6",
	}

	//Act
	MockRepo := MockRepository{
		Data: ListBuyers,
	}
	serv := NewService(&MockRepo)
	ctx := new(context.Context)
	result, err := serv.Save(*ctx, newBuyer)

	//arrange
	assert.Nil(t, err)
	assert.Equal(t, newBuyer, result)
}

// TestSaveFailCardNumberIdExist passes when domain.Buyer.ID already exists
// (return domain.Product{} and error buyer.ErrAlreadyExists)
func TestSaveFailCardNumberIdExist(t *testing.T) {
	//arrange
	newBuyer := domain.Buyer{
		ID:           6,
		CardNumberID: "004",
		FirstName:    "Comprador 6",
		LastName:     "Vendedor 6",
	}
	expectedResult, expectedErr := domain.Buyer{}, errors.New("card_number_id already exists")

	//Act
	MockRepo := MockRepository{
		Data: ListBuyers,
	}
	serv := NewService(&MockRepo)
	ctx := new(context.Context)
	result, err := serv.Save(*ctx, newBuyer)

	//arrange
	assert.Equal(t, expectedResult, result)
	assert.EqualError(t, expectedErr, err.Error())

}

// TestSaveFailInternalError passes when serv.Save() return error
// (return domain.Product{} and error buyer.ErrInternal)
func TestSaveFailInternalError(t *testing.T) {
	//arrange
	newBuyer := domain.Buyer{
		ID:           6,
		CardNumberID: "006",
		FirstName:    "Comprador 6",
		LastName:     "Vendedor 6",
	}
	expectedErr := errors.New("database internal error")

	//Act
	MockRepo := MockRepository{
		Data: ListBuyers,
		Err:  errors.New("database internal error"),
	}
	serv := NewService(&MockRepo)
	ctx := new(context.Context)
	result, err := serv.Save(*ctx, newBuyer)

	//arrange
	assert.Empty(t, result)
	assert.EqualError(t, expectedErr, err.Error())
}

// TestDeleteSuccess passes when id exists and deletion is successful
// (return nil error)
func TestDeleteSuccess(t *testing.T) {
	//arrange
	id := 4
	//Act
	MockRepo := MockRepository{
		Data: ListBuyers,
	}
	serv := NewService(&MockRepo)
	ctx := new(context.Context)
	err := serv.Delete(*ctx, id)

	//arrange
	assert.Nil(t, err)
}

// TestDeleteFailIdNotExist passes when the given id is not in database
// (return error buyer.ErrNotFound)
func TestDeleteFailIdNotExist(t *testing.T) {
	//arrange
	id := 15
	expectedError := errors.New("buyer not found")
	//Act
	MockRepo := MockRepository{
		Data: ListBuyers,
	}
	serv := NewService(&MockRepo)
	ctx := new(context.Context)
	err := serv.Delete(*ctx, id)

	//arrange
	assert.Equal(t, expectedError, err)
}

// TestUpdateSuccess passes when data is correct
// (return updated domain.Buyer and error nil)
func TestUpdateSuccess(t *testing.T) {
	//arrange
	buyer := domain.Buyer{
		ID:           4,
		CardNumberID: "005",
		FirstName:    "Comprador 6",
		LastName:     "Vendedor 6",
	}

	//Act
	MockRepo := MockRepository{
		Data: ListBuyers,
	}
	serv := NewService(&MockRepo)
	ctx := new(context.Context)
	result, err := serv.Update(*ctx, buyer)

	//arrange
	assert.Nil(t, err)
	assert.Equal(t, buyer, result)
}

// TestUpdateFailWrongId passes when id is not correct
// (return updated domain.Buyer{} and error Buyer.ErrNotFound)
func TestUpdateFailWrongId(t *testing.T) {
	//arrange
	buyer := domain.Buyer{
		ID:           15,
		CardNumberID: "004",
		FirstName:    "Comprador 6",
		LastName:     "Vendedor 6",
	}
	errExpected := errors.New("buyer not found")
	//Act
	MockRepo := MockRepository{
		Data: ListBuyers,
	}
	serv := NewService(&MockRepo)
	ctx := new(context.Context)
	result, err := serv.Update(*ctx, buyer)

	//arrange
	assert.Empty(t, result)
	assert.EqualError(t, errExpected, err.Error())
}

// TestUpdateFailInternalError passes when Update fails
// (return domain.Buyer{} and error Buyer.ErrInternal)
func TestUpdateFailInternalError(t *testing.T) {
	//arrange
	buyer := domain.Buyer{
		ID:           4,
		CardNumberID: "004",
		FirstName:    "Comprador 6",
		LastName:     "Vendedor 6",
	}
	errExpected := ErrInternal
	//Act
	MockRepo := MockRepository{
		Data: ListBuyers,
		Err:  ErrInternal,
	}
	serv := NewService(&MockRepo)
	ctx := new(context.Context)
	result, err := serv.Update(*ctx, buyer)

	//arrange
	assert.Empty(t, result)
	assert.EqualError(t, errExpected, err.Error())
}
