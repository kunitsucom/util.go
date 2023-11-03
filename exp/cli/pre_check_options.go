package cliz

import errorz "github.com/kunitsucom/util.go/errors"

func (cmd *Command) preCheckOptions() error {
	// NOTE: duplicate check
	if err := cmd.preCheckDuplicateOptions(make(map[string]bool)); err != nil {
		return errorz.Errorf("%s: %w", cmd.Name, err)
	}

	return nil
}

//nolint:cyclop
func (cmd *Command) preCheckDuplicateOptions(envs map[string]bool) error {
	names := make(map[string]bool)
	shorts := make(map[string]bool)

	for _, opt := range cmd.Options {
		if name := opt.GetName(); name != "" {
			TraceLog.Printf("preCheckDuplicateOptions: %s: option: %s", cmd.Name, name)
			if names[name] {
				err := ErrDuplicateOptionName
				return errorz.Errorf("option: %s%s: %w", longOptionPrefix, name, err)
			}
			names[name] = true
		}

		if short := opt.GetShort(); short != "" {
			if shorts[short] {
				return errorz.Errorf("short option: %s%s: %w", shortOptionPrefix, short, ErrDuplicateOptionName)
			}
			shorts[short] = true
		}

		if env := opt.GetEnvironment(); env != "" {
			if envs[env] {
				return errorz.Errorf("environment: %s: %w", env, ErrDuplicateOptionName)
			}
			envs[env] = true
		}
	}

	for _, subcmd := range cmd.SubCommands {
		if err := subcmd.preCheckDuplicateOptions(envs); err != nil {
			return errorz.Errorf("%s: %w", subcmd.Name, err)
		}
	}

	return nil
}
