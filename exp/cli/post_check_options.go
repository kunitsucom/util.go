package cliz

import errorz "github.com/kunitsucom/util.go/errors"

func (cmd *Command) postCheckOptions() error {
	// NOTE: required options
	if err := cmd.postCheckOptionRequired(); err != nil {
		return errorz.Errorf("%s: %w", cmd.GetName(), err)
	}

	return nil
}

//nolint:cyclop
func (cmd *Command) postCheckOptionRequired() error {
	if len(cmd.calledCommands) > 0 {
		for _, opt := range cmd.Options {
			name := opt.GetName()
			TraceLog.Printf("postCheckOptionRequired: %s: option: %s", cmd.GetName(), name)

			if !opt.HasValue() {
				return errorz.Errorf("option: %s%s: %w", longOptionPrefix, name, ErrOptionRequired)
			}
		}
	}

	for _, subcmd := range cmd.SubCommands {
		if err := subcmd.postCheckOptionRequired(); err != nil {
			return errorz.Errorf("%s: %w", subcmd.GetName(), err)
		}
	}

	return nil
}
