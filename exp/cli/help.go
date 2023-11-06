package cliz

import "errors"

// HelpOptionName is the option name for help.
const HelpOptionName = "help"

// IsHelp returns whether the error is ErrHelp.
func IsHelp(err error) bool {
	return errors.Is(err, ErrHelp)
}

func (cmd *Command) initAppendHelpOption() {
	if _, ok := cmd.getHelpOption(); !ok {
		cmd.Options = append(cmd.Options, &BoolOption{
			Name:        HelpOptionName,
			Description: "show usage",
			Default:     Default(false),
		})
	}

	for _, subcmd := range cmd.SubCommands {
		subcmd.initAppendHelpOption()
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

func (cmd *Command) checkHelp() error {
	TraceLog.Printf("checkHelp: %s", cmd.Name)
	v, err := cmd.getBoolOption(HelpOptionName)
	if err == nil && v {
		cmd.ShowUsage()
		return ErrHelp
	}
	for _, subcmd := range cmd.SubCommands {
		if err := subcmd.checkHelp(); err != nil {
			return err
		}
	}
	return nil
}
