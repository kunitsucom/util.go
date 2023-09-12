package statz

import (
	"context"
	"errors"
	"io"
	"testing"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errorz "github.com/kunitsucom/util.go/errors"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		before := errorz.Errorf("err: %w", io.ErrUnexpectedEOF)
		err := New(ctx, ErrorLevel, codes.Internal, "my error message", before, WithDetails(&errdetails.DebugInfo{StackEntries: []string{"stack1", "stack2"}}), WithDetails(nil), WithLogger(DefaultLogger))
		t.Logf("🪲: [%%v]:\n%v", err)
		t.Logf("🪲: [%%+v]:\n%+v", err)

		var e *statusError
		if !errors.As(err, &e) {
			t.Errorf("❌: something wrong")
		}
		e.error = io.ErrUnexpectedEOF
		t.Logf("🪲: [%%v]:\n%v", e)

		s, ok := New(ctx, ErrorLevel, codes.Internal, "my error message", io.ErrUnexpectedEOF, WithLogger(DiscardLogger)).(interface{ GRPCStatus() *status.Status }) //nolint:errorlint
		if !ok {
			t.Errorf("❌: New: not implement interface{ GRPCStatus() *status.Status }")
		}
		t.Logf("🪲: [%%v]:\n%v", s.GRPCStatus())
	})
}
