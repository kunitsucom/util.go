package cliz

import (
	"context"
	"fmt"
	"io"
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
		// Description is the description of the command.
		Description string
		// Options is the options of the command.
		Options []Option
		// Func is the function to be executed when the command is executed.
		Func func(ctx context.Context, remainingArgs []string) error
		// SubCommands is the subcommands of the command.
		SubCommands []*Command

		calledCommands []string
		remainingArgs  []string
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
	return fmt.Sprintf("command %q description", strings.Join(cmd.calledCommands, " "))
}

func (cmd *Command) Next() *Command {
	if cmd == nil {
		return nil
	}
	if len(cmd.calledCommands) == 0 {
		return nil
	}
	for _, subcmd := range cmd.SubCommands {
		if len(subcmd.calledCommands) > 0 {
			return subcmd
		}
	}
	return nil
}

func (cmd *Command) GetCalledCommands() []string {
	if cmd == nil {
		return nil
	}

	for _, subcmd := range cmd.SubCommands {
		if len(subcmd.calledCommands) > 0 {
			return subcmd.GetCalledCommands()
		}
	}

	return cmd.calledCommands
}

// getSubcommand returns the subcommand if cmd contains the subcommand.
func (cmd *Command) getSubcommand(arg string) (subcmd *Command) {
	if cmd == nil {
		return nil
	}

	for _, subcmd := range cmd.SubCommands {
		if subcmd.Name == arg {
			return subcmd
		}
	}
	return nil
}

func equalOptionArg(o Option, arg string) bool {
	return longOptionPrefix+o.GetName() == arg || shortOptionPrefix+o.GetShort() == arg
}

func hasPrefixOptionEqualArg(o Option, arg string) bool {
	return strings.HasPrefix(arg, longOptionPrefix+o.GetName()+"=") || strings.HasPrefix(arg, shortOptionPrefix+o.GetShort()+"=")
}

func extractValueOptionEqualArg(arg string) string {
	return strings.Join(strings.Split(arg, "=")[1:], "=")
}

func hasOptionValue(args []string, i int) bool {
	lastIndex := len(args) - 1
	return i+1 > lastIndex
}

//nolint:funlen,gocognit,cyclop
func (cmd *Command) parseArgs(args []string) (remaining []string, err error) {
	cmd.calledCommands = append(cmd.calledCommands, cmd.Name)
	cmd.remainingArgs = make([]string, 0)

	i := 0
argsLoop:
	for ; i < len(args); i++ {
		arg := args[i]

		switch {
		case arg == breakArg:
			cmd.remainingArgs = append(cmd.remainingArgs, args[i+1:]...)
			break argsLoop
		case strings.HasPrefix(arg, shortOptionPrefix):
			for _, opt := range cmd.Options {
				switch o := opt.(type) {
				case *StringOption:
					switch {
					case equalOptionArg(o, arg):
						DebugLog.Printf("%s: option: %s: %s", cmd.Name, o.Name, arg)
						if hasOptionValue(args, i) {
							return nil, errorz.Errorf("%s: %w", arg, ErrMissingOptionValue)
						}
						o.value = ptr(args[i+1])
						i++
						TraceLog.Printf("%s: parsed option: %s: %v", cmd.Name, o.Name, *o.value)
						continue argsLoop
					case hasPrefixOptionEqualArg(o, arg):
						DebugLog.Printf("%s: option: %s: %s", cmd.Name, o.Name, arg)
						o.value = ptr(extractValueOptionEqualArg(arg))
						TraceLog.Printf("%s: parsed option: %s: %v", cmd.Name, o.Name, *o.value)
						continue argsLoop
					}
				case *BoolOption:
					switch {
					case equalOptionArg(o, arg):
						DebugLog.Printf("%s: option: %s: %s", cmd.Name, o.Name, arg)
						o.value = ptr(true)
						TraceLog.Printf("%s: parsed option: %s: %v", cmd.Name, o.Name, *o.value)
						continue argsLoop
					case hasPrefixOptionEqualArg(o, arg):
						DebugLog.Printf("%s: option: %s: %s", cmd.Name, o.Name, arg)
						optVal, err := strconv.ParseBool(extractValueOptionEqualArg(arg))
						if err != nil {
							return nil, errorz.Errorf("%s: %w", arg, err)
						}
						o.value = &optVal
						TraceLog.Printf("%s: parsed option: %s: %v", cmd.Name, o.Name, *o.value)
						continue argsLoop
					}
				case *IntOption:
					switch {
					case equalOptionArg(o, arg):
						DebugLog.Printf("%s: option: %s: %s", cmd.Name, o.Name, arg)
						if hasOptionValue(args, i) {
							return nil, errorz.Errorf("%s: %w", arg, ErrMissingOptionValue)
						}
						optVal, err := strconv.Atoi(args[i+1])
						if err != nil {
							return nil, errorz.Errorf("%s: %w", arg, err)
						}
						o.value = &optVal
						i++
						TraceLog.Printf("%s: parsed option: %s: %v", cmd.Name, o.Name, *o.value)
						continue argsLoop
					case hasPrefixOptionEqualArg(o, arg):
						DebugLog.Printf("%s: option: %s: %s", cmd.Name, o.Name, arg)
						optVal, err := strconv.Atoi(extractValueOptionEqualArg(arg))
						if err != nil {
							return nil, errorz.Errorf("%s: %w", arg, err)
						}
						o.value = &optVal
						TraceLog.Printf("%s: parsed option: %s: %v", cmd.Name, o.Name, *o.value)
						continue argsLoop
					}
				case *Float64Option:
					switch {
					case equalOptionArg(o, arg):
						DebugLog.Printf("%s: option: %s: %s", cmd.Name, o.Name, arg)
						if hasOptionValue(args, i) {
							return nil, errorz.Errorf("%s: %w", arg, ErrMissingOptionValue)
						}
						optVal, err := strconv.ParseFloat(args[i+1], 64)
						if err != nil {
							return nil, errorz.Errorf("%s: %w", arg, err)
						}
						o.value = &optVal
						i++
						TraceLog.Printf("%s: parsed option: %s: %v", cmd.Name, o.Name, *o.value)
						continue argsLoop
					case hasPrefixOptionEqualArg(o, arg):
						DebugLog.Printf("%s: option: %s: %s", cmd.Name, o.Name, arg)
						optVal, err := strconv.ParseFloat(extractValueOptionEqualArg(arg), 64)
						if err != nil {
							return nil, errorz.Errorf("%s: %w", arg, err)
						}
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
			if subcmd := cmd.getSubcommand(arg); subcmd != nil {
				TraceLog.Printf("parse: subcommand: %s", subcmd.Name)
				subcmd.calledCommands = append(subcmd.calledCommands, cmd.calledCommands...)
				cmd.remainingArgs, err = subcmd.parseArgs(args[i+1:])
				if err != nil {
					return nil, errorz.Errorf("%s: %w", arg, err)
				}
				return cmd.remainingArgs, nil
			}

			// If subcmd is nil, it is not a subcommand.
			cmd.remainingArgs = append(cmd.remainingArgs, arg)
			continue argsLoop
		}
	}

	return cmd.remainingArgs, nil
}

func (cmd *Command) initCommand() {
	cmd.calledCommands = make([]string, 0)
	cmd.remainingArgs = make([]string, 0)

	for _, subcmd := range cmd.SubCommands {
		subcmd.initCommand()
	}
}

// Parse parses the arguments as commands and sub commands and options.
//
// If the "--help" option is specified, it will be displayed and ErrHelp will be returned.
//
// If the option is not specified, the default value will be used.
//
// If the environment variable is specified, it will be used as the value of the option.
//
//nolint:cyclop
func (cmd *Command) Parse(args []string) (remainingArgs []string, err error) {
	if len(args) > 0 && (args[0] == os.Args[0] || args[0] == cmd.Name) {
		args = args[1:]
	}

	cmd.initCommand()
	cmd.initAppendHelpOption()

	if err := cmd.preCheckSubCommands(); err != nil {
		return nil, errorz.Errorf("failed to pre-check commands: %w", err)
	}

	if err := cmd.preCheckOptions(); err != nil {
		return nil, errorz.Errorf("failed to pre-check options: %w", err)
	}

	if err := cmd.loadDefaults(); err != nil {
		return nil, errorz.Errorf("failed to load default: %w", err)
	}

	if err := cmd.loadEnvironments(); err != nil {
		return nil, errorz.Errorf("failed to load environment: %w", err)
	}

	remaining, err := cmd.parseArgs(args)
	if err != nil {
		return nil, errorz.Errorf("failed to parse commands and options: %w", err)
	}

	// NOTE: help
	if err := cmd.checkHelp(); err != nil {
		return nil, err //nolint:wrapcheck
	}

	if err := cmd.postCheckOptions(); err != nil {
		return nil, errorz.Errorf("failed to post-check options: %w", err)
	}

	return remaining, nil
}

func (cmd *Command) Run(ctx context.Context, args []string) error {
	remainingArgs, err := cmd.Parse(args)
	if err != nil {
		return errorz.Errorf("%s: %w", cmd.Name, err)
	}

	execCmd := cmd
	for len(execCmd.Next().GetCalledCommands()) > 0 {
		execCmd = execCmd.Next()
	}

	if execCmd.Func == nil {
		return errorz.Errorf("%s: %w", strings.Join(execCmd.calledCommands, " "), ErrCommandFuncNotSet)
	}

	return execCmd.Func(WithContext(ctx, execCmd), remainingArgs)
}
