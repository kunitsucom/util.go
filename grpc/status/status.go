package statz

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/runtime/protoiface"

	errorz "github.com/kunitsucom/util.go/errors"
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

//nolint:gochecknoglobals
var (
	DiscardLogger Logger = func(_ context.Context, _ Level, _ *status.Status, _ error) {}
	DefaultLogger Logger = func(_ context.Context, level Level, stat *status.Status, err error) {
		log.Printf("level=%s code=%s message=%q details=%s error=%q stacktrace=%q", level, stat.Code(), stat.Message(), stat.Details(), fmt.Sprintf("%v", err), fmt.Sprintf("%+v", err))
	}
	_      interface{ GRPCStatus() *status.Status } = (*statusError)(nil)
	errorf                                          = errorz.NewErrorf(errorz.WithCallerSkip(1))
)

func New(ctx context.Context, level Level, code codes.Code, msg string, err error, opts ...ErrorOption) error {
	e := &statusError{
		error:  errorf("errgrpc.New: level=%s code=%s message=%s: %w", level, code, msg, err),
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
	if formatter, ok := e.error.(fmt.Formatter); ok { //nolint:errorlint
		formatter.Format(s, verb)
		return
	}

	_, _ = fmt.Fprintf(s, fmt.FormatString(s, verb), e.error)
}

func (e *statusError) GRPCStatus() *status.Status {
	return e.s
}
