package errorz

import (
	"fmt"
	"io"
	"reflect"
	"runtime"
	"strings"
)

type (
	ErrorfOption interface {
		apply(c *errorfConfig)
	}
	errorfConfig struct {
		callerSkip int
	}
)

type callerSkipOption int

func (o callerSkipOption) apply(c *errorfConfig) {
	c.callerSkip = int(o)
}

// WithCallerSkip returns an ErrorfOption that sets the number of stack frames to skip.
func WithCallerSkip(callerSkip int) ErrorfOption {
	return callerSkipOption(callerSkip)
}

// NewErrorf returns a function like xerrors.Errorf.
// It is possible to return a function with different behaviors by passing ErrorfOption as arguments.
func NewErrorf(opts ...ErrorfOption) func(format string, a ...interface{}) error {
	c := &errorfConfig{}

	for _, opt := range opts {
		opt.apply(c)
	}

	return newErrorf(c)
}

const (
	indent4 = "    "
	ln      = "\n"
)

//nolint:cyclop
func newErrorf(c *errorfConfig) func(format string, a ...interface{}) error {
	return func(format string, a ...interface{}) error {
		const (
			suffixS      = ": %s"
			suffixV      = ": %v"
			suffixPlusV  = ": %+v"
			suffixSharpV = ": %#v"
			suffixW      = ": %w"
		)
		var (
			hasSuffixS      = strings.HasSuffix(format, suffixS)
			hasSuffixV      = strings.HasSuffix(format, suffixV)
			hasSuffixPlusV  = strings.HasSuffix(format, suffixPlusV)
			hasSuffixSharpV = strings.HasSuffix(format, suffixSharpV)
			hasSuffixW      = strings.HasSuffix(format, suffixW)
		)

		if !hasSuffixS && !hasSuffixV && !hasSuffixPlusV && !hasSuffixSharpV && !hasSuffixW {
			return fmt.Errorf(format, a...) //nolint:goerr113
		}

		prefix := format[:len(format)-len(suffixW)]
		suffix := format[len(format)-len(suffixW):]
		head := a[:len(a)-1]
		tail := a[len(a)-1]

		var e wrapError
		runtime.Callers(1+c.callerSkip, e.frame[:])
		e.msg = fmt.Sprintf(prefix, head...)
		switch err := tail.(type) {
		case formatter:
			e.err = err
		case error:
			switch {
			case hasSuffixS:
				e.err = fmt.Errorf("%s", err) //nolint:errorlint,goerr113 // for compatibility with xerrors.Errorf
			case hasSuffixV || hasSuffixPlusV || hasSuffixSharpV:
				e.err = fmt.Errorf("%v", err) //nolint:errorlint,goerr113 // for compatibility with xerrors.Errorf
			// case hasSuffixPlusV: // FIXME: support %+v
			// 	e.err = fmt.Errorf("%+v", err) //nolint:errorlint,goerr113 // for compatibility with xerrors.Errorf
			// case hasSuffixSharpV: // FIXME: support %#v
			// 	e.err = fmt.Errorf("%+v", err) //nolint:errorlint,goerr113 // for compatibility with xerrors.Errorf
			case hasSuffixW:
				e.err = err
			}
		default:
			e.msg += fmt.Sprintf(suffix, tail)
			e.err = nil
		}

		return &e
	}
}

//nolint:gochecknoglobals
var errorf = NewErrorf(WithCallerSkip(1))

// Errorf is a function like xerrors.Errorf.
func Errorf(format string, a ...interface{}) error {
	return errorf(format, a...)
}

type wrapError struct {
	msg   string
	err   error
	frame [3]uintptr // See: https://go.googlesource.com/go/+/032678e0fb/src/runtime/extern.go#169
}

var (
	_ error                       = (*wrapError)(nil)
	_ formatter                   = (*wrapError)(nil)
	_ fmt.Formatter               = (*wrapError)(nil)
	_ fmt.GoStringer              = (*wrapError)(nil)
	_ interface{ Unwrap() error } = (*wrapError)(nil)
)

type formatter interface {
	error
	format(s fmt.State, verb rune)
	Unwrap() error
}

func (e *wrapError) writeCallers(w io.Writer) {
	frames := runtime.CallersFrames(e.frame[:])
	if _, ok := frames.Next(); !ok {
		return
	}
	target, ok := frames.Next()
	if !ok {
		return
	}

	if target.Function != "" {
		fmt.Fprintf(w, ":"+ln+indent4+"%s", target.Function)
		// NOTE:
		//              ^^^^^^^^^^^^^^^^^
		//              means a part of stacktrace:
		//
		// funcA:\n
		//      ^^^
		//     github.com/org/repo/pkg.funcA
		// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
		if target.File != "" {
			fmt.Fprintf(w, ln+indent4+indent4+"%s:%d", target.File, target.Line)
			// NOTE:
			//             ^^^^^^^^^^^^^^^^^^^^^^^^^
			//             means a part of stacktrace:
			//
			//     github.com/org/repo/pkg.funcA\n
			//                                  ^^
			//         github.com/org/repo/pkg.go:123
			// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
		}
	}
}

func (e *wrapError) Error() string {
	return fmt.Sprint(e) //nolint:perfsprint
}

func (e *wrapError) Format(s fmt.State, verb rune) {
	var err error = e
loop:
	for {
		switch fe := err.(type) { //nolint:errorlint
		case formatter:
			fe.format(s, verb)
			err = fe.Unwrap()
		case fmt.Formatter:
			fe.Format(s, verb)
			break loop
		default:
			_, _ = fmt.Fprintf(s, fmt.FormatString(s, verb), fe)
			break loop
		}
		if err == nil {
			break loop
		}
	}
}

func (e *wrapError) format(s fmt.State, verb rune) {
	var withStacktrace bool
Verb:
	switch verb {
	// FormatError() will not be called with the 'w' verb.
	// case 'w':
	case 'v':
		switch {
		case s.Flag('#'):
			_, _ = io.WriteString(s, e.GoString())
			return
		case s.Flag('+'):
			withStacktrace = true
			break Verb
		}
	default:
	}

	_, _ = io.WriteString(s, e.msg)
	if withStacktrace {
		e.writeCallers(s)
		if e.err != nil {
			_, _ = io.WriteString(s, ln+"  - ")
			// NOTE:
			//                        ^^^^^^
			//                        means a part of stacktrace:
			//
			//         github.com/org/repo/pkg.go:123\n
			//                                       ^^
			//   - funcB:
			// ^^^^
		}
	} else { //nolint:gocritic
		if e.err != nil {
			_, _ = io.WriteString(s, ": ")
			// NOTE:
			//                        ^^
			//                        means a part of error output:
			// funcA: funcB: funcC: error
			//      ^^     ^^     ^^
		}
	}
}

func (e *wrapError) GoString() string {
	typ := reflect.TypeOf(*e)
	val := reflect.ValueOf(*e)
	elems := make([]string, typ.NumField())
	for i := 0; typ.NumField() > i; i++ {
		elems[i] = fmt.Sprintf("%s:%#v", typ.Field(i).Name, val.Field(i))
	}
	return fmt.Sprintf("&%s{%s}", typ, strings.Join(elems, ", "))
}

func (e *wrapError) Unwrap() error {
	return e.err
}

// FormatError is intended to be used as follows:
//
//	func (e *customError) Format(s fmt.State, verb rune) {
//		errorz.FormatError(s, verb, e.Unwrap())
//	}
func FormatError(s fmt.State, verb rune, err error) {
	if formatter, ok := err.(fmt.Formatter); ok { //nolint:errorlint
		formatter.Format(s, verb)
		return
	}

	_, _ = fmt.Fprintf(s, fmt.FormatString(s, verb), err)
}
