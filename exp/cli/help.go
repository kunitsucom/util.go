package cliz

func appendHelpOption(cmd *Command) {
	if _, ok := cmd.getHelpOption(); !ok {
		cmd.Options = append(cmd.Options, &BoolOption{
			Name:        HelpOptionName,
			Description: "show usage",
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
