package timez

import (
	"context"
	"sync"
	"time"
)

type ctxKeyNow struct{}

func Now(ctx context.Context) time.Time {
	nowFuncRWMu.RLock()
	defer nowFuncRWMu.RUnlock()
	return nowFunc(ctx)
}

//nolint:gochecknoglobals
var (
	nowFunc     = DefaultNowFunc
	nowFuncRWMu sync.RWMutex
)

type NowFunc = func(ctx context.Context) time.Time

func DefaultNowFunc(ctx context.Context) time.Time {
	if now, ok := FromContext(ctx); ok {
		return now
	}

	return time.Now()
}

func SetNowFunc(now NowFunc) (backup NowFunc) {
	nowFuncRWMu.Lock()
	defer nowFuncRWMu.Unlock()
	backup = nowFunc
	nowFunc = now
	return backup
}

func WithContext(ctx context.Context, now time.Time) context.Context {
	return context.WithValue(ctx, ctxKeyNow{}, now)
}

func FromContext(ctx context.Context) (time.Time, bool) {
	now, ok := ctx.Value(ctxKeyNow{}).(time.Time)
	return now, ok
}
