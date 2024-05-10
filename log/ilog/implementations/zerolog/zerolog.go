package zerolog

import (
	"time"

	"github.com/rs/zerolog"

	"github.com/kunitsucom/ilog.go"
)

type implLogger struct {
	level         ilog.Level
	zerologLogger *zerolog.Logger
}

func New(level ilog.Level, l zerolog.Logger) ilog.Logger {
	return &implLogger{
		level:         level,
		zerologLogger: &l,
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
	logger := l.zerologLogger.With().Caller().CallerWithSkipFrameCount(skip).Logger()
	copied.zerologLogger = &logger
	return copied
}

func (l *implLogger) Copy() ilog.Logger {
	return l.copy()
}

func (l *implLogger) copy() *implLogger {
	copied := *l
	copiedZerologLogger := *l.zerologLogger
	copied.zerologLogger = &copiedZerologLogger
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
	}
}

//nolint:errname
type implLogEntry struct {
	logger *implLogger
	zCtxs  []func(e zerolog.Context) zerolog.Context
}

func (*implLogEntry) Error() string {
	return ilog.ErrLogEntryIsNotWritten.Error()
}

func (e *implLogEntry) Any(key string, value interface{}) ilog.LogEntry {
	e.zCtxs = append(e.zCtxs, func(e zerolog.Context) zerolog.Context {
		return e.Interface(key, value)
	})
	return e
}

func (e *implLogEntry) Bool(key string, value bool) ilog.LogEntry {
	e.zCtxs = append(e.zCtxs, func(e zerolog.Context) zerolog.Context {
		return e.Bool(key, value)
	})
	return e
}

func (e *implLogEntry) Bytes(key string, value []byte) ilog.LogEntry {
	e.zCtxs = append(e.zCtxs, func(e zerolog.Context) zerolog.Context {
		return e.Bytes(key, value)
	})
	return e
}

func (e *implLogEntry) Duration(key string, value time.Duration) ilog.LogEntry {
	e.zCtxs = append(e.zCtxs, func(e zerolog.Context) zerolog.Context {
		return e.Dur(key, value)
	})
	return e
}

func (e *implLogEntry) Err(err error) ilog.LogEntry {
	e.zCtxs = append(e.zCtxs, func(e zerolog.Context) zerolog.Context {
		return e.Err(err)
	})
	return e
}

func (e *implLogEntry) ErrWithKey(key string, err error) ilog.LogEntry {
	e.zCtxs = append(e.zCtxs, func(e zerolog.Context) zerolog.Context {
		return e.AnErr(key, err)
	})
	return e
}

func (e *implLogEntry) Float32(key string, value float32) ilog.LogEntry {
	e.zCtxs = append(e.zCtxs, func(e zerolog.Context) zerolog.Context {
		return e.Float32(key, value)
	})
	return e
}

func (e *implLogEntry) Float64(key string, value float64) ilog.LogEntry {
	e.zCtxs = append(e.zCtxs, func(e zerolog.Context) zerolog.Context {
		return e.Float64(key, value)
	})
	return e
}

func (e *implLogEntry) Int(key string, value int) ilog.LogEntry {
	e.zCtxs = append(e.zCtxs, func(e zerolog.Context) zerolog.Context {
		return e.Int(key, value)
	})
	return e
}

func (e *implLogEntry) Int32(key string, value int32) ilog.LogEntry {
	e.zCtxs = append(e.zCtxs, func(e zerolog.Context) zerolog.Context {
		return e.Int32(key, value)
	})
	return e
}

func (e *implLogEntry) Int64(key string, value int64) ilog.LogEntry {
	e.zCtxs = append(e.zCtxs, func(e zerolog.Context) zerolog.Context {
		return e.Int64(key, value)
	})
	return e
}

func (e *implLogEntry) String(key, value string) ilog.LogEntry {
	e.zCtxs = append(e.zCtxs, func(e zerolog.Context) zerolog.Context {
		return e.Str(key, value)
	})
	return e
}

func (e *implLogEntry) Time(key string, value time.Time) ilog.LogEntry {
	e.zCtxs = append(e.zCtxs, func(e zerolog.Context) zerolog.Context {
		return e.Time(key, value)
	})
	return e
}

func (e *implLogEntry) Uint(key string, value uint) ilog.LogEntry {
	e.zCtxs = append(e.zCtxs, func(e zerolog.Context) zerolog.Context {
		return e.Uint(key, value)
	})
	return e
}

func (e *implLogEntry) Uint32(key string, value uint32) ilog.LogEntry {
	e.zCtxs = append(e.zCtxs, func(e zerolog.Context) zerolog.Context {
		return e.Uint32(key, value)
	})
	return e
}

func (e *implLogEntry) Uint64(key string, value uint64) ilog.LogEntry {
	e.zCtxs = append(e.zCtxs, func(e zerolog.Context) zerolog.Context {
		return e.Uint64(key, value)
	})
	return e
}

func (e *implLogEntry) Logger() ilog.Logger {
	copied := e.logger.copy()
	c := copied.zerologLogger.With()
	for _, event := range e.zCtxs {
		c = event(c)
	}
	logger := c.Logger()
	copied.zerologLogger = &logger

	return copied
}

func (e *implLogEntry) Write(p []byte) (n int, err error) {
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

//nolint:cyclop,funlen
func (e *implLogEntry) logf(level ilog.Level, format string, args ...interface{}) {
	if level < e.logger.level {
		return
	}
	switch level { //nolint:exhaustive
	case ilog.InfoLevel:
		c := e.logger.zerologLogger.With()
		for i := range e.zCtxs {
			c = e.zCtxs[i](c)
		}
		zl := c.Logger()
		e := zl.Info()
		if len(args) > 0 {
			e.Msgf(format, args...)
			return
		}
		e.Msg(format)
		return
	case ilog.WarnLevel:
		c := e.logger.zerologLogger.With()
		for i := range e.zCtxs {
			c = e.zCtxs[i](c)
		}
		zl := c.Logger()
		e := zl.Warn()
		if len(args) > 0 {
			e.Msgf(format, args...)
			return
		}
		e.Msg(format)
		return
	case ilog.ErrorLevel:
		c := e.logger.zerologLogger.With()
		for i := range e.zCtxs {
			c = e.zCtxs[i](c)
		}
		zl := c.Logger()
		e := zl.Error()
		if len(args) > 0 {
			e.Msgf(format, args...)
			return
		}
		e.Msg(format)
		return
	default:
		c := e.logger.zerologLogger.With()
		for i := range e.zCtxs {
			c = e.zCtxs[i](c)
		}
		zl := c.Logger()
		e := zl.Debug()
		if len(args) > 0 {
			e.Msgf(format, args...)
			return
		}
		e.Msg(format)
		return
	}
}
