package mathz_test

import (
	"testing"

	mathz "github.com/kunitsuinc/util.go/pkg/math"
)

func TestIsPow10(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		n   string
		num float64
	}
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		for _, v := range []testStruct{{"1e-323", 1e-323}, {"1", 1}, {"10", 10}, {"100", 100}, {"1000", 1000}, {"10000", 10000}, {"1e+308", 1e+308}} {
			if !mathz.IsPow10(v.num) {
				t.Errorf("%s should be pow10(x)", v.n)
			}
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		for _, v := range []testStruct{{"2", 2}, {"1e-324", 1e-324}} {
			if mathz.IsPow10(v.num) {
				t.Errorf("%s should not be pow10(x)", v.n)
			}
		}
	})
}

func TestIsPow10Int32(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		n   string
		num int32
	}
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		for _, v := range []testStruct{{"1", 1}, {"10", 10}, {"100", 100}, {"1000", 1000}, {"10000", 10000}} {
			if !mathz.IsPow10Int32(v.num) {
				t.Errorf("%s should be pow10(x)", v.n)
			}
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		for _, v := range []testStruct{{"2", 2}} {
			if mathz.IsPow10Int32(v.num) {
				t.Errorf("%s should not be pow10(x)", v.n)
			}
		}
	})
}

func TestIsPow10Int64(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		n   string
		num int64
	}
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		for _, v := range []testStruct{{"1", 1}, {"10", 10}, {"100", 100}, {"1000", 1000}, {"10000", 10000}} {
			if !mathz.IsPow10Int64(v.num) {
				t.Errorf("%s should be pow10(x)", v.n)
			}
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		for _, v := range []testStruct{{"2", 2}} {
			if mathz.IsPow10Int64(v.num) {
				t.Errorf("%s should not be pow10(x)", v.n)
			}
		}
	})
}

func TestIsPow10Uint32(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		n   string
		num uint32
	}
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		for _, v := range []testStruct{{"1", 1}, {"10", 10}, {"100", 100}, {"1000", 1000}, {"10000", 10000}} {
			if !mathz.IsPow10Uint32(v.num) {
				t.Errorf("%s should be pow10(x)", v.n)
			}
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		for _, v := range []testStruct{{"2", 2}} {
			if mathz.IsPow10Uint32(v.num) {
				t.Errorf("%s should not be pow10(x)", v.n)
			}
		}
	})
}

func TestIsPow10Uint64(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		n   string
		num uint64
	}
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		for _, v := range []testStruct{{"1", 1}, {"10", 10}, {"100", 100}, {"1000", 1000}, {"10000", 10000}} {
			if !mathz.IsPow10Uint64(v.num) {
				t.Errorf("%s should be pow10(x)", v.n)
			}
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		for _, v := range []testStruct{{"2", 2}} {
			if mathz.IsPow10Uint64(v.num) {
				t.Errorf("%s should not be pow10(x)", v.n)
			}
		}
	})
}
