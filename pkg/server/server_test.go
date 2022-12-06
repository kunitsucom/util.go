package server_test

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/kunitsuinc/util.go/pkg/server"
)

func TestServer_Serve(t *testing.T) {
	t.Parallel()

	t.Run("success(signal)", func(t *testing.T) {
		t.Parallel()
		signalChannel := make(chan os.Signal, 1)
		testMu, testNotCalled := sync.Mutex{}, uint64(0)
		testChan := make(chan error, 1)
		s := server.New(
			[]server.Serve{func(errorChannel chan error) {
				errorChannel <- <-testChan
				log.Printf("shutdown: %s", t.Name())
			}},
			func(ctx context.Context) error {
				testMu.Lock()
				defer testMu.Unlock()
				if atomic.LoadUint64(&testNotCalled) == 0 {
					atomic.AddUint64(&testNotCalled, 1)
					log.Printf("starting shutdown: %s", t.Name())
					testChan <- nil
				}
				return nil
			},
			server.WithSignalChannel(signalChannel),
		)
		go func() {
			signalChannel <- syscall.SIGHUP
			signalChannel <- syscall.SIGINT
		}()
		err := s.Serve(context.Background())
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
	})

	t.Run("failure(signal)", func(t *testing.T) {
		t.Parallel()
		signalChannel := make(chan os.Signal, 1)
		testMu, testNotCalled := sync.Mutex{}, uint64(0)
		testChan := make(chan error, 1)
		s := server.New(
			[]server.Serve{func(errorChannel chan error) {
				<-testChan
				errorChannel <- http.ErrServerClosed
				log.Printf("shutdown: %s", t.Name())
			}},
			func(ctx context.Context) error {
				testMu.Lock()
				defer testMu.Unlock()
				if atomic.LoadUint64(&testNotCalled) == 0 {
					atomic.AddUint64(&testNotCalled, 1)
					log.Printf("starting shutdown: %s", t.Name())
					testChan <- nil
				}
				return http.ErrServerClosed
			},
			server.WithSignalChannel(signalChannel),
			server.WithContinueSignalHandler(func(sig os.Signal) bool { log.Println("cache the signal:", sig); return sig == syscall.SIGHUP }),
		)
		go func() {
			signalChannel <- syscall.SIGHUP
			signalChannel <- syscall.SIGINT
		}()
		err := s.Serve(context.Background())
		if err == nil {
			t.Errorf("err == nil")
		}
	})

	t.Run("success(ctx)", func(t *testing.T) {
		t.Parallel()
		testMu, testNotCalled := sync.Mutex{}, uint64(0)
		testChan := make(chan error, 1)
		s := server.New(
			[]server.Serve{func(errorChannel chan error) {
				errorChannel <- <-testChan
				log.Printf("shutdown: %s", t.Name())
			}},
			func(ctx context.Context) error {
				testMu.Lock()
				defer testMu.Unlock()
				if atomic.LoadUint64(&testNotCalled) == 0 {
					atomic.AddUint64(&testNotCalled, 1)
					log.Printf("starting shutdown: %s", t.Name())
					testChan <- nil
				}
				return nil
			},
			server.WithShutdownContext(context.Background()),
			server.WithShutdownTimeout(10*time.Second),
		)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := s.Serve(ctx)
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
	})

	t.Run("failure(ctx)", func(t *testing.T) {
		t.Parallel()
		testMu, testNotCalled := sync.Mutex{}, uint64(0)
		testChan := make(chan error, 1)
		s := server.New(
			[]server.Serve{func(errorChannel chan error) {
				errorChannel <- <-testChan
			}},
			func(ctx context.Context) error {
				testMu.Lock()
				defer testMu.Unlock()
				if atomic.LoadUint64(&testNotCalled) == 0 {
					atomic.AddUint64(&testNotCalled, 1)
					log.Printf("starting shutdown: %s", t.Name())
					testChan <- nil
				}
				return http.ErrServerClosed
			},
			server.WithShutdownErrorHandler(func(err error) { log.Println("shutdown error:", err) }),
		)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := s.Serve(ctx)
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
	})
}
