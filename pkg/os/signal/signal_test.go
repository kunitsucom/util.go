package signalz //nolint:testpackage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestNotify(t *testing.T) {
	t.Parallel()
	source := make(chan os.Signal, 1)
	c := Notify(source, os.Interrupt)
	close(c)
}

func TestNotifyContext(t *testing.T) {
	t.Parallel()
	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		ctx, stop := NotifyContext(ctx, nil, os.Interrupt)
		stop(nil)
		if !errors.Is(ctx.Err(), context.Canceled) {
			t.Errorf("actual != expect: %v != %v", ctx.Err(), context.Canceled)
		}
	})
}

func Test_notifyContext(t *testing.T) {
	t.Parallel()
	t.Run("success(DefaultContextSignalHandler)", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		ch := make(chan os.Signal, 1)
		ctx, _ = notifyContext(ctx, ch, nil, os.Interrupt)
		ch <- os.Interrupt
		<-ctx.Done()
		if err := context.Cause(ctx); !errors.Is(err, context.Canceled) {
			t.Errorf("actual != expect: %v != %v", err, context.Canceled)
		}
	})
	t.Run("success(CustomContextSignalHandler)", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		ch := make(chan os.Signal, 1)
		ctx, _ = notifyContext(ctx, ch, func(sig os.Signal, stop context.CancelCauseFunc) {
			stop(fmt.Errorf("signal=%s: %w", sig, context.Canceled))
		}, os.Interrupt)
		ch <- os.Interrupt
		<-ctx.Done()
		if err := context.Cause(ctx); !errors.Is(err, context.Canceled) {
			t.Errorf("actual != expect: %v != %v", err, context.Canceled)
		}
	})
}
