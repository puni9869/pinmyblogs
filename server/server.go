package server

import (
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/server/home"
)

// RegisterRoutes configures the available Web server routes.
func RegisterRoutes(r *gin.Engine) {
	// diagnose url
	r.GET("/health", home.Health)

	// navbar handler
	r.GET("/", home.Home)
	r.GET("/home", home.Home)
	r.GET("/favourite", home.Favourite)
	r.GET("/archived", home.Archived)

	// this route will accept all the params
	r.NoRoute(home.OK)
}
