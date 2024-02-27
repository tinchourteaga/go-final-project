package logging

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

var logMsg = LogAttributes{
	timestamp:      "test timestamp",
	user:           "test user",
	file:           "test file",
	line:           "test line",
	callerFunction: "test function",
	msg:            "test message",
}

func TestSingleLogger_Save_OK(t *testing.T) {
	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	InitLog(db)
	mock.ExpectPrepare(
		regexp.QuoteMeta(SaveLog)).
		ExpectExec().
		WithArgs(logMsg.timestamp, logMsg.user, logMsg.file, logMsg.line, logMsg.callerFunction, logMsg.msg).
		WillReturnResult(sqlmock.NewResult(1, 1))
	_, errSave := singleLoggerInstance.Save(logMsg)

	// Assert
	assert.NoError(t, errSave)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSingleLogger_Save_FailQuery(t *testing.T) {
	// Arrange
	expectedErr := errors.New("forced query error")

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	InitLog(db)
	mock.ExpectPrepare(regexp.QuoteMeta(SaveLog)).WillReturnError(expectedErr)
	_, errSave := singleLoggerInstance.Save(logMsg)

	// Assert
	assert.EqualError(t, errSave, expectedErr.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSingleLogger_Save_FailParsingError(t *testing.T) {
	// Assert
	expectedErr := errors.New("forced non mysql error")

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	InitLog(db)
	mock.ExpectPrepare(regexp.QuoteMeta(SaveLog)).
		ExpectExec().
		WithArgs(logMsg.timestamp, logMsg.user, logMsg.file, logMsg.line, logMsg.callerFunction, logMsg.msg).
		WillReturnError(expectedErr)
	_, errSave := singleLoggerInstance.Save(logMsg)

	// Assert
	assert.EqualError(t, errSave, expectedErr.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSingleLogger_Save_FailLastID(t *testing.T) {
	// Arrange
	expectedErr := errors.New("forced last id error")

	// Act
	db, mock, errSql := sqlmock.New()
	assert.NoError(t, errSql)
	defer db.Close()
	InitLog(db)
	mock.ExpectPrepare(regexp.QuoteMeta(SaveLog)).
		ExpectExec().
		WillReturnResult(sqlmock.NewErrorResult(expectedErr))
	saveResultID, errSave := singleLoggerInstance.Save(logMsg)

	// Assert
	assert.EqualError(t, errSave, expectedErr.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, saveResultID, 0)
}
