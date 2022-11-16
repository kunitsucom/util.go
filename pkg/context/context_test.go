package contextz_test

import (
	"context"
	"errors"
	"os"
	"reflect"
	"testing"

	contextz "github.com/kunitsuinc/util.go/pkg/context"
)

func TestWithSignalChannel(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		expect := make(chan os.Signal, 1)

		ctx = contextz.WithSignalChannel(ctx, expect)

		actual := contextz.MustSignalChannel(ctx)

		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("expect != actual: %v", actual)
		}
	})

	t.Run("failure()", func(t *testing.T) {
		t.Parallel()

		defer func() {
			recov := recover()

			if recov == nil {
				t.Errorf("recov == nil")
			}

			err, ok := recov.(error)
			if !ok {
				t.Errorf("!ok")
			}

			if !errors.Is(err, contextz.ErrValueNotSet) {
				t.Errorf("err != contextz.ErrValueNotSet: %v", err)
			}
		}()

		_ = contextz.MustSignalChannel(context.Background())
	})
}
