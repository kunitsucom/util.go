package cliz

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
