package mapz_test

import (
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
