// nolint: testpackage
package httpz

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

type TestHTTPServer struct {
	Err error
}

func (s *TestHTTPServer) ListenAndServe() error {
	return s.Err
}

func TestListenAndServe(t *testing.T) {
	t.Parallel()
	t.Run("success(0)", func(t *testing.T) {
		t.Parallel()
		const (
			portrangeFirst = 49152
			portrangeLast  = 65535
		)
		// nolint: gosec
		port := rand.Intn(portrangeLast-portrangeFirst+1) + portrangeFirst
		s := &http.Server{
			Addr:              fmt.Sprintf("127.0.0.1:%d", port),
			ReadHeaderTimeout: 100 * time.Millisecond,
		}
		go func() {
			time.Sleep(200 * time.Millisecond)
			_ = s.Shutdown(context.Background())
		}()
		if err := ListenAndServe(s); err != nil {
			t.Errorf("err != nil: %v", err)
		}
	})

	t.Run("success(1)", func(t *testing.T) {
		t.Parallel()
		s := &TestHTTPServer{Err: http.ErrServerClosed}
		if err := listenAndServe(s); err != nil {
			t.Errorf("err != nil: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		s := &TestHTTPServer{io.ErrUnexpectedEOF}
		if err := listenAndServe(s); err == nil {
			t.Errorf("err == nil")
		}
	})
}
