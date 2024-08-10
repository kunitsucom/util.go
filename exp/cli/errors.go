package cliz

import "errors"

var (
	ErrHelp                   = errors.New("help requested")
	ErrNilContext             = errors.New("nil context")
	ErrCommandNotSetInContext = errors.New("command not set in context")
	ErrCommandFuncNotSet      = errors.New("command func not set")
	ErrMissingOptionValue     = errors.New("missing option value")
	ErrOptionRequired         = errors.New("option required")
	ErrUnknownOption          = errors.New("unknown option")
	ErrInvalidOptionType      = errors.New("invalid option type")
	ErrDuplicateOptionName    = errors.New("duplicate option name")
	ErrDuplicateSubCommand    = errors.New("duplicate sub command")
)
