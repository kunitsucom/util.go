package ilog //nolint:testpackage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"path"
	"regexp"
	"testing"
	"time"
)

type testFormatter struct{}

func (f *testFormatter) Format(s fmt.State, verb rune) {
	_, _ = fmt.Fprint(s, "testFormatter")
}

type testFormatterError struct {
	err error
}

func (f *testFormatterError) Format(s fmt.State, verb rune) {
	_, _ = fmt.Fprint(s, f.err)
}

func (f *testFormatterError) Error() string {
	return f.err.Error()
}

type testStringer string

func (s testStringer) String() string {
	return string(s)
}

func TestScenario(t *testing.T) {
	t.Parallel()
	t.Run("success,JSON", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		defer t.Logf("ℹ️: buf:\n%s", buf)

		expected := regexp.MustCompilePOSIX(`{"severity":"DEBUG","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+_test\.go:[0-9]+","message":"Logf: format string","bool":true,"boolPointer":false,"boolPointer2":null,"byte":"\\u0001","bytes":"bytes","time\.Duration":"1h1m1.001001001s","error":"ilog: log entry not written","errorFormatter":"ilog: log entry not written","errorNil":null,"float32":1\.234567,"float64":1\.23456789,"float64NaN":"NaN","float64\+Inf":"\+Inf","float64-Inf":"-Inf","int":-1,"int8":-1,"int16":-1,"int32":123456789,"int64":123456789,"string":"string","stringEscaped":"\\b\\f\\n\\r\\t","time\.Time":"2023-08-13T04:38:39\.123456789\+09:00","uint":1,"uint16":1,"uint32":123456789,"uint64":123456789,"jsonSuccess":{"json":true},"jsonFailure":"json.Marshaler: v.MarshalJSON: unexpected EOF","jsonNull":null,"fmt\.Formatter":"testFormatter","fmt\.Stringer":"testStringer","fmt.StringerNull":null,"func":"0x[0-9a-f]+","mapSuccess":{"map":{"in":1}},"mapFailure":"map\[map:0x[0-9a-f]+\]","sliceSuccess":\["a","b"\],"sliceFailure":"\[0x[0-9a-f]+\]","append":"logger"}`)

		l := NewBuilder(DebugLevel, buf).
			SetTimestampZone(time.UTC).
			UseSyncWriter().
			Build().
			Any("bool", true).
			Any("boolPointer", new(bool)).
			Any("boolPointer2", (*bool)(nil)).
			Any("byte", byte(1)).
			Any("bytes", []byte("bytes")).
			Any("time.Duration", time.Hour+time.Minute+time.Second+time.Millisecond+time.Microsecond+time.Nanosecond).
			Any("error", ErrLogEntryIsNotWritten).
			Any("errorFormatter", &testFormatterError{ErrLogEntryIsNotWritten}).
			Any("errorNil", nil).
			Any("float32", float32(1.234567)).
			Any("float64", float64(1.23456789)).
			Any("float64NaN", math.NaN()).
			Any("float64+Inf", math.Inf(1)).
			Any("float64-Inf", math.Inf(-1)).
			Any("int", int(-1)).
			Any("int8", int8(-1)).
			Any("int16", int16(-1)).
			Any("int32", int32(123456789)).
			Any("int64", int64(123456789)).
			Any("string", "string").
			Any("stringEscaped", "\b\f\n\r\t").
			Any("time.Time", time.Date(2023, 8, 13, 4, 38, 39, 123456789, time.FixedZone("Asia/Tokyo", int(9*time.Hour/time.Second)))).
			Any("uint", uint(1)).
			Any("uint16", uint16(1)).
			Any("uint32", uint32(123456789)).
			Any("uint64", uint64(123456789)).
			Any("jsonSuccess", &testJSONMarshaler{MockMarshalJSON: func() ([]byte, error) { return []byte(`{"json":true}`), nil }}).
			Any("jsonFailure", &testJSONMarshaler{MockMarshalJSON: func() ([]byte, error) { return nil, io.ErrUnexpectedEOF }}).
			Any("jsonNull", (*testJSONMarshaler)(nil)).
			Any("fmt.Formatter", &testFormatter{}).
			Any("fmt.Stringer", testStringer("testStringer")).
			Any("fmt.StringerNull", (*testStringer)(nil)).
			Any("func", func() {}).
			Any("mapSuccess", map[string]interface{}{"map": map[string]interface{}{"in": 1}}).
			Any("mapFailure", map[string]interface{}{"map": func() {}}).
			Any("sliceSuccess", []string{"a", "b"}).
			Any("sliceFailure", []func(){func() {}}).
			Logger()

		l = l.String("append", "logger").Logger()

		l.Logf(DebugLevel, "Logf: %s", "format string")
		if !expected.Match(buf.Bytes()) {
			t.Errorf("❌: !expected.Match(buf.Bytes()):\n%s", buf)
		}

		decoded := make(map[string]interface{})
		if err := json.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(&decoded); err != nil {
			t.Errorf("❌: err: %+v", err)
		}

		if expected, actual := "DEBUG", decoded["severity"]; expected != actual {
			t.Errorf("❌: severity: expected(%s) != actual(%s)", expected, actual)
		}

		if expected, actual := "Logf: format string", decoded["message"]; expected != actual {
			t.Errorf("❌: message: expected(%s) != actual(%s)", expected, actual)
		}

		if expected, actual := "+Inf", decoded["float64+Inf"]; expected != actual {
			t.Errorf("❌: float64+Inf: expected(%s) != actual(%s)", expected, actual)
		}
	})
}

func TestLogger(t *testing.T) {
	t.Parallel()
	t.Run("success,Logger", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		defer t.Logf("ℹ️: buf:\n%s", buf)

		expected := regexp.MustCompilePOSIX(`{"level":"DEBUG","time":"[0-9]+-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}","file":".+/ilog\.go/[a-z_]+_test\.go:[0-9]+","msg":"Logf"}` + "\r\n")

		const expectedLevel = DebugLevel
		l := NewBuilder(ErrorLevel, NewSyncWriter(buf)).
			UseSyncWriter().
			SetLevelKey("level").
			SetLevels(copyLevels(defaultLevels)).
			SetTimestampKey("time").
			SetTimestampFormat("2006-01-02 15:04:05").
			SetTimestampZone(time.UTC).
			SetCallerKey("file").
			UseLongCaller(true).
			SetMessageKey("msg").
			SetSeparator("\r\n").
			Build().
			SetLevel(expectedLevel).
			AddCallerSkip(10).
			AddCallerSkip(-10)

		if expected, actual := expectedLevel, l.Level(); expected != actual {
			t.Errorf("❌: expected(%d) != actual(%d)", expected, actual)
		}

		l.Logf(DebugLevel, "Logf")
		if !expected.Match(buf.Bytes()) {
			t.Errorf("❌: !expected.Match(buf.Bytes()):\n%s", buf)
		}
	})

	t.Run("success,Logger,common", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		defer t.Logf("ℹ️: buf:\n%s", buf)

		l := NewBuilder(DebugLevel, NewSyncWriter(buf)).SetTimestampKey("").SetCallerKey("").Build()
		l.Any("any", "any").Debugf("Debugf")
		l.Bool("bool", true).Debugf("Debugf")
		l.Bytes("bytes", []byte("bytes")).Debugf("Debugf")
		l.Duration("time.Duration", time.Hour+time.Minute+time.Second+time.Millisecond+time.Microsecond+time.Nanosecond).Debugf("Debugf")
		l.Err(io.ErrUnexpectedEOF).Debugf("Debugf")
		l.ErrWithKey("err", io.ErrUnexpectedEOF).Debugf("Debugf")
		l.ErrWithKey("errNull", (error)(nil)).Debugf("Debugf")
		l.Float32("float32", float32(1.234567)).Debugf("Debugf")
		l.Float64("float64", float64(1.23456789)).Debugf("Debugf")
		l.Int("int", int(-1)).Debugf("Debugf")
		l.Int32("int32", int32(-1)).Debugf("Debugf")
		l.Int64("int64", int64(-1)).Debugf("Debugf")
		l.String("string", "string").Debugf("Debugf")
		l.Time("time.Time", time.Date(2023, 8, 13, 4, 38, 39, 123456789, time.FixedZone("Asia/Tokyo", int(9*time.Hour/time.Second)))).Debugf("Debugf")
		l.Uint("uint", uint(1)).Debugf("Debugf")
		l.Uint32("uint32", uint32(123456789)).Debugf("Debugf")
		l.Uint64("uint64", uint64(123456789)).Debugf("Debugf")
		l.Debugf("Debugf")
		l.Infof("Infof")
		l.Warnf("Warnf")
		l.Errorf("Errorf")
		l.Any("any", "any").Debugf("Debugf")
		l.Any("any", "any").Infof("Infof")
		l.Any("any", "any").Warnf("Warnf")
		l.Any("any", "any").Errorf("Errorf")
		l.Any("any", "any").Logf(DebugLevel, "Logf")
		_, _ = l.Any("any", "any").Write([]byte("Write"))

		const expect = `{"severity":"DEBUG","message":"Debugf","any":"any"}
{"severity":"DEBUG","message":"Debugf","bool":true}
{"severity":"DEBUG","message":"Debugf","bytes":"bytes"}
{"severity":"DEBUG","message":"Debugf","time.Duration":"1h1m1.001001001s"}
{"severity":"DEBUG","message":"Debugf","error":"unexpected EOF"}
{"severity":"DEBUG","message":"Debugf","err":"unexpected EOF"}
{"severity":"DEBUG","message":"Debugf","errNull":null}
{"severity":"DEBUG","message":"Debugf","float32":1.234567}
{"severity":"DEBUG","message":"Debugf","float64":1.23456789}
{"severity":"DEBUG","message":"Debugf","int":-1}
{"severity":"DEBUG","message":"Debugf","int32":-1}
{"severity":"DEBUG","message":"Debugf","int64":-1}
{"severity":"DEBUG","message":"Debugf","string":"string"}
{"severity":"DEBUG","message":"Debugf","time.Time":"2023-08-13T04:38:39.123456789+09:00"}
{"severity":"DEBUG","message":"Debugf","uint":1}
{"severity":"DEBUG","message":"Debugf","uint32":123456789}
{"severity":"DEBUG","message":"Debugf","uint64":123456789}
{"severity":"DEBUG","message":"Debugf"}
{"severity":"INFO","message":"Infof"}
{"severity":"WARN","message":"Warnf"}
{"severity":"ERROR","message":"Errorf"}
{"severity":"DEBUG","message":"Debugf","any":"any"}
{"severity":"INFO","message":"Infof","any":"any"}
{"severity":"WARN","message":"Warnf","any":"any"}
{"severity":"ERROR","message":"Errorf","any":"any"}
{"severity":"DEBUG","message":"Logf","any":"any"}
{"severity":"DEBUG","message":"Write","any":"any"}
`

		if expected, actual := expect, buf.String(); expected != actual {
			t.Errorf("❌: expected(%s) != actual(%s)", expected, actual)
		}
	})
}

type testJSONMarshaler struct {
	MockMarshalJSON func() ([]byte, error)
}

func (m *testJSONMarshaler) MarshalJSON() ([]byte, error) {
	return m.MockMarshalJSON()
}

func TestLogEntry(t *testing.T) {
	t.Parallel()
	t.Run("success,LogEntry", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		defer t.Logf("ℹ️: buf:\n%s", buf)

		expected := regexp.MustCompilePOSIX(`{"severity":"DEBUG","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+_test\.go:[0-9]+","message":"Logf: format string","bool":true,"boolPointer":false,"boolPointer2":null,"byte":"\\u0001","bytes":"bytes","time\.Duration":"1h1m1.001001001s","error":"ilog: log entry not written","errorFormatter":"ilog: log entry not written","errorNil":null,"float32":1\.234567,"float64":1\.23456789,"float64NaN":"NaN","float64\+Inf":"\+Inf","float64-Inf":"-Inf","int":-1,"int8":-1,"int16":-1,"int32":123456789,"int64":123456789,"string":"string","stringEscaped":"\\b\\f\\n\\r\\t","time\.Time":"2023-08-13T04:38:39\.123456789\+09:00","uint":1,"uint16":1,"uint32":123456789,"uint64":123456789,"jsonSuccess":{"json":true},"jsonFailure":"json.Marshaler: v.MarshalJSON: unexpected EOF","jsonNull":null,"fmt\.Formatter":"testFormatter","fmt\.Stringer":"testStringer","fmt.StringerNull":null,"func":"0x[0-9a-f]+","mapSuccess":{"map":{"in":1}},"mapFailure":"map\[map:0x[0-9a-f]+\]","sliceSuccess":\["a","b"\],"sliceFailure":"\[0x[0-9a-f]+\]"}`)

		le := NewBuilder(DebugLevel, NewSyncWriter(buf)).
			SetTimestampZone(time.UTC).
			Build().
			Any("bool", true).
			Any("boolPointer", new(bool)).
			Any("boolPointer2", (*bool)(nil)).
			Any("byte", byte(1)).
			Any("bytes", []byte("bytes")).
			Any("time.Duration", time.Hour+time.Minute+time.Second+time.Millisecond+time.Microsecond+time.Nanosecond).
			Any("error", ErrLogEntryIsNotWritten).
			Any("errorFormatter", &testFormatterError{ErrLogEntryIsNotWritten}).
			Any("errorNil", nil).
			Any("float32", float32(1.234567)).
			Any("float64", float64(1.23456789)).
			Any("float64NaN", math.NaN()).
			Any("float64+Inf", math.Inf(1)).
			Any("float64-Inf", math.Inf(-1)).
			Any("int", int(-1)).
			Any("int8", int8(-1)).
			Any("int16", int16(-1)).
			Any("int32", int32(123456789)).
			Any("int64", int64(123456789)).
			Any("string", "string").
			Any("stringEscaped", "\b\f\n\r\t").
			Any("time.Time", time.Date(2023, 8, 13, 4, 38, 39, 123456789, time.FixedZone("Asia/Tokyo", int(9*time.Hour/time.Second)))).
			Any("uint", uint(1)).
			Any("uint16", uint16(1)).
			Any("uint32", uint32(123456789)).
			Any("uint64", uint64(123456789)).
			Any("jsonSuccess", &testJSONMarshaler{MockMarshalJSON: func() ([]byte, error) { return []byte(`{"json":true}`), nil }}).
			Any("jsonFailure", &testJSONMarshaler{MockMarshalJSON: func() ([]byte, error) { return nil, io.ErrUnexpectedEOF }}).
			Any("jsonNull", (*testJSONMarshaler)(nil)).
			Any("fmt.Formatter", &testFormatter{}).
			Any("fmt.Stringer", testStringer("testStringer")).
			Any("fmt.StringerNull", (*testStringer)(nil)).
			Any("func", func() {}).
			Any("mapSuccess", map[string]interface{}{"map": map[string]interface{}{"in": 1}}).
			Any("mapFailure", map[string]interface{}{"map": func() {}}).
			Any("sliceSuccess", []string{"a", "b"}).
			Any("sliceFailure", []func(){func() {}})

		if expected, actual := ErrLogEntryIsNotWritten.Error(), le.Error(); expected != actual {
			t.Errorf("❌: expected(%s) != actual(%s)", expected, actual)
		}

		le.Logf(DebugLevel, "Logf: %s", "format string")
		if !expected.Match(buf.Bytes()) {
			t.Errorf("❌: !expected.Match(buf.Bytes()):\n%s", buf)
		}
	})

	t.Run("success,default,DEBUG", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		defer t.Logf("ℹ️: buf:\n%s", buf)

		expected := regexp.MustCompilePOSIX(`{"severity":"DEBUG","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+_test\.go:[0-9]+","message":"default"}`)

		NewBuilder(-128, NewSyncWriter(buf)).SetTimestampZone(time.UTC).Build().Logf(-128, "default")

		if !expected.Match(buf.Bytes()) {
			t.Errorf("❌: !expected.Match(buf.Bytes()):\n%s", buf)
		}
	})

	t.Run("success,Debugf", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		defer t.Logf("ℹ️: buf:\n%s", buf)

		expected := regexp.MustCompilePOSIX(`{"severity":"DEBUG","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+_test\.go:[0-9]+","message":"Debugf"}`)

		NewBuilder(DebugLevel, NewSyncWriter(buf)).SetTimestampZone(time.UTC).Build().Debugf("Debugf")

		if !expected.Match(buf.Bytes()) {
			t.Errorf("❌: !expected.Match(buf.Bytes()):\n%s", buf)
		}
	})

	t.Run("success,Infof", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		defer t.Logf("ℹ️: buf:\n%s", buf)

		expected := regexp.MustCompilePOSIX(`{"severity":"INFO","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+_test\.go:[0-9]+","message":"Infof"}`)

		NewBuilder(DebugLevel, NewSyncWriter(buf)).SetTimestampZone(time.UTC).Build().Infof("Infof")

		if !expected.Match(buf.Bytes()) {
			t.Errorf("❌: !expected.Match(buf.Bytes()):\n%s", buf)
		}
	})

	t.Run("success,Warnf", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		defer t.Logf("ℹ️: buf:\n%s", buf)

		expected := regexp.MustCompilePOSIX(`{"severity":"WARN","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+_test\.go:[0-9]+","message":"Warnf"}`)

		NewBuilder(DebugLevel, NewSyncWriter(buf)).SetTimestampZone(time.UTC).Build().Warnf("Warnf")

		if !expected.Match(buf.Bytes()) {
			t.Errorf("❌: !expected.Match(buf.Bytes()):\n%s", buf)
		}
	})

	t.Run("success,Errorf", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		defer t.Logf("ℹ️: buf:\n%s", buf)

		expected := regexp.MustCompilePOSIX(`{"severity":"ERROR","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+_test\.go:[0-9]+","message":"Errorf"}`)

		NewBuilder(DebugLevel, NewSyncWriter(buf)).SetTimestampZone(time.UTC).Build().Errorf("Errorf")

		if !expected.Match(buf.Bytes()) {
			t.Errorf("❌: !expected.Match(buf.Bytes()):\n%s", buf)
		}
	})
}

func TestLogger_logf(t *testing.T) {
	t.Parallel()
	t.Run("success,empty", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)

		const expected = ""

		NewBuilder(InfoLevel, NewSyncWriter(buf)).SetTimestampZone(time.UTC).Build().Debugf("Debugf")
		if expected != buf.String() {
			t.Errorf("❌: expected(%s) != actual(%s)", expected, buf.String())
		}
		buf.Reset()

		NewBuilder(WarnLevel, NewSyncWriter(buf)).SetTimestampZone(time.UTC).Build().Infof("Infof")
		if expected != buf.String() {
			t.Errorf("❌: expected(%s) != actual(%s)", expected, buf.String())
		}
		buf.Reset()

		NewBuilder(ErrorLevel, NewSyncWriter(buf)).SetTimestampZone(time.UTC).Build().Warnf("Warnf")
		if expected != buf.String() {
			t.Errorf("❌: expected(%s) != actual(%s)", expected, buf.String())
		}
		buf.Reset()
	})

	t.Run("success,{}", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		defer t.Logf("ℹ️: buf:\n%s", buf)

		const expected = "{}\n"

		NewBuilder(DebugLevel, NewSyncWriter(buf)).
			SetLevelKey("").
			SetTimestampKey("").
			SetCallerKey("").
			SetMessageKey("").Build().Debugf("{}")

		if expected != buf.String() {
			t.Errorf("❌: expected(%s) != actual(%s)", expected, buf.String())
		}
	})
}

type testWriter struct {
	err error
}

func (w *testWriter) Write(p []byte) (n int, err error) {
	return 0, w.err
}

//nolint:paralleltest,tparallel
func TestLogger_Write(t *testing.T) {
	//nolint:paralleltest,tparallel
	t.Run("failure,Logger,Write", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		defer t.Logf("ℹ️: buf:\n%s", buf)

		defer SetGlobal(NewBuilder(DebugLevel, NewSyncWriter(buf)).SetTimestampZone(time.UTC).Build())()

		i, err := NewBuilder(DebugLevel, &testWriter{err: io.ErrUnexpectedEOF}).SetTimestampZone(time.UTC).Build().Write([]byte("ERROR"))
		if expected := regexp.MustCompilePOSIX(`w.logf: w.logger.writer.Write: p={"severity":"DEBUG","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+_test\.go:[0-9]+","message":"ERROR"}: unexpected EOF`); err == nil || !expected.MatchString(err.Error()) {
			t.Errorf("❌: err != nil: %v", err)
		}
		if i != 0 {
			t.Errorf("❌: i != 0: %d", i)
		}
		if expected := regexp.MustCompilePOSIX(`{"severity":"ERROR","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+\.go:[0-9]+","message":"w.logger.writer.Write: p={\\"severity\\":\\"DEBUG\\",\\"timestamp\\":\\"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z\\",\\"caller\\":\\"ilog\.go/[a-z_]+_test\.go:[0-9]+\\",\\"message\\":\\"ERROR(\\n)?\\"}: unexpected EOF"}`); !expected.Match(buf.Bytes()) {
			t.Errorf("❌: !expected.Match(buf.Bytes()):\n%s", buf)
		}
	})

	t.Run("failure,LogEntry,Write", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		defer t.Logf("ℹ️: buf:\n%s", buf)

		defer SetGlobal(NewBuilder(DebugLevel, NewSyncWriter(buf)).SetTimestampZone(time.UTC).Build())()

		i, err := NewBuilder(DebugLevel, &testWriter{err: io.ErrUnexpectedEOF}).SetTimestampZone(time.UTC).Build().Any("any", "any").Write([]byte("ERROR"))
		if expected := regexp.MustCompilePOSIX(`w.logf: w.logger.writer.Write: p={"severity":"DEBUG","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+_test\.go:[0-9]+","message":"ERROR","any":"any"}: unexpected EOF`); err == nil || !expected.MatchString(err.Error()) {
			t.Errorf("❌: err != nil: %v", err)
		}
		if i != 0 {
			t.Errorf("❌: i != 0: %d", i)
		}
		if expected := regexp.MustCompilePOSIX(`{"severity":"ERROR","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+\.go:[0-9]+","message":"w.logger.writer.Write: p={\\"severity\\":\\"DEBUG\\",\\"timestamp\\":\\"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z\\",\\"caller\\":\\"ilog\.go/[a-z_]+_test\.go:[0-9]+\\",\\"message\\":\\"ERROR(\\n)?\\",\\"any\\":\\"any\\"}: unexpected EOF"}`); !expected.Match(buf.Bytes()) {
			t.Errorf("❌: !expected.Match(buf.Bytes()):\n%s", buf)
		}
	})
}

func Test_extractShortPath(t *testing.T) {
	t.Parallel()
	t.Run("success,noIndex", func(t *testing.T) {
		t.Parallel()
		const expected = "expected"
		actual := extractShortPath(expected)
		if expected != actual {
			t.Errorf("❌: expected(%s) != actual(%s)", expected, actual)
		}
	})
	t.Run("success,1Index", func(t *testing.T) {
		t.Parallel()
		const expected = "expected/expected"
		actual := extractShortPath(expected)
		if expected != actual {
			t.Errorf("❌: expected(%s) != actual(%s)", expected, actual)
		}
	})

	t.Run("success,1Index", func(t *testing.T) {
		t.Parallel()
		const expected = "expected/expected"
		actual := extractShortPath(path.Join(expected, expected))
		if expected != actual {
			t.Errorf("❌: expected(%s) != actual(%s)", expected, actual)
		}
	})
}
