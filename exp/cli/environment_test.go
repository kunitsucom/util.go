//nolint:testpackage
package cliz

import (
	"errors"
	"testing"

	errorz "github.com/kunitsucom/util.go/errors"
)

//nolint:paralleltest
func TestCommand_loadEnvironment(t *testing.T) {
	//nolint:stylecheck
	const (
		fooString   = "foo-string"
		fooBool     = "foo-bool"
		fooInt      = "foo-int"
		fooFloat64  = "foo-float64"
		FOO_STRING  = "FOO_STRING"
		FOO_BOOL    = "BAR"
		FOO_INT     = "BAZ"
		FOO_FLOAT64 = "FOO_FLOAT64"
	)
	t.Run("success,ALL", func(t *testing.T) {
		c := &Command{
			Name: "main-cli",
			Options: []Option{
				&StringOption{
					Name:        fooString,
					Environment: FOO_STRING,
				},
				&BoolOption{
					Name:        fooBool,
					Environment: FOO_BOOL,
				},
				&IntOption{
					Name:        fooInt,
					Environment: FOO_INT,
				},
				&Float64Option{
					Name:        fooFloat64,
					Environment: FOO_FLOAT64,
				},
			},
		}
		t.Setenv(FOO_STRING, fooString)
		t.Setenv(FOO_BOOL, "true")
		t.Setenv(FOO_INT, "100")
		t.Setenv(FOO_FLOAT64, "1.11")
		if err := c.loadEnvironments(); err != nil {
			t.Fatalf("❌: c.loadEnvironments: err != nil: %+v", err)
		}
	})

	t.Run("failure,BoolOption", func(t *testing.T) {
		c := &Command{
			Name: "main-cli",
			Options: []Option{
				&BoolOption{
					Name:        fooBool,
					Environment: FOO_BOOL,
				},
			},
		}
		t.Setenv(FOO_BOOL, "string")
		if err := c.loadEnvironments(); !errorz.Contains(err, "invalid syntax") {
			t.Fatalf("❌: c.loadEnvironments: err != \"invalid syntax\": %+v", err)
		}
	})

	t.Run("failure,IntOption", func(t *testing.T) {
		c := &Command{
			Name: "main-cli",
			Options: []Option{
				&IntOption{
					Name:        fooInt,
					Environment: FOO_INT,
				},
			},
		}
		t.Setenv(FOO_INT, "string")
		if err := c.loadEnvironments(); !errorz.Contains(err, "invalid syntax") {
			t.Fatalf("❌: c.loadEnvironments: err != \"invalid syntax\": %+v", err)
		}
	})

	t.Run("failure,Float64Option", func(t *testing.T) {
		c := &Command{
			Name: "main-cli",
			Options: []Option{
				&Float64Option{
					Name:        "foo-float64",
					Environment: FOO_FLOAT64,
				},
			},
		}
		t.Setenv(FOO_FLOAT64, "string")
		if err := c.loadEnvironments(); !errorz.Contains(err, "invalid syntax") {
			t.Fatalf("❌: c.loadEnvironments: err != \"invalid syntax\": %+v", err)
		}
	})

	t.Run("failure,ErrInvalidOptionType", func(t *testing.T) {
		c := &Command{
			Name: "main-cli",
			SubCommands: []*Command{
				{
					Name: "sub-cmd",
					Options: []Option{
						&testOption{
							Name:        "foo-string",
							Environment: FOO_STRING,
						},
					},
				},
			},
		}
		if err := c.loadEnvironments(); !errors.Is(err, ErrInvalidOptionType) {
			t.Fatalf("❌: c.loadEnvironments: err != ErrInvalidOptionType: %+v", err)
		}
	})
}
