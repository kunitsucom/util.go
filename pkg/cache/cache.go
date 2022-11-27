package cache

import (
	"sync"
	"time"
)

type cache[T interface{}] struct {
	expirationTime time.Time
	value          T
}

func (c cache[T]) expired(now time.Time) bool {
	return c.expirationTime.Before(now)
}

type Store[T interface{}] struct {
	defaultTTL time.Duration
	cacheMap   map[string]cache[T]
	cacheMutex sync.Mutex
}

type StoreOption[T interface{}] func(*Store[T])

func NewStore[T interface{}](opts ...StoreOption[T]) *Store[T] {
	c := &Store[T]{
		defaultTTL: 1 * time.Minute,
		cacheMap:   make(map[string]cache[T]),
		cacheMutex: sync.Mutex{},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithDefaultTTL[T interface{}](ttl time.Duration) StoreOption[T] {
	return func(s *Store[T]) { s.defaultTTL = ttl }
}

// GetOrSet gets cache value T, or set the value T that returns getValue.
// If getValue does not return err, cache the value T.
func (c *Store[T]) GetOrSet(key string, getValue func() (T, error)) (T, error) { //nolint:ireturn
	return c.GetOrSetWithTTL(key, getValue, c.defaultTTL)
}

// GetOrSet gets cache value T, or set the value T that returns getValue with TTL.
// If getValue does not return err, cache the value T.
func (c *Store[T]) GetOrSetWithTTL(key string, getValue func() (T, error), ttl time.Duration) (T, error) { //nolint:ireturn
	return c.getOrSet(key, getValue, ttl, time.Now())
}

func (c *Store[T]) getOrSet(key string, getValue func() (T, error), ttl time.Duration, now time.Time) (T, error) { //nolint:ireturn
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	if !c.cacheMap[key].expired(now) {
		return c.cacheMap[key].value, nil
	}

	value, err := getValue()
	if err != nil {
		var zero T
		return zero, err
	}

	c.cacheMap[key] = cache[T]{
		expirationTime: now.Add(ttl),
		value:          value,
	}

	return c.cacheMap[key].value, nil
}
