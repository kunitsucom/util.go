// nolint: paralleltest
package rands_test

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/kunitsuinc/util.go/math/rands"
)

func TestRange(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		rand.Seed(0)
		expect := []int{2, 2, 1, 2, 3, 0}
		actual := []int{
			rands.Range(0, 3),
			rands.Range(0, 3),
			rands.Range(0, 3),
			rands.Range(0, 3),
			rands.Range(0, 3),
			rands.Range(0, 3),
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestRangeRand(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// nolint: gosec
		r := rand.New(rand.NewSource(0))
		expect := []int{2, 2, 1, 2, 3, 0}
		actual := []int{
			rands.RangeRand(r, 0, 3),
			rands.RangeRand(r, 0, 3),
			rands.RangeRand(r, 0, 3),
			rands.RangeRand(r, 0, 3),
			rands.RangeRand(r, 0, 3),
			rands.RangeRand(r, 0, 3),
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestRange32(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		rand.Seed(0)
		expect := []int32{2, 2, 1, 2, 3, 0}
		actual := []int32{
			rands.Range32(0, 3),
			rands.Range32(0, 3),
			rands.Range32(0, 3),
			rands.Range32(0, 3),
			rands.Range32(0, 3),
			rands.Range32(0, 3),
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestRange32Rand(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// nolint: gosec
		r := rand.New(rand.NewSource(0))
		expect := []int32{2, 2, 1, 2, 3, 0}
		actual := []int32{
			rands.Range32Rand(r, 0, 3),
			rands.Range32Rand(r, 0, 3),
			rands.Range32Rand(r, 0, 3),
			rands.Range32Rand(r, 0, 3),
			rands.Range32Rand(r, 0, 3),
			rands.Range32Rand(r, 0, 3),
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestRange64(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		rand.Seed(0)
		expect := []int64{1, 0, 3, 2, 2, 3}
		actual := []int64{
			rands.Range64(0, 3),
			rands.Range64(0, 3),
			rands.Range64(0, 3),
			rands.Range64(0, 3),
			rands.Range64(0, 3),
			rands.Range64(0, 3),
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestRange64Rand(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// nolint: gosec
		r := rand.New(rand.NewSource(0))
		expect := []int64{1, 0, 3, 2, 2, 3}
		actual := []int64{
			rands.Range64Rand(r, 0, 3),
			rands.Range64Rand(r, 0, 3),
			rands.Range64Rand(r, 0, 3),
			rands.Range64Rand(r, 0, 3),
			rands.Range64Rand(r, 0, 3),
			rands.Range64Rand(r, 0, 3),
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestRangeDuration(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		rand.Seed(0)
		expect := []time.Duration{1, 0, 3, 2, 2, 3}
		actual := []time.Duration{
			rands.RangeDuration(0, 3),
			rands.RangeDuration(0, 3),
			rands.RangeDuration(0, 3),
			rands.RangeDuration(0, 3),
			rands.RangeDuration(0, 3),
			rands.RangeDuration(0, 3),
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestRangeDurationRand(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// nolint: gosec
		r := rand.New(rand.NewSource(0))
		expect := []time.Duration{1, 0, 3, 2, 2, 3}
		actual := []time.Duration{
			rands.RangeDurationRand(r, 0, 3),
			rands.RangeDurationRand(r, 0, 3),
			rands.RangeDurationRand(r, 0, 3),
			rands.RangeDurationRand(r, 0, 3),
			rands.RangeDurationRand(r, 0, 3),
			rands.RangeDurationRand(r, 0, 3),
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}
