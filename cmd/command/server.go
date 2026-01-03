package command

import (
	"fmt"
	tf "github.com/puni9869/pinmyblogs/pkg/template_functions"
	"html/template"
	"io/fs"
	"net/http"
	"os"

	"github.com/puni9869/pinmyblogs/server/middlewares"

	"gorm.io/gorm"

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
		if r := recover(); r != nil {
			log.Errorf("panic during startup: %v", r)
			panic(r) // or os.Exit(1)
		}
	}()
	e := os.Getenv(environmentKey)
	if err = config.LoadConfig(e); err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	log.Infof("Loading environment... %s", config.GetEnv())
	log.Infoln("App config loaded...")

	// initiate the db connection
	var db *gorm.DB

	dbType := config.C.AppConfig.DataBaseType
	dbConfig := config.C.Database[dbType]
	switch dbConfig.Type {
	case "sqlite":
		db, err = database.NewSqliteConnection(&dbConfig)
	case "postgres":
		db, err = database.NewPostgresConnection(&dbConfig)
	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}
	if err != nil {
		return err //nolint:wrapcheck
	}
	log.Infoln("App database loaded...")
	log.Infoln("App database migrations begin...")
	database.RegisterModels(db)
	log.Infoln("App database migrations end...")

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

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middlewares.CacheMiddleware())

	// TODO: remove this after the start2 in the prod.
	// This is new mindset of bookmarks,
	// Motivation is to keep it simple and dumb
	exemptedRoutes := []string{"/start2"}
	if config.GetEnv() == config.ProdEnv {
		r.Use(middlewares.CSP(exemptedRoutes))
		r.Use(middlewares.SecurityHeaders(exemptedRoutes))
	}
	// --- Load embedded templates ---
	// Load the template first because they are not thread-safe

	tmplFS, err := fs.Sub(pinmyblogs.Files, "templates")
	if err != nil {
		return fmt.Errorf("load templates fs: %w", err)
	}
	var tmpl = template.Must(template.New("").
		Funcs(template.FuncMap{
			"add":          tf.Add,
			"sub":          tf.Sub,
			"relativeTime": tf.FormatRelativeTime,
			"domainName":   tf.DomainName,
			"asset":        tf.Asset(BuildVersion),
		}).ParseFS(tmplFS, "**/*.tmpl"))
	r.SetHTMLTemplate(tmpl)

	// --- Serve embedded static files ---
	staticFS, err := fs.Sub(pinmyblogs.Files, "frontend")
	if err != nil {
		return fmt.Errorf("load frontend fs: %w", err)
	}
	r.StaticFS("/statics", http.FS(staticFS))

	// ____ Serve the icons
	iconsFS, err := fs.Sub(pinmyblogs.Files, "frontend/icons")
	if err != nil {
		return fmt.Errorf("load frontend/icons fs: %w", err)
	}
	r.StaticFS("/icons", http.FS(iconsFS))

	// register all the server routes
	server.RegisterRoutes(r, sessionStore)
	err = r.Run(":" + config.C.AppConfig.CustomPort)
	if err != nil {
		return fmt.Errorf("server init failed: %w", err)
	}
	return nil
}
