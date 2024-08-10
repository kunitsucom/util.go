package discard_test

import (
	"io"
	"testing"

	"github.com/kunitsucom/util.go/discard"
)

func TestDiscard(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		discard.Discard(io.EOF)
	})
}

func TestOne(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		const e1 = 1
		a1 := discard.One(e1, io.EOF)
		if e1 != a1 {
			t.Errorf("❌: discard.One: e1 != a1")
		}
	})
}

func TestTwo(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		const (
			e1 = 1
			e2 = 2
		)
		a1, a2 := discard.Two(e1, e2, io.EOF)
		if e1 != a1 {
			t.Errorf("❌: discard.Two: e1 != a1")
		}
		if e2 != a2 {
			t.Errorf("❌: discard.Two: e2 != a2")
		}
	})
}

func TestThree(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		const (
			e1 = 1
			e2 = 2
			e3 = 3
		)
		a1, a2, a3 := discard.Three(e1, e2, e3, io.EOF)
		if e1 != a1 {
			t.Errorf("❌: discard.Three: e1 != a1")
		}
		if e2 != a2 {
			t.Errorf("❌: discard.Three: e2 != a2")
		}
		if e3 != a3 {
			t.Errorf("❌: discard.Three: e3 != a3")
		}
	})
}

func TestFour(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		const (
			e1 = 1
			e2 = 2
			e3 = 3
			e4 = 4
		)
		a1, a2, a3, a4 := discard.Four(e1, e2, e3, e4, io.EOF)
		if e1 != a1 {
			t.Errorf("❌: discard.Four: e1 != a1")
		}
		if e2 != a2 {
			t.Errorf("❌: discard.Four: e2 != a2")
		}
		if e3 != a3 {
			t.Errorf("❌: discard.Four: e3 != a3")
		}
		if e4 != a4 {
			t.Errorf("❌: discard.Four: e4 != a4")
		}
	})
}

func TestFive(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		const (
			e1 = 1
			e2 = 2
			e3 = 3
			e4 = 4
			e5 = 5
		)
		a1, a2, a3, a4, a5 := discard.Five(e1, e2, e3, e4, e5, io.EOF)
		if e1 != a1 {
			t.Errorf("❌: discard.Five: e1 != a1")
		}
		if e2 != a2 {
			t.Errorf("❌: discard.Five: e2 != a2")
		}
		if e3 != a3 {
			t.Errorf("❌: discard.Five: e3 != a3")
		}
		if e4 != a4 {
			t.Errorf("❌: discard.Five: e4 != a4")
		}
		if e5 != a5 {
			t.Errorf("❌: discard.Five: e5 != a5")
		}
	})
}
