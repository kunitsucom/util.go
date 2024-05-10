package ilog

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

//nolint:gochecknoglobals
var defaultLevels = map[Level]string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
}

func copyLevels(levels map[Level]string) map[Level]string {
	copied := make(map[Level]string, len(levels))
	for level, value := range levels {
		copied[level] = value
	}

	return copied
}

type implLoggerConfig struct {
	levelKey        string
	level           Level
	levels          map[Level]string
	timestampKey    string
	timestampFormat string
	timestampZone   *time.Location
	callerKey       string
	callerSkip      int
	useLongCaller   bool
	messageKey      string
	separator       string
	writer          io.Writer
}

type implLogger struct {
	config implLoggerConfig
	fields []byte
}

type syncWriter interface {
	io.Writer
	Lock()
	Unlock()
}

type _syncWriter struct {
	mu sync.Mutex
	w  io.Writer
}

func (w *_syncWriter) Write(p []byte) (int, error) {
	w.Lock()
	defer w.Unlock()
	return w.w.Write(p) //nolint:wrapcheck
}
func (w *_syncWriter) Lock()   { w.mu.Lock() }
func (w *_syncWriter) Unlock() { w.mu.Unlock() }

func NewSyncWriter(w io.Writer) io.Writer {
	return &_syncWriter{w: w}
}

// NewBuilder returns a new Builder of ilog.Logger with the specified level and writer.
func NewBuilder(level Level, w io.Writer) implLoggerConfig { //nolint:revive
	const defaultCallerSkip = 4
	return implLoggerConfig{
		levelKey:        "severity",
		level:           level,
		levels:          copyLevels(defaultLevels),
		timestampKey:    "timestamp",
		timestampFormat: time.RFC3339Nano,
		timestampZone:   time.Local, //nolint:gosmopolitan
		callerKey:       "caller",
		callerSkip:      defaultCallerSkip,
		useLongCaller:   false,
		messageKey:      "message",
		separator:       "\n",
		writer:          w,
	}
}

// SetLevelKey sets the key of the level field of the logger.
// If empty, the level field is not output.
// Default is "severity".
func (c implLoggerConfig) SetLevelKey(key string) implLoggerConfig { //nolint:revive
	c.levelKey = key
	return c
}

// SetLevel sets the level of the logger.
func (c implLoggerConfig) SetLevels(levels map[Level]string) implLoggerConfig { //nolint:revive
	c.levels = levels
	return c
}

// SetTimestampKey sets the key of the timestamp field of the logger.
// If empty, the timestamp field is not output.
// Default is "timestamp".
func (c implLoggerConfig) SetTimestampKey(key string) implLoggerConfig { //nolint:revive
	c.timestampKey = key
	return c
}

// SetTimestampFormat sets the format of the timestamp field of the logger.
// Default is time.RFC3339Nano.
func (c implLoggerConfig) SetTimestampFormat(format string) implLoggerConfig { //nolint:revive
	c.timestampFormat = format
	return c
}

// SetTimestampZone sets the time zone of the timestamp field of the logger.
// Default is time.Local.
func (c implLoggerConfig) SetTimestampZone(zone *time.Location) implLoggerConfig { //nolint:revive
	c.timestampZone = zone
	return c
}

// SetCallerKey sets the key of the caller field of the logger.
// If empty, the caller field is not output.
// Default is "caller".
func (c implLoggerConfig) SetCallerKey(key string) implLoggerConfig { //nolint:revive
	c.callerKey = key
	return c
}

// UseLongCaller sets whether to use long caller of the logger.
// If true, the long caller is used.
// Default caller is short caller.
func (c implLoggerConfig) UseLongCaller(useLongCaller bool) implLoggerConfig { //nolint:revive
	c.useLongCaller = useLongCaller
	return c
}

// SetMessageKey sets the key of the message field of the logger.
// If empty, the message field is not output.
// Default is "message".
func (c implLoggerConfig) SetMessageKey(key string) implLoggerConfig { //nolint:revive
	c.messageKey = key
	return c
}

// SetSeparator sets the log entry separator.
// Default is "\n".
func (c implLoggerConfig) SetSeparator(separator string) implLoggerConfig { //nolint:revive
	c.separator = separator
	return c
}

// UseSyncWriter sets whether to use sync writer of the logger.
func (c implLoggerConfig) UseSyncWriter() implLoggerConfig { //nolint:revive
	switch v := c.writer.(type) {
	case syncWriter:
		c.writer = v
	default:
		c.writer = NewSyncWriter(c.writer)
	}

	return c
}

// Build returns a new ilog.Logger with the specified configuration.
func (c implLoggerConfig) Build() Logger { //nolint:ireturn
	const fieldsCap = 1024
	return &implLogger{
		config: c,
		fields: make([]byte, 0, fieldsCap),
	}
}

func (l *implLogger) Level() Level {
	return l.config.level
}

func (l *implLogger) SetLevel(level Level) Logger { //nolint:ireturn
	copied := l.copy()
	copied.config.level = level
	return copied
}

func (l *implLogger) AddCallerSkip(skip int) Logger { //nolint:ireturn
	copied := l.copy()
	copied.config.callerSkip += skip
	return copied
}

func (l *implLogger) Copy() Logger { //nolint:ireturn
	return l.copy()
}

func (l *implLogger) copy() *implLogger {
	copied := *l
	copied.fields = make([]byte, len(l.fields))
	copy(copied.fields, l.fields)
	return &copied
}

func (l *implLogger) Any(key string, value interface{}) LogEntry { //nolint:ireturn
	return l.new().Any(key, value)
}

func (l *implLogger) Bool(key string, value bool) LogEntry { //nolint:ireturn
	return l.new().Bool(key, value)
}

func (l *implLogger) Bytes(key string, value []byte) LogEntry { //nolint:ireturn
	return l.new().Bytes(key, value)
}

func (l *implLogger) Duration(key string, value time.Duration) LogEntry { //nolint:ireturn
	return l.new().Duration(key, value)
}

func (l *implLogger) Err(err error) LogEntry { //nolint:ireturn
	return l.new().Err(err)
}

func (l *implLogger) ErrWithKey(key string, err error) LogEntry { //nolint:ireturn
	return l.new().ErrWithKey(key, err)
}

func (l *implLogger) Float32(key string, value float32) LogEntry { //nolint:ireturn
	return l.new().Float32(key, value)
}

func (l *implLogger) Float64(key string, value float64) LogEntry { //nolint:ireturn
	return l.new().Float64(key, value)
}

func (l *implLogger) Int(key string, value int) LogEntry { //nolint:ireturn
	return l.new().Int(key, value)
}

func (l *implLogger) Int32(key string, value int32) LogEntry { //nolint:ireturn
	return l.new().Int32(key, value)
}

func (l *implLogger) Int64(key string, value int64) LogEntry { //nolint:ireturn
	return l.new().Int64(key, value)
}

func (l *implLogger) String(key, value string) LogEntry { //nolint:ireturn
	return l.new().String(key, value)
}

func (l *implLogger) Time(key string, value time.Time) LogEntry { //nolint:ireturn
	return l.new().Time(key, value)
}

func (l *implLogger) Uint(key string, value uint) LogEntry { //nolint:ireturn
	return l.new().Uint(key, value)
}

func (l *implLogger) Uint32(key string, value uint32) LogEntry { //nolint:ireturn
	return l.new().Uint32(key, value)
}

func (l *implLogger) Uint64(key string, value uint64) LogEntry { //nolint:ireturn
	return l.new().Uint64(key, value)
}

func (l *implLogger) Debugf(format string, args ...interface{}) {
	_ = l.new().logf(DebugLevel, format, args...)
}

func (l *implLogger) Infof(format string, args ...interface{}) {
	_ = l.new().logf(InfoLevel, format, args...)
}

func (l *implLogger) Warnf(format string, args ...interface{}) {
	_ = l.new().logf(WarnLevel, format, args...)
}

func (l *implLogger) Errorf(format string, args ...interface{}) {
	_ = l.new().logf(ErrorLevel, format, args...)
}

func (l *implLogger) Logf(level Level, format string, args ...interface{}) {
	_ = l.new().logf(level, format, args...)
}

func (l *implLogger) Write(p []byte) (int, error) {
	if err := l.new().logf(l.config.level, string(p)); err != nil {
		return 0, fmt.Errorf("w.logf: %w", err)
	}
	return len(p), nil
}

func (l *implLogger) new() *implLogEntry {
	buffer, put := getBytesBuffer()
	return &implLogEntry{
		logger:      l,
		bytesBuffer: buffer,
		put:         put,
	}
}

//nolint:errname
type implLogEntry struct {
	logger      *implLogger
	bytesBuffer *bytesBuffer
	put         func()
}

func (*implLogEntry) Error() string {
	return ErrLogEntryIsNotWritten.Error()
}

func (e *implLogEntry) null(key string) LogEntry { //nolint:ireturn
	e.bytesBuffer.bytes = appendKey(e.bytesBuffer.bytes, key)
	e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, null...)
	e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, ',')
	return e
}

//nolint:cyclop,funlen
func (e *implLogEntry) Any(key string, value interface{}) (le LogEntry) { //nolint:ireturn
	switch v := value.(type) {
	case bool:
		return e.Bool(key, v)
	case *bool:
		if v != nil {
			return e.Bool(key, *v)
		}
		return e.null(key)
	case byte:
		return e.String(key, string(v))
	case []byte:
		return e.Bytes(key, v)
	case time.Duration:
		return e.Duration(key, v)
	case error:
		return e.ErrWithKey(key, v)
	case float32:
		return e.Float32(key, v)
	case float64:
		return e.Float64(key, v)
	case int:
		return e.Int(key, v)
	case int8:
		return e.Int(key, int(v))
	case int16:
		return e.Int(key, int(v))
	case int32:
		return e.Int32(key, v)
	case int64:
		return e.Int64(key, v)
	case string:
		return e.String(key, v)
	case time.Time:
		return e.Time(key, v)
	case uint:
		return e.Uint(key, v)
	// NOTE: uint8 == byte
	// case uint8:
	// 	return w.Uint(key, uint(v))
	case uint16:
		return e.Uint(key, uint(v))
	case uint32:
		return e.Uint32(key, v)
	case uint64:
		return e.Uint64(key, v)
	case json.Marshaler:
		defer func() {
			if p := recover(); p != nil {
				le = e.null(key)
			}
		}()
		// NOTE: Even if v is nil, it is not judged as nil because it has type information. Calling v.MarshalJSON() causes panic.
		b, err := v.MarshalJSON()
		if err != nil {
			return e.ErrWithKey(key, fmt.Errorf("json.Marshaler: v.MarshalJSON: %w", err))
		}
		e.bytesBuffer.bytes = appendKey(e.bytesBuffer.bytes, key)
		e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, b...)
		e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, ',')
		return e
	case fmt.Formatter:
		return e.String(key, fmt.Sprintf("%+v", v))
	case fmt.Stringer:
		defer func() {
			if p := recover(); p != nil {
				le = e.null(key)
			}
		}()
		// NOTE: Even if v is nil, it is not judged as nil because it has type information. Calling v.String() causes panic.
		return e.String(key, v.String())
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return e.String(key, fmt.Sprintf("%v", v))
		}
		e.bytesBuffer.bytes = appendKey(e.bytesBuffer.bytes, key)
		e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, b...)
		e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, ',')
		return e
	}
}

func (e *implLogEntry) Bool(key string, value bool) LogEntry { //nolint:ireturn
	e.bytesBuffer.bytes = appendKey(e.bytesBuffer.bytes, key)
	e.bytesBuffer.bytes = strconv.AppendBool(e.bytesBuffer.bytes, value)
	e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, ',')
	return e
}

func (e *implLogEntry) Bytes(key string, value []byte) LogEntry { //nolint:ireturn
	e.bytesBuffer.bytes = appendKey(e.bytesBuffer.bytes, key)
	e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, '"')
	e.bytesBuffer.bytes = appendJSONEscapedString(e.bytesBuffer.bytes, string(value))
	e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, '"', ',')
	return e
}

func (e *implLogEntry) Duration(key string, value time.Duration) LogEntry { //nolint:ireturn
	e.bytesBuffer.bytes = appendKey(e.bytesBuffer.bytes, key)
	e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, '"')
	e.bytesBuffer.bytes = appendJSONEscapedString(e.bytesBuffer.bytes, value.String())
	e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, '"', ',')
	return e
}

func (e *implLogEntry) Err(err error) LogEntry { //nolint:ireturn
	return e.ErrWithKey("error", err)
}

func (e *implLogEntry) ErrWithKey(key string, err error) (le LogEntry) { //nolint:ireturn
	defer func() {
		if p := recover(); p != nil {
			le = e.null(key)
		}
	}()

	// NOTE: Even if err is your unique error type and nil, it is not judged as nil because it has type information. Calling err.Error() causes panic.
	var v string
	formatter, ok := err.(fmt.Formatter) //nolint:errorlint
	if ok && formatter != nil {
		v = fmt.Sprintf("%+v", formatter)
	} else {
		v = err.Error()
	}
	e.bytesBuffer.bytes = appendKey(e.bytesBuffer.bytes, key)
	e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, '"')
	e.bytesBuffer.bytes = appendJSONEscapedString(e.bytesBuffer.bytes, v)
	e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, '"')
	e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, ',')
	return e
}

func (e *implLogEntry) Float32(key string, value float32) LogEntry { //nolint:ireturn
	e.bytesBuffer.bytes = appendKey(e.bytesBuffer.bytes, key)
	const bitSize = 32
	e.bytesBuffer.bytes = appendFloatFieldValue(e.bytesBuffer.bytes, float64(value), bitSize)
	e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, ',')
	return e
}

func (e *implLogEntry) Float64(key string, value float64) LogEntry { //nolint:ireturn
	e.bytesBuffer.bytes = appendKey(e.bytesBuffer.bytes, key)
	const bitSize = 64
	e.bytesBuffer.bytes = appendFloatFieldValue(e.bytesBuffer.bytes, value, bitSize)
	e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, ',')
	return e
}

func (e *implLogEntry) Int(key string, value int) LogEntry { //nolint:ireturn
	e.bytesBuffer.bytes = appendKey(e.bytesBuffer.bytes, key)
	const base = 10
	e.bytesBuffer.bytes = strconv.AppendInt(e.bytesBuffer.bytes, int64(value), base)
	e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, ',')
	return e
}

func (e *implLogEntry) Int32(key string, value int32) LogEntry { //nolint:ireturn
	e.bytesBuffer.bytes = appendKey(e.bytesBuffer.bytes, key)
	const base = 10
	e.bytesBuffer.bytes = strconv.AppendInt(e.bytesBuffer.bytes, int64(value), base)
	e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, ',')
	return e
}

func (e *implLogEntry) Int64(key string, value int64) LogEntry { //nolint:ireturn
	e.bytesBuffer.bytes = appendKey(e.bytesBuffer.bytes, key)
	const base = 10
	e.bytesBuffer.bytes = strconv.AppendInt(e.bytesBuffer.bytes, value, base)
	e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, ',')
	return e
}

func (e *implLogEntry) String(key string, value string) LogEntry { //nolint:ireturn
	e.bytesBuffer.bytes = appendKey(e.bytesBuffer.bytes, key)
	e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, '"')
	e.bytesBuffer.bytes = appendJSONEscapedString(e.bytesBuffer.bytes, value)
	e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, '"', ',')
	return e
}

func (e *implLogEntry) Time(key string, value time.Time) LogEntry { //nolint:ireturn
	e.bytesBuffer.bytes = appendKey(e.bytesBuffer.bytes, key)
	e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, '"')
	e.bytesBuffer.bytes = appendJSONEscapedString(e.bytesBuffer.bytes, value.Format(e.logger.config.timestampFormat))
	e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, '"', ',')
	return e
}

func (e *implLogEntry) Uint(key string, value uint) LogEntry { //nolint:ireturn
	e.bytesBuffer.bytes = appendKey(e.bytesBuffer.bytes, key)
	const base = 10
	e.bytesBuffer.bytes = strconv.AppendUint(e.bytesBuffer.bytes, uint64(value), base)
	e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, ',')
	return e
}

func (e *implLogEntry) Uint32(key string, value uint32) LogEntry { //nolint:ireturn
	e.bytesBuffer.bytes = appendKey(e.bytesBuffer.bytes, key)
	const base = 10
	e.bytesBuffer.bytes = strconv.AppendUint(e.bytesBuffer.bytes, uint64(value), base)
	e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, ',')
	return e
}

func (e *implLogEntry) Uint64(key string, value uint64) LogEntry { //nolint:ireturn
	e.bytesBuffer.bytes = appendKey(e.bytesBuffer.bytes, key)
	const base = 10
	e.bytesBuffer.bytes = strconv.AppendUint(e.bytesBuffer.bytes, value, base)
	e.bytesBuffer.bytes = append(e.bytesBuffer.bytes, ',')
	return e
}

func (e *implLogEntry) Logger() Logger { //nolint:ireturn
	copied := e.logger.copy()
	copied.fields = append(copied.fields, e.bytesBuffer.bytes...)
	return copied
}

func (e *implLogEntry) Debugf(format string, args ...interface{}) {
	_ = e.logf(DebugLevel, format, args...)
}

func (e *implLogEntry) Infof(format string, args ...interface{}) {
	_ = e.logf(InfoLevel, format, args...)
}

func (e *implLogEntry) Warnf(format string, args ...interface{}) {
	_ = e.logf(WarnLevel, format, args...)
}

func (e *implLogEntry) Errorf(format string, args ...interface{}) {
	_ = e.logf(ErrorLevel, format, args...)
}

func (e *implLogEntry) Logf(level Level, format string, args ...interface{}) {
	_ = e.logf(level, format, args...)
}

func (e *implLogEntry) Write(p []byte) (int, error) {
	if err := e.logf(e.logger.config.level, string(p)); err != nil {
		return 0, fmt.Errorf("w.logf: %w", err)
	}
	return len(p), nil
}

//nolint:cyclop
func (e *implLogEntry) logf(level Level, format string, args ...interface{}) error {
	defer e.put()
	if level < e.logger.config.level {
		return nil
	}

	b, put := getBytesBuffer()
	defer put()

	b.bytes = append(b.bytes, '{')

	if len(e.logger.config.levelKey) > 0 {
		b.bytes = appendKey(b.bytes, e.logger.config.levelKey)
		b.bytes = appendLevelField(b.bytes, e.logger.config.levels, level)
		b.bytes = append(b.bytes, ',')
	}
	if len(e.logger.config.timestampKey) > 0 {
		b.bytes = appendKey(b.bytes, e.logger.config.timestampKey)
		b.bytes = append(b.bytes, '"')
		b.bytes = appendJSONEscapedString(b.bytes, time.Now().In(e.logger.config.timestampZone).Format(e.logger.config.timestampFormat))
		b.bytes = append(b.bytes, '"', ',')
	}
	if len(e.logger.config.callerKey) > 0 {
		b.bytes = appendKey(b.bytes, e.logger.config.callerKey)
		b.bytes = append(b.bytes, '"')
		b.bytes = appendCaller(b.bytes, e.logger.config.callerSkip, e.logger.config.useLongCaller)
		b.bytes = append(b.bytes, '"', ',')
	}
	if len(e.logger.config.messageKey) > 0 {
		b.bytes = appendKey(b.bytes, e.logger.config.messageKey)
		b.bytes = append(b.bytes, '"')
		if len(args) > 0 {
			b.bytes = appendJSONEscapedString(b.bytes, fmt.Sprintf(format, args...))
		} else {
			b.bytes = appendJSONEscapedString(b.bytes, format)
		}
		b.bytes = append(b.bytes, '"', ',')
	}

	if len(e.logger.fields) > 0 {
		b.bytes = append(b.bytes, e.logger.fields...)
	}

	if len(e.bytesBuffer.bytes) > 0 {
		b.bytes = append(b.bytes, e.bytesBuffer.bytes...)
	}

	if b.bytes[len(b.bytes)-1] == ',' {
		b.bytes[len(b.bytes)-1] = '}'
	} else {
		b.bytes = append(b.bytes, '}')
	}

	if _, err := e.logger.config.writer.Write(append(b.bytes, e.logger.config.separator...)); err != nil {
		err = fmt.Errorf("w.logger.writer.Write: p=%s: %w", b.bytes, err)
		defer Global().Errorf(err.Error())
		return err
	}

	return nil
}

type (
	bytesBuffer struct {
		bytes []byte
	}
	pcBuffer struct {
		pc []uintptr
	}
)

// nolint: gochecknoglobals
var (
	_bufferPool = &sync.Pool{New: func() interface{} {
		const bufferCap = 1024
		return &bytesBuffer{make([]byte, 0, bufferCap)}
	}}
	_pcBufferPool = &sync.Pool{New: func() interface{} {
		const bufferCap = 64
		return &pcBuffer{make([]uintptr, bufferCap)}
	}} // NOTE: both len and cap are needed.
)

func getBytesBuffer() (buf *bytesBuffer, put func()) {
	b := _bufferPool.Get().(*bytesBuffer) //nolint:forcetypeassert
	b.bytes = b.bytes[:0]
	return b, func() {
		_bufferPool.Put(b)
	}
}

func getPCBuffer() (buf *pcBuffer, put func()) {
	b := _pcBufferPool.Get().(*pcBuffer) //nolint:forcetypeassert
	return b, func() {
		_pcBufferPool.Put(b)
	}
}

const null = "null"

// nolint: cyclop
// appendJSONEscapedString.
func appendJSONEscapedString(dst []byte, s string) []byte {
	for i := 0; i < len(s); i++ {
		if s[i] != '"' && s[i] != '\\' && s[i] > 0x1F {
			dst = append(dst, s[i])

			continue
		}

		// cf. https://tools.ietf.org/html/rfc8259#section-7
		// ... MUST be escaped: quotation mark, reverse solidus, and the control characters (U+0000 through U+001F).
		switch s[i] {
		case '"', '\\':
			dst = append(dst, '\\', s[i])
		case '\b' /* 0x08 */ :
			dst = append(dst, '\\', 'b')
		case '\f' /* 0x0C */ :
			dst = append(dst, '\\', 'f')
		case '\n' /* 0x0A */ :
			dst = append(dst, '\\', 'n')
		case '\r' /* 0x0D */ :
			dst = append(dst, '\\', 'r')
		case '\t' /* 0x09 */ :
			dst = append(dst, '\\', 't')
		default:
			const hexTable string = "0123456789abcdef"
			// cf. https://github.com/golang/go/blob/70deaa33ebd91944484526ab368fa19c499ff29f/src/encoding/hex/hex.go#L28-L29
			dst = append(dst, '\\', 'u', '0', '0', hexTable[s[i]>>4], hexTable[s[i]&0x0f])
		}
	}

	return dst
}

func appendFloatFieldValue(dst []byte, value float64, bitSize int) []byte {
	switch {
	case math.IsNaN(value):
		return append(dst, `"NaN"`...)
	case math.IsInf(value, 1):
		return append(dst, `"+Inf"`...)
	case math.IsInf(value, -1):
		return append(dst, `"-Inf"`...)
	}

	return strconv.AppendFloat(dst, value, 'f', -1, bitSize)
}

func appendCaller(dst []byte, callerSkip int, useLongCaller bool) []byte {
	pc, put := getPCBuffer()
	defer put()

	var frame runtime.Frame
	if runtime.Callers(callerSkip, pc.pc) > 0 {
		frame, _ = runtime.CallersFrames(pc.pc).Next()
	}

	return appendCallerFromFrame(dst, frame, useLongCaller)
}

// appendCallerFromFrame was split off from appendCaller in order to test different behaviors depending on the contents of the `runtime.Frame`.
func appendCallerFromFrame(dst []byte, frame runtime.Frame, useLongCaller bool) []byte {
	if useLongCaller {
		dst = appendJSONEscapedString(dst, frame.File)
	} else {
		dst = appendJSONEscapedString(dst, extractShortPath(frame.File))
	}

	dst = append(dst, ':')
	const base = 10
	dst = strconv.AppendInt(dst, int64(frame.Line), base)

	return dst
}

func extractShortPath(path string) string {
	// path == /path/to/directory/file
	//                           ~ <- idx
	idx := strings.LastIndexByte(path, '/')
	if idx == -1 {
		return path
	}

	// path[:idx] == /path/to/directory
	//                       ~ <- idx
	idx = strings.LastIndexByte(path[:idx], '/')
	if idx == -1 {
		return path
	}

	// path == /path/to/directory/file
	//                  ~~~~~~~~~~~~~~ <- filepath[idx+1:]
	return path[idx+1:]
}

func appendKey(dst []byte, key string) []byte {
	dst = append(dst, '"')
	dst = appendJSONEscapedString(dst, key)
	dst = append(dst, '"', ':')

	return dst
}

func appendLevelField(dst []byte, levels map[Level]string, level Level) []byte {
	v, ok := levels[level]
	if !ok {
		v = "DEBUG"
	}

	dst = append(dst, '"')
	dst = appendJSONEscapedString(dst, v)
	dst = append(dst, '"')
	return dst
}
