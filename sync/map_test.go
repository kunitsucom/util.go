package syncz

import (
	"context"
	"testing"
	"time"
)

func TestNewMap(t *testing.T) {
	t.Parallel()

	t.Run("success,all,foreground", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		m := NewMap[[]string](ctx, WithNewMapOptionTTL(1*time.Second), WithNewMapOptionCleanerInterval(1*time.Millisecond))
		if m == nil {
			t.Errorf("❌: NewMap: m == nil")
		}

		m.Store("key", []string{"value"})
		m.Store("key2", []string{"value2"})
		if actual, ok := m.LoadOrStore("key3", []string{"value3"}); ok || actual == nil {
			t.Errorf("❌: m.LoadOrStore(): expect(%v, %v) != actual(%v, %v)", []string{"value3"}, false, actual, ok)
		}
		if actual, ok := m.LoadOrStore("key3", []string{"value3"}); !ok || actual == nil {
			t.Errorf("❌: m.LoadOrStore(): expect(%v, %v) != actual(%v, %v)", []string{"value3"}, true, actual, ok)
		}
		stored := [][]string{}
		m.Range(func(key interface{}, value []string) bool {
			stored = append(stored, value)
			return true
		})
		m.Range(func(key interface{}, value []string) bool { return false })
		if len(stored) != 3 {
			t.Errorf("❌: m.Range(): expect(%v) != actual(%v)", 3, len(stored))
		}
		m.StoreTTL("ItWillBeDeletedSoon", []string{"value"}, 1*time.Nanosecond)
		time.Sleep(10 * time.Millisecond)
		if expect, actual := 3, m.Len(); expect != actual {
			t.Errorf("❌: m.Len(): expect(%v) != actual(%v)", expect, actual)
		}
		if actual := m.IsExpired("ItWillBeDeletedSoon"); !actual {
			t.Errorf("❌: m.IsExpired(): ItWillBeDeletedSoon: %v", actual)
		}
		if actual, ok := m.Load("key"); !ok || actual == nil {
			t.Errorf("❌: m.Load(): expect(%v, %v) != actual(%v, %v)", []string{"value"}, true, actual, ok)
		}
		if actual, ok := m.LoadAndDelete("key"); !ok || actual == nil {
			t.Errorf("❌: m.LoadAndDelete(): expect(%v, %v) != actual(%v, %v)", []string{"value"}, true, actual, ok)
		}
		if expect, actual := 2, m.Len(); expect != actual {
			t.Errorf("❌: m.Len(): expect(%v) != actual(%v)", expect, actual)
		}
		if actual, ok := m.LoadAndDelete("NotFound"); ok || actual != nil {
			t.Errorf("❌: m.LoadAndDelete(): expect(%v, %v) != actual(%v, %v)", []string{"value"}, true, actual, ok)
		}
		if expect, actual := 2, m.Len(); expect != actual {
			t.Errorf("❌: m.Len(): expect(%v) != actual(%v)", expect, actual)
		}
		m.Delete("key2")
		if expect, actual := 1, m.Len(); expect != actual {
			t.Errorf("❌: m.Len(): expect(%v) != actual(%v)", expect, actual)
		}
		m.Clear()
		if v := m.Len(); v != 0 {
			t.Errorf("❌: m.Len(): expect(%v) != actual(%v)", 0, v)
		}
		cancel()
		time.Sleep(10 * time.Millisecond)
	})

	t.Run("success,all,background", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		m := NewMap[[]string](ctx, WithNewMapOptionTTL(1*time.Second))
		if m == nil {
			t.Errorf("❌: NewMap: m == nil")
		}
		m.StoreTTL("key", []string{"value"}, 1*time.Nanosecond)
		time.Sleep(10 * time.Millisecond)
		if expect, actual := 1, m.Len(); expect != actual {
			t.Errorf("❌: m.Len(): expect(%v) != actual(%v)", expect, actual)
		}
		if isExpired := m.IsExpired("key"); !isExpired {
			t.Errorf("❌: m.IsExpired(): key: %v", isExpired)
		}
		if v, ok := m.Load("key"); ok || v != nil {
			t.Errorf("❌: m.Load(): expect(%v, %v) != actual(%v, %v)", nil, false, v, ok)
		}
		if v := m.Len(); v != 1 {
			t.Errorf("❌: m.Len(): expect(%v) != actual(%v)", 1, v)
		}
		if v, loaded := m.LoadAndDelete("key"); loaded || v != nil {
			t.Errorf("❌: m.LoadAndDelete(): expect(%v, %v) != actual(%v, %v)", nil, false, v, loaded)
		}
		if v := m.Len(); v != 0 {
			t.Errorf("❌: m.Len(): expect(%v) != actual(%v)", 0, v)
		}
		if v, loaded := m.LoadOrStore("key", []string{"value"}); loaded || v == nil {
			t.Errorf("❌: m.LoadOrStore(): expect(%v, %v) != actual(%v, %v)", []string{"value"}, false, v, loaded)
		}
		if v := m.Len(); v != 1 {
			t.Errorf("❌: m.Len(): expect(%v) != actual(%v)", 1, v)
		}
		if v, loaded := m.LoadOrStore("key", []string{"value"}); !loaded || v == nil {
			t.Errorf("❌: m.LoadOrStore(): expect(%v, %v) != actual(%v, %v)", []string{"value"}, true, v, loaded)
		}
		if v := m.Len(); v != 1 {
			t.Errorf("❌: m.Len(): expect(%v) != actual(%v)", 1, v)
		}
		m.Delete("key")
		if v := m.Len(); v != 0 {
			t.Errorf("❌: m.Len(): expect(%v) != actual(%v)", 0, v)
		}
	})

	t.Run("success,misc", func(t *testing.T) {
		t.Parallel()
		m := &_Map[string]{}
		m.private()
	})
}
