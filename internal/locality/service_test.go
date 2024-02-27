package locality

import (
	"context"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/stretchr/testify/assert"
)

var ctx context.Context

func init() {
	logging.InitLog(nil)
}

type MockRepositoryLocality struct {
	ReportCarry  []domain.ReportCarries
	Locality     domain.Locality
	DataMock     []domain.Locality
	Report       []domain.ReportSellers
	ErrorMock    error
	ErrorIdExist error
}

func (r *MockRepositoryLocality) ReportCarries(ctx context.Context) ([]domain.ReportCarries, error) {
	if r.ErrorMock != nil {
		return []domain.ReportCarries{}, r.ErrorMock
	}
	return r.ReportCarry, nil
}

func (r *MockRepositoryLocality) ReportCarriesByLocationID(ctx context.Context, id string) ([]domain.ReportCarries, error) {
	if r.ErrorMock != nil {
		return []domain.ReportCarries{}, r.ErrorMock
	}
	return r.ReportCarry, nil
}

func (r *MockRepositoryLocality) Get(ctx context.Context, id string) (l domain.Locality, err error) {
	if r.ErrorMock != nil {
		err = r.ErrorMock
		return
	}
	l = r.Locality
	return
}

func (r *MockRepositoryLocality) Save(ctx context.Context, l domain.Locality) (id string, err error) {
	if r.ErrorMock != nil {
		err = r.ErrorMock
		return
	}
	id = "5700"
	return
}

func (r *MockRepositoryLocality) Exists(ctx context.Context, id string) (exist bool) {
	exist = r.ErrorIdExist != nil
	return
}

func (r *MockRepositoryLocality) ReportSellers(ctx context.Context) (report []domain.ReportSellers, err error) {
	if r.ErrorMock != nil {
		err = r.ErrorMock
		return
	}
	report = r.Report
	return
}

func (r *MockRepositoryLocality) ReportSellersByLocationID(ctx context.Context, id string) (report []domain.ReportSellers, err error) {
	if r.ErrorMock != nil {
		err = r.ErrorMock
		return
	}
	report = r.Report
	return

}

// * ---------------------- ReportCarries --------------------------

// TestReportCarries checks the correct operation of the ReportCarries service method when a non nil locality id is recieved
func TestReportCarries(t *testing.T) {
	// Arrange
	testReport := domain.ReportCarries{
		LocalityID:   "0001",
		LocalityName: "Santiago",
		CarriesCount: 2,
	}
	expectedReport := []domain.ReportCarries{testReport}
	mockRepo := MockRepositoryLocality{
		ReportCarry: expectedReport,
		ErrorMock:   nil,
	}
	service := NewService(&mockRepo)

	// Act
	result, err := service.ReportCarries(ctx, &testReport.LocalityID)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expectedReport, result)
}

// TestReportCarriesNil checks the correct operation of the ReportCarries service method when a nil locality id is recieved
func TestReportCarriesNil(t *testing.T) {
	// Arrange
	testReport := domain.ReportCarries{
		LocalityID:   "0001",
		LocalityName: "Santiago",
		CarriesCount: 2,
	}
	expectedReport := []domain.ReportCarries{testReport}
	mockRepo := MockRepositoryLocality{
		ReportCarry: expectedReport,
		ErrorMock:   nil,
	}
	service := NewService(&mockRepo)

	// Act
	result, err := service.ReportCarries(ctx, nil)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expectedReport, result)
}

// * ---------------------- Create --------------------------
// TestCreate_Locality_OK passes when return locality created
func TestCreate_Locality_OK(t *testing.T) {
	// Arrange
	localityToCreate := domain.Locality{
		ID:           "5700",
		LocalityName: "Capital",
		ProvinceName: "San Luis",
		CountryName:  "Argentina",
	}

	localityExpected := domain.Locality{
		ID:           "5700",
		LocalityName: "Capital",
		ProvinceName: "San Luis",
		CountryName:  "Argentina",
	}

	mockRepo := MockRepositoryLocality{
		Locality:  localityExpected,
		ErrorMock: nil,
	}
	service := NewService(&mockRepo)

	// Act
	result, err := service.Create(ctx, localityToCreate)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, localityExpected, result)
}

// TestCreateLocality_Fail passes when return an error for failed Create function
func TestCreateLocality_Fail(t *testing.T) {
	// Arrange
	expectedError := ErrAlreadyExists
	localityToCreate := domain.Locality{}

	mockRepo := MockRepositoryLocality{
		ErrorMock: expectedError,
	}
	service := NewService(&mockRepo)

	// Act
	_, err := service.Create(ctx, localityToCreate)

	// Assert
	assert.EqualError(t, expectedError, err.Error())
}

// TestCreateLocality_ExistFail passes when return an error for failed Exist function
func TestCreateLocality_ExistFail(t *testing.T) {
	// Arrange
	expectedError := ErrAlreadyExists
	localityToCreate := domain.Locality{}

	mockRepo := MockRepositoryLocality{
		ErrorIdExist: expectedError,
	}
	service := NewService(&mockRepo)

	// Act
	_, err := service.Create(ctx, localityToCreate)

	// Assert
	assert.EqualError(t, expectedError, err.Error())
}

// * ---------------------- Get ----------------------------
// TestGetLocality_OK passes when return correct locality
func TestGetLocality_OK(t *testing.T) {
	// Arrange
	database := []domain.Locality{{
		ID:           "5700",
		LocalityName: "Capital",
		ProvinceName: "San Luis",
		CountryName:  "Argentina",
	}, {
		ID:           "1234",
		LocalityName: "Capital",
		ProvinceName: "Neuquen",
		CountryName:  "Argentina",
	},
	}

	expectedLocality := domain.Locality{
		ID:           "5700",
		LocalityName: "Capital",
		ProvinceName: "San Luis",
		CountryName:  "Argentina",
	}

	mockRepo := MockRepositoryLocality{
		Locality:  expectedLocality,
		DataMock:  database,
		ErrorMock: nil,
	}
	service := NewService(&mockRepo)

	// Act
	result, err := service.Get(ctx, expectedLocality.ID)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expectedLocality, result)
}

// * ---------------------- Report --------------------------
// TestReportSellers_ok passes when return correct ReportSellers
func TestReportSellers_ok(t *testing.T) {
	// Arrange
	testReport := domain.ReportSellers{
		LocalityID:   "5700",
		LocalityName: "San Luis",
		SellersCount: 3,
	}
	expectedReport := []domain.ReportSellers{testReport}
	mockRepo := MockRepositoryLocality{
		Report: expectedReport,
	}
	service := NewService(&mockRepo)

	// Act
	result, err := service.ReportSellers(ctx, &testReport.LocalityID)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expectedReport, result)
}

// TestReportSellersNil passes when return correct ReportSellers without locality id
func TestReportSellersNil(t *testing.T) {
	// Arrange
	testReport := domain.ReportSellers{
		LocalityID:   "5700",
		LocalityName: "San Luis",
		SellersCount: 3,
	}
	expectedReport := []domain.ReportSellers{testReport}
	mockRepo := MockRepositoryLocality{
		Report:    expectedReport,
		ErrorMock: nil,
	}
	service := NewService(&mockRepo)

	// Act
	result, err := service.ReportSellers(ctx, nil)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expectedReport, result)
}
