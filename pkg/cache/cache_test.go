package cache_test

import (
	"io"
	"testing"

	"github.com/kunitsuinc/util.go/pkg/cache"
)

func TestStore_GetOrSet(t *testing.T) {
	t.Parallel()
	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		store := cache.New[string]()
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

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		store := cache.New(cache.WithDefaultTTL[string](0))
		const key = "test_key"
		const value = "test value"
		got, err := store.GetOrSet(key, func() (string, error) { return value, nil })
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
		if got != value {
			t.Errorf("expect != actual: %v != %v", value, got)
		}
		const value2 = "not cached"
		notCached, err := store.GetOrSet(key, func() (string, error) { return value2, nil })
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
		if notCached != value2 {
			t.Errorf("expect != actual: %v != %v", value2, notCached)
		}
	})
}
