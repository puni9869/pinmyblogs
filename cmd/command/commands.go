package command

import "github.com/urfave/cli"

// Pinmyblogs contains the pinmyblogs CLI (sub-)commands.
var PinmyblogsCommands = []cli.Command{
	Start,
	Version,
}
