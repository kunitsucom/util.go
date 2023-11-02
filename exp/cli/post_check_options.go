package cliz

import errorz "github.com/kunitsucom/util.go/errors"

func (cmd *Command) postCheckOptions() error {
	// NOTE: required check
	if err := cmd.postCheckOptionRequired(); err != nil {
		return errorz.Errorf("%s: %w", cmd.Name, err)
	}

	return nil
}

//nolint:cyclop
func (cmd *Command) postCheckOptionRequired() error {
	if len(cmd.cmdStack) > 0 {
		for _, opt := range cmd.Options {
			name := opt.GetName()
			TraceLog.Printf("postCheckOptionRequired: %s: option: %s", cmd.Name, name)

			if !opt.HasValue() {
				return errorz.Errorf("option: %s%s: %w", longOptionPrefix, name, ErrOptionRequired)
			}
		}
	}

	for _, subcmd := range cmd.SubCommands {
		if err := subcmd.postCheckOptionRequired(); err != nil {
			return errorz.Errorf("%s: %w", subcmd.Name, err)
		}
	}

	return nil
}
