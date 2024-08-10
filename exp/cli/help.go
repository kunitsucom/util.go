package cliz

import "errors"

// HelpOptionName is the option name for help.
const HelpOptionName = "help"

// IsHelp returns whether the error is ErrHelp.
func IsHelp(err error) bool {
	return errors.Is(err, ErrHelp)
}

func (cmd *Command) initAppendHelpOption() {
	// If help option is already set, do nothing.
	if _, ok := cmd.getHelpOption(); !ok {
		cmd.Options = append(cmd.Options, &BoolOption{
			Name:        HelpOptionName,
			Description: "show usage",
			Default:     Default(false),
		})
	}

	// Recursively initialize help option for subcommands.
	for _, subcmd := range cmd.SubCommands {
		subcmd.initAppendHelpOption()
	}
}

func (cmd *Command) getHelpOption() (helpOption *BoolOption, ok bool) {
	// Find help option in the command options.
	for _, opt := range cmd.Options {
		if o, ok := opt.(*BoolOption); ok {
			if o.Name == HelpOptionName {
				return o, true
			}
		}
	}

	return nil, false
}

func (cmd *Command) checkHelp() error {
	TraceLog.Printf("checkHelp: %s", cmd.GetName())

	// If help option is set, show usage and return ErrHelp.
	v, err := cmd.getOptionBool(HelpOptionName)
	if err == nil && v {
		cmd.ShowUsage()
		return ErrHelp
	}

	// Recursively check help option for subcommands.
	for _, subcmd := range cmd.SubCommands {
		if err := subcmd.checkHelp(); err != nil {
			return err
		}
	}
	return nil
}
