package locality

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

var locality_test = domain.Locality{
	ID:           "5700",
	LocalityName: "Capital",
	ProvinceName: "San Luis",
	CountryName:  "Argentina",
}

func TestExist_OK(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"id"}
	rows := sqlmock.NewRows(columns)

	rows.AddRow(locality_test.ID)
	mock.ExpectQuery(regexp.QuoteMeta(EXIST_LOCALITY)).WithArgs(locality_test.ID).WillReturnRows(rows)

	// Act
	repo := NewRepository(db)
	result := repo.Exists(context.TODO(), locality_test.ID)

	// Assert
	assert.True(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// * ---------------------- Save --------------------------
// TestSave_Locality_OK passes when return locality created
func TestSave_Locality_OK(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(SAVE_LOCALITY))
	mock.ExpectExec(regexp.QuoteMeta(SAVE_LOCALITY)).WillReturnResult(sqlmock.NewResult(1, 1))

	columns := []string{"id", "locality_name", "province_name", "country_name"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(locality_test.ID, locality_test.LocalityName, locality_test.ProvinceName, locality_test.CountryName)

	// Act
	repository := NewRepository(db)

	newID, err := repository.Save(context.TODO(), locality_test)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, locality_test.ID, newID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSave_Locality_Fail passes when locality id doesn't exist in the database
func TestSave_Locality_FailForeignKeyConstraint(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(SAVE_LOCALITY))
	mock.ExpectExec(regexp.QuoteMeta(SAVE_LOCALITY)).WillReturnError(&mysql.MySQLError{Number: MySqlNumberForeignKeyConstraint})

	// Act
	repository := NewRepository(db)

	newID, err := repository.Save(context.TODO(), locality_test)

	// Assert
	assert.EqualError(t, err, ErrForeignKeyConstraint.Error())
	assert.Equal(t, "0", newID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSave_Locality_FailPrepare passes when function Prepare returns an error
func TestSave_Locality_FailPrepare(t *testing.T) {
	// Arrange
	locality := domain.Locality{
		ID:           "5700",
		LocalityName: "San Luis",
		ProvinceName: "San Luis",
		CountryName:  "Argentina",
	}
	expectedError := ErrInternal

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(SAVE_LOCALITY)).WillReturnError(expectedError)

	// Act
	repository := NewRepository(db)

	_, err = repository.Save(context.TODO(), locality)

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	// assert.Empty(t, newID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSave_Locality_FailErrInternal(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta(SAVE_LOCALITY))
	mock.ExpectExec(regexp.QuoteMeta(SAVE_LOCALITY)).WillReturnError(sql.ErrConnDone)

	// Act
	repository := NewRepository(db)

	newID, err := repository.Save(context.TODO(), locality_test)

	// Assert
	assert.EqualError(t, err, ErrInternal.Error())
	assert.Equal(t, "0", newID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// * ---------------------- Get ----------------------------
// TestGet_Locality_OK passes when return correct locality
func TestGet_Locality_OK(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	column := []string{"id", "locality_name", "province_name", "country_name"}
	rows := sqlmock.NewRows(column)
	locality := domain.Locality{ID: locality_test.ID, LocalityName: locality_test.LocalityName, ProvinceName: locality_test.ProvinceName, CountryName: locality_test.CountryName}

	rows.AddRow(locality.ID, locality.LocalityName, locality.ProvinceName, locality.CountryName)

	mock.ExpectQuery(regexp.QuoteMeta(GET_LOCALITY)).WillReturnRows(rows)

	// Act
	repo := NewRepository(db)
	result, err := repo.Get(context.TODO(), locality.ID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, locality, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGet_Locality_FailErrNotFound(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	locality := domain.Locality{ID: locality_test.ID, LocalityName: locality_test.LocalityName, ProvinceName: locality_test.ProvinceName, CountryName: locality_test.CountryName}

	mock.ExpectQuery(regexp.QuoteMeta(GET_LOCALITY)).WithArgs(locality_test.ID).WillReturnError(sql.ErrNoRows)

	// Act
	repo := NewRepository(db)
	_, err = repo.Get(context.TODO(), locality.ID)

	// Assert
	assert.EqualError(t, err, ErrNotFound.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestGet_Locality_FailErrInternal(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	locality := domain.Locality{ID: locality_test.ID, LocalityName: locality_test.LocalityName, ProvinceName: locality_test.ProvinceName, CountryName: locality_test.CountryName}

	mock.ExpectQuery(regexp.QuoteMeta(GET_LOCALITY)).WithArgs(locality_test.ID).WillReturnError(sql.ErrConnDone)

	// Act
	repo := NewRepository(db)
	res, err := repo.Get(context.TODO(), locality.ID)

	// Assert
	assert.EqualError(t, err, ErrInternal.Error())
	assert.Empty(t, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// * ---------------------- Report --------------------------
// TestReportSellers_OK passes when return corrects sellers by locality
func TestReportSellers_OK(t *testing.T) {
	// Arrange
	testReport := domain.ReportSellers{
		LocalityID:   "5700",
		LocalityName: "San Luis",
		SellersCount: 2,
	}
	expectedReport := []domain.ReportSellers{testReport}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"locality_id", "locality_name", "sellers_count"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(testReport.LocalityID, testReport.LocalityName, testReport.SellersCount)

	mock.ExpectQuery(regexp.QuoteMeta(GET_SELLERS)).WillReturnRows(rows)

	// Act
	repository := NewRepository(db)
	report, err := repository.ReportSellers(context.TODO())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedReport, report)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReportSellers_FailQuery passes when function Query returns an error
func TestReportSellers_FailQuery(t *testing.T) {
	// Arrange
	expectedError := ErrInternal

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(GET_SELLERS)).WillReturnError(expectedError)

	// Act
	repository := NewRepository(db)
	report, err := repository.ReportSellers(context.TODO())

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, report)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReportSellers_FailScan passes when function Scan returns an error
func TestReportSellers_FailScan(t *testing.T) {
	// Arrange
	expectedError := ErrInternal

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"locality_id", "locality_name", "sellers_count"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(nil, nil, nil) // this can't be parsed by Scan function

	mock.ExpectQuery(regexp.QuoteMeta(GET_SELLERS)).WillReturnRows(rows)

	// Act
	repository := NewRepository(db)
	report, err := repository.ReportSellers(context.TODO())

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, report)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReportSellersByLocationID_OK passes when return corrects sellers by locality id
func TestReportSellersByLocationID_OK(t *testing.T) {
	// Arrange
	testReport := domain.ReportSellers{
		LocalityID:   "5700",
		LocalityName: "San Luis",
		SellersCount: 2,
	}
	expectedReport := []domain.ReportSellers{testReport}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"locality_id", "locality_name", "sellers_count"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(testReport.LocalityID, testReport.LocalityName, testReport.SellersCount)

	mock.ExpectQuery(regexp.QuoteMeta(GET_SELLERS_BY_LOCATION_ID)).WillReturnRows(rows)

	// Act
	repository := NewRepository(db)
	report, err := repository.ReportSellersByLocationID(context.TODO(), testReport.LocalityID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedReport, report)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReportSellersByLocationID_FailNoRows passes when there is no rows to Scan
func TestReportSellersByLocationID_FailNoRows(t *testing.T) {
	// Arrange
	expectedError := ErrNotFound

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"locality_id", "locality_name", "sellers_count"}
	rows := sqlmock.NewRows(columns)

	mock.ExpectQuery(regexp.QuoteMeta(GET_SELLERS_BY_LOCATION_ID)).WillReturnRows(rows)

	// Act
	repository := NewRepository(db)
	report, err := repository.ReportSellersByLocationID(context.TODO(), "5700")

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, report)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReportSellersByLocationID_FailInternal passes when function Scan cant parsed the values
func TestReportSellersByLocationID_FailInternal(t *testing.T) {
	// Arrange
	expectedError := ErrInternal

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"locality_id", "locality_name", "sellers_count"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(nil, nil, nil)

	mock.ExpectQuery(regexp.QuoteMeta(GET_SELLERS_BY_LOCATION_ID)).WillReturnRows(rows)

	// Act
	repository := NewRepository(db)
	report, err := repository.ReportSellersByLocationID(context.TODO(), "5700")

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, report)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ---------- ReportCarries ----------

// TestRepositoryReportCarries checks the correct operation of the ReportCarries repository method
func TestRepositoryReportCarries(t *testing.T) {
	// Arrange
	testReport := domain.ReportCarries{
		LocalityID:   "0001",
		LocalityName: "Santiago",
		CarriesCount: 2,
	}
	expectedReport := []domain.ReportCarries{testReport}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"locality_id", "locality_name", "carries_count"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(testReport.LocalityID, testReport.LocalityName, testReport.CarriesCount)

	mock.ExpectQuery(regexp.QuoteMeta(GET_CARRIES)).WillReturnRows(rows)

	// Act
	repository := NewRepository(db)
	report, err := repository.ReportCarries(context.TODO())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedReport, report)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepositoryReportCarriesFailQuery is correct when function Query returns an error
func TestRepositoryReportCarriesFailQuery(t *testing.T) {
	// Arrange
	expectedError := ErrInternal

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(GET_CARRIES)).WillReturnError(expectedError)

	// Act
	repository := NewRepository(db)
	report, err := repository.ReportCarries(context.TODO())

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, report)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepositoryReportCarriesFailScan is correct when function Scan returns an error
func TestRepositoryReportCarriesFailScan(t *testing.T) {
	// Arrange
	expectedError := ErrInternal

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"locality_id", "locality_name", "carries_count"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(nil, nil, nil) // this can't be parsed by Scan function

	mock.ExpectQuery(regexp.QuoteMeta(GET_CARRIES)).WillReturnRows(rows)

	// Act
	repository := NewRepository(db)
	report, err := repository.ReportCarries(context.TODO())

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, report)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ---------- ReportCarriesByLocationID ----------

// TestRepositoryReportCarriesByLocationID checks the correct operation of the ReportCarries repository method
func TestRepositoryReportCarriesByLocationID(t *testing.T) {
	// Arrange
	testReport := domain.ReportCarries{
		LocalityID:   "0001",
		LocalityName: "Santiago",
		CarriesCount: 2,
	}
	expectedReport := []domain.ReportCarries{testReport}
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"locality_id", "locality_name", "carries_count"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(testReport.LocalityID, testReport.LocalityName, testReport.CarriesCount)

	mock.ExpectQuery(regexp.QuoteMeta(GET_CARRIES_BY_LOCATION_ID)).WillReturnRows(rows)

	// Act
	repository := NewRepository(db)
	report, err := repository.ReportCarriesByLocationID(context.TODO(), testReport.LocalityID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedReport, report)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepositoryReportCarriesByLocationIDFailNoRows is correct when there is no rows to Scan
func TestRepositoryReportCarriesByLocationIDFailNoRows(t *testing.T) {
	// Arrange
	expectedError := ErrNotFound

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"locality_id", "locality_name", "carries_count"}
	rows := sqlmock.NewRows(columns)

	mock.ExpectQuery(regexp.QuoteMeta(GET_CARRIES_BY_LOCATION_ID)).WillReturnRows(rows)

	// Act
	repository := NewRepository(db)
	report, err := repository.ReportCarriesByLocationID(context.TODO(), "0001")

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, report)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestRepositoryReportCarriesByLocationIDFailInternal is correct when function Scan cant parsed the values
func TestRepositoryReportCarriesByLocationIDFailInternal(t *testing.T) {
	// Arrange
	expectedError := ErrInternal

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	columns := []string{"locality_id", "locality_name", "carries_count"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(nil, nil, nil) // this can't be parsed by Scan function

	mock.ExpectQuery(regexp.QuoteMeta(GET_CARRIES_BY_LOCATION_ID)).WillReturnRows(rows)

	// Act
	repository := NewRepository(db)
	report, err := repository.ReportCarriesByLocationID(context.TODO(), "0001")

	// Assert
	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, report)
	assert.NoError(t, mock.ExpectationsWereMet())
}
