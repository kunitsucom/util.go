package cliz

import errorz "github.com/kunitsucom/util.go/errors"

func (cmd *Command) preCheckSubCommands() error {
	if err := cmd.preCheckDuplicateSubCommands(); err != nil {
		return errorz.Errorf("%s: %w", cmd.Name, err)
	}

	return nil
}

func (cmd *Command) preCheckDuplicateSubCommands() error {
	names := make(map[string]bool)

	for _, cmd := range cmd.SubCommands {
		name := cmd.Name

		TraceLog.Printf("preCheckDuplicateSubCommands: %s", name)

		if name != "" && names[name] {
			return errorz.Errorf("sub command: %s: %w", name, ErrDuplicateSubCommand)
		}
		names[name] = true
	}

	for _, subcmd := range cmd.SubCommands {
		if err := subcmd.preCheckDuplicateSubCommands(); err != nil {
			return errorz.Errorf("%w", err)
		}
	}
	return nil
}
