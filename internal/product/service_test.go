package product

import (
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func newIntPointer(value int) *int {
	return &value
}

func setupProductServiceTest() (ctx *gin.Context) {
	logging.InitLog(nil)
	gin.SetMode("test")
	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	return
}

// TestService_Save_OK passes when data is correct (return domain.Product.ID and nil error)
func TestService_Save_OK(t *testing.T) {
	// Arrange
	id := 2
	save := domain.Product{
		Description:                    "Tomatoes",
		ExpirationRate:                 70,
		FreezingRate:                   20,
		Height:                         20.4,
		Length:                         10.3,
		NetWeight:                      0.0,
		ProductCode:                    "kasbkj8ats9aka9",
		RecommendedFreezingTemperature: -12.6,
		Width:                          6.7,
		ProductTypeID:                  3,
		SellerID:                       newIntPointer(5),
	}
	expected := save
	expected.ID = id

	// Act
	ctx := setupProductServiceTest()
	mockProductRepository := RepositoryMock{db: []domain.Product{}, ExpectedID: id}
	productService := NewService(&mockProductRepository)
	result, err := productService.Save(ctx, save)

	// Assert
	assert.True(t, mockProductRepository.FlagSave)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

// TestService_Save_Conflict passes when ProductCode already exists (return empty domain.Product and error ServiceErrAlreadyExists)
func TestService_Save_Conflict(t *testing.T) {
	// Arrange
	save := domain.Product{
		Description:                    "Tomatoes",
		ExpirationRate:                 70,
		FreezingRate:                   20,
		Height:                         20.4,
		Length:                         10.3,
		NetWeight:                      0.0,
		ProductCode:                    "kasbkj8ats9aka9",
		RecommendedFreezingTemperature: -12.6,
		Width:                          6.7,
		ProductTypeID:                  3,
		SellerID:                       newIntPointer(5),
	}
	expectedErr := ServiceErrAlreadyExists

	// Act
	ctx := setupProductServiceTest()
	mockProductRepository := RepositoryMock{db: []domain.Product{}, ForcedErrExists: expectedErr}
	productService := NewService(&mockProductRepository)
	result, err := productService.Save(ctx, save)

	// Assert
	assert.False(t, mockProductRepository.FlagSave)
	assert.EqualError(t, err, expectedErr.Error())
	assert.Equal(t, domain.Product{}, result)
}

// TestService_Save_InternalServerError passes when Save fails (return empty domain.Product and error ServiceErrInternal)
func TestService_Save_InternalServerError(t *testing.T) {
	// Arrange
	save := domain.Product{
		Description:                    "Tomatoes",
		ExpirationRate:                 70,
		FreezingRate:                   20,
		Height:                         20.4,
		Length:                         10.3,
		NetWeight:                      0.0,
		ProductCode:                    "kasbkj8ats9aka9",
		RecommendedFreezingTemperature: -12.6,
		Width:                          6.7,
		ProductTypeID:                  3,
		SellerID:                       newIntPointer(5),
	}
	expectedErr := ServiceErrInternal

	// Act
	ctx := setupProductServiceTest()
	mockProductRepository := RepositoryMock{db: []domain.Product{}, ForcedErrSave: RepositoryErrInternal}
	productService := NewService(&mockProductRepository)
	result, err := productService.Save(ctx, save)

	// Assert
	assert.True(t, mockProductRepository.FlagSave)
	assert.EqualError(t, err, expectedErr.Error())
	assert.Equal(t, domain.Product{}, result)
}

// TestService_Save_FailForeignKeyNotFound passes when Save fails because seller id is not found in database (return empty domain.Product and error ServiceErrForeignKeyNotFound)
func TestService_Save_FailForeignKeyNotFound(t *testing.T) {
	// Arrange
	save := domain.Product{}
	expectedErr := ServiceErrForeignKeyNotFound

	// Act
	ctx := setupProductServiceTest()
	mockProductRepository := RepositoryMock{db: []domain.Product{}, ForcedErrSave: RepositoryErrForeignKeyConstraint}
	productService := NewService(&mockProductRepository)
	result, err := productService.Save(ctx, save)

	// Assert
	assert.True(t, mockProductRepository.FlagSave)
	assert.EqualError(t, err, expectedErr.Error())
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

// TestService_GetAll_OK passes when there are no errors on repository layer (return slice of all domain.Product and nil error)
func TestService_GetAll_OK(t *testing.T) {
	// Arrange
	expected := []domain.Product{
		{
			ID:                             1,
			Description:                    "Tomatoes",
			ExpirationRate:                 70,
			FreezingRate:                   20,
			Height:                         20.4,
			Length:                         10.3,
			NetWeight:                      0.0,
			ProductCode:                    "kasbkj8ats9aka9",
			RecommendedFreezingTemperature: -12.6,
			Width:                          6.7,
			ProductTypeID:                  3,
			SellerID:                       newIntPointer(5),
		},
		{
			ID:                             2,
			Description:                    "Potatoes",
			ExpirationRate:                 70,
			FreezingRate:                   20,
			Height:                         20.4,
			Length:                         10.3,
			NetWeight:                      0.0,
			ProductCode:                    "kasbsts9aka9",
			RecommendedFreezingTemperature: -12.6,
			Width:                          6.7,
			ProductTypeID:                  3,
			SellerID:                       newIntPointer(5),
		},
	}

	// Act
	ctx := setupProductServiceTest()
	mockProductRepository := RepositoryMock{db: expected}
	productService := NewService(&mockProductRepository)
	result, err := productService.GetAll(ctx)

	// Assert
	assert.True(t, mockProductRepository.FlagGetAll)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

// TestService_GetAll_FailInternalErr passes when there is a problem in repository layer (return an empty slice of domain.Product and error ServiceErrInternal)
func TestService_GetAll_FailInternalErr(t *testing.T) {
	// Arrange
	expected := ServiceErrInternal

	// Act
	ctx := setupProductServiceTest()
	mockProductRepository := RepositoryMock{db: []domain.Product{}, ForcedErrGetAll: RepositoryErrInternal}
	productService := NewService(&mockProductRepository)
	result, err := productService.GetAll(ctx)

	// Assert
	assert.True(t, mockProductRepository.FlagGetAll)
	assert.EqualError(t, err, expected.Error())
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

// TestService_Get_OK passes when id exists (return domain.Product and nil error)
func TestService_Get_OK(t *testing.T) {
	// Arrange
	searchID := 1
	expected := []domain.Product{
		{
			ID:                             1,
			Description:                    "Tomatoes",
			ExpirationRate:                 70,
			FreezingRate:                   20,
			Height:                         20.4,
			Length:                         10.3,
			NetWeight:                      0.0,
			ProductCode:                    "kasbkj8ats9aka9",
			RecommendedFreezingTemperature: -12.6,
			Width:                          6.7,
			ProductTypeID:                  3,
			SellerID:                       newIntPointer(5),
		},
	}

	// Act
	ctx := setupProductServiceTest()
	mockProductRepository := RepositoryMock{db: expected}
	productService := NewService(&mockProductRepository)
	result, err := productService.Get(ctx, searchID)

	// Assert
	assert.True(t, mockProductRepository.FlagGet)
	assert.Nil(t, err)
	assert.Equal(t, expected[0], result)
}

// TestService_Get_IDNonExistent passes when the given id is not in database (return empty domain.Product and error ServiceErrNotFound)
func TestService_Get_IDNonExistent(t *testing.T) {
	// Arrange
	searchID := 2
	expectedErr := ServiceErrNotFound

	// Act
	ctx := setupProductServiceTest()
	mockProductRepository := RepositoryMock{db: []domain.Product{}, ForcedErrGet: RepositoryErrNotFound}
	productService := NewService(&mockProductRepository)
	result, err := productService.Get(ctx, searchID)

	// Assert
	assert.True(t, mockProductRepository.FlagGet)
	assert.EqualError(t, err, expectedErr.Error())
	assert.Equal(t, domain.Product{}, result)
}

// TestService_Get_FailErrInternal passes when there is a problem in repository layer (return empty domain.Product and error ServiceErrInternal)
func TestService_Get_FailErrInternal(t *testing.T) {
	// Arrange
	searchID := 2
	expectedErr := ServiceErrInternal

	// Act
	ctx := setupProductServiceTest()
	mockProductRepository := RepositoryMock{db: []domain.Product{}, ForcedErrGet: RepositoryErrInternal}
	productService := NewService(&mockProductRepository)
	result, err := productService.Get(ctx, searchID)

	// Assert
	assert.True(t, mockProductRepository.FlagGet)
	assert.EqualError(t, err, expectedErr.Error())
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

// TestService_PartialUpdate_OK passes when data is correct (return updated domain.Product and error nil)
func TestService_PartialUpdate_OK(t *testing.T) {
	// Arrange
	db := []domain.Product{
		{
			ID:                             1,
			Description:                    "Tomatoes",
			ExpirationRate:                 70,
			FreezingRate:                   20,
			Height:                         20.4,
			Length:                         10.3,
			NetWeight:                      0.0,
			ProductCode:                    "kasbkj8ats9aka9",
			RecommendedFreezingTemperature: -12.6,
			Width:                          6.7,
			ProductTypeID:                  3,
			SellerID:                       newIntPointer(5),
		},
	}
	searchID := 1
	expected := domain.Product{
		ID:                             1,
		Description:                    "Bananas",
		ExpirationRate:                 70,
		FreezingRate:                   20,
		Height:                         20.4,
		Length:                         10.3,
		NetWeight:                      13.0,
		ProductCode:                    "kasbkj8ats9aka9",
		RecommendedFreezingTemperature: -12.6,
		Width:                          6.7,
		ProductTypeID:                  3,
		SellerID:                       newIntPointer(5),
	}

	// Act
	ctx := setupProductServiceTest()
	mockProductRepository := RepositoryMock{db: db}
	productService := NewService(&mockProductRepository)
	result, err := productService.PartialUpdate(ctx, searchID, expected)

	// Assert
	assert.True(t, mockProductRepository.FlagUpdate)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

// TestService_PartialUpdate_OKDifferentProductCode passes when data is correct (return updated domain.Product and error nil)
func TestService_PartialUpdate_OKDifferentProductCode(t *testing.T) {
	// Arrange
	db := []domain.Product{
		{
			ID:                             1,
			Description:                    "Tomatoes",
			ExpirationRate:                 70,
			FreezingRate:                   20,
			Height:                         20.4,
			Length:                         10.3,
			NetWeight:                      0.0,
			ProductCode:                    "kasbkj8ats9aka9",
			RecommendedFreezingTemperature: -12.6,
			Width:                          6.7,
			ProductTypeID:                  3,
			SellerID:                       newIntPointer(5),
		},
	}
	searchID := 1
	expected := domain.Product{
		ID:                             1,
		Description:                    "Bananas",
		ExpirationRate:                 70,
		FreezingRate:                   20,
		Height:                         20.4,
		Length:                         10.3,
		NetWeight:                      13.0,
		ProductCode:                    "kasbkj8ats9aka10",
		RecommendedFreezingTemperature: -12.6,
		Width:                          6.7,
		ProductTypeID:                  3,
		SellerID:                       newIntPointer(5),
	}

	// Act
	ctx := setupProductServiceTest()
	mockProductRepository := RepositoryMock{db: db}
	productService := NewService(&mockProductRepository)
	result, err := productService.PartialUpdate(ctx, searchID, expected)

	// Assert
	assert.True(t, mockProductRepository.FlagUpdate)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

// TestService_PartialUpdate_IDNonExistent passes when the given id is not in database (return empty domain.Product and error ServiceErrNotFound)
func TestService_PartialUpdate_IDNonExistent(t *testing.T) {
	// Arrange
	searchID := 1
	expectedErr := ServiceErrNotFound
	expected := domain.Product{}

	// Act
	ctx := setupProductServiceTest()
	mockProductRepository := RepositoryMock{db: []domain.Product{}, ForcedErrGet: RepositoryErrNotFound}
	productService := NewService(&mockProductRepository)
	result, err := productService.PartialUpdate(ctx, searchID, expected)

	// Assert
	assert.False(t, mockProductRepository.FlagUpdate)
	assert.EqualError(t, err, expectedErr.Error())
	assert.Equal(t, expected, result)
}

// TestService_PartialUpdate_DifferentButRepeatedProductCode passes when ProductCode is different to the current one but another product already has it (return empty domain.Product and error ServiceErrAlreadyExists)
func TestService_PartialUpdate_DifferentButRepeatedProductCode(t *testing.T) {
	// Arrange
	update := domain.Product{
		ProductCode: "oiaois",
	}
	searchID := 1
	expectedErr := ServiceErrAlreadyExists

	// Act
	ctx := setupProductServiceTest()
	mockProductRepository := RepositoryMock{db: []domain.Product{}, ForcedErrExists: expectedErr}
	productService := NewService(&mockProductRepository)
	result, err := productService.PartialUpdate(ctx, searchID, update)

	// Assert
	assert.False(t, mockProductRepository.FlagUpdate)
	assert.EqualError(t, err, expectedErr.Error())
	assert.Equal(t, domain.Product{}, result)
}

// TestService_PartialUpdate_InternalServerError passes when Update fails (return empty domain.Product and error ServiceErrInternal)
func TestService_PartialUpdate_InternalServerError(t *testing.T) {
	// Arrange
	searchID := 1
	expectedErr := ServiceErrInternal

	// Act
	ctx := setupProductServiceTest()
	mockProductRepository := RepositoryMock{db: []domain.Product{}, ForcedErrUpdate: RepositoryErrInternal}
	productService := NewService(&mockProductRepository)
	result, err := productService.PartialUpdate(ctx, searchID, domain.Product{})

	// Assert
	assert.True(t, mockProductRepository.FlagUpdate)
	assert.EqualError(t, err, expectedErr.Error())
	assert.Equal(t, domain.Product{}, result)
}

// TestService_PartialUpdate_FailForeignKeyNotFound passes when Update fails because seller id is not found in database (return empty domain.Product and error ServiceErrForeignKeyNotFound)
func TestService_PartialUpdate_FailForeignKeyNotFound(t *testing.T) {
	// Arrange
	searchID := 1
	expectedErr := ServiceErrForeignKeyNotFound

	// Act
	ctx := setupProductServiceTest()
	mockProductRepository := RepositoryMock{db: []domain.Product{}, ForcedErrUpdate: RepositoryErrForeignKeyConstraint}
	productService := NewService(&mockProductRepository)
	result, err := productService.PartialUpdate(ctx, searchID, domain.Product{})

	// Assert
	assert.True(t, mockProductRepository.FlagUpdate)
	assert.EqualError(t, err, expectedErr.Error())
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

// TestService_Delete_OK passes when id exists and deletion is successful (return nil error)
func TestService_Delete_OK(t *testing.T) {
	// Arrange
	searchID := 1

	// Act
	ctx := setupProductServiceTest()
	mockProductRepository := RepositoryMock{db: []domain.Product{}}
	productService := NewService(&mockProductRepository)
	err := productService.Delete(ctx, searchID)

	// Assert
	assert.True(t, mockProductRepository.FlagDelete)
	assert.Nil(t, err)
}

// TestService_Delete_IDNonExistent passes when the given id is not in database (return error ServiceErrNotFound)
func TestService_Delete_IDNonExistent(t *testing.T) {
	// Arrange
	searchID := 3
	expectedErr := ServiceErrNotFound

	// Act
	ctx := setupProductServiceTest()
	mockProductRepository := RepositoryMock{db: []domain.Product{}, ForcedErrDelete: RepositoryErrNotFound}
	productService := NewService(&mockProductRepository)
	err := productService.Delete(ctx, searchID)

	// Assert
	assert.True(t, mockProductRepository.FlagDelete)
	assert.EqualError(t, err, expectedErr.Error())
}

// TestService_Delete_FailInternalServerError passes when Delete fails (return error ServiceErrInternal)
func TestService_Delete_FailInternalServerError(t *testing.T) {
	// Arrange
	searchID := 3
	expectedErr := ServiceErrInternal

	// Act
	ctx := setupProductServiceTest()
	mockProductRepository := RepositoryMock{db: []domain.Product{}, ForcedErrDelete: RepositoryErrInternal}
	productService := NewService(&mockProductRepository)
	err := productService.Delete(ctx, searchID)

	// Assert
	assert.True(t, mockProductRepository.FlagDelete)
	assert.EqualError(t, err, expectedErr.Error())
}
