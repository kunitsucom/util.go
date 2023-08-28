package errorz

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"runtime"
	"strings"
	"testing"

	// "github.com/kunitsucom/ilog.go"
	// "go.uber.org/zap"
	// "go.uber.org/zap/zapcore"
	// "golang.org/x/xerrors".

	stringz "github.com/kunitsucom/util.go/strings"
	testingz "github.com/kunitsucom/util.go/testing"
)

const (
	_indent        = "    "
	_stackTail     = `:\n` + _indent + `[^ ]*github.com/kunitsucom/util.go/errors.case[0-9]+[^ ]*\n` + _indent + _indent + `[^ ]*github.com/kunitsucom/util.go/errors/fmt_test.go:[0-9]+`
	_stack         = _stackTail + `\n  - `
	_nilErrW       = "%!w(<nil>)"
	_nilErrWEscape = `%!w\(<nil>\)`
	_nilErrS       = "%!s(<nil>)"
	_nilErrSEscape = `%!s\(<nil>\)`
	_nilErrV       = "<nil>"
	_nilErrVEscape = "<nil>"
)

//nolint:thelper // because this is a test func for stacktrace.
func case1(t *testing.T, errorfFunc func(format string, a ...interface{}) error, orig error, verb string, nilErr string, nilErrEscape string, compare func(actual error, expect error) bool) {
	bufS := bytes.NewBuffer(nil)
	bufV := bytes.NewBuffer(nil)
	bufPlusV := bytes.NewBuffer(nil)
	case1FuncA := func() error { return orig }
	case1FuncB := func() error { return errorfFunc("case1FuncA: %"+verb, case1FuncA()) }
	case1FuncC := func() error { return errorfFunc("case1FuncB: %"+verb, case1FuncB()) }
	case1FuncD := func() error { return errorfFunc("case1FuncC: %"+verb, case1FuncC()) }
	actual := errorfFunc("case1FuncD: %"+verb, case1FuncD())
	expectS := `case1FuncD: case1FuncC: case1FuncB: case1FuncA: ` + nilErr
	expectV := `case1FuncD: case1FuncC: case1FuncB: case1FuncA: ` + nilErr
	expectPlusV := regexp.MustCompile(stringz.Join("", `case1FuncD`, _stack, `case1FuncC`, _stack, `case1FuncB`, _stack, `case1FuncA: `+nilErrEscape, _stackTail))
	if orig != nil {
		errStr := orig.Error()
		expectS = `case1FuncD: case1FuncC: case1FuncB: case1FuncA: ` + errStr
		expectV = `case1FuncD: case1FuncC: case1FuncB: case1FuncA: ` + errStr
		expectPlusV = regexp.MustCompile(stringz.Join("", `case1FuncD`, _stack, `case1FuncC`, _stack, `case1FuncB`, _stack, `case1FuncA`, _stack, errStr))
	}
	_, _ = fmt.Fprintf(bufS, "%s", actual)
	_, _ = fmt.Fprintf(bufV, "%v", actual)
	_, _ = fmt.Fprintf(bufPlusV, "%+v", actual)

	if expect := orig; !HasSuffix(actual, nilErr) && !compare(actual, expect) {
		t.Errorf("❌: [%s]: compare:\n[EXPECT]:\n%v\n[ACTUAL]:\n%v\n", t.Name(), expect, actual)
	}

	if expect, actual := expectS, bufS.String(); expect != actual {
		t.Errorf("❌: [%s]: [%%s]:\n[EXPECT]:\n%v\n[ACTUAL]:\n%v\n", t.Name(), expect, actual)
	}

	if expect, actual := expectV, bufV.String(); expect != actual {
		t.Errorf("❌: [%s]: [%%v]:\n[EXPECT]:\n%v\n[ACTUAL]:\n%v\n", t.Name(), expect, actual)
	}

	if expect, actual := expectPlusV, bufPlusV.String(); !expect.MatchString(actual) {
		t.Errorf("❌: [%s]: [%%+v]:\n[EXPECT]:\n%v\n[ACTUAL]:\n%v\n", t.Name(), strings.ReplaceAll(expect.String(), "\\n", "\n"), actual)
	}

	t.Logf("ℹ️: [%s]:\nError():\n%s\n", t.Name(), actual.Error())
	t.Logf("ℹ️: [%s]:\n[%%s]:\n%s\n", t.Name(), bufS)
	t.Logf("ℹ️: [%s]:\n[%%v]:\n%s\n", t.Name(), bufV)
	t.Logf("ℹ️: [%s]:\n[%%+v]:\n%s\n", t.Name(), bufPlusV)
}

//nolint:paralleltest
func TestErrorf(t *testing.T) {
	formatterError := &testingz.FormatterError{ErrorFunc: func() string { return io.ErrUnexpectedEOF.Error() }, FormatFunc: func(s fmt.State, verb rune) { fmt.Fprintf(s, "%"+string(verb), io.ErrUnexpectedEOF) }}
	//
	// Errorf
	//
	t.Run("success,case1,io.ErrUnexpectedEOF,w,errorz.Errorf", func(t *testing.T) {
		case1(t, Errorf, io.ErrUnexpectedEOF, "w", _nilErrW, _nilErrWEscape, errors.Is)
	})
	t.Run("success,case1,formatterError,w,errorz.Errorf", func(t *testing.T) {
		case1(t, Errorf, formatterError, "w", _nilErrW, _nilErrWEscape, errors.Is)
	})
	t.Run("success,case1,<nil>,w,errorz.Errorf", func(t *testing.T) {
		case1(t, Errorf, error(nil), "w", _nilErrW, _nilErrWEscape, errors.Is)
	})
	t.Run("success,case1,io.ErrUnexpectedEOF,s,errorz.Errorf", func(t *testing.T) {
		case1(t, Errorf, io.ErrUnexpectedEOF, "s", _nilErrS, _nilErrSEscape, func(actual, expect error) bool { return expect != nil && Contains(actual, expect.Error()) })
	})
	t.Run("success,case1,formatterError,s,errorz.Errorf", func(t *testing.T) {
		case1(t, Errorf, formatterError, "s", _nilErrS, _nilErrSEscape, func(actual, expect error) bool { return expect != nil && Contains(actual, expect.Error()) })
	})
	t.Run("success,case1,<nil>,s,errorz.Errorf", func(t *testing.T) {
		case1(t, Errorf, error(nil), "s", _nilErrS, _nilErrSEscape, func(actual, expect error) bool { return expect != nil && Contains(actual, expect.Error()) })
	})
	t.Run("success,case1,io.ErrUnexpectedEOF,v,errorz.Errorf", func(t *testing.T) {
		case1(t, Errorf, io.ErrUnexpectedEOF, "v", _nilErrV, _nilErrVEscape, func(actual, expect error) bool { return expect != nil && Contains(actual, expect.Error()) })
	})
	t.Run("success,case1,formatterError,v,errorz.Errorf", func(t *testing.T) {
		case1(t, Errorf, formatterError, "v", _nilErrV, _nilErrVEscape, func(actual, expect error) bool { return expect != nil && Contains(actual, expect.Error()) })
	})
	t.Run("success,case1,<nil>,v,errorz.Errorf", func(t *testing.T) {
		case1(t, Errorf, error(nil), "v", _nilErrV, _nilErrVEscape, func(actual, expect error) bool { return expect != nil && Contains(actual, expect.Error()) })
	})
	// TODO: support %+v
	// t.Run("success,case1,io.ErrUnexpectedEOF,+v,errorz.Errorf", func(t *testing.T) {
	// 	case1(t, Errorf, io.ErrUnexpectedEOF, "+v", _nilErrV, _nilErrVEscape, func(actual, expect error) bool { return expect != nil && Contains(actual, expect.Error()) })
	// })
	// t.Run("success,case1,formatterError,+v,errorz.Errorf", func(t *testing.T) {
	// 	case1(t, Errorf, formatterError, "+v", _nilErrV, _nilErrVEscape, func(actual, expect error) bool { return expect != nil && Contains(actual, expect.Error()) })
	// })
	// t.Run("success,case1,<nil>,+v,errorz.Errorf", func(t *testing.T) {
	// 	case1(t, Errorf, error(nil), "+v", _nilErrV, _nilErrVEscape, func(actual, expect error) bool { return expect != nil && Contains(actual, expect.Error()) })
	// })
	//
	// NewErrorf()
	//
	t.Run("success,case1,io.ErrUnexpectedEOF,w,errorz.NewErrorf()", func(t *testing.T) {
		case1(t, NewErrorf(), io.ErrUnexpectedEOF, "w", _nilErrW, _nilErrWEscape, errors.Is)
	})
	t.Run("success,case1,formatterError,w,errorz.NewErrorf()", func(t *testing.T) {
		case1(t, NewErrorf(), formatterError, "w", _nilErrW, _nilErrWEscape, errors.Is)
	})
	t.Run("success,case1,<nil>,w,errorz.NewErrorf()", func(t *testing.T) {
		case1(t, NewErrorf(), error(nil), "w", _nilErrW, _nilErrWEscape, errors.Is)
	})
	t.Run("success,case1,io.ErrUnexpectedEOF,s,errorz.NewErrorf()", func(t *testing.T) {
		case1(t, NewErrorf(), io.ErrUnexpectedEOF, "s", _nilErrS, _nilErrSEscape, func(actual, expect error) bool { return expect != nil && Contains(actual, expect.Error()) })
	})
	t.Run("success,case1,formatterError,s,errorz.NewErrorf()", func(t *testing.T) {
		case1(t, NewErrorf(), formatterError, "s", _nilErrS, _nilErrSEscape, func(actual, expect error) bool { return expect != nil && Contains(actual, expect.Error()) })
	})
	t.Run("success,case1,<nil>,s,errorz.NewErrorf()", func(t *testing.T) {
		case1(t, NewErrorf(), error(nil), "s", _nilErrS, _nilErrSEscape, func(actual, expect error) bool { return expect != nil && Contains(actual, expect.Error()) })
	})
	t.Run("success,case1,io.ErrUnexpectedEOF,v,errorz.NewErrorf()", func(t *testing.T) {
		case1(t, NewErrorf(), io.ErrUnexpectedEOF, "v", _nilErrV, _nilErrVEscape, func(actual, expect error) bool { return expect != nil && Contains(actual, expect.Error()) })
	})
	t.Run("success,case1,formatterError,v,errorz.NewErrorf()", func(t *testing.T) {
		case1(t, NewErrorf(), formatterError, "v", _nilErrV, _nilErrVEscape, func(actual, expect error) bool { return expect != nil && Contains(actual, expect.Error()) })
	})
	t.Run("success,case1,<nil>,v,errorz.NewErrorf()", func(t *testing.T) {
		case1(t, NewErrorf(), error(nil), "v", _nilErrV, _nilErrVEscape, func(actual, expect error) bool { return expect != nil && Contains(actual, expect.Error()) })
	})
	//
	// xerrors.Errorf
	//
	// t.Run("success,case1,io.ErrUnexpectedEOF,w,xerrors.Errorf", func(t *testing.T) {
	// 	case1(t, xerrors.Errorf, io.ErrUnexpectedEOF, "w", _nilErrW, _nilErrWEscape, errors.Is)
	// })
	// t.Run("success,case1,formatterError,w,xerrors.Errorf", func(t *testing.T) {
	// 	case1(t, xerrors.Errorf, formatterError, "w", _nilErrW, _nilErrWEscape, errors.Is)
	// })
	// t.Run("success,case1,<nil>,w,xerrors.Errorf", func(t *testing.T) {
	// 	case1(t, xerrors.Errorf, error(nil), "w", _nilErrW, _nilErrWEscape, errors.Is)
	// })
	// t.Run("success,case1,io.ErrUnexpectedEOF,s,xerrors.Errorf", func(t *testing.T) {
	// 	case1(t, xerrors.Errorf, io.ErrUnexpectedEOF, "s", _nilErrS, _nilErrSEscape, func(actual, expect error) bool { return expect != nil && Contains(actual, expect.Error()) })
	// })
	// t.Run("success,case1,formatterError,s,xerrors.Errorf", func(t *testing.T) {
	// 	case1(t, xerrors.Errorf, formatterError, "s", _nilErrS, _nilErrSEscape, func(actual, expect error) bool { return expect != nil && Contains(actual, expect.Error()) })
	// })
	// t.Run("success,case1,<nil>,s,xerrors.Errorf", func(t *testing.T) {
	// 	case1(t, xerrors.Errorf, error(nil), "s", _nilErrS, _nilErrSEscape, func(actual, expect error) bool { return expect != nil && Contains(actual, expect.Error()) })
	// })
	// t.Run("success,case1,io.ErrUnexpectedEOF,v,xerrors.Errorf", func(t *testing.T) {
	// 	case1(t, xerrors.Errorf, io.ErrUnexpectedEOF, "v", _nilErrV, _nilErrVEscape, func(actual, expect error) bool { return expect != nil && Contains(actual, expect.Error()) })
	// })
	// t.Run("success,case1,formatterError,v,xerrors.Errorf", func(t *testing.T) {
	// 	case1(t, xerrors.Errorf, formatterError, "v", _nilErrV, _nilErrVEscape, func(actual, expect error) bool { return expect != nil && Contains(actual, expect.Error()) })
	// })
	// t.Run("success,case1,<nil>,v,xerrors.Errorf", func(t *testing.T) {
	// 	case1(t, xerrors.Errorf, error(nil), "v", _nilErrV, _nilErrVEscape, func(actual, expect error) bool { return expect != nil && Contains(actual, expect.Error()) })
	// })
	// t.Run("success,case1,io.ErrUnexpectedEOF,+v,xerrors.Errorf", func(t *testing.T) {
	// 	case1(t, xerrors.Errorf, io.ErrUnexpectedEOF, "+v", _nilErrV, _nilErrVEscape, func(actual, expect error) bool { return expect != nil && Contains(actual, expect.Error()) })
	// })
	// t.Run("success,case1,formatterError,+v,xerrors.Errorf", func(t *testing.T) {
	// 	case1(t, xerrors.Errorf, formatterError, "+v", _nilErrV, _nilErrVEscape, func(actual, expect error) bool { return expect != nil && Contains(actual, expect.Error()) })
	// })
	// t.Run("success,case1,<nil>,+v,xerrors.Errorf", func(t *testing.T) {
	// 	case1(t, xerrors.Errorf, error(nil), "+v", _nilErrV, _nilErrVEscape, func(actual, expect error) bool { return expect != nil && Contains(actual, expect.Error()) })
	// })

	t.Run("failure,%d", func(t *testing.T) {
		actual := Errorf("%d", 123456)
		expect := "123456"
		if actual.Error() != expect {
			t.Errorf("❌: [%s]: [%%d]:\n[EXPECT]:\n%v\n[ACTUAL]:\n%v\n", t.Name(), expect, actual)
		}
		t.Logf("ℹ️: [%s]:\nError():\n%s\n", t.Name(), actual.Error())
	})
}

func Test_wrapError_FormatError(t *testing.T) {
	t.Parallel()
	t.Run("success,case1", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)

		e1 := Errorf("wrap: %w", io.ErrUnexpectedEOF)
		err := Errorf("wrap: %w", e1)
		_, _ = fmt.Fprintf(buf, "%#v", err)
		if expect, actual := regexp.MustCompile(`&errorz.wrapError{msg:"wrap", err:.+, frame:\[3\]uintptr{.+, .+, .+}}`), buf.String(); !expect.MatchString(actual) {
			t.Errorf("❌: expect(%v) != actual(%v)", expect, actual)
		}
	})
}

func Test_wrapError_writeCallers(t *testing.T) {
	t.Parallel()
	t.Run("failure,!ok", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)

		e := &wrapError{msg: "test", err: nil}
		e.writeCallers(buf)
		if expect, actual := "", buf.String(); expect != actual {
			t.Errorf("❌: expect(%v) != actual(%v)", expect, actual)
		}

		runtime.Callers(0, e.frame[:])
		e.frame[0] = 0
		e.writeCallers(buf)
		if expect, actual := "", buf.String(); expect != actual {
			t.Errorf("❌: expect(%v) != actual(%v)", expect, actual)
		}
	})
}
