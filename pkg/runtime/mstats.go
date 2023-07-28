package runtimez

import (
	"context"
	"runtime"
	"time"
)

type MemStatsTicker struct {
	stats    *runtime.MemStats
	interval time.Duration
	ticker   *time.Ticker
}

func NewMemStatsTicker(ctx context.Context, interval time.Duration) *MemStatsTicker {
	s := &MemStatsTicker{
		stats:    &runtime.MemStats{},
		interval: interval,
		ticker:   time.NewTicker(interval),
	}

	s.start(ctx)

	return s
}

func (s *MemStatsTicker) ReadMemStats() {
	runtime.ReadMemStats(s.stats)
}

func (s *MemStatsTicker) MemStats() *runtime.MemStats {
	return s.stats
}

func (s *MemStatsTicker) start(ctx context.Context) {
	s.ReadMemStats()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-s.ticker.C:
				s.ReadMemStats()
			}
		}
	}()
}

func (s *MemStatsTicker) Stop() {
	s.ticker.Stop()
}

func (s *MemStatsTicker) Restart() {
	s.ReadMemStats()
	s.ticker.Reset(s.interval)
}

func (s *MemStatsTicker) Reset(interval time.Duration) {
	s.interval = interval

	s.Restart()
}
