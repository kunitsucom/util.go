package cliz

import "errors"

var (
	ErrHelp                        = errors.New("help requested")
	ErrMissingOptionValue          = errors.New("missing option value")
	ErrOptionRequired              = errors.New("option required")
	ErrNoOption                    = errors.New("no option")
	ErrUnknownOption               = errors.New("unknown option")
	ErrInvalidOptionType           = errors.New("invalid option type")
	ErrUnexpectedError             = errors.New("unexpected error")
	ErrDuplicateOptionName         = errors.New("duplicate option name")
	ErrDuplicateSubCommand         = errors.New("duplicate sub command")
	ErrMultipleOptionsDefaultValue = errors.New("multiple options with the same name, only one option can have a default value")
)
