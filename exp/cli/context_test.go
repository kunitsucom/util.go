package cliz_test

import (
	"context"
	"errors"
	"testing"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

func TestFromContext(t *testing.T) {
	t.Parallel()

	t.Run("failure,ErrNilContext", func(t *testing.T) {
		t.Parallel()

		if _, err := cliz.FromContext(nil); !errors.Is(err, cliz.ErrNilContext) {
			t.Errorf("❌: err != cliz.ErrNilContext: %+v", err)
		}
	})

	t.Run("failure,ErrCommandNotSetInContext", func(t *testing.T) {
		t.Parallel()

		if _, err := cliz.FromContext(context.Background()); !errors.Is(err, cliz.ErrCommandNotSetInContext) {
			t.Errorf("❌: err != cliz.ErrCommandNotSetInContext: %+v", err)
		}
	})
}

func TestMustFromContext(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := cliz.WithContext(context.Background(), &cliz.Command{})
		if c := cliz.MustFromContext(ctx); c == nil {
			t.Errorf("❌: c == nil")
		}
	})

	t.Run("failure,panic", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("❌: panic did not occur")

				if err, ok := r.(error); ok {
					t.Errorf("❌: err: %+v", err)

					if !errors.Is(err, cliz.ErrNilContext) {
						t.Errorf("❌: err != cliz.ErrNilContext: %+v", err)
					}
				}
			}
		}()

		cliz.MustFromContext(nil)
	})
}
