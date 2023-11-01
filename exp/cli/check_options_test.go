package cliz

import (
	"errors"
	"testing"
)

func TestCommand_checkOptions(t *testing.T) {
	t.Parallel()
	t.Run("success,", func(t *testing.T) {
		t.Parallel()
		c := &Command{
			Name: "main-cli",
			SubCommands: []*Command{
				{
					Name: "sub-cmd",
					Options: []Option{
						&StringOption{
							Name: "foo",
						},
					},
				},
			},
		}

		if err := c.checkOptions(); err != nil {
			t.Fatalf("❌: %+v", err)
		}
	})

	t.Run("failure,ErrDuplicateOptionName,name", func(t *testing.T) {
		t.Parallel()
		c := &Command{
			Name: "main-cli",
			SubCommands: []*Command{
				{
					Name: "sub-cmd",
					Options: []Option{
						&StringOption{
							Name: "foo",
						},
						&StringOption{
							Name: "foo",
						},
					},
				},
			},
		}

		if err := c.checkOptions(); !errors.Is(err, ErrDuplicateOptionName) {
			t.Fatalf("❌: err != ErrDuplicateOptionName: %+v", err)
		}
	})

	t.Run("failure,ErrDuplicateOptionName,short", func(t *testing.T) {
		t.Parallel()
		c := &Command{
			Name: "main-cli",
			SubCommands: []*Command{
				{
					Name: "sub-cmd",
					Options: []Option{
						&StringOption{
							Short: "f",
						},
						&StringOption{
							Short: "f",
						},
					},
				},
			},
		}

		{
			err := c.checkOptions()
			if !errors.Is(err, ErrDuplicateOptionName) {
				t.Fatalf("❌: err != ErrDuplicateOptionName: %+v", err)
			}
		}
	})

	t.Run("failure,ErrDuplicateOptionName,environment", func(t *testing.T) {
		t.Parallel()
		c := &Command{
			Name: "main-cli",
			SubCommands: []*Command{
				{
					Name: "sub-cmd",
					Options: []Option{
						&StringOption{
							Environment: "FOO",
						},
						&StringOption{
							Environment: "FOO",
						},
					},
				},
			},
		}

		{
			err := c.checkOptions()
			if !errors.Is(err, ErrDuplicateOptionName) {
				t.Fatalf("❌: err != ErrDuplicateOptionName: %+v", err)
			}
		}
	})
}
