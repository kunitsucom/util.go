package runtimez_test

import (
	"context"
	"testing"
	"time"

	runtimez "github.com/kunitsuinc/util.go/pkg/runtime"
)

func TestMemStatsTicker(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		ticker := runtimez.NewMemStatsTicker(ctx, 100*time.Millisecond)
		m := ticker.MemStats()
		if m.Alloc == 0 {
			t.Errorf("❌: MemStats.Alloc == 0")
		}
		time.Sleep(1 * time.Second)
		ticker.Stop()
		ticker.Restart()
		ticker.Reset(1 * time.Hour)
	})

	t.Run("success(cancel)", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		ticker := runtimez.NewMemStatsTicker(ctx, 1*time.Hour)
		m := ticker.MemStats()
		if m.Alloc == 0 {
			t.Errorf("❌: MemStats.Alloc == 0")
		}
		cancel()
	})
}
