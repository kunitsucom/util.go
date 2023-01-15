// nolint: paralleltest
package randz_test

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	randz "github.com/kunitsuinc/util.go/pkg/math/rand"
)

func TestRange(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		rand.Seed(0)
		expect := []int{2, 2, 1, 2, 3, 0}
		actual := []int{
			randz.Range(0, 3),
			randz.Range(0, 3),
			randz.Range(0, 3),
			randz.Range(0, 3),
			randz.Range(0, 3),
			randz.Range(0, 3),
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestRangeRand(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// nolint: gosec
		r := rand.New(rand.NewSource(0))
		expect := []int{2, 2, 1, 2, 3, 0}
		actual := []int{
			randz.RangeRand(r, 0, 3),
			randz.RangeRand(r, 0, 3),
			randz.RangeRand(r, 0, 3),
			randz.RangeRand(r, 0, 3),
			randz.RangeRand(r, 0, 3),
			randz.RangeRand(r, 0, 3),
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestRange32(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		rand.Seed(0)
		expect := []int32{2, 2, 1, 2, 3, 0}
		actual := []int32{
			randz.Range31(0, 3),
			randz.Range31(0, 3),
			randz.Range31(0, 3),
			randz.Range31(0, 3),
			randz.Range31(0, 3),
			randz.Range31(0, 3),
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestRange32Rand(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// nolint: gosec
		r := rand.New(rand.NewSource(0))
		expect := []int32{2, 2, 1, 2, 3, 0}
		actual := []int32{
			randz.Range31Rand(r, 0, 3),
			randz.Range31Rand(r, 0, 3),
			randz.Range31Rand(r, 0, 3),
			randz.Range31Rand(r, 0, 3),
			randz.Range31Rand(r, 0, 3),
			randz.Range31Rand(r, 0, 3),
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestRange64(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		rand.Seed(0)
		expect := []int64{1, 0, 3, 2, 2, 3}
		actual := []int64{
			randz.Range63(0, 3),
			randz.Range63(0, 3),
			randz.Range63(0, 3),
			randz.Range63(0, 3),
			randz.Range63(0, 3),
			randz.Range63(0, 3),
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestRange64Rand(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// nolint: gosec
		r := rand.New(rand.NewSource(0))
		expect := []int64{1, 0, 3, 2, 2, 3}
		actual := []int64{
			randz.Range63Rand(r, 0, 3),
			randz.Range63Rand(r, 0, 3),
			randz.Range63Rand(r, 0, 3),
			randz.Range63Rand(r, 0, 3),
			randz.Range63Rand(r, 0, 3),
			randz.Range63Rand(r, 0, 3),
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestRangeDuration(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		rand.Seed(0)
		expect := []time.Duration{1, 0, 3, 2, 2, 3}
		actual := []time.Duration{
			randz.RangeDuration(0, 3),
			randz.RangeDuration(0, 3),
			randz.RangeDuration(0, 3),
			randz.RangeDuration(0, 3),
			randz.RangeDuration(0, 3),
			randz.RangeDuration(0, 3),
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestRangeDurationRand(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// nolint: gosec
		r := rand.New(rand.NewSource(0))
		expect := []time.Duration{1, 0, 3, 2, 2, 3}
		actual := []time.Duration{
			randz.RangeDurationRand(r, 0, 3),
			randz.RangeDurationRand(r, 0, 3),
			randz.RangeDurationRand(r, 0, 3),
			randz.RangeDurationRand(r, 0, 3),
			randz.RangeDurationRand(r, 0, 3),
			randz.RangeDurationRand(r, 0, 3),
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}
