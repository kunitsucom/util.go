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
	Stdout io.Writer = os.Stdout
	Stderr io.Writer = os.Stderr
)

//nolint:revive,stylecheck
const (
	UTIL_GO_CLI_TRACE = "UTIL_GO_CLI_TRACE"
	UTIL_GO_CLI_DEBUG = "UTIL_GO_CLI_DEBUG"
)

//nolint:gochecknoglobals
var (
	TraceLog = newLogger(Stderr, UTIL_GO_CLI_TRACE, "TRACE: ")
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
	HelpOptionName    = "help"
)

type (
	Command struct {
		Name        string
		Description string
		SubCommands []*Command
		Options     []Option
		Usage       func(c *Command)

		cmdStack []string
	}

	Option interface {
		GetName() string
		GetShort() string
		GetEnvironment() string
		HasDefault() bool
		getDefault() interface{}
		GetDescription() string
		private()
	}

	StringOption struct {
		Name        string
		Short       string
		Environment string
		Default     *string
		Description string

		value *string
	}

	BoolOption struct {
		Name        string
		Short       string
		Environment string
		Default     *bool
		Description string

		value *bool
	}

	IntOption struct {
		Name        string
		Short       string
		Environment string
		Default     *int
		Description string

		value *int
	}
)

func (cmd *Command) getDescription() string {
	if cmd.Description != "" {
		return cmd.Description
	}
	return fmt.Sprintf("command %q description", strings.Join(cmd.cmdStack, " "))
}

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
						if hasOptionValue(args, i) {
							return nil, errorz.Errorf("%s: %w", arg, ErrMissingOptionValue)
						}
						o.value = ptr(args[i+1])
						i++
						TraceLog.Printf("%s: parse: option: %s: %v", cmd.Name, o.Name, *o.value)
						continue argsLoop
					case optionEqualArg(o, arg):
						o.value = ptr(optionEqualArgExtractValue(arg))
						TraceLog.Printf("%s: parse: option: %s: %v", cmd.Name, o.Name, *o.value)
						continue argsLoop
					}
				case *BoolOption:
					switch {
					case optionArg(o, arg):
						o.value = ptr(true)
						TraceLog.Printf("%s: parse: option: %s: %v", cmd.Name, o.Name, *o.value)
						continue argsLoop
					case optionEqualArg(o, arg):
						optVal, _ := strconv.ParseBool(optionEqualArgExtractValue(arg))
						o.value = &optVal
						TraceLog.Printf("%s: parse: option: %s: %v", cmd.Name, o.Name, *o.value)
						continue argsLoop
					}
				case *IntOption:
					switch {
					case optionArg(o, arg):
						if hasOptionValue(args, i) {
							return nil, errorz.Errorf("%s: %w", arg, ErrMissingOptionValue)
						}
						optVal, _ := strconv.Atoi(args[i+1])
						o.value = &optVal
						i++
						TraceLog.Printf("%s: parse: option: %s: %v", cmd.Name, o.Name, *o.value)
						continue argsLoop
					case optionEqualArg(o, arg):
						optVal, _ := strconv.Atoi(optionEqualArgExtractValue(arg))
						o.value = &optVal
						TraceLog.Printf("%s: parse: option: %s: %v", cmd.Name, o.Name, *o.value)
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

func (cmd *Command) usage() {
	if cmd.Usage != nil {
		cmd.Usage(cmd)
		return
	}
	defaultUsage(cmd, Stderr)
}

//nolint:gocognit,cyclop
func defaultUsage(cmd *Command, w io.Writer) {
	const indent = "    "

	usage := "Usage:" + "\n"
	usage += indent + fmt.Sprintf("%s [options] [arguments]\n", strings.Join(cmd.cmdStack, " ")) + "\n"
	usage += "Description:" + "\n"
	usage += indent + cmd.getDescription() + "\n"

	{
		if len(cmd.SubCommands) > 0 {
			usage += "\n"
			usage += "sub commands:\n"
			for _, subcmd := range cmd.SubCommands {
				usage += indent + fmt.Sprintf("%s: %s", subcmd.Name, subcmd.getDescription()) + "\n"
			}
		}
		if len(cmd.Options) > 0 { //nolint:nestif
			usage += "\n"
			usage += "options:\n"
			for _, opt := range cmd.Options {
				name := opt.GetName()
				short := opt.GetShort()
				env := opt.GetEnvironment()
				usage += indent
				if name != "" {
					usage += fmt.Sprintf("%s%s", longOptionPrefix, name)
				}
				if short != "" {
					if name != "" {
						usage += ", "
					}
					usage += fmt.Sprintf("%s%s", shortOptionPrefix, short)
				}
				if env != "" {
					usage += " ("
					usage += fmt.Sprintf("env: %s", env)
				}
				if opt.HasDefault() {
					if env != "" {
						usage += ", "
					} else {
						usage += " ("
					}
					usage += fmt.Sprintf("default: %v", opt.getDefault())
				}
				if env != "" || opt.HasDefault() {
					usage += ")"
				}

				usage += "\n"
				usage += indent + indent + opt.GetDescription() + "\n"
			}
		}
	}

	_, _ = fmt.Fprint(w, usage)
}

func (cmd *Command) Parse(args []string) (remaining []string, err error) {
	appendHelpOption(cmd)

	if err := cmd.checkSubCommands(); err != nil {
		return nil, errorz.Errorf("failed to check commands: %w", err)
	}

	if err := cmd.checkOptions(); err != nil {
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
