package cliz

import errorz "github.com/kunitsucom/util.go/errors"

func (cmd *Command) postCheckOptions() error {
	// NOTE: duplicate check
	if err := cmd.postCheckDuplicateOptions(make(map[string]bool)); err != nil {
		return errorz.Errorf("%s: %w", cmd.Name, err)
	}

	return nil
}

//nolint:cyclop
func (cmd *Command) postCheckDuplicateOptions(envs map[string]bool) error {
	for _, opt := range cmd.Options {
		name := opt.GetName()
		TraceLog.Printf("postCheckDuplicateOptions: %s: option: %s", cmd.Name, name)

		if !opt.HasValue() {
			return errorz.Errorf("option: %s%s: %w", longOptionPrefix, name, ErrOptionRequired)
		}
	}

	for _, subcmd := range cmd.SubCommands {
		if err := subcmd.postCheckDuplicateOptions(envs); err != nil {
			return errorz.Errorf("%s: %w", subcmd.Name, err)
		}
	}

	return nil
}
