package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type (
	Serve    func(errorChannel chan error)
	Shutdown func(ctx context.Context) error

	Server struct {
		serves                []Serve
		shutdown              Shutdown
		shutdownContext       context.Context //nolint:containedctx
		shutdownTimeout       time.Duration
		shutdownErrorHandler  func(err error)
		signalChannel         chan os.Signal
		continueSignalHandler func(sig os.Signal) bool
	}
	Option func(s *Server)
)

func WithShutdownContext(ctx context.Context) Option {
	return func(s *Server) { s.shutdownContext = ctx }
}

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) { s.shutdownTimeout = timeout }
}

func WithShutdownErrorHandler(shutdownErrorHandler func(err error)) Option {
	return func(s *Server) { s.shutdownErrorHandler = shutdownErrorHandler }
}

func WithSignalChannel(signalChannel chan os.Signal) Option {
	return func(s *Server) { s.signalChannel = signalChannel }
}

func WithContinueSignalHandler(continueSignalHandler func(sig os.Signal) bool) Option {
	return func(s *Server) { s.continueSignalHandler = continueSignalHandler }
}

func New(serves []Serve, shutdown Shutdown, opts ...Option) *Server {
	s := &Server{
		serves:                serves,
		shutdown:              shutdown,
		shutdownContext:       context.Background(),
		shutdownTimeout:       10 * time.Second,
		shutdownErrorHandler:  func(err error) { log.Println("shutdown error:", err) },
		signalChannel:         make(chan os.Signal, 1),
		continueSignalHandler: func(sig os.Signal) bool { log.Println("cache the signal:", sig); return sig == syscall.SIGHUP },
	}
	signal.Notify(s.signalChannel, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM)

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Serve(ctx context.Context) error {
	errorChannel := make(chan error, len(s.serves))

	for _, serve := range s.serves {
		go serve(errorChannel)
	}

	for {
		select {
		case <-ctx.Done():
			ctx, cancel := context.WithTimeout(s.shutdownContext, s.shutdownTimeout)
			if err := s.shutdown(ctx); err != nil { //nolint:contextcheck
				s.shutdownErrorHandler(err)
			}
			cancel()
		case sig := <-s.signalChannel:
			if s.continueSignalHandler(sig) {
				continue
			}
			ctx, cancel := context.WithTimeout(s.shutdownContext, s.shutdownTimeout)
			if err := s.shutdown(ctx); err != nil { //nolint:contextcheck
				s.shutdownErrorHandler(err)
			}
			cancel()
		case err := <-errorChannel:
			if err != nil {
				return err
			}
			return nil
		}
	}
}
