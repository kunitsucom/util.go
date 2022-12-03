package cache_test

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/kunitsuinc/util.go/pkg/cache"
)

func TestStore_GetOrSet(t *testing.T) {
	t.Parallel()
	t.Run("success(cached)", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		store := cache.NewStore[string](ctx)
		store.StopRefresher()
		store.ResetRefresher(10 * time.Millisecond)
		const key = "test_key"
		const value = "test value"
		zero, err := store.GetOrSet(key, func() (string, error) { return value, io.ErrUnexpectedEOF })
		if err == nil {
			t.Errorf("err != nil: %v", err)
		}
		if zero != "" {
			t.Errorf("expect != actual: %v != %v", value, zero)
		}
		got, err := store.GetOrSet(key, func() (string, error) { return value, nil })
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
		if got != value {
			t.Errorf("expect != actual: %v != %v", value, got)
		}
		cached, err := store.GetOrSet(key, func() (string, error) { return value, nil })
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
		if cached != value {
			t.Errorf("expect != actual: %v != %v", value, cached)
		}
	})

	t.Run("success(undead)", func(t *testing.T) {
		t.Parallel()
		store := cache.NewStore(context.Background(), cache.WithDefaultTTL[string](0))
		const key = "test_key"
		const value = "test value"
		got, err := store.GetOrSet(key, func() (string, error) { return value, nil })
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
		if got != value {
			t.Errorf("expect != actual: %v != %v", value, got)
		}
		time.Sleep(50 * time.Millisecond)
		cached, err := store.GetOrSet(key, func() (string, error) { return value, nil })
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
		if expect, actual := value, cached; expect != actual {
			t.Errorf("expect != actual: %v != %v", value, cached)
		}
		store.Delete(key)
		deleted, err := store.GetOrSet(key, func() (string, error) { return "", io.ErrUnexpectedEOF })
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Errorf("err != io.ErrUnexpectedEOF: %v", err)
		}
		if expect, actual := "", deleted; expect != actual {
			t.Errorf("expect != actual: %v != %v", value, cached)
		}
	})

	t.Run("success(expired)", func(t *testing.T) {
		t.Parallel()
		store := cache.NewStore(context.Background(), cache.WithDefaultTTL[string](50*time.Millisecond), cache.WithRefreshInterval[string](10*time.Millisecond))
		const key = "test_key"
		const value = "test value"
		got, err := store.GetOrSet(key, func() (string, error) { return value, nil })
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
		if got != value {
			t.Errorf("expect != actual: %v != %v", value, got)
		}
		time.Sleep(100 * time.Millisecond)
		const value2 = "notCached"
		notCached, err := store.GetOrSet(key, func() (string, error) { return value2, nil })
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
		if notCached != value2 {
			t.Errorf("expect != actual: %v != %v", value2, notCached)
		}
	})

	t.Run("success(Flush)", func(t *testing.T) {
		t.Parallel()
		store := cache.NewStore[string](context.Background())
		const key = "test_key"
		const value = "test value"
		got, err := store.GetOrSetWithTTL(key, func() (string, error) { return value, nil }, 0)
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
		if got != value {
			t.Errorf("expect != actual: %v != %v", value, got)
		}
		store.Flush()
		const value2 = "notCached"
		notCached, err := store.GetOrSet(key, func() (string, error) { return value2, nil })
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
		if notCached != value2 {
			t.Errorf("expect != actual: %v != %v", value2, notCached)
		}
	})
}
