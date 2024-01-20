package command

import (
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/pkg/config"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"github.com/puni9869/pinmyblogs/server"
	"github.com/urfave/cli"
	"os"
)

var environmentKey = "ENVIRONMENT"

// Server configures the command name and action.
var Server = cli.Command{
	Name:    "server",
	Aliases: []string{"up"},
	Usage:   "Starts the Web server",
	Action:  startAction,
}

// versionAction prints the current version
func startAction(ctx *cli.Context) error {
	log := logger.NewLogger()
	e := os.Getenv(environmentKey)
	if err := config.LoadConfig(e); err != nil {
		return err
	}
	log.Infof("Loading environment... %s", config.GetEnv())
	log.Infoln("App config loaded...")

	if config.C.AppConfig.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
		log.Infof("Web server will listen on port: %s", config.C.AppConfig.DefaultPort)
	}
	router := gin.Default()
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
