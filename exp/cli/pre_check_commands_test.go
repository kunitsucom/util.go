package cliz

import (
	"errors"
	"testing"
)

func TestCommand_checkCommands(t *testing.T) {
	t.Parallel()
	t.Run("success,", func(t *testing.T) {
		t.Parallel()
		c := &Command{
			Name: "main-cli",
			SubCommands: []*Command{
				{
					Name: "sub-cmd",
				},
			},
		}

		if err := c.preCheckSubCommands(); err != nil {
			t.Fatalf("❌: %+v", err)
		}
	})

	t.Run("failure,", func(t *testing.T) {
		t.Parallel()

		c := &Command{
			Name: "main-cli",
			SubCommands: []*Command{
				{
					Name: "sub-cmd",
					SubCommands: []*Command{
						{
							Name: "sub-sub-cmd",
						},
						{
							Name: "sub-sub-cmd",
						},
					},
				},
			},
		}

		if err := c.preCheckSubCommands(); !errors.Is(err, ErrDuplicateSubCommand) {
			t.Fatalf("❌: err != ErrDuplicateSubCommand: %+v", err)
		}
	})
}
