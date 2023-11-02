package cliz

import (
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
	breakArg          = "--"
	longOptionPrefix  = "--"
	shortOptionPrefix = "-"
)

type (
	// Command is a structure for building command lines. Please fill in each field for the structure you are facing.
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
		//
		// If you want to use the default usage, remain empty.
		// Otherwise, set the custom usage.
		Usage string
		// UsageFunc is custom usage function.
		//
		// If you want to use the default usage function, remain nil.
		// Otherwise, set the custom usage function.
		UsageFunc func(c *Command)

		called    []string
		remaining []string
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
		// HasValue returns whether the option has a value.
		HasValue() bool

		// private is the private method for internal interface.
		private()
	}
)

func (cmd *Command) getDescription() string {
	if cmd.Description != "" {
		return cmd.Description
	}
	return fmt.Sprintf("command %q description", strings.Join(cmd.called, " "))
}

// Default is the helper function to create a default value.
func Default[T interface{}](v T) *T { return ptr[T](v) }

func ptr[T interface{}](v T) *T { return &v }

// getSubcommand returns the subcommand if cmd contains the subcommand.
func (cmd *Command) getSubcommand(arg string) (subcmd *Command) {
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
func (cmd *Command) parseArgs(args []string) (calledCommands []string, remainingArgs []string, err error) {
	cmd.called = append(cmd.called, cmd.Name)
	cmd.remaining = make([]string, 0)

	i := 0
argsLoop:
	for ; i < len(args); i++ {
		arg := args[i]

		switch {
		case arg == breakArg:
			cmd.remaining = append(cmd.remaining, args[i+1:]...)
			break argsLoop
		case strings.HasPrefix(arg, shortOptionPrefix):
			for _, opt := range cmd.Options {
				switch o := opt.(type) {
				case *StringOption:
					switch {
					case optionArg(o, arg):
						DebugLog.Printf("%s: option: %s: %s", cmd.Name, o.Name, arg)
						if hasOptionValue(args, i) {
							return nil, nil, errorz.Errorf("%s: %w", arg, ErrMissingOptionValue)
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
						optVal, err := strconv.ParseBool(optionEqualArgExtractValue(arg))
						if err != nil {
							return nil, nil, errorz.Errorf("%s: %w", arg, err)
						}
						o.value = &optVal
						TraceLog.Printf("%s: parsed option: %s: %v", cmd.Name, o.Name, *o.value)
						continue argsLoop
					}
				case *IntOption:
					switch {
					case optionArg(o, arg):
						DebugLog.Printf("%s: option: %s: %s", cmd.Name, o.Name, arg)
						if hasOptionValue(args, i) {
							return nil, nil, errorz.Errorf("%s: %w", arg, ErrMissingOptionValue)
						}
						optVal, err := strconv.Atoi(args[i+1])
						if err != nil {
							return nil, nil, errorz.Errorf("%s: %w", arg, err)
						}
						o.value = &optVal
						i++
						TraceLog.Printf("%s: parsed option: %s: %v", cmd.Name, o.Name, *o.value)
						continue argsLoop
					case optionEqualArg(o, arg):
						DebugLog.Printf("%s: option: %s: %s", cmd.Name, o.Name, arg)
						optVal, err := strconv.Atoi(optionEqualArgExtractValue(arg))
						if err != nil {
							return nil, nil, errorz.Errorf("%s: %w", arg, err)
						}
						o.value = &optVal
						TraceLog.Printf("%s: parsed option: %s: %v", cmd.Name, o.Name, *o.value)
						continue argsLoop
					}
				case *Float64Option:
					switch {
					case optionArg(o, arg):
						DebugLog.Printf("%s: option: %s: %s", cmd.Name, o.Name, arg)
						if hasOptionValue(args, i) {
							return nil, nil, errorz.Errorf("%s: %w", arg, ErrMissingOptionValue)
						}
						optVal, err := strconv.ParseFloat(args[i+1], 64)
						if err != nil {
							return nil, nil, errorz.Errorf("%s: %w", arg, err)
						}
						o.value = &optVal
						i++
						TraceLog.Printf("%s: parsed option: %s: %v", cmd.Name, o.Name, *o.value)
						continue argsLoop
					case optionEqualArg(o, arg):
						DebugLog.Printf("%s: option: %s: %s", cmd.Name, o.Name, arg)
						optVal, err := strconv.ParseFloat(optionEqualArgExtractValue(arg), 64)
						if err != nil {
							return nil, nil, errorz.Errorf("%s: %w", arg, err)
						}
						o.value = &optVal
						TraceLog.Printf("%s: parsed option: %s: %v", cmd.Name, o.Name, *o.value)
						continue argsLoop
					}
				default:
					return nil, nil, errorz.Errorf("%s: %w", arg, ErrInvalidOptionType)
				}
			}
			return nil, nil, errorz.Errorf("%s: %w", arg, ErrUnknownOption)
		default:
			subcmd := cmd.getSubcommand(arg)
			// If subcmd is nil, it is not a subcommand.
			if subcmd == nil {
				cmd.remaining = append(cmd.remaining, arg)
				continue argsLoop
			}

			TraceLog.Printf("parse: subcommand: %s", arg)
			subcmd.called = append(subcmd.called, cmd.called...)
			called, remaining, err := subcmd.parseArgs(args[i+1:])
			if err != nil {
				return nil, nil, errorz.Errorf("%s: %w", arg, err)
			}
			return called, remaining, nil
		}
	}

	return cmd.called, cmd.remaining, nil
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

// Parse parses the arguments as commands and sub commands and options.
//
// If the "--help" option is specified, it will be displayed and ErrHelp will be returned.
//
// If the option is not specified, the default value will be used.
//
// If the environment variable is specified, it will be used as the value of the option.
func (cmd *Command) Parse(args []string) (calledCommands []string, remainingArgs []string, err error) {
	appendHelpOption(cmd)

	if err := cmd.preCheckSubCommands(); err != nil {
		return nil, nil, errorz.Errorf("failed to pre-check commands: %w", err)
	}

	if err := cmd.preCheckOptions(); err != nil {
		return nil, nil, errorz.Errorf("failed to pre-check options: %w", err)
	}

	if err := cmd.loadDefaults(); err != nil {
		return nil, nil, errorz.Errorf("failed to load default: %w", err)
	}

	if err := cmd.loadEnvironments(); err != nil {
		return nil, nil, errorz.Errorf("failed to load environment: %w", err)
	}

	called, remaining, err := cmd.parseArgs(args)
	if err != nil {
		return nil, nil, errorz.Errorf("failed to parse commands and options: %w", err)
	}

	if err := cmd.checkHelp(); err != nil {
		return nil, nil, err //nolint:wrapcheck
	}

	if err := cmd.postCheckOptions(); err != nil {
		return nil, nil, errorz.Errorf("failed to post-check options: %w", err)
	}

	return called, remaining, nil
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
	if len(cmd.called) > 0 { //nolint:nestif
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
	if len(cmd.called) > 0 { //nolint:nestif
		for _, opt := range cmd.Options {
			if o, ok := opt.(*BoolOption); ok {
				TraceLog.Printf("getBoolOption: %s: option: %#v", cmd.Name, o)
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
	if len(cmd.called) > 0 { //nolint:nestif
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

func (cmd *Command) GetFloat64Option(name string) (float64, error) {
	// Search the contents of the subcommand in reverse order and prioritize the options of the descendant commands.
	for i := range cmd.SubCommands {
		subcmd := cmd.SubCommands[len(cmd.SubCommands)-1-i]
		v, err := subcmd.GetFloat64Option(name)
		if err == nil {
			return v, nil
		}
	}

	v, err := cmd.getFloat64Option(name)
	if err == nil {
		return v, nil
	}

	return 0, errorz.Errorf("%s: %w", name, ErrUnknownOption)
}

//nolint:cyclop
func (cmd *Command) getFloat64Option(name string) (float64, error) {
	if len(cmd.called) > 0 { //nolint:nestif
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
