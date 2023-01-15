package mapz_test

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	mapz "github.com/kunitsuinc/util.go/pkg/maps"
	slicez "github.com/kunitsuinc/util.go/pkg/slices"
)

func TestSortStringKey(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()

		src := map[string]int{
			"2": 2,
			"4": 4,
			"a": 10,
			"c": 12,
			"b": 11,
			"5": 5,
			"3": 3,
			"1": 1,
		}

		expect := []struct {
			Key   string
			Value int
		}{
			{Key: "1", Value: 1},
			{Key: "2", Value: 2},
			{Key: "3", Value: 3},
			{Key: "4", Value: 4},
			{Key: "5", Value: 5},
			{Key: "a", Value: 10},
			{Key: "b", Value: 11},
			{Key: "c", Value: 12},
		}

		actual := mapz.SortStringKey(src)

		if !slicez.Equal(expect, actual) {
			t.Errorf("expect(%v) != actual(%v)", expect, actual)
		}
	})
}

func TestSortIntKey(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()

		src := map[int]string{
			2:  "2",
			4:  "4",
			10: "a",
			12: "c",
			11: "b",
			5:  "5",
			3:  "3",
			1:  "1",
		}

		expect := []struct {
			Key   int
			Value string
		}{
			{Key: 1, Value: "1"},
			{Key: 2, Value: "2"},
			{Key: 3, Value: "3"},
			{Key: 4, Value: "4"},
			{Key: 5, Value: "5"},
			{Key: 10, Value: "a"},
			{Key: 11, Value: "b"},
			{Key: 12, Value: "c"},
		}

		actual := mapz.SortIntKey(src)

		if !slicez.Equal(expect, actual) {
			t.Errorf("expect(%v) != actual(%v)", expect, actual)
		}
	})
}

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		const testKey = "testKey"
		expect := "testValue"
		h := make(map[string]any)
		h[testKey] = expect
		var actual string
		if err := mapz.Get(h, testKey, &actual); err != nil {
			t.Fatalf("❌: mapz.Get: err != nil: %v", err)
		}
		if expect != actual {
			t.Fatalf("❌: mapz.Get: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		const testKey = "testKey"
		type Expect struct {
			expect string
			if1    *http.Client
		}
		expect := &Expect{expect: "test", if1: http.DefaultClient}
		h := make(map[string]any)
		h[testKey] = expect
		var actual *Expect
		if err := mapz.Get(h, testKey, &actual); err != nil {
			t.Fatalf("❌: mapz.Get: err != nil: %v", err)
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Fatalf("❌: mapz.Get: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(mapz.ErrVIsNotPointerOrInterface)", func(t *testing.T) {
		t.Parallel()
		const testKey = "testKey"
		h := make(map[string]any)
		if err := mapz.Get(h, testKey, nil); err == nil || !errors.Is(err, mapz.ErrVIsNotPointerOrInterface) {
			t.Fatalf("❌: mapz.Get: err: %v", err)
		}
	})

	t.Run("failure(mapz.ErrKeyIsNotFound)", func(t *testing.T) {
		t.Parallel()
		const testKey = "testKey"
		h := make(map[string]any)
		var v string
		if err := mapz.Get(h, testKey, &v); err == nil || !errors.Is(err, mapz.ErrKeyIsNotFound) {
			t.Fatalf("❌: mapz.Get: err: %v", err)
		}
	})

	t.Run("failure(mapz.ErrValueTypeIsNotMatch)", func(t *testing.T) {
		t.Parallel()
		const testKey = "testKey"
		type Expect struct {
			expect string
			if1    interface{}
		}
		expect := &Expect{expect: "test", if1: "test"}
		h := make(map[string]any)
		h[testKey] = expect
		var actual string
		if err := mapz.Get(h, testKey, &actual); err == nil || !errors.Is(err, mapz.ErrValueTypeIsNotMatch) {
			t.Fatalf("❌: mapz.Get: err: %v", err)
		}
	})
}
