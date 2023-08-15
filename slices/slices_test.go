package slicez_test

import (
	"math"
	"reflect"
	"strconv"
	"testing"

	slicez "github.com/kunitsucom/util.go/slices"
)

func TestSort(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := []int{0, 1, 2, 3}
		s := []int{1, 3, 2, 0}
		actual := slicez.Sort(s, func(a, b int) bool {
			return a < b
		})
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestCopy(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		src := []byte("TestString")
		dst := slicez.Copy(src)
		dst[0] = byte('t')
		expect := byte('T')
		actual := src[0]
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestContains(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := true
		s := []int{0, 1, 2, 3}
		value := 1
		actual := slicez.Contains(s, value)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		expect := false
		s := []int{0, 1, 2, 3}
		value := math.MaxInt
		actual := slicez.Contains(s, value)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
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
		actual := slicez.DeepContains(s, value)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		expect := false
		s := [][]int{{0}, {1}, {2}, {3}}
		value := []int{}
		actual := slicez.DeepContains(s, value)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
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
		actual := slicez.Equal(a, b)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(len)", func(t *testing.T) {
		t.Parallel()
		expect := false
		a := []int{0, 1, 2, 3}
		b := []int{1, 2, 3}
		actual := slicez.Equal(a, b)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(contents)", func(t *testing.T) {
		t.Parallel()
		expect := false
		a := []int{0, 1, 2, 3}
		b := []int{1, 2, 3, 0}
		actual := slicez.Equal(a, b)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
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
		actual := slicez.Exclude(s, value)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		expect := []int{0, 1, 2, 3}
		s := []int{0, 1, 2, 3}
		value := []int{math.MaxInt}
		actual := slicez.Exclude(s, value)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
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
		actual := slicez.DeepExclude(s, value)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		expect := [][]int{{0}, {1}, {2}, {3}}
		s := [][]int{{0}, {1}, {2}, {3}}
		value := []int{math.MaxInt}
		actual := slicez.DeepExclude(s, value)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
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
		slicez.Each(s, func(_, i int) {
			actual = append(actual, i-1)
		})
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestFilter(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := []int{0, 1, 2}
		s := []int{0, 1, 2, 3}
		actual := slicez.Filter(s, func(_, i int) bool {
			return i != 3
		})
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		expect := []int{0, 1, 2, 3}
		s := []int{0, 1, 2, 3}
		actual := slicez.Filter(s, func(_, i int) bool {
			return i != math.MaxInt
		})
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestToMap(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := map[string]int{"0": 0, "1": 1, "2": 2}
		s := []int{0, 1, 2}
		actual := slicez.ToMap(s, func(_, value int) string {
			key := strconv.Itoa(value)
			return key
		})
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestSelect(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := []string{"foo0", "foo1", "foo2"}
		s := []int{0, 1, 2}
		actual := slicez.Select(s, func(_, value int) string {
			return "foo" + strconv.Itoa(value)
		})
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestFirst(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := "1"
		s := []string{"1", "2", "3"}
		actual := slicez.First(s)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		expect := ""
		s := []string{}
		actual := slicez.First(s)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestLast(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := "3"
		s := []string{"1", "2", "3"}
		actual := slicez.Last(s)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		expect := ""
		s := []string{}
		actual := slicez.Last(s)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestSplit(t *testing.T) {
	t.Parallel()
	t.Run("success(1)", func(t *testing.T) {
		t.Parallel()
		expect := [][]string{{"1", "2"}, {"3", "4"}, {"5", "6"}}
		s := []string{"1", "2", "3", "4", "5", "6"}
		actual := slicez.Split(s, 2)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("success(2)", func(t *testing.T) {
		t.Parallel()
		expect := [][]string{{"1", "2"}, {"3", "4"}, {"5"}}
		s := []string{"1", "2", "3", "4", "5"}
		actual := slicez.Split(s, 2)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("success(3)", func(t *testing.T) {
		t.Parallel()
		expect := [][]string{{"1", "2", "3", "4", "5"}}
		s := []string{"1", "2", "3", "4", "5"}
		actual := slicez.Split(s, math.MaxInt)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestUniq(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := []string{"1", "2", "3"}
		s := []string{"1", "2", "3", "2", "3"}
		actual := slicez.Uniq(s)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestReverse(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expect := []string{"3", "2", "1", "0"}
		s := []string{"0", "1", "2", "3"}
		actual := slicez.Reverse(s)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestCutOff(t *testing.T) {
	t.Parallel()
	t.Run("success(more)", func(t *testing.T) {
		t.Parallel()
		expect := []string{"0", "1", "2"}
		s := []string{"0", "1", "2", "3"}
		actual := slicez.CutOff(s, 3)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
	t.Run("success(equal)", func(t *testing.T) {
		t.Parallel()
		expect := []string{"0", "1", "2", "3"}
		s := []string{"0", "1", "2", "3"}
		actual := slicez.CutOff(s, 4)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
	t.Run("success(less)", func(t *testing.T) {
		t.Parallel()
		expect := []string{"0", "1", "2", "3"}
		s := []string{"0", "1", "2", "3"}
		actual := slicez.CutOff(s, 5)
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}
