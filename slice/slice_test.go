package slice_test

import (
	"math"
	"reflect"
	"strconv"
	"testing"

	"github.com/kunitsuinc/util.go/slice"
)

func TestContains(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := true
		s := []int{0, 1, 2, 3}
		value := 1
		actual := slice.Contains(s, value)
		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		expect := false
		s := []int{0, 1, 2, 3}
		value := math.MaxInt
		actual := slice.Contains(s, value)
		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestDeepContains(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := true
		s := [][]int{{0}, {1}, {2}, {3}}
		value := []int{1}
		actual := slice.DeepContains(s, value)
		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		expect := false
		s := [][]int{{0}, {1}, {2}, {3}}
		value := []int{}
		actual := slice.DeepContains(s, value)
		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestEqual(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := true
		a := []int{0, 1, 2, 3}
		b := []int{0, 1, 2, 3}
		actual := slice.Equal(a, b)
		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(len)", func(t *testing.T) {
		t.Parallel()
		expect := false
		a := []int{0, 1, 2, 3}
		b := []int{1, 2, 3}
		actual := slice.Equal(a, b)
		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(contents)", func(t *testing.T) {
		t.Parallel()
		expect := false
		a := []int{0, 1, 2, 3}
		b := []int{1, 2, 3, 0}
		actual := slice.Equal(a, b)
		if expect != actual {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestExclude(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := []int{0, 2, 3}
		s := []int{0, 1, 2, 3}
		value := []int{1}
		actual := slice.Exclude(s, value)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		expect := []int{0, 1, 2, 3}
		s := []int{0, 1, 2, 3}
		value := []int{math.MaxInt}
		actual := slice.Exclude(s, value)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestDeepExclude(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := [][]int{{0}, {2}, {3}}
		s := [][]int{{0}, {1}, {2}, {3}}
		value := []int{1}
		actual := slice.DeepExclude(s, value)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		expect := [][]int{{0}, {1}, {2}, {3}}
		s := [][]int{{0}, {1}, {2}, {3}}
		value := []int{math.MaxInt}
		actual := slice.DeepExclude(s, value)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestEach(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := []int{0, 1, 2}
		s := []int{1, 2, 3}
		actual := make([]int, 0)
		slice.Each(s, func(_, i int) {
			actual = append(actual, i-1)
		})
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestFilter(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := []int{0, 1, 2}
		s := []int{0, 1, 2, 3}
		actual := slice.Filter(s, func(_, i int) bool {
			return i != 3
		})
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		expect := []int{0, 1, 2, 3}
		s := []int{0, 1, 2, 3}
		actual := slice.Filter(s, func(_, i int) bool {
			return i != math.MaxInt
		})
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestToMap(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := map[string]int{"0": 0, "1": 1, "2": 2}
		s := []int{0, 1, 2}
		actual := slice.ToMap(s, func(_, value int) string {
			key := strconv.Itoa(value)
			return key
		})
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestSelect(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := []string{"foo0", "foo1", "foo2"}
		s := []int{0, 1, 2}
		actual := slice.Select(s, func(_, value int) string {
			return "foo" + strconv.Itoa(value)
		})
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestFirst(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := "1"
		s := []string{"1", "2", "3"}
		actual := slice.First(s)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		expect := ""
		s := []string{}
		actual := slice.First(s)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestLast(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := "3"
		s := []string{"1", "2", "3"}
		actual := slice.Last(s)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		expect := ""
		s := []string{}
		actual := slice.Last(s)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestSplit(t *testing.T) {
	t.Parallel()
	t.Run("success(1)", func(t *testing.T) {
		t.Parallel()
		expect := [][]string{{"1", "2"}, {"3", "4"}, {"5", "6"}}
		s := []string{"1", "2", "3", "4", "5", "6"}
		actual := slice.Split(s, 2)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("success(2)", func(t *testing.T) {
		t.Parallel()
		expect := [][]string{{"1", "2"}, {"3", "4"}, {"5"}}
		s := []string{"1", "2", "3", "4", "5"}
		actual := slice.Split(s, 2)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("success(3)", func(t *testing.T) {
		t.Parallel()
		expect := [][]string{{"1", "2", "3", "4", "5"}}
		s := []string{"1", "2", "3", "4", "5"}
		actual := slice.Split(s, math.MaxInt)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestUniq(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := []string{"1", "2", "3"}
		s := []string{"1", "2", "3", "2", "3"}
		actual := slice.Uniq(s)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v != %v", expect, actual)
		}
	})
}
