package ilog //nolint:testpackage

import (
	"bytes"
	"context"
	"regexp"
	"testing"
	"time"
)

//nolint:paralleltest,tparallel
func TestContext(t *testing.T) {
	//nolint:paralleltest,tparallel
	t.Run("success", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		ctx := WithContext(context.Background(), NewBuilder(DebugLevel, NewSyncWriter(buf)).SetTimestampZone(time.UTC).Build())
		l := FromContext(ctx)
		l.Debugf("Debugf")
		if expected := regexp.MustCompilePOSIX(`{"severity":"DEBUG","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+_test\.go:[0-9]+","message":"Debugf"}`); !expected.Match(buf.Bytes()) {
			t.Errorf("❌: !expected.Match(buf.Bytes()):\n%s", buf)
		}
		t.Logf("ℹ️: buf:\n%s", buf)
	})

	t.Run("failure,nilContext", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		defer SetGlobal(NewBuilder(DebugLevel, NewSyncWriter(buf)).SetTimestampZone(time.UTC).Build())()
		_ = FromContext(nil) //nolint:staticcheck
		if expected := regexp.MustCompilePOSIX(`{"severity":"ERROR","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+\.go:[0-9]+","message":"ilog: nil context"}`); !expected.Match(buf.Bytes()) {
			t.Errorf("❌: !expected.Match(buf.Bytes()):\n%s", buf)
		}
		t.Logf("ℹ️: buf:\n%s", buf)
	})

	t.Run("failure,invalidType", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		defer SetGlobal(NewBuilder(DebugLevel, NewSyncWriter(buf)).SetTimestampZone(time.UTC).Build())()
		_ = FromContext(context.WithValue(context.Background(), contextKeyLogger{}, "invalid")) //nolint:staticcheck
		if expected := regexp.MustCompilePOSIX(`{"severity":"ERROR","timestamp":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.?[0-9]*Z","caller":"ilog\.go/[a-z_]+\.go:[0-9]+","message":"ilog: type assertion failed: expected=ilog.Logger, actual=string, value=\\"invalid\\""}`); !expected.Match(buf.Bytes()) {
			t.Errorf("❌: !expected.Match(buf.Bytes()):\n%s", buf)
		}
		t.Logf("ℹ️: buf:\n%s", buf)
	})
}
