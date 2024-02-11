package errorz

import "errors"

type retryableError struct {
	error
	retryable bool
}

var (
	_ interface{ Error() string } = (*retryableError)(nil)
	_ interface{ Unwrap() error } = (*retryableError)(nil)
)

func (e *retryableError) Unwrap() error {
	return e.error
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
