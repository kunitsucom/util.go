package zap

import (
	"time"

	"go.uber.org/zap"

	"github.com/kunitsucom/ilog.go"
)

type implLogger struct {
	level     ilog.Level
	zapLogger *zap.Logger
}

func New(level ilog.Level, logger *zap.Logger) ilog.Logger {
	return &implLogger{
		level:     level,
		zapLogger: logger.WithOptions(zap.AddCallerSkip(2)),
	}
}

func (l *implLogger) Level() ilog.Level {
	return l.level
}

func (l *implLogger) SetLevel(level ilog.Level) ilog.Logger {
	copied := l.copy()
	copied.level = level
	return copied
}

func (l *implLogger) AddCallerSkip(skip int) ilog.Logger {
	copied := l.copy()
	copied.zapLogger = l.zapLogger.WithOptions(zap.AddCallerSkip(skip))
	return copied
}

func (l *implLogger) Copy() ilog.Logger {
	return l.copy()
}

func (l *implLogger) copy() *implLogger {
	copied := *l
	copied.zapLogger = l.zapLogger.WithOptions() // NOTE: call (*zap.Logger).clone() internally
	return &copied
}

func (l *implLogger) Any(key string, value interface{}) ilog.LogEntry {
	return l.new().Any(key, value)
}

func (l *implLogger) Bool(key string, value bool) ilog.LogEntry {
	return l.new().Bool(key, value)
}

func (l *implLogger) Bytes(key string, value []byte) ilog.LogEntry {
	return l.new().Bytes(key, value)
}

func (l *implLogger) Duration(key string, value time.Duration) ilog.LogEntry {
	return l.new().Duration(key, value)
}

func (l *implLogger) Err(err error) ilog.LogEntry {
	return l.new().Err(err)
}

func (l *implLogger) ErrWithKey(key string, err error) ilog.LogEntry {
	return l.new().ErrWithKey(key, err)
}

func (l *implLogger) Float32(key string, value float32) ilog.LogEntry {
	return l.new().Float32(key, value)
}

func (l *implLogger) Float64(key string, value float64) ilog.LogEntry {
	return l.new().Float64(key, value)
}

func (l *implLogger) Int(key string, value int) ilog.LogEntry {
	return l.new().Int(key, value)
}

func (l *implLogger) Int32(key string, value int32) ilog.LogEntry {
	return l.new().Int32(key, value)
}

func (l *implLogger) Int64(key string, value int64) ilog.LogEntry {
	return l.new().Int64(key, value)
}

func (l *implLogger) String(key, value string) ilog.LogEntry {
	return l.new().String(key, value)
}

func (l *implLogger) Time(key string, value time.Time) ilog.LogEntry {
	return l.new().Time(key, value)
}

func (l *implLogger) Uint(key string, value uint) ilog.LogEntry {
	return l.new().Uint(key, value)
}

func (l *implLogger) Uint32(key string, value uint32) ilog.LogEntry {
	return l.new().Uint32(key, value)
}

func (l *implLogger) Uint64(key string, value uint64) ilog.LogEntry {
	return l.new().Uint64(key, value)
}

func (l *implLogger) Debugf(format string, args ...interface{}) {
	l.new().logf(ilog.DebugLevel, format, args...)
}

func (l *implLogger) Infof(format string, args ...interface{}) {
	l.new().logf(ilog.InfoLevel, format, args...)
}

func (l *implLogger) Warnf(format string, args ...interface{}) {
	l.new().logf(ilog.WarnLevel, format, args...)
}

func (l *implLogger) Errorf(format string, args ...interface{}) {
	l.new().logf(ilog.ErrorLevel, format, args...)
}

func (l *implLogger) Logf(level ilog.Level, format string, args ...interface{}) {
	l.new().logf(level, format, args...)
}

func (l *implLogger) Write(p []byte) (int, error) {
	l.new().logf(l.level, string(p))
	return len(p), nil
}

func (l *implLogger) new() *implLogEntry {
	return &implLogEntry{
		logger: l,
		fields: make([]zap.Field, 0),
	}
}

//nolint:errname
type implLogEntry struct {
	logger *implLogger
	fields []zap.Field
}

func (*implLogEntry) Error() string {
	return ilog.ErrLogEntryIsNotWritten.Error()
}

func (e *implLogEntry) Any(key string, value interface{}) ilog.LogEntry {
	e.fields = append(e.fields, zap.Any(key, value))
	return e
}

func (e *implLogEntry) Bool(key string, value bool) ilog.LogEntry {
	e.fields = append(e.fields, zap.Bool(key, value))
	return e
}

func (e *implLogEntry) Bytes(key string, value []byte) ilog.LogEntry {
	e.fields = append(e.fields, zap.ByteString(key, value))
	return e
}

func (e *implLogEntry) Duration(key string, value time.Duration) ilog.LogEntry {
	e.fields = append(e.fields, zap.Duration(key, value))
	return e
}

func (e *implLogEntry) Err(err error) ilog.LogEntry {
	e.fields = append(e.fields, zap.Error(err))
	return e
}

func (e *implLogEntry) ErrWithKey(key string, err error) ilog.LogEntry {
	e.fields = append(e.fields, zap.NamedError(key, err))
	return e
}

func (e *implLogEntry) Float32(key string, value float32) ilog.LogEntry {
	e.fields = append(e.fields, zap.Float32(key, value))
	return e
}

func (e *implLogEntry) Float64(key string, value float64) ilog.LogEntry {
	e.fields = append(e.fields, zap.Float64(key, value))
	return e
}

func (e *implLogEntry) Int(key string, value int) ilog.LogEntry {
	e.fields = append(e.fields, zap.Int(key, value))
	return e
}

func (e *implLogEntry) Int32(key string, value int32) ilog.LogEntry {
	e.fields = append(e.fields, zap.Int32(key, value))
	return e
}

func (e *implLogEntry) Int64(key string, value int64) ilog.LogEntry {
	e.fields = append(e.fields, zap.Int64(key, value))
	return e
}

func (e *implLogEntry) String(key, value string) ilog.LogEntry {
	e.fields = append(e.fields, zap.String(key, value))
	return e
}

func (e *implLogEntry) Time(key string, value time.Time) ilog.LogEntry {
	e.fields = append(e.fields, zap.Time(key, value))
	return e
}

func (e *implLogEntry) Uint(key string, value uint) ilog.LogEntry {
	e.fields = append(e.fields, zap.Uint(key, value))
	return e
}

func (e *implLogEntry) Uint32(key string, value uint32) ilog.LogEntry {
	e.fields = append(e.fields, zap.Uint32(key, value))
	return e
}

func (e *implLogEntry) Uint64(key string, value uint64) ilog.LogEntry {
	e.fields = append(e.fields, zap.Uint64(key, value))
	return e
}

func (e *implLogEntry) Logger() ilog.Logger {
	copied := e.logger.copy()
	copied.zapLogger = copied.zapLogger.With(e.fields...)
	return copied
}

func (e *implLogEntry) Write(p []byte) (int, error) {
	e.logf(e.logger.level, string(p))
	return len(p), nil
}

func (e *implLogEntry) Debugf(format string, args ...interface{}) {
	e.logf(ilog.DebugLevel, format, args...)
}

func (e *implLogEntry) Infof(format string, args ...interface{}) {
	e.logf(ilog.InfoLevel, format, args...)
}

func (e *implLogEntry) Warnf(format string, args ...interface{}) {
	e.logf(ilog.WarnLevel, format, args...)
}

func (e *implLogEntry) Errorf(format string, args ...interface{}) {
	e.logf(ilog.ErrorLevel, format, args...)
}

func (e *implLogEntry) Logf(level ilog.Level, format string, args ...interface{}) {
	e.logf(level, format, args...)
}

func (e *implLogEntry) logf(level ilog.Level, format string, args ...interface{}) {
	if level < e.logger.level {
		return
	}
	defer func() {
		e.fields = make([]zap.Field, 0)
	}()
	switch level { //nolint:exhaustive
	case ilog.InfoLevel:
		if len(args) > 0 {
			e.logger.zapLogger.With(e.fields...).Sugar().Infof(format, args...)
			return
		}
		e.logger.zapLogger.Info(format, e.fields...)
		return
	case ilog.WarnLevel:
		if len(args) > 0 {
			e.logger.zapLogger.With(e.fields...).Sugar().Warnf(format, args...)
			return
		}
		e.logger.zapLogger.Warn(format, e.fields...)
		return
	case ilog.ErrorLevel:
		if len(args) > 0 {
			e.logger.zapLogger.With(e.fields...).Sugar().Errorf(format, args...)
			return
		}
		e.logger.zapLogger.Error(format, e.fields...)
		return
	default:
		if len(args) > 0 {
			e.logger.zapLogger.With(e.fields...).Sugar().Debugf(format, args...)
			return
		}
		e.logger.zapLogger.Debug(format, e.fields...)
		return
	}
}
