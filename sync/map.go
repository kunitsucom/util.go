package syncz

import (
	"context"
	"sync"
	"time"
)

type val[T any] struct {
	val T
	exp time.Time
}

func (v val[T]) isExpired() bool { return v.exp.Before(time.Now()) }

type (
	Key        = interface{}
	Map[T any] interface {
		private()
		Load(key Key) (v T, ok bool)
		Len() int
		IsExpired(key Key) bool
		Store(key Key, value T)
		StoreTTL(key Key, value T, ttl time.Duration)
		LoadOrStore(key Key, value T) (v T, loaded bool)
		LoadAndDelete(key Key) (v T, loaded bool)
		Delete(key Key)
		Clear()
		Range(func(key Key, value T) bool)
	}
)

type (
	syncMapConfig struct {
		cleanerInterval time.Duration
		ttl             time.Duration
		useGoroutine    bool
	}
	syncMapConfigCleanerInterval time.Duration
	syncMapConfigDefaultTTL      time.Duration
)

func (c syncMapConfigDefaultTTL) apply(cfg *syncMapConfig) { cfg.ttl = time.Duration(c) }
func WithNewMapOptionTTL(ttl time.Duration) NewMapOption   { return syncMapConfigDefaultTTL(ttl) } //nolint:ireturn

func (c syncMapConfigCleanerInterval) apply(cfg *syncMapConfig) {
	cfg.cleanerInterval = time.Duration(c)
	cfg.useGoroutine = true
}

func WithNewMapOptionCleanerInterval(interval time.Duration) NewMapOption { //nolint:ireturn
	return syncMapConfigCleanerInterval(interval)
}

type NewMapOption interface{ apply(*syncMapConfig) }

type _Map[T any] struct {
	mu     sync.RWMutex
	kv     map[interface{}]*val[T]
	cfg    *syncMapConfig
	ticker *time.Ticker
}

const defaultTTL = time.Minute

func NewMap[T any](ctx context.Context, opts ...NewMapOption) Map[T] {
	c := &syncMapConfig{
		cleanerInterval: time.Minute,
		ttl:             defaultTTL,
	}

	for _, opt := range opts {
		opt.apply(c)
	}

	m := &_Map[T]{
		mu:  sync.RWMutex{},
		kv:  make(map[interface{}]*val[T]),
		cfg: c,
	}
	m.backgroundCleaner(ctx)
	return m
}

func (m *_Map[T]) private() {}

func (m *_Map[T]) Load(key Key) (v T, ok bool) { //nolint:ireturn
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.load(key)
}

func (m *_Map[T]) load(key Key) (v T, ok bool) { //nolint:ireturn
	if v, ok := m.kv[key]; ok && !v.isExpired() {
		return v.val, true
	}

	return v, false
}

func (m *_Map[T]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.kv)
}

func (m *_Map[T]) IsExpired(key Key) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.load(key)
	return !ok
}

func (m *_Map[T]) LoadOrStore(key Key, value T) (v T, loaded bool) { //nolint:ireturn
	m.mu.Lock()
	defer m.mu.Unlock()
	m.foregroundCleaner()
	if v, ok := m.load(key); ok {
		return v, true
	}
	m.kv[key] = &val[T]{
		val: value,
		exp: time.Now().Add(m.cfg.ttl),
	}
	return value, false
}

func (m *_Map[T]) LoadAndDelete(key Key) (v T, loaded bool) { //nolint:ireturn
	m.mu.Lock()
	defer m.mu.Unlock()
	m.foregroundCleaner()
	if v, ok := m.load(key); ok {
		delete(m.kv, key)
		return v, true
	}
	return v, false
}

func (m *_Map[T]) Store(key Key, value T) {
	// m.check() // not needed because of running in StoreTTL
	m.StoreTTL(key, value, m.cfg.ttl)
}

func (m *_Map[T]) StoreTTL(key Key, value T, ttl time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.foregroundCleaner()
	m.kv[key] = &val[T]{
		val: value,
		exp: time.Now().Add(ttl),
	}
}

func (m *_Map[T]) Delete(key Key) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.foregroundCleaner()
	delete(m.kv, key)
}

func (m *_Map[T]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	// m.check() // not need
	m.kv = make(map[interface{}]*val[T])
}

func (m *_Map[T]) Range(f func(key Key, value T) bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.foregroundCleaner()
	for k, v := range m.kv {
		if !f(k, v.val) {
			return
		}
	}
}

func (m *_Map[T]) foregroundCleaner() {
	if m.cfg.useGoroutine {
		return
	}
	for k, v := range m.kv {
		if v.isExpired() {
			delete(m.kv, k)
		}
	}
}

func (m *_Map[T]) backgroundCleaner(ctx context.Context) {
	if !m.cfg.useGoroutine {
		return
	}
	m.ticker = time.NewTicker(m.cfg.cleanerInterval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				m.ticker.Stop()
				return
			case <-m.ticker.C:
				m.mu.Lock()
				for k, v := range m.kv {
					if v.isExpired() {
						delete(m.kv, k)
					}
				}
				m.mu.Unlock()
			}
		}
	}()
}
