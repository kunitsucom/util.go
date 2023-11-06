package cliz

import (
	errorz "github.com/kunitsucom/util.go/errors"
)

// Default is the helper function to create a default value.
func Default[T interface{}](v T) *T { return ptr[T](v) }

func ptr[T interface{}](v T) *T { return &v }

func (cmd *Command) GetName() string {
	if cmd == nil {
		return ""
	}
	return cmd.Name
}

//nolint:cyclop
func (cmd *Command) loadDefaults() error {
	for _, opt := range cmd.Options {
		if !opt.HasDefault() {
			// If default value is not set, o.value remains.
			continue
		}

		switch o := opt.(type) {
		case *StringOption:
			DebugLog.Printf("%s: %s=%s", cmd.Name, o.Name, *o.Default)
			o.value = o.Default
		case *BoolOption:
			DebugLog.Printf("%s: %s=%t", cmd.Name, o.Name, *o.Default)
			o.value = o.Default
		case *IntOption:
			DebugLog.Printf("%s: %s=%d", cmd.Name, o.Name, *o.Default)
			o.value = o.Default
		case *Float64Option:
			DebugLog.Printf("%s: %s=%f", cmd.Name, o.Name, *o.Default)
			o.value = o.Default
		default:
			return errorz.Errorf("%s: %w", o.GetName(), ErrInvalidOptionType)
		}
	}

	for _, subcmd := range cmd.SubCommands {
		if err := subcmd.loadDefaults(); err != nil {
			return err
		}
	}
	return nil
}
