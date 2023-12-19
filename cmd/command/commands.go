package command

import "github.com/urfave/cli"

// PinmyblogsCommands contains the pinmyblogs CLI (sub-)commands.
var PinmyblogsCommands = []cli.Command{
	Server,
	Version,
}
