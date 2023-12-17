package command

import (
	"fmt"

	"github.com/urfave/cli"
)

// Version configures the command name, flags, and action.
var Version = cli.Command{
	Name:   "version",
	Usage:  "Shows version information",
	Action: versionAction,
}

const version = "v1.0"

// versionAction prints the current version
func versionAction(ctx *cli.Context) error {
	fmt.Println(version)
	return nil
}
