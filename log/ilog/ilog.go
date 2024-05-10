// ilog is a simple logging interface library package for Go.
// By defining only the logger interface `ilog.Logger`, users can easily swap out the underlying logging implementation without changing their application code.
package ilog

import (
	"errors"
	"time"
)

// ErrLogEntryIsNotWritten is the error returned when a log entry is not written.
var ErrLogEntryIsNotWritten = errors.New("ilog: log entry not written")

// Level is the type for logging level.
type Level int8

const (
	DebugLevel Level = -8
	InfoLevel  Level = 0
	WarnLevel  Level = 8
	ErrorLevel Level = 16
)

// Logger is the interface that has the basic logging methods.
type Logger interface {
	// Level returns the current logging level of the logger.
	Level() (currentLoggerLevel Level)
	// SetLevel sets the logging level of the logger.
	SetLevel(level Level) (copied Logger)
	// AddCallerSkip adds the number of stack frames to skip to the logger.
	AddCallerSkip(skip int) (copied Logger)
	// Copy returns a copy of the logger.
	Copy() (copied Logger)

	// common is the interface that has the common logging methods for both ilog.Logger and ilog.LogEntry.
	common
}

// common is the interface that has the common logging methods for both ilog.Logger and ilog.LogEntry.
type common interface {
	Any(key string, value interface{}) (entry LogEntry)
	Bool(key string, value bool) (entry LogEntry)
	Bytes(key string, value []byte) (entry LogEntry)
	Duration(key string, value time.Duration) (entry LogEntry)
	Err(err error) (entry LogEntry)
	ErrWithKey(key string, err error) (entry LogEntry)
	Float32(key string, value float32) (entry LogEntry)
	Float64(key string, value float64) (entry LogEntry)
	Int(key string, value int) (entry LogEntry)
	Int32(key string, value int32) (entry LogEntry)
	Int64(key string, value int64) (entry LogEntry)
	String(key, value string) (entry LogEntry)
	Time(key string, value time.Time) (entry LogEntry)
	Uint(key string, value uint) (entry LogEntry)
	Uint32(key string, value uint32) (entry LogEntry)
	Uint64(key string, value uint64) (entry LogEntry)

	// Debugf logs a message at debug level.
	// If the argument is one, it is treated 1st argument as a simple string.
	// If the argument is more than one, it is treated 1st argument as a format string.
	Debugf(format string, args ...interface{})
	// Infof logs a message at info level.
	// If the argument is one, it is treated 1st argument as a simple string.
	// If the argument is more than one, it is treated 1st argument as a format string.
	Infof(format string, args ...interface{})
	// Warnf logs a message at warn level.
	// If the argument is one, it is treated 1st argument as a simple string.
	// If the argument is more than one, it is treated 1st argument as a format string.
	Warnf(format string, args ...interface{})
	// Errorf logs a message at error level.
	// If the argument is one, it is treated 1st argument as a simple string.
	// If the argument is more than one, it is treated 1st argument as a format string.
	Errorf(format string, args ...interface{})
	// Logf logs a message at the specified level.
	// If the argument is one, it is treated 1st argument as a simple string.
	// If the argument is more than one, it is treated 1st argument as a format string.
	Logf(level Level, format string, args ...interface{})

	// Write implements io.Writer interface.
	// It outputs logs at the logger's log level.
	Write(p []byte) (n int, err error)
}

// LogEntry is the interface that has the logging methods for a single log entry.
type LogEntry interface {
	// common is the interface that has the common logging methods for both ilog.Logger and ilog.LogEntry.
	common

	// Logger returns a new logger with the same fields of the log entry.
	Logger() (copied Logger)

	// error: for considering undispatched LogEntry as error so that they can be detected by Go static analysis.
	error
}
