package statz_test

import (
	"context"
	"io"
	"testing"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"

	errorz "github.com/kunitsucom/util.go/errors"
	statz "github.com/kunitsucom/util.go/grpc/status"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		before := errorz.Errorf("err: %w", io.ErrUnexpectedEOF)
		err := statz.New(ctx, statz.ErrorLevel, codes.Internal, "my error message", before, statz.WithDetails(&errdetails.DebugInfo{StackEntries: []string{"stack1", "stack2"}}), statz.WithDetails(nil), statz.WithLogger(statz.DefaultLogger))
		t.Logf("🪲: [%%v]:\n%v", err)
		t.Logf("🪲: [%%+v]:\n%+v", err)
		_ = statz.New(ctx, statz.ErrorLevel, codes.Internal, "my error message", before, statz.WithLogger(statz.DiscardLogger))
	})
}
