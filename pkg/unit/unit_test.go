package unit_test

import (
	"testing"

	"github.com/kunitsuinc/util.go/pkg/unit"
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

func TestKiTo(t *testing.T) {
	t.Parallel()
	var expect uint64 = 1 << 10
	actual := unit.KiTo(1)
	if expect != actual {
		t.Errorf("expect != actual: %v != %v", expect, actual)
	}
}

func TestMiTo(t *testing.T) {
	t.Parallel()
	var expect uint64 = 1 << 20
	actual := unit.MiTo(1)
	if expect != actual {
		t.Errorf("expect != actual: %v != %v", expect, actual)
	}
}

func TestGiTo(t *testing.T) {
	t.Parallel()
	var expect uint64 = 1 << 30
	actual := unit.GiTo(1)
	if expect != actual {
		t.Errorf("expect != actual: %v != %v", expect, actual)
	}
}

func TestTiTo(t *testing.T) {
	t.Parallel()
	var expect uint64 = 1 << 40
	actual := unit.TiTo(1)
	if expect != actual {
		t.Errorf("expect != actual: %v != %v", expect, actual)
	}
}

func TestPiTo(t *testing.T) {
	t.Parallel()
	var expect uint64 = 1 << 50
	actual := unit.PiTo(1)
	if expect != actual {
		t.Errorf("expect != actual: %v != %v", expect, actual)
	}
}

func TestEiTo(t *testing.T) {
	t.Parallel()
	var expect uint64 = 1 << 60
	actual := unit.EiTo(1)
	if expect != actual {
		t.Errorf("expect != actual: %v != %v", expect, actual)
	}
}
