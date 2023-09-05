package testingz

import (
	"errors"
	"fmt"
)

var ErrTestError = errors.New("testingz: test error")

type FormatterError struct {
	ErrorFunc  func() string
	FormatFunc func(s fmt.State, verb rune)
}

var (
	_ error         = (*FormatterError)(nil)
	_ fmt.Formatter = (*FormatterError)(nil)
)

func (f *FormatterError) Error() string {
	return f.ErrorFunc()
}

func (f *FormatterError) Format(s fmt.State, verb rune) {
	f.FormatFunc(s, verb)
}
