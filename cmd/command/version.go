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

var BuildVersion = "v1.0"

// versionAction prints the current version
func versionAction(_ *cli.Context) error {
	fmt.Println(BuildVersion)
	return nil
}
