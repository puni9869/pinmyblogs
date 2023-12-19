package command

import (
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/server"
	"github.com/urfave/cli"
)

// Start configures the command name, flags, and action.
var Server = cli.Command{
	Name:    "server",
	Aliases: []string{"up"},
	Usage:   "Starts the Web server",
	Action:  startAction,
}

// versionAction prints the current version
func startAction(ctx *cli.Context) error {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*/**")
	server.RegisterRoutes(router)
	err := router.Run()
	if err != nil {
		panic(err)
	}
	return nil
}
