package errorz

import (
	"errors"
	"strings"
)

func HasPrefix(err error, prefix string) bool {
	return err != nil && strings.HasPrefix(err.Error(), prefix)
}

func HasSuffix(err error, suffix string) bool {
	return err != nil && strings.HasSuffix(err.Error(), suffix)
}

func Contains(err error, substr string) bool {
	return err != nil && strings.Contains(err.Error(), substr)
}

type retryableError struct {
	error
	retryable bool
}

func (e *retryableError) IsRetryable() bool {
	return e.retryable
}

func WithRetryable(err error, retryable bool) error {
	return &retryableError{
		error:     err,
		retryable: retryable,
	}
}

func IsRetryable(err error) bool {
	var target interface{ IsRetryable() bool }
	if errors.As(err, &target) {
		return target.IsRetryable()
	}

	return false
}

func PanicOrIgnore(err error, ignores ...error) {
	if err == nil {
		return
	}

	for _, ignore := range ignores {
		if errors.Is(err, ignore) {
			return
		}
	}

	panic(err)
}
