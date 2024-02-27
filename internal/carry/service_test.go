package carry

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

// TestSave checks the correct operation of the Save service method
func TestSave(t *testing.T) {
	// arrange
	carry := domain.Carry{
		ID:          1,
		CID:         "CID#1",
		CompanyName: "some name",
		Address:     "corrientes 800",
		Telephone:   "4567-4567",
		Locality_id: "1",
	}
	mockRepo := MockRepo{}
	service := NewService(&mockRepo)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	// act
	result, _ := service.Save(ctx, carry.CID, carry.CompanyName, carry.Address, carry.Telephone, carry.Locality_id)
	// assert
	assert.Equal(t, carry, result)
}

// TestSaveFailureSave is correct when repository function Save returns an error
func TestSaveFailureSave(t *testing.T) {
	//arrange
	carry := domain.Carry{
		ID:          1,
		CID:         "CID#1",
		CompanyName: "some name",
		Address:     "corrientes 800",
		Telephone:   "4567-4567",
		Locality_id: "1",
	}
	expectedError := ErrInternal
	mockRepo := MockRepo{mockError: expectedError}
	service := NewService(&mockRepo)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	//act
	_, err := service.Save(ctx, carry.CID, carry.CompanyName, carry.Address, carry.Telephone, carry.Locality_id)
	//assert
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
}
