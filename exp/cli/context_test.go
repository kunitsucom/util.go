package cliz_test

import (
	"context"
	"errors"
	"testing"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

func TestFromContext(t *testing.T) {
	t.Parallel()

	t.Run("failure,ErrCommandNotSetInContext", func(t *testing.T) {
		t.Parallel()

		if _, err := cliz.FromContext(context.Background()); !errors.Is(err, cliz.ErrCommandNotSetInContext) {
			t.Errorf("❌: err != cliz.ErrCommandNotSetInContext: %+v", err)
		}
	})
}
