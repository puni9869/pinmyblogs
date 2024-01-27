package server

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/server/auth"
	"github.com/puni9869/pinmyblogs/server/home"
	"github.com/puni9869/pinmyblogs/server/setting"
)

// RegisterRoutes configures the available Web server routes.
func RegisterRoutes(r *gin.Engine, sessionStore gorm.Store) {
	r.Use(sessions.Sessions("session", sessionStore))
	// diagnose url
	r.GET("/health", home.Health)

	// auth urls
	r.Any("/logout", auth.Logout)
	r.POST("/login", auth.Login)
	r.GET("/login", auth.Login)
	r.GET("/signup", auth.Signup)
	r.GET("/index", auth.Index)

	// navbar handler
	r.GET("/", auth.Login)
	r.GET("/home", home.Home)
	r.GET("/favourite", home.Favourite)
	r.GET("/archived", home.Archived)
	r.GET("/trash", home.Trash)

	// setting handler
	r.GET("/setting", setting.Setting)

	// this route will accept all the params
	r.NoRoute(home.Home)
}
