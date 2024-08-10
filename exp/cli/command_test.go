package cliz

import (
	"bytes"
	"context"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"

	errorz "github.com/kunitsucom/util.go/errors"
)

//nolint:paralleltest
func Test_newLogger(t *testing.T) {
	t.Run("success,", func(t *testing.T) {
		t.Setenv("ASDFGHJK", "true")
		buf := bytes.NewBufferString("")
		logger := logger(buf, "ASDFGHJK", "")
		if logger == nil {
			t.Errorf("❌: logger == nil")
		}
	})
}

func TestOption_private(t *testing.T) {
	t.Parallel()
	(&StringOption{}).private()
	(&BoolOption{}).private()
	(&IntOption{}).private()
	(&Float64Option{}).private()
}

//nolint:paralleltest,tparallel
func TestCommand(t *testing.T) {
	const (
		HOST = "HOST"
		PORT = "PORT"
	)

	newCmd := func() *Command {
		return &Command{
			Name:        "my-cli",
			Description: `My awesome CLI tool.`,
			Options: []Option{
				&BoolOption{
					Name:        "version",
					Short:       "v",
					Description: "show version",
					Default:     Default(false),
				},
				&IntOption{
					Name:        "priority",
					Description: "priority number",
					Default:     Default(1),
				},
				&BoolOption{
					Name:        "verbose",
					Environment: "VERBOSE",
					Description: "output verbose",
					Default:     Default(false),
				},
				&StringOption{
					Name:        "annotation",
					Description: "annotate value",
					Default:     Default(""),
				},
				&Float64Option{
					Name:        "ratio",
					Description: "ratio value",
					Default:     Default(0.99),
				},
			},
			SubCommands: []*Command{
				{
					Name:        "sub-cmd",
					Description: `My awesome CLI tool's sub command.`,
					SubCommands: []*Command{
						{
							Name: "sub-sub-cmd",
							RunFunc: func(ctx context.Context, cmd *Command, remainingArgs []string) error {
								called := cmd.GetCalledCommands()
								if !reflect.DeepEqual(called, []string{"my-cli", "sub-cmd", "sub-sub-cmd"}) {
									return errorz.Errorf("unexpected command name: %v", called)
								}
								return nil
							},
							Options: []Option{
								&BoolOption{
									Name:        HelpOptionName,
									Description: "show usage",
									Default:     Default(false),
								},
								&BoolOption{
									Name:        "recursive",
									Short:       "r",
									Environment: "MAIN_CLI_SUB_SUB_RECURSIVE",
									Default:     Default(false),
									Description: "show recursive",
								},
								&StringOption{
									Name:        "annotation",
									Description: "annotate command",
								},
								&StringOption{
									Name:        "foo-string",
									Description: "FOO_STRING",
									Default:     Default("string"),
								},
								&BoolOption{
									Name:        "foo-bool",
									Description: "FOO_BOOL",
									Default:     Default(true),
								},
								&IntOption{
									Name:        "foo-int",
									Description: "FOO_INT",
									Default:     Default(100),
								},
								&Float64Option{
									Name:        "foo-float64",
									Description: "FOO_FLOAT64",
									Default:     Default(1.11),
								},
							},
						},
					},
					Options: []Option{
						&StringOption{
							Name:        "host",
							Environment: HOST,
							Description: "host name",
						},
						&IntOption{
							Name:        "port",
							Environment: PORT,
							Default:     Default(8080),
							Description: "port number",
						},
						&StringOption{
							Name:        "annotation",
							Environment: "MAIN_CLI_SUB_ANNOTATION",
							Default:     Default("annotated-value"),
							Description: "my annotate opt",
						},
						&StringOption{
							Name:        "foo-string",
							Short:       "fs",
							Default:     Default("foo-value"),
							Description: "my foo string opt",
						},
						&BoolOption{
							Name:        "foo-bool",
							Short:       "fb",
							Default:     Default(true),
							Description: "my foo bool opt",
						},
						&IntOption{
							Name:        "foo-int",
							Short:       "fi",
							Default:     Default(100),
							Description: "my foo int opt",
						},
						&Float64Option{
							Name:        "foo-float64",
							Short:       "ff",
							Default:     Default(float64(1.11)),
							Description: "my foo float64 opt",
						},
						&StringOption{
							Name: "bar-string",
						},
						&BoolOption{
							Name: "bar-bool",
						},
						&IntOption{
							Name: "bar-int",
						},
						&Float64Option{
							Name: "bar-float64",
						},
					},
				},
			},
		}
	}

	t.Run("success,", func(t *testing.T) {
		t.Setenv(PORT, "8000")

		c := newCmd()
		args := []string{"my-cli", "-v", "--priority=1", "--annotation=4main", "--verbose=false", "--ratio", "0.98", "sub-cmd", "--host", "localhost", "--port", "8081", "--annotation=4sub", "--bar-string=bar", "--bar-bool=true", "--bar-int=100", "--bar-float64=1.11", "sub-sub-cmd", "--annotation=4subsub", "path/to/source", "path/to/destination", "--recursive", "--", "path/to/abc"}

		if _, err := c.Parse(args[1:]); err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}
		remaining, err := c.Parse(args[1:])
		if err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}

		if (*Command)(nil).Next() != nil {
			t.Errorf("❌: expect != actual: %v != %v", nil, (*Command)(nil).Next())
		}
		if (&Command{}).Next() != nil {
			t.Errorf("❌: expect != actual: %v != %v", nil, (&Command{}).Next())
		}
		switch c.GetName() {
		case "my-cli":
			switch c := c.Next(); c.GetName() {
			case "sub-cmd":
				switch c := c.Next(); c.GetName() {
				case "sub-sub-cmd":
					switch c := c.Next(); c.GetName() {
					case "":
						// OK
					default:
						t.Errorf("❌: expect != actual: %v != %v", "", c.GetName())
					}
				default:
					t.Errorf("❌: expect != actual: %v != %v", "sub-sub-cmd", c.GetName())
				}
			default:
				t.Errorf("❌: expect != actual: %v != %v", "sub-cmd", c.GetName())
			}
		default:
			t.Errorf("❌: expect != actual: %v != %v", "my-cli", c.GetName())
		}

		if expect, actual := []string{"my-cli", "sub-cmd", "sub-sub-cmd"}, c.GetCalledCommands(); !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}

		if expect, actual := []string{"path/to/source", "path/to/destination", "path/to/abc"}, remaining; !reflect.DeepEqual(expect, actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}

		version, err := c.GetOptionBool("version")
		if err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}
		if !version {
			t.Errorf("❌: %v: unexpected value: %s=%t", "version", args, version)
		}

		annotation4main, err := c.GetOptionString("annotation")
		if err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}
		if annotation4main != "4subsub" {
			t.Errorf("❌: %v: unexpected value: %s=%s", "annotation", args, annotation4main)
		}

		ratio, err := c.GetOptionFloat64("ratio")
		if err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}
		if ratio != 0.98 {
			t.Errorf("❌: %v: unexpected value: %s=%f", "ratio", args, ratio)
		}

		host, err := c.GetOptionString("host")
		if err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}
		if host != "localhost" {
			t.Errorf("❌: %v: unexpected value: %s=%s", "host", args, host)
		}

		port, err := c.GetOptionInt("port")
		if err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}
		if port != 8081 {
			t.Errorf("❌: %v: unexpected value: %s=%d", "port", args, port)
		}

		priority, err := c.GetOptionInt("priority")
		if err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}
		if priority != 1 {
			t.Errorf("❌: %v: unexpected value: %s=%d", "priority", args, priority)
		}

		verbose, err := c.GetOptionBool("verbose")
		if err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}
		if verbose {
			t.Errorf("❌: %v: unexpected value: %s=%t", "verbose", args, verbose)
		}

		recursive, err := c.GetOptionBool("recursive")
		if err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}
		if !recursive {
			t.Errorf("❌: %v: unexpected value: %s=%t", "recursive", args, recursive)
		}

		fs, err := c.GetOptionString("foo-string")
		if err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}
		if fs != "string" {
			t.Errorf("❌: %v: unexpected value: %s=%s", "foo", args, fs)
		}

		fb, err := c.GetOptionBool("foo-bool")
		if err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}
		if !fb {
			t.Errorf("❌: %v: unexpected value: %s=%t", "bar", args, fb)
		}

		fi, err := c.GetOptionInt("foo-int")
		if err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}
		if fi != 100 {
			t.Errorf("❌: %v: unexpected value: %s=%d", "baz", args, fi)
		}

		ff, err := c.GetOptionFloat64("foo-float64")
		if err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}
		if ff != 1.11 {
			t.Errorf("❌: %v: unexpected value: %s=%f", "baz", args, ff)
		}
		if err := c.Run(context.Background(), args[1:]); err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}
	})

	t.Run("success,ErrHelp", func(t *testing.T) {
		c := newCmd()

		const golden = `Usage:
    my-cli sub-cmd [options] <subcommand>

Description:
    My awesome CLI tool's sub command.

sub commands:
    sub-sub-cmd: command "my-cli sub-cmd sub-sub-cmd" description

options:
    --host (env: HOST, required)
        host name
    --port (env: PORT, default: 8080)
        port number
    --annotation (env: MAIN_CLI_SUB_ANNOTATION, default: annotated-value)
        my annotate opt
    --foo-string, -fs (default: foo-value)
        my foo string opt
    --foo-bool, -fb (default: true)
        my foo bool opt
    --foo-int, -fi (default: 100)
        my foo int opt
    --foo-float64, -ff (default: 1.11)
        my foo float64 opt
    --bar-string (required)
        string value
    --bar-bool (required)
        bool value
    --bar-int (required)
        int value
    --bar-float64 (required)
        float64 value
    --help (default: false)
        show usage
`

		backup := Stderr
		t.Cleanup(func() { Stderr = backup })
		Stderr = bytes.NewBufferString("")

		args := []string{"my-cli", "sub-cmd", "not-subcmd", "--help", "sub-sub-cmd"}
		remaining, err := c.Parse(args)
		if !IsHelp(err) {
			t.Errorf("❌: expect != actual: %v != %+v", ErrHelp, err)
		}
		if len(remaining) > 0 {
			t.Errorf("❌: expect != actual: %v != %v", []string{}, remaining)
		}
		if expect, actual := golden, Stderr.(*bytes.Buffer).String(); expect != actual { //nolint:forcetypeassert
			t.Errorf("❌: expect != actual:\n--- EXPECT\n%v\n--- ACTUAL\n%v", expect, actual)
		}
		t.Logf("✅: %v: remaining: %+v", args, remaining)
	})

	t.Run("success,Usage", func(t *testing.T) {
		backup := Stderr
		t.Cleanup(func() { Stderr = backup })
		Stderr = bytes.NewBufferString("")

		c := newCmd()
		c.Usage = "my awesome CLI tool"
		remaining, err := c.Parse([]string{"my-cli", "--help"})
		if !IsHelp(err) {
			t.Fatalf("❌: expect != actual: %v != %+v", ErrHelp, err)
		}
		if len(remaining) > 0 {
			t.Errorf("❌: expect != actual: %v != %v", []string{}, remaining)
		}
		if !strings.Contains(Stderr.(*bytes.Buffer).String(), c.Usage) { //nolint:forcetypeassert
			t.Errorf("❌: not contains: %v", c.Usage)
		}
	})

	t.Run("success,UsageFunc", func(t *testing.T) {
		c := newCmd()
		called := false
		c.UsageFunc = func(c *Command) {
			called = true
		}

		remaining, err := c.Parse([]string{"my-cli", "--help"})
		if !IsHelp(err) {
			t.Errorf("❌: expect != actual: %v != %+v", ErrHelp, err)
		}
		if len(remaining) > 0 {
			t.Errorf("❌: expect != actual: %v != %v", []string{}, remaining)
		}
		if !called {
			t.Errorf("❌: expect != actual: %v != %v", true, called)
		}
	})

	t.Run("failure,ErrMissingOptionValue", func(t *testing.T) {
		t.Parallel()

		c := newCmd()

		if _, err := c.Parse([]string{"my-cli", "sub-cmd", "--host"}); !errors.Is(err, ErrMissingOptionValue) {
			t.Errorf("❌: expect != actual: %v != %+v", ErrMissingOptionValue, err)
		}

		if _, err := c.Parse([]string{"my-cli", "--priority"}); !errors.Is(err, ErrMissingOptionValue) {
			t.Errorf("❌: expect != actual: %v != %+v", ErrMissingOptionValue, err)
		}

		if _, err := c.Parse([]string{"my-cli", "--ratio"}); !errors.Is(err, ErrMissingOptionValue) {
			t.Errorf("❌: expect != actual: %v != %+v", ErrMissingOptionValue, err)
		}
	})

	t.Run("failure,ErrInvalidOptionType", func(t *testing.T) {
		t.Parallel()

		c := newCmd()
		c.Options = append(c.Options, &testOption{
			Name: "foo",
		})

		if _, err := c.Parse([]string{"my-cli", "--foo", "string"}); !errors.Is(err, ErrInvalidOptionType) {
			t.Errorf("❌: expect != actual: %v != %+v", ErrInvalidOptionType, err)
		}
	})

	t.Run("failure,ErrUnknownOption", func(t *testing.T) {
		t.Parallel()

		c := newCmd()

		if _, err := c.Parse([]string{"my-cli", "--foo", "string"}); !errors.Is(err, ErrUnknownOption) {
			t.Errorf("❌: expect != actual: %v != %+v", ErrUnknownOption, err)
		}
	})

	t.Run("failure,ErrDuplicateSubCommand", func(t *testing.T) {
		t.Parallel()

		c := newCmd()
		c.SubCommands = append(c.SubCommands, c.SubCommands...)

		if _, err := c.Parse([]string{"my-cli", "sub-cmd"}); !errors.Is(err, ErrDuplicateSubCommand) {
			t.Errorf("❌: expect != actual: %v != %+v", ErrDuplicateSubCommand, err)
		}
	})

	t.Run("failure,ErrDuplicateOptionName", func(t *testing.T) {
		t.Parallel()

		c := newCmd()
		c.Options = append(c.Options, c.Options...)

		if _, err := c.Parse([]string{"my-cli", "sub-cmd"}); !errors.Is(err, ErrDuplicateOptionName) {
			t.Errorf("❌: expect != actual: %v != %+v", ErrDuplicateOptionName, err)
		}
	})

	t.Run("failure,ErrInvalidOptionType,Environment", func(t *testing.T) {
		t.Parallel()

		c := newCmd()
		c.Options = append(c.Options, &testOption{
			Name:        "foo",
			Environment: "FOO",
		})

		if _, err := c.Parse([]string{"my-cli", "sub-cmd"}); !errors.Is(err, ErrInvalidOptionType) {
			t.Errorf("❌: expect != actual: %v != %+v", ErrInvalidOptionType, err)
		}
	})

	t.Run("failure,ErrInvalidOptionType,Default", func(t *testing.T) {
		t.Parallel()

		c := newCmd()
		c.SubCommands[0].Options = append(c.SubCommands[0].Options, &testOption{
			Name:    "test-option",
			Default: Default("test-option"),
		})

		if _, err := c.Parse([]string{"my-cli", "sub-cmd", "--host=host"}); !errors.Is(err, ErrInvalidOptionType) {
			t.Errorf("❌: expect != actual: %v != %+v", ErrInvalidOptionType, err)
		}
	})

	t.Run("failure,ErrOptionRequired", func(t *testing.T) {
		t.Parallel()

		c := newCmd()
		if _, err := c.Parse([]string{"my-cli", "sub-cmd", "--bar-string=bar", "--bar-bool=true", "--bar-int=100", "--bar-float64=1.11"}); !errors.Is(err, ErrOptionRequired) {
			t.Errorf("❌: expect != actual: %v != %+v", ErrOptionRequired, err)
		}
	})

	t.Run("failure,ErrOptionRequired", func(t *testing.T) {
		t.Parallel()

		c := newCmd()
		if _, err := c.Parse([]string{"my-cli", "sub-cmd", "--host=host", "--bar-string=bar", "--bar-bool=INVALID", "--bar-int=100", "--bar-float64=1.11"}); !errorz.Contains(err, "invalid syntax") {
			t.Errorf("❌: expect != actual: err != \"invalid syntax\": %+v", err)
		}
		if _, err := c.Parse([]string{"my-cli", "sub-cmd", "--host=host", "--bar-string=bar", "--bar-bool=true", "--bar-int", "INVALID", "--bar-float64=1.11"}); !errorz.Contains(err, "invalid syntax") {
			t.Errorf("❌: expect != actual: err != \"invalid syntax\": %+v", err)
		}
		if _, err := c.Parse([]string{"my-cli", "sub-cmd", "--host=host", "--bar-string=bar", "--bar-bool=true", "--bar-int=INVALID", "--bar-float64=1.11"}); !errorz.Contains(err, "invalid syntax") {
			t.Errorf("❌: expect != actual: err != \"invalid syntax\": %+v", err)
		}
		if _, err := c.Parse([]string{"my-cli", "sub-cmd", "--host=host", "--bar-string=bar", "--bar-bool=true", "--bar-int=100", "--bar-float64", "INVALID"}); !errorz.Contains(err, "invalid syntax") {
			t.Errorf("❌: expect != actual: err != \"invalid syntax\": %+v", err)
		}
		if _, err := c.Parse([]string{"my-cli", "sub-cmd", "--host=host", "--bar-string=bar", "--bar-bool=true", "--bar-int=100", "--bar-float64=INVALID"}); !errorz.Contains(err, "invalid syntax") {
			t.Errorf("❌: expect != actual: err != \"invalid syntax\": %+v", err)
		}
	})
}

func TestCommand_GetName(t *testing.T) {
	t.Parallel()

	t.Run("success,nil", func(t *testing.T) {
		t.Parallel()
		expected := ""
		actual := (*Command)(nil).GetName()
		if expected != actual {
			t.Errorf("❌: expected != actual: %v != %v", expected, actual)
		}
	})
}

func TestCommand_IsCommand(t *testing.T) {
	t.Parallel()

	t.Run("success,nil", func(t *testing.T) {
		t.Parallel()
		expected := false
		actual := (*Command)(nil).IsCommand("")
		if expected != actual {
			t.Errorf("❌: expected != actual: %v != %v", expected, actual)
		}
	})

	t.Run("success,alias", func(t *testing.T) {
		t.Parallel()
		alias := "short"
		cmd := &Command{
			Name:    "my-cli",
			Aliases: []string{"short"},
		}
		expected := true
		actual := cmd.IsCommand(alias)
		if expected != actual {
			t.Errorf("❌: expected != actual: %v != %v", expected, actual)
		}
	})
}

func TestCommand_getDescription(t *testing.T) {
	t.Parallel()

	t.Run("success,nil", func(t *testing.T) {
		t.Parallel()
		expected := ""
		actual := (*Command)(nil).getDescription()
		if expected != actual {
			t.Errorf("❌: expected != actual: %v != %v", expected, actual)
		}
	})

	t.Run("success,Description", func(t *testing.T) {
		t.Parallel()
		cmd := &Command{
			Description: "my awesome CLI tool",
		}
		expected := "my awesome CLI tool"
		actual := cmd.getDescription()
		if expected != actual {
			t.Errorf("❌: expected != actual: %v != %v", expected, actual)
		}
	})

	t.Run("success,default", func(t *testing.T) {
		t.Parallel()
		cmd := &Command{
			Name: "my-cli",
		}
		expected := `command "my-cli" description`
		actual := cmd.getDescription()
		if expected != actual {
			t.Errorf("❌: expected != actual: %v != %v", expected, actual)
		}
	})
}

func TestCommand_getSubcommand(t *testing.T) {
	t.Parallel()

	t.Run("success,nil", func(t *testing.T) {
		t.Parallel()
		if actual := (*Command)(nil).getSubcommand(""); nil != actual {
			t.Errorf("❌: expected != actual: %v != %v", nil, actual)
		}
	})
}

func TestCommand_GetCalledCommands(t *testing.T) {
	t.Parallel()

	t.Run("success,nil", func(t *testing.T) {
		t.Parallel()
		if actual := (*Command)(nil).GetCalledCommands(); actual != nil {
			t.Errorf("❌: expected != actual: %v != %v", nil, actual)
		}
	})
}

func TestCommand_Run(t *testing.T) {
	t.Parallel()

	t.Run("success,Run,ErrHelp", func(t *testing.T) {
		t.Parallel()
		args := []string{"my-cli", "--help"}
		c := &Command{}
		if err := c.Run(context.Background(), args[1:]); !IsHelp(err) {
			t.Errorf("❌: expect != actual: %v != %+v", ErrHelp, err)
		}
	})

	t.Run("success,Run,", func(t *testing.T) {
		t.Parallel()
		c := &Command{
			Name: "my-cli",
			PreRunFunc: func(ctx context.Context, cmd *Command, remainingArgs []string) error {
				return nil
			},
			RunFunc: func(ctx context.Context, cmd *Command, remainingArgs []string) error {
				return nil
			},
			PostRunFunc: func(ctx context.Context, cmd *Command, remainingArgs []string) error {
				return nil
			},
			SubCommands: []*Command{
				{
					Name: "sub-cmd",
				},
			},
		}
		if err := c.Run(context.Background(), []string{"my-cli", "sub-cmd"}[1:]); !errors.Is(err, ErrCommandFuncNotSet) {
			t.Errorf("❌: expect != actual: %v != %+v", ErrCommandFuncNotSet, err)
		}
		if err := c.Run(context.Background(), []string{"my-cli"}[1:]); err != nil {
			t.Errorf("❌: err != nil: %v != %+v", nil, err)
		}
	})

	t.Run("failure,Run,PreRunFunc", func(t *testing.T) {
		t.Parallel()
		c := &Command{
			Name: "my-cli",
			PreRunFunc: func(ctx context.Context, cmd *Command, remainingArgs []string) error {
				return io.ErrUnexpectedEOF
			},
			RunFunc: func(ctx context.Context, cmd *Command, remainingArgs []string) error {
				return nil
			},
			PostRunFunc: func(ctx context.Context, cmd *Command, remainingArgs []string) error {
				return nil
			},
		}
		if err := c.Run(context.Background(), []string{"my-cli"}[1:]); !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Errorf("❌: expect != actual: %v != %+v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("failure,Run,RunFunc", func(t *testing.T) {
		t.Parallel()
		c := &Command{
			Name: "my-cli",
			PreRunFunc: func(ctx context.Context, cmd *Command, remainingArgs []string) error {
				return nil
			},
			RunFunc: func(ctx context.Context, cmd *Command, remainingArgs []string) error {
				return io.ErrUnexpectedEOF
			},
			PostRunFunc: func(ctx context.Context, cmd *Command, remainingArgs []string) error {
				return nil
			},
		}
		if err := c.Run(context.Background(), []string{"my-cli"}[1:]); !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Errorf("❌: expect != actual: %v != %+v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("failure,Run,PostRunFunc", func(t *testing.T) {
		t.Parallel()
		c := &Command{
			Name: "my-cli",
			PreRunFunc: func(ctx context.Context, cmd *Command, remainingArgs []string) error {
				return nil
			},
			RunFunc: func(ctx context.Context, cmd *Command, remainingArgs []string) error {
				return nil
			},
			PostRunFunc: func(ctx context.Context, cmd *Command, remainingArgs []string) error {
				return io.ErrUnexpectedEOF
			},
		}
		if err := c.Run(context.Background(), []string{"my-cli"}[1:]); !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Errorf("❌: expect != actual: %v != %+v", io.ErrUnexpectedEOF, err)
		}
	})
}
