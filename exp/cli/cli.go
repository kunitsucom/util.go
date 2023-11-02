package cliz

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	errorz "github.com/kunitsucom/util.go/errors"
)

//nolint:gochecknoglobals
var (
	// Stdout is the writer to be used for standard output.
	Stdout io.Writer = os.Stdout
	// Stderr is the writer to be used for standard error.
	Stderr io.Writer = os.Stderr
)

//nolint:revive,stylecheck
const (
	// UTIL_GO_CLI_TRACE is the environment variable name for trace log.
	UTIL_GO_CLI_TRACE = "UTIL_GO_CLI_TRACE"
	// UTIL_GO_CLI_DEBUG is the environment variable name for debug log.
	UTIL_GO_CLI_DEBUG = "UTIL_GO_CLI_DEBUG"
)

//nolint:gochecknoglobals
var (
	// TraceLog is the logger to be used for trace log.
	TraceLog = newLogger(Stderr, UTIL_GO_CLI_TRACE, "TRACE: ")
	// DebugLog is the logger to be used for debug log.
	DebugLog = newLogger(Stderr, UTIL_GO_CLI_DEBUG, "DEBUG: ")
)

func newLogger(w io.Writer, environ string, prefix string) *log.Logger {
	if v := os.Getenv(environ); v == "true" {
		return log.New(w, prefix, log.LstdFlags|log.Lshortfile)
	}

	return log.New(io.Discard, prefix, log.LstdFlags)
}

const (
	// HelpOptionName is the option name for help.
	HelpOptionName    = "help"
	breakArg          = "--"
	longOptionPrefix  = "--"
	shortOptionPrefix = "-"
)

type (
	Command struct {
		// Name is the name of the command.
		Name string
		// Description is the description of the command.
		Description string
		// SubCommands is the subcommands of the command.
		SubCommands []*Command
		// Options is the options of the command.
		Options []Option
		// Usage is the usage of the command.
		Usage func(c *Command)

		cmdStack []string
	}

	// Option is the interface for the option.
	Option interface {
		// GetName returns the name of the option.
		GetName() string
		// GetShort returns the short name of the option.
		GetShort() string
		// GetEnvironment returns the environment variable name of the option.
		GetEnvironment() string
		// GetDescription returns the description of the option.
		GetDescription() string
		// HasDefault returns whether the option has a default value.
		HasDefault() bool
		// getDefault returns the default value of the option.
		getDefault() interface{}

		// private is the private method for internal interface.
		private()
	}

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
)

func (cmd *Command) getDescription() string {
	if cmd.Description != "" {
		return cmd.Description
	}
	return fmt.Sprintf("command %q description", strings.Join(cmd.cmdStack, " "))
}

// Default is the helper function to create a default value.
func Default[T interface{}](v T) *T {
	return ptr[T](v)
}

func ptr[T interface{}](v T) *T {
	return &v
}

func (o *StringOption) GetName() string         { return o.Name }
func (o *StringOption) GetShort() string        { return o.Short }
func (o *StringOption) GetEnvironment() string  { return o.Environment }
func (o *StringOption) HasDefault() bool        { return o.Default != nil }
func (o *StringOption) getDefault() interface{} { return *o.Default }
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
func (o *IntOption) GetDescription() string {
	if o.Description != "" {
		return o.Description
	}
	return "int value"
}
func (*IntOption) private() {}

func (cmd *Command) GetSubcommand(arg string) (subcmd *Command) {
	for _, subcmd := range cmd.SubCommands {
		if subcmd.Name == arg {
			return subcmd
		}
	}
	return nil
}

func optionArg(o Option, arg string) bool {
	return longOptionPrefix+o.GetName() == arg || shortOptionPrefix+o.GetShort() == arg
}

func optionEqualArg(o Option, arg string) bool {
	return strings.HasPrefix(arg, longOptionPrefix+o.GetName()+"=") || strings.HasPrefix(arg, shortOptionPrefix+o.GetShort()+"=")
}

func optionEqualArgExtractValue(arg string) string {
	return strings.Join(strings.Split(arg, "=")[1:], "=")
}

func hasOptionValue(args []string, i int) bool {
	lastIndex := len(args) - 1
	return i+1 > lastIndex
}

//nolint:funlen,gocognit,cyclop
func (cmd *Command) parse(args []string) (remaining []string, err error) {
	cmd.cmdStack = append(cmd.cmdStack, cmd.Name)
	remaining = make([]string, 0)

	i := 0
argsLoop:
	for ; i < len(args); i++ {
		arg := args[i]

		switch {
		case arg == breakArg:
			remaining = append(remaining, args[i+1:]...)
			break argsLoop
		case strings.HasPrefix(arg, shortOptionPrefix):
			for _, opt := range cmd.Options {
				switch o := opt.(type) {
				case *StringOption:
					switch {
					case optionArg(o, arg):
						DebugLog.Printf("%s: option: %s: %s", cmd.Name, o.Name, arg)
						if hasOptionValue(args, i) {
							return nil, errorz.Errorf("%s: %w", arg, ErrMissingOptionValue)
						}
						o.value = ptr(args[i+1])
						i++
						TraceLog.Printf("%s: parsed option: %s: %v", cmd.Name, o.Name, *o.value)
						continue argsLoop
					case optionEqualArg(o, arg):
						DebugLog.Printf("%s: option: %s: %s", cmd.Name, o.Name, arg)
						o.value = ptr(optionEqualArgExtractValue(arg))
						TraceLog.Printf("%s: parsed option: %s: %v", cmd.Name, o.Name, *o.value)
						continue argsLoop
					}
				case *BoolOption:
					switch {
					case optionArg(o, arg):
						DebugLog.Printf("%s: option: %s: %s", cmd.Name, o.Name, arg)
						o.value = ptr(true)
						TraceLog.Printf("%s: parsed option: %s: %v", cmd.Name, o.Name, *o.value)
						continue argsLoop
					case optionEqualArg(o, arg):
						DebugLog.Printf("%s: option: %s: %s", cmd.Name, o.Name, arg)
						optVal, _ := strconv.ParseBool(optionEqualArgExtractValue(arg))
						o.value = &optVal
						TraceLog.Printf("%s: parsed option: %s: %v", cmd.Name, o.Name, *o.value)
						continue argsLoop
					}
				case *IntOption:
					switch {
					case optionArg(o, arg):
						DebugLog.Printf("%s: option: %s: %s", cmd.Name, o.Name, arg)
						if hasOptionValue(args, i) {
							return nil, errorz.Errorf("%s: %w", arg, ErrMissingOptionValue)
						}
						optVal, _ := strconv.Atoi(args[i+1])
						o.value = &optVal
						i++
						TraceLog.Printf("%s: parsed option: %s: %v", cmd.Name, o.Name, *o.value)
						continue argsLoop
					case optionEqualArg(o, arg):
						DebugLog.Printf("%s: option: %s: %s", cmd.Name, o.Name, arg)
						optVal, _ := strconv.Atoi(optionEqualArgExtractValue(arg))
						o.value = &optVal
						TraceLog.Printf("%s: parsed option: %s: %v", cmd.Name, o.Name, *o.value)
						continue argsLoop
					}
				default:
					return nil, errorz.Errorf("%s: %w", arg, ErrInvalidOptionType)
				}
			}
			return nil, errorz.Errorf("%s: %w", arg, ErrUnknownOption)
		default:
			subcmd := cmd.GetSubcommand(arg)
			// If subcmd is nil, it is not a subcommand.
			if subcmd == nil {
				remaining = append(remaining, arg)
				continue argsLoop
			}

			TraceLog.Printf("parse: subcommand: %s", arg)
			subcmd.cmdStack = append(subcmd.cmdStack, cmd.cmdStack...)
			remaining, err := subcmd.parse(args[i+1:])
			if err != nil {
				return nil, errorz.Errorf("%s: %w", arg, err)
			}
			return remaining, nil
		}
	}

	return remaining, nil
}

func (cmd *Command) checkHelp() error {
	TraceLog.Printf("checkHelp: %s", cmd.Name)
	v, err := cmd.getBoolOption(HelpOptionName)
	if err == nil && v {
		cmd.usage()
		return ErrHelp
	}
	for _, subcmd := range cmd.SubCommands {
		if err := subcmd.checkHelp(); err != nil {
			return err
		}
	}
	return nil
}

// Parse parses the commands and options.
//
// If the help option is specified, it will be displayed and ErrHelp will be returned.
//
// If the option is not specified, the default value will be used.
func (cmd *Command) Parse(args []string) (remaining []string, err error) {
	appendHelpOption(cmd)

	if err := cmd.preCheckSubCommands(); err != nil {
		return nil, errorz.Errorf("failed to check commands: %w", err)
	}

	if err := cmd.preCheckOptions(); err != nil {
		return nil, errorz.Errorf("failed to check options: %w", err)
	}

	if err := cmd.loadEnvironments(); err != nil {
		return nil, errorz.Errorf("failed to load environment: %w", err)
	}

	r, err := cmd.parse(args)
	if err != nil {
		return nil, errorz.Errorf("failed to parse commands and options: %w", err)
	}

	if err := cmd.checkHelp(); err != nil {
		return nil, err //nolint:wrapcheck
	}

	return r, nil
}

// IsHelp returns whether the error is ErrHelp.
func IsHelp(err error) bool {
	return errors.Is(err, ErrHelp)
}

func (cmd *Command) GetStringOption(name string) (string, error) {
	// Search the contents of the subcommand in reverse order and prioritize the options of the descendant commands.
	for i := range cmd.SubCommands {
		subcmd := cmd.SubCommands[len(cmd.SubCommands)-1-i]
		v, err := subcmd.GetStringOption(name)
		if err == nil {
			return v, nil
		}
	}

	v, err := cmd.getStringOption(name)
	if err == nil {
		return v, nil
	}

	return "", errorz.Errorf("%s: %w", name, ErrUnknownOption)
}

//nolint:cyclop
func (cmd *Command) getStringOption(name string) (string, error) {
	if len(cmd.cmdStack) > 0 { //nolint:nestif
		for _, opt := range cmd.Options {
			if o, ok := opt.(*StringOption); ok {
				if (o.Name != "" && o.Name == name) || (o.Short != "" && o.Short == name) || (o.Environment != "" && o.Environment == name) {
					if o.value != nil {
						return *o.value, nil
					}
					if o.Default != nil {
						return *o.Default, nil
					}
				}
			}
		}
	}
	return "", errorz.Errorf("%s: %s: %w", cmd.Name, name, ErrUnknownOption)
}

func (cmd *Command) GetBoolOption(name string) (bool, error) {
	// Search the contents of the subcommand in reverse order and prioritize the options of the descendant commands.
	for i := range cmd.SubCommands {
		subcmd := cmd.SubCommands[len(cmd.SubCommands)-1-i]
		v, err := subcmd.GetBoolOption(name)
		if err == nil {
			return v, nil
		}
	}

	v, err := cmd.getBoolOption(name)
	if err == nil {
		return v, nil
	}

	return false, errorz.Errorf("%s: %w", name, ErrUnknownOption)
}

//nolint:cyclop
func (cmd *Command) getBoolOption(name string) (bool, error) {
	if len(cmd.cmdStack) > 0 { //nolint:nestif
		for _, opt := range cmd.Options {
			if o, ok := opt.(*BoolOption); ok {
				TraceLog.Printf("getBoolOption: %s: option: %#v", cmd.Name, o)
				// If Name, Short, or Environment is matched, return the value.
				if (o.Name != "" && o.Name == name) || (o.Short != "" && o.Short == name) || (o.Environment != "" && o.Environment == name) {
					if o.value != nil {
						return *o.value, nil
					}
					// If o.value is nil, use o.Default.
					if o.Default != nil {
						return *o.Default, nil
					}
				}
			}
		}
	}
	return false, errorz.Errorf("%s: %s: %w", cmd.Name, name, ErrUnknownOption)
}

func (cmd *Command) GetIntOption(name string) (int, error) {
	// Search the contents of the subcommand in reverse order and prioritize the options of the descendant commands.
	for i := range cmd.SubCommands {
		subcmd := cmd.SubCommands[len(cmd.SubCommands)-1-i]
		v, err := subcmd.GetIntOption(name)
		if err == nil {
			return v, nil
		}
	}

	v, err := cmd.getIntOption(name)
	if err == nil {
		return v, nil
	}

	return 0, errorz.Errorf("%s: %w", name, ErrUnknownOption)
}

//nolint:cyclop
func (cmd *Command) getIntOption(name string) (int, error) {
	if len(cmd.cmdStack) > 0 { //nolint:nestif
		for _, opt := range cmd.Options {
			if o, ok := opt.(*IntOption); ok {
				if (o.Name != "" && o.Name == name) || (o.Short != "" && o.Short == name) || (o.Environment != "" && o.Environment == name) {
					if o.value != nil {
						return *o.value, nil
					}
					if o.Default != nil {
						return *o.Default, nil
					}
				}
			}
		}
	}
	return 0, errorz.Errorf("%s: %s: %w", cmd.Name, name, ErrUnknownOption)
}
