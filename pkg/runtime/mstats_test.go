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
		s := runtimez.NewMemStatsTicker(ctx, 1*time.Hour)
		MemStats := s.MemStats()
		if MemStats.Alloc == 0 {
			t.Errorf("❌: MemStats.Alloc == 0")
		}
		s.Stop()
		s.Restart()
		s.Reset(1 * time.Hour)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		s := runtimez.NewMemStatsTicker(ctx, 1*time.Hour)
		MemStats := s.MemStats()
		if MemStats.Alloc == 0 {
			t.Errorf("❌: MemStats.Alloc == 0")
		}
		cancel()
	})
}
