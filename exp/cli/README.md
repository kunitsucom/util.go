# package cliz

## Usage

```go
import (
    cliz "github.com/kunitsucom/util.go/exp/cli"
)

func main() {
    cmd := &cliz.Command{
        Name: "my-cli",
        Description: "My awesome CLI tool",
        Usage: "my-cli [options] <subcommand> [arguments...]",
        Options: []Option{
            &BoolOption{
                Name:        "version",
                Short:       "v",
                Description: "show version",
                Default:     Default(false),
            },
        },
        SubCommands: []*Command{
            {
                Name:        "sub-cmd",
                Description: `My awesome CLI tool's sub command.`,
            },
        },
    }

    remaining, err := cmd.Parse(os.Args[0:])
    if err != nil {
        if errors.Is(err, cliz.ErrHelp) {
            return
        }
        log.Fatalf("failed to parse command line arguments: %+v", err)
    }
```
