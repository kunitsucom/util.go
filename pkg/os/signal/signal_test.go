package signalz_test

import (
	"os"
	"testing"

	signalz "github.com/kunitsuinc/util.go/pkg/os/signal"
)

func TestNotify(t *testing.T) {
	t.Parallel()
	source := make(chan os.Signal, 1)
	c := signalz.Notify(source, os.Interrupt)
	close(c)
}
