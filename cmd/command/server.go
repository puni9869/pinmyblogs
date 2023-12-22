package command

import (
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/server"
	"github.com/urfave/cli"
)

// Server configures the command name and action.
var Server = cli.Command{
	Name:    "server",
	Aliases: []string{"up"},
	Usage:   "Starts the Web server",
	Action:  startAction,
}

// versionAction prints the current version
func startAction(ctx *cli.Context) error {
	router := gin.Default()
	gin.SetMode(gin.DebugMode)

	// Serve the static content like *.js, *.css, *.icon, *.img
	router.Static("/statics", "./frontend")

	// Serve the templates strict to the extension *.tmpl
	router.LoadHTMLGlob("templates/*/**.tmpl")
	server.RegisterRoutes(router)
	err := router.Run()
	if err != nil {
		panic(err)
	}
	return nil
}
