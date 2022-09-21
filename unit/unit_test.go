package unit_test

import (
	"testing"

	"github.com/kunitsuinc/util.go/unit"
)

func TestKi(t *testing.T) {
	t.Parallel()
	var expect uint64 = 1
	actual := unit.ToKi(1 << 10)
	if expect != actual {
		t.Errorf("expect != actual: %v != %v", expect, actual)
	}
}

func TestMi(t *testing.T) {
	t.Parallel()
	var expect uint64 = 1
	actual := unit.ToMi(1 << 20)
	if expect != actual {
		t.Errorf("expect != actual: %v != %v", expect, actual)
	}
}

func TestGi(t *testing.T) {
	t.Parallel()
	var expect uint64 = 1
	actual := unit.ToGi(1 << 30)
	if expect != actual {
		t.Errorf("expect != actual: %v != %v", expect, actual)
	}
}

func TestTi(t *testing.T) {
	t.Parallel()
	var expect uint64 = 1
	actual := unit.ToTi(1 << 40)
	if expect != actual {
		t.Errorf("expect != actual: %v != %v", expect, actual)
	}
}

func TestPi(t *testing.T) {
	t.Parallel()
	var expect uint64 = 1
	actual := unit.ToPi(1 << 50)
	if expect != actual {
		t.Errorf("expect != actual: %v != %v", expect, actual)
	}
}

func TestEi(t *testing.T) {
	t.Parallel()
	var expect uint64 = 1
	actual := unit.ToEi(1 << 60)
	if expect != actual {
		t.Errorf("expect != actual: %v != %v", expect, actual)
	}
}
