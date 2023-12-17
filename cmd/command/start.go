package command

import "github.com/urfave/cli"

// Start configures the command name, flags, and action.
var Start = cli.Command{
	Name:    "start",
	Aliases: []string{"up"},
	Usage:   "Starts the Web server",
	Flags:   startFlags,
	Action:  startAction,
}

// startFlags specifies the start command parameters.
var startFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "detach-server, d",
		Usage: "detach from the console (daemon mode)",
	},
	cli.BoolFlag{
		Name:  "config, c",
		Usage: "show config",
	},
}

// versionAction prints the current version
func startAction(ctx *cli.Context) error {
	return nil
}
