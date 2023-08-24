package timez

import (
	"context"
	"sync"
	"time"
)

type ctxKeyNow struct{}

func Now(ctx context.Context) time.Time {
	_nowFuncRWMu.RLock()
	t := _nowFunc(ctx)
	_nowFuncRWMu.RUnlock()
	return t
}

//nolint:gochecknoglobals
var (
	_nowFunc     NowFunc = DefaultNowFunc //nolint:revive
	_nowFuncRWMu sync.RWMutex
)

type NowFunc = func(ctx context.Context) time.Time

func DefaultNowFunc(ctx context.Context) time.Time {
	if now, ok := FromContext(ctx); ok {
		return now
	}

	return time.Now()
}

func SetNowFunc(nowFunc NowFunc) (backup NowFunc) {
	_nowFuncRWMu.Lock()
	backup = _nowFunc
	_nowFunc = nowFunc
	_nowFuncRWMu.Unlock()
	return backup
}

func WithContext(ctx context.Context, now time.Time) context.Context {
	return context.WithValue(ctx, ctxKeyNow{}, now)
}

func FromContext(ctx context.Context) (time.Time, bool) {
	now, ok := ctx.Value(ctxKeyNow{}).(time.Time)
	return now, ok
}
