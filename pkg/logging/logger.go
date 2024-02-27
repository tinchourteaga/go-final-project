package logging

import (
	"database/sql"
	"fmt"
	"log"
	"os/user"
	"runtime"
	"strconv"
	"time"
)

const (
	SaveLog = `INSERT INTO logs (time_stamp, user, file_path, function_line, caller_function, msg) VALUES (?, ?, ?, ?, ?, ?);`
)

type LogAttributes struct {
	timestamp      string
	user           string
	file           string
	line           string
	callerFunction string
	msg            string
}

type Logger interface {
	Log(msg interface{})
	Save(logMsg LogAttributes) (int, error)
}

type SingleLogger struct {
	db *sql.DB
}

var singleLoggerInstance Logger

func InitLog(db *sql.DB) {
	singleLoggerInstance = SingleLogger{db: db}
}

func Default() Logger {
	return singleLoggerInstance
}

func Log(message interface{}) {
	Default().Log(message)
}

// Log receives the message to be logged. It appends the message to the timestamp, user and caller function
func (sl SingleLogger) Log(msg interface{}) {
	msgString := fmt.Sprintf("Error: %s", msg)

	// Get current datetime
	timestamp := time.Now().Format(time.RFC3339)

	// Get current os user
	currentUser, err := user.Current()

	if err != nil {
		log.Fatal("cannot obtain current user")
	}

	username := currentUser.Username

	// Get caller function, including file and line
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	file := frame.File
	line := strconv.Itoa(frame.Line)
	callerFunction := frame.Function

	logMsg := LogAttributes{timestamp: timestamp, user: username, file: file, line: line, callerFunction: callerFunction, msg: msgString}

	if sl.db != nil {
		_, err = sl.Save(logMsg)

		if err != nil {
			log.Fatal("cannot save in database")
		}
	}

	fmt.Printf("%s - %s - %s:%s %s \"%s\"\n", timestamp, username, file, line, callerFunction, msg)
}

func (sl SingleLogger) Save(logMsg LogAttributes) (int, error) {
	stmt, err := sl.db.Prepare(SaveLog)

	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(&logMsg.timestamp, &logMsg.user, &logMsg.file, &logMsg.line, &logMsg.callerFunction, &logMsg.msg)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(id), nil
}
