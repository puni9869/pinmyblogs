package command

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/pkg/config"
	"github.com/puni9869/pinmyblogs/pkg/database"
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
	var err error
	e := os.Getenv(environmentKey)
	if err = config.LoadConfig(e); err != nil {
		return err
	}
	log.Infof("Loading environment... %s", config.GetEnv())
	log.Infoln("App config loaded...")

	// initiate the db connection
	db, err := database.NewConnection(&config.C.Database)
	if err != nil {
		log.WithError(err)
		return err
	}

	_ = database.RegisterModels(db)

	// webapp apply debug level
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
	// register all the server routes
	server.RegisterRoutes(router)
	//router.Use(setCors())
	err = router.Run()
	if err != nil {
		panic(err)
	}
	return nil
}
func setCors() gin.HandlerFunc {
	// - No origin allowed by default
	// - GET,POST, PUT, HEAD methods
	// - Credentials share disabled
	// - Preflight requests cached for 12 hours
	corsConfig := cors.DefaultConfig()
	//corsConfig.AllowOrigins = []string{"http://google.com"}
	// config.AllowOrigins = []string{"http://google.com", "http://facebook.com"}
	//config.AllowAllOrigins = true
	return cors.New(corsConfig)
}
