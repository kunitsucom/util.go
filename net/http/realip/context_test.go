package realip_test

import (
	"context"
	"testing"

	"github.com/kunitsuinc/util.go/net/http/realip"
)

func TestContextXRealIP(t *testing.T) {
	t.Parallel()
	expect := ""
	actual := realip.ContextXRealIP(context.Background())
	if expect != actual {
		t.Errorf("expect != actual: %s", actual)
	}
}
