package command

import (
	"github.com/puni9869/pinmyblogs/pkg/utils"
	"gorm.io/gorm"
	"html/template"
	"io/fs"
	"net/http"
	"os"

	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/pkg/database"
	"github.com/puni9869/pinmyblogs/server"

	"github.com/puni9869/pinmyblogs"
	"github.com/puni9869/pinmyblogs/pkg/config"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"github.com/urfave/cli"
)

var environmentKey = "ENVIRONMENT"

// Server configures the command name and action.
var Server = cli.Command{
	Name:    "server",
	Aliases: []string{"up"},
	Usage:   "Starts the web server",
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
	var db *gorm.DB
	dbConfig := config.C.Database["sqlite"]

	switch dbConfig.Type {
	case "sqlite":
		db, err = database.NewSqliteConnection(&dbConfig)
	case "postgres":
		db, err = database.NewPostgresConnection(&dbConfig)
	}

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
		log.Infof("Web server will listen on port: %s", config.C.AppConfig.CustomPort)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// --- Serve embedded static files ---
	staticFS, err := fs.Sub(pinmyblogs.Files, "frontend")
	if err != nil {
		return nil
	}
	r.StaticFS("/statics", http.FS(staticFS))

	// --- Load embedded templates ---
	tmplFS, err := fs.Sub(pinmyblogs.Files, "templates")
	if err != nil {
		return nil
	}
	var tmpl = template.Must(template.New("").Funcs(template.FuncMap{"relativeTime": utils.FormatRelativeTime}).Funcs(template.FuncMap{"domainName": utils.DomainName}).ParseFS(tmplFS, "**/*.tmpl"))
	r.SetHTMLTemplate(tmpl)

	// register all the server routes
	server.RegisterRoutes(r, sessionStore)
	err = r.Run(":" + config.C.AppConfig.CustomPort)
	if err != nil {
		return err
	}
	return nil
}
