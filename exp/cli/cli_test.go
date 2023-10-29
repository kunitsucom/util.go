package cliz

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

//nolint:paralleltest
func Test_newLogger(t *testing.T) {
	t.Run("success,", func(t *testing.T) {
		t.Setenv("ASDFGHJK", "true")
		buf := bytes.NewBufferString("")
		logger := newLogger(buf, "ASDFGHJK", "")
		if logger == nil {
			t.Errorf("❌: logger == nil")
		}
	})
}

//nolint:paralleltest,tparallel
func TestCommand(t *testing.T) {
	const (
		HOST = "HOST"
		PORT = "PORT"
	)

	(&StringOption{}).private()
	(&BoolOption{}).private()
	(&IntOption{}).private()

	newCmd := func() *Command {
		return &Command{
			Name:        "main-cli",
			Description: `My awesome CLI tool.`,
			SubCommands: []*Command{
				{
					Name:        "sub-cmd",
					Description: `My awesome CLI tool's sub command.`,
					SubCommands: []*Command{
						{
							Name: "sub-sub-cmd",
							Options: []Option{
								&BoolOption{
									Name:        HelpOptionName,
									Description: "show usage",
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
									Name:        "foo",
									Description: "FOO",
									Default:     Default("foo"),
								},
								&BoolOption{
									Name:        "bar",
									Description: "BAR",
									Default:     Default(true),
								},
								&IntOption{
									Name:        "baz",
									Description: "BAZ",
									Default:     Default(100),
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
							Description: "port number",
							Default:     Default(8080),
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
						&IntOption{
							Name:        "foo-int",
							Short:       "fi",
							Default:     Default(100),
							Description: "my foo int opt",
						},
						&BoolOption{
							Name:        "foo-bool",
							Short:       "fb",
							Default:     Default(true),
							Description: "my foo bool opt",
						},
						&StringOption{
							Name: "bar-string",
						},
						&IntOption{
							Name: "bar-int",
						},
						&BoolOption{
							Name: "bar-bool",
						},
					},
				},
			},
			Options: []Option{
				&BoolOption{
					Name:        "version",
					Short:       "v",
					Description: "show version",
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
				},
				&StringOption{
					Name:        "annotation",
					Description: "show version",
				},
			},
		}
	}

	t.Run("success,", func(t *testing.T) {
		t.Setenv(PORT, "8000")

		c := newCmd()
		args := []string{"main-cli", "-v", "--priority=1", "--annotation=4main", "--verbose=false", "sub-cmd", "--host", "localhost", "--port", "8081", "--annotation=4sub", "sub-sub-cmd", "--annotation=4subsub", "path/to/source", "path/to/destination", "--recursive", "--", "path/to/abc"}
		remaining, err := c.Parse(args[1:])
		if err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}
		t.Logf("✅: %v: remaining: %+v", args, remaining)

		version, err := c.GetBoolOption("version")
		if err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}
		if !version {
			t.Errorf("❌: %v: unexpected value: %s=%t", "version", args, version)
		}

		annotation4main, err := c.GetStringOption("annotation")
		if err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}
		if annotation4main != "4subsub" {
			t.Errorf("❌: %v: unexpected value: %s=%s", "annotation", args, annotation4main)
		}

		host, err := c.GetStringOption("host")
		if err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}
		if host != "localhost" {
			t.Errorf("❌: %v: unexpected value: %s=%s", "host", args, host)
		}

		port, err := c.GetIntOption("port")
		if err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}
		if port != 8081 {
			t.Errorf("❌: %v: unexpected value: %s=%d", "port", args, port)
		}

		priority, err := c.GetIntOption("priority")
		if err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}
		if priority != 1 {
			t.Errorf("❌: %v: unexpected value: %s=%d", "priority", args, priority)
		}

		verbose, err := c.GetBoolOption("verbose")
		if err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}
		if verbose {
			t.Errorf("❌: %v: unexpected value: %s=%t", "verbose", args, verbose)
		}

		recursive, err := c.GetBoolOption("recursive")
		if err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}
		if !recursive {
			t.Errorf("❌: %v: unexpected value: %s=%t", "recursive", args, recursive)
		}

		foo, err := c.GetStringOption("foo")
		if err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}
		if foo != "foo" {
			t.Errorf("❌: %v: unexpected value: %s=%s", "foo", args, foo)
		}

		bar, err := c.GetBoolOption("bar")
		if err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}
		if !bar {
			t.Errorf("❌: %v: unexpected value: %s=%t", "bar", args, bar)
		}

		baz, err := c.GetIntOption("baz")
		if err != nil {
			t.Fatalf("❌: %v: %+v", args, err)
		}
		if baz != 100 {
			t.Errorf("❌: %v: unexpected value: %s=%d", "baz", args, baz)
		}

		if expect, actual := "path/to/source path/to/destination path/to/abc", strings.Join(remaining, " "); expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("success,ErrHelp", func(t *testing.T) {
		c := newCmd()

		const golden = `Usage:
    main-cli sub-cmd [options] [arguments]

Description:
    My awesome CLI tool's sub command.

sub commands:
    sub-sub-cmd: command "main-cli sub-cmd sub-sub-cmd" description

options:
    --host (env: HOST)
        host name
    --port (env: PORT, default: 8080)
        port number
    --annotation (env: MAIN_CLI_SUB_ANNOTATION, default: annotated-value)
        my annotate opt
    --foo-string, -fs (default: foo-value)
        my foo string opt
    --foo-int, -fi (default: 100)
        my foo int opt
    --foo-bool, -fb (default: true)
        my foo bool opt
    --bar-string
        string value
    --bar-int
        int value
    --bar-bool
        bool value
    --help
        show usage
`

		Stderr = bytes.NewBufferString("")

		args := []string{"main-cli", "sub-cmd", "not-subcmd", "--help", "sub-sub-cmd"}
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
		c := newCmd()
		called := false
		c.Usage = func(c *Command) {
			called = true
		}
		remaining, err := c.Parse([]string{"main-cli", "--help"})
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

		if _, err := c.Parse([]string{"main-cli", "sub-cmd", "--host"}); !errors.Is(err, ErrMissingOptionValue) {
			t.Errorf("❌: expect != actual: %v != %+v", ErrMissingOptionValue, err)
		}

		if _, err := c.Parse([]string{"main-cli", "--priority"}); !errors.Is(err, ErrMissingOptionValue) {
			t.Errorf("❌: expect != actual: %v != %+v", ErrMissingOptionValue, err)
		}
	})

	t.Run("failure,ErrInvalidOptionType", func(t *testing.T) {
		t.Parallel()

		c := newCmd()
		c.Options = append(c.Options, &testOption{
			Name: "foo",
		})

		if _, err := c.Parse([]string{"main-cli", "--foo", "string"}); !errors.Is(err, ErrInvalidOptionType) {
			t.Errorf("❌: expect != actual: %v != %+v", ErrInvalidOptionType, err)
		}
	})

	t.Run("failure,ErrUnknownOption", func(t *testing.T) {
		t.Parallel()

		c := newCmd()

		if _, err := c.Parse([]string{"main-cli", "--foo", "string"}); !errors.Is(err, ErrUnknownOption) {
			t.Errorf("❌: expect != actual: %v != %+v", ErrUnknownOption, err)
		}
	})

	t.Run("failure,ErrDuplicateSubCommand", func(t *testing.T) {
		t.Parallel()

		c := newCmd()
		c.SubCommands = append(c.SubCommands, c.SubCommands...)

		if _, err := c.Parse([]string{"main-cli", "sub-cmd"}); !errors.Is(err, ErrDuplicateSubCommand) {
			t.Errorf("❌: expect != actual: %v != %+v", ErrDuplicateSubCommand, err)
		}
	})

	t.Run("failure,ErrDuplicateOptionName", func(t *testing.T) {
		t.Parallel()

		c := newCmd()
		c.Options = append(c.Options, c.Options...)

		if _, err := c.Parse([]string{"main-cli", "sub-cmd"}); !errors.Is(err, ErrDuplicateOptionName) {
			t.Errorf("❌: expect != actual: %v != %+v", ErrDuplicateOptionName, err)
		}
	})

	t.Run("failure,ErrInvalidOptionType", func(t *testing.T) {
		t.Parallel()

		c := newCmd()
		c.Options = append(c.Options, &testOption{
			Name:        "foo",
			Environment: "FOO",
		})

		if _, err := c.Parse([]string{"main-cli", "sub-cmd"}); !errors.Is(err, ErrInvalidOptionType) {
			t.Errorf("❌: expect != actual: %v != %+v", ErrInvalidOptionType, err)
		}
	})
}
