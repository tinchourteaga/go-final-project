package seller

import (
	"context"
	"errors"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/stretchr/testify/assert"
)

var ctx context.Context

func init() {
	logging.InitLog(nil)
}

// * ---------------------- GetAll -------------------------
// TestGetAll_Seller passes when return correct sellers
func TestGetAll_Seller(t *testing.T) {
	// Arrange
	database := []domain.Seller{{
		ID:          1,
		CID:         1,
		CompanyName: "Kiosco 1",
		Address:     "Junin 323",
		Telephone:   "2664727336",
	}, {
		ID:          2,
		CID:         2,
		CompanyName: "Kiosco 2",
		Address:     "Colon 664",
		Telephone:   "2664233242",
	},
	}

	mockRepo := MockRepositorySeller{
		DataMock:  database,
		ErrorMock: nil,
	}
	service := NewService(&mockRepo)

	// Act
	result, err := service.GetAll(ctx)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, database, result)
}

// TestGetAllFail_Sellers passes when return an error in GetAll from db
func TestGetAllFail_Seller(t *testing.T) {
	// Arrange
	expectedError := errors.New("Error on GetAll")

	mockRepo := MockRepositorySeller{
		ErrorMock: expectedError,
	}
	service := NewService(&mockRepo)

	// Act
	_, err := service.GetAll(ctx)

	// Assert
	assert.EqualError(t, expectedError, err.Error())
}

// * ---------------------- Get ----------------------------
// TestGet_Seller passes when return correct seller
func TestGet_Seller(t *testing.T) {
	// Arrange
	database := []domain.Seller{{
		ID:          1,
		CID:         1,
		CompanyName: "Kiosco 1",
		Address:     "Junin 323",
		Telephone:   "2664727336",
	}, {
		ID:          2,
		CID:         2,
		CompanyName: "Kiosco 2",
		Address:     "Colon 664",
		Telephone:   "2664233242",
	},
	}

	expectedSeller := domain.Seller{
		ID:          1,
		CID:         1,
		CompanyName: "Kiosco 1",
		Address:     "Junin 323",
		Telephone:   "2664727336",
	}

	mockRepo := MockRepositorySeller{
		Seller:    expectedSeller,
		DataMock:  database,
		ErrorMock: nil,
	}
	service := NewService(&mockRepo)

	// Act
	result, err := service.Get(ctx, expectedSeller.ID)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expectedSeller, result)
}

// TestSGet_IDNonExistent_Seller passes when the given id is not in database (return domain.Seller{} and error seller.ErrNotFound)
func TestSGet_IDNonExistent_Seller(t *testing.T) {
	// Arrange
	expectedError := ServiceErrNotFound
	mockRepo := MockRepositorySeller{
		ErrorMock: expectedError,
	}
	service := NewService(&mockRepo)

	// Act
	result, err := service.Get(ctx, 1)

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.Equal(t, domain.Seller{}, result)
}

// TestGetFail_Seller passes when return error invalid id
func TestGetFail_Seller(t *testing.T) {
	// Arrange
	expectedError := errors.New("Id not exist")

	mockRepo := MockRepositorySeller{
		ErrorMock: expectedError,
	}
	service := NewService(&mockRepo)

	// Act
	_, err := service.Get(ctx, 1)

	// Assert
	assert.EqualError(t, expectedError, err.Error())
}

// * ---------------------- Create --------------------------
// TestCreate_Seller passes when return seller created
func TestCreate_Seller(t *testing.T) {
	// Arrange
	sellerToCreate := domain.Seller{
		CID:         10,
		CompanyName: "Kiosco 10",
		Address:     "Mitre 1323",
		Telephone:   "26647273999",
	}

	sellerExpected := domain.Seller{
		ID:          1,
		CID:         10,
		CompanyName: "Kiosco 10",
		Address:     "Mitre 1323",
		Telephone:   "26647273999",
	}

	mockRepo := MockRepositorySeller{
		Seller:    sellerExpected,
		ErrorMock: nil,
	}
	service := NewService(&mockRepo)

	// Act
	result, err := service.Create(ctx, sellerToCreate)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, sellerExpected, result)
}

// TestCreateFail_Seller passes when return an error for failed Create function
func TestCreateFail_Seller(t *testing.T) {
	// Arrange
	expectedError := errors.New("Cid already exists")
	sellerToCreate := domain.Seller{}

	mockRepo := MockRepositorySeller{
		ErrorMock: expectedError,
	}
	service := NewService(&mockRepo)

	// Act
	result, err := service.Create(ctx, sellerToCreate)

	// Assert
	assert.EqualError(t, expectedError, err.Error())
	assert.Equal(t, domain.Seller{}, result)
}

// TestCreateExistFail_Seller passes when return an error for failed Exist function
func TestCreateExistFail_Seller(t *testing.T) {
	// Arrange
	expectedError := errors.New("cid already exists")
	sellerToCreate := domain.Seller{}

	mockRepo := MockRepositorySeller{
		ErrorCidExist: expectedError,
	}
	service := NewService(&mockRepo)

	// Act
	result, err := service.Create(ctx, sellerToCreate)

	// Assert
	assert.EqualError(t, expectedError, err.Error())
	assert.Equal(t, domain.Seller{}, result)
}

// * ---------------------- Delete ---------------------------
// TestDelete_Seller passes when return Delete function is successful
func TestDelete_Seller(t *testing.T) {
	// Arrange
	database := []domain.Seller{{
		ID:          1,
		CID:         1,
		CompanyName: "Kiosco 1",
		Address:     "Junin 323",
		Telephone:   "2664727336",
	},
	}

	mockRepo := MockRepositorySeller{
		DataMock: database,
	}

	service := NewService(&mockRepo)

	// Act
	err := service.Delete(ctx, 1)

	// Assert
	assert.Nil(t, err)
}

// TestSDelete_IDNonExistent_Seller passes when the given id is not in database (return domain.Seller{} and error seller.ErrNotFound)
func TestSDelete_IDNonExistent_Seller(t *testing.T) {
	// Arrange
	expectedError := ServiceErrNotFound
	mockRepo := MockRepositorySeller{
		ErrorMock: expectedError,
	}
	service := NewService(&mockRepo)

	// Act
	result, err := service.Get(ctx, 1)

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.Equal(t, domain.Seller{}, result)
}

// TestDeleteFail_Seller passes when return an error for failed Delete function
func TestDeleteFail_Seller(t *testing.T) {
	// Arrange
	expectedError := errors.New("Id not exist")

	mockRepo := MockRepositorySeller{
		ErrorMock: expectedError,
	}
	service := NewService(&mockRepo)

	// Act
	err := service.Delete(ctx, 1)

	// Assert
	assert.EqualError(t, expectedError, err.Error())
}

// * ---------------------- Update ---------------------------
// TestUpdate_Seller passes when return seller updated
func TestUpdate_Seller(t *testing.T) {
	// Arrange
	database := []domain.Seller{{
		ID:          1,
		CID:         1,
		CompanyName: "Kiosco 1",
		Address:     "Junin 323",
		Telephone:   "2664727336",
		Locality_id: "5700",
	},
	}
	expectedSeller := domain.Seller{
		ID:          1,
		CID:         11,
		CompanyName: "Maxi Kiosco",
		Address:     "Junin 323",
		Telephone:   "2664727336",
		Locality_id: "5700",
	}

	mockRepo := MockRepositorySeller{
		Seller:   expectedSeller,
		DataMock: database,
	}
	service := NewService(&mockRepo)

	// Act
	result, err := service.Update(ctx, expectedSeller.ID, &expectedSeller.CID, &expectedSeller.CompanyName, &expectedSeller.Address, &expectedSeller.Telephone, &expectedSeller.Locality_id)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expectedSeller, result)
}

// TestUpdateFail_Service passes when return an error for db
func TestUpdateFail_Service(t *testing.T) {
	// Arrange
	expectedError := errors.New("Error on update")
	expectedSeller := domain.Seller{
		ID:          1,
		CID:         11,
		CompanyName: "Maxi Kiosco",
		Address:     "Junin 323",
		Telephone:   "2664727336",
		Locality_id: "5700",
	}

	mockRepo := MockRepositorySeller{
		ErrorMock: expectedError,
	}
	service := NewService(&mockRepo)

	// Act
	result, err := service.Update(ctx, expectedSeller.ID, &expectedSeller.CID, &expectedSeller.CompanyName, &expectedSeller.Address, &expectedSeller.Telephone, &expectedSeller.Locality_id)

	// Assert
	assert.EqualError(t, expectedError, err.Error())
	assert.Equal(t, domain.Seller{}, result)
}

// TestUpdateExistFail_Seller passes when return an error for failed Exist function
func TestUpdateExistFail_Seller(t *testing.T) {
	// Arrange
	expectedError := errors.New("cid already exists")
	expectedSeller := domain.Seller{
		ID:          1,
		CID:         11,
		CompanyName: "Maxi Kiosco",
		Address:     "Junin 323",
		Telephone:   "2664727336",
		Locality_id: "5700",
	}

	mockRepo := MockRepositorySeller{
		ErrorCidExist: expectedError,
	}
	service := NewService(&mockRepo)

	// Act
	result, err := service.Update(ctx, expectedSeller.ID, &expectedSeller.CID, &expectedSeller.CompanyName, &expectedSeller.Address, &expectedSeller.Telephone, &expectedSeller.Locality_id)

	// Assert
	assert.EqualError(t, expectedError, err.Error())
	assert.Equal(t, domain.Seller{}, result)
}

// TestUpdateErrorFail_Seller passes when return an error for failed Update function
func TestUpdateErrorFail_Seller(t *testing.T) {
	// Arrange
	expectedError := errors.New("Error on update")
	database := []domain.Seller{{
		ID:          1,
		CID:         1,
		CompanyName: "Kiosco 1",
		Address:     "Junin 323",
		Telephone:   "2664727336",
		Locality_id: "5700",
	}, {
		ID:          2,
		CID:         2,
		CompanyName: "Kiosco 2",
		Address:     "Colon 664",
		Telephone:   "2664233242",
		Locality_id: "5700",
	},
	}

	expectedSeller := domain.Seller{
		ID:          1,
		CID:         11,
		CompanyName: "Maxi Kiosco",
		Address:     "Junin 323",
		Telephone:   "2664727336",
		Locality_id: "5700",
	}

	mockRepo := MockRepositorySeller{
		DataMock:  database,
		ErrorMock: expectedError,
	}
	service := NewService(&mockRepo)

	// Act
	result, err := service.Update(ctx, expectedSeller.ID, &expectedSeller.CID, &expectedSeller.CompanyName, &expectedSeller.Address, &expectedSeller.Telephone, &expectedSeller.Locality_id)

	// Assert
	assert.EqualError(t, expectedError, err.Error())
	assert.Equal(t, domain.Seller{}, result)
}
