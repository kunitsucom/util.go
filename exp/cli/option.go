package cliz

import errorz "github.com/kunitsucom/util.go/errors"

type (
	// StringOption is the option for string value.
	StringOption struct {
		// Name is the name of the option.
		Name string
		// Short is the short name of the option.
		Short string
		// Environment is the environment variable name of the option.
		Environment string
		// Description is the description of the option.
		Description string
		// Default is the default value of the option.
		Default *string

		// value is the value of the option.
		value *string
	}

	// BoolOption is the option for bool value.
	BoolOption struct {
		// Name is the name of the option.
		Name string
		// Short is the short name of the option.
		Short string
		// Environment is the environment variable name of the option.
		Environment string
		// Description is the description of the option.
		Description string
		// Default is the default value of the option.
		Default *bool

		// value is the value of the option.
		value *bool
	}

	// IntOption is the option for int value.
	IntOption struct {
		// Name is the name of the option.
		Name string
		// Short is the short name of the option.
		Short string
		// Environment is the environment variable name of the option.
		Environment string
		// Description is the description of the option.
		Description string
		// Default is the default value of the option.
		Default *int

		// value is the value of the option.
		value *int
	}

	// Float64Option is the option for float value.
	Float64Option struct {
		// Name is the name of the option.
		Name string
		// Short is the short name of the option.
		Short string
		// Environment is the environment variable name of the option.
		Environment string
		// Description is the description of the option.
		Description string
		// Default is the default value of the option.
		Default *float64

		// value is the value of the option.
		value *float64
	}
)

func (o *StringOption) GetName() string         { return o.Name }
func (o *StringOption) GetShort() string        { return o.Short }
func (o *StringOption) GetEnvironment() string  { return o.Environment }
func (o *StringOption) HasDefault() bool        { return o.Default != nil }
func (o *StringOption) getDefault() interface{} { return *o.Default }
func (o *StringOption) HasValue() bool          { return o.value != nil }
func (o *StringOption) GetDescription() string {
	if o.Description != "" {
		return o.Description
	}
	return "string value"
}
func (*StringOption) private() {}

func (o *BoolOption) GetName() string         { return o.Name }
func (o *BoolOption) GetShort() string        { return o.Short }
func (o *BoolOption) GetEnvironment() string  { return o.Environment }
func (o *BoolOption) HasDefault() bool        { return o.Default != nil }
func (o *BoolOption) getDefault() interface{} { return *o.Default }
func (o *BoolOption) HasValue() bool          { return o.value != nil }
func (o *BoolOption) GetDescription() string {
	if o.Description != "" {
		return o.Description
	}
	return "bool value"
}
func (*BoolOption) private() {}

func (o *IntOption) GetName() string         { return o.Name }
func (o *IntOption) GetShort() string        { return o.Short }
func (o *IntOption) GetEnvironment() string  { return o.Environment }
func (o *IntOption) HasDefault() bool        { return o.Default != nil }
func (o *IntOption) getDefault() interface{} { return *o.Default }
func (o *IntOption) HasValue() bool          { return o.value != nil }
func (o *IntOption) GetDescription() string {
	if o.Description != "" {
		return o.Description
	}
	return "int value"
}
func (*IntOption) private() {}

func (o *Float64Option) GetName() string         { return o.Name }
func (o *Float64Option) GetShort() string        { return o.Short }
func (o *Float64Option) GetEnvironment() string  { return o.Environment }
func (o *Float64Option) HasDefault() bool        { return o.Default != nil }
func (o *Float64Option) getDefault() interface{} { return *o.Default }
func (o *Float64Option) HasValue() bool          { return o.value != nil }
func (o *Float64Option) GetDescription() string {
	if o.Description != "" {
		return o.Description
	}
	return "float64 value"
}
func (*Float64Option) private() {}

func (cmd *Command) GetOptionString(name string) (string, error) {
	// Search the contents of the subcommand in reverse order and prioritize the options of the descendant commands.
	for i := range cmd.SubCommands {
		subcmd := cmd.SubCommands[len(cmd.SubCommands)-1-i]
		v, err := subcmd.GetOptionString(name)
		if err == nil {
			return v, nil
		}
	}

	v, err := cmd.getOptionString(name)
	if err == nil {
		return v, nil
	}

	return "", errorz.Errorf("%s: %w", name, ErrUnknownOption)
}

//nolint:cyclop
func (cmd *Command) getOptionString(name string) (string, error) {
	if len(cmd.calledCommands) > 0 { //nolint:nestif
		for _, opt := range cmd.Options {
			if o, ok := opt.(*StringOption); ok {
				if (o.Name != "" && o.Name == name) || (o.Short != "" && o.Short == name) || (o.Environment != "" && o.Environment == name) {
					if o.value != nil {
						return *o.value, nil
					}
				}
			}
		}
	}
	return "", errorz.Errorf("%s: %s: %w", cmd.Name, name, ErrUnknownOption)
}

func (cmd *Command) GetOptionBool(name string) (bool, error) {
	// Search the contents of the subcommand in reverse order and prioritize the options of the descendant commands.
	for i := range cmd.SubCommands {
		subcmd := cmd.SubCommands[len(cmd.SubCommands)-1-i]
		v, err := subcmd.GetOptionBool(name)
		if err == nil {
			return v, nil
		}
	}

	v, err := cmd.getOptionBool(name)
	if err == nil {
		return v, nil
	}

	return false, errorz.Errorf("%s: %w", name, ErrUnknownOption)
}

//nolint:cyclop
func (cmd *Command) getOptionBool(name string) (bool, error) {
	if len(cmd.calledCommands) > 0 { //nolint:nestif
		for _, opt := range cmd.Options {
			if o, ok := opt.(*BoolOption); ok {
				TraceLog.Printf("getOptionBool: %s: option: %#v", cmd.Name, o)
				// If Name, Short, or Environment is matched, return the value.
				if (o.Name != "" && o.Name == name) || (o.Short != "" && o.Short == name) || (o.Environment != "" && o.Environment == name) {
					if o.value != nil {
						return *o.value, nil
					}
				}
			}
		}
	}
	return false, errorz.Errorf("%s: %s: %w", cmd.Name, name, ErrUnknownOption)
}

func (cmd *Command) GetOptionInt(name string) (int, error) {
	// Search the contents of the subcommand in reverse order and prioritize the options of the descendant commands.
	for i := range cmd.SubCommands {
		subcmd := cmd.SubCommands[len(cmd.SubCommands)-1-i]
		v, err := subcmd.GetOptionInt(name)
		if err == nil {
			return v, nil
		}
	}

	v, err := cmd.getOptionInt(name)
	if err == nil {
		return v, nil
	}

	return 0, errorz.Errorf("%s: %w", name, ErrUnknownOption)
}

//nolint:cyclop
func (cmd *Command) getOptionInt(name string) (int, error) {
	if len(cmd.calledCommands) > 0 { //nolint:nestif
		for _, opt := range cmd.Options {
			if o, ok := opt.(*IntOption); ok {
				if (o.Name != "" && o.Name == name) || (o.Short != "" && o.Short == name) || (o.Environment != "" && o.Environment == name) {
					if o.value != nil {
						return *o.value, nil
					}
				}
			}
		}
	}
	return 0, errorz.Errorf("%s: %s: %w", cmd.Name, name, ErrUnknownOption)
}

func (cmd *Command) GetOptionFloat64(name string) (float64, error) {
	// Search the contents of the subcommand in reverse order and prioritize the options of the descendant commands.
	for i := range cmd.SubCommands {
		subcmd := cmd.SubCommands[len(cmd.SubCommands)-1-i]
		v, err := subcmd.GetOptionFloat64(name)
		if err == nil {
			return v, nil
		}
	}

	v, err := cmd.getOptionFloat64(name)
	if err == nil {
		return v, nil
	}

	return 0, errorz.Errorf("%s: %w", name, ErrUnknownOption)
}

//nolint:cyclop
func (cmd *Command) getOptionFloat64(name string) (float64, error) {
	if len(cmd.calledCommands) > 0 { //nolint:nestif
		for _, opt := range cmd.Options {
			if o, ok := opt.(*Float64Option); ok {
				if (o.Name != "" && o.Name == name) || (o.Short != "" && o.Short == name) || (o.Environment != "" && o.Environment == name) {
					if o.value != nil {
						return *o.value, nil
					}
				}
			}
		}
	}
	return 0, errorz.Errorf("%s: %s: %w", cmd.Name, name, ErrUnknownOption)
}
