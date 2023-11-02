package cliz

import (
	"fmt"
	"io"
	"strings"
)

func (cmd *Command) usage() {
	if cmd.Usage != nil {
		cmd.Usage(cmd)
		return
	}
	defaultUsage(cmd, Stderr)
}

//nolint:gocognit,cyclop
func defaultUsage(cmd *Command, w io.Writer) {
	const indent = "    "

	usage := "Usage:" + "\n"
	usage += indent + fmt.Sprintf("%s [options] [arguments]\n", strings.Join(cmd.cmdStack, " ")) + "\n"
	usage += "Description:" + "\n"
	usage += indent + cmd.getDescription() + "\n"

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
				if env != "" {
					usage += " ("
					usage += fmt.Sprintf("env: %s", env)
				}
				if opt.HasDefault() {
					if env != "" {
						usage += ", "
					} else {
						usage += " ("
					}
					usage += fmt.Sprintf("default: %v", opt.getDefault())
				}
				if env != "" || opt.HasDefault() {
					usage += ")"
				}

				usage += "\n"
				usage += indent + indent + opt.GetDescription() + "\n"
			}
		}
	}

	_, _ = fmt.Fprint(w, usage)
}
