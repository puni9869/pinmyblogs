/*
Copyright (c) 2017 - 2023 Pinmyblogs. All rights reserved.

	This program is free software: you can redistribute it and/or modify
	it under Version 3 of the GNU General Public License (the "GPL"):

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

Feel free to send an email to hello@pinmyblogs.com if you have questions,
want to support our work, or just want to say hello.
*/
package main

import (
	"github.com/puni9869/pinmyblogs/cmd/command"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"github.com/urfave/cli"
	"os"
)

const version = "development"

const appName = "Pinmyblogs"
const appAbout = "Pinmyblogs®"
const appEdition = "ce"
const appDescription = "Pinmyblogs is the url saving website."
const appCopyright = "(c) 2017-2024 Pinmyblogs. All rights reserved."

// Metadata contains build specific information.
var Metadata = map[string]interface{}{
	"Name":        appName,
	"About":       appAbout,
	"Edition":     appEdition,
	"Description": appDescription,
	"Version":     version,
}

func main() {
	log := logger.NewLogger()

	defer func() {
		if r := recover(); r != nil {
			os.Exit(1)
		}
	}()

	app := cli.NewApp()
	app.Usage = appAbout
	app.Description = appDescription
	app.Version = version
	app.Copyright = appCopyright
	app.Commands = command.PinmyblogsCommands
	app.EnableBashCompletion = true
	if err := app.Run(os.Args); err != nil {
		log.Error(err)
	}
}
