package cliz

import (
	"fmt"
	"io"
	"strings"
)

func (cmd *Command) ShowUsage() {
	if cmd.UsageFunc != nil {
		cmd.UsageFunc(cmd)
		return
	}
	showUsage(Stderr, cmd)
}

//nolint:cyclop,funlen,gocognit
func showUsage(w io.Writer, cmd *Command) {
	const indent = "    "

	// Usage
	usage := "Usage:" + "\n"
	if cmd.Usage != "" {
		usage += indent + cmd.Usage + "\n"
	} else {
		usage += indent + strings.Join(cmd.calledCommands, " ")
		if len(cmd.Options) > 0 {
			usage += " [options]"
		}
		if len(cmd.SubCommands) > 0 {
			usage += " <subcommand>"
		}
		usage += "\n"
	}
	usage += "\n"

	// Description
	usage += "Description:" + "\n"
	usage += indent + cmd.getDescription() + "\n"

	// Options
	{
		if len(cmd.SubCommands) > 0 {
			usage += "\n"
			usage += "sub commands:\n"
			for _, subcmd := range cmd.SubCommands {
				usage += indent + fmt.Sprintf("%s: %s", subcmd.Name, subcmd.getDescription()) + "\n"
			}
		}
		if len(cmd.Options) > 0 { //nolint:nestif
			usage += "\n"
			usage += "options:\n"
			for _, opt := range cmd.Options {
				name := opt.GetName()
				short := opt.GetShort()
				env := opt.GetEnvironment()
				usage += indent
				if name != "" {
					usage += fmt.Sprintf("%s%s", longOptionPrefix, name)
				}
				if short != "" {
					if name != "" {
						usage += ", "
					}
					usage += fmt.Sprintf("%s%s", shortOptionPrefix, short)
				}

				usage += " ("

				if env != "" {
					usage += fmt.Sprintf("env: %s, ", env)
				}

				if opt.HasDefault() {
					usage += fmt.Sprintf("default: %v", opt.getDefault())
				} else {
					usage += "required"
				}

				usage += ")"

				usage += "\n"
				usage += indent + indent + opt.GetDescription() + "\n"
			}
		}
	}

	// Output
	_, _ = fmt.Fprint(w, usage)
}
