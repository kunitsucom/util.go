package cliz

import (
	"os"
	"strconv"

	errorz "github.com/kunitsucom/util.go/errors"
)

//nolint:cyclop
func (cmd *Command) loadEnvironments() error {
	for _, opt := range cmd.Options {
		if opt.GetEnvironment() == "" {
			// If v is an empty string, o.value remains.
			continue
		}

		switch o := opt.(type) {
		case *StringOption:
			if s := os.Getenv(o.Environment); s != "" {
				o.value = &s
			}
		case *BoolOption:
			if s := os.Getenv(o.Environment); s != "" {
				v, err := strconv.ParseBool(s)
				if err != nil {
					return errorz.Errorf("%s: %w", o.Environment, err)
				}
				o.value = &v
			}
		case *IntOption:
			if s := os.Getenv(o.Environment); s != "" {
				v, err := strconv.Atoi(s)
				if err != nil {
					return errorz.Errorf("%s: %w", o.Environment, err)
				}
				o.value = &v
			}
		default:
			return errorz.Errorf("%s: %w", o.GetName(), ErrInvalidOptionType)
		}
	}

	for _, subcmd := range cmd.SubCommands {
		if err := subcmd.loadEnvironments(); err != nil {
			return err
		}
	}
	return nil
}
