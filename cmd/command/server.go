package command

import (
	gormsessions "github.com/gin-contrib/sessions/gorm"
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
	var err error
	log := logger.NewLogger()
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
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

	database.RegisterModels(db)

	// sessionStore is store in database for session values
	sessionStore := gormsessions.NewStore(db, true, []byte(config.C.AppConfig.SecretKey))

	//CSRF := csrf.Protect([]byte("32-byte-long-auth-key"))
	// webapp apply debug level
	if config.C.AppConfig.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
		log.Infof("Web server will listen on port: %s", config.C.AppConfig.DefaultPort)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Serve the static content like *.js, *.css, *.icon, *.img
	r.Static("/statics", "./frontend")
	// Serve the templates strict to the extension *.tmpl
	r.LoadHTMLGlob("templates/*/**.tmpl")
	// register all the server routes
	server.RegisterRoutes(r, sessionStore)
	err = r.Run(":" + config.C.AppConfig.CustomPort)
	if err != nil {
		return err
	}
	return nil
}
