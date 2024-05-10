package ilog //nolint:testpackage

import (
	"bytes"
	"io"
	"log"
	"regexp"
	"testing"
	"time"
)

//nolint:paralleltest,tparallel
func TestGlobal(t *testing.T) {
	//nolint:paralleltest,tparallel
	t.Run("success", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		defer SetGlobal(NewBuilder(DebugLevel, NewSyncWriter(buf)).SetTimestampZone(time.UTC).Build())()
		expected := regexp.MustCompilePOSIX(`{"severity":"DEBUG","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+_test\.go:[0-9]+","message":"Logf","any":"any","bool":true,"bytes":"bytes","duration":"1h1m1.001001001s","error":"unexpected EOF","eof":"EOF","float32":1\.234567,"float64":1\.23456789,"int":-1,"int32":123456789,"int64":123456789,"string":"string","time":"2023-08-13T04:38:39\.123456789\+09:00","uint":1,"uint32":123456789,"uint64":123456789}`)
		l := Global().Copy().
			Any("any", "any").
			Bool("bool", true).
			Bytes("bytes", []byte("bytes")).
			Duration("duration", time.Hour+time.Minute+time.Second+time.Millisecond+time.Microsecond+time.Nanosecond).
			Err(io.ErrUnexpectedEOF).
			ErrWithKey("eof", io.EOF).
			Float32("float32", 1.234567).
			Float64("float64", 1.23456789).
			Int("int", -1).
			Int32("int32", 123456789).
			Int64("int64", 123456789).
			String("string", "string").
			Time("time", time.Date(2023, 8, 13, 4, 38, 39, 123456789, time.FixedZone("Asia/Tokyo", int(9*time.Hour/time.Second)))).
			Uint("uint", 1).
			Uint32("uint32", 123456789).
			Uint64("uint64", 123456789).
			Logger()
		l.Logf(DebugLevel, "Logf")
		if !expected.Match(buf.Bytes()) {
			t.Errorf("❌: !expected.Match(buf.Bytes()):\n%s", buf)
		}
		t.Logf("ℹ️: buf:\n%s", buf)
	})
}

//nolint:paralleltest,tparallel
func TestSetStdLogger(t *testing.T) {
	//nolint:paralleltest,tparallel
	t.Run("success", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		expected := regexp.MustCompilePOSIX(`{"severity":"DEBUG","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog.go/global_test.go:[0-9]+","message":"Print\\n","any":"any","bool":true,"bytes":"bytes","duration":"1h1m1.001001001s","error":"unexpected EOF","eof":"EOF","float32":1\.234567,"float64":1\.23456789,"int":-1,"int32":123456789,"int64":123456789,"string":"string","time":"2023-08-13T04:38:39\.123456789\+09:00","uint":1,"uint32":123456789,"uint64":123456789}` + "\n")
		l := NewBuilder(DebugLevel, NewSyncWriter(buf)).
			SetTimestampZone(time.UTC).
			Build().
			Any("any", "any").
			Bool("bool", true).
			Bytes("bytes", []byte("bytes")).
			Duration("duration", time.Hour+time.Minute+time.Second+time.Millisecond+time.Microsecond+time.Nanosecond).
			Err(io.ErrUnexpectedEOF).
			ErrWithKey("eof", io.EOF).
			Float32("float32", 1.234567).
			Float64("float64", 1.23456789).
			Int("int", -1).
			Int32("int32", 123456789).
			Int64("int64", 123456789).
			String("string", "string").
			Time("time", time.Date(2023, 8, 13, 4, 38, 39, 123456789, time.FixedZone("Asia/Tokyo", int(9*time.Hour/time.Second)))).
			Uint("uint", 1).
			Uint32("uint32", 123456789).
			Uint64("uint64", 123456789).
			Logger()
		defer SetStdLogger(l)()
		log.Print("Print")
		if !expected.Match(buf.Bytes()) {
			t.Errorf("❌: !expected.Match(buf.Bytes()):\n%s", buf)
		}
		t.Logf("ℹ️: buf:\n%s", buf)
	})
}
