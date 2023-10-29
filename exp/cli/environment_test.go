//nolint:testpackage
package cliz

import (
	"errors"
	"testing"

	errorz "github.com/kunitsucom/util.go/errors"
)

//nolint:paralleltest
func TestCommand_loadEnvironment(t *testing.T) {
	const (
		FOO = "FOO"
		BAR = "BAR"
		BAZ = "BAZ"
	)
	t.Run("success,ALL", func(t *testing.T) {
		c := &Command{
			Name: "main-cli",
			Options: []Option{
				&StringOption{
					Name:        "foo",
					Environment: FOO,
				},
				&BoolOption{
					Name:        "bar",
					Environment: BAR,
				},
				&IntOption{
					Name:        "baz",
					Environment: BAZ,
				},
			},
		}
		t.Setenv(FOO, "foo")
		t.Setenv(BAR, "true")
		t.Setenv(BAZ, "100")
		if err := c.loadEnvironments(); err != nil {
			t.Fatalf("❌: c.loadEnvironments: err != nil: %+v", err)
		}
	})

	t.Run("failure,BoolOption", func(t *testing.T) {
		c := &Command{
			Name: "main-cli",
			Options: []Option{
				&BoolOption{
					Name:        "bar",
					Environment: BAR,
				},
			},
		}
		t.Setenv(BAR, "string")
		if err := c.loadEnvironments(); !errorz.Contains(err, "invalid syntax") {
			t.Fatalf("❌: c.loadEnvironments: err != \"invalid syntax\": %+v", err)
		}
	})

	t.Run("failure,IntOption", func(t *testing.T) {
		c := &Command{
			Name: "main-cli",
			Options: []Option{
				&IntOption{
					Name:        "baz",
					Environment: BAZ,
				},
			},
		}
		t.Setenv(BAZ, "string")
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
							Name:        "foo",
							Environment: FOO,
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
