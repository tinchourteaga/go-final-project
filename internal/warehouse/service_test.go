package warehouse

import (
	"net/http/httptest"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	logging.InitLog(nil)
}

// TestGetAll checks the correct operation of the GetAll service method
func TestGetAll(t *testing.T) {
	//arrange
	expectedWarehouses := []domain.Warehouse{}
	mockRepo := MockRepo{mockWarehouses: expectedWarehouses}
	service := NewService(&mockRepo)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	//act
	result, _ := service.GetAll(ctx)
	//assert
	assert.Equal(t, expectedWarehouses, result)
}

// TestGetAllFailure is correct when repository returns an error
func TestGetAllFailure(t *testing.T) {
	//arrange
	expectedError := ErrAlreadyExists
	mockRepo := MockRepo{mockErrorInternal: expectedError}
	service := NewService(&mockRepo)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	//act
	_, err := service.GetAll(ctx)
	//assert
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
}

// TestGet checks the correct operation of the Get service method
func TestGet(t *testing.T) {
	//arrange
	expectedWarehouse := domain.Warehouse{
		ID:                 1,
		Address:            "Monroe 1230",
		Telephone:          "47470000",
		WarehouseCode:      "DHM1",
		MinimumCapacity:    10,
		MinimumTemperature: 0,
	}
	mockRepo := MockRepo{mockWarehouse: expectedWarehouse}
	service := NewService(&mockRepo)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	//act
	result, _ := service.Get(ctx, 1)
	//assert
	assert.Equal(t, expectedWarehouse, result)
}

// TestGetFailure is correct when repository returns an error
func TestGetFailure(t *testing.T) {
	//arrange
	expectedError := ErrInternal
	mockRepo := MockRepo{mockErrorInternal: expectedError}
	service := NewService(&mockRepo)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	//act
	_, err := service.Get(ctx, 1)
	//assert
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
}

// TestCreate checks the correct operation of the Create service method
func TestCreate(t *testing.T) {
	// arrange
	expectedWarehouse := domain.Warehouse{
		ID:                 1,
		Address:            "Monroe 1230",
		Telephone:          "47470000",
		WarehouseCode:      "DHM1",
		MinimumCapacity:    10,
		MinimumTemperature: 0,
	}
	mockRepo := MockRepo{mockWarehouse: expectedWarehouse}
	service := NewService(&mockRepo)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	// act
	result, _ := service.Create(ctx, expectedWarehouse.Address, expectedWarehouse.Telephone, expectedWarehouse.WarehouseCode, expectedWarehouse.MinimumCapacity, expectedWarehouse.MinimumTemperature)
	// assert
	assert.Equal(t, expectedWarehouse, result)
}

// TestCreateFailureSave is correct when repository function Save returns an error
func TestCreateFailureSave(t *testing.T) {
	//arrange
	warehouse := domain.Warehouse{
		ID:                 1,
		Address:            "Monroe 1230",
		Telephone:          "47470000",
		WarehouseCode:      "DHM1",
		MinimumCapacity:    10,
		MinimumTemperature: 0,
	}
	expectedError := ErrInternal
	mockRepo := MockRepo{mockErrorInternal: expectedError}
	service := NewService(&mockRepo)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	//act
	_, err := service.Create(ctx, warehouse.Address, warehouse.Telephone, warehouse.WarehouseCode, warehouse.MinimumCapacity, warehouse.MinimumTemperature)
	//assert
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
}

// TestCreateFailureExists is correct when repository function Exists returns true
func TestCreateFailureExists(t *testing.T) {
	//arrange
	warehouse := domain.Warehouse{
		ID:                 1,
		Address:            "Monroe 1230",
		Telephone:          "47470000",
		WarehouseCode:      "DHM1",
		MinimumCapacity:    10,
		MinimumTemperature: 0,
	}
	expectedError := ErrAlreadyExists
	mockRepo := MockRepo{mockErrorExists: expectedError}
	service := NewService(&mockRepo)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	//act
	_, err := service.Create(ctx, warehouse.Address, warehouse.Telephone, warehouse.WarehouseCode, warehouse.MinimumCapacity, warehouse.MinimumTemperature)
	//assert
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
}

// TestUpdate checks the correct operation of the Update service method
func TestUpdate(t *testing.T) {
	// arrange
	warehouse := domain.Warehouse{
		ID:                 1,
		Address:            "Monroe 1230",
		Telephone:          "47470000",
		WarehouseCode:      "DHM1",
		MinimumCapacity:    10,
		MinimumTemperature: 0,
	}
	mockRepo := MockRepo{mockWarehouse: warehouse}
	service := NewService(&mockRepo)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	// act
	result, _ := service.Update(ctx, warehouse.ID, &warehouse.Address, &warehouse.Telephone, &warehouse.WarehouseCode, &warehouse.MinimumCapacity, &warehouse.MinimumTemperature)
	// assert
	assert.Equal(t, warehouse, result)
}

// TestUpdateFailureGet is correct when repository function Get returns an error
func TestUpdateFailureGet(t *testing.T) {
	// arrange
	warehouse := domain.Warehouse{
		ID:                 1,
		Address:            "Monroe 1230",
		Telephone:          "47470000",
		WarehouseCode:      "DHM1",
		MinimumCapacity:    10,
		MinimumTemperature: 0,
	}
	expectedError := ErrInternal
	mockRepo := MockRepo{mockErrorInternal: expectedError}
	service := NewService(&mockRepo)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	// act
	_, err := service.Update(ctx, warehouse.ID, &warehouse.Address, &warehouse.Telephone, &warehouse.WarehouseCode, &warehouse.MinimumCapacity, &warehouse.MinimumTemperature)
	// assert
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
}

// TestUpdateFailureExists is correct when repository function Exists returns true
func TestUpdateFailureExists(t *testing.T) {
	// arrange
	warehouse := domain.Warehouse{
		ID:                 1,
		Address:            "Monroe 1230",
		Telephone:          "47470000",
		WarehouseCode:      "DHM1",
		MinimumCapacity:    10,
		MinimumTemperature: 0,
	}
	expectedError := ErrAlreadyExists
	mockRepo := MockRepo{mockErrorExists: expectedError}
	service := NewService(&mockRepo)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	// act
	_, err := service.Update(ctx, warehouse.ID, &warehouse.Address, &warehouse.Telephone, &warehouse.WarehouseCode, &warehouse.MinimumCapacity, &warehouse.MinimumTemperature)
	// assert
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
}

// TestUpdateFailureUpdate is correct when repository function Update returns an error
func TestUpdateFailureUpdate(t *testing.T) {
	// arrange
	warehouse := domain.Warehouse{
		ID:                 1,
		Address:            "Monroe 1230",
		Telephone:          "47470000",
		WarehouseCode:      "DHM1",
		MinimumCapacity:    10,
		MinimumTemperature: 0,
	}
	expectedError := ErrInternal
	mockRepo := MockRepo{mockErrorUpdate: expectedError}
	service := NewService(&mockRepo)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	// act
	_, err := service.Update(ctx, warehouse.ID, &warehouse.Address, &warehouse.Telephone, &warehouse.WarehouseCode, &warehouse.MinimumCapacity, &warehouse.MinimumTemperature)
	// assert
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
}

// TestDelete checks the correct operation of the Update service method
func TestDelete(t *testing.T) {
	// arrange
	mockRepo := MockRepo{}
	service := NewService(&mockRepo)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	// act
	err := service.Delete(ctx, 1)
	// assert
	assert.Nil(t, err)
}

// TestDeleteFailureInternal is correct when repository returns an error
func TestDeleteFailureInternal(t *testing.T) {
	// arrange
	expectedError := ErrInternal
	mockRepo := MockRepo{mockErrorInternal: expectedError}
	service := NewService(&mockRepo)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	// act
	err := service.Delete(ctx, 1)
	// assert
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
}
