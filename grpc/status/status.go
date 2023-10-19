package statz

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/runtime/protoiface"

	errorz "github.com/kunitsucom/util.go/errors"
	filepathz "github.com/kunitsucom/util.go/path/filepath"
)

type Level string

const (
	DebugLevel Level = "DEBUG"
	InfoLevel  Level = "INFO"
	WarnLevel  Level = "WARN"
	ErrorLevel Level = "ERROR"
)

type statusError struct {
	error
	s      *status.Status
	level  Level
	logger Logger
}

type (
	ErrorOption interface {
		apply(*statusError)
	}
	Logger           = func(ctx context.Context, level Level, stat *status.Status, err error)
	newOptionLogger  struct{ handler Logger }
	newOptionDetails struct{ details []protoiface.MessageV1 }
)

func (o *newOptionLogger) apply(e *statusError) { e.logger = o.handler }

func WithLogger(handler Logger) ErrorOption { //nolint:ireturn
	return &newOptionLogger{handler: handler}
}

func (o *newOptionDetails) apply(e *statusError) {
	statusWithDetails, err := e.s.WithDetails(o.details...)
	if err != nil {
		err = errorz.NewErrorf(errorz.WithCallerSkip(2))("e.s.WithDetails: details=%v: %w", o.details, err)
		e.logger(context.Background(), ErrorLevel, status.New(codes.Unknown, "WithDetails failed"), err)
		return
	}

	e.s = statusWithDetails
}

func WithDetails(details ...protoiface.MessageV1) ErrorOption { //nolint:ireturn
	return &newOptionDetails{details: details}
}

// DiscardLogger is a logger to discard logs. If you want to disable logging, assign DiscardLogger to DefaultLogger.
func DiscardLogger(_ context.Context, _ Level, _ *status.Status, _ error) {}

//nolint:gochecknoglobals
var (
	// DefaultLogger is a logger to record the location of statz.New() calls.
	DefaultLogger Logger = func(_ context.Context, level Level, stat *status.Status, err error) {
		_, file, line, _ := runtime.Caller(2)
		log.Printf("level=%s caller=%s:%d code=%s message=%q details=%s error=%q stacktrace=%q", level, filepathz.Short(file), line, stat.Code(), stat.Message(), stat.Details(), fmt.Sprintf("%v", err), fmt.Sprintf("%+v", err))
	}

	errorf = errorz.NewErrorf(errorz.WithCallerSkip(1))
)

var (
	_ error                                    = (*statusError)(nil)
	_ fmt.Formatter                            = (*statusError)(nil)
	_ interface{ Unwrap() error }              = (*statusError)(nil)
	_ interface{ GRPCStatus() *status.Status } = (*statusError)(nil)
)

// New is a function like errors.New() for gRPC status.Status.
//
// The location where statz.New() is called is logged by DefaultLogger.
func New(ctx context.Context, level Level, code codes.Code, msg string, err error, opts ...ErrorOption) error {
	e := &statusError{
		error:  errorf("statz.New: level=%s code=%s message=%s: %w", level, code, msg, err),
		s:      status.New(code, msg),
		level:  level,
		logger: DefaultLogger,
	}
	for _, opt := range opts {
		opt.apply(e)
	}
	e.logger(ctx, e.level, e.s, err)
	return e
}

func (e *statusError) Format(s fmt.State, verb rune) {
	errorz.FormatError(s, verb, e.Unwrap())
}

func (e *statusError) Unwrap() error {
	return e.error
}

func (e *statusError) GRPCStatus() *status.Status {
	return e.s
}
