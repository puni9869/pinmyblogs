package server

import (
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/server/home"
)

// RegisterRoutes configures the available Web server routes.
func RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", home.Home)
}
