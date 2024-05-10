package zap_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/kunitsucom/ilog.go"
	ilogzap "github.com/kunitsucom/ilog.go/implementations/zap"
)

func TestNew(t *testing.T) {
	t.Parallel()
	buf := bytes.NewBuffer(nil)
	l := ilogzap.New(ilog.DebugLevel, zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(buf), zapcore.DebugLevel))).
		Any("any", "any").
		Bool("bool", true).
		Bytes("bytes", []byte("bytes")).
		Duration("duration", time.Hour+time.Minute+time.Second+time.Millisecond+time.Microsecond+time.Nanosecond).
		Err(io.ErrUnexpectedEOF).
		ErrWithKey("err", io.ErrUnexpectedEOF).
		Float32("float32", 1.1).
		Float64("float64", 1.1).
		Int("int", 1).
		Int32("int32", 1).
		Int64("int64", 1).
		String("string", "string").
		Time("time", time.Date(2023, 8, 13, 4, 38, 39, 123456789, time.FixedZone("Asia/Tokyo", int(9*time.Hour/time.Second)))).
		Uint("uint", 1).
		Uint32("uint32", 1).
		Uint64("uint64", 1).
		Logger()

	l = l.String("append", "logger").Logger()

	l.String("string", "new logger").Debugf("debug message")

	t.Logf("ℹ️: buf:\n%s", buf)
}
