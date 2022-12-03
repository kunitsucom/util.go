package cache

import (
	"context"
	"sync"
	"time"
)

type cache[T interface{}] struct {
	undead         bool
	expirationTime time.Time
	value          T
}

func (c cache[T]) expired(now time.Time) bool {
	return !c.undead && c.expirationTime.Before(now)
}

type Store[T interface{}] struct {
	defaultTTL time.Duration
	cache      map[string]cache[T]
	mu         sync.Mutex
	ticker     *time.Ticker
}

type StoreOption[T interface{}] func(*Store[T])

func NewStore[T interface{}](ctx context.Context, opts ...StoreOption[T]) *Store[T] {
	s := &Store[T]{
		defaultTTL: 1 * time.Minute,
		cache:      make(map[string]cache[T]),
		mu:         sync.Mutex{},
		ticker:     time.NewTicker(1 * time.Second),
	}

	for _, opt := range opts {
		opt(s)
	}

	s.startRefresher(ctx)

	return s
}

func WithDefaultTTL[T interface{}](ttl time.Duration) StoreOption[T] {
	return func(s *Store[T]) { s.defaultTTL = ttl }
}

func WithRefreshInterval[T interface{}](interval time.Duration) StoreOption[T] {
	return func(s *Store[T]) { s.ticker = time.NewTicker(interval) }
}

func (s *Store[T]) startRefresher(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				now := time.Now()
				s.mu.Lock()
				for k, v := range s.cache {
					if v.expired(now) {
						s.Delete(k)
					}
				}
				s.mu.Unlock()
				<-s.ticker.C
			}
		}
	}()
}

func (s *Store[T]) ResetRefresher(interval time.Duration) {
	s.ticker.Reset(interval)
}

func (s *Store[T]) StopRefresher() {
	s.ticker.Stop()
}

// GetOrSet gets cache value T, or set the value T that returns getValue.
// If getValue does not return err, cache the value T.
func (s *Store[T]) GetOrSet(key string, getValue func() (T, error)) (T, error) { //nolint:ireturn
	return s.GetOrSetWithTTL(key, getValue, s.defaultTTL)
}

// GetOrSet gets cache value T, or set the value T that returns getValue with TTL.
// If getValue does not return err, cache the value T.
func (s *Store[T]) GetOrSetWithTTL(key string, getValue func() (T, error), ttl time.Duration) (T, error) { //nolint:ireturn
	return s.getOrSet(key, getValue, ttl, time.Now())
}

func (s *Store[T]) getOrSet(key string, getValue func() (T, error), ttl time.Duration, now time.Time) (T, error) { //nolint:ireturn
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.cache[key].expired(now) {
		return s.cache[key].value, nil
	}

	value, err := getValue()
	if err != nil {
		var zero T
		return zero, err
	}

	s.cache[key] = cache[T]{
		undead:         ttl == 0,
		expirationTime: now.Add(ttl),
		value:          value,
	}

	return s.cache[key].value, nil
}

func (s *Store[T]) Delete(key string) {
	if s.mu.TryLock() {
		defer s.mu.Unlock()
	}
	delete(s.cache, key)
}

func (s *Store[T]) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache = make(map[string]cache[T])
}
