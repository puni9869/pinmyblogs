package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "home.tmpl", nil)
}

func Logout(c *gin.Context) {
	c.String(http.StatusOK, "Logging out...")
}
