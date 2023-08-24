package timez_test

import (
	"context"
	"testing"
	"time"

	timez "github.com/kunitsucom/util.go/time"
)

func TestNow(t *testing.T) {
	t.Parallel()
	t.Run("success,", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		constant := time.Date(2023, 8, 24, 21, 32, 0, 0, time.UTC)
		ctx = timez.WithContext(ctx, constant)

		actual := timez.Now(ctx)
		if constant != actual {
			t.Errorf("❌: constant(%v) != actual(%v)", constant, actual)
		}

		backup := timez.SetNowFunc(func(_ context.Context) time.Time { return time.Now() })
		current := timez.Now(context.Background())
		if current.After(constant) {
			t.Errorf("❌: current(%v) after constant(%v)", current, constant)
		}

		timez.SetNowFunc(backup)
		current2 := timez.Now(context.Background())
		if !current2.After(current) {
			t.Errorf("❌: current2(%v) after current(%v)", current2, current)
		}
	})
}
