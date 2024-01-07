package server

import (
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/server/home"
)

// RegisterRoutes configures the available Web server routes.
func RegisterRoutes(r *gin.Engine) {
	r.GET("/home", home.Home)
	r.GET("/favourite", home.Favourite)
	r.GET("/archived", home.Archived)
	r.GET("/", home.Home)
}
