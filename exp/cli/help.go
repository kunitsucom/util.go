package cliz

import "errors"

// HelpOptionName is the option name for help.
const HelpOptionName = "help"

// IsHelp returns whether the error is ErrHelp.
func IsHelp(err error) bool {
	return errors.Is(err, ErrHelp)
}

func appendHelpOption(cmd *Command) {
	if _, ok := cmd.getHelpOption(); !ok {
		cmd.Options = append(cmd.Options, &BoolOption{
			Name:        HelpOptionName,
			Description: "show usage",
			Default:     Default(false),
		})
	}

	for _, subcmd := range cmd.SubCommands {
		appendHelpOption(subcmd)
	}
}

func (cmd *Command) getHelpOption() (helpOption *BoolOption, ok bool) {
	for _, opt := range cmd.Options {
		if o, ok := opt.(*BoolOption); ok {
			if o.Name == HelpOptionName {
				return o, true
			}
		}
	}

	return nil, false
}
