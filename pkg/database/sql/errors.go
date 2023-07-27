package sqlz

import "errors"

var (
	ErrMustBePointer        = errors.New("sql: must be pointer")
	ErrMustNotNil           = errors.New("sql: must not nil")
	ErrDataTypeNotSupported = errors.New("sql: data type not supported")
)
